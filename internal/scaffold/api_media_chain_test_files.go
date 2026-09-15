package scaffold

// Tests generated into a project's internal/media for the chainable API.

func mediaImageTestGo() string {
	return `package media

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"{{MODULE}}/internal/safefetch"
)

var (
	red    = color.NRGBA{220, 30, 30, 255}
	green  = color.NRGBA{30, 200, 60, 255}
	blue   = color.NRGBA{40, 60, 220, 255}
	yellow = color.NRGBA{240, 220, 40, 255}
)

// quadrants is the reference image: four solid quarters, so every flip and
// rotation moves a colour somewhere checkable.
func quadrants(w, h int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := red
			switch {
			case x >= w/2 && y < h/2:
				c = green
			case x < w/2 && y >= h/2:
				c = blue
			case x >= w/2 && y >= h/2:
				c = yellow
			}
			img.SetNRGBA(x, y, c)
		}
	}
	return img
}

// mapped builds a w x h image whose pixel (x, y) is src at at(x, y): the
// reference for a flip, rotation or crop, written from the definition rather
// than by calling the code under test.
func mapped(src image.Image, w, h int, at func(x, y int) (int, int)) *image.NRGBA {
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx, sy := at(x, y)
			out.Set(x, y, src.At(sx, sy))
		}
	}
	return out
}

func pngBytes(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func jpegBytes(t *testing.T, img image.Image, quality int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality}); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func mustOpen(t *testing.T, data []byte, opts ...Option) *Image {
	t.Helper()
	img, err := Open(bytes.NewReader(data), opts...)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

func decoded(t *testing.T, data []byte) image.Image {
	t.Helper()
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	return img
}

// meanDiff is the mean absolute difference per RGB channel, 0 to 255.
func meanDiff(t *testing.T, got, want image.Image) float64 {
	t.Helper()
	gb, wb := got.Bounds(), want.Bounds()
	if gb.Dx() != wb.Dx() || gb.Dy() != wb.Dy() {
		t.Fatalf("got %dx%d, want %dx%d", gb.Dx(), gb.Dy(), wb.Dx(), wb.Dy())
	}
	var sum float64
	for y := 0; y < wb.Dy(); y++ {
		for x := 0; x < wb.Dx(); x++ {
			r1, g1, b1, _ := got.At(gb.Min.X+x, gb.Min.Y+y).RGBA()
			r2, g2, b2, _ := want.At(wb.Min.X+x, wb.Min.Y+y).RGBA()
			sum += math.Abs(float64(r1>>8)-float64(r2>>8)) +
				math.Abs(float64(g1>>8)-float64(g2>>8)) +
				math.Abs(float64(b1>>8)-float64(b2>>8))
		}
	}
	return sum / float64(3*wb.Dx()*wb.Dy())
}

func meanColour(img image.Image) [3]float64 {
	b := img.Bounds()
	var out [3]float64
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			out[0] += float64(r >> 8)
			out[1] += float64(g >> 8)
			out[2] += float64(bl >> 8)
		}
	}
	n := float64(b.Dx() * b.Dy())
	return [3]float64{out[0] / n, out[1] / n, out[2] / n}
}

// Each operation against a reference image built independently of it, with a
// tolerance per operation: 0 where the result is exact, a few levels where a
// resampling filter rings at a hard edge.
func TestOperationsMatchReferenceImages(t *testing.T) {
	const w, h = 256, 128
	src := quadrants(w, h)
	base := mustOpen(t, pngBytes(t, src))

	gray := mapped(src, w, h, func(x, y int) (int, int) { return x, y })
	for i := 0; i < len(gray.Pix); i += 4 {
		l := uint8(math.Round(0.299*float64(gray.Pix[i]) + 0.587*float64(gray.Pix[i+1]) + 0.114*float64(gray.Pix[i+2])))
		gray.Pix[i], gray.Pix[i+1], gray.Pix[i+2] = l, l, l
	}

	cases := []struct {
		name      string
		chain     *Image
		want      image.Image
		tolerance float64
	}{
		{"FlipH", base.FlipH(), mapped(src, w, h, func(x, y int) (int, int) { return w - 1 - x, y }), 0},
		{"FlipV", base.FlipV(), mapped(src, w, h, func(x, y int) (int, int) { return x, h - 1 - y }), 0},
		{"Rotate(90) is clockwise", base.Rotate(90, nil), mapped(src, h, w, func(x, y int) (int, int) { return y, h - 1 - x }), 0},
		{"Rotate(180)", base.Rotate(180, nil), mapped(src, w, h, func(x, y int) (int, int) { return w - 1 - x, h - 1 - y }), 0},
		{"Rotate(-90)", base.Rotate(-90, nil), mapped(src, h, w, func(x, y int) (int, int) { return w - 1 - y, x }), 0},
		{"Crop", base.Crop(100, 20, 60, 50), mapped(src, 60, 50, func(x, y int) (int, int) { return 100 + x, 20 + y }), 0},
		{"Grayscale", base.Grayscale(), gray, 1},
		// Measured 0.2 to 1.0 where the Lanczos filter rings at the colour
		// edges. A one pixel offset in the framing measures about 4.7, so
		// these catch it.
		{"Resize(128, 0)", base.Resize(128, 0), quadrants(128, 64), 1},
		{"Resize(0, 32)", base.Resize(0, 32), quadrants(64, 32), 2.5},
		{"Resize(300, 100) stretches", base.Resize(300, 100), quadrants(300, 100), 1},
		{"Fit(100, 100)", base.Fit(100, 100), quadrants(100, 50), 1.5},
		{"Fit(0, 64)", base.Fit(0, 64), quadrants(128, 64), 1},
		{"Fit(1000, 0) does not enlarge", base.Fit(1000, 0), src, 0},
		// The centre 128x128 of the source holds a quarter of each colour.
		{"Cover(64, 64)", base.Cover(64, 64), quadrants(64, 64), 1.5},
		{"Cover(128, 128)", base.Cover(128, 128), quadrants(128, 128), 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.chain.Err(); err != nil {
				t.Fatal(err)
			}
			out, err := tc.chain.ToPNG().Encode()
			if err != nil {
				t.Fatal(err)
			}
			d := meanDiff(t, decoded(t, out), tc.want)
			t.Logf("mean difference %.2f, tolerance %.2f", d, tc.tolerance)
			if d > tc.tolerance {
				t.Errorf("mean difference %.2f, tolerance %.2f", d, tc.tolerance)
			}
		})
	}
}

// Blur and sharpen have no closed-form reference, so each is held to what it
// must do: change the edges within a bound and leave the average colour alone.
func TestBlurAndSharpen(t *testing.T) {
	src := quadrants(256, 128)
	base := mustOpen(t, pngBytes(t, src))
	want := meanColour(src)
	for _, tc := range []struct {
		name     string
		chain    *Image
		min, max float64
	}{
		{"Blur(10)", base.Blur(10), 1, 12},
		{"Blur(50) softens more", base.Blur(50), 5, 30},
		{"Sharpen(60)", base.Sharpen(60), 0.3, 12},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.chain.ToPNG().Encode()
			if err != nil {
				t.Fatal(err)
			}
			img := decoded(t, out)
			if d := meanDiff(t, img, src); d < tc.min || d > tc.max {
				t.Errorf("mean difference from the source %.2f, want %.1f to %.1f", d, tc.min, tc.max)
			}
			got := meanColour(img)
			for c := range got {
				if math.Abs(got[c]-want[c]) > 3 {
					t.Errorf("average colour moved: %v, want %v", got, want)
					break
				}
			}
		})
	}
	if base.Blur(0).Err() != nil || base.Blur(101).Err() == nil || base.Sharpen(-1).Err() == nil {
		t.Error("amounts run from 0 to 100")
	}
}

func TestEncoders(t *testing.T) {
	src := quadrants(64, 32)
	base := mustOpen(t, pngBytes(t, src))

	out, err := base.ToPNG().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if d := meanDiff(t, decoded(t, out), src); d != 0 {
		t.Errorf("PNG must be exact, mean difference %.2f", d)
	}

	jpg, err := base.Quality(90).ToJPEG().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if d := meanDiff(t, decoded(t, jpg), src); d > 6 {
		t.Errorf("JPEG q90 mean difference %.2f, tolerance 6", d)
	}
	if m, e := base.ToJPEG().MIME(), base.ToJPEG().Extension(); m != "image/jpeg" || e != ".jpg" {
		t.Errorf("JPEG reports %s %s", m, e)
	}

	webp, err := base.ToWebP().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if base.ToWebP().MIME() != "image/webp" {
		t.Errorf("WebP reports %s", base.ToWebP().MIME())
	}
	if !SupportsLossyWebP() {
		// Lossless on this backend: exact, whatever the quality says.
		if d := meanDiff(t, decoded(t, webp), src); d != 0 {
			t.Errorf("pure-Go WebP is lossless, mean difference %.2f", d)
		}
		low, err := base.Quality(10).ToWebP().Encode()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(low, webp) {
			t.Error("pure-Go WebP is documented to ignore quality")
		}
	}

	// Quality works on a photograph, where there is something to lose.
	noisy := image.NewNRGBA(image.Rect(0, 0, 400, 300))
	rng := rand.New(rand.NewSource(3))
	rng.Read(noisy.Pix)
	photo := mustOpen(t, pngBytes(t, noisy))
	small, err := photo.Quality(30).ToJPEG().Encode()
	if err != nil {
		t.Fatal(err)
	}
	large, err := photo.Quality(95).ToJPEG().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if len(small) >= len(large) {
		t.Errorf("quality 30 gave %d bytes, quality 95 gave %d", len(small), len(large))
	}
	if base.Quality(0).Err() == nil || base.Quality(101).Err() == nil {
		t.Error("quality runs from 1 to 100")
	}

	// Auto: an opaque image is a JPEG, a transparent one keeps its alpha.
	if got := base.MIME(); got != "image/jpeg" && !SupportsLossyWebP() {
		t.Errorf("an opaque image with no To call should be JPEG, got %s", got)
	}
	if got := base.Rotate(30, nil).MIME(); got != "image/webp" {
		t.Errorf("transparent corners should keep alpha, got %s", got)
	}
}

// Asking the pure-Go backend for lossy WebP or AVIF is an error naming the
// build tag, not a lossless file twenty times the size.
func TestLossyWebPAndAVIFNeedVips(t *testing.T) {
	if SupportsLossyWebP() {
		t.Skip("built with -tags vips")
	}
	base := mustOpen(t, pngBytes(t, quadrants(64, 32)))
	for name, chain := range map[string]*Image{"AVIF": base.ToAVIF(), "lossy WebP": base.ToLossyWebP()} {
		_, err := chain.Quality(80).Encode()
		if !errors.Is(err, ErrNeedsVips) {
			t.Errorf("%s: expected ErrNeedsVips, got %v", name, err)
			continue
		}
		if !strings.Contains(err.Error(), "-tags vips") || !strings.Contains(err.Error(), name) {
			t.Errorf("%s: the error should say what and how, got: %v", name, err)
		}
	}
}

// exifSegment is an APP1 block holding one orientation tag.
func exifSegment(orientation uint16, littleEndian bool) []byte {
	seg := []byte{0xFF, 0xE1, 0x00, 0x22, 'E', 'x', 'i', 'f', 0, 0}
	if littleEndian {
		return append(seg,
			'I', 'I', 0x2A, 0x00, 0x08, 0x00, 0x00, 0x00,
			0x01, 0x00,
			0x12, 0x01, 0x03, 0x00, 0x01, 0x00, 0x00, 0x00, byte(orientation), byte(orientation >> 8), 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00)
	}
	return append(seg,
		'M', 'M', 0x00, 0x2A, 0x00, 0x00, 0x00, 0x08,
		0x00, 0x01,
		0x01, 0x12, 0x00, 0x03, 0x00, 0x00, 0x00, 0x01, byte(orientation >> 8), byte(orientation), 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00)
}

// storedFor lays out an upright image the way a camera stores it under each
// EXIF orientation, from the definitions in the EXIF specification.
func storedFor(upright image.Image, orientation int) *image.NRGBA {
	b := upright.Bounds()
	W, H := b.Dx(), b.Dy()
	switch orientation {
	case 2: // mirrored horizontally
		return mapped(upright, W, H, func(x, y int) (int, int) { return W - 1 - x, y })
	case 3: // rotated 180
		return mapped(upright, W, H, func(x, y int) (int, int) { return W - 1 - x, H - 1 - y })
	case 4: // mirrored vertically
		return mapped(upright, W, H, func(x, y int) (int, int) { return x, H - 1 - y })
	case 5: // mirrored along the top-left diagonal
		return mapped(upright, H, W, func(x, y int) (int, int) { return y, x })
	case 6: // shown after turning 90 clockwise
		return mapped(upright, H, W, func(x, y int) (int, int) { return W - 1 - y, x })
	case 7: // mirrored along the top-right diagonal
		return mapped(upright, H, W, func(x, y int) (int, int) { return W - 1 - y, H - 1 - x })
	case 8: // shown after turning 90 counter-clockwise
		return mapped(upright, H, W, func(x, y int) (int, int) { return y, H - 1 - x })
	}
	return mapped(upright, W, H, func(x, y int) (int, int) { return x, y })
}

func withEXIF(jpg, segment []byte) []byte {
	out := append([]byte{}, jpg[:2]...)
	out = append(out, segment...)
	return append(out, jpg[2:]...)
}

// All eight orientations come out upright, and the EXIF does not come out at
// all.
func TestEXIFOrientationsComeOutUpright(t *testing.T) {
	upright := quadrants(96, 48)
	for orientation := 1; orientation <= 8; orientation++ {
		for _, le := range []bool{false, true} {
			data := withEXIF(jpegBytes(t, storedFor(upright, orientation), 95), exifSegment(uint16(orientation), le))
			img := mustOpen(t, data)
			if img.Orientation() != orientation {
				t.Errorf("orientation %d (little endian %v): read as %d", orientation, le, img.Orientation())
			}
			out, err := img.ToPNG().Encode()
			if err != nil {
				t.Fatal(err)
			}
			if d := meanDiff(t, decoded(t, out), upright); d > 4 {
				t.Errorf("orientation %d (little endian %v): mean difference from upright %.2f, tolerance 4", orientation, le, d)
			}
			jpg, err := img.ToJPEG().Encode()
			if err != nil {
				t.Fatal(err)
			}
			if bytes.Contains(jpg, []byte("Exif")) {
				t.Errorf("orientation %d: EXIF survived the encode", orientation)
			}
		}
	}
}

// Opened without auto-orientation, the stored pixels come back as stored, and
// Orient is what turns them upright, once.
func TestOrientIsExplicitWhenAsked(t *testing.T) {
	upright := quadrants(96, 48)
	data := withEXIF(jpegBytes(t, storedFor(upright, 6), 95), exifSegment(6, false))
	raw := mustOpen(t, data, WithoutAutoOrient())
	if w, h := raw.Dimensions(); w != 48 || h != 96 {
		t.Fatalf("stored pixels should be 48x96, got %dx%d", w, h)
	}
	fixed := raw.Orient().Orient()
	if w, h := fixed.Dimensions(); w != 96 || h != 48 {
		t.Fatalf("Orient should give 96x48, got %dx%d", w, h)
	}
	out, err := fixed.ToPNG().Encode()
	if err != nil {
		t.Fatal(err)
	}
	if d := meanDiff(t, decoded(t, out), upright); d > 4 {
		t.Errorf("mean difference from upright %.2f", d)
	}
}

func TestDominantColor(t *testing.T) {
	solid := image.NewNRGBA(image.Rect(0, 0, 300, 200))
	for i := 0; i < len(solid.Pix); i += 4 {
		copy(solid.Pix[i:], []uint8{108, 92, 231, 255})
	}
	got, err := mustOpen(t, pngBytes(t, solid)).DominantColor()
	if err != nil {
		t.Fatal(err)
	}
	if got != (color.RGBA{108, 92, 231, 255}) || HexColor(got) != "#6c5ce7" {
		t.Errorf("solid: got %v (%s)", got, HexColor(got))
	}

	// Seventy per cent one colour, thirty the other, through a lossy JPEG.
	two := image.NewNRGBA(image.Rect(0, 0, 500, 300))
	for y := 0; y < 300; y++ {
		for x := 0; x < 500; x++ {
			if x < 350 {
				two.SetNRGBA(x, y, color.NRGBA{0, 184, 148, 255})
			} else {
				two.SetNRGBA(x, y, color.NRGBA{255, 107, 107, 255})
			}
		}
	}
	got, err = mustOpen(t, jpegBytes(t, two, 90)).DominantColor()
	if err != nil {
		t.Fatal(err)
	}
	if absDiff(got.R, 0) > 4 || absDiff(got.G, 184) > 4 || absDiff(got.B, 148) > 4 {
		t.Errorf("two-tone: got %v, want about {0 184 148}", got)
	}

	if _, err := mustOpen(t, pngBytes(t, image.NewNRGBA(image.Rect(0, 0, 10, 10)))).DominantColor(); err == nil {
		t.Error("a fully transparent image has no dominant colour")
	}
}

func absDiff(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

// The pixel limit holds on every entry point: the header on the way in, and
// the output of any operation that can grow an image.
func TestPixelLimitHoldsThroughTheChain(t *testing.T) {
	const dim = 12000 // 144 megapixels, about 165 KB on the wire
	bomb := image.NewGray(image.Rect(0, 0, dim, dim))
	for i := range bomb.Pix {
		bomb.Pix[i] = 200
	}
	data := pngBytes(t, bomb)
	_, err := Open(bytes.NewReader(data))
	if !errors.Is(err, ErrPixelLimit) || !strings.Contains(err.Error(), "megapixel") {
		t.Fatalf("a %d KB, 144 megapixel PNG must be refused from its header, got %v", len(data)/1024, err)
	}
	if _, err := FromDisk(context.Background(), newMemDisk(map[string][]byte{"bomb.png": data}), "bomb.png"); !errors.Is(err, ErrPixelLimit) {
		t.Errorf("FromDisk must refuse it too, got %v", err)
	}

	small := mustOpen(t, pngBytes(t, quadrants(64, 32)))
	for name, chain := range map[string]*Image{
		"Resize": small.Resize(20000, 0),
		"Cover":  small.Cover(9000, 9000),
		"Rotate": mustOpen(t, pngBytes(t, quadrants(64, 32)), WithMaxPixels(2500)).Rotate(45, nil),
	} {
		if _, err := chain.Encode(); !errors.Is(err, ErrPixelLimit) {
			t.Errorf("%s past the limit: expected ErrPixelLimit, got %v", name, err)
		}
	}
	if _, err := Open(bytes.NewReader(pngBytes(t, quadrants(64, 32))), WithMaxPixels(1000)); !errors.Is(err, ErrPixelLimit) {
		t.Errorf("WithMaxPixels must lower the limit, got %v", err)
	}
}

func TestByteLimit(t *testing.T) {
	data := pngBytes(t, quadrants(64, 32))
	if _, err := Open(bytes.NewReader(data), WithMaxBytes(int64(len(data)-1))); !errors.Is(err, ErrTooLarge) {
		t.Errorf("expected ErrTooLarge, got %v", err)
	}
	if _, err := Open(bytes.NewReader(data), WithMaxBytes(int64(len(data)))); err != nil {
		t.Errorf("exactly the limit must pass: %v", err)
	}
	disk := newMemDisk(map[string][]byte{"a.png": data})
	if _, err := FromDisk(context.Background(), disk, "a.png", WithMaxBytes(10)); !errors.Is(err, ErrTooLarge) {
		t.Errorf("FromDisk: expected ErrTooLarge, got %v", err)
	}
}

// GIF is refused rather than reduced to its first frame.
func TestGIFIsRefused(t *testing.T) {
	var buf bytes.Buffer
	if err := gif.Encode(&buf, quadrants(16, 16), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(&buf); !errors.Is(err, ErrUnsupported) {
		t.Errorf("expected ErrUnsupported, got %v", err)
	}
}

// FromURL goes through safefetch: a loopback server is never contacted.
func TestFromURLRefusesLoopback(t *testing.T) {
	var hits atomic.Int32
	body := pngBytes(t, quadrants(64, 32))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		if _, err := w.Write(body); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()

	for _, url := range []string{server.URL + "/a.png", "http://localhost/a.png", "http://169.254.169.254/latest/meta-data/", "file:///etc/passwd"} {
		if _, err := FromURL(context.Background(), url); !errors.Is(err, safefetch.ErrBlocked) {
			t.Errorf("%s: expected safefetch.ErrBlocked, got %v", url, err)
		}
	}
	if n := hits.Load(); n != 0 {
		t.Errorf("the loopback server was contacted %d times", n)
	}
}

// The response checks FromURL runs once safefetch has allowed a target.
func TestFromURLResponseLimits(t *testing.T) {
	ok := &http.Response{StatusCode: http.StatusOK, ContentLength: -1, Body: io.NopCloser(bytes.NewReader(make([]byte, 100)))}
	if _, err := readResponse(ok, 50); !errors.Is(err, ErrTooLarge) {
		t.Errorf("a body over the limit with no Content-Length: expected ErrTooLarge, got %v", err)
	}
	declared := &http.Response{StatusCode: http.StatusOK, ContentLength: 1 << 30, Body: io.NopCloser(strings.NewReader(""))}
	if _, err := readResponse(declared, 50); !errors.Is(err, ErrTooLarge) {
		t.Errorf("a declared 1 GB body: expected ErrTooLarge before reading, got %v", err)
	}
	missing := &http.Response{StatusCode: http.StatusNotFound, Status: "404 Not Found", Body: io.NopCloser(strings.NewReader(""))}
	if _, err := readResponse(missing, 50); err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("a 404 should be an error naming it, got %v", err)
	}
}

// memDisk has storage.Disk's Put shape with its own options type, which is
// all putFile relies on.
type memDisk struct {
	files map[string][]byte
	types map[string]string
	vis   map[string]string
}

type memPutOptions struct {
	ContentType string
	Visibility  string
}

func newMemDisk(files map[string][]byte) *memDisk {
	return &memDisk{files: files, types: map[string]string{}, vis: map[string]string{}}
}

func (d *memDisk) Put(ctx context.Context, key string, r io.Reader, opts memPutOptions) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	d.files[key], d.types[key], d.vis[key] = data, opts.ContentType, opts.Visibility
	return nil
}

func (d *memDisk) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	data, ok := d.files[key]
	if !ok {
		return nil, errors.New("not found")
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (d *memDisk) URL(key string) string { return "https://files.example.com/" + key }

type noPutDisk struct{}

func (noPutDisk) Get(ctx context.Context, key string) (io.ReadCloser, error) { return nil, errors.New("no") }
func (noPutDisk) URL(key string) string                                       { return "" }

func TestStoreAndFromDisk(t *testing.T) {
	ctx := context.Background()
	disk := newMemDisk(map[string][]byte{})
	img := mustOpen(t, pngBytes(t, quadrants(640, 320)))

	key, err := img.Orient().Cover(150, 150).ToPNG().Store(ctx, disk, "/avatars/", PublicFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(key, "uploads/avatars/") || !strings.HasSuffix(key, ".png") || len(key) != len("uploads/avatars/")+24+4 {
		t.Errorf("unexpected key %q", key)
	}
	if disk.types[key] != "image/png" || disk.vis[key] != "public" {
		t.Errorf("stored with content type %q and visibility %q", disk.types[key], disk.vis[key])
	}
	back, err := FromDisk(ctx, disk, key)
	if err != nil {
		t.Fatal(err)
	}
	if w, h := back.Dimensions(); w != 150 || h != 150 {
		t.Errorf("stored image is %dx%d, want 150x150", w, h)
	}

	private, err := img.ToJPEG().Store(ctx, disk, "originals/2026", PrivateFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(private, "private/originals/2026/") || !strings.HasSuffix(private, ".jpg") {
		t.Errorf("unexpected private key %q", private)
	}
	again, err := img.ToJPEG().Store(ctx, disk, "originals/2026", PrivateFile)
	if err != nil || again == private {
		t.Errorf("two stores must get two keys: %q %q %v", private, again, err)
	}

	for _, dir := range []string{"../etc", "a/../../b", "a//b"} {
		if _, err := img.Store(ctx, disk, dir, PublicFile); err == nil {
			t.Errorf("directory %q must be refused", dir)
		}
	}
	if _, err := img.Store(ctx, disk, "a", Visibility("world")); err == nil {
		t.Error("an unknown visibility must be refused")
	}
	if _, err := img.Store(ctx, noPutDisk{}, "a", PublicFile); err == nil || !strings.Contains(err.Error(), "Put") {
		t.Errorf("a disk without Put should say so, got %v", err)
	}
	// A failed chain stores nothing.
	before := len(disk.files)
	if _, err := img.Crop(600, 0, 100, 100).Store(ctx, disk, "a", PublicFile); err == nil || len(disk.files) != before {
		t.Errorf("a failed chain must not store: %v", err)
	}
}

// Operations return new images, so one decode can feed several outputs.
func TestChainBranchesFromOneDecode(t *testing.T) {
	base := mustOpen(t, pngBytes(t, quadrants(256, 128)))
	thumb := base.Cover(32, 32)
	wide := base.Fit(200, 0).ToJPEG()
	if w, h := base.Dimensions(); w != 256 || h != 128 {
		t.Errorf("the base changed to %dx%d", w, h)
	}
	if thumb.Width() != 32 || wide.Width() != 200 || wide.Height() != 100 {
		t.Errorf("branches: thumb %dx%d, wide %dx%d", thumb.Width(), thumb.Height(), wide.Width(), wide.Height())
	}
	if base.format != "" || wide.MIME() != "image/jpeg" {
		t.Errorf("a To call on a branch must not change the base: base %q, branch %s", base.format, wide.MIME())
	}
	var nilImage *Image
	if _, err := nilImage.Cover(10, 10).Encode(); err == nil {
		t.Error("a nil image must be an error, not a panic")
	}
}

// Every entry point waits for a transform slot, so a burst of chains cannot
// decode more images at once than there are CPUs.
func TestTheChainWaitsForATransformSlot(t *testing.T) {
	data := pngBytes(t, quadrants(64, 32))
	for i := 0; i < cap(transformSlots); i++ {
		transformSlots <- struct{}{}
	}
	done := make(chan error, 1)
	go func() {
		img, err := Open(bytes.NewReader(data))
		if err == nil {
			_, err = img.Resize(32, 0).Encode()
		}
		done <- err
	}()
	ranEarly := false
	select {
	case <-done:
		ranEarly = true
	case <-time.After(100 * time.Millisecond):
	}
	// Released before any failure is reported, or every later test would
	// block on a slot.
	for i := 0; i < cap(transformSlots); i++ {
		<-transformSlots
	}
	if ranEarly {
		t.Fatal("the chain ran with every slot taken")
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

// Fit and Fill with one side 0 keep the aspect ratio. Before, the zero reached
// imaging.Fit, which returned an empty 0x0 image that was encoded and stored
// without an error.
func TestOneSidedSizesKeepTheAspectRatio(t *testing.T) {
	src := jpegBytes(t, quadrants(2000, 1000), 90)
	res, err := Transform(bytes.NewReader(src), Profile{
		Max: Fit(800, 0),
		Renditions: map[string]Size{
			"tall": Fit(0, 250),
			"fill": Fill(400, 0),
			"big":  Fit(5000, 0),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary.Width != 800 || res.Primary.Height != 400 {
		t.Errorf("Fit(800, 0) gave %dx%d, want 800x400", res.Primary.Width, res.Primary.Height)
	}
	want := map[string][2]int{"tall": {500, 250}, "fill": {400, 200}, "big": {2000, 1000}}
	for _, r := range res.Extra {
		if w := want[r.Name]; r.Width != w[0] || r.Height != w[1] {
			t.Errorf("%s gave %dx%d, want %dx%d", r.Name, r.Width, r.Height, w[0], w[1])
		}
	}
	if len(res.Extra) != len(want) {
		t.Errorf("expected %d renditions, got %d", len(want), len(res.Extra))
	}
}
`
}

func mediaImageStorageTestGo() string {
	return `package media_test

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

	"{{MODULE}}/internal/media"
	"{{MODULE}}/internal/storage"
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
`
}
