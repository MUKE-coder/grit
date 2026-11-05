package scaffold

import (
	"encoding/json"
	"strings"
	"testing"
)

// grit#152: which upload implementation is this project running?
//
// There are two, and they are not the same artefact. packages/upload in a
// generated project is the vendored SOURCE, resolved through workspace:* with
// no build step; @gritframework/upload on npm is a tsup build of the same code
// under a different name. So there is no tarball to compare a tree against,
// and the issue's "a test that the vendored tree equals the published tarball"
// is comparing two different shapes.
//
// What was real: the vendored copy said "version": "0.0.0", so a bug report
// about an upload named no release. Uploads are the feature with the most
// moving parts, which is the one where that matters most.
func TestTheVendoredUploadPackageNamesItsGritVersion(t *testing.T) {
	src := uploadPackageJSON("3.470.0")

	var pkg struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal([]byte(src), &pkg); err != nil {
		t.Fatalf("the generated package.json does not parse: %v\n\n%s", err, src)
	}

	if pkg.Version == "0.0.0" || pkg.Version == "" {
		t.Error(`packages/upload says version "0.0.0", so a bug report about an upload ` +
			"cannot name the release that wrote the code it is about")
	}
	if pkg.Version != "3.470.0" {
		t.Errorf("version = %q, want the Grit version it was written by", pkg.Version)
	}

	// The workspace name, not the published one. A generated project resolves
	// this from disk, where @gritframework/upload would mean nothing.
	if pkg.Name != "@repo/upload" {
		t.Errorf("name = %q, want @repo/upload", pkg.Name)
	}
}

// An empty version falls back rather than writing an empty string, because a
// package.json with "version": "" fails pnpm install outright.
func TestTheVendoredUploadPackageFallsBackToTheDefaultVersion(t *testing.T) {
	src := uploadPackageJSON("")

	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal([]byte(src), &pkg); err != nil {
		t.Fatalf("the generated package.json does not parse: %v", err)
	}
	if pkg.Version != DefaultVersion {
		t.Errorf("version = %q with no version given, want DefaultVersion (%q). "+
			`An empty string here fails "pnpm install" with an invalid package.json.`,
			pkg.Version, DefaultVersion)
	}
	if strings.Contains(src, `"version": ""`) {
		t.Error("the package.json carries an empty version")
	}
}
