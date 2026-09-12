// Package docscheck reads the commands the documentation shows and checks them
// against the CLI that has to run them.
//
// Docs drift is the most common root cause in this project's own changelog: a
// page told somebody to run a command that had been renamed, or passed a flag
// that never existed, or showed a field type the generator would refuse. Every
// one of those was found by a person typing it. Nothing in the build knew the
// docs existed.
//
// So this extracts every grit command line out of the docs pages and asks the
// real command tree and the real field parser about it. It is deliberately
// conservative: a line it cannot read confidently is skipped rather than
// reported, because a checker that cries wolf gets switched off.
package docscheck

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Sample is one command line a docs page shows.
type Sample struct {
	File    string // relative to the docs root
	Line    int    // 1-indexed, so an editor can jump to it
	Command string // cleaned of the JSX around it
	// Verify groups the runnable samples: blocks marked verify="<id>" in the
	// docs are executed in order against a scaffolded project by the docs
	// workflow. Empty for a block nobody claimed is runnable.
	Verify string
}

// Problem is a command the docs show that the CLI would refuse.
type Problem struct {
	Sample  Sample
	Message string
}

func (p Problem) String() string {
	return fmt.Sprintf("%s:%d: %s\n    %s", p.Sample.File, p.Sample.Line, p.Message, p.Sample.Command)
}

// Extract reads every docs page under root and returns the grit commands they
// show, in file order.
func Extract(root string) ([]Sample, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if name := info.Name(); name == "node_modules" || name == ".next" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".tsx") || strings.HasSuffix(path, ".mdx") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	var samples []Sample
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		samples = append(samples, extractFile(filepath.ToSlash(rel), string(data))...)
	}
	return samples, nil
}

// extractFile pulls the command lines out of one page.
//
// A docs page is TSX, and a command sits in one of two places: inside a template
// literal spanning several lines, or inside a single-line prop such as
// code="grit migrate" or cmd="grit migrate --fresh". Only those count. Prose is
// not a command, however it starts, and the first version of this learned that
// from a changelog sentence opening "grit can keep itself current" and a heading
// reading "grit sync: Manual Type Generation".
func extractFile(rel, src string) []Sample {
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")
	verify := ""
	inLiteral := false
	var out []Sample

	for i := 0; i < len(lines); i++ {
		raw := lines[i]
		openedInLiteral := inLiteral
		inLiteral = inLiteral != (strings.Count(raw, "`")%2 == 1)

		// verify="<id>" on a CodeBlock marks the block that follows as runnable.
		// Inside a literal it is a page showing what a marker looks like, not one.
		if !openedInLiteral {
			if id, ok := verifyID(raw); ok {
				verify = id
			}
		}

		var candidates []string
		if openedInLiteral {
			// Inside a block: the line is the command, with any closing markup
			// after it.
			candidates = []string{raw}
		} else {
			candidates = quotedPayloads(raw)
			// A block usually starts its first command on the same line as
			// code={`, so the text after the opening backtick is a command too.
			// Missing it meant missing the first line of every multi-line block,
			// which is most of the commands in these docs.
			if inLiteral {
				if at := strings.LastIndex(raw, "`"); at >= 0 {
					candidates = append(candidates, raw[at+1:])
				}
			}
		}
		// A JS string can carry an escaped newline, in which case one payload holds
		// two commands one after the other.
		candidates = splitEscapedNewlines(candidates)

		found := false
		for _, candidate := range candidates {
			command, ok := cleanCommand(candidate)
			if !ok {
				continue
			}
			// A line ending in a backslash continues on the next one, which is how
			// the docs show a long --fields spec.
			for strings.HasSuffix(command, "\\") && i+1 < len(lines) {
				i++
				inLiteral = inLiteral != (strings.Count(lines[i], "`")%2 == 1)
				command = strings.TrimSpace(strings.TrimSuffix(command, "\\")) + " " +
					strings.TrimSpace(stripTrailingJSX(lines[i]))
			}
			out = append(out, Sample{File: rel, Line: i + 1, Command: command, Verify: verify})
			found = true
		}

		// The block is over once its literal closes, so the next block has to
		// claim verify for itself.
		if verify != "" && openedInLiteral && !inLiteral {
			verify = ""
		} else if verify != "" && found && !openedInLiteral && !inLiteral {
			verify = ""
		}
	}
	return out
}

// quotedPayloads returns the quoted and backticked strings on one line, which is
// where a single-line code block keeps its command.
func quotedPayloads(line string) []string {
	var out []string
	for _, quote := range []byte{'"', '`'} {
		rest := line
		for {
			start := strings.IndexByte(rest, quote)
			if start < 0 {
				break
			}
			end := strings.IndexByte(rest[start+1:], quote)
			if end < 0 {
				break
			}
			out = append(out, rest[start+1:start+1+end])
			rest = rest[start+1+end+1:]
		}
	}
	return out
}

// splitEscapedNewlines breaks payloads that hold more than one line.
func splitEscapedNewlines(candidates []string) []string {
	var out []string
	for _, candidate := range candidates {
		out = append(out, strings.Split(candidate, `\n`)...)
	}
	return out
}

// verifyID reads verify="<id>" off a CodeBlock prop.
func verifyID(line string) (string, bool) {
	const marker = `verify="`
	at := strings.Index(line, marker)
	if at < 0 {
		return "", false
	}
	rest := line[at+len(marker):]
	end := strings.Index(rest, `"`)
	if end <= 0 {
		return "", false
	}
	return rest[:end], true
}

// cleanCommand decides whether a line is a grit command, and strips the JSX and
// shell noise around it.
func cleanCommand(raw string) (string, bool) {
	line := decodeEntities(strings.TrimSpace(strings.TrimRight(raw, "\r")))
	line = strings.TrimPrefix(line, "$ ")
	if !strings.HasPrefix(line, "grit ") {
		return "", false
	}
	line = stripTrailingJSX(line)

	// A comment after the command is prose, and so is anything after a pipe or a
	// chain: each piece would be its own command line.
	for _, cut := range []string{" #", " |", " &&", " ;", "  # "} {
		if at := strings.Index(line, cut); at > 0 {
			line = line[:at]
		}
	}
	line = strings.TrimSpace(line)

	// Lines that are still markup, or that trail off, say nothing checkable.
	if line == "" || strings.ContainsAny(line, "<>{}") || strings.HasSuffix(line, "...") {
		return "", false
	}
	// A comma followed by a space is prose. A field spec separates with a bare
	// comma, so "grit seed, faker, and per-resource seeders" is a sentence that
	// happens to start with a command name.
	if strings.Contains(line, ", ") {
		return "", false
	}
	return line, true
}

// decodeEntities turns the entities JSX text needs back into the characters a
// shell would see. A --fields value written with &quot; around it is a valid docs
// page and an invalid command until this runs.
func decodeEntities(line string) string {
	for _, pair := range [][2]string{
		{"&quot;", `"`}, {"&#34;", `"`}, {"&apos;", "'"}, {"&#39;", "'"},
		{"&lt;", "<"}, {"&gt;", ">"}, {"&nbsp;", " "}, {"&amp;", "&"},
	} {
		line = strings.ReplaceAll(line, pair[0], pair[1])
	}
	return line
}

// stripTrailingJSX removes what surrounds a command inside a template literal.
func stripTrailingJSX(line string) string {
	line = strings.TrimSpace(line)
	// A trailing quote is left alone: it closes the command's own --fields value,
	// and cutting it turned every quoted spec into an unbalanced one.
	for _, suffix := range []string{"`} />", "`}/>", "`}", "`", `"}`, "/>"} {
		line = strings.TrimSpace(strings.TrimSuffix(line, suffix))
	}
	return line
}

// Fields returns the value of a --fields or --items flag in a sample, which the
// generator's own parser can then be asked about.
func Fields(command, flag string) (string, bool) {
	tokens, err := Tokenize(command)
	if err != nil {
		return "", false
	}
	for i, token := range tokens {
		if token == "--"+flag && i+1 < len(tokens) {
			return tokens[i+1], true
		}
		if strings.HasPrefix(token, "--"+flag+"=") {
			return strings.TrimPrefix(token, "--"+flag+"="), true
		}
	}
	return "", false
}

// Tokenize splits a command line the way a shell would, honouring quotes.
func Tokenize(command string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	quote := rune(0)
	started := false

	for _, r := range command {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
				continue
			}
			current.WriteRune(r)
		case r == '\'' || r == '"':
			quote = r
			started = true
		case r == ' ' || r == '\t':
			if started || current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
				started = false
			}
		default:
			current.WriteRune(r)
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("unclosed %c quote", quote)
	}
	if started || current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens, nil
}
