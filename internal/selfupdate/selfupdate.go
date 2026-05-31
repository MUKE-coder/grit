// Package selfupdate replaces the running grit binary with the latest GitHub
// release. It detects OS+arch, picks the matching release artifact, downloads
// it, extracts the binary, and atomically swaps it via inconshreveable/go-update
// (which handles the Windows "can't overwrite a running binary" dance for us).
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

// Run performs the self-update flow. `current` is the version string compiled
// into the binary (e.g. "3.24.0"). Returns nil when the update succeeded or
// the binary was already on the latest version.
func Run(current string) error {
	fmt.Println("Checking GitHub for the latest release...")
	rel, err := fetchLatest()
	if err != nil {
		return fmt.Errorf("fetch latest release: %w", err)
	}

	latest := strings.TrimPrefix(rel.TagName, "v")
	if normalize(current) == normalize(latest) {
		fmt.Printf("Already on the latest version: v%s\n", current)
		return nil
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

func fetchLatest() (*release, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(http.MethodGet, releaseAPI, nil)
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
