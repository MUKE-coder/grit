package project

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// PluginInstalled reports whether a plugin is recorded in the project's
// lockfile.
//
// Reading the lockfile rather than looking for a file the plugin writes: the
// lockfile is what `grit plugin remove` maintains, so this answers "is it
// installed now" and not "was it ever".
//
// It lives here rather than in internal/generate because two packages need the
// same answer: the generator, deciding what to emit for a new resource, and
// the upgrade in internal/scaffold, deciding what to repair in one that
// already exists. internal/scaffold cannot import internal/generate, so the
// alternative was a second reader, and a second reader is how the two start
// disagreeing about whether a project is multitenant.
func PluginInstalled(root, name string) bool {
	data, err := os.ReadFile(filepath.Join(root, ".grit", "plugins.lock.json"))
	if err != nil {
		return false
	}
	var lock struct {
		Plugins []struct {
			Name string `json:"name"`
		} `json:"plugins"`
	}
	if json.Unmarshal(data, &lock) != nil {
		return false
	}
	for _, p := range lock.Plugins {
		if p.Name == name {
			return true
		}
	}
	return false
}
