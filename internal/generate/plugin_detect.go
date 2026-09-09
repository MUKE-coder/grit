package generate

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
