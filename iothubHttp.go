package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	maxIdleConnections int    = 100
	requestTimeout     int    = 10
	tokenValidSecs     int    = 3600
	apiVersion         string = "2016-11-14"
)

type sharedAccessKey string
type sharedAccessKeyName string
type hostName string
type deviceID string

// IotHubHTTPClient is a simple client to connect to Azure IoT Hub
type IotHubHTTPClient struct {
	sharedAccessKeyName sharedAccessKeyName
	sharedAccessKey     sharedAccessKey
	hostName            hostName
	deviceID            deviceID
	client              *http.Client
}

func parseConnectionString(connString string) (hostName, sharedAccessKey, sharedAccessKeyName, deviceID, error) {
	url, err := url.ParseQuery(connString)
	if err != nil {
		return "", "", "", "", err
	}

	h := tryGetKeyByName(url, "HostName")
	kn := tryGetKeyByName(url, "SharedAccessKeyName")
	k := tryGetKeyByName(url, "SharedAccessKey")
	d := tryGetKeyByName(url, "DeviceId")

	return hostName(h), sharedAccessKey(k), sharedAccessKeyName(kn), deviceID(d), nil
}

func tryGetKeyByName(v url.Values, key string) string {
	if len(v[key]) == 0 {
		return ""
	}

	return strings.Replace(v[key][0], " ", "+", -1)
}

// NewIotHubHTTPClient is a constructor of IutHubClient
func NewIotHubHTTPClient(hostNameStr string, sharedAccessKeyNameStr string, sharedAccessKeyStr string, deviceIDStr string) *IotHubHTTPClient {
	return &IotHubHTTPClient{
		sharedAccessKeyName: sharedAccessKeyName(sharedAccessKeyNameStr),
		sharedAccessKey:     sharedAccessKey(sharedAccessKeyStr),
		hostName:            hostName(hostNameStr),
		deviceID:            deviceID(deviceIDStr),
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: maxIdleConnections,
			},
			Timeout: time.Duration(requestTimeout) * time.Second,
		},
	}
}

// NewIotHubHTTPClientFromConnectionString creates new client from connection string
func NewIotHubHTTPClientFromConnectionString(connectionString string) (*IotHubHTTPClient, error) {
	h, k, kn, d, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}

	return NewIotHubHTTPClient(string(h), string(kn), string(k), string(d)), nil
}

// IsDevice tell either device id was specified when client created.
// If device id was specified in connection string this will enabled device scoped requests.
func (c *IotHubHTTPClient) IsDevice() bool {
	return c.deviceID != ""
}

// Service API

// CreateDeviceID creates record for for given device identifier in Azure IoT Hub
func (c *IotHubHTTPClient) CreateDeviceID(deviceID string) (string, string) {
	url := fmt.Sprintf("%s/devices/%s?api-version=%s", c.hostName, deviceID, apiVersion)
	data := fmt.Sprintf(`{"deviceId":"%s"}`, deviceID)
	return c.performRequest("PUT", url, data)
}

// GetDeviceID retrieves device by id
func (c *IotHubHTTPClient) GetDeviceID(deviceID string) (string, string) {
	url := fmt.Sprintf("%s/devices/%s?api-version=%s", c.hostName, deviceID, apiVersion)
	return c.performRequest("GET", url, "")
}

// DeleteDeviceID deletes device by id
func (c *IotHubHTTPClient) DeleteDeviceID(deviceID string) (string, string) {
	url := fmt.Sprintf("%s/devices/%s?api-version=%s", c.hostName, deviceID, apiVersion)
	return c.performRequest("DELETE", url, "")
}

// PurgeCommandsForDeviceID removed commands for specified device
func (c *IotHubHTTPClient) PurgeCommandsForDeviceID(deviceID string) (string, string) {
	url := fmt.Sprintf("%s/devices/%s/commands?api-version=%s", c.hostName, deviceID, apiVersion)
	return c.performRequest("DELETE", url, "")
}

// ListDeviceIDs list all device ids
func (c *IotHubHTTPClient) ListDeviceIDs(top int) (string, string) {
	url := fmt.Sprintf("%s/devices?top=%d&api-version=%s", c.hostName, top, apiVersion)
	return c.performRequest("GET", url, "")
}

// TODO: SendMessageToDevice as soon as that endpoint is exposed via HTTP

// Device API

// SendMessage from a logged in device
func (c *IotHubHTTPClient) SendMessage(message string) (string, string) {
	url := fmt.Sprintf("%s/devices/%s/messages/events?api-version=%s", c.hostName, c.deviceID, apiVersion)
	return c.performRequest("POST", url, message)
}

// ReceiveMessage to a logged in device
func (c *IotHubHTTPClient) ReceiveMessage() (string, string) {
	url := fmt.Sprintf("%s/devices/%s/messages/deviceBound?api-version=%s", c.hostName, c.deviceID, apiVersion)
	return c.performRequest("GET", url, "")
}

func (c *IotHubHTTPClient) buildSasToken(uri string) string {
	timestamp := time.Now().Unix() + int64(3600)
	encodedURI := template.URLQueryEscaper(uri)

	toSign := encodedURI + "\n" + strconv.FormatInt(timestamp, 10)

	binKey, _ := base64.StdEncoding.DecodeString(string(c.sharedAccessKey))
	mac := hmac.New(sha256.New, []byte(binKey))
	mac.Write([]byte(toSign))

	encodedSignature := template.URLQueryEscaper(base64.StdEncoding.EncodeToString(mac.Sum(nil)))

	if c.sharedAccessKeyName != "" {
		return fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%d&skn=%s", encodedURI, encodedSignature, timestamp, c.sharedAccessKeyName)
	}

	return fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%d", encodedURI, encodedSignature, timestamp)
}

func (c *IotHubHTTPClient) performRequest(method string, uri string, data string) (string, string) {
	token := c.buildSasToken(uri)
	log.Printf("%s https://%s\n", method, uri)
	req, _ := http.NewRequest(method, "https://"+uri, bytes.NewBufferString(data))
	log.Println(data)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "golang-iot-client")
	req.Header.Set("Authorization", token)

	log.Println("Authorization:", token)

	if method == "DELETE" {
		req.Header.Set("If-Match", "*")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	// read the entire reply to ensure connection re-use
	text, _ := ioutil.ReadAll(resp.Body)

	io.Copy(ioutil.Discard, resp.Body)
	defer resp.Body.Close()

	return string(text), resp.Status
}
