<div align="center">

<br/>

# pixelate

**Convert images to ASCII art in your terminal.**

True color output -- image filters -- animated GIFs -- multiple charsets -- SVG and HTML export

<br/>

![Go](https://img.shields.io/badge/Go-0d0d0d?style=flat-square&logo=go&logoColor=00ADD8)
![License](https://img.shields.io/badge/License-MIT-0d0d0d?style=flat-square)

<br/>

> **Early prototype.** Some edge cases may not be handled. Contributions and feedback are welcome.

</div>

<br/>

## About

pixelate is a command-line tool written in Go that converts images into ASCII art with full true color support. It reads JPG, PNG, GIF, BMP, and WebP files from local paths, URLs, or stdin, applies optional image filters and brightness/contrast adjustments, then maps each pixel to a character based on its luminance. Output can be rendered directly in the terminal with 24-bit ANSI colors or exported to TXT, ANSI, HTML, and SVG formats.

## Features

- **True Color Output** - 24-bit ANSI color codes for terminal rendering that preserves the original image palette
- **Multiple Color Modes** - Choose between true color (24-bit), 256-color, 16-color ANSI, or no color at all
- **Six Character Presets** - Built-in charsets optimized for different visual styles, from simple ASCII to Unicode braille
- **Custom Character Sets** - Define your own light-to-dark character gradient with the `--charset` flag
- **Six Image Filters** - Edge detection, negative, sepia, blur, sharpen, and pixelate with configurable intensity
- **Animated GIF Playback** - Detects animated GIFs and plays them frame-by-frame as ASCII in the terminal
- **Four Export Formats** - Save output as plain text, ANSI-encoded text, standalone HTML, or scalable SVG
- **URL Support** - Fetch and convert images directly from HTTP/HTTPS URLs
- **Stdin Piping** - Accepts image data from pipes for integration with other tools
- **Brightness and Contrast** - Fine-tune image appearance with multiplier-based adjustments
- **Floyd-Steinberg Dithering** - Error-diffusion dithering for smoother gradients in grayscale mode
- **Invert Mode** - Reverse the character mapping for light-on-dark or dark-on-light output
- **Multi-Image Batch** - Process multiple images in a single command
- **Lanczos3 Resampling** - High-quality image downscaling before character mapping

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

### Basic conversion

```bash
pixelate photo.jpg
```

### Set output width

```bash
pixelate -w 100 photo.jpg
pixelate -w 200 landscape.png
```

### Grayscale mode

```bash
pixelate -g photo.jpg
pixelate --grayscale photo.jpg
```

### Apply filters

```bash
pixelate --filter edge photo.jpg
pixelate --filter sepia photo.jpg
pixelate --filter negative photo.jpg
pixelate --filter blur --filter-intensity 3 photo.jpg
pixelate --filter sharpen photo.jpg
pixelate --filter pixelate --filter-intensity 12 photo.jpg
pixelate --filter grayscale photo.jpg
```

### Use different charsets

```bash
pixelate --preset block photo.jpg
pixelate --preset braille photo.jpg
pixelate --preset detailed photo.jpg
pixelate --preset dots photo.jpg
pixelate --preset ascii photo.jpg
```

### Custom character set

```bash
pixelate --charset " .oO@" photo.jpg
pixelate --charset " ░▒▓█" photo.jpg
```

### Color modes

```bash
pixelate --color-mode true photo.jpg      # 24-bit true color (default)
pixelate --color-mode 256 photo.jpg       # 256-color ANSI palette
pixelate --color-mode 16 photo.jpg        # 16-color ANSI
pixelate --color-mode none photo.jpg      # Characters only, no color
```

### Export to file

```bash
pixelate -o art.txt photo.jpg             # Plain text
pixelate -o art.ans photo.jpg             # ANSI escape codes
pixelate -o art.html photo.jpg            # Standalone HTML page
pixelate -o art.svg photo.jpg             # Scalable vector graphic
```

### Load from URL

```bash
pixelate https://example.com/photo.jpg
pixelate -w 80 https://example.com/landscape.png
```

### Pipe from stdin

```bash
curl -s https://example.com/photo.jpg | pixelate -
cat photo.png | pixelate -w 60 -
```

### Adjust brightness and contrast

```bash
pixelate --brightness 1.3 photo.jpg
pixelate --contrast 1.5 photo.jpg
pixelate --brightness 1.2 --contrast 1.3 photo.jpg
```

### Dithering

```bash
pixelate -g --dither photo.jpg
pixelate --grayscale --dither -w 120 photo.jpg
```

### Invert characters

```bash
pixelate --invert photo.jpg
pixelate --invert --preset block photo.jpg
```

### Animated GIF playback

```bash
pixelate animation.gif
pixelate -w 60 loading.gif
```

### Combine options

```bash
pixelate -w 120 --filter edge --preset block --color-mode 256 photo.jpg
pixelate -w 80 --filter sepia --brightness 1.1 --preset braille -o art.html photo.jpg
pixelate -g --dither --invert -w 100 --preset detailed photo.jpg
```

### Batch processing

```bash
pixelate photo1.jpg photo2.png photo3.bmp
```

## Example Output

The `block` preset with a simple gradient:

```
████████████████████████████████████████
████████████████████████████████████████
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒
░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
```

The `simple` preset with a circle shape:

```
          ,,,,iiii,,,,
      ,,ii11ttfftt11ii,,
    ,,11ttfffffftttt11,,
  ,,ii1ttffffffffffff1ii,,
  ,,11ttffffffffffttt11,,
  ,,i1ttfffffffffffft1i,,
  ,,11ttffffffffffttt11,,
  ,,ii1ttffffffffffff1ii,,
    ,,11ttfffffftttt11,,
      ,,ii11ttfftt11ii,,
          ,,,,iiii,,,,
```

The `ascii` preset with a diamond:

```
          ..
        ..--..
      ..--==--..
    ..--==++==--..
  ..--==++**++==--..
..--==++*##*++==--..
  ..--==++**++==--..
    ..--==++==--..
      ..--==--..
        ..--..
          ..
```

## Presets

| Preset | Characters | Description |
|--------|-----------|-------------|
| `simple` | ` .,:;i1tfLCG08@` | General-purpose 15-character gradient. Good default for most images. |
| `detailed` | ` .'`^",:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$` | 70-character gradient for maximum detail. Best at high widths. |
| `block` | `░▒▓█` | Unicode block elements. Produces a pixel-art mosaic effect. |
| `braille` | `⠀⠁⠃⠇⡇⣇⣧⣷⣿` | Unicode braille patterns. Creates smooth, dense gradients. |
| `dots` | `⠀⠄⠆⠖⠶⡶⣶⣿` | Sparse braille dot patterns. Lighter, more minimal look. |
| `ascii` | ` .:-=+*#%@` | Classic ASCII-only characters. Works in any terminal or font. |

All presets are ordered from lightest (space) to darkest (most filled). The `--invert` flag reverses this order.

## Filters

| Filter | Intensity Parameter | Default Intensity | Description |
|--------|-------------------|-------------------|-------------|
| `edge` | Not used | -- | Sobel edge detection using 3x3 kernels. Extracts edges per RGB channel and computes gradient magnitude. Produces dramatic outline effects. |
| `negative` | Not used | -- | Inverts all RGB channels (`255 - value`). Creates a photographic negative. |
| `sepia` | Not used | -- | Applies a warm brownish tone using the standard sepia matrix transform. Produces a vintage photograph look. |
| `blur` | Radius (pixels) | 2 | Box blur averaging all pixels within the given radius. Higher values produce stronger smoothing. |
| `sharpen` | Not used | -- | Unsharp mask via 3x3 convolution kernel. Enhances edges and fine detail. |
| `pixelate` | Block size (pixels) | 8 | Averages each block of pixels into a single color. Higher values produce larger mosaic blocks. |
| `grayscale` | Not used | -- | Converts to grayscale using ITU-R BT.601 luma coefficients while keeping the RGB color model. |

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--width` | `-w` | Terminal width | Output width in characters. Auto-detected from terminal. Capped at 500. |
| `--grayscale` | `-g` | `false` | Disable colors and output characters only based on luminance. |
| `--invert` | -- | `false` | Reverse the character set so light pixels map to dense characters and dark pixels map to sparse characters. |
| `--output` | `-o` | stdout | Save output to a file. Format is determined by extension: `.txt`, `.ans`, `.html`, `.svg`. |
| `--preset` | -- | `simple` | Character preset to use. One of: `simple`, `detailed`, `block`, `braille`, `dots`, `ascii`. |
| `--charset` | -- | -- | Custom character string ordered from light to dark. Overrides `--preset` when set. |
| `--filter` | -- | `none` | Image filter to apply before conversion. One of: `none`, `edge`, `negative`, `sepia`, `blur`, `sharpen`, `pixelate`, `grayscale`. |
| `--filter-intensity` | -- | `0` (uses filter default) | Strength parameter for filters that support it. Used by `blur` (radius) and `pixelate` (block size). |
| `--color-mode` | -- | `true` | Terminal color mode. One of: `true` (24-bit), `256` (256-color), `16` (16-color), `none` (no color). |
| `--brightness` | -- | `1.0` | Brightness multiplier applied to all pixels. Values above 1.0 brighten, below 1.0 darken. |
| `--contrast` | -- | `1.0` | Contrast multiplier applied around the midpoint (128). Values above 1.0 increase contrast. |
| `--dither` | -- | `false` | Apply Floyd-Steinberg error-diffusion dithering. Only effective in grayscale mode. |

## Export Formats

### Plain Text (`.txt`)

Writes raw characters with no color information. Compatible with any text editor or viewer. Useful for embedding in documents or source code comments.

### ANSI (`.ans`)

Writes characters with ANSI escape sequences for color. Can be viewed with `cat` in a compatible terminal. Preserves the exact terminal output including all color codes.

### HTML (`.html`)

Generates a standalone HTML page with inline color styles on each character using `<span>` elements. Opens in any browser. Uses a monospace font and dark background for accurate rendering.

### SVG (`.svg`)

Generates a scalable vector graphic with each character positioned as a `<text>` element with its original color. Resolution-independent and suitable for printing, embedding in web pages, or further editing in vector graphics tools. Uses a 10px font size grid.

## Color Modes

### True Color (`true`)

Uses 24-bit ANSI escape sequences (`\033[38;2;R;G;Bm`) to render each character in its exact original pixel color. Requires a terminal that supports true color (Windows Terminal, iTerm2, Alacritty, Kitty, most modern terminals). This is the default mode and produces the most accurate color representation.

### 256-Color (`256`)

Maps each pixel color to the nearest entry in the standard 256-color ANSI palette. Compatible with a wider range of terminals. Slight color quantization is visible but the overall image structure is preserved.

### 16-Color (`16`)

Maps each pixel to one of the 16 standard ANSI colors. High color loss but maximum terminal compatibility. Works in virtually every terminal emulator.

### No Color (`none`)

Outputs characters only with no escape sequences. Equivalent to combining with `--grayscale`. Useful for plain text export or terminals with no color support.

## How the Algorithm Works

1. **Load** - The image is decoded from the source (file, URL, or stdin). Supported formats: JPEG, PNG, GIF, BMP, WebP. Animated GIFs are detected and handled separately for frame-by-frame playback.

2. **Filter** - If a filter is specified, it is applied to the full-resolution image before any resizing. Sobel edge detection uses two 3x3 convolution kernels (horizontal and vertical) and computes the gradient magnitude per channel. Blur uses a box kernel with configurable radius. Sharpen uses an unsharp mask kernel. Pixelate averages blocks of pixels.

3. **Resize** - The image is scaled to the target width using Lanczos3 resampling (high-quality sinc-based interpolation). The height is calculated from the aspect ratio with a 0.5 factor to compensate for terminal characters being approximately twice as tall as they are wide.

4. **Adjust** - Brightness is applied as a linear multiplier to each RGB channel. Contrast is applied around the midpoint value of 128: `(channel - 128) * contrast + 128`.

5. **Map** - For each pixel, luminance is calculated using ITU-R BT.601 coefficients: `L = 0.299R + 0.587G + 0.114B`. The luminance value (0.0 to 1.0) is used as an index into the character set, mapping dark pixels to sparse characters and bright pixels to dense characters.

6. **Dither** (optional) - When enabled in grayscale mode, Floyd-Steinberg error diffusion distributes quantization error to neighboring pixels using the standard 7/16, 3/16, 5/16, 1/16 kernel, producing smoother gradients.

7. **Output** - Characters are rendered to the terminal with ANSI color codes based on the selected color mode, or serialized to the chosen export format (TXT, ANSI, HTML, SVG).

## Requirements

- Go 1.21 or later
- A terminal with color support is recommended for the best experience
- True color terminal for full 24-bit output (Windows Terminal, iTerm2, Alacritty, Kitty)
- Non-color modes work in any terminal

## License

MIT

---

<div align="center">

Built by [nullfeel](https://github.com/nullfeel)

</div>
