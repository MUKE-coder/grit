package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeDesktopClientFiles scaffolds `apps/desktop/` — a Wails window that
// consumes the shared Go API over HTTP (same API the web/admin/expo apps use).
//
// This is different from `grit new-desktop` (standalone offline-first app):
//   - grit new-desktop:  Wails + embedded Go + local SQLite (offline)
//   - --desktop flag:    Wails shell + remote API (always-online client)
//
// Frontend stack:  React + Vite + TanStack Router + Tailwind + TanStack Query.
// Wails bindings:  only native OS stuff (window controls, file dialogs, keychain).
// Token storage:   OS keychain via github.com/99designs/keyring (not localStorage).
// Design:          premium desktop UX — frameless, command palette, keyboard shortcuts.
func writeDesktopClientFiles(root string, opts Options) error {
	desktopRoot := filepath.Join(root, "apps", "desktop")
	// Desktop module is a peer of apps/api, not a child.
	// opts.Module() returns "<project>/apps/api" for monorepo architectures,
	// so we build the desktop path from the project name directly.
	module := opts.ProjectName + "/apps/desktop"

	files := map[string]string{
		// Go side (minimal — native OS only)
		// keychain.go stays at the top level (same `package main` as main.go/app.go).
		filepath.Join(desktopRoot, "main.go"):     desktopClientMainGo(opts),
		filepath.Join(desktopRoot, "app.go"):      desktopClientAppGo(),
		filepath.Join(desktopRoot, "keychain.go"): desktopClientKeychainGo(opts),
		filepath.Join(desktopRoot, "go.mod"):      desktopClientGoMod(module),
		filepath.Join(desktopRoot, "wails.json"):  desktopClientWailsJSON(opts),
		filepath.Join(desktopRoot, ".gitignore"):  desktopClientGitignore(),
		filepath.Join(desktopRoot, "README.md"):   desktopClientReadme(opts),

		// Frontend config
		filepath.Join(desktopRoot, "frontend", "package.json"):         desktopClientPackageJSON(opts),
		filepath.Join(desktopRoot, "frontend", "tsconfig.json"):        desktopClientTSConfig(),
		filepath.Join(desktopRoot, "frontend", "tsconfig.node.json"):   desktopClientTSConfigNode(),
		filepath.Join(desktopRoot, "frontend", "vite.config.ts"):       desktopClientViteConfig(),
		filepath.Join(desktopRoot, "frontend", "tailwind.config.ts"):   desktopClientTailwindConfig(),
		filepath.Join(desktopRoot, "frontend", "postcss.config.js"):    desktopClientPostCSSConfig(),
		filepath.Join(desktopRoot, "frontend", "index.html"):           desktopClientIndexHTML(opts),
		filepath.Join(desktopRoot, "frontend", "src", "main.tsx"):      desktopClientMainTSX(),
		filepath.Join(desktopRoot, "frontend", "src", "globals.css"):   desktopClientGlobalsCSS(),
		filepath.Join(desktopRoot, "frontend", "src", "vite-env.d.ts"): desktopClientViteEnvDTS(),

		// Routes (TanStack Router file-based)
		filepath.Join(desktopRoot, "frontend", "src", "routes", "__root.tsx"):            desktopClientRootRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "index.tsx"):             desktopClientIndexRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "auth.tsx"):             desktopClientAuthLayoutRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "auth", "login.tsx"):    desktopClientLoginRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "auth", "register.tsx"): desktopClientRegisterRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "app.tsx"):              desktopClientAppLayoutRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "app", "index.tsx"):     desktopClientDashboardRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "app", "profile.tsx"):   desktopClientProfileRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "app", "settings.tsx"):  desktopClientSettingsRoute(),

		// Layout components
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "title-bar.tsx"):       desktopClientTitleBar(opts),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "sidebar.tsx"):         desktopClientSidebarV2(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "topbar.tsx"):          desktopClientTopbar(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "page-header.tsx"):     desktopClientPageHeader(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "command-palette.tsx"): desktopClientCommandPalette(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "app-shell.tsx"):       desktopClientAppShell(),

		// UI components
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "button.tsx"): desktopClientButton(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "input.tsx"):  desktopClientInput(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "card.tsx"):   desktopClientCard(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "dialog.tsx"): desktopClientDialog(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "avatar.tsx"): desktopClientAvatar(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "badge.tsx"):  desktopClientBadge(),

		// Desktop primitives — list/detail container, form fields, filter chips.
		// These solve the "every desktop CRUD page reinvents the same 200 LOC" problem.
		filepath.Join(desktopRoot, "frontend", "src", "components", "two-pane.tsx"):    desktopClientTwoPane(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "form.tsx"):        desktopClientForm(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "filter-chip.tsx"): desktopClientFilterChip(),

		// Frontend stdlib + form primitives (v3.15)
		filepath.Join(desktopRoot, "frontend", "src", "lib", "format.ts"):                    desktopClientFormatLib(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "nav-config.ts"):                desktopClientNavConfig(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "currency-field.tsx"):    desktopClientCurrencyField(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "searchable-select.tsx"): desktopClientSearchableSelect(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "date-field.tsx"):        desktopClientDateField(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "drawer.tsx"):            desktopClientDrawer(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "status-badge.tsx"):      desktopClientStatusBadge(),

		// Lib (API client, auth, utils, shortcuts, Wails bridge)
		filepath.Join(desktopRoot, "frontend", "src", "lib", "api-client.ts"):      desktopClientApiClientTS(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "auth-provider.tsx"):  desktopClientAuthProvider(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "query-client.ts"):    desktopClientQueryClient(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "theme-provider.tsx"): desktopClientThemeProvider(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "use-shortcuts.ts"):   desktopClientUseShortcuts(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "wails-bridge.ts"):    desktopClientWailsBridge(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "utils.ts"):           desktopClientUtils(),

		// Hooks
		filepath.Join(desktopRoot, "frontend", "src", "hooks", "use-auth.ts"):          desktopClientUseAuth(),
		filepath.Join(desktopRoot, "frontend", "src", "hooks", "use-online-status.ts"): desktopClientUseOnlineStatus(),
		filepath.Join(desktopRoot, "frontend", "src", "hooks", "use-realtime.ts"):      desktopClientUseRealtime(),
		filepath.Join(desktopRoot, "frontend", "src", "hooks", "use-sync.ts"):          desktopClientUseSync(),

		// Sync client wrappers + UI
		filepath.Join(desktopRoot, "frontend", "src", "lib", "sync-client.ts"):             desktopClientSyncClientTS(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "sync-button.tsx"):     desktopClientSyncButton(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "pending-changes.tsx"): desktopClientPendingChanges(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "conflict-dialog.tsx"): desktopClientConflictDialog(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "offline-mode-toggle.tsx"): desktopClientOfflineToggle(),
		filepath.Join(desktopRoot, "frontend", "src", "hooks", "use-sync-status.ts"):          desktopClientUseSyncStatus(),

		// Theming — shares packages/shared/themes.ts with the admin panel so
		// --theme=atlas|aurora|pulse styles both apps identically.
		filepath.Join(desktopRoot, "frontend", "src", "lib", "theme-tokens.ts"):             desktopClientThemeTokens(opts),
		filepath.Join(desktopRoot, "frontend", "src", "components", "auth", "AuthShell.tsx"):       desktopClientAuthShell(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "auth", "AuthField.tsx"):       desktopClientAuthField(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "auth", "BrandMark.tsx"):       desktopClientBrandMark(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "auth", "AtlasAuthShell.tsx"):  desktopAtlasAuthShell(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "auth", "AuroraAuthShell.tsx"): desktopAuroraAuthShell(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "auth", "PulseAuthShell.tsx"):  desktopPulseAuthShell(),

		// Realtime client (WebSocket + reconnect + EventTarget bus)
		filepath.Join(desktopRoot, "frontend", "src", "lib", "realtime.ts"): desktopClientRealtimeTS(),

		// Offline-first sync engine (local SQLite + outbox + push/pull orchestration)
		filepath.Join(desktopRoot, "sync", "engine.go"): desktopSyncEngineGo(),
		filepath.Join(desktopRoot, "sync", "outbox.go"): desktopSyncOutboxGo(),
		filepath.Join(desktopRoot, "sync", "local.go"):  desktopSyncLocalGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		content = strings.ReplaceAll(content, "{{PROJECT_NAME}}", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// ═══════════════════════════════════════════════════════════════════
// Go side — minimal, native OS only
// ═══════════════════════════════════════════════════════════════════

func desktopClientMainGo(opts Options) string {
	return `package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "` + opts.ProjectName + `",
		Width:     1280,
		Height:    800,
		MinWidth:  1000,
		MinHeight: 700,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 10, G: 10, B: 15, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		// Drag region: any HTML element with style "--wails-draggable: drag"
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",

		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			BackdropType:                      windows.Mica,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			About: &mac.AboutInfo{
				Title:   "` + opts.ProjectName + `",
				Message: "Built with Grit",
			},
		},
		Linux: &linux.Options{
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
`
}

func desktopClientAppGo() string {
	return `package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	goruntime "runtime"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"{{MODULE}}/sync"
)

// syncTables is the single source of truth for which models the background
// auto-sync loop and manual Sync cover. The frontend reads it via
// GetSyncTables(), and ~grit generate resource~ appends new resources at the
// marker below — so a new resource joins offline sync automatically.
var syncTables = []string{
	"users",
	"uploads",
	// grit:sync-tables
}

// App exposes native OS methods + the offline sync engine to the React
// frontend via Wails bindings. Online business logic still goes through
// HTTP to the shared API; offline-first writes go through the Local*
// methods below, which queue them in the sync engine's outbox until the
// user explicitly clicks "Sync".
type App struct {
	ctx    context.Context
	kc     *Keychain
	sync   *sync.Engine
	apiURL string
}

func NewApp() *App {
	return &App{
		kc: NewKeychain(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Configure the sync engine. The local SQLite lives under the OS-
	// standard user data dir so it survives app updates.
	a.apiURL = os.Getenv("VITE_API_URL")
	if a.apiURL == "" {
		a.apiURL = "http://localhost:8080/api"
	}
	dataDir, err := os.UserConfigDir()
	if err != nil {
		dataDir, _ = os.UserHomeDir()
	}
	dbPath := filepath.Join(dataDir, "{{PROJECT_NAME}}", "sync.db")
	engine, err := sync.Open(dbPath, a.apiURL, func() (string, error) {
		v, _ := a.kc.Get("access_token")
		return v, nil
	})
	if err != nil {
		// Don't crash the whole app; offline features will simply fail
		// closed. The frontend can warn the user via SyncStatus().
		fmt.Println("[sync] open failed:", err)
		return
	}
	a.sync = engine

	// Start the background mirror loop. While the server is reachable and the
	// user hasn't chosen "Work offline", this pulls fresh server data into the
	// local mirror and pushes any queued offline edits every 30s — so going
	// offline always has recent data, and coming back online auto-reconciles.
	engine.StartAutoSync(syncTables, 30*time.Second)
}

// ─── Token storage (OS keychain) ──────────────────────────────────

func (a *App) SetToken(key, value string) error {
	return a.kc.Set(key, value)
}

func (a *App) GetToken(key string) string {
	v, _ := a.kc.Get(key)
	return v
}

func (a *App) DeleteToken(key string) error {
	return a.kc.Delete(key)
}

// ─── Window controls ──────────────────────────────────────────────

func (a *App) MinimiseWindow()   { wailsruntime.WindowMinimise(a.ctx) }
func (a *App) MaximiseWindow()   { wailsruntime.WindowMaximise(a.ctx) }
func (a *App) UnmaximiseWindow() { wailsruntime.WindowUnmaximise(a.ctx) }
func (a *App) ToggleMaximise()   { wailsruntime.WindowToggleMaximise(a.ctx) }
func (a *App) CloseWindow()      { wailsruntime.Quit(a.ctx) }
func (a *App) IsMaximised() bool { return wailsruntime.WindowIsMaximised(a.ctx) }

// ─── File dialogs ─────────────────────────────────────────────────

func (a *App) OpenFileDialog(title string) (string, error) {
	return wailsruntime.OpenFileDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: title,
	})
}

func (a *App) SaveFileDialog(title, defaultFilename string) (string, error) {
	return wailsruntime.SaveFileDialog(a.ctx, wailsruntime.SaveDialogOptions{
		Title:           title,
		DefaultFilename: defaultFilename,
	})
}

// ─── System info ──────────────────────────────────────────────────

// GetPlatform returns "darwin" | "windows" | "linux".
func (a *App) GetPlatform() string {
	return goruntime.GOOS
}

// GetAppVersion returns the build version.
func (a *App) GetAppVersion() string {
	return "0.1.0"
}

// ─── Offline sync (local-first CRUD + push/pull orchestration) ────
//
// Frontend usage:
//   await LocalCreate("buildings", "", { name: "Foo" })   // generates UUID
//   await LocalUpdate("buildings", id, { name: "Bar" })
//   await LocalDelete("buildings", id)
//   const items = await LocalList("buildings")
//   const result = await Sync(["buildings", "tenants"])
//   if (result.conflicts > 0) { /* open conflict dialog */ }
//
// Reads come from the local SQLite mirror, populated by Sync's pull
// phase. Writes are queued in the outbox and only hit the network when
// the user clicks Sync.

func (a *App) LocalCreate(table, id string, data map[string]interface{}) error {
	if a.sync == nil {
		return fmt.Errorf("sync engine not initialized")
	}
	return a.sync.LocalCreate(table, id, data)
}

func (a *App) LocalUpdate(table, id string, data map[string]interface{}) error {
	if a.sync == nil {
		return fmt.Errorf("sync engine not initialized")
	}
	return a.sync.LocalUpdate(table, id, data)
}

func (a *App) LocalDelete(table, id string) error {
	if a.sync == nil {
		return fmt.Errorf("sync engine not initialized")
	}
	return a.sync.LocalDelete(table, id)
}

func (a *App) LocalGet(table, id string) (map[string]interface{}, error) {
	if a.sync == nil {
		return nil, fmt.Errorf("sync engine not initialized")
	}
	return a.sync.LocalGet(table, id)
}

func (a *App) LocalList(table string) ([]map[string]interface{}, error) {
	if a.sync == nil {
		return nil, fmt.Errorf("sync engine not initialized")
	}
	return a.sync.LocalList(table)
}

// Sync runs Pull (for the listed tables) then Push. Returns counts
// for the UI to render.
func (a *App) Sync(tables []string) (*sync.SyncResult, error) {
	if a.sync == nil {
		return nil, fmt.Errorf("sync engine not initialized")
	}
	return a.sync.Sync(tables)
}

// GetSyncTables returns the models covered by offline sync. The frontend uses
// this instead of a hardcoded list, so a newly generated resource is picked up
// automatically.
func (a *App) GetSyncTables() []string { return syncTables }

// SetOfflineMode toggles the manual "Work offline" switch shown in the
// dashboard. Turning it OFF (going back online) triggers an immediate
// background reconcile so queued edits push and fresh data pulls right away.
func (a *App) SetOfflineMode(offline bool) error {
	if a.sync == nil {
		return fmt.Errorf("sync engine not initialized")
	}
	if err := a.sync.SetForceOffline(offline); err != nil {
		return err
	}
	if !offline {
		go func() { _, _ = a.sync.SyncNow() }()
	}
	return nil
}

// GetSyncStatus returns the reachable/offline/pending snapshot for the
// dashboard indicator.
func (a *App) GetSyncStatus() sync.SyncStatus {
	if a.sync == nil {
		return sync.SyncStatus{}
	}
	return a.sync.Status()
}

// SyncNow forces an immediate Pull+Push (the dashboard "Sync now" action).
func (a *App) SyncNow() (*sync.SyncResult, error) {
	if a.sync == nil {
		return nil, fmt.Errorf("sync engine not initialized")
	}
	return a.sync.SyncNow()
}

// PendingCount returns the number of unpushed entries — wired to the
// title-bar Sync button badge.
func (a *App) PendingCount() (int64, error) {
	if a.sync == nil {
		return 0, nil
	}
	return a.sync.PendingCount()
}

// GetPendingChanges returns the full outbox for the review panel.
func (a *App) GetPendingChanges() ([]sync.Outbox, error) {
	if a.sync == nil {
		return nil, fmt.Errorf("sync engine not initialized")
	}
	return a.sync.GetPendingChanges()
}

// ResolveConflict accepts the user's merged data for a conflicted entry.
// serverVersion is the version the user is overwriting (so the next
// push's optimistic-lock check matches). The engine clears HasConflict
// so the entry is replayed on the next Sync.
func (a *App) ResolveConflict(table, entityID string, mergedData map[string]interface{}, serverVersion int) error {
	if a.sync == nil {
		return fmt.Errorf("sync engine not initialized")
	}
	return a.sync.ResolveConflict(table, entityID, mergedData, serverVersion)
}
`
}

func desktopClientKeychainGo(opts Options) string {
	return `package main

import (
	"github.com/99designs/keyring"
)

// Keychain is a thin wrapper around go-keyring that uses the best available
// OS-native secret store (macOS Keychain, Windows Credential Manager,
// Linux Secret Service). Falls back to an encrypted file in ~/.config on
// systems without a keyring daemon.
type Keychain struct {
	ring keyring.Keyring
}

func NewKeychain() *Keychain {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: "` + opts.ProjectName + `",
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,           // macOS
			keyring.WinCredBackend,            // Windows
			keyring.SecretServiceBackend,      // Linux (GNOME)
			keyring.KWalletBackend,            // Linux (KDE)
			keyring.FileBackend,               // fallback
		},
		FilePasswordFunc: keyring.FixedStringPrompt(""),
		FileDir:          "~/.config/` + opts.ProjectName + `/keyring",
	})
	if err != nil {
		// Best-effort: fall back to in-memory (not persisted).
		return &Keychain{ring: nil}
	}
	return &Keychain{ring: ring}
}

func (k *Keychain) Set(key, value string) error {
	if k.ring == nil {
		return nil
	}
	return k.ring.Set(keyring.Item{
		Key:  key,
		Data: []byte(value),
	})
}

func (k *Keychain) Get(key string) (string, error) {
	if k.ring == nil {
		return "", nil
	}
	item, err := k.ring.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

func (k *Keychain) Delete(key string) error {
	if k.ring == nil {
		return nil
	}
	return k.ring.Remove(key)
}
`
}

func desktopClientGoMod(module string) string {
	return `module ` + module + `

go 1.24.2

require (
	github.com/99designs/keyring v1.2.2
	github.com/glebarez/sqlite v1.11.0
	github.com/google/uuid v1.6.0
	github.com/wailsapp/wails/v2 v2.9.2
	gorm.io/gorm v1.25.12
)
`
}

func desktopClientWailsJSON(opts Options) string {
	return `{
  "$schema": "https://wails.io/schemas/config.v2.json",
  "name": "` + opts.ProjectName + `",
  "outputfilename": "` + opts.ProjectName + `",
  "frontend:install": "pnpm install",
  "frontend:build": "pnpm build",
  "frontend:dev:watcher": "pnpm dev",
  "frontend:dev:serverUrl": "auto",
  "author": {
    "name": "",
    "email": ""
  },
  "info": {
    "productName": "` + opts.ProjectName + `",
    "productVersion": "0.1.0",
    "copyright": "",
    "comments": "Built with Grit"
  }
}
`
}

func desktopClientGitignore() string {
	return `# Build output
build/bin/
frontend/dist/
frontend/node_modules/
frontend/wailsjs/

# Go
*.exe
*.test

# Editor
.vscode/
.idea/
*.swp

# OS
.DS_Store
Thumbs.db

# Env
.env
.env.local
`
}

func desktopClientReadme(opts Options) string {
	return `# ` + opts.ProjectName + ` — Desktop

Native desktop app built with [Wails v2](https://wails.io). Shares the same
Go API (over HTTP) as the web app, admin panel, and mobile app in this monorepo.

## Prerequisites

- Go 1.24+
- pnpm 8+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation): ` + "`" + `go install github.com/wailsapp/wails/v2/cmd/wails@latest` + "`" + `

## Development

From the project root, make sure the API is running:

` + "```bash" + `
cd apps/api && go run cmd/server/main.go
` + "```" + `

Then start the desktop app:

` + "```bash" + `
cd apps/desktop
wails dev
` + "```" + `

This opens the app window with hot reload for both Go and React.

## Build

` + "```bash" + `
wails build
` + "```" + `

Outputs a native binary in ` + "`" + `build/bin/` + "`" + `.

## Architecture

This desktop app is a **thin client** of the shared API. Unlike ` + "`" + `grit new-desktop` + "`" + `
(standalone offline-first app), this app:

- Has NO embedded Go business logic
- Has NO local SQLite
- Calls the same API that the web/mobile apps call
- Uses OS keychain (via ` + "`" + `99designs/keyring` + "`" + `) for secure JWT storage
- Uses Wails bindings only for native OS stuff (window controls, file dialogs, keychain)

See [GRIT_STYLE_GUIDE.md](../../GRIT_STYLE_GUIDE.md) §14.5 for desktop design patterns.
`
}

// ═══════════════════════════════════════════════════════════════════
// Frontend config
// ═══════════════════════════════════════════════════════════════════

// desktopClientPackageJSON emits the Wails frontend package.json.
//
// The "build" script runs vite BEFORE tsc on purpose: the TanStack Router
// plugin generates src/routeTree.gen.ts during vite's build, and tsc can't
// typecheck the routes until that file exists. With "tsc -b && vite build" a
// fresh clone fails with "Cannot find module './routeTree.gen'" plus a
// createFileRoute error on every route — which breaks `wails build`, since it
// shells out to this script. tsc still gates the build: a type error exits
// non-zero and wails aborts.
func desktopClientPackageJSON(opts Options) string {
	return `{
  "name": "@` + opts.ProjectName + `/desktop",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build && tsc -b",
    "typecheck": "tsc -b",
    "preview": "vite preview",
    "format": "prettier --write ."
  },
  "dependencies": {
    "@tanstack/react-query": "^5.0.0",
    "@tanstack/react-router": "^1.95.0",
    "axios": "^1.7.9",
    "class-variance-authority": "^0.7.1",
    "clsx": "^2.1.1",
    "lucide-react": "^0.469.0",
    "react": "^19.0.0",
    "react-dom": "^19.0.0",
    "react-hook-form": "^7.54.0",
    "@hookform/resolvers": "^3.9.0",
    "tailwind-merge": "^2.5.5",
    "zod": "^3.23.0"
  },
  "devDependencies": {
    "@tanstack/router-plugin": "^1.95.0",
    "@types/node": "^22.10.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "@vitejs/plugin-react": "^4.3.4",
    "autoprefixer": "^10.4.20",
    "postcss": "^8.4.49",
    "prettier": "^3.3.0",
    "prettier-plugin-tailwindcss": "^0.6.0",
    "tailwindcss": "^3.4.0",
    "typescript": "^5.7.0",
    "vite": "^6.0.0"
  }
}
`
}

func desktopClientTSConfig() string {
	return `{
  "compilerOptions": {
    "target": "ES2022",
    "useDefineForClassFields": true,
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "skipLibCheck": true,
    "moduleResolution": "bundler",
    "allowImportingTsExtensions": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx",
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noFallthroughCasesInSwitch": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["./src/*"],
      "@repo/shared": ["../../../packages/shared"],
      "@repo/shared/*": ["../../../packages/shared/*"]
    }
  },
  "include": ["src"],
  "references": [{ "path": "./tsconfig.node.json" }]
}
`
}

func desktopClientTSConfigNode() string {
	return `{
  "compilerOptions": {
    "composite": true,
    "skipLibCheck": true,
    "module": "ESNext",
    "moduleResolution": "bundler",
    "allowSyntheticDefaultImports": true,
    "strict": true
  },
  "include": ["vite.config.ts"]
}
`
}

func desktopClientViteConfig() string {
	return `import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";
import path from "node:path";

export default defineConfig({
  plugins: [
    TanStackRouterVite({ routesDirectory: "./src/routes", generatedRouteTree: "./src/routeTree.gen.ts" }),
    react(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      "@repo/shared": path.resolve(__dirname, "../../../packages/shared"),
    },
  },
  server: {
    port: 5174,
    strictPort: true,
    // In dev, proxy /api to the Go API server. Wails dev mode runs this Vite
    // server and wraps it in a native window.
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
});
`
}

func desktopClientTailwindConfig() string {
	return `import type { Config } from "tailwindcss";

// GRIT_STYLE_GUIDE-aligned config, plus desktop-specific overrides
// (larger base padding, tighter focus rings, command palette tokens).
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  darkMode: "class",
  theme: {
    extend: {
      fontFamily: {
        // Driven by the active theme (atlas: Inter, aurora: Geist,
        // pulse: Onest + DM Serif Display). See lib/theme-tokens.ts.
        sans: ["var(--font-ui)", "system-ui", "sans-serif"],
        display: ["var(--font-display)", "var(--font-ui)", "system-ui", "sans-serif"],
        mono: ["JetBrains Mono", "ui-monospace", "monospace"],
      },
      colors: {
        // Grit purple
        primary: {
          50: "#F3F0FE",
          100: "#E6DFFD",
          200: "#CCC0FB",
          400: "#9B8BF5",
          500: "#7C6CE9",
          600: "#6C5CE7",
          700: "#5B4BD6",
          800: "#4A3DB5",
          900: "#3B2F8F",
        },
        // CSS variable bindings for runtime theming
        background: "var(--bg-primary)",
        surface: "var(--bg-secondary)",
        "surface-2": "var(--bg-tertiary)",
        "surface-3": "var(--bg-elevated)",
        "surface-hover": "var(--bg-hover)",
        border: "var(--border)",
        "border-subtle": "var(--border-subtle)",
        foreground: "var(--text-foreground)",
        "foreground-secondary": "var(--text-secondary)",
        "foreground-muted": "var(--text-muted)",
        accent: "var(--accent)",
        "accent-hover": "var(--accent-hover)",
        success: "#059669",
        warning: "#D97706",
        danger: "#DC2626",
      },
      boxShadow: {
        xs: "0 1px 2px 0 rgba(0, 0, 0, 0.3)",
        "focus": "0 0 0 2px rgba(108, 92, 231, 0.25)",
      },
      borderRadius: {
        // Each theme ships its own radius token (atlas .625rem, aurora .75rem,
        // pulse .5rem).
        DEFAULT: "var(--radius)",
      },
      spacing: {
        // Desktop uses larger padding than web — more breathing room
        "content": "2rem",     // 32px main content padding
        "sidebar": "15rem",    // 240px fixed sidebar
        "titlebar": "3rem",    // 48px custom title bar
        "listpane": "22rem",   // 352px list pane in TwoPane layout
      },
    },
  },
  plugins: [],
} satisfies Config;
`
}

func desktopClientPostCSSConfig() string {
	return `export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
};
`
}

func desktopClientIndexHTML(opts Options) string {
	return `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>` + opts.ProjectName + `</title>
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link href="` + desktopThemeFontURL(opts.Theme) + `" rel="stylesheet" />
  </head>
  <body>
    <div id="app"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
`
}

func desktopClientMainTSX() string {
	return `import React from "react";
import ReactDOM from "react-dom/client";
import { QueryClientProvider } from "@tanstack/react-query";
import { RouterProvider, createRouter, createHashHistory } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";
import { queryClient } from "./lib/query-client";
import { AuthProvider } from "./lib/auth-provider";
import { ThemeProvider } from "./lib/theme-provider";
import "./globals.css";

// Hash history works inside Wails' single-page webview context.
const router = createRouter({
  routeTree,
  history: createHashHistory(),
  context: { auth: undefined! }, // populated by AuthProvider
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

ReactDOM.createRoot(document.getElementById("app")!).render(
  <React.StrictMode>
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <RouterProvider router={router} />
        </AuthProvider>
      </QueryClientProvider>
    </ThemeProvider>
  </React.StrictMode>
);
`
}

func desktopClientGlobalsCSS() string {
	return `@tailwind base;
@tailwind components;
@tailwind utilities;

/* ═══════════════════════════════════════════════════════════════════
   Grit Desktop — theme tokens (GRIT_STYLE_GUIDE §17)
   ═══════════════════════════════════════════════════════════════════ */

:root {
  /* Light mode — rarely used on desktop but supported */
  --bg-primary: #FAFAFA;
  --bg-secondary: #FFFFFF;
  --bg-tertiary: #F4F4F5;
  --bg-elevated: #FFFFFF;
  --bg-hover: #F4F4F5;

  --text-foreground: #18181B;
  --text-secondary: #52525B;
  --text-muted: #71717A;

  --border: #E4E4E7;
  --border-subtle: #F4F4F5;

  --accent: #6C5CE7;
  --accent-hover: #5B4BD6;
}

.dark {
  /* Dark mode — default on desktop */
  --bg-primary: #0a0a0f;
  --bg-secondary: #111118;
  --bg-tertiary: #1a1a24;
  --bg-elevated: #22222e;
  --bg-hover: #2a2a38;

  --text-foreground: #e8e8f0;
  --text-secondary: #9090a8;
  --text-muted: #606078;

  --border: #2a2a3a;
  --border-subtle: #1a1a24;

  --accent: #6c5ce7;
  --accent-hover: #7c6cf7;
}

html, body, #app {
  height: 100%;
  margin: 0;
  font-family: Onest, system-ui, sans-serif;
  background-color: var(--bg-primary);
  color: var(--text-foreground);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  overflow: hidden; /* Desktop window is already fixed; no page scroll */
}

/* Drag region for Wails frameless window */
.drag-region {
  --wails-draggable: drag;
}

.no-drag {
  --wails-draggable: no-drag;
}

/* Custom scrollbar — desktop users see scrollbars more */
::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 999px;
  border: 2px solid var(--bg-primary);
}

::-webkit-scrollbar-thumb:hover {
  background: var(--text-muted);
}

/* Focus ring matches GRIT_STYLE_GUIDE */
*:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
  border-radius: 4px;
}
`
}

func desktopClientViteEnvDTS() string {
	return `/// <reference types="vite/client" />

// Wails bindings are injected at runtime. This is a placeholder type so
// TypeScript compiles before the first build generates wailsjs/.
//
// IMPORTANT: every method bound on the Go App struct (app.go) must be declared
// here, or any file that calls it fails typechecking — which is what
// "pnpm build" (and therefore "wails build") runs. Add a binding in app.go?
// Add it here too.

// Structural mirrors of the Go sync types. Kept structural (not imported) so
// this .d.ts has no module dependencies; lib/sync-client.ts declares the same
// shapes and TypeScript matches them by structure.
type WailsSyncResult = {
  pushed: number;
  pulled: number;
  conflicts: number;
  errors?: string[];
  started_at: string;
  finished_at: string;
};

type WailsOutboxEntry = {
  ID: number;
  Model: string;
  EntityID: string;
  Op: "create" | "update" | "delete";
  Data: string | null;
  Version: number;
  CreatedAt: number;
  HasConflict: boolean;
  ServerData: string | null;
  ServerVersion: number;
  ConflictMsg: string;
};

type WailsSyncStatus = {
  reachable: boolean;
  force_offline: boolean;
  pending: number;
  last_sync?: string;
  last_error?: string;
};

declare global {
  interface Window {
    runtime?: {
      EventsOn: (event: string, callback: (...args: unknown[]) => void) => void;
      EventsEmit: (event: string, ...args: unknown[]) => void;
    };
    go?: {
      main: {
        App: {
          SetToken: (key: string, value: string) => Promise<void>;
          GetToken: (key: string) => Promise<string>;
          DeleteToken: (key: string) => Promise<void>;
          MinimiseWindow: () => Promise<void>;
          MaximiseWindow: () => Promise<void>;
          UnmaximiseWindow: () => Promise<void>;
          ToggleMaximise: () => Promise<void>;
          CloseWindow: () => Promise<void>;
          IsMaximised: () => Promise<boolean>;
          OpenFileDialog: (title: string) => Promise<string>;
          SaveFileDialog: (title: string, defaultFilename: string) => Promise<string>;
          GetPlatform: () => Promise<"darwin" | "windows" | "linux">;
          GetAppVersion: () => Promise<string>;

          // Offline sync engine (local mirror + outbox)
          LocalCreate: (table: string, id: string, data: Record<string, unknown>) => Promise<void>;
          LocalUpdate: (table: string, id: string, data: Record<string, unknown>) => Promise<void>;
          LocalDelete: (table: string, id: string) => Promise<void>;
          LocalGet: (table: string, id: string) => Promise<Record<string, unknown> | null>;
          LocalList: (table: string) => Promise<Record<string, unknown>[]>;
          Sync: (tables: string[]) => Promise<WailsSyncResult>;
          PendingCount: () => Promise<number>;
          GetPendingChanges: () => Promise<WailsOutboxEntry[]>;
          ResolveConflict: (
            table: string,
            entityID: string,
            mergedData: Record<string, unknown>,
            serverVersion: number,
          ) => Promise<void>;

          // Online/offline mode
          GetSyncTables: () => Promise<string[]>;
          SetOfflineMode: (offline: boolean) => Promise<void>;
          GetSyncStatus: () => Promise<WailsSyncStatus>;
          SyncNow: () => Promise<WailsSyncResult>;
        };
      };
    };
  }
}

export {};
`
}

// ═══════════════════════════════════════════════════════════════════
// Routes (TanStack Router file-based)
// ═══════════════════════════════════════════════════════════════════

func desktopClientRootRoute() string {
	return `import { Outlet, createRootRoute } from "@tanstack/react-router";

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return <Outlet />;
}
`
}

func desktopClientIndexRoute() string {
	return `import { createFileRoute, redirect } from "@tanstack/react-router";

// Root redirect: if authenticated go to /app, else /auth/login
export const Route = createFileRoute("/")({
  beforeLoad: async () => {
    // Check for token via Wails bridge (falls back to localStorage in dev)
    const { getToken } = await import("@/lib/wails-bridge");
    const token = await getToken("access_token");
    if (token) {
      throw redirect({ to: "/app" });
    }
    throw redirect({ to: "/auth/login" });
  },
});
`
}

func desktopClientAuthLayoutRoute() string {
	return `import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";
import { TitleBar } from "@/components/layout/title-bar";

export const Route = createFileRoute("/auth")({
  beforeLoad: async () => {
    const { getToken } = await import("@/lib/wails-bridge");
    const token = await getToken("access_token");
    if (token) {
      throw redirect({ to: "/app" });
    }
  },
  component: AuthLayout,
});

function AuthLayout() {
  return (
    <div className="flex flex-col h-screen bg-background">
      <TitleBar showSidebarControls={false} />
      {/* The themed AuthShell owns its own layout (split hero / centered card),
          so this wrapper only provides the scroll container. */}
      <div className="flex-1 overflow-auto">
        <Outlet />
      </div>
    </div>
  );
}
`
}

func desktopClientLoginRoute() string {
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useState } from "react";
import { Loader2, Eye, EyeOff } from "lucide-react";
import { useLogin } from "@/hooks/use-auth";
import { AuthShell } from "@/components/auth/AuthShell";
import { authInputCls, authInputStyle, AuthSubmit } from "@/components/auth/AuthField";

export const Route = createFileRoute("/auth/login")({
  component: LoginPage,
});

const LoginSchema = z.object({
  email: z.string().email("Invalid email"),
  password: z.string().min(6, "Minimum 6 characters"),
});

type LoginInput = z.infer<typeof LoginSchema>;

function LoginPage() {
  const navigate = useNavigate();
  const { mutate: login, isPending, error } = useLogin();
  const [showPassword, setShowPassword] = useState(false);

  const { register, handleSubmit, formState: { errors } } = useForm<LoginInput>({
    resolver: zodResolver(LoginSchema),
  });

  const onSubmit = (data: LoginInput) => {
    login(data, { onSuccess: () => navigate({ to: "/app" }) });
  };

  return (
    <AuthShell
      mode="login"
      title="Welcome back"
      subtitle="Sign in to your account"
      errorMessage={error ? ((error as Error).message || "Invalid email or password") : undefined}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="space-y-1.5">
          <label htmlFor="email" className="block text-sm font-medium">Email</label>
          <input
            id="email"
            type="email"
            autoComplete="email"
            placeholder="you@example.com"
            className={authInputCls}
            style={authInputStyle}
            {...register("email")}
          />
          {errors.email && <p className="text-xs text-[#dc2626]">{errors.email.message}</p>}
        </div>

        <div className="space-y-1.5">
          <label htmlFor="password" className="block text-sm font-medium">Password</label>
          <div className="relative">
            <input
              id="password"
              type={showPassword ? "text" : "password"}
              autoComplete="current-password"
              placeholder="Enter your password"
              className={authInputCls + " pr-10"}
              style={authInputStyle}
              {...register("password")}
            />
            <button
              type="button"
              onClick={() => setShowPassword((v) => !v)}
              className="absolute right-3 top-1/2 -translate-y-1/2 opacity-60 hover:opacity-100"
              aria-label={showPassword ? "Hide password" : "Show password"}
            >
              {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </button>
          </div>
          {errors.password && <p className="text-xs text-[#dc2626]">{errors.password.message}</p>}
        </div>

        <AuthSubmit disabled={isPending}>
          {isPending ? (<><Loader2 className="h-4 w-4 animate-spin" /> Signing in…</>) : "Sign In"}
        </AuthSubmit>
      </form>
    </AuthShell>
  );
}
`
}

func desktopClientRegisterRoute() string {
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";
import { useRegister } from "@/hooks/use-auth";
import { AuthShell } from "@/components/auth/AuthShell";
import { authInputCls, authInputStyle, AuthSubmit } from "@/components/auth/AuthField";

export const Route = createFileRoute("/auth/register")({
  component: RegisterPage,
});

const RegisterSchema = z.object({
  first_name: z.string().min(1, "Required"),
  last_name: z.string().min(1, "Required"),
  email: z.string().email("Invalid email"),
  password: z.string().min(8, "Minimum 8 characters"),
});

type RegisterInput = z.infer<typeof RegisterSchema>;

function RegisterPage() {
  const navigate = useNavigate();
  const { mutate: registerUser, isPending, error } = useRegister();

  const { register, handleSubmit, formState: { errors } } = useForm<RegisterInput>({
    resolver: zodResolver(RegisterSchema),
  });

  const onSubmit = (data: RegisterInput) => {
    registerUser(data, { onSuccess: () => navigate({ to: "/app" }) });
  };

  return (
    <AuthShell
      mode="sign-up"
      title="Create your account"
      subtitle="Get started in less than a minute"
      errorMessage={error ? ((error as Error).message || "Could not create your account") : undefined}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="grid grid-cols-2 gap-3">
          <div className="space-y-1.5">
            <label className="block text-sm font-medium">First name</label>
            <input type="text" placeholder="Jane" className={authInputCls} style={authInputStyle} {...register("first_name")} />
            {errors.first_name && <p className="text-xs text-[#dc2626]">{errors.first_name.message}</p>}
          </div>
          <div className="space-y-1.5">
            <label className="block text-sm font-medium">Last name</label>
            <input type="text" placeholder="Doe" className={authInputCls} style={authInputStyle} {...register("last_name")} />
            {errors.last_name && <p className="text-xs text-[#dc2626]">{errors.last_name.message}</p>}
          </div>
        </div>

        <div className="space-y-1.5">
          <label className="block text-sm font-medium">Email</label>
          <input type="email" autoComplete="email" placeholder="you@example.com" className={authInputCls} style={authInputStyle} {...register("email")} />
          {errors.email && <p className="text-xs text-[#dc2626]">{errors.email.message}</p>}
        </div>

        <div className="space-y-1.5">
          <label className="block text-sm font-medium">Password</label>
          <input type="password" autoComplete="new-password" placeholder="At least 8 characters" className={authInputCls} style={authInputStyle} {...register("password")} />
          {errors.password && <p className="text-xs text-[#dc2626]">{errors.password.message}</p>}
        </div>

        <AuthSubmit disabled={isPending}>
          {isPending ? (<><Loader2 className="h-4 w-4 animate-spin" /> Creating account…</>) : "Create account"}
        </AuthSubmit>
      </form>
    </AuthShell>
  );
}
`
}
func desktopClientAppLayoutRoute() string {
	return `import { Outlet, createFileRoute, redirect } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { TitleBar } from "@/components/layout/title-bar";
import { Sidebar } from "@/components/layout/sidebar";
import { Topbar } from "@/components/layout/topbar";
import { CommandPalette } from "@/components/layout/command-palette";
import { useShortcuts } from "@/lib/use-shortcuts";

export const Route = createFileRoute("/app")({
  beforeLoad: async () => {
    const { getToken } = await import("@/lib/wails-bridge");
    const token = await getToken("access_token");
    if (!token) {
      throw redirect({ to: "/auth/login" });
    }
  },
  component: AppLayout,
});

function AppLayout() {
  const [paletteOpen, setPaletteOpen] = useState(false);

  // Global keyboard shortcuts
  useShortcuts({
    "mod+k": () => setPaletteOpen(true),
  });

  // Close palette on escape
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") setPaletteOpen(false);
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, []);

  return (
    <div className="flex flex-col h-screen bg-background">
      {/* Frameless title bar */}
      <TitleBar />

      <div className="flex flex-1 overflow-hidden">
        {/* Fixed desktop sidebar (NOT collapsible — desktop has room) */}
        <Sidebar />

        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Topbar: search, notifications, theme, user */}
          <Topbar onOpenPalette={() => setPaletteOpen(true)} />

          {/* Main content — 32px padding (more negative space than web) */}
          <main className="flex-1 overflow-auto p-content">
            <Outlet />
          </main>
        </div>
      </div>

      {/* Command palette (Cmd+K) */}
      <CommandPalette open={paletteOpen} onOpenChange={setPaletteOpen} />
    </div>
  );
}
`
}

func desktopClientDashboardRoute() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "@/components/layout/page-header";
import { Card, CardContent } from "@/components/ui/card";
import { Activity, Users, Zap, Clock } from "lucide-react";

export const Route = createFileRoute("/app/")({
  component: DashboardPage,
});

function DashboardPage() {
  return (
    <div>
      <PageHeader
        title="Dashboard"
        description="Welcome back. Here's what's happening."
        stats={[
          { label: "Total Users", value: "—", icon: Users, color: "default" },
          { label: "Active Now", value: "—", icon: Activity, color: "success" },
          { label: "Requests/min", value: "—", icon: Zap, color: "default" },
          { label: "Uptime", value: "—", icon: Clock, color: "default" },
        ]}
      />

      <div className="mt-8 grid gap-4 md:grid-cols-2">
        <Card>
          <CardContent className="p-6">
            <h3 className="text-[15px] font-semibold text-foreground">Recent Activity</h3>
            <p className="mt-1 text-[13px] text-foreground-secondary">
              Connect to your API to show recent events here.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <h3 className="text-[15px] font-semibold text-foreground">Quick Actions</h3>
            <p className="mt-1 text-[13px] text-foreground-secondary">
              Press{" "}
              <kbd className="rounded border border-border bg-surface-2 px-1.5 py-0.5 text-[11px] font-mono">
                ⌘K
              </kbd>{" "}
              to open the command palette.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
`
}

func desktopClientProfileRoute() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "@/components/layout/page-header";
import { Card, CardContent } from "@/components/ui/card";
import { useMe } from "@/hooks/use-auth";

export const Route = createFileRoute("/app/profile")({
  component: ProfilePage,
});

function ProfilePage() {
  const { data: user } = useMe();

  return (
    <div>
      <PageHeader title="Profile" description="Manage your account details" />

      <Card className="mt-6 max-w-2xl">
        <CardContent className="p-6">
          <div className="flex items-center gap-4">
            <div className="h-16 w-16 rounded-full bg-accent/20 flex items-center justify-center">
              <span className="text-xl font-semibold text-accent">
                {user?.first_name?.charAt(0)?.toUpperCase() || "?"}
              </span>
            </div>
            <div>
              <p className="text-[16px] font-semibold text-foreground">
                {user?.first_name} {user?.last_name}
              </p>
              <p className="text-[13px] text-foreground-secondary">{user?.email}</p>
              <p className="mt-1 inline-block text-[11px] font-semibold uppercase tracking-wider text-accent">
                {user?.role}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
`
}

func desktopClientSettingsRoute() string {
	return `import { createFileRoute } from "@tanstack/react-router";
import { PageHeader } from "@/components/layout/page-header";
import { Card, CardContent } from "@/components/ui/card";
import { useTheme } from "@/lib/theme-provider";
import { Moon, Sun } from "lucide-react";
import { OfflineModeToggle } from "@/components/offline-mode-toggle";

export const Route = createFileRoute("/app/settings")({
  component: SettingsPage,
});

function SettingsPage() {
  const { theme, setTheme } = useTheme();

  return (
    <div>
      <PageHeader title="Settings" description="Configure your desktop app" />

      <div className="mt-6 max-w-2xl space-y-4">
        <Card>
          <CardContent className="p-6">
            <div>
              <h3 className="text-[15px] font-semibold text-foreground">Sync &amp; Offline</h3>
              <p className="mt-1 mb-4 text-[13px] text-foreground-secondary">
                By default the app works online and mirrors your data locally.
                Switch to offline to keep working against the local copy — your
                changes queue up and sync automatically when you switch back.
              </p>
              <OfflineModeToggle />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-[15px] font-semibold text-foreground">Appearance</h3>
                <p className="mt-1 text-[13px] text-foreground-secondary">
                  Choose your preferred theme
                </p>
              </div>
              <button
                onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
                className="flex items-center gap-2 rounded-lg border border-border bg-surface-2 px-3 h-9 text-[13px] font-medium text-foreground hover:bg-surface-hover transition-colors"
              >
                {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
                {theme === "dark" ? "Light" : "Dark"}
              </button>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <h3 className="text-[15px] font-semibold text-foreground">Keyboard Shortcuts</h3>
            <p className="mt-1 text-[13px] text-foreground-secondary">
              Global shortcuts available everywhere in the app.
            </p>
            <div className="mt-4 space-y-2 text-[13px]">
              <div className="flex items-center justify-between">
                <span className="text-foreground-secondary">Command palette</span>
                <kbd className="rounded border border-border bg-surface-2 px-1.5 py-0.5 text-[11px] font-mono">⌘K</kbd>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-foreground-secondary">Settings</span>
                <kbd className="rounded border border-border bg-surface-2 px-1.5 py-0.5 text-[11px] font-mono">⌘,</kbd>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-foreground-secondary">Logout</span>
                <kbd className="rounded border border-border bg-surface-2 px-1.5 py-0.5 text-[11px] font-mono">⌘L</kbd>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
`
}

// ═══════════════════════════════════════════════════════════════════
// Layout components
// ═══════════════════════════════════════════════════════════════════

func desktopClientTitleBar(opts Options) string {
	return `import { Minus, Square, X } from "lucide-react";
import { useEffect, useState } from "react";
import {
  minimise,
  toggleMaximise,
  closeWindow,
  getPlatform,
} from "@/lib/wails-bridge";
import { useOnlineStatus } from "@/hooks/use-online-status";
import { SyncButton } from "@/components/sync-button";
import { getSyncTables } from "@/lib/wails-bridge";

interface TitleBarProps {
  showSidebarControls?: boolean;
}

// Frameless title bar with drag region and OS-specific window controls.
// macOS: traffic lights on the left. Windows/Linux: min/max/close on the right.
export function TitleBar({ showSidebarControls: _ = true }: TitleBarProps) {
  const [platform, setPlatform] = useState<"darwin" | "windows" | "linux">("windows");
  // Sync tables are owned by the Go side (grit generate resource appends to
  // them), so read them at runtime rather than hardcoding a list here.
  const [syncTables, setSyncTables] = useState<string[]>(["users", "uploads"]);

  useEffect(() => {
    getPlatform().then(setPlatform);
    getSyncTables().then(setSyncTables).catch(() => {});
  }, []);

  const isMac = platform === "darwin";

  return (
    <div
      className="drag-region flex h-titlebar shrink-0 items-center justify-between border-b border-border-subtle bg-surface px-3 select-none"
    >
      {/* Left: macOS traffic lights OR spacer on Win/Linux */}
      <div className="flex items-center gap-2">
        {isMac ? (
          <div className="flex items-center gap-2">
            <MacTrafficLight color="close" onClick={closeWindow} />
            <MacTrafficLight color="minimise" onClick={minimise} />
            <MacTrafficLight color="maximise" onClick={toggleMaximise} />
          </div>
        ) : (
          <div className="flex items-center gap-2 pl-2">
            <div className="h-6 w-6 rounded-md bg-accent flex items-center justify-center">
              <span className="text-[11px] font-bold text-white">G</span>
            </div>
            <span className="text-[13px] font-medium text-foreground">` + opts.ProjectName + `</span>
          </div>
        )}
      </div>

      {/* Center: app name (macOS only, since controls are on left) */}
      {isMac && (
        <div className="text-[13px] font-medium text-foreground-secondary">` + opts.ProjectName + `</div>
      )}

      {/* Right: sync button + connection indicator + (on Windows/Linux) window controls */}
      <div className="no-drag flex items-center">
        <SyncButton tables={syncTables} />
        <ConnectionIndicator />
        {!isMac && (
          <>
            <WinControl icon={<Minus className="h-3.5 w-3.5" />} onClick={minimise} />
            <WinControl icon={<Square className="h-3 w-3" />} onClick={toggleMaximise} />
            <WinControl icon={<X className="h-4 w-4" />} onClick={closeWindow} variant="close" />
          </>
        )}
      </div>
    </div>
  );
}

// ConnectionIndicator renders a small colored dot reflecting whether the
// API is reachable. Green = healthy, amber = no network or API unreachable.
// Hover for last-checked timestamp.
function ConnectionIndicator() {
  const { isOnline, lastCheckedAt } = useOnlineStatus();
  const label = isOnline ? "Connected" : "Reconnecting...";
  const checked = lastCheckedAt ? lastCheckedAt.toLocaleTimeString() : "...";
  return (
    <div
      className="flex h-titlebar items-center px-3"
      title={` + "`" + `${label} (last check: ${checked})` + "`" + `}
      aria-label={label}
    >
      <span
        className={` + "`" + `h-2 w-2 rounded-full transition-colors ${
          isOnline ? "bg-success animate-none" : "bg-warning animate-pulse"
        }` + "`" + `}
      />
    </div>
  );
}

function MacTrafficLight({
  color,
  onClick,
}: {
  color: "close" | "minimise" | "maximise";
  onClick: () => void;
}) {
  const bg = {
    close: "bg-[#ff5f57] hover:bg-[#ff5f57]/80",
    minimise: "bg-[#febc2e] hover:bg-[#febc2e]/80",
    maximise: "bg-[#28c840] hover:bg-[#28c840]/80",
  }[color];

  return (
    <button
      onClick={onClick}
      className={` + "`" + `no-drag h-3 w-3 rounded-full transition-colors ${bg}` + "`" + `}
      aria-label={color}
    />
  );
}

function WinControl({
  icon,
  onClick,
  variant = "default",
}: {
  icon: React.ReactNode;
  onClick: () => void;
  variant?: "default" | "close";
}) {
  const hoverClass =
    variant === "close"
      ? "hover:bg-danger hover:text-white"
      : "hover:bg-surface-hover";

  return (
    <button
      onClick={onClick}
      className={` + "`" + `flex h-titlebar w-12 items-center justify-center text-foreground-secondary transition-colors ${hoverClass}` + "`" + `}
    >
      {icon}
    </button>
  );
}
`
}

func desktopClientSidebar() string {
	return `import { Link, useRouterState } from "@tanstack/react-router";
import { Home, User, Settings, LayoutGrid } from "lucide-react";

// Fixed-width desktop sidebar. NOT collapsible — desktop windows are wide
// enough, and a fixed sidebar is more predictable for power users.
export function Sidebar() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });

  const items = [
    { to: "/app", label: "Dashboard", icon: Home },
    { to: "/app/profile", label: "Profile", icon: User },
    { to: "/app/settings", label: "Settings", icon: Settings },
  ];

  return (
    <aside className="w-sidebar shrink-0 border-r border-border-subtle bg-surface flex flex-col">
      {/* Workspace header */}
      <div className="px-4 py-3 border-b border-border-subtle">
        <div className="flex items-center gap-2">
          <div className="h-7 w-7 rounded-md bg-accent/10 flex items-center justify-center">
            <LayoutGrid className="h-3.5 w-3.5 text-accent" />
          </div>
          <span className="text-[13px] font-semibold text-foreground">Workspace</span>
        </div>
      </div>

      {/* Nav */}
      <nav className="flex-1 p-3 space-y-0.5">
        {items.map((item) => {
          const Icon = item.icon;
          const active = pathname === item.to || (item.to !== "/app" && pathname.startsWith(item.to));
          return (
            <Link
              key={item.to}
              to={item.to}
              className={` + "`" + `flex items-center gap-3 rounded-lg px-3 py-2 text-[13px] font-medium transition-colors ${
                active
                  ? "bg-accent/10 text-accent"
                  : "text-foreground-secondary hover:bg-surface-hover hover:text-foreground"
              }` + "`" + `}
            >
              <Icon className="h-4 w-4 shrink-0" />
              {item.label}
            </Link>
          );
        })}
      </nav>

      {/* Footer */}
      <div className="px-4 py-3 border-t border-border-subtle text-[11px] text-foreground-muted">
        Built with Grit
      </div>
    </aside>
  );
}
`
}

func desktopClientTopbar() string {
	return `import { useState } from "react";
import { Bell, Search, Sun, Moon, LogOut, Command } from "lucide-react";
import { useTheme } from "@/lib/theme-provider";
import { useMe, useLogout } from "@/hooks/use-auth";
import { useNavigate } from "@tanstack/react-router";

interface TopbarProps {
  onOpenPalette: () => void;
}

// Topbar for desktop app. Same structure as admin v3.8.0 but without the
// sidebar collapse toggle (desktop sidebar is fixed).
export function Topbar({ onOpenPalette }: TopbarProps) {
  const { theme, setTheme } = useTheme();
  const { data: user } = useMe();
  const { mutate: logout } = useLogout();
  const navigate = useNavigate();
  const [userMenuOpen, setUserMenuOpen] = useState(false);

  const handleLogout = () => {
    setUserMenuOpen(false);
    logout(undefined, { onSuccess: () => navigate({ to: "/auth/login" }) });
  };

  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-border-subtle bg-background px-6">
      {/* Search trigger (opens command palette) */}
      <button
        onClick={onOpenPalette}
        className="flex items-center gap-2 rounded-lg border border-border bg-surface-2 px-3 h-9 w-80 text-[13px] text-foreground-muted hover:bg-surface-hover transition-colors"
      >
        <Search className="h-3.5 w-3.5" />
        <span>Search or jump to...</span>
        <div className="ml-auto flex items-center gap-1 text-[11px] font-mono">
          <kbd className="rounded bg-surface px-1 py-0.5">⌘K</kbd>
        </div>
      </button>

      {/* Right cluster */}
      <div className="flex items-center gap-1">
        <button className="flex items-center justify-center h-9 w-9 rounded-lg hover:bg-surface-hover text-foreground-secondary transition-colors" aria-label="Notifications">
          <Bell className="h-4 w-4" />
        </button>

        <button
          onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
          className="flex items-center justify-center h-9 w-9 rounded-lg hover:bg-surface-hover text-foreground-secondary transition-colors"
          aria-label="Toggle theme"
        >
          {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </button>

        <div className="relative">
          <button
            onClick={() => setUserMenuOpen(!userMenuOpen)}
            className="flex items-center gap-2 rounded-lg pl-1 pr-2 py-1 hover:bg-surface-hover transition-colors"
          >
            <div className="h-7 w-7 rounded-full bg-accent/20 flex items-center justify-center">
              <span className="text-[13px] font-medium text-accent">
                {user?.first_name?.charAt(0)?.toUpperCase() || "?"}
              </span>
            </div>
          </button>

          {userMenuOpen && (
            <>
              <div className="fixed inset-0 z-40" onClick={() => setUserMenuOpen(false)} />
              <div className="absolute right-0 top-full mt-2 w-64 rounded-xl border border-border bg-surface-3 shadow-xl z-50 overflow-hidden">
                <div className="px-4 py-3 border-b border-border-subtle">
                  <p className="text-[13px] font-semibold text-foreground">
                    {user?.first_name} {user?.last_name}
                  </p>
                  <p className="text-[12px] text-foreground-muted truncate">{user?.email}</p>
                </div>
                <div className="p-1">
                  <button
                    onClick={() => { setUserMenuOpen(false); navigate({ to: "/app/profile" }); }}
                    className="flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-[13px] text-foreground-secondary hover:bg-surface-hover hover:text-foreground transition-colors"
                  >
                    <Command className="h-4 w-4" />
                    Profile
                  </button>
                </div>
                <div className="p-1 border-t border-border-subtle">
                  <button
                    onClick={handleLogout}
                    className="flex w-full items-center gap-2.5 rounded-md px-3 py-2 text-[13px] text-danger hover:bg-danger/10 transition-colors"
                  >
                    <LogOut className="h-4 w-4" />
                    Log out
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
`
}

func desktopClientPageHeader() string {
	return `import type { LucideIcon } from "lucide-react";
import { TrendingUp, TrendingDown } from "lucide-react";

export interface StatCard {
  label: string;
  value: string | number;
  icon?: LucideIcon;
  color?: "default" | "success" | "warning" | "danger";
  trend?: { value: number; direction: "up" | "down" };
}

interface PageHeaderProps {
  title: string;
  description?: string;
  actions?: React.ReactNode;
  stats?: StatCard[];
}

export function PageHeader({ title, description, actions, stats }: PageHeaderProps) {
  return (
    <div className="mb-8">
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
        <div className="min-w-0 flex-1">
          <h1 className="text-2xl font-bold text-foreground tracking-tight">{title}</h1>
          {description && (
            <p className="mt-1 text-[14px] text-foreground-secondary">{description}</p>
          )}
        </div>
        {actions && <div className="flex items-center gap-2 shrink-0">{actions}</div>}
      </div>

      {stats && stats.length > 0 && (
        <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {stats.map((stat, i) => (
            <StatCardItem key={i} stat={stat} />
          ))}
        </div>
      )}
    </div>
  );
}

const colorClasses: Record<string, { bg: string; text: string }> = {
  default: { bg: "bg-accent/10", text: "text-accent" },
  success: { bg: "bg-success/10", text: "text-success" },
  warning: { bg: "bg-warning/10", text: "text-warning" },
  danger: { bg: "bg-danger/10", text: "text-danger" },
};

function StatCardItem({ stat }: { stat: StatCard }) {
  const color = colorClasses[stat.color || "default"];
  const Icon = stat.icon;

  return (
    <div className="rounded-xl border border-border bg-surface p-5 transition-colors hover:border-foreground-muted/30">
      <div className="flex items-start justify-between">
        <div className="min-w-0 flex-1">
          <p className="text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">
            {stat.label}
          </p>
          <div className="mt-2 flex items-baseline gap-2">
            <p className="text-2xl font-bold text-foreground tabular-nums">
              {typeof stat.value === "number" ? stat.value.toLocaleString() : stat.value}
            </p>
            {stat.trend && (
              <span
                className={` + "`" + `flex items-center gap-0.5 text-xs font-medium ${
                  stat.trend.direction === "up" ? "text-success" : "text-danger"
                }` + "`" + `}
              >
                {stat.trend.direction === "up" ? (
                  <TrendingUp className="h-3 w-3" />
                ) : (
                  <TrendingDown className="h-3 w-3" />
                )}
                {stat.trend.value}%
              </span>
            )}
          </div>
        </div>
        {Icon && (
          <div className={` + "`" + `flex h-9 w-9 items-center justify-center rounded-lg ${color.bg}` + "`" + `}>
            <Icon className={` + "`" + `h-4 w-4 ${color.text}` + "`" + `} />
          </div>
        )}
      </div>
    </div>
  );
}
`
}

func desktopClientCommandPalette() string {
	return `import { useNavigate } from "@tanstack/react-router";
import { useState, useEffect, useRef } from "react";
import { Search, Home, User, Settings, LogOut, ArrowRight } from "lucide-react";
import { useLogout } from "@/hooks/use-auth";

interface CommandPaletteProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CommandPalette({ open, onOpenChange }: CommandPaletteProps) {
  const navigate = useNavigate();
  const { mutate: logout } = useLogout();
  const [query, setQuery] = useState("");
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);

  const go = (to: string) => {
    navigate({ to });
    onOpenChange(false);
    setQuery("");
  };

  const commands = [
    { id: "dashboard", label: "Dashboard", shortcut: "G D", icon: Home, action: () => go("/app") },
    { id: "profile", label: "Profile", shortcut: "G P", icon: User, action: () => go("/app/profile") },
    { id: "settings", label: "Settings", shortcut: "⌘,", icon: Settings, action: () => go("/app/settings") },
    { id: "logout", label: "Log out", shortcut: "⌘L", icon: LogOut, action: () => { logout(undefined, { onSuccess: () => go("/auth/login") }); } },
  ];

  const filtered = commands.filter((c) =>
    c.label.toLowerCase().includes(query.toLowerCase())
  );

  useEffect(() => {
    if (open) {
      setTimeout(() => inputRef.current?.focus(), 0);
      setQuery("");
      setSelectedIndex(0);
    }
  }, [open]);

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (!open) return;
      if (e.key === "ArrowDown") {
        e.preventDefault();
        setSelectedIndex((i) => Math.min(i + 1, filtered.length - 1));
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setSelectedIndex((i) => Math.max(i - 1, 0));
      } else if (e.key === "Enter") {
        e.preventDefault();
        filtered[selectedIndex]?.action();
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [open, filtered, selectedIndex]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-[100] flex items-start justify-center pt-[15vh] bg-black/40 backdrop-blur-sm">
      <div
        onClick={(e) => e.stopPropagation()}
        className="w-full max-w-[560px] mx-4 rounded-xl border border-border bg-surface-3 shadow-2xl overflow-hidden"
      >
        {/* Input */}
        <div className="flex items-center gap-3 border-b border-border-subtle px-4 h-12">
          <Search className="h-4 w-4 text-foreground-muted shrink-0" />
          <input
            ref={inputRef}
            value={query}
            onChange={(e) => { setQuery(e.target.value); setSelectedIndex(0); }}
            placeholder="Type a command or search..."
            className="flex-1 bg-transparent text-[14px] text-foreground placeholder:text-foreground-muted focus:outline-none"
          />
          <kbd className="rounded border border-border bg-surface-2 px-1.5 py-0.5 text-[10px] font-mono text-foreground-muted">ESC</kbd>
        </div>

        {/* Results */}
        <div className="max-h-[400px] overflow-y-auto p-2">
          {filtered.length === 0 ? (
            <div className="px-4 py-8 text-center text-[13px] text-foreground-muted">
              No results for "{query}"
            </div>
          ) : (
            filtered.map((cmd, i) => {
              const Icon = cmd.icon;
              const isSelected = i === selectedIndex;
              return (
                <button
                  key={cmd.id}
                  onClick={cmd.action}
                  onMouseEnter={() => setSelectedIndex(i)}
                  className={` + "`" + `flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-left transition-colors ${
                    isSelected ? "bg-accent/10" : "hover:bg-surface-hover"
                  }` + "`" + `}
                >
                  <Icon className={` + "`" + `h-4 w-4 shrink-0 ${isSelected ? "text-accent" : "text-foreground-secondary"}` + "`" + `} />
                  <span className={` + "`" + `flex-1 text-[14px] ${isSelected ? "text-foreground" : "text-foreground-secondary"}` + "`" + `}>
                    {cmd.label}
                  </span>
                  <div className="flex items-center gap-1">
                    {cmd.shortcut && (
                      <kbd className="rounded border border-border bg-surface-2 px-1.5 py-0.5 text-[11px] font-mono text-foreground-muted">
                        {cmd.shortcut}
                      </kbd>
                    )}
                    {isSelected && <ArrowRight className="h-3.5 w-3.5 text-accent" />}
                  </div>
                </button>
              );
            })
          )}
        </div>
      </div>
      <div
        className="fixed inset-0 -z-10"
        onClick={() => onOpenChange(false)}
      />
    </div>
  );
}
`
}

// ═══════════════════════════════════════════════════════════════════
// UI components (minimal — GRIT_STYLE_GUIDE tokens)
// ═══════════════════════════════════════════════════════════════════

func desktopClientButton() string {
	return `import { forwardRef } from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-1.5 rounded-lg text-[13px] font-medium transition-colors focus-visible:ring-2 focus-visible:ring-accent/20 focus-visible:outline-none disabled:opacity-50 disabled:cursor-not-allowed",
  {
    variants: {
      variant: {
        primary: "bg-accent text-white hover:bg-accent-hover",
        secondary: "border border-border bg-surface-2 text-foreground hover:bg-surface-hover",
        ghost: "text-foreground-secondary hover:bg-surface-hover hover:text-foreground",
        danger: "bg-danger text-white hover:bg-danger/90",
      },
      size: {
        sm: "h-8 px-2.5",
        md: "h-9 px-3",
        lg: "h-11 px-4 text-[14px]",
      },
    },
    defaultVariants: {
      variant: "primary",
      size: "md",
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(buttonVariants({ variant, size }), className)}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";
`
}

func desktopClientInput() string {
	return `import { forwardRef } from "react";
import { cn } from "@/lib/utils";

export const Input = forwardRef<
  HTMLInputElement,
  React.InputHTMLAttributes<HTMLInputElement>
>(({ className, ...props }, ref) => {
  return (
    <input
      ref={ref}
      className={cn(
        "w-full h-10 rounded-lg border border-border bg-surface-2 px-3 text-[14px] text-foreground placeholder:text-foreground-muted focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/20 transition-colors",
        className
      )}
      {...props}
    />
  );
});
Input.displayName = "Input";
`
}

func desktopClientCard() string {
	return `import { forwardRef } from "react";
import { cn } from "@/lib/utils";

export const Card = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cn("rounded-xl border border-border bg-surface shadow-xs", className)}
      {...props}
    />
  )
);
Card.displayName = "Card";

export const CardHeader = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cn("flex flex-col space-y-1.5 p-6 pb-3", className)}
      {...props}
    />
  )
);
CardHeader.displayName = "CardHeader";

export const CardTitle = forwardRef<HTMLHeadingElement, React.HTMLAttributes<HTMLHeadingElement>>(
  ({ className, ...props }, ref) => (
    <h3
      ref={ref}
      className={cn("text-[15px] font-semibold leading-none tracking-tight text-foreground", className)}
      {...props}
    />
  )
);
CardTitle.displayName = "CardTitle";

export const CardDescription = forwardRef<HTMLParagraphElement, React.HTMLAttributes<HTMLParagraphElement>>(
  ({ className, ...props }, ref) => (
    <p
      ref={ref}
      className={cn("text-[13px] text-foreground-secondary", className)}
      {...props}
    />
  )
);
CardDescription.displayName = "CardDescription";

export const CardContent = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div ref={ref} className={cn("p-6 pt-3", className)} {...props} />
  )
);
CardContent.displayName = "CardContent";

export const CardFooter = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
  ({ className, ...props }, ref) => (
    <div
      ref={ref}
      className={cn("flex items-center p-6 pt-0", className)}
      {...props}
    />
  )
);
CardFooter.displayName = "CardFooter";
`
}

func desktopClientDialog() string {
	return `import { useEffect } from "react";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";

interface DialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  children: React.ReactNode;
  className?: string;
}

export function Dialog({ open, onOpenChange, children, className }: DialogProps) {
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") onOpenChange(false);
    };
    if (open) window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [open, onOpenChange]);

  if (!open) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4"
      onClick={() => onOpenChange(false)}
    >
      <div
        onClick={(e) => e.stopPropagation()}
        className={cn(
          "relative w-full max-w-lg rounded-2xl border border-border bg-surface shadow-xl",
          className
        )}
      >
        <button
          onClick={() => onOpenChange(false)}
          className="absolute right-4 top-4 rounded-md p-1 text-foreground-muted hover:bg-surface-hover hover:text-foreground transition-colors"
          aria-label="Close"
        >
          <X className="h-4 w-4" />
        </button>
        {children}
      </div>
    </div>
  );
}

export function DialogHeader({ children }: { children: React.ReactNode }) {
  return <div className="px-6 pt-6 pb-3">{children}</div>;
}

export function DialogTitle({ children }: { children: React.ReactNode }) {
  return <h2 className="text-[17px] font-semibold text-foreground">{children}</h2>;
}

export function DialogDescription({ children }: { children: React.ReactNode }) {
  return <p className="mt-1 text-[13px] text-foreground-secondary">{children}</p>;
}

export function DialogBody({ children }: { children: React.ReactNode }) {
  return <div className="px-6 py-3">{children}</div>;
}

export function DialogFooter({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex justify-end gap-2 px-6 pb-6 pt-3">{children}</div>
  );
}
`
}

func desktopClientAvatar() string {
	return `import { cn } from "@/lib/utils";

interface AvatarProps {
  src?: string;
  alt?: string;
  fallback?: string;
  size?: "sm" | "md" | "lg";
  className?: string;
}

const sizeClasses = {
  sm: "h-7 w-7 text-[11px]",
  md: "h-9 w-9 text-[13px]",
  lg: "h-12 w-12 text-[15px]",
};

export function Avatar({ src, alt, fallback, size = "md", className }: AvatarProps) {
  if (src) {
    return (
      <img
        src={src}
        alt={alt}
        className={cn("rounded-full object-cover", sizeClasses[size], className)}
      />
    );
  }

  return (
    <div
      className={cn(
        "rounded-full bg-accent/20 flex items-center justify-center font-medium text-accent",
        sizeClasses[size],
        className
      )}
    >
      {fallback?.charAt(0)?.toUpperCase() || "?"}
    </div>
  );
}
`
}

func desktopClientBadge() string {
	return `import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[11px] font-medium",
  {
    variants: {
      variant: {
        default: "bg-accent/10 text-accent",
        success: "bg-success/10 text-success",
        warning: "bg-warning/10 text-warning",
        danger: "bg-danger/10 text-danger",
        muted: "bg-surface-2 text-foreground-secondary",
      },
    },
    defaultVariants: { variant: "default" },
  }
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return <span className={cn(badgeVariants({ variant }), className)} {...props} />;
}
`
}

// ═══════════════════════════════════════════════════════════════════
// Lib (API client, auth, shortcuts, Wails bridge)
// ═══════════════════════════════════════════════════════════════════

func desktopClientWailsBridge() string {
	return `// Wails bridge — wraps Wails runtime calls with safe fallbacks for dev mode.
//
// In production (Wails-built binary), window.go.main.App is injected.
// In dev mode (Vite running outside Wails), fall back to localStorage for
// token storage and no-op for window controls so the app still works when
// you run "pnpm dev" directly from the frontend folder.

const isWails = typeof window !== "undefined" && !!window.go?.main?.App;

// ─── Token storage (OS keychain or localStorage fallback) ─────────

export async function setToken(key: string, value: string): Promise<void> {
  if (isWails) {
    return window.go!.main.App.SetToken(key, value);
  }
  localStorage.setItem(key, value);
}

export async function getToken(key: string): Promise<string> {
  if (isWails) {
    return window.go!.main.App.GetToken(key);
  }
  return localStorage.getItem(key) || "";
}

export async function deleteToken(key: string): Promise<void> {
  if (isWails) {
    return window.go!.main.App.DeleteToken(key);
  }
  localStorage.removeItem(key);
}

// ─── Window controls ──────────────────────────────────────────────

export async function minimise(): Promise<void> {
  if (isWails) return window.go!.main.App.MinimiseWindow();
}

export async function toggleMaximise(): Promise<void> {
  if (isWails) return window.go!.main.App.ToggleMaximise();
}

export async function closeWindow(): Promise<void> {
  if (isWails) return window.go!.main.App.CloseWindow();
}

// ─── System info ──────────────────────────────────────────────────

export async function getPlatform(): Promise<"darwin" | "windows" | "linux"> {
  if (isWails) return window.go!.main.App.GetPlatform();
  // In dev (browser), detect from user agent.
  const ua = navigator.userAgent.toLowerCase();
  if (ua.includes("mac")) return "darwin";
  if (ua.includes("linux")) return "linux";
  return "windows";
}

export async function getAppVersion(): Promise<string> {
  if (isWails) return window.go!.main.App.GetAppVersion();
  return "0.1.0-dev";
}

// ─── File dialogs ─────────────────────────────────────────────────

export async function openFileDialog(title: string): Promise<string> {
  if (isWails) return window.go!.main.App.OpenFileDialog(title);
  throw new Error("openFileDialog requires Wails runtime");
}

export async function saveFileDialog(title: string, defaultFilename: string): Promise<string> {
  if (isWails) return window.go!.main.App.SaveFileDialog(title, defaultFilename);
  throw new Error("saveFileDialog requires Wails runtime");
}

// ─── Offline sync mode ────────────────────────────────────────────

export interface SyncStatus {
  reachable: boolean;
  force_offline: boolean;
  pending: number;
  last_sync?: string;
  last_error?: string;
}

// getSyncTables returns the models covered by offline sync (owned by the Go
// side so a newly generated resource is included automatically).
export async function getSyncTables(): Promise<string[]> {
  if (isWails) return window.go!.main.App.GetSyncTables();
  return ["users", "uploads"];
}

// setOfflineMode flips the manual "Work offline" switch. Turning it off
// triggers an immediate background reconcile on the Go side.
export async function setOfflineMode(offline: boolean): Promise<void> {
  if (isWails) return window.go!.main.App.SetOfflineMode(offline);
}

// getSyncStatus returns the reachable/offline/pending snapshot.
export async function getSyncStatus(): Promise<SyncStatus> {
  if (isWails) return window.go!.main.App.GetSyncStatus();
  return { reachable: true, force_offline: false, pending: 0 };
}

// syncNow forces an immediate Pull+Push.
export async function syncNow(): Promise<unknown> {
  if (isWails) return window.go!.main.App.SyncNow();
  return null;
}
`
}

// desktopClientUseSyncStatus polls the Go engine's sync status so the dashboard
// toggle and any status pill stay live.
func desktopClientUseSyncStatus() string {
	return `import { useCallback, useEffect, useState } from "react";
import { getSyncStatus, setOfflineMode as setOfflineModeBridge, syncNow, type SyncStatus } from "@/lib/wails-bridge";

// useSyncStatus polls the engine every 4s (and on demand) for the
// reachable/offline/pending snapshot, and exposes actions to toggle offline
// mode and force a sync.
export function useSyncStatus() {
  const [status, setStatus] = useState<SyncStatus>({
    reachable: true,
    force_offline: false,
    pending: 0,
  });
  const [busy, setBusy] = useState(false);

  const refresh = useCallback(async () => {
    try {
      setStatus(await getSyncStatus());
    } catch {
      /* engine not ready yet — keep last snapshot */
    }
  }, []);

  useEffect(() => {
    refresh();
    const id = setInterval(refresh, 4000);
    return () => clearInterval(id);
  }, [refresh]);

  const setOffline = useCallback(
    async (offline: boolean) => {
      setBusy(true);
      try {
        await setOfflineModeBridge(offline);
        await refresh();
      } finally {
        setBusy(false);
      }
    },
    [refresh],
  );

  const forceSync = useCallback(async () => {
    setBusy(true);
    try {
      await syncNow();
      await refresh();
    } finally {
      setBusy(false);
    }
  }, [refresh]);

  return { status, busy, setOffline, forceSync, refresh };
}
`
}

// desktopClientOfflineToggle is the dashboard "Work offline" switch + status
// pill. Toggling it OFF (back online) triggers an immediate reconcile.
func desktopClientOfflineToggle() string {
	return `import { Cloud, CloudOff, RefreshCw } from "lucide-react";
import { useSyncStatus } from "@/hooks/use-sync-status";

// OfflineModeToggle is the dashboard control for switching between online
// (mirror + sync) and offline (work against the local copy) mode. When online
// it shows reachability + pending count; when offline it shows how many edits
// are queued to push on reconnect.
export function OfflineModeToggle() {
  const { status, busy, setOffline, forceSync } = useSyncStatus();
  const offline = status.force_offline;

  return (
    <div className="flex items-center gap-3 rounded-lg border border-border bg-surface px-3 py-2">
      {offline ? (
        <CloudOff className="h-4 w-4 text-warning" />
      ) : (
        <Cloud className={"h-4 w-4 " + (status.reachable ? "text-success" : "text-danger")} />
      )}

      <div className="min-w-0 flex-1">
        <div className="text-[13px] font-medium text-foreground">
          {offline ? "Working offline" : status.reachable ? "Online" : "Server unreachable"}
        </div>
        <div className="text-[11px] text-foreground-muted">
          {status.pending > 0
            ? status.pending + " change" + (status.pending === 1 ? "" : "s") + " waiting to sync"
            : status.last_sync
              ? "Last synced " + new Date(status.last_sync).toLocaleTimeString()
              : "All changes synced"}
        </div>
      </div>

      {!offline && status.reachable && (
        <button
          type="button"
          onClick={forceSync}
          disabled={busy}
          title="Sync now"
          className="rounded p-1.5 text-foreground-secondary hover:bg-surface-hover disabled:opacity-50"
        >
          <RefreshCw className={"h-4 w-4 " + (busy ? "animate-spin" : "")} />
        </button>
      )}

      <button
        type="button"
        onClick={() => setOffline(!offline)}
        disabled={busy}
        role="switch"
        aria-checked={offline}
        title={offline ? "Switch back online" : "Switch to offline mode"}
        className={
          "relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors " +
          (offline ? "bg-warning" : "bg-accent")
        }
      >
        <span
          className={
            "inline-block h-4 w-4 rounded-full bg-white transition-transform " +
            (offline ? "translate-x-6" : "translate-x-1")
          }
        />
      </button>
    </div>
  );
}
`
}

func desktopClientApiClientTS() string {
	return `import axios, { type AxiosRequestConfig, type InternalAxiosRequestConfig } from "axios";
import { getToken, setToken, deleteToken } from "./wails-bridge";

// In Wails dev, "/api" is proxied to http://localhost:8080 via Vite.
// In Wails production, the frontend is served from file:// — we need the
// full API URL. Configure via VITE_API_URL or default to localhost:8080.
const API_URL =
  import.meta.env.VITE_API_URL ||
  (typeof window !== "undefined" && !!window.go?.main?.App
    ? "http://localhost:8080/api"
    : "/api");

export const apiClient = axios.create({
  baseURL: API_URL,
  headers: { "Content-Type": "application/json" },
});

// Attach JWT from OS keychain (or localStorage in dev) and auto-generate an
// Idempotency-Key for unsafe methods so the server can dedupe retries (e.g.
// the 401 refresh path below re-issues the same request — without a stable
// key, a network blip mid-write could double-charge or double-create).
apiClient.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
  const token = await getToken("access_token");
  if (token && config.headers) {
    config.headers.Authorization = ` + "`" + `Bearer ${token}` + "`" + `;
  }
  if (config.headers) {
    const method = (config.method || "get").toUpperCase();
    const unsafe = method === "POST" || method === "PUT" || method === "PATCH" || method === "DELETE";
    if (unsafe && !config.headers["Idempotency-Key"]) {
      config.headers["Idempotency-Key"] = crypto.randomUUID();
    }
  }
  return config;
});

// Refresh on 401
let isRefreshing = false;
let refreshQueue: Array<(token: string) => void> = [];

apiClient.interceptors.response.use(
  (res) => res,
  async (error) => {
    const original = error.config as AxiosRequestConfig & { _retry?: boolean };

    // Don't try to refresh on the auth endpoints themselves — a wrong
    // password is a real 401 that should bubble up cleanly. Refreshing
    // here would 401 again, wipe tokens, and leave the user stuck in
    // a login loop.
    const url = original.url || "";
    const isAuthEndpoint =
      url.includes("/auth/login") ||
      url.includes("/auth/register") ||
      url.includes("/auth/refresh");

    if (error.response?.status === 401 && !original._retry && !isAuthEndpoint) {
      original._retry = true;

      if (isRefreshing) {
        return new Promise((resolve) => {
          refreshQueue.push((token: string) => {
            if (original.headers) original.headers.Authorization = ` + "`" + `Bearer ${token}` + "`" + `;
            resolve(apiClient(original));
          });
        });
      }

      isRefreshing = true;

      try {
        const refreshToken = await getToken("refresh_token");
        if (!refreshToken) throw new Error("No refresh token");

        const { data } = await axios.post(` + "`" + `${API_URL}/auth/refresh` + "`" + `, { refresh_token: refreshToken });

        // Same wrapper as login: { data: { tokens: { access_token, ... } } }
        const tokens = data.data.tokens;

        await setToken("access_token", tokens.access_token);
        await setToken("refresh_token", tokens.refresh_token);

        refreshQueue.forEach((cb) => cb(tokens.access_token));
        refreshQueue = [];

        if (original.headers) original.headers.Authorization = ` + "`" + `Bearer ${tokens.access_token}` + "`" + `;
        return apiClient(original);
      } catch (refreshErr) {
        await deleteToken("access_token");
        await deleteToken("refresh_token");
        // Let the UI handle redirect to login via auth state
        return Promise.reject(refreshErr);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);
`
}

func desktopClientAuthProvider() string {
	return `import { createContext, useContext, useEffect, useState } from "react";

export interface User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  role: string;
}

interface AuthContextValue {
  isHydrated: boolean;
}

const AuthContext = createContext<AuthContextValue>({ isHydrated: false });

export function useAuthContext() {
  return useContext(AuthContext);
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isHydrated, setIsHydrated] = useState(false);

  // Give useAuth hooks a chance to initialize before rendering routes.
  useEffect(() => {
    setIsHydrated(true);
  }, []);

  return (
    <AuthContext.Provider value={{ isHydrated }}>
      {children}
    </AuthContext.Provider>
  );
}
`
}

func desktopClientQueryClient() string {
	return `import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      retry: 1,
      refetchOnWindowFocus: false,
    },
    mutations: {
      retry: 0,
    },
  },
});
`
}

func desktopClientThemeProvider() string {
	return `import { createContext, useContext, useEffect, useState } from "react";
import { activeTheme, applyThemeVars } from "@/lib/theme-tokens";

type Theme = "light" | "dark";

interface ThemeContextValue {
  /** Light/dark colour mode. */
  theme: Theme;
  setTheme: (theme: Theme) => void;
  /** The active brand theme (atlas | aurora | pulse) from @repo/shared. */
  tokens: typeof activeTheme;
}

const ThemeContext = createContext<ThemeContextValue>({
  theme: "dark",
  setTheme: () => {},
  tokens: activeTheme,
});

export function useTheme() {
  return useContext(ThemeContext);
}

// ThemeProvider owns two orthogonal things:
//   1. the brand theme (atlas/aurora/pulse) — fixed at scaffold time, shared
//      with the admin panel via packages/shared/themes.ts
//   2. the light/dark colour mode — a desktop affordance the user toggles
//
// It writes both onto <html> as CSS variables, which every Tailwind colour
// utility in this app already reads. That's why changing the theme restyles
// the dashboard, settings, sidebar and generated resource screens without
// touching a single page.
export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof window !== "undefined") {
      const saved = localStorage.getItem("grit-theme") as Theme | null;
      if (saved === "light" || saved === "dark") return saved;
    }
    // Light-first mirrors the admin panel, whose themes are light palettes.
    return "light";
  });

  useEffect(() => {
    const root = document.documentElement;
    if (theme === "dark") root.classList.add("dark");
    else root.classList.remove("dark");
    applyThemeVars(theme);
    localStorage.setItem("grit-theme", theme);
  }, [theme]);

  return (
    <ThemeContext.Provider value={{ theme, setTheme: setThemeState, tokens: activeTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}
`
}

func desktopClientUseShortcuts() string {
	return `import { useEffect } from "react";

// Register global keyboard shortcuts. "mod+k" = Cmd+K on Mac, Ctrl+K on Win/Linux.
export function useShortcuts(bindings: Record<string, () => void>) {
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      const mod = e.metaKey || e.ctrlKey;
      const key = e.key.toLowerCase();

      for (const [combo, callback] of Object.entries(bindings)) {
        const parts = combo.toLowerCase().split("+");
        const needsMod = parts.includes("mod");
        const needsShift = parts.includes("shift");
        const needsAlt = parts.includes("alt");
        const targetKey = parts[parts.length - 1];

        if (needsMod && !mod) continue;
        if (needsShift !== e.shiftKey) continue;
        if (needsAlt !== e.altKey) continue;
        if (targetKey !== key) continue;

        e.preventDefault();
        callback();
        return;
      }
    };

    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [bindings]);
}
`
}

func desktopClientUtils() string {
	return `import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
`
}

// ═══════════════════════════════════════════════════════════════════
// Desktop primitives — TwoPane / Form / FilterChip
// Pattern lifted from the rental-manager project. Every desktop CRUD
// page reinvents these otherwise.
// ═══════════════════════════════════════════════════════════════════

func desktopClientTwoPane() string {
	return `import { Search, Plus, Inbox } from "lucide-react";
import { cn } from "@/lib/utils";

// TwoPane is the master-detail layout that almost every desktop CRUD
// page wants: a fixed-width list on the left (TwoPane.List) and a
// detail view on the right (TwoPane.Detail) that fills the rest.
//
// Usage:
//   <TwoPane>
//     <ListPane title="Buildings" search={s} onSearch={setS} onNew={...}>
//       {items.map((b) => <ListRow ... />)}
//     </ListPane>
//     <DetailPane empty={!selected}>{selected ? ... : null}</DetailPane>
//   </TwoPane>
export function TwoPane({ children }: { children: React.ReactNode }) {
  return <div className="flex h-full overflow-hidden">{children}</div>;
}

interface ListPaneProps {
  title: string;
  count?: number;
  onNew?: () => void;
  newLabel?: string;
  search: string;
  onSearch: (v: string) => void;
  searchPlaceholder?: string;
  filters?: React.ReactNode;
  children: React.ReactNode;
  footer?: React.ReactNode;
  // Optional slot rendered before the New button (e.g. a refresh button).
  toolbar?: React.ReactNode;
}

// ListPane: sticky header (title + search + new) + scrollable body.
export function ListPane({
  title,
  count,
  onNew,
  newLabel = "New",
  search,
  onSearch,
  searchPlaceholder = "Search...",
  filters,
  children,
  footer,
  toolbar,
}: ListPaneProps) {
  return (
    <aside className="w-listpane shrink-0 border-r border-border bg-surface flex flex-col">
      <div className="px-4 pt-4 pb-2 space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-baseline gap-2">
            <h2 className="text-[15px] font-semibold text-foreground">{title}</h2>
            {count !== undefined && (
              <span className="text-[12px] text-foreground-muted">{count}</span>
            )}
          </div>
          <div className="flex items-center gap-1.5">
            {toolbar}
            {onNew && (
              <button
                type="button"
                onClick={onNew}
                className="h-8 px-2.5 rounded-lg bg-accent text-white text-[12px] font-medium hover:bg-accent-hover transition-colors inline-flex items-center gap-1"
              >
                <Plus className="h-3.5 w-3.5" />
                {newLabel}
              </button>
            )}
          </div>
        </div>

        <div className="relative">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-foreground-muted" />
          <input
            type="search"
            value={search}
            onChange={(e) => onSearch(e.target.value)}
            placeholder={searchPlaceholder}
            className="w-full h-8 pl-8 pr-2.5 rounded-lg border border-border bg-surface-2 text-[12.5px] placeholder:text-foreground-muted focus:border-accent focus:outline-none"
          />
        </div>

        {filters}
      </div>
      <div className="flex-1 overflow-y-auto">{children}</div>
      {footer && <div className="px-4 py-2 border-t border-border-subtle">{footer}</div>}
    </aside>
  );
}

interface ListRowProps {
  onClick?: () => void;
  selected?: boolean;
  icon?: React.ReactNode;
  title: React.ReactNode;
  subtitle?: React.ReactNode;
  rightTop?: React.ReactNode;
  rightBottom?: React.ReactNode;
}

// ListRow: avatar/icon + title/subtitle + right meta. Selected state
// shows a 2px accent bar on the left edge.
export function ListRow({
  onClick,
  selected,
  icon,
  title,
  subtitle,
  rightTop,
  rightBottom,
}: ListRowProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={cn(
        "relative w-full flex items-start gap-3 py-3 px-4 text-left transition-colors border-b border-border-subtle",
        "hover:bg-surface-hover",
        selected && "bg-surface-hover"
      )}
    >
      {selected && (
        <span aria-hidden className="absolute inset-y-0 left-0 w-[2px] bg-accent" />
      )}
      {icon && (
        <div className="h-9 w-9 rounded-full bg-surface-2 shrink-0 flex items-center justify-center text-[12px] font-semibold text-foreground-secondary overflow-hidden">
          {icon}
        </div>
      )}
      <div className="flex-1 min-w-0">
        <div className="text-[13.5px] font-medium text-foreground truncate">{title}</div>
        {subtitle && (
          <div className="text-[12px] text-foreground-muted truncate mt-0.5">
            {subtitle}
          </div>
        )}
      </div>
      {(rightTop || rightBottom) && (
        <div className="flex flex-col items-end gap-1 shrink-0">
          {rightTop && <div className="text-[11.5px] text-foreground-muted">{rightTop}</div>}
          {rightBottom}
        </div>
      )}
    </button>
  );
}

// DetailPane: sticky header + scrollable content. When empty=true,
// shows a centered EmptyState instead.
export function DetailPane({
  header,
  children,
  empty,
  emptyTitle = "Nothing selected",
  emptyHint = "Pick an item from the list, or create a new one.",
}: {
  header?: React.ReactNode;
  children?: React.ReactNode;
  empty?: boolean;
  emptyTitle?: string;
  emptyHint?: string;
}) {
  if (empty) {
    return (
      <section className="flex-1 flex items-center justify-center bg-background">
        <EmptyState title={emptyTitle} hint={emptyHint} />
      </section>
    );
  }
  return (
    <section className="flex-1 flex flex-col min-w-0 bg-background overflow-hidden">
      {header && (
        <div className="border-b border-border bg-surface px-6 py-4">{header}</div>
      )}
      <div className="flex-1 overflow-y-auto p-6">{children}</div>
    </section>
  );
}

// EmptyState: inbox icon + title + hint + optional action.
export function EmptyState({
  title,
  hint,
  action,
}: {
  title: string;
  hint?: string;
  action?: React.ReactNode;
}) {
  return (
    <div className="flex flex-col items-center text-center max-w-sm px-6">
      <div className="h-14 w-14 rounded-full bg-surface-2 flex items-center justify-center mb-4">
        <Inbox className="h-6 w-6 text-foreground-muted" />
      </div>
      <div className="text-[14px] font-semibold text-foreground">{title}</div>
      {hint && <div className="text-[13px] text-foreground-secondary mt-1">{hint}</div>}
      {action && <div className="mt-4">{action}</div>}
    </div>
  );
}

// DetailSection: small caps section header inside a detail pane.
export function DetailSection({
  title,
  action,
  children,
}: {
  title: string;
  action?: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <section className="mb-6">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">
          {title}
        </h3>
        {action}
      </div>
      {children}
    </section>
  );
}

// DetailField: labelled value row, used heavily in detail panes.
export function DetailField({
  label,
  value,
}: {
  label: string;
  value: React.ReactNode;
}) {
  return (
    <div>
      <div className="text-[11px] font-medium uppercase tracking-wider text-foreground-muted mb-0.5">
        {label}
      </div>
      <div className="text-[13.5px] text-foreground">
        {value || <span className="text-foreground-muted">--</span>}
      </div>
    </div>
  );
}
`
}

func desktopClientForm() string {
	return `import { forwardRef } from "react";
import { cn } from "@/lib/utils";

// FieldWrap is the shared chrome around every form field: label,
// required marker, child input, hint OR error message.
function FieldWrap({
  label,
  hint,
  error,
  required,
  className,
  children,
}: {
  label?: string;
  hint?: string;
  error?: string;
  required?: boolean;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <label className={cn("block space-y-1", className)}>
      {label && (
        <span className="block text-[12.5px] font-medium text-foreground-secondary">
          {label}
          {required && <span className="text-danger ml-0.5">*</span>}
        </span>
      )}
      {children}
      {error ? (
        <span className="block text-[11.5px] text-danger">{error}</span>
      ) : hint ? (
        <span className="block text-[11.5px] text-foreground-muted">{hint}</span>
      ) : null}
    </label>
  );
}

const baseInput =
  "w-full h-10 px-3 rounded-lg border border-border bg-surface text-[13.5px] text-foreground placeholder:text-foreground-muted focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15 disabled:bg-surface-2 disabled:text-foreground-muted";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  hint?: string;
  error?: string;
}

export const TextField = forwardRef<HTMLInputElement, InputProps>(
  ({ label, hint, error, required, className, ...props }, ref) => (
    <FieldWrap label={label} hint={hint} error={error} required={required}>
      <input
        ref={ref}
        required={required}
        className={cn(baseInput, className)}
        {...props}
      />
    </FieldWrap>
  )
);
TextField.displayName = "TextField";

interface TextAreaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  hint?: string;
  error?: string;
}

export const TextAreaField = forwardRef<HTMLTextAreaElement, TextAreaProps>(
  ({ label, hint, error, required, className, ...props }, ref) => (
    <FieldWrap label={label} hint={hint} error={error} required={required}>
      <textarea
        ref={ref}
        required={required}
        className={cn(
          baseInput,
          "h-auto min-h-[90px] py-2 resize-y",
          className
        )}
        {...props}
      />
    </FieldWrap>
  )
);
TextAreaField.displayName = "TextAreaField";

interface SelectOpt {
  value: string;
  label: string;
}
interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label?: string;
  hint?: string;
  error?: string;
  options: SelectOpt[];
  placeholder?: string;
}

export const SelectField = forwardRef<HTMLSelectElement, SelectProps>(
  (
    { label, hint, error, required, className, options, placeholder, ...props },
    ref
  ) => (
    <FieldWrap label={label} hint={hint} error={error} required={required}>
      <select
        ref={ref}
        required={required}
        className={cn(baseInput, "appearance-none pr-9 bg-no-repeat", className)}
        style={{
          backgroundImage:
            "url(\"data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%238F95A3' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'/%3E%3C/svg%3E\")",
          backgroundPosition: "right 12px center",
        }}
        {...props}
      >
        {placeholder && <option value="">{placeholder}</option>}
        {options.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
    </FieldWrap>
  )
);
SelectField.displayName = "SelectField";

// FormGrid lays fields out in 1, 2, or 3 columns on >= sm breakpoints,
// stacking to one column on small screens.
export function FormGrid({
  columns = 2,
  children,
}: {
  columns?: 1 | 2 | 3;
  children: React.ReactNode;
}) {
  const cls =
    columns === 1
      ? "grid grid-cols-1 gap-4"
      : columns === 3
      ? "grid grid-cols-1 sm:grid-cols-3 gap-4"
      : "grid grid-cols-1 sm:grid-cols-2 gap-4";
  return <div className={cls}>{children}</div>;
}

// FormSection groups related fields under a small caps title.
export function FormSection({
  title,
  description,
  children,
}: {
  title?: string;
  description?: string;
  children: React.ReactNode;
}) {
  return (
    <section className="space-y-3">
      {(title || description) && (
        <div>
          {title && (
            <h3 className="text-[11px] font-semibold uppercase tracking-wider text-foreground-muted">
              {title}
            </h3>
          )}
          {description && (
            <p className="text-[12.5px] text-foreground-secondary mt-0.5">
              {description}
            </p>
          )}
        </div>
      )}
      {children}
    </section>
  );
}

// FormActions is the sticky-footer cancel/submit pair every form needs.
// Pass isPending from your TanStack mutation to disable + show "Saving..."
export function FormActions({
  onCancel,
  submitLabel = "Save",
  isPending,
  extra,
}: {
  onCancel?: () => void;
  submitLabel?: string;
  isPending?: boolean;
  extra?: React.ReactNode;
}) {
  return (
    <div className="flex items-center justify-between pt-4 border-t border-border-subtle">
      <div>{extra}</div>
      <div className="flex gap-2">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="h-9 px-3.5 rounded-lg border border-border bg-surface text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover"
          >
            Cancel
          </button>
        )}
        <button
          type="submit"
          disabled={isPending}
          className="h-9 px-3.5 rounded-lg bg-accent text-white text-[13px] font-medium hover:bg-accent-hover disabled:opacity-60"
        >
          {isPending ? "Saving..." : submitLabel}
        </button>
      </div>
    </div>
  );
}
`
}

func desktopClientFilterChip() string {
	return `import { X } from "lucide-react";
import { cn } from "@/lib/utils";

// FilterChip is a toggleable pill used to compose filter bars on
// desktop list views. Active chips show a subtle accent background
// and an X to clear; inactive chips look like neutral tags.
//
// Pair with FilterBar for the standard horizontal scrollable layout:
//
//   <FilterBar>
//     <FilterChip active={status === "all"} onClick={() => setStatus("all")}>All</FilterChip>
//     <FilterChip active={status === "open"} onClick={() => setStatus("open")} onClear={() => setStatus("all")}>Open</FilterChip>
//   </FilterBar>
export function FilterChip({
  children,
  active,
  onClick,
  onClear,
  count,
}: {
  children: React.ReactNode;
  active?: boolean;
  onClick?: () => void;
  // When provided AND active, an X button appears that calls onClear instead
  // of toggling the chip. Lets users clear a single filter without opening a menu.
  onClear?: () => void;
  // Optional count rendered after the label (e.g. "Open (3)").
  count?: number;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1 h-7 pl-2.5 rounded-full text-[12px] font-medium border transition-colors select-none",
        onClear ? "pr-1" : "pr-2.5",
        active
          ? "bg-accent/15 border-accent/30 text-accent"
          : "bg-surface border-border text-foreground-secondary hover:bg-surface-hover hover:text-foreground cursor-pointer"
      )}
      onClick={!active ? onClick : undefined}
      role={!active ? "button" : undefined}
      tabIndex={!active ? 0 : undefined}
      onKeyDown={(e) => {
        if (!active && onClick && (e.key === "Enter" || e.key === " ")) {
          e.preventDefault();
          onClick();
        }
      }}
    >
      <span>{children}</span>
      {count !== undefined && (
        <span className={cn("text-[11px]", active ? "text-accent/70" : "text-foreground-muted")}>
          {count}
        </span>
      )}
      {active && onClear && (
        <button
          type="button"
          onClick={onClear}
          className="ml-0.5 h-5 w-5 rounded-full hover:bg-accent/20 inline-flex items-center justify-center"
          aria-label="Clear filter"
        >
          <X className="h-3 w-3" />
        </button>
      )}
    </span>
  );
}

// FilterBar wraps a row of FilterChips with horizontal scroll for
// when there are too many to fit. Use directly inside a ListPane
// via the filters prop.
export function FilterBar({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex items-center gap-1.5 overflow-x-auto -mx-1 px-1 pb-1 scrollbar-thin">
      {children}
    </div>
  );
}
`
}

// ═══════════════════════════════════════════════════════════════════
// Realtime client (WebSocket + reconnect + EventTarget bus)
// ═══════════════════════════════════════════════════════════════════

func desktopClientRealtimeTS() string {
	return `// Tiny WebSocket client for the API's realtime hub. Auto-reconnects with
// exponential backoff. Events are fanned out via a global EventTarget bus
// so any component can subscribe with realtimeBus.addEventListener(...) or
// the useRealtimeEvent hook.
import { getToken } from "./wails-bridge";

export interface RealtimeEvent<T = unknown> {
  type: string;
  payload: T;
}

// Resolve the ws:// URL from the same env the API client uses.
function resolveWsUrl(): string {
  const apiUrl =
    (import.meta.env.VITE_API_URL as string | undefined) ||
    (typeof window !== "undefined" && (window as any).go?.main?.App
      ? "http://localhost:8080/api"
      : "/api");
  if (apiUrl.startsWith("http://") || apiUrl.startsWith("https://")) {
    return apiUrl.replace(/^http/, "ws") + "/ws";
  }
  // Relative — resolve against window.location.
  const base =
    window.location.protocol === "https:"
      ? "wss://" + window.location.host
      : "ws://" + window.location.host;
  return base + apiUrl + "/ws";
}

// Global event bus — any component can listen for "chat.message.new",
// "notification.new", "system.connected", or your own resource events.
// The "*" event fires for every message, useful for debugging.
export const realtimeBus = new EventTarget();

class RealtimeClient {
  private ws: WebSocket | null = null;
  private retries = 0;
  private retryTimer: number | null = null;
  private stopped = true;

  // Call from your AuthProvider once the user has tokens. Idempotent.
  async start() {
    this.stopped = false;
    await this.connect();
  }

  // Call on logout to tear down the connection without auto-reconnecting.
  stop() {
    this.stopped = true;
    if (this.retryTimer) {
      window.clearTimeout(this.retryTimer);
      this.retryTimer = null;
    }
    if (this.ws) {
      try {
        this.ws.close();
      } catch {
        // ignore
      }
      this.ws = null;
    }
  }

  private async connect() {
    if (this.stopped) return;
    const token = await getToken("access_token");
    if (!token) {
      // Not logged in; try again in a bit.
      this.scheduleReconnect();
      return;
    }
    const url = resolveWsUrl() + "?token=" + encodeURIComponent(token);
    try {
      this.ws = new WebSocket(url);
    } catch {
      this.scheduleReconnect();
      return;
    }

    this.ws.addEventListener("open", () => {
      this.retries = 0;
    });

    this.ws.addEventListener("message", (ev) => {
      try {
        const data = JSON.parse(ev.data) as RealtimeEvent;
        // Dispatch the typed event AND a "*" wildcard for debugging/logging.
        realtimeBus.dispatchEvent(
          new CustomEvent(data.type, { detail: data.payload })
        );
        realtimeBus.dispatchEvent(new CustomEvent("*", { detail: data }));
      } catch {
        // ignore malformed messages
      }
    });

    this.ws.addEventListener("close", () => {
      this.ws = null;
      this.scheduleReconnect();
    });

    this.ws.addEventListener("error", () => {
      // The close handler will fire next and trigger reconnect.
    });
  }

  private scheduleReconnect() {
    if (this.stopped) return;
    if (this.retryTimer) return;
    // Exponential backoff: 1s, 2s, 4s, 8s, capped at 15s.
    const delay = Math.min(15_000, 1000 * Math.pow(2, this.retries));
    this.retries++;
    this.retryTimer = window.setTimeout(() => {
      this.retryTimer = null;
      this.connect();
    }, delay);
  }
}

// Singleton client. Start it from AuthProvider after tokens land,
// stop it on logout.
export const realtime = new RealtimeClient();
`
}

// ═══════════════════════════════════════════════════════════════════
// Hooks
// ═══════════════════════════════════════════════════════════════════

func desktopClientUseRealtime() string {
	return `import { useEffect } from "react";
import { realtimeBus, type RealtimeEvent } from "@/lib/realtime";

// useRealtimeEvent subscribes a callback to a single realtime event type.
// The callback is wrapped so it always sees the latest closure (no stale
// state), and unsubscribed on unmount.
//
// Usage:
//   useRealtimeEvent<ChatMessage>("chat.message.new", (msg) => {
//     queryClient.invalidateQueries({ queryKey: ["chats"] });
//   });
export function useRealtimeEvent<T = unknown>(
  type: string,
  callback: (payload: T) => void,
) {
  useEffect(() => {
    const handler = (e: Event) => {
      const ce = e as CustomEvent<T>;
      callback(ce.detail);
    };
    realtimeBus.addEventListener(type, handler);
    return () => realtimeBus.removeEventListener(type, handler);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [type, callback]);
}

// useRealtimeAny fires for EVERY message — useful for an in-app toast bar
// or a debug console. Receives the full envelope, not just the payload.
export function useRealtimeAny(callback: (event: RealtimeEvent) => void) {
  useEffect(() => {
    const handler = (e: Event) => {
      const ce = e as CustomEvent<RealtimeEvent>;
      callback(ce.detail);
    };
    realtimeBus.addEventListener("*", handler);
    return () => realtimeBus.removeEventListener("*", handler);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [callback]);
}
`
}

func desktopClientUseOnlineStatus() string {
	return `import { useEffect, useState } from "react";
import { apiClient } from "@/lib/api-client";

// useOnlineStatus tracks whether the desktop client can reach the API.
//
// "Online" here means two things together:
//   1) the OS reports a network connection (navigator.onLine), AND
//   2) the API has answered a heartbeat in the last HEARTBEAT_INTERVAL_MS
//
// We need both because navigator.onLine only reflects the OS-level network
// link — it returns true even when the API is down or the user is on a
// captive-portal Wi-Fi that hasn't been authenticated. The API heartbeat
// is the truth signal; navigator.onLine is just a cheap pre-check that
// stops us from making a doomed request when the laptop lid was just
// closed.
//
// The heartbeat hits GET /api/health which Grit's API exposes for free —
// no auth required, no DB hit, fast 200 OK.

const HEARTBEAT_INTERVAL_MS = 15_000;
const HEARTBEAT_TIMEOUT_MS = 5_000;

export function useOnlineStatus() {
  const [isOnline, setIsOnline] = useState<boolean>(
    typeof navigator !== "undefined" ? navigator.onLine : true,
  );
  const [lastCheckedAt, setLastCheckedAt] = useState<Date | null>(null);

  useEffect(() => {
    let cancelled = false;
    let timer: ReturnType<typeof setInterval> | null = null;

    const ping = async () => {
      // OS says no network → don't even try.
      if (typeof navigator !== "undefined" && !navigator.onLine) {
        if (!cancelled) {
          setIsOnline(false);
          setLastCheckedAt(new Date());
        }
        return;
      }
      try {
        const controller = new AbortController();
        const timeoutId = setTimeout(() => controller.abort(), HEARTBEAT_TIMEOUT_MS);
        await apiClient.get("/health", { signal: controller.signal, timeout: HEARTBEAT_TIMEOUT_MS });
        clearTimeout(timeoutId);
        if (!cancelled) {
          setIsOnline(true);
          setLastCheckedAt(new Date());
        }
      } catch {
        if (!cancelled) {
          setIsOnline(false);
          setLastCheckedAt(new Date());
        }
      }
    };

    ping();
    timer = setInterval(ping, HEARTBEAT_INTERVAL_MS);

    const handleOnline = () => ping();
    const handleOffline = () => {
      if (!cancelled) {
        setIsOnline(false);
        setLastCheckedAt(new Date());
      }
    };

    if (typeof window !== "undefined") {
      window.addEventListener("online", handleOnline);
      window.addEventListener("offline", handleOffline);
    }

    return () => {
      cancelled = true;
      if (timer) clearInterval(timer);
      if (typeof window !== "undefined") {
        window.removeEventListener("online", handleOnline);
        window.removeEventListener("offline", handleOffline);
      }
    };
  }, []);

  return { isOnline, lastCheckedAt };
}
`
}

func desktopClientUseAuth() string {
	return `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { setToken, deleteToken } from "@/lib/wails-bridge";
import type { User } from "@/lib/auth-provider";

interface LoginInput {
  email: string;
  password: string;
}

interface RegisterInput {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
}

// The API wraps auth payloads as:
//   { "data": { "user": {...}, "tokens": { access_token, refresh_token, expires_at } } }
// apiClient hands us the body, so mutationFn returns data.data — i.e. this shape.
// Reading access_token off the top level (the old bug) stored an undefined token,
// which made /app's beforeLoad bounce straight back to /auth/login after a
// *successful* login.
interface AuthTokens {
  access_token: string;
  refresh_token: string;
  expires_at: number;
}

interface AuthResponse {
  user: User;
  tokens: AuthTokens;
}

export function useMe() {
  return useQuery<User>({
    queryKey: ["me"],
    queryFn: async () => {
      const { data } = await apiClient.get("/auth/me");
      return data.data;
    },
    retry: false,
  });
}

export function useLogin() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: LoginInput): Promise<AuthResponse> => {
      const { data } = await apiClient.post("/auth/login", input);
      return data.data;
    },
    onSuccess: async (data) => {
      await setToken("access_token", data.tokens.access_token);
      await setToken("refresh_token", data.tokens.refresh_token);
      qc.setQueryData(["me"], data.user);
    },
  });
}

export function useRegister() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: RegisterInput): Promise<AuthResponse> => {
      const { data } = await apiClient.post("/auth/register", input);
      return data.data;
    },
    onSuccess: async (data) => {
      await setToken("access_token", data.tokens.access_token);
      await setToken("refresh_token", data.tokens.refresh_token);
      qc.setQueryData(["me"], data.user);
    },
  });
}

export function useLogout() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      try {
        await apiClient.post("/auth/logout");
      } catch {
        // Best-effort — continue cleanup even if server call fails
      }
      await deleteToken("access_token");
      await deleteToken("refresh_token");
    },
    onSuccess: () => {
      qc.clear();
    },
  });
}
`
}
