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

// PDFStep tracks the conversion wizard step
type PDFStep int

const (
	PDFStepSelectFile PDFStep = iota
	PDFStepSelectFormat
	PDFStepSetDensity
	PDFStepSetQuality
	PDFStepSetPrefix
	PDFStepConfirm
	PDFStepConverting
	PDFStepDone
)

// PDFConverterModel handles PDF to image conversion
type PDFConverterModel struct {
	step         PDFStep
	filePicker   *FilePickerModel
	
	// Settings
	inputFile    string
	outputFormat string
	density      int
	quality      int
	prefix       string
	outputDir    string
	
	// Format selection
	formats      []string
	formatCursor int
	
	// Text inputs
	densityInput textinput.Model
	qualityInput textinput.Model
	prefixInput  textinput.Model
	
	// Results
	converting   bool
	progress     string
	result       string
	isError      bool
	outputFiles  []string
	
	// Navigation
	done         bool
	backToMenu   bool
}

// NewPDFConverterModel creates a new PDF converter
func NewPDFConverterModel() *PDFConverterModel {
	fp := NewFilePickerModel()
	fp.SetMode(FilePickerPDF)

	densityInput := textinput.New()
	densityInput.Placeholder = fmt.Sprintf("%d", config.DefaultDensity)
	densityInput.CharLimit = 4
	densityInput.Width = 10

	qualityInput := textinput.New()
	qualityInput.Placeholder = fmt.Sprintf("%d", config.DefaultQuality)
	qualityInput.CharLimit = 3
	qualityInput.Width = 10

	prefixInput := textinput.New()
	prefixInput.Placeholder = config.DefaultPrefix
	prefixInput.CharLimit = 50
	prefixInput.Width = 30

	return &PDFConverterModel{
		step:         PDFStepSelectFile,
		filePicker:   fp,
		formats:      config.SupportedPDFOutputFormats,
		formatCursor: 0,
		outputFormat: config.DefaultOutputFormat,
		density:      config.DefaultDensity,
		quality:      config.DefaultQuality,
		prefix:       config.DefaultPrefix,
		densityInput: densityInput,
		qualityInput: qualityInput,
		prefixInput:  prefixInput,
	}
}

// Update handles input
func (m *PDFConverterModel) Update(msg tea.Msg) (*PDFConverterModel, tea.Cmd) {
	var cmd tea.Cmd

	switch m.step {
	case PDFStepSelectFile:
		m.filePicker, cmd = m.filePicker.Update(msg)
		if m.filePicker.IsDone() {
			if m.filePicker.IsCancelled() {
				m.backToMenu = true
				m.done = true
			} else {
				m.inputFile = m.filePicker.SelectedFile()
				m.step = PDFStepSelectFormat
				// Set default output directory
				m.outputDir = filepath.Join(filepath.Dir(m.inputFile), 
					strings.TrimSuffix(filepath.Base(m.inputFile), filepath.Ext(m.inputFile))+"_images")
			}
		}
		return m, cmd

	case PDFStepSelectFormat:
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
				m.outputFormat = m.formats[m.formatCursor]
				m.step = PDFStepSetDensity
				m.densityInput.Focus()
				return m, textinput.Blink
			case key.Matches(msg, keys.Back):
				m.step = PDFStepSelectFile
				m.filePicker.Reset()
			}
		}
		return m, nil

	case PDFStepSetDensity:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				val := m.densityInput.Value()
				if val == "" {
					m.density = config.DefaultDensity
				} else {
					fmt.Sscanf(val, "%d", &m.density)
					if m.density < 72 {
						m.density = 72
					}
					if m.density > 600 {
						m.density = 600
					}
				}
				m.step = PDFStepSetQuality
				m.densityInput.Blur()
				m.qualityInput.Focus()
				return m, textinput.Blink
			case "esc":
				m.step = PDFStepSelectFormat
				m.densityInput.Blur()
				return m, nil
			}
		}
		m.densityInput, cmd = m.densityInput.Update(msg)
		return m, cmd

	case PDFStepSetQuality:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				val := m.qualityInput.Value()
				if val == "" {
					m.quality = config.DefaultQuality
				} else {
					fmt.Sscanf(val, "%d", &m.quality)
					if m.quality < 1 {
						m.quality = 1
					}
					if m.quality > 100 {
						m.quality = 100
					}
				}
				m.step = PDFStepSetPrefix
				m.qualityInput.Blur()
				m.prefixInput.Focus()
				return m, textinput.Blink
			case "esc":
				m.step = PDFStepSetDensity
				m.qualityInput.Blur()
				m.densityInput.Focus()
				return m, textinput.Blink
			}
		}
		m.qualityInput, cmd = m.qualityInput.Update(msg)
		return m, cmd

	case PDFStepSetPrefix:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				val := m.prefixInput.Value()
				if val == "" {
					m.prefix = config.DefaultPrefix
				} else {
					m.prefix = val
				}
				m.step = PDFStepConfirm
				m.prefixInput.Blur()
				return m, nil
			case "esc":
				m.step = PDFStepSetQuality
				m.prefixInput.Blur()
				m.qualityInput.Focus()
				return m, textinput.Blink
			}
		}
		m.prefixInput, cmd = m.prefixInput.Update(msg)
		return m, cmd

	case PDFStepConfirm:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y", "enter":
				m.step = PDFStepConverting
				return m, m.runConversion
			case "n", "N", "esc":
				m.step = PDFStepSetPrefix
				m.prefixInput.Focus()
				return m, textinput.Blink
			case "b":
				m.backToMenu = true
				m.done = true
			}
		}
		return m, nil

	case PDFStepConverting:
		switch msg := msg.(type) {
		case conversionResultMsg:
			m.step = PDFStepDone
			m.result = msg.message
			m.isError = msg.isError
			m.outputFiles = msg.files
			return m, nil
		}
		return m, nil

	case PDFStepDone:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "m":
				m.backToMenu = true
				m.done = true
			case "o": // Open output folder
				if m.outputDir != "" {
					exec.Command("explorer", m.outputDir).Start()
				}
			case "q":
				return m, tea.Quit
			}
		}
		return m, nil
	}

	return m, nil
}

// conversionResultMsg contains conversion results
type conversionResultMsg struct {
	message string
	isError bool
	files   []string
}

// runConversion executes the ImageMagick command
func (m *PDFConverterModel) runConversion() tea.Msg {
	// Create output directory
	if err := os.MkdirAll(m.outputDir, 0755); err != nil {
		return conversionResultMsg{
			message: fmt.Sprintf("Failed to create output directory: %v", err),
			isError: true,
		}
	}

	// Build output pattern
	outputPattern := filepath.Join(m.outputDir, m.prefix+"%d."+m.outputFormat)

	// Run ImageMagick
	cmd := exec.Command("magick",
		"-density", fmt.Sprintf("%d", m.density),
		m.inputFile,
		"-quality", fmt.Sprintf("%d", m.quality),
		outputPattern,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return conversionResultMsg{
			message: fmt.Sprintf("Conversion failed: %v\n%s", err, string(output)),
			isError: true,
		}
	}

	// Count output files
	var files []string
	pattern := filepath.Join(m.outputDir, m.prefix+"*."+m.outputFormat)
	matches, _ := filepath.Glob(pattern)
	files = append(files, matches...)

	return conversionResultMsg{
		message: fmt.Sprintf("Successfully converted %d page(s)", len(files)),
		isError: false,
		files:   files,
	}
}

// View renders the PDF converter
func (m *PDFConverterModel) View() string {
	var b strings.Builder

	// Header
	header := headerStyle.Render(" " + IconPDF + " PDF to Image Converter ")
	b.WriteString("\n")
	b.WriteString(header)
	b.WriteString("\n\n")

	switch m.step {
	case PDFStepSelectFile:
		b.WriteString(m.filePicker.View())

	case PDFStepSelectFormat:
		b.WriteString(inputLabelStyle.Render("Select output format:"))
		b.WriteString("\n\n")
		
		for i, format := range m.formats {
			cursor := "  "
			style := menuItemStyle
			if i == m.formatCursor {
				cursor = IconPointer + " "
				style = selectedItemStyle
			}
			b.WriteString(style.Render(cursor + strings.ToUpper(format)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("↑↓ Navigate • Enter Select • Esc Back"))

	case PDFStepSetDensity:
		b.WriteString(inputLabelStyle.Render("Set DPI/Density (72-600):"))
		b.WriteString("\n\n")
		b.WriteString(m.densityInput.View())
		b.WriteString("\n\n")
		b.WriteString(descriptionStyle.Render(fmt.Sprintf("Higher = better quality, larger files. Default: %d", config.DefaultDensity)))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter to confirm • Esc Back"))

	case PDFStepSetQuality:
		b.WriteString(inputLabelStyle.Render("Set Quality (1-100):"))
		b.WriteString("\n\n")
		b.WriteString(m.qualityInput.View())
		b.WriteString("\n\n")
		b.WriteString(descriptionStyle.Render(fmt.Sprintf("100 = best quality. Default: %d", config.DefaultQuality)))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter to confirm • Esc Back"))

	case PDFStepSetPrefix:
		b.WriteString(inputLabelStyle.Render("Set filename prefix:"))
		b.WriteString("\n\n")
		b.WriteString(m.prefixInput.View())
		b.WriteString("\n\n")
		b.WriteString(descriptionStyle.Render(fmt.Sprintf("Output: %s0.%s, %s1.%s, ... Default: %s", 
			m.prefix, m.outputFormat, m.prefix, m.outputFormat, config.DefaultPrefix)))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter to confirm • Esc Back"))

	case PDFStepConfirm:
		b.WriteString(inputLabelStyle.Render("Conversion Summary"))
		b.WriteString("\n\n")
		
		summaryBox := boxStyle.Render(
			fmt.Sprintf("Input:    %s\n", filepath.Base(m.inputFile)) +
			fmt.Sprintf("Format:   %s\n", strings.ToUpper(m.outputFormat)) +
			fmt.Sprintf("Density:  %d DPI\n", m.density) +
			fmt.Sprintf("Quality:  %d\n", m.quality) +
			fmt.Sprintf("Prefix:   %s\n", m.prefix) +
			fmt.Sprintf("Output:   %s", m.outputDir),
		)
		b.WriteString(summaryBox)
		b.WriteString("\n\n")
		b.WriteString(warningStyle.Render("Proceed with conversion? (Y/n)"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Y/Enter Proceed • N/Esc Back • B Menu"))

	case PDFStepConverting:
		b.WriteString("\n")
		b.WriteString(progressStyle.Render("⏳ Converting... Please wait"))
		b.WriteString("\n")

	case PDFStepDone:
		if m.isError {
			b.WriteString(errorStyle.Render(IconError + " " + m.result))
		} else {
			b.WriteString(successStyle.Render(IconSuccess + " " + m.result))
			b.WriteString("\n\n")
			b.WriteString(descriptionStyle.Render("Output folder: " + m.outputDir))
		}
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter/M Menu • O Open Folder • Q Quit"))
	}

	return b.String()
}

// IsDone returns true if conversion flow is complete
func (m *PDFConverterModel) IsDone() bool {
	return m.done
}

// BackToMenu returns true if user wants to go back to menu
func (m *PDFConverterModel) BackToMenu() bool {
	return m.backToMenu
}
