#!/usr/bin/env sh
# Grit installer — macOS + Linux
#
#   curl -fsSL https://gritframework.dev/install.sh | sh
#
# Behaviour mirrors `winget install`:
#   - If grit is already on PATH, runs `grit update` (which is idempotent —
#     no-ops when already on latest, otherwise self-updates in place).
#   - If grit is not installed, downloads the matching release binary for
#     your OS/arch from GitHub and installs it to /usr/local/bin/grit
#     (writable) or ~/.local/bin/grit (fallback).
#
# Honors:
#   GRIT_INSTALL_DIR    install destination (overrides auto-detection)
#   GRIT_VERSION        pin a specific version (default: latest)

set -eu

REPO="MUKE-coder/grit"
COLOR_BLUE="$(printf '\033[34m')"
COLOR_GREEN="$(printf '\033[32m')"
COLOR_YELLOW="$(printf '\033[33m')"
COLOR_RED="$(printf '\033[31m')"
COLOR_RESET="$(printf '\033[0m')"

info()  { printf "%s>>%s %s\n" "$COLOR_BLUE"  "$COLOR_RESET" "$*"; }
ok()    { printf "%s✓%s  %s\n" "$COLOR_GREEN" "$COLOR_RESET" "$*"; }
warn()  { printf "%s!%s  %s\n" "$COLOR_YELLOW" "$COLOR_RESET" "$*" >&2; }
die()   { printf "%sx%s  %s\n" "$COLOR_RED"   "$COLOR_RESET" "$*" >&2; exit 1; }

# ─── 1. Already installed? Hand off to `grit update` ──────────────────
if command -v grit >/dev/null 2>&1; then
  current="$(grit version 2>/dev/null | awk '{print $NF}' | sed 's/^v//')"
  info "Found existing grit v${current:-unknown} — running \`grit update\`"
  printf "\n"
  exec grit update
fi

# ─── 2. Detect OS + arch ──────────────────────────────────────────────
os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  darwin) ;;
  linux)  ;;
  *)      die "Unsupported OS: $os. Open an issue at https://github.com/$REPO/issues." ;;
esac

arch="$(uname -m)"
case "$arch" in
  x86_64 | amd64) arch="amd64" ;;
  arm64  | aarch64) arch="arm64" ;;
  *) die "Unsupported architecture: $arch. Open an issue at https://github.com/$REPO/issues." ;;
esac

# ─── 3. Resolve version ───────────────────────────────────────────────
version="${GRIT_VERSION:-latest}"
if [ "$version" = "latest" ]; then
  info "Resolving latest release..."
  tag="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep -m1 '"tag_name":' \
    | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
  [ -n "$tag" ] || die "Could not determine latest version from GitHub API."
else
  tag="$version"
  case "$tag" in v*) ;; *) tag="v$tag" ;; esac
fi
ok "Installing grit $tag (${os}/${arch})"

# ─── 4. Pick install dir ──────────────────────────────────────────────
install_dir=""
if [ -n "${GRIT_INSTALL_DIR:-}" ]; then
  install_dir="$GRIT_INSTALL_DIR"
elif [ -w "/usr/local/bin" ]; then
  install_dir="/usr/local/bin"
elif [ -d "$HOME/.local/bin" ] || mkdir -p "$HOME/.local/bin" 2>/dev/null; then
  install_dir="$HOME/.local/bin"
else
  die "Cannot find a writable install directory. Set GRIT_INSTALL_DIR=/some/path and retry."
fi

# ─── 5. Download + extract ────────────────────────────────────────────
archive="grit-${os}-${arch}.tar.gz"
url="https://github.com/$REPO/releases/download/$tag/$archive"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

info "Downloading $archive"
if ! curl -fsSL "$url" -o "$tmp/$archive"; then
  die "Download failed: $url"
fi

info "Extracting..."
tar -xzf "$tmp/$archive" -C "$tmp"

# Release archives contain a binary named "grit-${os}-${arch}" — find it.
src="$(find "$tmp" -maxdepth 2 -name "grit-${os}-${arch}" -type f | head -n1)"
[ -n "$src" ] || src="$(find "$tmp" -maxdepth 2 -name 'grit' -type f | head -n1)"
[ -n "$src" ] || die "Could not find grit binary inside the archive."

chmod +x "$src"

# ─── 6. Move into place ───────────────────────────────────────────────
dest="$install_dir/grit"
info "Installing to $dest"
if ! mv "$src" "$dest" 2>/dev/null; then
  # /usr/local/bin frequently needs sudo on Linux even when w-checked.
  if command -v sudo >/dev/null 2>&1; then
    info "Elevating with sudo for $install_dir"
    sudo mv "$src" "$dest"
  else
    die "Failed to install to $dest and sudo is unavailable. Set GRIT_INSTALL_DIR to a writable directory."
  fi
fi

# ─── 7. PATH sanity check ─────────────────────────────────────────────
case ":$PATH:" in
  *":$install_dir:"*) ;;
  *)
    warn "$install_dir is not in your PATH."
    warn "Add this line to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
    printf "\n    export PATH=\"%s:\$PATH\"\n\n" "$install_dir"
    ;;
esac

# ─── 8. Done ──────────────────────────────────────────────────────────
ok "grit installed to $dest"
printf "\n"
"$dest" version 2>/dev/null || true
printf "\n"
info "Next: run \`grit new my-app\` to scaffold your first project."
info "Update any time with: \`grit update\`"
