package qrgen

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		name    string
		hex     string
		want    color.RGBA
		wantErr bool
	}{
		{"valid 6 hex", "ff0000", color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, false},
		{"valid with #", "#00ff00", color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, false},
		{"short 3 hex", "fff", color.RGBA{}, true}, // fail: len!=6
		{"invalid chars", "zz0000", color.RGBA{}, true},
		{"empty", "", color.RGBA{}, true},
		{"too long", "fffffffff", color.RGBA{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHexColor(tt.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHexColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseHexColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    int
		wantErr bool
	}{
		{"default", "", 256, false},
		{"valid", "512", 512, false},
		{"min", "64", 64, false},
		{"max", "1024", 1024, false},
		{"too small", "32", 0, true},
		{"too big", "2048", 0, true},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSize(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseSize() got = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestParseParams(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		size    string
		fg      string
		bg      string
		wantErr bool
	}{
		{"happy path", "https://example.com", "256", "ff0000", "ffffff", false},
		{"no url", "", "256", "000000", "ffffff", true},
		{"invalid size", "https://example.com", "abc", "000000", "ffffff", true},
		{"invalid fg", "https://example.com", "256", "zz0000", "ffffff", true},
		{"defaults", "https://example.com", "", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseParams(tt.url, tt.size, tt.fg, tt.bg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid", "https://example.com", false},
		{"http", "http://localhost:8080", false},
		{"empty", "", true},
		{"no scheme", "example.com", true},
		{"too long", strings.Repeat("a", 4001), true},
		{"invalid", "://invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestValidateFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"png", "png", false},
		{"jpeg", "jpeg", false},
		{"jpg", "jpg", false},
		{"invalid", "gif", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFormat(%q) error = %v, wantErr %v", tt.format, err, tt.wantErr)
			}
		})
	}
}

func TestGenerateQR(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		size    int
		wantErr bool
		wantW   int
		wantH   int
	}{
		{"happy path", "https://example.com", 256, false, 256, 256},
		{"default size", "https://go.dev", 0, false, 256, 256},
		{"invalid url", "", 256, true, 0, 0},
		{"min size", "https://a.co", 64, false, 64, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fg := color.RGBA{0x00, 0x00, 0x00, 0xFF}
			bg := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

			img, err := GenerateQR(tt.url, tt.size, fg, bg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateQR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				b := img.Bounds()
				if b.Dx() != tt.wantW || b.Dy() != tt.wantH {
					t.Errorf("GenerateQR() size = %dx%d, want %dx%d", b.Dx(), b.Dy(), tt.wantW, tt.wantH)
				}
			}
		})
	}
}

func TestGenerateQRToFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.png")

	fg := color.RGBA{0xFF, 0x00, 0x00, 0xFF} // red
	bg := color.RGBA{0xFF, 0xFF, 0xFF, 0xFF} // white

	err := GenerateQRToFile("https://example.com", tmpFile, 128, fg, bg, "png")
	if err != nil {
		t.Fatalf("GenerateQRToFile() error = %v", err)
	}

	// Проверяем, что файл создался
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Errorf("file %q not created", tmpFile)
	}

	// Проверяем размер (примерно)
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("stat %q: %v", tmpFile, err)
	}
	if info.Size() < 500 {
		t.Errorf("file too small: %d bytes", info.Size())
	}
}

func TestValidateParams(t *testing.T) {
	tests := []struct {
		name    string
		params  QRParams
		wantErr bool
	}{
		{"valid", QRParams{URL: "https://example.com", Size: 256}, false},
		{"invalid url", QRParams{URL: "", Size: 256}, true},
		{"invalid size", QRParams{URL: "https://a.co", Size: 32}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateParams(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
