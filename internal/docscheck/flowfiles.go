package docscheck

// A verified flow declared in a file of its own.
//
// grit#168. Executed coverage sat at four blocks while the docs grew to 193
// pages, and that is not a discipline problem: adding a verify= block means
// editing a React component, finding the right place in it, and getting a
// template literal inside JSX right. The two numbers were the finding.
//
// So a flow can also be a small markdown file: frontmatter naming the page it
// checks, prose saying why it exists, and one or more fenced shell blocks. The
// runner reads both sources and does not care which a flow came from, so
// `tests/docs/run-flows.sh` is unchanged and the cost of adding a check is now
// a ten-line file.
//
// Why not MDX for the pages themselves, which is what the issue asked for: the
// 289 existing pages are TSX built from a dozen components, two separate search
// indexers parse page.tsx, and so do three check scripts. Migrating them is a
// project. A flow file gets the coverage growing today and takes nothing away
// from doing the migration later, because this reader does not care what the
// page is written in.
//
// The frontmatter's page: is checked against the docs tree, because a flow
// pointing at a page that has moved is a check nobody will think to look at
// again.

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FlowFileDir is where flow files live, relative to the repository root.
const FlowFileDir = "docs/content/verified"

// FlowFile is one declared flow, as read from disk.
type FlowFile struct {
	// ID is the file name without its extension, and the flow's name.
	ID string
	// Page is the docs route this flow checks, as written in the frontmatter.
	Page string
	// Title is the one-line description `-list` shows.
	Title string
	// Path is the file, relative to the repository root where possible.
	Path string
	// Blocks are the fenced shell blocks, in file order.
	Blocks []Block
}

// ReadFlowFiles reads every flow file in dir.
//
// A missing directory is not an error: a repository with no flow files has no
// flows declared this way, which is the state this started from.
func ReadFlowFiles(dir string) ([]FlowFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", filepath.ToSlash(dir), err)
	}

	var out []FlowFile
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || strings.HasPrefix(name, "_") {
			continue
		}
		path := filepath.Join(dir, name)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil, readErr
		}
		flow, parseErr := parseFlowFile(filepath.ToSlash(path), strings.TrimSuffix(name, ".md"), string(data))
		if parseErr != nil {
			return nil, parseErr
		}
		out = append(out, flow)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// parseFlowFile reads the frontmatter and the fenced blocks.
//
// The frontmatter is key: value lines and nothing more. A YAML parser would be
// a dependency for three keys, and a flow file that needs anything richer is a
// flow file doing too much.
func parseFlowFile(path, id, src string) (FlowFile, error) {
	text := strings.ReplaceAll(src, "\r\n", "\n")
	lines := strings.Split(text, "\n")
	flow := FlowFile{ID: id, Path: path}

	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return flow, fmt.Errorf("%s: a flow file opens with --- and a page: line", path)
	}
	body := 0
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			body = i + 1
			break
		}
		key, value, ok := strings.Cut(lines[i], ":")
		if !ok {
			continue
		}
		value = strings.TrimSpace(value)
		switch strings.TrimSpace(key) {
		case "page":
			flow.Page = value
		case "title":
			flow.Title = value
		}
	}
	if body == 0 {
		return flow, fmt.Errorf("%s: the frontmatter is never closed", path)
	}
	if flow.Page == "" {
		return flow, fmt.Errorf("%s: no page: in the frontmatter. A flow that names no page is a "+
			"script nobody will find when the page it checks changes", path)
	}
	if flow.Title == "" {
		return flow, fmt.Errorf("%s: no title: in the frontmatter", path)
	}

	blocks, err := fencedShellBlocks(path, lines[body:], body)
	if err != nil {
		return flow, err
	}
	if len(blocks) == 0 {
		return flow, fmt.Errorf("%s: no ```bash block. A flow with no commands passes without "+
			"running anything, which is the one outcome worse than no flow at all", path)
	}
	for i := range blocks {
		blocks[i].Flow = id
	}
	flow.Blocks = blocks
	return flow, nil
}

// fencedShellBlocks returns the ```bash (or ```sh) blocks, in order.
//
// offset is the line number the slice starts at, so a failure names the line in
// the file rather than in the slice.
func fencedShellBlocks(path string, lines []string, offset int) ([]Block, error) {
	var out []Block
	for i := 0; i < len(lines); i++ {
		fence := strings.TrimSpace(lines[i])
		if fence != "```bash" && fence != "```sh" {
			// A block in another language is documentation inside the flow
			// file, not something to run.
			if strings.HasPrefix(fence, "```") {
				for i++; i < len(lines); i++ {
					if strings.TrimSpace(lines[i]) == "```" {
						break
					}
				}
			}
			continue
		}
		var script []string
		closed := false
		for j := i + 1; j < len(lines); j++ {
			if strings.TrimSpace(lines[j]) == "```" {
				i = j
				closed = true
				break
			}
			script = append(script, lines[j])
		}
		if !closed {
			return nil, fmt.Errorf("%s:%d: the block is never closed", path, offset+i+1)
		}
		out = append(out, Block{
			File:   path,
			Line:   offset + i + 1,
			Script: strings.Join(script, "\n"),
		})
	}
	return out, nil
}

// AllFlows is every flow the runner knows about: the verify= blocks in the
// pages, and the flow files.
//
// An id used by both is an error. Merging them silently would run a page's
// blocks and a file's blocks as one flow in an order neither author chose.
func AllFlows(appRoot, flowDir string) (map[string][]Block, error) {
	flows, err := Flows(appRoot)
	if err != nil {
		return nil, err
	}
	files, err := ReadFlowFiles(flowDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if existing, clash := flows[file.ID]; clash {
			return nil, fmt.Errorf("flow %q is declared twice: in %s and in %s. Rename one, or "+
				"the two sets of commands run as one flow in an order neither author chose",
				file.ID, file.Path, existing[0].File)
		}
		flows[file.ID] = file.Blocks
	}
	return flows, nil
}

// FlowFilePages reports a flow file whose page: does not exist under appRoot.
//
// Returned rather than logged, so the check script can fail on it. A flow
// pointing at a page that has moved still runs, and still passes, and has
// stopped being evidence about anything a reader can find.
func FlowFilePages(appRoot, flowDir string) ([]string, error) {
	files, err := ReadFlowFiles(flowDir)
	if err != nil {
		return nil, err
	}
	var dangling []string
	for _, file := range files {
		rel := strings.TrimPrefix(file.Page, "/")
		candidates := []string{
			filepath.Join(appRoot, rel, "page.tsx"),
			filepath.Join(appRoot, rel, "page.mdx"),
			filepath.Join(appRoot, rel+".tsx"),
			filepath.Join(appRoot, rel+".mdx"),
		}
		found := false
		for _, candidate := range candidates {
			if _, statErr := os.Stat(candidate); statErr == nil {
				found = true
				break
			}
		}
		if !found {
			dangling = append(dangling, fmt.Sprintf("%s names page %s, which does not exist", file.Path, file.Page))
		}
	}
	return dangling, nil
}
