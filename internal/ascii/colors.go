package ascii

import (
	"fmt"
	"math"
)

// ColorMode defines the color output mode.
type ColorMode string

const (
	ColorTrue ColorMode = "true" // 24-bit true color
	Color256  ColorMode = "256"  // 256-color ANSI
	Color16   ColorMode = "16"   // 16-color ANSI
	ColorNone ColorMode = "none" // No color (characters only)
)

// FormatCharTrue formats a character with 24-bit true color ANSI escape.
func FormatCharTrue(ch rune, r, g, b uint8) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%c", r, g, b, ch)
}

// FormatChar256 formats a character with 256-color ANSI escape.
func FormatChar256(ch rune, r, g, b uint8) string {
	idx := RGBTo256(r, g, b)
	return fmt.Sprintf("\033[38;5;%dm%c", idx, ch)
}

// FormatChar16 formats a character with 16-color ANSI escape.
func FormatChar16(ch rune, r, g, b uint8) string {
	idx := RGBTo16(r, g, b)
	if idx < 8 {
		return fmt.Sprintf("\033[%dm%c", 30+idx, ch)
	}
	return fmt.Sprintf("\033[%dm%c", 90+(idx-8), ch)
}

// FormatCharNone formats a character with no color escape.
func FormatCharNone(ch rune) string {
	return string(ch)
}

// FormatChar formats a character using the specified color mode.
func FormatChar(ch rune, r, g, b uint8, mode ColorMode) string {
	switch mode {
	case Color256:
		return FormatChar256(ch, r, g, b)
	case Color16:
		return FormatChar16(ch, r, g, b)
	case ColorNone:
		return FormatCharNone(ch)
	default:
		return FormatCharTrue(ch, r, g, b)
	}
}

// RGBTo256 maps an RGB color to the nearest 256-color palette index.
// The 256-color palette:
//
//	0-7:     standard colors
//	8-15:    bright colors
//	16-231:  6x6x6 color cube
//	232-255: grayscale ramp
func RGBTo256(r, g, b uint8) uint8 {
	// Check if it's close to a grayscale value
	if r == g && g == b {
		if r < 8 {
			return 16 // black end of color cube
		}
		if r > 248 {
			return 231 // white end of color cube
		}
		// Map to grayscale ramp (232-255, 24 shades from 8 to 238)
		return uint8(math.Round(float64(r-8)/230.0*23.0)) + 232
	}

	// Map to 6x6x6 color cube (indices 16-231)
	ri := uint8(math.Round(float64(r) / 255.0 * 5.0))
	gi := uint8(math.Round(float64(g) / 255.0 * 5.0))
	bi := uint8(math.Round(float64(b) / 255.0 * 5.0))
	return 16 + 36*ri + 6*gi + bi
}

// ansi16Colors defines the RGB values for the 16 standard ANSI colors.
var ansi16Colors = [16][3]uint8{
	{0, 0, 0},       // 0: Black
	{128, 0, 0},     // 1: Red
	{0, 128, 0},     // 2: Green
	{128, 128, 0},   // 3: Yellow
	{0, 0, 128},     // 4: Blue
	{128, 0, 128},   // 5: Magenta
	{0, 128, 128},   // 6: Cyan
	{192, 192, 192}, // 7: White
	{128, 128, 128}, // 8: Bright Black (Gray)
	{255, 0, 0},     // 9: Bright Red
	{0, 255, 0},     // 10: Bright Green
	{255, 255, 0},   // 11: Bright Yellow
	{0, 0, 255},     // 12: Bright Blue
	{255, 0, 255},   // 13: Bright Magenta
	{0, 255, 255},   // 14: Bright Cyan
	{255, 255, 255}, // 15: Bright White
}

// RGBTo16 maps an RGB color to the nearest 16-color ANSI index using
// Euclidean distance in RGB space.
func RGBTo16(r, g, b uint8) uint8 {
	bestIdx := uint8(0)
	bestDist := math.MaxFloat64

	for i, c := range ansi16Colors {
		dr := float64(r) - float64(c[0])
		dg := float64(g) - float64(c[1])
		db := float64(b) - float64(c[2])
		dist := dr*dr + dg*dg + db*db
		if dist < bestDist {
			bestDist = dist
			bestIdx = uint8(i)
		}
	}
	return bestIdx
}
