package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/nullfeel/pixelate/internal/ascii"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	width      int
	grayscale  bool
	invert     bool
	output     string
	charset    string
	preset     string
	brightness float64
	contrast   float64
	dither     bool
)

var rootCmd = &cobra.Command{
	Use:   "pixelate [flags] <image> [image...]",
	Short: "Convert images to colored ASCII art",
	Long: `pixelate converts images to stunning ASCII art with true color support.

Supports JPG, PNG, GIF, BMP, and WebP formats.
Accepts local files, URLs, or piped stdin.

Examples:
  pixelate photo.jpg
  pixelate -w 120 photo.png
  pixelate --grayscale photo.jpg
  pixelate --preset block photo.jpg
  pixelate --preset braille photo.jpg
  pixelate -o art.html photo.jpg
  pixelate https://example.com/photo.jpg
  cat photo.jpg | pixelate -`,
	Args: cobra.MinimumNArgs(0),
	RunE: runPixelate,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().IntVarP(&width, "width", "w", 0, "Output width in characters (default: terminal width)")
	rootCmd.Flags().BoolVarP(&grayscale, "grayscale", "g", false, "Output in grayscale (no colors)")
	rootCmd.Flags().BoolVar(&invert, "invert", false, "Invert light/dark characters")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Save output to file (.txt, .ans, .html)")
	rootCmd.Flags().StringVar(&charset, "charset", "", "Custom character set (light to dark)")
	rootCmd.Flags().StringVar(&preset, "preset", "simple", "Character preset: simple, block, braille, detailed")
	rootCmd.Flags().Float64Var(&brightness, "brightness", 1.0, "Brightness multiplier (e.g., 1.2)")
	rootCmd.Flags().Float64Var(&contrast, "contrast", 1.0, "Contrast multiplier (e.g., 1.5)")
	rootCmd.Flags().BoolVar(&dither, "dither", false, "Apply Floyd-Steinberg dithering (grayscale)")
}

func getTerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

func runPixelate(cmd *cobra.Command, args []string) error {
	// Determine target width
	targetWidth := width
	if targetWidth <= 0 {
		targetWidth = getTerminalWidth()
	}

	// Resolve character set
	chars := charset
	if chars == "" {
		chars = ascii.GetPreset(preset)
	}

	if invert {
		chars = reverseString(chars)
	}

	opts := ascii.ConvertOptions{
		Width:      targetWidth,
		Charset:    chars,
		Grayscale:  grayscale,
		Brightness: brightness,
		Contrast:   contrast,
		Dither:     dither,
	}

	// Read from stdin if no args or arg is "-"
	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		info, _ := os.Stdin.Stat()
		if info.Mode()&os.ModeCharDevice != 0 {
			return fmt.Errorf("no input provided. Usage: pixelate <image>")
		}
		return processStdin(opts)
	}

	// Process each argument
	for i, arg := range args {
		if i > 0 {
			fmt.Println()
		}
		if err := processArg(arg, opts); err != nil {
			fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", arg, err)
		}
	}

	return nil
}

func processArg(arg string, opts ascii.ConvertOptions) error {
	if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
		return processURL(arg, opts)
	}
	return processFile(arg, opts)
}

func processFile(path string, opts ascii.ConvertOptions) error {
	// Check for animated GIF
	if strings.HasSuffix(strings.ToLower(path), ".gif") {
		isAnimated, err := ascii.IsAnimatedGIF(path)
		if err == nil && isAnimated && output == "" {
			return ascii.PlayAnimatedGIF(path, opts)
		}
	}

	img, err := ascii.LoadFromFile(path)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	art := ascii.Convert(img, opts)
	return outputResult(art)
}

func processURL(url string, opts ascii.ConvertOptions) error {
	img, err := ascii.LoadFromURL(url)
	if err != nil {
		return fmt.Errorf("failed to load image from URL: %w", err)
	}

	art := ascii.Convert(img, opts)
	return outputResult(art)
}

func processStdin(opts ascii.ConvertOptions) error {
	img, err := ascii.LoadFromStdin()
	if err != nil {
		return fmt.Errorf("failed to load image from stdin: %w", err)
	}

	art := ascii.Convert(img, opts)
	return outputResult(art)
}

func outputResult(art [][]ascii.AsciiChar) error {
	if output == "" {
		ascii.PrintToTerminal(art)
		return nil
	}

	lower := strings.ToLower(output)
	switch {
	case strings.HasSuffix(lower, ".html"):
		return ascii.SaveToHTML(art, output)
	case strings.HasSuffix(lower, ".ans"):
		return ascii.SaveToANSI(art, output)
	default:
		return ascii.SaveToText(art, output)
	}
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
