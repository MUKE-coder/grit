package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeFrameworkOwnedFiles writes the code the framework owns.
//
// These are not your files. Nothing generates them, nothing in a resource
// definition describes them, and they only make sense as a set: the webhook
// model, the receiver that fills it and the dispatch package all read the same
// columns by name.
//
// They live in their own function because they have to travel on upgrade. A
// framework fix that changes a constraint or a field type is useless if half
// the cluster stays at the version the project was scaffolded with,
// and the failure is quiet: the code compiles, the migration reports nothing
// to do, and the constraint everything assumes is simply absent. That is
// exactly how webhook deduplication spent several versions being decoration.
//
// Your own models, in the same package, are never touched by this. So is
// models/user.go, which holds the registry people add to. The API key model
// is not here either: it travels with its seeder in writeMigrateSeedFiles,
// because the two break each other when they drift apart.
//
// writeFile is manifest-guarded, so a file someone has edited is reported as a
// conflict rather than overwritten.
func writeFrameworkOwnedFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		// The webhook cluster. The model, the receiver that builds it and the
		// dispatch package all read the same columns, so upgrading one without
		// the others gives a project that does not compile: the receiver
		// assigning a string to a *string is how this was found.
		filepath.Join(apiRoot, "internal", "models", "webhook_event.go"): apiWebhookEventModelGo(),
		filepath.Join(apiRoot, "internal", "handlers", "webhooks.go"):    apiWebhooksHandlerGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "webhooks.go"):    apiWebhooksGo(),
		filepath.Join(apiRoot, "internal", "webhooks", "dedup_test.go"):  apiWebhookDedupTestGo(),

		filepath.Join(apiRoot, "internal", "models", "outbox_message.go"): apiOutboxModelGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}
