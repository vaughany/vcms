package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	vcms "vcms/internal"
)

func (ch *ContextHandler) apiAnnounceHandler(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}

	if ch.debug {
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

	// Check for API key, ensure it matches if length is non-zero.
	if len(ch.APIKey) > 0 {
		if data.Meta.APIKey != ch.APIKey {
			error := fmt.Sprintf("ERROR: %s: Collector API key '%s' does not match Receiver API key (not shown). Ignoring data.", r.Host, data.Meta.APIKey)
			log.Println(error)
			http.Error(w, error, http.StatusForbidden)
			return
		}
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
		cpu := vcms.CPU{
			Count: data.CPU.Count,
			Speed: data.CPU.Speed,
		}
		disk := vcms.TNF{
			Total: data.Disk.Total,
			Free:  data.Disk.Free,
		}
		memory := vcms.TNF{
			Total: data.Memory.Total,
			Free:  data.Memory.Free,
		}
		swap := vcms.TNF{
			Total: data.Swap.Total,
			Free:  data.Swap.Free,
		}
		meta := vcms.Meta{
			AppVersion: data.Meta.AppVersion,
			AppUptime:  data.Meta.AppUptime,
			Errors:     data.Meta.Errors,
		}
		nodes[data.Hostname] = &vcms.SystemData{
			Username:       data.Username, // These things pretty much stay the same.
			Hostname:       data.Hostname,
			IPAddress:      data.IPAddress,
			FirstSeen:      time.Now(),
			OSVersion:      data.OSVersion,
			CPU:            cpu,
			LastSeen:       time.Now(), // These things change.
			HostUptime:     data.HostUptime,
			RebootRequired: data.RebootRequired,
			LoadAvgs:       data.LoadAvgs,
			Memory:         memory,
			Swap:           swap,
			Disk:           disk,
			Meta:           meta,
		}
	} else {
		nodes[data.Hostname].LastSeen = time.Now()
		nodes[data.Hostname].HostUptime = data.HostUptime
		nodes[data.Hostname].RebootRequired = data.RebootRequired
		nodes[data.Hostname].LoadAvgs = data.LoadAvgs
		nodes[data.Hostname].Memory.Total = data.Memory.Total
		nodes[data.Hostname].Memory.Free = data.Memory.Free
		nodes[data.Hostname].Swap.Total = data.Swap.Total
		nodes[data.Hostname].Swap.Free = data.Swap.Free
		nodes[data.Hostname].Disk.Total = data.Disk.Total
		nodes[data.Hostname].Disk.Free = data.Disk.Free

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

func exportJSONHandler(w http.ResponseWriter, r *http.Request) {
	newNodes := vcms.SystemDataPlusDateTime{
		SaveDateTime: time.Now(),
		SystemData:   nodes,
	}

	jsonBytes, err := json.Marshal(newNodes)
	if err != nil {
		log.Println("JSON data could not be marshalled for some reason, so could not export data.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", jsonBytes)

	log.Println("Exporting nodes as JSON.")
}
