package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ensureRecoveryWiring registers the recovery model and mounts its routes in a
// project that already exists.
//
// Writing the files is not enough on an upgrade. Two things live in files the
// upgrade deliberately does not rewrite, because people edit them: the model
// registry inside internal/models/user.go, and the route table. Without both,
// an upgraded project gets the handler and the model source and then reports
// "could not load your security settings", because the table was never created
// and the endpoint was never mounted.
//
// Both files carry markers the generators already inject against, so this uses
// the same mechanism rather than rewriting anything.
func ensureRecoveryWiring(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)

	if err := ensureRecoveryModelRegistered(filepath.Join(apiRoot, "internal", "models", "user.go")); err != nil {
		return err
	}
	return ensureRecoveryRoutes(filepath.Join(apiRoot, "internal", "routes", "routes.go"))
}

func ensureRecoveryModelRegistered(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		// A project without this file is not one we can wire; the scaffold
		// path will have written it correctly anyway.
		return nil
	}
	content := string(body)
	if strings.Contains(content, "&RecoveryContact{}") && strings.Contains(content, "&Passkey{}") {
		return nil
	}
	const marker = "// grit:models"
	if !strings.Contains(content, marker) {
		fmt.Println("  ⚠ models registry marker missing; add &RecoveryContactToken{} and &RecoveryContact{} to Models() by hand")
		return nil
	}
	// Passkeys ride along: their models live in the same registry, and an
	// upgraded project with the handler but not the table answers every
	// passkey request with "no such table: passkeys".
	inject := "\t\t&RecoveryContactToken{},\n\t\t&RecoveryContact{},\n" +
		"\t\t&Passkey{},\n\t\t&WebAuthnSession{},\n\t\t" + marker
	content = strings.Replace(content, marker, inject, 1)
	return os.WriteFile(path, []byte(content), 0o644)
}

func ensureRecoveryRoutes(path string) error {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(body)
	if strings.Contains(content, "recoveryHandler") {
		return nil
	}
	const marker = "// grit:routes:protected"
	if !strings.Contains(content, marker) {
		fmt.Println("  ⚠ protected routes marker missing; mount the recovery routes by hand")
		return nil
	}

	// svc.Mailer is what every other handler in this block is given, so the
	// name is already in scope wherever this marker sits.
	block := strings.Join([]string{
		"// Recovery contacts. Every write takes the account password, because a",
		"// recovery address is a second way in and a live session is exactly",
		"// what somebody on a borrowed laptop already has.",
		"recoveryHandler := handlers.NewRecoveryHandler(db, svc.Mailer)",
		`protected.GET("/auth/security", recoveryHandler.Overview)`,
		`protected.POST("/auth/recovery/email", recoveryHandler.SetEmail)`,
		`protected.POST("/auth/recovery/email/verify", recoveryHandler.VerifyEmail)`,
		`protected.DELETE("/auth/recovery/email", recoveryHandler.ClearEmail)`,
		`protected.POST("/auth/recovery/phone", recoveryHandler.SetPhone)`,
		`protected.POST("/auth/recovery/phone/verify", recoveryHandler.VerifyPhone)`,
		`protected.DELETE("/auth/recovery/phone", recoveryHandler.ClearPhone)`,
		"",
		marker,
	}, "\n\t\t")

	content = strings.Replace(content, marker, block, 1)
	return os.WriteFile(path, []byte(content), 0o644)
}

// ensurePasskeyWiring mounts the passkey routes in a project that already
// exists.
//
// Same reason as the recovery wiring: upgrade writes the handler and the
// service but does not rewrite routes.go, so without this an upgraded project
// gets a 404 on every passkey endpoint and a card that cannot do anything.
//
// The relying party is constructed next to the other handlers rather than at a
// route marker, because the sign-in pair is public and the management routes
// are protected, and one declaration has to serve both.
func ensurePasskeyWiring(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "routes", "routes.go")
	body, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	content := string(body)
	if strings.Contains(content, "passkeyHandler") {
		return nil
	}

	const authAnchor = "\tauthHandler := &handlers.AuthHandler{"
	if !strings.Contains(content, authAnchor) {
		fmt.Println("  ⚠ could not find where handlers are built; mount the passkey routes by hand")
		return nil
	}
	construct := strings.Join([]string{
		"\t// The passkey relying party, built from the origins the frontends run",
		"\t// on. No usable origin means no relying party, and every passkey route",
		"\t// answers 501 rather than panicking: passkeys are optional, a broken",
		"\t// boot is not.",
		"\tpasskeys, passkeyErr := services.NewPasskeys(db, cfg.AppName, cfg.CORSOrigins)",
		"\tif passkeyErr != nil {",
		"\t\tlog.Printf(\"Passkeys disabled: %v\", passkeyErr)",
		"\t\tpasskeys = nil",
		"\t}",
		"\tpasskeyHandler := handlers.NewPasskeyHandler(db, passkeys, authHandler)",
		"",
		authAnchor,
	}, "\n")
	content = strings.Replace(content, authAnchor, construct, 1)

	// The construction reads authHandler, so it has to sit after it, not
	// before. Move it down past the literal that builds authHandler.
	content = moveAfterAuthHandler(content)

	// Anchored on the login route, not on grit:routes:public. That marker
	// sits inside the API-key-guarded catalogue group, where an auth route
	// does not belong and where the auth group is not even in scope.
	loginAnchor := "\t\tauth.POST(\"/login\", authHandler.Login)"
	if strings.Contains(content, loginAnchor) {
		pub := strings.Join([]string{
			loginAnchor,
			"",
			"\t\t// Passkey sign-in. Public because there is no session yet; the",
			"\t\t// server-side challenge is what makes it safe.",
			"\t\tauth.POST(\"/passkeys/login/begin\", passkeyHandler.BeginLogin)",
			"\t\tauth.POST(\"/passkeys/login/finish\", passkeyHandler.FinishLogin)",
		}, "\n")
		content = strings.Replace(content, loginAnchor, pub, 1)
	} else {
		fmt.Println("  ⚠ could not find the login route; mount the passkey sign-in routes by hand")
	}

	const protectedMarker = "// grit:routes:protected"
	if strings.Contains(content, protectedMarker) {
		prot := strings.Join([]string{
			"// Passkey management, on an account you are already signed in to.",
			"protected.GET(\"/auth/passkeys\", passkeyHandler.List)",
			"protected.POST(\"/auth/passkeys/register/begin\", passkeyHandler.BeginRegistration)",
			"protected.POST(\"/auth/passkeys/register/finish\", passkeyHandler.FinishRegistration)",
			"protected.PATCH(\"/auth/passkeys/:id\", passkeyHandler.Rename)",
			"protected.DELETE(\"/auth/passkeys/:id\", passkeyHandler.Delete)",
			"",
			protectedMarker,
		}, "\n\t\t")
		content = strings.Replace(content, protectedMarker, prot, 1)
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

// moveAfterAuthHandler relocates the passkey construction to just after the
// authHandler literal closes.
//
// Injecting before the anchor is the simple thing and produces code that does
// not compile, because the construction takes authHandler as an argument. This
// walks to the end of the composite literal and puts it there instead.
func moveAfterAuthHandler(content string) string {
	start := strings.Index(content, "\t// The passkey relying party")
	if start < 0 {
		return content
	}
	end := strings.Index(content[start:], "\tauthHandler := &handlers.AuthHandler{")
	if end < 0 {
		return content
	}
	block := content[start : start+end]
	rest := content[start+end:]

	// Find the line that closes the authHandler literal: the first "\t}" at the
	// start of a line.
	closeIdx := strings.Index(rest, "\n\t}\n")
	if closeIdx < 0 {
		return content
	}
	insertAt := closeIdx + len("\n\t}\n")
	return content[:start] + rest[:insertAt] + block + rest[insertAt:]
}
