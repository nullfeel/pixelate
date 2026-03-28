[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

# pixelate

Convert images to stunning ASCII art with **true color** terminal output.

```
    @@@@@@@@@@@@@@@@@@@@
  @@88888888888888888888@@
 @8888GGGGGGGGGGGGGGg888@
 @888GL11111111111111fG88@
@888GC1;;;;;;;;;;;;;1tG88@
@888GL1;,,,,,,,,,,,;1tG88@
@888GL1;,...........;1tG88@
@888GL1;,...........;1tG88@
@888GL1;,,,,,,,,,,,;1tG88@
@888GC1;;;;;;;;;;;;;1tG88@
 @888GLt1111111111ttfG88@
 @8888GGGGGGGGGGGGGGg888@
  @@88888888888888888888@@
    @@@@@@@@@@@@@@@@@@@@
```

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
# Basic — prints colored ASCII art to terminal
pixelate photo.jpg

# Custom width
pixelate -w 120 photo.png

# Grayscale mode
pixelate -g photo.jpg

# Block characters (great for high-density look)
pixelate --preset block photo.jpg

# Braille characters (ultra-detailed)
pixelate --preset braille photo.jpg

# Detailed ASCII charset
pixelate --preset detailed photo.jpg

# Invert colors (for light terminals)
pixelate --invert photo.jpg

# Adjust brightness and contrast
pixelate --brightness 1.3 --contrast 1.5 photo.jpg

# Dithering for smoother grayscale gradients
pixelate -g --dither photo.jpg

# Save as plain text
pixelate -o art.txt photo.jpg

# Save with ANSI colors (viewable with cat)
pixelate -o art.ans photo.jpg

# Save as HTML
pixelate -o art.html photo.jpg

# From URL
pixelate https://example.com/photo.jpg

# Pipe from stdin
curl -s https://example.com/photo.jpg | pixelate -

# Multiple files
pixelate *.jpg

# Animated GIF (plays in terminal)
pixelate animation.gif
```

## Presets

| Preset | Characters | Best for |
|--------|-----------|----------|
| `simple` | ` .,:;i1tfLCG08@` | General use |
| `detailed` | 70 ASCII characters | Maximum detail |
| `block` | ` ░▒▓█` | Dense, blocky look |
| `braille` | `⠀⠁⠃⠇⡇⣇⣧⣷⣿` | Ultra-high resolution feel |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--width` | `-w` | Output width in characters (default: terminal width) |
| `--grayscale` | `-g` | Grayscale output (no colors) |
| `--invert` | | Invert light/dark mapping |
| `--output` | `-o` | Save to file (.txt, .ans, .html) |
| `--charset` | | Custom character set (light to dark) |
| `--preset` | | Preset: simple, detailed, block, braille |
| `--brightness` | | Brightness multiplier (default: 1.0) |
| `--contrast` | | Contrast multiplier (default: 1.0) |
| `--dither` | | Floyd-Steinberg dithering (grayscale) |

## Supported Formats

- JPEG
- PNG
- GIF (static and animated)
- BMP
- WebP

## How It Works

1. Load and decode the image
2. Resize to target width, applying 0.5 height factor to correct for terminal character aspect ratio
3. For each pixel, calculate luminance and map to a character from the charset
4. Wrap each character in a 24-bit ANSI true color escape sequence (`\033[38;2;R;G;Bm`)
5. Print row by row to the terminal

The result is ASCII art that preserves the original colors of the image, producing vivid output on any terminal that supports true color (most modern terminals do).

## Note

Early prototype, may have edge cases. Contributions welcome.

---

Built by [nullfeel](https://github.com/nullfeel)
