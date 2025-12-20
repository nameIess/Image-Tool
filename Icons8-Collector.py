import os
import argparse
import requests
import re
import asyncio
import shutil
from PIL import Image

# Default settings
DEFAULT_SIZE = 256
BROWSER_DATA_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '.browser_data')

HEADERS = {
    'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
    'Accept': 'application/json',
}


async def scrape_collection_icons(url, size=DEFAULT_SIZE, email=None, password=None, headless=True):
    try:
        from playwright.async_api import async_playwright
    except ImportError:
        print("Playwright is not installed. Installing...")
        os.system("pip install playwright")
        os.system("python -m playwright install chromium")
        from playwright.async_api import async_playwright
    
    icons = []
    os.makedirs(BROWSER_DATA_DIR, exist_ok=True)
    
    async with async_playwright() as p:
        mode_str = "headless" if headless else "visible"
        print(f"Launching browser ({mode_str})...")
        
        try:
            context = await p.chromium.launch_persistent_context(
                user_data_dir=BROWSER_DATA_DIR,
                headless=headless,
                channel="chrome"
            )
        except:
            print("Could not find Chrome, using Chromium instead...")
            context = await p.chromium.launch_persistent_context(
                user_data_dir=BROWSER_DATA_DIR,
                headless=headless
            )
        
        page = await context.new_page()
        
        print(f"Opening collection page: {url}")
        await page.goto(url, timeout=60000)
        await asyncio.sleep(5)
        
        print("Waiting for page to fully load...")
        await asyncio.sleep(3)
        
        # Check if icons are already visible (meaning we're logged in)
        icons_visible = await page.locator('div.app-grid-icon__image img, img[srcset*="icons8.com"]').count()
        print(f"Icons already visible: {icons_visible}")
        
        if icons_visible > 0:
            print("Already logged in! Skipping login process...")
        elif email and password:
            print("Not logged in - attempting login...")
            
            try:
                btn_found = await page.evaluate('''() => {
                    const btns = document.querySelectorAll('.login-button, [class*="login"], button');
                    for (const btn of btns) {
                        if (btn.textContent.includes('Sign in') || btn.classList.contains('login-button')) {
                            return true;
                        }
                    }
                    return false;
                }''')
                
                if btn_found:
                    print("Clicking Sign in button...")
                    await page.evaluate('''() => {
                        const btns = document.querySelectorAll('.login-button, [class*="login"], button');
                        for (const btn of btns) {
                            if (btn.textContent.includes('Sign in') || btn.classList.contains('login-button')) {
                                btn.click();
                                return true;
                            }
                        }
                        return false;
                    }''')
                    await asyncio.sleep(4)
                    
                    print("Waiting for login form...")
                    await page.wait_for_selector('input[type="email"], input[placeholder*="mail"]', timeout=15000)
                    await asyncio.sleep(1)
                    
                    print(f"Filling email: {email}")
                    email_input = page.locator('input[type="email"], input[placeholder*="mail"]').first
                    await email_input.fill(email)
                    await asyncio.sleep(1)
                    
                    print("Filling password...")
                    password_input = page.locator('input[type="password"]').first
                    await password_input.fill(password)
                    await asyncio.sleep(2)
                    
                    print("Waiting for captcha verification (if any)...")
                    await asyncio.sleep(5)
                    
                    print("Clicking Log in button...")
                    await page.evaluate('''() => {
                        const btns = document.querySelectorAll('button');
                        for (const btn of btns) {
                            if (btn.textContent.trim() === 'Log in' || btn.classList.contains('i8-login-form__submit')) {
                                btn.click();
                                return true;
                            }
                        }
                        const form = document.querySelector('form');
                        if (form) form.submit();
                        return false;
                    }''')
                    
                    print("Waiting for login to complete...")
                    await asyncio.sleep(8)
                    
                    print(f"Reloading collection page...")
                    await page.goto(url, timeout=60000)
                    await asyncio.sleep(5)
                    
                    print("Login completed!")
                else:
                    print("No Sign in button found - trying login page...")
                    await page.goto("https://icons8.com/login", timeout=60000)
                    await asyncio.sleep(5)
                    
                    print("Waiting for login form...")
                    await page.wait_for_selector('input[type="email"], input[placeholder*="mail"], input[placeholder="Email"]', timeout=20000)
                    await asyncio.sleep(2)
                    
                    print(f"Filling email: {email}")
                    email_input = page.locator('input[type="email"], input[placeholder*="mail"], input[placeholder="Email"]').first
                    await email_input.fill(email)
                    await asyncio.sleep(2)
                    
                    print("Filling password...")
                    password_input = page.locator('input[type="password"]').first
                    await password_input.fill(password)
                    await asyncio.sleep(2)
                    
                    print("Waiting for captcha verification (if any)...")
                    await asyncio.sleep(5)
                    
                    print("Clicking Log in button...")
                    await page.evaluate('''() => {
                        const btns = document.querySelectorAll('button');
                        for (const btn of btns) {
                            if (btn.textContent.trim() === 'Log in' || btn.classList.contains('i8-login-form__submit')) {
                                btn.click();
                                return true;
                            }
                        }
                        return false;
                    }''')
                    
                    await asyncio.sleep(8)
                    
                    print(f"Going to collection page...")
                    await page.goto(url, timeout=60000)
                    await asyncio.sleep(5)
                    
                    print("Login completed!")
                    
            except Exception as e:
                print(f"Auto-login failed: {e}")
                print("Error: Could not log in automatically. Please check your credentials.")
                await context.close()
                return []
        else:
            print("No icons visible and no credentials provided.")
            print("Please provide email and password to log in.")
            await context.close()
            return []
        
        print(f"Loading collection page...")
        await asyncio.sleep(2)
        
        print("Waiting for icons to load...")
        try:
            await page.wait_for_selector('.app-grid-icon__image, .collection-icon, img[srcset*="icons8"]', timeout=15000)
        except:
            print("Warning: Could not find icon elements with standard selectors")
        
        await asyncio.sleep(2)
        
        title = await page.title()
        print(f"Page title: {title}")
        
        # Scroll to load all icons
        print("Scrolling to load icons...")
        prev_count = 0
        for i in range(20):
            await page.evaluate('window.scrollTo(0, document.body.scrollHeight)')
            await asyncio.sleep(1.5)
            
            icon_elements = await page.locator('div.app-grid-icon__image img').count()
            if icon_elements == 0:
                icon_elements = await page.locator('img[srcset*="icons8.com"]').count()
            
            print(f"  Scroll {i+1}: Found {icon_elements} icon elements...")
            
            if icon_elements == prev_count and i > 5 and icon_elements > 0:
                break
            prev_count = icon_elements
        
        print("\nExtracting icons from page...")
        
        icon_imgs = page.locator('div.app-grid-icon__image img')
        count = await icon_imgs.count()
        
        if count == 0:
            print("Trying alternative selector: img[srcset*='icons8']")
            icon_imgs = page.locator('img[srcset*="icons8.com"]')
            count = await icon_imgs.count()
        
        # Fallback: regex extraction
        if count == 0:
            print("Trying regex extraction from page content...")
            content = await page.content()
            
            id_matches = re.findall(r'img\.icons8\.com/?\?[^"\'>\s]*id=([A-Za-z0-9_-]+)[^"\'>\s]*', content)
            print(f"Found {len(id_matches)} icon IDs via regex")
            
            seen_ids = set()
            for icon_id in id_matches:
                if icon_id not in seen_ids:
                    seen_ids.add(icon_id)
                    icons.append({
                        'id': icon_id,
                        'name': f'icon-{icon_id}',
                        'url': f"https://img.icons8.com/?size={size}&id={icon_id}&format=png"
                    })
            
            if icons:
                print(f"Extracted {len(icons)} icons via regex")
                await context.close()
                return icons
        
        print(f"Found {count} icon images via DOM")
        
        seen_ids = set()
        for i in range(count):
            try:
                img = icon_imgs.nth(i)
                srcset = await img.get_attribute('srcset')
                alt = await img.get_attribute('alt') or f'icon_{i}'
                
                if srcset:
                    id_match = re.search(r'id=([A-Za-z0-9_-]+)', srcset)
                    if id_match:
                        icon_id = id_match.group(1)
                        if icon_id not in seen_ids:
                            seen_ids.add(icon_id)
                            name = alt.replace(' icon', '').strip()
                            
                            icons.append({
                                'id': icon_id,
                                'name': name,
                                'url': f"https://img.icons8.com/?size={size}&id={icon_id}&format=png"
                            })
                            print(f"  Found: {name} (ID: {icon_id})")
            except Exception as e:
                print(f"  Error extracting icon {i}: {e}")
                continue
        
        print(f"\nTotal unique icons found: {len(icons)}")
        await context.close()
    
    return icons


def get_collection_icons(url, size=DEFAULT_SIZE, email=None, password=None, headless=True):
    return asyncio.run(scrape_collection_icons(url, size, email, password, headless))


def download_icon(url, output_path):
    try:
        response = requests.get(url, headers=HEADERS, stream=True)
        response.raise_for_status()
        
        content_type = response.headers.get('content-type', '')
        if 'image' not in content_type and len(response.content) < 100:
            print(f"  Warning: Not an image or empty response")
            return False
        
        with open(output_path, 'wb') as f:
            for chunk in response.iter_content(chunk_size=8192):
                f.write(chunk)
        return True
    except requests.RequestException as e:
        print(f"  Failed to download: {e}")
        return False


def convert_to_ico(png_path, ico_path):
    try:
        img = Image.open(png_path)
        if img.mode != 'RGBA':
            img = img.convert('RGBA')
        img.save(ico_path, format='ICO', sizes=[(img.width, img.height)])
        return True
    except Exception as e:
        print(f"  Failed to convert to ICO: {e}")
        return False


def clear_screen():
    os.system('cls' if os.name == 'nt' else 'clear')


def print_header():
    clear_screen()
    print("\n")
    print("  ╔══════════════════════════════════════════════════════════╗")
    print("  ║                                                          ║")
    print("  ║             🎨  ICONS8 DOWNLOADER  🎨                   ║")
    print("  ║                                                          ║")
    print("  ║        Download icons from your Icons8 collections       ║")
    print("  ║                                                          ║")
    print("  ╚══════════════════════════════════════════════════════════╝")
    print("\n")


def print_section(title):
    print(f"\n  ┌─ {title} " + "─" * (50 - len(title)) + "┐")


def print_option(num, text, default=False):
    default_str = " (default)" if default else ""
    print(f"  │  [{num}] {text}{default_str}")


def print_section_end():
    print("  └" + "─" * 54 + "┘")


def get_input(prompt, default=None):
    if default:
        result = input(f"  │  {prompt} [{default}]: ").strip()
        return result if result else default
    return input(f"  │  {prompt}: ").strip()


def get_user_input():
    print_header()
    
    # Collection URL
    print_section("COLLECTION URL")
    print("  │")
    url = get_input("Enter collection URL")
    print("  │")
    print_section_end()
    
    # Login credentials
    print_section("LOGIN CREDENTIALS")
    print("  │")
    print("  │  (Leave empty if already logged in)")
    print("  │")
    email = get_input("Email")
    if email:
        password = get_input("Password")
    else:
        password = None
    print("  │")
    print_section_end()
    
    # Output format
    print_section("OUTPUT FORMAT")
    print("  │")
    print_option(1, "PNG only")
    print_option(2, "ICO only (deletes PNG after conversion)", default=True)
    print_option(3, "Both PNG and ICO")
    print("  │")
    while True:
        format_choice = get_input("Select format", "2")
        if format_choice in ['1', '2', '3']:
            break
        print("  │  ⚠ Invalid choice. Please enter 1, 2, or 3.")
    print("  │")
    print_section_end()
    
    # Icon size
    print_section("ICON SIZE")
    print("  │")
    print_option(1, "64px  - Small")
    print_option(2, "128px - Medium")
    print_option(3, "256px - Large (Best Quality)", default=True)
    print_option(4, "512px - Extra Large")
    print_option(5, "Custom size")
    print("  │")
    size_choice = get_input("Select size", "3")
    size_map = {'1': 64, '2': 128, '3': 256, '4': 512}
    if size_choice in size_map:
        size = size_map[size_choice]
    elif size_choice == '5':
        custom = get_input("Enter custom size (px)", "256")
        size = int(custom) if custom.isdigit() else 256
    else:
        size = 256
    print("  │")
    print_section_end()
    
    # Browser mode
    print_section("BROWSER MODE")
    print("  │")
    print_option(1, "Headless (invisible)", default=True)
    print_option(2, "Visible (show browser window)")
    print("  │")
    browser_choice = get_input("Select mode", "1")
    headless = browser_choice != '2'
    print("  │")
    print_section_end()
    
    # Confirmation
    print("\n")
    print("  ┌─ SUMMARY " + "─" * 43 + "┐")
    print("  │")
    print(f"  │  Collection: {url[:45]}..." if len(url) > 45 else f"  │  Collection: {url}")
    print(f"  │  Email:      {email if email else '(using saved session)'}")
    format_names = {'1': 'PNG only', '2': 'ICO only', '3': 'Both PNG and ICO'}
    print(f"  │  Format:     {format_names[format_choice]}")
    print(f"  │  Size:       {size}px")
    print(f"  │  Browser:    {'Headless' if headless else 'Visible'}")
    print("  │")
    print("  └" + "─" * 54 + "┘")
    print("\n")
    
    confirm = input("  Press ENTER to start download (or 'q' to quit): ").strip().lower()
    if confirm == 'q':
        print("\n  Cancelled.\n")
        return None
    
    return {
        'url': url,
        'email': email if email else None,
        'password': password if password else None,
        'format_choice': format_choice,
        'size': size,
        'headless': headless
    }


def main():
    parser = argparse.ArgumentParser(description='Download icons from Icons8.com collections')
    parser.add_argument('--url', '-u', type=str, help='Collection URL to scrape icons from')
    parser.add_argument('--email', '-e', type=str, help='Icons8 account email for login')
    parser.add_argument('--password', '-P', type=str, help='Icons8 account password for login')
    parser.add_argument('--size', '-z', type=int, default=DEFAULT_SIZE,
                        help=f'Icon size in pixels (default: {DEFAULT_SIZE})')
    parser.add_argument('--output', '-o', type=str, default='data',
                        help='Output directory (default: data)')
    parser.add_argument('--format', '-f', type=str, choices=['png', 'ico', 'both'], default='ico',
                        help='Output format: png, ico, or both (default: ico)')
    parser.add_argument('--visible', '-v', action='store_true',
                        help='Show browser window (default: headless)')
    parser.add_argument('--interactive', '-i', action='store_true',
                        help='Run in interactive mode (prompts for input)')
    
    args = parser.parse_args()
    headless = not args.visible
    
    # Interactive mode if no URL provided or explicitly requested
    if args.interactive or not args.url:
        user_input = get_user_input()
        
        if user_input is None:
            return
        
        args.url = user_input['url']
        args.email = user_input['email']
        args.password = user_input['password']
        args.size = user_input['size']
        headless = user_input.get('headless', True)
        
        format_map = {'1': 'png', '2': 'ico', '3': 'both'}
        args.format = format_map[user_input['format_choice']]
    
    # Scrape from collection URL
    print(f"\n  📂 Scraping collection from: {args.url}")
    icons = get_collection_icons(args.url, args.size, args.email, args.password, headless)
    
    if not icons:
        print("No icons found!")
        return
    
    print(f"Found {len(icons)} icons")
    
    # Create output directories
    png_path = None
    ico_path = None
    
    if args.format in ['png', 'both']:
        png_path = os.path.join(args.output, "Collection_PNG")
        os.makedirs(png_path, exist_ok=True)
    
    if args.format in ['ico', 'both']:
        ico_path = os.path.join(args.output, "Collection_ICO")
        os.makedirs(ico_path, exist_ok=True)
    
    # Temp directory for ICO-only mode
    temp_png_path = None
    if args.format == 'ico':
        temp_png_path = os.path.join(args.output, '.temp_png')
        os.makedirs(temp_png_path, exist_ok=True)
    
    output_dir = png_path if png_path else temp_png_path
    print(f"\nDownloading icons...")
    downloaded = 0
    converted = 0
    
    for i, icon in enumerate(icons, 1):
        name = icon.get('name', f'icon_{i}')
        safe_name = "".join(c for c in name if c.isalnum() or c in (' ', '-', '_')).rstrip()
        safe_name = safe_name.replace(' ', '_')
        
        if not safe_name:
            safe_name = f'icon_{i}'
        
        url = icon['url']
        png_file = os.path.join(output_dir, f"{safe_name}.png")
        
        print(f"[{i}/{len(icons)}] {name}...")
        
        if download_icon(url, png_file):
            downloaded += 1
            
            if args.format in ['ico', 'both']:
                ico_file = os.path.join(ico_path, f"{safe_name}.ico")
                if convert_to_ico(png_file, ico_file):
                    converted += 1
                
                if args.format == 'ico':
                    os.remove(png_file)
    
    # Clean up temp directory
    if temp_png_path and os.path.exists(temp_png_path):
        shutil.rmtree(temp_png_path)
    
    # Print summary
    print("\n")
    print("  ╔══════════════════════════════════════════════════════════╗")
    print("  ║                                                          ║")
    print("  ║              ✅  DOWNLOAD COMPLETE!  ✅                 ║")
    print("  ║                                                          ║")
    print("  ╚══════════════════════════════════════════════════════════╝")
    print("\n")
    
    if args.format == 'png':
        print(f"  📁 Downloaded {downloaded} PNG files")
        print(f"  📂 Location: {png_path}")
    elif args.format == 'ico':
        print(f"  📁 Converted {converted} ICO files")
        print(f"  📂 Location: {ico_path}")
    else:
        print(f"  📁 Downloaded {downloaded} PNG files to: {png_path}")
        print(f"  📁 Converted {converted} ICO files to: {ico_path}")
    print("\n")


if __name__ == "__main__":
    main()
