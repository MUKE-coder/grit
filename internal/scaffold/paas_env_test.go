package scaffold

import (
	"strings"
	"testing"
)

// grit#92 item 1: every platform that routes to a container injects PORT, while
// Grit read only APP_PORT, so the deploy went green with nothing answering.
func TestConfigReadsPORTBeforeAPPPORT(t *testing.T) {
	src := apiConfigGo()
	want := `Port:        firstNonEmpty(os.Getenv("PORT"), os.Getenv("APP_PORT"), "8080"),`
	if !strings.Contains(src, want) {
		t.Errorf("config does not prefer PORT:\n%s", section(src, "Port:"))
	}
	if strings.Contains(src, `Port:        getEnv("APP_PORT"`) {
		t.Error("the old APP_PORT-only read is still there")
	}
}

// grit#92 item 4: a managed bucket injects the AWS names, and the s3 driver read
// only the S3_ ones, so attaching one needed the same two values typed twice.
func TestS3ConfigAcceptsTheAWSNames(t *testing.T) {
	src := apiConfigGo()
	for _, want := range []string{
		`Endpoint:  firstNonEmpty(os.Getenv("S3_ENDPOINT"), os.Getenv("AWS_ENDPOINT_URL")),`,
		`Bucket:    firstNonEmpty(os.Getenv("S3_BUCKET"), os.Getenv("AWS_BUCKET"), "uploads"),`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the s3 driver is missing:\n%s", want)
		}
	}
}

// grit#92 item 7: on a platform whose domain is a public suffix, two apps of one
// project are cross-site, the SameSite=Lax cookies are never sent, and sign-in
// returns 200 followed by 401s with nothing logged.
func TestConfigWarnsAboutCrossSiteAuth(t *testing.T) {
	src := apiConfigGo()
	for _, want := range []string{
		"func (c *Config) warnCrossSiteAuth()",
		"cfg.warnCrossSiteAuth()",
		"func registrableDomain(raw string) string",
		`publicsuffix.EffectiveTLDPlusOne(host)`,
		`"golang.org/x/net/publicsuffix"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the cross-site warning is missing: %s", want)
		}
	}
	// localhost and a bare IP say nothing: dev has the API and the frontends on
	// different ports of one host, which is the same site.
	if !strings.Contains(src, `host == "localhost" || net.ParseIP(host) != nil`) {
		t.Error("the check would fire in development")
	}
}

// grit#92 item 5: the API signs an upload for one host and the browser refuses
// it because connect-src names only the host files are read from.
func TestNextCSPAdmitsASeparateUploadHost(t *testing.T) {
	src := nextSecurityHeaders()
	if !strings.Contains(src, "NEXT_PUBLIC_STORAGE_UPLOAD_URL") {
		t.Error("the Next.js CSP has no upload origin")
	}
	if !strings.Contains(src, nextConnectSrcUpload) {
		t.Errorf("connect-src does not include STORAGE_UPLOAD_ORIGIN:\n%s", section(src, "connect-src"))
	}
	// Reads still come from STORAGE_ORIGIN, so the upload host has no business
	// in img-src: it would widen where a stored <img> may be loaded from for
	// nothing.
	if strings.Contains(section(src, "img-src"), "STORAGE_UPLOAD_ORIGIN") {
		t.Error("the upload origin leaked into img-src")
	}
}

func TestViteCSPAdmitsASeparateUploadHost(t *testing.T) {
	src := viteSecurityHeaders()
	if !strings.Contains(src, "VITE_STORAGE_UPLOAD_URL") {
		t.Error("the Vite CSP has no upload origin")
	}
	if !strings.Contains(src, viteConnectSrcNew) {
		t.Errorf("connect-src does not include STORAGE_UPLOAD_ORIGIN:\n%s", section(src, "connect-src"))
	}
}

func TestRepairCSPUploadOriginNext(t *testing.T) {
	fresh := nextSecurityHeaders()
	old := strings.Replace(fresh, nextUploadOriginBlock, "", 1)
	old = strings.Replace(old, nextConnectSrcUpload, nextConnectSrcNew, 1)
	if old == fresh {
		t.Fatal("the fixture is the fresh template, so the repair is not being exercised")
	}

	out, fixed, warn := repairCSPUploadOriginNextSource(old)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if out != fresh {
		t.Errorf("the repaired config differs from a fresh one:\n%s", section(out, "STORAGE_UPLOAD_ORIGIN"))
	}
	if again, fixed, _ := repairCSPUploadOriginNextSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the config again")
	}
	if out, fixed, _ := repairCSPUploadOriginNextSource(fresh); out != fresh || len(fixed) > 0 {
		t.Error("a fresh config still needs the repair")
	}
}

func TestRepairCSPUploadOriginVite(t *testing.T) {
	fresh := viteSecurityHeaders()
	old := strings.Replace(fresh, viteUploadOriginBlock, "", 1)
	old = strings.Replace(old, viteConnectSrcNew, viteConnectSrcOld, 1)
	if old == fresh {
		t.Fatal("the fixture is the fresh template, so the repair is not being exercised")
	}

	out, fixed, warn := repairCSPUploadOriginViteSource(old)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if out != fresh {
		t.Errorf("the repaired config differs from a fresh one:\n%s", section(out, "STORAGE_UPLOAD_ORIGIN"))
	}
	if again, fixed, _ := repairCSPUploadOriginViteSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the config again")
	}
}

func TestRepairCSPUploadOriginLeavesAnEditedPolicyAlone(t *testing.T) {
	src := strings.Replace(nextSecurityHeaders(), nextUploadOriginBlock, "", 1)
	src = strings.Replace(src, nextConnectSrcUpload, `"connect-src 'self' " + MY_ORIGIN + (isDev`, 1)
	out, fixed, warn := repairCSPUploadOriginNextSource(src)
	if out != src || len(fixed) > 0 {
		t.Error("a hand-written policy was edited")
	}
	if len(warn) != 1 || !strings.Contains(warn[0], "connect-src") {
		t.Errorf("the developer was not told what to add by hand: %v", warn)
	}
}

// section returns the lines around the first occurrence of needle, for a
// failure message that shows what is actually there.
func section(src, needle string) string {
	lines := strings.Split(src, "\n")
	for i, line := range lines {
		if !strings.Contains(line, needle) {
			continue
		}
		lo, hi := i-2, i+3
		if lo < 0 {
			lo = 0
		}
		if hi > len(lines) {
			hi = len(lines)
		}
		return strings.Join(lines[lo:hi], "\n")
	}
	return "(" + needle + " not found)"
}
