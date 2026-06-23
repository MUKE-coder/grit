import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The admin panel auto-protects its dashboard routes — visit one
        without a session and you bounce to <code>/login</code>{' '}
        automatically. The customer-facing web app at{' '}
        <code>apps/web/</code> doesn&apos;t do this by default;
        it&apos;s designed to host public marketing pages alongside
        member-only areas, so blanket protection would break the
        marketing side. <strong><code>grit add web-auth</code></strong>{' '}
        (shipped in <strong>v3.31.22</strong>) gives you a one-shot
        command that drops two opt-in tools into the web app: a
        cookie-checking middleware and a client-side wrapper.
      </p>

      <h2>Running the command</h2>

      <CodeBlock
        terminal
        code={`grit add web-auth`}
      />

      <p>You&apos;ll see:</p>

      <CodeBlock
        language="text"
        code={`  Adding web-auth helpers to apps/web/
  ✓ wrote apps/web/middleware.ts
  ✓ wrote apps/web/components/ProtectedWebRoute.tsx

  Next steps:
    1. Open apps/web/middleware.ts and add protected paths to the matcher.
    2. Or wrap a page client-side with <ProtectedWebRoute>.
    3. See /docs/concepts/protecting-web-pages for both patterns.`}
      />

      <p>
        Re-run it any time — it&apos;s idempotent. Existing files are
        skipped with a notice unless you pass <code>--force</code>.
      </p>

      <h2>Two patterns, one command</h2>
      <p>
        Web-auth is intentionally not a single mechanism. Different
        pages need different protection guarantees, and forcing one
        approach onto everything leads to bad trade-offs:
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Pattern</th>
              <th className="text-left px-3 py-2 font-medium">Lives in</th>
              <th className="text-left px-3 py-2 font-medium">When to use</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 font-medium">Middleware (SSR)</td>
              <td className="px-3 py-2 font-mono text-[12px]">middleware.ts</td>
              <td className="px-3 py-2 text-[12px]">Whole-section gates: <code>/account</code>, <code>/checkout</code>. Fast, no flash.</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-medium">Wrapper (client)</td>
              <td className="px-3 py-2 font-mono text-[12px]">&lt;ProtectedWebRoute&gt;</td>
              <td className="px-3 py-2 text-[12px]">Role-gated content. Per-page protection without editing middleware.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <h2>Pattern 1 — Middleware (SSR cookie gate)</h2>
      <p>
        Open <code>apps/web/middleware.ts</code> and you&apos;ll find
        two arrays at the top:
      </p>

      <CodeBlock
        language="ts"
        filename="apps/web/middleware.ts (excerpt)"
        code={`// Add paths here that require an authenticated visitor.
const PROTECTED_PATHS: string[] = [
  "/account",
  "/account/:path*",
  // "/checkout",
  // "/dashboard",
];

// Paths that should be inaccessible to ALREADY-signed-in users
// (the login form, sign-up). Sends them to /account instead.
const AUTH_PATHS: string[] = [
  "/login",
  "/register",
  "/forgot-password",
];`}
      />

      <p>
        The middleware runs on every Next.js request matching its
        config matcher. It checks for the <code>grit_access</code>{' '}
        HttpOnly cookie:
      </p>
      <ul>
        <li>
          <strong>No cookie + protected path</strong> → redirect to{' '}
          <code>/login?next=&lt;original-path&gt;</code>. After login
          succeeds, the auth flow can read <code>?next=</code> and send
          the visitor back to where they were headed.
        </li>
        <li>
          <strong>Has cookie + auth path</strong> → bounce to{' '}
          <code>/account</code>. No reason to show the login form to
          someone who&apos;s already signed in.
        </li>
        <li>
          <strong>Anything else</strong> → pass through unchanged.
        </li>
      </ul>

      <TipBox tone="info">
        <strong>Why not verify the JWT?</strong> Verifying the cookie
        would require a network call to the API on every page request
        — adding 50-200ms of latency to your homepage too. The
        middleware checks <em>cookie existence</em>, which is
        instant. Invalid or expired cookies still bounce — but at
        the next API call (which fails with 401, which is caught by{' '}
        <code>useMe()</code>, which redirects). The middleware
        prevents the unauthenticated-visitor-sees-account-shell flash
        without the per-request cost.
      </TipBox>

      <h3>Editing the matcher</h3>
      <p>
        The <code>config.matcher</code> at the bottom of{' '}
        <code>middleware.ts</code> tells Next.js which paths to run the
        middleware on. Keep it in sync with <code>PROTECTED_PATHS</code>{' '}
        + <code>AUTH_PATHS</code> — anything outside the matcher
        bypasses the middleware entirely (good for performance on
        static assets, blog pages, marketing routes):
      </p>

      <CodeBlock
        language="ts"
        code={`export const config = {
  matcher: [
    "/account/:path*",
    "/checkout",          // ← add new paths here when you add them to PROTECTED_PATHS
    "/login",
    "/register",
    "/forgot-password",
  ],
};`}
      />

      <h2>Pattern 2 — ProtectedWebRoute (client wrapper)</h2>
      <p>
        The middleware is fast but can only check cookie existence.
        Three cases it can&apos;t handle:
      </p>
      <ul>
        <li>
          <strong>Expired cookies</strong> — the cookie is present
          (middleware passes through) but invalid. The page shell
          renders before the first API call exposes the problem.
        </li>
        <li>
          <strong>Role gates</strong> — &quot;this page is only for
          ADMIN users.&quot; The cookie doesn&apos;t carry the role; you
          need a real <code>/api/auth/me</code> probe.
        </li>
        <li>
          <strong>Per-page protection</strong> — you don&apos;t want
          to add yet another path to the middleware matcher.
        </li>
      </ul>

      <CodeBlock
        language="tsx"
        filename="apps/web/app/billing/page.tsx"
        code={`import { ProtectedWebRoute } from "@/components/ProtectedWebRoute";

export default function BillingPage() {
  return (
    <ProtectedWebRoute>
      <h1>Billing</h1>
      <BillingDetails />
    </ProtectedWebRoute>
  );
}`}
      />

      <p>
        The wrapper calls <code>useMe()</code>, shows a spinner while
        the probe is in flight, redirects to <code>/login?next=…</code>{' '}
        on null user, and renders children when authenticated.
      </p>

      <h3>Role gating</h3>
      <p>
        Pass a <code>roles</code> prop to restrict to specific roles.
        Visitors with a session but the wrong role are redirected to{' '}
        <code>/</code> (the landing page) — not <code>/login</code>,
        since they&apos;re already signed in.
      </p>

      <CodeBlock
        language="tsx"
        code={`<ProtectedWebRoute roles={["ADMIN", "EDITOR"]}>
  <CMSPage />
</ProtectedWebRoute>

{/* Single role works too — pass a string instead of an array */}
<ProtectedWebRoute roles="ADMIN">
  <SuperUserOnly />
</ProtectedWebRoute>`}
      />

      <h3>Custom loading view</h3>
      <p>
        The default spinner is intentionally bland. Pass a{' '}
        <code>fallback</code> prop to replace it with a skeleton that
        matches the page&apos;s real layout:
      </p>

      <CodeBlock
        language="tsx"
        code={`<ProtectedWebRoute
  fallback={
    <div className="space-y-3 p-8">
      <div className="h-8 w-48 rounded bg-slate-200 animate-pulse" />
      <div className="h-32 rounded bg-slate-100 animate-pulse" />
    </div>
  }
>
  <AccountPage />
</ProtectedWebRoute>`}
      />

      <h2>Combining both patterns</h2>
      <p>
        Middleware and wrapper aren&apos;t alternatives — they layer:
      </p>

      <CodeBlock
        language="text"
        code={`HTTP request → /account/billing
                  │
                  ▼
        middleware.ts checks grit_access cookie
            │ missing? → redirect /login?next=/account/billing  ✗
            │ present? → pass through                            ✓
                  │
                  ▼
        page hydrates: <ProtectedWebRoute roles="ADMIN">
            │ useMe() probes /api/auth/me
            │   200 + correct role  → render <BillingPage />     ✓
            │   200 + wrong role    → redirect /                 ✗
            │   401                 → redirect /login            ✗
            │`}
      />

      <p>
        The middleware catches the "no cookie at all" case (which is
        99% of unauthorized visitors), shaving network calls and flash.
        The wrapper catches the rare "cookie present but
        invalid/expired/wrong-role" cases.
      </p>

      <KnowledgeCheck
        question="You added `/checkout` to PROTECTED_PATHS in middleware.ts but forgot to add it to the config.matcher at the bottom. What happens when a logged-out visitor hits /checkout?"
        choices={[
          {
            label: "They're redirected to /login, as intended",
            feedback:
              "Wrong — the middleware never RUNS on a path that isn't in the matcher. Next.js skips the file entirely for unmatched paths.",
          },
          {
            label: "The page renders normally — middleware was skipped because /checkout isn't in the matcher",
            correct: true,
            feedback:
              "Right. The matcher is the gate that decides whether middleware runs at all. PROTECTED_PATHS only matters once we're inside middleware(). Always update both in tandem.",
          },
          {
            label: "Next.js throws a build error",
            feedback:
              "Wrong — the matcher and the PROTECTED_PATHS array are independent lists; mismatches are silent.",
          },
          {
            label: "The visitor sees an infinite redirect loop",
            feedback:
              "Wrong — the middleware would only redirect to /login if it ran, which it doesn't if the path isn't matched.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Add web-auth to your contact-app and protect a new
              page end-to-end:
            </p>
            <ol>
              <li>
                Run <code>grit add web-auth</code>.
              </li>
              <li>
                Create a new page at{' '}
                <code>apps/web/app/account/page.tsx</code> that
                renders a simple heading.
              </li>
              <li>
                Add <code>/account</code> to both{' '}
                <code>PROTECTED_PATHS</code> and{' '}
                <code>config.matcher</code> in middleware.ts.
              </li>
              <li>
                In an incognito tab, visit{' '}
                <code>http://localhost:3000/account</code> — you
                should bounce to{' '}
                <code>/login?next=%2Faccount</code>.
              </li>
              <li>
                Sign in. After redirect, you should land back on{' '}
                <code>/account</code> (if your login page reads the{' '}
                <code>next</code> param) or on the default
                post-login page.
              </li>
              <li>
                Wrap the page contents in{' '}
                <code>&lt;ProtectedWebRoute roles=&quot;ADMIN&quot;&gt;</code>.
                Visit again. If your seed user is ADMIN, you should
                see the page. If you make your role something else
                (e.g. via the admin user list), you should bounce to{' '}
                <code>/</code>.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For step 5 to work, the login page must read{' '}
            <code>?next=</code> from the URL and pass it to the
            success redirect. Default Grit login pages do this. If
            yours doesn&apos;t, hand-edit{' '}
            <code>apps/web/app/(auth)/login/page.tsx</code> to add{' '}
            <code>const next = useSearchParams().get(&quot;next&quot;) || &quot;/&quot;</code>{' '}
            and use it as the redirect target.
          </>
        }
        solution={
          <>
            <p>
              <code>apps/web/app/account/page.tsx</code>:
            </p>
            <CodeBlock
              language="tsx"
              code={`import { ProtectedWebRoute } from "@/components/ProtectedWebRoute";

export default function AccountPage() {
  return (
    <ProtectedWebRoute roles="ADMIN">
      <main className="p-8">
        <h1 className="text-2xl font-semibold">Welcome to your account</h1>
        <p className="text-slate-600 mt-2">
          Admin-only content lives here.
        </p>
      </main>
    </ProtectedWebRoute>
  );
}`}
            />
            <p>
              The middleware bounces unauthenticated visitors at the
              edge (no flash). The wrapper enforces the ADMIN role
              after hydration. The two work together — neither does
              the other&apos;s job.
            </p>
          </>
        }
      />

      <h2>Chapter recap</h2>
      <p>
        This is the last lesson in the resource-lifecycle chapter. You
        can now:
      </p>
      <ul>
        <li>Generate full-stack resources (<code>grit generate</code>)</li>
        <li>Sync types and propagate Go changes to TypeScript (<code>grit sync</code>)</li>
        <li>Remove resources cleanly (<code>grit remove</code>)</li>
        <li>Customise the admin form (render mode, groups + PATCH, field types)</li>
        <li>Customise the admin table (formats, packed columns, filters)</li>
        <li>Expose forms and tables to the web app (<code>grit expose</code>)</li>
        <li>Share resources publicly with a token (FormShare + /forms/[token])</li>
        <li>Protect customer-facing pages (<code>grit add web-auth</code>)</li>
      </ul>
      <p>
        A complete model lifecycle from Go struct to authenticated
        customer-facing page, all from a single resource definition.
        That&apos;s the whole pitch.
      </p>
    </>
  )
}
