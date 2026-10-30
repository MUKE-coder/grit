package selfupdate

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

// Verifying what `grit update` is about to run.
//
// Until v3.450.0 the updater downloaded the release artifact and handed it
// straight to update.Apply with no checksum and no signature, while
// release.yml was already publishing SHA256SUMS, a cosign signature over it,
// and the certificate that signature was made with. Everything needed to
// check the download existed; nothing read it.
//
// What is checked here, in order, and all of it before a single byte is
// written over the running binary:
//
//  1. The certificate chains to Fulcio's root. Pinned below, not fetched,
//     because a trust root you download at verification time is not one.
//  2. Its subject is THIS repository's release workflow, at a tag. A valid
//     Fulcio certificate for some other identity is not evidence about Grit.
//  3. Its issuer extension is GitHub's OIDC provider, so the identity came
//     from an Actions run rather than from somebody's laptop with a browser.
//  4. The signature over SHA256SUMS verifies against that certificate.
//  5. The artifact's SHA-256 appears in those verified sums.
//
// What is NOT checked, said plainly because the alternative is implying more
// than is true: there is no Rekor transparency-log inclusion proof, so there
// is no trusted timestamp for the signature. Full sigstore verification wants
// one. Adding sigstore-go would take this CLI from 56 modules to 129, which
// is a poor trade in the binary whose supply chain this is protecting, and
// the property that actually matters here is held without it: to forge a
// signature you need Fulcio's CA key or control of this repository's release
// workflow, and either of those is already game over.

//go:embed fulcio_roots.pem
var fulcioChainPEM []byte

const (
	// The only identity allowed to sign a Grit release: this repository's
	// release workflow, running on a tag. The "@refs/tags/" suffix matters as
	// much as the path, because a workflow run from a branch could otherwise
	// sign a release.
	wantSANPrefix = "https://github.com/MUKE-coder/grit/.github/workflows/release.yml@refs/tags/"

	// Fulcio records the OIDC issuer in this extension. Checked so the
	// identity is one GitHub Actions minted, not one anybody could obtain
	// from another provider that happens to allow the same subject string.
	wantIssuer = "https://token.actions.githubusercontent.com"
)

// oidFulcioIssuerV2 is 1.3.6.1.4.1.57264.1.8, the issuer as a DER-encoded
// UTF8String. Fulcio also writes 1.3.6.1.4.1.57264.1.1 with the raw value;
// the v2 extension is the one it documents as current.
var oidFulcioIssuerV2 = []int{1, 3, 6, 1, 4, 1, 57264, 1, 8}

// verifiedChecksums parses SHA256SUMS after proving it was signed by this
// repository's release workflow.
//
// sumsPEM is the certificate as cosign writes it, which is base64 over the
// PEM rather than the PEM itself. sig is base64 over a DER ECDSA signature.
func verifiedChecksums(sums, sig, certFile []byte) (map[string]string, error) {
	cert, err := parseSigningCert(certFile)
	if err != nil {
		return nil, err
	}
	if err := checkIdentity(cert); err != nil {
		return nil, err
	}
	if err := checkChain(cert); err != nil {
		return nil, err
	}
	if err := checkSignature(cert, sums, sig); err != nil {
		return nil, err
	}
	return parseChecksums(sums)
}

func parseSigningCert(certFile []byte) (*x509.Certificate, error) {
	raw := certFile
	// cosign --output-certificate writes base64 over the PEM. Accept both, so
	// a release that ever stops doing that does not break the updater.
	if decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(certFile))); err == nil &&
		strings.HasPrefix(string(decoded), "-----BEGIN") {
		raw = decoded
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("the release's signing certificate is not PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing the release's signing certificate: %w", err)
	}
	return cert, nil
}

// checkIdentity refuses a certificate that is valid but not ours.
func checkIdentity(cert *x509.Certificate) error {
	var san string
	for _, u := range cert.URIs {
		if strings.HasPrefix(u.String(), wantSANPrefix) {
			san = u.String()
			break
		}
	}
	if san == "" {
		var got []string
		for _, u := range cert.URIs {
			got = append(got, u.String())
		}
		return fmt.Errorf("the release was signed by %q, not by this repository's "+
			"release workflow on a tag.\n\nThis is what a tampered release looks "+
			"like. Do not install it; open an issue at https://github.com/%s/issues",
			strings.Join(got, ", "), repoSlug)
	}

	for _, ext := range cert.Extensions {
		if !ext.Id.Equal(oidFulcioIssuerV2) {
			continue
		}
		// A DER UTF8String: two bytes of tag and length, then the value. The
		// issuer URL is short, so the length is always single-byte.
		value := string(ext.Value)
		if i := strings.Index(value, "https://"); i >= 0 {
			value = value[i:]
		}
		if value != wantIssuer {
			return fmt.Errorf("the release's signing identity came from %q, not from "+
				"GitHub Actions", value)
		}
		return nil
	}
	return fmt.Errorf("the release's signing certificate does not say which OIDC " +
		"issuer minted it, so the identity cannot be trusted")
}

// checkChain verifies the certificate against the pinned Fulcio root.
//
// At a time inside the certificate's own validity window, deliberately. A
// Fulcio certificate lives ten minutes: it is expired by the time anybody
// runs `grit update`, and refusing on that basis would refuse every real
// release. Proving WHEN the signature was made needs the transparency log,
// which this does not read; what this proves is that Fulcio issued the
// certificate to the identity above.
func checkChain(cert *x509.Certificate) error {
	roots := x509.NewCertPool()
	intermediates := x509.NewCertPool()

	rest := fulcioChainPEM
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		c, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			continue
		}
		// The chain file holds the intermediate first and the root second.
		// Self-signed is the root.
		if c.Subject.String() == c.Issuer.String() {
			roots.AddCert(c)
		} else {
			intermediates.AddCert(c)
		}
	}

	at := cert.NotBefore.Add(time.Minute)
	if at.After(cert.NotAfter) {
		at = cert.NotBefore
	}
	if _, err := cert.Verify(x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
		CurrentTime:   at,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageCodeSigning},
	}); err != nil {
		return fmt.Errorf("the release's signing certificate does not chain to "+
			"Fulcio: %w", err)
	}
	return nil
}

func checkSignature(cert *x509.Certificate, sums, sig []byte) error {
	der, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(sig)))
	if err != nil {
		return fmt.Errorf("the release's signature is not base64: %w", err)
	}
	pub, ok := cert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("the release's signing key is %T, not ECDSA", cert.PublicKey)
	}
	digest := sha256.Sum256(sums)
	if !ecdsa.VerifyASN1(pub, digest[:], der) {
		return fmt.Errorf("the release's checksums do not match their signature.\n\n" +
			"The file has been changed since it was signed. Do not install this " +
			"release")
	}
	return nil
}

// parseChecksums reads the "<hex>  <filename>" lines sha256sum writes.
func parseChecksums(sums []byte) (map[string]string, error) {
	out := map[string]string{}
	for _, line := range strings.Split(string(sums), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		sum, name := fields[0], strings.TrimPrefix(fields[1], "*")
		if len(sum) != 64 {
			continue
		}
		if _, err := hex.DecodeString(sum); err != nil {
			continue
		}
		out[name] = sum
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("the release's SHA256SUMS has no checksums in it")
	}
	return out, nil
}

// matchesChecksum reports whether the downloaded artifact is the one the
// signed checksums describe.
func matchesChecksum(artifact []byte, want string) error {
	got := sha256.Sum256(artifact)
	if hex.EncodeToString(got[:]) != want {
		return fmt.Errorf("the downloaded file is not the one this release signed.\n\n"+
			"  expected  %s\n  got       %s\n\n"+
			"Something changed it between GitHub and here. Do not install it",
			want, hex.EncodeToString(got[:]))
	}
	return nil
}
