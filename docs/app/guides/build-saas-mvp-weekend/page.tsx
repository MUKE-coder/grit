import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  Step,
  H3,
  Outcome,
  Note,
  Warn,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Build a SaaS MVP in a weekend with Go and React',
  description:
    'A two-day plan: tenants before anything else, auth, the one feature, billing and a deploy, with the decisions that are expensive to reverse marked.',
  alternates: { canonical: 'https://gritframework.dev/guides/build-saas-mvp-weekend' },
  openGraph: {
    title: 'Build a SaaS MVP in a weekend with Go and React',
    description: 'What to build on Saturday, what to build on Sunday, and what to skip.',
    url: 'https://gritframework.dev/guides/build-saas-mvp-weekend',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="build-saas-mvp-weekend"
      lead={
        <>
          <p className="mb-4">
            A weekend is enough for a SaaS that takes money, if you spend it on the right things.
            Most weekend attempts fail on Saturday afternoon, rebuilding authentication, and never
            reach the part that makes the product worth anything.
          </p>
          <p>
            This is a plan with an order. The order matters more than the speed: two decisions
            here are cheap now and very expensive in three months, and both are made before lunch
            on Saturday.
          </p>
        </>
      }
    >
      <Outcome>
        <p>
          A deployed multi-tenant application with sign-up, a subscription, and one feature that
          does something worth paying for.
        </p>
        <p>
          The two irreversible decisions made deliberately rather than discovered later: tenancy
          and the billing model.
        </p>
        <p>A list of what you deliberately did not build, which is most of the plan.</p>
      </Outcome>

      <Section id="rules" title="Three rules for the weekend">
        <Checklist
          items={[
            <>
              <strong>One feature.</strong> Not a platform. The thing somebody would pay for,
              working, and nothing beside it. If you cannot say it in one sentence, you are not
              ready to start building.
            </>,
            <>
              <strong>Nothing bespoke that exists already.</strong> Auth, billing, email and file
              uploads are solved. Building them is how the weekend disappears, and yours will be
              worse.
            </>,
            <>
              <strong>Deploy on Saturday night,</strong> not Sunday evening. A thing that is not
              deployed is not finished, and discovering a deployment problem with two hours left
              is how a weekend produces nothing.
            </>,
          ]}
        />
      </Section>

      <Section id="decisions" title="The two decisions to make before writing anything">
        <H3>1. What is a tenant, and is it a person or an organisation?</H3>
        <p>
          This is the decision that is nearly impossible to reverse. If accounts belong to
          individual users, and in six months a customer asks to add three colleagues, you are
          retrofitting an organisation layer through every query, every permission check and every
          billing record in the system.
        </p>
        <p>
          <strong>Build organisations from the start, even with one member each.</strong> A user
          belongs to an organisation, data belongs to the organisation, and the subscription
          belongs to the organisation. The cost today is one extra table and one extra column. The
          cost later is a migration you will put off until it is a quarter of work.
        </p>
        <CodeBlock
          code={`type Organisation struct {
	ID   uint   \`gorm:"primaryKey"\`
	Name string \`gorm:"not null"\`
	Plan string \`gorm:"not null;default:'free'"\`
}

type User struct {
	ID             uint \`gorm:"primaryKey"\`
	OrganisationID uint \`gorm:"index;not null"\`   // never nullable
	Email          string \`gorm:"uniqueIndex;not null"\`
	Role           string \`gorm:"not null;default:'owner'"\`
}

// And on everything else you store:
type Project struct {
	ID             uint \`gorm:"primaryKey"\`
	OrganisationID uint \`gorm:"index;not null"\`
	Name           string
}`}
        />
        <Warn title="The scope goes in the query, from the first query">
          <p>
            Every read and every write is scoped to the organisation, in the <code>WHERE</code>{' '}
            clause rather than in a check afterwards. Add it later and you will miss one, and the
            one you miss is a customer seeing another customer&apos;s data, which ends the
            business you spent the weekend starting.
          </p>
          <CodeBlock
            code={`db.Where("organisation_id = ?", orgID).Find(&projects)
db.Where("id = ? AND organisation_id = ?", id, orgID).First(&project)`}
          />
        </Warn>
        <H3>2. What are you charging for?</H3>
        <p>
          Per seat, per usage, or a flat tier. Decide now, because it determines what you must
          count from the first write, and counting retroactively is impossible: the events are
          gone.
        </p>
        <Checklist
          items={[
            <>
              <strong>Flat tiers</strong> are the simplest to build and the simplest to explain.
              Start here unless you have a reason.
            </>,
            <>
              <strong>Per seat</strong> needs a member count and a webhook that updates the
              subscription quantity when it changes.
            </>,
            <>
              <strong>Per usage</strong> needs a usage table from day one, written in the same
              transaction as the work. Add it later and your first invoices are guesses.
            </>,
          ]}
        />
      </Section>

      <Step n={1} title="Saturday morning: the foundation">
        <p>
          Two hours, and nothing in it is specific to your product. This is exactly the part worth
          not writing by hand.
        </p>
        <CodeBlock
          terminal
          code={`grit new myapp
cd myapp
docker compose up -d        # Postgres, Redis, MinIO, Mailhog
grit add multitenant        # organisations, membership, scoped queries
grit migrate
grit seed
grit start`}
        />
        <p>
          You now have sign-up, sign-in, password reset, email verification, sessions with
          per-device revoke, an admin panel and organisation scoping. Running, on your machine,
          before coffee.
        </p>
        <p>
          <strong>Checkpoint.</strong> Register two accounts in two different browsers, create
          something in each, and confirm neither can see the other&apos;s. That is the isolation
          test, and it is the only one that matters today. Repeat it on Sunday after you have
          written your own endpoints.
        </p>
        <Note title="If you would rather assemble it yourself">
          <p>
            The same two hours by hand is a Gin project, GORM, bcrypt, a JWT library, refresh
            token rotation, an email provider, a compose file and the organisation model. Each
            step is covered in these guides, and the point of the plan stands either way: this is
            foundation work that is identical in every SaaS, and the weekend should be spent
            elsewhere.
          </p>
        </Note>
      </Step>

      <Step n={2} title="Saturday midday: your one feature">
        <p>
          The part only you can write. Model it as a resource, generate the scaffolding, then
          spend your time on the logic that is actually yours.
        </p>
        <CodeBlock
          terminal
          code={`grit generate resource Project \\
  --fields "name:string:required,status:enum(active|paused|done),due:date,notes:text" \\
  --owned-by organisation`}
        />
        <p>
          That gives you the model, the API with pagination and filtering, the TypeScript types,
          the React hooks and an admin screen, all scoped to the organisation. The hours you have
          left go on the behaviour that makes somebody pay: the calculation, the integration, the
          thing nobody else does.
        </p>
        <p>
          Resist adding the second feature. A product with one feature that works is sellable; a
          product with three half-features is a demo.
        </p>
      </Step>

      <Step n={3} title="Saturday afternoon: billing">
        <p>Two hours, and most of it is the webhook rather than the checkout.</p>
        <CodeBlock
          terminal
          code={`grit plugin add stripe`}
        />
        <Checklist
          items={[
            <>
              <strong>Checkout.</strong> Use the hosted page. Card fields on your own domain mean
              a compliance conversation you do not want on day one.
            </>,
            <>
              <strong>The webhook is the source of truth.</strong> Not the redirect back from
              checkout, which the user can close, refresh or never reach. Payment state changes
              when Stripe says so.
            </>,
            <>
              <strong>Verify the signature,</strong> with a constant-time compare, and store the
              event id so a redelivery is ignored rather than processed twice. Stripe retries, and
              it will deliver the same event more than once.
            </>,
            <>
              <strong>Enforce the plan in the API,</strong> not in the interface. A limit that
              lives in a React component is a limit anybody can skip with curl.
            </>,
          ]}
        />
        <CodeBlock
          code={`func (s *BillingService) Allow(ctx context.Context, orgID uint, feature string) error {
	org, err := s.orgs.ByID(ctx, orgID)
	if err != nil {
		return fmt.Errorf("load org: %w", err)
	}
	// The check is here, in the service, so every caller gets it: the
	// API, the admin, a scheduled job and whatever you add next month.
	if !plans[org.Plan].Includes(feature) {
		return ErrPlanRequired
	}
	return nil
}`}
        />
        <p>
          Use Stripe&apos;s test mode and their CLI to replay webhooks at your local machine. You
          can exercise the entire subscription lifecycle, including cancellation and failed
          payment, without a real card.
        </p>
      </Step>

      <Step n={4} title="Saturday evening: deploy">
        <p>
          Before you are tired, and before anything depends on it. A deploy that works tonight is
          a deploy you can repeat on Sunday without thinking.
        </p>
        <CodeBlock
          terminal
          code={`grit deploy --host myapp.com`}
        />
        <Checklist
          items={[
            <>
              A real domain, with HTTPS. A product on an IP address does not get signed up to.
            </>,
            <>
              A sending domain for email, with SPF and DKIM configured. Without them your welcome
              email goes to spam and you conclude nobody wanted the product.
            </>,
            <>
              Error reporting switched on, so Sunday&apos;s crashes are visible rather than
              anecdotal.
            </>,
            <>
              A database backup, taken and <strong>restored once</strong> to confirm it is real.
            </>,
          ]}
        />
      </Step>

      <Step n={5} title="Sunday morning: the landing page">
        <p>
          Nobody signs up to an application. They sign up to a promise, and the landing page is
          the promise. Give it real hours; it is not decoration.
        </p>
        <Checklist
          items={[
            <>
              <strong>One sentence saying what it does,</strong> for somebody who has never heard
              of you. Not what it is built with, not your mission.
            </>,
            <>
              <strong>A screenshot of the actual product.</strong> Real screens, with plausible
              data. People look at the picture before they read anything.
            </>,
            <>
              <strong>Three benefits, in the customer&apos;s words,</strong> taken from how they
              describe the problem rather than how you describe the solution.
            </>,
            <>
              <strong>Prices, visible.</strong> &quot;Contact us&quot; at this stage means nobody
              contacts you.
            </>,
            <>
              <strong>One call to action,</strong> repeated. Not four competing ones.
            </>,
          ]}
        />
        <Warn title="Never invent a number">
          <p>
            No fabricated user counts, no invented testimonials, no logos of companies that are
            not customers. It is the one mistake on a landing page that is not a matter of taste:
            it is a claim the business has to stand behind, and it is the first thing a sharp
            visitor checks. Leave the social proof section out until you have some.
          </p>
        </Warn>
      </Step>

      <Step n={6} title="Sunday afternoon: the things that lose customers">
        <p>
          Not features. These are the gaps that make a working product feel unfinished, and each
          is under an hour.
        </p>
        <Checklist
          items={[
            <>
              <strong>The empty state for a new account.</strong> The first screen every single
              customer sees, and usually the last one anybody designs. It should say what to do
              first, with a button that does it.
            </>,
            <>
              <strong>The welcome email.</strong> Sent, arriving, not in spam, and saying one
              useful thing rather than &quot;thanks for signing up&quot;.
            </>,
            <>
              <strong>Error states that are honest.</strong> &quot;Something went wrong&quot; with
              no next step is worse than the error.
            </>,
            <>
              <strong>A way to contact you,</strong> visible. One email address is enough.
            </>,
            <>
              <strong>Working cancellation.</strong> People check before they subscribe, and a
              subscription you cannot cancel is a chargeback.
            </>,
            <>
              <strong>Mobile.</strong> Open every screen at 375 pixels. Half your first visitors
              arrive on a phone, from a link somebody sent them.
            </>,
          ]}
        />
      </Step>

      <Section id="skip" title="What to deliberately not build">
        <p>Each of these is a day, and none changes whether the first ten customers appear.</p>
        <Checklist
          items={[
            <>
              <strong>Social login.</strong> Email and password is fine. Add Google when somebody
              asks.
            </>,
            <>
              <strong>A roles and permissions system.</strong> Owner and member is enough for an
              organisation of one to five.
            </>,
            <>
              <strong>An onboarding tour.</strong> A good empty state does more, in a tenth of the
              time.
            </>,
            <>
              <strong>Dark mode,</strong> unless it comes free with what you are using.
            </>,
            <>
              <strong>An analytics dashboard</strong> for customers who have no data yet.
            </>,
            <>
              <strong>A public API,</strong> before anybody has asked for one.
            </>,
            <>
              <strong>Kubernetes.</strong> One server. One.
            </>,
          ]}
        />
        <Note title="The test for a feature this weekend">
          <p>
            Would its absence stop somebody paying? If no, it is next weekend&apos;s work. Almost
            everything fails that test, which is why the plan fits in two days.
          </p>
        </Note>
      </Section>

      <Section id="monday" title="Monday">
        <p>
          The weekend produced a product. It has not produced a customer, and the next bit is
          harder than the building was.
        </p>
        <Checklist
          items={[
            <>
              Show it to ten people who have the problem. Not developers, unless developers are
              the customer. Watch them use it without helping, which is uncomfortable and is where
              the real list comes from.
            </>,
            <>
              Charge from the first customer. A free product tells you nothing about whether it is
              worth money, and converting free users later is a different and harder product.
            </>,
            <>
              Fix the three things everybody trips over before building the fourth feature.
            </>,
          ]}
        />
      </Section>

      <Faq
        slug="build-saas-mvp-weekend"
        items={[
          {
            q: 'Is a weekend realistic?',
            a: (
              <>
                For a focused product with one feature, yes, if you do not build authentication,
                billing or a tenancy model from scratch. People who fail at this spend Saturday
                on login and Sunday on deployment, and never reach the feature. The weekend is
                realistic precisely because most of what a SaaS needs is identical in every SaaS.
              </>
            ),
            text:
              'For a focused single-feature product, yes, provided you do not build authentication, billing or tenancy from scratch. The common failure is spending Saturday on login and Sunday on deployment and never reaching the feature. It is realistic because most of what a SaaS needs is identical in every SaaS.',
          },
          {
            q: 'Organisations or individual user accounts?',
            a: (
              <>
                Organisations, always, even when every organisation has one member. Retrofitting a
                tenancy layer means touching every query, every permission check and every billing
                record in the system. Building it on day one costs one table and one column.
              </>
            ),
            text:
              'Organisations, even when each has one member. Retrofitting tenancy means touching every query, permission check and billing record in the system, while building it on day one costs one table and one column.',
          },
          {
            q: 'How do I keep one customer from seeing another customer’s data?',
            a: (
              <>
                Put the organisation id in the <code>WHERE</code> clause of every query rather
                than checking after the fetch, so a row outside the scope cannot be returned at
                all. Then test it: two accounts, two browsers, and an attempt to reach the
                other&apos;s records by id. Run that test again whenever you add an endpoint,
                especially an export or a bulk action, which are the ones that get written without
                the scope.
              </>
            ),
            text:
              'Put the organisation id in the WHERE clause of every query rather than checking after the fetch, so an out-of-scope row cannot be returned. Then test it with two accounts in two browsers, trying to reach the other by id, and repeat whenever you add an endpoint. Exports and bulk actions are the ones usually written without the scope.',
          },
          {
            q: 'Should the MVP be free to start with?',
            a: (
              <>
                Charge from the first customer. A free product measures curiosity, not demand, and
                the feedback from people paying nothing is systematically different from the
                feedback from people paying. A free trial with a card up front is fine; free
                forever is a different product, and a harder one.
              </>
            ),
            text:
              'Charge from the first customer. A free product measures curiosity rather than demand, and feedback from non-paying users is systematically different. A free trial with a card up front is fine; free forever is a different and harder product.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/deploy-go-app-vps-https',
            label: 'Deploy to a VPS with HTTPS',
            blurb: 'Saturday evening, in detail.',
          },
          {
            href: '/guides/rbac-explained-golang',
            label: 'Roles and permissions',
            blurb: 'For when the first customer asks to add a colleague.',
          },
          {
            href: '/docs/plugins/multitenant',
            label: 'Multitenancy',
            blurb: 'Organisations, membership and scoped queries as generated code.',
          },
          {
            href: '/blueprints',
            label: 'Blueprints',
            blurb: 'Whole applications to start from, when the shape is already known.',
          },
        ]}
      />
    </GuidePage>
  )
}
