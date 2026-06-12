import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Build the login + register screens — forms, validation, error
        states. The Grit mobile scaffold ships these pre-built; this lesson
        explains every piece so you can customize them.
      </p>

      <h2>The screens</h2>
      <CodeBlock
        language="text"
        code={`app/(auth)/_layout.tsx       wrapper — centered card, no tabs
app/(auth)/login.tsx         email + password
app/(auth)/register.tsx      email + name + password`}
      />

      <h2>The form pattern</h2>
      <CodeBlock
        language="tsx"
        filename="app/(auth)/login.tsx"
        code={`import { useState } from 'react'
import { View, TextInput, Text, Pressable } from 'react-native'
import { router } from 'expo-router'
import { useAuth } from '@/hooks/use-auth'

export default function LoginScreen() {
  const { login, isLoading } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)

  async function onSubmit() {
    setError(null)
    try {
      await login(email, password)
      router.replace('/(tabs)')
    } catch (err) {
      if (err.code === 'INVALID_CREDENTIALS') setError('Wrong email or password.')
      else setError(err.message)
    }
  }

  return (
    <View className="flex-1 p-6">
      <Text className="text-2xl font-semibold mb-6">Sign in</Text>

      <TextInput
        autoCapitalize="none"
        keyboardType="email-address"
        textContentType="emailAddress"
        autoComplete="email"
        placeholder="Email"
        value={email}
        onChangeText={setEmail}
        className="border rounded p-3 mb-3"
      />
      <TextInput
        secureTextEntry
        autoComplete="password"
        textContentType="password"
        placeholder="Password"
        value={password}
        onChangeText={setPassword}
        className="border rounded p-3 mb-3"
      />

      {error && <Text className="text-red-500 mb-3">{error}</Text>}

      <Pressable
        onPress={onSubmit}
        disabled={isLoading}
        className="bg-blue-600 rounded p-3 items-center"
      >
        <Text className="text-white">{isLoading ? 'Signing in…' : 'Sign in'}</Text>
      </Pressable>
    </View>
  )
}`}
      />

      <h2>The bits that matter for mobile UX</h2>
      <ul>
        <li>
          <code>autoCapitalize=&quot;none&quot;</code> on the email field —
          otherwise iOS capitalises the first letter
        </li>
        <li>
          <code>keyboardType=&quot;email-address&quot;</code> — shows the
          email keyboard with <code>@</code> readily accessible
        </li>
        <li>
          <code>autoComplete + textContentType</code> — triggers the OS
          password autofill from 1Password / iCloud Keychain / Bitwarden
        </li>
        <li>
          <code>secureTextEntry</code> on the password field — masks input,
          disables autocorrect, disables screenshots in some OS configs
        </li>
      </ul>

      <TipBox tone="info">
        <strong>Hardware keyboard hint:</strong> on a physical phone, the
        keyboard takes half the screen. Wrap the form in a{' '}
        <code>KeyboardAvoidingView</code> + <code>ScrollView</code> so the
        focused input is always visible above the keyboard. The scaffold
        does this automatically.
      </TipBox>

      <h2>Validation with shared Zod schemas</h2>
      <p>
        From the previous chapter, you have <code>LoginSchema</code> in{' '}
        <code>packages/shared/src/schemas/</code>. Use it client-side too:
      </p>
      <CodeBlock
        language="tsx"
        code={`import { LoginSchema } from '@grit/shared/schemas/auth'

async function onSubmit() {
  const result = LoginSchema.safeParse({ email, password })
  if (!result.success) {
    setError(result.error.flatten().fieldErrors.email?.[0] ?? 'Invalid input')
    return
  }
  await login(email, password)
}`}
      />
      <p>
        Now the same validation runs on mobile, web, and on the API. Bad
        input never reaches the server.
      </p>

      <h2>Useful UX touches</h2>
      <ul>
        <li>
          Show a loading spinner on the button instead of the text. Use{' '}
          <code>ActivityIndicator</code> from React Native.
        </li>
        <li>
          Disable the button while submitting to prevent double-taps.
        </li>
        <li>
          Auto-focus the email field on mount with{' '}
          <code>autoFocus={`{true}`}</code>.
        </li>
        <li>
          Show inline field errors below each input, not just a single
          error string at the top.
        </li>
      </ul>

      <KnowledgeCheck
        question="A user enters `alex@example.com` and a wrong password. The API returns 401 INVALID_CREDENTIALS. Your screen shows the system error 'Network request returned 401'. What's the right fix?"
        choices={[
          {
            label: 'Catch ApiError, branch on err.code, show "Wrong email or password"',
            correct: true,
            feedback:
              "Right — that's the whole point of structured error codes from your API. The system error is for developers; the friendly string is for users.",
          },
          {
            label: 'Add an alert in onError that says \'Login failed\'',
            feedback:
              "Better than the raw error, but generic — the user can't tell if they typed wrong or the server is down. Branching on code is precise.",
          },
          {
            label: 'Disable the login form after a 401',
            feedback:
              "Punishes the user for one typo. They'll try the password again — let them.",
          },
          {
            label: 'Log them out',
            feedback:
              "They're not logged in yet — there's nothing to log out from.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              On your scaffolded mobile app, sign up a new user and sign in:
            </p>
            <ol>
              <li>Run the API (<code>cd apps/api && go run ./cmd/server</code>)</li>
              <li>Run the mobile app (<code>cd apps/mobile && pnpm dev</code>)</li>
              <li>Open the app on simulator or Expo Go</li>
              <li>Navigate to <code>/register</code>, create an account</li>
              <li>You should land on the home screen with the user list</li>
            </ol>
            <p>Paste a screenshot in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If the register call fails, check the API logs — it&apos;ll tell
            you why (validation error, duplicate email, etc.).
          </>
        }
        solution={
          <>
            <p>
              The flow:
            </p>
            <ol>
              <li>Form submitted → POST /api/auth/register</li>
              <li>API returns 201 + access_token + refresh_token</li>
              <li>useAuth saves tokens, sets user state</li>
              <li>Route changes to (tabs)/index — home screen renders</li>
            </ol>
            <p>Three screens, one auth round-trip, you&apos;re signed in.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Tokens stored where? <strong>SecureStore + Keychain</strong> — the
        OS-level secret store. Next lesson.
      </p>
    </>
  )
}
