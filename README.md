<div align="center">

# pixelate

**Convert images to stunning ASCII art in your terminal.**

True color · Image filters · Animated GIFs · Multiple charsets · SVG/HTML export

![Go](https://img.shields.io/badge/Go-0d0d0d?style=flat-square&logo=go&logoColor=00ADD8)
![License](https://img.shields.io/badge/License-MIT-0d0d0d?style=flat-square)

</div>

## Features

- **True Color Output** — 24-bit ANSI colors for stunning terminal art
- **Multiple Color Modes** — True color, 256-color, 16-color, or no color
- **Image Filters** — Edge detection, negative, sepia, blur, sharpen, pixelate
- **Character Presets** — Simple, detailed, block, braille, dots, ascii
- **Animated GIFs** — Play animated GIFs as ASCII in the terminal
- **Multiple Export Formats** — TXT, ANSI, HTML, SVG
- **URL Support** — Convert images directly from URLs
- **Stdin Piping** — Works with pipes: `curl ... | pixelate -`
- **Brightness/Contrast** — Fine-tune image appearance
- **Floyd-Steinberg Dithering** — Better gradients in grayscale mode

## Install

```bash
go install github.com/nullfeel/pixelate@latest
```

Or build from source:

```bash
git clone https://github.com/nullfeel/pixelate.git
cd pixelate
go build -o pixelate .
```

## Usage

```bash
# Basic conversion
pixelate photo.jpg

# Set output width
pixelate -w 100 photo.jpg

# Grayscale mode
pixelate -g photo.jpg

# Apply filters
pixelate --filter edge photo.jpg
pixelate --filter sepia photo.jpg
pixelate --filter negative photo.jpg
pixelate --filter blur --filter-intensity 3 photo.jpg
pixelate --filter pixelate --filter-intensity 12 photo.jpg
pixelate --filter sharpen photo.jpg

# Different charsets
pixelate --preset block photo.jpg
pixelate --preset braille photo.jpg
pixelate --preset detailed photo.jpg
pixelate --preset dots photo.jpg

# Color modes
pixelate --color-mode true photo.jpg     # 24-bit (default)
pixelate --color-mode 256 photo.jpg      # 256-color ANSI
pixelate --color-mode 16 photo.jpg       # 16-color ANSI
pixelate --color-mode none photo.jpg     # characters only

# Save to file
pixelate -o art.txt photo.jpg            # plain text
pixelate -o art.ans photo.jpg            # ANSI colors
pixelate -o art.html photo.jpg           # HTML
pixelate -o art.svg photo.jpg            # SVG vector

# From URL
pixelate https://example.com/photo.jpg

# Pipe from stdin
curl -s https://example.com/photo.jpg | pixelate -

# Adjust image
pixelate --brightness 1.3 --contrast 1.2 photo.jpg

# Animated GIF playback
pixelate animation.gif

# Combine options
pixelate -w 120 --filter edge --preset block --color-mode 256 photo.jpg
```

## Presets

| Preset | Characters | Best For |
|--------|-----------|----------|
| `simple` | ` .,:;i1tfLCG08@` | General use |
| `detailed` | 70 character gradient | High detail |
| `block` | `░▒▓█` | Pixel art feel |
| `braille` | `⠀⠁⠃⠇⡇⣇⣧⣷⣿` | Smooth gradients |
| `dots` | `⠀⠄⠆⠖⠶⡶⣶⣿` | Minimal dots |
| `ascii` | ` .:-=+*#%@` | Classic ASCII |

## Filters

| Filter | Description |
|--------|-------------|
| `edge` | Sobel edge detection — dramatic outlines |
| `negative` | Invert all colors |
| `sepia` | Warm brownish vintage tone |
| `blur` | Box blur (intensity = radius, default 2) |
| `sharpen` | Enhance edges and details |
| `pixelate` | Block mosaic effect (intensity = block size, default 8) |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--width` | `-w` | terminal | Output width in characters |
| `--grayscale` | `-g` | false | No colors, characters only |
| `--invert` | | false | Invert light/dark |
| `--output` | `-o` | stdout | Save to file (.txt .ans .html .svg) |
| `--preset` | | simple | Character preset |
| `--charset` | | | Custom characters (light to dark) |
| `--filter` | | none | Image filter to apply |
| `--filter-intensity` | | 2.0 | Filter strength |
| `--color-mode` | | true | Color: true, 256, 16, none |
| `--brightness` | | 1.0 | Brightness multiplier |
| `--contrast` | | 1.0 | Contrast multiplier |
| `--dither` | | false | Floyd-Steinberg dithering |

## Export Formats

| Extension | Description |
|-----------|-------------|
| `.txt` | Plain text, no colors |
| `.ans` | Text with ANSI escape codes |
| `.html` | HTML with inline color styles |
| `.svg` | Scalable vector graphic |

## How It Works

1. Load and decode the image (JPG, PNG, GIF, BMP, WebP)
2. Apply image filter if specified (edge, sepia, blur, etc.)
3. Resize to target width with Lanczos3 resampling
4. Apply 0.5 height factor for terminal character aspect ratio
5. Calculate per-pixel luminance (ITU-R BT.601)
6. Map luminance to character from selected charset
7. Apply color formatting based on color mode
8. Output to terminal or save to file

## Requirements

- Go 1.21+
- Terminal with color support recommended
- True color terminal for best results (Windows Terminal, iTerm2, Alacritty, Kitty)

> **Note:** Early prototype. Some edge cases may not be handled. Contributions welcome.

## License

MIT

---

<div align="center">

Built by [nullfeel](https://github.com/nullfeel)

</div>
