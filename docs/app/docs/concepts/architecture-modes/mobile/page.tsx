import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Files, Folder, File } from '@/components/files'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/concepts/architecture-modes/mobile')

export default function MobileArchitecturePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="max-w-4xl mx-auto py-12 px-6 lg:px-8">
          {/* Header */}
          <div className="mb-14">
            <p className="text-sm font-mono font-medium text-primary mb-3 tracking-wide uppercase">
              Architecture Modes
            </p>
            <h1 className="text-3xl lg:text-4xl font-bold tracking-tight text-foreground mb-4">
              Mobile Architecture: API + Expo React Native
            </h1>
            <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
              Go API paired with an Expo React Native app in a Turborepo monorepo.
              Shared types between backend and mobile, with encrypted token storage
              and file-based navigation.
            </p>
          </div>

          <div className="prose-grit">
            {/* Overview */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Overview
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The mobile architecture gives you a Go backend and an Expo React Native app
                inside a Turborepo monorepo. Types and schemas are shared through a{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared/</code>{' '}
                directory, just like the double and triple architectures. The mobile app uses
                Expo Router for file-based navigation (similar to Next.js App Router but for
                native screens) and SecureStore for encrypted token storage.
              </p>
              <p className="text-muted-foreground leading-relaxed mb-4">
                The Go API is identical to every other architecture -- same handler-service-model
                pattern, same batteries, same code generation. The difference is on the client side:
                instead of a web browser, your users interact with a native iOS/Android app.
              </p>

              <div className="rounded-lg border border-border/40 bg-accent/20 p-5 mb-6">
                <h4 className="text-sm font-semibold text-foreground mb-3">Scaffold command</h4>
                <CodeBlock language="bash" code={`grit new myapp --mobile`} />
              </div>
            </div>

            {/* Key Characteristics */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Key Characteristics
              </h2>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Property</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Mobile Architecture</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Backend</td>
                      <td className="px-4 py-2.5">Go API (Gin + GORM) -- same as all architectures</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Client</td>
                      <td className="px-4 py-2.5">Expo React Native (iOS + Android)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Navigation</td>
                      <td className="px-4 py-2.5">Expo Router (file-based, like Next.js App Router)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Token storage</td>
                      <td className="px-4 py-2.5">SecureStore (encrypted, not localStorage)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Monorepo</td>
                      <td className="px-4 py-2.5">Turborepo + pnpm workspaces</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">Shared types</td>
                      <td className="px-4 py-2.5">packages/shared/ (Zod schemas + TypeScript types)</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            {/* Full Folder Structure */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Full Folder Structure
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                A Turborepo monorepo with the Go API in{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api/</code>,
                the Expo app in{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/expo/</code>,
                and shared types in{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared/</code>.
              </p>
              <Files title="myapp/">
                <Folder name="apps" defaultOpen>
                  <Folder name="api" comment="Go backend (same as other architectures)" defaultOpen>
                    <File name="go.mod" comment="Module: myapp/apps/api" />
                    <File name="go.sum" />
                    <File name="Dockerfile" />
                    <Folder name="cmd" defaultOpen>
                      <File name="server/main.go" />
                      <File name="migrate/main.go" />
                      <File name="seed/main.go" />
                    </Folder>
                    <Folder name="internal" comment="Full Go backend" defaultOpen>
                      <Folder name="config" />
                      <Folder name="database" />
                      <Folder name="models" />
                      <Folder name="handlers" />
                      <Folder name="services" />
                      <Folder name="middleware" />
                      <Folder name="routes" />
                      <Folder name="mail" />
                      <Folder name="storage" />
                      <Folder name="jobs" />
                      <Folder name="cache" />
                      <Folder name="ai" />
                      <Folder name="auth" />
                    </Folder>
                  </Folder>
                  <Folder name="expo" comment="Expo React Native app" defaultOpen>
                    <File name="package.json" />
                    <File name="app.json" comment="Expo configuration" />
                    <File name="babel.config.js" />
                    <File name="tsconfig.json" />
                    <Folder name="app" comment="Expo Router screens" defaultOpen>
                      <File name="_layout.tsx" comment="Root layout" />
                      <File name="login.tsx" comment="Login screen" />
                      <File name="register.tsx" comment="Register screen" />
                      <Folder name="(tabs)" comment="Tab navigation" defaultOpen>
                        <File name="_layout.tsx" comment="Tab bar layout" />
                        <File name="index.tsx" comment="Home tab" />
                        <File name="profile.tsx" comment="Profile tab" />
                        <File name="settings.tsx" comment="Settings tab" />
                      </Folder>
                    </Folder>
                    <Folder name="components" comment="React Native components" />
                    <Folder name="hooks" comment="useAuth, useQuery hooks" />
                    <Folder name="lib" comment="API client, SecureStore helpers" />
                    <Folder name="assets" comment="Images, fonts" />
                  </Folder>
                </Folder>
                <Folder name="packages" defaultOpen>
                  <Folder name="shared" comment="Shared types + schemas" defaultOpen>
                    <File name="package.json" />
                    <File name="tsconfig.json" />
                    <Folder name="schemas" comment="Zod validation schemas" />
                    <Folder name="types" comment="TypeScript interfaces" />
                    <Folder name="constants" comment="API routes, config" />
                  </Folder>
                </Folder>
                <File name=".env" />
                <File name=".env.example" />
                <File name=".gitignore" />
                <File name="docker-compose.yml" comment="PostgreSQL, Redis, MinIO, Mailhog" />
                <File name="docker-compose.prod.yml" />
                <File name="turbo.json" />
                <File name="pnpm-workspace.yaml" />
                <File name="grit.json" comment={'architecture: "mobile"'} />
                <Folder name=".claude/skills/grit" defaultOpen>
                  <File name="SKILL.md" />
                  <File name="reference.md" />
                </Folder>
              </Files>
            </div>

            {/* Directory Breakdown */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Directory Breakdown
              </h2>

              <div className="space-y-6">
                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">apps/api/</code> -- Go Backend
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    Identical to every other Grit architecture. The full Go backend with Gin router,
                    GORM models, handler-service-model pattern, and all batteries (auth, storage,
                    email, jobs, cache, AI). The API does not know or care whether the client is
                    a web browser or a mobile app -- it just serves JSON over REST.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">apps/expo/</code> -- Expo React Native App
                  </h3>
                  <p className="text-muted-foreground leading-relaxed mb-3">
                    The Expo app uses the managed workflow with Expo Router for file-based navigation.
                    Screens live in the{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/</code> directory
                    and follow the same conventions as Next.js App Router:{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">_layout.tsx</code> for
                    layout wrappers, group folders like{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">(tabs)/</code> for
                    navigation groups, and regular{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.tsx</code> files
                    for individual screens.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">apps/expo/app/(tabs)/</code> -- Tab Navigation
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    The tab group provides the main navigation structure. The{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">_layout.tsx</code> inside
                    this folder defines the tab bar with icons. Each file becomes a tab:{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">index.tsx</code> is Home,{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">profile.tsx</code> is Profile,
                    and{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">settings.tsx</code> is Settings.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">packages/shared/</code> -- Shared Types
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    Zod schemas and TypeScript interfaces shared between the API types (via{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code>) and
                    the Expo app. Import with{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">@shared/schemas</code> or{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">@shared/types</code>{' '}
                    from anywhere in the monorepo.
                  </p>
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    <code className="text-sm font-mono text-primary">grit.json</code> -- Project Config
                  </h3>
                  <CodeBlock language="json" filename="grit.json" code={`{
  "architecture": "mobile",
  "go_module": "myapp/apps/api"
}`} />
                </div>
              </div>
            </div>

            {/* SecureStore vs localStorage */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                SecureStore vs localStorage
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                On the web, JWT tokens are stored in{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localStorage</code>.
                On mobile, this is insecure because any app with root access can read localStorage.
                Expo provides{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">expo-secure-store</code>,
                which uses the iOS Keychain and Android Keystore to encrypt tokens at rest.
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Feature</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Web (localStorage)</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Mobile (SecureStore)</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Encryption</td>
                      <td className="px-4 py-2.5">None (plaintext)</td>
                      <td className="px-4 py-2.5">AES-256 (hardware-backed)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Storage</td>
                      <td className="px-4 py-2.5">Browser storage</td>
                      <td className="px-4 py-2.5">iOS Keychain / Android Keystore</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">Access</td>
                      <td className="px-4 py-2.5">Any JS in the page</td>
                      <td className="px-4 py-2.5">Only your app (sandboxed)</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <CodeBlock language="typescript" filename="apps/expo/lib/secure-store.ts" code={`import * as SecureStore from 'expo-secure-store'

const TOKEN_KEY = 'auth_token'
const REFRESH_KEY = 'refresh_token'

export async function saveTokens(access: string, refresh: string) {
  await SecureStore.setItemAsync(TOKEN_KEY, access)
  await SecureStore.setItemAsync(REFRESH_KEY, refresh)
}

export async function getAccessToken(): Promise<string | null> {
  return SecureStore.getItemAsync(TOKEN_KEY)
}

export async function clearTokens() {
  await SecureStore.deleteItemAsync(TOKEN_KEY)
  await SecureStore.deleteItemAsync(REFRESH_KEY)
}`} />
            </div>

            {/* API URL */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                API URL Configuration
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                On the web,{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localhost:8080</code>{' '}
                works fine because the browser and the API are on the same machine. On a physical
                mobile device,{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localhost</code>{' '}
                refers to the phone itself, not your development machine. You need to use your
                machine&apos;s local IP address instead.
              </p>

              <CodeBlock language="typescript" filename="apps/expo/lib/api-client.ts" code={`import Constants from 'expo-constants'
import { getAccessToken } from './secure-store'

// In development, use your machine's local IP (not localhost)
// Find it with: ipconfig (Windows) or ifconfig (macOS/Linux)
const API_URL = __DEV__
  ? 'http://192.168.1.100:8080'  // Your machine's IP
  : 'https://api.yourapp.com'     // Production URL

export async function apiFetch(path: string, options: RequestInit = {}) {
  const token = await getAccessToken()

  const res = await fetch(API_URL + path, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: 'Bearer ' + token } : {}),
      ...options.headers,
    },
  })

  if (!res.ok) throw new Error(await res.text())
  return res.json()
}`} />
              <p className="text-sm text-muted-foreground/60 mt-3">
                On the iOS simulator and Android emulator, you can use{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localhost</code>{' '}
                (iOS) or{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">10.0.2.2</code>{' '}
                (Android emulator). But on a real device connected via USB or Wi-Fi, use the
                machine&apos;s local IP.
              </p>
            </div>

            {/* Data Flow */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Data Flow
              </h2>
              <CodeBlock language="bash" filename="data flow" code={`Mobile App (Expo)
    \u2502
    \u251c\u2500\u2500 SecureStore → reads JWT token
    \u2502
    \u251c\u2500\u2500 fetch('https://api.yourapp.com/api/posts')
    \u2502       \u2502
    \u2502       \u2514\u2500\u2500 Go API → JWT middleware → handler → service → PostgreSQL
    \u2502
    \u251c\u2500\u2500 FlatList renders data (not HTML tables)
    \u251c\u2500\u2500 Pull-to-refresh triggers refetch
    \u2514\u2500\u2500 Infinite scroll loads next page`} />
            </div>

            {/* Mobile-Specific Patterns */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Mobile-Specific Patterns
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Mobile apps use different UI patterns than web apps. Here are the key differences
                in the Expo scaffold.
              </p>

              <div className="space-y-6">
                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    FlatList Instead of HTML Tables
                  </h3>
                  <p className="text-muted-foreground leading-relaxed mb-3">
                    React Native doesn&apos;t have{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`<table>`}</code> or{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`<div>`}</code>.
                    Data lists use{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">FlatList</code>{' '}
                    which virtualizes rendering for performance -- only visible items are rendered.
                  </p>
                  <CodeBlock language="typescript" filename="example: FlatList" code={`<FlatList
  data={posts}
  keyExtractor={(item) => item.id.toString()}
  renderItem={({ item }) => (
    <Pressable onPress={() => router.push('/post/' + item.id)}>
      <Text style={styles.title}>{item.title}</Text>
      <Text style={styles.date}>{item.created_at}</Text>
    </Pressable>
  )}
  onRefresh={refetch}
  refreshing={isRefetching}
  onEndReached={fetchNextPage}
  onEndReachedThreshold={0.5}
/>`} />
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    Tab Navigation Layout
                  </h3>
                  <p className="text-muted-foreground leading-relaxed mb-3">
                    The tab bar is the primary navigation pattern on mobile. Expo Router makes this
                    declarative with a{' '}
                    <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">(tabs)/</code>{' '}
                    group folder.
                  </p>
                  <CodeBlock language="typescript" filename="apps/expo/app/(tabs)/_layout.tsx" code={`import { Tabs } from 'expo-router'
import { Home, User, Settings } from 'lucide-react-native'

export default function TabLayout() {
  return (
    <Tabs screenOptions={{
      tabBarActiveTintColor: '#6c5ce7',
      tabBarStyle: { backgroundColor: '#0a0a0f', borderTopColor: '#2a2a3a' },
    }}>
      <Tabs.Screen
        name="index"
        options={{ title: 'Home', tabBarIcon: ({ color }) => <Home color={color} size={22} /> }}
      />
      <Tabs.Screen
        name="profile"
        options={{ title: 'Profile', tabBarIcon: ({ color }) => <User color={color} size={22} /> }}
      />
      <Tabs.Screen
        name="settings"
        options={{ title: 'Settings', tabBarIcon: ({ color }) => <Settings color={color} size={22} /> }}
      />
    </Tabs>
  )
}`} />
                </div>

                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    Push Notifications
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    Expo Notifications provides a unified API for push notifications on both iOS
                    and Android. Register the device&apos;s push token with your Go API, then send
                    notifications from the backend using Expo&apos;s push service. This is useful for
                    real-time alerts, chat messages, order updates, and marketing notifications.
                  </p>
                </div>
              </div>
            </div>

            {/* Code Generation */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Code Generation
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Running{' '}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>{' '}
                in a mobile project scaffolds the full stack &mdash; Go backend, shared TypeScript
                types, <strong>and</strong> a working set of Expo Router screens (list, detail, create,
                edit) plus a typed React Query hook and a form component. You get a functioning CRUD
                flow on the device immediately, then style it to taste.
              </p>

              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Generated File</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Location</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Go model</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/models/post.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Go service</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/services/post.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Go handler (+ CSV import)</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/api/internal/handlers/post.go, post_import.go</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Zod schema + TS types</td>
                      <td className="px-4 py-2.5 font-mono text-xs">packages/shared/schemas/post.ts, types/post.ts</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Expo list + detail screens</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/expo/app/posts/index.tsx, [id].tsx</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Expo create + edit screens</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/expo/app/posts/new.tsx, edit/[id].tsx</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Expo form component</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/expo/components/resource-forms/posts-form.tsx</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">React Query hook</td>
                      <td className="px-4 py-2.5 font-mono text-xs">apps/expo/hooks/use-posts.ts</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div className="rounded-lg border border-primary/20 bg-primary/5 p-5 mt-6">
                <h4 className="text-sm font-semibold text-foreground mb-2">Working screens, then make them yours</h4>
                <p className="text-sm text-muted-foreground">
                  The generated Expo screens are real, navigable, and wired to the API through the typed
                  hook &mdash; a full create/read/update/delete flow the moment you generate the resource.
                  They&apos;re a starting point, not a cage: restyle the <code className="text-xs font-mono bg-accent/50 px-1 rounded">FlatList</code>,
                  swap components, or rework the form in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1 rounded">resource-forms/posts-form.tsx</code> &mdash;
                  the shared types keep everything in sync with the backend.
                </p>
              </div>
            </div>

            {/* Building & Distribution */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Building and Distribution
              </h2>
              <p className="text-muted-foreground leading-relaxed mb-4">
                Expo provides two ways to distribute your app: EAS Build for full native builds
                (submitted to App Store / Google Play), and OTA updates for JavaScript-only changes
                that skip the app store review process.
              </p>

              <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                EAS Build
              </h3>
              <CodeBlock language="bash" filename="build for stores" code={`# Install EAS CLI
npm install -g eas-cli

# Configure build profiles
cd apps/expo && eas build:configure

# Build for iOS (requires Apple Developer account)
eas build --platform ios

# Build for Android
eas build --platform android

# Submit to stores
eas submit --platform ios
eas submit --platform android`} />

              <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                OTA Updates
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-3">
                For JavaScript-only changes (bug fixes, UI tweaks, new screens), you can push
                updates directly to users without going through the app store review process.
              </p>
              <CodeBlock language="bash" code={`# Push an update to all users
eas update --branch production --message "Fix login screen layout"`} />

              <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                API Deployment
              </h3>
              <p className="text-muted-foreground leading-relaxed mb-3">
                The Go API is deployed separately using Docker, the same as the API-only architecture.
              </p>
              <CodeBlock language="bash" code={`# Deploy the API
cd apps/api && docker build -t myapp-api .
docker push myapp-api:latest`} />
            </div>

            {/* When to Choose */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                When to Choose Mobile
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="rounded-lg border border-primary/20 bg-primary/5 p-5">
                  <h4 className="text-sm font-semibold text-foreground mb-3">Good fit</h4>
                  <ul className="space-y-2 text-sm text-muted-foreground">
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Mobile-first products (social, delivery, fitness, etc.)
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Cross-platform apps (one codebase for iOS + Android)
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Apps that need push notifications, camera, GPS
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Teams that know React and want to build native apps
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary mt-0.5">-</span>
                      Apps where the primary interface is a phone
                    </li>
                  </ul>
                </div>
                <div className="rounded-lg border border-border/30 bg-card/30 p-5">
                  <h4 className="text-sm font-semibold text-foreground mb-3">Not ideal for</h4>
                  <ul className="space-y-2 text-sm text-muted-foreground">
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Web-first apps (use double or triple instead)
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Apps that need both web and mobile (add web to this scaffold later)
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Games or heavy GPU apps (use Unity or native SDKs)
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-muted-foreground/50 mt-0.5">-</span>
                      Projects that need an admin panel (add admin app separately)
                    </li>
                  </ul>
                </div>
              </div>
            </div>

            {/* Example Project */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Example Project
              </h2>
              <div className="rounded-lg border border-border/40 bg-accent/20 p-5">
                <p className="text-muted-foreground mb-3">
                  The Job Portal example built with the mobile architecture. Includes the Go API,
                  Expo app with tab navigation, SecureStore auth, shared types, and EAS build config.
                </p>
                <a
                  href="https://github.com/MUKE-coder/grit/tree/main/examples/job-portal-mobile-expo"
                  target="_blank"
                  rel="noreferrer"
                  className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
                >
                  View Job Portal (Mobile + Expo) on GitHub &rarr;
                </a>
              </div>
            </div>
          </div>

          {/* Nav */}
          <div className="flex items-center justify-between pt-8 border-t border-border/40">
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes/api-only" className="gap-1.5">
                <ArrowLeft className="h-3.5 w-3.5" />
                API Only Architecture
              </Link>
            </Button>
            <Button variant="ghost" size="sm" asChild className="text-muted-foreground/70 hover:text-foreground">
              <Link href="/docs/concepts/architecture-modes" className="gap-1.5">
                Architecture Modes Overview
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
