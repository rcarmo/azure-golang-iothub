# azure-golang-iot-hub

Experimental (minimalistic) Azure IoT Hub client in Go (HTTPS only).
This has been tested using Go version 1.8.5 on a linux/amd64 architecture.

## TODO

* [x] Device registration/enumeration
* [x] Device-to-cloud messages
* [x] HTTP connection re-use
* [x] Proper testing, built-in example reading connection string from environment variable
+ [x] Support both named and unnamed (`DeviceId`) connection strings
* [x] Refactor as library
* [ ] Implement AMQP client

## HOWTO

```bash
export CONNECTION_STRING='HostName=myhub.azure-devices.net;SharedAccessKeyName=iothubowner;SharedAccessKey=SxiN78h8tdN3yQXMBhmV193ZxKWBHhmJptGcvheA3dg='
make run
```

or...

```bash
export CONNECTION_STRING='HostName=myhub.azure-devices.net;DeviceId=testdevice;SharedAccessKey=SxiN78h8tdN3yQXMBhmV193ZxKWBHhmJptGcvheA3dg='
make run
```
