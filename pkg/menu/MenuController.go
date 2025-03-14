package menu

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"dui/colours"
	"golang.org/x/term"
)

// TreeView manages the tree display and navigation
type TreeView struct {
	Root          *TreeNode
	VisibleNodes  []*TreeNode
	SelectedIdx   int
	TermWidth     int
	TermHeight    int
	Running       bool
	OldState      *term.State
	StatusMessage string
}

// NewTreeView creates a new tree view
func NewTreeView(root *TreeNode) *TreeView {
	width, height, _ := term.GetSize(int(os.Stdout.Fd()))
	if width <= 0 {
		width = 80 // Default width
	}
	if height <= 0 {
		height = 24 // Default height
	}

	tv := &TreeView{
		Root:        root,
		SelectedIdx: 0,
		TermWidth:   width,
		TermHeight:  height,
		Running:     true,
	}

	// Expand root by default
	root.Expanded = true

	// Build initial visible nodes list
	tv.rebuildVisibleNodes()

	return tv
}

// rebuildVisibleNodes rebuilds the list of nodes that should be displayed
func (tv *TreeView) rebuildVisibleNodes() {
	tv.VisibleNodes = []*TreeNode{}
	tv.traverseVisible(tv.Root, func(node *TreeNode) {
		tv.VisibleNodes = append(tv.VisibleNodes, node)
	})

	// Ensure selected index is valid
	if len(tv.VisibleNodes) == 0 {
		tv.SelectedIdx = 0
	} else if tv.SelectedIdx >= len(tv.VisibleNodes) {
		tv.SelectedIdx = len(tv.VisibleNodes) - 1
	}
}

// traverseVisible traverses nodes that should be visible
func (tv *TreeView) traverseVisible(node *TreeNode, fn func(*TreeNode)) {
	fn(node)
	if node.Expanded {
		for _, child := range node.Children {
			tv.traverseVisible(child, fn)
		}
	}
}

// Draw renders the tree view
func (tv *TreeView) Draw() {
	clearScreen()

	// Display title and controls
	colours.Fprint(colours.Cyan, "Developer User Interface")
	colours.Fprint(colours.Yellow, "Controls: ↑/↓ = navigate, Enter/Space = expand/collapse, x = execute script, q = quit")
	if runtime.GOOS == "windows" {
		colours.Fprint(colours.Yellow, "Alternate controls: w/s = up/down, e = expand/collapse")
	}
	fmt.Println(strings.Repeat("─", tv.TermWidth))

	// Calculate visible range based on current selection
	startIdx := 0
	maxVisible := tv.TermHeight - 5 // Account for title, controls, separator, and status line

	if len(tv.VisibleNodes) > maxVisible {
		// Center the selected item if possible
		half := maxVisible / 2
		if tv.SelectedIdx > half {
			startIdx = tv.SelectedIdx - half
		}

		if startIdx+maxVisible > len(tv.VisibleNodes) {
			startIdx = len(tv.VisibleNodes) - maxVisible
		}
	}

	endIdx := startIdx + maxVisible
	if endIdx > len(tv.VisibleNodes) {
		endIdx = len(tv.VisibleNodes)
	}

	// Draw visible nodes
	for i := startIdx; i < endIdx; i++ {
		node := tv.VisibleNodes[i]

		// Determine prefix based on expanded state
		prefix := "  "
		if len(node.Children) > 0 {
			if node.Expanded {
				prefix = "- "
			} else {
				prefix = "+ "
			}
		}

		// Create indentation
		indent := ""
		for j := 0; j < node.Depth; j++ {
			indent += "  "
		}
		label := fmt.Sprintf("%s%s%s%s", indent, prefix, "", node.Label)
		// Add script indicator if executable
		if node.IsExecutable {
			execIndicator := "[x] "
			label = fmt.Sprint(colours.Yellow, fmt.Sprintf("%s%s%s%s", indent, prefix, execIndicator, node.Label), colours.Reset)
		}

		// Format with highlighting if selected
		if i == tv.SelectedIdx {
			// Try to use reverse video, or fallback to a simple indicator
			fmt.Printf("\033[7m%s\033[0m\n", label)
		} else {
			fmt.Println(label)
		}
	}

	// Status line
	fmt.Println(strings.Repeat("─", tv.TermWidth))
	if tv.StatusMessage != "" {
		fmt.Println(tv.StatusMessage)
		tv.StatusMessage = "" // Clear status message after displaying
	} else if len(tv.VisibleNodes) > 0 {
		selectedNode := tv.VisibleNodes[tv.SelectedIdx]
		status := fmt.Sprintf("Selected: %s", selectedNode.Label)
		if selectedNode.IsExecutable {
			status += fmt.Sprintf(" (Press 'x' to execute: %s)", selectedNode.ScriptPath)
		}
		fmt.Println(status)
	} else {
		colours.Err("No nodes available")
	}
}

// MoveUp moves selection up
func (tv *TreeView) MoveUp() {
	if tv.SelectedIdx > 0 {
		tv.SelectedIdx--
		tv.Draw()
	}
}

// MoveDown moves selection down
func (tv *TreeView) MoveDown() {
	if tv.SelectedIdx < len(tv.VisibleNodes)-1 {
		tv.SelectedIdx++
		tv.Draw()
	}
}

// ToggleExpand expands or collapses the selected node
func (tv *TreeView) ToggleExpand() {
	if len(tv.VisibleNodes) == 0 {
		return
	}

	node := tv.VisibleNodes[tv.SelectedIdx]
	if len(node.Children) > 0 {
		node.Expanded = !node.Expanded
		tv.rebuildVisibleNodes()
		tv.Draw()
	}
}

// ExecuteSelected executes the script of the selected node if it's executable
func (tv *TreeView) ExecuteSelected() {
	if len(tv.VisibleNodes) == 0 {
		return
	}

	node := tv.VisibleNodes[tv.SelectedIdx]
	if node.IsExecutable {
		// Restore terminal to normal mode for script execution
		term.Restore(int(os.Stdin.Fd()), tv.OldState)

		// Clear screen and show what's being executed
		clearScreen()
		colours.Info(fmt.Sprintf("Executing script: %s\n\n", node.ScriptPath))

		// Execute the script
		err := node.ExecuteScript()

		// After execution, wait for user to press a key
		colours.Info(fmt.Sprintf("\nScript execution completed. Press enter to return to the menu..."))
		waitForKeyPress()

		// Return to raw mode for tree view
		var err2 error
		tv.OldState, err2 = term.MakeRaw(int(os.Stdin.Fd()))
		if err2 != nil {
			fmt.Println("Error returning to raw mode:", err2)
			tv.Running = false
			return
		}

		// Set status message based on execution result
		if err != nil {
			tv.StatusMessage = colours.SErr(fmt.Sprintf("Error executing script: %v", err))
		} else {
			tv.StatusMessage = colours.SOk(fmt.Sprintf("Successfully executed: %s", node.ScriptPath))
		}

		tv.Draw()
	}
}

// waitForKeyPress waits for a single keypress
func waitForKeyPress() {
	// Create a temporary buffer to read a single key
	buf := make([]byte, 1)
	os.Stdin.Read(buf)
}

// readInput reads a sequence of bytes from stdin with timeout
func readInputWithTimeout(timeout time.Duration) []byte {
	// Make a buffer to read into
	buf := make([]byte, 10)
	result := []byte{}

	// Set stdin to non-blocking mode
	os.Stdin.SetReadDeadline(time.Now().Add(timeout))

	// Read bytes until timeout
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			break
		}
		result = append(result, buf[:n]...)

		// If we've read enough for most key sequences, we can stop
		if len(result) >= 3 {
			break
		}

		// Reset the deadline for more bytes
		os.Stdin.SetReadDeadline(time.Now().Add(timeout))
	}

	// Reset stdin to blocking mode
	os.Stdin.SetReadDeadline(time.Time{})

	return result
}

// Run starts the tree view
func (tv *TreeView) Run() {
	// Save current terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return
	}
	tv.OldState = oldState

	// Restore terminal when done
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	tv.Draw()

	for tv.Running {
		// Read a single byte first
		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			continue
		}

		// Process input
		switch buf[0] {
		case 3: // Ctrl+C
			tv.Running = false

		case 'q', 'Q': // Quit
			tv.Running = false

		case 'w', 'W': // Alternative up
			tv.MoveUp()

		case 's', 'S': // Alternative down
			tv.MoveDown()

		case 'e', 'E', ' ': // Alternative expand/collapse
			tv.ToggleExpand()

		case 'x', 'X': // Execute script
			tv.ExecuteSelected()

		case 13: // Enter
			tv.ToggleExpand()

		case 27: // Escape sequence (might be arrow keys)
			// Try to read more bytes with a short timeout
			moreBuf := make([]byte, 2)
			n, _ := os.Stdin.Read(moreBuf)

			if n >= 2 {
				// Check for arrow keys
				if moreBuf[0] == 91 { // [
					switch moreBuf[1] {
					case 65: // Up arrow
						tv.MoveUp()
					case 66: // Down arrow
						tv.MoveDown()
					case 67, 68: // Right or Left arrow
						tv.ToggleExpand()
					}
				}
			}
		}
	}
}

// TerminalCleanup restores the terminal
func (tv *TreeView) TerminalCleanup() {
	if tv.OldState != nil {
		term.Restore(int(os.Stdin.Fd()), tv.OldState)
	}
}

// clearScreen clears the terminal screen
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
