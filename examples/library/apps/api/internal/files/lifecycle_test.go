package files

import (
	"context"
	"testing"
)

// fakeStorage records every Delete call so tests can assert on which
// keys were purged. Implements the Storage interface.
type fakeStorage struct {
	deleted []string
}

func (f *fakeStorage) Delete(_ context.Context, key string) error {
	f.deleted = append(f.deleted, key)
	return nil
}

func TestDiffSingle(t *testing.T) {
	t.Run("nil old returns empty", func(t *testing.T) {
		if got := DiffSingle(nil, &FileRef{Key: "k"}); got != "" {
			t.Errorf("nil old should return empty, got %q", got)
		}
	})
	t.Run("same key returns empty", func(t *testing.T) {
		a := &FileRef{Key: "k"}
		b := &FileRef{Key: "k"}
		if got := DiffSingle(a, b); got != "" {
			t.Errorf("same key should return empty, got %q", got)
		}
	})
	t.Run("different key returns old", func(t *testing.T) {
		a := &FileRef{Key: "old"}
		b := &FileRef{Key: "new"}
		if got := DiffSingle(a, b); got != "old" {
			t.Errorf("different key should return old, got %q", got)
		}
	})
	t.Run("cleared (new nil) returns old", func(t *testing.T) {
		a := &FileRef{Key: "old"}
		if got := DiffSingle(a, nil); got != "old" {
			t.Errorf("cleared should return old, got %q", got)
		}
	})
}

func TestDiffMulti(t *testing.T) {
	old := FileRefs{{Key: "a"}, {Key: "b"}, {Key: "c"}}
	t.Run("shrunk to subset", func(t *testing.T) {
		removed := DiffMulti(old, FileRefs{{Key: "a"}, {Key: "c"}})
		if len(removed) != 1 || removed[0] != "b" {
			t.Errorf("expected [b], got %v", removed)
		}
	})
	t.Run("cleared entirely", func(t *testing.T) {
		removed := DiffMulti(old, FileRefs{})
		if len(removed) != 3 {
			t.Errorf("expected 3 removed, got %v", removed)
		}
	})
	t.Run("unchanged", func(t *testing.T) {
		removed := DiffMulti(old, old)
		if len(removed) != 0 {
			t.Errorf("expected 0 removed, got %v", removed)
		}
	})
	t.Run("nil old", func(t *testing.T) {
		if removed := DiffMulti(nil, FileRefs{{Key: "a"}}); len(removed) != 0 {
			t.Errorf("nil old should produce no removed keys, got %v", removed)
		}
	})
}

func TestCleanupRemovedReflectionSingle(t *testing.T) {
	type fakeProduct struct {
		ID    string
		Name  string
		Image *FileRef
	}
	old := fakeProduct{ID: "1", Image: &FileRef{Key: "old.jpg"}}
	neu := fakeProduct{ID: "1", Image: &FileRef{Key: "new.jpg"}}

	st := &fakeStorage{}
	CleanupRemoved(context.Background(), st, &old, &neu)
	if len(st.deleted) != 1 || st.deleted[0] != "old.jpg" {
		t.Errorf("expected old.jpg deleted, got %v", st.deleted)
	}
}

func TestCleanupRemovedReflectionMulti(t *testing.T) {
	type fakeGallery struct {
		ID     string
		Photos FileRefs
	}
	old := fakeGallery{ID: "1", Photos: FileRefs{{Key: "a"}, {Key: "b"}, {Key: "c"}}}
	neu := fakeGallery{ID: "1", Photos: FileRefs{{Key: "a"}, {Key: "c"}, {Key: "d"}}}

	st := &fakeStorage{}
	CleanupRemoved(context.Background(), st, &old, &neu)
	if len(st.deleted) != 1 || st.deleted[0] != "b" {
		t.Errorf("expected only b deleted, got %v", st.deleted)
	}
}

func TestCleanupRemovedNoFileFields(t *testing.T) {
	type plainRecord struct {
		ID   string
		Name string
	}
	old := plainRecord{ID: "1", Name: "before"}
	neu := plainRecord{ID: "1", Name: "after"}

	st := &fakeStorage{}
	CleanupRemoved(context.Background(), st, &old, &neu)
	if len(st.deleted) != 0 {
		t.Errorf("plain record should not trigger any deletes, got %v", st.deleted)
	}
}

func TestCleanupRemovedNilStorage(t *testing.T) {
	type record struct {
		Image *FileRef
	}
	old := record{Image: &FileRef{Key: "x"}}
	neu := record{Image: nil}
	// Must not panic when storage is nil -- e.g. when storage env vars
	// aren't configured.
	CleanupRemoved(context.Background(), nil, &old, &neu)
}
