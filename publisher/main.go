/* OLD VERSION package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/koron/go-ssdp"
)

// GetLocalIP retrieves the non-loopback local IP address of the host
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// Check if the address is an IP address and not a loopback
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// Ensure the IP is IPv4
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no non-loopback IPv4 address found")
}

func main() {
	localIP, err := GetLocalIP()
	if err != nil {
		log.Fatalf("Error retrieving local IP address: %v", err)
	}

	// Construct the service description URL using the local IP address
	serviceDescURL := fmt.Sprintf("http://%s:8080/description.xml", localIP)

	// Create a new SSDP advertiser
	ad, err := ssdp.Advertise(
		"urn:schemas-upnp-org:service:YourService:1", // Service Type (ST)
		"uuid:grok",                // Unique Service Name (USN)
		serviceDescURL,                               // Location of service description
		"OS/version UPnP/1.1 product/version",        // Server information
		1800,                                         // Cache-Control max-age in seconds
	)
	if err != nil {
		log.Fatalf("Failed to advertise service: %v", err)
	}
	defer ad.Close()

	// Set up a ticker to send periodic alive messages
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()

	// Channel to handle OS interrupts (e.g., Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Main loop to handle alive messages and graceful shutdown
	for {
		select {
		case <-ticker.C:
			// Send an ssdp:alive message
			if err := ad.Alive(); err != nil {
				log.Printf("Failed to send alive message: %v", err)
			}
		case <-quit:
			// Send an ssdp:byebye message before exiting
			if err := ad.Bye(); err != nil {
				log.Printf("Failed to send byebye message: %v", err)
			}
			return
		}
	}
}
*/


package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/koron/go-ssdp"
)

// XML content to be served at /description.xml
const descriptionXML = `<?xml version="1.0" encoding="utf-8"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
  <specVersion>
    <major>1</major>
    <minor>0</minor>
  </specVersion>
  <device>
    <deviceType>urn:schemas-upnp-org:device:MediaServer:1</deviceType>
    <friendlyName>My Media Server</friendlyName>
    <manufacturer>Example Manufacturer</manufacturer>
    <modelName>Model XYZ</modelName>
    <UDN>uuid:unique-device-id</UDN>
    <iconList>
      <icon>
        <mimetype>image/png</mimetype>
        <width>48</width>
        <height>48</height>
        <depth>24</depth>
        <url>/icons/icon.png</url>
      </icon>
    </iconList>
    <serviceList>
      <service>
        <serviceType>urn:schemas-upnp-org:service:ContentDirectory:1</serviceType>
        <serviceId>urn:upnp-org:serviceId:ContentDirectory</serviceId>
        <SCPDURL>/service/ContentDirectory1.xml</SCPDURL>
        <controlURL>/service/ContentDirectory_control</controlURL>
        <eventSubURL>/service/ContentDirectory_event</eventSubURL>
      </service>
      <service>
        <serviceType>urn:schemas-upnp-org:service:ConnectionManager:1</serviceType>
        <serviceId>urn:upnp-org:serviceId:ConnectionManager</serviceId>
        <SCPDURL>/service/ConnectionManager1.xml</SCPDURL>
        <controlURL>/service/ConnectionManager_control</controlURL>
        <eventSubURL>/service/ConnectionManager_event</eventSubURL>
      </service>
    </serviceList>
  </device>
</root>`

// GetLocalIP retrieves the non-loopback local IP address of the host
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// Check if the address is an IP address and not a loopback
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// Ensure the IP is IPv4
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no non-loopback IPv4 address found")
}

// handler serves the description.xml content
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, descriptionXML)
}

func main() {
	localIP, err := GetLocalIP()
	if err != nil {
		log.Fatalf("Error retrieving local IP address: %v", err)
	}

	// Construct the service description URL using the local IP address
	serviceDescURL := fmt.Sprintf("http://%s:8080/description.xml", localIP)

	// Create a new SSDP advertiser
	ad, err := ssdp.Advertise(
		"urn:schemas-upnp-org:service:YourService:1", // Service Type (ST)
		"uuid:your-unique-service-id",                // Unique Service Name (USN)
		serviceDescURL,                               // Location of service description
		"OS/version UPnP/1.1 product/version",        // Server information
		1800,                                         // Cache-Control max-age in seconds
	)
	if err != nil {
		log.Fatalf("Failed to advertise service: %v", err)
	}
	defer ad.Close()

	// Set up a ticker to send periodic alive messages
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()

	// Channel to handle OS interrupts (e.g., Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Start the HTTP server in a new goroutine
	go func() {
		http.HandleFunc("/description.xml", handler)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Main loop to handle alive messages and graceful shutdown
	for {
		select {
		case <-ticker.C:
			// Send an ssdp:alive message
			if err := ad.Alive(); err != nil {
				log.Printf("Failed to send alive message: %v", err)
			}
		case <-quit:
			// Send an ssdp:byebye message before exiting
			if err := ad.Bye(); err != nil {
				log.Printf("Failed to send byebye message: %v", err)
			}
			return
		}
	}
}

