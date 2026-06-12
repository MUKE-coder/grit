import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The whole point of the mobile kit is type-safety from Go struct to
        React Native screen. <code>grit sync</code> is the command that
        makes that real. This lesson covers what it does and how mobile
        consumes the output.
      </p>

      <h2>The sync</h2>
      <CodeBlock terminal code={`grit sync`} />
      <p>
        Reads every model in <code>apps/api/internal/models/</code>, writes
        a matching TypeScript type in{' '}
        <code>packages/shared/src/types/</code>. Same Concepts course
        coverage, same command — works identically in the mobile kit.
      </p>

      <h2>What mobile gets</h2>
      <CodeBlock
        language="ts"
        filename="packages/shared/src/types/user.ts (generated)"
        code={`export interface User {
  id: string
  email: string
  name: string
  role: 'user' | 'staff' | 'admin'
  created_at: string
  updated_at: string
}`}
      />
      <p>
        The mobile app imports this directly:
      </p>
      <CodeBlock
        language="ts"
        filename="apps/mobile/hooks/use-users.ts"
        code={`import { useQuery } from '@tanstack/react-query'
import type { User } from '@grit/shared/types/user'
import { api } from '@/lib/api'

export function useUsers() {
  return useQuery<User[]>({
    queryKey: ['users'],
    queryFn: () => api.get('/api/users').then((r) => r.data.data),
  })
}`}
      />
      <p>
        TypeScript autocompletes <code>user.role</code> with the three
        valid string literals. If the Go side adds a new role,
        <code>grit sync</code> picks it up and TS catches every place that
        doesn&apos;t handle it.
      </p>

      <h2>Zod schemas too</h2>
      <p>
        <code>grit sync</code> also writes Zod schemas. Use them to validate
        form input before it leaves the device:
      </p>
      <CodeBlock
        language="ts"
        code={`import { CreateUserSchema } from '@grit/shared/schemas/user'

const result = CreateUserSchema.safeParse(formInput)
if (!result.success) {
  // result.error.flatten() gives per-field errors for the form
}`}
      />
      <p>
        Same Zod schema validates on mobile, web, and as a sanity check on
        the API before GORM writes.
      </p>

      <h2>When to re-run sync</h2>
      <ul>
        <li>After <code>grit generate resource ...</code> on the API</li>
        <li>After manually editing a Go struct in <code>internal/models/</code></li>
        <li>After pulling main if a teammate changed models</li>
      </ul>
      <p>
        Wire it into your dev script:
      </p>
      <CodeBlock
        language="json"
        filename="package.json (root)"
        code={`{
  "scripts": {
    "dev": "grit sync && turbo dev"
  }
}`}
      />
      <p>
        Now every <code>pnpm dev</code> starts fresh with synced types.
      </p>

      <TipBox tone="warning">
        <strong>Add <code>packages/shared/src/types/</code> to .gitignore?</strong>{' '}
        Don&apos;t. Even though it&apos;s generated, committing it means CI
        doesn&apos;t need a Go install to build the mobile app. Lock-step Go
        and TS changes via PR.
      </TipBox>

      <h2>What sync doesn&apos;t handle</h2>
      <ul>
        <li>
          <strong>Computed / derived fields.</strong> If your handler
          decorates the response with extra fields, those aren&apos;t in the Go
          struct — TS won&apos;t know. Add them manually to a type extension.
        </li>
        <li>
          <strong>Custom JSON marshalling.</strong> If you implement
          <code> MarshalJSON</code> on a Go type to reshape it, sync sees
          the struct, not the actual wire format.
        </li>
      </ul>

      <KnowledgeCheck
        question="You add `IsActive bool` to the Go User struct, run `grit migrate`, but the mobile app's `User` type doesn't have `is_active`. What did you skip?"
        choices={[
          {
            label: '`grit sync` — it regenerates packages/shared types from Go structs',
            correct: true,
            feedback:
              "Right — grit migrate updates the DB; grit sync updates the types. They're two different commands. Re-run sync any time the Go model changes.",
          },
          {
            label: 'Reinstall packages — pnpm i',
            feedback:
              "Wrong — no new dependency. The shared package's source files are out of date until sync regenerates them.",
          },
          {
            label: 'Restart Metro',
            feedback:
              "Sometimes helps for caching weirdness, but the underlying issue is that the .ts files in packages/shared don't have the new field. Sync fixes that.",
          },
          {
            label: 'Reboot the simulator',
            feedback:
              "Won't help — the simulator runs the same code. The bug is in the source types.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Trigger the sync flow end-to-end:</p>
            <ol>
              <li>
                Add a field to your Go User model:{' '}
                <code>Bio string &#96;json:&quot;bio&quot;&#96;</code>
              </li>
              <li>Run <code>grit migrate</code></li>
              <li>
                In <code>apps/mobile/app/index.tsx</code>, try to access{' '}
                <code>user.bio</code> — TypeScript yells
              </li>
              <li>Run <code>grit sync</code></li>
              <li>The error is gone — paste before/after in <code>notes.md</code></li>
            </ol>
          </>
        }
        hint={
          <>
            If you don&apos;t have a typed <code>user</code> handy in the screen
            yet, just import the type:{' '}
            <code>{`import type { User } from '@grit/shared/types/user'`}</code>{' '}
            and declare{' '}
            <code>{`const u: User = {} as User; u.bio`}</code>.
          </>
        }
        solution={
          <>
            <p>Before sync:</p>
            <CodeBlock language="text" code={`Property 'bio' does not exist on type 'User'.`} />
            <p>After sync:</p>
            <CodeBlock language="text" code={`(no errors)`} />
            <p>That round-trip — change Go, migrate, sync — is the loop you&apos;ll do daily.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Types are aligned. Last lesson of the chapter — wrap the API
        client in a thin TypeScript layer that mobile screens consume via
        React Query.
      </p>
    </>
  )
}
