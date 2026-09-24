package scaffold

import (
	"strings"
	"testing"
)

// Signing in with an emailed link.
//
// Three properties, each of them a way this feature is usually got wrong: the
// request form must not report who has an account, the token must not be spent
// by the GET that opens the link, and a link must not replace the second
// factor.

func TestAskingForALinkTellsNobodyWhoHasAnAccount(t *testing.T) {
	handler := magicLinkHandlerGo("example.com/app")

	req := handler[strings.Index(handler, "func (h *AuthHandler) RequestMagicLink("):]
	req = req[:strings.Index(req, "type MagicLinkConsumeRequest")]

	if strings.Count(req, "sameAnswer") < 3 {
		t.Error("the unknown-address and the known-address paths do not end at the same answer")
	}
	for _, leak := range []string{"no account", "not found", "No user", "unknown address"} {
		if strings.Contains(req, leak) {
			t.Errorf("the request form says %q, which turns it into a way to find out who is registered", leak)
		}
	}
	// A deployment that cannot send is the one exception, and it has to be one:
	// "on its way" would leave somebody watching an inbox forever.
	if !strings.Contains(req, "respond.CodeMailFailed") {
		t.Error("a deployment with no mailer still promises a link it cannot send")
	}
	// Every other auth lookup matches the address exactly as typed. Normalising
	// only here means an account that signs in with a password cannot be found
	// by a link, and since the form answers the same either way, nobody would
	// ever see why.
	if strings.Contains(req, "strings.ToLower(") {
		t.Error("the link lookup normalises the address differently from the password flow")
	}
	if !strings.Contains(req, "services.ErrMagicLinkTooSoon") {
		t.Error("a second request inside the rate window is silently swallowed, so the person waits for an email that is not coming")
	}
}

func TestALinkIsSpentByThePageNotByTheMailScanner(t *testing.T) {
	routes := apiRoutesGo()
	if !strings.Contains(routes, `auth.POST("/magic-link/consume"`) {
		t.Error("consuming a link is not a POST, so a mail scanner following the URL spends the token first")
	}
	if strings.Contains(routes, `auth.GET("/magic-link`) {
		t.Error("a GET route spends the token on whatever opens the link, scanner included")
	}

	page := adminMagicLinkPage()
	if !strings.Contains(page, `apiClient
      .post("/api/auth/magic-link/consume"`) {
		t.Error("the landing page does not spend the token itself")
	}
	if !strings.Contains(page, "spent.current") {
		t.Error("nothing stops the landing page consuming the token twice, which reports the second attempt as used")
	}
}

func TestALinkReplacesThePasswordNotTheSecondFactor(t *testing.T) {
	handler := magicLinkHandlerGo("example.com/app")

	consume := handler[strings.Index(handler, "func (h *AuthHandler) ConsumeMagicLink("):]
	if !strings.Contains(consume, "h.startTOTPChallenge(c, &user)") {
		t.Error("a link signs in past two-factor, so a mailbox alone is enough to get in")
	}
	if !strings.Contains(consume, "services.CreateSession(") {
		t.Error("signing in by link records no session, so the device never appears in Devices and cannot be revoked")
	}
	if !strings.Contains(consume, "services.LogLogin(") {
		t.Error("signing in by link is not logged, so the activity trail has a hole in exactly the shape of an attack")
	}
	if !strings.Contains(consume, "!user.Active") {
		t.Error("a disabled account can still sign in with a link issued before it was disabled")
	}
}

func TestAMagicLinkTokenIsStoredHashedAndSpentOnce(t *testing.T) {
	svc := magicLinkServiceGo("example.com/app")

	if strings.Contains(svc, "Token string `gorm") {
		t.Error("the raw token is stored, so a database read is a set of working sign-in links")
	}
	if !strings.Contains(svc, "sha256.Sum256") {
		t.Error("the token is not hashed")
	}
	// Spending has to be a conditional update, not read-then-write: two tabs
	// opening the same link at once would otherwise both succeed.
	if !strings.Contains(svc, "used_at IS NULL") {
		t.Error("spending a link is not conditional on it being unspent, so a racing pair of requests both consume it")
	}
	if !strings.Contains(svc, "RowsAffected") {
		t.Error("the conditional update's result is not checked, so a losing race reports success")
	}
}

func TestSignInOffersALinkAndTheRouteExists(t *testing.T) {
	login := adminThemedLoginPage()
	if !strings.Contains(login, "Email me a sign-in link") {
		t.Error("the sign-in page does not offer a link, which is the only way in for somebody who never set a password")
	}
	if !strings.Contains(login, `apiClient.post("/api/auth/magic-link"`) {
		t.Error("the button posts nowhere")
	}

	routes := apiRoutesGo()
	if !strings.Contains(routes, "authHandler.RequestMagicLink") {
		t.Error("the button's endpoint is not mounted, so it answers 404")
	}

	// And the table has to exist, or the first request is "no such table".
	if !strings.Contains(apiUserModelGo(), "&MagicLinkToken{}") {
		t.Error("MagicLinkToken is not in the model registry, so its table is never created")
	}
}

// An upgraded project gets routes.go left alone, so the repair is what mounts
// these. Without it the admin's button 404s with nothing to explain why.
func TestAnUpgradedProjectGetsTheLinkRoutes(t *testing.T) {
	before := `	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/verify-email", authHandler.VerifyEmail)
	}

	protected := v1.Group("")
	{
		protected.GET("/auth/sessions", sessionHandler.List)
	}
`
	after, changes, warnings := repairMagicLinkRoutesSource(before)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(changes) == 0 {
		t.Fatal("the repair reported no change")
	}
	if !strings.Contains(after, "authHandler.RequestMagicLink") || !strings.Contains(after, "authHandler.ConsumeMagicLink") {
		t.Error("the repair did not mount both routes")
	}
	if !strings.Contains(after, "authHandler.RecentMagicLinks") {
		t.Error("the repair did not mount the activity route, so the Sign-in links card stays empty")
	}

	again, changes, _ := repairMagicLinkRoutesSource(after)
	if len(changes) != 0 || again != after {
		t.Error("the repair is not idempotent, so a second upgrade mounts the routes twice")
	}
}

// The account screen's Security tab is where a request nobody made becomes
// visible. It has to be: the request form answers the same to every address,
// so it can never warn anybody itself.
func TestTheSecurityTabShowsWhoAskedForALink(t *testing.T) {
	page := adminAccountPageTSX()
	if !strings.Contains(page, "<SignInLinksCard />") {
		t.Error("the Security tab does not show sign-in link activity")
	}
	if !strings.Contains(page, `from "@/components/account/sign-in-links"`) {
		t.Error("the card is rendered but never imported")
	}

	card := adminSignInLinksCardTSX()
	if !strings.Contains(card, "/api/auth/magic-link/recent") {
		t.Error("the card reads from nowhere")
	}
	for _, leak := range []string{"token", "token_hash", "TokenHash"} {
		if strings.Contains(card, leak) {
			t.Errorf("the card mentions %q: this list must never carry a usable link", leak)
		}
	}

	routes := apiRoutesGo()
	if !strings.Contains(routes, "authHandler.RecentMagicLinks") {
		t.Error("the activity endpoint is not mounted")
	}
	// Behind a session, not public: it says which addresses have been asking
	// about this account.
	if strings.Contains(routes, `auth.GET("/magic-link/recent"`) {
		t.Error("the activity endpoint is public, so anyone can ask who has been signing in")
	}
}
