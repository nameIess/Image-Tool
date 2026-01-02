# ImageTool

ImageTool is a cross-platform terminal application for converting PDF files to images, converting between image formats, and compressing images or PDFs—all with an intuitive, interactive text user interface. Built in Go using the Bubble Tea TUI framework, it streamlines batch image processing and PDF conversion for developers, designers, and power users.

## Features

- **PDF to Image Converter** - Convert PDF pages to images (PNG, JPG, BMP, TIFF, GIF)
- **Image Format Converter** - Convert between image formats (PNG, JPG, WebP, AVIF, BMP, TIFF, GIF)
- **Image/PDF Compressor** - Reduce file size by percentage or target size
- **Interactive TUI** - Beautiful terminal interface with keyboard navigation
- **Built-in File Picker** - Browse and select files without leaving the app

## Screenshots

```
┌─────────────────────────────────────────────────────────────┐
│                      ImageTool v1.0                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   [PDF]  PDF to Image Converter                             │
│          Convert PDF pages to images (PNG, JPG, etc.)       │
│                                                             │
│   [IMG]  Convert Image Format                               │
│          Convert images between formats (WebP, AVIF, etc.)  │
│                                                             │
│   [ZIP]  Compress Image/PDF                                 │
│          Reduce file size by percentage or target size      │
│                                                             │
│   [X]    Exit                                               │
│          Quit the application                               │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  Up/k up | Down/j down | enter select | q quit              │
└─────────────────────────────────────────────────────────────┘
```

## Prerequisites

Before using ImageTool, install the following dependencies:

### 1. ImageMagick (Required)

ImageMagick is used for all image processing operations.

**Windows:**

- Download from [ImageMagick Downloads](https://imagemagick.org/script/download.php)
- Run the installer and check **"Add application directory to your system path"**

**macOS:**

```bash
brew install imagemagick
```

**Linux (Ubuntu/Debian):**

```bash
sudo apt install imagemagick
```

### 2. Ghostscript (Required for PDF operations)

Ghostscript enables PDF to image conversion.

**Windows:**

- Download from [Ghostscript Downloads](https://www.ghostscript.com/releases/gsdnld.html)
- Install the appropriate version for your system

**macOS:**

```bash
brew install ghostscript
```

**Linux (Ubuntu/Debian):**

```bash
sudo apt install ghostscript
```

### 3. Verify Installation

```bash
magick -version
gs --version  # or gswin64c --version on Windows
```

## Installation

### Option 1: Download Pre-built Binary

Download the latest release from the [Releases](https://github.com/nameIess/Image-Tool/releases) page.

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

### Option 3: Install with Go

```bash
go install github.com/nameIess/Image-Tool/cmd/imagetool@latest
```

## Usage

### Run the Application

```bash
# Windows
Image-Tool.exe

# macOS/Linux
./Image-Tool
```

### Keyboard Navigation

| Key                 | Action                                |
| ------------------- | ------------------------------------- |
| `Up` / `k`          | Move up                               |
| `Down` / `j`        | Move down                             |
| `Enter`             | Select / Confirm                      |
| `Esc` / `Backspace` | Go back                               |
| `q` / `Ctrl+C`      | Quit                                  |
| `o`                 | Open output folder (after conversion) |

## Features in Detail

### PDF to Image Converter

Convert multi-page PDFs to individual image files.

**Settings:**

- **Output Format:** PNG, JPG, JPEG, BMP, TIFF, GIF
- **Density (DPI):** Resolution quality (default: 180)
- **Quality:** Compression level 1-100 (default: 90)
- **Prefix:** Filename prefix for pages (default: `Page-`)

**Output:** Files are saved to `<PDF_name>_image/` folder

### Image Format Converter

Convert images between different formats.

**Supported Formats:** PNG, JPG, JPEG, WebP, AVIF, BMP, TIFF, GIF

**Output:** `<original_name>_conv.<new_format>`

### Image/PDF Compressor

Reduce file size using two methods:

1. **Percentage:** Target a percentage of original size (e.g., 50%)
2. **Fixed Size:** Target a specific file size (e.g., 500KB, 2MB)

**Output:** `<original_name>_comp.<ext>`

## Project Structure

```
imagetool/
├── cmd/
│   └── imagetool/
│       ├── main.go              # Application entry point
│       └── versioninfo.json     # Windows version info
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration and defaults
│   └── tui/
│       ├── app.go               # Main TUI application
│       ├── filepicker.go        # File browser component
│       ├── pdf_converter.go     # PDF to image converter
│       ├── format_converter.go  # Image format converter
│       ├── compressor.go        # File compressor
│       ├── styles.go            # UI styles and themes
│       └── utils.go             # Utility functions
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Development

### Building with Version Info (Windows)

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

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

## Troubleshooting

### "magick: command not found"

ImageMagick is not installed or not in your PATH.

- **Windows:** Reinstall and check "Add to PATH"
- **macOS/Linux:** Ensure `/usr/local/bin` is in your PATH

### "gs: command not found" (PDF conversion fails)

Ghostscript is not installed.

- Install Ghostscript from the prerequisites section

### Conversion produces blank images

- Increase the **density** (DPI) value
- Ensure the PDF is not password-protected

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
