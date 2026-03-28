package ascii

import (
	"fmt"
	"os"
	"strings"
)

// PrintToTerminal renders the ASCII art to stdout with 24-bit ANSI true color.
// Each character is wrapped in an escape sequence that sets the foreground color
// to the original pixel's RGB value, producing vivid full-color output.
func PrintToTerminal(art [][]AsciiChar) {
	var sb strings.Builder
	for _, row := range art {
		for _, ch := range row {
			// \033[38;2;R;G;Bm sets 24-bit foreground color
			sb.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c", ch.R, ch.G, ch.B, ch.Char))
		}
		sb.WriteString("\033[0m\n") // reset at end of each line
	}
	fmt.Print(sb.String())
}

// SaveToText writes the ASCII art as plain text (no colors) to a file.
func SaveToText(art [][]AsciiChar, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer f.Close()

	for _, row := range art {
		for _, ch := range row {
			fmt.Fprintf(f, "%c", ch.Char)
		}
		fmt.Fprintln(f)
	}
	return nil
}

// SaveToANSI writes the ASCII art with ANSI escape codes to a file.
// The resulting file can be displayed with `cat` in a true-color terminal.
func SaveToANSI(art [][]AsciiChar, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer f.Close()

	for _, row := range art {
		for _, ch := range row {
			fmt.Fprintf(f, "\033[38;2;%d;%d;%dm%c", ch.R, ch.G, ch.B, ch.Char)
		}
		fmt.Fprintln(f, "\033[0m")
	}
	return nil
}

// SaveToHTML writes the ASCII art as an HTML file with inline color styles.
// Each character is wrapped in a <span> with its RGB color.
func SaveToHTML(art [][]AsciiChar, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}
	defer f.Close()

	fmt.Fprintln(f, `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>pixelate ASCII art</title>
<style>
  body {
    background: #0d1117;
    display: flex;
    justify-content: center;
    padding: 2rem;
    margin: 0;
  }
  pre {
    font-family: 'Cascadia Code', 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
    font-size: 8px;
    line-height: 1.0;
    letter-spacing: 0.05em;
  }
</style>
</head>
<body>
<pre>`)

	for _, row := range art {
		for _, ch := range row {
			c := ch.Char
			// Escape HTML special characters
			var escaped string
			switch c {
			case '<':
				escaped = "&lt;"
			case '>':
				escaped = "&gt;"
			case '&':
				escaped = "&amp;"
			case '"':
				escaped = "&quot;"
			case ' ':
				escaped = "&nbsp;"
			default:
				escaped = string(c)
			}
			fmt.Fprintf(f, `<span style="color:rgb(%d,%d,%d)">%s</span>`, ch.R, ch.G, ch.B, escaped)
		}
		fmt.Fprintln(f)
	}

	fmt.Fprintln(f, `</pre>
</body>
</html>`)
	return nil
}
