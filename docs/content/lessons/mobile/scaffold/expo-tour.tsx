import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Tour the Expo project. By the end you&apos;ll know which folder holds
        what and how file-based routing works in Expo Router.
      </p>

      <h2>The tree</h2>
      <CodeBlock
        language="text"
        filename="apps/mobile/"
        code={`apps/mobile/
├── app/                       file-based routes (= URLs / screens)
│   ├── _layout.tsx            root layout, mounts providers
│   ├── index.tsx              the home screen
│   ├── (auth)/                grouped routes (parens = no URL segment)
│   │   ├── _layout.tsx
│   │   ├── login.tsx
│   │   └── register.tsx
│   └── (tabs)/                tabbed navigation
│       ├── _layout.tsx
│       ├── index.tsx          first tab
│       └── settings.tsx
├── components/                reusable UI (Button, Card, …)
├── hooks/                     useAuth, useUsers, …
├── lib/                       api client, secure-storage, utils
├── assets/                    images, fonts, icons
├── app.json                   Expo config (name, icons, splash, deep linking)
├── eas.json                   build profiles (preview, production)
├── tsconfig.json
└── package.json`}
      />

      <h2>File-based routing</h2>
      <p>
        Same pattern as Next.js App Router. The file becomes the URL.
      </p>
      <CodeBlock
        language="text"
        code={`app/index.tsx          → /
app/about.tsx          → /about
app/(auth)/login.tsx   → /login    ← (auth) doesn't appear in URL
app/users/[id].tsx     → /users/123
app/_layout.tsx        → wraps every child screen`}
      />

      <h2>Grouped routes — the (parens) trick</h2>
      <p>
        Folders in parens are <em>route groups</em>. They organise files
        without affecting the URL. Use them to share a layout across
        related screens:
      </p>
      <CodeBlock
        language="text"
        code={`app/
├── (auth)/_layout.tsx    auth-screen layout (centered card, no tabs)
│   ├── login.tsx
│   └── register.tsx
└── (tabs)/_layout.tsx    main app layout (with bottom tab bar)
    ├── index.tsx
    └── settings.tsx`}
      />
      <p>
        Same URLs as before — <code>/login</code>, <code>/register</code>,
        <code>/</code>, <code>/settings</code>. Different shells.
      </p>

      <h2>The root layout</h2>
      <CodeBlock
        language="tsx"
        filename="app/_layout.tsx"
        code={`import { Stack } from 'expo-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider } from '@/hooks/use-auth'

const queryClient = new QueryClient()

export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Stack screenOptions={{ headerShown: false }} />
      </AuthProvider>
    </QueryClientProvider>
  )
}`}
      />
      <p>
        Mounts React Query + AuthProvider once for the whole app. Every
        screen below it can call <code>useQuery</code> and{' '}
        <code>useAuth()</code>.
      </p>

      <TipBox tone="success">
        <strong>The screen that renders is decided by URL.</strong> Push
        navigation with <code>router.push(&apos;/users/123&apos;)</code>; the file
        at <code>app/users/[id].tsx</code> renders with{' '}
        <code>useLocalSearchParams()</code> returning <code>{`{ id: '123' }`}</code>.
      </TipBox>

      <KnowledgeCheck
        question="You have a checkout flow with 3 screens: cart, shipping, payment. You want all three to share a step-progress header but DIFFERENT URLs (/cart, /shipping, /payment). What's the right shape?"
        choices={[
          {
            label: 'Put them all in app/(checkout)/cart.tsx, /shipping.tsx, /payment.tsx with a shared _layout.tsx',
            correct: true,
            feedback:
              "Right — (checkout) doesn't appear in the URL, the _layout.tsx renders the step header, each file is its own screen with its own URL.",
          },
          {
            label: 'Put them in app/checkout/cart.tsx, /shipping.tsx, /payment.tsx',
            feedback:
              "URLs would be /checkout/cart, /checkout/shipping, /checkout/payment — not what you wanted.",
          },
          {
            label: 'One app/checkout.tsx file with internal state',
            feedback:
              "Possible but you lose deep-linkable URLs. Hard to bookmark or share. The route-group pattern is cleaner.",
          },
          {
            label: 'Three independent screens with no shared layout',
            feedback:
              "Works but you'd duplicate the step-progress header in three places. The shared _layout.tsx is the DRY answer.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Open <code>apps/mobile/app/_layout.tsx</code> in your scaffolded
              project. Read it. Then in <code>notes.md</code>:
            </p>
            <ol>
              <li>List what providers it wraps the app in</li>
              <li>List what screens are accessible (look at the file tree under <code>app/</code>)</li>
              <li>
                Find one routegroup (folder in parens). Write its purpose in
                one sentence.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            Expo Router infers routes from the directory tree at build time.
            What you see in <code>app/</code> is exactly what URLs exist.
          </>
        }
        solution={
          <>
            <p>For a fresh scaffold you should find:</p>
            <ul>
              <li>Providers: QueryClient (React Query) + AuthProvider</li>
              <li>Screens: <code>/</code> (home), <code>/login</code>, <code>/register</code>, <code>/settings</code></li>
              <li>
                Route groups: <code>(auth)</code> (login + register share a
                no-tab-bar layout), <code>(tabs)</code> (main app screens
                with bottom tabs)
              </li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of this chapter — actually run it. Simulator vs Expo Go
        vs dev build, and which to use when.
      </p>
    </>
  )
}
