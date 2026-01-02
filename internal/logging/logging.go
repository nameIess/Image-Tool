// Package logging provides error and activity logging for Image-Tool.
// Logs are written to a file for debugging and troubleshooting.
package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Level represents the severity of a log message.
type Level int

const (
	// LevelDebug is for detailed debugging information.
	LevelDebug Level = iota
	// LevelInfo is for general operational information.
	LevelInfo
	// LevelWarn is for warning conditions.
	LevelWarn
	// LevelError is for error conditions.
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}
	return "UNKNOWN"
}

// Logger handles logging to file.
type Logger struct {
	mu       sync.Mutex
	file     *os.File
	writer   io.Writer
	minLevel Level
	errors   []LogEntry // Collected errors for summary
}

// LogEntry represents a single log entry.
type LogEntry struct {
	Time    time.Time
	Level   Level
	Message string
	Context map[string]interface{}
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Init initializes the default logger. Safe to call multiple times.
func Init(logDir string) error {
	var initErr error
	once.Do(func() {
		logger, err := NewLogger(logDir)
		if err != nil {
			initErr = err
			return
		}
		defaultLogger = logger
	})
	return initErr
}

// NewLogger creates a new logger instance.
func NewLogger(logDir string) (*Logger, error) {
	// Create log directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, fmt.Sprintf("imagetool_%s.log", timestamp))

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := &Logger{
		file:     file,
		writer:   file,
		minLevel: LevelInfo,
		errors:   make([]LogEntry, 0),
	}

	// Log startup
	logger.Info("Image-Tool started", map[string]interface{}{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	})

	return logger, nil
}

// SetLevel sets the minimum log level.
func (l *Logger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.minLevel = level
}

// log writes a log entry.
func (l *Logger) log(level Level, message string, context map[string]interface{}) {
	if l == nil || l.writer == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if level < l.minLevel {
		return
	}

	entry := LogEntry{
		Time:    time.Now(),
		Level:   level,
		Message: message,
		Context: context,
	}

	// Collect errors for summary
	if level == LevelError {
		l.errors = append(l.errors, entry)
	}

	// Format log line
	line := fmt.Sprintf("[%s] %s: %s",
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Level.String(),
		entry.Message,
	)

	// Add context if present
	if len(context) > 0 {
		line += fmt.Sprintf(" %v", context)
	}

	line += "\n"

	l.writer.Write([]byte(line))
}

// Debug logs a debug message.
func (l *Logger) Debug(message string, context map[string]interface{}) {
	l.log(LevelDebug, message, context)
}

// Info logs an informational message.
func (l *Logger) Info(message string, context map[string]interface{}) {
	l.log(LevelInfo, message, context)
}

// Warn logs a warning message.
func (l *Logger) Warn(message string, context map[string]interface{}) {
	l.log(LevelWarn, message, context)
}

// Error logs an error message.
func (l *Logger) Error(message string, context map[string]interface{}) {
	l.log(LevelError, message, context)
}

// GetErrors returns all logged errors.
func (l *Logger) GetErrors() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]LogEntry{}, l.errors...)
}

// ClearErrors clears the error collection.
func (l *Logger) ClearErrors() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errors = make([]LogEntry, 0)
}

// GetErrorSummary returns a formatted summary of errors.
func (l *Logger) GetErrorSummary() string {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.errors) == 0 {
		return ""
	}

	summary := fmt.Sprintf("Error Summary (%d errors):\n", len(l.errors))
	for i, entry := range l.errors {
		summary += fmt.Sprintf("  %d. [%s] %s\n",
			i+1,
			entry.Time.Format("15:04:05"),
			entry.Message,
		)
	}
	return summary
}

// Close closes the log file.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Package-level functions for default logger

// Debug logs a debug message to the default logger.
func Debug(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(message, context)
	}
}

// Info logs an informational message to the default logger.
func Info(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(message, context)
	}
}

// Warn logs a warning message to the default logger.
func Warn(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(message, context)
	}
}

// Error logs an error message to the default logger.
func Error(message string, context map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(message, context)
	}
}

// GetErrors returns all logged errors from the default logger.
func GetErrors() []LogEntry {
	if defaultLogger != nil {
		return defaultLogger.GetErrors()
	}
	return nil
}

// GetErrorSummary returns a formatted error summary from the default logger.
func GetErrorSummary() string {
	if defaultLogger != nil {
		return defaultLogger.GetErrorSummary()
	}
	return ""
}

// Close closes the default logger.
func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}
