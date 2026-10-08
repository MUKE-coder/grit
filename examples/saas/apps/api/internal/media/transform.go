//go:build !vips

// Pure-Go image pipeline: the default.
//
// Encodes JPEG and lossless WebP with no system libraries, so the API stays a
// single static binary and "grit deploy" can keep cross-compiling to Linux
// from any machine. Build with -tags vips to swap in libvips, which is several
// times faster and can write lossy WebP and AVIF, at the cost of that.

package media

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"runtime"

	"github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
)

// transformSlots lets one image decode run per CPU. A decoded 50 megapixel image
// is about 200 MB, and uploads arriving together each decoded at once.
var transformSlots = make(chan struct{}, runtime.NumCPU())

// Transform decodes, orients, resizes and re-encodes an image according to a
// profile.
//
// Two things happen here that are not obvious from the name.
//
// EXIF orientation is applied on decode. Without it, a photo taken in portrait
// on a phone arrives with landscape pixel data plus a rotation flag, and every
// resize produces a sideways image.
//
// EXIF is then stripped, for free: the output is encoded from decoded pixels,
// so no source metadata survives. That matters more than it sounds, because a
// phone photo carries GPS coordinates, and a shop publishing product photos
// would otherwise publish the seller's home address with them.
func Transform(r io.ReadSeeker, p Profile) (Result, error) {
	transformSlots <- struct{}{}
	defer func() { <-transformSlots }()

	p = withDefaults(p)

	// Read the dimensions from the header first. DecodeConfig parses only the
	// header, so this costs nothing and it is the only chance to refuse a
	// decompression bomb: once Decode runs, the memory is already committed.
	cfg, _, err := image.DecodeConfig(r)
	if err != nil {
		return Result{}, fmt.Errorf("reading image header: %w", err)
	}
	if px := int64(cfg.Width) * int64(cfg.Height); px > int64(p.MaxPixels) {
		return Result{}, fmt.Errorf(
			"image is %dx%d (%d megapixels), over the %d megapixel limit",
			cfg.Width, cfg.Height, px/1000000, p.MaxPixels/1000000)
	}
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return Result{}, fmt.Errorf("rewinding after the header: %w", err)
	}

	src, err := imaging.Decode(r, imaging.AutoOrientation(true))
	if err != nil {
		return Result{}, fmt.Errorf("decoding image: %w", err)
	}

	b := src.Bounds()
	out := Result{
		Optimised:      true,
		OriginalWidth:  b.Dx(),
		OriginalHeight: b.Dy(),
	}

	primary := resize(src, p.Max)
	format := resolveFormat(p.Format, primary)
	enc, err := encode(primary, format, p.Quality)
	if err != nil {
		return Result{}, err
	}
	pb := primary.Bounds()
	out.Primary = Rendition{
		Name: "primary", Bytes: enc, Width: pb.Dx(), Height: pb.Dy(),
		MIME: format.mime(), Ext: format.ext(),
	}

	for name, size := range p.Renditions {
		// Derived from the source rather than from the primary, so a crop is
		// taken at full detail instead of from an already-downscaled copy.
		img := resize(src, size)
		// The same format decision, taken per rendition: a crop can lose the
		// transparent margin that made the primary a WebP.
		rf := resolveFormat(p.Format, img)
		rb, err := encode(img, rf, p.Quality)
		if err != nil {
			return Result{}, err
		}
		ib := img.Bounds()
		out.Extra = append(out.Extra, Rendition{
			Name: name, Bytes: rb, Width: ib.Dx(), Height: ib.Dy(),
			MIME: rf.mime(), Ext: rf.ext(),
		})
	}

	return out, nil
}

func resize(src image.Image, s Size) image.Image {
	if s.Width <= 0 && s.Height <= 0 {
		return src
	}
	if s.Crop && s.Width > 0 && s.Height > 0 {
		return imaging.Fill(src, s.Width, s.Height, imaging.Center, imaging.Lanczos)
	}
	// Fit never scales up: enlarging a small upload to the profile's box wastes
	// bytes on pixels the source never had. A side of 0 is unbounded, so
	// Fit(800, 0) keeps the aspect ratio; passed to imaging.Fit it produced an
	// empty 0x0 image. A Fill with one side 0 has no box to crop to and scales
	// the same way.
	b := src.Bounds()
	w, h := fitDimensions(b.Dx(), b.Dy(), s.Width, s.Height)
	if w == b.Dx() && h == b.Dy() {
		return src
	}
	return imaging.Resize(src, w, h, imaging.Lanczos)
}

// resolveFormat turns Auto into a concrete encoding.
//
// A profile asking for AVIF gets JPEG here, because this backend has no AVIF
// encoder. It is downgraded rather than refused so that one binary built
// without -tags vips still serves a project whose profiles assume it, and the
// ref records what was actually produced so nothing downstream has to guess.
func resolveFormat(f Format, img image.Image) Format {
	if f == AVIF {
		if hasAlpha(img) {
			return WebP
		}
		return JPEG
	}
	if f != Auto {
		return f
	}
	if hasAlpha(img) {
		return WebP
	}
	return JPEG
}

// hasAlpha reports whether any pixel is meaningfully transparent.
//
// The threshold is not 0xffff because alpha survives resampling as values a
// hair under opaque, and treating those as transparency would push every
// resized photograph down the lossless path and make it twenty times larger.
func hasAlpha(img image.Image) bool {
	if op, ok := img.(interface{ Opaque() bool }); ok && op.Opaque() {
		return false
	}
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if _, _, _, a := img.At(x, y).RGBA(); a < 0xff00 {
				return true
			}
		}
	}
	return false
}

func encode(img image.Image, f Format, quality float64) ([]byte, error) {
	var buf bytes.Buffer
	switch f {
	case PNG:
		e := png.Encoder{CompressionLevel: png.BestCompression}
		if err := e.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("encoding png: %w", err)
		}
	case WebP:
		if err := nativewebp.Encode(&buf, img, &nativewebp.Options{}); err != nil {
			return nil, fmt.Errorf("encoding webp: %w", err)
		}
	default:
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality(quality)}); err != nil {
			return nil, fmt.Errorf("encoding jpeg: %w", err)
		}
	}
	return buf.Bytes(), nil
}

// Backend names the compiled-in image backend, for logs and grit doctor.
func Backend() string { return "pure-go" }

// SupportsLossyWebP reports whether this backend can write lossy WebP or AVIF.
// The pure-Go encoder cannot: it is lossless (VP8L) only.
func SupportsLossyWebP() bool { return false }
