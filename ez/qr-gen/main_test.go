package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_HappyPath(t *testing.T) {
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "test.png")

	args := []string{
		"-url", "https://example.com",
		"-out", outFile,
		"-size", "128",
		"-fg", "ff0000", // red
		"-bg", "ffff00", // yellow
	}

	cmd := exec.Command(os.Args[0], args...) // go test ./...
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")

	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		t.Fatalf("CLI failed: %v\nstderr: %s", err, stderr.String())
	}

	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Fatalf("output file %q not created", outFile)
	}

	got := stdout.String()
	want := "QR сохранён: test.png"
	if !strings.Contains(got, want) {
		t.Errorf("stdout want %q, got %q", want, got)
	}
}

func TestCLI_NoURL(t *testing.T) {
	args := []string{"-out", "qr.png"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")

	output, _ := cmd.CombinedOutput()
	wantErr := "url required"
	if !strings.Contains(string(output), wantErr) {
		t.Errorf("want error %q, got %q", wantErr, string(output))
	}
}

func TestHTTP_HappyPath(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?url=https://example.com&size=128&fg=ff0000&format=png", nil)
	rr := httptest.NewRecorder()

	HandleQR(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: %d want %d", status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		t.Errorf("Content-Type = %q, want image/*", contentType)
	}

	body, _ := io.ReadAll(rr.Body)
	if len(body) < 1000 {
		t.Errorf("response too small: %d bytes", len(body))
	}
}

func TestHTTP_NoURL(t *testing.T) {
	req, _ := http.NewRequest("GET", "/?size=256", nil)
	rr := httptest.NewRecorder()

	HandleQR(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("status %d, want %d", status, http.StatusBadRequest)
	}

	body, _ := io.ReadAll(rr.Body)
	if !strings.Contains(string(body), "invalid params") {
		t.Errorf("want 'invalid params', got %q", string(body))
	}
}

func TestHTTP_InvalidParams(t *testing.T) {
	tests := []string{
		"/?url=invalid", // no scheme
		"/?size=abc",    // invalid size
		"/?fg=zz0000",   // invalid hex
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req, _ := http.NewRequest("GET", path, nil)
			rr := httptest.NewRecorder()

			HandleQR(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("%q status %d, want %d", path, status, http.StatusBadRequest)
			}
		})
	}
}
