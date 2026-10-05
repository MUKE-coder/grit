package scaffold

import (
	"strings"
	"testing"
)

// An .avif must upload, because three other lists already offer it.
//
// Reported by somebody adding a collection image through the admin: "File
// content does not match its declared type". Two faults, one behind the other.
//
// http.DetectContentType implements the WHATWG sniffing standard, which has no
// AVIF or HEIC matcher, so both arrive as application/octet-stream and the
// "declared an image, is not one" check rejects them. Behind that, the upload
// baseline did not list image/avif, while the "image" accept group in
// internal/files does and so does the admin's file picker: the form offered a
// file the server would not take.
func TestAVIFUploads(t *testing.T) {
	store := storageStoreGo()

	// The sniffer, and that it is consulted.
	if !strings.Contains(store, "func isobmffImageType(head []byte) string {") {
		t.Fatal("storage has no ISOBMFF sniffer, so an .avif still sniffs as octet-stream")
	}
	if !strings.Contains(store, "if t := isobmffImageType(head[:n]); t != \"\" {") {
		t.Error("DetectContentType does not consult the ISOBMFF sniffer")
	}
	// Only when Go's own sniffer gave up: a PNG must keep going through the
	// standard path.
	if !strings.Contains(store, `if detected == "application/octet-stream" {`) {
		t.Error("the ISOBMFF sniffer is not confined to what Go could not identify")
	}

	// The brands, and the priority between them. An AVIF commonly lists "mif1"
	// before "avif", so matching the first recognised brand calls it HEIF.
	avif := strings.Index(store, `{"image/avif", []string{"avif"`)
	heic := strings.Index(store, `{"image/heic", []string{"heic"`)
	heif := strings.Index(store, `{"image/heif", []string{"mif1"`)
	if avif < 0 || heic < 0 || heif < 0 {
		t.Fatal("the brand table is missing avif, heic or heif")
	}
	if !(avif < heic && heic < heif) {
		t.Error("the brands are not checked in priority order, so an AVIF listing " +
			"mif1 among its compatible brands is identified as HEIF")
	}

	// And the baseline accepts it, or the sniffing only changes which error.
	if !strings.Contains(uploadHandlerGo(), `"image/avif":      true,`) {
		t.Error("image/avif is not in defaultAllowedMIME, so the picker offers a " +
			"file the server refuses")
	}
}

// What the "image" accept group offers, the upload baseline must accept.
//
// Anything in that group that the baseline refuses is a file picker that
// accepts a file and a server that does not, which is how this was reported.
// SVG is the one exception, refused everywhere on purpose: served from your
// own origin it runs script.
func TestImageAcceptGroupAgreesWithTheUploadBaseline(t *testing.T) {
	group := filesAcceptsGo()
	baseline := uploadHandlerGo()

	start := strings.Index(group, `case "image":`)
	if start < 0 {
		t.Fatal("the image accept group is not where this test expects it")
	}
	end := strings.Index(group[start:], `case "video":`)
	if end < 0 {
		end = len(group) - start
	}
	block := group[start : start+end]

	for _, mime := range []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/avif"} {
		if !strings.Contains(block, `"`+mime+`"`) {
			continue // not offered, so nothing to agree about
		}
		if !strings.Contains(baseline, `"`+mime+`":`) {
			t.Errorf("the image accept group offers %s and defaultAllowedMIME does not "+
				"list it, so a field declared file:image offers a file the server refuses", mime)
		}
	}

	if strings.Contains(baseline, `"image/svg+xml":`) {
		t.Error("SVG is in the upload baseline; it runs script when served from your own origin")
	}
}
