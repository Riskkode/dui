package main

import (
	"flag"

	"dui/pkg/config"
	"dui/pkg/menu"
)

func main() {

	flag.Func("scripts", "Type a fully qualified directory location to import .sh scripts and their layout.json - eg dui -config -scripts C:/User/Name/MyScripts", config.SetScriptLocation)

	flag.Parse()

	ScriptsLayoutJSON := config.AppConfig.Section("").Key(config.ScriptsKey).String() + "/layout.json"

	// Retrieve custom user configs
	scripts := config.BuildMenuItemFromConfig(ScriptsLayoutJSON)
	root := menu.Builder(menu.DefaultMenu, nil)
	menu.Builder(scripts, root)

	treeView := menu.NewTreeView(root)
	defer treeView.TerminalCleanup()

	treeView.Run()
}
