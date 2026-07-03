package scaffold

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

// gritLogoPNG is the Grit brand mark, embedded into the CLI so every scaffold
// ships real icon/splash/favicon assets instead of dangling references.
//
//go:embed assets/grit_logo.png
var gritLogoPNG []byte

// writeBytes writes binary content (embedded PNGs) to path, creating parent
// directories as needed. The binary sibling of writeFile.
func writeBytes(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", path, err)
	}
	return os.WriteFile(path, data, 0644)
}

// writeBrandLogo drops the Grit logo into dir under each of the given filenames.
// Used for app icons, splash screens, and web favicons across scaffolds.
func writeBrandLogo(dir string, names ...string) error {
	for _, name := range names {
		if err := writeBytes(filepath.Join(dir, name), gritLogoPNG); err != nil {
			return fmt.Errorf("writing brand asset %s: %w", name, err)
		}
	}
	return nil
}
