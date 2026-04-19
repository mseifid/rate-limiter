package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	App struct {
		Port string `json:"port"`
	} `json:"app"`
}

func LoadConfig() (*Config, error) {
	// TODO: This should be kept in memory for app lifetime
	content, err := os.ReadFile("../configs/dev.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}