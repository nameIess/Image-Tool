package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"imagetool/internal/config"
)

// FormatStep tracks the conversion wizard step
type FormatStep int

const (
	FormatStepSelectFile FormatStep = iota
	FormatStepSelectFormat
	FormatStepConfirm
	FormatStepConverting
	FormatStepDone
)

// FormatConverterModel handles image format conversion
type FormatConverterModel struct {
	step         FormatStep
	filePicker   *FilePickerModel

	// Settings
	inputFile    string
	outputFormat string
	outputFile   string

	// Format selection
	formats       []string
	formatCursor  int
	customFormat  bool
	customInput   textinput.Model

	// Results
	result    string
	isError   bool
	fileSize  int64

	// Navigation
	done       bool
	backToMenu bool
}

// NewFormatConverterModel creates a new format converter
func NewFormatConverterModel() *FormatConverterModel {
	fp := NewFilePickerModel()
	fp.SetMode(FilePickerImage)

	customInput := textinput.New()
	customInput.Placeholder = "avif, webp, heic..."
	customInput.CharLimit = 10
	customInput.Width = 20

	// Include "custom" as last option
	formats := append([]string{}, config.SupportedImageFormats...)
	formats = append(formats, "custom")

	return &FormatConverterModel{
		step:         FormatStepSelectFile,
		filePicker:   fp,
		formats:      formats,
		formatCursor: 0,
		customInput:  customInput,
	}
}

// Update handles input
func (m *FormatConverterModel) Update(msg tea.Msg) (*FormatConverterModel, tea.Cmd) {
	var cmd tea.Cmd

	switch m.step {
	case FormatStepSelectFile:
		m.filePicker, cmd = m.filePicker.Update(msg)
		if m.filePicker.IsDone() {
			if m.filePicker.IsCancelled() {
				m.backToMenu = true
				m.done = true
			} else {
				m.inputFile = m.filePicker.SelectedFile()
				m.step = FormatStepSelectFormat
			}
		}
		return m, cmd

	case FormatStepSelectFormat:
		// Handle custom format input mode
		if m.customFormat {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "enter":
					val := strings.TrimSpace(m.customInput.Value())
					if val != "" {
						m.outputFormat = strings.TrimPrefix(val, ".")
						m.customFormat = false
						m.customInput.Blur()
						m.buildOutputPath()
						m.step = FormatStepConfirm
					}
					return m, nil
				case "esc":
					m.customFormat = false
					m.customInput.Blur()
					m.customInput.SetValue("")
					return m, nil
				}
			}
			m.customInput, cmd = m.customInput.Update(msg)
			return m, cmd
		}

		// Normal format selection
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Up):
				if m.formatCursor > 0 {
					m.formatCursor--
				} else {
					m.formatCursor = len(m.formats) - 1
				}
			case key.Matches(msg, keys.Down):
				if m.formatCursor < len(m.formats)-1 {
					m.formatCursor++
				} else {
					m.formatCursor = 0
				}
			case key.Matches(msg, keys.Enter):
				selected := m.formats[m.formatCursor]
				if selected == "custom" {
					m.customFormat = true
					m.customInput.Focus()
					return m, textinput.Blink
				}
				m.outputFormat = selected
				m.buildOutputPath()
				m.step = FormatStepConfirm
				return m, nil
			case key.Matches(msg, keys.Back):
				m.step = FormatStepSelectFile
				m.filePicker.Reset()
			}
		}
		return m, nil

	case FormatStepConfirm:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y", "enter":
				m.step = FormatStepConverting
				return m, m.runConversion
			case "n", "N", "esc":
				m.step = FormatStepSelectFormat
			case "b":
				m.backToMenu = true
				m.done = true
			}
		}
		return m, nil

	case FormatStepConverting:
		switch msg := msg.(type) {
		case formatConversionResultMsg:
			m.step = FormatStepDone
			m.result = msg.message
			m.isError = msg.isError
			m.fileSize = msg.fileSize
			return m, nil
		}
		return m, nil

	case FormatStepDone:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "m":
				m.backToMenu = true
				m.done = true
			case "a": // Convert another
				m.step = FormatStepSelectFile
				m.filePicker.Reset()
				m.inputFile = ""
				m.outputFile = ""
				m.result = ""
			case "q":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	return m, nil
}

// buildOutputPath creates the output file path
func (m *FormatConverterModel) buildOutputPath() {
	dir := filepath.Dir(m.inputFile)
	base := strings.TrimSuffix(filepath.Base(m.inputFile), filepath.Ext(m.inputFile))
	m.outputFile = filepath.Join(dir, base+"_conv."+m.outputFormat)
}

// formatConversionResultMsg contains conversion results
type formatConversionResultMsg struct {
	message  string
	isError  bool
	fileSize int64
}

// runConversion executes the ImageMagick command
func (m *FormatConverterModel) runConversion() tea.Msg {
	// Run ImageMagick
	cmd := exec.Command("magick", m.inputFile, m.outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return formatConversionResultMsg{
			message: fmt.Sprintf("Conversion failed: %v\n%s", err, string(output)),
			isError: true,
		}
	}

	// Get output file size
	info, err := os.Stat(m.outputFile)
	var size int64
	if err == nil {
		size = info.Size()
	}

	return formatConversionResultMsg{
		message:  "Image converted successfully",
		isError:  false,
		fileSize: size,
	}
}

// View renders the format converter
func (m *FormatConverterModel) View() string {
	var b strings.Builder

	// Header
	header := headerStyle.Render(" " + IconConvert + " Convert Image Format ")
	b.WriteString("\n")
	b.WriteString(header)
	b.WriteString("\n\n")

	switch m.step {
	case FormatStepSelectFile:
		b.WriteString(m.filePicker.View())

	case FormatStepSelectFormat:
		if m.customFormat {
			b.WriteString(inputLabelStyle.Render("Enter custom format:"))
			b.WriteString("\n\n")
			b.WriteString(m.customInput.View())
			b.WriteString("\n\n")
			b.WriteString(descriptionStyle.Render("Examples: avif, webp, heic, ico, svg"))
			b.WriteString("\n\n")
			b.WriteString(helpStyle.Render("Enter to confirm • Esc Back"))
		} else {
			b.WriteString(inputLabelStyle.Render("Select output format:"))
			b.WriteString("\n\n")

			for i, format := range m.formats {
				cursor := "  "
				style := menuItemStyle
				if i == m.formatCursor {
					cursor = IconPointer + " "
					style = selectedItemStyle
				}
				
				display := strings.ToUpper(format)
				if format == "custom" {
					display = "Custom (enter any format)"
				}
				b.WriteString(style.Render(cursor + display))
				b.WriteString("\n")
			}
			b.WriteString("\n")
			b.WriteString(helpStyle.Render("↑↓ Navigate • Enter Select • Esc Back"))
		}

	case FormatStepConfirm:
		b.WriteString(inputLabelStyle.Render("Conversion Summary"))
		b.WriteString("\n\n")

		// Get input file size
		inputSize := int64(0)
		if info, err := os.Stat(m.inputFile); err == nil {
			inputSize = info.Size()
		}

		summaryBox := boxStyle.Render(
			fmt.Sprintf("Input:   %s (%s)\n", filepath.Base(m.inputFile), formatSize(inputSize)) +
				fmt.Sprintf("Format:  %s → %s\n", strings.ToUpper(filepath.Ext(m.inputFile)[1:]), strings.ToUpper(m.outputFormat)) +
				fmt.Sprintf("Output:  %s", filepath.Base(m.outputFile)),
		)
		b.WriteString(summaryBox)
		b.WriteString("\n\n")
		b.WriteString(warningStyle.Render("Proceed with conversion? (Y/n)"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Y/Enter Proceed • N/Esc Back • B Menu"))

	case FormatStepConverting:
		b.WriteString("\n")
		b.WriteString(progressStyle.Render("⏳ Converting... Please wait"))
		b.WriteString("\n")

	case FormatStepDone:
		if m.isError {
			b.WriteString(errorStyle.Render(IconError + " " + m.result))
		} else {
			b.WriteString(successStyle.Render(IconSuccess + " " + m.result))
			b.WriteString("\n\n")
			b.WriteString(descriptionStyle.Render(fmt.Sprintf("Output: %s (%s)", m.outputFile, formatSize(m.fileSize))))
		}
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter/M Menu • A Convert Another • Q Quit"))
	}

	return b.String()
}

// IsDone returns true if conversion flow is complete
func (m *FormatConverterModel) IsDone() bool {
	return m.done
}

// BackToMenu returns true if user wants to go back to menu
func (m *FormatConverterModel) BackToMenu() bool {
	return m.backToMenu
}
