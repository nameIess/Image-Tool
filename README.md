# Image Converter Suite (Batch)

A Windows batch script powered by ImageMagick that lets you:

- Convert PDFs to images, and
- Convert images to a different format (e.g., PNG, JPG, AVIF, WEBP, TIFF, BMP).

Interactive prompts, sensible defaults, and a simple main menu make it fast to use.

## Features

- ✅ **Main Menu with Defaults**: Choose operation by number; pressing Enter selects the default option
- ✅ **Interactive Prompts**: All settings support defaults; just press Enter to accept
- ✅ **Local File Detection**:
  - PDFs in script folder are listed for quick selection (PDF flow)
  - Images in script folder are listed for quick selection (Image conversion flow)
- ✅ **Per‑PDF Output Folders**: PDF conversions save into `<PDF name>_image` (e.g., `Report_image`)
- ✅ **Auto‑open Output (PDF only)**: After a successful PDF conversion, the output folder opens in Explorer
- ✅ **Image Format Conversion**: Convert any supported image to another format
  - Options: `jpeg`, `jpg`, or `custom` (enter `avif`, `webp`, `tiff`, `bmp`, etc.)
  - Output file name: `<original_name>_conv.<format>`
  - Shows converted file size in bytes/KB/MB
  - Does NOT auto‑open the folder
- ✅ **Quoted/Absolute Path Support**: Paste full paths with or without quotes; the script normalizes them
- ✅ **Multiple Format Support**:
  - PDF export to: PNG, JPG, JPEG, BMP, TIFF, GIF
  - Image conversion to: any format supported by your ImageMagick build
- ✅ **Customizable Settings (PDF)**: Adjustable output quality (1-100) and density (DPI)
- ✅ **Input Validation**: Clear errors and helpful guidance
- ✅ **Multi‑page PDFs**: Custom filename prefixes for page outputs
- ✅ **Smart Fallback**: If folder creation fails, saves in current directory

## Installation Instructions

### Prerequisites

1. **Install ImageMagick:**

   - Visit: [ImageMagick](https://imagemagick.org/script/download.php)
   - Download the Windows version (choose x64 or x86 based on your system)
   - Run the installer as Administrator
   - ✅ **Important**: Check "Add application directory to your system path"

2. **Install Ghostscript (Required for PDF handling):**

   - Visit: [Ghostscript](https://www.ghostscript.com/releases/gsdnld.html)
   - Download and install the appropriate version for your system

3. **Clone or Download the Repository:**

   ```cmd
   git clone https://github.com/nameIess/Image-Tool.git
   cd Image-Tool
   ```

   Download: [Image-Tool](https://github.com/nameIess/Image-Tool/archive/refs/heads/master.zip)

4. **Verify Installation:**
   ```cmd
   magick -version
   gswin64c --version
   ```
   Both commands should display version information if installed correctly.

## Project Structure

```
Image-Tool/
├── Convert.bat          # Main conversion script
└── README.md            # This documentation file
```

## Usage Instructions

### Start

Run the script and pick an operation from the main menu (Enter defaults to 1):

```
=========================================================
*              Image Converter Suite                    *
*                Using ImageMagick                      *
=========================================================

Select an operation:
  1. PDF to Image Converter
  2. Convert Image Format
  0. Exit

Enter your choice (default: 1):
```

### PDF to Image

1. Optional: Pick a local PDF (auto‑detected):

   - If there are any `*.pdf` files in the same folder as `Convert.bat`, the script will list them like:

     ```
     -- Local PDF detection --
     Found 2 PDF(s) in script directory:
     1. Example1.pdf
     2. Example2.pdf

     Select a PDF to use (1-2) or press Enter to skip: 1
     Selected: F:\\Code\\Image-Tool\\Example1.pdf
     Using selected PDF: F:\\Code\\Image-Tool\\Example1.pdf
     ```

   - Enter a number to select and continue. Press Enter to skip and type a filename manually.

2. Follow the interactive prompts:

   - **PDF filename**: Enter the full path or filename (e.g., `document.pdf`)
   - **Output format**: Choose from png, jpg, jpeg, bmp, tiff, gif (default: png)
   - **Density**: Set DPI value for image resolution (default: 180)
   - **Quality**: Set compression quality 1-100 (default: 90)
   - **Filename prefix**: Set prefix for output files (default: Page-)

3. Review and confirm:
   - Check the conversion summary
   - Confirm to proceed with conversion

#### Example (PDF to PNG)

```
=========================================================
*                    PDF to Image Converter             *
*                      Using ImageMagick                *
=========================================================

-- Local PDF detection --
Found 1 PDF(s) in script directory:
1. Example.pdf

Select a PDF to use (1-1) or press Enter to skip: 1
Selected: F:\Code\Image-Tool\Example.pdf
Using selected PDF: F:\Code\Image-Tool\Example.pdf

-- Input Parameters --

Enter output image format (default: png): png
Enter density value (default: 180): 300
Enter quality value 1-100 (default: 90): 95
Enter output filename prefix (default: Page-): Doc-

..*..Conversion Summary..*..

Input PDF:      F:\Code\Image-Tool\Example.pdf
Output format:  .png
Density:        300
Quality:        95
Prefix:         Doc-
Output folder:  Example_image\
Command:        magick -density 300 "F:\\Code\\Image-Tool\\Example.pdf" -quality 95 "Example_image\Doc-%d.png"

Proceed with conversion? (y/N): y
```

Result: Creates `Example_image/Doc-0.png`, `Example_image/Doc-1.png`, etc., then opens that folder automatically.

### Convert Image Format

Convert any supported image to another format. Output naming: `<original_name>_conv.<format>`.

1. Optional: Pick a local image (auto‑detected in the script folder), or press Enter and paste a path.
2. Choose output format:
   - 1. `jpeg`
   - 2. `jpg`
   - 3. `custom` → type anything supported (e.g., `avif`, `webp`, `tiff`, `bmp`)
3. Conversion runs and shows the resulting file size in bytes, KB, and MB.

Notes:

- The image conversion flow does NOT auto‑open Explorer.
- Size display helps decide next actions (like manual compression outside the script).

## Notes and Troubleshooting

### Important Notes

- 📁 **Output Location (PDF)**: Converted images are saved in a per‑PDF folder named `<PDF name>_image` in the current working directory. If folder creation fails, output falls back to the current directory.
- �️ **Output Location (Image Conversion)**: Converted images are written next to the source image (same folder) with `_conv` in the filename.
- �🗂️ **Local File Detection**: At startup, the script can list local PDFs (PDF flow) or local images (image conversion flow). To use files in a different folder, enter their full paths.
- 🔄 **Automatic Folder Creation**: The script creates the per‑PDF output folder if it doesn't exist
- 📋 **Multi-page Support**: Each page of the PDF becomes a separate image file
- 🔢 **File Naming**: Output files use the format: `[prefix][page-number].[format]`
- ⚡ **Performance**: Higher density values create larger, higher-quality images but take more time

### Common Issues & Solutions

1. **"ImageMagick is not installed or not in PATH!"**

   - ✅ **Solution**: Install ImageMagick and ensure it's added to system PATH
   - ✅ **Verify**: Run `magick -version` in command prompt

2. **"Failed to create PDF_Images folder"**

   - ✅ **Cause**: Insufficient permissions or OneDrive sync issues
   - ✅ **Solution**: Run as Administrator or choose to continue (saves to current directory)

3. **"File does not exist!"**

   - ✅ **Solution**: Check filename spelling and file location
   - ✅ **Tip**: Use full file path if PDF is in different directory

4. **"Conversion failed!"**

   - ✅ **Common causes**:
     - PDF is password protected
     - PDF file is corrupted
     - Insufficient disk space
     - Missing Ghostscript installation
   - ✅ **Solution**: Install Ghostscript and verify PDF file integrity

5. **Poor output quality or large file sizes**
   - ✅ **For better quality**: Increase density (e.g., 300+ DPI)
   - ✅ **For smaller files**: Decrease quality (70-80) or density (72-150)

### Performance Tips

- **Large PDFs**: Use density 72-150 for faster processing
- **Print Quality**: Use density 300+ and quality 90-100
- **Web Use**: Use density 72-96 and quality 70-85
- **Disk Space**: Monitor available space for high-density conversions

### System Requirements

- **OS**: Windows 7 or later
- **ImageMagick**: Version 7.0 or later
- **Ghostscript**: Latest version recommended
- **Disk Space**: Varies based on PDF size and output settings

## Support

If you encounter any issues:

1. Verify ImageMagick installation
2. Ensure PDF file is not corrupted or password-protected

## Version History

- **v1.3** - Introduced main menu; added Image Format Conversion (jpeg/jpg/custom); shows converted image file size; PDF conversions auto‑open folder; image conversions do not.
- **v1.2** - Per‑PDF output folders (`<name>_image`), automatic opening of the output folder after success, improved robust handling of quoted/absolute input paths, and a single final exit instead of multiple pauses.
- **v1.1** - Added local PDF detection and selection at startup; selected file is used automatically for conversion (no auto-open).
- **v1.0** - Initial release with basic PDF to image conversion functionality.

## Roadmap

- Compression (next):
  - Optional post‑conversion compression for images
  - Presets (Light/Medium/Heavy/Maximum) and Lossless mode
  - Output naming like `<original_name>_conv_comp.<format>` or `<original_name>_comp.<ext>`
  - File size comparison before/after
