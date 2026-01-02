// Package config handles persistent configuration for Image-Tool.
// Configuration is stored in a JSON file in the user's config directory.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Default configuration values
const (
	DefaultOutputFormat    = "png"
	DefaultDensity         = 180
	DefaultQuality         = 90
	DefaultPrefix          = "Page-"
	DefaultCompressPercent = 75

	// Density limits for PDF conversion
	MinDensity = 72
	MaxDensity = 600

	// Config file name
	configFileName = "imagetool_config.json"
)

// SupportedImageFormats lists all formats for image conversion
var SupportedImageFormats = []string{"png", "jpg", "jpeg", "webp", "avif", "bmp", "tiff", "gif"}

// SupportedPDFOutputFormats lists formats for PDF export
var SupportedPDFOutputFormats = []string{"png", "jpg", "jpeg", "bmp", "tiff", "gif"}

// Config holds user preferences
type Config struct {
	// PDF conversion settings
	OutputFormat string `json:"output_format"`
	Density      int    `json:"density"`
	Quality      int    `json:"quality"`
	Prefix       string `json:"prefix"`

	// Compression settings
	CompressPercent int `json:"compress_percent"`

	// UI preferences
	LastDirectory string `json:"last_directory,omitempty"`

	// Internal: config file path (not persisted)
	filePath string `json:"-"`
}

// NewConfig creates a config with default values.
func NewConfig() *Config {
	return &Config{
		OutputFormat:    DefaultOutputFormat,
		Density:         DefaultDensity,
		Quality:         DefaultQuality,
		Prefix:          DefaultPrefix,
		CompressPercent: DefaultCompressPercent,
	}
}

// Load reads configuration from disk. Returns default config if file doesn't exist.
func Load() (*Config, error) {
	cfg := NewConfig()
	cfg.filePath = getConfigPath()

	data, err := os.ReadFile(cfg.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// No config file yet, return defaults
			return cfg, nil
		}
		return cfg, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		// Return defaults if config is corrupted
		return NewConfig(), nil
	}

	// Validate loaded values
	cfg.validate()

	return cfg, nil
}

// Save writes configuration to disk.
func (c *Config) Save() error {
	if c.filePath == "" {
		c.filePath = getConfigPath()
	}

	// Ensure directory exists
	dir := filepath.Dir(c.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.filePath, data, 0644)
}

// validate ensures all values are within acceptable ranges.
func (c *Config) validate() {
	if c.Density < MinDensity {
		c.Density = MinDensity
	}
	if c.Density > MaxDensity {
		c.Density = MaxDensity
	}
	if c.Quality < 1 {
		c.Quality = 1
	}
	if c.Quality > 100 {
		c.Quality = 100
	}
	if c.CompressPercent < 1 {
		c.CompressPercent = 1
	}
	if c.CompressPercent > 100 {
		c.CompressPercent = 100
	}
	if c.Prefix == "" {
		c.Prefix = DefaultPrefix
	}
	if c.OutputFormat == "" {
		c.OutputFormat = DefaultOutputFormat
	}
}

// Reset restores default values.
func (c *Config) Reset() {
	c.OutputFormat = DefaultOutputFormat
	c.Density = DefaultDensity
	c.Quality = DefaultQuality
	c.Prefix = DefaultPrefix
	c.CompressPercent = DefaultCompressPercent
	c.LastDirectory = ""
}

// getConfigPath returns the path to the config file.
func getConfigPath() string {
	// Try user's config directory first
	configDir, err := os.UserConfigDir()
	if err == nil {
		return filepath.Join(configDir, "ImageTool", configFileName)
	}

	// Fallback to executable directory
	execPath, err := os.Executable()
	if err == nil {
		return filepath.Join(filepath.Dir(execPath), configFileName)
	}

	// Last resort: current directory
	return configFileName
}

// GetConfigDir returns the directory containing the config file.
func GetConfigDir() string {
	return filepath.Dir(getConfigPath())
}
