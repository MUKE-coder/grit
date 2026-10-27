package scaffold

import (
	"path/filepath"
)

// An upgrade must not put back what a removal took out.
//
// The blog is an example resource: every project is scaffolded with one, and
// `grit remove resource Blog` takes it out, all fourteen files. Three of those
// files are written from the scaffold's own maps rather than by the generator,
// so the next `grit upgrade` wrote them back:
//
//	apps/api/internal/database/blogs_seeder.go
//	apps/web/hooks/use-blogs.ts
//	apps/web/lib/blog-api.ts
//
// The seeder refers to models.Blog, which is gone, so the API stopped
// compiling. Nothing referred to the other two. The upgrade turned a building
// project into one that did not build, in files the developer had never
// touched and had deliberately deleted, and without git nothing said which
// command had done it (grit#111).
//
// Keyed on the model rather than on a record of the removal. A record would
// only know about removals `grit remove` made, and a model deleted by hand
// leaves exactly the same hole.

// blogRemoved reports whether this project had the blog example and took it
// out.
//
// Only meaningful while upgrading. On a new project the model has not been
// written yet when the seeder map is built, so asking there would answer "the
// blog is gone" for every project Grit creates and no project would ever get
// one.
func blogRemoved(root string, opts Options) bool {
	if !opts.Upgrading {
		return false
	}
	return !fileExists(filepath.Join(opts.APIRoot(root), "internal", "models", "blog.go"))
}

// withoutRemoved drops the files of a resource this project no longer has.
//
// Takes the map a writer was about to write, so every writer gets the same
// answer from one place: the three files above are written from three
// different maps, and a condition repeated at three call sites is a condition
// that will be right at two of them.
func withoutRemoved(files map[string]string, root string, opts Options) map[string]string {
	if !blogRemoved(root, opts) {
		return files
	}
	for path := range files {
		if isBlogExampleFile(path) {
			delete(files, path)
		}
	}
	return files
}

// isBlogExampleFile reports whether a path is part of the blog example.
//
// Matched on the file name, because the same file is written to a different
// place in each architecture: apps/web/hooks for a triple, frontend/src/hooks
// for a single, apps/expo/hooks for the mobile app.
func isBlogExampleFile(path string) bool {
	switch filepath.Base(path) {
	case "blogs_seeder.go", "use-blogs.ts", "blog-api.ts":
		return true
	}
	return false
}
