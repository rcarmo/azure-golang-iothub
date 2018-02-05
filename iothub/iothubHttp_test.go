package iothub

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestConnectionStringParsing(t *testing.T) {
	sampleConnectionString := "HostName=blahblah.azure-devices.net;SharedAccessKeyName=device;SharedAccessKey=y2R1N8XvMBRjN9yl+r3Z4vuYhpHMuWc8zvUpF/1e2IM="
	hostName, sharedAccessKey, sharedAccessKeyName, deviceID, err := parseConnectionString(sampleConnectionString)
	if err != nil {
		t.Fatal("Error parsing correct connection string")
	}

	if hostName != "blahblah.azure-devices.net" {
		t.Error("Host name is not parsed correctly")
	}

	if sharedAccessKeyName != "device" {
		t.Error("SharedAccessKeyName is not parsed correctly")
	}

	if sharedAccessKey != "y2R1N8XvMBRjN9yl+r3Z4vuYhpHMuWc8zvUpF/1e2IM=" {
		t.Error("SharedAccessKey is not parsed correctly")
	}

	if deviceID != "" {
		t.Error("Missing device id must be returned as empty string")
	}
}

func TestCreateNewDeviceThenRetrieveItThenDeleteIt(t *testing.T) {
	defaultDeviceID := "fooTestDevice"
	connectionString := os.Getenv("CONNECTION_STRING")
	if connectionString == "" {
		t.Errorf("No CONNECTION_STRING in environment")
	}

	client, err := NewIotHubHTTPClientFromConnectionString(connectionString)
	if err != nil {
		t.Errorf("Error creating http client from connection string", err)
	}

	// create a new device
	respPut, statusPut := client.CreateDeviceID(defaultDeviceID)
	if statusPut != "200 OK" {
		t.Errorf("The HTTP response status is not 200, instead: %s", statusPut)
	}
	if !strings.Contains(respPut, fmt.Sprintf("\"deviceId\":\"%s\"", defaultDeviceID)) {
		t.Errorf("The retrieved device is not the expected one with device ID '%s': %s", defaultDeviceID, respPut)
	}
	// fmt.Printf(">>> %v\n", respPut)

	// retrieve device details
	respGet, statusGet := client.GetDeviceID(defaultDeviceID)
	if statusGet != "200 OK" {
		t.Errorf("The HTTP response status is not 200, instead: %s", statusGet)
	}
	if !strings.Contains(respGet, fmt.Sprintf("\"deviceId\":\"%s\"", defaultDeviceID)) {
		t.Errorf("The retrieved device is not the expected one with device ID '%s': %s", defaultDeviceID, respGet)
	}
	// fmt.Printf(">>> %v\n", respGet)

	// delete device details
	respDel, statusDel := client.DeleteDeviceID(defaultDeviceID)
	if statusDel != "204 No Content" {
		t.Errorf("The HTTP response status is not 204, instead: %s", statusDel)
	}
	if respDel != "" {
		t.Errorf("A response body for the HTTP delete was not expected, found: %s", respDel)
	}
}
