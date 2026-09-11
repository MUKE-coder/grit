package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeSSOFiles brings an upgraded project's enterprise SSO up to date.
//
// Upgrade does not regenerate API code in general. This is the exception
// because the old code let a customer's identity provider assert any email
// address, and linking by email then signed it in as whoever held that
// address, the app's own administrator included. The same release fixes group
// revocation, partial connection updates wiping fields, SAML request matching
// and transient NameIDs.
//
// Manifest-guarded like every upgrade write. The model goes first: the handler
// calls a method it adds, so an edited model that is held back means the
// handler is held back too, rather than written against a model it cannot use.
func writeSSOFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	write := func(path, content string) error {
		if err := writeFile(path, strings.ReplaceAll(content, "{{MODULE}}", opts.Module())); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		return nil
	}

	model := filepath.Join(apiRoot, "internal", "models", "sso.go")
	if err := write(model, apiSSOModelGo()); err != nil {
		return err
	}
	if !fileContains(model, "func (s *SSOConnection) OwnsEmail(") {
		fmt.Println("  ⚠ models/sso.go has been edited, so the SSO security fix was not applied.\n" +
			"    A customer's identity provider can still sign in as any account.\n" +
			"    Run grit upgrade --diff to see what it needs.")
		return nil
	}

	for _, f := range []struct{ path, content string }{
		{filepath.Join(apiRoot, "internal", "services", "sso.go"), apiSSOServiceGo()},
		{filepath.Join(apiRoot, "internal", "handlers", "sso.go"), apiSSOHandlerGo()},
		{filepath.Join(apiRoot, "internal", "handlers", "sso_test.go"), apiSSOTestGo()},
		{filepath.Join(apiRoot, "internal", "services", "saml.go"), apiSAMLServiceGo()},
		{filepath.Join(apiRoot, "internal", "handlers", "saml.go"), apiSAMLHandlerGo()},
		{filepath.Join(apiRoot, "internal", "handlers", "saml_test.go"), apiSAMLTestGo()},
	} {
		if err := write(f.path, f.content); err != nil {
			return err
		}
	}
	return nil
}
