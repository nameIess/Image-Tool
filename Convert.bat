@echo off
title PDF to Image Converter
setlocal enabledelayedexpansion

:: PDF to Image Converter using ImageMagick
:: =======================================

:: -----------------------------------------------------------------
:: Centralized Defaults (edit here to change script default values)
:: -----------------------------------------------------------------
set "DEF_OUTPUT_FORMAT=png"
set "DEF_DENSITY=180"
set "DEF_QUALITY=90"
set "DEF_PREFIX=Page-"
:: -----------------------------------------------------------------

echo.
echo =========================================================
echo *                 PDF to Image Converter                *        
echo *                   Using ImageMagick                   *        
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
echo No PDF files found in script directory.
echo.
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

:: Sanitize: remove surrounding quotes if user pasted a quoted path
set "input_pdf=%input_pdf:"=%"
:: Normalize to full path (handles relative paths and spaces)
for %%I in ("%input_pdf%") do set "input_pdf=%%~fI"

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
set /p output_format="Enter output image format (default: %DEF_OUTPUT_FORMAT%): "
if "%output_format%"=="" set "output_format=%DEF_OUTPUT_FORMAT%"

:: Validate format (add dot if not present)
if not "!output_format:~0,1!"=="." set output_format=.!output_format!

:: Get density value
:density_input
echo.
set /p density="Enter density value (default: %DEF_DENSITY%): "
if "%density%"=="" set "density=%DEF_DENSITY%"

:: Validate density is numeric
echo %density%| findstr /r "^[0-9][0-9]*$" >nul
if %errorlevel% neq 0 (
    echo [ERROR] Please enter a valid numeric density value!
    goto density_input
)

:: Get quality value
:quality_input
echo.
set /p quality="Enter quality value 1-100 (default: %DEF_QUALITY%): "
if "%quality%"=="" set "quality=%DEF_QUALITY%"

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
set /p prefix="Enter output filename prefix (default: %DEF_PREFIX%): "
if "%prefix%"=="" set "prefix=%DEF_PREFIX%"

:: Derive output folder name from input PDF filename (without extension)
for %%I in ("%input_pdf%") do set "pdf_base=%%~nI"
set "output_folder=!pdf_base!_image"

:: Check if output folder exists, create if not
if not exist "!output_folder!" (
    echo.
    echo [INFO] Creating output folder: "!output_folder!" ...
    mkdir "!output_folder!" 2>nul
    if exist "!output_folder!" (
        echo [SUCCESS] Folder created: "!output_folder!"
    ) else (
        echo [ERROR] Failed to create output folder: "!output_folder!"!
        echo [INFO] This might be due to permissions or disk space issues.
        echo [INFO] Falling back to current directory for output files.
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
    echo [INFO] Using existing output folder: "!output_folder!" ...
)

:: Set output path based on folder availability
if not defined use_subfolder set use_subfolder=true
if "!use_subfolder!"=="true" (
    set "output_path=!output_folder!\"
) else (
    set "output_path="
)

:: Compute absolute output directory path for opening in Explorer
set "output_dir=%CD%"
if "!use_subfolder!"=="true" set "output_dir=%CD%\!output_folder!"

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
set /p confirm="Proceed with conversion? (y/N): "
if /i not "%confirm%"=="Y" if /i not "%confirm%"=="YES" (
    echo .......Conversion cancelled.......
    goto end
)

:: Execute the conversion
echo.
echo ..*..Converting..*.. 
echo.
echo Running ImageMagick conversion...
echo.

magick -density %density% "%input_pdf%" -quality %quality% "!output_path!%prefix%%%d%output_format%"

:: Check if conversion was successful (avoid complex parentheses to prevent parser errors)
if errorlevel 1 goto conversion_failed

echo.
echo [SUCCESS] PDF conversion completed successfully!

if /i "!use_subfolder!"=="true" goto success_subfolder
goto success_cwd

:success_subfolder
echo Output files: !output_folder!\%prefix%*%output_format%
:: Count and display created files
for /f %%i in ('dir /b "!output_folder!\%prefix%*%output_format%" 2^>nul ^| find /c /v ""') do set file_count=%%i
if defined file_count echo Created !file_count! image file(s) in !output_folder! folder.
:: Automatically open the output folder
echo.
echo Opening output folder...
explorer "!output_dir!" >nul 2>&1
goto post_success

:success_cwd
echo Output files: %prefix%*%output_format%
:: Count and display created files
for /f %%i in ('dir /b "%prefix%*%output_format%" 2^>nul ^| find /c /v ""') do set file_count=%%i
if defined file_count echo Created !file_count! image file(s) in current directory.
:: Automatically open the output folder
echo.
echo Opening output folder...
explorer "!output_dir!" >nul 2>&1
goto post_success

:conversion_failed
echo.
echo [ERROR] Conversion failed! Please check the error messages above.
echo.
echo Common issues:
echo - PDF file might be corrupted or password protected
echo - Insufficient disk space
echo - Invalid parameters
echo - ImageMagick configuration issues

:post_success

echo.

:end
echo.
echo ==========================================================
echo Thank you for using PDF to Image Converter!
echo ==========================================================
echo.
timeout /t 5 >nul
exit