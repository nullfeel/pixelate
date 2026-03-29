package ascii

import (
	"bytes"
	"fmt"
	"image"
	// Register standard image decoders via blank imports.
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)

// LoadFromFile loads an image from a local file path.
// Supports JPEG, PNG, GIF (first frame), BMP, and WebP.
func LoadFromFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("cannot decode image: %w", err)
	}
	return img, nil
}

// LoadFromURL downloads an image from an HTTP/HTTPS URL and decodes it.
// Uses a 30-second timeout and limits response body to 50MB.
func LoadFromURL(url string) (image.Image, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		return nil, fmt.Errorf("URL does not point to an image (content-type: %s)", ct)
	}

	// Limit response body to 50MB to prevent memory exhaustion
	limited := io.LimitReader(resp.Body, 50*1024*1024)
	img, _, err := image.Decode(limited)
	if err != nil {
		return nil, fmt.Errorf("cannot decode image from URL: %w", err)
	}
	return img, nil
}

// LoadFromStdin reads an image from standard input.
func LoadFromStdin() (image.Image, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("cannot read stdin: %w", err)
	}

	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("cannot decode image from stdin: %w", err)
	}
	return img, nil
}
