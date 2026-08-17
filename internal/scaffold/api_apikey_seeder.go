package scaffold

// apiAPIKeySeederGo emits internal/database/api_keys_seeder.go.
//
// A new project gets two keys the moment it is seeded: a publishable one for
// the client apps and a secret one for server-to-server work. Both are printed
// with the rest of the seed output.
//
// Seeded rather than created by the CLI because the CLI does not have a
// database connection at `grit new` time, and because a key belongs to a user,
// which the seeder has just created. It also means `grit seed` on a fresh
// database always leaves you with working keys, including in CI.
func apiAPIKeySeederGo() string {
	return `package database

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/services"
)

// SeedAPIKeys issues the two keys every project starts with.
//
// Idempotent by name: running the seeder twice does not mint a second pair,
// because the second pair would be just as valid and you would have no way to
// know which one your app is using.
func SeedAPIKeys(db *gorm.DB) error {
	var owner models.User
	if err := db.Where("role = ?", "ADMIN").Order("created_at asc").First(&owner).Error; err != nil {
		log.Println("Skipping API keys: no admin user to own them")
		return nil
	}

	publishable, err := ensureKey(db, models.APIKey{
		UserID: owner.ID,
		Name:   "Client apps (publishable)",
		Kind:   models.KindPublishable,
		// No endpoint allowlist. The kind alone already restricts it to
		// public routes, and narrowing further on a fresh project would mean
		// editing this list every time somebody marks a resource --public.
		// Add one here when you want a specific client narrowed.
	})
	if err != nil {
		return err
	}

	secret, err := ensureKey(db, models.APIKey{
		UserID: owner.ID,
		Name:   "Server to server (secret)",
		Kind:   models.KindSecret,
	})
	if err != nil {
		return err
	}

	if publishable != "" || secret != "" {
		log.Println("================================================================")
		log.Println("API keys")
		if publishable != "" {
			log.Printf("  Publishable  %s", publishable)
			log.Println("               Safe in a browser or a mobile app. Reaches")
			log.Println("               endpoints marked public, and nothing else.")
		}
		if secret != "" {
			log.Printf("  Secret       %s", secret)
			log.Println("               Server side only. Shown once, right now.")
		}
		log.Println("  Manage both in the admin at Settings, API Keys.")
		log.Println("================================================================")
		writeClientEnv(publishable)
	}

	return nil
}

// ensureKey creates a key if one with that name does not exist. Returns the
// token when it minted one, and "" when it did not.
func ensureKey(db *gorm.DB, want models.APIKey) (string, error) {
	var existing models.APIKey
	err := db.Where("user_id = ? AND name = ?", want.UserID, want.Name).First(&existing).Error
	if err == nil {
		return "", nil
	}
	if err != gorm.ErrRecordNotFound {
		return "", err
	}

	issued, err := services.GenerateAPIKey(db, services.KeyOptions{
		UserID:    want.UserID,
		Name:      want.Name,
		Kind:      want.Kind,
		Endpoints: want.Endpoints,
		Origins:   want.Origins,
	})
	if err != nil {
		return "", err
	}
	return issued.Token, nil
}

// writeClientEnv drops the publishable key into the web app's local env, so a
// fresh project's storefront can call the API without anyone copying anything.
//
// Only ever written when absent. Overwriting somebody's env file because a
// seeder ran is the kind of helpfulness that loses an afternoon.
func writeClientEnv(publishable string) {
	if publishable == "" {
		return
	}
	for _, rel := range []string{
		filepath.Join("..", "web", ".env.local"),
		filepath.Join("..", "admin", ".env.local"),
	} {
		if _, err := os.Stat(filepath.Dir(rel)); err != nil {
			continue // that app is not part of this project
		}
		if _, err := os.Stat(rel); err == nil {
			continue // already there
		}
		body := "# Publishable API key, written by the seeder.\n" +
			"# Safe to ship to a browser: it reaches public endpoints only.\n" +
			"NEXT_PUBLIC_API_KEY=" + publishable + "\n"
		if err := os.WriteFile(rel, []byte(body), 0o644); err == nil {
			log.Printf("  Wrote %s", rel)
		}
	}
}
`
}
