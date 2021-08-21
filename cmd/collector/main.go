// Package Collector collects information about the computer, and sends it to the Receiver app.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	"vcms"
)

func main() {
	const (
		cmdName     string = "VCMS - Collector"
		cmdDesc     string = "Collects information about the computer, and sends it to the Receiver app."
		cmdCodename string = "vcms-collector"
	)

	var (
		version     = false
		debug       = false
		testing     = false
		receiverURL = "http://127.0.0.1:8080"
	)

	flag.BoolVar(&debug, "d", debug, "Shows debugging info")
	flag.BoolVar(&testing, "t", testing, "Creates a random hostname, username and IP address")
	flag.StringVar(&receiverURL, "r", receiverURL, "URL of the 'Receiver' application")
	flag.BoolVar(&version, "v", version, "Show version info and quit")
	flag.Parse()

	if version {
		fmt.Println(vcms.Version(cmdName))
		os.Exit(0)
	}

	log.Println(vcms.Version(cmdName))
	log.Printf("%s \n", cmdDesc)
	log.Printf("%s \n", vcms.AppDesc)

	sendAnnounce(debug, testing, receiverURL)
}

func sendAnnounce(debug bool, testing bool, receiverURL string) {
	var (
		watchDelay         = 10
		startTime          = time.Now()
		lastSuccessfulSend time.Time
	)

	data := vcms.SystemData{}

	for {
		var errors []string

		// Data that will not change. // lol
		data.Hostname = getHostname()
		data.IPAddress = getIPAddress()
		data.Username = getUsername()
		data.OSVersion = getOSVersion()
		data.CPUCount, data.CPUSpeed = getCPUDetails()

		memoryDetails := getMemoryDetails()
		diskDetails := getDiskDetails()

		// Data that will change.
		data.HostUptime = getHostUptime()
		data.RebootRequired = getRebootRequired()
		data.MemoryTotal = memoryDetails[0]
		data.MemoryFree = memoryDetails[1]
		data.SwapTotal = memoryDetails[2]
		data.SwapFree = memoryDetails[3]
		data.DiskTotal = diskDetails[0]
		data.DiskFree = diskDetails[1]
		data.LoadAvgs = getLoadAvgs()
		data.Meta.AppVersion = vcms.AppVersion
		data.Meta.AppUptime = getAppUptime(startTime)
		data.Meta.Errors = errors

		// Adjust some of the core data if we're testing.
		if testing {
			data.Hostname = getRandomHostname()
			data.IPAddress = getRandomIPAddress()
			data.Username = getRandomUsername()
		}

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
			log.Print(getLastSuccessfulSend(lastSuccessfulSend))
		} else {
			log.Printf("Response: %s", response.Status)
			if response.StatusCode == 200 {
				lastSuccessfulSend = time.Now()
			} else {
				body, _ := io.ReadAll(response.Body)
				log.Print(string(body))
				log.Print(getLastSuccessfulSend(lastSuccessfulSend))
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

func getHostUptime() string {
	if runtime.GOOS == "windows" {
		log.Println("WARNING: 'getHostUptime()' function not yet implemented for Windows.")
		return ""
	}
	if runtime.GOOS == "solaris" {
		log.Println("WARNING: 'getHostUptime()' function not yet implemented for Solaris.")
		return ""
	}

	contents, err := os.ReadFile("/proc/uptime")
	if err != nil {
		log.Panic(err)
	}

	uptimeInt, err := strconv.Atoi(strings.Split(strings.Split(string(contents), " ")[0], ".")[0])
	if err != nil {
		log.Println(err)
	}

	return time.Since(time.Unix(time.Now().Unix()-int64(uptimeInt), 0)).Round(time.Second).String()
}

func getOSVersion() string {
	switch runtime.GOOS {
	case "windows":
		// https://en.wikipedia.org/wiki/List_of_Microsoft_Windows_versions
		// https://gist.github.com/flxxyz/ae3ef071dc4ffb0c55daedc7f0740611
		// log.Println("WARNING: 'getOsVersion()' function not yet implemented for Windows.")
		// return ""

		cmd := exec.Command("cmd.exe")
		out, _ := cmd.StdoutPipe()
		buffer := bytes.NewBuffer(make([]byte, 0))
		cmd.Start()
		buffer.ReadFrom(out)
		str, _ := buffer.ReadString(']')
		cmd.Wait()

		remove := []string{"[", "Version", "]"}
		for _, r := range remove {
			str = strings.ReplaceAll(str, r, "")
		}
		return str

	case "linux":
		release, err := os.ReadFile("/etc/os-release")
		if err != nil {
			log.Panic(err)
		}

		regexName := regexp.MustCompile(`PRETTY_NAME=.*`)
		name := regexName.FindString(string(release))

		return name[13 : len(name)-1]

	case "solaris":
		release, err := os.ReadFile("/etc/release")
		if err != nil {
			log.Panic(err)
		}

		regexName := regexp.MustCompile(`^\s+Oracle\s.*`)
		name := regexName.FindString(string(release))

		return strings.Join(strings.Fields(name)[0:3], " ")

	default:
		return ""
	}
}

func getRebootRequired() bool {
	if runtime.GOOS == "windows" {
		log.Println("WARNING: 'getRebootRequired()' function not yet implemented for Windows.")
		return false
	}

	filename := "/var/run/reboot-required"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func getMemoryDetails() [4]int {
	if runtime.GOOS == "windows" {
		log.Println("WARNING: 'getMemoryDetails()' function not yet implemented for Windows.")
		return [4]int{0, 0, 0, 0}
	}
	if runtime.GOOS == "solaris" {
		log.Println("WARNING: 'getMemoryDetails()' function not yet implemented for Solaris.")
		return [4]int{0, 0, 0, 0}
	}

	memory, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		log.Panic(err)
	}

	// https://stackoverflow.com/a/30483899/254146
	regexMemTotal := regexp.MustCompile(`MemTotal:\s+(?P<MemTotal>\d+).+`)
	regexMemFree := regexp.MustCompile(`MemFree:\s+(?P<MemFree>\d+).+`)
	regexSwapTotal := regexp.MustCompile(`SwapTotal:\s+(?P<SwapTotal>\d+).+`)
	regexSwapFree := regexp.MustCompile(`SwapFree:\s+(?P<SwapFree>\d+).+`)

	memTotal, _ := strconv.Atoi(regexMemTotal.FindStringSubmatch(string(memory))[1])
	memFree, _ := strconv.Atoi(regexMemFree.FindStringSubmatch(string(memory))[1])
	swapTotal, _ := strconv.Atoi(regexSwapTotal.FindStringSubmatch(string(memory))[1])
	swapFree, _ := strconv.Atoi(regexSwapFree.FindStringSubmatch(string(memory))[1])

	return [4]int{memTotal, memFree, swapTotal, swapFree}
}

func getDiskDetails() [2]int {
	if runtime.GOOS == "windows" {
		log.Println("WARNING: 'getDiskDetails()' function not yet implemented for Windows.")
		return [2]int{0, 0}
	}

	disk, err := exec.Command("df", "-k", "/").Output()
	if err != nil {
		log.Panic(err)
	}

	regexDisk := regexp.MustCompile(`\s+(?P<Size>\d+)\s+(?P<Used>\d+)\s+(?P<Avail>\d+)`)

	diskTotal, _ := strconv.Atoi(regexDisk.FindStringSubmatch(string(disk))[1])
	diskFree, _ := strconv.Atoi(regexDisk.FindStringSubmatch(string(disk))[3])

	return [2]int{diskTotal, diskFree}
}

func getLoadAvgs() [3]float64 {
	if runtime.GOOS == "windows" {
		log.Println("WARNING: 'getLoadAvgs()' function not yet implemented for Windows.")
		return [3]float64{0, 0, 0}
	}
	if runtime.GOOS == "solaris" {
		log.Println("WARNING: 'getLoadAvgs()' function not yet implemented for Solaris.")
		return [3]float64{0, 0, 0}
	}

	load, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		log.Panic(err)
	}

	loads := strings.Split(string(load), " ")
	one, _ := strconv.ParseFloat(loads[0], 64)
	five, _ := strconv.ParseFloat(loads[1], 64)
	fifteen, _ := strconv.ParseFloat(loads[2], 64)

	return [3]float64{one, five, fifteen}
}

func getLastSuccessfulSend(t time.Time) string {
	msg := "Last successful data send:"
	if t.IsZero() {
		return fmt.Sprintf("%s Never.\n", msg)
	} else {
		return fmt.Sprintf("%s %s (%s ago)\n", msg, t.Format(time.RFC1123Z), time.Since(t).Round(time.Second))
	}
}

func getAppUptime(startTime time.Time) string {
	return time.Since(startTime).Round(time.Second).String()
}

func getCPUDetails() (int, string) {
	if runtime.GOOS == "windows" {
		log.Println("WARNING: 'getCPUDetails()' function not yet implemented for Windows.")
		return 0, ""
	}
	if runtime.GOOS == "solaris" {
		log.Println("WARNING: 'getCPUDetails()' function not yet implemented for Solaris.")
		return 0, ""
	}

	cpuinfo, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		log.Panic(err)
	}

	regex := regexp.MustCompile(`model name\s+:\s.*`)
	cpuDetails := regex.FindAllSubmatch(cpuinfo, -1)
	cpuDetailsSlice := strings.Split(string(cpuDetails[0][0]), " ")

	count := len(cpuDetails)
	speed := fmt.Sprint(cpuDetailsSlice[len(cpuDetailsSlice)-1])

	return count, speed
}
