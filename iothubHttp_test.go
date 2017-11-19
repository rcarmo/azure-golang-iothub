package main

import "testing"

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
