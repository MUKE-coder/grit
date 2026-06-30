import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Mobile Auth & Navigation — Grit Mobile Course',
  description: 'Implement secure token storage, Expo Router navigation, tab and stack navigation, protected routes, and token refresh in a React Native app.',
}

export default function AuthNavigationCourse() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-mobile" className="hover:text-foreground transition-colors">Grit Mobile</Link>
          <span>/</span>
          <span className="text-foreground">Mobile Auth & Navigation</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 2 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Mobile Auth & Navigation
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Mobile authentication is fundamentally different from the web. There are no cookies and
            no localStorage. In this course, you will learn how to store tokens securely, build login
            and logout flows, protect routes, and create tab and stack navigation with Expo Router.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Auth in Mobile Apps ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Auth in Mobile Apps</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On the web, you store JWT tokens in cookies or localStorage. On mobile, neither exists.
            Instead, mobile apps use <strong className="text-foreground">encrypted device storage</strong> that
            is protected by the operating system itself.
          </p>

          <Definition term="SecureStore (expo-secure-store)">
            Encrypted key-value storage provided by Expo. On iOS, it stores data in the
            <strong className="text-foreground"> iOS Keychain</strong>. On Android, it uses the
            <strong className="text-foreground"> Android Keystore</strong>. Tokens stored here are
            encrypted at rest and only accessible to your app — even if the device is stolen,
            the data cannot be read without the device passcode.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This means your authentication tokens are <strong className="text-foreground">more secure</strong> on
            mobile than on web. The Keychain/Keystore is hardware-backed on modern devices, making
            it extremely difficult to extract tokens.
          </p>

          <Challenge number={1} title="Understand the Difference">
            <p>Explain in your own words why mobile apps cannot use cookies or localStorage. What does <Code>expo-secure-store</Code> use instead on iOS vs Android?</p>
          </Challenge>
        </section>

        {/* ═══ Login Flow ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Login Flow</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The login flow on mobile follows four steps:
          </p>

          <ol className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">1.</span> User enters email and password on the login screen</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">2.</span> App calls <Code>POST /api/auth/login</Code> with the credentials</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">3.</span> Store the returned tokens in SecureStore</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">4.</span> Navigate to the home screen</li>
          </ol>

          <CodeBlock language="tsx">{`import * as SecureStore from 'expo-secure-store'

async function handleLogin(email: string, password: string) {
  const response = await fetch(API_URL + "/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  const data = await response.json()

  // Store tokens in encrypted storage
  await SecureStore.setItemAsync('access_token', data.access_token)
  await SecureStore.setItemAsync('refresh_token', data.refresh_token)

  // Navigate to home
  router.replace('/(tabs)')
}`}</CodeBlock>

          <Note>
            <Code>router.replace()</Code> is used instead of <Code>router.push()</Code> so the user
            cannot press the back button to return to the login screen after logging in.
          </Note>

          <Challenge number={2} title="Find Token Storage">
            <p>Find where tokens are stored in the scaffolded code. What key names are used for the access token and refresh token? What file contains this logic?</p>
          </Challenge>

          <Challenge number={3} title="Test the Login">
            <p>Register a user (if you haven{"'"}t already), then log in. Check the API server logs — do you see the login request? Does the app navigate away from the login screen?</p>
          </Challenge>
        </section>

        {/* ═══ Auth Context ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Auth Context</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The auth state needs to be accessible from any screen in the app. Grit uses the
            <strong className="text-foreground"> React Context</strong> pattern — an <Code>AuthProvider</Code> wraps
            the entire app and provides auth state to all screens through a <Code>useAuth</Code> hook.
          </p>

          <CodeBlock language="tsx">{`// hooks/useAuth.tsx
interface AuthContextType {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  logout: () => Promise<void>
}

// In any screen:
const { user, isAuthenticated, logout } = useAuth()`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>AuthProvider</Code> reads tokens from SecureStore on app launch. If valid tokens
            exist, the user is automatically logged in. If not, they see the login screen.
          </p>

          <Challenge number={4} title="Find the Auth Context">
            <p>Find the auth context or <Code>useAuth</Code> hook in the scaffolded code. What values does it provide? Where is the <Code>AuthProvider</Code> mounted?</p>
          </Challenge>
        </section>

        {/* ═══ Protected Routes ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Protected Routes</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Screens inside <Code>(tabs)/</Code> should only be accessible to logged-in users.
            The tab layout checks the auth state and redirects to login if the user is not authenticated:
          </p>

          <CodeBlock language="tsx">{`// app/(tabs)/_layout.tsx
import { Redirect } from 'expo-router'
import { useAuth } from '@/hooks/useAuth'

export default function TabLayout() {
  const { user, isLoading } = useAuth()

  // Show loading screen while checking auth
  if (isLoading) return <LoadingScreen />

  // Redirect to login if not authenticated
  if (!user) return <Redirect href="/login" />

  return (
    <Tabs>
      <Tabs.Screen name="index" options={{ title: "Home" }} />
      <Tabs.Screen name="profile" options={{ title: "Profile" }} />
      <Tabs.Screen name="settings" options={{ title: "Settings" }} />
    </Tabs>
  )
}`}</CodeBlock>

          <Definition term="Protected Route">
            A screen that requires authentication to access. If a user who is not logged in tries
            to open a protected route, they are automatically redirected to the login screen. After
            logging in, they return to the screen they originally tried to access.
          </Definition>

          <Challenge number={5} title="Test Protected Routes">
            <p>Log out of the app and try to access a protected screen (like the home tab). Does it redirect you to the login screen? Log in again — do you return to the protected screen?</p>
          </Challenge>
        </section>

        {/* ═══ Expo Router Navigation ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Expo Router Navigation</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Mobile apps use two primary navigation patterns:
          </p>

          <Definition term="Tab Navigation">
            A bar at the bottom of the screen with icons for each main section (Home, Search,
            Profile, etc.). The user taps a tab to switch between sections. This is the primary
            navigation pattern in most mobile apps — think Instagram, Twitter, or Spotify.
          </Definition>

          <Definition term="Stack Navigation">
            Screens that stack on top of each other. When you tap an item in a list, a detail
            screen slides in from the right. Press the back button to go back. This is used for
            drill-down flows within a tab.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how the file structure maps to navigation:
          </p>

          <CodeBlock language="text">{`app/
├── _layout.tsx          ← Root layout (Stack)
├── login.tsx            ← Login screen
├── register.tsx         ← Register screen
└── (tabs)/
    ├── _layout.tsx      ← Tab bar layout
    ├── index.tsx         ← Home tab
    ├── profile.tsx       ← Profile tab
    └── settings.tsx      ← Settings tab`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The root layout uses a <strong className="text-foreground">Stack</strong> navigator. Login and
            register are stack screens. The <Code>(tabs)/</Code> group contains a
            <strong className="text-foreground"> Tab</strong> navigator with three tabs. This means you
            can have stack navigation (push/pop screens) at the root level, and tab navigation inside
            the authenticated area.
          </p>

          <Challenge number={6} title="Add a New Tab">
            <p>Create a new file at <Code>app/(tabs)/explore.tsx</Code> with a simple screen component. Then add it to the <Code>(tabs)/_layout.tsx</Code> tab configuration. Does it appear in the tab bar?</p>
          </Challenge>

          <Challenge number={7} title="Stack Inside a Tab">
            <p>Inside one of your tabs, add a button that navigates to a detail screen using <Code>router.push()</Code>. Does the back button appear? Can you navigate back?</p>
          </Challenge>
        </section>

        {/* ═══ Token Refresh ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Token Refresh</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The access token expires in 15 minutes. When it expires, the API returns a
            <Code>401 Unauthorized</Code> response. Instead of logging the user out, the API client
            automatically uses the refresh token to get a new access token.
          </p>

          <CodeBlock language="tsx">{`// lib/api.ts — simplified interceptor pattern
async function fetchWithAuth(url: string, options: RequestInit = {}) {
  const token = await SecureStore.getItemAsync('access_token')

  const response = await fetch(url, {
    ...options,
    headers: { ...options.headers, Authorization: "Bearer " + token },
  })

  // If token expired, refresh and retry
  if (response.status === 401) {
    const newToken = await refreshAccessToken()
    if (newToken) {
      return fetch(url, {
        ...options,
        headers: { ...options.headers, Authorization: "Bearer " + newToken },
      })
    }
  }

  return response
}

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = await SecureStore.getItemAsync('refresh_token')
  const response = await fetch(API_URL + "/api/auth/refresh", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
  })
  const data = await response.json()
  if (data.access_token) {
    await SecureStore.setItemAsync('access_token', data.access_token)
    return data.access_token
  }
  return null
}`}</CodeBlock>

          <Definition term="Token Refresh">
            The process of exchanging an expired access token for a new one using a longer-lived
            refresh token. This happens automatically in the background — the user never sees a
            login screen unless their refresh token has also expired (typically after 7 days).
          </Definition>

          <Challenge number={8} title="Find the Refresh Logic">
            <p>Look at the API client in the scaffolded code. Can you find the refresh token logic? What happens if the refresh token itself has expired?</p>
          </Challenge>
        </section>

        {/* ═══ Logout ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Logout</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Logging out involves three steps: clear tokens from SecureStore, reset the auth state,
            and navigate to the login screen.
          </p>

          <CodeBlock language="tsx">{`async function logout() {
  // 1. Clear tokens from encrypted storage
  await SecureStore.deleteItemAsync('access_token')
  await SecureStore.deleteItemAsync('refresh_token')

  // 2. Reset auth state
  setUser(null)

  // 3. Navigate to login
  router.replace('/login')
}`}</CodeBlock>

          <Tip>
            Always use <Code>router.replace()</Code> when navigating to login after logout. This
            prevents the user from pressing the back button to return to a protected screen.
          </Tip>

          <Challenge number={9} title="Test Logout">
            <p>Log in, then log out. After logging out, try pressing the back button. Can you access the protected screens? You should not be able to.</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s what you learned in this course:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Mobile apps use SecureStore (encrypted device storage) instead of cookies</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> The auth context provides user state to all screens via <Code>useAuth</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Protected routes redirect unauthenticated users to login</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Tab navigation is the primary pattern; stack navigation is for drill-downs</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Token refresh happens automatically — users stay logged in seamlessly</li>
          </ul>

          <Challenge number={10} title="Build a Complete Auth Flow">
            <p>Verify your app has: login screen, register screen, protected tab navigation (Home, Profile, Settings), and a logout button in the settings or profile screen.</p>
          </Challenge>

          <Challenge number={11} title="Add a New Protected Screen">
            <p>Add a {'"'}Workouts{'"'} tab to the tab bar. It should show a simple list screen. Verify that it{"'"}s protected — logging out and accessing it directly should redirect to login.</p>
          </Challenge>

          <Challenge number={12} title="Test Token Expiry">
            <p>If possible, set the access token expiry to a short duration (e.g., 1 minute) in the API config. Use the app until the token expires. Does the refresh happen transparently? Does the app stay logged in?</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-mobile/first-app', label: 'Prev: Your First Mobile App' }}
            next={{ href: '/courses/grit-mobile/api-offline', label: 'Next: API Integration & Offline' }}
          />
        </div>
      </main>
    </div>
  )
}
