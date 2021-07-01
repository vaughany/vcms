package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"
	"vcms"

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
	data.Title = template.HTML(vcms.AppTitle)
	// data.Subtitle = template.HTML("Something Something Darkside")
	data.Footer = template.HTML(cmdFooterHTML)
	for _, key := range keys {
		var row rowData
		row.Hostname = template.HTML(nodes[key].Hostname)
		row.IPAddress = template.HTML(nodes[key].IPAddress)
		row.Username = template.HTML(nodes[key].Username)
		row.FirstSeen = template.HTML(fmt.Sprintf("%s <span class=\"has-text-grey-light\"><small>(%s ago)</small></span>", nodes[key].FirstSeen.Format(conciseDateTimeFormat), durafmt.Parse(time.Since(nodes[key].FirstSeen).Round(time.Second))))
		row.LastSeen = template.HTML(fmt.Sprintf("%s <span class=\"has-text-grey-light\"><small>(%s ago)</small></span>", nodes[key].LastSeen.Format(conciseDateTimeFormat), durafmt.Parse(time.Since(nodes[key].LastSeen).Round(time.Second))))
		data.Rows = append(data.Rows, row)
	}

	tmpl := template.Must(template.New("layout.gohtml").Funcs(funcMap).ParseFS(embeddedFiles, "templates/layout.gohtml", "templates/dashboard.gohtml"))
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Dashboard accessed.")
}

func logoHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/img/logo.png")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./assets/img/favicon.png")
}
