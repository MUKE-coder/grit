package scaffold

import (
	"fmt"
	"strings"
)

// apiUIComponentSeedExtendedGo returns Go source code for the extended UI component seed.
// This seeds all 96 new components beyond the original 4.
func apiUIComponentSeedExtendedGo() string {
	// nameToConst converts "api-key-row-01" → "apiKeyRow01"
	nameToConst := func(name string) string {
		parts := strings.Split(name, "-")
		for i := 1; i < len(parts); i++ {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
			}
		}
		return strings.Join(parts, "")
	}

	type entry struct {
		name, displayName, desc, cat, tags, preview, deps string
		tsx                                               func() string
	}

	// 96 new components (the original 4 are seeded in SeedUIComponents)
	entries := []entry{
		// ── Marketing (10) ──────────────────────────────────────────────────────
		{"hero-split-01", "Split Hero", "Side-by-side hero with left copy and right terminal code block", "marketing", `["hero","landing","split"]`, `<HeroSplit01 />`, `[]`, gritUIHeroSplit01TSX},
		{"logo-cloud-01", "Logo Cloud", "Infinite-scroll marquee of customer and partner logos", "marketing", `["logos","trust","marquee"]`, `<LogoCloud01 />`, `[]`, gritUILogoCLoud01TSX},
		{"testimonial-card-01", "Testimonial Card", "Single customer testimonial with avatar, stars, and quote", "marketing", `["testimonial","review","trust"]`, `<TestimonialCard01 quote="..." name="..." role="..." company="..." />`, `[]`, gritUITestimonialCard01TSX},
		{"testimonial-grid-01", "Testimonial Grid", "3-column masonry grid of customer testimonials", "marketing", `["testimonials","grid","social-proof"]`, `<TestimonialGrid01 />`, `[]`, gritUITestimonialGrid01TSX},
		{"stats-banner-01", "Stats Banner", "Dark stats strip with large numbers and radial glow", "marketing", `["stats","numbers","marketing"]`, `<StatsBanner01 />`, `[]`, gritUIStatsBanner01TSX},
		{"faq-accordion-01", "FAQ Accordion", "Click-to-expand FAQ accordion with animated chevron", "marketing", `["faq","accordion","marketing"]`, `<FaqAccordion01 />`, `[]`, gritUIFaqAccordion01TSX},
		{"cta-banner-01", "CTA Banner", "Full-width CTA section with email capture and success state", "marketing", `["cta","email","conversion"]`, `<CtaBanner01 title="Get started free" />`, `[]`, gritUICtaBanner01TSX},
		{"newsletter-signup-01", "Newsletter Signup", "Newsletter card with subscriber count badge and success state", "marketing", `["newsletter","email","signup"]`, `<NewsletterSignup01 />`, `[]`, gritUINewsletterSignup01TSX},
		{"changelog-item-01", "Changelog Item", "Version changelog entry with timeline dot, tags, and change list", "marketing", `["changelog","version","release"]`, `<ChangelogItem01 version="v1.0" date="Mar 2025" title="Initial Release" />`, `[]`, gritUIChangelogItem01TSX},
		{"pricing-toggle-01", "Pricing Toggle", "Full pricing section with monthly/annual toggle and 3 tier cards", "marketing", `["pricing","plans","toggle","billing"]`, `<PricingToggle01 />`, `[]`, gritUIPricingToggle01TSX},
		// ── Auth (10) ────────────────────────────────────────────────────────────
		{"login-card-01", "Login Card", "Premium login card with social OAuth buttons, remember me, forgot password", "auth", `["login","auth","oauth"]`, `<LoginCard01 />`, `[]`, gritUILoginCard01TSX},
		{"signup-card-01", "Signup Card", "Split signup form with feature list and password strength indicator", "auth", `["signup","register","auth"]`, `<SignupCard01 />`, `[]`, gritUISignupCard01TSX},
		{"forgot-password-card-01", "Forgot Password", "Email recovery card with clean success state", "auth", `["forgot","password","recovery"]`, `<ForgotPasswordCard01 />`, `[]`, gritUIForgotPasswordCard01TSX},
		{"otp-input-01", "OTP Input", "6-digit OTP input with auto-advance, paste support, and resend timer", "auth", `["otp","2fa","verification"]`, `<OtpInput01 email="user@example.com" onComplete={(code) => {}} />`, `[]`, gritUIOtpInput01TSX},
		{"social-login-buttons-01", "Social Login Buttons", "Google and GitHub OAuth sign-in buttons with divider", "auth", `["oauth","google","github","social"]`, `<SocialLoginButtons01 />`, `[]`, gritUISocialLoginButtons01TSX},
		{"onboarding-wizard-01", "Onboarding Wizard", "3-step setup wizard with role, workspace, and plan selection", "auth", `["onboarding","wizard","setup"]`, `<OnboardingWizard01 onComplete={() => {}} />`, `[]`, gritUIOnboardingWizard01TSX},
		{"two-factor-setup-01", "Two-Factor Setup", "2FA setup card with QR scan, verify code, and backup codes steps", "auth", `["2fa","security","totp"]`, `<TwoFactorSetup />`, `[]`, authUITwoFactorSetup01},
		{"profile-avatar-01", "Profile Avatar", "User avatar component with image upload and initials fallback", "auth", `["avatar","profile","upload"]`, `<ProfileAvatar name="Sarah Chen" email="sarah@example.com" />`, `[]`, authUIProfileAvatar01},
		{"account-deletion-01", "Account Deletion", "Confirmation dialog for irreversible account/workspace deletion", "auth", `["delete","danger","confirmation"]`, `<AccountDeletion workspaceName="Acme Corp" />`, `[]`, authUIAccountDeletion01},
		{"oauth-button-01", "OAuth Button", "OAuth provider sign-in button supporting GitHub, Google, Discord, X", "auth", `["oauth","social","login"]`, `<OAuthButton provider="github" />`, `[]`, authUIOAuthButton01},
		// ── SaaS (30) ────────────────────────────────────────────────────────────
		{"billing-card-01", "Billing Card", "Current plan overview card with usage bar and upgrade button (Stripe-style)", "saas", `["billing","plan","subscription"]`, `<BillingCard name="Pro" price="29" period="month" />`, `[]`, saasUIBillingCard01},
		{"usage-meter-01", "Usage Meter", "Resource usage bar with limit indicator and upgrade link (Vercel-style)", "saas", `["usage","quota","limits"]`, `<UsageMeter label="API Requests" used={42000} limit={100000} />`, `[]`, saasUIUsageMeter01},
		{"team-member-row-01", "Team Member Row", "Team roster row with avatar, role badge, last active, and actions", "saas", `["team","members","roles"]`, `<TeamMemberRow name="Sarah Chen" email="sarah@example.com" role="Admin" />`, `[]`, saasUITeamMemberRow01},
		{"api-key-row-01", "API Key Row", "API key row with masked key, copy button, and revoke (Vercel-style)", "saas", `["api","keys","security"]`, `<APIKeyRow name="Production" prefix="grit_live" suffix="xK9m" />`, `[]`, saasUIApiKeyRow01},
		{"kanban-card-01", "Kanban Card", "Issue/task card with priority dot, labels, and assignee (Linear-style)", "saas", `["kanban","tasks","issues","linear"]`, `<KanbanCard id="GRT-42" title="Implement OAuth" priority="high" />`, `[]`, saasUIKanbanCard01},
		{"metric-card-01", "Metric Card", "KPI metric card with trend indicator and sparkline (Vercel Analytics)", "saas", `["metrics","kpi","dashboard","analytics"]`, `<MetricCard label="Monthly Revenue" value="$12,430" change="+18.2%" trend="up" />`, `[]`, saasUIMetricCard01},
		{"activity-item-01", "Activity Feed Item", "Activity log item with actor avatar, action, and timestamp (GitHub-style)", "saas", `["activity","feed","log","github"]`, `<ActivityItem actor="sarah" action="merged pull request" target="#142" />`, `[]`, saasUIActivityItem01},
		{"notification-item-01", "Notification Item", "Notification row with unread indicator and mark-read button (Linear-style)", "saas", `["notifications","inbox","alerts"]`, `<NotificationItem title="Sarah commented" timestamp="5m ago" />`, `[]`, saasUINotificationItem01},
		{"project-card-01", "Project Card", "Project/deployment card with framework, status, and commit info (Vercel)", "saas", `["project","deployment","vercel"]`, `<ProjectCard name="grit-app" domain="grit-app.vercel.app" status="ready" />`, `[]`, saasUIProjectCard01},
		{"feature-flag-row-01", "Feature Flag Row", "Feature flag toggle row with key and rollout percentage (LaunchDarkly)", "saas", `["feature-flags","toggles","rollout"]`, `<FeatureFlagRow name="new-dashboard" displayName="New Dashboard" />`, `[]`, saasUIFeatureFlagRow01},
		{"webhook-row-01", "Webhook Row", "Webhook endpoint row with success rate and test/edit/delete (Stripe-style)", "saas", `["webhooks","integrations","stripe"]`, `<WebhookRow url="https://api.example.com/webhooks" events={["payment.success"]} />`, `[]`, saasUIWebhookRow01},
		{"plan-upgrade-banner-01", "Plan Upgrade Banner", "Upgrade prompt banner with feature callout and dismiss (Intercom-style)", "saas", `["upgrade","upsell","conversion"]`, `<PlanUpgradeBanner title="You've used 90% of your quota" />`, `[]`, saasUIPlanUpgradeBanner01},
		{"trial-countdown-01", "Trial Countdown", "Trial expiry countdown with days remaining and upgrade CTA", "saas", `["trial","countdown","conversion"]`, `<TrialCountdown daysLeft={7} planName="Pro" />`, `[]`, saasUITrialCountdown01},
		{"integration-card-01", "Integration Card", "Integration card with logo, description, and connect button (Linear/Notion)", "saas", `["integrations","apps","marketplace"]`, `<IntegrationCard name="GitHub" isConnected={false} />`, `[]`, saasUIIntegrationCard01},
		{"empty-state-01", "Empty State", "Empty state placeholder with icon, title, description, and CTA (Linear)", "saas", `["empty","placeholder","zero-state"]`, `<EmptyState title="No issues yet" icon="inbox" ctaLabel="Create issue" />`, `[]`, saasUIEmptyState01},
		{"command-item-01", "Command Item", "Command palette item with icon, label, and keyboard shortcut (Linear Cmd+K)", "saas", `["command","palette","search","keyboard"]`, `<CommandItem label="Create new issue" shortcut={["⌘", "N"]} />`, `[]`, saasUICommandItem01},
		{"onboarding-step-01", "Onboarding Step", "Onboarding checklist step with done/active/pending states (Stripe)", "saas", `["onboarding","checklist","setup"]`, `<OnboardingStep number={1} title="Set up your workspace" status="current" />`, `[]`, saasUIOnboardingStep01},
		{"quota-alert-01", "Quota Alert", "Usage quota warning with progress bar and upgrade link (Vercel/Netlify)", "saas", `["quota","limits","warning","alert"]`, `<QuotaAlert resource="Build minutes" used={450} limit={500} />`, `[]`, saasUIQuotaAlert01},
		{"audit-log-row-01", "Audit Log Row", "Audit trail entry with user, action, IP, and timestamp (Stripe-style)", "saas", `["audit","log","security","compliance"]`, `<AuditLogRow actor="sarah@example.com" action="Updated" resource="API key" />`, `[]`, saasUIAuditLogRow01},
		{"settings-section-01", "Settings Section", "Settings card section with title, description, and action area", "saas", `["settings","form","configuration"]`, `<SettingsSection title="Workspace Name" description="Edit your workspace name."><input /></SettingsSection>`, `[]`, saasUISettingsSection01},
		{"invitation-link-01", "Invitation Link", "Workspace invite link with copy button, expiry, and disable option", "saas", `["invite","link","sharing"]`, `<InvitationLink inviteUrl="https://app.example.com/invite/abc123" />`, `[]`, saasUIInvitationLink01},
		{"file-upload-zone-01", "File Upload Zone", "Drag-and-drop file upload zone with type hints (Linear attachments)", "saas", `["upload","file","drag-drop"]`, `<FileUploadZone accept=".png,.jpg,.pdf" maxSizeMB={10} />`, `[]`, saasUIFileUploadZone01},
		{"progress-tracker-01", "Progress Tracker", "Multi-step completion tracker with step names and states (Stripe)", "saas", `["progress","steps","setup","onboarding"]`, `<ProgressTracker steps={[{label:"Account",status:"done"},{label:"Team",status:"current"}]} />`, `[]`, saasUIProgressTracker01},
		{"search-filter-bar-01", "Search Filter Bar", "Search input with active filter pills (GitHub issues-style)", "saas", `["search","filter","issues"]`, `<SearchFilterBar placeholder="Search issues..." />`, `[]`, saasUISearchFilterBar01},
		{"data-table-row-01", "Data Table Row", "Data table row with checkbox, status badge, and hover actions", "saas", `["table","row","list","data"]`, `<DataTableRow cells={["user@example.com","Admin","Mar 1","Active"]} />`, `[]`, saasUIDataTableRow01},
		{"session-row-01", "Session Row", "Active session/device row with location, last active, and revoke (GitHub)", "saas", `["session","security","devices"]`, `<SessionRow deviceName="MacBook Pro" browser="Chrome 122" isCurrent={true} />`, `[]`, saasUISessionRow01},
		{"billing-history-row-01", "Billing History Row", "Invoice row with amount, status badge, and PDF download (Stripe)", "saas", `["billing","invoices","payment","history"]`, `<BillingHistoryRow invoiceId="INV-2025-003" amount="$29.00" status="paid" />`, `[]`, saasUIBillingHistoryRow01},
		{"workspace-switcher-01", "Workspace Switcher", "Workspace/org switcher dropdown trigger (Linear/Notion-style)", "saas", `["workspace","org","switcher","navigation"]`, `<WorkspaceSwitcher currentId="1" />`, `[]`, saasUIWorkspaceSwitcher01},
		{"permission-row-01", "Permission Row", "Permission scope toggle row with Viewer/Member/Admin buttons", "saas", `["permissions","roles","access-control"]`, `<PermissionRow scope="Members" currentRole="Member" />`, `[]`, saasUIPermissionRow01},
		{"sso-button-01", "SSO Button", "SSO provider sign-in button with org name and domain (Okta-style)", "saas", `["sso","enterprise","okta","auth"]`, `<SSOButton provider="Google Workspace" orgName="Acme Corp" orgDomain="acme.com" />`, `[]`, saasUISSOButton01},
		// ── Ecommerce (19) ───────────────────────────────────────────────────────
		{"product-card-list-01", "Product List Row", "Product list row with image, rating, price, and stock level (Shopify admin)", "ecommerce", `["product","list","ecommerce","shopify"]`, `<ProductCardList name="Wireless Headphones" price="$129.99" stock={42} />`, `[]`, ecomUIProductCardList01},
		{"product-gallery-01", "Product Gallery", "Product image gallery with clickable thumbnail strip (Shopify PDP)", "ecommerce", `["gallery","images","product","pdp"]`, `<ProductGallery productName="Wireless Headphones" />`, `[]`, ecomUIProductGallery01},
		{"product-variant-01", "Product Variant Selector", "Color swatch and size variant selector buttons (Shopify PDP)", "ecommerce", `["variants","colors","sizes","product"]`, `<ProductVariantSelector />`, `[]`, ecomUIProductVariant01},
		{"cart-item-01", "Cart Item", "Cart line item with image, quantity spinner, and remove button (Shopify cart)", "ecommerce", `["cart","line-item","checkout","ecommerce"]`, `<CartItem name="Wireless Headphones" price={129.99} quantity={1} />`, `[]`, ecomUICartItem01},
		{"cart-summary-01", "Cart Summary", "Order summary with subtotal, discount, tax, total, and checkout CTA (Shopify)", "ecommerce", `["cart","summary","checkout","order"]`, `<CartSummary subtotal={259.98} />`, `[]`, ecomUICartSummary01},
		{"order-status-01", "Order Status Tracker", "Order progress timeline from placed to delivered (Shopify/Amazon)", "ecommerce", `["order","tracking","shipping","status"]`, `<OrderStatus orderId="#1042" />`, `[]`, ecomUIOrderStatus01},
		{"order-card-01", "Order Card", "Order history card with items, status badge, and action buttons", "ecommerce", `["order","history","ecommerce"]`, `<OrderCard orderId="#1042" status="shipped" total="$169.97" />`, `[]`, ecomUIOrderCard01},
		{"review-card-01", "Review Card", "Product review with star rating, verified badge, and helpful vote", "ecommerce", `["review","rating","social-proof","product"]`, `<ReviewCard author="Alex M." rating={5} title="Best product ever" />`, `[]`, ecomUIReviewCard01},
		{"discount-badge-01", "Discount Badge", "Sale percentage-off badge and flash sale pill variants", "ecommerce", `["discount","badge","sale","promotion"]`, `<DiscountBadge type="percent" value="20" /><SaleBadge text="Flash Sale" />`, `[]`, ecomUIDiscountBadge01},
		{"stock-indicator-01", "Stock Indicator", "In stock / low stock / out of stock indicator with count (Shopify PDP)", "ecommerce", `["stock","inventory","product"]`, `<StockIndicator stock={8} lowThreshold={10} />`, `[]`, ecomUIStockIndicator01},
		{"flash-sale-timer-01", "Flash Sale Timer", "Live countdown timer with urgency styling for limited-time offers", "ecommerce", `["timer","countdown","flash-sale","urgency"]`, `<FlashSaleTimer label="Flash Sale ends in" />`, `[]`, ecomUIFlashSaleTimer01},
		{"bundle-offer-01", "Bundle Offer", "Frequently bought together bundle offer card (Amazon-style)", "ecommerce", `["bundle","upsell","cross-sell","ecommerce"]`, `<BundleOffer discountPercent={15} />`, `[]`, ecomUIBundleOffer01},
		{"category-filter-01", "Category Filter", "Filter sidebar with category checkboxes, price range, and star rating", "ecommerce", `["filter","sidebar","facets","search","ecommerce"]`, `<CategoryFilter />`, `[]`, ecomUICategoryFilter01},
		{"checkout-step-indicator-01", "Checkout Step Indicator", "Checkout progress breadcrumb (Cart → Information → Shipping → Payment)", "ecommerce", `["checkout","steps","progress","ecommerce"]`, `<CheckoutStepIndicator currentStep={1} />`, `[]`, ecomUICheckoutStepIndicator01},
		{"payment-method-01", "Payment Method Selector", "Payment selector with card form, PayPal, and Apple Pay (Stripe checkout)", "ecommerce", `["payment","checkout","stripe","ecommerce"]`, `<PaymentMethodSelector />`, `[]`, ecomUIPaymentMethod01},
		{"address-card-01", "Address Card", "Saved address card with select, edit, set-as-default, and delete actions", "ecommerce", `["address","shipping","checkout","ecommerce"]`, `<AddressCard name="Sarah Chen" city="San Francisco" isDefault={true} />`, `[]`, ecomUIAddressCard01},
		{"wishlist-button-01", "Wishlist Button", "Heart/bookmark toggle button for save-for-later / wishlist (Shopify/Amazon)", "ecommerce", `["wishlist","save","heart","product"]`, `<WishlistButton productName="Wireless Headphones" />`, `[]`, ecomUIWishlistButton01},
		{"shipping-estimate-01", "Shipping Estimate", "Shipping option selector with carrier, delivery estimate, and price", "ecommerce", `["shipping","delivery","checkout","ecommerce"]`, `<ShippingEstimate />`, `[]`, ecomUIShippingEstimate01},
		{"review-summary-01", "Review Summary", "Aggregate rating summary with star distribution bars (Amazon-style)", "ecommerce", `["reviews","rating","summary","social-proof"]`, `<ReviewSummary average={4.3} total={247} />`, `[]`, ecomUIReviewSummary01},
		// ── Layout (18) ──────────────────────────────────────────────────────────
		{"app-navbar-01", "App Navbar", "Dark app top navbar with logo, nav links, and user menu (Vercel/Linear)", "layout", `["navbar","navigation","header","app"]`, `<AppNavbar appName="Grit" />`, `[]`, layoutUIAppNavbar01},
		{"marketing-navbar-01", "Marketing Navbar", "Marketing site navbar with links, login, and CTA button", "layout", `["navbar","marketing","header","navigation"]`, `<MarketingNavbar brand="Grit" />`, `[]`, layoutUIMarketingNavbar01},
		{"sidebar-nav-01", "Sidebar Nav", "App sidebar with icon + label navigation and badge counts (Linear/Notion)", "layout", `["sidebar","navigation","app","layout"]`, `<SidebarNav />`, `[]`, layoutUISidebarNav01},
		{"breadcrumb-01", "Breadcrumb", "Breadcrumb navigation with chevron separators (Vercel dashboard-style)", "layout", `["breadcrumb","navigation","wayfinding"]`, `<Breadcrumb items={[{label:"Projects",href:"/projects"},{label:"App"}]} />`, `[]`, layoutUIBreadcrumb01},
		{"page-header-01", "Page Header", "Page header with title, description, optional badge, and action button slot", "layout", `["header","page","title","layout"]`, `<PageHeader title="Settings" description="Manage your workspace." />`, `[]`, layoutUIPageHeader01},
		{"marketing-footer-01", "Marketing Footer", "Footer with 4 column links, brand tagline, and copyright (Vercel-style)", "layout", `["footer","marketing","links"]`, `<MarketingFooter brand="Grit" />`, `[]`, layoutUIMarketingFooter01},
		{"alert-banner-01", "Alert Banner", "Inline alert/warning/success/error banner with dismiss button", "layout", `["alert","banner","notification","feedback"]`, `<AlertBanner variant="success" message="Changes saved." />`, `[]`, layoutUIAlertBanner01},
		{"modal-shell-01", "Modal Shell", "Modal dialog shell with Escape close, backdrop, and size variants", "layout", `["modal","dialog","overlay","ui"]`, `<ModalShell title="Confirm" onClose={() => {}}>...</ModalShell>`, `[]`, layoutUIModalShell01},
		{"code-block-01", "Code Block", "Code display block with language badge, copy button, optional line numbers (Vercel docs)", "layout", `["code","syntax","copy","developer"]`, `<CodeBlock code="grit new my-app" language="bash" />`, `[]`, layoutUICodeBlock01},
		{"stat-card-01", "Stat Card", "Simple KPI stat card with icon, value, and trend indicator (Stripe dashboard)", "layout", `["stat","metric","kpi","dashboard"]`, `<StatCard label="Total Revenue" value="$48,352" trend="up" />`, `[]`, layoutUIStatCard01},
		{"section-header-01", "Section Header", "Centered marketing section header with badge, title, and subtitle", "layout", `["section","heading","marketing","landing"]`, `<SectionHeader badge="Features" title="Everything you need" />`, `[]`, layoutUISectionHeader01},
		{"tab-nav-01", "Tab Nav", "Tab navigation bar with badge counts and active indicator (GitHub repo tabs)", "layout", `["tabs","navigation","ui","layout"]`, `<TabNav defaultTab="overview" />`, `[]`, layoutUITabNav01},
		{"loading-skeleton-01", "Loading Skeleton", "Animated skeleton placeholder with card and table row presets", "layout", `["skeleton","loading","placeholder","ui"]`, `<CardSkeleton /><TableRowSkeleton cols={4} />`, `[]`, layoutUILoadingSkeleton01},
		{"toast-notification-01", "Toast Notification", "Toast with success/error/info/warning variants, action, and dismiss", "layout", `["toast","notification","feedback","ui"]`, `<Toast message="Saved!" variant="success" />`, `[]`, layoutUIToastNotification01},
		{"data-list-01", "Data List", "Key-value data list / details panel (Stripe invoice details-style)", "layout", `["data","list","details","description"]`, `<DataList items={[{label:"Status",value:"Active"}]} />`, `[]`, layoutUIDataList01},
		{"error-page-01", "Error Page", "404/500 error page with code, message, and back/home navigation", "layout", `["error","404","500","page"]`, `<ErrorPage code="404" title="Page not found" />`, `[]`, layoutUIErrorPage01},
		{"card-grid-01", "Card Grid", "Responsive 2/3/4-column card grid wrapper with header and action slot", "layout", `["grid","cards","layout","responsive"]`, `<CardGrid title="Projects" columns={3}>...</CardGrid>`, `[]`, layoutUICardGrid01},
		{"floating-action-01", "Floating Action Button", "Fixed floating action button for primary creation actions (Notion-style)", "layout", `["fab","action","button","floating"]`, `<FloatingActionButton label="New" />`, `[]`, layoutUIFloatingAction01},
		// ── Extra Marketing (7) ───────────────────────────────────────────────────
		{"social-proof-bar-01", "Social Proof Bar", "Avatar stack with member count and star rating trust bar", "marketing", `["social-proof","trust","testimonial","avatar"]`, `<SocialProofBar count="50,000+" />`, `[]`, mktSocialProofBar01},
		{"feature-comparison-01", "Feature Comparison", "Feature comparison table across plans (Linear pricing-style)", "marketing", `["comparison","table","pricing","features"]`, `<FeatureComparison />`, `[]`, mktFeatureComparison01},
		{"announcement-bar-01", "Announcement Bar", "Dismissible top-of-page announcement banner (Vercel-style)", "marketing", `["announcement","banner","bar","notification"]`, `<AnnouncementBar message="v2.0 is here" />`, `[]`, mktAnnouncementBar01},
		{"waitlist-form-01", "Waitlist Form", "Pre-launch waitlist signup with position counter and join count", "marketing", `["waitlist","signup","pre-launch","email"]`, `<WaitlistForm />`, `[]`, mktWaitlistForm01},
		{"case-study-card-01", "Case Study Card", "Customer case study card with metric callout and quote", "marketing", `["case-study","customer","testimonial","social-proof"]`, `<CaseStudyCard company="Acme" metric="10x" metricLabel="faster deploys" />`, `[]`, mktCaseStudyCard01},
		{"roadmap-item-01", "Roadmap Item", "Public roadmap card with status badge and upvote button", "marketing", `["roadmap","upvote","feedback","public"]`, `<RoadmapItem title="Dark mode" status="planned" votes={142} />`, `[]`, mktRoadmapItem01},
		{"partner-card-01", "Partner Card", "Technology partner showcase card with official badge", "marketing", `["partner","technology","integration","badge"]`, `<PartnerCard name="Vercel" category="Deployment" />`, `[]`, mktPartnerCard01},
		// ── Extra Layout (2) ─────────────────────────────────────────────────────
		{"settings-nav-01", "Settings Nav", "Settings sidebar with grouped navigation sections (Vercel/GitHub-style)", "layout", `["settings","sidebar","navigation","grouped"]`, `<SettingsNav />`, `[]`, layoutSettingsNav01},
		{"command-palette-shell-01", "Command Palette Shell", "Full command palette modal with search, groups, and keyboard shortcuts", "layout", `["command","palette","search","keyboard","modal"]`, `<CommandPaletteShell isOpen={true} onClose={() => {}} />`, `[]`, layoutCommandPaletteShell01},
	}

	// Build TSX constants section
	var constsBuf strings.Builder
	for _, e := range entries {
		constName := nameToConst(e.name) + "TSX"
		tsxContent := e.tsx()
		constsBuf.WriteString("\nconst " + constName + " = `" + tsxContent + "`\n")
	}

	// Build seed entries section
	var entryBuf strings.Builder
	for _, e := range entries {
		constName := nameToConst(e.name) + "TSX"
		entryBuf.WriteString(fmt.Sprintf(`
		{
			Name:         %q,
			DisplayName:  %q,
			Description:  %q,
			Category:     %q,
			Tags:         %q,
			Files:        filesJSON(%q, %s),
			Dependencies: %q,
			RegistryDeps: "[]",
			PreviewCode:  %q,
			IsPublic:     true,
		},`, e.name, e.displayName, e.desc, e.cat, e.tags, e.name, constName, e.deps, e.preview))
	}

	header := `package models

import (
	"encoding/json"
	"log"

	"gorm.io/gorm"
)
`
	footer := `

// SeedUIComponentsExtended inserts the 96 extended Grit UI components.
// Safe to call on every startup — skips silently if components already exist.
// Must be called AFTER SeedUIComponents().
func SeedUIComponentsExtended(db *gorm.DB) {
	var count int64
	db.Model(&UIComponent{}).Count(&count)
	if count >= 100 {
		return
	}

	filesJSON := func(name, tsx string) string {
		data, _ := json.Marshal([]map[string]string{{
			"path":    "components/grit-ui/" + name + ".tsx",
			"type":    "registry:component",
			"target":  "components/grit-ui/" + name + ".tsx",
			"content": tsx,
		}})
		return string(data)
	}

	components := []UIComponent{` + entryBuf.String() + `
	}

	if err := db.Create(&components).Error; err != nil {
		log.Printf("Warning: failed to seed extended Grit UI components: %v", err)
	} else {
		log.Printf("Seeded %d extended Grit UI components", len(components))
	}
}
`

	return header + constsBuf.String() + footer
}
