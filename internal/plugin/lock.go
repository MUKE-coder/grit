package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LockPath is where the record of installed plugins lives, relative to the
// project root. Commit it — it's how a teammate knows what's installed, and how
// removal knows what to undo.
const LockPath = ".grit/plugins.lock.json"

// Lock is the on-disk record of installed plugins.
type Lock struct {
	// Version of the lockfile format itself.
	Version int               `json:"version"`
	Plugins []InstalledPlugin `json:"plugins"`
}

// InstalledPlugin records exactly what one install did.
//
// This is the whole point of the design: removal reads this, not a
// hand-maintained list in the plugin's own code. A plugin author cannot forget
// to update their removal logic, because there isn't any.
type InstalledPlugin struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	InstalledAt time.Time `json:"installed_at"`

	// Requires is copied from the plugin at install time so removal can honour
	// dependencies even for a plugin this CLI doesn't have in its registry.
	Requires []string `json:"requires,omitempty"`

	// Files written by the install, relative to the project root.
	Files []string `json:"files"`

	// Injections records the exact text inserted, so removal can take out
	// precisely that and nothing else.
	Injections []LockedInjection `json:"injections"`

	GoDeps   []Dependency `json:"go_deps,omitempty"`
	NodeDeps []Dependency `json:"node_deps,omitempty"`
}

// LockedInjection is one applied edit.
type LockedInjection struct {
	File   string `json:"file"`
	Marker string `json:"marker"`
	// Code is the exact snippet inserted. Removal deletes this literal text.
	Code string `json:"code"`
}

// LoadLock reads the lockfile. A missing file is not an error — it just means
// no plugins are installed.
func LoadLock(root string) (*Lock, error) {
	path := filepath.Join(root, LockPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Lock{Version: 1}, nil
		}
		return nil, fmt.Errorf("reading %s: %w", LockPath, err)
	}

	var lock Lock
	if err := json.Unmarshal(data, &lock); err != nil {
		// Refuse to guess. Silently starting from an empty lock would orphan
		// every file the previous install wrote.
		return nil, fmt.Errorf("parsing %s (fix or delete it): %w", LockPath, err)
	}
	if lock.Version == 0 {
		lock.Version = 1
	}
	return &lock, nil
}

// Save writes the lockfile.
func (l *Lock) Save(root string) error {
	path := filepath.Join(root, LockPath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating %s: %w", filepath.Dir(path), err)
	}

	data, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

// Find returns the record for an installed plugin.
func (l *Lock) Find(name string) (*InstalledPlugin, bool) {
	for i := range l.Plugins {
		if l.Plugins[i].Name == name {
			return &l.Plugins[i], true
		}
	}
	return nil, false
}

// Add records an install, replacing any previous record for the same plugin.
func (l *Lock) Add(p InstalledPlugin) {
	for i := range l.Plugins {
		if l.Plugins[i].Name == p.Name {
			l.Plugins[i] = p
			return
		}
	}
	l.Plugins = append(l.Plugins, p)
}

// Remove drops a plugin's record.
func (l *Lock) Remove(name string) {
	out := l.Plugins[:0]
	for _, p := range l.Plugins {
		if p.Name != name {
			out = append(out, p)
		}
	}
	l.Plugins = out
}

// Installed lists installed plugin names.
func (l *Lock) Installed() []string {
	names := make([]string, 0, len(l.Plugins))
	for _, p := range l.Plugins {
		names = append(names, p.Name)
	}
	return names
}
