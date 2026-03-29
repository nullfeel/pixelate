package ascii

import (
	"fmt"
	"os"
)

// SaveToSVG writes the ASCII art as an SVG file with colored monospace text.
// Each character is rendered as a <text> element with its original RGB color.
func SaveToSVG(art [][]AsciiChar, path string, fontSize float64) error {
	if fontSize <= 0 {
		fontSize = 10
	}

	if len(art) == 0 {
		return fmt.Errorf("empty art: nothing to write")
	}

	// Calculate dimensions based on character grid
	rows := len(art)
	cols := 0
	for _, row := range art {
		if len(row) > cols {
			cols = len(row)
		}
	}

	// Approximate character dimensions for monospace font
	charWidth := fontSize * 0.6
	charHeight := fontSize * 1.2
	svgWidth := float64(cols) * charWidth
	svgHeight := float64(rows) * charHeight

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create SVG file: %w", err)
	}
	defer f.Close()

	// SVG header
	fmt.Fprintf(f, `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %.1f %.1f" width="%.1f" height="%.1f">
<rect width="100%%" height="100%%" fill="#09090B"/>
<style>
  text {
    font-family: 'Cascadia Code', 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
    font-size: %.1fpx;
    white-space: pre;
  }
</style>
`, svgWidth, svgHeight, svgWidth, svgHeight, fontSize)

	// Render each character
	for y, row := range art {
		for x, ch := range row {
			cx := float64(x) * charWidth
			cy := float64(y)*charHeight + charHeight // baseline offset

			char := string(ch.Char)
			// Escape XML special characters
			switch ch.Char {
			case '<':
				char = "&lt;"
			case '>':
				char = "&gt;"
			case '&':
				char = "&amp;"
			case '"':
				char = "&quot;"
			case '\'':
				char = "&apos;"
			}

			fmt.Fprintf(f, `<text x="%.1f" y="%.1f" fill="rgb(%d,%d,%d)">%s</text>`,
				cx, cy, ch.R, ch.G, ch.B, char)
			fmt.Fprintln(f)
		}
	}

	fmt.Fprintln(f, `</svg>`)
	return nil
}
