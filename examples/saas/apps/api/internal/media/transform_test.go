package media

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"math/rand"
	"strings"
	"testing"
)

// photo builds something with photographic entropy: smooth gradients plus
// per-pixel noise. A flat colour block would compress so well that it would
// hide exactly the regressions these tests exist to catch.
func photo(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	rng := rand.New(rand.NewSource(7))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			fx, fy := float64(x)/float64(w), float64(y)/float64(h)
			n := func(v float64) uint8 {
				return uint8(math.Max(0, math.Min(255, v+rng.Float64()*30-15)))
			}
			img.Set(x, y, color.RGBA{
				n(128 + 110*math.Sin(fx*5)*math.Cos(fy*3)),
				n(128 + 110*math.Sin(fx*2+fy*6)),
				n(128 + 110*math.Cos(fx*7-fy*3)), 255})
		}
	}
	return img
}

func transparent(w, h int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx, dy := float64(x-w/2), float64(y-h/2)
			if math.Sqrt(dx*dx+dy*dy) < float64(w)/3 {
				img.Set(x, y, color.RGBA{108, 92, 231, 255})
			}
		}
	}
	return img
}

func encodeJPEG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95}); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// A photograph must come out very much smaller. This is the whole point of the
// pipeline, so it is asserted as a ratio rather than a fixed size.
func TestTransformShrinksAPhotograph(t *testing.T) {
	src := encodeJPEG(t, photo(3000, 2000))
	res, err := Transform(bytes.NewReader(src), DefaultProfile())
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary.MIME != "image/jpeg" {
		t.Errorf("a photograph should encode as JPEG, got %s", res.Primary.MIME)
	}
	if ratio := len(src) / len(res.Primary.Bytes); ratio < 5 {
		t.Errorf("expected at least 5x smaller, got %dx (%d -> %d bytes)",
			ratio, len(src), len(res.Primary.Bytes))
	}
	if res.OriginalWidth != 3000 || res.OriginalHeight != 2000 {
		t.Errorf("source dimensions not reported: %dx%d", res.OriginalWidth, res.OriginalHeight)
	}
}

// The failure this prevents is a transparent logo silently gaining a black
// background because it was encoded as JPEG.
func TestTransparencySurvives(t *testing.T) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, transparent(600, 600)); err != nil {
		t.Fatal(err)
	}
	res, err := Transform(bytes.NewReader(buf.Bytes()), DefaultProfile())
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary.MIME != "image/webp" {
		t.Fatalf("an image with alpha must stay lossless, got %s", res.Primary.MIME)
	}
}

// Resampling leaves alpha values a hair under opaque. Treating those as
// transparency would send every resized photograph down the lossless path,
// where it is roughly twenty times larger.
func TestResizedPhotoIsNotTreatedAsTransparent(t *testing.T) {
	src := encodeJPEG(t, photo(2400, 1600))
	res, err := Transform(bytes.NewReader(src), DefaultProfile())
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary.MIME != "image/jpeg" {
		t.Errorf("resized photo took the lossless path: %s", res.Primary.MIME)
	}
}

// Fit must never scale up. Enlarging a small upload spends bytes on pixels the
// source never had.
func TestSmallImageIsNotEnlarged(t *testing.T) {
	src := encodeJPEG(t, photo(200, 150))
	res, err := Transform(bytes.NewReader(src), DefaultProfile())
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary.Width != 200 || res.Primary.Height != 150 {
		t.Errorf("expected 200x150 untouched, got %dx%d", res.Primary.Width, res.Primary.Height)
	}
}

func TestRenditionsAreProduced(t *testing.T) {
	src := encodeJPEG(t, photo(2000, 2000))
	res, err := Transform(bytes.NewReader(src), DefaultProfile())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Extra) != 1 || res.Extra[0].Name != "thumb" {
		t.Fatalf("expected one thumb rendition, got %+v", res.Extra)
	}
	// Fill crops to exactly the requested box.
	if res.Extra[0].Width != 400 || res.Extra[0].Height != 400 {
		t.Errorf("thumb should be exactly 400x400, got %dx%d",
			res.Extra[0].Width, res.Extra[0].Height)
	}
}

// A profile that sets one field keeps the defaults for the rest, or every
// profile in a project has to restate every number.
func TestPartialProfileKeepsDefaults(t *testing.T) {
	Define("tiny", Profile{Max: Fit(100, 100)})
	p := Get("tiny")
	if p.Max.Width != 100 {
		t.Errorf("declared field lost: %+v", p.Max)
	}
	// The reason DiscardOriginal is phrased as a negative. As KeepOriginal,
	// this assertion failed: withDefaults cannot tell a bool left unset from
	// one set to false, so a profile that never mentioned it threw the
	// original away while the docs promised it was kept.
	if p.DiscardOriginal {
		t.Error("a profile that does not mention the original must keep it")
	}
	if p.Quality != DefaultProfile().Quality {
		t.Errorf("quality should fall back to the default, got %v", p.Quality)
	}
	if p.Format != Auto {
		t.Errorf("format should fall back to Auto, got %v", p.Format)
	}
}

// An unknown name downgrades to the default rather than failing an upload.
func TestUnknownProfileFallsBack(t *testing.T) {
	if got := Get("no-such-profile"); got.Quality != DefaultProfile().Quality {
		t.Errorf("unknown profile should return the default, got %+v", got)
	}
	if Known("no-such-profile") {
		t.Error("Known must report an undefined profile as unknown")
	}
}

// Garbage in must be an error rather than a panic, because the upload handler
// turns that error into "store the original and mark it unoptimised".
func TestGarbageIsAnError(t *testing.T) {
	if _, err := Transform(bytes.NewReader([]byte("not an image")), DefaultProfile()); err == nil {
		t.Error("expected an error decoding garbage")
	}
}

// A decompression bomb is refused from its header, before pixels are
// allocated.
//
// The one test here about an attacker rather than a mistake. A solid-colour
// PNG compresses to almost nothing whatever its dimensions, so an upload that
// passes every file-size check on the way in can still be hundreds of
// megabytes once decoded.
func TestDecompressionBombIsRefused(t *testing.T) {
	const dim = 12000 // 144 megapixels, about 165 KB on the wire
	img := image.NewGray(image.Rect(0, 0, dim, dim))
	for i := range img.Pix {
		img.Pix[i] = 200
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	_, err := Transform(bytes.NewReader(buf.Bytes()), DefaultProfile())
	if err == nil {
		t.Fatal("a 144 megapixel image must be refused, not decoded")
	}
	if !strings.Contains(err.Error(), "megapixel") {
		t.Errorf("the error should name the limit, got: %v", err)
	}
}

// The limit must not refuse a real camera frame.
func TestLargeCameraFrameIsAccepted(t *testing.T) {
	src := encodeJPEG(t, photo(6000, 4000)) // 24 megapixels
	if _, err := Transform(bytes.NewReader(src), DefaultProfile()); err != nil {
		t.Fatalf("a 24 megapixel photo must be accepted: %v", err)
	}
}

// An animated GIF must not be silently reduced to its first frame.
func TestGIFIsNotOptimisable(t *testing.T) {
	if IsOptimisable("image/gif") {
		t.Error("GIF must be left alone: decoding one keeps only the first frame")
	}
}
