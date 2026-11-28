package main

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Every flag the binary accepts is in the Command Explorer.
//
// Seven flags of `grit generate resource` were not: --owned-by,
// --tenant-owned, --confidential, --append-only, --audit-reads, --retain and
// --unique existed, worked, and appeared nowhere in docs/config/cli-commands.ts.
// Nothing said so. The page rendered a shorter table and looked complete, and
// the flags were unfindable by search, because search reads that config.
//
// That is a two-sided fault, which is the kind a test on one side cannot
// catch: the CLI is complete by construction and the config is valid TypeScript
// whatever it leaves out. So this reads both and compares them.
//
// Not the reverse direction. The config documents a few spellings the binary
// does not register as flags of their own, shorthand pairs written as
// "--seed / --faker" among them, and a test that insisted on an exact match
// would be a test about punctuation.
func TestEveryFlagIsDocumented(t *testing.T) {
	documented := documentedFlags(t)

	// The commands the explorer covers. A command absent from the config
	// entirely is a separate gap, reported by TestEveryCommandIsDocumented.
	var missing []string
	walk(rootCommand(), "grit", func(path string, flag *pflag.Flag) {
		if flag.Hidden || flag.Deprecated != "" {
			return
		}
		// The root's own flags and cobra's built-in help are not per-command
		// documentation.
		if flag.Name == "help" || path == "grit" {
			return
		}
		command := commandKey(path)
		if _, covered := documented[command]; !covered {
			return
		}
		if !documented[command]["--"+flag.Name] {
			missing = append(missing, path+" --"+flag.Name)
		}
	})

	// A budget rather than a verdict.
	//
	// There were eighty-one of these when the check was written, most of them
	// on grit deploy, and failing the suite for all of them would have meant
	// either a docs sweep nobody asked for or a test somebody deletes. So the
	// known gaps are recorded, this fails on a new one, and the recorded list
	// may only shrink: a line whose flag has since been documented is also a
	// failure, so the file cannot be padded and cannot go stale.
	known := knownUndocumented(t)
	var fresh, fixed []string
	for _, gap := range missing {
		if !known[gap] {
			fresh = append(fresh, gap)
		}
		delete(known, gap)
	}
	for gap := range known {
		fixed = append(fixed, gap)
	}

	if len(fresh) > 0 {
		sort.Strings(fresh)
		t.Errorf("%d flag(s) the binary accepts are not in docs/config/cli-commands.ts:\n  %s\n\n"+
			"    Add each one to that command's flags array. The Command Explorer renders\n"+
			"    it, and both search indexes read it, so an undocumented flag is also an\n"+
			"    unfindable one: searching for it returns whatever page happens to share\n"+
			"    a word with it.",
			len(fresh), strings.Join(fresh, "\n  "))
	}
	if len(fixed) > 0 {
		sort.Strings(fixed)
		t.Errorf("%d line(s) in testdata/undocumented-flags.txt name a flag that IS now\n"+
			"documented. Delete them: the budget only counts while it is true.\n  %s",
			len(fixed), strings.Join(fixed, "\n  "))
	}
}

// knownUndocumented reads the budget file: "grit deploy --port" per line,
// with # comments and blank lines ignored.
func knownUndocumented(t *testing.T) map[string]bool {
	t.Helper()
	path := filepath.Join("testdata", "undocumented-flags.txt")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	out := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out[line] = true
	}
	return out
}

// documentedFlags is command name -> the flags its entry lists.
//
// Parsed rather than evaluated, the way every other reader of this config
// works: it is TypeScript, and the alternative is a JavaScript runtime in a
// Go test.
func documentedFlags(t *testing.T) map[string]map[string]bool {
	t.Helper()
	path := filepath.Join("..", "..", "docs", "config", "cli-commands.ts")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	source := string(data)

	// \r? because the file is CRLF on Windows: an anchor that forgets it
	// matches nothing, and this reader would then report "the config has no
	// commands" rather than a mistake in its own regex.
	nameRe := regexp.MustCompile(`(?m)^    name: '((?:[^'\\]|\\.)*)',\r?$`)
	flagRe := regexp.MustCompile(`\bflag: '((?:[^'\\]|\\.)*)'`)
	tokenRe := regexp.MustCompile(`--[a-z0-9-]+`)

	starts := nameRe.FindAllStringSubmatchIndex(source, -1)
	out := map[string]map[string]bool{}
	for i, loc := range starts {
		end := len(source)
		if i+1 < len(starts) {
			end = starts[i+1][0]
		}
		name := source[loc[2]:loc[3]]
		flags := map[string]bool{}
		for _, m := range flagRe.FindAllStringSubmatch(source[loc[0]:end], -1) {
			for _, token := range tokenRe.FindAllString(m[1], -1) {
				flags[token] = true
			}
		}
		out[commandKey(name)] = flags
	}
	if len(out) < 30 {
		t.Fatalf("read %d commands from cli-commands.ts, which is far fewer than it has: "+
			"the reader has stopped matching and this test is checking nothing", len(out))
	}
	return out
}

// commandKey is the comparable form of a command's name: "grit generate
// resource" from either the config's name or the walked path.
func commandKey(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

// walk visits every flag of every subcommand, with the full path to it.
func walk(cmd *cobra.Command, path string, visit func(string, *pflag.Flag)) {
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) { visit(path, f) })
	for _, child := range cmd.Commands() {
		if child.Hidden || child.Name() == "help" || child.Name() == "completion" {
			continue
		}
		walk(child, path+" "+child.Name(), visit)
	}
}

// Every command in the binary has an entry, so the explorer is the whole CLI
// rather than the part somebody remembered to write down.
//
// A warning rather than a failure: a command added today and documented
// tomorrow is normal, and failing the suite for it would make the next person
// delete the test instead of the gap.
func TestEveryCommandIsDocumented(t *testing.T) {
	documented := documentedFlags(t)
	seen := map[string]bool{}
	var undocumented []string
	walk(rootCommand(), "grit", func(path string, _ *pflag.Flag) {
		key := commandKey(path)
		if key == "grit" || seen[key] {
			return
		}
		seen[key] = true
		if _, ok := documented[key]; !ok {
			undocumented = append(undocumented, key)
		}
	})
	if len(undocumented) > 0 {
		sort.Strings(undocumented)
		t.Logf("%d command(s) are not in the Command Explorer:\n  %s",
			len(undocumented), strings.Join(undocumented, "\n  "))
	}
}
