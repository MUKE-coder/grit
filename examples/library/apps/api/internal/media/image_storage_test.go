package media_test

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"library/apps/api/internal/media"
	"library/apps/api/internal/storage"
)

// The chain against the real local disk driver: an avatar is stored through a
// storage.Disk and served back by the driver's own file route.
func TestStoreOnTheLocalDiskIsServedBack(t *testing.T) {
	var local *storage.LocalDisk
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		local.ServeHTTP(w, r)
	}))
	defer server.Close()
	d, err := storage.NewLocalDisk(storage.LocalConfig{
		Root:      t.TempDir(),
		PublicURL: server.URL + "/files",
		Secret:    "a-test-secret-long-enough-to-sign-with",
	})
	if err != nil {
		t.Fatal(err)
	}
	local = d
	// Assigned to the interface, so this fails to compile if storage.Disk
	// ever stops satisfying media.Disk.
	var disk storage.Disk = local

	src := image.NewNRGBA(image.Rect(0, 0, 900, 600))
	for y := 0; y < 600; y++ {
		for x := 0; x < 900; x++ {
			src.SetNRGBA(x, y, color.NRGBA{uint8(x / 4), uint8(y / 3), 160, 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, src); err != nil {
		t.Fatal(err)
	}
	img, err := media.Open(&buf)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	key, err := img.Orient().Cover(150, 150).ToWebP().Store(ctx, disk, "avatars", media.PublicFile)
	if err != nil {
		t.Fatal(err)
	}
	served := get(t, disk.URL(key), http.StatusOK)
	back, err := media.Open(bytes.NewReader(served))
	if err != nil {
		t.Fatal(err)
	}
	if w, h := back.Dimensions(); w != 150 || h != 150 || back.SourceFormat() != media.WebP {
		t.Errorf("served back %dx%d %s, want 150x150 webp", w, h, back.SourceFormat())
	}

	private, err := img.Fit(300, 0).ToJPEG().Store(ctx, disk, "avatars", media.PrivateFile)
	if err != nil {
		t.Fatal(err)
	}
	get(t, disk.URL(private), http.StatusNotFound)
	signed, err := disk.TemporaryURL(ctx, private, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	get(t, signed, http.StatusOK)
}

func get(t *testing.T, url string, status int) []byte {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != status {
		t.Fatalf("GET %s: %d, want %d", url, resp.StatusCode, status)
	}
	return body
}
