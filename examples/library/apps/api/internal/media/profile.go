package media

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
//	1600  covers a 2x retina display at a typical content width. Storing more
//	      is paying to keep pixels no browser will draw.
//	0.82  the point at which JPEG is visually indistinguishable from source.
//	      Below about 0.75 artifacts appear on gradients and skin.
//	400   the thumbnail an admin table row or a card grid needs at 2x.
//
// Measured on a 5.3 MB phone photo: 147 KB out, including the thumbnail.
func DefaultProfile() Profile {
	return Profile{
		Max:        Fit(1600, 1600),
		Quality:    0.82,
		Format:     Auto,
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
	Name       string                `json:"name"`
	MaxWidth   int                   `json:"max_width"`
	MaxHeight  int                   `json:"max_height"`
	Crop       bool                  `json:"crop"`
	Quality    float64               `json:"quality"`
	Format     string                `json:"format"`
	MaxPixels  int                   `json:"max_pixels"`
	Renditions map[string]PublicSize `json:"renditions,omitempty"`
}

// PublicSize is one rendition target on the wire.
type PublicSize struct {
	Width  int  `json:"width"`
	Height int  `json:"height"`
	Crop   bool `json:"crop"`
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
