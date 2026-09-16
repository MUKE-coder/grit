package scaffold

// The chainable image API (issue 002, workstream D): media.Open and friends.
//
// Written by writeMediaFiles next to transform.go, so grit upgrade delivers it
// with the rest of the framework-owned media package and no repair is needed.
//
//   image.go          the chain: Open, operations, inspection, encoders
//   image_load.go     FromDisk, FromURL and Store
//   image_purego.go   the pure-Go encoders the chain uses (default build)
//   image_vips.go     the libvips encoders the chain uses (-tags vips)
//   image_test.go     reference-image tests per operation, EXIF, limits
//   image_storage_test.go  Store on the real local disk, served back

func mediaImageGo() string {
	return `package media

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"time"

	// Registers the WebP decoder in both backends. transform.go imports it to
	// encode, but a -tags vips build leaves that file out.
	_ "github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
)

// Image is one image being worked on. Every operation returns a new Image, so
// a decoded original can be branched into several outputs:
//
//	img, err := media.Open(file)
//	if err != nil { ... }
//	thumb, err := img.Cover(150, 150).ToWebP().Encode()
//	large, err := img.Fit(1600, 0).Quality(80).ToJPEG().Encode()
//
// An operation that fails does not panic and does not return an error of its
// own: the error travels down the chain and comes out of Encode or Store, so
// a chain reads as one expression. Err reports it early when you want to.
type Image struct {
	img         image.Image
	err         error
	source      Format
	orientation int
	oriented    bool
	format      Format
	quality     int
	maxPixels   int
}

var (
	// ErrPixelLimit is returned, wrapped, for an image whose decoded size, or
	// the size an operation would produce, is over the pixel limit.
	ErrPixelLimit = errors.New("media: image is over the pixel limit")
	// ErrTooLarge is returned, wrapped, when the encoded input is over the
	// byte limit.
	ErrTooLarge = errors.New("media: image is over the byte limit")
	// ErrUnsupported is returned, wrapped, for input that is not JPEG, PNG or
	// WebP.
	ErrUnsupported = errors.New("media: unsupported image format")
	// ErrNeedsVips is returned, wrapped, for lossy WebP and AVIF on the pure-Go
	// backend, which has no encoder for either. Build with -tags vips.
	ErrNeedsVips = errors.New("media: this encoding needs the libvips backend, build with -tags vips")

	errNilImage = errors.New("media: no image, check the error Open returned")
)

const (
	// DefaultMaxBytes is the most encoded input Open, FromDisk and FromURL
	// read before refusing. A phone photo is 3 to 8 MB.
	DefaultMaxBytes int64 = 20 << 20
	// DefaultFetchTimeout bounds FromURL, connection and body together.
	DefaultFetchTimeout = 10 * time.Second
	// defaultQuality matches DefaultProfile's 0.82.
	defaultQuality = 82
)

// Option adjusts how an image is loaded.
type Option func(*options)

type options struct {
	maxPixels  int
	maxBytes   int64
	timeout    time.Duration
	autoOrient bool
}

func collectOptions(opts []Option) options {
	o := options{
		maxPixels:  DefaultProfile().MaxPixels,
		maxBytes:   DefaultMaxBytes,
		timeout:    DefaultFetchTimeout,
		autoOrient: true,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return o
}

// WithMaxPixels sets the pixel limit, 50 megapixels by default. It is checked
// against the header before decoding, and against the output of every
// operation that can make an image larger.
func WithMaxPixels(n int) Option {
	return func(o *options) {
		if n > 0 {
			o.maxPixels = n
		}
	}
}

// WithMaxBytes sets the most encoded input read, DefaultMaxBytes by default.
func WithMaxBytes(n int64) Option {
	return func(o *options) {
		if n > 0 {
			o.maxBytes = n
		}
	}
}

// WithTimeout bounds FromURL, DefaultFetchTimeout by default.
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		if d > 0 {
			o.timeout = d
		}
	}
}

// WithoutAutoOrient keeps the pixels as stored, ignoring the EXIF orientation
// until Orient is called. Without this option Open applies it, because a
// portrait phone photo is stored sideways with a flag saying so.
func WithoutAutoOrient() Option {
	return func(o *options) { o.autoOrient = false }
}

// Open reads and decodes an image: JPEG, PNG or WebP.
//
// The header is read first and anything over the pixel limit is refused
// before a pixel is allocated, which is what stops a decompression bomb: a
// 165 KB PNG can decode to 144 megapixels. GIF is refused rather than reduced
// to its first frame.
func Open(r io.Reader, opts ...Option) (*Image, error) {
	o := collectOptions(opts)
	data, err := readLimited(r, o.maxBytes)
	if err != nil {
		return nil, err
	}
	return decodeBytes(data, o)
}

func readLimited(r io.Reader, limit int64) ([]byte, error) {
	if r == nil {
		return nil, errors.New("media: nothing to read")
	}
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, fmt.Errorf("media: reading image: %w", err)
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("%w: more than %d bytes", ErrTooLarge, limit)
	}
	return data, nil
}

var sourceFormats = map[string]Format{"jpeg": JPEG, "png": PNG, "webp": WebP}

func decodeBytes(data []byte, o options) (*Image, error) {
	transformSlots <- struct{}{}
	defer func() { <-transformSlots }()

	cfg, name, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("media: reading image header: %w", err)
	}
	source, ok := sourceFormats[name]
	if !ok {
		if name == "gif" {
			return nil, fmt.Errorf("%w: gif, which is usually animated and would keep only its first frame", ErrUnsupported)
		}
		return nil, fmt.Errorf("%w: %s", ErrUnsupported, name)
	}
	if err := checkPixels(cfg.Width, cfg.Height, o.maxPixels); err != nil {
		return nil, err
	}
	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(false))
	if err != nil {
		return nil, fmt.Errorf("media: decoding image: %w", err)
	}
	out := &Image{
		img:         img,
		source:      source,
		orientation: exifOrientation(data),
		maxPixels:   o.maxPixels,
	}
	if o.autoOrient {
		out.img = orient(out.img, out.orientation)
		out.oriented = true
	}
	return out, nil
}

func checkPixels(width, height, limit int) error {
	if width <= 0 || height <= 0 {
		return fmt.Errorf("media: %dx%d is not an image size", width, height)
	}
	if px := int64(width) * int64(height); px > int64(limit) {
		return fmt.Errorf("%w: %dx%d is %d megapixels, over the %d megapixel limit",
			ErrPixelLimit, width, height, px/1000000, limit/1000000)
	}
	return nil
}

// apply runs one operation inside a transform slot and returns the next Image.
func (i *Image) apply(op string, fn func(src image.Image) (image.Image, error)) *Image {
	if i == nil {
		return &Image{err: errNilImage}
	}
	if i.err != nil {
		return i
	}
	next := *i
	transformSlots <- struct{}{}
	out, err := fn(i.img)
	<-transformSlots
	if err != nil {
		next.err = fmt.Errorf("media: %s: %w", op, err)
		return &next
	}
	next.img = out
	return &next
}

// Cover scales and crops to exactly width x height, keeping the centre. It is
// the avatar and card-image operation: nothing is letterboxed.
func (i *Image) Cover(width, height int) *Image {
	return i.apply("cover", func(src image.Image) (image.Image, error) {
		if width <= 0 || height <= 0 {
			return nil, fmt.Errorf("needs a width and a height above zero, got %dx%d (Fit and Resize take a 0 to keep the aspect ratio)", width, height)
		}
		if err := checkPixels(width, height, i.maxPixels); err != nil {
			return nil, err
		}
		return imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos), nil
	})
}

// Resize scales to width x height. A 0 for either side keeps the aspect
// ratio. Unlike Fit it will enlarge, up to the pixel limit.
func (i *Image) Resize(width, height int) *Image {
	return i.apply("resize", func(src image.Image) (image.Image, error) {
		if width < 0 || height < 0 || (width == 0 && height == 0) {
			return nil, fmt.Errorf("needs a width or a height above zero, got %dx%d", width, height)
		}
		b := src.Bounds()
		w, h := width, height
		if w == 0 {
			w = scaleSide(b.Dx(), height, b.Dy())
		}
		if h == 0 {
			h = scaleSide(b.Dy(), width, b.Dx())
		}
		if err := checkPixels(w, h, i.maxPixels); err != nil {
			return nil, err
		}
		return imaging.Resize(src, w, h, imaging.Lanczos), nil
	})
}

// Fit scales down to sit inside width x height, keeping the aspect ratio. A 0
// for either side leaves that side unbounded, so Fit(800, 0) means "at most
// 800 wide". It never enlarges: pixels the source never had only cost bytes.
func (i *Image) Fit(width, height int) *Image {
	return i.apply("fit", func(src image.Image) (image.Image, error) {
		if width < 0 || height < 0 || (width == 0 && height == 0) {
			return nil, fmt.Errorf("needs a width or a height above zero, got %dx%d", width, height)
		}
		b := src.Bounds()
		w, h := fitDimensions(b.Dx(), b.Dy(), width, height)
		if w == b.Dx() && h == b.Dy() {
			return src, nil
		}
		return imaging.Resize(src, w, h, imaging.Lanczos), nil
	})
}

// Crop keeps the width x height rectangle whose top-left corner is at x, y.
// A rectangle reaching outside the image is an error rather than being
// clipped, so a wrong offset is noticed.
func (i *Image) Crop(x, y, width, height int) *Image {
	return i.apply("crop", func(src image.Image) (image.Image, error) {
		b := src.Bounds()
		rect := image.Rect(x, y, x+width, y+height).Add(b.Min)
		if width <= 0 || height <= 0 || !rect.In(b) {
			return nil, fmt.Errorf("%dx%d at %d,%d is not inside the %dx%d image", width, height, x, y, b.Dx(), b.Dy())
		}
		return imaging.Crop(src, rect), nil
	})
}

// Rotate turns the image clockwise by degrees. A right angle is exact; any
// other angle grows the canvas to fit and fills the corners with background,
// transparent when it is nil.
func (i *Image) Rotate(degrees float64, background color.Color) *Image {
	return i.apply("rotate", func(src image.Image) (image.Image, error) {
		if background == nil {
			background = color.Transparent
		}
		b := src.Bounds()
		rad := degrees * math.Pi / 180
		sin, cos := math.Abs(math.Sin(rad)), math.Abs(math.Cos(rad))
		w := int(math.Round(float64(b.Dx())*cos + float64(b.Dy())*sin))
		h := int(math.Round(float64(b.Dx())*sin + float64(b.Dy())*cos))
		if err := checkPixels(w, h, i.maxPixels); err != nil {
			return nil, err
		}
		// imaging turns counter-clockwise.
		return imaging.Rotate(src, -degrees, background), nil
	})
}

// FlipH mirrors the image left to right.
func (i *Image) FlipH() *Image {
	return i.apply("flip", func(src image.Image) (image.Image, error) { return imaging.FlipH(src), nil })
}

// FlipV mirrors the image top to bottom.
func (i *Image) FlipV() *Image {
	return i.apply("flip", func(src image.Image) (image.Image, error) { return imaging.FlipV(src), nil })
}

// Grayscale removes the colour.
func (i *Image) Grayscale() *Image {
	return i.apply("grayscale", func(src image.Image) (image.Image, error) { return imaging.Grayscale(src), nil })
}

// Blur softens the image. amount runs from 0 (untouched) to 100 (a Gaussian
// sigma of 20 pixels, enough to hide a face in a 400 pixel thumbnail).
func (i *Image) Blur(amount float64) *Image {
	return i.apply("blur", func(src image.Image) (image.Image, error) {
		if amount < 0 || amount > 100 {
			return nil, fmt.Errorf("amount runs from 0 to 100, got %v", amount)
		}
		if amount == 0 {
			return src, nil
		}
		return imaging.Blur(src, amount*0.2), nil
	})
}

// Sharpen crisps edges, typically after a large downscale. amount runs from 0
// (untouched) to 100 (a sigma of 5 pixels); 10 to 30 suits a thumbnail.
func (i *Image) Sharpen(amount float64) *Image {
	return i.apply("sharpen", func(src image.Image) (image.Image, error) {
		if amount < 0 || amount > 100 {
			return nil, fmt.Errorf("amount runs from 0 to 100, got %v", amount)
		}
		if amount == 0 {
			return src, nil
		}
		return imaging.Sharpen(src, amount*0.05), nil
	})
}

// Orient applies the EXIF orientation. Open already does unless it was given
// WithoutAutoOrient, so calling it again changes nothing; it is there to make
// the step visible in a chain, and for images opened without it.
func (i *Image) Orient() *Image {
	if i == nil {
		return &Image{err: errNilImage}
	}
	if i.err != nil || i.oriented {
		return i
	}
	next := i.apply("orient", func(src image.Image) (image.Image, error) {
		return orient(src, i.orientation), nil
	})
	if next.err == nil {
		next.oriented = true
	}
	return next
}

// orient maps an EXIF orientation to the transform that makes it upright.
func orient(img image.Image, orientation int) image.Image {
	switch orientation {
	case 2:
		return imaging.FlipH(img)
	case 3:
		return imaging.Rotate180(img)
	case 4:
		return imaging.FlipV(img)
	case 5:
		return imaging.Transpose(img)
	case 6:
		return imaging.Rotate270(img)
	case 7:
		return imaging.Transverse(img)
	case 8:
		return imaging.Rotate90(img)
	}
	return img
}

// Err is the first error in the chain so far, or nil.
func (i *Image) Err() error {
	if i == nil {
		return errNilImage
	}
	return i.err
}

// Width is the current width in pixels, 0 once the chain has failed.
func (i *Image) Width() int {
	w, _ := i.Dimensions()
	return w
}

// Height is the current height in pixels, 0 once the chain has failed.
func (i *Image) Height() int {
	_, h := i.Dimensions()
	return h
}

// Dimensions is the current width and height, 0 and 0 once the chain has
// failed.
func (i *Image) Dimensions() (int, int) {
	if i.Err() != nil {
		return 0, 0
	}
	b := i.img.Bounds()
	return b.Dx(), b.Dy()
}

// Orientation is the EXIF orientation the source carried, 1 to 8, and 1 when
// it carried none.
func (i *Image) Orientation() int {
	if i == nil || i.orientation == 0 {
		return 1
	}
	return i.orientation
}

// SourceFormat is what the input was: JPEG, PNG or WebP.
func (i *Image) SourceFormat() Format {
	if i == nil {
		return ""
	}
	return i.source
}

// MIME is the type Encode will produce, for a Content-Type header.
func (i *Image) MIME() string {
	if i.Err() != nil {
		return ""
	}
	return i.outputFormat().mime()
}

// Extension is the file extension Encode will produce, with its dot: ".webp".
func (i *Image) Extension() string {
	if i.Err() != nil {
		return ""
	}
	return i.outputFormat().ext()
}

// outputFormat is the encoding chosen with a To method, or, when none was,
// the one Auto picks on this backend.
func (i *Image) outputFormat() Format {
	if i.format != "" && i.format != Auto {
		return i.format
	}
	return autoFormat(i.img)
}

// DominantColor is the most common colour, for a placeholder behind an image
// that is still loading.
//
// The image is sampled down to 64x64 and every opaque pixel is counted into
// one of 4096 buckets (4 bits a channel); the answer is the average of the
// pixels in the fullest bucket, so a solid colour comes back exactly rather
// than rounded to its bucket.
func (i *Image) DominantColor() (color.RGBA, error) {
	if err := i.Err(); err != nil {
		return color.RGBA{}, err
	}
	transformSlots <- struct{}{}
	defer func() { <-transformSlots }()

	small := imaging.Resize(i.img, 64, 64, imaging.NearestNeighbor)
	type bucket struct{ n, r, g, b int }
	var buckets [4096]bucket
	for y := 0; y < 64; y++ {
		for x := 0; x < 64; x++ {
			c := small.NRGBAAt(x, y)
			if c.A < 128 {
				continue
			}
			k := int(c.R>>4)<<8 | int(c.G>>4)<<4 | int(c.B>>4)
			buckets[k].n++
			buckets[k].r += int(c.R)
			buckets[k].g += int(c.G)
			buckets[k].b += int(c.B)
		}
	}
	var best bucket
	for _, bk := range &buckets {
		if bk.n > best.n {
			best = bk
		}
	}
	if best.n == 0 {
		return color.RGBA{}, errors.New("media: the image is fully transparent and has no dominant colour")
	}
	// #nosec G115 -- the mean of 8-bit samples is itself between 0 and 255.
	return color.RGBA{R: uint8(best.r / best.n), G: uint8(best.g / best.n), B: uint8(best.b / best.n), A: 255}, nil
}

// HexColor formats a colour as #rrggbb, for CSS.
func HexColor(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// ToJPEG encodes as JPEG at Quality.
func (i *Image) ToJPEG() *Image { return i.to(JPEG, "") }

// ToPNG encodes as PNG. Lossless, so Quality has no effect.
func (i *Image) ToPNG() *Image { return i.to(PNG, "") }

// ToWebP encodes as WebP. On the pure-Go backend this is lossless (VP8L),
// keeps transparency and ignores Quality, because there is no pure-Go lossy
// WebP encoder; built with -tags vips it is lossy at Quality. Use
// ToLossyWebP to insist on lossy output.
func (i *Image) ToWebP() *Image { return i.to(WebP, "") }

// ToLossyWebP encodes as lossy WebP at Quality. It needs -tags vips; the
// pure-Go backend returns ErrNeedsVips rather than quietly writing a lossless
// file twenty times the size.
func (i *Image) ToLossyWebP() *Image { return i.to(WebP, "lossy WebP") }

// ToAVIF encodes as AVIF at Quality. It needs -tags vips and a libvips built
// with libheif; the pure-Go backend returns ErrNeedsVips.
func (i *Image) ToAVIF() *Image { return i.to(AVIF, "AVIF") }

func (i *Image) to(f Format, needsVips string) *Image {
	return i.set(func(next *Image) error {
		if needsVips != "" && !SupportsLossyWebP() {
			return fmt.Errorf("%w: the pure-Go backend has no %s encoder", ErrNeedsVips, needsVips)
		}
		next.format = f
		return nil
	})
}

// Quality sets the lossy quality, 1 to 100, 82 by default. It has no effect
// on PNG, or on WebP from the pure-Go backend.
func (i *Image) Quality(q int) *Image {
	return i.set(func(next *Image) error {
		if q < 1 || q > 100 {
			return fmt.Errorf("media: quality runs from 1 to 100, got %d", q)
		}
		next.quality = q
		return nil
	})
}

// set changes a setting without touching pixels, so it needs no slot.
func (i *Image) set(fn func(next *Image) error) *Image {
	if i == nil {
		return &Image{err: errNilImage}
	}
	if i.err != nil {
		return i
	}
	next := *i
	if err := fn(&next); err != nil {
		next.err = err
	}
	return &next
}

// Encode returns the encoded bytes, or the first error in the chain. The
// output carries no EXIF: it is written from pixels, so a phone photo's GPS
// coordinates do not survive.
func (i *Image) Encode() ([]byte, error) {
	if err := i.Err(); err != nil {
		return nil, err
	}
	transformSlots <- struct{}{}
	defer func() { <-transformSlots }()

	q := i.quality
	if q == 0 {
		q = defaultQuality
	}
	return encodeImage(i.img, i.outputFormat(), q)
}

// fitDimensions is the size that fits srcW x srcH inside boxW x boxH without
// enlarging, keeping the aspect ratio. A box side of 0 is unbounded.
func fitDimensions(srcW, srcH, boxW, boxH int) (int, int) {
	w, h := srcW, srcH
	if w <= 0 || h <= 0 {
		return w, h
	}
	if boxW > 0 && w > boxW {
		h = scaleSide(h, boxW, w)
		w = boxW
	}
	if boxH > 0 && h > boxH {
		w = scaleSide(w, boxH, h)
		h = boxH
	}
	return w, h
}

// scaleSide is side * target / other, rounded, and never below one pixel.
func scaleSide(side, target, other int) int {
	v := int(math.Round(float64(side) * float64(target) / float64(other)))
	if v < 1 {
		return 1
	}
	return v
}

// exifOrientation reads the orientation tag from a JPEG's EXIF segment, and
// returns 1 for anything else. It walks every APP1 segment rather than only
// the first, because an XMP block can come before the EXIF one.
func exifOrientation(data []byte) int {
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		return 1
	}
	pos := 2
	for pos+4 <= len(data) {
		if data[pos] != 0xFF {
			return 1
		}
		marker := data[pos+1]
		if marker == 0xFF {
			pos++
			continue
		}
		if marker == 0xDA || marker == 0xD9 {
			return 1
		}
		size := int(data[pos+2])<<8 | int(data[pos+3])
		if size < 2 || pos+2+size > len(data) {
			return 1
		}
		segment := data[pos+4 : pos+2+size]
		if marker == 0xE1 && len(segment) >= 6 && string(segment[:6]) == "Exif\x00\x00" {
			if o := tiffOrientation(segment[6:]); o != 0 {
				return o
			}
		}
		pos += 2 + size
	}
	return 1
}

func tiffOrientation(tiff []byte) int {
	if len(tiff) < 8 {
		return 0
	}
	var order binary.ByteOrder
	switch string(tiff[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 0
	}
	ifd := int64(order.Uint32(tiff[4:8]))
	if ifd < 8 || ifd+2 > int64(len(tiff)) {
		return 0
	}
	count := int(order.Uint16(tiff[ifd:]))
	for k := 0; k < count; k++ {
		entry := int(ifd) + 2 + k*12
		if entry+12 > len(tiff) {
			return 0
		}
		if order.Uint16(tiff[entry:]) != 0x0112 {
			continue
		}
		if v := int(order.Uint16(tiff[entry+8:])); v >= 1 && v <= 8 {
			return v
		}
		return 0
	}
	return 0
}
`
}

func mediaImageLoadGo() string {
	return `package media

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"{{MODULE}}/internal/safefetch"
)

// Disk is where FromDisk reads and Store writes: pass a storage.Disk, from
// storage.Storage's Disk().
//
// It is not storage.Disk itself because the storage package imports media for
// its image helpers, and Go refuses an import cycle. So media names the
// methods it can name, which storage.Disk satisfies at compile time, and calls
// Put, whose options type lives in storage, through reflection in putFile.
type Disk interface {
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	URL(key string) string
}

// Visibility says who may read a stored image.
type Visibility string

const (
	// PublicFile is readable by anyone with the URL: an avatar, a product
	// photo.
	PublicFile Visibility = "public"
	// PrivateFile is read only through a temporary URL.
	PrivateFile Visibility = "private"
)

// Key prefixes for each visibility. The public one matches the first of
// storage.PublicPrefixes, which is what the local driver serves without a
// signature and what the bucket policy opens to anonymous reads; everything
// outside those prefixes is private on every driver.
const (
	publicKeyPrefix  = "uploads/"
	privateKeyPrefix = "private/"
)

// FromDisk opens the image stored at key, reading at most the byte limit.
func FromDisk(ctx context.Context, disk Disk, key string, opts ...Option) (*Image, error) {
	if disk == nil {
		return nil, errors.New("media: FromDisk needs a disk")
	}
	o := collectOptions(opts)
	body, err := disk.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("media: opening %q: %w", key, err)
	}
	data, err := readLimited(body, o.maxBytes)
	if closeErr := body.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf("media: closing %q: %w", key, closeErr)
	}
	if err != nil {
		return nil, err
	}
	return decodeBytes(data, o)
}

// FromURL downloads and opens an image from a URL somebody else chose.
//
// The request goes through safefetch, so a loopback, private or cloud metadata
// address is refused, including one a public hostname resolves to or
// redirects to. The body is read up to the byte limit and the whole fetch is
// bounded by the timeout (WithMaxBytes, WithTimeout).
func FromURL(ctx context.Context, rawURL string, opts ...Option) (*Image, error) {
	o := collectOptions(opts)
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	resp, err := safefetch.Get(ctx, rawURL)
	if err != nil {
		return nil, fmt.Errorf("media: fetching image: %w", err)
	}
	data, err := readResponse(resp, o.maxBytes)
	if closeErr := resp.Body.Close(); err == nil && closeErr != nil {
		err = fmt.Errorf("media: closing the response: %w", closeErr)
	}
	if err != nil {
		return nil, err
	}
	return decodeBytes(data, o)
}

func readResponse(resp *http.Response, limit int64) ([]byte, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("media: fetching image: the server answered %s", resp.Status)
	}
	if resp.ContentLength > limit {
		return nil, fmt.Errorf("%w: the server says %d bytes", ErrTooLarge, resp.ContentLength)
	}
	return readLimited(resp.Body, limit)
}

// Store encodes the image and writes it to disk under dir, returning the key.
//
//	key, err := img.Orient().Cover(150, 150).ToWebP().
//	    Store(ctx, svc.Storage.Disk(), "avatars", media.PublicFile)
//	url := svc.Storage.Disk().URL(key)
//
// A PublicFile lands under uploads/<dir>/, a PrivateFile under private/<dir>/,
// named with 24 random hex characters and the extension of what was encoded.
func (i *Image) Store(ctx context.Context, disk Disk, dir string, visibility Visibility) (string, error) {
	if err := i.Err(); err != nil {
		return "", err
	}
	if disk == nil {
		return "", errors.New("media: Store needs a disk")
	}
	var prefix string
	switch visibility {
	case PublicFile:
		prefix = publicKeyPrefix
	case PrivateFile:
		prefix = privateKeyPrefix
	default:
		return "", fmt.Errorf("media: visibility must be media.PublicFile or media.PrivateFile, got %q", visibility)
	}
	dir, err := cleanDir(dir)
	if err != nil {
		return "", err
	}
	data, err := i.Encode()
	if err != nil {
		return "", err
	}
	name := make([]byte, 12)
	if _, err := rand.Read(name); err != nil {
		return "", fmt.Errorf("media: naming the file: %w", err)
	}
	format := i.outputFormat()
	key := prefix + dir + hex.EncodeToString(name) + format.ext()
	if err := putFile(ctx, disk, key, data, format.mime(), visibility); err != nil {
		return "", err
	}
	return key, nil
}

// cleanDir returns dir with one trailing slash, or "" for none, refusing a
// directory that could climb out of its prefix.
func cleanDir(dir string) (string, error) {
	dir = strings.Trim(dir, "/")
	if dir == "" {
		return "", nil
	}
	for _, segment := range strings.Split(dir, "/") {
		if segment == "" || segment == "." || segment == ".." || strings.ContainsAny(segment, "\\\x00") {
			return "", fmt.Errorf("media: %q is not a directory Store can write to", dir)
		}
	}
	return dir + "/", nil
}

var (
	contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	readerType  = reflect.TypeOf((*io.Reader)(nil)).Elem()
	errorType   = reflect.TypeOf((*error)(nil)).Elem()
)

// putFile calls disk.Put(ctx, key, body, storage.PutOptions{ContentType: ...}).
//
// Reflection, because media cannot import storage (see Disk). The shape is
// checked before the call, so a disk without that method is a clear error
// rather than a panic. If the options gain a Visibility string field, it is
// set too.
func putFile(ctx context.Context, disk Disk, key string, data []byte, contentType string, visibility Visibility) error {
	method := reflect.ValueOf(disk).MethodByName("Put")
	if !method.IsValid() {
		return fmt.Errorf("media: %T has no Put method, pass a storage.Disk", disk)
	}
	t := method.Type()
	if t.NumIn() != 4 || t.NumOut() != 1 || t.In(0) != contextType || t.In(1).Kind() != reflect.String ||
		t.In(2) != readerType || t.In(3).Kind() != reflect.Struct || t.Out(0) != errorType {
		return fmt.Errorf("media: %T.Put is not Put(ctx, key, reader, options) error, pass a storage.Disk", disk)
	}
	opts := reflect.New(t.In(3)).Elem()
	if field := opts.FieldByName("ContentType"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(contentType)
	}
	if field := opts.FieldByName("Visibility"); field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
		field.SetString(string(visibility))
	}
	if ctx == nil {
		ctx = context.Background()
	}
	out := method.Call([]reflect.Value{
		reflect.ValueOf(&ctx).Elem(),
		reflect.ValueOf(key).Convert(t.In(1)),
		reflect.ValueOf(io.Reader(bytes.NewReader(data))),
		opts,
	})
	if err, ok := out[0].Interface().(error); ok && err != nil {
		return fmt.Errorf("media: storing %q: %w", key, err)
	}
	return nil
}
`
}

func mediaImagePureGo() string {
	return `//go:build !vips

package media

import "image"

// encodeImage writes a chain's pixels with the encoders transform.go uses.
// WebP here is lossless: quality has no effect on it.
func encodeImage(img image.Image, f Format, quality int) ([]byte, error) {
	return encode(img, f, float64(quality)/100)
}

// autoFormat is what an Image with no To call encodes as: lossless WebP when
// it has transparency, JPEG otherwise, the same choice a profile's Auto makes.
func autoFormat(img image.Image) Format {
	return resolveFormat(Auto, img)
}
`
}

func mediaImageVipsGo() string {
	return `//go:build vips

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
`
}
