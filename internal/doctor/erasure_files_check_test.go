package doctor

import (
	"strings"
	"testing"
)

// erasedInvoice is the clean Invoice with a file column and the erasure
// registration a resource generated with --owned-by carries. cleanup says
// whether that registration takes the files with the rows.
func erasedInvoice(cleanup bool) string {
	src := cleanProject()[invoiceModel]
	src = strings.Replace(src,
		"\tUserID string                 `gorm:\"size:36;index\" json:\"user_id\"`\n",
		"\tUserID string                 `gorm:\"size:36;index\" json:\"user_id\"`\n"+
			"\tAttachment *files.FileRef    `gorm:\"type:json\" json:\"attachment\"`\n", 1)
	if cleanup {
		return src + "\nfunc init() { erasure.RegisterWithCleanup(&Invoice{}, \"user_id\", eraseInvoiceFiles) }\n"
	}
	return src + "\nfunc init() { erasure.Register(&Invoice{}, \"user_id\") }\n"
}

// The case the check is for: the rows go with their owner and the files do
// not, so a file an administrator uploaded onto one of them is in storage for
// good with nothing pointing at it (grit#311).
func TestAnOwnedResourceWithFilesAndNoCleanupIsFound(t *testing.T) {
	report := run(t, map[string]string{invoiceModel: erasedInvoice(false)})
	found := fired(report, "owned-files-outlive-erasure")
	if len(found) != 1 {
		t.Fatalf("got %d findings, want 1: %v", len(found), found)
	}
	f := found[0]
	// An error, not a warning: the erasure is certified complete, and the
	// file is still there.
	if f.Level != "error" || f.Resource != "Invoice" {
		t.Errorf("got a %s on %s, want an error on Invoice", f.Level, f.Resource)
	}
	if !strings.Contains(f.Fix, "erasure.RegisterWithCleanup(&Invoice{}, \"user_id\"") ||
		!strings.Contains(f.Fix, "files.EraseOwnedFiles") {
		t.Errorf("the fix should name the edit and the function that does the work: %q", f.Fix)
	}
}

// A registration that carries the cleanup is the fixed state, and the two are
// told apart by one regex that must not match the other.
func TestACleanupAnswersTheCheck(t *testing.T) {
	if found := fired(run(t, map[string]string{invoiceModel: erasedInvoice(true)}),
		"owned-files-outlive-erasure"); len(found) > 0 {
		t.Errorf("a resource that already deletes its files was reported: %q", found[0].Message)
	}
}

// No file column, nothing in storage to leave behind.
func TestAnOwnedResourceWithoutFilesIsNotAsked(t *testing.T) {
	model := cleanProject()[invoiceModel] + "\nfunc init() { erasure.Register(&Invoice{}, \"user_id\") }\n"
	if found := fired(run(t, map[string]string{invoiceModel: model}),
		"owned-files-outlive-erasure"); len(found) > 0 {
		t.Errorf("a resource with no file column was asked about its files: %q", found[0].Message)
	}
}

// A resource whose rows are not erased at all is an older state and a bigger
// one. Saying "and its files too" about it would be a second complaint about
// the same event.
func TestAResourceThatIsNotErasedIsNotAsked(t *testing.T) {
	model := strings.Replace(erasedInvoice(false),
		"\nfunc init() { erasure.Register(&Invoice{}, \"user_id\") }\n", "", 1)
	if found := fired(run(t, map[string]string{invoiceModel: model}),
		"owned-files-outlive-erasure"); len(found) > 0 {
		t.Errorf("a resource with no erasure registration was reported: %q", found[0].Message)
	}
}
