package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Contact-app review M27: the public share page at /forms/[token] asked the API
// for the share in a useEffect, with a raw axios import, so somebody opening a
// link had to wait for the bundle and then for a round trip before the form
// existed. It is a server component now, and the interactive half moved to
// public-form.tsx next to it.
//
// apps/web has no general upgrade path, so this is written here rather than
// picked up from the file map. The manifest guard still applies: a page
// somebody edited is left alone and reported as a conflict.

// publicFormUseEffectMarker is the fetch the old page ran on mount.
const publicFormUseEffectMarker = `axios.get(API_URL + "/api/public/forms/"`

// repairPublicFormPage moves the share fetch onto the server.
func repairPublicFormPage(root string, opts Options) error {
	dir := filepath.Join(root, "apps", "web", "app", "forms", "[token]")
	page := filepath.Join(dir, "page.tsx")
	if !fileExists(page) {
		return nil
	}
	body, err := os.ReadFile(page)
	if err != nil {
		return fmt.Errorf("reading %s: %w", page, err)
	}
	if !strings.Contains(string(body), publicFormUseEffectMarker) {
		// Already the server component, or a page somebody replaced entirely.
		return nil
	}
	if err := writeFile(filepath.Join(dir, "public-form.tsx"), webPublicFormClient()); err != nil {
		return err
	}
	if err := writeFile(page, webPublicFormPage()); err != nil {
		return err
	}
	fmt.Printf("  ✓ apps/web/app/forms/[token]: the share is read on the server, not in a useEffect\n")
	return nil
}
