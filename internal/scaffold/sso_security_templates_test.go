package scaffold

import (
	"strings"
	"testing"
)

// A customer's identity provider could assert any address, and linking by
// email then signed it in as that account. Reproduced live against a mock IdP:
// it signed in as the app's ADMIN.
func TestSSOTrustsAnIdPOnlyForItsDomains(t *testing.T) {
	if !strings.Contains(apiSSOModelGo(), "func (s *SSOConnection) OwnsEmail(email string) bool") {
		t.Error("the connection cannot say which addresses it may vouch for")
	}
	h := apiSSOHandlerGo()
	for _, want := range []string{
		"if !conn.OwnsEmail(email) {",      // before linking by email or provisioning
		"if !conn.OwnsEmail(user.Email) {", // before honouring an existing link
	} {
		if !strings.Contains(h, want) {
			t.Errorf("resolveUser is missing %s", want)
		}
	}
	// The email check has to come before the lookup that links by email.
	if strings.Index(h, "if !conn.OwnsEmail(email) {") > strings.Index(h, `err = h.DB.Where("email = ?", email).First(&user).Error`) {
		t.Error("the domain check runs after the account is already found by email")
	}
}

func TestSSOLeavingEveryGroupRevokes(t *testing.T) {
	if !strings.Contains(apiSSOHandlerGo(), "names = []string{models.RoleUser}") {
		t.Error("leaving every mapped group does not revoke the mapped role")
	}
	if !strings.Contains(apiSSOTestGo(), "func TestSSO_LeavingEveryMappedGroupRevokes(") {
		t.Error("the shipped tests do not cover leaving every mapped group")
	}
}

// A partial PUT cleared the metadata URL and took a SAML connection offline.
func TestSSOUpdateKeepsWhatWasNotSent(t *testing.T) {
	h := apiSSOHandlerGo()
	for _, field := range []string{"metadata_url", "group_mappings", "discovery_url", "groups_attribute"} {
		if !strings.Contains(h, `setIfSent("`+field+`"`) {
			t.Errorf("an update that leaves out %s still clears it", field)
		}
	}
}

// With IdP-initiated sign-in off, every login was refused, because the ACS
// had no request ids to match a response against.
func TestSAMLMatchesResponsesToRequests(t *testing.T) {
	h := apiSAMLHandlerGo()
	for _, want := range []string{
		`c.SetCookie(samlRequestCookie, slug+":"+req.ID`,
		"sp.ParseResponse(c.Request, requestIDs)",
	} {
		if !strings.Contains(h, want) {
			t.Errorf("the SAML handler is missing %s", want)
		}
	}
	if strings.Contains(h, "sp.ParseResponse(c.Request, []string{})") {
		t.Error("the ACS still checks responses against no request ids")
	}
}

// The IdP-initiated toggle wrote a column that does not exist: GORM named it
// allow_id_p_initiated. And two connections could claim one domain.
func TestSSOToggleAndDomainOwnership(t *testing.T) {
	if !strings.Contains(apiSSOModelGo(), "column:allow_id_p_initiated") {
		t.Error("the IdP-initiated flag's column is left to GORM's guess")
	}
	h := apiSSOHandlerGo()
	if !strings.Contains(h, `updates["allow_id_p_initiated"]`) || strings.Contains(h, `updates["allow_idp_initiated"]`) {
		t.Error("the update still writes a column that does not exist")
	}
	if strings.Count(h, "h.domainTaken(in.Domains") != 2 {
		t.Error("create and update do not both refuse a domain another connection claims")
	}
}

func TestSAMLPassesOverTransientNameIDs(t *testing.T) {
	if !strings.Contains(apiSAMLServiceGo(), "id.Format != TransientNameID") {
		t.Error("a transient NameID is still used as the subject")
	}
	if !strings.Contains(apiSSOServiceGo(), `(protocol IS NULL OR protocol <> ?)`) {
		t.Error("the OIDC registry still tries to build SAML connections")
	}
}
