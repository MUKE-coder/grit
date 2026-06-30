import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Grit UI: Using the 100-Component Registry',
  description: 'Learn to use Grit UI, a 100-component shadcn-compatible registry. Install marketing, auth, SaaS, ecommerce, and layout components into your project.',
}

export default function GritUIComponentsCourse() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Grit UI Components</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">10 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Grit UI: Using the 100-Component Registry
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Grit ships 100 shadcn-compatible React components across 5 categories. Install them into
            your project with a single command, customize the source code freely, and build complete
            pages from pre-built building blocks.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is Grit UI? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Grit UI?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Building UI from scratch is slow. Every project needs hero sections, login forms, pricing
            tables, dashboards, and navigation — and they all look roughly the same. Grit UI gives
            you 100 pre-built components that you install, own, and customize.
          </p>

          <Definition term="Component Registry">
            A collection of pre-built, installable React components served via an API. Instead of
            installing an npm package (where you can{"'"}t modify the source), you download the actual
            TSX source code into your project. You own it completely — change colors, layout, text,
            add features. No version lock-in.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit UI includes 100 components across 5 categories:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Marketing (21)</strong> — hero sections, feature grids, pricing, testimonials, CTAs, footers</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Auth (10)</strong> — login, register, forgot password, OTP, social login</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">SaaS (30)</strong> — dashboards, billing, settings, onboarding, analytics, teams</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Ecommerce (20)</strong> — product cards, cart, checkout, orders, categories</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Layout (20)</strong> — navbars, sidebars, footers, headers, app shells, page layouts</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The registry is served by your Grit API at the <Code>/r</Code> endpoints. Every Grit
            project includes it automatically.
          </p>

          <Challenge number={1} title="Browse the Registry">
            <p>Start your Grit API and open <Code>http://localhost:8080/r.json</Code> in your browser.
            How many components are listed? What fields does each component have? Can you find
            components in all 5 categories?</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: How the Registry Works ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How the Registry Works</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your Grit API exposes two registry endpoints. The first lists all components. The second
            returns a specific component with its full source code.
          </p>

          <CodeBlock filename="Registry API">
{`# List all 100 components
GET /r.json
# Returns: array of components with name, description, category

# Get a specific component
GET /r/hero-split.json
# Returns: full component data with embedded TSX source code`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Each component JSON contains:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">name</strong> — unique identifier (e.g., {'"'}hero-split{'"'})</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">description</strong> — what the component does</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">category</strong> — marketing, auth, saas, ecommerce, or layout</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">dependencies</strong> — npm packages the component needs (shadcn/ui primitives, lucide-react)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">files</strong> — the full TSX source code, ready to copy into your project</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The format is shadcn-compatible. This means any tool that works with the shadcn registry
            format (including the official CLI) can install Grit UI components.
          </p>

          <Challenge number={2} title="Fetch a Component">
            <p>Fetch <Code>/r/hero-split.json</Code> from your API. What fields does the response
            contain? Look at the files field — can you see the full TSX source code? What
            dependencies does the component require?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Installing Components ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Installing Components</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Installing a Grit UI component uses the shadcn CLI. One command downloads the component
            source into your project. You own the code — it{"'"}s not a dependency you import from
            node_modules.
          </p>

          <CodeBlock filename="Terminal">
{`# Install a component from your local Grit API
npx shadcn@latest add hero-split --url http://localhost:8080/r`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This downloads the <Code>hero-split.tsx</Code> file into your components directory
            (typically <Code>components/ui/</Code> or wherever your shadcn config points). The file
            is now part of your project — edit it freely.
          </p>

          <Definition term="shadcn/ui">
            A component system where you copy source code into your project instead of installing
            an npm package. You get full control over every line of code. No version conflicts, no
            breaking updates, no fighting with a library{"'"}s API. The tradeoff: you maintain the
            code yourself. For most teams, this is a much better deal.
          </Definition>

          <Tip>
            Your API must be running when you install components. The shadcn CLI fetches the component
            JSON from your API{"'"}s <Code>/r</Code> endpoint. Start your API with <Code>grit dev</Code>
            or <Code>go run cmd/server/main.go</Code> first.
          </Tip>

          <Challenge number={3} title="Install Your First Component">
            <p>Install a hero component using the shadcn CLI. Find the downloaded file in your project.
            Open it — can you read the TSX source code? Try importing it in a page and rendering it.
            Does it display correctly?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Marketing Components (21) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Marketing Components (21)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Marketing components are the building blocks of landing pages. Hero sections grab attention.
            Feature grids explain your product. Pricing tables convert visitors. Testimonials build trust.
            CTAs drive action.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key marketing components:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">hero-split</strong> — headline on the left, image on the right. Classic SaaS hero layout.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">hero-centered</strong> — centered headline with CTA buttons below. Clean and focused.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">features-grid</strong> — 3-column grid of feature cards with icons and descriptions.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">pricing-three-tier</strong> — three pricing plans side by side with feature comparison.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">testimonials-carousel</strong> — rotating customer quotes with avatars and company names.</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A typical landing page stacks these components top to bottom: hero, features, testimonials,
            pricing, CTA, footer. Five components, one complete landing page.
          </p>

          <CodeBlock filename="Landing Page Structure">
{`// A complete landing page from Grit UI components:
//
// HeroSplit - "Build faster with Grit"
// FeaturesGrid - 6 features with icons
// TestimonialsCarousel - 5 customer quotes
// PricingThreeTier - Free, Pro, Enterprise
// CTABanner - "Start building today"
// FooterLinks - Links, social, copyright`}
          </CodeBlock>

          <Challenge number={4} title="Install 3 Marketing Components">
            <p>Install <Code>hero-split</Code>, <Code>features-grid</Code>, and <Code>pricing-three-tier</Code>.
            Import all three into a page component and render them in order. You{"'"}ve just built a
            landing page in under 5 minutes.</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Auth Components (10) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Auth Components (10)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Auth pages are the first thing users see. They need to look polished, work on mobile,
            and handle edge cases (loading states, error messages, password visibility toggle).
            Grit UI{"'"}s auth components handle all of this.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key auth components:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">login-card</strong> — email/password form with remember me, forgot password link, and social login buttons</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">register-card</strong> — name, email, password, confirm password with validation</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">forgot-password</strong> — email input with submit and back-to-login link</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">otp-input</strong> — 6-digit code input with auto-focus and paste support</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">social-buttons</strong> — Google, GitHub, Twitter login buttons styled consistently</li>
          </ul>

          <Note>
            These components handle the UI only — form layout, validation feedback, and styling.
            You still wire them up to your auth API endpoints. The scaffolded auth pages already
            call the API, so compare them side by side.
          </Note>

          <Challenge number={5} title="Install Login Card">
            <p>Install the <Code>login-card</Code> component. Compare it to your project{"'"}s default
            login page. What{"'"}s different? What{"'"}s the same? Try replacing the default login page
            with the Grit UI version.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: SaaS Components (30) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">SaaS Components (30)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            SaaS components power the authenticated app experience. Dashboards, settings pages,
            billing flows, onboarding wizards, analytics — the screens your paying users interact
            with daily. This is the largest category with 30 components.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key SaaS components:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">dashboard-stats</strong> — 4-card grid showing key metrics (users, revenue, growth, active)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">billing-plans</strong> — subscription plan selector with monthly/yearly toggle</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">settings-profile</strong> — profile form with avatar upload, name, email, bio</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">onboarding-wizard</strong> — multi-step setup flow with progress indicator</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">analytics-chart</strong> — line/bar chart component for visualizing data over time</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            SaaS components often need data from your API. The components render the UI, and you
            connect them to your React Query hooks for data fetching. The pattern is always the same:
            fetch data, pass it as props, render the component.
          </p>

          <Tip>
            Start with <Code>dashboard-stats</Code> — it{"'"}s the most impactful. Replace static numbers
            with real data from your API and your dashboard instantly looks professional.
          </Tip>

          <Challenge number={6} title="Install SaaS Components">
            <p>Install <Code>dashboard-stats</Code>, <Code>billing-plans</Code>, and <Code>settings-profile</Code>.
            Use <Code>dashboard-stats</Code> on your dashboard page. Pass in real numbers from your API
            (total users, total posts, etc.). How does it compare to the default dashboard?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Ecommerce Components (20) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Ecommerce Components (20)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Building an online store? Ecommerce components give you product displays, shopping cart,
            checkout flow, and order management. These are the most complex components — they handle
            state management, quantity updates, price calculations, and form validation.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key ecommerce components:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">product-card</strong> — image, title, price, add-to-cart button. Responsive grid item.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">cart-drawer</strong> — slide-out cart panel with quantity controls and subtotal</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">checkout-form</strong> — shipping address, payment info, order summary</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">order-history</strong> — table of past orders with status badges and view details</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">category-grid</strong> — grid of product categories with images and counts</li>
          </ul>

          <Note>
            Ecommerce components pair with Grit{"'"}s generated resources. Generate a Product resource
            with name, price, description, image, and category_id — then display them using the
            product-card component. The API provides the data, the component provides the UI.
          </Note>

          <Challenge number={7} title="Install Product Card">
            <p>Install the <Code>product-card</Code> component. If you have a Product resource, display
            your products using the component. If not, create a simple array of product data and
            render a grid of product cards. Customize the card{"'"}s colors to match your app.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Layout Components (20) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Layout Components (20)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Layout components define the structure of your app. Navbars, sidebars, footers, and app
            shells create the container that all other components live inside. Get these right and
            everything else falls into place.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key layout components:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">navbar-main</strong> — top navigation bar with logo, nav links, auth buttons, mobile menu</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">sidebar-collapsible</strong> — sidebar that collapses to icons on small screens</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">footer-links</strong> — multi-column footer with link groups, social icons, copyright</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">app-shell</strong> — complete app layout with sidebar, header, and content area</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">page-header</strong> — page title with breadcrumb and action buttons</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>app-shell</Code> component is especially powerful — it combines a sidebar,
            header, and content area into a single component. Replace your entire app layout with
            one install.
          </p>

          <Challenge number={8} title="Install Navbar and Footer">
            <p>Install <Code>navbar-main</Code> and <Code>footer-links</Code>. Replace your app{"'"}s
            default navigation with these components. Update the nav links to match your routes.
            Does the mobile menu work correctly?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Customizing Components ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Customizing Components</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The whole point of the copy-paste model is customization. After installing a component,
            the source code is in your project. Change anything: colors, layout, text, animations,
            add new props, remove features you don{"'"}t need.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Common customizations:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Colors</strong> — swap Tailwind classes to match your brand. Change <Code>bg-primary</Code> to <Code>bg-blue-600</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Text</strong> — replace placeholder copy with your product{"'"}s real messaging</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Layout</strong> — change from 3-column to 4-column grid, swap left/right positions</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Props</strong> — add new props for dynamic data (prices, feature lists, user info)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Remove</strong> — delete sections you don{"'"}t need (social login, testimonials, etc.)</li>
          </ul>

          <Tip>
            Start by changing the text and colors. Then move to structural changes. Don{"'"}t be afraid
            to delete code you don{"'"}t need — the component is yours. You can always re-install the
            original version and start over.
          </Tip>

          <Challenge number={9} title="Customize a Component">
            <p>Install a pricing component. Change the plan names and prices to match a fictional SaaS
            product (e.g., Starter $9/mo, Pro $29/mo, Enterprise $99/mo). Change the feature lists
            for each plan. Update the colors to match your app{"'"}s theme. The component should look
            completely custom.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s everything you learned in this course:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'Grit UI is a 100-component registry served by your API at /r endpoints',
              'Components are shadcn-compatible: install with npx shadcn@latest add --url',
              'You own the source code — copy-paste model, not dependency model',
              'Marketing (21): heroes, features, pricing, testimonials, CTAs for landing pages',
              'Auth (10): login, register, forgot password, OTP, social login forms',
              'SaaS (30): dashboards, billing, settings, onboarding, analytics for app pages',
              'Ecommerce (20): product cards, cart, checkout, orders for online stores',
              'Layout (20): navbars, sidebars, footers, app shells for app structure',
              'Customize freely: change colors, text, layout, add props, remove sections',
              'Build complete pages by composing multiple registry components together',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={10} title="Build a Landing Page from Components">
            <p>Build a complete landing page for a fictional SaaS product using ONLY Grit UI components.
            Install and customize: (1) a hero component with your product{"'"}s headline, (2) a features
            grid with 6 features, (3) a testimonials section with 3 quotes, (4) a pricing table with
            3 plans, (5) a footer with links and social icons. Every section should be a Grit UI
            component, customized with your product{"'"}s branding and copy. The entire page should look
            like a real product — not a template.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
