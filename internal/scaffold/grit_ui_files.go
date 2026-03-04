package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeGritUIFiles writes the packages/grit-ui package — a shadcn-compatible
// component registry with 100 components across marketing, saas, ecommerce,
// auth, and layout categories.
func writeGritUIFiles(root string, opts Options) error {
	uiRoot := filepath.Join(root, "packages", "grit-ui")

	files := map[string]string{
		filepath.Join(uiRoot, "package.json"):  gritUIPackageJSON(opts),
		filepath.Join(uiRoot, "registry.json"): gritUIRegistryJSON(opts),
	}

	// All 100 components: name → {title, description, category, deps, tsxFn}
	type compInfo struct {
		title, description, category string
		deps                         []string
		tsx                          string
	}

	allComponents := map[string]compInfo{
		// ── Original 4 ────────────────────────────────────────────────────────
		"hero-01":          {"Hero Section", "Bold hero with gradient text, announcement badge, and dual CTA buttons", "marketing", nil, ""},
		"pricing-card-01":  {"Pricing Card", "Pricing tier card with feature list, popular badge, and gradient CTA", "marketing", []string{"lucide-react"}, ""},
		"product-card-01":  {"Product Card", "E-commerce product card with image, star rating, and add-to-cart", "ecommerce", []string{"lucide-react"}, ""},
		"feature-grid-01":  {"Feature Grid", "Feature showcase section with icon cards in a responsive grid", "marketing", []string{"lucide-react"}, ""},

		// ── Marketing (10 new) ─────────────────────────────────────────────────
		"hero-split-01":       {"Split Hero", "Side-by-side hero with left copy + right terminal code block (Vercel-style)", "marketing", nil, gritUIHeroSplit01TSX()},
		"logo-cloud-01":       {"Logo Cloud", "Infinite-scroll marquee of customer/partner logos", "marketing", nil, gritUILogoCLoud01TSX()},
		"testimonial-card-01": {"Testimonial Card", "Single customer testimonial with avatar, stars, and quote", "marketing", nil, gritUITestimonialCard01TSX()},
		"testimonial-grid-01": {"Testimonial Grid", "3-column masonry grid of customer testimonials", "marketing", nil, gritUITestimonialGrid01TSX()},
		"stats-banner-01":     {"Stats Banner", "Dark stats strip with large numbers and radial glow", "marketing", nil, gritUIStatsBanner01TSX()},
		"faq-accordion-01":    {"FAQ Accordion", "Click-to-expand FAQ accordion with animated chevron", "marketing", nil, gritUIFaqAccordion01TSX()},
		"cta-banner-01":       {"CTA Banner", "Full-width CTA with email capture and success state", "marketing", nil, gritUICtaBanner01TSX()},
		"newsletter-signup-01": {"Newsletter Signup", "Newsletter card with subscriber count and success state", "marketing", nil, gritUINewsletterSignup01TSX()},
		"changelog-item-01":   {"Changelog Item", "Version changelog entry with timeline, tags, and change list", "marketing", nil, gritUIChangelogItem01TSX()},
		"pricing-toggle-01":   {"Pricing Toggle", "Full pricing section with monthly/annual toggle and 3 tier cards", "marketing", nil, gritUIPricingToggle01TSX()},

		// ── Auth (6 from agent + 4 extra) ──────────────────────────────────────
		"login-card-01":           {"Login Card", "Premium login card with social OAuth buttons and remember me", "auth", nil, gritUILoginCard01TSX()},
		"signup-card-01":          {"Signup Card", "Split signup form with feature list and password strength", "auth", nil, gritUISignupCard01TSX()},
		"forgot-password-card-01": {"Forgot Password", "Email recovery card with success state", "auth", nil, gritUIForgotPasswordCard01TSX()},
		"otp-input-01":            {"OTP Input", "6-digit OTP input with auto-advance, paste, and resend timer", "auth", nil, gritUIOtpInput01TSX()},
		"social-login-buttons-01": {"Social Login Buttons", "Google and GitHub OAuth buttons with divider", "auth", nil, gritUISocialLoginButtons01TSX()},
		"onboarding-wizard-01":    {"Onboarding Wizard", "3-step setup wizard with role, workspace, and plan selection", "auth", nil, gritUIOnboardingWizard01TSX()},
		"two-factor-setup-01":     {"Two-Factor Setup", "2FA setup with QR scan, verify code, and backup codes steps", "auth", nil, authUITwoFactorSetup01()},
		"profile-avatar-01":       {"Profile Avatar", "User avatar with image upload/preview and initials fallback", "auth", nil, authUIProfileAvatar01()},
		"account-deletion-01":     {"Account Deletion", "Confirmation dialog for irreversible account/workspace deletion", "auth", nil, authUIAccountDeletion01()},
		"oauth-button-01":         {"OAuth Button", "OAuth provider sign-in button (GitHub, Google, Discord, X)", "auth", nil, authUIOAuthButton01()},

		// ── SaaS (30) ─────────────────────────────────────────────────────────
		"billing-card-01":       {"Billing Card", "Current plan card with usage bar and upgrade button (Stripe-style)", "saas", nil, saasUIBillingCard01()},
		"usage-meter-01":        {"Usage Meter", "Resource usage bar with limit and percentage (Vercel/Linear-style)", "saas", nil, saasUIUsageMeter01()},
		"team-member-row-01":    {"Team Member Row", "Team roster row with avatar, role badge, and actions", "saas", nil, saasUITeamMemberRow01()},
		"api-key-row-01":        {"API Key Row", "API key row with masked key, copy, and revoke (Vercel-style)", "saas", nil, saasUIApiKeyRow01()},
		"kanban-card-01":        {"Kanban Card", "Task card with priority, labels, and assignee (Linear-style)", "saas", nil, saasUIKanbanCard01()},
		"metric-card-01":        {"Metric Card", "KPI card with trend arrow and sparkline (Vercel Analytics)", "saas", nil, saasUIMetricCard01()},
		"activity-item-01":      {"Activity Feed Item", "Activity log entry with actor, action, and target (GitHub-style)", "saas", nil, saasUIActivityItem01()},
		"notification-item-01":  {"Notification Item", "Notification row with unread dot and mark-read (Linear-style)", "saas", nil, saasUINotificationItem01()},
		"project-card-01":       {"Project Card", "Project/deployment card with status and last deploy (Vercel-style)", "saas", nil, saasUIProjectCard01()},
		"feature-flag-row-01":   {"Feature Flag Row", "Feature flag toggle row with key and rollout (LaunchDarkly)", "saas", nil, saasUIFeatureFlagRow01()},
		"webhook-row-01":        {"Webhook Row", "Webhook endpoint row with success rate and actions (Stripe-style)", "saas", nil, saasUIWebhookRow01()},
		"plan-upgrade-banner-01": {"Plan Upgrade Banner", "Upgrade prompt banner with feature callout (Intercom-style)", "saas", nil, saasUIPlanUpgradeBanner01()},
		"trial-countdown-01":    {"Trial Countdown", "Trial expiry countdown with days remaining and urgency styling", "saas", nil, saasUITrialCountdown01()},
		"integration-card-01":   {"Integration Card", "Integration connect card with logo and status (Notion/Linear)", "saas", nil, saasUIIntegrationCard01()},
		"empty-state-01":        {"Empty State", "Empty state with icon, title, description, and CTA (Linear-style)", "saas", nil, saasUIEmptyState01()},
		"command-item-01":       {"Command Item", "Command palette item with icon and keyboard shortcut (Linear)", "saas", nil, saasUICommandItem01()},
		"onboarding-step-01":    {"Onboarding Step", "Onboarding checklist step with status and action (Stripe)", "saas", nil, saasUIOnboardingStep01()},
		"quota-alert-01":        {"Quota Alert", "Usage quota warning with progress bar and upgrade (Vercel)", "saas", nil, saasUIQuotaAlert01()},
		"audit-log-row-01":      {"Audit Log Row", "Audit trail entry with actor, action, IP, and timestamp (Stripe)", "saas", nil, saasUIAuditLogRow01()},
		"settings-section-01":   {"Settings Section", "Settings card with title, description, and action area", "saas", nil, saasUISettingsSection01()},
		"invitation-link-01":    {"Invitation Link", "Workspace invite link with copy and expiry (Notion-style)", "saas", nil, saasUIInvitationLink01()},
		"file-upload-zone-01":   {"File Upload Zone", "Drag-and-drop file upload zone (Linear attachments-style)", "saas", nil, saasUIFileUploadZone01()},
		"progress-tracker-01":   {"Progress Tracker", "Multi-step completion tracker with step names (Stripe setup)", "saas", nil, saasUIProgressTracker01()},
		"search-filter-bar-01":  {"Search Filter Bar", "Search input with filter pills (GitHub issues-style)", "saas", nil, saasUISearchFilterBar01()},
		"data-table-row-01":     {"Data Table Row", "Table row with checkbox, status badge, and actions", "saas", nil, saasUIDataTableRow01()},
		"session-row-01":        {"Session Row", "Active session row with device, location, and revoke (GitHub)", "saas", nil, saasUISessionRow01()},
		"billing-history-row-01": {"Billing History Row", "Invoice row with status badge and PDF download (Stripe)", "saas", nil, saasUIBillingHistoryRow01()},
		"workspace-switcher-01": {"Workspace Switcher", "Workspace/org switcher dropdown (Linear/Notion-style)", "saas", nil, saasUIWorkspaceSwitcher01()},
		"permission-row-01":     {"Permission Row", "Permission scope row with role toggle buttons (Linear)", "saas", nil, saasUIPermissionRow01()},
		"sso-button-01":         {"SSO Button", "SSO provider sign-in button with org info (Okta-style)", "saas", nil, saasUISSOButton01()},

		// ── Ecommerce (19 new + product-card-01 existing = 20) ────────────────
		"product-card-list-01":      {"Product List Row", "Product list row with image, rating, price, and stock (Shopify)", "ecommerce", nil, ecomUIProductCardList01()},
		"product-gallery-01":        {"Product Gallery", "Product image gallery with thumbnail strip (Shopify PDP)", "ecommerce", nil, ecomUIProductGallery01()},
		"product-variant-01":        {"Product Variant Selector", "Color swatch and size variant selectors (Shopify PDP)", "ecommerce", nil, ecomUIProductVariant01()},
		"cart-item-01":              {"Cart Item", "Cart line item with image, quantity spinner, and remove (Shopify)", "ecommerce", nil, ecomUICartItem01()},
		"cart-summary-01":           {"Cart Summary", "Order summary card with subtotal, tax, and checkout (Shopify)", "ecommerce", nil, ecomUICartSummary01()},
		"order-status-01":           {"Order Status Tracker", "Order progress timeline from placed to delivered (Shopify)", "ecommerce", nil, ecomUIOrderStatus01()},
		"order-card-01":             {"Order Card", "Order history card with items, status, and track button", "ecommerce", nil, ecomUIOrderCard01()},
		"review-card-01":            {"Review Card", "Product review with stars, verified badge, and helpful vote", "ecommerce", nil, ecomUIReviewCard01()},
		"discount-badge-01":         {"Discount Badge", "Sale / percent off badge and flash sale pill", "ecommerce", nil, ecomUIDiscountBadge01()},
		"stock-indicator-01":        {"Stock Indicator", "In stock / low stock / out of stock indicator (Shopify)", "ecommerce", nil, ecomUIStockIndicator01()},
		"flash-sale-timer-01":       {"Flash Sale Timer", "Countdown timer with urgency styling for limited offers", "ecommerce", nil, ecomUIFlashSaleTimer01()},
		"bundle-offer-01":           {"Bundle Offer", "Frequently bought together bundle card (Amazon-style)", "ecommerce", nil, ecomUIBundleOffer01()},
		"category-filter-01":        {"Category Filter", "Filter sidebar with categories, price range, and ratings", "ecommerce", nil, ecomUICategoryFilter01()},
		"checkout-step-indicator-01": {"Checkout Step Indicator", "Checkout progress bar (Cart → Shipping → Payment)", "ecommerce", nil, ecomUICheckoutStepIndicator01()},
		"payment-method-01":         {"Payment Method Selector", "Payment method selector (card, PayPal, Apple Pay)", "ecommerce", nil, ecomUIPaymentMethod01()},
		"address-card-01":           {"Address Card", "Saved address card with select, edit, and default actions", "ecommerce", nil, ecomUIAddressCard01()},
		"wishlist-button-01":        {"Wishlist Button", "Heart toggle button for wishlist / save for later", "ecommerce", nil, ecomUIWishlistButton01()},
		"shipping-estimate-01":      {"Shipping Estimate", "Shipping option selector with carrier and delivery estimate", "ecommerce", nil, ecomUIShippingEstimate01()},
		"review-summary-01":         {"Review Summary", "Rating summary with star distribution bars (Amazon-style)", "ecommerce", nil, ecomUIReviewSummary01()},

		// ── Extra Marketing (7) ────────────────────────────────────────────────
		"social-proof-bar-01":      {"Social Proof Bar", "Avatar stack with member count and star rating trust bar", "marketing", nil, mktSocialProofBar01()},
		"feature-comparison-01":    {"Feature Comparison", "Feature comparison table across plans (Linear pricing-style)", "marketing", nil, mktFeatureComparison01()},
		"announcement-bar-01":      {"Announcement Bar", "Dismissible top-of-page announcement banner (Vercel-style)", "marketing", nil, mktAnnouncementBar01()},
		"waitlist-form-01":          {"Waitlist Form", "Pre-launch waitlist signup with position counter and join count", "marketing", nil, mktWaitlistForm01()},
		"case-study-card-01":       {"Case Study Card", "Customer case study card with metric callout and quote", "marketing", nil, mktCaseStudyCard01()},
		"roadmap-item-01":          {"Roadmap Item", "Public roadmap card with status badge and upvote button", "marketing", nil, mktRoadmapItem01()},
		"partner-card-01":          {"Partner Card", "Technology partner showcase card with official badge", "marketing", nil, mktPartnerCard01()},

		// ── Layout (18 + 2 extra = 20) ─────────────────────────────────────────
		"app-navbar-01":         {"App Navbar", "App top navbar with logo, nav links, and user menu (Vercel/Linear)", "layout", nil, layoutUIAppNavbar01()},
		"marketing-navbar-01":   {"Marketing Navbar", "Marketing site navbar with CTA button and mobile menu", "layout", nil, layoutUIMarketingNavbar01()},
		"sidebar-nav-01":        {"Sidebar Nav", "App sidebar with icon + label nav items (Linear/Notion)", "layout", nil, layoutUISidebarNav01()},
		"breadcrumb-01":         {"Breadcrumb", "Breadcrumb navigation with chevron separators (Vercel)", "layout", nil, layoutUIBreadcrumb01()},
		"page-header-01":        {"Page Header", "Page header with title, description, badge, and action button", "layout", nil, layoutUIPageHeader01()},
		"marketing-footer-01":   {"Marketing Footer", "Footer with 4 column links and copyright (Vercel-style)", "layout", nil, layoutUIMarketingFooter01()},
		"alert-banner-01":       {"Alert Banner", "Alert/warning/success/error inline banner component", "layout", nil, layoutUIAlertBanner01()},
		"modal-shell-01":        {"Modal Shell", "Modal dialog shell with Escape key close and backdrop", "layout", nil, layoutUIModalShell01()},
		"code-block-01":         {"Code Block", "Code display with language badge, copy button, line numbers (Vercel docs)", "layout", nil, layoutUICodeBlock01()},
		"stat-card-01":          {"Stat Card", "Simple KPI stat card with icon and trend (Stripe dashboard)", "layout", nil, layoutUIStatCard01()},
		"section-header-01":     {"Section Header", "Centered marketing section header with badge and subtitle", "layout", nil, layoutUISectionHeader01()},
		"tab-nav-01":            {"Tab Nav", "Tab navigation bar with badges (GitHub repo tabs-style)", "layout", nil, layoutUITabNav01()},
		"loading-skeleton-01":   {"Loading Skeleton", "Animated skeleton placeholder and pre-built card/table layouts", "layout", nil, layoutUILoadingSkeleton01()},
		"toast-notification-01": {"Toast Notification", "Toast notification with variants and dismiss (Sonner-style)", "layout", nil, layoutUIToastNotification01()},
		"data-list-01":          {"Data List", "Key-value data list / details panel (Stripe details-style)", "layout", nil, layoutUIDataList01()},
		"error-page-01":         {"Error Page", "404/500 error page with back and home buttons (Vercel-style)", "layout", nil, layoutUIErrorPage01()},
		"card-grid-01":          {"Card Grid", "Responsive card grid wrapper with header and action slot", "layout", nil, layoutUICardGrid01()},
		"floating-action-01":    {"Floating Action Button", "Fixed floating action button (Notion new page-style)", "layout", nil, layoutUIFloatingAction01()},
		"settings-nav-01":           {"Settings Nav", "Settings sidebar with grouped navigation sections (Vercel/GitHub-style)", "layout", nil, layoutSettingsNav01()},
		"command-palette-shell-01":  {"Command Palette Shell", "Full command palette modal with search, groups, and keyboard shortcuts", "layout", nil, layoutCommandPaletteShell01()},
	}

	// Write per-component JSON files and TSX source files
	for name, info := range allComponents {
		jsonPath := filepath.Join(uiRoot, "registry", name+".json")
		files[jsonPath] = gritUIComponentRegistryItem(name, info.title, info.description, info.category, info.deps)

		// Write TSX source if available (new components only; original 4 rely on DB seed)
		if info.tsx != "" {
			tsxPath := filepath.Join(uiRoot, "registry", name+".tsx")
			files[tsxPath] = info.tsx
		}
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// gritUIComponentRegistryItem generates the shadcn-compatible per-component JSON.
func gritUIComponentRegistryItem(name, title, description, category string, deps []string) string {
	depsJSON := "[]"
	if len(deps) > 0 {
		quoted := make([]string, len(deps))
		for i, d := range deps {
			quoted[i] = fmt.Sprintf("%q", d)
		}
		depsJSON = "[" + strings.Join(quoted, ", ") + "]"
	}

	return fmt.Sprintf(`{
  "$schema": "https://ui.shadcn.com/schema/registry-item.json",
  "name": %q,
  "type": "registry:block",
  "title": %q,
  "description": %q,
  "files": [
    {
      "path": "components/grit-ui/%s.tsx",
      "type": "registry:component",
      "target": "components/grit-ui/%s.tsx"
    }
  ],
  "dependencies": %s,
  "registryDependencies": []
}
`, name, title, description, name, name, depsJSON)
}

func gritUIPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "@%s/grit-ui",
  "version": "0.1.0",
  "private": true,
  "description": "Grit UI component registry — shadcn-compatible blocks for Grit apps"
}
`, opts.ProjectName)
}

func gritUIRegistryJSON(_ Options) string {
	// All 100 components listed in the registry root
	items := []struct{ name, title, description, category string }{
		// Original 4
		{"hero-01", "Hero Section", "Bold hero with gradient text and dual CTA buttons", "marketing"},
		{"pricing-card-01", "Pricing Card", "Pricing tier card with feature list and popular badge", "marketing"},
		{"product-card-01", "Product Card", "E-commerce product card with image, rating, and add-to-cart", "ecommerce"},
		{"feature-grid-01", "Feature Grid", "Feature showcase grid with icon cards", "marketing"},
		// Marketing (10 new)
		{"hero-split-01", "Split Hero", "Side-by-side hero with copy and terminal code block", "marketing"},
		{"logo-cloud-01", "Logo Cloud", "Infinite-scroll marquee of customer logos", "marketing"},
		{"testimonial-card-01", "Testimonial Card", "Single customer testimonial with avatar and stars", "marketing"},
		{"testimonial-grid-01", "Testimonial Grid", "3-column grid of customer testimonials", "marketing"},
		{"stats-banner-01", "Stats Banner", "Stats strip with large numbers and radial glow", "marketing"},
		{"faq-accordion-01", "FAQ Accordion", "Click-to-expand FAQ accordion", "marketing"},
		{"cta-banner-01", "CTA Banner", "Full-width CTA with email capture", "marketing"},
		{"newsletter-signup-01", "Newsletter Signup", "Newsletter card with subscriber count", "marketing"},
		{"changelog-item-01", "Changelog Item", "Version changelog entry with timeline", "marketing"},
		{"pricing-toggle-01", "Pricing Toggle", "Full pricing section with monthly/annual toggle", "marketing"},
		// Auth (10)
		{"login-card-01", "Login Card", "Premium login card with social OAuth buttons", "auth"},
		{"signup-card-01", "Signup Card", "Split signup form with feature list", "auth"},
		{"forgot-password-card-01", "Forgot Password", "Email recovery card with success state", "auth"},
		{"otp-input-01", "OTP Input", "6-digit OTP input with auto-advance and resend timer", "auth"},
		{"social-login-buttons-01", "Social Login Buttons", "Google and GitHub OAuth buttons with divider", "auth"},
		{"onboarding-wizard-01", "Onboarding Wizard", "3-step setup wizard with role and plan selection", "auth"},
		{"two-factor-setup-01", "Two-Factor Setup", "2FA setup with QR scan, verify, and backup codes", "auth"},
		{"profile-avatar-01", "Profile Avatar", "User avatar with image upload and initials fallback", "auth"},
		{"account-deletion-01", "Account Deletion", "Confirmation dialog for account deletion", "auth"},
		{"oauth-button-01", "OAuth Button", "OAuth provider button (GitHub, Google, Discord, X)", "auth"},
		// SaaS (30)
		{"billing-card-01", "Billing Card", "Current plan card with usage bar and upgrade button", "saas"},
		{"usage-meter-01", "Usage Meter", "Resource usage bar with limit indicator", "saas"},
		{"team-member-row-01", "Team Member Row", "Team roster row with avatar, role badge, and actions", "saas"},
		{"api-key-row-01", "API Key Row", "API key row with masked key, copy, and revoke", "saas"},
		{"kanban-card-01", "Kanban Card", "Task card with priority, labels, and assignee", "saas"},
		{"metric-card-01", "Metric Card", "KPI card with trend arrow and sparkline", "saas"},
		{"activity-item-01", "Activity Feed Item", "Activity log entry with actor, action, and target", "saas"},
		{"notification-item-01", "Notification Item", "Notification row with unread dot and mark-read", "saas"},
		{"project-card-01", "Project Card", "Project/deployment card with status and last deploy", "saas"},
		{"feature-flag-row-01", "Feature Flag Row", "Feature flag toggle row with key and rollout", "saas"},
		{"webhook-row-01", "Webhook Row", "Webhook endpoint row with success rate and actions", "saas"},
		{"plan-upgrade-banner-01", "Plan Upgrade Banner", "Upgrade prompt banner with feature callout", "saas"},
		{"trial-countdown-01", "Trial Countdown", "Trial expiry countdown with days remaining", "saas"},
		{"integration-card-01", "Integration Card", "Integration connect card with logo and status", "saas"},
		{"empty-state-01", "Empty State", "Empty state with icon, title, description, and CTA", "saas"},
		{"command-item-01", "Command Item", "Command palette item with icon and shortcut", "saas"},
		{"onboarding-step-01", "Onboarding Step", "Onboarding checklist step with status and action", "saas"},
		{"quota-alert-01", "Quota Alert", "Usage quota warning with progress bar", "saas"},
		{"audit-log-row-01", "Audit Log Row", "Audit trail entry with actor, action, and IP", "saas"},
		{"settings-section-01", "Settings Section", "Settings card with title, description, and action", "saas"},
		{"invitation-link-01", "Invitation Link", "Workspace invite link with copy and expiry", "saas"},
		{"file-upload-zone-01", "File Upload Zone", "Drag-and-drop file upload zone", "saas"},
		{"progress-tracker-01", "Progress Tracker", "Multi-step completion tracker with step names", "saas"},
		{"search-filter-bar-01", "Search Filter Bar", "Search input with filter pills", "saas"},
		{"data-table-row-01", "Data Table Row", "Table row with checkbox, status badge, and actions", "saas"},
		{"session-row-01", "Session Row", "Active session row with device, location, and revoke", "saas"},
		{"billing-history-row-01", "Billing History Row", "Invoice row with status badge and PDF download", "saas"},
		{"workspace-switcher-01", "Workspace Switcher", "Workspace/org switcher dropdown", "saas"},
		{"permission-row-01", "Permission Row", "Permission scope row with role toggle buttons", "saas"},
		{"sso-button-01", "SSO Button", "SSO provider sign-in button with org info", "saas"},
		// Ecommerce (20)
		{"product-card-list-01", "Product List Row", "Product list row with image, rating, and stock", "ecommerce"},
		{"product-gallery-01", "Product Gallery", "Product image gallery with thumbnail strip", "ecommerce"},
		{"product-variant-01", "Product Variant Selector", "Color swatch and size variant selectors", "ecommerce"},
		{"cart-item-01", "Cart Item", "Cart line item with quantity spinner and remove", "ecommerce"},
		{"cart-summary-01", "Cart Summary", "Order summary card with subtotal, tax, and checkout", "ecommerce"},
		{"order-status-01", "Order Status Tracker", "Order progress timeline from placed to delivered", "ecommerce"},
		{"order-card-01", "Order Card", "Order history card with items, status, and track button", "ecommerce"},
		{"review-card-01", "Review Card", "Product review with stars and verified badge", "ecommerce"},
		{"discount-badge-01", "Discount Badge", "Sale / percent off badge and flash sale pill", "ecommerce"},
		{"stock-indicator-01", "Stock Indicator", "In stock / low stock / out of stock indicator", "ecommerce"},
		{"flash-sale-timer-01", "Flash Sale Timer", "Countdown timer for limited-time offers", "ecommerce"},
		{"bundle-offer-01", "Bundle Offer", "Frequently bought together bundle card", "ecommerce"},
		{"category-filter-01", "Category Filter", "Filter sidebar with categories and price range", "ecommerce"},
		{"checkout-step-indicator-01", "Checkout Step Indicator", "Checkout progress bar", "ecommerce"},
		{"payment-method-01", "Payment Method Selector", "Payment method selector (card, PayPal, Apple Pay)", "ecommerce"},
		{"address-card-01", "Address Card", "Saved address card with select, edit, and default", "ecommerce"},
		{"wishlist-button-01", "Wishlist Button", "Heart toggle button for wishlist", "ecommerce"},
		{"shipping-estimate-01", "Shipping Estimate", "Shipping option selector with delivery estimate", "ecommerce"},
		{"review-summary-01", "Review Summary", "Rating summary with star distribution bars", "ecommerce"},
		// Extra Marketing (7)
		{"social-proof-bar-01", "Social Proof Bar", "Avatar stack with member count and star rating trust bar", "marketing"},
		{"feature-comparison-01", "Feature Comparison", "Feature comparison table across plans (Linear pricing-style)", "marketing"},
		{"announcement-bar-01", "Announcement Bar", "Dismissible top-of-page announcement banner (Vercel-style)", "marketing"},
		{"waitlist-form-01", "Waitlist Form", "Pre-launch waitlist signup with position counter and join count", "marketing"},
		{"case-study-card-01", "Case Study Card", "Customer case study card with metric callout and quote", "marketing"},
		{"roadmap-item-01", "Roadmap Item", "Public roadmap card with status badge and upvote button", "marketing"},
		{"partner-card-01", "Partner Card", "Technology partner showcase card with official badge", "marketing"},
		// Layout (18 + 2 extra = 20)
		{"app-navbar-01", "App Navbar", "App top navbar with logo, nav links, and user menu", "layout"},
		{"marketing-navbar-01", "Marketing Navbar", "Marketing site navbar with CTA button", "layout"},
		{"sidebar-nav-01", "Sidebar Nav", "App sidebar with icon + label nav items", "layout"},
		{"breadcrumb-01", "Breadcrumb", "Breadcrumb navigation with chevron separators", "layout"},
		{"page-header-01", "Page Header", "Page header with title, description, and action button", "layout"},
		{"marketing-footer-01", "Marketing Footer", "Footer with column links and copyright", "layout"},
		{"alert-banner-01", "Alert Banner", "Alert/warning/success/error inline banner", "layout"},
		{"modal-shell-01", "Modal Shell", "Modal dialog shell with Escape key close", "layout"},
		{"code-block-01", "Code Block", "Code display with language badge and copy button", "layout"},
		{"stat-card-01", "Stat Card", "Simple KPI stat card with icon and trend", "layout"},
		{"section-header-01", "Section Header", "Centered marketing section header with badge", "layout"},
		{"tab-nav-01", "Tab Nav", "Tab navigation bar with badges", "layout"},
		{"loading-skeleton-01", "Loading Skeleton", "Animated skeleton placeholder", "layout"},
		{"toast-notification-01", "Toast Notification", "Toast notification with variants and dismiss", "layout"},
		{"data-list-01", "Data List", "Key-value data list / details panel", "layout"},
		{"error-page-01", "Error Page", "404/500 error page with navigation buttons", "layout"},
		{"card-grid-01", "Card Grid", "Responsive card grid wrapper with header", "layout"},
		{"floating-action-01", "Floating Action Button", "Fixed floating action button", "layout"},
		{"settings-nav-01", "Settings Nav", "Settings sidebar with grouped navigation sections (Vercel/GitHub-style)", "layout"},
		{"command-palette-shell-01", "Command Palette Shell", "Full command palette modal with search, groups, and keyboard shortcuts", "layout"},
	}

	// Build items JSON
	var itemsJSON strings.Builder
	for i, item := range items {
		if i > 0 {
			itemsJSON.WriteString(",\n    ")
		}
		itemsJSON.WriteString(fmt.Sprintf(
			`{"name": %q, "type": "registry:block", "title": %q, "description": %q, "category": %q, "registryUrl": "http://localhost:8080/r/%s.json"}`,
			item.name, item.title, item.description, item.category, item.name,
		))
	}

	return fmt.Sprintf(`{
  "$schema": "https://ui.shadcn.com/schema/registry.json",
  "name": "grit-ui",
  "homepage": "https://gritframework.dev/components",
  "items": [
    %s
  ]
}
`, itemsJSON.String())
}
