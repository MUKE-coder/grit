import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        One tiny module owns all your HTTP calls. Screens use React Query
        hooks; the hooks use a typed <code>api</code> client. That&apos;s the
        whole pattern. Get it right once, you write screen code without
        thinking about fetch.
      </p>

      <h2>The API client</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/lib/api.ts"
        code={`import Constants from 'expo-constants'

const API_URL = Constants.expoConfig?.extra?.apiUrl ?? 'http://localhost:8080'

async function request<T>(method: string, path: string, body?: unknown, token?: string): Promise<T> {
  const res = await fetch(\`\${API_URL}\${path}\`, {
    method,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: \`Bearer \${token}\` } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
  })
  const json = await res.json()
  if (!res.ok) throw new ApiError(json.error?.code, json.error?.message, res.status)
  return json as T
}

export const api = {
  get:    <T,>(path: string, token?: string) => request<T>('GET', path, undefined, token),
  post:   <T,>(path: string, body: unknown, token?: string) => request<T>('POST', path, body, token),
  put:    <T,>(path: string, body: unknown, token?: string) => request<T>('PUT', path, body, token),
  del:    <T,>(path: string, token?: string) => request<T>('DELETE', path, undefined, token),
}

export class ApiError extends Error {
  constructor(public code: string, message: string, public status: number) {
    super(message)
  }
}`}
      />
      <p>
        Tiny: 25 lines. Returns the envelope&apos;s <code>data</code> field
        wrapped in your generated TS type. Errors become a typed{' '}
        <code>ApiError</code> with the code from Grit&apos;s response envelope.
      </p>

      <h2>The React Query hook</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/hooks/use-users.ts"
        code={`import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { User } from '@grit/shared/types/user'
import { api } from '@/lib/api'
import { useAuth } from './use-auth'

export function useUsers() {
  const { token } = useAuth()
  return useQuery({
    queryKey: ['users'],
    queryFn: () => api.get<{ data: User[] }>(\`/api/users\`, token).then((r) => r.data),
    enabled: !!token,
  })
}

export function useCreateUser() {
  const qc = useQueryClient()
  const { token } = useAuth()
  return useMutation({
    mutationFn: (input: { email: string; name: string }) =>
      api.post<{ data: User }>('/api/users', input, token),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['users'] }),
  })
}`}
      />

      <h2>Using it in a screen</h2>
      <CodeBlock
        language="tsx"
        filename="apps/mobile/app/(tabs)/index.tsx"
        code={`import { FlatList, Text, View } from 'react-native'
import { useUsers } from '@/hooks/use-users'

export default function HomeScreen() {
  const { data: users, isLoading, error } = useUsers()

  if (isLoading) return <Text>Loading...</Text>
  if (error) return <Text>Error: {error.message}</Text>

  return (
    <FlatList
      data={users ?? []}
      keyExtractor={(u) => u.id}
      renderItem={({ item }) => (
        <View>
          <Text>{item.name}</Text>
          <Text>{item.email}</Text>
        </View>
      )}
    />
  )
}`}
      />
      <p>
        That&apos;s the whole pattern. Screen state lives in React Query; the
        API client is type-safe; mutations invalidate the queries that
        depend on them. Add a new endpoint by adding a hook.
      </p>

      <TipBox tone="success">
        <strong>Don&apos;t put fetch in your screens.</strong> Every screen
        importing the api client directly is how codebases drift into
        inconsistency. Hooks are the single layer between screens and the
        network — typed, cached, predictable.
      </TipBox>

      <h2>Error UI — toast + boundary</h2>
      <CodeBlock
        language="tsx"
        code={`import Toast from 'react-native-toast-message'

export function useCreateUser() {
  return useMutation({
    mutationFn: (input) => api.post(...),
    onError: (err) => {
      if (err instanceof ApiError && err.code === 'VALIDATION_ERROR') {
        // form-level handling
        return
      }
      Toast.show({ type: 'error', text1: err.message })
    },
  })
}`}
      />
      <p>
        Toast for unrecoverable errors; <code>onError</code> branches by{' '}
        <code>err.code</code> for ones the form should display inline. Same
        pattern you saw in the Concepts course.
      </p>

      <KnowledgeCheck
        question="You add `useUsers()` to a screen. It works on the simulator but fails on your physical phone with 'network request failed'. The hook itself is fine. What's the cause?"
        choices={[
          {
            label: 'Token is missing',
            feedback:
              "Wrong — that'd be a 401, not a network failure. The request isn't reaching the server.",
          },
          {
            label: 'API_URL is localhost — works in the simulator (which shares the host network) but not on the phone',
            correct: true,
            feedback:
              "Right — same gotcha as the first-run lesson. Use your computer's LAN IP via Constants.expoConfig.extra in app.json so simulator AND phone both reach it.",
          },
          {
            label: 'React Query needs a network provider',
            feedback:
              "Wrong — React Query handles HTTP via your provided queryFn. No network provider needed.",
          },
          {
            label: 'The phone is on Wi-Fi but the computer is on Ethernet',
            feedback:
              "Could be — they need to be on the same LAN segment. But the underlying fix is the same: use a routable address, not localhost.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Build the chapter assignment:</p>
            <ol>
              <li>
                In <code>apps/mobile/hooks/</code>, add{' '}
                <code>use-users.ts</code> using the pattern above.
              </li>
              <li>
                In <code>apps/mobile/app/(tabs)/index.tsx</code>, render the
                user list.
              </li>
              <li>
                Confirm TypeScript autocompletes <code>user.email</code>{' '}
                without an explicit type annotation — that&apos;s the proof the
                sync is working.
              </li>
              <li>Paste a screenshot of the rendered list in <code>notes.md</code>.</li>
            </ol>
          </>
        }
        hint={
          <>
            If you see &quot;401 unauthorized&quot;, you don&apos;t have a token yet
            — chapter 3 wires login. For this exercise, manually grab a
            token from your API (curl POST /api/auth/login) and hard-code it
            in useAuth temporarily.
          </>
        }
        solution={
          <>
            <p>A working screen looks like:</p>
            <CodeBlock
              language="text"
              code={`Users
─────────────
alex@example.com   Alex
bob@example.com    Bob
carol@example.com  Carol`}
            />
            <p>
              You&apos;ve now got a typed, cached, type-safe API client. Every
              new endpoint is one new hook.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 3 — <strong>Mobile Auth</strong>. Login screens, secure
        token storage with SecureStore / Keychain, and silent refresh-on-401.
      </p>
    </>
  )
}
