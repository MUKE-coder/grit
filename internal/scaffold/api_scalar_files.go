package scaffold

// gin-docs replaces the old Scalar-based API documentation.
// Documentation is now auto-generated from routes and GORM models
// via gindocs.Mount() in routes.go — no separate docs files needed.
//
// The writeScalarFiles function is kept for upgrade compatibility
// but is now a no-op.

func writeScalarFiles(root string, opts Options) error {
	// gin-docs is mounted directly in routes.go via gindocs.Mount()
	// No separate docs files are needed anymore.
	return nil
}
