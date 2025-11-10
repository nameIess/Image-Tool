@echo off
setlocal enabledelayedexpansion

:: PDF to Image Converter using ImageMagick
:: =======================================

echo.
echo =========================================================
echo *                    PDF to Image Converter             *        
echo *                      Using ImageMagick                *        
echo =========================================================
echo.

:: Check if ImageMagick is available
where magick >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] ImageMagick is not installed or not in PATH!
    echo Please install ImageMagick from: https://imagemagick.org/script/download.php
    echo.
    pause
    exit /b 1
)

:: Check for PDF files in the script directory and offer to open one
echo.
echo -- Local PDF detection --
set "script_dir=%~dp0"
set "found=0"
pushd "%script_dir%" >nul 2>&1
for %%F in (*.pdf) do (
    set /a found+=1
    set "pdf!found!=%%F"
)
popd >nul 2>&1

if %found% gtr 0 goto show_pdfs
goto after_pdf_list

:show_pdfs
echo Found %found% PDF(s) in script directory:
for /l %%i in (1,1,%found%) do echo %%i. !pdf%%i!
echo.
set /p sel="Select a PDF to use (1-%found%) or press Enter to skip: "
if "%sel%"=="" goto after_pdf_list
echo %sel%| findstr /r "^[0-9][0-9]*$" >nul || goto after_pdf_list
if %sel% lss 1 goto after_pdf_list
if %sel% gtr %found% goto after_pdf_list
call set "file_to_open=%%pdf%sel%%%"
set "input_pdf=%script_dir%!file_to_open!"
echo Selected: %input_pdf%

:after_pdf_list

:input_section
echo -- Input Parameters --
echo.

:: Get input PDF filename
:input_pdf
:: If a PDF was selected above, reuse it and skip prompt
if defined input_pdf (
    echo Using selected PDF: %input_pdf%
    goto validate_input_file
)
set /p input_pdf="Enter PDF filename (with extension): "
if "%input_pdf%"=="" (
    echo [ERROR] Please enter a valid filename!
    goto input_pdf
)

:: Check if file exists
:validate_input_file
if not exist "%input_pdf%" (
    echo [ERROR] File '%input_pdf%' does not exist!
    echo.
    goto input_pdf
)

:: Get output image format
:output_format
echo.
echo Available formats: png, jpg, jpeg, bmp, tiff, gif
set /p output_format="Enter output image format (default: png): "
if "%output_format%"=="" set output_format=png

:: Validate format (add dot if not present)
if not "!output_format:~0,1!"=="." set output_format=.!output_format!

:: Get density value
:density_input
echo.
set /p density="Enter density value (default: 180): "
if "%density%"=="" set density=180

:: Validate density is numeric
echo %density%| findstr /r "^[0-9][0-9]*$" >nul
if %errorlevel% neq 0 (
    echo [ERROR] Please enter a valid numeric density value!
    goto density_input
)

:: Get quality value
:quality_input
echo.
set /p quality="Enter quality value 1-100 (default: 90): "
if "%quality%"=="" set quality=90

:: Validate quality is numeric and in range
echo %quality%| findstr /r "^[0-9][0-9]*$" >nul
if %errorlevel% neq 0 (
    echo [ERROR] Please enter a valid numeric quality value!
    goto quality_input
)

if %quality% lss 1 (
    echo [ERROR] Quality must be between 1 and 100!
    goto quality_input
)
if %quality% gtr 100 (
    echo [ERROR] Quality must be between 1 and 100!
    goto quality_input
)

:: Get output filename prefix
echo.
set /p prefix="Enter output filename prefix (default: Page-): "
if "%prefix%"=="" set prefix=Page-

:: Check if PDF_Images folder exists, create if not
if not exist "PDF_Images" (
    echo.
    echo [INFO] Creating PDF_Images folder...
    mkdir "PDF_Images" 2>nul
    if exist "PDF_Images" (
        echo [SUCCESS] PDF_Images folder created successfully!
    ) else (
        echo [ERROR] Failed to create PDF_Images folder!
        echo [INFO] This might be due to permissions or disk space issues.
        echo [INFO] Trying to create folder in current directory...
        echo Current directory: %CD%
        echo.
        set /p continue_anyway="Continue with conversion anyway? Images will be saved in current directory (Y/N): "
        if /i not "!continue_anyway!"=="Y" if /i not "!continue_anyway!"=="YES" (
            echo Conversion cancelled.
            pause
            exit /b 1
        )
        set use_subfolder=false
    )
) else (
    echo.
    echo [INFO] Using existing PDF_Images folder...
)

:: Set output path based on folder availability
if not defined use_subfolder set use_subfolder=true
if "!use_subfolder!"=="true" (
    set "output_path=PDF_Images\"
) else (
    set "output_path="
)

:: Display summary
echo.
echo ..*..Conversion Summary..*..
echo.
echo Input PDF:      %input_pdf%
echo Output format:  %output_format%
echo Density:        %density%
echo Quality:        %quality%
echo Prefix:         %prefix%
echo Output folder:  !output_path!
echo Command:        magick -density %density% "%input_pdf%" -quality %quality% "!output_path!%prefix%%%d%output_format%"
echo.

:: Confirm before proceeding
set /p confirm="Proceed with conversion? (Y/N): "
if /i not "%confirm%"=="Y" if /i not "%confirm%"=="YES" (
    echo Conversion cancelled.
    goto end
)

:: Execute the conversion
echo.
echo ..*..Converting..*.. 
echo.
echo Running ImageMagick conversion...
echo.

magick -density %density% "%input_pdf%" -quality %quality% "!output_path!%prefix%%%d%output_format%"

:: Check if conversion was successful
if %errorlevel% equ 0 (
    echo.
    echo [SUCCESS] PDF conversion completed successfully!
    if "!use_subfolder!"=="true" (
        echo Output files: PDF_Images\%prefix%*%output_format%
        :: Count and display created files
        for /f %%i in ('dir /b "PDF_Images\%prefix%*%output_format%" 2^>nul ^| find /c /v ""') do set file_count=%%i
        if defined file_count (
            echo Created !file_count! image file(s) in PDF_Images folder.
        )
        :: Ask if user wants to open the output folder
        echo.
        set /p open_folder="Open PDF_Images folder? (Y/N): "
        if /i "!open_folder!"=="Y" if /i "!open_folder!"=="YES" (
            start PDF_Images
        )
    ) else (
        echo Output files: %prefix%*%output_format%
        :: Count and display created files
        for /f %%i in ('dir /b "%prefix%*%output_format%" 2^>nul ^| find /c /v ""') do set file_count=%%i
        if defined file_count (
            echo Created !file_count! image file(s) in current directory.
        )
        :: Ask if user wants to open the output folder
        echo.
        set /p open_folder="Open current folder? (Y/N): "
        if /i "!open_folder!"=="Y" if /i "!open_folder!"=="YES" (
            start .
        )
    )
) else (
    echo.
    echo [ERROR] Conversion failed! Please check the error messages above.
    echo.
    echo Common issues:
    echo - PDF file might be corrupted or password protected
    echo - Insufficient disk space
    echo - Invalid parameters
    echo - ImageMagick configuration issues
)

echo.

:: Ask if user wants to convert another file
:another_conversion
set /p another="Convert another PDF file? (Y/N): "
if /i "%another%"=="Y" if /i "%another%"=="YES" (
    echo.
    echo ==========================================================
    echo.
    goto input_section
)

:end
echo.
echo....................................................
echo Thank you for using PDF to Image Converter!
echo....................................................
echo.
pause