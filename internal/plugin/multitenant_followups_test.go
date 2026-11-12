package plugin

import (
	"strings"
	"testing"
)

// grit#199, the three items the reporter put first or near it. All three were
// found by exercising the v3.486.0 invitation flow on real generated projects,
// and none of them blocks ordinary use, which is most of why they needed
// writing down rather than noticing.

func mtHandlerSource(t *testing.T) string {
	t.Helper()
	src := mtHandler(Context{})
	if !strings.Contains(src, "func (h *OrganizationHandler) UploadLogo(") {
		t.Fatal("the organization handler no longer has UploadLogo, so this file is testing nothing")
	}
	return src
}

// Item 2. A replaced or removed logo left its object in the bucket for ever.
//
// The comment defending that said the storage driver owned its own cleanup.
// Nothing did: RunOrphanCleanup reads the uploads TABLE for rows nothing
// claimed, and a logo is uploaded straight to storage and never gets a row.
func TestASupersededLogoObjectIsDeleted(t *testing.T) {
	src := mtHandlerSource(t)

	for _, want := range []string{
		"func (h *OrganizationHandler) previousLogoKey(",
		"func (h *OrganizationHandler) dropLogoObject(",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("missing %s", want)
		}
	}
	// Both paths that move the column: replace and clear.
	if n := strings.Count(src, "h.previousLogoKey(orgID)"); n != 2 {
		t.Errorf("previousLogoKey is read %d times; UploadLogo and RemoveLogo are two, and the "+
			"one that is missing is a logo that leaks", n)
	}
	if n := strings.Count(src, "h.dropLogoObject(c,"); n != 2 {
		t.Errorf("dropLogoObject is called %d times, want 2", n)
	}
	// And the claim that was not true is gone.
	if strings.Contains(src, "storage driver owns its own cleanup") {
		t.Error("the handler still says the storage driver cleans these up, which nothing does")
	}
}

// The key is derived and then VERIFIED, which is what makes a string operation
// on a URL safe here: a project whose storage has moved declines rather than
// deleting something it guessed at.
func TestTheLogoKeyDerivationVerifiesItself(t *testing.T) {
	src := mtHandlerSource(t)
	body, ok := funcBodyOf(src, "func (h *OrganizationHandler) previousLogoKey(")
	if !ok {
		t.Fatal("previousLogoKey not found")
	}

	if !strings.Contains(body, "h.Storage.GetURL(key) != current.LogoURL") {
		t.Error("the derived key is not rendered back through GetURL and compared, so a storage " +
			"move would have it deleting an object it guessed the name of")
	}
	if !strings.Contains(body, `strings.IndexAny(key, "?#")`) {
		t.Error("a query or fragment is not trimmed, so a signed URL yields a key with the " +
			"signature stuck to it and the verification silently declines every time")
	}
	// Scoped to this organization's own prefix: a URL from somewhere else is
	// not something this handler may delete.
	if !strings.Contains(body, `"uploads/organizations/" + orgID + "/"`) {
		t.Error("the derivation is not anchored on the organization's own prefix")
	}
}

// Deleted after the column moves, never before: an object removed ahead of a
// write that then fails leaves the organization pointing at nothing.
func TestTheLogoObjectIsDeletedAfterTheColumnMoves(t *testing.T) {
	src := mtHandlerSource(t)
	upload, ok := funcBodyOf(src, "func (h *OrganizationHandler) UploadLogo(")
	if !ok {
		t.Fatal("UploadLogo not found")
	}
	write := strings.Index(upload, `Select("logo_url")`)
	drop := strings.Index(upload, "h.dropLogoObject(c,")
	if write < 0 || drop < 0 || drop < write {
		t.Errorf("the delete is not after the column update (update at %d, delete at %d)", write, drop)
	}
}

// Item 3. Slug and tax id were checked with a SELECT and then written, so a
// concurrent write hit the unique index and answered 500 INTERNAL_ERROR. The
// index is the truth and the check is a courtesy, so the violation has to
// answer what the check would have.
func TestASlugOrTaxIDRaceAnswersAConflict(t *testing.T) {
	src := mtHandlerSource(t)

	if !strings.Contains(src, "respond.DuplicateKey(c, err)") {
		t.Fatal("no unique-violation check, so a race is still a 500")
	}
	if n := strings.Count(src, "respond.DuplicateKey(c, err)"); n != 2 {
		t.Errorf("DuplicateKey is consulted %d times; Create and Update are the two write paths "+
			"that check first, and the one that is missing still answers 500", n)
	}
	// Update can be racing either column, and the driver names which.
	update, ok := funcBodyOf(src, "func (h *OrganizationHandler) Update(")
	if !ok {
		t.Fatal("Update not found")
	}
	if !strings.Contains(update, `column == "tax_id"`) {
		t.Error("Update answers the same code whichever column collided, so a tax id race is " +
			"reported as a slug that is taken")
	}
	// The module placeholder is already substituted by the time mtHandler
	// returns, so match the suffix rather than the template.
	if !strings.Contains(src, `/internal/respond"`) {
		t.Error("the respond package is not imported, so the generated file does not compile")
	}
}

// funcBodyOf is the text from a function's signature to the next top-level
// closing brace. Text, because the subject is a generated template string and
// go/parser would need the {{MODULE}} placeholders resolved first.
func funcBodyOf(src, signature string) (string, bool) {
	at := strings.Index(src, signature)
	if at < 0 {
		return "", false
	}
	rest := src[at:]
	if end := strings.Index(rest, "\n}\n"); end >= 0 {
		return rest[:end], true
	}
	return rest, true
}
