package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DataFilePath string `json:"dataFilePath"`
}

func GetConfig() (Config, error) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		return Config{}, err
	}

	var config *Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return Config{}, err
	}

	return *config, nil
}
