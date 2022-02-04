package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	Machines map[string]string `json:"machines"`
}

func readConfigurationFile(name *string) Configuration {
	content, err := os.ReadFile(*name)
	if err != nil {
		log.Fatalf("Could not read configuration file: %s\n", err)
	}

	valid := json.Valid(content)
	if !valid {
		log.Fatalf("The JSON of the configuration file is invalid")
	}

	var config Configuration
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Could not parse JSON: %s\n", err)
	}

	return config
}
