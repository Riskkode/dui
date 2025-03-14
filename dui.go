package main

import (
	"dui/pkg/menu"
)

func main() {
	// Create a tree structure with scripts
	rootItems := menu.ParseConfig()
	root := menu.Builder(rootItems, nil)

	// Create and run the tree view
	treeView := menu.NewTreeView(root)
	defer treeView.TerminalCleanup()

	treeView.Run()
}
