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
		log.Printf("Received %s from %s", jsonBytes, r.RemoteAddr)
	} else {
		log.Printf("Received data from %s", r.RemoteAddr)
	}

	// Unmarshall JSON.
	var data vcms.SystemData
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		log.Panic(err)
	}

	// Check for version parity.
	if data.Meta.AppVersion != vcms.AppVersion {
		error := fmt.Sprintf("ERROR: %s: Collector version %s does not match Receiver version %s. Ignoring data.", r.Host, data.Meta.AppVersion, vcms.AppVersion)
		log.Println(error)
		http.Error(w, error, http.StatusUnprocessableEntity)
		return
	}

	// It's all good.
	w.WriteHeader(200)

	// // Bail if any of the values are missing.
	// if !data.CheckPopulatedValues() {
	// 	return
	// }

	// If the hostname's not been seen before, do a full update. Else, just update the 'last seen' time.
	// TODO: Of course, now we're updating more than just the 'lsat seen' time and the code block below is repeating code.
	// TODO: We're assuming that hostnames are unique. They might not be. Maybe the map's key should be 'hostname+ipaddress'?
	if _, ok := nodes[data.Hostname]; !ok {
		nodes[data.Hostname] = &vcms.SystemData{
			Username:  data.Username, // These things pretty much stay the same.
			Hostname:  data.Hostname,
			IPAddress: data.IPAddress,
			FirstSeen: time.Now(),
			OsVersion: data.OsVersion,
			// LastSeen:       time.Now(), // These things change.
			// HostUptime:     data.HostUptime,
			// RebootRequired: data.RebootRequired,
			// LoadAvgs:       data.LoadAvgs,
			// MemoryTotal:    data.MemoryTotal,
			// MemoryFree:     data.MemoryFree,
			// SwapTotal:      data.SwapTotal,
			// SwapFree:       data.SwapFree,
			// DiskTotal:      data.DiskTotal,
			// DiskFree:       data.DiskFree,
		}
	} else {
		nodes[data.Hostname].LastSeen = time.Now()
		nodes[data.Hostname].HostUptime = data.HostUptime
		nodes[data.Hostname].RebootRequired = data.RebootRequired
		nodes[data.Hostname].LoadAvgs = data.LoadAvgs
		nodes[data.Hostname].MemoryTotal = data.MemoryTotal
		nodes[data.Hostname].MemoryFree = data.MemoryFree
		nodes[data.Hostname].SwapTotal = data.SwapTotal
		nodes[data.Hostname].SwapFree = data.SwapFree
		nodes[data.Hostname].DiskTotal = data.DiskTotal
		nodes[data.Hostname].DiskFree = data.DiskFree

		nodes[data.Hostname].Meta.AppVersion = data.Meta.AppVersion
		nodes[data.Hostname].Meta.AppUptime = data.Meta.AppUptime
		nodes[data.Hostname].Meta.Errors = data.Meta.Errors
	}
}

func apiPingHandler(w http.ResponseWriter, r *http.Request) {
	result := `{"result":"pong"}`
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, result)
	log.Println("Received ping. Returned " + result)
}
