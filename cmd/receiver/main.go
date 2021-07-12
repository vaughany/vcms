package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
	"vcms"
)

//go:embed templates/*.gohtml
//go:embed assets/img/*.png
var embeddedFiles embed.FS

const (
	cmdName     string = "VCMS - Receiver"
	cmdDesc     string = "Receives data from the Collector apps, creates a web page."
	cmdCodename string = "vcms-receiver"
)

var (
	debug                 bool   = false
	conciseDateTimeFormat string = "Mon Jan 2 2006, 15:04"
	nodes                        = make(map[string]*vcms.SystemData)
	cmdSubtitleHTML       string
	cmdFooterHTML         string
	receiverURL           string = "127.0.0.1:8080" // Don't put e.g. http:// at the start. Add this to docs.
	// logFile  string = vcms.MakeLogName(appCodename)
)

/*
HTMLData represents the data sent to the HTML template.
*/
type HTMLData struct {
	Title    template.HTML
	Subtitle template.HTML
	Footer   template.HTML
	Rows     []rowData
}

type rowData struct {
	Hostname       template.HTML
	Errors         template.HTML
	IPAddress      template.HTML
	Username       template.HTML
	FirstSeen      template.HTML
	LastSeen       template.HTML
	HostUptime     template.HTML
	OsVersion      template.HTML
	RebootRequired template.HTML
	LoadAvgs       template.HTML
	MemoryTotal    template.HTML
	MemoryFree     template.HTML
	SwapTotal      template.HTML
	SwapFree       template.HTML
	DiskTotal      template.HTML
	DiskFree       template.HTML
}

func init() {
	var version bool = false

	cmdSubtitleHTML = fmt.Sprintf("See <a href=\"https://%s\" target=\"_blank\">%s</a> for more info.", vcms.ProjectURL, vcms.ProjectURL)
	cmdFooterHTML = fmt.Sprintf("<strong>%s</strong> v%s (%s), built with %s, %s/%s. %s", vcms.AppTitle, vcms.AppVersion, vcms.AppDate, runtime.Version(), runtime.GOOS, runtime.GOARCH, cmdSubtitleHTML)

	flag.BoolVar(&debug, "d", debug, "Shows debugging info")
	flag.BoolVar(&version, "v", false, "Show version info and quit")
	flag.StringVar(&receiverURL, "r", receiverURL, "URL to run this application's web server on")
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

	if debug {
		go dumper(nodes)
	}

	// Handle files being served out of ./assets/img folder.
	// fileServer := http.FileServer(http.Dir("./assets/img/"))
	// http.Handle("/img/", http.StripPrefix("/img", fileServer))
	// ...or just handle the two files we actually want.
	http.HandleFunc("/img/logo.png", logoHandler)
	http.HandleFunc("/img/favicon.png", faviconHandler)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/api/announce", apiAnnounceHandler)
	http.HandleFunc("/api/ping", apiPingHandler)

	log.Printf("Running web server on http://%s.", receiverURL)
	log.Fatal(http.ListenAndServe(receiverURL, nil))
}

func dumper(nodes map[string]*vcms.SystemData) {
	for {
		log.Println("Dumping nodes:")
		if len(nodes) > 0 {
			var count int
			for index, node := range nodes {
				count++
				log.Printf("Node %d: %s: %v\n", count, index, node)
			}
		} else {
			log.Println("No nodes to dump.")
		}
		time.Sleep(time.Second * time.Duration(10))
	}
}
