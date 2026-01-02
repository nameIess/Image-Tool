package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"imagetool/internal/config"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CompressStep tracks the compression wizard step
type CompressStep int

const (
	CompressStepSelectFile CompressStep = iota
	CompressStepSelectMethod
	CompressStepSetPercent
	CompressStepSetFixedSize
	CompressStepConfirm
	CompressStepCompressing
	CompressStepDone
)

// CompressionMethod defines how to compress
type CompressionMethod int

const (
	CompressMethodPercent CompressionMethod = iota
	CompressMethodFixedSize
)

// CompressorModel handles file compression
type CompressorModel struct {
	step       CompressStep
	filePicker *FilePickerModel

	// Settings
	inputFile     string
	outputFile    string
	inputSize     int64
	method        CompressionMethod
	targetPercent int
	targetBytes   int64

	// Method selection
	methods      []string
	methodCursor int

	// Text inputs
	percentInput textinput.Model
	sizeInput    textinput.Model
	unitInput    textinput.Model

	// Fixed size state
	sizeValue int
	sizeUnit  string // B, KB, MB

	// Results
	result     string
	isError    bool
	outputSize int64

	// Navigation
	done       bool
	backToMenu bool
}

// NewCompressorModel creates a new compressor
func NewCompressorModel() *CompressorModel {
	fp := NewFilePickerModel()
	fp.SetMode(FilePickerAll) // Both images and PDFs

	percentInput := textinput.New()
	percentInput.Placeholder = fmt.Sprintf("%d", config.DefaultCompressPercent)
	percentInput.CharLimit = 3
	percentInput.Width = 10

	sizeInput := textinput.New()
	sizeInput.Placeholder = "100"
	sizeInput.CharLimit = 10
	sizeInput.Width = 15

	unitInput := textinput.New()
	unitInput.Placeholder = "KB"
	unitInput.CharLimit = 2
	unitInput.Width = 5

	return &CompressorModel{
		step:          CompressStepSelectFile,
		filePicker:    fp,
		methods:       []string{"By Percentage", "Fixed File Size"},
		methodCursor:  0,
		targetPercent: config.DefaultCompressPercent,
		percentInput:  percentInput,
		sizeInput:     sizeInput,
		unitInput:     unitInput,
		sizeUnit:      "KB",
	}
}

// Update handles input
func (m *CompressorModel) Update(msg tea.Msg) (*CompressorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch m.step {
	case CompressStepSelectFile:
		m.filePicker, cmd = m.filePicker.Update(msg)
		if m.filePicker.IsDone() {
			if m.filePicker.IsCancelled() {
				m.backToMenu = true
				m.done = true
			} else {
				m.inputFile = m.filePicker.SelectedFile()
				// Get input file size
				if info, err := os.Stat(m.inputFile); err == nil {
					m.inputSize = info.Size()
				}
				m.buildOutputPath()
				m.step = CompressStepSelectMethod
			}
		}
		return m, cmd

	case CompressStepSelectMethod:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.Up):
				if m.methodCursor > 0 {
					m.methodCursor--
				} else {
					m.methodCursor = len(m.methods) - 1
				}
			case key.Matches(msg, keys.Down):
				if m.methodCursor < len(m.methods)-1 {
					m.methodCursor++
				} else {
					m.methodCursor = 0
				}
			case key.Matches(msg, keys.Enter):
				m.method = CompressionMethod(m.methodCursor)
				if m.method == CompressMethodPercent {
					m.step = CompressStepSetPercent
					m.percentInput.Focus()
					return m, textinput.Blink
				} else {
					m.step = CompressStepSetFixedSize
					m.sizeInput.Focus()
					return m, textinput.Blink
				}
			case key.Matches(msg, keys.Back):
				m.step = CompressStepSelectFile
				m.filePicker.Reset()
			}
		}
		return m, nil

	case CompressStepSetPercent:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				val := m.percentInput.Value()
				if val == "" {
					m.targetPercent = config.DefaultCompressPercent
				} else {
					parsedPercent, err := strconv.Atoi(val)
					if err != nil {
						m.targetPercent = config.DefaultCompressPercent
					} else {
						m.targetPercent = parsedPercent
					}
					if m.targetPercent < 1 {
						m.targetPercent = 1
					}
					if m.targetPercent > 100 {
						m.targetPercent = 100
					}
				}
				// Calculate target bytes
				m.targetBytes = m.inputSize * int64(m.targetPercent) / 100
				m.step = CompressStepConfirm
				m.percentInput.Blur()
				return m, nil
			case "esc":
				m.step = CompressStepSelectMethod
				m.percentInput.Blur()
				return m, nil
			}
		}
		m.percentInput, cmd = m.percentInput.Update(msg)
		return m, cmd

	case CompressStepSetFixedSize:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Parse size value
				val := m.sizeInput.Value()
				if val == "" {
					m.sizeValue = 100
				} else {
					parsedSize, err := strconv.Atoi(val)
					if err != nil {
						m.sizeValue = 100
					} else {
						m.sizeValue = parsedSize
					}
				}
				if m.sizeValue < 1 {
					m.sizeValue = 1
				}

				// Ask for unit
				m.sizeInput.Blur()
				m.unitInput.Focus()
				return m, textinput.Blink

			case "esc":
				m.step = CompressStepSelectMethod
				m.sizeInput.Blur()
				return m, nil

			case "k", "K":
				if m.sizeInput.Value() != "" {
					parsedSize, err := strconv.Atoi(m.sizeInput.Value())
					if err == nil {
						m.sizeValue = parsedSize
						m.sizeUnit = "KB"
						m.targetBytes = int64(m.sizeValue) * 1024
						m.step = CompressStepConfirm
						m.sizeInput.Blur()
					}
					return m, nil
				}
			case "m", "M":
				if m.sizeInput.Value() != "" {
					parsedSize, err := strconv.Atoi(m.sizeInput.Value())
					if err == nil {
						m.sizeValue = parsedSize
						m.sizeUnit = "MB"
						m.targetBytes = int64(m.sizeValue) * 1024 * 1024
						m.step = CompressStepConfirm
						m.sizeInput.Blur()
					}
					return m, nil
				}
			case "b", "B":
				if m.sizeInput.Value() != "" {
					parsedSize, err := strconv.Atoi(m.sizeInput.Value())
					if err == nil {
						m.sizeValue = parsedSize
						m.sizeUnit = "B"
						m.targetBytes = int64(m.sizeValue)
						m.step = CompressStepConfirm
						m.sizeInput.Blur()
					}
					return m, nil
				}
			}
		}

		// Check if unit input is focused
		if m.unitInput.Focused() {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "enter":
					unit := strings.ToUpper(m.unitInput.Value())
					if unit == "" {
						unit = "KB"
					}
					m.sizeUnit = unit

					switch m.sizeUnit {
					case "MB":
						m.targetBytes = int64(m.sizeValue) * 1024 * 1024
					case "KB":
						m.targetBytes = int64(m.sizeValue) * 1024
					default:
						m.targetBytes = int64(m.sizeValue)
					}
					m.step = CompressStepConfirm
					m.unitInput.Blur()
					return m, nil
				case "esc":
					m.unitInput.Blur()
					m.sizeInput.Focus()
					return m, textinput.Blink
				}
			}
			m.unitInput, cmd = m.unitInput.Update(msg)
			return m, cmd
		}

		m.sizeInput, cmd = m.sizeInput.Update(msg)
		return m, cmd

	case CompressStepConfirm:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y", "Y", "enter":
				m.step = CompressStepCompressing
				return m, m.runCompression
			case "n", "N", "esc":
				if m.method == CompressMethodPercent {
					m.step = CompressStepSetPercent
					m.percentInput.Focus()
					return m, textinput.Blink
				} else {
					m.step = CompressStepSetFixedSize
					m.sizeInput.Focus()
					return m, textinput.Blink
				}
			case "b":
				m.backToMenu = true
				m.done = true
			}
		}
		return m, nil

	case CompressStepCompressing:
		switch msg := msg.(type) {
		case compressResultMsg:
			m.step = CompressStepDone
			m.result = msg.message
			m.isError = msg.isError
			m.outputSize = msg.outputSize
			return m, nil
		}
		return m, nil

	case CompressStepDone:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "m":
				m.backToMenu = true
				m.done = true
			case "a": // Compress another
				m.step = CompressStepSelectFile
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
func (m *CompressorModel) buildOutputPath() {
	dir := filepath.Dir(m.inputFile)
	base := strings.TrimSuffix(filepath.Base(m.inputFile), filepath.Ext(m.inputFile))
	ext := filepath.Ext(m.inputFile)

	// For compression, we output as JPG for images (better compression)
	// Keep PDF as PDF
	if strings.ToLower(ext) != ".pdf" {
		ext = ".jpg"
	}

	m.outputFile = filepath.Join(dir, base+"_comp"+ext)
}

// compressResultMsg contains compression results
type compressResultMsg struct {
	message    string
	isError    bool
	outputSize int64
}

// runCompression executes the ImageMagick command
func (m *CompressorModel) runCompression() tea.Msg {
	// Build ImageMagick arguments based on output format
	ext := strings.ToLower(filepath.Ext(m.outputFile))
	args := []string{m.inputFile}

	// Use jpeg:extent for target size compression when output is JPEG
	if ext == ".jpg" || ext == ".jpeg" {
		targetSize := fmt.Sprintf("%d", m.targetBytes)
		args = append(args, "-define", fmt.Sprintf("jpeg:extent=%s", targetSize))
	}
	args = append(args, m.outputFile)

	cmd := exec.Command("magick", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return compressResultMsg{
			message: fmt.Sprintf("Compression failed: %v\n%s", err, string(output)),
			isError: true,
		}
	}

	// Get output file size
	var outputSize int64
	if info, err := os.Stat(m.outputFile); err == nil {
		outputSize = info.Size()
	}

	// Check if we achieved target
	reduction := 0.0
	if m.inputSize > 0 {
		reduction = float64(m.inputSize-outputSize) / float64(m.inputSize) * 100
	}

	msg := fmt.Sprintf("Compressed successfully! Reduced by %.1f%%", reduction)
	if outputSize > m.targetBytes {
		msg = fmt.Sprintf("Compressed, but couldn't reach target. Best possible: %s", formatSize(outputSize))
	}

	return compressResultMsg{
		message:    msg,
		isError:    false,
		outputSize: outputSize,
	}
}

// View renders the compressor
func (m *CompressorModel) View() string {
	var b strings.Builder

	// Header
	header := headerStyle.Render(" " + IconCompress + " Compress Image/PDF ")
	b.WriteString("\n")
	b.WriteString(header)
	b.WriteString("\n\n")

	switch m.step {
	case CompressStepSelectFile:
		b.WriteString(m.filePicker.View())

	case CompressStepSelectMethod:
		b.WriteString(inputLabelStyle.Render("Select compression method:"))
		b.WriteString("\n\n")

		b.WriteString(descriptionStyle.Render(fmt.Sprintf("Input file: %s (%s)", filepath.Base(m.inputFile), formatSize(m.inputSize))))
		b.WriteString("\n\n")

		for i, method := range m.methods {
			cursor := "  "
			style := menuItemStyle
			if i == m.methodCursor {
				cursor = IconPointer + " "
				style = selectedItemStyle
			}
			b.WriteString(style.Render(cursor + method))
			b.WriteString("\n")

			// Description for selected
			if i == m.methodCursor {
				desc := ""
				if i == 0 {
					desc = "    Compress to a percentage of original size (e.g., 50%)"
				} else {
					desc = "    Compress to exact target size (e.g., 100KB)"
				}
				b.WriteString(descriptionStyle.Render(desc))
				b.WriteString("\n")
			}
		}
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("↑↓ Navigate • Enter Select • Esc Back"))

	case CompressStepSetPercent:
		b.WriteString(inputLabelStyle.Render("Target percentage (1-100):"))
		b.WriteString("\n\n")
		b.WriteString(m.percentInput.View())
		b.WriteString("\n\n")

		// Show preview
		preview := m.inputSize * int64(config.DefaultCompressPercent) / 100
		if m.percentInput.Value() != "" {
			pct, _ := strconv.Atoi(m.percentInput.Value())
			if pct > 0 && pct <= 100 {
				preview = m.inputSize * int64(pct) / 100
			}
		}
		b.WriteString(descriptionStyle.Render(fmt.Sprintf("Current: %s → Target: ~%s", formatSize(m.inputSize), formatSize(preview))))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter to confirm • Esc Back"))

	case CompressStepSetFixedSize:
		b.WriteString(inputLabelStyle.Render("Target file size:"))
		b.WriteString("\n\n")

		if m.unitInput.Focused() {
			b.WriteString(fmt.Sprintf("Size: %d ", m.sizeValue))
			b.WriteString(m.unitInput.View())
			b.WriteString("\n\n")
			b.WriteString(descriptionStyle.Render("Enter unit: B, KB, or MB"))
		} else {
			b.WriteString("Size: ")
			b.WriteString(m.sizeInput.View())
			b.WriteString("\n\n")
			b.WriteString(descriptionStyle.Render(fmt.Sprintf("Current: %s | Enter number then press K for KB, M for MB, or Enter for unit selection", formatSize(m.inputSize))))
		}
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter to continue • K=KB M=MB • Esc Back"))

	case CompressStepConfirm:
		b.WriteString(inputLabelStyle.Render("Compression Summary"))
		b.WriteString("\n\n")

		methodStr := "Percentage"
		targetStr := fmt.Sprintf("%d%% of original", m.targetPercent)
		if m.method == CompressMethodFixedSize {
			methodStr = "Fixed Size"
			targetStr = fmt.Sprintf("%d %s", m.sizeValue, m.sizeUnit)
		}

		summaryBox := boxStyle.Render(
			fmt.Sprintf("Input:   %s (%s)\n", filepath.Base(m.inputFile), formatSize(m.inputSize)) +
				fmt.Sprintf("Method:  %s\n", methodStr) +
				fmt.Sprintf("Target:  %s (%s)\n", targetStr, formatSize(m.targetBytes)) +
				fmt.Sprintf("Output:  %s", filepath.Base(m.outputFile)),
		)
		b.WriteString(summaryBox)
		b.WriteString("\n\n")

		if strings.ToLower(filepath.Ext(m.inputFile)) == ".pdf" {
			b.WriteString(warningStyle.Render("⚠️  PDF compression may rasterize content"))
			b.WriteString("\n\n")
		}

		b.WriteString(warningStyle.Render("Proceed with compression? (Y/n)"))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Y/Enter Proceed • N/Esc Back • B Menu"))

	case CompressStepCompressing:
		b.WriteString("\n")
		b.WriteString(progressStyle.Render("⏳ Compressing... Please wait"))
		b.WriteString("\n")

	case CompressStepDone:
		if m.isError {
			b.WriteString(errorStyle.Render(IconError + " " + m.result))
		} else {
			b.WriteString(successStyle.Render(IconSuccess + " " + m.result))
			b.WriteString("\n\n")
			b.WriteString(descriptionStyle.Render(fmt.Sprintf("Original: %s → Compressed: %s", formatSize(m.inputSize), formatSize(m.outputSize))))
			b.WriteString("\n")
			b.WriteString(descriptionStyle.Render(fmt.Sprintf("Output: %s", m.outputFile)))
		}
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("Enter/M Menu • A Compress Another • Q Quit"))
	}

	return b.String()
}

// IsDone returns true if compression flow is complete
func (m *CompressorModel) IsDone() bool {
	return m.done
}

// BackToMenu returns true if user wants to go back to menu
func (m *CompressorModel) BackToMenu() bool {
	return m.backToMenu
}
