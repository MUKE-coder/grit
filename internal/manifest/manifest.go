// Package manifest records which Grit generator wrote which file, and what it
// wrote, so a later command can tell an untouched generated file from one you
// have edited since.
//
// Without this, grit upgrade has only two options and both are wrong. It can
// overwrite everything, which is what it did before v3.147.0 and which silently
// discards your edits. Or it can overwrite nothing, which means no project ever
// receives a framework fix. The manifest is what makes the third option
// possible: overwrite the files nobody has touched, and leave the rest alone
// with a diff to look at.
package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FormatVersion is the manifest's own schema version, not Grit's. It changes
// only when the shape of this file changes in a way older CLIs cannot read.
const FormatVersion = 1

// Dir is the project-local directory holding Grit's own state.
const Dir = ".grit"

// FileName is the manifest inside Dir.
const FileName = "manifest.json"

// Manifest is the whole file.
type Manifest struct {
	Format    int              `json:"format"`
	Grit      string           `json:"grit"`
	UpdatedAt string           `json:"updatedAt"`
	Files     map[string]Entry `json:"files"`
}

// Entry is one generated file's provenance.
type Entry struct {
	// Generator is who wrote it: "scaffold", "resource:Product",
	// "plugin:webhooks". Kept as a free string so a new generator does not need
	// a change here to start recording.
	Generator string `json:"generator"`
	// Grit is the version that last wrote this file, which is not necessarily
	// the version that last wrote the manifest.
	Grit string `json:"grit"`
	// Hash is of the normalised content. See Hash.
	Hash string `json:"hash"`
}

// Status describes how a file on disk relates to what Grit last wrote there.
type Status int

const (
	// Untracked means Grit has no record of writing it: either you created it,
	// or it predates the manifest.
	Untracked Status = iota
	// Missing means Grit wrote it and it is not there any more.
	Missing
	// Unchanged means the bytes on disk are the bytes Grit wrote.
	Unchanged
	// Modified means it has been edited since Grit wrote it.
	Modified
)

func (s Status) String() string {
	switch s {
	case Missing:
		return "missing"
	case Unchanged:
		return "unchanged"
	case Modified:
		return "modified"
	default:
		return "untracked"
	}
}

// Path is where the manifest lives for a project root.
func Path(root string) string {
	return filepath.Join(root, Dir, FileName)
}

// Load reads a project's manifest. A project without one is not an error: it
// is every project generated before v3.147.0, and it gets an empty manifest in
// which every file reads as Untracked.
func Load(root string) (*Manifest, error) {
	m := &Manifest{Format: FormatVersion, Files: map[string]Entry{}}

	data, err := os.ReadFile(Path(root))
	if os.IsNotExist(err) {
		return m, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}

	if err := json.Unmarshal(data, m); err != nil {
		// A corrupt manifest must not stop the command that read it. Losing
		// provenance costs a needlessly cautious upgrade; refusing to run costs
		// the user their upgrade entirely.
		return &Manifest{Format: FormatVersion, Files: map[string]Entry{}}, nil
	}
	if m.Files == nil {
		m.Files = map[string]Entry{}
	}
	return m, nil
}

// Save writes the manifest, creating .grit/ if needed. Keys are sorted by
// encoding/json's map handling, so the file is stable across runs and a diff of
// it shows only what actually changed.
func (m *Manifest) Save(root, gritVersion string) error {
	m.Format = FormatVersion
	m.Grit = gritVersion
	m.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	dir := filepath.Join(root, Dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding manifest: %w", err)
	}
	if err := os.WriteFile(Path(root), append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}
	return nil
}

// Record notes that generator wrote content to rel.
func (m *Manifest) Record(rel, generator, gritVersion, content string) {
	if m.Files == nil {
		m.Files = map[string]Entry{}
	}
	m.Files[Key(rel)] = Entry{
		Generator: generator,
		Grit:      gritVersion,
		Hash:      Hash(content),
	}
}

// Forget drops a file from the manifest, for when a generator deletes what it
// previously wrote.
func (m *Manifest) Forget(rel string) {
	delete(m.Files, Key(rel))
}

// Entries returns the tracked paths in a stable order.
func (m *Manifest) Entries() []string {
	paths := make([]string, 0, len(m.Files))
	for p := range m.Files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}

// StatusOf compares one file on disk against what Grit last wrote there.
func (m *Manifest) StatusOf(root, rel string) Status {
	entry, tracked := m.Files[Key(rel)]
	if !tracked {
		return Untracked
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(Key(rel))))
	if err != nil {
		return Missing
	}
	if Hash(string(data)) == entry.Hash {
		return Unchanged
	}
	return Modified
}

// Key normalises a project-relative path into the form used as a map key:
// forward slashes, no leading "./".
func Key(rel string) string {
	k := filepath.ToSlash(rel)
	k = strings.TrimPrefix(k, "./")
	return k
}

// Rel turns an absolute path into a manifest key, reporting whether the path is
// inside the project at all. A generator writing outside the root (a global
// config, a temp file) is simply not recorded.
func Rel(root, abs string) (string, bool) {
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return "", false
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return Key(rel), true
}

// Hash is the SHA-256 of the content with line endings normalised.
//
// Normalising is not cosmetic. Git on Windows checks files out with CRLF while
// the generator writes LF, so hashing the raw bytes would report every file in
// a freshly cloned project as edited, and an upgrade that trusted that would
// refuse to update anything.
func Hash(content string) string {
	normalised := strings.ReplaceAll(content, "\r\n", "\n")
	sum := sha256.Sum256([]byte(normalised))
	return "sha256:" + hex.EncodeToString(sum[:])
}
