@echo off
title Image Converter Suite
setlocal enabledelayedexpansion

:: Image Converter Suite using ImageMagick
:: ========================================

:: -----------------------------------------------------------------
:: Centralized Defaults (edit here to change script default values)
:: -----------------------------------------------------------------
set "DEF_OUTPUT_FORMAT=png"
set "DEF_DENSITY=180"
set "DEF_QUALITY=90"
set "DEF_PREFIX=Page-"
set "DEF_CONVERT_FORMAT=png"
:: -----------------------------------------------------------------

:main_menu
cls
echo.
echo =========================================================
echo *              Image Converter Suite                    *        
echo *                Using ImageMagick                      *        
echo =========================================================
echo.
echo Select an operation:
echo.
echo   1. PDF to Image Converter
echo   2. Convert Image Format
echo   0. Exit
echo.
set /p menu_choice="Enter your choice (default: 1): "
if "%menu_choice%"=="" set "menu_choice=1"

if "%menu_choice%"=="0" goto exit_script
if "%menu_choice%"=="1" goto pdf_converter
if "%menu_choice%"=="2" goto image_format_converter


echo [ERROR] Invalid choice!
timeout /t 2 >nul
goto main_menu

:exit_script
echo.
echo Exiting... Goodbye!
timeout /t 1 >nul
exit /b 0

:pdf_converter
cls
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
echo.
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
echo.
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
echo.
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
set /p return_menu="Return to main menu? (Y/n): "
if /i "%return_menu%"=="n" goto final_exit
if /i "%return_menu%"=="no" goto final_exit
goto main_menu

:final_exit
echo.
echo ==========================================================
echo Thank you for using Image Converter Suite!
echo ==========================================================
echo.
timeout /t 3 >nul
exit

:: =========================================================
:: FEATURE 2: Convert Image Format
:: =========================================================
:image_format_converter
cls
echo.
echo =========================================================
echo *              Convert Image Format                     *
echo =========================================================
echo.

:: Check if ImageMagick is available
where magick >nul 2>nul
if %errorlevel% neq 0 (
    echo [ERROR] ImageMagick is not installed or not in PATH!
    echo Please install ImageMagick from: https://imagemagick.org/script/download.php
    echo.
    pause
    goto main_menu
)

:: Check for image files in the script directory
echo.
echo -- Local Image Detection --
set "script_dir=%~dp0"
set "found=0"
pushd "%script_dir%" >nul 2>&1
for %%F in (*.jpg *.jpeg *.png *.bmp *.gif *.tiff *.tif *.webp *.avif) do (
    set /a found+=1
    set "img!found!=%%F"
)
popd >nul 2>&1

if %found% gtr 0 goto show_images
echo No image files found in script directory.
echo.
goto after_image_list

:show_images
echo Found %found% image(s) in script directory:
for /l %%i in (1,1,%found%) do echo %%i. !img%%i!
echo.
set /p img_sel="Select an image (1-%found%) or press Enter to specify path: "
if "%img_sel%"=="" goto after_image_list
echo %img_sel%| findstr /r "^[0-9][0-9]*$" >nul || goto after_image_list
if %img_sel% lss 1 goto after_image_list
if %img_sel% gtr %found% goto after_image_list
call set "file_to_convert=%%img%img_sel%%%"
set "input_image=%script_dir%!file_to_convert!"
echo Selected: %input_image%

:after_image_list

:: Get input image filename
:input_image_convert
if defined input_image (
    echo Using selected image: %input_image%
    goto validate_input_image
)
set /p input_image="Enter image filename (with extension): "
if "%input_image%"=="" (
    echo [ERROR] Please enter a valid filename!
    goto input_image_convert
)

:: Sanitize: remove surrounding quotes
set "input_image=%input_image:"=%"
:: Normalize to full path
for %%I in ("%input_image%") do set "input_image=%%~fI"

:validate_input_image
if not exist "%input_image%" (
    echo [ERROR] File '%input_image%' does not exist!
    echo.
    set "input_image="
    goto input_image_convert
)

:: Get output format
echo.
echo Available format options:
echo   1) jpeg
echo   2) jpg
echo   3) custom (avif, webp, tiff, bmp, etc.)
echo.
set /p format_choice="Select output format (default: png): "
if "%format_choice%"=="" set "output_format_conv=png"
if "%format_choice%"=="1" set "output_format_conv=jpeg"
if "%format_choice%"=="2" set "output_format_conv=jpg"
if "%format_choice%"=="3" (
    set /p custom_format="Enter custom format: "
    if "!custom_format!"=="" (
        echo [ERROR] Custom format cannot be empty!
        goto input_image_convert
    )
    set "output_format_conv=!custom_format!"
)
if not defined output_format_conv set "output_format_conv=png"

:: Remove dot if present
if "!output_format_conv:~0,1!"=="." set "output_format_conv=!output_format_conv:~1!"

:: Generate output filename
for %%I in ("%input_image%") do (
    set "img_dir=%%~dpI"
    set "img_base=%%~nI"
    set "img_ext=%%~xI"
)
set "output_image=!img_dir!!img_base!_conv.!output_format_conv!"

:: Display conversion summary
echo.
echo ..*..Conversion Summary..*..
echo.
echo Input image:    %input_image%
echo Output format:  !output_format_conv!
echo Output file:    !output_image!
echo.
set /p confirm_conv="Proceed with conversion? (Y/n): "
if /i "%confirm_conv%"=="n" goto end
if /i "%confirm_conv%"=="no" goto end

:: Execute conversion
echo.
echo ..*..Converting..*..
echo.
magick "%input_image%" "!output_image!"

if errorlevel 1 (
    echo.
    echo [ERROR] Conversion failed!
    pause
    goto end
)

echo.
echo [SUCCESS] Image converted successfully!
echo Output: !output_image!

:: Display file size of converted image
for %%A in ("!output_image!") do set "converted_size=%%~zA"
set /a converted_size_kb=!converted_size! / 1024
set /a converted_size_mb=!converted_size_kb! / 1024
echo File size: !converted_size! bytes (!converted_size_kb! KB / !converted_size_mb! MB)
goto end

:: (Compression features removed)