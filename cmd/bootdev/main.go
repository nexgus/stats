package main

import (
	"fmt"
	"log"
	"stats/pkg/disk"

	"github.com/google/uuid"
)

func main() {
	device, err := disk.FindBootDevice()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Boot device: %s\n", device)

	id := disk.DeviceUniqueID(device)
	fmt.Printf("Device ID: \"%s\"\n", id)

	devUUID := uuid.NewMD5(uuid.Nil, []byte(id))
	fmt.Printf("Device UUID: %s\n", devUUID.String())
}
