import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Your First Mobile App — Grit Mobile Course',
  description: 'Scaffold a Go API + Expo React Native project, set up your emulator, run on a device, and understand the mobile project structure.',
}

export default function FirstMobileAppCourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-mobile" className="hover:text-foreground transition-colors">Grit Mobile</Link>
          <span>/</span>
          <span className="text-foreground">Your First Mobile App</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 1 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Your First Mobile App
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will scaffold a mobile project with Grit, set up your development
            environment, run the app on a device or emulator, and understand how the Go API and
            Expo app work together. By the end, you will have a working mobile app connected to
            a live backend.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ What is Mobile Development with Grit? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Mobile Development with Grit?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit scaffolds a <strong className="text-foreground">monorepo</strong> that contains both a
            Go API and an Expo React Native app. You write React components that compile to native iOS
            and Android code. The Go API is exactly the same as in web projects — same routes, same
            services, same database. Only the frontend is different.
          </p>

          <Definition term="React Native">
            A framework by Meta that lets you build native mobile apps using React and JavaScript.
            Instead of rendering to the browser DOM, React Native renders to real native iOS and
            Android views. You write one codebase, and it runs on both platforms.
          </Definition>

          <Definition term="Expo">
            A platform built on top of React Native that handles the complexity of native iOS and
            Android development. Expo provides dev tools, build services, and over-the-air updates
            so you can focus on writing React code instead of configuring Xcode and Android Studio.
          </Definition>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> The <strong className="text-foreground">Go API</strong> runs on your server — same as a web project</li>
            <li className="flex gap-2"><span className="text-primary">•</span> The <strong className="text-foreground">Expo app</strong> runs on a phone (or emulator) and calls the API over HTTP</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Shared types</strong> keep the API and app in sync</li>
            <li className="flex gap-2"><span className="text-primary">•</span> One codebase, two platforms (iOS + Android)</li>
          </ul>
        </section>

        {/* ═══ Prerequisites ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Prerequisites</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You need the same tools as a web project, plus a few mobile-specific ones:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Tool</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What it does</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Check</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Go 1.21+</td>
                  <td className="px-4 py-3">Runs the backend API</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">go version</code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Node.js 18+</td>
                  <td className="px-4 py-3">Runs the Expo dev server</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">node --version</code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">pnpm 8+</td>
                  <td className="px-4 py-3">Installs JavaScript packages</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">pnpm --version</code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Docker</td>
                  <td className="px-4 py-3">Runs PostgreSQL, Redis, MinIO</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">docker --version</code></td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Expo Go app</td>
                  <td className="px-4 py-3">Runs your app on a physical phone</td>
                  <td className="px-4 py-3">App Store / Play Store</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Definition term="Emulator / Simulator">
            Software that mimics a phone on your computer. On macOS, the <strong className="text-foreground">iOS Simulator</strong> comes
            with Xcode. On any OS, the <strong className="text-foreground">Android Emulator</strong> comes with Android Studio.
            If you don{"'"}t want to install either, you can use the <strong className="text-foreground">Expo Go</strong> app
            on your physical phone instead.
          </Definition>

          <Tip>
            The fastest way to get started is to install the <strong>Expo Go</strong> app on your phone.
            No emulator setup required — just scan a QR code to run the app on your device.
          </Tip>

          <Challenge number={1} title="Install Expo Go">
            <p>Install the Expo Go app on your phone from the App Store (iOS) or Play Store (Android). Open it and confirm it launches successfully.</p>
          </Challenge>

          <Challenge number={2} title="Check Your Tools">
            <p>Open your terminal and run all four check commands from the table above. Confirm each tool is installed and returns a version number.</p>
          </Challenge>
        </section>

        {/* ═══ Scaffold a Mobile Project ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold a Mobile Project</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Use the <Code>--mobile</Code> flag to scaffold a project with an Expo app instead of (or alongside) a web frontend:
          </p>

          <CodeBlock language="bash">{`grit new fitness --mobile`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This creates a monorepo with the following structure:
          </p>

          <CodeBlock language="text">{`fitness/
├── apps/
│   ├── api/           ← Go backend (same as web)
│   └── expo/          ← React Native app (Expo)
├── packages/
│   └── shared/        ← Shared types + schemas
├── docker-compose.yml
└── turbo.json`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>apps/api/</Code> directory is identical to a web project — same Go code, same
            routes, same services. The only difference is <Code>apps/expo/</Code> replacing
            <Code>apps/web/</Code>.
          </p>

          <Definition term="Expo Router">
            File-based routing for React Native. Just like Next.js App Router maps files to web pages,
            Expo Router maps files to mobile screens. Files in the <Code>app/</Code> directory become
            navigable screens in your app.
          </Definition>

          <Challenge number={3} title="Scaffold a Mobile Project">
            <p>Run <Code>grit new fitness --mobile</Code>. Explore the folder structure. How many directories are inside <Code>apps/</Code>? What{"'"}s inside <Code>packages/shared/</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Start the API ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start the API</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The API is a standard Grit Go backend. Start the infrastructure services first, then
            run the API server:
          </p>

          <CodeBlock language="bash">{`# Start PostgreSQL, Redis, MinIO
docker compose up -d

# Run the Go API
cd apps/api && go run cmd/server/main.go`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The API will start on <Code>http://localhost:8080</Code>. You can verify it{"'"}s running by
            opening that URL in your browser — you should see a JSON health check response.
          </p>

          <Note>
            The mobile app connects to this exact same API. There is no special {'"'}mobile backend{'"'} —
            the Go API serves both web and mobile clients identically.
          </Note>

          <Challenge number={4} title="Start the API">
            <p>Start the infrastructure with <Code>docker compose up -d</Code>, then run the API. Open <Code>http://localhost:8080</Code> in your browser. What JSON response do you see?</p>
          </Challenge>
        </section>

        {/* ═══ Start Expo ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start Expo</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With the API running, open a new terminal and start the Expo dev server:
          </p>

          <CodeBlock language="bash">{`cd apps/expo && npx expo start`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Expo will display a QR code in the terminal. You have three options:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> <strong className="text-foreground">Physical phone:</strong> Scan the QR code with Expo Go (Android) or your camera app (iOS)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> <strong className="text-foreground">iOS Simulator:</strong> Press <Code>i</Code> in the terminal (macOS only, requires Xcode)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> <strong className="text-foreground">Android Emulator:</strong> Press <Code>a</Code> in the terminal (requires Android Studio)</li>
          </ul>

          <Tip>
            Hot reloading works automatically. When you save a file, the app updates instantly on
            your device or emulator — no need to restart.
          </Tip>

          <Challenge number={5} title="Run Expo">
            <p>Start Expo and open the app on your phone or emulator. You should see the default home screen. Try editing <Code>apps/expo/app/index.tsx</Code> — does the change appear immediately?</p>
          </Challenge>
        </section>

        {/* ═══ Project Structure ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Project Structure</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Expo app follows the same conventions as a Next.js project. Here{"'"}s what{"'"}s inside
            <Code>apps/expo/</Code>:
          </p>

          <CodeBlock language="text">{`apps/expo/
├── app/                 ← Screens (Expo Router)
│   ├── _layout.tsx      ← Root layout (wraps all screens)
│   ├── index.tsx        ← Home screen
│   ├── login.tsx        ← Login screen
│   ├── register.tsx     ← Register screen
│   └── (tabs)/          ← Tab navigation group
│       ├── _layout.tsx  ← Tab bar configuration
│       ├── index.tsx    ← Home tab
│       └── profile.tsx  ← Profile tab
├── components/          ← Reusable UI components
├── hooks/               ← Custom hooks (useAuth, etc.)
├── lib/                 ← API client, utilities
└── package.json`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key Expo Router conventions:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>_layout.tsx</Code> — Wraps child screens. Defines navigation structure (stack, tabs).</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>index.tsx</Code> — The default screen for a directory (like <Code>/</Code> in a URL).</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>(tabs)/</Code> — A route group. The parentheses mean it doesn{"'"}t appear in the URL.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Any <Code>.tsx</Code> file in <Code>app/</Code> becomes a screen automatically.</li>
          </ul>

          <Challenge number={6} title="Explore the Structure">
            <p>Open <Code>apps/expo/app/</Code>. List all the files. Which file is the root layout? What does the <Code>(tabs)/</Code> directory contain?</p>
          </Challenge>

          <Challenge number={7} title="Understand the Layout">
            <p>Open <Code>apps/expo/app/_layout.tsx</Code>. What component does it use to wrap the app? Does it set up a Stack navigator or a Tab navigator at the root level?</p>
          </Challenge>
        </section>

        {/* ═══ Shared Types ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Shared Types</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>packages/shared/</Code> directory is used by <strong className="text-foreground">both</strong> the
            Go API and the Expo app. It contains Zod schemas and TypeScript types that keep the
            frontend and backend in sync.
          </p>

          <CodeBlock language="text">{`packages/shared/
├── types/
│   ├── user.ts       ← User type + schemas
│   ├── auth.ts       ← Login/register schemas
│   └── api.ts        ← API response types
├── schemas/
│   └── index.ts      ← Zod validation schemas
└── constants/
    └── routes.ts     ← API route constants`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you use <Code>grit generate resource</Code> to add a new API resource, the
            TypeScript types are automatically generated in <Code>packages/shared/</Code>. Your
            mobile app gets them for free — no manual syncing required.
          </p>

          <Challenge number={8} title="Explore Shared Types">
            <p>Open <Code>packages/shared/types/</Code>. What types are defined? Compare them to what the API returns at <Code>/api/auth/login</Code>. Are they the same?</p>
          </Challenge>
        </section>

        {/* ═══ The API Client ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The API Client</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The mobile app calls the Go API using a configured HTTP client. Here{"'"}s how a typical
            API call looks:
          </p>

          <CodeBlock language="tsx">{`// lib/api.ts
const API_URL = "http://localhost:8080"

export async function login(email: string, password: string) {
  const response = await fetch(API_URL + "/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  return response.json()
}`}</CodeBlock>

          <Note>
            On a physical device, <Code>localhost</Code> won{"'"}t work — your phone can{"'"}t reach
            your computer{"'"}s localhost. You need to use your computer{"'"}s local IP address instead
            (e.g., <Code>http://192.168.1.42:8080</Code>). Expo shows your local IP in the terminal
            output when you run <Code>npx expo start</Code>.
          </Note>

          <Tip>
            On Android Emulator, you can use <Code>http://10.0.2.2:8080</Code> to reach the host
            machine{"'"}s localhost. This is a special alias built into the Android emulator.
          </Tip>

          <Challenge number={9} title="Find the API Client">
            <p>Find the API client file in the Expo app (check <Code>lib/</Code> or <Code>hooks/</Code>). What base URL does it use? What headers does it set?</p>
          </Challenge>

          <Challenge number={10} title="Test the Connection">
            <p>With the API running, open the mobile app and try to register a new user. Check the API logs in your terminal — do you see the incoming request?</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s what you learned in this course:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Grit scaffolds a monorepo with a Go API and an Expo React Native app</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> The Go API is identical to web projects — same code, same routes</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Expo Router provides file-based navigation (like Next.js for mobile)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Shared types keep the API and mobile app in sync</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> The API client connects the mobile app to the Go backend over HTTP</li>
          </ul>

          <Challenge number={11} title="Register a User">
            <p>Using the mobile app, register a new user with an email and password. Then log in with those credentials. Does the app navigate to the home screen?</p>
          </Challenge>

          <Challenge number={12} title="End-to-End Test">
            <p>Create a blog post via the API docs at <Code>http://localhost:8080/swagger</Code>. Then open the mobile app and navigate to the content screen. Does the new post appear? Pull down to refresh if needed.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-mobile', label: 'Back to Grit Mobile' }}
            next={{ href: '/courses/grit-mobile/auth-navigation', label: 'Next: Mobile Auth & Navigation' }}
          />
        </div>
      </main>
    </div>
  )
}
