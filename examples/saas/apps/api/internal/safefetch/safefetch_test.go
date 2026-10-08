package safefetch

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestBlocksPrivateHosts(t *testing.T) {
	cases := []string{
		"http://127.0.0.1/",
		"http://localhost/",
		"http://10.0.0.1/",
		"http://192.168.1.1/",
		"http://169.254.169.254/latest/meta-data/",
		"http://[::1]/",
		"http://metadata.google.internal/computeMetadata/v1/",
		"http://100.64.0.5/",
	}
	for _, raw := range cases {
		_, err := Get(context.Background(), raw)
		if !errors.Is(err, ErrBlocked) {
			t.Errorf("Get(%q): expected ErrBlocked, got %v", raw, err)
		}
	}
}

func TestBlocksUnknownSchemes(t *testing.T) {
	cases := []string{
		"file:///etc/passwd",
		"gopher://evil.example.com/x",
		"ftp://example.com/",
	}
	for _, raw := range cases {
		_, err := Get(context.Background(), raw)
		if err == nil || !strings.Contains(err.Error(), "scheme") {
			t.Errorf("Get(%q): expected scheme rejection, got %v", raw, err)
		}
	}
}
