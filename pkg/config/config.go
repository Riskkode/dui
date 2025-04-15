package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"dui/pkg/menu"
)

const (
	ScriptsKey = "scripts"
)

func BuildMenuItemFromConfig(location string) menu.MenuItem {

	configData, err := os.ReadFile(location)

	var Item menu.MenuItem
	err = json.Unmarshal(configData, &Item)
	if err != nil {
		panic(err)
	}

	return Item
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func SetScriptLocation(location string) error {
	if !PathExists(location) {
		return fmt.Errorf("path <%s> does not exist", location)
	}

	updateConfig(ScriptsKey, location)
	return nil
}

func updateConfig(key, value string) {
	AppConfig.Section("").Key(key).SetValue(value)

	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, configFileName)

	if err := AppConfig.SaveTo(configPath); err != nil {
		fmt.Println("Error saving updated config:", err)
	} else {
		fmt.Println("Config updated successfully.")
	}
}
