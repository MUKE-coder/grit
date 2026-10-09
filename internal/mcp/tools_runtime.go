package mcp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/devlog"
)

// The runtime tools: what the application printed, and what went wrong.
//
// Everything else this server exposes is parsed from the checkout, and says so.
// These two are the only ones that describe what actually happened, and they
// exist because a week of building a Grit application produced exactly one
// lesson worth generalising: every bug that mattered was found by reading
// output, and none of them showed in an exit code.
//
// A compile failure after generating a resource, a panic on the first request,
// a 500 from a handler that looks fine: all of them are a line in a log and
// invisible to anything that only reads source.

func (s *Server) devLogTool(raw json.RawMessage) (string, error) {
	var args struct {
		Lines int    `json:"lines"`
		Grep  string `json:"grep"`
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &args); err != nil {
			return "", fmt.Errorf("reading arguments: %w", err)
		}
	}
	if args.Lines <= 0 || args.Lines > 500 {
		args.Lines = 120
	}

	if !devlog.Exists(s.Root) {
		return asJSON(map[string]interface{}{
			"available": false,
			"reason": "There is no development log yet. It is written by `grit start`, " +
				"which tees everything the API and the frontends print into " +
				devlog.Dir + "/" + devlog.Name + ". Ask the developer to run `grit start`, " +
				"or to leave it running, and try again.",
		})
	}

	lines, err := devlog.Tail(s.Root, args.Lines)
	if err != nil {
		return "", err
	}

	if args.Grep != "" {
		needle := strings.ToLower(args.Grep)
		kept := make([]string, 0, len(lines))
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), needle) {
				kept = append(kept, line)
			}
		}
		lines = kept
	}

	return asJSON(map[string]interface{}{
		"available": true,
		"path":      devlog.Dir + "/" + devlog.Name,
		"lines":     lines,
		"note": "The tail of the current `grit start` run. The previous run is kept " +
			"alongside it as " + devlog.Name + ".1 for when something worked a minute ago.",
	})
}

func (s *Server) lastErrorsTool(raw json.RawMessage) (string, error) {
	var args struct {
		Limit int `json:"limit"`
	}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &args); err != nil {
			return "", fmt.Errorf("reading arguments: %w", err)
		}
	}
	if args.Limit <= 0 || args.Limit > 100 {
		args.Limit = 20
	}

	problems, err := devlog.Errors(s.Root, args.Limit)
	if err != nil {
		// No log is not an error, it is a thing to say.
		return asJSON(map[string]interface{}{
			"problems": []interface{}{},
			"note": "Nothing to read. The development log is written by `grit start`; " +
				"compile failures are also read from the API's tmp/build-errors.log, " +
				"which only exists once the hot reloader has run.",
		})
	}

	out := map[string]interface{}{
		"problems": problems,
		"count":    len(problems),
	}
	if len(problems) == 0 {
		out["note"] = "No errors in the current run. If the application is misbehaving " +
			"without logging anything, call grit_dev_log and read the tail: not every " +
			"failure says the word error."
	} else {
		out["note"] = "Newest last. A build failure means the binary never started; a panic " +
			"means it started and died. Check which before assuming the code is wrong."
	}
	return asJSON(out)
}
