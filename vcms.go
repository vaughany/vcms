package vcms

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

/*
AppXxx stores strings common to both the Collector and Receiver.
ProjectURL is the URL of the project on GitHub.
*/
const (
	AppDate    string = "2021-07-01"
	AppVersion string = "0.0.1"
	AppTitle   string = "Vaughany's Computer Monitoring System"
	ProjectURL string = "github.com/vaughany/vcms"
	AppDesc    string = "Description of the whole system goes here."
)

/*
SystemData represents the data we're marshalling to JSON in the Collector, and from JSON in the Receiver.
*/
type SystemData struct {
	Hostname  string `json:"hostname"`
	IPAddress string `json:"ip_address"`
	Username  string `json:"username"`
	FirstSeen time.Time
	LastSeen  time.Time
	// HostUptime     string   `json:"host_uptime"`
	// UpdateRequired bool     `json:"update_required"`
	// RebootRequired bool     `json:"reboot_required"`
	// OsVersion      string   `json:"os_version"`
	// LoadAvgs       []string `json:"load_avgs"`
	// AppUptime      string   `json:"app_uptime"`
	// AppVersion     string   `json:"app_version"`
	// MemoryTotal    string   `json:"memory_total"`
	// MemoryFree     string   `json:"memory_free"`
	// SwapTotal      string   `json:"swap_total"`
	// SwapFree       string   `json:"swap_free"`
	// DiskTotal      string   `json:"disk_total"`
	// DiskFree       string   `json:"disk_free"`
}

/*
CheckPopulatedValues checks each item in the SystemData struct for emptiness.
*/
func (s SystemData) CheckPopulatedValues() bool {
	switch {
	case len(s.Hostname) == 0:
		log.Panic("hostname is empty")
		return false
	case len(s.IPAddress) == 0:
		log.Panic("IP address is empty.")
		return false
	}
	return true
}

/*
Version formats and prints to the log the version details of the application, then quits.
*/
func Version(appName string) string {
	// return fmt.Sprintf("%s v%s (%s), %s.\n%s\n", appName, AppVersion, AppDate, runtime.Version(), AppDesc)
	return fmt.Sprintf("%s v%s (%s), %s.", appName, AppVersion, AppDate, runtime.Version())
}
