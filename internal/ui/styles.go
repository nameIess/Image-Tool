// Package ui provides the TUI (Terminal User Interface) for Image-Tool.
// This layer handles rendering and input only - no system command execution.
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	secondaryColor = lipgloss.Color("#10B981") // Green
	accentColor    = lipgloss.Color("#F59E0B") // Amber
	errorColor     = lipgloss.Color("#EF4444") // Red
	subtleColor    = lipgloss.Color("#6B7280") // Gray
	textColor      = lipgloss.Color("#F3F4F6") // Light gray
)

// Styles
var (
	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Header box style
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor).
			Background(primaryColor).
			Padding(0, 2).
			MarginBottom(1)

	// Menu item styles
	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(primaryColor).
				Bold(true)

	// Description style
	descriptionStyle = lipgloss.NewStyle().
				Foreground(subtleColor).
				Italic(true).
				MarginTop(1)

	// Success style
	successStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	// Error style
	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	// Warning style
	warningStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	// Input style
	inputLabelStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Box style for sections
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			MarginTop(1)

	// Progress style
	progressStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	// File item style
	fileItemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedFileStyle = lipgloss.NewStyle().
				PaddingLeft(4).
				Foreground(secondaryColor).
				Bold(true)

	// Directory style
	dirStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	// Dependency status styles
	depOKStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	depErrorStyle = lipgloss.NewStyle().
			Foreground(errorColor)
)

// Icons
const (
	IconFolder   = "📁"
	IconFile     = "📄"
	IconImage    = "🖼️"
	IconPDF      = "📑"
	IconCheck    = "✓"
	IconCross    = "✗"
	IconArrow    = "→"
	IconPointer  = "▶"
	IconCompress = "📦"
	IconConvert  = "🔄"
	IconSettings = "⚙️"
	IconExit     = "❌"
	IconSuccess  = "✅"
	IconError    = "❌"
	IconWarning  = "⚠️"
	IconSpinner  = "◐"
)
