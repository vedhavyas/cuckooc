package cuckooc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config represents the configuration for the service
type Config struct {
	Debug        bool   `json:"debug"`
	BackupFolder string `json:"backup_folder"`
	TCP          string `json:"tcp"`
	UDP          string `json:"udp"`
}

// LoadConfig loads the service configuration from the file provided
func LoadConfig(file string) (config Config, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return config, fmt.Errorf("failed to read the config file: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshall config data: %v", err)
	}

	return config, nil
}
