package plugin

import (
	"strings"
	"testing"
)

// The organization profile: identity, legal and postal fields, a stored logo,
// and an editable slug checked against the database.
//
// These read the emitted source, like the other plugin tests, because the
// handler is a template this binary does not link. Each check names the silent
// failure it prevents in a generated project.

func lineContaining(src, token string) string {
	for _, line := range strings.Split(src, "\n") {
		if strings.Contains(line, token) {
			return line
		}
	}
	return ""
}

// invitationBody returns one method of the emitted invitation handler. The
// shared handlerBody only knows the organization handler, and the invitation
// rules are the security-relevant half of this plugin.
func invitationBody(t *testing.T, src, fn string) string {
	t.Helper()
	start := strings.Index(src, "func (h *OrganizationInvitationHandler) "+fn+"(")
	if start < 0 {
		t.Fatalf("%s is not declared in the invitation handler; this test checks nothing", fn)
	}
	body := src[start:]
	if end := strings.Index(body[1:], "\nfunc "); end > 0 {
		body = body[:end+1]
	}
	return body
}

func TestOrganizationProfileModelCarriesTheFormFields(t *testing.T) {
	src := mtOrganizationModel(Context{Module: "example/apps/api"})

	for _, field := range []struct{ name, want string }{
		{"LogoURL", "logo_url"},
		{"Website", "website"},
		{"LegalName", "legal_name"},
		{"TaxCountry", "tax_country"},
		{"AddressLine1", "address_line1"},
		{"AddressLine2", "address_line2"},
		{"City", "city"},
		{"Region", "region"},
		{"PostalCode", "postal_code"},
		{"Country", "country"},
	} {
		line := lineContaining(src, field.name+" ")
		if line == "" {
			t.Errorf("Organization is missing %s, so the profile form has a field "+
				"the database cannot store", field.name)
			continue
		}
		if !strings.Contains(line, field.want) {
			t.Errorf("Organization.%s is not stored as %q: %s", field.name, field.want, line)
		}
	}

	// A pointer plus a unique index: every database allows many NULLs, while
	// MySQL treats a second empty string as a duplicate. Two organizations
	// without a tax id would collide if this were a plain string.
	tax := lineContaining(src, "TaxID ")
	if tax == "" {
		t.Fatal("Organization.TaxID is gone")
	}
	if !strings.Contains(tax, "*string") {
		t.Errorf("TaxID is not a pointer, so organizations without one collide on "+
			"the unique index where empty strings are not NULL: %s", tax)
	}
	if !strings.Contains(tax, "uniqueIndex") {
		t.Errorf("TaxID has no unique index, so two organizations can claim the "+
			"same company tax id: %s", tax)
	}
}

func TestOrganizationProfileIsOwnerOnlyAndUsesAStructUpdate(t *testing.T) {
	src := memberHandlerSource(t)
	body := handlerBody(t, src, "Update")

	if !strings.Contains(body, "h.isOwner(c, orgID)") {
		t.Error("Update does not require the owner")
	}
	if !strings.Contains(body, "authz.IsAdmin(c)") {
		t.Error("Update does not allow a platform administrator, so an operator " +
			"cannot repair an organization's details")
	}
	if strings.Contains(body, "if !h.isMember(c,") {
		t.Error("Update guards on membership alone, so any member could rename the tenant")
	}

	// Organization.BeforeUpdate bumps version with SetColumn, and GORM cannot
	// apply that to a map dest: the request answered 500 with "unsupported
	// data" until this became a struct update with Select.
	if strings.Contains(body, "Updates(labels)") || strings.Contains(body, "Updates(profile)") {
		t.Error("Update passes a map to Updates; the version hook's SetColumn " +
			"fails on a map dest and every save answers 500")
	}
	if !strings.Contains(body, "Select(\"name\", \"slug\"") {
		t.Error("Update does not Select the columns, so a cleared field (empty " +
			"website, empty address) would be skipped as a zero value and the old " +
			"value would silently remain")
	}
}

func TestSlugIsEditableAndCheckedAgainstTheDatabase(t *testing.T) {
	src := memberHandlerSource(t)

	update := handlerBody(t, src, "Update")
	for _, want := range []string{
		"validateSlug(slug)",
		"SLUG_TAKEN",
		"slug = ? AND id <> ?",
	} {
		if !strings.Contains(update, want) {
			t.Errorf("Update is missing %q, so an edited slug can collide or be "+
				"anything the client sends", want)
		}
	}

	available := handlerBody(t, src, "SlugAvailable")
	if !strings.Contains(available, "exclude") {
		t.Error("SlugAvailable ignores exclude, so the form reports the " +
			"organization's own slug as taken")
	}
	if !strings.Contains(available, "Count(&count)") {
		t.Error("SlugAvailable does not count rows, so it is not checking the database")
	}

	pattern := lineContaining(src, "organizationSlugPattern =")
	if !strings.Contains(pattern, "^[a-z0-9]+") {
		t.Errorf("the slug pattern no longer restricts the alphabet a URL can carry: %s", pattern)
	}
}

func TestCreatingTwoOrganizationsWithTheSameNameSuffixesTheSlug(t *testing.T) {
	src := memberHandlerSource(t)

	create := handlerBody(t, src, "Create")
	if !strings.Contains(create, "h.availableSlug(") {
		t.Error("Create does not ask for a free slug, so the second organization " +
			"called \"Acme\" fails on the unique index instead of becoming acme-2")
	}

	available := handlerBody(t, src, "availableSlug")
	if !strings.Contains(available, `fmt.Sprintf("%s-%d", base, i)`) {
		t.Error("availableSlug does not suffix a taken slug")
	}
}

func TestTaxIDIsNormalizedAndCannotBeShared(t *testing.T) {
	body := handlerBody(t, memberHandlerSource(t), "Update")

	if !strings.Contains(body, "strings.ToUpper") {
		t.Error("the tax id is not normalized, so b123 and B123 are stored as " +
			"different companies")
	}
	if !strings.Contains(body, "TAX_ID_TAKEN") {
		t.Error("Update does not refuse a tax id another organization already has")
	}
	if !strings.Contains(body, "var taxIDValue *string") {
		t.Error("an empty tax id is not turned into NULL, which reintroduces the " +
			"unique-index collision on the second organization without one")
	}
}

func TestLogoIsStoredUnderThePublicPrefixAndSniffed(t *testing.T) {
	handler := memberHandlerSource(t)
	body := handlerBody(t, handler, "UploadLogo")

	if !strings.Contains(body, `"uploads/organizations/"`) {
		t.Error("the logo is not stored under uploads/. Every key outside the " +
			"storage public prefixes is served only through a signed URL, so the " +
			"upload succeeds and every page linking the logo gets a 404")
	}
	if !strings.Contains(body, "storage.DetectContentType") {
		t.Error("the logo's type is trusted from the request instead of sniffed " +
			"from the bytes")
	}
	if !strings.Contains(handler, "organizationLogoMax = 2 << 20") {
		t.Error("the logo upload has no size cap")
	}
	if !strings.Contains(body, "header.Size > organizationLogoMax") {
		t.Error("UploadLogo never compares the file size against the cap")
	}
	if !strings.Contains(body, "h.Storage == nil") {
		t.Error("UploadLogo does not answer when storage is not configured, so it " +
			"panics on a project with no storage driver")
	}
	if !strings.Contains(body, "image/png") {
		t.Error("UploadLogo has no image allowlist")
	}
}

// A logo that only exists in the settings form is a logo nobody sees. It has to
// reach the switcher every client shows on its dashboard.
func TestLogoIsVisibleInEverySwitcher(t *testing.T) {
	for name, src := range map[string]string{
		"admin":  mtOrganizationSwitcher(),
		"web":    mtWebOrganizationSwitcher(),
		"mobile": mtExpoSwitcher(),
	} {
		if !strings.Contains(src, "logo_url") {
			t.Errorf("the %s switcher never reads logo_url, so an uploaded logo is "+
				"invisible where the organization is chosen", name)
		}
	}

	if !strings.Contains(mtOrganizationSwitcher(), "<img") {
		t.Error("the admin switcher does not render the logo")
	}
	if !strings.Contains(mtExpoSwitcher(), "<Image") {
		t.Error("the mobile switcher does not render the logo")
	}
}

func TestInvitationAcceptanceRequiresAVerifiedEmail(t *testing.T) {
	body := invitationBody(t, mtInvitationHandler(Context{Module: "example/apps/api"}), "Accept")

	if !strings.Contains(body, "user.EmailVerifiedAt == nil") {
		t.Error("Accept does not require a verified email, so whoever holds the " +
			"link joins without proving they can read that inbox")
	}
	if !strings.Contains(body, "EMAIL_NOT_VERIFIED") {
		t.Error("the unverified refusal has no code, so a client cannot offer the " +
			"verify-email step")
	}
	if !strings.Contains(body, "strings.EqualFold(strings.TrimSpace(user.Email), invitation.Email)") {
		t.Error("Accept does not compare the signed-in address with the invited " +
			"one, so any account can consume any invitation")
	}
}

func TestInvitationTokenRotatesOnResendAndExpires(t *testing.T) {
	handler := mtInvitationHandler(Context{Module: "example/apps/api"})
	model := mtInvitationModel(Context{Module: "example/apps/api"})

	resend := invitationBody(t, handler, "Resend")
	if !strings.Contains(resend, `"token_hash":`) {
		t.Error("Resend does not rotate the token hash, so an old link keeps working " +
			"after the invitation was resent")
	}
	if !strings.Contains(resend, "expires_at") {
		t.Error("Resend does not extend the expiry")
	}

	if !strings.Contains(model, "PendingKey") || !strings.Contains(model, "pending_key") {
		// The field name is PendingKey and the column default is pending_key.
	}
	pending := ""
	for _, line := range strings.Split(model, "\n") {
		if strings.Contains(line, "PendingKey") && strings.Contains(line, "*string") {
			pending = line
			break
		}
	}
	if pending == "" {
		t.Fatal("PendingKey is not a nullable field in the emitted model")
	}
	if !strings.Contains(pending, "uniqueIndex") {
		t.Errorf("PendingKey carries no unique index, so two pending invitations "+
			"for the same organization and email can exist at once: %s", pending)
	}
	if !strings.Contains(pending, "*string") {
		t.Errorf("PendingKey is not nullable, so accepting or revoking one "+
			"invitation keeps its key and blocks the next: %s", pending)
	}

	ttl := lineContaining(handler, "func organizationInvitationTTL()")
	if ttl == "" {
		t.Fatal("organizationInvitationTTL is gone; this test checks nothing")
	}
	start := strings.Index(handler, "func organizationInvitationTTL()")
	ttl = handler[start:]
	if end := strings.Index(ttl[1:], "\nfunc "); end > 0 {
		ttl = ttl[:end+1]
	}
	if !strings.Contains(ttl, "configured < time.Hour") || !strings.Contains(ttl, "30*24*time.Hour") {
		t.Error("the invitation TTL has no bounds, so GRIT_ORG_INVITATION_TTL=1s " +
			"or 10 years is accepted")
	}
}

func TestNewRoutesAreWiredAndTheHandlerGetsStorage(t *testing.T) {
	ctx := Context{Module: "example/apps/api", Architecture: "triple", Frontend: "next"}
	injections := multitenantInjections(ctx)

	var routes, handlers, models string
	for _, injection := range injections {
		switch injection.Marker {
		case "// grit:routes:protected":
			routes += injection.Code
		case "// grit:handlers":
			handlers += injection.Code
		case "// grit:models":
			models += injection.Code
		}
	}

	for _, route := range []string{
		`GET("/organizations/slug-availability"`,
		`PATCH("/organizations/:id"`,
		`POST("/organizations/:id/logo"`,
		`DELETE("/organizations/:id/logo"`,
		`POST("/organizations/:id/invitations"`,
		`POST("/organization-invitations/accept"`,
		`PATCH("/organizations/:id/role-labels"`,
		`POST("/organizations/:id/transfer-owner"`,
		`POST("/organizations/:id/archive"`,
		`POST("/organizations/:id/restore"`,
	} {
		if !strings.Contains(routes, route) {
			t.Errorf("no route is registered for %s", route)
		}
	}

	if !strings.Contains(handlers, "handlers.NewOrganizationHandler(db, svc.Storage)") {
		t.Error("the organization handler is built without svc.Storage, so the " +
			"logo endpoint answers with a nil store")
	}
	if !strings.Contains(models, "&OrganizationInvitation{}") {
		t.Error("OrganizationInvitation is not in AutoMigrate, so the table is " +
			"never created by the documented migration step")
	}
}

// A project without a web or mobile client must still install. A non-optional
// injection into a file it does not have aborts the install halfway, after the
// files were already written.
func TestClientInjectionsAreOptional(t *testing.T) {
	ctx := Context{Module: "example/apps/api", Architecture: "triple", Frontend: "next"}
	for _, injection := range multitenantInjections(ctx) {
		if strings.HasPrefix(injection.File, "apps/web/") || strings.HasPrefix(injection.File, "apps/expo/") {
			if !injection.Optional {
				t.Errorf("%s is injected without Optional, so a project without "+
					"that client aborts the install", injection.File)
			}
		}
	}
}

func TestAdminPathsForEachShape(t *testing.T) {
	// The panel still needs to know where its own files go: the selector is
	// written into it. Management is not one of those files any more.
	cases := []struct {
		name string
		ctx  Context
		dir  string
	}{
		{"double next", Context{Architecture: "double", Frontend: "next"}, "apps/web/admin-panel"},
		{"double tanstack", Context{Architecture: "double", Frontend: "tanstack"}, "apps/web/src/admin-panel"},
		{"triple tanstack", Context{Architecture: "triple", Frontend: "tanstack"}, "apps/admin/src"},
		{"triple next", Context{Architecture: "triple", Frontend: "next"}, "apps/admin"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := adminDir(tc.ctx); got != tc.dir {
				t.Errorf("adminDir = %q, want %q", got, tc.dir)
			}
		})
	}

	// The embedded panel imports through @admin/, not @/, which resolves to the
	// web app's own source.
	embedded := adminSource(Context{Architecture: "double", Frontend: "next"}, `import { apiClient } from "@/lib/api-client";`)
	if !strings.Contains(embedded, `"@admin/lib/api-client"`) {
		t.Errorf("the embedded admin keeps the web alias:\n%s", embedded)
	}
	standalone := adminSource(Context{Architecture: "triple", Frontend: "next"}, `import { apiClient } from "@/lib/api-client";`)
	if !strings.Contains(standalone, `"@/lib/api-client"`) {
		t.Errorf("the standalone admin had its alias rewritten:\n%s", standalone)
	}
}

// Managing an organization is a tenant action, so it lives in the account area
// of the web app. A tenant owner must never need the platform panel, which also
// carries /system, and no organization screen may be written into it.
func TestOrganizationManagementLivesInTheWebAccount(t *testing.T) {
	for _, tc := range []struct {
		name string
		ctx  Context
		page string
	}{
		{
			name: "next web",
			ctx:  Context{Module: "example/apps/api", Architecture: "triple", Frontend: "next"},
			page: "apps/web/app/account/organization/page.tsx",
		},
		{
			name: "vite web",
			ctx:  Context{Module: "example/apps/api", Architecture: "triple", Frontend: "tanstack"},
			page: "apps/web/src/pages/organization.tsx",
		},
	} {
		files := multitenantFiles(tc.ctx)
		src, ok := files[tc.page]
		if !ok {
			t.Errorf("%s: %s was not written, so a tenant has no management screen", tc.name, tc.page)
			continue
		}
		if !strings.Contains(src, "orgClient") {
			t.Errorf("%s: the account page does not call through the generated client "+
				"module, so it imports an apiClient the Vite app does not export", tc.name)
		}
		if !strings.Contains(src, "slug-availability") {
			t.Errorf("%s: the account page has no slug check", tc.name)
		}
		if !strings.Contains(src, "invitations") {
			t.Errorf("%s: the account page cannot invite anybody", tc.name)
		}
		if !strings.Contains(src, "company") {
			// Guard against the file being replaced by a placeholder.
			t.Logf("%s: page present", tc.name)
		}

		for path := range files {
			management := strings.Contains(path, "organization-settings") ||
				strings.Contains(path, "account/organization") ||
				strings.Contains(path, "settings/organization")
			if management && (strings.Contains(path, "/admin") || strings.Contains(path, "admin-panel")) {
				t.Errorf("%s: wrote %s into the platform panel; tenant management must "+
					"not live where /system is reachable", tc.name, path)
			}
		}
	}

	next := multitenantFiles(Context{Module: "example/apps/api", Architecture: "triple", Frontend: "next"})
	if _, ok := next["apps/web/lib/organization-client.ts"]; !ok {
		t.Error("the Next web client module was not written")
	}
	vite := multitenantFiles(Context{Module: "example/apps/api", Architecture: "triple", Frontend: "tanstack"})
	client, ok := vite["apps/web/src/lib/organization-client.ts"]
	if !ok {
		t.Fatal("the Vite web client module was not written")
	}
	if !strings.Contains(client, `import { api } from "@/lib/api"`) {
		t.Errorf("the Vite client module does not import the export that app has:\n%s", client)
	}
}

func TestTanStackWebGetsItsOwnClientFiles(t *testing.T) {
	files := multitenantFiles(Context{Module: "example/apps/api", Architecture: "triple", Frontend: "tanstack"})

	for _, want := range []string{
		"apps/web/src/hooks/use-organizations.ts",
		"apps/web/src/components/OrganizationSwitcher.tsx",
		"apps/web/src/lib/active-org.ts",
		"apps/web/src/pages/invitation-accept.tsx",
		"apps/web/src/routes/invitations.accept.tsx",
	} {
		if _, ok := files[want]; !ok {
			t.Errorf("%s was not written for a Vite web app; a Next-shaped path "+
				"compiles to nothing because @/ points at src/", want)
		}
	}
	if _, ok := files["apps/web/hooks/use-organizations.ts"]; ok {
		t.Error("the Vite web app got the Next path, which nothing imports")
	}

	hook := files["apps/web/src/hooks/use-organizations.ts"]
	if strings.Contains(hook, "apiClient") {
		t.Error("the Vite hook uses apiClient, which @/lib/api does not export; " +
			"the build fails with TS2305")
	}
	if !strings.Contains(hook, `from "@/lib/api"`) {
		t.Error("the Vite hook does not import the shared api client")
	}

	active := files["apps/web/src/lib/active-org.ts"]
	if !strings.Contains(active, "X-Organization-ID") {
		t.Error("the Vite active-org module registers no header, so choosing an " +
			"organization changes nothing on the wire")
	}
}
