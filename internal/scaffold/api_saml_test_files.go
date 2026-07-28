package scaffold

import "strings"

// apiSAMLTestGo emits internal/handlers/saml_test.go.
//
// These cover attribute extraction, which is where SAML integrations actually
// fail: every provider names its claims differently, so "your identity provider
// did not release an email address" is the error an operator hits when the
// fallbacks are wrong. Pinned names, the Okta-style short names and the Entra
// ID claim URIs are all exercised here.
func apiSAMLTestGo() string {
	src := `package handlers

import (
	"testing"

	"github.com/crewjam/saml"
	"github.com/stretchr/testify/assert"

	"{{MODULE}}/internal/models"
)

func attrStmt(pairs map[string][]string) []saml.AttributeStatement {
	attrs := []saml.Attribute{}
	for name, vals := range pairs {
		vs := []saml.AttributeValue{}
		for _, v := range vals {
			vs = append(vs, saml.AttributeValue{Value: v})
		}
		attrs = append(attrs, saml.Attribute{Name: name, Values: vs})
	}
	return []saml.AttributeStatement{{Attributes: attrs}}
}

func assertionWith(nameID string, pairs map[string][]string) *saml.Assertion {
	a := &saml.Assertion{AttributeStatements: attrStmt(pairs)}
	if nameID != "" {
		a.Subject = &saml.Subject{NameID: &saml.NameID{Value: nameID}}
	}
	return a
}

// With no attribute names pinned, the conventional ones must be found — this is
// what makes an Okta connection work with an empty config.
func TestSAMLIdentity_ConventionalAttributeNames(t *testing.T) {
	conn := &models.SSOConnection{Slug: "acme", Protocol: "saml"}
	a := assertionWith("okta-sub-1", map[string][]string{
		"email":     {"bob@acme.com"},
		"firstName": {"Bob"},
		"lastName":  {"Smith"},
		"groups":    {"it-admins", "everyone"},
	})

	id := samlIdentity(conn, a)
	assert.Equal(t, "okta-sub-1", id.Subject, "NameID is the subject")
	assert.Equal(t, "bob@acme.com", id.Email)
	assert.Equal(t, "Bob", id.FirstName)
	assert.Equal(t, "Smith", id.LastName)
	assert.ElementsMatch(t, []string{"it-admins", "everyone"}, id.Groups)
}

// Entra ID sends long Microsoft claim URIs rather than friendly names.
func TestSAMLIdentity_MicrosoftClaimURIs(t *testing.T) {
	conn := &models.SSOConnection{Slug: "acme", Protocol: "saml"}
	a := assertionWith("entra-sub", map[string][]string{
		"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress": {"sue@acme.com"},
		"http://schemas.xmlsoap.org/ws/2005/05/identity/claims/givenname":    {"Sue"},
		"http://schemas.microsoft.com/ws/2008/06/identity/claims/groups":     {"Admins"},
	})

	id := samlIdentity(conn, a)
	assert.Equal(t, "sue@acme.com", id.Email)
	assert.Equal(t, "Sue", id.FirstName)
	assert.ElementsMatch(t, []string{"Admins"}, id.Groups)
}

// A pinned attribute name must win over the conventional ones.
func TestSAMLIdentity_PinnedAttributeWins(t *testing.T) {
	conn := &models.SSOConnection{
		Slug: "acme", Protocol: "saml",
		EmailAttribute:  "urn:custom:workEmail",
		GroupsAttribute: "urn:custom:teams",
	}
	a := assertionWith("sub", map[string][]string{
		"email":                {"wrong@acme.com"},
		"urn:custom:workEmail": {"right@acme.com"},
		"groups":               {"ignored"},
		"urn:custom:teams":     {"platform"},
	})

	id := samlIdentity(conn, a)
	assert.Equal(t, "right@acme.com", id.Email, "the pinned attribute must win")
	assert.ElementsMatch(t, []string{"platform"}, id.Groups)
}

// Email is normalized so it matches an existing account regardless of casing.
func TestSAMLIdentity_EmailLowercased(t *testing.T) {
	conn := &models.SSOConnection{Slug: "acme", Protocol: "saml"}
	id := samlIdentity(conn, assertionWith("sub", map[string][]string{"email": {"Bob@ACME.com"}}))
	assert.Equal(t, "bob@acme.com", id.Email)
}

// Without a NameID the email stands in as the subject, so a transient-NameID
// IdP still works rather than failing every login.
func TestSAMLIdentity_FallsBackToEmailSubject(t *testing.T) {
	conn := &models.SSOConnection{Slug: "acme", Protocol: "saml"}
	id := samlIdentity(conn, assertionWith("", map[string][]string{"email": {"bob@acme.com"}}))
	assert.Equal(t, "bob@acme.com", id.Subject)
}

// An assertion carrying nothing usable must produce empty fields, which the
// handler turns into a clean "your IdP released no email" rather than a crash.
func TestSAMLIdentity_EmptyAssertionIsSafe(t *testing.T) {
	conn := &models.SSOConnection{Slug: "acme", Protocol: "saml"}
	id := samlIdentity(conn, assertionWith("", map[string][]string{}))
	assert.Equal(t, "", id.Email)
	assert.Equal(t, "", id.Subject)
	assert.Empty(t, id.Groups)
}
`
	return strings.ReplaceAll(src, "~", "`")
}
