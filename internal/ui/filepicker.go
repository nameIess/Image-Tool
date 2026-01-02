package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"imagetool/internal/core"
	"imagetool/internal/logging"
)

// FilePickerMode determines what files to show
type FilePickerMode int

const (
	FilePickerPDF FilePickerMode = iota
	FilePickerImage
	FilePickerAll
)

// FileEntry represents a file or directory
type FileEntry struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
}

// FilePickerModel handles file selection
type FilePickerModel struct {
	mode         FilePickerMode
	currentDir   string
	entries      []FileEntry
	cursor       int
	selectedFile string
	done         bool
	cancelled    bool
	err          error

	// For manual path input
	showInput bool
	pathInput textinput.Model
}

// NewFilePickerModel creates a new file picker
func NewFilePickerModel() *FilePickerModel {
	ti := textinput.New()
	ti.Placeholder = "Enter file path..."
	ti.CharLimit = 500
	ti.Width = 60

	// Get executable directory instead of working directory
	execDir := getExecutableDir()

	fp := &FilePickerModel{
		mode:       FilePickerAll,
		currentDir: execDir,
		pathInput:  ti,
	}
	fp.loadFiles()
	return fp
}

// getExecutableDir returns the directory where the executable is located
func getExecutableDir() string {
	// Try to get executable path
	execPath, err := os.Executable()
	if err == nil {
		return filepath.Dir(execPath)
	}
	// Fallback to working directory
	cwd, err := os.Getwd()
	if err == nil && cwd != "" {
		return cwd
	}
	// As a final fallback, use the current directory explicitly
	return "."
}

// SetMode sets what files to display
func (fp *FilePickerModel) SetMode(mode FilePickerMode) {
	fp.mode = mode
	fp.loadFiles()
}

// SetDirectory changes the current directory
func (fp *FilePickerModel) SetDirectory(dir string) {
	fp.currentDir = dir
	fp.loadFiles()
}

// loadFiles reads only matching files from current directory (no subdirectories)
func (fp *FilePickerModel) loadFiles() {
	fp.entries = []FileEntry{}
	fp.cursor = 0

	entries, err := os.ReadDir(fp.currentDir)
	if err != nil {
		fp.err = err
		return
	}

	// Only collect files that match the filter (no directories)
	var files []FileEntry

	for _, entry := range entries {
		// Skip directories - we only want files
		if entry.IsDir() {
			continue
		}

		// Check if file matches filter
		if !fp.matchesFilter(entry.Name()) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		fe := FileEntry{
			Name:  entry.Name(),
			Path:  filepath.Join(fp.currentDir, entry.Name()),
			IsDir: false,
			Size:  info.Size(),
		}
		files = append(files, fe)
	}

	// Sort files alphabetically
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	fp.entries = files
}

// looksLikePath checks if the input string looks like a file path
// This helps detect drag-and-drop or pasted paths
func looksLikePath(s string) bool {
	// Must be reasonably long to be a path
	if len(s) < 3 {
		return false
	}
	// Windows absolute path: C:\ or D:\ etc
	if len(s) >= 3 && s[1] == ':' && (s[2] == '\\' || s[2] == '/') {
		return true
	}
	// UNC path: \\server\share
	if strings.HasPrefix(s, "\\\\") {
		return true
	}
	// Path with backslashes (Windows)
	if strings.Contains(s, "\\") && strings.Contains(s, ".") {
		return true
	}
	return false
}

// matchesFilter checks if file matches current mode
func (fp *FilePickerModel) matchesFilter(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))

	switch fp.mode {
	case FilePickerPDF:
		return ext == ".pdf"
	case FilePickerImage:
		return core.IsImageFile(name)
	case FilePickerAll:
		return core.IsImageFile(name) || core.IsPDFFile(name)
	}
	return true
}

// Update handles input
func (fp *FilePickerModel) Update(msg tea.Msg) (*FilePickerModel, tea.Cmd) {
	var cmd tea.Cmd

	if fp.showInput {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				path := fp.pathInput.Value()
				if path == "" {
					return fp, nil
				}
				// Strip quotes from path (handles "path/to/file" or 'path/to/file')
				path = strings.Trim(path, "\"'`")
				path = strings.TrimSpace(path)

				if info, err := os.Stat(path); err == nil {
					if info.IsDir() {
						// If directory entered, show error - we need a file
						fp.err = fmt.Errorf("please enter a file path, not a directory")
						logging.Warn("User entered directory instead of file", map[string]interface{}{"path": path})
					} else {
						// Validate file matches the filter
						if fp.matchesFilter(filepath.Base(path)) {
							fp.selectedFile = path
							fp.done = true
							logging.Debug("File selected", map[string]interface{}{"path": path})
						} else {
							fp.err = fmt.Errorf("file type not supported for this operation")
							logging.Warn("Unsupported file type selected", map[string]interface{}{
								"path": path,
								"mode": fp.mode,
							})
						}
					}
				} else {
					fp.err = fmt.Errorf("file not found: %s", path)
					logging.Warn("File not found", map[string]interface{}{"path": path, "error": err.Error()})
				}
				return fp, nil
			case "esc":
				fp.showInput = false
				fp.pathInput.SetValue("")
				fp.err = nil
				return fp, nil
			}
		}
		fp.pathInput, cmd = fp.pathInput.Update(msg)
		return fp, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyStr := msg.String()

		// Auto-detect pasted path (drag and drop or paste)
		// If it looks like a path (contains : or \ or /), switch to input mode
		if !fp.showInput && looksLikePath(keyStr) {
			fp.showInput = true
			fp.err = nil
			fp.pathInput.Focus()
			fp.pathInput.SetValue(keyStr)
			// Move cursor to end
			fp.pathInput.CursorEnd()
			logging.Debug("Auto-detected path input", map[string]interface{}{"input": keyStr})
			return fp, textinput.Blink
		}

		switch {
		case key.Matches(msg, keys.Up):
			if fp.cursor > 0 {
				fp.cursor--
			} else if len(fp.entries) > 0 {
				fp.cursor = len(fp.entries) - 1
			}

		case key.Matches(msg, keys.Down):
			if fp.cursor < len(fp.entries)-1 {
				fp.cursor++
			} else {
				fp.cursor = 0
			}

		case key.Matches(msg, keys.Enter):
			if len(fp.entries) > 0 {
				entry := fp.entries[fp.cursor]
				fp.selectedFile = entry.Path
				fp.done = true
			}

		case keyStr == "p": // Manual path input
			fp.showInput = true
			fp.err = nil
			fp.pathInput.Focus()
			return fp, textinput.Blink

		case key.Matches(msg, keys.Back):
			fp.cancelled = true
			fp.done = true
		}
	}

	return fp, nil
}

// View renders the file picker
func (fp *FilePickerModel) View() string {
	var b strings.Builder

	// Header
	header := headerStyle.Render(" " + IconFolder + " Select File ")
	b.WriteString("\n")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Current directory info
	dirLine := inputLabelStyle.Render("Directory: ") + fp.currentDir
	b.WriteString(dirLine)
	b.WriteString("\n\n")

	// Show error if any
	if fp.err != nil {
		b.WriteString(errorStyle.Render("Error: " + fp.err.Error()))
		b.WriteString("\n\n")
	}

	// Manual input mode
	if fp.showInput {
		b.WriteString(inputLabelStyle.Render("Enter file path: "))
		b.WriteString(fp.pathInput.View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter full path to file • Enter to confirm • Esc to cancel"))
		return b.String()
	}

	// File list - show count like batch file
	fileTypeDesc := fp.getFileTypeDescription()
	if len(fp.entries) == 0 {
		b.WriteString(warningStyle.Render(fmt.Sprintf("No %s file(s) found in this directory.", fileTypeDesc)))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Press 'p' to enter a file path manually • Esc to go back"))
	} else {
		// Show file count
		countMsg := fmt.Sprintf("Found %d %s file(s):", len(fp.entries), fileTypeDesc)
		b.WriteString(lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(countMsg))
		b.WriteString("\n\n")

		// Show limited entries with scroll
		visibleCount := 15
		start := 0
		if fp.cursor >= visibleCount {
			start = fp.cursor - visibleCount + 1
		}
		end := start + visibleCount
		if end > len(fp.entries) {
			end = len(fp.entries)
		}

		for i := start; i < end; i++ {
			entry := fp.entries[i]
			cursor := "  "
			style := fileItemStyle

			if i == fp.cursor {
				cursor = IconPointer + " "
				style = selectedFileStyle
			}

			// File number
			numStr := fmt.Sprintf("%2d. ", i+1)

			icon := IconFile
			if strings.HasSuffix(strings.ToLower(entry.Name), ".pdf") {
				icon = IconPDF
			} else {
				icon = IconImage
			}

			// Format file size
			sizeStr := fmt.Sprintf(" (%s)", core.FormatSize(entry.Size))

			line := style.Render(cursor + numStr + icon + " " + entry.Name + sizeStr)
			b.WriteString(line)
			b.WriteString("\n")
		}

		// Scroll indicator
		if len(fp.entries) > visibleCount {
			scrollInfo := lipgloss.NewStyle().Foreground(subtleColor).Render(fmt.Sprintf("\n  Showing %d-%d of %d items", start+1, end, len(fp.entries)))
			b.WriteString(scrollInfo)
		}

		// Help
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("↑↓ Navigate • Enter Select • p Enter path manually • Esc Back"))
	}

	return b.String()
}

// getFileTypeDescription returns description based on mode
func (fp *FilePickerModel) getFileTypeDescription() string {
	switch fp.mode {
	case FilePickerPDF:
		return "PDF"
	case FilePickerImage:
		return "image"
	case FilePickerAll:
		return "image/PDF"
	}
	return "file"
}

// SelectedFile returns the selected file path
func (fp *FilePickerModel) SelectedFile() string {
	return fp.selectedFile
}

// IsDone returns true if selection is complete
func (fp *FilePickerModel) IsDone() bool {
	return fp.done
}

// IsCancelled returns true if user cancelled
func (fp *FilePickerModel) IsCancelled() bool {
	return fp.cancelled
}

// Reset resets the file picker state
func (fp *FilePickerModel) Reset() {
	fp.selectedFile = ""
	fp.done = false
	fp.cancelled = false
	fp.cursor = 0
}
