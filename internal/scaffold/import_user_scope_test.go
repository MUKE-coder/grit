package scaffold

import (
	"strings"
	"testing"
)

// grit#139: the generated importer resolved an owner email against every user
// on the platform. In a multitenant project that is a cross-tenant read
// through an importer: an import running in organization A could file rows
// under somebody who belongs only to organization B, and the row error said
// whether that address exists at all. The shape GHSA-4j45-98gm-mv9v is about.

// unscopedImporter is the importer as the generator wrote it before v3.445.0,
// for an --owned-by resource.
const unscopedImporter = `package services

func (s *DealService) ImportCSV(ctx context.Context, jobID, path string) {
	db := s.db(ctx)
	for {
		if v, ok := get(rec, "owner"); ok && v != "" && canAssignOwner {
			var rel models.User
			if err := db.Where("email = ?", v).First(&rel).Error; err != nil {
				message := fmt.Sprintf("no user with email %q", v)
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					message = err.Error()
				}
				failed++
				if len(rowErrors) < 50 {
					rowErrors = append(rowErrors, map[string]interface{}{"row": rowNum, "message": message})
				}
				continue
			}
			item.OwnerID = rel.ID
		}
	}
}
`

func TestTheUpgradeScopesAnUnscopedUserLookup(t *testing.T) {
	out := unscopedImporter
	n := scopeImportUserLookups(&out, "example/apps/api")
	if n != 1 {
		t.Fatalf("scoped %d lookups, want 1", n)
	}

	for _, want := range []string{
		"tenant.FromContext(ctx)",
		"Model(&models.OrganizationMember{})",
		`Where("organization_id = ?", orgID)`,
		// The condition the resource already had must survive: without it an
		// --owned-by import would let anybody name another owner.
		`ok && v != "" && canAssignOwner {`,
		"item.OwnerID = rel.ID",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the scoped lookup is missing %q", want)
		}
	}
	if strings.Contains(out, `if err := db.Where("email = ?", v).First(&rel)`) {
		t.Error("the unscoped lookup survived, so the cross-tenant read is still there")
	}

	// Upgrading twice must not nest the scoping.
	again := out
	if n := scopeImportUserLookups(&again, "example/apps/api"); n != 0 || again != out {
		t.Errorf("a second upgrade changed it again (%d)", n)
	}
}

// An importer somebody rewrote by hand no longer matches the generated text,
// and is left alone rather than half-rewritten.
func TestTheUpgradeLeavesAHandWrittenLookupAlone(t *testing.T) {
	const mine = `package services

func (s *DealService) ImportCSV(ctx context.Context, jobID, path string) {
	if v, ok := get(rec, "owner"); ok && v != "" {
		var rel models.User
		// my own lookup, by username
		if err := s.findUser(ctx, v, &rel); err != nil {
			continue
		}
		item.OwnerID = rel.ID
	}
}
`
	out := mine
	if n := scopeImportUserLookups(&out, "example/apps/api"); n != 0 || out != mine {
		t.Errorf("a hand-written importer was rewritten (%d)", n)
	}
}

// And the generator's own output, scoped and unscoped, comes from this one
// function, so a project generated today and one upgraded today cannot differ.
func TestImportUserLookupScopesOnlyWhenAsked(t *testing.T) {
	plain := ImportUserLookup("owner", "email", "OwnerID", "", false, false)
	if strings.Contains(plain, "tenant.FromContext") {
		t.Error("a project without the multitenant plugin got the scoping, which " +
			"references a models.OrganizationMember it does not have")
	}

	scoped := ImportUserLookup("owner", "email", "OwnerID", "", false, true)
	if !strings.Contains(scoped, "tenant.FromContext(ctx)") {
		t.Error("the multitenant form is not scoped")
	}
	// The message must not say "in this organization" when nothing scoped the
	// query, or it tells somebody importing platform-wide the wrong thing.
	if !strings.Contains(scoped, `"no user with email %q%s", v, inOrg`) {
		t.Error("the message does not follow whether the query was actually scoped")
	}
}
