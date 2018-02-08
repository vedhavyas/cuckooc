package cuckooc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config represents the configuration for the service
type Config struct {
	BackupFolder string `json:"backup_folder"`
}

// loadConfig loads the service configuration from the file provided
func loadConfig(file string) (config Config, err error) {
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
