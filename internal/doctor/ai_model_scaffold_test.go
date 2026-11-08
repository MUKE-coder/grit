package doctor

import (
	"os"
	"regexp"
	"sort"
	"testing"
)

// scaffoldedAIModels is every AI_GATEWAY_MODEL value the scaffold writes into a
// new project, and the default compiled into its config.
//
// Read as text rather than by importing internal/scaffold. The only thing
// needed is the literal, the files are in this repository, and a text read
// cannot drag the whole scaffold package into the doctor's tests.
func scaffoldedAIModels(t *testing.T) []string {
	t.Helper()

	// The .env.example lines, and the getEnv fallback in the generated config.
	patterns := map[string]*regexp.Regexp{
		"../scaffold/root_files.go": regexp.MustCompile(`AI_GATEWAY_MODEL=([^\s"'` + "`" + `]+)`),
		"../scaffold/api_files.go":  regexp.MustCompile(`getEnv\("AI_GATEWAY_MODEL",\s*"([^"]+)"\)`),
		"../scaffold/docs_files.go": regexp.MustCompile(`AI_GATEWAY_MODEL=([^\s"'` + "`" + `]+)`),
	}

	seen := map[string]bool{}
	for path, re := range patterns {
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("reading %s: %v", path, err)
		}
		found := re.FindAllStringSubmatch(string(body), -1)
		if len(found) == 0 {
			t.Fatalf("no AI_GATEWAY_MODEL value in %s: this test has to be "+
				"pointed at wherever the scaffold writes it now, not quietly pass", path)
		}
		for _, m := range found {
			seen[m[1]] = true
		}
	}

	out := make([]string, 0, len(seen))
	for m := range seen {
		out = append(out, m)
	}
	sort.Strings(out)
	return out
}
