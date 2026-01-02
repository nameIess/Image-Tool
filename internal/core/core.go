// Package core provides the business logic for image and PDF processing.
// This layer handles conversion, compression, and workflows without any TUI code.
package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ImageFormat represents a supported image format.
type ImageFormat string

const (
	FormatPNG  ImageFormat = "png"
	FormatJPG  ImageFormat = "jpg"
	FormatJPEG ImageFormat = "jpeg"
	FormatWebP ImageFormat = "webp"
	FormatAVIF ImageFormat = "avif"
	FormatBMP  ImageFormat = "bmp"
	FormatTIFF ImageFormat = "tiff"
	FormatGIF  ImageFormat = "gif"
)

// SupportedImageFormats lists all formats for image conversion.
var SupportedImageFormats = []ImageFormat{
	FormatPNG, FormatJPG, FormatJPEG, FormatWebP,
	FormatAVIF, FormatBMP, FormatTIFF, FormatGIF,
}

// SupportedPDFOutputFormats lists formats for PDF export.
var SupportedPDFOutputFormats = []ImageFormat{
	FormatPNG, FormatJPG, FormatJPEG, FormatBMP, FormatTIFF, FormatGIF,
}

// Result represents the outcome of a processing operation.
type Result struct {
	Success     bool
	Message     string
	OutputPath  string
	OutputPaths []string
	OutputSize  int64
	Error       error
}

// ConvertImageOptions contains options for image format conversion.
type ConvertImageOptions struct {
	InputPath    string
	OutputFormat ImageFormat
	OutputPath   string // Optional, will be auto-generated if empty
}

// ConvertImage converts an image to a different format using ImageMagick.
func ConvertImage(opts ConvertImageOptions) Result {
	if opts.OutputPath == "" {
		opts.OutputPath = generateOutputPath(opts.InputPath, "_conv", string(opts.OutputFormat))
	}

	cmd := exec.Command("magick", opts.InputPath, opts.OutputPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Conversion failed: %v", err),
			Error:   fmt.Errorf("%v: %s", err, string(output)),
		}
	}

	info, err := os.Stat(opts.OutputPath)
	var size int64
	if err == nil {
		size = info.Size()
	}

	return Result{
		Success:    true,
		Message:    "Image converted successfully",
		OutputPath: opts.OutputPath,
		OutputSize: size,
	}
}

// ConvertPDFOptions contains options for PDF to image conversion.
type ConvertPDFOptions struct {
	InputPath    string
	OutputFormat ImageFormat
	OutputDir    string // Directory for output images
	Density      int    // DPI (72-600)
	Quality      int    // Output quality (1-100)
	Prefix       string // Filename prefix for output images
}

// ConvertPDFToImages converts a PDF to images using ImageMagick.
func ConvertPDFToImages(opts ConvertPDFOptions) Result {
	// Set defaults
	if opts.Density < 72 {
		opts.Density = 180
	}
	if opts.Density > 600 {
		opts.Density = 600
	}
	if opts.Quality < 1 || opts.Quality > 100 {
		opts.Quality = 90
	}
	if opts.Prefix == "" {
		opts.Prefix = "Page-"
	}
	if opts.OutputDir == "" {
		base := strings.TrimSuffix(filepath.Base(opts.InputPath), filepath.Ext(opts.InputPath))
		opts.OutputDir = filepath.Join(filepath.Dir(opts.InputPath), base+"_images")
	}

	// Create output directory
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to create output directory: %v", err),
			Error:   err,
		}
	}

	// Build output pattern
	outputPattern := filepath.Join(opts.OutputDir, opts.Prefix+"%d."+string(opts.OutputFormat))

	// Run ImageMagick
	cmd := exec.Command("magick",
		"-density", fmt.Sprintf("%d", opts.Density),
		opts.InputPath,
		"-quality", fmt.Sprintf("%d", opts.Quality),
		outputPattern,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Conversion failed: %v", err),
			Error:   fmt.Errorf("%v: %s", err, string(output)),
		}
	}

	// Count output files
	pattern := filepath.Join(opts.OutputDir, opts.Prefix+"*."+string(opts.OutputFormat))
	matches, _ := filepath.Glob(pattern)

	return Result{
		Success:     true,
		Message:     fmt.Sprintf("Successfully converted %d page(s)", len(matches)),
		OutputPath:  opts.OutputDir,
		OutputPaths: matches,
	}
}

// CompressMethod defines how compression targets are specified.
type CompressMethod int

const (
	// CompressMethodPercent compresses to a percentage of original size.
	CompressMethodPercent CompressMethod = iota
	// CompressMethodFixedSize compresses to a specific file size.
	CompressMethodFixedSize
)

// CompressOptions contains options for file compression.
type CompressOptions struct {
	InputPath     string
	Method        CompressMethod
	TargetPercent int   // For CompressMethodPercent (1-100)
	TargetBytes   int64 // For CompressMethodFixedSize
	OutputPath    string
}

// CompressFile compresses an image or PDF using ImageMagick.
func CompressFile(opts CompressOptions) Result {
	// Get input file size
	inputInfo, err := os.Stat(opts.InputPath)
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Failed to read input file: %v", err),
			Error:   err,
		}
	}
	inputSize := inputInfo.Size()

	// Calculate target bytes if using percentage
	var targetBytes int64
	if opts.Method == CompressMethodPercent {
		if opts.TargetPercent < 1 {
			opts.TargetPercent = 1
		}
		if opts.TargetPercent > 100 {
			opts.TargetPercent = 100
		}
		targetBytes = inputSize * int64(opts.TargetPercent) / 100
	} else {
		targetBytes = opts.TargetBytes
	}

	// Generate output path if not provided
	if opts.OutputPath == "" {
		ext := filepath.Ext(opts.InputPath)
		// For images, output as JPG for better compression
		if strings.ToLower(ext) != ".pdf" {
			ext = ".jpg"
		}
		opts.OutputPath = generateOutputPath(opts.InputPath, "_comp", ext[1:])
	}

	// Build ImageMagick arguments
	ext := strings.ToLower(filepath.Ext(opts.OutputPath))
	args := []string{opts.InputPath}

	// Use jpeg:extent for target size compression when output is JPEG
	if ext == ".jpg" || ext == ".jpeg" {
		targetSize := fmt.Sprintf("%d", targetBytes)
		args = append(args, "-define", fmt.Sprintf("jpeg:extent=%s", targetSize))
	}
	args = append(args, opts.OutputPath)

	cmd := exec.Command("magick", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Result{
			Success: false,
			Message: fmt.Sprintf("Compression failed: %v", err),
			Error:   fmt.Errorf("%v: %s", err, string(output)),
		}
	}

	// Get output file size
	outputInfo, err := os.Stat(opts.OutputPath)
	var outputSize int64
	if err == nil {
		outputSize = outputInfo.Size()
	}

	// Calculate reduction
	reduction := 0.0
	if inputSize > 0 {
		reduction = float64(inputSize-outputSize) / float64(inputSize) * 100
	}

	msg := fmt.Sprintf("Compressed successfully! Reduced by %.1f%%", reduction)
	if outputSize > targetBytes {
		msg = fmt.Sprintf("Compressed, but couldn't reach target. Best possible: %s", FormatSize(outputSize))
	}

	return Result{
		Success:    true,
		Message:    msg,
		OutputPath: opts.OutputPath,
		OutputSize: outputSize,
	}
}

// BatchResult contains results for batch operations.
type BatchResult struct {
	TotalFiles      int
	SuccessCount    int
	FailCount       int
	Results         []Result
	TotalInputSize  int64
	TotalOutputSize int64
}

// BatchConvertImages converts multiple images to a different format.
func BatchConvertImages(inputPaths []string, outputFormat ImageFormat) BatchResult {
	batch := BatchResult{
		TotalFiles: len(inputPaths),
		Results:    make([]Result, 0, len(inputPaths)),
	}

	for _, inputPath := range inputPaths {
		result := ConvertImage(ConvertImageOptions{
			InputPath:    inputPath,
			OutputFormat: outputFormat,
		})
		batch.Results = append(batch.Results, result)

		if result.Success {
			batch.SuccessCount++
			batch.TotalOutputSize += result.OutputSize
		} else {
			batch.FailCount++
		}

		// Track input size
		if info, err := os.Stat(inputPath); err == nil {
			batch.TotalInputSize += info.Size()
		}
	}

	return batch
}

// generateOutputPath creates an output path with suffix and extension.
func generateOutputPath(inputPath, suffix, ext string) string {
	dir := filepath.Dir(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	return filepath.Join(dir, base+suffix+"."+ext)
}

// FormatSize formats bytes to human readable string.
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// OpenFolder opens a folder in Windows Explorer.
func OpenFolder(path string) error {
	cmd := exec.Command("explorer", path)
	return cmd.Start()
}

// IsImageFile checks if a file is a supported image format.
func IsImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp", ".avif"}
	for _, e := range imageExts {
		if ext == e {
			return true
		}
	}
	return false
}

// IsPDFFile checks if a file is a PDF.
func IsPDFFile(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".pdf"
}

// GetFilesInDirectory returns files matching the filter in a directory.
func GetFilesInDirectory(dir string, includeImages, includePDFs bool) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		if includeImages && IsImageFile(path) {
			files = append(files, path)
		}
		if includePDFs && IsPDFFile(path) {
			files = append(files, path)
		}
	}

	return files, nil
}
