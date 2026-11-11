package generate

import (
	"strings"
	"testing"
)

// grit#203, second half: an administrator who cannot read customer content.
//
// Until --confidential there was no way to say it. ADMIN was exempt from
// ownership scoping everywhere, so a product whose rows are medical records,
// legal files or coursework had an operator who could read all of them, and
// the only alternative was hand-editing framework files.

func confidentialDef() *ResourceDefinition {
	return &ResourceDefinition{
		Name: "MedicalRecord",
		Fields: []Field{
			{Name: "patient", Type: "belongs_to", RelatedModel: "User"},
			{Name: "notes", Type: "text"},
		},
		OwnedBy:      "patient",
		Confidential: true,
	}
}

// The service uses the scope that does not exempt an ADMIN.
//
// This is where the guarantee lives: every route, job, command and test goes
// through these methods, so a rule enforced here is enforced everywhere
// rather than on the routes somebody remembered.
func TestAConfidentialResourceScopesWithoutTheAdminExemption(t *testing.T) {
	g := &Generator{Definition: confidentialDef(), Module: "app"}
	src := g.serviceSource(g.Names())

	if strings.Contains(src, "authz.ScopeOwned(") {
		t.Error("a confidential resource still uses ScopeOwned, which lets every ADMIN through")
	}
	if !strings.Contains(src, `authz.ScopeConfidential(ctx, query, "patient_id", "medical_records.confidential")`) {
		t.Error("the list is not scoped confidentially")
	}
	if !strings.Contains(src, `authz.OwnsConfidential(ctx, &item, "medical_records.confidential")`) {
		t.Error("the single-row guard still exempts an ADMIN")
	}
	// The export is the list without pages. Leaving it on the old scope would
	// mean an admin who cannot read the rows can download them as a CSV,
	// which is the same leak wearing a hat.
	if n := strings.Count(src, "authz.ScopeConfidential("); n < 3 {
		t.Errorf("ScopeConfidential appears %d times; the list, the export and the bulk "+
			"scope all need it, or one of them is a way around the other two", n)
	}
}

// An ordinary owned resource is untouched. --confidential is opt-in, and an
// existing project must behave exactly as it did.
func TestAnOwnedResourceKeepsTheAdminExemption(t *testing.T) {
	def := confidentialDef()
	def.Confidential = false
	g := &Generator{Definition: def, Module: "app"}
	src := g.serviceSource(g.Names())

	if strings.Contains(src, "ScopeConfidential") {
		t.Error("a plain --owned-by resource gained confidential scoping, which changes " +
			"behaviour for every existing project")
	}
	if !strings.Contains(src, `authz.ScopeOwned(ctx, query, "patient_id")`) {
		t.Error("ordinary ownership scoping was lost")
	}
}

// The two refusals, both of which produce a resource nobody can read.
func TestConfidentialRefusesTheCombinationsThatCannotWork(t *testing.T) {
	// No owner column: the scope has nothing to compare against, so every
	// caller but a grant holder matches nothing. Never what anybody typed.
	def := confidentialDef()
	def.OwnedBy = ""
	g := &Generator{Definition: def, Module: "app"}
	err := g.checkFlagCombinations()
	if err == nil {
		t.Error("--confidential without --owned-by was accepted, producing a resource " +
			"readable by nobody")
	} else if !strings.Contains(err.Error(), "--owned-by") {
		t.Errorf("the error does not say what to add: %v", err)
	}

	// Public reads are anonymous and unscoped, so marking a resource
	// confidential and then publishing it hands every row to the internet.
	def2 := confidentialDef()
	def2.Public = true
	g2 := &Generator{Definition: def2, Module: "app"}
	if err := g2.checkFlagCombinations(); err == nil {
		t.Error("--confidential and --public were accepted together")
	}
}

// Confidential access has to be accountable: that is the whole reason it is a
// grant against a named person rather than a role. Implied rather than
// required, because somebody who typed --confidential and not --audit-reads
// did not mean "and keep no record".
func TestConfidentialImpliesAuditedReads(t *testing.T) {
	d := &ResourceDefinition{Name: "MedicalRecord"}
	d.ApplyFlags(false, false, "patient", false, false, false, true)
	if !d.AuditReads {
		t.Error("--confidential does not turn on read auditing, so a grant holder can read " +
			"everything and leave no trace, which defeats granting it to a named person")
	}
}

// The permission is per resource, so a project can hand out patient records
// without handing out payroll.
func TestTheConfidentialPermissionIsPerResource(t *testing.T) {
	if got := confidentialDef().ConfidentialPermission(); got != "medical_records.confidential" {
		t.Errorf("permission is %q, want medical_records.confidential", got)
	}
	other := &ResourceDefinition{Name: "PayrollEntry"}
	if got := other.ConfidentialPermission(); got == confidentialDef().ConfidentialPermission() {
		t.Error("two resources share one permission, so granting access to either grants both")
	}
}
