package config

import (
	"encoding/json"
	"log"
	"os"
)

var appConfig = AppConfig{
	Port: "",
}

type AppConfig struct {
	Port string `json:"port"`
}

func init() {
	content, err := os.ReadFile("../configs/dev.json")
	if err != nil {
		log.Fatal("could not read config file")
	}

	err = json.Unmarshal(content, &appConfig)
	if err != nil {
		log.Fatal("could not parse config file content")
	}
}

func GetConfig() AppConfig {
	return AppConfig(appConfig)
}