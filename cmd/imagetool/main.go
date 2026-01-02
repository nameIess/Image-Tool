//go:generate goversioninfo -o resource.syso versioninfo.json

package main

import (
	"fmt"
	"os"

	"imagetool/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Set terminal title
	fmt.Print("\033]0;Image-Tool\007")

	p := tea.NewProgram(tui.NewApp(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
