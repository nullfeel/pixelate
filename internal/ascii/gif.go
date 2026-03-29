package ascii

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
	"os/signal"
	"strings"
	"time"
)

// IsAnimatedGIF checks whether a GIF file contains multiple frames.
func IsAnimatedGIF(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return false, err
	}
	return len(g.Image) > 1, nil
}

// PlayAnimatedGIF decodes all frames of a GIF, converts each to ASCII art,
// and plays them in the terminal with the original frame delays.
// It loops indefinitely until interrupted with Ctrl+C.
func PlayAnimatedGIF(path string, opts ConvertOptions) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open GIF: %w", err)
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		return fmt.Errorf("cannot decode GIF: %w", err)
	}

	if len(g.Image) == 0 {
		return fmt.Errorf("GIF contains no frames")
	}

	// Cap frames to prevent memory exhaustion
	maxFrames := 500
	frameCount := len(g.Image)
	if frameCount > maxFrames {
		frameCount = maxFrames
	}

	// Pre-render all frames to ASCII
	frames := make([][][]AsciiChar, frameCount)

	// Build a composite canvas to handle GIF disposal methods properly.
	// Each frame may only cover a sub-rectangle of the full image.
	if g.Config.Width <= 0 || g.Config.Height <= 0 {
		return fmt.Errorf("GIF has invalid dimensions: %dx%d", g.Config.Width, g.Config.Height)
	}
	bounds := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	canvas := image.NewPaletted(bounds, palette.Plan9)

	for i := 0; i < frameCount; i++ {
		frame := g.Image[i]
		// Validate frame bounds are within canvas
		fb := frame.Bounds()
		if !fb.In(image.Rect(0, 0, g.Config.Width, g.Config.Height)) {
			fb = fb.Intersect(bounds)
		}

		// Draw the frame onto the canvas
		draw.Draw(canvas, fb, frame, frame.Bounds().Min, draw.Over)

		// Convert the full canvas to a regular RGBA image for consistent processing
		rgba := image.NewRGBA(bounds)
		draw.Draw(rgba, bounds, canvas, bounds.Min, draw.Src)

		frames[i] = Convert(rgba, opts)
	}

	// Handle Ctrl+C gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	// Hide cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h\033[0m") // show cursor and reset colors on exit

	loopCount := g.LoopCount // 0 means infinite
	iteration := 0

	for {
		for i, frame := range frames {
			select {
			case <-sigCh:
				return nil
			default:
			}

			// Move cursor to top-left
			fmt.Print("\033[H")

			// Render frame
			var sb strings.Builder
			for _, row := range frame {
				for _, ch := range row {
					sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c", ch.R, ch.G, ch.B, ch.Char))
				}
				sb.WriteString("\033[0m\n")
			}
			fmt.Print(sb.String())

			// Delay: GIF delay is in 100ths of a second
			delay := 10 // default ~100ms
			if i < len(g.Delay) && g.Delay[i] > 0 {
				delay = g.Delay[i]
			}
			time.Sleep(time.Duration(delay) * 10 * time.Millisecond)
		}

		iteration++
		if loopCount > 0 && iteration >= loopCount {
			break
		}
	}

	return nil
}
