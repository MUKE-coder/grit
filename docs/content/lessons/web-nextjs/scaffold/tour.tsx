import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Tour all three apps. We already covered <code>apps/api</code> in
        the Go API course; this lesson focuses on <code>apps/web</code> and{' '}
        <code>apps/admin</code> — what&apos;s the same, what&apos;s different.
      </p>

      <h2>apps/web — the public Next.js site</h2>
      <CodeBlock
        language="text"
        filename="apps/web/"
        code={`apps/web/
├── app/
│   ├── (marketing)/         public pages — landing, pricing, blog
│   │   ├── page.tsx
│   │   ├── pricing/page.tsx
│   │   └── layout.tsx       header + footer
│   ├── (auth)/              login, signup, forgot-password
│   ├── (app)/               logged-in pages — dashboard, settings
│   │   ├── dashboard/page.tsx
│   │   └── layout.tsx       sidebar + topbar
│   ├── api/                 Next.js route handlers (mostly empty)
│   └── layout.tsx           root layout — theme, providers
├── components/              reusable React components
├── hooks/                   useUsers, useCurrentUser, …
├── lib/                     api client, auth helpers, utils
└── public/                  static assets`}
      />
      <p>
        Three route groups: <code>(marketing)</code>, <code>(auth)</code>,{' '}
        <code>(app)</code>. Each has its own layout. The parens don&apos;t
        appear in URLs.
      </p>

      <h2>apps/admin — the staff panel</h2>
      <CodeBlock
        language="text"
        filename="apps/admin/"
        code={`apps/admin/
├── app/
│   ├── (auth)/login/page.tsx
│   ├── dashboard/page.tsx
│   ├── resources/
│   │   ├── users/page.tsx       defineResource() — generates CRUD
│   │   └── products/page.tsx
│   └── layout.tsx               admin shell — sidebar, top nav
├── components/
│   ├── admin/                   DataTable, FormBuilder, etc.
│   └── ui/                      shadcn primitives
├── hooks/
└── lib/`}
      />
      <p>
        The killer file is <code>app/resources/<em>x</em>/page.tsx</code> —
        we cover <code>defineResource()</code> in chapter 4. Almost every
        admin page is one of these.
      </p>

      <h2>What&apos;s the same across both apps</h2>
      <ul>
        <li>
          Both use Next.js 14 App Router (no Pages Router)
        </li>
        <li>
          Both import types from <code>@workspace/shared</code>
        </li>
        <li>
          Both use React Query for data
        </li>
        <li>
          Both use shadcn/ui + Tailwind CSS
        </li>
      </ul>

      <h2>What differs</h2>
      <ul>
        <li>
          <strong>Auth scope</strong> — web cookies are scoped to{' '}
          <code>customer-token</code>; admin to <code>staff-token</code>.
          One can&apos;t accidentally log into the other.
        </li>
        <li>
          <strong>Layout</strong> — web is marketing-shaped (header / hero /
          footer); admin is dashboard-shaped (sidebar / topbar).
        </li>
        <li>
          <strong>Robots</strong> — admin has{' '}
          <code>noindex,nofollow</code> on every page.
        </li>
      </ul>

      <TipBox tone="info">
        <strong>The shared API client trick:</strong> both apps import the
        same <code>lib/api.ts</code> shape but configure it with different
        cookie names. We&apos;ll see this in chapter 3.
      </TipBox>

      <KnowledgeCheck
        question="A teammate suggests deleting apps/admin and gating admin URLs inside apps/web by role. Trade-off check — what's the biggest downside?"
        choices={[
          {
            label: 'Larger bundle for marketing users',
            feedback:
              "True but solvable with code-splitting. Not the biggest issue.",
          },
          {
            label: 'Sharing one auth cookie between customer + staff means a customer XSS / token theft compromises staff too',
            correct: true,
            feedback:
              "Right — same cookie scope means same blast radius. Two apps = two cookie scopes = compromise of one doesn't grant access to the other.",
          },
          {
            label: 'You\'d need to learn another framework',
            feedback:
              "Wrong — both apps are Next.js.",
          },
          {
            label: 'SEO would suffer',
            feedback:
              "Admin shouldn't be indexed anyway. SEO isn't the main concern; security separation is.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In your scaffolded project, open both <code>apps/web/app/layout.tsx</code> and{' '}
              <code>apps/admin/app/layout.tsx</code>. In <code>notes.md</code>, write:
            </p>
            <ul>
              <li>What providers does each app wrap children in?</li>
              <li>What&apos;s in each app&apos;s root metadata (title, OG, etc.)?</li>
              <li>What one component or import is DIFFERENT between them?</li>
            </ul>
          </>
        }
        hint={
          <>
            The differing piece is usually the auth provider (different cookie
            name) and the shell layout (marketing vs. dashboard).
          </>
        }
        solution={
          <>
            <p>You should find:</p>
            <ul>
              <li>Both: QueryClient + ThemeProvider + ToastProvider</li>
              <li>Web layout: marketing header + footer; admin: sidebar + topbar</li>
              <li>Different cookie name in the AuthProvider</li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of the chapter — <strong>packages/shared</strong>, the
        glue between the three apps.
      </p>
    </>
  )
}
