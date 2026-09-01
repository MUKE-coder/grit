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
	if strings.Contains(content, "&RecoveryContact{}") {
		return nil
	}
	const marker = "// grit:models"
	if !strings.Contains(content, marker) {
		fmt.Println("  ⚠ models registry marker missing; add &RecoveryContactToken{} and &RecoveryContact{} to Models() by hand")
		return nil
	}
	inject := "\t\t&RecoveryContactToken{},\n\t\t&RecoveryContact{},\n\t\t" + marker
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
