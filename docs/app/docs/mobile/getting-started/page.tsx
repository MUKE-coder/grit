import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/mobile/getting-started')

export default function MobileGettingStartedPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Mobile</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Getting Started with Mobile
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Scaffold a Go API and an Expo React Native app with one command. This guide
                covers prerequisites, project creation, the development workflow, and the
                device-vs-emulator API URL matrix so your phone can actually reach the backend.
              </p>
            </div>

            {/* Prerequisites */}
            <div className="prose-grit mb-10">
              <h2>Prerequisites</h2>
              <p>
                Make sure the following tools are installed before creating a mobile project.
                Go and Node power the monorepo; the Expo Go app and platform SDKs are what run
                the app on a phone or emulator.
              </p>
            </div>

            <div className="grid gap-3 sm:grid-cols-2 mb-10">
              {[
                { name: 'Go', version: '1.21+', check: 'go version' },
                { name: 'Node.js', version: '18+', check: 'node --version' },
                { name: 'pnpm', version: '8+', check: 'pnpm --version' },
                { name: 'Grit CLI', version: 'Latest', check: 'grit --help' },
                { name: 'Expo Go', version: 'iOS / Android', check: 'App Store / Play Store' },
                { name: 'EAS CLI', version: 'For builds', check: 'npm i -g eas-cli' },
              ].map((tool) => (
                <div
                  key={tool.name}
                  className="rounded-lg border border-border/30 bg-card/30 px-4 py-3"
                >
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-[15px] font-semibold">{tool.name}</span>
                    <span className="text-sm font-mono text-primary/60">{tool.version}</span>
                  </div>
                  <code className="text-sm font-mono text-muted-foreground/50">{tool.check}</code>
                </div>
              ))}
            </div>

            <div className="prose-grit mb-10">
              <p>
                To run the app on a <strong>physical device</strong>, install the free{' '}
                <a href="https://expo.dev/go" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Expo Go</a>{' '}
                app from the App Store or Google Play — no native toolchain required. To run on
                a <strong>simulator/emulator</strong> instead, install Xcode (iOS, macOS only)
                or Android Studio (Android, any OS). You only need the emulators if you want to
                test without a phone; Expo Go on a real device is the fastest path.
              </p>
            </div>

            {/* Step 1 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  1
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Scaffold the Project</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Create a new mobile project with the Grit CLI. This generates a Turborepo
                  monorepo with the full Go backend in <code>apps/api</code>, an Expo React
                  Native app in <code>apps/expo</code>, and shared Zod schemas and TypeScript
                  types in <code>packages/shared</code>.
                </p>
              </div>
              <CodeBlock terminal code="grit new myapp --mobile" className="mb-4 glow-purple-sm" />
              <div className="prose-grit mb-4">
                <p>
                  <code>--mobile</code> is shorthand for <code>--arch=mobile</code>. Both produce
                  the same result. If you already have another architecture and want to{' '}
                  <em>add</em> Expo alongside it, use the <code>--expo</code> flag instead:
                </p>
              </div>
              <CodeBlock terminal code={`grit new myapp --mobile          # Go API + Expo app
grit new myapp --arch=mobile     # identical to --mobile
grit new myapp --arch=triple --expo   # web + admin + API, plus Expo`} className="mb-0" />
            </div>

            {/* Step 2 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  2
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Project Structure</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  The Expo app lives at <code>apps/expo</code> and uses{' '}
                  <a href="https://docs.expo.dev/router/introduction/" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Expo Router</a>{' '}
                  for file-based navigation — the same mental model as the Next.js App Router,
                  but for native screens. Files under <code>app/</code> become screens; the{' '}
                  <code>(tabs)/</code> group is the tab bar.
                </p>
              </div>
              <CodeBlock language="text" filename="myapp/apps/expo/" code={`apps/expo/
├── app.json                  # Expo configuration (name, icons, scheme)
├── package.json
├── babel.config.js
├── metro.config.js
├── tailwind.config.js        # NativeWind (Tailwind for React Native)
├── tsconfig.json
├── app/                      # Expo Router screens (file-based)
│   ├── _layout.tsx           # Root layout — providers (Query, Auth, Theme)
│   ├── (auth)/               # Public auth stack
│   │   ├── login.tsx
│   │   └── register.tsx
│   ├── (tabs)/               # Tab navigation
│   │   ├── _layout.tsx       # Tab bar definition
│   │   ├── index.tsx         # Home tab
│   │   ├── explore.tsx       # More tab — generated resources link here
│   │   ├── profile.tsx
│   │   └── settings.tsx
│   └── <resource>/           # Generated resource screens (per resource)
│       ├── index.tsx         #   list
│       ├── [id].tsx          #   detail
│       ├── new.tsx           #   create
│       └── edit/[id].tsx     #   edit
├── components/
│   ├── ui/                   # ScreenHeader, FormSheet, pickers, …
│   └── resource-forms/       # Generated <resource>-form.tsx components
├── hooks/                    # Typed React Query hooks (use-<resource>.ts)
├── lib/
│   ├── api.ts                # API client (URL resolution + token refresh)
│   ├── auth.tsx              # Auth provider (SecureStore-backed)
│   ├── theme.tsx             # Light/dark palette
│   ├── query-client.ts       # React Query client
│   ├── images.ts             # resolveImageUrl() host rewrite
│   └── upload.ts             # Native file upload helper
└── assets/                   # Icons, splash, fonts`} className="mb-0" />
            </div>

            {/* Step 3 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  3
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Start Development</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  You need two processes running: the Go API and the Expo dev server. From the
                  project root, <code>grit start</code> launches everything in parallel, or you
                  can run each app on its own.
                </p>
              </div>
              <CodeBlock terminal code={`cd myapp

grit start          # Go API + every frontend app (incl. Expo) in parallel
grit start server   # Go API only  →  http://localhost:8080
grit start expo     # Expo dev server only  (runs \`expo start\` in apps/expo)`} className="mb-0 glow-purple-sm" />
              <div className="prose-grit mt-4">
                <p>
                  Under the hood, <code>grit start expo</code> runs <code>expo start</code> in{' '}
                  <code>apps/expo</code>. Expo prints a QR code and starts Metro (the JS bundler)
                  on port <code>8081</code>. Leave the Go API running on <code>8080</code> in a
                  second terminal so the app has a backend to talk to.
                </p>
                <blockquote>
                  First launch installs Expo dependencies and warms the Metro cache, so it takes
                  a moment. Subsequent starts are fast.
                </blockquote>
              </div>
            </div>

            {/* Step 4: API URL matrix */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  4
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">
                  Point the App at the API
                </h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  The single biggest gotcha in mobile development: <code>localhost</code> on a
                  phone means <em>the phone itself</em>, not your dev machine. Grit&apos;s
                  scaffolded <code>lib/api.ts</code> resolves the right host automatically, so
                  you rarely edit anything — but it helps to know the rules.
                </p>
              </div>

              <div className="overflow-x-auto mb-6">
                <table className="w-full text-sm border border-border rounded-lg">
                  <thead>
                    <tr className="border-b border-border bg-muted/50">
                      <th className="text-left p-3 font-medium">Target</th>
                      <th className="text-left p-3 font-medium">Base URL used</th>
                      <th className="text-left p-3 font-medium">Why</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-border text-muted-foreground">
                    {[
                      ['iOS simulator', 'http://localhost:8080/api', 'Shares the host machine’s network'],
                      ['Android emulator', 'http://10.0.2.2:8080/api', '10.0.2.2 is the emulator’s alias for the host'],
                      ['Physical device (Expo Go)', 'http://<machine-LAN-IP>:8080/api', 'Derived from the host Metro is served on'],
                      ['Any target (override)', 'EXPO_PUBLIC_API_URL', 'Wins if set — e.g. a deployed backend'],
                    ].map(([target, url, why]) => (
                      <tr key={target}>
                        <td className="p-3">{target}</td>
                        <td className="p-3 font-mono text-xs text-primary">{url}</td>
                        <td className="p-3">{why}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="prose-grit mb-4">
                <p>
                  Here is the resolver Grit ships in <code>lib/api.ts</code>. On a physical
                  device it reads the LAN host the phone <em>already</em> used to reach Metro —
                  so it just works over Wi-Fi without you hard-coding an IP:
                </p>
              </div>
              <CodeBlock language="typescript" filename="apps/expo/lib/api.ts" code={`const API_PORT = 8080;

function resolveApiUrl(): string {
  // 1. An explicit env var always wins (e.g. a deployed backend).
  const explicit = process.env.EXPO_PUBLIC_API_URL;
  if (explicit) return explicit.replace(/\\/$/, "") + "/api";

  // 2. Derive the dev machine's LAN IP from the host the phone used
  //    to reach Metro. A real device can't reach "localhost"/"10.0.2.2".
  const hostUri =
    Constants.expoConfig?.hostUri ||
    (Constants.expoGoConfig as any)?.debuggerHost;
  const host = hostUri ? String(hostUri).split(":")[0] : undefined;
  if (host && host !== "localhost" && host !== "127.0.0.1") {
    return "http://" + host + ":" + API_PORT + "/api";
  }

  // 3. Fall back to the platform loopback (simulator / web).
  return Platform.select({
    android: "http://10.0.2.2:" + API_PORT + "/api",
    default: "http://localhost:" + API_PORT + "/api",
  }) as string;
}`} className="mb-4" />

              <div className="rounded-lg border border-primary/20 bg-primary/5 p-4">
                <p className="text-sm text-muted-foreground leading-relaxed">
                  <strong className="text-foreground">Deploying?</strong> Set{' '}
                  <code className="text-xs font-mono">EXPO_PUBLIC_API_URL=https://api.yourapp.com</code>{' '}
                  in the environment and the resolver uses it verbatim (appending <code>/api</code>).
                  That is how release builds reach your production backend instead of a LAN IP.
                </p>
              </div>
            </div>

            {/* Step 5: First run */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  5
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Run It on Your Phone</h2>
              </div>
              <div className="prose-grit mb-0">
                <ol>
                  <li>Start the Go API: <code>grit start server</code> (leave it running).</li>
                  <li>Start Expo: <code>grit start expo</code> — a QR code appears in the terminal.</li>
                  <li>Open <strong>Expo Go</strong> on your phone and scan the QR code (iOS: Camera app; Android: the Expo Go scanner).</li>
                  <li>The app bundles and opens. Register an account, or log in with the seeded admin (<code>admin@example.com</code> / <code>admin123</code>).</li>
                </ol>
                <blockquote>
                  Phone and computer must be on the <strong>same Wi-Fi network</strong>. If the
                  app hangs on the splash screen, the API URL is usually the culprit — confirm the
                  Go API is reachable at your machine&apos;s LAN IP on port 8080. Grit&apos;s client
                  fails fast after 15s rather than hanging indefinitely.
                </blockquote>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/architecture-modes/mobile" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Mobile Overview
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/mobile/first-app" className="gap-1.5">
                  Your First Mobile App
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
