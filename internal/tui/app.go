package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View represents different screens in the app
type View int

const (
	ViewMenu View = iota
	ViewPDFConverter
	ViewFormatConverter
	ViewCompressor
	ViewFilePicker
)

// MenuItem represents a main menu option
type MenuItem struct {
	Title       string
	Description string
	Icon        string
}

// App is the main application model
type App struct {
	currentView View
	menuItems   []MenuItem
	menuCursor  int
	width       int
	height      int
	quitting    bool

	// Sub-models
	pdfConverter    *PDFConverterModel
	formatConverter *FormatConverterModel
	compressor      *CompressorModel
	filePicker      *FilePickerModel

	// Shared state
	statusMessage string
	isError       bool
}

// KeyMap defines key bindings
type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding
	Help  key.Binding
}

var keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{
		currentView: ViewMenu,
		menuItems: []MenuItem{
			{Title: "PDF to Image Converter", Description: "Convert PDF pages to images (PNG, JPG, etc.)", Icon: IconPDF},
			{Title: "Convert Image Format", Description: "Convert images between formats (WebP, AVIF, etc.)", Icon: IconConvert},
			{Title: "Compress Image/PDF", Description: "Reduce file size by percentage or target size", Icon: IconCompress},
			{Title: "Exit", Description: "Quit the application", Icon: IconExit},
		},
		menuCursor:      0,
		pdfConverter:    NewPDFConverterModel(),
		formatConverter: NewFormatConverterModel(),
		compressor:      NewCompressorModel(),
		filePicker:      NewFilePickerModel(),
	}
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		checkDependencies,
	)
}

// checkDependencies verifies ImageMagick and Ghostscript are available
func checkDependencies() tea.Msg {
	return dependencyCheckMsg{
		magickOK: checkMagick(),
		gsOK:     checkGhostscript(),
	}
}

type dependencyCheckMsg struct {
	magickOK bool
	gsOK     bool
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit
		if key.Matches(msg, keys.Quit) && a.currentView == ViewMenu {
			a.quitting = true
			return a, tea.Quit
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case dependencyCheckMsg:
		if !msg.magickOK {
			a.statusMessage = "⚠️  ImageMagick not found - install from imagemagick.org"
			a.isError = true
		} else if !msg.gsOK {
			a.statusMessage = "⚠️  Ghostscript not found - PDF features limited"
			a.isError = false
		} else {
			a.statusMessage = "✓ ImageMagick & Ghostscript detected"
			a.isError = false
		}
		return a, nil
	}

	// Route to current view
	switch a.currentView {
	case ViewMenu:
		return a.updateMenu(msg)
	case ViewPDFConverter:
		return a.updatePDFConverter(msg)
	case ViewFormatConverter:
		return a.updateFormatConverter(msg)
	case ViewCompressor:
		return a.updateCompressor(msg)
	case ViewFilePicker:
		return a.updateFilePicker(msg)
	}

	return a, nil
}

// updateMenu handles main menu navigation
func (a *App) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Up):
			if a.menuCursor > 0 {
				a.menuCursor--
			} else {
				a.menuCursor = len(a.menuItems) - 1 // Wrap to bottom
			}

		case key.Matches(msg, keys.Down):
			if a.menuCursor < len(a.menuItems)-1 {
				a.menuCursor++
			} else {
				a.menuCursor = 0 // Wrap to top
			}

		case key.Matches(msg, keys.Enter):
			switch a.menuCursor {
			case 0: // PDF to Image
				a.currentView = ViewPDFConverter
				a.pdfConverter = NewPDFConverterModel()
				return a, nil
			case 1: // Convert Format
				a.currentView = ViewFormatConverter
				a.formatConverter = NewFormatConverterModel()
				return a, nil
			case 2: // Compress
				a.currentView = ViewCompressor
				a.compressor = NewCompressorModel()
				return a, nil
			case 3: // Exit
				a.quitting = true
				return a, tea.Quit
			}
		}
	}
	return a, nil
}

// updatePDFConverter handles PDF converter view
func (a *App) updatePDFConverter(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.pdfConverter, cmd = a.pdfConverter.Update(msg)

	if a.pdfConverter.IsDone() {
		if a.pdfConverter.BackToMenu() {
			a.currentView = ViewMenu
		}
	}
	return a, cmd
}

// updateFormatConverter handles format converter view
func (a *App) updateFormatConverter(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.formatConverter, cmd = a.formatConverter.Update(msg)

	if a.formatConverter.IsDone() {
		if a.formatConverter.BackToMenu() {
			a.currentView = ViewMenu
		}
	}
	return a, cmd
}

// updateCompressor handles compressor view
func (a *App) updateCompressor(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.compressor, cmd = a.compressor.Update(msg)

	if a.compressor.IsDone() {
		if a.compressor.BackToMenu() {
			a.currentView = ViewMenu
		}
	}
	return a, cmd
}

// updateFilePicker handles file picker view
func (a *App) updateFilePicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	a.filePicker, cmd = a.filePicker.Update(msg)
	return a, cmd
}

// View implements tea.Model
func (a *App) View() string {
	if a.quitting {
		return "\n  👋 Thanks for using Image Tool!\n\n"
	}

	switch a.currentView {
	case ViewMenu:
		return a.viewMenu()
	case ViewPDFConverter:
		return a.pdfConverter.View()
	case ViewFormatConverter:
		return a.formatConverter.View()
	case ViewCompressor:
		return a.compressor.View()
	case ViewFilePicker:
		return a.filePicker.View()
	}

	return ""
}

// viewMenu renders the main menu
func (a *App) viewMenu() string {
	var b strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(primaryColor).
		Padding(0, 2).
		Render(" 🖼️  Image Tool ")

	subtitle := lipgloss.NewStyle().
		Foreground(subtleColor).
		Render("Powered by Go + ImageMagick")

	b.WriteString("\n")
	b.WriteString(header)
	b.WriteString("\n")
	b.WriteString(subtitle)
	b.WriteString("\n\n")

	// Navigation hint
	hint := helpStyle.Render("Use ↑↓ arrows to navigate, Enter to select, q to quit")
	b.WriteString(hint)
	b.WriteString("\n\n")

	// Menu items
	for i, item := range a.menuItems {
		cursor := "  "
		style := menuItemStyle

		if i == a.menuCursor {
			cursor = IconPointer + " "
			style = selectedItemStyle
		}

		line := style.Render(cursor + item.Icon + "  " + item.Title)
		b.WriteString(line)
		b.WriteString("\n")

		// Show description for selected item
		if i == a.menuCursor {
			desc := descriptionStyle.Render("    " + item.Description)
			b.WriteString(desc)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Status bar
	if a.statusMessage != "" {
		var statusStyle lipgloss.Style
		if a.isError {
			statusStyle = errorStyle
		} else {
			statusStyle = successStyle
		}
		status := statusStyle.Render(a.statusMessage)
		b.WriteString("\n")
		b.WriteString(status)
	}

	return b.String()
}
