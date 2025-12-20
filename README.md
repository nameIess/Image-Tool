# 🎨 Icons8 Collector

Download icons from your Icons8 collections with ease. Supports PNG and ICO formats with automatic conversion.

## 📥 Set-Up

**[⬇️ Download ](https://github.com/nameIess/Icons8-Collector/archive/refs/heads/master.zip)**

Or clone the repository:

```bash
git clone https://github.com/nameIess/Icons8-Collector.git
cd Icons8-Collector
```

## ✨ Features

- 🔐 **Automatic login** - Logs into your Icons8 account automatically
- 💾 **Session caching** - Remembers your login (no need to login every time)
- 🖼️ **Multiple formats** - Download as PNG, ICO, or both
- 📐 **Custom sizes** - Choose from 64px to 512px (or custom)
- 🤖 **Headless mode** - Runs invisibly in the background
- 🎛️ **Interactive UI** - Beautiful terminal interface

## 🚀 Quick Start

### Installation

```bash
# Install dependencies
pip install -r requirements.txt

# Install Playwright browser
python -m playwright install chromium
```

### Usage

#### One-liner (recommended for scripts):

```bash
python Icons8-Collector.py --url "https://icons8.com/icons/collections/YOUR_COLLECTION_ID" --email "your@email.com" --password "yourpassword"
```

#### Interactive mode (just run without arguments):

```bash
python Icons8-Collector.py
```

This opens a nice terminal UI where you can enter all options step by step.

## 📋 Command Line Options

| Option          | Alias | Default | Description                             |
| --------------- | ----- | ------- | --------------------------------------- |
| `--url`         | `-u`  | -       | Collection URL (required)               |
| `--email`       | `-e`  | -       | Icons8 account email                    |
| `--password`    | `-P`  | -       | Icons8 account password                 |
| `--format`      | `-f`  | `ico`   | Output format: `png`, `ico`, or `both`  |
| `--size`        | `-z`  | `256`   | Icon size in pixels (64, 128, 256, 512) |
| `--output`      | `-o`  | `data`  | Output directory path                   |
| `--visible`     | `-v`  | `false` | Show browser window (for debugging)     |
| `--interactive` | `-i`  | `false` | Force interactive mode with terminal UI |
| `--help`        | `-h`  | -       | Show help message and exit              |

### Option Details

```
--url, -u           Required. The Icons8 collection URL
                    Example: https://icons8.com/icons/collections/abc123

--email, -e         Your Icons8 account email address
                    Leave empty to use cached login session

--password, -P      Your Icons8 account password (capital P)
                    Only needed if not already logged in

--format, -f        Choose output format:
                    • png  - PNG files only
                    • ico  - ICO files only (default, deletes PNG after conversion)
                    • both - Keep both PNG and ICO files

--size, -z          Icon size in pixels. Common values:
                    • 64   - Small
                    • 128  - Medium
                    • 256  - Large (default, best quality)
                    • 512  - Extra large

--output, -o        Directory where icons will be saved
                    Default: ./data/

--visible, -v       Show the browser window during scraping
                    Useful for debugging login issues

--interactive, -i   Launch interactive terminal UI
                    Prompts for all options step by step
```

## 📝 Examples

### Download collection as ICO files (default):

```bash
python Icons8-Collector.py --url "https://icons8.com/icons/collections/abc123" --email "me@email.com" --password "mypass"
```

### Download as PNG with custom size:

```bash
python Icons8-Collector.py --url "https://icons8.com/icons/collections/abc123" --email "me@email.com" --password "mypass" --format png --size 512
```

### Download both PNG and ICO:

```bash
python Icons8-Collector.py --url "https://icons8.com/icons/collections/abc123" --email "me@email.com" --password "mypass" --format both
```

### Run with visible browser (for debugging):

```bash
python Icons8-Collector.py --url "https://icons8.com/icons/collections/abc123" --email "me@email.com" --password "mypass" --visible
```

## 📁 Output Structure

```
data/
├── Collection_PNG/     # PNG files (if format is png or both)
│   ├── icon_name_1.png
│   └── icon_name_2.png
└── Collection_ICO/     # ICO files (if format is ico or both)
    ├── icon_name_1.ico
    └── icon_name_2.ico
```

## 🔒 Privacy & Security

- Your login session is stored locally in `.browser_data/` folder
- Credentials are never saved to disk
- The `.browser_data/` folder is gitignored
- Run with `--visible` to see exactly what the script is doing `Headless Recommended`

## 🛠️ Requirements

- Python 3.8+
- Icons8 account (free or paid)
- Chrome browser (optional, falls back to Chromium)

## 📜 License

MIT License - See [LICENSE](License) file for details.

### Examples

```bash
# Download 20 folder icons in fluent style
python Icons8-Collector.py --search "folder" --style fluent --amount 20

# Download 100 home icons in color style, size 256px
python Icons8-Collector.py --search "home" --style color --amount 100 --size 256

# Download icons without ICO conversion
python Icons8-Collector.py --search "arrow" --style material --no-ico

# Download icons from a collection (login required)
python Icons8-Collector.py --url "https://icons8.com/icons/collections/xxx" --email "your@email.com" --password "yourpassword"
```

## Output

Icons are saved to the output directory (default: `data`):

- `data/{Style}_PNG/` - PNG format icons
- `data/{Style}_ICO/` - ICO format icons (unless `--no-ico` is used)
- For collections, output folders are named `Collection_PNG` and `Collection_ICO`

## Notes

- Collection URLs require user authentication and browser automation. You may need to log in manually if credentials are not provided.
- The script only downloads preview images (max resolution 550px) available on Icons8.
- If you only want PNG files, use the `--no-ico` option.
- The script uses Playwright for scraping collections. Make sure Chromium is installed (`python -m playwright install chromium`).

## Troubleshooting

- If Playwright is not installed, the script will attempt to install it automatically.
- For collection downloads, if auto-login fails, log in manually in the browser window and press ENTER in the terminal to continue.
- If you encounter issues, check the output and debug logs for details.

## License

MIT License. See License file for details.
