//go:build vips

package media

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"runtime"

	"github.com/davidbyttow/govips/v2/vips"
)

// transformSlots caps concurrent pixel work for the chain, one per CPU, as the
// pure-Go backend's transform.go does.
var transformSlots = make(chan struct{}, runtime.NumCPU())

// autoFormat is lossy WebP on this backend, as a profile's Auto is: it carries
// transparency and is smaller than JPEG at the same quality.
func autoFormat(image.Image) Format { return WebP }

// encodeImage hands the chain's pixels to libvips as an uncompressed PNG and
// exports them, with metadata stripped.
func encodeImage(img image.Image, f Format, quality int) ([]byte, error) {
	startup()

	var raw bytes.Buffer
	enc := png.Encoder{CompressionLevel: png.NoCompression}
	if err := enc.Encode(&raw, img); err != nil {
		return nil, fmt.Errorf("handing pixels to libvips: %w", err)
	}
	ref, err := vips.NewImageFromBuffer(raw.Bytes())
	if err != nil {
		return nil, fmt.Errorf("loading pixels into libvips: %w", err)
	}
	defer ref.Close()

	var out []byte
	switch f {
	case JPEG:
		ep := vips.NewJpegExportParams()
		ep.Quality = quality
		ep.StripMetadata = true
		out, _, err = ref.ExportJpeg(ep)
	case PNG:
		ep := vips.NewPngExportParams()
		ep.StripMetadata = true
		out, _, err = ref.ExportPng(ep)
	case AVIF:
		ep := vips.NewAvifExportParams()
		ep.Quality = quality
		ep.StripMetadata = true
		out, _, err = ref.ExportAvif(ep)
	default:
		ep := vips.NewWebpExportParams()
		ep.Quality = quality
		ep.StripMetadata = true
		out, _, err = ref.ExportWebp(ep)
	}
	if err != nil {
		return nil, fmt.Errorf("encoding %s: %w", f, err)
	}
	return out, nil
}
