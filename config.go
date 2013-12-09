package main

import (
	"encoding/json"
	"os"
	"log"
)

type StorageConfiguration struct {
	Backend  *string
	Location *string
	Port     int
}

type Configuration struct {
	Storage     *StorageConfiguration
	WwwRoot     *string
}

func ReadConfig(filename string) *Configuration {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	config := &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return config
}
