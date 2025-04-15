package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const configFileName = ".riskkode/dui/config.ini"

var AppConfig *ini.File

func init() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	configPath := filepath.Join(homeDir, configFileName)

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create config file and parent dir if needed
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			fmt.Println("Error creating config directory:", err)
			os.Exit(1)
		}

		cfg := ini.Empty()
		cfg.Section("").Key("data_dir").SetValue("/path/to/data")
		cfg.Section("").Key("logs_dir").SetValue("/path/to/logs")

		if err := cfg.SaveTo(configPath); err != nil {
			fmt.Println("Error writing default config:", err)
			os.Exit(1)
		}

		AppConfig = cfg
		fmt.Println("Default config file created at:", configPath)
		return
	}

	// Load existing config
	cfg, err := ini.Load(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	AppConfig = cfg
	fmt.Println("Config loaded from:", configPath)
}
