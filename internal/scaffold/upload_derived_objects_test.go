package scaffold

import (
	"path/filepath"
	"strings"
	"testing"
)

// grit#221: an upload is more than one object, and only one of them was ever
// deleted.
//
// The media pipeline writes the stored file, one object per rendition the
// profile asked for, and the original under a private prefix for
// reprocessing. Four objects for an ordinary image. Every delete path in the
// framework worked from a single key:
//
//   - clearing or replacing a file field deleted the stored file and left the
//     renditions and the original, which is the issue as reported;
//   - DELETE /api/uploads/:id did the same, and derived a thumbnails/ key by
//     prefix swap, which is the async worker's scheme and not the pipeline's,
//     so for an optimised image it deleted a key that did not exist;
//   - the nightly orphan sweep walks uploads rows, and a rendition has none;
//   - and the right-to-erasure cleanup deleted each row's path column, so what
//     survived an erasure was the untouched original of the photo somebody
//     asked to have deleted.
//
// Nothing could find them afterwards either. There is no row to sweep, and the
// stored file's row is claimed for good once a record referenced it (grit#219).
//
// So the row records what it owns, and these are the two sides that have to
// agree: the handler that writes the column and the paths that read it. A test
// per side passes while the pair is broken, which is why both are here.
func TestAnUploadRecordsTheObjectsItOwns(t *testing.T) {
	handler := uploadHandlerGo()

	if !strings.Contains(handler, "Derived:      ref.DerivedKeys(),") {
		t.Error("the upload handler does not record what the pipeline derived, so " +
			"nothing downstream can delete the renditions or the original")
	}

	model := apiUploadModelGo()
	if !strings.Contains(model, "Derived      files.KeyList") {
		t.Error("the uploads model has no Derived column, so there is nowhere to " +
			"record them")
	}
	// Not in the API response. These are storage keys, and one of them names
	// the original, which is kept under a private prefix precisely because it
	// is not for serving.
	if !strings.Contains(model, `json:"-"`) {
		t.Error("Derived is serialised to JSON, which publishes the private " +
			"original's key in every upload response")
	}
}

func TestEveryDeletePathReadsWhatTheUploadOwns(t *testing.T) {
	for name, body := range map[string]string{
		"the uploads delete handler": uploadHandlerGo(),
		"the erasure cleanup":        apiErasableModelsGo(),
		"the orphan sweep":           filesLifecycleGo(),
	} {
		if !strings.Contains(body, "Derived") {
			t.Errorf("%s does not read the derived objects, so it deletes one of "+
				"the three or four objects an upload owns", name)
		}
	}

	// And the field-level cleanup, which has the refs rather than the row, so
	// it reads them off the ref instead.
	lifecycle := filesLifecycleGo()
	for _, want := range []string{"RemovedKeys(oldRef, newRef)", "RemovedKeysMulti(oldRefs, newRefs)"} {
		if !strings.Contains(lifecycle, want) {
			t.Errorf("CleanupRemoved does not call %s, so it is back to deleting "+
				"primary keys only", want)
		}
	}
	// The old helpers stay for projects that call them, but nothing in the
	// framework may delete by them again: that is the bug.
	for _, gone := range []string{"st.Delete(ctx, DiffSingle(", "range DiffMulti(oldRefs"} {
		if strings.Contains(lifecycle, gone) {
			t.Errorf("something deletes by %s, which returns primary keys only", gone)
		}
	}
}

// The two schemes for a derived key must not both be guessed at once.
//
// The async worker writes a thumbnail by swapping the uploads/ prefix for
// thumbnails/. The media pipeline appends -<name> to the stem instead. The
// delete handler knew only the first, which is why it deleted a nonexistent
// key for every optimised image. It may still use that rule for rows written
// before there was a column to record one, and must not use it for a row that
// records its own.
func TestTheLegacyThumbnailRuleOnlyAppliesToRowsThatRecordNothing(t *testing.T) {
	handler := uploadHandlerGo()
	guard := `if upload.ThumbnailURL != "" && len(upload.Derived) == 0 {`
	if !strings.Contains(handler, guard) {
		t.Error("the thumbnails/ prefix swap is applied unconditionally, so a row " +
			"that records its renditions also gets a key derived by the wrong rule")
	}

	// And the worker records the key it wrote, so new rows need no rule.
	workers := jobsWorkersGo()
	if !strings.Contains(workers, `updates["derived"] = files.KeyList{thumbKey}`) {
		t.Error("the image worker does not record the thumbnail it wrote, so " +
			"deleting it stays a guess")
	}
	// A map, not a struct: Upload's BeforeUpdate bumps version with SetColumn,
	// and a struct Save discards what such a hook set.
	if !strings.Contains(workers, "Updates(updates)") {
		t.Error("the worker writes the thumbnail with something other than a map " +
			"update, which discards the version bump its own hook sets")
	}
}

// The column may be absent, and naming it anyway fails the whole statement.
//
// grit upgrade rewrites internal/files and models/erasable.go and does not run
// grit migrate, so a project can take this code before it has the column. The
// sweep stopping and an erasure refusing to complete are both worse than the
// leak this closes, so both callers ask the schema first.
func TestTheDerivedColumnIsOptionalForEveryReader(t *testing.T) {
	if !strings.Contains(filesLifecycleGo(), "func UploadsRecordDerived(db *gorm.DB) bool") {
		t.Fatal("there is no way to ask whether the column exists")
	}
	for name, body := range map[string]string{
		"the orphan sweep":    filesLifecycleGo(),
		"the erasure cleanup": apiErasableModelsGo(),
		"the image worker":    jobsWorkersGo(),
	} {
		if !strings.Contains(body, "UploadsRecordDerived(") {
			t.Errorf("%s names the derived column unconditionally, so it fails "+
				"outright on a project that has not migrated yet", name)
		}
	}
}

// models/upload.go has to travel on upgrade, with its readers.
//
// internal/files and models/erasable.go both go out on upgrade and both refer
// to the column. Leaving the model behind gives an upgraded project a
// lifecycle package naming a field its own model does not declare, which is a
// build failure rather than a missing fix. writeFrameworkOwnedFiles is called
// from writeAPIFiles and from upgrade, so one entry serves both.
func TestTheUploadModelTravelsWithItsReaders(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: "triple"}
	apiRoot := opts.APIRoot(root)

	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatalf("writeFrameworkOwnedFiles: %v", err)
	}
	if !fileExists(filepath.Join(apiRoot, "internal", "models", "upload.go")) {
		t.Error("models/upload.go is not in the set that travels on upgrade, so an " +
			"upgraded project gets readers of a column its model does not declare")
	}
}
