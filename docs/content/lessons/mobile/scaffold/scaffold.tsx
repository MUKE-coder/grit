import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Welcome to <strong>Building Mobile with Go API</strong>. By the end
        of this chapter you&apos;ll have an Expo app running on a simulator AND
        on your physical phone, talking to your own Grit API.
      </p>

      <h2>The command</h2>
      <CodeBlock terminal code={`grit new field-app --mobile`} />
      <p>
        Produces a monorepo with two apps: <code>apps/api</code> (your Go
        backend) and <code>apps/mobile</code> (Expo React Native), plus the
        shared <code>packages/shared</code> for Zod schemas + TS types.
      </p>

      <h2>What gets scaffolded</h2>
      <CodeBlock
        language="text"
        filename="field-app/"
        code={`field-app/
├── apps/
│   ├── api/                          The Grit Go API
│   └── mobile/                       The Expo app
│       ├── app/                      file-based routes
│       │   ├── (auth)/login.tsx
│       │   ├── (tabs)/index.tsx
│       │   └── _layout.tsx
│       ├── components/
│       ├── hooks/
│       ├── lib/                      api client, auth helpers
│       ├── app.json                  Expo config
│       ├── eas.json                  EAS Build config
│       └── package.json
├── packages/shared/                  Zod + TS types
└── docker-compose.yml`}
      />

      <h2>Prerequisites your machine needs</h2>
      <ul>
        <li>Node 18+ and pnpm</li>
        <li>
          <strong>Expo Go app</strong> on your physical phone (App Store / Play Store)
        </li>
        <li>
          For iOS simulator on macOS:{' '}
          <code>xcode-select --install</code> + Xcode
        </li>
        <li>
          For Android emulator: Android Studio with an AVD created
        </li>
      </ul>

      <TipBox tone="info">
        <strong>You don&apos;t need Xcode to ship!</strong> EAS Build (lesson
        5.1) builds your iOS app in Expo&apos;s cloud. You can develop and ship
        an iOS app from a Windows or Linux machine if you use a physical
        iPhone for testing.
      </TipBox>

      <h2>Mobile is monorepo-shaped — the API ships with it</h2>
      <p>
        <code>--mobile</code> includes <code>apps/api</code> by default. You
        get both surfaces. Reasons:
      </p>
      <ul>
        <li>Mobile and API share types through <code>packages/shared</code></li>
        <li>You can iterate API + mobile in lock-step on one branch</li>
        <li>One <code>grit start</code> spins up both</li>
      </ul>

      <h2>What didn&apos;t get scaffolded</h2>
      <ul>
        <li>No <code>apps/web</code> — if you want a web companion, use{' '}
        <code>grit new app --triple --mobile</code> (covered in the
        Multi-Platform course)</li>
        <li>No <code>apps/admin</code> — same as above</li>
      </ul>

      <KnowledgeCheck
        question="You want to build a mobile-only product, no web companion at all, but you do want a server-side admin to manage users. What's the right command?"
        choices={[
          {
            label: 'grit new my-app --mobile',
            feedback:
              "Misses the admin requirement. --mobile gives you mobile + API but no admin panel.",
          },
          {
            label: 'grit new my-app --triple --mobile',
            correct: true,
            feedback:
              "Right — --triple adds the Filament-style admin panel; --mobile adds the Expo client. Same monorepo, all three surfaces.",
          },
          {
            label: 'grit new my-app --api',
            feedback:
              "Wrong — that's API-only. You'd have to add the mobile app yourself.",
          },
          {
            label: 'grit new my-app --double + manually add mobile',
            feedback:
              "Possible but more work. The composition flags exist exactly so you don't have to do this manually.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Scaffold the mobile project on your machine and inspect what
              you got:
            </p>
            <ol>
              <li><code>grit new field-app --mobile</code></li>
              <li><code>cd field-app && ls apps</code></li>
              <li>Open <code>apps/mobile/app.json</code> and find the <code>name</code> + <code>slug</code> fields</li>
            </ol>
            <p>
              Paste both the directory listing and the app.json snippet in{' '}
              <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            If <code>pnpm install</code> takes forever, you have the wrong
            Node version. Run <code>node -v</code> — anything &lt;18 won&apos;t
            work with Expo SDK 50+.
          </>
        }
        solution={
          <>
            <p>You should see:</p>
            <CodeBlock language="text" code={`api  mobile`} />
            <p>And in app.json:</p>
            <CodeBlock
              language="json"
              code={`{
  "expo": {
    "name": "field-app",
    "slug": "field-app",
    "version": "1.0.0",
    "orientation": "portrait",
    ...
  }
}`}
            />
            <p>That&apos;s the scaffold done.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Tour the Expo project next — what each folder is for and how Expo
        Router handles navigation.
      </p>
    </>
  )
}
