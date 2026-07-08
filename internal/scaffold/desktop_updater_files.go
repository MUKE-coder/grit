package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeDesktopUpdaterFiles drops the in-app auto-updater into the desktop
// scaffold: the Go binary-swap updater bound to Wails, the version helper
// the release script bumps, and the frontend modal + banner + hook that
// drive the UI.
func writeDesktopUpdaterFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "updater.go"):                                         desktopUpdaterGo(opts),
		filepath.Join(root, "version.go"):                                         desktopVersionGo(opts),
		filepath.Join(root, "frontend", "src", "lib", "version.ts"):               desktopFrontendVersionTS(opts),
		filepath.Join(root, "frontend", "src", "hooks", "use-update-checker.ts"):  desktopUseUpdateCheckerTS(opts),
		filepath.Join(root, "frontend", "src", "components", "update-modal.tsx"):  desktopUpdateModalTSX(opts),
		filepath.Join(root, "frontend", "src", "components", "update-banner.tsx"): desktopUpdateBannerTSX(opts),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

/* ─── Go: updater.go ──────────────────────────────────────────────── */

func desktopUpdaterGo(opts DesktopOptions) string {
	tpl := `package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Updater is the in-app auto-updater bound to Wails. It is exposed to
// the React side so the UI can drive the full lifecycle:
//
//	CheckForUpdates()            → hit GitHub releases (or proxy), return latest
//	StartDownload(url, version)  → fetch the new binary to %TEMP%
//	GetProgress()                → polled by the modal for the progress bar
//	CancelDownload()             → abort the in-flight download
//	InstallAndRestart()          → swap the binary and relaunch
//
// The binary swap uses the classic Windows trick: rename the running
// .exe to .exe.old (allowed even while the file is locked, because
// Windows holds the file *open* for read/exec but lets the directory
// entry move), copy the new binary into the original location, spawn
// it, then exit. On next startup, CleanupOldOnStartup deletes the
// leftover .old (retries a few times to ride out antivirus locks).
//
// On macOS/Linux POSIX inode semantics mean you can overwrite a
// running binary directly — see the platform branches in
// InstallAndRestart for the simpler path.
type Updater struct {
	ctx context.Context

	mu      sync.RWMutex
	status  string // idle | checking | downloading | downloaded | installing | error
	target  string // path to the downloaded binary
	version string // release version being downloaded
	errMsg  string

	bytesDownloaded atomic.Int64
	bytesTotal      atomic.Int64

	cancel context.CancelFunc

	// Settings the release pipeline + ` + "`runtime.envOverride`" + ` populate.
	GithubOwner string // e.g. "MUKE-coder"
	GithubRepo  string // e.g. "my-desktop-app"
}

// UpdateRelease is the JSON the frontend receives from CheckForUpdates.
// Field names use snake_case so React side can read it without re-mapping.
type UpdateRelease struct {
	Version      string ` + "`json:\"version\"`" + `       // 1.2.3 (no leading v)
	DownloadURL  string ` + "`json:\"download_url\"`" + `  // direct .exe / .tar.gz URL
	GithubURL    string ` + "`json:\"github_url\"`" + `    // human-readable release page
	ReleasedAt   string ` + "`json:\"released_at\"`" + `   // ISO-8601
	ReleaseNotes string ` + "`json:\"release_notes\"`" + ` // markdown
	IsNewer      bool   ` + "`json:\"is_newer\"`" + `      // true if newer than the running version
}

// UpdaterState is the JSON-serialisable snapshot the modal polls each tick.
type UpdaterState struct {
	Status          string ` + "`json:\"status\"`" + `
	Version         string ` + "`json:\"version\"`" + `
	BytesDownloaded int64  ` + "`json:\"bytes_downloaded\"`" + `
	BytesTotal      int64  ` + "`json:\"bytes_total\"`" + `
	Error           string ` + "`json:\"error,omitempty\"`" + `
}

// NewUpdater constructs the updater. Wire GitHub owner / repo via env
// (UPDATER_GITHUB_OWNER, UPDATER_GITHUB_REPO) so the same binary can
// point at a fork without a recompile.
func NewUpdater() *Updater {
	return &Updater{
		status:      "idle",
		GithubOwner: envOr("UPDATER_GITHUB_OWNER", "<OWNER>"),
		GithubRepo:  envOr("UPDATER_GITHUB_REPO", "<REPO>"),
	}
}

// SetContext is called by app.startup so the updater can emit Wails
// events and drive runtime.Quit on install.
func (u *Updater) SetContext(ctx context.Context) {
	u.ctx = ctx
}

// CheckForUpdates queries the GitHub releases API for the latest tag
// and picks the right asset for the current OS / arch. Cheap (one HTTP
// call); safe to call from a polling React Query loop.
//
// For a SaaS deploy you can put your own /api/desktop/latest proxy in
// front of GitHub (private repo, custom version gating, override env
// vars) — set UPDATER_PROXY_URL to override.
func (u *Updater) CheckForUpdates() (*UpdateRelease, error) {
	u.setStatus("checking")
	defer func() {
		// CheckForUpdates is idempotent — if we were idle we stay idle;
		// in-flight downloads aren't touched.
		u.mu.Lock()
		if u.status == "checking" {
			u.status = "idle"
		}
		u.mu.Unlock()
	}()

	proxy := os.Getenv("UPDATER_PROXY_URL")
	var url string
	if proxy != "" {
		url = strings.TrimRight(proxy, "/") + "/latest"
	} else {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.GithubOwner, u.GithubRepo)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", AppName+"-updater")
	if tok := os.Getenv("UPDATER_GITHUB_TOKEN"); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github releases API: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("github releases API: HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var raw struct {
		TagName     string ` + "`json:\"tag_name\"`" + `
		HTMLURL     string ` + "`json:\"html_url\"`" + `
		Body        string ` + "`json:\"body\"`" + `
		PublishedAt string ` + "`json:\"published_at\"`" + `
		Assets      []struct {
			Name string ` + "`json:\"name\"`" + `
			URL  string ` + "`json:\"browser_download_url\"`" + `
		} ` + "`json:\"assets\"`" + `
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode release: %w", err)
	}

	asset := pickAsset(raw.Assets, runtime.GOOS, runtime.GOARCH)
	if asset == "" {
		return nil, fmt.Errorf("no asset for %s/%s in tag %s", runtime.GOOS, runtime.GOARCH, raw.TagName)
	}

	rel := &UpdateRelease{
		Version:      strings.TrimPrefix(raw.TagName, "v"),
		DownloadURL:  asset,
		GithubURL:    raw.HTMLURL,
		ReleasedAt:   raw.PublishedAt,
		ReleaseNotes: raw.Body,
	}
	rel.IsNewer = isNewer(rel.Version, AppVersion)
	return rel, nil
}

// pickAsset finds the artifact for the current platform.
//   - Windows: prefer the raw .exe (auto-updater uses this for the swap),
//     fall back to a Setup .exe if a project ships only installers.
//   - macOS: .dmg or .tar.gz
//   - Linux: .tar.gz or .AppImage
//
// Naming convention enforced by scripts/release-desktop.sh:
//
//	<AppName>-v<version>.exe                     (Windows raw, used for swap)
//	<AppName>-Setup-v<version>.exe               (Windows NSIS — manual install)
//	<AppName>-Setup-Slim-v<version>.exe          (Windows slim — online bootstrapper)
//	<AppName>-darwin-<arch>-v<version>.tar.gz    (macOS)
//	<AppName>-linux-<arch>-v<version>.tar.gz     (Linux)
func pickAsset(assets []struct {
	Name string ` + "`json:\"name\"`" + `
	URL  string ` + "`json:\"browser_download_url\"`" + `
}, goos, goarch string) string {
	matchers := map[string][]string{
		"windows": {".exe"},
		"darwin":  {"-darwin-" + goarch, ".dmg", "-mac-"},
		"linux":   {"-linux-" + goarch, ".AppImage"},
	}[goos]

	// Pass 1: prefer non-installer assets (raw binary the swap needs).
	for _, a := range assets {
		lower := strings.ToLower(a.Name)
		if strings.Contains(lower, "setup") || strings.Contains(lower, "installer") {
			continue
		}
		for _, m := range matchers {
			if strings.Contains(strings.ToLower(a.Name), strings.ToLower(m)) {
				return a.URL
			}
		}
	}
	// Pass 2: any matching asset (legacy projects without naming convention).
	for _, a := range assets {
		for _, m := range matchers {
			if strings.Contains(strings.ToLower(a.Name), strings.ToLower(m)) {
				return a.URL
			}
		}
	}
	return ""
}

// CleanupOldOnStartup removes the .exe.old left behind by a previous
// in-app update. Best-effort: if the old process is still holding the
// handle briefly we retry a handful of times. Called from main.go via
// a goroutine so app startup isn't blocked.
func (u *Updater) CleanupOldOnStartup() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	old := exe + ".old"
	if _, err := os.Stat(old); errors.Is(err, os.ErrNotExist) {
		return
	}
	for i := 0; i < 25; i++ {
		if err := os.Remove(old); err == nil || errors.Is(err, os.ErrNotExist) {
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func (u *Updater) setStatus(s string) {
	u.mu.Lock()
	u.status = s
	u.mu.Unlock()
}

// GetProgress returns the snapshot the frontend's polling loop reads.
func (u *Updater) GetProgress() UpdaterState {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return UpdaterState{
		Status:          u.status,
		Version:         u.version,
		BytesDownloaded: u.bytesDownloaded.Load(),
		BytesTotal:      u.bytesTotal.Load(),
		Error:           u.errMsg,
	}
}

// StartDownload kicks off the binary fetch in a goroutine. Returns
// immediately. Idempotent — calling it while a download is in flight
// returns an error rather than starting a second one.
func (u *Updater) StartDownload(url, version string) error {
	if url == "" {
		return errors.New("download url is empty")
	}
	if version == "" {
		return errors.New("version is empty")
	}
	u.mu.Lock()
	if u.status == "downloading" {
		u.mu.Unlock()
		return errors.New("download already in progress")
	}
	ctx, cancel := context.WithCancel(context.Background())
	u.cancel = cancel
	u.status = "downloading"
	u.version = version
	u.errMsg = ""
	u.target = ""
	u.bytesDownloaded.Store(0)
	u.bytesTotal.Store(0)
	u.mu.Unlock()

	go u.downloadWorker(ctx, url, version)
	return nil
}

func (u *Updater) downloadWorker(ctx context.Context, url, version string) {
	fail := func(err error) {
		u.mu.Lock()
		u.status = "error"
		u.errMsg = err.Error()
		u.mu.Unlock()
		if u.ctx != nil {
			wailsRuntime.EventsEmit(u.ctx, "updater:error", err.Error())
		}
	}

	ext := ".exe"
	if runtime.GOOS != "windows" {
		ext = ".tar.gz"
	}
	tmpDir := os.TempDir()
	partial := filepath.Join(tmpDir, fmt.Sprintf("%s-v%s%s.partial", AppName, version, ext))
	final := filepath.Join(tmpDir, fmt.Sprintf("%s-v%s%s", AppName, version, ext))

	_ = os.Remove(partial)
	_ = os.Remove(final)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fail(err)
		return
	}
	// No overall timeout — release binaries can be 20-200 MB on shaky
	// links. The context drives the user-cancel path.
	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		fail(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fail(fmt.Errorf("download HTTP %d", resp.StatusCode))
		return
	}
	u.bytesTotal.Store(resp.ContentLength)

	out, err := os.Create(partial)
	if err != nil {
		fail(err)
		return
	}
	prog := &progressWriter{u: u}
	if _, err := io.Copy(out, io.TeeReader(resp.Body, prog)); err != nil {
		out.Close()
		_ = os.Remove(partial)
		if ctx.Err() != nil {
			u.mu.Lock()
			u.status = "idle"
			u.mu.Unlock()
			return
		}
		fail(err)
		return
	}
	if err := out.Close(); err != nil {
		fail(err)
		return
	}
	if err := os.Rename(partial, final); err != nil {
		fail(err)
		return
	}

	u.mu.Lock()
	u.status = "downloaded"
	u.target = final
	u.mu.Unlock()
	if u.ctx != nil {
		wailsRuntime.EventsEmit(u.ctx, "updater:downloaded", final)
	}
}

type progressWriter struct{ u *Updater }

func (p *progressWriter) Write(b []byte) (int, error) {
	p.u.bytesDownloaded.Add(int64(len(b)))
	return len(b), nil
}

// CancelDownload aborts an in-flight download. Safe to call when
// nothing's running.
func (u *Updater) CancelDownload() {
	u.mu.Lock()
	if u.cancel != nil {
		u.cancel()
	}
	u.mu.Unlock()
}

// InstallAndRestart swaps the running binary for the downloaded one
// and relaunches. The actual swap strategy is OS-specific — see the
// platform comments below.
func (u *Updater) InstallAndRestart() error {
	u.mu.Lock()
	src := u.target
	curStatus := u.status
	u.mu.Unlock()

	if curStatus != "downloaded" || src == "" {
		return fmt.Errorf("no completed download to install (status=%s)", curStatus)
	}
	u.setStatus("installing")

	currentExe, err := os.Executable()
	if err != nil {
		u.setStatus("error")
		return fmt.Errorf("resolving current exe: %w", err)
	}

	if runtime.GOOS == "windows" {
		if err := u.swapWindows(currentExe, src); err != nil {
			u.setStatus("error")
			return err
		}
	} else {
		// Linux/macOS — POSIX keeps the running process's inode alive
		// when we overwrite the file at the path, so a straight copy
		// works. For Linux .AppImage / macOS .dmg the user usually
		// installs via the file manager, but we still support a raw
		// binary swap for portable installs.
		if err := copyFile(src, currentExe); err != nil {
			u.setStatus("error")
			return fmt.Errorf("copying new binary: %w", err)
		}
	}

	_ = os.Remove(src)

	// Spawn new version. cmd.Start spawns and returns — the new
	// process survives our exit because the OS doesn't tie child
	// lifetime to the parent unless we explicitly Wait().
	cmd := exec.Command(currentExe)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	if err := cmd.Start(); err != nil {
		u.setStatus("error")
		return fmt.Errorf("starting new version: %w", err)
	}

	// Small delay so the modal's "Installing…" state paints before
	// the window vanishes — feels less abrupt to the operator.
	go func() {
		time.Sleep(250 * time.Millisecond)
		if u.ctx != nil {
			wailsRuntime.Quit(u.ctx)
		} else {
			os.Exit(0)
		}
	}()
	return nil
}

func (u *Updater) swapWindows(currentExe, src string) error {
	oldPath := currentExe + ".old"
	_ = os.Remove(oldPath)

	// Windows allows renaming a locked .exe — the kernel object stays
	// mapped, the directory entry moves.
	if err := os.Rename(currentExe, oldPath); err != nil {
		return fmt.Errorf("renaming current exe: %w", err)
	}
	// Copy (not rename) because %TEMP% is often on a different volume
	// than %LOCALAPPDATA% — Rename across volumes fails with EXDEV.
	// If the copy fails we roll back so the operator isn't stranded.
	if err := copyFile(src, currentExe); err != nil {
		_ = os.Rename(oldPath, currentExe)
		return fmt.Errorf("copying new exe: %w", err)
	}
	return nil
}

// copyFile is a no-frills copy. Used because %TEMP% (where the
// downloaded binary lives) and the install dir are usually on
// different volumes — os.Rename would fail with EXDEV.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// isNewer returns true when "latest" semver > "current". Strips a
// leading "v" and anything after the first "-" (so "1.2.3-rc1"
// compares as "1.2.3").
func isNewer(latest, current string) bool {
	a := strings.SplitN(strings.TrimPrefix(latest, "v"), "-", 2)[0]
	b := strings.SplitN(strings.TrimPrefix(current, "v"), "-", 2)[0]
	ap := splitInts(a)
	bp := splitInts(b)
	for i := 0; i < 3; i++ {
		if i >= len(ap) {
			ap = append(ap, 0)
		}
		if i >= len(bp) {
			bp = append(bp, 0)
		}
		if ap[i] != bp[i] {
			return ap[i] > bp[i]
		}
	}
	return false
}

func splitInts(v string) []int {
	parts := strings.Split(v, ".")
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		var n int
		_, _ = fmt.Sscanf(p, "%d", &n)
		out = append(out, n)
	}
	return out
}
`
	tpl = strings.ReplaceAll(tpl, "<OWNER>", opts.ProjectName)
	tpl = strings.ReplaceAll(tpl, "<REPO>", opts.ProjectName)
	return tpl
}

/* ─── Go: version.go ──────────────────────────────────────────────── */

func desktopVersionGo(opts DesktopOptions) string {
	return fmt.Sprintf(`package main

// AppName + AppVersion are the canonical version metadata.
// scripts/release-desktop.sh bumps AppVersion + wails.json's
// productVersion + frontend/src/lib/version.ts together so the
// running binary, the installer filename, and the UI always agree.
//
// Don't import this from anything that the frontend touches — the
// React side reads APP_VERSION from frontend/src/lib/version.ts.
const (
	AppName    = %q
	AppVersion = "0.1.0"
)
`, opts.ProjectName)
}

/* ─── Frontend: lib/version.ts ────────────────────────────────────── */

func desktopFrontendVersionTS(opts DesktopOptions) string {
	_ = opts
	return `/**
 * APP_VERSION is the build-time injected version string. vite.config
 * reads ../wails.json's productVersion and defines __APP_VERSION__
 * at build time, so this file never needs to be hand-edited — the
 * release script bumps wails.json once and the whole frontend picks
 * it up.
 */
declare const __APP_VERSION__: string

export const APP_VERSION: string = (typeof __APP_VERSION__ !== 'undefined' ? __APP_VERSION__ : '0.0.0')

/**
 * isNewerVersion returns true when 'latest' semver > 'current'.
 * Strips a leading 'v' and anything after the first '-' (so
 * "1.2.3-rc1" compares as "1.2.3").
 */
export function isNewerVersion(latest: string, current: string): boolean {
  const norm = (v: string) => v.replace(/^v/, '').split('-')[0]
  const a = norm(latest).split('.').map(Number)
  const b = norm(current).split('.').map(Number)
  for (let i = 0; i < 3; i++) {
    const ai = a[i] ?? 0
    const bi = b[i] ?? 0
    if (ai !== bi) return ai > bi
  }
  return false
}
`
}

/* ─── Frontend: hooks/use-update-checker.ts ───────────────────────── */

func desktopUseUpdateCheckerTS(opts DesktopOptions) string {
	_ = opts
	return `import { useQuery } from '@tanstack/react-query'
import { APP_VERSION, isNewerVersion } from '@/lib/version'

// Generated by Wails at build time — see wailsjs/go/main/Updater.d.ts.
// We type just the calls we use to keep the hook framework-agnostic.
declare global {
  interface Window {
    go: {
      main: {
        Updater: {
          CheckForUpdates(): Promise<{
            version: string
            download_url: string
            github_url: string
            released_at: string
            release_notes: string
            is_newer: boolean
          }>
        }
      }
    }
  }
}

export interface AvailableUpdate {
  version: string
  downloadUrl: string
  githubUrl: string
  releasedAt: string
  releaseNotes: string
}

/**
 * useUpdateChecker polls the Go updater's CheckForUpdates every 6h
 * (cheap — one HTTP call to GitHub). Returns isUpdateAvailable +
 * the latest release metadata. The caller wires it to the modal +
 * banner.
 */
export function useUpdateChecker() {
  const q = useQuery({
    queryKey: ['updater', 'latest'],
    queryFn: async () => window.go.main.Updater.CheckForUpdates(),
    // 6h between automatic checks. The "Check for updates" button
    // in Settings should call refetch() to force-bypass this.
    refetchInterval: 6 * 60 * 60 * 1000,
    refetchOnWindowFocus: false,
    staleTime: 60 * 60 * 1000,
    retry: 1,
  })

  const isUpdateAvailable =
    !!q.data && q.data.is_newer && isNewerVersion(q.data.version, APP_VERSION)

  const update: AvailableUpdate | null = isUpdateAvailable && q.data
    ? {
        version: q.data.version,
        downloadUrl: q.data.download_url,
        githubUrl: q.data.github_url,
        releasedAt: q.data.released_at,
        releaseNotes: q.data.release_notes,
      }
    : null

  return {
    update,
    isUpdateAvailable,
    isChecking: q.isFetching,
    error: q.error as Error | null,
    checkNow: () => q.refetch(),
  }
}
`
}

/* ─── Frontend: components/update-modal.tsx ───────────────────────── */

func desktopUpdateModalTSX(opts DesktopOptions) string {
	_ = opts
	return `import { useEffect, useState } from 'react'
import { Download, Check, X, RotateCcw } from 'lucide-react'
import { APP_VERSION } from '@/lib/version'
import type { AvailableUpdate } from '@/hooks/use-update-checker'

declare global {
  interface Window {
    go: {
      main: {
        Updater: {
          StartDownload(url: string, version: string): Promise<void>
          GetProgress(): Promise<{
            status: 'idle' | 'checking' | 'downloading' | 'downloaded' | 'installing' | 'error'
            version: string
            bytes_downloaded: number
            bytes_total: number
            error?: string
          }>
          CancelDownload(): Promise<void>
          InstallAndRestart(): Promise<void>
        }
      }
    }
  }
}

interface Props {
  update: AvailableUpdate
  onClose: () => void
}

/**
 * UpdateModal — kicks off the download on mount, polls progress every
 * 500ms, swaps to a green "Install & restart" button when complete,
 * and surfaces errors with a "Try again" recovery.
 *
 * The release-notes section is intentionally bounded (max 3 bullets
 * collapsed, scrollable "Show all" expansion) so a chatty changelog
 * can never hide the action button below the fold.
 */
export function UpdateModal({ update, onClose }: Props) {
  const [progress, setProgress] = useState({
    status: 'idle' as 'idle' | 'checking' | 'downloading' | 'downloaded' | 'installing' | 'error',
    bytesDownloaded: 0,
    bytesTotal: 0,
    error: '',
  })
  const [showAllNotes, setShowAllNotes] = useState(false)

  // Auto-start the download on mount.
  useEffect(() => {
    let mounted = true
    window.go.main.Updater.StartDownload(update.downloadUrl, update.version).catch((err) => {
      if (!mounted) return
      setProgress((p) => ({ ...p, status: 'error', error: String(err) }))
    })
    return () => {
      mounted = false
    }
  }, [update.downloadUrl, update.version])

  // Poll progress every 500ms while the modal is open.
  useEffect(() => {
    let mounted = true
    const id = setInterval(async () => {
      try {
        const p = await window.go.main.Updater.GetProgress()
        if (!mounted) return
        setProgress({
          status: p.status,
          bytesDownloaded: p.bytes_downloaded,
          bytesTotal: p.bytes_total,
          error: p.error || '',
        })
      } catch {
        // Polling errors are non-fatal — try again next tick.
      }
    }, 500)
    return () => {
      mounted = false
      clearInterval(id)
    }
  }, [])

  const pct = progress.bytesTotal > 0 ? (progress.bytesDownloaded / progress.bytesTotal) * 100 : 0
  const mbDown = (progress.bytesDownloaded / (1024 * 1024)).toFixed(1)
  const mbTotal = progress.bytesTotal > 0 ? (progress.bytesTotal / (1024 * 1024)).toFixed(1) : '?'

  const notes = update.releaseNotes
    .split('\n')
    .map((l) => l.trim())
    .filter((l) => l.startsWith('-') || l.startsWith('*'))
  const shownNotes = showAllNotes ? notes : notes.slice(0, 3)

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="w-[480px] max-h-[90vh] flex flex-col rounded-xl border border-border bg-card shadow-2xl overflow-hidden">
        {/* Header */}
        <div className="flex items-center justify-between px-5 py-4 border-b border-border">
          <div>
            <h2 className="text-base font-semibold">Update available</h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              v{APP_VERSION} <span className="mx-1">→</span> v{update.version}
            </p>
          </div>
          <button
            onClick={onClose}
            className="text-muted-foreground hover:text-foreground rounded p-1"
            aria-label="Close"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {/* Release notes */}
        {notes.length > 0 && (
          <div className="px-5 py-3 border-b border-border max-h-[180px] overflow-y-auto">
            <p className="text-xs font-mono uppercase tracking-wider text-muted-foreground mb-2">
              What's new
            </p>
            <ul className="space-y-1 text-sm text-foreground/85">
              {shownNotes.map((n, i) => (
                <li key={i} className="leading-snug">
                  {n.replace(/^[-*]\s*/, '• ')}
                </li>
              ))}
            </ul>
            {notes.length > 3 && (
              <button
                onClick={() => setShowAllNotes((v) => !v)}
                className="mt-2 text-xs text-primary hover:underline"
              >
                {showAllNotes ? 'Show less' : 'Show all (' + notes.length + ')'}
              </button>
            )}
          </div>
        )}

        {/* Progress */}
        <div className="px-5 py-4 flex-1">
          {progress.status === 'error' && (
            <div className="rounded-md border border-destructive/40 bg-destructive/10 p-3 mb-3 text-sm text-destructive flex items-start gap-2">
              <X className="h-4 w-4 mt-0.5 shrink-0" />
              <div>
                <p className="font-medium">Download failed</p>
                <p className="text-xs mt-1 text-foreground/80">{progress.error}</p>
              </div>
            </div>
          )}

          {progress.status === 'downloading' && (
            <>
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-foreground/80 flex items-center gap-1.5">
                  <Download className="h-3.5 w-3.5" /> Downloading…
                </span>
                <span className="font-mono text-xs text-muted-foreground">
                  {mbDown} / {mbTotal} MB
                </span>
              </div>
              <div className="h-2 rounded-full bg-muted overflow-hidden">
                <div
                  className="h-full bg-primary transition-all duration-300"
                  style={{ width: pct + '%' }}
                />
              </div>
            </>
          )}

          {progress.status === 'downloaded' && (
            <div className="flex items-center gap-2 text-sm text-success">
              <Check className="h-4 w-4" />
              <span>Downloaded — ready to install.</span>
            </div>
          )}

          {progress.status === 'installing' && (
            <div className="text-sm text-foreground/85">
              Installing — the app will restart momentarily…
            </div>
          )}
        </div>

        {/* Action buttons */}
        <div className="flex items-center justify-end gap-2 px-5 py-3 border-t border-border bg-muted/20">
          {progress.status === 'downloading' && (
            <button
              onClick={() => window.go.main.Updater.CancelDownload()}
              className="rounded-md px-3 py-1.5 text-sm border border-border hover:bg-muted"
            >
              Cancel
            </button>
          )}
          {progress.status === 'downloaded' && (
            <button
              onClick={() => window.go.main.Updater.InstallAndRestart()}
              className="rounded-md px-4 py-1.5 text-sm bg-primary text-primary-foreground hover:bg-primary/90 flex items-center gap-1.5"
            >
              <RotateCcw className="h-3.5 w-3.5" />
              Install &amp; restart
            </button>
          )}
          {progress.status === 'error' && (
            <button
              onClick={() => {
                setProgress({ status: 'idle', bytesDownloaded: 0, bytesTotal: 0, error: '' })
                window.go.main.Updater.StartDownload(update.downloadUrl, update.version)
              }}
              className="rounded-md px-3 py-1.5 text-sm bg-primary text-primary-foreground hover:bg-primary/90"
            >
              Try again
            </button>
          )}
        </div>
      </div>
    </div>
  )
}
`
}

/* ─── Frontend: components/update-banner.tsx ──────────────────────── */

func desktopUpdateBannerTSX(opts DesktopOptions) string {
	_ = opts
	return `import { useState } from 'react'
import { Download, X } from 'lucide-react'
import { useUpdateChecker } from '@/hooks/use-update-checker'
import { UpdateModal } from './update-modal'

/**
 * UpdateBanner sits in the app shell and shows a dismissible
 * notification when an update is available. Clicking "Update now"
 * opens the download / install modal.
 *
 * Mount once in the top-level layout. Dismissal is per-session — the
 * banner re-appears on next launch if the update is still pending.
 */
export function UpdateBanner() {
  const { update, isUpdateAvailable } = useUpdateChecker()
  const [dismissed, setDismissed] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)

  if (!isUpdateAvailable || !update || dismissed) return null

  return (
    <>
      <div className="flex items-center justify-between gap-3 px-4 py-2 border-b border-primary/30 bg-primary/[0.06] text-sm">
        <div className="flex items-center gap-2 min-w-0">
          <Download className="h-3.5 w-3.5 text-primary shrink-0" />
          <span className="text-foreground/85 truncate">
            Update available: <strong>v{update.version}</strong>
          </span>
        </div>
        <div className="flex items-center gap-2 shrink-0">
          <button
            onClick={() => setModalOpen(true)}
            className="rounded-md px-3 py-1 text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90"
          >
            Update now
          </button>
          <button
            onClick={() => setDismissed(true)}
            className="text-muted-foreground hover:text-foreground rounded p-1"
            aria-label="Dismiss"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        </div>
      </div>

      {modalOpen && <UpdateModal update={update} onClose={() => setModalOpen(false)} />}
    </>
  )
}
`
}
