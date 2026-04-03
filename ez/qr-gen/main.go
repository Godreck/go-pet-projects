package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	qrgen "qr-gen/internal/qrgen"
)

func HandleQR(w http.ResponseWriter, r *http.Request) {
	// Получаем params
	query := r.URL.Query()
	url := query.Get("url")
	sizeStr := query.Get("size")
	fgHex := query.Get("fg")
	bgHex := query.Get("bg")
	format := query.Get("format")

	// ParseParams, парсим и валидируем
	params, err := qrgen.ParseParams(url, sizeStr, fgHex, bgHex)
	if err != nil {
		http.Error(w, fmt.Sprintf("invalid params: %v", err), http.StatusBadRequest)
		return
	}

	img, err := qrgen.GenerateQR(params.URL, params.Size, params.FGColor, params.BGColor)
	if err != nil {
		http.Error(w, fmt.Sprintf("generate QR: %v", err), http.StatusInternalServerError)
		return
	}

	format = strings.ToLower(format)
	contentType := "image/png"
	if format == "jpeg" || format == "jpg" {
		contentType = "image/jpeg"
	}
	w.Header().Set("Content-Type", contentType)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		http.Error(w, fmt.Sprintf("encoding png: %v", err), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	switch contentType {
	case "image/jpeg":
		if err := jpeg.Encode(w, img, &jpeg.Options{Quality: 90}); err != nil {
			http.Error(w, "jpeg encode: "+err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		if err := png.Encode(w, img); err != nil {
			http.Error(w, "png encode: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RunCLI(cli *flag.FlagSet) {
	// Parse flags
	url := cli.String("url", "", "URL (required)")
	out := cli.String("out", "qrcode.png", "output file")
	sizeStr := cli.String("size", "", "size (64-1024)")
	fg := cli.String("fg", "", "fg color HEX")
	bg := cli.String("bg", "", "bg color HEX")

	cli.Parse(os.Args[1:])

	// ParseParams
	params, err := qrgen.ParseParams(*url, *sizeStr, *fg, *bg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ParseParams error: %v\n", err)
		cli.Usage()
		os.Exit(1)
	}

	// GenerateQRToFile
	if err := qrgen.GenerateQRToFile(params.URL, *out, params.Size, params.FGColor, params.BGColor, "auto"); err != nil {
		log.Fatalf("generate error: %v", err)
	}

	// Success
	fmt.Printf("QR сохранён: %s\n", *out)
	fmt.Printf("   URL: %s | Size: %d | fg=#%02x%02x%02x bg=#%02x%02x%02x\n",
		params.URL, params.Size,
		params.FGColor.R, params.FGColor.G, params.FGColor.B,
		params.BGColor.R, params.BGColor.G, params.BGColor.B)
}

func main() {
	usage := `qr — QR code generator (CLI + HTTP)

CLI:
    qr -url https://example.com -out qr.png -fg 3498db -bg ffffff

HTTP:
    qr -http
    -> http://localhost:8080/?url=...&fg=...&format=png|jpg
    `

	cli := flag.NewFlagSet("qr", flag.ExitOnError)
	cli.Usage = func() { fmt.Print(usage) }

	httpMode := cli.Bool("http", false, "start HTTP server")
	cli.Parse(os.Args[1:])

	// if *httpMode {
	// 	// HTTP mode
	// 	port := os.Getenv("PORT")
	// 	if port == "" {
	// 		port = "8080"
	// 	}
	// 	log.Printf("Server: http://localhost:%s/?url=...", port)
	// 	http.HandleFunc("/", handleQR) // root = /qr
	// 	log.Fatal(http.ListenAndServe(":"+port, nil))
	// }
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if *httpMode {
		srv := &http.Server{Addr: ":" + port, Handler: nil}
		log.Printf("🚀 http://localhost:%s/?url=...", port)
		http.HandleFunc("/", HandleQR)
		log.Fatal(srv.ListenAndServe())
	}

	// CLI mode
	RunCLI(cli)
}
