# PDF to Image Converter

A Windows batch script that converts PDF files to various image formats using ImageMagick. This tool provides an interactive command-line interface for easy PDF to image conversion with customizable settings.

## Features

- ✅ **Interactive User Prompts**: Prompts user for input PDF filename, output format, density, quality, and filename prefix
- ✅ **Local PDF Detection & Selection**: Automatically lists PDF files in the script's folder at startup and lets you select one to use (skips the filename prompt)
- ✅ **Automatic Folder Management**: Automatically creates `PDF_Images` folder if it does not exist
- ✅ **Organized Output**: Saves all converted images inside the `PDF_Images` directory
- ✅ **Multiple Format Support**: Convert PDF files to PNG, JPG, JPEG, BMP, TIFF, GIF formats
- ✅ **Customizable Settings**: Adjustable output quality (1-100) and density (DPI) settings
- ✅ **Input Validation**: Comprehensive error checking and user-friendly feedback
- ✅ **Batch Processing**: Convert multi-page PDFs with custom filename prefixes
- ✅ **Smart Fallback**: Graceful handling when folder creation fails

## Installation Instructions

### Prerequisites

1. **Install ImageMagick:**

   - Visit: [ImageMagick](https://imagemagick.org/script/download.php)
   - Download the Windows version (choose x64 or x86 based on your system)
   - Run the installer as Administrator
   - ✅ **Important**: Check "Install ImageMagick for all users"
   - ✅ **Important**: Check "Add application directory to your system path"

2. **Install Ghostscript (Required for PDF handling):**

   - Visit: [Ghostscript](https://www.ghostscript.com/releases/gsdnld.html)
   - Download and install the appropriate version for your system

3. **Clone or Download the Repository:**

   ```cmd
   git clone <repository-url>
   cd pdf_converter
   ```

4. **Verify Installation:**
   ```cmd
   magick -version
   gswin64c --version
   ```
   Both commands should display version information if installed correctly.

## Project Structure

```
pdf_converter/
├── Convert.bat          # Main conversion script
├── README.md            # This documentation file
└── PDF_Images/          # Output directory for converted images (created on first run)
```

## Usage Instructions

### Quick Start

1. **Run the Script:**

   ```cmd
   Convert.bat
   ```

2. **Optional: Pick a local PDF (auto-detected):**

   - If there are any `*.pdf` files in the same folder as `Convert.bat`, the script will list them like:

     ```
     -- Local PDF detection --
     Found 2 PDF(s) in script directory:
     1. Example1.pdf
     2. Example2.pdf

     Select a PDF to use (1-2) or press Enter to skip: 1
     Selected: F:\Code\pdf_converter\Example1.pdf
     Using selected PDF: F:\Code\pdf_converter\Example1.pdf
     ```

   - Enter a number to select and continue. Press Enter to skip and type a filename manually.

3. **Follow the Interactive Prompts:**

   - **PDF filename**: Enter the full path or filename (e.g., `document.pdf`)
   - **Output format**: Choose from png, jpg, jpeg, bmp, tiff, gif (default: png)
   - **Density**: Set DPI value for image resolution (default: 180)
   - **Quality**: Set compression quality 1-100 (default: 90)
   - **Filename prefix**: Set prefix for output files (default: Page-)

4. **Review and Confirm:**
   - Check the conversion summary
   - Confirm to proceed with conversion

### Example Usage

```
=========================================================
*                    PDF to Image Converter             *
*                      Using ImageMagick                *
=========================================================

-- Local PDF detection --
Found 1 PDF(s) in script directory:
1. Example.pdf

Select a PDF to use (1-1) or press Enter to skip: 1
Selected: F:\Code\pdf_converter\Example.pdf
Using selected PDF: F:\Code\pdf_converter\Example.pdf

┌─ Input Parameters ─┐

Enter output image format (default: png): png
Enter density value (default: 180): 300
Enter quality value 1-100 (default: 90): 95
Enter output filename prefix (default: Page-): Doc-

┌─ Conversion Summary ─┐

Input PDF:      F:\Code\pdf_converter\Example.pdf
Output format:  .png
Density:        300
Quality:        95
Prefix:         Doc-
Output folder:  PDF_Images\
Command:        magick -density 300 "F:\\Code\\pdf_converter\\Example.pdf" -quality 95 "PDF_Images\Doc-%d.png"

Proceed with conversion? (Y/N): Y
```

**Result**: Creates `PDF_Images/Doc-0.png`, `PDF_Images/Doc-1.png`, etc.

## Notes and Troubleshooting

### Important Notes

- 📁 **Output Location**: All converted images are automatically saved in the `PDF_Images` folder
- 🗂️ **Local PDF Detection**: At startup, the script scans the script folder for `*.pdf` and lets you pick one. To use a PDF in a different folder, either move it next to `Convert.bat` or enter its full path when prompted.
- 🔄 **Automatic Folder Creation**: The script creates the `PDF_Images` directory if it doesn't exist
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

## License

This project is open source and available under the [MIT License](LICENSE).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues:

1. Check the troubleshooting section above
2. Verify ImageMagick installation
3. Ensure PDF file is not corrupted or password-protected

## Version History

- **v1.1** - Added local PDF detection and selection at startup; selected file is used automatically for conversion (no auto-open)
- **v1.0** - Initial release with basic PDF to image conversion functionality
