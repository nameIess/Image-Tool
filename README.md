# 🖼️ Image-Tool

A Windows TUI (Terminal User Interface) application for image and PDF processing. Built in Go using the Bubble Tea framework, it provides an intuitive interactive interface for batch image processing, format conversion, and compression.

## ✨ Features

- 📄 **PDF to Image Converter** - Convert PDF pages to images (PNG, JPG, BMP, TIFF, GIF)
- 🖼️ **Image Format Converter** - Convert between image formats (PNG, JPG, WebP, AVIF, BMP, TIFF, GIF)
- 🗜️ **Image/PDF Compressor** - Reduce file size by percentage or target size
- 🖥️ **Interactive TUI** - Beautiful terminal interface with keyboard navigation
- 📁 **Built-in File Picker** - Browse and select files without leaving the app
- 📁 **Batch Processing** - Process entire folders of files
- 🔄 **Drag-and-Drop Support** - Windows drag-and-drop functionality

### ImageMagick (Required)

ImageMagick v7.x is required for all image processing operations.

**Detection:** `magick -version`

**Download:** [imagemagick.org/script/download.php](https://imagemagick.org/script/download.php)

> ⚠️ Install manually and ensure it's in your system PATH. This application does not install dependencies automatically.

### Ghostscript (Required for PDF)

Ghostscript is required for PDF processing operations.

**Detection:** `gswin64c -version`

**Download:** [ghostscript.com/releases/gsdnld.html](https://ghostscript.com/releases/gsdnld.html)

> ⚠️ Install manually and ensure it's in your system PATH.

### Startup Dependency Check

On startup, the application verifies all dependencies:

```
Dependencies:
  ✔ ImageMagick (7.1.0-62)
  ✔ Ghostscript (10.02.1)
```

If dependencies are missing, clear instructions and download links are provided.

## 🛠️ Installation

### Option 1: Download Pre-built Binary

Download `Image-Tool.exe` from the [Releases](https://github.com/nameIess/Image-Tool/releases) page.

### Option 2: Build from Source

**Requirements:** Go 1.21 or higher

```bash
# Clone the repository
git clone https://github.com/nameIess/Image-Tool.git
cd Image-Tool

# Download dependencies
go mod download

# Build the executable
go build -o Image-Tool.exe ./cmd/imagetool

# Or build with a custom name
go build -o <custom_name>.exe ./cmd/imagetool
```

**[⬇️ Download ZIP](https://github.com/nameIess/Image-Tool/archive/refs/heads/master.zip)**

**Cross-compilation examples:**

```bash
# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o Image-Tool-windows-amd64.exe ./cmd/imagetool

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o Image-Tool-darwin-arm64 ./cmd/imagetool

# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o Image-Tool-linux-amd64 ./cmd/imagetool
```

### 3️⃣ Option 3: Install with Go

```bash
go install github.com/nameIess/Image-Tool/cmd/imagetool@latest
```

## 🚀 Usage

### ▶️ Run the Application

```bash
# Windows
Image-Tool.exe

# macOS/Linux
./Image-Tool
```

### ⌨️ Keyboard Navigation

| Key                 | Action                                |
| ------------------- | ------------------------------------- |
| `Up` / `k`          | Move up                               |
| `Down` / `j`        | Move down                             |
| `Enter`             | Select / Confirm                      |
| `Esc` / `Backspace` | Go back                               |
| `q` / `Ctrl+C`      | Quit                                  |
| `o`                 | Open output folder (after conversion) |

## 🔍 Features in Detail

### 📄 PDF to Image Converter

Convert multi-page PDFs to individual image files.

**Settings:**

- **Output Format:** PNG, JPG, JPEG, BMP, TIFF, GIF
- **Density (DPI):** Resolution quality (default: 180)
- **Quality:** Compression level 1-100 (default: 90)
- **Prefix:** Filename prefix for pages (default: `Page-`)

**Output:** Files are saved to `<PDF_name>_image/` folder

### 🖼️ Image Format Converter

Convert images between different formats.

**Supported Formats:** PNG, JPG, JPEG, WebP, AVIF, BMP, TIFF, GIF

**Output:** `<original_name>_conv.<new_format>`

### 🗜️ Image/PDF Compressor

Reduce file size using two methods:

1. **Percentage:** Target a percentage of original size (e.g., 50%)
2. **Fixed Size:** Target a specific file size (e.g., 500KB, 2MB)

**Output:** `<original_name>_comp.<ext>`

## 🔒 Security

This application follows strict security principles:

- ❌ No automatic installation of external tools
- ❌ No silent downloads
- ❌ No privilege escalation
- ❌ No PATH or registry modification
- ✅ User-managed dependencies only

## 👨‍💻 Development

### 🏗️ Building with Version Info (Windows)

To embed version information in the Windows executable:

```bash
# Install goversioninfo
go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest

# Generate resource file
cd cmd/imagetool
goversioninfo -o resource.syso versioninfo.json

# Build with embedded version info
cd ../..
go build -o Image-Tool.exe ./cmd/imagetool
```

### 🧪 Running Tests

```bash
go test ./...
```

### 🎨 Code Formatting

```bash
go fmt ./...
```

## 🛠️ Troubleshooting

### ❌ "magick: command not found"

ImageMagick is not installed or not in your PATH.

- **Windows:** Reinstall and check "Add to PATH"
- **macOS/Linux:** Ensure `/usr/local/bin` is in your PATH

### ❌ "gs: command not found" (PDF conversion fails)

Ghostscript is not installed.

- Install Ghostscript from the prerequisites section

### ⚠️ Conversion produces blank images

- Increase the **density** (DPI) value
- Ensure the PDF is not password-protected

## 📦 Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

Image-Tool is an independent project and is not affiliated with ImageMagick Studio LLC or Artifex Software, Inc.

All trademarks and software names belong to their respective owners.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
