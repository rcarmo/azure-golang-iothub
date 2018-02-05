package iothub

import "testing"

func TestDeviceJSONUnmarshalling(t *testing.T) {

	responsePutJson := `
{
  "deviceId": "fooTestDevice",
  "generationId": "GENERATION_ID_VALUE",
  "etag": "ETAG_VALUE",
  "connectionState": "Disconnected",
  "status": "enabled",
  "statusReason": null,
  "connectionStateUpdatedTime": "0001-01-01T00:00:00",
  "statusUpdatedTime": "0001-01-01T00:00:00",
  "lastActivityTime": "0001-01-01T00:00:00",
  "cloudToDeviceMessageCount": 0,
  "authentication": {
    "symmetricKey": {
      "primaryKey": "PRIMARY_KEY_VALUE",
      "secondaryKey": "SECONDARY_KEY_VALUE"
    },
    "x509Thumbprint": {
      "primaryThumbprint": null,
      "secondaryThumbprint": null
    }
  }
}
`
	var iotDevice Device

	err := iotDevice.Unmarshal(responsePutJson)
	if err != nil {
		t.Error(err)
	}
	if iotDevice.Authentication.SymmetricKey.PrimaryKey != "PRIMARY_KEY_VALUE" {
		t.Error("The primary key is not recognized")
	}
	if iotDevice.Authentication.SymmetricKey.SecondaryKey != "SECONDARY_KEY_VALUE" {
		t.Error("The secondary key is not recognized")
	}
}
