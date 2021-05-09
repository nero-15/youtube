package config

import (
	"log"
	"os"

	ini "gopkg.in/ini.v1"
)

// ConfigList is api key struct
type ConfigList struct {
	APIKey string
	Port   string
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
		APIKey: cfg.Section("youtube").Key("key").String(),
		Port:   cfg.Section("web").Key("port").String(),
	}
}
