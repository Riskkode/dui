package menu

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var DefaultMenu = MenuItem{
	Name:     "dui",
	ItemType: ItemTypeExe,
	Children: []MenuItem{
		{
			Name:     "Settings",
			ItemType: ItemTypeNode,
			Children: []MenuItem{},
		},
	},
}

const (
	ItemTypeExe  = "executable"
	ItemTypeNode = "node"
)

type MenuItem struct {
	Name     string     `json:"name"`
	ItemType string     `json:"type"`
	Path     string     `json:"path"`
	Children []MenuItem `json:"children"`
}

// TreeNode represents a node in the tree structure
type TreeNode struct {
	Label        string
	Children     []*TreeNode
	Expanded     bool
	Parent       *TreeNode
	Depth        int
	Data         interface{} // Optional additional data
	IsExecutable bool        // Flag to indicate if node is an executable script
	ScriptPath   string      // Path to the script to execute
}

// Builder uses a menu item and appends it to the parent node.
func Builder(item MenuItem, parent *TreeNode) *TreeNode {
	var node *TreeNode

	if parent == nil {
		node = NewTreeNode(item.Name, nil)
	} else {
		if item.ItemType == "executable" {
			node = parent.AddExecutableChild(item.Name, item.Path)
		} else {
			node = parent.AddChild(item.Name)
		}
	}

	for _, child := range item.Children {
		Builder(child, node)
	}
	return node
}

// NewTreeNode creates a new tree node
func NewTreeNode(label string, parent *TreeNode) *TreeNode {
	depth := 0
	if parent != nil {
		depth = parent.Depth + 1
	}
	return &TreeNode{
		Label:        label,
		Children:     []*TreeNode{},
		Expanded:     false,
		Parent:       parent,
		Depth:        depth,
		IsExecutable: false,
	}
}

// AddChild adds a child node to the current node
func (n *TreeNode) AddChild(label string) *TreeNode {
	child := NewTreeNode(label, n)
	n.Children = append(n.Children, child)
	return child
}

// AddExecutableChild adds a child node that can execute a bash script
func (n *TreeNode) AddExecutableChild(label string, scriptPath string) *TreeNode {
	child := n.AddChild(label)
	child.IsExecutable = true
	child.ScriptPath = scriptPath
	return child
}

// ExecuteScript runs the bash script associated with this node
func (n *TreeNode) ExecuteScript() error {
	if !n.IsExecutable {
		return fmt.Errorf("node %s is not executable", n.Label)
	}

	// Determine shell to use based on OS
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("bash", n.ScriptPath)
	} else {
		cmd = exec.Command(n.ScriptPath)
	}

	// Set up command to run in the current terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
