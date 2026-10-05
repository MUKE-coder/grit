//go:build vips

// libvips image pipeline: opt-in, via -tags vips.
//
// What it buys, over the pure-Go default:
//
//   - Several times faster, and it streams, so a large image is not fully
//     decoded into the heap before it can be resized.
//   - Lossy WebP, which is typically 25-35%% smaller than JPEG at the same
//     visual quality, and unlike the pure-Go encoder it keeps transparency
//     without falling back to lossless.
//   - AVIF, smaller again, when the linked libvips was built with libheif.
//
// What it costs:
//
//   - cgo. The API is no longer a single static binary, so "grit deploy",
//     which cross-compiles to linux/amd64 with CGO_ENABLED=0, cannot build it
//     from a machine that is not the target. Build in Docker instead.
//   - libvips 8.10+ must be present at build time and on the host at runtime.
//
// Requires: libvips-dev (Debian/Ubuntu) or vips (Homebrew), then
//   go get github.com/davidbyttow/govips/v2/vips
//
// The same Profile drives both backends. What changes is the encoding Auto
// resolves to, which is why FileRef records the format actually produced
// rather than leaving a client to infer it.

package media

import (
	"fmt"
	"io"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"
)

var startOnce sync.Once

func startup() {
	startOnce.Do(func() {
		// Quiet: libvips logs at info level are noise in an API server.
		vips.LoggingSettings(nil, vips.LogLevelError)
		vips.Startup(nil)
	})
}

// Backend names the compiled-in image backend, for logs and grit doctor.
func Backend() string { return "libvips " + vips.Version }

// SupportsLossyWebP reports whether this backend can write lossy WebP or AVIF.
func SupportsLossyWebP() bool { return true }

// Transform decodes, orients, resizes and re-encodes an image, entirely inside
// libvips.
//
// The work is done with Thumbnail rather than a decode-then-resize pair,
// because libvips can use a format's own reduced-resolution decoding to avoid
// ever materialising the full frame. That, and not the encoder, is where the
// memory difference against the pure-Go path comes from.
func Transform(r io.ReadSeeker, p Profile) (Result, error) {
	startup()
	p = withDefaults(p)

	buf, err := io.ReadAll(r)
	if err != nil {
		return Result{}, fmt.Errorf("reading image: %w", err)
	}

	// Dimensions first, so a decompression bomb is refused before any pixels
	// are committed. libvips is lazy: loading a buffer parses the header and
	// reading Width/Height does not decode the image.
	probe, err := vips.NewImageFromBuffer(buf)
	if err != nil {
		return Result{}, fmt.Errorf("reading image header: %w", err)
	}
	srcW, srcH := probe.Width(), probe.Height()
	probe.Close()
	if px := int64(srcW) * int64(srcH); px > int64(p.MaxPixels) {
		return Result{}, fmt.Errorf(
			"image is %dx%d (%d megapixels), over the %d megapixel limit",
			srcW, srcH, px/1000000, p.MaxPixels/1000000)
	}

	out := Result{
		Optimised:      true,
		OriginalWidth:  srcW,
		OriginalHeight: srcH,
	}

	primary, err := renderOne(buf, p.Max, p)
	if err != nil {
		return Result{}, err
	}
	primary.Name = "primary"
	out.Primary = primary

	for name, size := range p.Renditions {
		// From the source bytes rather than from the primary, so a crop is
		// taken at full detail instead of from an already-downscaled copy.
		rend, err := renderOne(buf, size, p)
		if err != nil {
			return Result{}, err
		}
		rend.Name = name
		out.Extra = append(out.Extra, rend)
	}

	return out, nil
}

func renderOne(buf []byte, s Size, p Profile) (Rendition, error) {
	img, err := vips.NewImageFromBuffer(buf)
	if err != nil {
		return Rendition{}, fmt.Errorf("loading image: %w", err)
	}
	defer img.Close()

	// EXIF orientation applied here. Metadata is stripped at export, so a
	// phone photo neither arrives sideways nor leaves carrying GPS.
	if err := img.AutoRotate(); err != nil {
		return Rendition{}, fmt.Errorf("orienting image: %w", err)
	}

	if s.Crop && s.Width > 0 && s.Height > 0 {
		// Centre, matching what the pure-Go backend does, so the two produce
		// the same framing.
		if err := img.Thumbnail(s.Width, s.Height, vips.InterestingCentre); err != nil {
			return Rendition{}, fmt.Errorf("resizing image: %w", err)
		}
	} else if w, h := fitDimensions(img.Width(), img.Height(), s.Width, s.Height); w != img.Width() || h != img.Height() {
		// Fit, and never upscale: enlarging a small upload spends bytes on
		// pixels the source never had. A side of 0 is unbounded, so Fit(800, 0)
		// keeps the aspect ratio instead of asking libvips for a zero height.
		if err := img.Thumbnail(w, h, vips.InterestingNone); err != nil {
			return Rendition{}, fmt.Errorf("resizing image: %w", err)
		}
	}

	format := p.Format
	if format == Auto {
		// Lossy WebP for everything, which this backend can do and the pure-Go
		// one cannot. It carries alpha, so a transparent logo needs no special
		// case, and it is smaller than JPEG at the same quality.
		format = WebP
	}

	quality := jpegQuality(p.Quality)
	var encoded []byte
	switch format {
	case JPEG:
		ep := vips.NewJpegExportParams()
		ep.Quality = quality
		ep.StripMetadata = true
		encoded, _, err = img.ExportJpeg(ep)
	case PNG:
		ep := vips.NewPngExportParams()
		ep.StripMetadata = true
		encoded, _, err = img.ExportPng(ep)
	case AVIF:
		ep := vips.NewAvifExportParams()
		ep.Quality = quality
		ep.StripMetadata = true
		encoded, _, err = img.ExportAvif(ep)
	default:
		ep := vips.NewWebpExportParams()
		ep.Quality = quality
		ep.StripMetadata = true
		encoded, _, err = img.ExportWebp(ep)
	}
	if err != nil {
		return Rendition{}, fmt.Errorf("encoding %s: %w", format, err)
	}

	return Rendition{
		Bytes: encoded, Width: img.Width(), Height: img.Height(),
		MIME: format.mime(), Ext: format.ext(),
	}, nil
}
