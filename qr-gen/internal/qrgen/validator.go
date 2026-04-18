package qrgen

import (
	"fmt"
	u "net/url"
	"regexp"
	"strings"
)

var (
	urlRegex = regexp.MustCompile(`^https?://`)

	SupportedFormats = []string{"png", "jpeg", "jpg"}
)

func ValidateParams(p QRParams) error {
	if err := ValidateURL(p.URL); err != nil {
		return err
	}
	if p.Size < 64 || p.Size > 1024 {
		return fmt.Errorf("size out of range [64-1024]: %d", p.Size)
	}
	return nil
}

// ValidateURL проверяет URL для QR (http(s), valid host, len <=4000).
func ValidateURL(url string) error {
	if len(url) == 0 {
		return fmt.Errorf("url required")
	}
	if len(url) > 4000 {
		return fmt.Errorf("url too long: %d > 4000", len(url))
	}
	if !urlRegex.MatchString(url) {
		return fmt.Errorf("url must start with http:// or https://: %q", url)
	}

	parsed, err := u.Parse(url)
	if err != nil {
		return fmt.Errorf("invalid url %q: %w", url, err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("url must have scheme and host: %q", url)
	}
	return nil
}

func ValidateFormat(format string) error {
	format = strings.ToLower(format)
	for _, f := range SupportedFormats {
		if format == f {
			return nil
		}
	}
	return fmt.Errorf("unsupported format %q, want %v", format, SupportedFormats)
}
