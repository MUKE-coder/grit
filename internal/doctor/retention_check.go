package doctor

import (
	"fmt"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/generate"
)

// A resource owned by a user, with no retention window.
//
// Its rows stay until somebody deletes them or the user is erased, which is
// fine for a user's own notes and wrong for a chat log a privacy policy says
// goes after ninety days. Only the project knows which, so this says what the
// model shows, a user's rows with nothing that removes them, and does not call
// the data personal (grit#260). Owned is the signal because it is the one the
// model carries; a column that merely names a user is checked by
// resource-could-be-owned first.
//
// A warning, because "nothing to do" is a correct answer for a lot of these,
// and the fix says so.
func checkOwnedRetention(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		m := p.models[name]
		if m.OwnerColumn == "" || m.Retained {
			continue
		}
		// Append-only rows can never be deleted, so --retain is refused for
		// them. Suggesting it would be advice the generator turns down.
		if m.AppendOnly {
			continue
		}
		message := "is owned by a user and has no retention window: its rows stay until they are deleted or the user is erased"
		// --confidential is the project saying the rows are sensitive, so the
		// question is worth asking more plainly. It is generated with
		// ScopeConfidential in place of ScopeOwned.
		if strings.Contains(p.services[p.snake[name]], "authz.ScopeConfidential(") {
			message = "is confidential and owned by a user, and has no retention window: its rows stay until they are deleted or the user is erased"
		}
		out = append(out, Finding{
			Level:    "warning",
			Resource: name,
			Message:  message,
			Fix: "if the rows should go after a fixed time, add func init() { retention.Register(&" + name +
				"{}, 90*24*time.Hour" + retentionOptions(p, m) + ") } to its model file, which is what --retain 90d writes. " +
				"If they should stay as long as the account does, nothing to do.",
		})
	}
	return out
}

// retentionOptions is what the suggested line needs for the resource's shape,
// the same options --retain writes. Without them the line compiles and the
// database refuses the sweep every night: a record cannot go before its lines,
// nor a tree node before its children.
func retentionOptions(p *project, m *model) string {
	// Both shapes are read by internal/generate, beside the code that writes
	// them, so this cannot drift from the model the generator produces.
	// TestTheModelReadersMatchWhatTheGeneratorWrites runs both ends together.
	var opts string
	if child := generate.ItemsChildModel(m.Source); child != "" {
		opts += fmt.Sprintf(", retention.WithLines(&%s{}, %q)", child, p.snake[m.Name]+"_id")
	}
	if generate.IsTreeModel(m.Source) {
		opts += ", retention.AsTree()"
	}
	return opts
}
