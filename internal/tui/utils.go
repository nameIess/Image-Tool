package tui

import (
	"os/exec"
	"runtime"
)

// checkMagick verifies ImageMagick is installed
func checkMagick() bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "magick")
	} else {
		cmd = exec.Command("which", "magick")
	}
	return cmd.Run() == nil
}

// checkGhostscript verifies Ghostscript is installed
func checkGhostscript() bool {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Try both 64-bit and 32-bit versions
		cmd = exec.Command("where", "gswin64c")
		if cmd.Run() == nil {
			return true
		}
		cmd = exec.Command("where", "gswin32c")
	} else {
		cmd = exec.Command("which", "gs")
	}
	return cmd.Run() == nil
}
