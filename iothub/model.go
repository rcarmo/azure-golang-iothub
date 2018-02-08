package iothub

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Device is the struct representing the JSON received from the REST API
// Spec: https://docs.microsoft.com/en-us/rest/api/iothub/deviceapi/getdevice#device
type Device struct {
	DeviceId                   string                  `json:"deviceId"`
	GenerationId               string                  `json:"generationId"`
	Etag                       string                  `json:"etag"`
	ConnectionState            string                  `json:"connectionState"`
	Status                     string                  `json:"status"`
	StatusReason               string                  `json:"statusReason"`
	ConnectionStateUpdatedTime string                  `json:"connectionStateUpdatedTime"`
	StatusUpdatedTime          string                  `json:"statusUpdatedTime"`
	LastActivityTime           string                  `json:"lastActivityTime"`
	CloudToDeviceMessageCount  int64                   `json:"cloudToDeviceMessageCount"`
	Authentication             AuthenticationMechanism `json:"authentication"`
}

var connectionStates = []string{"disconnected", "connected"}
var stata = []string{"disabled", "enabled"}

// Unmarshal allows to populate the fields of the Device struct taking those
// values from the received JSON payload of the HTTP response
func (d *Device) Unmarshal(deviceJSON string) error {
	err := json.Unmarshal([]byte(deviceJSON), &d)
	if err != nil {
		return err
	}
	// Spec: https://docs.microsoft.com/en-us/rest/api/iothub/deviceapi/getdevice#deviceconnectionstate
	connStateValid := false
	for _, state := range connectionStates {
		if strings.ToLower(d.ConnectionState) == state {
			connStateValid = true
			break
		}
	}
	if !connStateValid {
		return errors.New(fmt.Sprintf("The connection state is not recognized: %s", d.ConnectionState))
	}
	// Spec: https://docs.microsoft.com/en-us/rest/api/iothub/deviceapi/getdevice#devicestatus
	statusValid := false
	for _, currStatus := range stata {
		if strings.ToLower(d.Status) == currStatus {
			statusValid = true
			break
		}
	}
	if !statusValid {
		return errors.New(fmt.Sprintf("This status is not recognized: %s", d.Status))
	}
	return nil
}

// AuthenticationMechanism is the struct representing some of the nested fields of the JSON received from the REST API
// Spec: https://docs.microsoft.com/en-us/rest/api/iothub/deviceapi/getdevice#authenticationmechanism
type AuthenticationMechanism struct {
	SymmetricKey   SymmetricKey   `json:"symmetricKey"`
	X509Thumbprint X509Thumbprint `json:"x509Thumbprint"`
}

// SymmetricKey is the struct representing some of the nested fields of the JSON received from the REST API
// Spec: https://docs.microsoft.com/en-us/rest/api/iothub/deviceapi/getdevice#symmetrickey
type SymmetricKey struct {
	PrimaryKey   string `json:"primaryKey"`
	SecondaryKey string `json:"secondaryKey"`
}

// X509Thumbprint is the struct representing some of the nested fields of the JSON received from the REST API
// Spec: https://docs.microsoft.com/en-us/rest/api/iothub/deviceapi/getdevice#X509Thumbprint
type X509Thumbprint struct {
	PrimaryThumbprint   string `json:"primaryThumbprint"`
	SecondaryThumbprint string `json:"secondaryThumbprint"`
}
