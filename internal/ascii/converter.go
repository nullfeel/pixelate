package ascii

import (
	"image"
	"image/color"
	"math"

	"github.com/nfnt/resize"
)

// AsciiChar represents a single character in the ASCII art output,
// holding both the character to display and its original pixel color.
type AsciiChar struct {
	Char  rune
	R, G, B uint8
}

// ConvertOptions controls how an image is converted to ASCII art.
type ConvertOptions struct {
	Width      int
	Charset    string
	Grayscale  bool
	Brightness float64
	Contrast   float64
	Dither     bool
}

// Convert transforms an image into a 2D grid of AsciiChar.
// Each pixel is mapped to a character based on its luminance,
// and retains its RGB color for true-color terminal output.
func Convert(img image.Image, opts ConvertOptions) [][]AsciiChar {
	if opts.Width <= 0 {
		opts.Width = 80
	}
	if opts.Charset == "" {
		opts.Charset = GetPreset("simple")
	}
	if opts.Brightness == 0 {
		opts.Brightness = 1.0
	}
	if opts.Contrast == 0 {
		opts.Contrast = 1.0
	}

	chars := []rune(opts.Charset)
	numChars := len(chars)

	// Resize image to target width.
	// Height factor 0.5 compensates for terminal characters being ~2x taller than wide.
	bounds := img.Bounds()
	imgW := bounds.Dx()
	imgH := bounds.Dy()

	targetW := uint(opts.Width)
	ratio := float64(imgH) / float64(imgW)
	targetH := uint(float64(targetW) * ratio * 0.5)

	if targetH < 1 {
		targetH = 1
	}

	resized := resize.Resize(targetW, targetH, img, resize.Lanczos3)
	rb := resized.Bounds()
	w := rb.Dx()
	h := rb.Dy()

	// If dithering is enabled, build a luminance matrix and apply Floyd-Steinberg
	if opts.Dither && opts.Grayscale {
		return convertWithDithering(resized, w, h, chars, numChars, opts)
	}

	art := make([][]AsciiChar, h)
	for y := 0; y < h; y++ {
		row := make([]AsciiChar, w)
		for x := 0; x < w; x++ {
			r, g, b := pixelRGB(resized.At(x+rb.Min.X, y+rb.Min.Y))

			// Apply brightness and contrast
			rf, gf, bf := adjustColor(float64(r), float64(g), float64(b), opts.Brightness, opts.Contrast)
			r8 := clampUint8(rf)
			g8 := clampUint8(gf)
			b8 := clampUint8(bf)

			lum := luminance(r8, g8, b8)
			idx := int(lum * float64(numChars-1))
			if idx >= numChars {
				idx = numChars - 1
			}

			row[x] = AsciiChar{
				Char: chars[idx],
				R:    r8,
				G:    g8,
				B:    b8,
			}
		}
		art[y] = row
	}

	return art
}

// convertWithDithering applies Floyd-Steinberg dithering for smoother grayscale gradients.
func convertWithDithering(img image.Image, w, h int, chars []rune, numChars int, opts ConvertOptions) [][]AsciiChar {
	bounds := img.Bounds()

	// Build floating-point luminance grid
	lum := make([][]float64, h)
	type rgb struct{ R, G, B uint8 }
	colors := make([][]rgb, h)
	for y := 0; y < h; y++ {
		lum[y] = make([]float64, w)
		colors[y] = make([]rgb, w)
		for x := 0; x < w; x++ {
			r, g, b := pixelRGB(img.At(x+bounds.Min.X, y+bounds.Min.Y))
			rf, gf, bf := adjustColor(float64(r), float64(g), float64(b), opts.Brightness, opts.Contrast)
			r8 := clampUint8(rf)
			g8 := clampUint8(gf)
			b8 := clampUint8(bf)
			colors[y][x] = rgb{r8, g8, b8}
			lum[y][x] = luminance(r8, g8, b8)
		}
	}

	// Floyd-Steinberg dithering
	art := make([][]AsciiChar, h)
	for y := 0; y < h; y++ {
		row := make([]AsciiChar, w)
		for x := 0; x < w; x++ {
			oldVal := lum[y][x]
			// Quantize to nearest character level
			idx := int(math.Round(oldVal * float64(numChars-1)))
			if idx < 0 {
				idx = 0
			}
			if idx >= numChars {
				idx = numChars - 1
			}
			newVal := float64(idx) / float64(numChars-1)
			quantErr := oldVal - newVal

			// Distribute error to neighbors
			if x+1 < w {
				lum[y][x+1] += quantErr * 7.0 / 16.0
			}
			if y+1 < h {
				if x-1 >= 0 {
					lum[y+1][x-1] += quantErr * 3.0 / 16.0
				}
				lum[y+1][x] += quantErr * 5.0 / 16.0
				if x+1 < w {
					lum[y+1][x+1] += quantErr * 1.0 / 16.0
				}
			}

			row[x] = AsciiChar{
				Char: chars[idx],
				R:    colors[y][x].R,
				G:    colors[y][x].G,
				B:    colors[y][x].B,
			}
		}
		art[y] = row
	}

	return art
}

// pixelRGB extracts 8-bit R, G, B from a color.Color.
func pixelRGB(c color.Color) (uint8, uint8, uint8) {
	r, g, b, _ := c.RGBA()
	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
}

// luminance calculates perceived brightness (0.0 to 1.0) using the
// ITU-R BT.601 luma coefficients.
func luminance(r, g, b uint8) float64 {
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255.0
}

// adjustColor applies brightness and contrast adjustments.
// Brightness is a multiplier (1.0 = no change).
// Contrast is applied around the midpoint 128.
func adjustColor(r, g, b, brightness, contrast float64) (float64, float64, float64) {
	// Apply brightness
	r *= brightness
	g *= brightness
	b *= brightness

	// Apply contrast around midpoint
	r = (r-128)*contrast + 128
	g = (g-128)*contrast + 128
	b = (b-128)*contrast + 128

	return r, g, b
}

// clampUint8 clamps a float64 to [0, 255] and returns as uint8.
func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
