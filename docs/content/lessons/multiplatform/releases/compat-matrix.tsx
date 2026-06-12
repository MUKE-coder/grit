import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The web app updates instantly. Mobile takes a week to be on
        most users&apos; devices. Desktop&apos;s a slow trickle.
        Your API has to support all three versions simultaneously.
        The compatibility matrix is how you keep track.
      </p>

      <h2>The problem in one paragraph</h2>
      <p>
        You ship the API at v1.0. Web is at v1.0. You add a field,
        deploy API + web to v1.1 together — fine, they updated in
        lockstep. Now mobile pushes a v1.1 build to TestFlight, takes
        24h to review, another week for the user to update. During
        that week, the user is on mobile v1.0 + API v1.1. Did you
        BREAK them, or are they FINE?
      </p>

      <h2>The matrix — write it down</h2>
      <p>
        Maintain a real document. <code>docs/compat-matrix.md</code> at the
        repo root:
      </p>
      <CodeBlock
        language="markdown"
        filename="docs/compat-matrix.md"
        code={`# Client / API compatibility

| API ver | Web ver  | Mobile ver | Desktop ver | Notes                  |
|---------|----------|------------|-------------|------------------------|
| 1.0     | >= 1.0   | >= 1.0     | >= 1.0      | initial release        |
| 1.1     | >= 1.0   | >= 1.0     | >= 1.0      | added optional field   |
| 1.2     | >= 1.2   | >= 1.0     | >= 1.0      | web-only feature flag  |
| 2.0     | >= 2.0   | >= 1.9     | >= 1.9      | BREAKING — see below   |

## 2.0 breaking changes
- Removed deprecated /api/legacy-search
- Renamed Product.code → Product.sku (was deprecated since 1.5)
- Clients must send X-Client-Version header`}
      />
      <p>
        Three things this gives you:
      </p>
      <ul>
        <li>A single source of truth before deploys.</li>
        <li>A way to refuse old clients (via X-Client-Version + middleware).</li>
        <li>A historical record when debugging &quot;why did v1.2 break for some users?&quot;</li>
      </ul>

      <h2>Two API rules that make compat easy</h2>

      <h3>1. Adding fields is fine. Removing them is breaking.</h3>
      <p>
        Old clients ignore new fields. So adding{' '}
        <code>is_featured: bool</code> to <code>/api/products</code>{' '}
        does not break mobile v1.0 — it just doesn&apos;t use the field.
        Removing or renaming a field breaks every client that reads it.
      </p>

      <h3>2. New endpoints are fine. Changed semantics are breaking.</h3>
      <p>
        A new <code>POST /api/search</code> endpoint costs nothing —
        old clients call the old <code>/api/products?q=</code>. But if
        you change what <code>POST /api/products</code> returns (now
        it&apos;s 202 instead of 201), old clients break.
      </p>

      <TipBox tone="info">
        <strong>Add then deprecate, don&apos;t mutate.</strong> When
        you need to rename a field: add the new field, populate both,
        wait for clients to migrate, then remove the old field. This
        is how you ship breaking changes without breaking anyone.
      </TipBox>

      <h2>X-Client-Version — the runtime check</h2>
      <CodeBlock
        language="ts"
        filename="apps/web/lib/api.ts (and same in mobile/desktop)"
        code={`const APP_VERSION = process.env.NEXT_PUBLIC_APP_VERSION  // baked at build time

export async function api(path: string, opts?: RequestInit) {
  return fetch(API_URL + path, {
    ...opts,
    headers: {
      ...opts?.headers,
      'X-Client-Version': APP_VERSION,
      'X-Client-Surface': 'web',  // or 'mobile' / 'desktop'
    },
  })
}`}
      />
      <CodeBlock
        language="go"
        filename="apps/api/internal/middleware/version_gate.go"
        code={`func VersionGate(minVersions map[string]string) gin.HandlerFunc {
  return func(c *gin.Context) {
    surface := c.GetHeader("X-Client-Surface")
    version := c.GetHeader("X-Client-Version")
    min, ok := minVersions[surface]
    if ok && semver.Compare("v"+version, "v"+min) < 0 {
      c.AbortWithStatusJSON(426, gin.H{
        "error": "client_too_old",
        "min_version": min,
        "upgrade_url": "https://app.example.com/upgrade",
      })
      return
    }
    c.Next()
  }
}`}
      />
      <p>
        The HTTP 426 (&quot;Upgrade Required&quot;) is the standard
        status code for &quot;your client is too old&quot;. Each
        surface handles it differently: web reloads; mobile shows a
        forced-update screen; desktop triggers the auto-updater.
      </p>

      <h2>Per-surface upgrade UX</h2>
      <ul>
        <li>
          <strong>Web:</strong> on 426, show a banner &quot;a new
          version is available, refresh&quot;. The next page load
          gets the new JS. Easiest.
        </li>
        <li>
          <strong>Desktop:</strong> on 426, trigger the auto-updater
          (from the Desktop course). User clicks &quot;Update&quot;,
          app restarts in &lt; 30s.
        </li>
        <li>
          <strong>Mobile:</strong> hardest. You can show a forced
          screen with a link to the App Store, but if the app needs
          review, the user is stuck for hours. Plan ahead: keep mobile
          backwards-compat for at least 2-3 minor versions.
        </li>
      </ul>

      <h2>Real example — adding a required field</h2>
      <p>
        Marketing wants every product to have a <code>category</code>.
        How to roll it out without breaking mobile?
      </p>
      <ol>
        <li>
          <strong>API v1.5:</strong> Add the <code>category</code>{' '}
          field as <strong>optional</strong>. Default to empty
          string. Existing rows backfilled to empty.
        </li>
        <li>
          <strong>Web + admin v1.5:</strong> Adds a UI to set it on
          new products.
        </li>
        <li>
          <strong>Mobile v1.5:</strong> Reads it if present, hides
          the section if empty.
        </li>
        <li>
          <strong>API v1.6 (months later):</strong> Make it required
          server-side. By now all active products have a category.
        </li>
      </ol>
      <p>
        Three releases over weeks. Slow, but no client ever 500&apos;d.
      </p>

      <KnowledgeCheck
        question="You ship API v2.0 with a breaking change. Mobile v1.0 users (still 30% of your install base) start hitting 500s. What's the FIRST thing you should do?"
        choices={[
          {
            label: 'Email all users to update',
            feedback:
              "Slow and most won't read it. Email is the long-term fix, not the now-fix.",
          },
          {
            label: 'Roll back API v2.0 immediately. The breaking change should have been gated, deprecated first, or a v1/v2 split.',
            correct: true,
            feedback:
              "Right — never let 30% of users be on 500. Roll back, then plan: bring the v1 endpoints back as a separate /v1/* prefix, deprecate gracefully, push mobile to update, and ship v2 again when adoption is high.",
          },
          {
            label: 'Force-update mobile users via the OS update mechanism',
            feedback:
              "You can't force OS updates — and even if you nudge, the app review delay is days. Roll back first.",
          },
          {
            label: 'Hope they update soon',
            feedback: "Hope isn't a strategy.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Set up a real compatibility gate:</p>
            <ol>
              <li>
                Create <code>docs/compat-matrix.md</code> in your repo
                with the current versions.
              </li>
              <li>
                Add <code>X-Client-Version</code> + <code>X-Client-Surface</code>{' '}
                to every API call from web / admin / mobile / desktop.
              </li>
              <li>
                Wire the VersionGate middleware in the API. Set the
                current min versions in a map.
              </li>
              <li>
                Hand-set web&apos;s version to v0.0.1 in <code>.env</code>.
                Hit the API. Should get a 426 + upgrade JSON.
              </li>
              <li>
                Implement the &quot;new version available&quot; banner
                on web that shows when a 426 comes back.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the web banner, store the 426 state in a React
            Context. Any API call that gets 426 sets a global flag;
            your <code>AppLayout</code> reads it and shows the banner.
          </>
        }
        solution={
          <>
            <p>
              You now have explicit version control between API and
              clients. When you ship a breaking v2.0, you bump the
              min for mobile to 1.9 (not 2.0) — giving real users a
              week or two to update naturally before getting 426&apos;d.
            </p>
            <p>
              The discipline of the matrix is more valuable than the
              code. Most projects fail at multi-platform because
              there&apos;s no shared understanding of who supports
              what.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Final lesson — <strong>Staggered releases</strong>. With the
        compat matrix in hand, the actual deploy order: API first,
        then web, then desktop, then mobile (with reasoning at each
        step).
      </p>
    </>
  )
}
