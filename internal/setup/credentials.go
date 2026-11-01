package setup

import (
	"os"
	"path/filepath"
)

// adminCredentials is the login the seeder creates. It never returns the value
// of SEED_ADMIN_PASSWORD: that is a secret the person already knows, and a
// terminal summary is not the place to echo it.
func adminCredentials(root string) string {
	if os.Getenv("SEED_ADMIN_PASSWORD") != "" || envValue(filepath.Join(root, ".env"), "SEED_ADMIN_PASSWORD") != "" {
		return "admin@example.com / the SEED_ADMIN_PASSWORD you set"
	}
	return "admin@example.com / admin123"
}

// localCredentials is the login for one of the generated local tools. The
// values are written to .env by the scaffold and are only meaningful there:
// they are not the admin account, and they are not printed when the tool is
// disabled.
func localCredentials(root, prefix string) string {
	env := filepath.Join(root, ".env")
	if envValue(env, prefix+"_ENABLED") == "false" {
		return ""
	}
	user := envValue(env, prefix+"_USERNAME")
	pass := envValue(env, prefix+"_PASSWORD")
	if user == "" || pass == "" {
		return ""
	}
	return user + " / " + pass
}
