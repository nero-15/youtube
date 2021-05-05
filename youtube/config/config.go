package config

import (
	"log"
	"os"

	ini "gopkg.in/ini.v1"
)

// ConfigList is api key struct
type ConfigList struct {
	ApiKey  string
	BaseURL string
	Port    string
}

// Config is ConfigList
var Config ConfigList

func init() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		ApiKey:  cfg.Section("youtube").Key("key").String(),
		BaseURL: cfg.Section("youtube").Key("baseURL").String(),
		Port:    cfg.Section("web").Key("port").String(),
	}
}
