// This 'vcms' package contains everything common to both programs. DRY, right?

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
	AppDate    string = "2021-08-07"
	AppVersion string = "0.0.6"
	AppTitle   string = "Vaughany's Computer Monitoring System"
	ProjectURL string = "github.com/vaughany/vcms"
	AppDesc    string = "Description of the whole system goes here."
)

/*
SystemData represents the data we're marshalling to JSON in the Collector, and from JSON in the Receiver.
*/
type SystemData struct {
	Meta struct {
		AppVersion string   `json:"app_version"`
		AppUptime  string   `json:"app_uptime"`
		Errors     []string `json:"errors"`
	} `json:"meta"`
	Hostname       string    `json:"hostname"`
	IPAddress      string    `json:"ip_address"`
	Username       string    `json:"username"`
	FirstSeen      time.Time `json:"first_seen"`
	LastSeen       time.Time `json:"last_seen"`
	HostUptime     string    `json:"host_uptime"`
	OSVersion      string    `json:"os_version"`
	RebootRequired bool      `json:"reboot_required"`
	MemoryTotal    int       `json:"memory_total"`
	MemoryFree     int       `json:"memory_free"`
	SwapTotal      int       `json:"swap_total"`
	SwapFree       int       `json:"swap_free"`
	DiskTotal      int       `json:"disk_total"`
	DiskFree       int       `json:"disk_free"`
	LoadAvgs       []float64 `json:"load_avgs"`
	CPUCount       int       `json:"cpu_count"`
	CPUSpeed       string    `json:"cpu_speed"`
	// UpdateRequired bool     `json:"update_required"`
}

/*
SystemDataPlusDateTime represents SystemData and a timestamp - used for file-persisted data.
*/
type SystemDataPlusDateTime struct {
	SaveDateTime time.Time              `json:"save_datetime"`
	SystemData   map[string]*SystemData `json:"system_data"`
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
