package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/MUKE-coder/grit/v3/internal/docs"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// The documentation tools.
//
// The gap these close is version drift, which is Grit's sharpest exposure to an
// agent. Grit ships several releases a week; an agent working in a project
// pinned to v3.350 that reads the docs site is reading v3.393, and the symptom
// is generated code calling a helper that did not exist yet or a flag that has
// since been renamed. Neither is a compile error in the agent's head, so it
// writes the code confidently and the developer finds out.
//
// Every answer carries the version it is true of, and says so when the project
// disagrees. An agent that knows the documentation is newer than the project
// can tell the developer to upgrade instead of generating something that will
// not build.

// searchDocs answers a question from the embedded index.
func (s *Server) searchDocs(raw json.RawMessage) (string, error) {
	var args struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &args); err != nil {
			return "", fmt.Errorf("reading arguments: %w", err)
		}
	}
	if args.Query == "" {
		return "", fmt.Errorf("query is required: what do you want to know")
	}
	if args.Limit <= 0 || args.Limit > 20 {
		args.Limit = 5
	}

	hits, err := docs.Search(args.Query, args.Limit)
	if err != nil {
		return "", err
	}

	results := make([]map[string]interface{}, 0, len(hits))
	for _, hit := range hits {
		results = append(results, map[string]interface{}{
			"url":     hit.URL,
			"title":   hit.Title,
			"kind":    hit.Kind,
			"snippet": hit.Snippet,
			"score":   hit.Score,
		})
	}

	out := map[string]interface{}{
		"documentation_version": docs.Version(),
		"query":                 args.Query,
		"results":               results,
		"next":                  "Call grit_read_doc with a url for the whole page.",
	}
	if len(results) == 0 {
		out["note"] = "Nothing matched. Try fewer or different words; this index is prose " +
			"from the docs site and does not include code samples."
	}
	s.noteVersionDrift(out)
	return asJSON(out)
}

// readDoc returns one page.
func (s *Server) readDoc(raw json.RawMessage) (string, error) {
	var args struct {
		URL string `json:"url"`
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &args); err != nil {
			return "", fmt.Errorf("reading arguments: %w", err)
		}
	}
	if args.URL == "" {
		return "", fmt.Errorf("url is required, e.g. /docs/concepts/field-types; " +
			"use grit_search_docs to find one")
	}

	page, err := docs.Read(args.URL)
	if err != nil {
		return "", err
	}

	out := map[string]interface{}{
		"documentation_version": docs.Version(),
		"url":                   page.URL,
		"title":                 page.Title,
		"description":           page.Description,
		"kind":                  page.Kind,
		"headings":              page.Headings,
		"text":                  page.Text,
	}
	s.noteVersionDrift(out)
	return asJSON(out)
}

// noteVersionDrift adds the warning when the project was scaffolded with a
// different version from the one these docs describe.
//
// A field in the payload rather than a log line, because an agent reads the
// payload and nothing else. It is the most useful thing in the response when it
// is present: it is the difference between "this is how it works" and "this is
// how it works in a version you are not running".
func (s *Server) noteVersionDrift(out map[string]interface{}) {
	if s.Root == "" {
		return
	}
	projectVersion := scaffold.ProjectVersion(s.Root)
	if projectVersion == "" || projectVersion == docs.Version() {
		return
	}
	out["project_version"] = projectVersion
	out["version_warning"] = fmt.Sprintf(
		"This documentation is for v%s and this project was scaffolded with v%s. "+
			"Anything described here that is newer than v%s will not exist in this project. "+
			"Tell the developer to run `grit update`, or check grit_cli_reference for what "+
			"this project's CLI actually supports.",
		docs.Version(), projectVersion, projectVersion)
}
