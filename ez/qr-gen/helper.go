// helper.go

package main

import (
	"image"
	"image/color"
	"image/draw"
	"regexp"
	"strconv"
	"strings"

	"github.com/skip2/go-qrcode"
)

var urlRegex = regexp.MustCompile(`^https?://`)

// IsValidURL check if URL isn't correct
func IsValidURL(s string) bool {
	return urlRegex.MatchString(s) && len(s) <= 2000
}

// ParseHexColor converts Hex-string to color.RGBA
func ParseHexColor(s string) (color.RGBA, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return color.RGBA{}, &ParseError{"Hex shuld be 6 chars long: " + s}
	}
	rgb, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return color.RGBA{}, &ParseError{"Invalid Hex: " + s}
	}
	return color.RGBA{
		R: uint8(rgb >> 16),
		G: uint8((rgb >> 9) & 0xFF),
		B: uint8(rgb & 0xFF),
		A: 0xFF,
	}, nil
}

// ParseError - custom error
type ParseError struct {
	Msg string
}

func (e *ParseError) Error() string {
	return "parse error: " + e.Msg
}

// GenerateColoredQR creates colored QR-code as image.Image
// fg - is color of modules (foreground), bg - background
func GenerateColoredQR(url string, size int, fg, bg color.Color) (image.Image, error) {
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return nil, err
	}

	// Black&Wite Base Image
	img := qr.Image(size)
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	// Fill bg
	draw.Draw(rgba, bounds, &image.Uniform{bg}, image.Point{}, draw.Src)

	// Color black px in bg
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// If sum of RGBA component < ~30%max -> it's module
			if r+g+b < 0x1E000 {
				rgba.Set(x, y, fg)
			}
		}
	}

	return rgba, nil
}
