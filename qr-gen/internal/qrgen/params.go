package qrgen

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

// ParseHexColor парсит HEX (#RRGGBB или RRGGBB) в color.RGBA.
func ParseHexColor(hex string) (color.RGBA, error) {
	if hex == "" {
		return color.RGBA{}, fmt.Errorf("hex required")
	}

	s := strings.TrimPrefix(hex, "#")
	s = strings.ToUpper(s)

	if len(s) != 6 {
		return color.RGBA{}, fmt.Errorf("hex should be 6 chars long: %q", s)
	}

	rgb, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return color.RGBA{}, fmt.Errorf("invalid hex %q: %w", s, err)
	}

	return color.RGBA{
		R: uint8(rgb >> 16),
		G: uint8((rgb >> 8) & 0xFF),
		B: uint8(rgb & 0xFF),
		A: 0xFF,
	}, nil
}

func ParseSize(s string) (int, error) {
	if s == "" {
		return 256, nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("invalid size %q: %w", s, err)
	}
	if i < 64 || i > 1024 {
		return 0, fmt.Errorf("size out of range [64-1024]: %d", i)
	}
	return i, nil
}

// QRParams — структура для всех параметров.
type QRParams struct {
	URL     string
	Size    int
	FGColor color.RGBA
	BGColor color.RGBA
}

// ParseParams парсит query string или CLI args в QRParams.
func ParseParams(rawURL string, rawSize, rawFG, rawBG string) (QRParams, error) {
	if rawURL == "" {
		return QRParams{}, fmt.Errorf("url required")
	}

	size, err := ParseSize(rawSize)
	if err != nil {
		return QRParams{}, err
	}

	var fgColor color.RGBA
	if rawFG != "" {
		fgColor, err = ParseHexColor(rawFG)
		if err != nil {
			return QRParams{}, fmt.Errorf("invalid fg color %q: %w", rawFG, err)
		}
	} else {
		fgColor = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	}

	var bgColor color.RGBA
	if rawBG != "" {
		bgColor, err = ParseHexColor(rawBG)
		if err != nil {
			return QRParams{}, fmt.Errorf("invalid bg color %q: %w", rawBG, err)
		}
	} else {
		bgColor = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	}

	return QRParams{
		URL:     rawURL,
		Size:    size,
		FGColor: fgColor,
		BGColor: bgColor,
	}, nil
}
