package main

import (
	"fmt"
	"log"

	"github.com/huin/goupnp"
	"github.com/huin/goupnp/ssdp"
)

func main() {
	// Discover all UPnP root devices
	devices, err := goupnp.DiscoverDevices(ssdp.UPNPRootDevice)
	if err != nil {
		log.Fatalf("Error discovering devices: %v", err)
	}

	// Iterate over discovered devices
	for _, device := range devices {
		if device.Err != nil {
			log.Printf("Error with device at %s: %v", device.Location.String(), device.Err)
			continue
		}

		fmt.Printf("Device found at %s\n", device.Location.String())
		fmt.Printf("  Friendly Name: %s\n", device.Root.Device.FriendlyName)
		fmt.Printf("  Manufacturer: %s\n", device.Root.Device.Manufacturer)
		fmt.Printf("  Device Type: %s\n", device.Root.Device.DeviceType)
	}
}

