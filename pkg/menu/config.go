package menu

import (
	"encoding/json"
	"os"
)

func ParseConfig() MenuItem {
	configData, err := os.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}

	var rootItem MenuItem
	err = json.Unmarshal(configData, &rootItem)
	if err != nil {
		panic(err)
	}

	return rootItem
}
