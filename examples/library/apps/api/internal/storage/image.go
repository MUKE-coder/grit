package storage

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"library/apps/api/internal/media"
)

// These helpers are the path an image takes when the upload pipeline did not
// handle it: a presigned upload, which goes straight to storage, or a type the
// pipeline skipped. They run through media.Transform, the pipeline itself, so
// EXIF orientation is applied and metadata stripped the same way on both paths.
// Before, they decoded without orientation, and a portrait phone photo got a
// sideways thumbnail.

// MaxImageWidth is the maximum width for processed images.
const MaxImageWidth = 1920

// ThumbnailSize is the edge of the square thumbnail GenerateThumbnail makes:
// the "thumb" rendition of media's default profile, so a thumbnail is the same
// size whichever path made it.
const ThumbnailSize = 400

// MaxImageBytes is the most the helpers read of an image. Uploads are capped at
// 50 MB, and the one byte over tells a larger file from one exactly at the cap.
const MaxImageBytes = 50<<20 + 1

var (
	// ErrImageTooLarge is returned for an image refused before decoding: over
	// MaxImageBytes, or with a header claiming more pixels than the media
	// profile allows. A few kilobytes of PNG can claim 30000x30000, which
	// decodes to 3.6 GB and takes the process with it.
	ErrImageTooLarge = errors.New("image is too large to decode")
	// ErrUnreadableImage is returned for data that is not an image this
	// package can decode. Like ErrImageTooLarge it fails the same way every
	// time, so a job should not retry it.
	ErrUnreadableImage = errors.New("image cannot be decoded")
)

// readImage reads an image and checks its header, refusing it when the header
// claims more pixels than media.Get("").MaxPixels. Decoding commits the memory
// for every pixel before it reads them, so the header is the only place a
// decompression bomb can be stopped, and it is stopped here with an error a job
// knows not to retry.
func readImage(reader io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(reader, MaxImageBytes))
	if err != nil {
		return nil, fmt.Errorf("reading image: %w", err)
	}
	if len(data) >= MaxImageBytes {
		return nil, fmt.Errorf("image is over %d MB: %w", (MaxImageBytes-1)>>20, ErrImageTooLarge)
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("reading image header: %v: %w", err, ErrUnreadableImage)
	}
	limit := int64(media.Get("").MaxPixels)
	if px := int64(cfg.Width) * int64(cfg.Height); limit > 0 && px > limit {
		return nil, fmt.Errorf("image is %dx%d, over the %d megapixel limit: %w",
			cfg.Width, cfg.Height, limit/1000000, ErrImageTooLarge)
	}
	return data, nil
}

// transformImage runs one image through media.Transform at size, keeping the
// format it came in: PNG stays PNG, anything else becomes JPEG.
func transformImage(reader io.Reader, mimeType string, size media.Size) ([]byte, error) {
	data, err := readImage(reader)
	if err != nil {
		return nil, err
	}
	profile := media.Get("")
	profile.Max = size
	profile.Quality = 0.85
	profile.Format = media.JPEG
	if strings.EqualFold(mimeType, "image/png") {
		profile.Format = media.PNG
	}
	// No extra renditions: the caller wants this one image.
	profile.Renditions = map[string]media.Size{}
	result, err := media.Transform(bytes.NewReader(data), profile)
	if err != nil {
		return nil, fmt.Errorf("processing image: %v: %w", err, ErrUnreadableImage)
	}
	return result.Primary.Bytes, nil
}

// ProcessImage resizes an image wider than MaxImageWidth, keeping its aspect
// ratio and applying its EXIF orientation, and returns the encoded bytes.
func ProcessImage(reader io.Reader, mimeType string) ([]byte, error) {
	// The bound is the width: the height limit is far past any image under the
	// pixel limit that is also this wide.
	return transformImage(reader, mimeType, media.Fit(MaxImageWidth, MaxImageWidth*64))
}

// GenerateThumbnail makes a ThumbnailSize square thumbnail, cropped from the
// centre of the image the right way up.
func GenerateThumbnail(reader io.Reader, mimeType string) ([]byte, error) {
	return transformImage(reader, mimeType, media.Fill(ThumbnailSize, ThumbnailSize))
}

// IsImageMimeType returns true if the MIME type is a supported image format.
func IsImageMimeType(mimeType string) bool {
	switch strings.ToLower(mimeType) {
	case "image/jpeg", "image/png", "image/gif":
		return true
	}
	return false
}
