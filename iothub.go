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
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	maxIdleConnections int    = 100
	requestTimeout     int    = 10
	apiVersion         string = "2016-02-03"
	tokenValidSecs     int    = 365 * 24 * 60 * 60
	tokenFormat        string = "SharedAccessSignature sr=%s&sig=%s&se=%s&skn=%s"
	urlFormat          string = "https://%s/devices/%s?api-version=%s"
	bulkFormat         string = "https://%s/devices/?top=%d&api-version=%s"
	eventFormat        string = "https://%s/devices/%s/messages/events?api-version=%s"
)

// IoTHub representation
type IoTHub struct {
	HostName             string
	SharedAccessKeyName  string
	SharedAccessKeyValue string
	Client               *http.Client
}

func NewIoTHub(conn string) (hub *IoTHub, err error) {
	fields, err := url.ParseQuery(conn)
	if err != nil {
		log.Fatal(err)
	}
	hub = new(IoTHub)
	t := reflect.ValueOf(hub).Elem()
	for k, v := range fields {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v[0]))
	}
	hub.Client = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConnections,
		},
		Timeout: time.Duration(requestTimeout) * time.Second,
	}
	return hub, nil
}

func buildSasToken(hub *IoTHub) string {
	key, err := base64.StdEncoding.DecodeString(hub.SharedAccessKeyValue)
	if err != nil {
		log.Fatal(err)
	}
	mac := hmac.New(sha256.New, key)
	timestamp := strconv.FormatInt(time.Now().Unix()+int64(tokenValidSecs), 10)
	signature := url.QueryEscape(base64.URLEncoding.EncodeToString(mac.Sum([]byte(strings.ToLower(hub.HostName) + "\n" + timestamp))))
	return fmt.Sprintf(tokenFormat, strings.ToLower(hub.HostName), signature, timestamp, hub.SharedAccessKeyName)
}

// Perform individual requests (we assume these won't require a persistent session)
func performRequest(hub *IoTHub, method string, url string, body string) (string, string) {
	payload := []byte(body)
	token := buildSasToken(hub)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	if method == "DELETE" {
		req.Header.Set("If-Match", "*")
	}

	resp, err := hub.Client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	text, _ := ioutil.ReadAll(resp.Body)
	return string(text), resp.Status
}

// CreateDeviceID adds a given device to an IoTHub and
// returns the HTTP request data
func CreateDeviceID(hub *IoTHub, deviceID string) (string, string) {
	url := fmt.Sprintf(urlFormat, hub.HostName, deviceID, apiVersion)
	body := fmt.Sprintf("{deviceId: \"%s\"}", deviceID)
	return performRequest(hub, "PUT", url, body)
}

func GetDeviceID(hub *IoTHub, deviceID string) (string, string) {
	url := fmt.Sprintf(urlFormat, hub.HostName, deviceID, apiVersion)
	return performRequest(hub, "GET", url, "")
}

func DeleteDeviceID(hub *IoTHub, deviceID string) (string, string) {
	url := fmt.Sprintf(urlFormat, hub.HostName, deviceID, apiVersion)
	return performRequest(hub, "DELETE", url, "")
}

func ListDeviceIDs(hub *IoTHub, top int) (string, string) {
	url := fmt.Sprintf(bulkFormat, hub.HostName, top, apiVersion)
	return performRequest(hub, "GET", url, "")
}

func SendMessage(hub *IoTHub, deviceID string, message string) (string, string) {
	url := fmt.Sprintf(eventFormat, hub.HostName, deviceID, apiVersion)
	return performRequest(hub, "POST", url, message)
}

func main() {
	hub, _ := NewIoTHub("HostName=foobar.com;SharedAccessKeyName=blabla;SharedAccessKeyValue=puppies")
	fmt.Println(buildSasToken(hub))
}
