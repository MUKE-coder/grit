# Grit installer — Windows PowerShell 5.1+ / PowerShell 7+
#
#   iwr -useb https://gritframework.dev/install.ps1 | iex
#
# Behaviour mirrors `winget install`:
#   - If grit is already on PATH, runs `grit update` (idempotent — no-ops if
#     already on latest, otherwise self-updates in place).
#   - If grit is not installed, downloads the matching release zip for your
#     architecture from GitHub and installs it to %USERPROFILE%\.grit\bin,
#     adding that directory to your user PATH if it isn't already.
#
# Honors:
#   $env:GRIT_INSTALL_DIR    install destination (overrides default)
#   $env:GRIT_VERSION        pin a specific version (default: latest)

$ErrorActionPreference = 'Stop'

$repo = 'MUKE-coder/grit'

function Info($msg)  { Write-Host ">> $msg" -ForegroundColor Cyan }
function Ok($msg)    { Write-Host "OK $msg" -ForegroundColor Green }
function Warn($msg)  { Write-Host "!  $msg" -ForegroundColor Yellow }
function Die($msg)   { Write-Host "x  $msg" -ForegroundColor Red; exit 1 }

# ─── 1. Already installed? Hand off to `grit update` ──────────────────
$existing = Get-Command grit -ErrorAction SilentlyContinue
if ($existing) {
    try {
        $current = (& $existing.Source version) -split ' ' | Select-Object -Last 1
        $current = $current.TrimStart('v')
    } catch { $current = 'unknown' }
    Info "Found existing grit v$current  running 'grit update'"
    Write-Host ''
    & $existing.Source update
    exit $LASTEXITCODE
}

# ─── 2. Detect arch ───────────────────────────────────────────────────
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }
} else {
    Die "32-bit Windows is not supported by grit releases."
}

# ─── 3. Resolve version ───────────────────────────────────────────────
$version = if ($env:GRIT_VERSION) { $env:GRIT_VERSION } else { 'latest' }
if ($version -eq 'latest') {
    Info 'Resolving latest release...'
    try {
        $rel = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest" `
            -Headers @{ 'User-Agent' = 'grit-install' }
        $tag = $rel.tag_name
    } catch { Die "Could not query GitHub releases API: $($_.Exception.Message)" }
} else {
    $tag = if ($version.StartsWith('v')) { $version } else { "v$version" }
}
Ok "Installing grit $tag (windows/$arch)"

# ─── 4. Pick install dir ──────────────────────────────────────────────
$installDir = if ($env:GRIT_INSTALL_DIR) {
    $env:GRIT_INSTALL_DIR
} else {
    Join-Path $env:USERPROFILE '.grit\bin'
}
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# ─── 5. Download + extract ────────────────────────────────────────────
$archive = "grit-windows-$arch.zip"
$url     = "https://github.com/$repo/releases/download/$tag/$archive"
$tmp     = Join-Path ([IO.Path]::GetTempPath()) ("grit-install-" + [Guid]::NewGuid())
New-Item -ItemType Directory -Path $tmp | Out-Null

try {
    Info "Downloading $archive"
    Invoke-WebRequest -Uri $url -OutFile (Join-Path $tmp $archive) -UseBasicParsing

    Info 'Extracting...'
    Expand-Archive -Path (Join-Path $tmp $archive) -DestinationPath $tmp -Force

    # Release archives ship the binary as grit-windows-{arch}.exe — fall back
    # to bare grit.exe if a future release ever re-packages it.
    $src = Get-ChildItem -Path $tmp -Recurse -File `
        | Where-Object { $_.Name -in @("grit-windows-$arch.exe", 'grit.exe') } `
        | Select-Object -First 1
    if (-not $src) { Die 'Could not find grit.exe inside the archive.' }

    $dest = Join-Path $installDir 'grit.exe'
    # Atomic-ish: write to .new, then swap. Avoids leaving a half-written file
    # if the user happens to be running another grit process.
    Copy-Item -Path $src.FullName -Destination "$dest.new" -Force
    if (Test-Path $dest) {
        try { Move-Item -Path $dest -Destination "$dest.old" -Force } catch {}
    }
    Move-Item -Path "$dest.new" -Destination $dest -Force
    if (Test-Path "$dest.old") { Remove-Item -Path "$dest.old" -Force -ErrorAction SilentlyContinue }
    Ok "grit installed to $dest"
} finally {
    Remove-Item -Path $tmp -Recurse -Force -ErrorAction SilentlyContinue
}

# ─── 6. Ensure install dir is on user PATH ────────────────────────────
$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$onPath   = ($userPath -split ';') -contains $installDir
if (-not $onPath) {
    Info "Adding $installDir to your user PATH"
    $newPath = if ([string]::IsNullOrEmpty($userPath)) { $installDir } else { "$userPath;$installDir" }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    # The current shell still has the old PATH cached; refresh process-local too.
    $env:Path = "$env:Path;$installDir"
    Warn 'Restart your terminal (or open a new one) for the PATH change to fully apply.'
}

# ─── 7. Verify + next steps ───────────────────────────────────────────
Write-Host ''
& (Join-Path $installDir 'grit.exe') version
Write-Host ''
Info "Next: run 'grit new my-app' to scaffold your first project."
Info "Update any time with: 'grit update'"
