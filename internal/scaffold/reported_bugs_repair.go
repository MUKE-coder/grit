package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Repairs for bugs reported against projects already in the field. Each is
// written so a project that already has the fix, or has rewritten the file
// itself, is left alone.

// repairReportedBugs runs them all. Called from Upgrade.
func repairReportedBugs(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	apiRoot := opts.APIRoot(root)

	for _, r := range []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "middleware", "cors.go"), repairCORSOptimisticLockingSource},
		{filepath.Join(apiRoot, "internal", "routes", "routes.go"), repairAliasStatusSource},
		{filepath.Join(apiRoot, "internal", "mail", "mime.go"), repairCallerMessageIDSource},
		{filepath.Join(apiRoot, "internal", "sequence", "sequence.go"), repairSequenceSaveSource},
	} {
		if !fileExists(r.path) {
			continue
		}
		if err := repairSourceFile(root, m, r.path, r.fn); err != nil {
			return err
		}
	}
	return nil
}

// grit#124. Optimistic locking could not be used by any browser on another
// origin: If-Match was stripped during preflight, so the request never
// arrived, and ETag was not exposed, so a client could not read the version
// it had to send back. The existing CORS repair converts a string literal to
// a list and stops there, so a project already on the list shape needs this.
func repairCORSOptimisticLockingSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "var allowedRequestHeaders = []string{") {
		// Still the literal, or rewritten by hand. repairCORSHeaderList deals
		// with the first; the second is the developer's business.
		return src, nil, nil
	}

	var changes []string
	out := src

	if !strings.Contains(out, `"If-Match"`) {
		const anchor = "var allowedRequestHeaders = []string{\n"
		out = strings.Replace(out, anchor, anchor+
			"\t// Optimistic locking. internal/concurrency reads the version the\n"+
			"\t// client last saw from If-Match and answers 409 rather than\n"+
			"\t// overwriting a newer save. Without it here the preflight strips\n"+
			"\t// the header and the write is never sent. See grit#124.\n"+
			"\t\"If-Match\",\n", 1)
		changes = append(changes, "If-Match is allowed, so a browser on another origin can use optimistic locking")
	}

	if !strings.Contains(out, "exposedResponseHeaders") {
		const anchor = "// allowedRequestHeadersValue is the joined header value, built once."
		if !strings.Contains(out, anchor) {
			return out, changes, []string{
				"no place found for the exposed-header list, so ETag is still unreadable " +
					"cross-origin and If-Match has no value to carry (grit#124)"}
		}
		out = strings.Replace(out, anchor, corsExposeRegistrySource+anchor, 1)
		changes = append(changes, "ETag, Deprecation and Link are exposed, so a cross-origin client can read them")
	}

	if !strings.Contains(out, "Access-Control-Expose-Headers") {
		const anchor = `c.Header("Access-Control-Allow-Headers", allowedRequestHeadersValue)`
		if !strings.Contains(out, anchor) {
			return out, changes, []string{
				"the allow-headers line was not found, so Access-Control-Expose-Headers " +
					"was not added (grit#124)"}
		}
		out = strings.Replace(out, anchor,
			anchor+"\n\t\tc.Header(\"Access-Control-Expose-Headers\", exposedResponseHeadersValue)", 1)
	}

	return out, changes, nil
}

// The same text the scaffold writes, so an upgraded project and a new one
// read identically.
const corsExposeRegistrySource = `// exposedResponseHeaders is what cross-origin JavaScript may READ back.
//
// A response header not named here is invisible to fetch() on another
// origin, whatever the server sent. ETag carries the version a client puts
// in If-Match on its next write, so without this the optimistic locking in
// internal/concurrency is only half usable from a browser: it can send the
// header and has no way to learn the value. See grit#124.
var exposedResponseHeaders = []string{
	"ETag",
	"Deprecation",
	"Link",
}

// ExposeResponseHeader adds a header to that list, for the same reason
// AllowRequestHeader exists.
func ExposeResponseHeader(name string) {
	for _, h := range exposedResponseHeaders {
		if strings.EqualFold(h, name) {
			return
		}
	}
	exposedResponseHeaders = append(exposedResponseHeaders, name)
	exposedResponseHeadersValue = strings.Join(exposedResponseHeaders, ", ")
}

// exposedResponseHeadersValue is the joined header value, built once.
var exposedResponseHeadersValue = strings.Join(exposedResponseHeaders, ", ")

`

// grit#128. Every generated /export answered 404 through the unversioned
// /api path, with the whole file in the body: the alias is a NoRoute handler,
// gin has already set 404 on the way in, and HandleContext does not reset the
// writer, so a handler that streams bytes flushes the leftover status.
func repairAliasStatusSource(src string) (string, []string, []string) {
	i := strings.Index(src, "func mountLegacyAPIAlias")
	if i < 0 {
		return src, nil, nil
	}
	const anchor = "\t\t\tr.HandleContext(c)"
	if !strings.Contains(src[i:], anchor) {
		return src, nil, []string{
			"mountLegacyAPIAlias does not re-dispatch the way it used to, so the " +
				"pending 404 was not cleared: downloads through /api may still " +
				"answer 404 with the file in the body (grit#128)"}
	}
	if strings.Contains(src[i:], "c.Status(http.StatusOK)") {
		return src, nil, nil
	}

	head, tail := src[:i], src[i:]
	tail = strings.Replace(tail, anchor,
		"\t\t\t// Clear the 404 gin set on the way into NoRoute. A handler that\n"+
			"\t\t\t// streams bytes without calling WriteHeader flushes that status\n"+
			"\t\t\t// instead of its own: every generated /export answered 404 with\n"+
			"\t\t\t// the whole file in the body. See grit#128.\n"+
			"\t\t\tc.Status(http.StatusOK)\n"+anchor, 1)
	out := head + tail

	withImport, ok := addImportGroup(out, "net/http")
	if !ok {
		return out, nil, []string{
			`the alias fix needs "net/http" and this file has no import block to add it to: add the import by hand`}
	}
	return withImport, []string{
		"the unversioned /api alias clears the pending 404, so downloads through it answer 200"}, nil
}

// grit#122. A Message-ID set by the caller was dropped by the SMTP driver and
// kept by the API drivers, so the same message got a different id depending
// on which one sent it.
func repairCallerMessageIDSource(src string) (string, []string, []string) {
	if strings.Contains(src, "resolveMessageID(") {
		return src, nil, nil
	}
	const anchor = "\tid, err := messageID(from.Address)"
	if !strings.Contains(src, anchor) {
		return src, nil, nil
	}
	out := strings.Replace(src, anchor, "\tid, err := resolveMessageID(m.Headers, from.Address)", 1)

	const funcAnchor = "func messageID(from string) (string, error) {"
	if !strings.Contains(out, funcAnchor) {
		return src, nil, []string{
			"messageID is not declared here, so a caller-set Message-ID is still " +
				"dropped by the SMTP driver (grit#122)"}
	}
	out = strings.Replace(out, funcAnchor, resolveMessageIDSource+funcAnchor, 1)

	return out, []string{
		"a caller-set Message-ID is kept, so a reply can be matched to what was sent"}, nil
}

const resolveMessageIDSource = `// resolveMessageID keeps the caller's Message-ID when they set one, and
// generates it otherwise.
//
// A Message-ID is the only thing a reply carries back, in In-Reply-To and
// References, so an app that logs what it sent needs to know the id that went
// out. This used to always generate: Message-Id is in reservedHeaders, so the
// SMTP driver replaced the caller's silently while the API drivers passed
// theirs through. See grit#122.
//
// A caller value is used only when it is <local@domain> with no CR or LF. A
// header value containing either is header injection.
func resolveMessageID(headers map[string]string, from string) (string, error) {
	for name, v := range headers {
		if !strings.EqualFold(name, "Message-ID") {
			continue
		}
		v = strings.TrimSpace(v)
		if validMessageID(v) {
			return v, nil
		}
		break
	}
	return messageID(from)
}

// validMessageID reports whether v is a well formed <local@domain>.
func validMessageID(v string) bool {
	if len(v) < 5 || !strings.HasPrefix(v, "<") || !strings.HasSuffix(v, ">") {
		return false
	}
	if strings.ContainsAny(v, "\r\n") {
		return false
	}
	inner := v[1 : len(v)-1]
	at := strings.LastIndex(inner, "@")
	if at <= 0 || at >= len(inner)-1 {
		return false
	}
	return !strings.ContainsAny(inner, "<> ")
}

`

// grit#132. A sequence with Reset "never" issued one number and then failed
// forever: its bucket is "", the zero value of the string half of the
// counter's primary key, and GORM's Save treats a zero primary key as a new
// row, so the second call inserted a duplicate.
func repairSequenceSaveSource(src string) (string, []string, []string) {
	const anchor = "\t\treturn tx.Save(&c).Error"
	if !strings.Contains(src, anchor) {
		return src, nil, nil
	}
	out := strings.Replace(src, anchor,
		"\t\t// An explicit UPDATE, not Save. With Reset \"never\" the bucket is \"\",\n"+
			"\t\t// which is the zero value of the string half of the primary key, and\n"+
			"\t\t// Save treats a zero primary key as a new row: it INSERTed a second\n"+
			"\t\t// (\"name\", \"\") and every number after the first failed on the unique\n"+
			"\t\t// constraint. See grit#132.\n"+
			"\t\treturn tx.Model(&Counter{}).\n"+
			"\t\t\tWhere(\"name = ? AND bucket = ?\", cfg.Name, bucket).\n"+
			"\t\t\tUpdates(map[string]any{\n"+
			"\t\t\t\t\"next_value\": c.NextValue,\n"+
			"\t\t\t\t\"updated_at\": time.Now(),\n"+
			"\t\t\t}).Error", 1)

	return out, []string{
		"a sequence that never resets issues more than one number"}, nil
}
