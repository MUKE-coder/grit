package doctor

import (
	"regexp"
	"strings"
)

// An AI_GATEWAY_MODEL that no gateway will ever resolve.
//
// Grit's own default was `anthropic/claude-sonnet-4-6` for months, and that is
// not a model: the gateway's catalogue spells every Anthropic minor with a dot,
// `anthropic/claude-sonnet-4.6`. So the AI endpoints answered with the
// gateway's model-not-found for anybody who had not set the variable, which is
// everybody who read the .env.example and left it alone (grit#186).
//
// This is a shape check, not a lookup. `grit doctor` reads files and
// compiled-in floors and never the network, which is what makes it safe in CI
// and on a machine with no credentials, and asking the gateway would also need
// the project's API key. The shape is enough for the mistake that actually
// happens: a dashed minor is never right, and an id with no provider in front
// of it is never right either.
//
// It deliberately does not police the model name itself. A catalogue of 498
// models that grows every month cannot be compiled into a binary people
// upgrade quarterly, and a check that calls a real new model wrong is worse
// than no check.
var (
	// A minor version written with a dash, at the end of the name:
	// claude-sonnet-4-6. Anchored, so claude-3-5-sonnet-20241022 and any name
	// ending in a date are left alone.
	dashedMinorRe = regexp.MustCompile(`-\d+-\d+$`)
	// provider/model, which is the format the gateway documents.
	providerModelRe = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*/[A-Za-z0-9][A-Za-z0-9._-]*$`)
)

func checkAIModelID(p *project) []Finding {
	model := strings.TrimSpace(p.env["AI_GATEWAY_MODEL"])
	if model == "" {
		// Unset: the compiled-in default applies, and that is Grit's to get
		// right rather than this project's.
		return nil
	}

	// A project with no key is not calling the gateway, so a bad id there is
	// something to fix before it is used, not something that is failing now.
	level := "warning"
	if strings.TrimSpace(p.env["AI_GATEWAY_API_KEY"]) != "" {
		level = "error"
	}

	if !providerModelRe.MatchString(model) {
		return []Finding{{
			Level:    level,
			Resource: "AI_GATEWAY_MODEL",
			Message: "is " + model + ", which is not in the gateway's provider/model form, " +
				"so every AI request answers with a model error rather than a completion",
			Fix: "write the provider in front of it, as in anthropic/claude-sonnet-5.5; " +
				"the catalogue is at https://ai-gateway.vercel.sh/v1/models",
		}}
	}

	provider, name, _ := strings.Cut(model, "/")
	if strings.EqualFold(provider, "anthropic") && dashedMinorRe.MatchString(name) {
		dotted := dashedMinorRe.ReplaceAllStringFunc(name, func(s string) string {
			// -4-6 becomes -4.6: only the last dash moves.
			i := strings.LastIndex(s, "-")
			return s[:i] + "." + s[i+1:]
		})
		return []Finding{{
			Level:    level,
			Resource: "AI_GATEWAY_MODEL",
			Message: "is " + model + ", and the gateway has no such model: every Anthropic " +
				"id there writes the minor version with a dot, so AI requests answer with " +
				"a model error rather than a completion",
			Fix: "use anthropic/" + dotted + ", or anthropic/claude-sonnet-5.5, which is " +
				"newer and cheaper; the catalogue is at https://ai-gateway.vercel.sh/v1/models",
		}}
	}
	return nil
}
