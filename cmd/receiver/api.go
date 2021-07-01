package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"vcms"
)

func apiAnnounceHandler(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	if debug {
		log.Printf("Received %s from %s", jsonBytes, r.Host)
	} else {
		log.Printf("Received data from %s", r.Host)
	}

	// Unmarshall JSON.
	var data vcms.SystemData
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		log.Panic(err)
	}

	// // Bail if any of the values are missing.
	// if !data.CheckPopulatedValues() {
	// 	return
	// }

	// If the hostname's not been seen before, do a full update. Else, just update the 'last seen' time.
	// TODO: We're assuming that hostnames are unique. They might not be. Maybe the map's key should be 'hostname+ipaddress'?
	if _, ok := nodes[data.Hostname]; !ok {
		nodes[data.Hostname] = &vcms.SystemData{
			Username:  data.Username,
			Hostname:  data.Hostname,
			IPAddress: data.IPAddress,
			FirstSeen: time.Now(),
			LastSeen:  time.Now(),
		}
	} else {
		nodes[data.Hostname].LastSeen = time.Now()
	}
}

func apiPingHandler(w http.ResponseWriter, r *http.Request) {
	result := `{"result":"pong"}`
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, result)
	log.Println("Received ping. Returned " + result)
}
