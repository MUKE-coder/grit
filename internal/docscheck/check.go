package docscheck

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/generate"
)

// Check asks the command tree about every command the docs show.
//
// What it catches is the drift this project keeps paying for: a command that was
// renamed, a flag that never existed or was removed, a subcommand the docs
// invented. What it deliberately does not do is guess: a line it cannot read as a
// command is skipped by Extract, and a placeholder argument is accepted as given.
func Check(samples []Sample, root *cobra.Command) []Problem {
	var problems []Problem
	for _, sample := range samples {
		if isHistory(sample.File) {
			continue
		}
		if problem, bad := check(sample, root); bad {
			problems = append(problems, problem)
		}
	}
	return problems
}

// isHistory marks pages that describe the past on purpose.
//
// The changelog says what a release did at the time, including commands and flags
// that have since been renamed or removed. Checking it against today's CLI would
// force the history to be edited, which is the opposite of what it is for.
func isHistory(file string) bool {
	return strings.HasPrefix(file, "docs/changelog/") || strings.HasPrefix(file, "docs/versioning/")
}

func check(sample Sample, root *cobra.Command) (Problem, bool) {
	tokens, err := Tokenize(sample.Command)
	if err != nil {
		// An unbalanced quote in a docs sample is a real problem: copied into a
		// shell, it hangs waiting for the rest.
		return Problem{sample, err.Error()}, true
	}
	if len(tokens) == 0 || tokens[0] != "grit" {
		return Problem{}, false
	}
	tokens = tokens[1:]

	// Walk as deep into the tree as the leading words go. The first word that is
	// not a subcommand is an argument, and arguments are the docs' business.
	cmd := root
	for len(tokens) > 0 && !strings.HasPrefix(tokens[0], "-") {
		child := childNamed(cmd, tokens[0])
		if child == nil {
			if cmd == root {
				return Problem{sample, fmt.Sprintf("no such command: grit %s", tokens[0])}, true
			}
			// A subcommand of a command that takes none is worth saying: the docs
			// are describing a command tree that is not there.
			if len(cmd.Commands()) > 0 && looksLikeSubcommand(cmd, tokens[0]) {
				return Problem{sample, fmt.Sprintf("grit %s has no %q subcommand",
					strings.TrimPrefix(cmd.CommandPath(), "grit "), tokens[0])}, true
			}
			break
		}
		cmd = child
		tokens = tokens[1:]
	}

	// cobra adds --help and -h when a command runs, not when it is built, so
	// without this every documented `grit --help` reads as an unknown flag.
	cmd.InitDefaultHelpFlag()

	// Then the flags, which is where docs age fastest.
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if !strings.HasPrefix(token, "-") || token == "-" || token == "--" {
			continue
		}
		name := strings.TrimLeft(token, "-")
		if at := strings.Index(name, "="); at >= 0 {
			name = name[:at]
		}
		if name == "" || isPlaceholder(name) {
			continue
		}

		flag := cmd.Flags().Lookup(name)
		if flag == nil {
			flag = cmd.InheritedFlags().Lookup(name)
		}
		if flag == nil && !strings.HasPrefix(token, "--") && len(name) == 1 {
			flag = cmd.Flags().ShorthandLookup(name)
		}
		if flag == nil {
			where := cmd.CommandPath()
			return Problem{sample, fmt.Sprintf("%s has no --%s flag", where, name)}, true
		}
		// A flag that takes a value eats the next token, so the value is not
		// mistaken for anything else.
		if flag.NoOptDefVal == "" && !strings.Contains(token, "=") {
			i++
		}
	}
	return Problem{}, false
}

// childNamed finds a subcommand by name or alias.
func childNamed(cmd *cobra.Command, name string) *cobra.Command {
	for _, child := range cmd.Commands() {
		if child.Name() == name {
			return child
		}
		for _, alias := range child.Aliases {
			if alias == name {
				return child
			}
		}
	}
	return nil
}

// looksLikeSubcommand distinguishes a misspelled subcommand from an argument.
//
// `grit generate resource Post` has Post as an argument; `grit generate resourcs`
// is a typo. Only lower-case words with no punctuation are treated as an intended
// subcommand, because every real one is spelled that way and the arguments the
// docs pass are names, specs or paths.
func looksLikeSubcommand(cmd *cobra.Command, token string) bool {
	if cmd.Args != nil {
		// The command accepts arguments of its own, so this may well be one.
		return false
	}
	for _, r := range token {
		if !(r >= 'a' && r <= 'z') && r != '-' {
			return false
		}
	}
	return true
}

// isPlaceholder spots the stand-ins docs use for a value.
func isPlaceholder(token string) bool {
	return strings.ContainsAny(token, "<>[]") || token == "..."
}

// CheckFieldSpecs runs the docs' field specs through the parser the generator
// uses, so a documented field type that does not exist fails the build.
//
// This is the other half of docs drift, and the worse half: a command that does
// not exist fails immediately and loudly, while a field spec the generator
// refuses wastes somebody's afternoon on a tutorial that cannot work.
func CheckFieldSpecs(samples []Sample) []Problem {
	var problems []Problem
	for _, sample := range samples {
		spec, resource, ok := fieldSpec(sample.Command)
		if !ok || spec == "" {
			continue
		}
		if strings.ContainsAny(spec, "<>{}") {
			continue // a placeholder, not a spec
		}
		if _, err := generate.ParseInlineFields(resource, spec); err != nil {
			problems = append(problems, Problem{sample, "field spec is not valid: " + err.Error()})
		}
	}
	return problems
}

// fieldSpec returns the field spec a command carries, and the resource it is for.
func fieldSpec(command string) (spec, resource string, ok bool) {
	tokens, err := Tokenize(command)
	if err != nil || len(tokens) < 3 {
		return "", "", false
	}
	// grit generate resource <Name> --fields "..."  (g is the documented alias)
	if (tokens[1] == "generate" || tokens[1] == "g") && tokens[2] == "resource" {
		if len(tokens) > 3 {
			resource = tokens[3]
		}
		if value, found := Fields(command, "fields"); found {
			return value, orDefault(resource, "Thing"), true
		}
		return "", "", false
	}
	// grit generate field <Resource> <spec>
	if (tokens[1] == "generate" || tokens[1] == "g") && tokens[2] == "field" && len(tokens) > 4 {
		return tokens[4], orDefault(tokens[3], "Thing"), true
	}
	return "", "", false
}

func orDefault(value, fallback string) string {
	if strings.TrimSpace(value) == "" || isPlaceholder(value) {
		return fallback
	}
	return value
}
