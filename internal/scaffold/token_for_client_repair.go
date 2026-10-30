package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Review 2.10: the refresh token was returned in the JSON body to every
// client, beside the HttpOnly cookie carrying the same value, while the docs
// said a script never sees it. For a browser those cannot both be true.
//
// This gives an older project TokenPair.ForClient, which the repaired sign-in
// handler is about to call. The ordering matters for the same reason
// repairAuthServiceQueries runs before the disclosure repair: that repair
// rewrites Login to what the template says today, and what it says today is
// tokens.ForClient(c). A project whose services/auth.go has no such method
// would be handed a call to a method that is not there and would stop
// compiling on an upgrade that reported success.
//
// The text comes from the template, not from a copy kept here. A copy drifts,
// and the drift is invisible because the two halves are joined by a Go "+"
// that no search for the text can follow.

// repairTokenForClient adds ForClient and its helpers to an older project's
// auth service, and stops the refresh token being serialised when it is empty.
func repairTokenForClient(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "services", "auth.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairTokenForClientSource)
}

func repairTokenForClientSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func (t TokenPair) ForClient(") {
		return src, nil, nil
	}
	if !strings.Contains(src, "type TokenPair struct") {
		return src, nil, []string{"has no TokenPair, so the refresh token cannot be withheld from browsers: add ForClient by hand, or take Grit's copy of internal/services/auth.go"}
	}

	// omitempty, so a browser's response has no refresh_token key at all
	// rather than an empty one. A client checking for the field's presence is
	// then told the truth.
	const oldTag = "\tRefreshToken string `json:\"refresh_token\"`\n"
	out := src
	if strings.Count(out, oldTag) == 1 {
		out = strings.Replace(out, oldTag, "\tRefreshToken string `json:\"refresh_token,omitempty\"`\n", 1)
	} else if !strings.Contains(out, "refresh_token,omitempty") {
		return src, nil, []string{"its TokenPair does not declare RefreshToken the way Grit wrote it, so a browser would still be sent the token: add `json:\"refresh_token,omitempty\"` and ForClient by hand"}
	}

	// strings, for the header compare. gin is already imported: this file sets
	// the auth cookies.
	if !strings.Contains(out, "\n\t\"strings\"\n") {
		const anchor = "import (\n"
		i := strings.Index(out, anchor)
		if i < 0 {
			return src, nil, []string{"has no import block, so ForClient cannot be added to it: take Grit's copy of internal/services/auth.go"}
		}
		out = out[:i+len(anchor)] + "\t\"strings\"\n" + out[i+len(anchor):]
	}

	out = strings.TrimRight(out, "\n") + "\n\n" + repairBlock("api/services/auth.go", "token-for-client")
	return out, []string{"the refresh token is withheld from a browser's response body, which is what the HttpOnly cookie is for, and returned to a native client that asks for it"}, nil
}

// repairTokenForClientHandlers routes every token response through ForClient.
//
// The disclosure repair rewrites Login only when it matches one of the older
// shapes it knows. A project upgraded often enough has the current shape, so
// that repair has nothing to do and the ForClient call never arrives: verified
// on a real project, where services/auth.go gained the method and
// handlers/auth.go kept handing out the raw pair.
//
// So this is a sweep rather than a rewrite. One substitution, in the four
// handlers that answer with a token pair, and it is its own inverse check: a
// file that already calls ForClient is left alone.
func repairTokenForClientHandlers(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, name := range []string{"auth.go", "magic_link.go", "passkey.go", "totp.go"} {
		path := filepath.Join(apiRoot, "internal", "handlers", name)
		if !fileExists(path) {
			continue
		}
		if err := repairSourceFile(root, m, path, repairTokenForClientHandlerSource); err != nil {
			return err
		}
	}
	return nil
}

func repairTokenForClientHandlerSource(src string) (string, []string, []string) {
	const old = `"tokens": tokens,`
	n := strings.Count(src, old)
	if n == 0 {
		return src, nil, nil
	}
	out := strings.ReplaceAll(src, old, `"tokens": tokens.ForClient(c),`)
	return out, []string{"its token responses go through ForClient, so a browser is not handed the refresh token it already has in an HttpOnly cookie"}, nil
}
