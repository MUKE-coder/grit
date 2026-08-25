package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// writeMediaFiles writes internal/media: the image optimisation pipeline.
func writeMediaFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "media", "profile.go"):        mediaProfileGo(),
		filepath.Join(apiRoot, "internal", "media", "transform.go"):      mediaTransformGo(),
		filepath.Join(apiRoot, "internal", "media", "transform_vips.go"): mediaTransformVipsGo(),
		filepath.Join(apiRoot, "internal", "media", "transform_test.go"): mediaTransformTestGo(module),
	}

	// Yours to edit, so it is written once and never overwritten.
	once := map[string]string{
		filepath.Join(apiRoot, "internal", "media", "profiles.go"): mediaProfilesGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	for path, content := range once {
		if _, err := os.Stat(path); err == nil {
			continue
		}
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func mediaProfileGo() string {
	return `package media

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Format is what a transformed image is encoded as.
type Format string

const (
	// Auto picks per image, and is the default because the right answer is
	// decidable from the pixels rather than something a developer should have
	// to know:
	//
	//   has meaningful alpha -> WebP lossless, which keeps transparency and
	//                           beats PNG on both size and encode time
	//   otherwise            -> JPEG, which is roughly 20x smaller than any
	//                           lossless encoding of the same photograph
	//
	// The failure this prevents is a transparent logo encoded as JPEG, which
	// silently gains a black box. Under Auto that is unrepresentable.
	Auto Format = "auto"
	JPEG Format = "jpeg"
	PNG  Format = "png"
	// WebP is lossless (VP8L) on the pure-Go backend, because there is no
	// pure-Go lossy WebP encoder. Built with -tags vips it is lossy, carries
	// alpha, and is smaller than JPEG at the same quality.
	WebP Format = "webp"
	// AVIF requires -tags vips and a libvips built with libheif. The pure-Go
	// backend cannot produce it and falls back to its own Auto choice, which
	// is why FileRef records the format actually produced.
	AVIF Format = "avif"
)

// Size is a target box. Fit scales down to sit inside it and keeps the aspect
// ratio; Fill crops to exactly these dimensions.
type Size struct {
	Width  int
	Height int
	Crop   bool
}

// Fit returns a Size that scales down inside the box, preserving aspect ratio.
func Fit(w, h int) Size { return Size{Width: w, Height: h} }

// Fill returns a Size cropped to exactly these dimensions, centred.
func Fill(w, h int) Size { return Size{Width: w, Height: h, Crop: true} }

// FailurePolicy decides what happens when an image cannot be transformed.
type FailurePolicy string

const (
	// StoreOriginal keeps the upload and records optimised:false on the ref.
	// The default, because losing somebody's file because an encoder choked is
	// worse than storing a large one, and "grit media reprocess" can retry it.
	StoreOriginal FailurePolicy = "store_original"
	// Reject refuses the upload outright.
	Reject FailurePolicy = "reject"
)

// Profile is the optimisation configuration for one file field.
type Profile struct {
	// Max is the bounding box for the primary rendition.
	Max Size
	// Quality is the lossy quality, 0 to 1. Ignored for lossless formats.
	Quality float64
	// Format is the encoding target. Auto is almost always right.
	Format Format
	// Renditions are extra sizes derived alongside the primary one.
	Renditions map[string]Size
	// DiscardOriginal throws the untouched upload away instead of keeping it
	// under a private prefix.
	//
	// Inverted, so that the zero value is the recommended behaviour. As
	// KeepOriginal it was a promise the code could not keep: withDefaults
	// restores any field left at its zero value, and a bool has no way to say
	// "unset", so every profile that did not mention it silently discarded the
	// original while the documentation said the opposite.
	DiscardOriginal bool
	// OnError decides what a failed transform does.
	OnError FailurePolicy
	// MaxPixels refuses an image whose decoded size would exceed this, read
	// from the header before any pixels are allocated.
	//
	// The attack this closes is a decompression bomb: a solid-colour PNG
	// compresses to almost nothing whatever its dimensions, so 165 KB on the
	// wire can be 144 megapixels decoded. Measured, that one upload allocated
	// 224 MB, and ten concurrent would have been 2.2 GB. The file size limit
	// upstream cannot see it, because on disk it is small.
	//
	// 50 megapixels passes a 48 MP professional camera frame and refuses the
	// bomb.
	MaxPixels int
}

// DefaultProfile applies to every image field that does not name one.
//
// Every number here is a default somebody would otherwise have to research:
//
//   1600  covers a 2x retina display at a typical content width. Storing more
//         is paying to keep pixels no browser will draw.
//   0.82  the point at which JPEG is visually indistinguishable from source.
//         Below about 0.75 artifacts appear on gradients and skin.
//   400   the thumbnail an admin table row or a card grid needs at 2x.
//
// Measured on a 5.3 MB phone photo: 147 KB out, including the thumbnail.
func DefaultProfile() Profile {
	return Profile{
		Max:          Fit(1600, 1600),
		Quality:      0.82,
		Format:       Auto,
		Renditions: map[string]Size{"thumb": Fill(400, 400)},
		OnError:    StoreOriginal,
		MaxPixels:  50_000_000,
	}
}

var (
	mu       sync.RWMutex
	profiles = map[string]Profile{}
)

// Define registers a named profile. Call it from an init function or from
// main before the server starts serving.
//
//	media.Define("product-image", media.Profile{
//	    Max:     media.Fit(800, 800),
//	    Quality: 0.8,
//	    Format:  media.Auto,
//	})
//
// Fields left zero fall back to the corresponding default, so a profile that
// only wants a different size says only that.
func Define(name string, p Profile) {
	mu.Lock()
	defer mu.Unlock()
	profiles[name] = withDefaults(p)
}

// Get returns a profile by name, falling back to the default when the name is
// empty or unknown.
//
// Unknown rather than an error on purpose: a stale profile name in a client
// request should downgrade to sensible behaviour, not fail an upload. Use
// Known to validate a name when you do want to reject one.
func Get(name string) Profile {
	if name == "" {
		return DefaultProfile()
	}
	mu.RLock()
	defer mu.RUnlock()
	if p, ok := profiles[strings.ToLower(name)]; ok {
		return p
	}
	return DefaultProfile()
}

// Known reports whether a profile name has been defined.
func Known(name string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := profiles[strings.ToLower(name)]
	return ok
}

// Names lists the registered profiles, sorted. Used by the admin to offer them
// and by "grit doctor" to report what a project defines.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(profiles))
	for name := range profiles {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

func withDefaults(p Profile) Profile {
	d := DefaultProfile()
	if p.Max.Width == 0 && p.Max.Height == 0 {
		p.Max = d.Max
	}
	if p.Quality <= 0 || p.Quality > 1 {
		p.Quality = d.Quality
	}
	if p.Format == "" {
		p.Format = d.Format
	}
	if p.Renditions == nil {
		p.Renditions = d.Renditions
	}
	if p.OnError == "" {
		p.OnError = d.OnError
	}
	if p.MaxPixels <= 0 {
		p.MaxPixels = d.MaxPixels
	}
	return p
}

// jpegQuality converts the 0-1 scale to the 1-100 the encoder wants.
func jpegQuality(q float64) int {
	v := int(q * 100)
	if v < 1 {
		v = 1
	}
	if v > 100 {
		v = 100
	}
	return v
}

func (f Format) mime() string {
	switch f {
	case PNG:
		return "image/png"
	case WebP:
		return "image/webp"
	case AVIF:
		return "image/avif"
	default:
		return "image/jpeg"
	}
}

func (f Format) ext() string {
	switch f {
	case PNG:
		return ".png"
	case WebP:
		return ".webp"
	case AVIF:
		return ".avif"
	default:
		return ".jpg"
	}
}

// Public is the wire form of a profile, for clients that do the optimising.
//
// Uploads go browser to storage directly, so the resizing happens on the
// client. Publishing the profile means the numbers live in one place: a client
// that hardcodes 1600 and a server default of 1600 are the same number until
// somebody changes one of them.
type Public struct {
	Name       string                   ` + "`" + `json:"name"` + "`" + `
	MaxWidth   int                      ` + "`" + `json:"max_width"` + "`" + `
	MaxHeight  int                      ` + "`" + `json:"max_height"` + "`" + `
	Crop       bool                     ` + "`" + `json:"crop"` + "`" + `
	Quality    float64                  ` + "`" + `json:"quality"` + "`" + `
	Format     string                   ` + "`" + `json:"format"` + "`" + `
	MaxPixels  int                      ` + "`" + `json:"max_pixels"` + "`" + `
	Renditions map[string]PublicSize    ` + "`" + `json:"renditions,omitempty"` + "`" + `
}

// PublicSize is one rendition target on the wire.
type PublicSize struct {
	Width  int  ` + "`" + `json:"width"` + "`" + `
	Height int  ` + "`" + `json:"height"` + "`" + `
	Crop   bool ` + "`" + `json:"crop"` + "`" + `
}

// ToPublic converts a profile to its wire form.
func ToPublic(name string, p Profile) Public {
	p = withDefaults(p)
	out := Public{
		Name:      name,
		MaxWidth:  p.Max.Width,
		MaxHeight: p.Max.Height,
		Crop:      p.Max.Crop,
		Quality:   p.Quality,
		Format:    string(p.Format),
		MaxPixels: p.MaxPixels,
	}
	if len(p.Renditions) > 0 {
		out.Renditions = map[string]PublicSize{}
		for k, v := range p.Renditions {
			out.Renditions[k] = PublicSize{Width: v.Width, Height: v.Height, Crop: v.Crop}
		}
	}
	return out
}

// AllPublic returns every registered profile plus the default, for
// GET /api/v1/media/profiles.
func AllPublic() []Public {
	out := []Public{ToPublic("default", DefaultProfile())}
	for _, name := range Names() {
		out = append(out, ToPublic(name, Get(name)))
	}
	return out
}

// Rendition is one encoded output.
type Rendition struct {
	Name   string
	Bytes  []byte
	Width  int
	Height int
	MIME   string
	Ext    string
}

// Result is everything a transform produced.
type Result struct {
	// Primary is the rendition the field points at.
	Primary Rendition
	// Extra holds the profile's named renditions, e.g. "thumb".
	Extra []Rendition
	// Optimised is false when the pipeline fell back to storing the original,
	// which is recorded on the ref so the admin can show what needs attention
	// and reprocessing can find it later.
	Optimised bool
	// OriginalWidth and OriginalHeight are the source dimensions, before any
	// resizing, so "6 MB 4000x3000 -> 147 KB 1600x1200" can be reported.
	OriginalWidth  int
	OriginalHeight int
}

// IsOptimisable reports whether the pipeline can handle this MIME type.
//
// GIF is excluded deliberately: it is usually animated, and decoding one
// yields the first frame only, so "optimising" it would silently throw the
// animation away.
func IsOptimisable(mime string) bool {
	switch mime {
	case "image/jpeg", "image/png", "image/webp":
		return true
	}
	return false
}

func (p Profile) String() string {
	return fmt.Sprintf("%dx%d q%.2f %s", p.Max.Width, p.Max.Height, p.Quality, p.Format)
}
`
}

func mediaTransformGo() string {
	return `//go:build !vips

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

	"github.com/HugoSmits86/nativewebp"
	"github.com/disintegration/imaging"
)

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
	p = withDefaults(p)

	// Read the dimensions from the header first. DecodeConfig parses only the
	// header, so this costs nothing and it is the only chance to refuse a
	// decompression bomb: once Decode runs, the memory is already committed.
	cfg, _, err := image.DecodeConfig(r)
	if err != nil {
		return Result{}, fmt.Errorf("reading image header: %w", err)
	}
	if px := cfg.Width * cfg.Height; px > p.MaxPixels {
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
	if s.Crop {
		return imaging.Fill(src, s.Width, s.Height, imaging.Center, imaging.Lanczos)
	}
	// Fit never scales up: enlarging a small upload to the profile's box wastes
	// bytes on pixels the source never had.
	b := src.Bounds()
	if b.Dx() <= s.Width && b.Dy() <= s.Height {
		return src
	}
	return imaging.Fit(src, s.Width, s.Height, imaging.Lanczos)
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
`
}

func mediaTransformTestGo(module string) string {
	return `package media

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
`
}

func mediaProfilesGo() string {
	return `package media

// Image optimisation profiles for this project.
//
// This file is yours. It is written once and never regenerated, so anything
// you add here survives grit generate and grit upgrade.
//
// You do not have to define anything. Every image field without a profile gets
// DefaultProfile(): fit inside 1600x1600, quality 0.82, format chosen per
// image, a 400x400 thumbnail, the original kept privately for reprocessing,
// and EXIF stripped. On a 6 MB phone photo that is about 150 KB out.
//
// Define a profile when a particular field wants something different, then
// name it from the upload: POST /api/v1/uploads?profile=product-image
//
// Fields you leave zero keep the default, so a profile that only wants a
// different size says only that.

func init() {
	// A product photograph on a storefront. Smaller than the default, because
	// a catalogue page shows a dozen of them at once.
	Define("product-image", Profile{
		Max:     Fit(1000, 1000),
		Quality: 0.8,
		Renditions: map[string]Size{
			"thumb": Fill(300, 300),
			"card":  Fit(600, 600),
		},
	})

	// An avatar is always square and always small, so it crops rather than
	// fits: a portrait photo shown in a round frame should not letterbox.
	Define("avatar", Profile{
		Max:        Fill(400, 400),
		Quality:    0.85,
		Renditions: map[string]Size{"thumb": Fill(80, 80)},
		// Nothing to reprocess later: the source is a user's selfie, not
		// artwork you will re-crop, and keeping every original costs storage
		// for a file nobody will ask for again.
		DiscardOriginal: true,
	})

	// A hero or cover image, where width matters more than total pixels.
	Define("cover", Profile{
		Max:        Fit(2400, 1200),
		Quality:    0.82,
		Renditions: map[string]Size{"thumb": Fill(600, 300)},
	})
}
`
}

func mediaTransformVipsGo() string {
	return `//go:build vips

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
	if px := srcW * srcH; px > p.MaxPixels {
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

	if s.Width > 0 || s.Height > 0 {
		crop := vips.InterestingNone
		if s.Crop {
			// Centre, matching what the pure-Go backend does, so the two
			// produce the same framing.
			crop = vips.InterestingCentre
		}
		if s.Crop {
			if err := img.Thumbnail(s.Width, s.Height, crop); err != nil {
				return Rendition{}, fmt.Errorf("resizing image: %w", err)
			}
		} else if img.Width() > s.Width || img.Height() > s.Height {
			// Fit, and never upscale: enlarging a small upload spends bytes on
			// pixels the source never had.
			if err := img.Thumbnail(s.Width, s.Height, vips.InterestingNone); err != nil {
				return Rendition{}, fmt.Errorf("resizing image: %w", err)
			}
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
`
}
