package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// ─── the claims table on a project that predates internal/notify ────────────
//
// internal/notify and models/notification_claim.go travel with
// writeFrameworkOwnedFiles. The AutoMigrate list does not: models/user.go is
// the file people add their own models to, so Grit never rewrites it, and a
// table missing from that list is a table grit migrate never creates.
//
// The failure without this repair is the worst shape available. The code
// compiles, the package imports, and the first call to notify.Once fails at
// runtime against a table that does not exist. In a scheduled sweep that means
// a job that errors every run, which is exactly the place nobody is watching.
//
// It is also the reason the table is registered rather than created on first
// use: a table that appears when a job runs is a table missing from every
// backup taken before then, and a backup that restores a deployment with no
// record of what it has already sent makes the next sweep send all of it again.

// repairNotifyModel adds the claims table to an existing project's model list.
func repairNotifyModel(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	models := filepath.Join(apiRoot, "internal", "models", "user.go")
	if !fileExists(models) {
		return nil
	}
	// Only once the model file is actually there. writeFrameworkOwnedFiles runs
	// before the repairs, so it will be, but a project whose own upgrade was
	// interrupted should not have a registration for a type it has not got:
	// that is a project that does not compile, from a repair meant to help.
	if !fileExists(filepath.Join(apiRoot, "internal", "models", "notification_claim.go")) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, models, repairNotifyModelSource)
}

func repairNotifyModelSource(src string) (string, []string, []string) {
	// Guarded on the registration itself, not on prose about it.
	if strings.Contains(src, "&NotificationClaim{}") {
		return src, nil, nil
	}
	if strings.Count(src, notifyModelAnchor) != 1 {
		return src, nil, []string{
			"the model list has no // grit:models marker, so the notification claims " +
				"table was not registered: add &NotificationClaim{} to it, or notify.Once " +
				"writes to a table AutoMigrate never created",
		}
	}
	// After the marker, not before it, and the reason is the saga repair.
	//
	// Two repairs insert into this list, and each has to arrive at the same
	// bytes the template has whether the other has run or not. Both inserting
	// before // grit:models cannot do that: a project missing only the saga
	// block gets its registration on the far side of this one, and the two
	// files differ for ever after. Inserting on opposite sides of the marker
	// makes them independent. Nothing depends on the order of this table in
	// the list: it has no foreign keys.
	out := strings.Replace(src, notifyModelAnchor, notifyModelAnchor+notifyModelRegistration, 1)
	return out, []string{
		"the notification claims table is in the model list, so migrate creates it " +
			"and a backup includes it",
	}, nil
}
