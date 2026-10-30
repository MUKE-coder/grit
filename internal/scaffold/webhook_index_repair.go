package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairWebhookIndexMigration makes an existing project drop the unique index
// that keyed webhook deliveries on (provider, external_id) alone.
//
// The key is (provider, tenant, external_id) as of v3.446.0, so a gateway
// where every organization connects its own account can receive the same
// reference from two of them (grit#138). AutoMigrate creates the new index,
// because it has a different name, and leaves the old one in place: while that
// one exists the second organization's delivery is refused as a duplicate it
// is not, so the feature is present and does not work.
//
// The drop itself lives in models/webhook_event.go, which upgrades do deliver.
// What does NOT get delivered is Migrate, which lives in models/user.go and is
// written by writeAPIFiles for new projects only, so the call has to be
// injected. That is the trap this exact shape has sprung twice before: a
// migration put somewhere only a new project reads reaches nobody who needs
// it.
func repairWebhookIndexMigration(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "internal", "models", "user.go")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairWebhookIndexMigrationSource)
}

// webhookIndexMigrateHook is the same text the scaffold writes, so an upgraded
// project and a new one read identically.
const webhookIndexMigrateHook = `
	// A webhook delivery is unique per (provider, tenant, event id). The index
	// this replaces keyed on (provider, event id) alone, which refused a
	// second organization's delivery of the same reference. See grit#138.
	if err := MigrateWebhookEventIndexes(db); err != nil {
		return err
	}
`

func repairWebhookIndexMigrationSource(src string) (string, []string, []string) {
	// Not every project has the webhook receiver, and one without it has no
	// model to migrate.
	if !strings.Contains(src, "func Migrate(db *gorm.DB) error {") {
		return src, nil, nil
	}
	if strings.Contains(src, "MigrateWebhookEventIndexes(db)") {
		return src, nil, nil
	}

	// After the roles are seeded, which is the last thing Migrate does and the
	// anchor the append-only triggers already hang off.
	const anchor = `	if err := SeedRoles(db); err != nil {
		return fmt.Errorf("seeding default roles: %w", err)
	}
`
	if !strings.Contains(src, anchor) {
		return src, nil, []string{
			"could not find where to drop the old webhook index in Migrate, so " +
				"deliveries stay unique per (provider, event id): two organizations " +
				"cannot receive the same reference. Add a call to " +
				"MigrateWebhookEventIndexes(db) at the end of Migrate (grit#138)"}
	}

	return strings.Replace(src, anchor, anchor+webhookIndexMigrateHook, 1),
		[]string{"a webhook delivery is unique per (provider, tenant, event id), " +
			"so two organizations can receive the same provider reference"},
		nil
}
