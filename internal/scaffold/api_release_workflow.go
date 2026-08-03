package scaffold

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

func releaseCIYAML() string {
	return `name: Release

# Tag-triggered. Push v1.2.3 and this builds, signs and publishes.
on:
  push:
    tags:
      - "v*"

permissions:
  contents: write
  # Required for keyless signing and provenance. Without id-token the cosign
  # step fails with a confusing OIDC error rather than an obvious one.
  id-token: write
  attestations: write

jobs:
  release:
    name: Build & Release
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: apps/api/go.mod
          cache-dependency-path: apps/api/go.sum

      # Tests gate the release. A tag that ships a failing build is worse than
      # a tag that never ships.
      - name: Test
        working-directory: apps/api
        run: go test ./... -race

      - name: Build API binaries
        working-directory: apps/api
        run: |
          mkdir -p ../../dist
          VERSION="${GITHUB_REF_NAME}"
          for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64; do
            os="${target%/*}"; arch="${target#*/}"
            out="../../dist/api-${os}-${arch}"
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
      - name: Generate SBOM
        run: |
          curl -sSfL https://get.anchore.io/syft | sh -s -- -b /usr/local/bin
          syft dir:. \
            --exclude './**/node_modules/**' \
            --exclude './**/.next/**' \
            -o spdx-json > dist/sbom.spdx.json
          # An empty SBOM passes a scanner and reports nothing, which is worse
          # than having none at all.
          python3 -c "import json,sys;d=json.load(open('dist/sbom.spdx.json'));n=len(d.get('packages',[]));print(f'SBOM packages: {n}');sys.exit(1 if n==0 else 0)"

      - name: Install cosign
        uses: sigstore/cosign-installer@v3

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
        uses: actions/attest-build-provenance@v2
        with:
          subject-path: dist/api-*

      - name: Publish
        uses: softprops/action-gh-release@v2
        with:
          generate_release_notes: true
          files: |
            dist/api-*
            dist/SHA256SUMS
            dist/SHA256SUMS.pem
            dist/SHA256SUMS.sig
            dist/sbom.spdx.json

  # Desktop installers are a separate job because they need Windows and macOS
  # runners, and because signing them needs credentials this workflow cannot
  # invent. Delete this job if you do not ship a desktop app.
  desktop:
    name: Desktop installer (${{ matrix.os }})
    # Only runs when the repo actually has a desktop app.
    if: hashFiles('apps/desktop/wails.json') != ''
    strategy:
      matrix:
        os: [windows-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: apps/api/go.mod

      - name: Install Wails
        run: go install github.com/wailsapp/wails/v2/cmd/wails@latest

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
        uses: softprops/action-gh-release@v2
        with:
          files: |
            apps/desktop/build/bin/*
`
}
