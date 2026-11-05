package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Whose file is this, and what will `grit upgrade` do to it.
//
// The answer already existed twice: the manifest holds it, and guardAllows
// below acts on it. What there was no way to do was ask. An agent could, through
// grit_file_ownership over MCP; a person at a terminal read the manifest by hand
// or inferred it from the directory, and inferring it is how people get this
// wrong (grit#161).
//
// This mirrors guardAllows deliberately and completely, including the two cases
// that look like oversights and are not: a deleted file whose directory still
// exists is a request for the stock file back, and a project with no manifest at
// all is every project generated before v3.147.0, where there is no basis for
// claiming any file was edited.
//
// Mirrored, not shared, because guardAllows counts and records as it decides
// and this answers a question. TestOwnershipAgreesWithTheUpgradeGuard drives
// both over the same files and fails if they ever disagree, which is the only
// way a mirror stays one: a test per side would pass while the pair was broken.

// Ownership is what Grit knows about one file.
type Ownership struct {
	// Rel is the path as the manifest keys it: forward slashes, project-relative.
	Rel string
	// Status is the manifest's verdict.
	Status manifest.Status
	// Tracked is whether the PROJECT has a manifest, which is a different
	// question from whether this file has an entry in it.
	Tracked bool
	// WillReplace is what `grit upgrade` does, where that does not depend on
	// the bytes the new template would hold. False with Certain false means
	// "left alone unless the new template happens to match what is there".
	WillReplace bool
	// Certain is whether WillReplace can be known without the new template.
	Certain bool
	// Headline is the one-line answer.
	Headline string
	// Detail is what follows from it, including what to do instead.
	Detail string
}

// OwnershipOf answers for one path, absolute or relative to the working
// directory. root is the project root.
func OwnershipOf(root, path string) (Ownership, error) {
	abs := path
	if !filepath.IsAbs(abs) {
		var err error
		abs, err = filepath.Abs(path)
		if err != nil {
			return Ownership{}, err
		}
	}

	rel, inside := manifest.Rel(root, abs)
	if !inside {
		return Ownership{}, fmt.Errorf("%s is outside the project at %s", path, root)
	}

	m, err := manifest.Load(root)
	if err != nil {
		return Ownership{}, err
	}

	o := Ownership{
		Rel:     rel,
		Status:  m.StatusOf(root, rel),
		Tracked: manifest.Exists(root),
	}

	switch o.Status {
	case manifest.Unchanged:
		o.WillReplace, o.Certain = true, true
		o.Headline = "framework-owned, and unchanged since Grit wrote it"
		o.Detail = "grit upgrade will replace it. Keep your changes somewhere it will not:\n" +
			"a *.custom overlay where one exists, or a file of your own that imports this one."

	case manifest.Modified:
		// Certain is false: an edit that happens to leave the file identical to
		// the new template is adopted rather than reported, and that cannot be
		// known from here without rendering the template.
		o.WillReplace, o.Certain = false, false
		o.Headline = "framework-owned, and edited since Grit wrote it"
		o.Detail = "grit upgrade will leave it alone and name it in the report.\n" +
			"grit upgrade --diff shows what the new version would have changed."

	case manifest.Missing:
		// Read from the directory, as the guard does. A file gone from a
		// directory that is still there was deleted; a whole directory gone is
		// a feature the project removed.
		if dirExists(filepath.Dir(abs)) {
			o.WillReplace, o.Certain = true, true
			o.Headline = "framework-owned, and deleted"
			o.Detail = "grit upgrade will write it back, because its directory is still here.\n" +
				"To keep it gone, remove the feature rather than the file."
		} else {
			o.WillReplace, o.Certain = false, true
			o.Headline = "framework-owned, and removed along with its directory"
			o.Detail = "grit upgrade will leave it out: a missing directory reads as a\n" +
				"feature this project does not have, rather than a file to restore."
		}

	default: // manifest.Untracked
		if !o.Tracked {
			o.WillReplace, o.Certain = true, true
			o.Headline = "unknown, because this project has no manifest"
			o.Detail = "Generated before v3.147.0, so there is no record of what Grit wrote\n" +
				"and no basis for claiming any file was edited. grit upgrade writes them all.\n" +
				"Commit first, and read the diff afterwards."
		} else if _, err := os.Stat(abs); err != nil {
			o.WillReplace, o.Certain = true, true
			o.Headline = "not here, and not Grit's"
			o.Detail = "Nothing is recorded at this path and nothing is on disk. If a template\n" +
				"writes here, grit upgrade will create it."
		} else {
			o.WillReplace, o.Certain = false, false
			o.Headline = "yours"
			o.Detail = "Grit has no record of writing it: you created it, a command did not\n" +
				"record it, or it was copied in. grit upgrade leaves it alone, and names it\n" +
				"only if a template wants that path and the contents differ."
		}
	}

	return o, nil
}
