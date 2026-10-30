// Package selfupdate replaces the running grit binary with a GitHub release.
// It detects OS+arch, picks the matching artifact, verifies it, extracts the
// binary, and atomically swaps it via inconshreveable/go-update (which handles
// the Windows "can't overwrite a running binary" dance for us).
//
// Verification is not optional and happens before anything is written: the
// release's checksums must carry a cosign signature from this repository's
// release workflow, and the downloaded artifact must be the one those
// checksums name. See verify.go for what that does and does not prove.
//
// Until v3.450.0 none of this happened. The download went straight to
// update.Apply, while release.yml had been publishing the signature all
// along.
package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/inconshreveable/go-update"
)

const (
	releaseAPI = "https://api.github.com/repos/MUKE-coder/grit/releases/latest"
	repoSlug   = "MUKE-coder/grit"
)

type asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
}

type release struct {
	TagName string  `json:"tag_name"`
	HTMLURL string  `json:"html_url"`
	Assets  []asset `json:"assets"`
}

// LatestVersion returns the latest published tag (without leading "v")
// alongside a flag indicating whether `current` matches it. Lets the caller
// decide what to do without committing to a full download cycle.
func LatestVersion(current string) (latest string, upToDate bool, err error) {
	rel, err := fetchLatest()
	if err != nil {
		return "", false, err
	}
	latest = strings.TrimPrefix(rel.TagName, "v")
	return latest, normalize(current) == normalize(latest), nil
}

// Options narrows what Run does. The zero value is "move to the latest
// release, refusing to go backwards", which is what `grit update` means.
type Options struct {
	// Version pins the release to install, with or without a leading "v".
	// Empty means the latest.
	Version string

	// AllowDowngrade permits installing a release older than the running
	// binary. Off by default: an update that silently moves you backwards is
	// how a fixed vulnerability comes back, and the only reason to do it on
	// purpose is to reproduce a bug.
	AllowDowngrade bool
}

// Run performs the self-update flow. `current` is the version string compiled
// into the binary (e.g. "3.24.0"). Returns nil when the update succeeded or
// the binary was already on the requested version.
func Run(current string) error { return RunWith(current, Options{}) }

// RunWith is Run with the choices made explicit.
func RunWith(current string, opts Options) error {
	var rel *release
	var err error

	if opts.Version != "" {
		want := "v" + strings.TrimPrefix(opts.Version, "v")
		fmt.Printf("Looking for %s...\n", want)
		rel, err = fetchTag(want)
		if err != nil {
			return fmt.Errorf("fetch release %s: %w", want, err)
		}
	} else {
		fmt.Println("Checking GitHub for the latest release...")
		rel, err = fetchLatest()
		if err != nil {
			return fmt.Errorf("fetch latest release: %w", err)
		}
	}

	latest := strings.TrimPrefix(rel.TagName, "v")
	if normalize(current) == normalize(latest) {
		fmt.Printf("Already on v%s\n", current)
		return nil
	}

	// Backwards is a decision, not a default. Without this, a release that is
	// yanked or a tag that moves takes every machine that runs `grit update`
	// back past a fix.
	if !opts.AllowDowngrade && isOlder(latest, current) {
		return fmt.Errorf("v%s is older than the v%s you are running.\n\n"+
			"If you mean it: grit update --version %s --allow-downgrade",
			latest, current, latest)
	}

	fmt.Printf("Updating grit  v%s → v%s\n", current, latest)
	fmt.Printf("Platform:       %s/%s\n", runtime.GOOS, runtime.GOARCH)

	a := pickAsset(rel.Assets)
	if a == nil {
		return fmt.Errorf(
			"no release artifact found for %s/%s — open an issue at https://github.com/%s/issues",
			runtime.GOOS, runtime.GOARCH, repoSlug,
		)
	}
	fmt.Printf("Downloading:    %s (%s)\n", a.Name, humanSize(a.Size))

	archive, err := download(a.DownloadURL)
	if err != nil {
		return fmt.Errorf("download artifact: %w", err)
	}

	// Before anything is written over the running binary.
	fmt.Println("Verifying signature...")
	sums, err := signedChecksums(rel)
	if err != nil {
		return err
	}
	want, ok := sums[a.Name]
	if !ok {
		return fmt.Errorf("%s is not in the release's signed checksums, so there "+
			"is nothing to check it against. Do not install it", a.Name)
	}
	if err := matchesChecksum(archive, want); err != nil {
		return err
	}
	fmt.Printf("Verified:       %s, signed by the release workflow\n", a.Name)

	bin, err := extractBinary(archive, a.Name)
	if err != nil {
		return fmt.Errorf("extract binary from archive: %w", err)
	}

	fmt.Println("Swapping binary...")
	if err := update.Apply(bytes.NewReader(bin), update.Options{}); err != nil {
		// Apply rolls back on its own when possible. If rollback also failed,
		// guide the user to recover manually.
		if rerr := update.RollbackError(err); rerr != nil {
			return fmt.Errorf("update failed AND rollback failed (%v); reinstall from %s",
				rerr, rel.HTMLURL)
		}
		return fmt.Errorf("apply update: %w", err)
	}

	fmt.Printf("\nUpdated to v%s.  Run `grit version` to verify.\n", latest)
	return nil
}

func fetchLatest() (*release, error) { return fetchRelease(releaseAPI) }

// fetchTag reads one release by its tag, for `grit update --version`.
func fetchTag(tag string) (*release, error) {
	return fetchRelease(fmt.Sprintf(
		"https://api.github.com/repos/%s/releases/tags/%s", repoSlug, tag))
}

// signedChecksums downloads SHA256SUMS and the cosign artifacts beside it and
// returns the checksums only if the signature holds.
//
// A release without them is refused rather than installed unverified. Every
// release since the workflow gained cosign has all three; one that does not
// is either very old or not a real release, and neither is worth running.
func signedChecksums(rel *release) (map[string]string, error) {
	find := func(name string) *asset {
		for i := range rel.Assets {
			if rel.Assets[i].Name == name {
				return &rel.Assets[i]
			}
		}
		return nil
	}

	var parts [3][]byte
	for i, name := range []string{"SHA256SUMS", "SHA256SUMS.sig", "SHA256SUMS.pem"} {
		a := find(name)
		if a == nil {
			return nil, fmt.Errorf("this release has no %s, so the download cannot "+
				"be verified. Install it by hand from %s if you trust it",
				name, rel.HTMLURL)
		}
		b, err := download(a.DownloadURL)
		if err != nil {
			return nil, fmt.Errorf("download %s: %w", name, err)
		}
		parts[i] = b
	}
	return verifiedChecksums(parts[0], parts[1], parts[2])
}

func fetchRelease(apiURL string) (*release, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	// Soft-rate-limit identifier so anonymous calls aren't pooled with random scripts.
	req.Header.Set("User-Agent", "grit-self-update")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("github API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, fmt.Errorf("decode release JSON: %w", err)
	}
	if rel.TagName == "" {
		return nil, fmt.Errorf("github API returned an empty tag_name")
	}
	return &rel, nil
}

// pickAsset finds the artifact whose filename matches the current OS+arch.
// Filenames follow the release.yml convention: grit-{os}-{arch}.{tar.gz|zip}.
func pickAsset(assets []asset) *asset {
	wantOS := runtime.GOOS
	wantArch := runtime.GOARCH
	// release.yml builds darwin/linux as .tar.gz and windows as .zip.
	wantExt := ".tar.gz"
	if wantOS == "windows" {
		wantExt = ".zip"
	}
	suffix := fmt.Sprintf("-%s-%s%s", wantOS, wantArch, wantExt)
	for i := range assets {
		if strings.HasSuffix(assets[i].Name, suffix) {
			return &assets[i]
		}
	}
	return nil
}

func download(url string) ([]byte, error) {
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download returned %d", resp.StatusCode)
	}
	// Bound the read so a corrupt or malicious release can't OOM the process.
	// 200 MiB is well above any grit binary we'd ship.
	const maxArtifact = 200 << 20
	return io.ReadAll(io.LimitReader(resp.Body, maxArtifact))
}

// extractBinary pulls the grit executable out of the downloaded archive.
// Returns the binary bytes ready to be handed to update.Apply.
//
// The release.yml workflow packages the binary as `grit-{goos}-{goarch}{.exe}`
// inside the archive — we accept that primary shape plus a few looser
// fallbacks (`grit`, `grit.exe`, anything ending in those suffixes) so a
// hand-packaged release still works.
func extractBinary(archive []byte, archiveName string) ([]byte, error) {
	candidates := binaryCandidates()
	if strings.HasSuffix(archiveName, ".zip") {
		return extractFromZip(archive, candidates)
	}
	return extractFromTarGz(archive, candidates)
}

// binaryCandidates returns the filenames we'll accept when scanning an
// archive, ordered most-specific-first.
func binaryCandidates() []string {
	exe := ""
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}
	platform := fmt.Sprintf("grit-%s-%s%s", runtime.GOOS, runtime.GOARCH, exe)
	bare := "grit" + exe
	return []string{platform, bare}
}

// matchCandidate returns true when the archive entry name matches any of the
// accepted candidate names — either exactly, or as a trailing path segment.
func matchCandidate(name string, candidates []string) bool {
	for _, c := range candidates {
		if name == c || strings.HasSuffix(name, "/"+c) {
			return true
		}
	}
	return false
}

func extractFromTarGz(archive []byte, candidates []string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, fmt.Errorf("gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("tar reader: %w", err)
		}
		if matchCandidate(hdr.Name, candidates) {
			return io.ReadAll(tr)
		}
	}
	return nil, fmt.Errorf("no matching binary in archive (looked for %v)", candidates)
}

func extractFromZip(archive []byte, candidates []string) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, fmt.Errorf("zip reader: %w", err)
	}
	for _, f := range zr.File {
		if matchCandidate(f.Name, candidates) {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return io.ReadAll(rc)
		}
	}
	return nil, fmt.Errorf("no matching binary in archive (looked for %v)", candidates)
}

// normalize strips leading "v" and trailing whitespace so "v3.24.0" and
// "3.24.0" compare equal.
func normalize(v string) string {
	return strings.TrimSpace(strings.TrimPrefix(v, "v"))
}

// isOlder reports whether version a is behind version b.
//
// Compared field by field as integers, because "3.9.0" sorts after "3.10.0"
// as a string and that is the comparison a downgrade guard cannot get wrong.
// A field that is not a number (a "-rc1" suffix, say) makes the answer false:
// refusing to guess is better than refusing a release somebody asked for.
func isOlder(a, b string) bool {
	as, bs := strings.Split(normalize(a), "."), strings.Split(normalize(b), ".")
	for i := 0; i < len(as) && i < len(bs); i++ {
		ai, err1 := strconv.Atoi(as[i])
		bi, err2 := strconv.Atoi(bs[i])
		if err1 != nil || err2 != nil {
			return false
		}
		if ai != bi {
			return ai < bi
		}
	}
	return len(as) < len(bs)
}

func humanSize(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// ExecutablePath returns the path of the currently-running binary, useful
// for logging and error messages. Falls back to "grit" if Executable() errors.
func ExecutablePath() string {
	p, err := os.Executable()
	if err != nil {
		return "grit"
	}
	return p
}
