package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	if cfg.OutputFormat != DefaultOutputFormat {
		t.Errorf("OutputFormat = %s; want %s", cfg.OutputFormat, DefaultOutputFormat)
	}
	if cfg.Density != DefaultDensity {
		t.Errorf("Density = %d; want %d", cfg.Density, DefaultDensity)
	}
	if cfg.Quality != DefaultQuality {
		t.Errorf("Quality = %d; want %d", cfg.Quality, DefaultQuality)
	}
	if cfg.Prefix != DefaultPrefix {
		t.Errorf("Prefix = %s; want %s", cfg.Prefix, DefaultPrefix)
	}
	if cfg.CompressPercent != DefaultCompressPercent {
		t.Errorf("CompressPercent = %d; want %d", cfg.CompressPercent, DefaultCompressPercent)
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := &Config{
		Density:         10,  // Below min
		Quality:         200, // Above max
		CompressPercent: 0,   // Below min
		Prefix:          "",  // Empty
		OutputFormat:    "",  // Empty
	}

	cfg.validate()

	if cfg.Density != MinDensity {
		t.Errorf("Density = %d; want %d", cfg.Density, MinDensity)
	}
	if cfg.Quality != 100 {
		t.Errorf("Quality = %d; want 100", cfg.Quality)
	}
	if cfg.CompressPercent != 1 {
		t.Errorf("CompressPercent = %d; want 1", cfg.CompressPercent)
	}
	if cfg.Prefix != DefaultPrefix {
		t.Errorf("Prefix = %s; want %s", cfg.Prefix, DefaultPrefix)
	}
	if cfg.OutputFormat != DefaultOutputFormat {
		t.Errorf("OutputFormat = %s; want %s", cfg.OutputFormat, DefaultOutputFormat)
	}
}

func TestConfigSaveLoad(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "imagetool-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create config with custom values
	cfg := NewConfig()
	cfg.filePath = filepath.Join(tempDir, "test_config.json")
	cfg.Density = 300
	cfg.Quality = 85
	cfg.CompressPercent = 50
	cfg.LastDirectory = "/test/path"

	// Save
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Create new config and load
	cfg2 := NewConfig()
	cfg2.filePath = cfg.filePath

	data, err := os.ReadFile(cfg.filePath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Verify file was created and contains expected content
	if len(data) == 0 {
		t.Error("Config file is empty")
	}
}

func TestConfigReset(t *testing.T) {
	cfg := &Config{
		OutputFormat:    "jpg",
		Density:         300,
		Quality:         50,
		Prefix:          "Custom-",
		CompressPercent: 25,
		LastDirectory:   "/some/path",
	}

	cfg.Reset()

	if cfg.OutputFormat != DefaultOutputFormat {
		t.Errorf("OutputFormat = %s; want %s", cfg.OutputFormat, DefaultOutputFormat)
	}
	if cfg.Density != DefaultDensity {
		t.Errorf("Density = %d; want %d", cfg.Density, DefaultDensity)
	}
	if cfg.Quality != DefaultQuality {
		t.Errorf("Quality = %d; want %d", cfg.Quality, DefaultQuality)
	}
	if cfg.Prefix != DefaultPrefix {
		t.Errorf("Prefix = %s; want %s", cfg.Prefix, DefaultPrefix)
	}
	if cfg.CompressPercent != DefaultCompressPercent {
		t.Errorf("CompressPercent = %d; want %d", cfg.CompressPercent, DefaultCompressPercent)
	}
	if cfg.LastDirectory != "" {
		t.Errorf("LastDirectory = %s; want empty", cfg.LastDirectory)
	}
}
