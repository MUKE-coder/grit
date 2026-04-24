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
	module := opts.Module() + "/apps/desktop"

	files := map[string]string{
		// Go side (minimal — native OS only)
		filepath.Join(desktopRoot, "main.go"):               desktopClientMainGo(opts),
		filepath.Join(desktopRoot, "app.go"):                desktopClientAppGo(),
		filepath.Join(desktopRoot, "go.mod"):                desktopClientGoMod(module),
		filepath.Join(desktopRoot, "wails.json"):            desktopClientWailsJSON(opts),
		filepath.Join(desktopRoot, "internal", "keychain.go"): desktopClientKeychainGo(opts),
		filepath.Join(desktopRoot, ".gitignore"):            desktopClientGitignore(),
		filepath.Join(desktopRoot, "README.md"):             desktopClientReadme(opts),

		// Frontend config
		filepath.Join(desktopRoot, "frontend", "package.json"):     desktopClientPackageJSON(opts),
		filepath.Join(desktopRoot, "frontend", "tsconfig.json"):    desktopClientTSConfig(),
		filepath.Join(desktopRoot, "frontend", "tsconfig.node.json"): desktopClientTSConfigNode(),
		filepath.Join(desktopRoot, "frontend", "vite.config.ts"):   desktopClientViteConfig(),
		filepath.Join(desktopRoot, "frontend", "tailwind.config.ts"): desktopClientTailwindConfig(),
		filepath.Join(desktopRoot, "frontend", "postcss.config.js"): desktopClientPostCSSConfig(),
		filepath.Join(desktopRoot, "frontend", "index.html"):       desktopClientIndexHTML(opts),
		filepath.Join(desktopRoot, "frontend", "src", "main.tsx"):  desktopClientMainTSX(),
		filepath.Join(desktopRoot, "frontend", "src", "globals.css"): desktopClientGlobalsCSS(),
		filepath.Join(desktopRoot, "frontend", "src", "vite-env.d.ts"): desktopClientViteEnvDTS(),

		// Routes (TanStack Router file-based)
		filepath.Join(desktopRoot, "frontend", "src", "routes", "__root.tsx"):      desktopClientRootRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "index.tsx"):       desktopClientIndexRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "_auth.tsx"):       desktopClientAuthLayoutRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "_auth", "login.tsx"):    desktopClientLoginRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "_auth", "register.tsx"): desktopClientRegisterRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "_app.tsx"):        desktopClientAppLayoutRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "_app", "index.tsx"):    desktopClientDashboardRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "_app", "profile.tsx"):  desktopClientProfileRoute(),
		filepath.Join(desktopRoot, "frontend", "src", "routes", "_app", "settings.tsx"): desktopClientSettingsRoute(),

		// Layout components
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "title-bar.tsx"):       desktopClientTitleBar(opts),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "sidebar.tsx"):         desktopClientSidebar(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "topbar.tsx"):          desktopClientTopbar(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "page-header.tsx"):     desktopClientPageHeader(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "layout", "command-palette.tsx"): desktopClientCommandPalette(),

		// UI components
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "button.tsx"):       desktopClientButton(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "input.tsx"):        desktopClientInput(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "card.tsx"):         desktopClientCard(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "dialog.tsx"):       desktopClientDialog(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "avatar.tsx"):       desktopClientAvatar(),
		filepath.Join(desktopRoot, "frontend", "src", "components", "ui", "badge.tsx"):        desktopClientBadge(),

		// Lib (API client, auth, utils, shortcuts, Wails bridge)
		filepath.Join(desktopRoot, "frontend", "src", "lib", "api-client.ts"):     desktopClientApiClientTS(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "auth-provider.tsx"): desktopClientAuthProvider(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "query-client.ts"):   desktopClientQueryClient(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "theme-provider.tsx"): desktopClientThemeProvider(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "use-shortcuts.ts"):  desktopClientUseShortcuts(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "wails-bridge.ts"):   desktopClientWailsBridge(),
		filepath.Join(desktopRoot, "frontend", "src", "lib", "utils.ts"):          desktopClientUtils(),

		// Hooks
		filepath.Join(desktopRoot, "frontend", "src", "hooks", "use-auth.ts"): desktopClientUseAuth(),
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
	goruntime "runtime"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App exposes a minimal set of native OS methods to the React frontend
// via Wails bindings. All business logic (CRUD, auth, etc.) goes through
// HTTP to the shared API — NOT through these methods.
type App struct {
	ctx context.Context
	kc  *Keychain
}

func NewApp() *App {
	return &App{
		kc: NewKeychain(),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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
	github.com/wailsapp/wails/v2 v2.9.2
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

func desktopClientPackageJSON(opts Options) string {
	return `{
  "name": "@` + opts.ProjectName + `/desktop",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
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
        sans: ["Onest", "system-ui", "sans-serif"],
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
        DEFAULT: "0.5rem",
      },
      spacing: {
        // Desktop uses larger padding than web — more breathing room
        "content": "2rem",     // 32px main content padding
        "sidebar": "15rem",    // 240px fixed sidebar
        "titlebar": "3rem",    // 48px custom title bar
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
<html lang="en" class="dark">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>` + opts.ProjectName + `</title>
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link href="https://fonts.googleapis.com/css2?family=Onest:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet" />
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

// Wails bindings are injected at runtime. This is a placeholder type
// so TypeScript doesn't complain before the first build.
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

export const Route = createFileRoute("/_auth")({
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
      <div className="flex-1 flex items-center justify-center overflow-auto">
        <Outlet />
      </div>
    </div>
  );
}
`
}

func desktopClientLoginRoute() string {
	return `import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useState } from "react";
import { Loader2, Eye, EyeOff } from "lucide-react";
import { useLogin } from "@/hooks/use-auth";

export const Route = createFileRoute("/_auth/login")({
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
    <div className="relative w-full max-w-[420px] px-4">
      {/* Subtle radial glow */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-x-0 -top-24 h-[400px] opacity-60"
        style={{
          background: "radial-gradient(600px at 50% 0%, rgba(108, 92, 231, 0.12), transparent 70%)",
        }}
      />

      <div className="relative rounded-2xl border border-border bg-surface p-10 shadow-sm">
        <div className="flex justify-center">
          <div className="h-10 w-10 rounded-xl bg-accent flex items-center justify-center shadow-sm">
            <span className="text-[17px] font-bold text-white">G</span>
          </div>
        </div>

        <div className="mt-5 text-center">
          <h1 className="text-[22px] font-semibold text-foreground tracking-tight">
            Welcome back
          </h1>
          <p className="mt-1.5 text-[13px] text-foreground-secondary">
            Sign in to continue
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="mt-7 space-y-5">
          {error && (
            <div className="rounded-lg border border-danger/30 bg-danger/10 px-3.5 py-2.5 text-[13px] text-danger">
              {(error as Error).message || "Invalid email or password"}
            </div>
          )}

          <div className="space-y-1.5">
            <label htmlFor="email" className="block text-[13px] font-medium text-foreground-secondary">
              Email
            </label>
            <input
              id="email"
              type="email"
              autoComplete="email"
              placeholder="you@example.com"
              className="w-full h-11 rounded-lg border border-border bg-surface-2 px-3.5 text-[14px] text-foreground placeholder:text-foreground-muted focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/20"
              {...register("email")}
            />
            {errors.email && <p className="text-[12px] text-danger">{errors.email.message}</p>}
          </div>

          <div className="space-y-1.5">
            <label htmlFor="password" className="block text-[13px] font-medium text-foreground-secondary">
              Password
            </label>
            <div className="relative">
              <input
                id="password"
                type={showPassword ? "text" : "password"}
                autoComplete="current-password"
                placeholder="Enter your password"
                className="w-full h-11 rounded-lg border border-border bg-surface-2 px-3.5 pr-10 text-[14px] text-foreground placeholder:text-foreground-muted focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/20"
                {...register("password")}
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-foreground-muted hover:text-foreground"
              >
                {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
            {errors.password && <p className="text-[12px] text-danger">{errors.password.message}</p>}
          </div>

          <button
            type="submit"
            disabled={isPending}
            className="w-full h-11 rounded-lg bg-accent text-[14px] font-medium text-white hover:bg-accent-hover transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
          >
            {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            {isPending ? "Signing in..." : "Sign in"}
          </button>
        </form>

        <p className="mt-7 text-center text-[13px] text-foreground-secondary">
          Don't have an account?{" "}
          <Link to="/auth/register" className="font-medium text-accent hover:underline">
            Sign up
          </Link>
        </p>
      </div>
    </div>
  );
}
`
}

func desktopClientRegisterRoute() string {
	return `import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";
import { useRegister } from "@/hooks/use-auth";

export const Route = createFileRoute("/_auth/register")({
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

  const input = "w-full h-11 rounded-lg border border-border bg-surface-2 px-3.5 text-[14px] text-foreground placeholder:text-foreground-muted focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/20";

  return (
    <div className="relative w-full max-w-[420px] px-4">
      <div className="relative rounded-2xl border border-border bg-surface p-10 shadow-sm">
        <div className="flex justify-center">
          <div className="h-10 w-10 rounded-xl bg-accent flex items-center justify-center shadow-sm">
            <span className="text-[17px] font-bold text-white">G</span>
          </div>
        </div>

        <div className="mt-5 text-center">
          <h1 className="text-[22px] font-semibold text-foreground tracking-tight">
            Create account
          </h1>
          <p className="mt-1.5 text-[13px] text-foreground-secondary">
            Get started in 30 seconds
          </p>
        </div>

        <form onSubmit={handleSubmit(onSubmit)} className="mt-7 space-y-5">
          {error && (
            <div className="rounded-lg border border-danger/30 bg-danger/10 px-3.5 py-2.5 text-[13px] text-danger">
              {(error as Error).message || "Registration failed"}
            </div>
          )}

          <div className="grid grid-cols-2 gap-3">
            <div className="space-y-1.5">
              <label className="block text-[13px] font-medium text-foreground-secondary">
                First name
              </label>
              <input type="text" placeholder="Jane" className={input} {...register("first_name")} />
              {errors.first_name && <p className="text-[12px] text-danger">{errors.first_name.message}</p>}
            </div>
            <div className="space-y-1.5">
              <label className="block text-[13px] font-medium text-foreground-secondary">
                Last name
              </label>
              <input type="text" placeholder="Doe" className={input} {...register("last_name")} />
              {errors.last_name && <p className="text-[12px] text-danger">{errors.last_name.message}</p>}
            </div>
          </div>

          <div className="space-y-1.5">
            <label className="block text-[13px] font-medium text-foreground-secondary">Email</label>
            <input type="email" placeholder="you@example.com" className={input} {...register("email")} />
            {errors.email && <p className="text-[12px] text-danger">{errors.email.message}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="block text-[13px] font-medium text-foreground-secondary">Password</label>
            <input type="password" placeholder="Min 8 characters" className={input} {...register("password")} />
            {errors.password && <p className="text-[12px] text-danger">{errors.password.message}</p>}
          </div>

          <button
            type="submit"
            disabled={isPending}
            className="w-full h-11 rounded-lg bg-accent text-[14px] font-medium text-white hover:bg-accent-hover transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
          >
            {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            {isPending ? "Creating..." : "Create account"}
          </button>
        </form>

        <p className="mt-7 text-center text-[13px] text-foreground-secondary">
          Already have an account?{" "}
          <Link to="/auth/login" className="font-medium text-accent hover:underline">
            Sign in
          </Link>
        </p>
      </div>
    </div>
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

export const Route = createFileRoute("/_app")({
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

export const Route = createFileRoute("/_app/")({
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

export const Route = createFileRoute("/_app/profile")({
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

export const Route = createFileRoute("/_app/settings")({
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

interface TitleBarProps {
  showSidebarControls?: boolean;
}

// Frameless title bar with drag region and OS-specific window controls.
// macOS: traffic lights on the left. Windows/Linux: min/max/close on the right.
export function TitleBar({ showSidebarControls: _ = true }: TitleBarProps) {
  const [platform, setPlatform] = useState<"darwin" | "windows" | "linux">("windows");

  useEffect(() => {
    getPlatform().then(setPlatform);
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

      {/* Right: Windows/Linux window controls */}
      {!isMac ? (
        <div className="no-drag flex items-center">
          <WinControl icon={<Minus className="h-3.5 w-3.5" />} onClick={minimise} />
          <WinControl icon={<Square className="h-3 w-3" />} onClick={toggleMaximise} />
          <WinControl icon={<X className="h-4 w-4" />} onClick={closeWindow} variant="close" />
        </div>
      ) : (
        <div />
      )}
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

// Attach JWT from OS keychain (or localStorage in dev)
apiClient.interceptors.request.use(async (config: InternalAxiosRequestConfig) => {
  const token = await getToken("access_token");
  if (token && config.headers) {
    config.headers.Authorization = ` + "`" + `Bearer ${token}` + "`" + `;
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

    if (error.response?.status === 401 && !original._retry) {
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

        await setToken("access_token", data.access_token);
        await setToken("refresh_token", data.refresh_token);

        refreshQueue.forEach((cb) => cb(data.access_token));
        refreshQueue = [];

        if (original.headers) original.headers.Authorization = ` + "`" + `Bearer ${data.access_token}` + "`" + `;
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
  id: number;
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

type Theme = "light" | "dark";

interface ThemeContextValue {
  theme: Theme;
  setTheme: (theme: Theme) => void;
}

const ThemeContext = createContext<ThemeContextValue>({
  theme: "dark",
  setTheme: () => {},
});

export function useTheme() {
  return useContext(ThemeContext);
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(() => {
    if (typeof window !== "undefined") {
      const saved = localStorage.getItem("grit-theme") as Theme | null;
      if (saved === "light" || saved === "dark") return saved;
    }
    return "dark";
  });

  useEffect(() => {
    const root = document.documentElement;
    if (theme === "dark") root.classList.add("dark");
    else root.classList.remove("dark");
    localStorage.setItem("grit-theme", theme);
  }, [theme]);

  return (
    <ThemeContext.Provider value={{ theme, setTheme: setThemeState }}>
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
// Hooks
// ═══════════════════════════════════════════════════════════════════

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

interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
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
      await setToken("access_token", data.access_token);
      await setToken("refresh_token", data.refresh_token);
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
      await setToken("access_token", data.access_token);
      await setToken("refresh_token", data.refresh_token);
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
