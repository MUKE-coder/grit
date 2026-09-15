package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The chainable image API ships with the media package, which writeMediaFiles
// delivers on both grit new and grit upgrade, so no repair exists for it.
func TestMediaChainIsWrittenWithTheMediaPackage(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchAPI}
	if err := writeMediaFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(opts.APIRoot(root), "internal", "media")
	for _, name := range []string{"image.go", "image_load.go", "image_purego.go", "image_vips.go", "image_test.go", "image_storage_test.go"} {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Errorf("%s was not written: %v", name, err)
			continue
		}
		if strings.Contains(string(raw), "{{MODULE}}") {
			t.Errorf("%s still has a {{MODULE}} placeholder", name)
		}
		mustFormatGo(t, name, string(raw))
	}
}

func TestMediaChainTemplates(t *testing.T) {
	image := mediaImageGo()
	// Decode, every pixel operation, the dominant colour and the encode each
	// take a transform slot.
	if n := strings.Count(image, "transformSlots <- struct{}{}"); n < 4 {
		t.Errorf("expected the decode, apply, DominantColor and Encode to take a slot, found %d", n)
	}
	for _, want := range []string{
		"func Open(r io.Reader, opts ...Option) (*Image, error)",
		"checkPixels(cfg.Width, cfg.Height, o.maxPixels)",
		"func (i *Image) Cover(width, height int) *Image",
		"func (i *Image) Resize(width, height int) *Image",
		"func (i *Image) Fit(width, height int) *Image",
		"func (i *Image) Crop(x, y, width, height int) *Image",
		"func (i *Image) Rotate(degrees float64, background color.Color) *Image",
		"func (i *Image) FlipH() *Image",
		"func (i *Image) FlipV() *Image",
		"func (i *Image) Grayscale() *Image",
		"func (i *Image) Blur(amount float64) *Image",
		"func (i *Image) Sharpen(amount float64) *Image",
		"func (i *Image) Orient() *Image",
		"func (i *Image) DominantColor() (color.RGBA, error)",
		"func (i *Image) ToAVIF() *Image",
		"func (i *Image) ToLossyWebP() *Image",
		"ErrNeedsVips",
		"func (i *Image) Encode() ([]byte, error)",
	} {
		if !strings.Contains(image, want) {
			t.Errorf("image.go is missing %q", want)
		}
	}

	load := mediaImageLoadGo()
	for _, want := range []string{
		"safefetch.Get(ctx, rawURL)",
		"context.WithTimeout(ctx, o.timeout)",
		"readLimited(body, o.maxBytes)",
		"func (i *Image) Store(ctx context.Context, disk Disk, dir string, visibility Visibility) (string, error)",
	} {
		if !strings.Contains(load, want) {
			t.Errorf("image_load.go is missing %q", want)
		}
	}
	// storage imports media, so media importing storage would be a cycle.
	if strings.Contains(load, "internal/storage\"") || strings.Contains(image, "internal/storage\"") {
		t.Error("media must not import storage: storage imports media")
	}

	if !strings.HasPrefix(mediaImagePureGo(), "//go:build !vips") || !strings.HasPrefix(mediaImageVipsGo(), "//go:build vips") {
		t.Error("the chain's encoders must be split by the vips build tag")
	}
	if !strings.Contains(mediaImageVipsGo(), "var transformSlots") {
		t.Error("the vips build needs its own transformSlots: transform.go is left out of it")
	}
	for name, src := range map[string]string{
		"image.go": image, "image_load.go": load, "image_purego.go": mediaImagePureGo(),
		"image_vips.go": mediaImageVipsGo(), "image_test.go": mediaImageTestGo(),
		"image_storage_test.go": mediaImageStorageTestGo(),
	} {
		mustFormatGo(t, name, src)
		if strings.ContainsRune(src, rune(0x2014)) {
			t.Errorf("%s contains an em dash", name)
		}
	}
}

// A one-sided size keeps the aspect ratio in both backends. The zero used to
// reach imaging.Fit, which returns an empty image, and libvips' Thumbnail.
func TestOneSidedResizeKeepsAspectRatioInTemplates(t *testing.T) {
	pure := mediaTransformGo()
	if strings.Contains(pure, "imaging.Fit(src, s.Width, s.Height") {
		t.Error("transform.go still passes a possibly zero side to imaging.Fit")
	}
	if !strings.Contains(pure, "fitDimensions(b.Dx(), b.Dy(), s.Width, s.Height)") ||
		!strings.Contains(pure, "if s.Crop && s.Width > 0 && s.Height > 0 {") {
		t.Error("transform.go does not size a one-sided box by the aspect ratio")
	}
	vips := mediaTransformVipsGo()
	if !strings.Contains(vips, "fitDimensions(img.Width(), img.Height(), s.Width, s.Height)") ||
		strings.Contains(vips, "img.Thumbnail(s.Width, s.Height, vips.InterestingNone)") {
		t.Error("transform_vips.go does not size a one-sided box by the aspect ratio")
	}
	mustFormatGo(t, "transform_vips.go", vips)
}
