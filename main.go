package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	connectionString := os.Getenv("CONNECTION_STRING")
	if connectionString == "" {
		log.Fatal("No CONNECTION_STRING in environment")
	}

	client, err := NewIotHubHTTPClientFromConnectionString(connectionString)
	if err != nil {
		log.Fatalln("Error creating http client from connection string", err)
	}

	defaultDeviceID := "testDevice1"

	if !client.IsDevice() {
		log.Printf("No DeviceId in connection string, running provisioning test.")
		resp, status := client.ListDeviceIDs(10)
		log.Printf("%s, %s\n\n", resp, status)
		resp, status = client.CreateDeviceID(defaultDeviceID)
		log.Printf("%s, %s\n\n", resp, status)
		resp, status = client.GetDeviceID(defaultDeviceID)
		log.Printf("%s, %s\n\n", resp, status)
		resp, status = client.PurgeCommandsForDeviceID(defaultDeviceID)
		log.Printf("%s, %s\n\n", resp, status)
		//resp, status = DeleteDeviceID(hub, defaultDeviceID)
		//log.Printf("%s, %s\n\n", resp, status)
	} else {
		log.Printf("DeviceID defined in connection string, running message test.")
		resp, status := client.ReceiveMessage()
		log.Printf("%s, %s\n\n", resp, status)
		for i := 0; i < 10; i++ {
			resp, status = client.SendMessage(fmt.Sprintf(`{"count":%d}`, i))
			log.Printf("%s, %s\n\n", resp, status)
		}
	}
}
