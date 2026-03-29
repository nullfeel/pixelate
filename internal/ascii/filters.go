package ascii

import (
	"image"
	"image/color"
	"math"
)

// FilterType defines the available image filter types.
type FilterType string

const (
	FilterNone      FilterType = "none"
	FilterEdge      FilterType = "edge"
	FilterNegative  FilterType = "negative"
	FilterSepia     FilterType = "sepia"
	FilterBlur      FilterType = "blur"
	FilterSharpen   FilterType = "sharpen"
	FilterPixelate  FilterType = "pixelate"
	FilterGrayscale FilterType = "grayscale"
)

// ApplyFilter applies the specified filter to the image with the given intensity.
// Returns a new image with the filter applied.
func ApplyFilter(img image.Image, filter FilterType, intensity float64) image.Image {
	switch filter {
	case FilterEdge:
		return applyEdge(img)
	case FilterNegative:
		return applyNegative(img)
	case FilterSepia:
		return applySepia(img)
	case FilterBlur:
		radius := int(intensity)
		if radius < 1 {
			radius = 2
		}
		return applyBlur(img, radius)
	case FilterSharpen:
		return applySharpen(img)
	case FilterPixelate:
		blockSize := int(intensity)
		if blockSize < 2 {
			blockSize = 8
		}
		return applyPixelateFilter(img, blockSize)
	case FilterGrayscale:
		return applyGrayscale(img)
	default:
		return img
	}
}

// applyEdge applies Sobel edge detection using 3x3 kernels.
func applyEdge(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))

	// Sobel kernels
	gx := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	gy := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sumXr, sumXg, sumXb float64
			var sumYr, sumYg, sumYb float64

			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := clampInt(x+kx, 0, w-1) + bounds.Min.X
					py := clampInt(y+ky, 0, h-1) + bounds.Min.Y
					r, g, b := pixelRGB(img.At(px, py))
					rf, gf, bf := float64(r), float64(g), float64(b)

					sumXr += rf * gx[ky+1][kx+1]
					sumXg += gf * gx[ky+1][kx+1]
					sumXb += bf * gx[ky+1][kx+1]

					sumYr += rf * gy[ky+1][kx+1]
					sumYg += gf * gy[ky+1][kx+1]
					sumYb += bf * gy[ky+1][kx+1]
				}
			}

			rr := clampUint8(math.Sqrt(sumXr*sumXr + sumYr*sumYr))
			gg := clampUint8(math.Sqrt(sumXg*sumXg + sumYg*sumYg))
			bb := clampUint8(math.Sqrt(sumXb*sumXb + sumYb*sumYb))

			out.SetNRGBA(x, y, color.NRGBA{R: rr, G: gg, B: bb, A: 255})
		}
	}
	return out
}

// applyNegative inverts all RGB channels.
func applyNegative(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b := pixelRGB(img.At(x+bounds.Min.X, y+bounds.Min.Y))
			out.SetNRGBA(x, y, color.NRGBA{R: 255 - r, G: 255 - g, B: 255 - b, A: 255})
		}
	}
	return out
}

// applySepia applies a classic warm brownish tone.
func applySepia(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b := pixelRGB(img.At(x+bounds.Min.X, y+bounds.Min.Y))
			rf, gf, bf := float64(r), float64(g), float64(b)

			nr := clampUint8(rf*0.393 + gf*0.769 + bf*0.189)
			ng := clampUint8(rf*0.349 + gf*0.686 + bf*0.168)
			nb := clampUint8(rf*0.272 + gf*0.534 + bf*0.131)

			out.SetNRGBA(x, y, color.NRGBA{R: nr, G: ng, B: nb, A: 255})
		}
	}
	return out
}

// applyBlur applies a box blur with the given radius.
func applyBlur(img image.Image, radius int) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sumR, sumG, sumB float64
			var count float64

			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					px := clampInt(x+kx, 0, w-1) + bounds.Min.X
					py := clampInt(y+ky, 0, h-1) + bounds.Min.Y
					r, g, b := pixelRGB(img.At(px, py))
					sumR += float64(r)
					sumG += float64(g)
					sumB += float64(b)
					count++
				}
			}

			out.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8(sumR / count),
				G: clampUint8(sumG / count),
				B: clampUint8(sumB / count),
				A: 255,
			})
		}
	}
	return out
}

// applySharpen applies an unsharp mask (sharpen kernel convolution).
func applySharpen(img image.Image) image.Image {
	kernel := [3][3]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
	return applyConvolution(img, kernel)
}

// applyPixelateFilter applies a block pixelation effect.
func applyPixelateFilter(img image.Image, blockSize int) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))

	for by := 0; by < h; by += blockSize {
		for bx := 0; bx < w; bx += blockSize {
			var sumR, sumG, sumB float64
			var count float64

			// Calculate the average color for this block
			endY := by + blockSize
			if endY > h {
				endY = h
			}
			endX := bx + blockSize
			if endX > w {
				endX = w
			}

			for y := by; y < endY; y++ {
				for x := bx; x < endX; x++ {
					r, g, b := pixelRGB(img.At(x+bounds.Min.X, y+bounds.Min.Y))
					sumR += float64(r)
					sumG += float64(g)
					sumB += float64(b)
					count++
				}
			}

			avgR := clampUint8(sumR / count)
			avgG := clampUint8(sumG / count)
			avgB := clampUint8(sumB / count)

			// Fill the block with the average color
			for y := by; y < endY; y++ {
				for x := bx; x < endX; x++ {
					out.SetNRGBA(x, y, color.NRGBA{R: avgR, G: avgG, B: avgB, A: 255})
				}
			}
		}
	}
	return out
}

// applyGrayscale converts to grayscale while keeping as RGB image.
func applyGrayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b := pixelRGB(img.At(x+bounds.Min.X, y+bounds.Min.Y))
			gray := clampUint8(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
			out.SetNRGBA(x, y, color.NRGBA{R: gray, G: gray, B: gray, A: 255})
		}
	}
	return out
}

// applyConvolution applies a 3x3 convolution kernel to the image.
func applyConvolution(img image.Image, kernel [3][3]float64) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	out := image.NewNRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			var sumR, sumG, sumB float64

			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := clampInt(x+kx, 0, w-1) + bounds.Min.X
					py := clampInt(y+ky, 0, h-1) + bounds.Min.Y
					r, g, b := pixelRGB(img.At(px, py))

					k := kernel[ky+1][kx+1]
					sumR += float64(r) * k
					sumG += float64(g) * k
					sumB += float64(b) * k
				}
			}

			out.SetNRGBA(x, y, color.NRGBA{
				R: clampUint8(sumR),
				G: clampUint8(sumG),
				B: clampUint8(sumB),
				A: 255,
			})
		}
	}
	return out
}

// clampInt clamps an integer to [min, max].
func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
