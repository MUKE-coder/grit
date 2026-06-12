import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Triple is the SaaS shape — three apps in one monorepo. It&apos;s the most
        common mode you&apos;ll use, the one this course&apos;s scaffolded project is
        built in, and the recommended starting point for almost every
        multi-tenant product.
      </p>

      <h2>The three apps</h2>
      <CodeBlock
        language="text"
        code={`apps/web/      Public site + customer dashboard. Built with Next.js.
apps/admin/    Filament-style admin panel. Built with Next.js.
apps/api/      Go API serving both. Gin + GORM.`}
      />

      <h2>Who uses each surface?</h2>
      <ul>
        <li>
          <strong>web</strong> — your customers. Marketing pages + the
          logged-in dashboard.
        </li>
        <li>
          <strong>admin</strong> — your staff. Sales sees customers, support
          edits orders, ops triggers refunds.
        </li>
        <li>
          <strong>api</strong> — both of the above + future mobile + desktop.
          One backend, multiple surfaces.
        </li>
      </ul>

      <h2>Why split web and admin?</h2>
      <p>Separate domains, separate concerns:</p>
      <ul>
        <li>
          <strong>Different auth scopes.</strong> Customers can&apos;t log into the
          admin panel; staff doesn&apos;t need to sign up for customer accounts.
        </li>
        <li>
          <strong>Different SEO needs.</strong> The public site is
          marketing-optimized; the admin is behind auth, no SEO.
        </li>
        <li>
          <strong>Different shipping cadence.</strong> Marketing changes
          weekly; admin features change on a different rhythm.
        </li>
        <li>
          <strong>Different attack surface.</strong> Admin can sit behind a
          VPN or IP allow-list; the public site can&apos;t.
        </li>
      </ul>

      <h2>What&apos;s shared</h2>
      <p>One thing matters most: <strong>packages/shared</strong>.</p>
      <CodeBlock
        language="text"
        code={`packages/shared/src/schemas/   Zod schemas — web and admin both validate the same way.
packages/shared/src/types/     Generated TS types — both apps know what a Product looks like.
packages/shared/src/constants/ Route paths, enum values, shared constants.`}
      />
      <p>
        Both Next.js apps import <code>@workspace/shared</code> and they
        always agree on schemas + types. No drift.
      </p>

      <TipBox tone="info">
        <strong>The admin panel writes itself.</strong> When you generate a
        resource, <code>defineResource()</code> in{' '}
        <code>apps/admin/app/resources/</code> auto-renders a CRUD page —
        you don&apos;t write HTML or form code. The deep dive on this is in the
        Web (Next.js) course.
      </TipBox>

      <h2>Deployment shape</h2>
      <p>You can ship triple in two ways:</p>
      <ol>
        <li>
          <strong>All three on one box</strong> with a reverse proxy
          (Traefik) routing <code>web.example.com</code> → apps/web,{' '}
          <code>admin.example.com</code> → apps/admin,{' '}
          <code>api.example.com</code> → apps/api. Simplest for small teams.
        </li>
        <li>
          <strong>Apps deployed independently.</strong> apps/web on Vercel,
          apps/admin on a private box, apps/api on a VPS. Most flexible.
        </li>
      </ol>

      <KnowledgeCheck
        question="Your customer wants a feature where staff can refund any order. Where do you put each part of it?"
        choices={[
          {
            label: 'All in apps/web — staff log in there too',
            feedback:
              "Wrong — staff and customers should NOT use the same auth surface. That's why apps/admin exists.",
          },
          {
            label: 'Handler in apps/api, page in apps/admin, audit-log model in apps/api',
            correct: true,
            feedback:
              "Right — POST /api/orders/:id/refund handler in apps/api, the 'Refund' button in apps/admin/app/resources/orders, the audit log row in apps/api models.",
          },
          {
            label: 'Refund logic in apps/web, exposed only when role=staff',
            feedback:
              "Wrong — business logic doesn't live in the frontend. Refund processing belongs in a service in apps/api.",
          },
          {
            label: 'Stripe handles it — no Grit code needed',
            feedback:
              'Stripe processes the money, but YOUR app needs to record the refund, audit who did it, and update the order status. That logic lives in apps/api.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Look at your scaffolded triple project. Pick one feature
              (signup, login, the admin&apos;s users page) and trace it across all
              three apps. In <code>notes.md</code>, list:
            </p>
            <ul>
              <li>Where in <code>apps/web</code> is the UI?</li>
              <li>Where in <code>apps/admin</code> is the matching admin view?</li>
              <li>Where in <code>apps/api</code> is the handler + service?</li>
              <li>Where in <code>packages/shared</code> is the Zod schema?</li>
            </ul>
          </>
        }
        hint={<>Signup is the simplest one to trace.</>}
        solution={
          <>
            <p>For signup:</p>
            <ul>
              <li>apps/web/app/(auth)/signup/page.tsx</li>
              <li>apps/admin/app/resources/users/page.tsx (admin lists users)</li>
              <li>apps/api/internal/handlers/auth.go + internal/services/auth.go</li>
              <li>packages/shared/src/schemas/user.ts (SignupSchema)</li>
            </ul>
            <p>
              Four files, three apps, one Zod schema validating signup in
              all of them. That&apos;s the triple mode flow.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Single and triple cover the two extremes. Next we cover the three
        specialized modes — API-only, Mobile, Desktop — for when you want
        one slice instead of the full triple.
      </p>
    </>
  )
}
