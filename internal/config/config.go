package config

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
)

// SupportedImageFormats lists all formats for image conversion
var SupportedImageFormats = []string{"png", "jpg", "jpeg", "webp", "avif", "bmp", "tiff", "gif"}

// SupportedPDFOutputFormats lists formats for PDF export
var SupportedPDFOutputFormats = []string{"png", "jpg", "jpeg", "bmp", "tiff", "gif"}

// Config holds user preferences
type Config struct {
	OutputFormat    string
	Density         int
	Quality         int
	Prefix          string
	CompressPercent int
}

// NewConfig creates a config with default values
func NewConfig() *Config {
	return &Config{
		OutputFormat:    DefaultOutputFormat,
		Density:         DefaultDensity,
		Quality:         DefaultQuality,
		Prefix:          DefaultPrefix,
		CompressPercent: DefaultCompressPercent,
	}
}
