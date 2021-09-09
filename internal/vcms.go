// Package vcms contains everything common to both programs. DRY, right?
package vcms

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

/*
AppXxx stores strings common to both the Collector and Receiver.
ProjectURL is the URL of the project on GitHub.
*/
const (
	AppDate    string = "2021-09-09"
	AppVersion string = "0.0.9"
	AppTitle   string = "Vaughany's Computer Monitoring System"
	ProjectURL string = "github.com/vaughany/vcms"
	AppDesc    string = "Description of the whole system goes here."
	LogFolder  string = "logs"
)

/*
SystemData represents the data we're marshalling to JSON in the Collector, and from JSON in the Receiver.
*/
type SystemData struct {
	Meta           Meta       `json:"meta"`
	Hostname       string     `json:"hostname"`
	IPAddress      string     `json:"ip_address"`
	Username       string     `json:"username"`
	FirstSeen      time.Time  `json:"first_seen"`
	LastSeen       time.Time  `json:"last_seen"`
	HostUptime     string     `json:"host_uptime"`
	OSVersion      string     `json:"os_version"`
	RebootRequired bool       `json:"reboot_required"`
	Memory         TNF        `json:"memory"`
	Swap           TNF        `json:"swap"`
	Disk           TNF        `json:"disk"`
	LoadAvgs       [3]float64 `json:"load_avgs"`
	CPU            CPU        `json:"cpu"`
	// UpdateRequired bool     `json:"update_required"`
}

// Meta holds meta-information about the Collector app itself.
type Meta struct {
	AppVersion string   `json:"app_version"`
	AppUptime  string   `json:"app_uptime"`
	Errors     []string `json:"errors"`
}

// CPU struct holds details about the CPU.
type CPU struct {
	Count int    `json:"count"`
	Speed string `json:"speed"`
}

// TNF struct stores "total 'n' free" data.
type TNF struct {
	Total int `json:"total"`
	Free  int `json:"free"`
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
// func (s SystemData) CheckPopulatedValues() bool {
// 	switch {
// 	case len(s.Hostname) == 0:
// 		log.Panic("hostname is empty")
// 		return false
// 	case len(s.IPAddress) == 0:
// 		log.Panic("IP address is empty.")
// 		return false
// 	}
// 	return true
// }

/*
Version formats and prints to the log the version details of the application, then quits.
*/
func Version(appName string) string {
	// return fmt.Sprintf("%s v%s (%s), %s.\n%s\n", appName, AppVersion, AppDate, runtime.Version(), AppDesc)
	return fmt.Sprintf("%s v%s (%s), %s.", appName, AppVersion, AppDate, runtime.Version())
}

/*
CheckLogFolder checks for the existence of the log folder.
*/
func CheckLogFolder(logFolder string) bool {
	if stat, err := os.Stat(logFolder); err != nil || !stat.IsDir() {
		log.Printf("Log folder '%s' doesn't exist.\n", logFolder)
		if !CreateLogFolder(logFolder) {
			log.Printf("Error: Log folder '%s' could not be created.\n", logFolder)
			return false
		}
		log.Printf("Log folder '%s' created.\n", logFolder)
		return true
	}
	log.Printf("Log folder '%s' exists.\n", logFolder)
	return true
}

/*
CreateLogFolder creates a log folder.
*/
func CreateLogFolder(logFolder string) bool {
	err := os.Mkdir(logFolder, 0755)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

/*
CreateLogName creates and returns a name for the log file, based on the application's start date.
*/
func CreateLogName(logFolder string, cmdName string) string {
	// time := time.Now().Format("2006-01-02_15-04-05")
	// time := time.Now().Format("2006-01-02_15-04")
	time := time.Now().Format("2006-01-02")

	return fmt.Sprintf("%s/%s_%s.log", logFolder, cmdName, time)
}
