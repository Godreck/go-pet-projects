package qrgen

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQR генерирует QR-код как image.Image.
func GenerateQR(url string, size int, fgColor, bgColor color.RGBA) (image.Image, error) {
	if size <= 0 || size > 1024 {
		size = 256
	}

	params := QRParams{
		URL:     url,
		Size:    size,
		FGColor: fgColor,
		BGColor: bgColor,
	}

	if err := ValidateParams(params); err != nil { // новая func
		return nil, err
	}

	qr, err := qrcode.New(params.URL, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("qrcode.New: %w", err)
	}

	qrImg := qr.Image(params.Size)
	return recolor(qrImg, qr, params.FGColor, params.BGColor), nil
}

func recolor(src image.Image, qr *qrcode.QRCode, fg, bg color.RGBA) image.Image {
	b := src.Bounds()
	dst := image.NewRGBA(b)

	draw.Draw(dst, b, image.NewUniform(bg), image.Point{}, draw.Src)

	bitmap := qr.Bitmap()
	modules := len(bitmap)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			mx, my := scaleCoord(x, y, b, modules)
			if bitmap[my][mx] {
				dst.SetRGBA(x, y, fg)
			}
		}
	}
	return dst
}

func scaleCoord(x, y int, b image.Rectangle, modules int) (int, int) {
	moduleSizeX := float64(b.Dx()) / float64(modules)
	moduleSizeY := float64(b.Dy()) / float64(modules)
	return clamp(int(float64(x)/moduleSizeX), 0, modules-1),
		clamp(int(float64(y)/moduleSizeY), 0, modules-1)
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func detectFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".jpg" || ext == ".jpeg" {
		return "jpeg"
	}
	return "png"
}

// GenerateQRToFile сохраняет QR в файл (PNG/JPEG).
func GenerateQRToFile(url, filename string, size int, fgColor, bgColor color.RGBA, format string) error {
	if format == "auto" {
		format = detectFormat(filename)
	}

	if err := ValidateFormat(format); err != nil {
		return fmt.Errorf("format: %w", err)
	}

	img, err := GenerateQR(url, size, fgColor, bgColor)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	format = strings.ToLower(format)
	switch format {
	case "jpeg", "jpg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: 95})
	default:
		return png.Encode(f, img)
	}
}
