package main

import (
	"strings"
	"testing"
)

// The access table is generated from the route files, so hand-editing a guard
// makes it stale. grit doctor said so and suggested `grit upgrade`, which
// moves the whole project onto the CLI's templates, or `grit generate
// resource`, which rewrites the route file and drops the guard that made the
// table stale in the first place (grit#104). Neither is a way to refresh one
// table.
func TestRoutesHasAWriteFlag(t *testing.T) {
	c := routesCmd()
	if c.Name() != "routes" {
		t.Fatalf("command is named %q", c.Name())
	}
	f := c.Flags().Lookup("write")
	if f == nil {
		t.Fatal("grit routes has no --write, so there is no narrow way to refresh the access table")
	}
	if f.Value.String() != "false" {
		t.Error("--write must be opt-in: listing routes should never rewrite a file")
	}
	if !strings.Contains(c.Long, "internal/access/registry.go") {
		t.Error("the help should name the file it writes")
	}
}
