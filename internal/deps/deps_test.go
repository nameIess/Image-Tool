package deps

import (
	"testing"
)

func TestParseImageMagickVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Version: ImageMagick 7.1.0-62 Q16-HDRI x64 b0a9f2e:20231125 https://imagemagick.org",
			"7.1.0-62",
		},
		{
			"Version: ImageMagick 7.0.0-0 Q16 x64 https://imagemagick.org",
			"7.0.0-0",
		},
		{
			"Some random text",
			"",
		},
		{
			"",
			"",
		},
	}

	for _, tt := range tests {
		result := parseImageMagickVersion(tt.input)
		if result != tt.expected {
			t.Errorf("parseImageMagickVersion(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParseGhostscriptVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"GPL Ghostscript 10.02.1 (2023-11-01)",
			"10.02.1",
		},
		{
			"GPL Ghostscript 9.56.1",
			"9.56.1",
		},
		{
			"Some random text",
			"",
		},
		{
			"",
			"",
		},
	}

	for _, tt := range tests {
		result := parseGhostscriptVersion(tt.input)
		if result != tt.expected {
			t.Errorf("parseGhostscriptVersion(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDependencyFormatStatus(t *testing.T) {
	tests := []struct {
		dep      Dependency
		contains string
	}{
		{
			Dependency{Name: "ImageMagick", Status: StatusOK, Version: "7.1.0"},
			"✔ ImageMagick (7.1.0)",
		},
		{
			Dependency{Name: "ImageMagick", Status: StatusOK, Version: ""},
			"✔ ImageMagick",
		},
		{
			Dependency{Name: "ImageMagick", Status: StatusNotFound},
			"✗ ImageMagick - Not found",
		},
	}

	for _, tt := range tests {
		result := tt.dep.FormatStatus()
		if result != tt.contains {
			t.Errorf("FormatStatus() = %q; want %q", result, tt.contains)
		}
	}
}
