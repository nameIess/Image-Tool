@echo off
title Image Converter Suite
setlocal enabledelayedexpansion

:: =========================================================
:: Configuration
:: =========================================================
set "HEADER_BORDER========================================================="
set "HDR_MAIN_LINE1=*              Image Converter Suite                    *"
set "HDR_MAIN_LINE2=*                Using ImageMagick                      *"
set "HDR_PDF_LINE1=*                 PDF to Image Converter                *"
set "HDR_PDF_LINE2=*                   Using ImageMagick                   *"
set "HDR_FORMAT_LINE1=*              Convert Image Format                     *"
set "HDR_COMP_LINE1=*          Compress Image/PDF to Target Size            *"

set "DEF_OUTPUT_FORMAT=png"
set "DEF_DENSITY=180"
set "DEF_QUALITY=90"
set "DEF_PREFIX=Page-"
set "DEF_CONVERT_FORMAT=png"
set "DEF_COMPRESS_PERCENT=75"

set "AVAILABLE_FORMATS=png, jpg, jpeg, bmp, tiff, gif"

set "PATTERN_PDF=*.pdf"
set "PATTERN_IMAGES=*.jpg *.jpeg *.png *.bmp *.gif *.tiff *.tif *.webp *.avif"
set "PATTERN_COMPRESS=*.jpg *.jpeg *.png *.bmp *.gif *.tiff *.tif *.webp *.avif *.pdf"

set "SCRIPT_DIR=%~dp0"
set "MAGICK_AVAILABLE=0"

goto main_menu

:main_menu
call :print_header "%HDR_MAIN_LINE1%" "%HDR_MAIN_LINE2%"
echo Select an operation:
echo.
echo   1. PDF to Image Converter
echo   2. Convert Image Format
echo   3. Compress Image/PDF to Target Size
echo   0. Exit
echo.
set "menu_choice="
set /p menu_choice="Enter your choice (default: 1): "
if "%menu_choice%"=="" set "menu_choice=1"

if "%menu_choice%"=="0" goto final_exit
if "%menu_choice%"=="1" (
    call :feature_pdf_to_image
    goto main_menu
)
if "%menu_choice%"=="2" (
    call :feature_convert_format
    goto main_menu
)
if "%menu_choice%"=="3" (
    call :feature_compress_file
    goto main_menu
)

echo [ERROR] Invalid choice!
timeout /t 2 >nul
goto main_menu

:feature_pdf_to_image
call :print_header "%HDR_PDF_LINE1%" "%HDR_PDF_LINE2%"
call :ensure_magick
if errorlevel 1 exit /b 0

echo -- Local PDF detection --
echo.
set "selected_pdf="
call :select_file "PDF" selected_pdf %PATTERN_PDF%
call :resolve_input_file "Enter PDF filepath (with extension): " "%selected_pdf%" input_pdf

echo.
echo Available formats: %AVAILABLE_FORMATS%
echo.
set "output_format="
set /p output_format="Enter output image format (default: %DEF_OUTPUT_FORMAT%): "
if "%output_format%"=="" set "output_format=%DEF_OUTPUT_FORMAT%"
if not "%output_format:~0,1%"=="." set "output_format=.%output_format%"

call :prompt_numeric "Enter density value (default: %DEF_DENSITY%): " "%DEF_DENSITY%" "" "" density
call :prompt_numeric "Enter quality value 1-100 (default: %DEF_QUALITY%): " "%DEF_QUALITY%" "1" "100" quality

echo.
set "prefix="
set /p prefix="Enter output filename prefix (default: %DEF_PREFIX%): "
if "%prefix%"=="" set "prefix=%DEF_PREFIX%"

for %%I in ("%input_pdf%") do set "pdf_base=%%~nI"
set "output_folder=%pdf_base%_images"
set "use_subfolder=true"
call :ensure_directory "%output_folder%"
if errorlevel 1 (
    echo.
    echo [ERROR] Failed to create output folder: "%output_folder%"!
    echo [INFO] Images will be saved in current directory instead.
    call :prompt_yes_no "Continue with conversion? (Y/n): " "N" continue_anyway
    if /i "%continue_anyway%"=="N" (
        echo Conversion cancelled.
        pause
        exit /b 0
    )
    set "use_subfolder=false"
)

set "output_dir=%CD%"
set "output_path="
if /i "%use_subfolder%"=="true" (
    set "output_path=%output_folder%\"
    for %%I in ("%CD%\%output_folder%") do set "output_dir=%%~fI"
)

echo.
echo ..*..Conversion Summary..*..
echo.
echo Input PDF:      %input_pdf%
echo Output format:  %output_format%
echo Density:        %density%
echo Quality:        %quality%
echo Prefix:         %prefix%
echo Output folder:  !output_path!
echo Command:        magick -density %density% "%input_pdf%" -quality %quality% "!output_path!!prefix!%%d!output_format!"
echo.

call :prompt_yes_no "Proceed with conversion? (Y/n): " "Y" confirm_pdf
if /i "%confirm_pdf%"=="N" (
    echo .......Conversion cancelled.......
    exit /b 0
)

echo.
echo ..*..Converting..*..
echo.
magick -density %density% "%input_pdf%" -quality %quality% "!output_path!!prefix!%%d!output_format!"
if errorlevel 1 (
    echo.
    echo [ERROR] Conversion failed! Please check the error messages above.
    echo.
    echo Common issues:
    echo - PDF file might be corrupted or password protected
    echo - Insufficient disk space
    echo - Invalid parameters
    echo - ImageMagick configuration issues
    goto pdf_post
)

echo.
echo [SUCCESS] PDF conversion completed successfully!

set "file_count="
if /i "%use_subfolder%"=="true" (
    echo Output files: !output_folder!\!prefix!*!output_format!
    call :count_created_files "!output_dir!" "!prefix!*!output_format!" file_count
) else (
    echo Output files: !prefix!*!output_format!
    call :count_created_files "!output_dir!" "!prefix!*!output_format!" file_count
)
if defined file_count echo Created %file_count% image file(s).

echo.
echo Opening output folder...
explorer "!output_dir!" >nul 2>&1

:pdf_post
echo.
call :prompt_yes_no "Return to main menu? (Y/n): " "Y" return_menu_pdf
if /i "%return_menu_pdf%"=="Y" exit /b 0
goto final_exit

:feature_convert_format
call :print_header "%HDR_FORMAT_LINE1%"
call :ensure_magick
if errorlevel 1 exit /b 0

echo -- Local Image Detection --
echo.
set "selected_image="
call :select_file "image" selected_image %PATTERN_IMAGES%
call :resolve_input_file "Enter image filepath (with extension): " "%selected_image%" input_image

:format_prompt
echo.
echo Available format options:
echo   1) jpeg
echo   2) jpg
echo   3) custom (avif, webp, tiff, bmp, etc.)
echo.
set "output_format_conv="
set /p format_choice="Select output format (default: %DEF_CONVERT_FORMAT%): "
if "%format_choice%"=="" (
    set "output_format_conv=%DEF_CONVERT_FORMAT%"
) else (
    if "%format_choice%"=="1" (
        set "output_format_conv=jpeg"
    ) else (
        if "%format_choice%"=="2" (
            set "output_format_conv=jpg"
        ) else (
            if "%format_choice%"=="3" (
                set /p output_format_conv="Enter custom format: "
                if "%output_format_conv%"=="" (
                    echo [ERROR] Custom format cannot be empty!
                    goto format_prompt
                )
            ) else (
                set "output_format_conv=%format_choice%"
            )
        )
    )
)
if not defined output_format_conv (
    echo [ERROR] Output format cannot be empty!
    goto format_prompt
)
if "!output_format_conv:~0,1!"=="." set "output_format_conv=!output_format_conv:~1!"

for %%I in ("%input_image%") do (
    set "img_dir=%%~dpI"
    set "img_base=%%~nI"
)
set "output_image=!img_dir!!img_base!_conv.!output_format_conv!"

echo.
echo ..*..Conversion Summary..*..
echo.
echo Input image:    %input_image%
echo Output format:  !output_format_conv!
echo Output file:    !output_image!
echo.
call :prompt_yes_no "Proceed with conversion? (Y/n): " "Y" confirm_conv
if /i "%confirm_conv%"=="N" exit /b 0

echo.
echo ..*..Converting..*..
echo.
magick "%input_image%" "!output_image!"
if errorlevel 1 (
    echo.
    echo [ERROR] Conversion failed!
    pause
    exit /b 0
)

echo.
echo [SUCCESS] Image converted successfully!
echo Output: !output_image!
call :display_size "!output_image!" "File size"

echo.
call :prompt_yes_no "Return to main menu? (Y/n): " "Y" return_menu_format
if /i "%return_menu_format%"=="Y" exit /b 0
goto final_exit

:feature_compress_file
call :print_header "%HDR_COMP_LINE1%"
call :ensure_magick
if errorlevel 1 exit /b 0

echo -- Local File Detection --
echo.
set "selected_file="
call :select_file "file" selected_file %PATTERN_COMPRESS%
call :resolve_input_file "Enter filepath (with extension): " "%selected_file%" input_file

for %%A in ("%input_file%") do set "orig_size=%%~zA"
set /a orig_size_kb=orig_size/1024
if %orig_size_kb% lss 1 set orig_size_kb=1
echo Original size: %orig_size% bytes (~%orig_size_kb% KB)

:: Ask user for compression mode
echo.
echo Select compression mode:
echo   1. By percentage (default)
echo   2. To fixed file size
set "compress_mode="
set /p compress_mode="Enter choice (1 or 2, default: 1): "
if "%compress_mode%"=="" set "compress_mode=1"

if "%compress_mode%"=="2" goto compress_fixed_size

:: --- Existing percentage-based compression ---
echo.
echo Enter target file size percentage (1-100).
echo Example: 50 means the output file will be 50%% of the original size.
call :prompt_numeric "Target Percentage (default: %DEF_COMPRESS_PERCENT%): " "%DEF_COMPRESS_PERCENT%" "1" "100" target_percent

set /a target_size_kb = orig_size_kb * target_percent / 100
if %target_size_kb% lss 1 set target_size_kb=1
set "target_size=%target_size_kb%KB"
set /a target_size_mb = target_size_kb / 1024
echo Target size calculated: %target_size_kb%KB (~%target_size_mb% MB)

for %%I in ("%input_file%") do (
    set "file_dir=%%~dpI"
    set "file_base=%%~nI"
    set "file_ext=%%~xI"
)
set "output_ext=.jpg"
if /i "%file_ext%"==".pdf" set "output_ext=.pdf"
set "output_file=!file_dir!!file_base!_comp!output_ext!"

echo.
echo ..*..Compression Summary..*..
echo.
echo Input file:     %input_file%
echo Target Size:    %target_size_kb%KB (~%target_size_mb% MB)
echo Output file:    !output_file!
echo.
call :prompt_yes_no "Proceed with compression? (Y/n): " "Y" confirm_comp
if /i "%confirm_comp%"=="N" exit /b 0

echo.
echo ..*..Compressing..*..
echo.
if /i "%file_ext%"==".pdf" (
    echo [INFO] Processing PDF. This may rasterize the content to achieve the target size.
)

magick "%input_file%" -define jpeg:extent=%target_size% "!output_file!"
if errorlevel 1 (
    echo.
    echo [ERROR] Compression failed!
    pause
    exit /b 0
)

echo.
echo [SUCCESS] File compressed successfully!
echo Output: !output_file!
call :display_size "!output_file!" "New file size"

echo.
call :prompt_yes_no "Return to main menu? (Y/n): " "Y" return_menu_comp
if /i "%return_menu_comp%"=="Y" exit /b 0
goto final_exit

:compress_fixed_size
:: --- New fixed size compression ---
for %%I in ("%input_file%") do (
    set "file_dir=%%~dpI"
    set "file_base=%%~nI"
    set "file_ext=%%~xI"
)
set "output_ext=.jpg"
if /i "%file_ext%"==".pdf" (
    echo [ERROR] Fixed size compression for PDF is not supported yet.
    pause
    exit /b 0
)
set "output_file=!file_dir!!file_base!_comp!output_ext!"

:: Prompt for value and unit
set "fixed_value="
set "fixed_unit="
echo.
set /p fixed_value="Enter target file size value (e.g., 20): "
if "%fixed_value%"=="" (
    echo [ERROR] Value is required!
    pause
    exit /b 0
)
set /a fixed_value_num=%fixed_value% 2>nul
if "%fixed_value_num%"=="0" (
    echo [ERROR] Please enter a valid positive number!
    pause
    exit /b 0
)
echo Enter unit (B, KB, MB):
set /p fixed_unit="Unit: "
if /i "%fixed_unit%"=="KB" set /a target_bytes=%fixed_value%*1024
if /i "%fixed_unit%"=="MB" set /a target_bytes=%fixed_value%*1024*1024
if /i "%fixed_unit%"=="B" set /a target_bytes=%fixed_value%
if not defined target_bytes (
    echo [ERROR] Invalid unit! Please enter B, KB, or MB.
    pause
    exit /b 0
)
if %target_bytes% lss 1024 set target_bytes=1024

:: Check if target is less than original
if %target_bytes% geq %orig_size% (
    echo [ERROR] Target size must be less than original file size!
    pause
    exit /b 0
)

echo.
echo ..*..Compression Summary..*..
echo.
echo Input file:     %input_file%
echo Target Size:    %fixed_value%%fixed_unit% (%target_bytes% bytes)
echo Output file:    !output_file!
echo.
call :prompt_yes_no "Proceed with compression? (Y/n): " "Y" confirm_comp_fixed
if /i "%confirm_comp_fixed%"=="N" exit /b 0

echo.
echo ..*..Compressing to fixed size..*..
echo.

:: Iteratively compress by reducing quality
setlocal enabledelayedexpansion
set "quality=90"
set "min_quality=10"
set "step=5"
set "done=0"
if exist "!output_file!" del /f /q "!output_file!" >nul 2>&1

:compress_loop
magick "%input_file%" -quality !quality! "!output_file!"
for %%A in ("!output_file!") do set "out_size=%%~zA"
if not defined out_size set "out_size=0"
if !out_size! lss %target_bytes% (
    set "done=1"
    goto compress_done
)
if !quality! leq !min_quality! (
    echo [WARN] Minimum quality reached. Could not reach target size.
    set "done=2"
    goto compress_done
)
set /a quality-=step
goto compress_loop

:compress_done
if !done! equ 1 (
    echo.
    echo [SUCCESS] File compressed to !out_size! bytes (target: %target_bytes% bytes)
    echo Output: !output_file!
    call :display_size "!output_file!" "New file size"
) else (
    echo.
    echo [INFO] Final file size: !out_size! bytes (target: %target_bytes% bytes)
    echo Output: !output_file!
    call :display_size "!output_file!" "New file size"
)
endlocal

echo.
call :prompt_yes_no "Return to main menu? (Y/n): " "Y" return_menu_comp
if /i "%return_menu_comp%"=="Y" exit /b 0
goto final_exit

:ensure_magick
if "%MAGICK_AVAILABLE%"=="1" exit /b 0
where magick >nul 2>nul
if errorlevel 1 (
    echo [ERROR] ImageMagick is not installed or not in PATH!
    echo Please install ImageMagick from: https://imagemagick.org/script/download.php
    echo.
    pause
    exit /b 1
)
set "MAGICK_AVAILABLE=1"
exit /b 0

:select_file
setlocal enabledelayedexpansion
set "desc=%~1"
set "result="
set "found=0"
set "script_dir=!SCRIPT_DIR!"
shift
set "out_var=%~1"
shift
if not defined script_dir set "script_dir=%CD%"
pushd "!script_dir!" >nul 2>&1
:select_file_loop
if "%~1"=="" goto select_file_done
for %%F in (%1) do (
    if exist "%%~fF" (
        set /a found+=1
        set "file!found!=%%~fF"
    )
)
shift
goto select_file_loop
:select_file_done
popd >nul 2>&1
if !found! gtr 0 (
    echo Found !found! %desc% file^(s^) in script directory:
    for /l %%i in (1,1,!found!) do (
        for %%G in ("!file%%i!") do echo   %%i. %%~nxG
    )
    echo.
    set "sel="
    set /p sel="Select a %desc% (1-!found!) or press Enter to skip: "
    if defined sel (
        echo !sel!| findstr /r "^[0-9][0-9]*$" >nul
        if not errorlevel 1 (
            if !sel! geq 1 if !sel! leq !found! (
                for %%i in (!sel!) do set "result=!file%%i!"
                for %%G in ("!result!") do echo Selected: %%~nxG
            )
        )
    )
)
if defined result (
    endlocal & set "%out_var%=%result%" & exit /b 0
)
endlocal & set "%out_var%=" & exit /b 0

:resolve_input_file
setlocal enabledelayedexpansion
set "prompt=%~1"
set "preselected=%~2"
if defined preselected (
    if exist "!preselected!" (
        for %%A in ("!preselected!") do set "pre_disp=%%~nxA"
        if defined pre_disp (
            echo Using selected file: !pre_disp!
        ) else (
            echo Using selected file: !preselected!
        )
        endlocal & set "%~3=%preselected%" & exit /b 0
    ) else (
        echo [WARN] Preselected file '!preselected!' not found. Please choose a file.
    )
)
:resolve_loop
set "input="
set /p input="%prompt%"
if not defined input (
    echo [ERROR] Please enter a valid filename!
    goto resolve_loop
)
set "input=!input:"=!"
for %%I in ("!input!") do set "input=%%~fI"
if not exist "!input!" (
    echo [ERROR] File '!input!' does not exist!
    goto resolve_loop
)
endlocal & set "%~3=%input%" & exit /b 0

:prompt_numeric
setlocal enabledelayedexpansion
set "prompt=%~1"
set "default=%~2"
set "min=%~3"
set "max=%~4"
:prompt_numeric_loop
set "value="
set /p value="%prompt%"
if not defined value set "value=%default%"
if not defined value (
    echo [ERROR] Value is required!
    goto prompt_numeric_loop
)
set "value=!value:"=!"
echo !value!| findstr /r "^[0-9][0-9]*$" >nul
if errorlevel 1 (
    echo [ERROR] Please enter a numeric value!
    goto prompt_numeric_loop
)
if defined min (
    if !value! lss !min! (
        echo [ERROR] Value must be at least !min!!
        goto prompt_numeric_loop
    )
)
if defined max (
    if !value! gtr !max! (
        echo [ERROR] Value must be !max! or less!
        goto prompt_numeric_loop
    )
)
endlocal & set "%~5=%value%" & exit /b 0

:prompt_yes_no
setlocal enabledelayedexpansion
set "prompt=%~1"
set "default=%~2"
if not defined default set "default=Y"
:prompt_yes_no_loop
set "answer="
set /p answer="%prompt%"
if not defined answer set "answer=!default!"
set "first=!answer:~0,1!"
if /i "!first!"=="Y" (
    endlocal & set "%~3=Y" & exit /b 0
)
if /i "!first!"=="N" (
    endlocal & set "%~3=N" & exit /b 0
)
echo [ERROR] Please enter Y or N.
goto prompt_yes_no_loop

:display_size
if not exist "%~1" exit /b 0
setlocal enabledelayedexpansion
set "file=%~1"
set "label=%~2"
for %%A in ("!file!") do set "size=%%~zA"
set /a sizeKB=size/1024
if !sizeKB! lss 1 set sizeKB=1
set /a sizeMB=sizeKB/1024
if not defined label set "label=File size"
echo !label!: !size! bytes (!sizeKB! KB / !sizeMB! MB)
endlocal & exit /b 0

:ensure_directory
set "target=%~1"
if "%target%"=="" exit /b 1
if exist "%target%" exit /b 0
mkdir "%target%" 2>nul
if exist "%target%" exit /b 0
exit /b 1

:print_header
cls
echo.
echo %HEADER_BORDER%
if not "%~1"=="" echo %~1
if not "%~2"=="" echo %~2
if not "%~3"=="" echo %~3
echo %HEADER_BORDER%
echo.
exit /b 0

:count_created_files
setlocal enabledelayedexpansion
set "dir=%~1"
set "pattern=%~2"
set "count=0"
if not exist "!dir!" (
    endlocal & set "%~3=0" & exit /b 0
)
pushd "!dir!" >nul 2>&1
for /f %%i in ('dir /b "!pattern!" 2^>nul ^| find /c /v ""') do set "count=%%i"
popd >nul 2>&1
endlocal & set "%~3=%count%" & exit /b 0

:final_exit
echo.
echo %HEADER_BORDER%
echo Thank you for using Image Converter Suite!
echo %HEADER_BORDER%
echo.
timeout /t 3 >nul
endlocal
exit