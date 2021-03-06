// Package Receiver receives data from the Collector apps, creates a web page.
package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	vcms "vcms/internal"
)

//go:embed templates/*.gohtml
//go:embed assets/img/*.png
//go:embed assets/img/*.html
var embeddedFiles embed.FS

// Something like this, to ditch global state?
//   https://stackoverflow.com/a/46517000/254146
var (
	nodes = make(map[string]*vcms.SystemData)
)

const (
	conciseDateTimeFormat = "Mon Jan 2 2006, 15:04"
	persistentStorage     = "nodes.json"
)

/*
HTMLData represents the data sent to the HTML template.
*/
type HTMLData struct {
	Title    string
	Subtitle string
	Footer   template.HTML
	Rows     []rowData
	RowCount int
}

type rowData struct {
	Hostname       string
	Errors         string
	IPAddress      string
	Username       string
	FirstSeen      template.HTML
	LastSeen       template.HTML
	HostUptime     string
	OSVersion      string
	OSImage        string
	CPU            string
	RebootRequired string
	LoadAvgs       string
	MemoryTotal    string
	MemoryFree     string
	SwapTotal      string
	SwapFree       string
	DiskTotal      string
	DiskFree       template.HTML
}

type ContextHandler struct {
	debug  bool
	APIKey string
}

func main() {
	const (
		cmdName                       = "VCMS - Receiver"
		cmdDesc                       = "Receives data from the Collector apps, creates a web page."
		cmdCodename                   = "vcms-receiver"
		persistentStorageSaveInterval = 60 // TODO: make configurable.
	)

	var (
		debug       = false
		version     = false
		keyGen      = false
		receiverURL = "127.0.0.1:8080" // Don't put e.g. http:// at the start. Add this to docs.
		logName     = vcms.CreateLogName(vcms.LogFolder, cmdCodename)
		APIKey      = ""
	)

	flag.StringVar(&APIKey, "apikey", APIKey, "API key, if you want the Collector to prove it's legit")
	flag.BoolVar(&debug, "d", debug, "Shows debugging info, incl all JSON being sent")
	flag.BoolVar(&keyGen, "k", false, "Quickly generate a few random keys")
	flag.StringVar(&receiverURL, "r", receiverURL, "URL to run this application's web server on")
	flag.BoolVar(&version, "v", false, "Show version info and quit")
	flag.Parse()

	if version {
		fmt.Println(vcms.Version(cmdName))
		os.Exit(0)
	}

	if keyGen {
		vcms.ShowRandomKeys()
		os.Exit(0)
	}

	// Opens a logfile for writing and also outputs to stdout.
	vcms.CheckLogFolder(vcms.LogFolder)
	lf, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		log.Print(err)
		log.Printf("Error: Couldn't open '%s' for writing, so logging to stdout only.\n", logName)
	} else {
		defer lf.Close()
		log.SetOutput(io.MultiWriter(os.Stdout, lf))
		log.Printf("Startup: Logging to stdout and '%s'.\n", logName)
	}

	shutdownHandler()

	log.Println(vcms.Version(cmdName))
	log.Printf("%s \n", cmdDesc)
	log.Printf("%s \n", vcms.AppDesc)

	if debug {
		go dumper()
	}

	// Load the nodes from a file.
	loadFromPersistentStorage()

	// Save all the nodes out to a file regularly.
	// saveToPersistentStorageTicker := time.NewTicker(time.Second * time.Duration(persistentStorageSaveInterval))
	// defer saveToPersistentStorageTicker.Stop()
	// saveToPersistentStorageDone := make(chan bool)
	// go saveToPersistentStorageRegularly(saveToPersistentStorageTicker, saveToPersistentStorageDone)
	go saveToPersistentStorageRegularly(persistentStorageSaveInterval)

	// Handle files being served out of ./assets/img folder.
	fileServer := http.FileServer(http.Dir("./assets/img/"))
	http.Handle("/img/", http.StripPrefix("/img", fileServer))
	// ...or just handle the two files we actually want.
	// http.HandleFunc("/img/logo.png", logoHandler)
	// http.HandleFunc("/img/favicon.png", faviconHandler)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/dashboard", dashboardHandler)
	http.HandleFunc("/dashboard/full", dashboardHandler)
	http.HandleFunc("/hosts", hostsHandler)
	http.HandleFunc("/host/", hostHandler) // Note the trailing '/'.

	contextHandler := &ContextHandler{debug: debug, APIKey: APIKey}
	http.HandleFunc("/api/announce", contextHandler.apiAnnounceHandler)
	http.HandleFunc("/api/ping", apiPingHandler)

	http.HandleFunc("/save", saveToPersistentStorageHandler)
	http.HandleFunc("/load", loadFromPersistentStorageHandler)
	http.HandleFunc("/node/remove/", nodeRemoveHandler) // Note the trailing '/'.

	http.HandleFunc("/export/json", exportJSONHandler)

	APIKeyInfo := ""
	if len(APIKey) > 0 {
		APIKeyInfo = fmt.Sprintf(" --apikey %s", APIKey)
	}
	log.Printf("Running web server on http://%s.", receiverURL)
	log.Printf("To connect a Collector, run: './collector -r http://%s%s'.", receiverURL, APIKeyInfo)
	log.Fatal(http.ListenAndServe(receiverURL, nil))
}

func dumper() {
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

// https://stackoverflow.com/a/12571099/254146
func shutdownHandler() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		log.Println("Close signal detected: shutting down.")

		log.Println("Shutdown: saving nodes to persistent storage.")
		saveToPersistentStorage()
		// TODO: e.g. notify Slack, send an email etc.

		log.Println("Shutdown: bye bye.")
		log.Println()
		log.Println()
		log.Println()
		os.Exit(0)
	}()
}

func makeHTMLFooter() string {
	return fmt.Sprintf("<strong>%s</strong> v%s (%s), built with %s, %s/%s. See <a href=\"https://%s\" target=\"_blank\">%s</a> for more info.",
		vcms.AppTitle, vcms.AppVersion, vcms.AppDate, runtime.Version(), runtime.GOOS, runtime.GOARCH, vcms.ProjectURL, vcms.ProjectURL)
}
