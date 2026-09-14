package scaffold

import "strings"

// A release workflow for the generated project.
//
// Grit signs its own releases — SBOM, keyless cosign signature, provenance
// attestation — and gave projects it scaffolds none of that. The asymmetry was
// the whole problem: the framework met a supply-chain bar its output could not.
//
// Everything here works with no configuration. Keyless signing uses the
// workflow's GitHub OIDC identity, so there is no key to store, rotate or leak;
// the certificate binds the signature to this repository and this workflow
// file, and Rekor logs it publicly. A signature therefore cannot be produced
// outside a real run.
//
// Authenticode and notarization are the exception, because they need
// credentials only the project owner can buy. Those steps are present but
// guarded on their secrets, so a project without them releases unsigned
// binaries rather than failing — and the workflow says so in its output rather
// than skipping quietly.

func releaseCIYAML(opts Options) string {
	l := ciLayoutFor(opts)
	base := `name: Release

# Tag-triggered. Push v1.2.3 and this builds, signs and publishes.
on:
  push:
    tags:
      - "v*"

# Read-only by default. Each job that publishes asks for what it needs.
permissions:
  contents: read

jobs:
  release:
    name: Build & Release
    runs-on: ubuntu-latest
    # Signing and provenance need an OIDC token; without id-token the cosign
    # step fails with a confusing OIDC error rather than an obvious one.
    permissions:
      contents: write
      id-token: write
      attestations: write
    steps:
      - uses: actions/checkout@11d5960a326750d5838078e36cf38b85af677262 # v4

      - uses: actions/setup-go@40f1582b2485089dde7abd97c1529aa768e1baff # v5
        with:
          go-version-file: {{API_DIR}}/go.mod
          cache-dependency-path: {{API_DIR}}/go.sum

{{FRONTEND_BUILD}}      # Tests gate the release. A tag that ships a failing build is worse than
      # a tag that never ships.
      - name: Test
        working-directory: {{API_DIR}}
        run: go test ./... -race

      - name: Build API binaries
        working-directory: {{API_DIR}}
        run: |
          mkdir -p "$GITHUB_WORKSPACE/dist"
          VERSION="${GITHUB_REF_NAME}"
          for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
            os="${target%/*}"; arch="${target#*/}"
            out="$GITHUB_WORKSPACE/dist/api-${os}-${arch}"
            # CGO off so the binary runs on any distro, including scratch and
            # distroless images, without a matching libc.
            CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build \
              -trimpath -ldflags "-s -w -X main.version=${VERSION}" \
              -o "$out" ./cmd/server
          done

      - name: Checksums
        working-directory: dist
        run: sha256sum api-* > SHA256SUMS

      # An SBOM is what turns "are we affected by this CVE?" from an
      # afternoon into a grep.
      # Syft from Anchore's action, pinned like every other step, rather than a
      # shell script fetched from a URL at release time.
      - name: Install Syft
        id: syft
        uses: anchore/sbom-action/download-syft@e22c389904149dbc22b58101806040fa8d37a610 # v0

      - name: Generate SBOM
        run: |
          "${{ steps.syft.outputs.cmd }}" dir:. \
            --exclude './**/node_modules/**' \
            --exclude './**/.next/**' \
            -o spdx-json > dist/sbom.spdx.json
          # An empty SBOM passes a scanner and reports nothing, which is worse
          # than having none at all.
          python3 -c "import json,sys;d=json.load(open('dist/sbom.spdx.json'));n=len(d.get('packages',[]));print(f'SBOM packages: {n}');sys.exit(1 if n==0 else 0)"

      - name: Install cosign
        uses: sigstore/cosign-installer@398d4b0eeef1380460a10c8013a76f728fb906ac # v3

      - name: Sign the checksums
        working-directory: dist
        run: |
          # Keyless: no private key exists to steal. The short-lived certificate
          # binds this signature to this repo and workflow, and the whole thing
          # is recorded in the public Rekor transparency log.
          cosign sign-blob --yes \
            --output-certificate SHA256SUMS.pem \
            --output-signature   SHA256SUMS.sig \
            SHA256SUMS

      - name: Attest build provenance
        uses: actions/attest-build-provenance@e8998f949152b193b063cb0ec769d69d929409be # v2
        with:
          subject-path: dist/api-*

      - name: Publish
        uses: softprops/action-gh-release@3bb12739c298aeb8a4eeaf626c5b8d85266b0e65 # v2
        with:
          generate_release_notes: true
          files: |
            dist/api-*
            dist/SHA256SUMS
            dist/SHA256SUMS.pem
            dist/SHA256SUMS.sig
            dist/sbom.spdx.json

`
	frontend, desktop := "", ""
	if l.single {
		// The server embeds the built frontend, so it is built first.
		frontend = ciPnpmSteps(l) + "      - name: Build the frontend\n        working-directory: frontend\n        run: pnpm build\n\n"
		frontend = strings.Replace(frontend, "        run: if [ -f pnpm-lock.yaml ]", "        working-directory: frontend\n        run: if [ -f pnpm-lock.yaml ]", 1)
	}
	if l.desktop {
		desktop = releaseDesktopJob
	}
	return ciReplace(l, strings.Replace(base, "{{FRONTEND_BUILD}}", frontend, 1)+desktop)
}

// releaseDesktopJob builds the desktop installers, for a project that has a
// desktop app. It was once in every project's workflow behind a job-level
// hashFiles condition, which GitHub does not allow there, so the whole workflow
// was rejected on the first tag.
const releaseDesktopJob = `
  # Desktop installers are a separate job because they need Windows and macOS
  # runners, and because signing them needs credentials this workflow cannot
  # invent.
  desktop:
    name: Desktop installer (${{ matrix.os }})
    strategy:
      matrix:
        os: [windows-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@11d5960a326750d5838078e36cf38b85af677262 # v4
      - uses: actions/setup-go@40f1582b2485089dde7abd97c1529aa768e1baff # v5
        with:
          go-version-file: {{API_DIR}}/go.mod

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@v2.9.2

      - name: Build
        working-directory: apps/desktop
        run: wails build -trimpath -ldflags "-s -w"

      # ── Windows Authenticode ────────────────────────────────────────
      # Skipped, loudly, when the certificate is not configured. An unsigned
      # Windows installer shows a SmartScreen warning that most users read as
      # "this is malware", so this is worth setting up before a public launch.
      #
      #   WINDOWS_CERT_BASE64  — the .pfx, base64-encoded
      #   WINDOWS_CERT_PASSWORD
      - name: Sign (Windows)
        if: matrix.os == 'windows-latest'
        shell: pwsh
        env:
          CERT: ${{ secrets.WINDOWS_CERT_BASE64 }}
          CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
        run: |
          if (-not $env:CERT) {
            Write-Host "::warning::WINDOWS_CERT_BASE64 is not set - shipping an UNSIGNED installer. Users will see a SmartScreen warning."
            exit 0
          }
          [IO.File]::WriteAllBytes("cert.pfx", [Convert]::FromBase64String($env:CERT))
          $signtool = Get-ChildItem "C:\Program Files (x86)\Windows Kits\10\bin\*\x64\signtool.exe" |
            Sort-Object FullName -Descending | Select-Object -First 1
          $exe = Get-ChildItem apps/desktop/build/bin/*.exe
          & $signtool.FullName sign /f cert.pfx /p $env:CERT_PASSWORD /tr http://timestamp.digicert.com /td sha256 /fd sha256 $exe
          Remove-Item cert.pfx

      # ── macOS signing + notarization ────────────────────────────────
      # Without notarization, Gatekeeper refuses to open the app at all on a
      # machine that did not build it.
      #
      #   MACOS_CERT_BASE64 / MACOS_CERT_PASSWORD  — Developer ID Application
      #   MACOS_SIGN_IDENTITY                      — "Developer ID Application: You (TEAMID)"
      #   MACOS_NOTARY_APPLE_ID / MACOS_NOTARY_PASSWORD / MACOS_NOTARY_TEAM_ID
      - name: Sign and notarize (macOS)
        if: matrix.os == 'macos-latest'
        env:
          CERT: ${{ secrets.MACOS_CERT_BASE64 }}
          CERT_PASSWORD: ${{ secrets.MACOS_CERT_PASSWORD }}
          SIGN_IDENTITY: ${{ secrets.MACOS_SIGN_IDENTITY }}
          NOTARY_APPLE_ID: ${{ secrets.MACOS_NOTARY_APPLE_ID }}
          NOTARY_PASSWORD: ${{ secrets.MACOS_NOTARY_PASSWORD }}
          NOTARY_TEAM_ID: ${{ secrets.MACOS_NOTARY_TEAM_ID }}
        run: |
          if [ -z "$CERT" ]; then
            echo "::warning::MACOS_CERT_BASE64 is not set - shipping an UNSIGNED app. Gatekeeper will refuse to open it."
            exit 0
          fi
          # A throwaway keychain, so the runner's default keychain is untouched
          # and the certificate disappears with the VM.
          echo "$CERT" | base64 --decode > cert.p12
          security create-keychain -p runner build.keychain
          security default-keychain -s build.keychain
          security unlock-keychain -p runner build.keychain
          security import cert.p12 -k build.keychain -P "$CERT_PASSWORD" -T /usr/bin/codesign
          security set-key-partition-list -S apple-tool:,apple:,codesign: -s -k runner build.keychain

          APP=$(find apps/desktop/build/bin -name "*.app" -maxdepth 1 | head -1)
          # --options runtime enables the hardened runtime, which notarization
          # requires; without it the submission is rejected.
          codesign --deep --force --options runtime --timestamp \
            --sign "$SIGN_IDENTITY" "$APP"

          ditto -c -k --keepParent "$APP" notarize.zip
          xcrun notarytool submit notarize.zip \
            --apple-id "$NOTARY_APPLE_ID" \
            --password "$NOTARY_PASSWORD" \
            --team-id "$NOTARY_TEAM_ID" \
            --wait
          # Staple, so the app opens on a machine with no network.
          xcrun stapler staple "$APP"
          rm -f cert.p12 notarize.zip

      - name: Publish installers
        uses: softprops/action-gh-release@3bb12739c298aeb8a4eeaf626c5b8d85266b0e65 # v2
        with:
          files: |
            apps/desktop/build/bin/*
`
