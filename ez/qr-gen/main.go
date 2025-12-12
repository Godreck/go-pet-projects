package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func handleQR(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "?url = required", http.StatusBadRequest)
		return
	}
	if !IsValidURL(url) {
		http.Error(w, "URL shuld start with http:// or https://", http.StatusBadRequest)
		return
	}

	size := getIntParam(r.URL.Query().Get("size"), 256, 64, 1024)
	fgColor := getColorParam(r.URL.Query().Get("fg"), color.Black)
	bgColor := getColorParam(r.URL.Query().Get("bg"), color.White)

	img, err := GenerateColoredQR(url, size, fgColor, bgColor)
	if err != nil {
		http.Error(w, "QR generation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	format := strings.ToLower(r.URL.Query().Get("format"))
	w.Header().Set("Content-Type", "image/png")
	if format == "jpeg" || format == "jpg" {
		w.Header().Set("Content-Type", "image/jpeg")
		if err := jpeg.Encode(w, img, &jpeg.Options{Quality: 90}); err != nil {
			http.Error(w, "JPEG encode: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := png.Encode(w, img); err != nil {
		http.Error(w, "PNG encode: "+err.Error(), http.StatusInternalServerError)
	}
}

func getIntParam(s string, def, min, max int) int {
	if s == "" {
		return def
	}
	if v, err := strconv.Atoi(s); err == nil && v >= min && v <= max {
		return v
	}
	return def
}

func getColorParam(hex string, def color.Color) color.Color {
	if hex == "" {
		return def
	}
	if c, err := ParseHexColor(hex); err == nil {
		return c
	}
	return def
}

func main() {
	const usage = `qr -- for generating QR-codes (CLI and HTTP)

	CLI:
		qr -url https://exapmle.com -out qr.png -fg 3498db -bg ffffff
	
	HTTP:
		qr -http
		-> http://localhost:8080/qr?url=...&fg=...&bg=...&format=png|jpg

	`

	cli := flag.NewFlagSet("qr", flag.ExitOnError)
	cli.Usage = func() { fmt.Print(usage) }

	url := cli.String("url", "", "URL for transfoming in QR-code (required)")
	out := cli.String("out", "qrcode.png", "Name of output file (by default: qrcode.png)")
	size := cli.Int("size", 256, "QR-code sieze in pixels (range: 65-1024, by default: 256)")
	fg := cli.String("fg", "00000", "Module's color (Hex without '#')")
	bg := cli.String("bg", "fffff", "Background (Hex without '#')")
	httpMode := cli.Bool("http", false, "Lounch HTTP-server")

	cli.Parse(os.Args[1:])

	if *httpMode {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Follow this link: http://localhost:%s/qr?url=...", port)
		http.HandleFunc("/qr", handleQR)
		// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 	http.Redirect(w, r, "/qr?url=https://github.com/Godreck/go-pet-projects")
		// })
		log.Fatal(http.ListenAndServe(":"+port, nil))
	}

	if *url == "" {
		fmt.Fprintln(os.Stderr, "❌ Error: flag -url is requaired")
		cli.Usage()
		os.Exit(1)
	}
	if !IsValidURL(*url) {
		log.Fatal("❌ URL shuld start with http:// or https://")
	}

	fgColor, err := ParseHexColor(*fg)
	if err != nil {
		log.Fatalf("❌ fg HEX: %v", err)
	}
	bgColor, err := ParseHexColor(*bg)
	if err != nil {
		log.Fatalf("❌ bg HEX: %v", err)
	}

	img, err := GenerateColoredQR(*url, *size, fgColor, bgColor)
	if err != nil {
		log.Fatalf("❌ Generating: %v", err)
	}

	file, err := os.Create(*out)
	if err != nil {
		log.Fatalf("❌ Craeting file %s: %v", *out, err)
	}
	defer file.Close()

	var encodeErr error
	if strings.HasSuffix(strings.ToLower(*out), ".jpg") ||
		strings.HasSuffix(strings.ToLower(*out), ".jpeg") {
		encodeErr = jpeg.Encode(file, img, &jpeg.Options{Quality: 95})
	} else {
		encodeErr = png.Encode(file, img)
	}

	if encodeErr != nil {
		log.Fatalf("❌ Write: %v", encodeErr)
	}

	fmt.Printf("✅ QR сохранён: %s\n", *out)
	fmt.Printf("   URL: %s | Размер: %d | fg=#%s bg=#%s\n", *url, *size, *fg, *bg)

}
