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
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"

	"{{MODULE}}/internal/config"
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
	err := db.Where("role = ?", "ADMIN").Order("created_at asc").First(&owner).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Skipping API keys: no admin user to own them")
		return nil
	}
	// Anything else is the database failing, not a project without an admin.
	if err != nil {
		return err
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
	if !errors.Is(err, gorm.ErrRecordNotFound) {
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

// apiOrigin is where the browser should call this API.
// apiPort is the port the API listens on, as a string. The Expo app needs the
// port without a host: it derives the host from Metro, because a phone cannot
// reach localhost.
func apiPort() string {
	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Port == "" {
		return "8080"
	}
	return cfg.Port
}

func apiOrigin() string {
	cfg, err := config.Load()
	if err != nil || cfg == nil || cfg.Port == "" {
		return "http://localhost:8080"
	}
	return "http://localhost:" + cfg.Port
}

// storageOrigin is the browser-facing origin of stored files.
//
// Derived from whichever storage this project is configured for rather than
// assumed. The default MinIO port is only correct until somebody moves it, and
// in production it is never correct. Getting it wrong does not fail loudly:
// the upload is refused by the Content-Security-Policy, and the only trace is
// a console message nobody is watching.
func storageOrigin() string {
	cfg, err := config.Load()
	if err != nil || cfg == nil {
		return "http://localhost:9002"
	}
	for _, candidate := range []string{cfg.Storage.PublicURL, cfg.Storage.Endpoint} {
		if candidate == "" {
			continue
		}
		value := strings.TrimSuffix(candidate, "/")
		if !strings.Contains(value, "://") {
			value = "http://" + value
		}
		// A public URL usually carries the bucket path. The CSP wants an
		// origin, so everything after the host is dropped.
		if i := strings.Index(value, "://"); i >= 0 {
			if j := strings.Index(value[i+3:], "/"); j >= 0 {
				value = value[:i+3+j]
			}
		}
		return value
	}
	return "http://localhost:9002"
}

// writeClientEnv drops the publishable key into the web app's local env, so a
// fresh project's storefront can call the API without anyone copying anything.
//
// Only ever written when absent. Overwriting somebody's env file because a
// seeder ran is the kind of helpfulness that loses an afternoon.
//
// It does say so when the file holds a key this database will not accept,
// which is what grit migrate --fresh leaves behind: the keys table is dropped,
// the seeder mints a new one, and the file still names the old one. Every
// public request then answers INVALID_API_KEY with nothing to connect that to
// a file the seeder decided not to touch.
func writeClientEnv(publishable string) {
	if publishable == "" {
		return
	}
	// Every frontend this project might have, relative to where the seeder runs.
	//
	// In a monorepo that is apps/api, so the other apps are one level up.
	// In a single project the Go module is the root, so its one frontend is
	// frontend/ right here. A single project got neither path and so got no
	// .env.local at all: its frontend fell back to localhost:8080 whatever the
	// API was actually listening on, and had no publishable key, so every public
	// endpoint answered INVALID_API_KEY. A directory that is not part of this
	// project is skipped below, so listing all of them is safe in all of them.
	for _, rel := range []string{
		filepath.Join("..", "web", ".env.local"),
		filepath.Join("..", "admin", ".env.local"),
		filepath.Join("..", "expo", ".env.local"),
		// A single project's frontend is the project root, and the seeder runs from
		// api/ beside it. The old flat layout kept it in frontend/ from the root,
		// which is the second path: an upgraded project still gets its file.
		filepath.Join("..", ".env.local"),
		filepath.Join("frontend", ".env.local"),
	} {
		dir := filepath.Dir(rel)
		// A frontend, not merely a directory. "../" is the project root for a
		// single project and the apps/ folder for a monorepo, and only one of
		// those is somewhere an .env.local belongs: a package.json is what tells
		// them apart. Before this the test was "does the directory exist", which
		// "../" always satisfies.
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
			continue // that app is not part of this project
		}
		prefix := clientEnvPrefix(dir)
		if existing, err := os.ReadFile(rel); err == nil {
			if stale := storedKey(string(existing)); stale != "" && stale != publishable {
				log.Printf("  ! %s still names a key this database does not have.", rel)
				log.Printf("    Public requests will answer INVALID_API_KEY until you set")
				log.Printf("    %sAPI_KEY=%s", prefix, publishable)
			}
			continue // never overwritten
		}
		if err := os.WriteFile(rel, []byte(clientEnvBody(prefix, publishable)), 0o644); err == nil {
			log.Printf("  Wrote %s", rel)
		}
	}
}

// clientEnvPrefix is the prefix the app in dir reads its public variables
// under. Each bundler exposes only its own, so writing the wrong one produces
// a file that is read by nothing and explains nothing.
func clientEnvPrefix(dir string) string {
	for _, name := range []string{"vite.config.ts", "vite.config.js", "vite.config.mjs"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return "VITE_"
		}
	}
	// Metro is the Expo bundler and nothing else here has its config. A file
	// test rather than a search for a key inside app.json: the same shape as
	// the Vite test above, and it cannot be fooled by a string in a comment.
	for _, name := range []string{"metro.config.js", "metro.config.cjs"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			return "EXPO_PUBLIC_"
		}
	}
	return "NEXT_PUBLIC_"
}

// clientEnvBody is the file, in the names that app reads.
func clientEnvBody(prefix, publishable string) string {
	body := "# Written by the seeder. Only ever created, never overwritten." + "\n" +
		"\n" +
		"# Publishable API key. Safe to ship to a client: it reaches public" + "\n" +
		"# endpoints only." + "\n" +
		prefix + "API_KEY=" + publishable + "\n" +
		"\n"

	if prefix == "EXPO_PUBLIC_" {
		// A port, not a URL, because a phone cannot reach localhost: that
		// address is the phone. The app derives the dev machine's LAN address
		// from the host Metro is already served on and needs only the port to
		// go with it, so naming the port here moves the app with APP_PORT and
		// keeps working on a real device. Set EXPO_PUBLIC_API_URL instead to
		// point at a deployed backend, which wins over both.
		return body +
			"# The port the API listens on, to go with the LAN address the app" + "\n" +
			"# derives from Metro. Set EXPO_PUBLIC_API_URL instead to point at a" + "\n" +
			"# deployed backend." + "\n" +
			prefix + "API_PORT=" + apiPort() + "\n" +
			"\n" +
			"# Origin of stored files, for images the app displays." + "\n" +
			prefix + "STORAGE_URL=" + storageOrigin() + "\n"
	}

	return body +
		"# Where this app calls the API." + "\n" +
		prefix + "API_URL=" + apiOrigin() + "\n" +
		"\n" +
		"# Browser-facing origin of stored files." + "\n" +
		"#" + "\n" +
		"# Both of these end up in the Content-Security-Policy. Uploads are" + "\n" +
		"# presigned PUTs made straight from the browser to object storage," + "\n" +
		"# so an origin missing from connect-src is an upload the browser" + "\n" +
		"# refuses, reporting it only as a CSP violation in the console." + "\n" +
		"#" + "\n" +
		"# In production set this to your S3, R2 or CDN origin." + "\n" +
		prefix + "STORAGE_URL=" + storageOrigin() + "\n"
}

// storedKey pulls NEXT_PUBLIC_API_KEY out of an env file, or "" when it has
// none. Written by hand rather than with a parser: one key, one line, and a
// dependency for this would be absurd.
func storedKey(body string) string {
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		for _, prefix := range []string{"NEXT_PUBLIC_", "VITE_", "EXPO_PUBLIC_"} {
			if value, ok := strings.CutPrefix(line, prefix+"API_KEY="); ok {
				return strings.TrimSpace(value)
			}
		}
	}
	return ""
}
`
}
