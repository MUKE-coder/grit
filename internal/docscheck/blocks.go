package docscheck

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Block is a code block the docs claim is runnable: one marked verify="<flow>".
//
// Checking that a command exists is cheap and catches renames. It does not catch
// a page whose steps are each valid and which does not work when followed in
// order, which is the other half of docs drift and the half that wastes an
// afternoon. So a block can opt in by name, and CI runs the blocks of a flow in
// page order against a real scaffolded project.
type Block struct {
	File   string // relative to the docs root
	Line   int    // where the block's code starts
	Flow   string // the verify="<flow>" id
	Script string // the block's contents, exactly as the page shows them
}

// Flows returns the runnable blocks grouped by flow, each in document order.
func Flows(root string) (map[string][]Block, error) {
	blocks, err := ExtractBlocks(root)
	if err != nil {
		return nil, err
	}
	flows := map[string][]Block{}
	for _, block := range blocks {
		flows[block.Flow] = append(flows[block.Flow], block)
	}
	return flows, nil
}

// ExtractBlocks reads every docs page and returns the blocks marked verify.
func ExtractBlocks(root string) ([]Block, error) {
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

	var blocks []Block
	for _, path := range files {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		found, err := extractBlocks(filepath.ToSlash(rel), string(data))
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, found...)
	}
	return blocks, nil
}

// extractBlocks finds the verify-marked blocks in one page.
//
// A marked block is a CodeBlock with verify="<flow>" somewhere in its props and
// its script in a code={`...`} template literal. Anything else marked verify is an
// error rather than a silent skip: a block that believes it is covered and is not
// is worse than one that never claimed to be.
func extractBlocks(rel, src string) ([]Block, error) {
	text := strings.ReplaceAll(src, "\r\n", "\n")
	lines := strings.Split(text, "\n")

	var blocks []Block
	inLiteral := false
	for i, line := range lines {
		insideAtStart := inLiteral
		inLiteral = inLiteral != (strings.Count(line, "`")%2 == 1)

		// A marker inside a code block is a page documenting markers, not a claim
		// about itself. Without this, the testing page's example of a marked block
		// joined the flow it was describing.
		if insideAtStart {
			continue
		}
		flow, ok := verifyID(line)
		if !ok {
			continue
		}
		script, start, err := scriptAfter(lines, i)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: verify=%q: %w", rel, i+1, flow, err)
		}
		blocks = append(blocks, Block{File: rel, Line: start, Flow: flow, Script: script})
	}
	return blocks, nil
}

// scriptAfter reads the code={`...`} literal that follows a verify marker.
func scriptAfter(lines []string, from int) (script string, start int, err error) {
	for i := from; i < len(lines) && i < from+12; i++ {
		at := strings.Index(lines[i], "code={`")
		if at < 0 {
			continue
		}
		rest := lines[i][at+len("code={`"):]
		// A single-line block closes on the same line.
		if end := strings.Index(rest, "`}"); end >= 0 {
			return rest[:end], i + 1, nil
		}
		var out []string
		if strings.TrimSpace(rest) != "" {
			out = append(out, rest)
		}
		for j := i + 1; j < len(lines); j++ {
			if end := strings.Index(lines[j], "`}"); end >= 0 {
				if head := lines[j][:end]; strings.TrimSpace(head) != "" {
					out = append(out, head)
				}
				return strings.Join(out, "\n"), i + 1, nil
			}
			out = append(out, lines[j])
		}
		return "", 0, fmt.Errorf("the code literal is never closed")
	}
	return "", 0, fmt.Errorf("no code={`...`} block follows the marker")
}

// Script assembles a flow into one shell script, in page order.
//
// Every block keeps its own comments, and each is preceded by where it came from,
// so a failure in CI names the page to fix rather than a line number in generated
// output.
func Script(flow string, blocks []Block) string {
	var out strings.Builder
	out.WriteString("#!/usr/bin/env bash\n")
	out.WriteString("# Generated from the docs. Flow: " + flow + "\n")
	out.WriteString("# Do not edit: fix the page the block comes from.\n")
	out.WriteString("set -euo pipefail\n\n")
	for _, block := range blocks {
		out.WriteString(fmt.Sprintf("echo \"--- %s:%d\"\n", block.File, block.Line))
		out.WriteString(strings.TrimRight(block.Script, "\n"))
		out.WriteString("\n\n")
	}
	return out.String()
}
