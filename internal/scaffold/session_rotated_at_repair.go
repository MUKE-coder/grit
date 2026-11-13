package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairSessionRotatedAt adds Session.RotatedAt to a project whose model
// predates it.
//
// The shape of the problem, and it is not specific to this field:
// internal/services/session.go is framework code and an upgrade rewrites it.
// internal/models is the project's own, so an upgrade leaves it alone, which
// is right for every model a project wrote and wrong for the handful the
// framework ships. When the service learned about RotatedAt in v3.464.0, a
// project upgrading into it got a service referencing a field its model does
// not declare, and the package stops compiling with an error naming the
// service rather than the model.
//
// So the field is carried across by a repair, which is the only thing allowed
// to edit a file in internal/models.
//
// Idempotent: a model that already declares it is returned untouched.
func repairSessionRotatedAt(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	model := filepath.Join(apiRoot, "internal", "models", "session.go")
	service := filepath.Join(apiRoot, "internal", "services", "session.go")

	if !fileExists(model) {
		return nil
	}
	// Only when the service actually needs it. A project that has not taken
	// the service yet is not helped by a field nothing reads, and adding one
	// is a migration it did not ask for.
	//
	// Compared with line endings normalised, because this file may be CRLF
	// on a Windows checkout and the needle is not.
	if !strings.Contains(normaliseNewlines(readFileOrEmpty(service)), "RotatedAt") {
		return nil
	}

	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	// The edit itself runs inside repairSourceFile, which hands over the
	// source with CRLF already normalised and puts it back afterwards.
	// Matching a \n anchor against the raw bytes instead is a repair that
	// silently does nothing on a Windows checkout, which is exactly what the
	// first version of this did.
	return repairSourceFile(root, m, model, func(src string) (string, []string, []string) {
		if strings.Contains(src, "RotatedAt") {
			return src, nil, nil
		}
		// Anchored on the declaration rather than the comment above it,
		// because the comment is prose somebody may have rewritten.
		anchor := "\tPrevTokenHash string `gorm:\"size:64;index\" json:\"-\"`\n"
		if !strings.Contains(src, anchor) {
			return src, nil, nil
		}
		field := anchor + "\n" +
			"\t// RotatedAt is when PrevTokenHash stopped being current. The grace window\n" +
			"\t// is measured from here and nowhere else, so presenting the old token\n" +
			"\t// again cannot slide the window forward.\n" +
			"\tRotatedAt *time.Time `json:\"-\"`\n"

		out := strings.Replace(src, anchor, field, 1)
		if !strings.Contains(out, `"time"`) {
			out = strings.Replace(out, "import (\n", "import (\n\t\"time\"\n", 1)
		}
		return out, []string{
			"the Session model gained RotatedAt, which internal/services/session.go needs to tell " +
				"a replayed refresh token from two tabs refreshing in the same second",
		}, nil
	})
}

// RepairFrameworkModels carries fields the framework added to the models it
// ships into a project whose copies predate them.
//
// One entry point rather than a call per field, because this is a class of
// problem rather than one bug: every time a framework service learns about a
// column, the model beside it needs the same change and cannot get it from
// the delivery map.
func RepairFrameworkModels() error {
	root, err := FindProjectRoot()
	if err != nil {
		return err
	}
	opts := Options{
		Architecture:     detectArchitecture(root),
		LegacySingleFlat: dirExists(filepath.Join(root, "internal", "routes")),
	}
	opts.Normalize()
	return repairSessionRotatedAt(root, opts)
}
