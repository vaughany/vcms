package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
	"vcms"

	"code.cloudfoundry.org/bytefmt"
	"github.com/hako/durafmt"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Root accessed. Redirecting...")
	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	var data HTMLData

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	// Sorting the map into order by hostname.
	keys := make([]string, 0, len(nodes))
	for key := range nodes {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Build the HTML.
	data.Title = vcms.AppTitle
	// data.Subtitle = "Something Something Darkside"
	data.Footer = template.HTML(cmdFooterHTML)
	for _, key := range keys {
		var row rowData
		if len(nodes[key].Meta.Errors) > 0 {
			row.Errors = fmt.Sprintf(" <span style=\"color: red;\" title=\"%s\">ERROR!</span>", strings.Join(nodes[key].Meta.Errors, "\n"))
		}
		row.Hostname = nodes[key].Hostname
		row.IPAddress = nodes[key].IPAddress
		row.Username = nodes[key].Username
		row.FirstSeen = template.HTML(fmt.Sprintf("%s <span class=\"has-text-grey-light\"><small>(%s ago)</small></span>", nodes[key].FirstSeen.Format(conciseDateTimeFormat), durafmt.Parse(time.Since(nodes[key].FirstSeen).Round(time.Second))))
		row.LastSeen = template.HTML(fmt.Sprintf("%s <span class=\"has-text-grey-light\"><small>(%s ago)</small></span>", nodes[key].LastSeen.Format(conciseDateTimeFormat), durafmt.Parse(time.Since(nodes[key].LastSeen).Round(time.Second))))
		row.HostUptime = nodes[key].HostUptime
		row.OSVersion = nodes[key].OSVersion
		row.CPU = getCPUHTML(nodes[key])
		if nodes[key].RebootRequired {
			row.RebootRequired = "yes"
		} else {
			row.RebootRequired = "no"
		}

		loadAvgsString := []string{}
		for _, loadAvg := range nodes[key].LoadAvgs {
			loadAvgsString = append(loadAvgsString, fmt.Sprintf("%.2f", loadAvg))
		}
		row.LoadAvgs = strings.Join(loadAvgsString, " ")

		row.MemoryTotal = bytefmt.ByteSize(uint64(nodes[key].MemoryTotal * bytefmt.KILOBYTE))
		row.MemoryFree = bytefmt.ByteSize(uint64(nodes[key].MemoryFree * bytefmt.KILOBYTE))
		row.SwapTotal = bytefmt.ByteSize(uint64(nodes[key].SwapTotal * bytefmt.KILOBYTE))
		row.SwapFree = bytefmt.ByteSize(uint64(nodes[key].SwapFree * bytefmt.KILOBYTE))

		row.DiskTotal = bytefmt.ByteSize(uint64(nodes[key].DiskTotal * bytefmt.KILOBYTE))
		percentage := float64(nodes[key].DiskFree) / float64(nodes[key].DiskTotal) * 100
		row.DiskFree = template.HTML(fmt.Sprintf("%s <span class=\"has-text-grey-light\"><small>(%.1f%%)</small></span>", bytefmt.ByteSize(uint64(nodes[key].DiskFree*bytefmt.KILOBYTE)), percentage))

		row.OSImage = getOSImage(nodes[key])

		data.Rows = append(data.Rows, row)
	}

	tmpl := template.Must(template.New("layout.gohtml").Funcs(funcMap).ParseFS(embeddedFiles, "templates/layout.gohtml", "templates/dashboard.gohtml"))
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Dashboard accessed.")
}

// func logoHandler(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./assets/img/logo.png")
// }

// func faviconHandler(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "./assets/img/favicon.png")
// }

func getCPUHTML(node *vcms.SystemData) string {
	if node.CPUCount == 0 {
		return "-"
	}

	return fmt.Sprintf("%d, %s", node.CPUCount, node.CPUSpeed)
}

func getOSImage(node *vcms.SystemData) string {
	input := strings.ToLower(node.OSVersion)
	switch {
	case strings.Contains(input, "arch"):
		return "arch"
	case strings.Contains(input, "centos"):
		return "centos"
	case strings.Contains(input, "debian"):
		return "debian"
	case strings.Contains(input, "elementary"):
		return "elementary"
	case strings.Contains(input, "fedora"):
		return "fedora"
	case strings.Contains(input, "kali"):
		return "kali"
	// case strings.Contains(input, "kubuntu"): // Kubuntu identifies itself as Ubuntu.
	// 	return "kubuntu"
	// case strings.Contains(input, "lubuntu"): // Lubuntu identifies itself as Ubuntu.
	// 	return "lubuntu"
	case strings.Contains(input, "manjaro"):
		return "manjaro"
	case strings.Contains(input, "mint"):
		return "mint"
	// case strings.Contains(input, "mx"): // MX identifies itself as Debian.
	// 	return "mx"
	case strings.Contains(input, "opensuse"):
		return "opensuse"
	case strings.Contains(input, "oracle solaris"):
		return "solaris"
	case strings.Contains(input, "oracle"):
		return "oracle"
	case strings.Contains(input, "red hat"):
		return "redhat"
	case strings.Contains(input, "ubuntu"):
		return "ubuntu"
	case strings.Contains(input, "windows"):
		return "windows"
	case strings.Contains(input, "zorin"):
		return "zorin"
	default:
		return "generic"
	}
}

func saveToPersistentStorageHandler(w http.ResponseWriter, r *http.Request) {
	saveToPersistentStorage()

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

func loadFromPersistentStorageHandler(w http.ResponseWriter, r *http.Request) {
	loadFromPersistentStorage()

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}

func nodeRemoveHandler(w http.ResponseWriter, r *http.Request) {
	URLNodes, ok := r.URL.Query()["node"]
	if !ok || len(URLNodes[0]) < 1 {
		log.Println("URL parameter 'node' is missing.")
		return
	}
	node := URLNodes[0]

	// log.Printf("URL param 'node' is: '%s'.\n", string(node))
	log.Printf("Removing node '%s'.\n", string(node))

	delete(nodes, node)

	http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
}
