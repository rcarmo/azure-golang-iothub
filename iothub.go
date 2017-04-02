package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"text/template"
	"time"
)

const (
	maxIdleConnections int    = 100
	requestTimeout     int    = 10
	tokenValidSecs     int    = 3600
	statsFormat        string = "%s/statistics/devices&api-version=2016-11-14"
	listFormat         string = "%s/devices?top=%d&api-version=2016-11-14"
	urlFormat          string = "%s/devices/%s?api-version=2016-11-14"
	eventFormat        string = "https://%s/devices/%s/messages/events?api-version=%s"
)

// IoTHub representation
type IoTHub struct {
	HostName            string
	SharedAccessKeyName string
	SharedAccessKey     string
	Client              *http.Client
}

func NewIoTHub(conn string) (hub *IoTHub, err error) {
	// hijack the ParseQuery function to split the connection string into a map
	fields, err := url.ParseQuery(conn)
	if err != nil {
		log.Fatal(err)
	}
	hub = new(IoTHub)

	// use reflection to match each connection string component with a struct field
	// TODO: make sure we have all required fields
	t := reflect.ValueOf(hub).Elem()
	for k, v := range fields {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v[0]))
	}

	// set up a shared client for all connections, with long timeouts
	hub.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}
	return hub, nil
}

func buildSasToken(hub *IoTHub, uri string) string {
	timestamp := time.Now().Unix() + int64(tokenValidSecs)
	encodedUri := template.URLQueryEscaper(uri)
	toSign := encodedUri + "\n" + strconv.FormatInt(timestamp, 10)
	binKey, _ := base64.StdEncoding.DecodeString(hub.SharedAccessKey)
	mac := hmac.New(sha256.New, []byte(binKey))
	mac.Write([]byte(toSign))
	encodedSignature := template.URLQueryEscaper(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	return fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%d&skn=%s", encodedUri, encodedSignature, timestamp, hub.SharedAccessKeyName)
}

// Perform individual requests (we assume these won't require a persistent session)
func performRequest(hub *IoTHub, method string, url string, body string) (string, string) {
	payload := []byte(body)
	token := buildSasToken(hub, url)
	req, _ := http.NewRequest(method, "https://"+url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/vnd.microsoft.iothub.json")
	req.Header.Set("User-Agent", "golang-iot-client")
	req.Header.Set("Authorization", token)
	if method == "DELETE" {
		req.Header.Set("If-Match", "*")
	}

	log.Println(url)
	log.Println(body)
	resp, err := hub.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// read the entire reply to ensure connection re-use
	text, _ := ioutil.ReadAll(resp.Body)
	return string(text), resp.Status
}

// CreateDeviceID adds a given device to an IoTHub and
// returns the HTTP request data
func CreateDeviceID(hub *IoTHub, deviceID string) (string, string) {
	url := fmt.Sprintf(urlFormat, hub.HostName, deviceID)
	body := fmt.Sprintf("{\"deviceId\": \"%s\"}", deviceID)
	return performRequest(hub, "PUT", url, body)
}

func GetDeviceID(hub *IoTHub, deviceID string) (string, string) {
	url := fmt.Sprintf(urlFormat, hub.HostName, deviceID)
	return performRequest(hub, "GET", url, "")
}

func DeleteDeviceID(hub *IoTHub, deviceID string) (string, string) {
	url := fmt.Sprintf(urlFormat, hub.HostName, deviceID)
	return performRequest(hub, "DELETE", url, "")
}

func ListDeviceIDs(hub *IoTHub, top int) (string, string) {
	url := fmt.Sprintf(listFormat, hub.HostName, top)
	return performRequest(hub, "GET", url, "")
}

func GetStatistics(hub *IoTHub) (string, string) {
	url := fmt.Sprintf(statsFormat, hub.HostName)
	return performRequest(hub, "GET", url, "")
}

func SendMessage(hub *IoTHub, deviceID string, message string) (string, string) {
	url := fmt.Sprintf(eventFormat, hub.HostName, deviceID)
	return performRequest(hub, "POST", url, message)
}

func main() {
	connectionString := os.Getenv("CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("No CONNECTION_STRING in environment")
	}
	hub, _ := NewIoTHub(connectionString)
	resp, status := GetStatistics(hub)
	log.Printf("%s, %s\n\n", resp, status)
	resp, status = ListDeviceIDs(hub, 10)
	log.Printf("%s, %s\n\n", resp, status)
	resp, status = CreateDeviceID(hub, "testDevice")
	log.Printf("%s, %s\n\n", resp, status)
}
