package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
	vcms "vcms/internal"
)

func saveToPersistentStorage() {
	log.Println("Saving nodes to persistent storage.")
	saveJSONToFile()
}

func saveToPersistentStorageRegularly(persistentStorageSaveInterval int) {
	for {
		time.Sleep(time.Second * time.Duration(persistentStorageSaveInterval))
		log.Println("Saving nodes to persistent storage (regularly scheduled task).")
		saveJSONToFile()
	}
}

func loadFromPersistentStorage() {
	log.Println("Loading nodes from persistent storage.")
	loadJSONFromFile()
}

func saveJSONToFile() {
	newNodes := vcms.SystemDataPlusDateTime{
		SaveDateTime: time.Now(),
		SystemData:   nodes,
	}

	jsonBytes, err := json.Marshal(newNodes)
	if err != nil {
		log.Println("JSON data could not be marshalled for some reason, so could not save data.")
		return
	}

	file, err := os.Create(persistentStorage)
	if err != nil {
		log.Printf("%s could not be created, so could not save node data.\n", persistentStorage)
		return
	}
	defer file.Close()

	numBytes, err := file.Write(jsonBytes)
	if err != nil {
		log.Printf("%s exists but could not be written to, so could not save node data.\n", persistentStorage)
		return
	}

	log.Printf("Saving nodes as JSON: wrote %d bytes to %s.\n", numBytes, persistentStorage)
}

func loadJSONFromFile() {
	jsonBytes, err := os.ReadFile(persistentStorage)
	if err != nil {
		log.Printf("%s could not be read, so could not load node data.\n", persistentStorage)
		return
	}

	newNodes := vcms.SystemDataPlusDateTime{}
	err = json.Unmarshal(jsonBytes, &newNodes)
	if err != nil {
		log.Printf("JSON data from %s was loaded but could not be understood.\n", persistentStorage)
		return
	}

	nodes = newNodes.SystemData
	log.Printf("Loading nodes as JSON: read %d bytes from %s.\n", len(jsonBytes), persistentStorage)
}
