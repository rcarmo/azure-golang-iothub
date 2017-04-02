# azure-golang-iot-hub

Experimental (minimalistic) Azure IoT Hub client in Go (HTTPS only)

## TODO

* [x] Device registration/enumeration
* [x] Device-to-cloud messages
* [x] HTTP connection re-use
* [x] Proper testing, built-in example reading connection string from environment variable
* [ ] Refactor as library

## HOWTO

```bash
export CONNECTION_STRING='HostName=myhub.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=SxiN78h8tdN3yQXMBhmV193ZxKWBHhmJptGcvheA3dg='
make run
```
