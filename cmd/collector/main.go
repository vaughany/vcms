package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/user"
	"sync"
	"time"
	"vcms"
)

const (
	cmdName     string = "VCMS - Collector"
	cmdDesc     string = "Collects information about the computer, and sends it to the Receiver app."
	cmdCodename string = "vcms-collector"
)

var (
	debug       bool   = false
	testing     bool   = false
	version     bool   = false
	receiverURL string = "http://127.0.0.1:8080"
)

func init() {
	flag.BoolVar(&debug, "d", debug, "Shows debugging info")
	flag.BoolVar(&testing, "t", testing, "Creates a random hostname, username and IP address")
	flag.StringVar(&receiverURL, "r", receiverURL, "URL of the 'Receiver' application")
	flag.BoolVar(&version, "v", version, "Show version info and quit")
	flag.Parse()

	if version {
		fmt.Println(vcms.Version(cmdName))
		os.Exit(0)
	}
}

func main() {
	log.Println(vcms.Version(cmdName))
	log.Printf("%s \n", cmdDesc)
	log.Printf("%s \n", vcms.AppDesc)

	var wg sync.WaitGroup

	wg.Add(1)
	go sendAnnounce()
	wg.Wait()
}

func sendAnnounce() {
	var (
		watchDelay         int = 10
		lastSuccessfulSend time.Time
	)

	data := vcms.SystemData{}

	if testing {
		data.Hostname = getRandomHostname()
		data.IPAddress = getRandomIPAddress()
		data.Username = getRandomUsername()
	} else {
		data.Hostname = getHostname()
		data.IPAddress = getIPAddress()
		data.Username = getUsername()
	}

	for {
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			log.Panic(err)
		}

		sendURL := receiverURL + "/api/announce"
		if debug {
			log.Printf("Sending %s to %s", jsonBytes, sendURL)
		} else {
			log.Printf("Sending data to %s", sendURL)
		}

		response, err := http.Post(sendURL, "application/json", bytes.NewBuffer(jsonBytes))

		if err != nil {
			log.Print(err)
			printLastSuccessfulSend(lastSuccessfulSend)
		} else {
			log.Printf("Response: %s", response.Status)
			if response.StatusCode != 200 {
				// This is fundamentally an error, so handle it.
				log.Print(response)
			} else {
				lastSuccessfulSend = time.Now()
			}
		}

		time.Sleep(time.Second * time.Duration(watchDelay))
	}
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return hostname
}

func getIPAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		// interface down
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// loopback interface
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			// return "", err
			panic(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			// return ip.String(), nil
			return ip.String()
		}
	}
	return ""
}

func getUsername() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	return user.Username
}

func printLastSuccessfulSend(t time.Time) {
	msg := "Last successful data send:"
	if t.IsZero() {
		log.Printf("%s Never.\n", msg)
	} else {
		log.Printf("%s %s (%s ago)\n", msg, t.Format(time.RFC1123Z), time.Since(t).Round(time.Second))
	}
}
