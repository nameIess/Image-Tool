package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1572864, "1.5 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		result := FormatSize(tt.bytes)
		if result != tt.expected {
			t.Errorf("FormatSize(%d) = %s; want %s", tt.bytes, result, tt.expected)
		}
	}
}

func TestIsImageFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"test.jpg", true},
		{"test.jpeg", true},
		{"test.png", true},
		{"test.gif", true},
		{"test.bmp", true},
		{"test.webp", true},
		{"test.avif", true},
		{"test.tiff", true},
		{"test.tif", true},
		{"test.pdf", false},
		{"test.txt", false},
		{"test.doc", false},
		{"TEST.JPG", true},
		{"TEST.PNG", true},
	}

	for _, tt := range tests {
		result := IsImageFile(tt.path)
		if result != tt.expected {
			t.Errorf("IsImageFile(%s) = %v; want %v", tt.path, result, tt.expected)
		}
	}
}

func TestIsPDFFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"test.pdf", true},
		{"test.PDF", true},
		{"test.Pdf", true},
		{"test.jpg", false},
		{"test.txt", false},
		{"document.pdf", true},
	}

	for _, tt := range tests {
		result := IsPDFFile(tt.path)
		if result != tt.expected {
			t.Errorf("IsPDFFile(%s) = %v; want %v", tt.path, result, tt.expected)
		}
	}
}

func TestGenerateOutputPath(t *testing.T) {
	tests := []struct {
		input    string
		suffix   string
		ext      string
		expected string
	}{
		{filepath.Join("dir", "test.jpg"), "_conv", "png", filepath.Join("dir", "test_conv.png")},
		{filepath.Join("path", "to", "file.png"), "_comp", "jpg", filepath.Join("path", "to", "file_comp.jpg")},
	}

	for _, tt := range tests {
		result := generateOutputPath(tt.input, tt.suffix, tt.ext)
		if result != tt.expected {
			t.Errorf("generateOutputPath(%s, %s, %s) = %s; want %s",
				tt.input, tt.suffix, tt.ext, result, tt.expected)
		}
	}
}

func TestGetFilesInDirectory(t *testing.T) {
	// Create temp directory with test files
	tempDir, err := os.MkdirTemp("", "imagetool-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{"test1.jpg", "test2.png", "document.pdf", "readme.txt"}
	for _, name := range testFiles {
		f, err := os.Create(filepath.Join(tempDir, name))
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		f.Close()
	}

	// Test images only
	files, err := GetFilesInDirectory(tempDir, true, false)
	if err != nil {
		t.Fatalf("GetFilesInDirectory failed: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("Expected 2 image files, got %d", len(files))
	}

	// Test PDFs only
	files, err = GetFilesInDirectory(tempDir, false, true)
	if err != nil {
		t.Fatalf("GetFilesInDirectory failed: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 PDF file, got %d", len(files))
	}

	// Test both
	files, err = GetFilesInDirectory(tempDir, true, true)
	if err != nil {
		t.Fatalf("GetFilesInDirectory failed: %v", err)
	}
	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
}
