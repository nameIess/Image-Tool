// Package deps handles detection and validation of external dependencies.
// It checks for required tools like ImageMagick and Ghostscript without
// modifying system PATH, registry, or installing anything automatically.
package deps

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Status represents the availability status of a dependency.
type Status int

const (
	// StatusOK means the dependency is available and meets requirements.
	StatusOK Status = iota
	// StatusNotFound means the dependency was not found in PATH.
	StatusNotFound
	// StatusVersionError means the dependency exists but version check failed.
	StatusVersionError
)

// Dependency represents an external tool dependency.
type Dependency struct {
	Name        string
	Command     string
	VersionArgs []string
	MinVersion  string
	Status      Status
	Version     string
	Error       error
	DownloadURL string
	Description string
}

// CheckResult contains the results of all dependency checks.
type CheckResult struct {
	ImageMagick Dependency
	Ghostscript Dependency
	AllOK       bool
}

// Check verifies all required dependencies are available.
func Check() CheckResult {
	result := CheckResult{
		ImageMagick: checkImageMagick(),
		Ghostscript: checkGhostscript(),
	}
	result.AllOK = result.ImageMagick.Status == StatusOK && result.Ghostscript.Status == StatusOK
	return result
}

// checkImageMagick verifies ImageMagick is installed and meets minimum version.
func checkImageMagick() Dependency {
	dep := Dependency{
		Name:        "ImageMagick",
		Command:     "magick",
		VersionArgs: []string{"-version"},
		MinVersion:  "7.0.0",
		DownloadURL: "https://imagemagick.org/script/download.php",
		Description: "Required for image format conversion and compression",
	}

	// Check if magick command exists
	cmd := exec.Command("where", "magick")
	if err := cmd.Run(); err != nil {
		dep.Status = StatusNotFound
		dep.Error = fmt.Errorf("ImageMagick not found in PATH")
		return dep
	}

	// Get version
	cmd = exec.Command(dep.Command, dep.VersionArgs...)
	output, err := cmd.Output()
	if err != nil {
		dep.Status = StatusVersionError
		dep.Error = fmt.Errorf("failed to get version: %w", err)
		return dep
	}

	// Parse version from output like "Version: ImageMagick 7.1.0-62 Q16-HDRI x64..."
	version := parseImageMagickVersion(string(output))
	dep.Version = version

	if version == "" {
		dep.Status = StatusVersionError
		dep.Error = fmt.Errorf("could not parse version from output")
		return dep
	}

	// Check minimum version (must be 7.x)
	if !strings.HasPrefix(version, "7.") {
		dep.Status = StatusVersionError
		dep.Error = fmt.Errorf("version %s is below minimum required %s", version, dep.MinVersion)
		return dep
	}

	dep.Status = StatusOK
	return dep
}

// checkGhostscript verifies Ghostscript is installed.
func checkGhostscript() Dependency {
	dep := Dependency{
		Name:        "Ghostscript",
		Command:     "gswin64c",
		VersionArgs: []string{"-version"},
		MinVersion:  "",
		DownloadURL: "https://ghostscript.com/releases/gsdnld.html",
		Description: "Required for PDF processing and conversion",
	}

	// Try 64-bit version first
	cmd := exec.Command("where", "gswin64c")
	if err := cmd.Run(); err != nil {
		// Try 32-bit version
		dep.Command = "gswin32c"
		cmd = exec.Command("where", "gswin32c")
		if err := cmd.Run(); err != nil {
			dep.Status = StatusNotFound
			dep.Error = fmt.Errorf("Ghostscript not found in PATH")
			return dep
		}
	}

	// Get version
	cmd = exec.Command(dep.Command, dep.VersionArgs...)
	output, err := cmd.Output()
	if err != nil {
		// Ghostscript might output version to stderr or have different behavior
		// If command exists in PATH, we consider it OK
		dep.Status = StatusOK
		dep.Version = "detected"
		return dep
	}

	// Parse version from output
	version := parseGhostscriptVersion(string(output))
	if version != "" {
		dep.Version = version
	} else {
		dep.Version = "detected"
	}

	dep.Status = StatusOK
	return dep
}

// parseImageMagickVersion extracts version number from ImageMagick output.
func parseImageMagickVersion(output string) string {
	// Match "Version: ImageMagick 7.1.0-62" pattern
	re := regexp.MustCompile(`Version:\s*ImageMagick\s+(\d+\.\d+\.\d+(?:-\d+)?)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// parseGhostscriptVersion extracts version number from Ghostscript output.
func parseGhostscriptVersion(output string) string {
	// Match "GPL Ghostscript 10.02.1" or similar patterns
	re := regexp.MustCompile(`Ghostscript\s+(\d+\.\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// FormatStatus returns a formatted string showing dependency status.
func (d *Dependency) FormatStatus() string {
	switch d.Status {
	case StatusOK:
		if d.Version != "" && d.Version != "detected" {
			return fmt.Sprintf("✔ %s (%s)", d.Name, d.Version)
		}
		return fmt.Sprintf("✔ %s", d.Name)
	case StatusNotFound:
		return fmt.Sprintf("✗ %s - Not found", d.Name)
	case StatusVersionError:
		return fmt.Sprintf("⚠ %s - Version issue: %v", d.Name, d.Error)
	}
	return fmt.Sprintf("? %s - Unknown status", d.Name)
}

// GetMissingDependencyMessage returns a detailed message for missing dependencies.
func GetMissingDependencyMessage(result CheckResult) string {
	var b strings.Builder

	b.WriteString("Missing Dependencies\n")
	b.WriteString("═══════════════════════════════════════════════════════════\n\n")

	if result.ImageMagick.Status != StatusOK {
		b.WriteString(formatMissingDependency(result.ImageMagick))
		b.WriteString("\n")
	}

	if result.Ghostscript.Status != StatusOK {
		b.WriteString(formatMissingDependency(result.Ghostscript))
		b.WriteString("\n")
	}

	b.WriteString("───────────────────────────────────────────────────────────\n")
	b.WriteString("After installation, restart this application.\n")
	b.WriteString("Make sure the tools are added to your system PATH.\n")

	return b.String()
}

// formatMissingDependency formats a single missing dependency message.
func formatMissingDependency(dep Dependency) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("  %s\n", dep.Name))
	b.WriteString(fmt.Sprintf("  ─────────────────────────────\n"))
	b.WriteString(fmt.Sprintf("  Status:   %s\n", dep.FormatStatus()))
	b.WriteString(fmt.Sprintf("  Purpose:  %s\n", dep.Description))
	b.WriteString(fmt.Sprintf("  Download: %s\n", dep.DownloadURL))

	if dep.MinVersion != "" {
		b.WriteString(fmt.Sprintf("  Minimum:  v%s\n", dep.MinVersion))
	}

	return b.String()
}

// GetGhostscriptCommand returns the correct Ghostscript command for the system.
func GetGhostscriptCommand() string {
	// Try 64-bit first
	cmd := exec.Command("where", "gswin64c")
	if cmd.Run() == nil {
		return "gswin64c"
	}
	// Fall back to 32-bit
	return "gswin32c"
}
