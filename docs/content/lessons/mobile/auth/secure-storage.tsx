import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Where you store the user&apos;s JWT matters. AsyncStorage is plaintext
        on disk — anyone with file-system access reads it. SecureStore (iOS
        Keychain, Android Keystore) is OS-encrypted. Use it. Here&apos;s how.
      </p>

      <h2>The right primitive — expo-secure-store</h2>
      <CodeBlock
        terminal
        code={`pnpm add expo-secure-store`}
      />
      <p>
        Already in the Grit mobile scaffold. Stores values in:
      </p>
      <ul>
        <li>
          <strong>iOS:</strong> Keychain — encrypted by the device&apos;s secure
          enclave, tied to the app&apos;s bundle ID
        </li>
        <li>
          <strong>Android:</strong> Keystore + EncryptedSharedPreferences —
          encrypted by hardware-backed AES, tied to app signature
        </li>
      </ul>
      <p>
        Both survive app updates. Both wipe on uninstall. Both refuse access
        from other apps.
      </p>

      <h2>The auth helper</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/lib/auth-storage.ts"
        code={`import * as SecureStore from 'expo-secure-store'

const ACCESS_KEY = 'grit.access_token'
const REFRESH_KEY = 'grit.refresh_token'

export async function saveTokens(access: string, refresh: string) {
  await SecureStore.setItemAsync(ACCESS_KEY, access)
  await SecureStore.setItemAsync(REFRESH_KEY, refresh)
}

export async function loadTokens() {
  return {
    access: await SecureStore.getItemAsync(ACCESS_KEY),
    refresh: await SecureStore.getItemAsync(REFRESH_KEY),
  }
}

export async function clearTokens() {
  await SecureStore.deleteItemAsync(ACCESS_KEY)
  await SecureStore.deleteItemAsync(REFRESH_KEY)
}`}
      />

      <h2>Wiring it into useAuth</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/hooks/use-auth.ts (excerpt)"
        code={`export function AuthProvider({ children }) {
  const [tokens, setTokens] = useState<{ access?: string; refresh?: string }>({})
  const [isLoading, setLoading] = useState(true)

  // Restore on app cold start
  useEffect(() => {
    loadTokens().then((stored) => {
      setTokens(stored)
      setLoading(false)
    })
  }, [])

  async function login(email: string, password: string) {
    const res = await api.post('/api/auth/login', { email, password })
    await saveTokens(res.data.access_token, res.data.refresh_token)
    setTokens({ access: res.data.access_token, refresh: res.data.refresh_token })
  }

  async function logout() {
    await clearTokens()
    setTokens({})
  }

  return (
    <AuthContext.Provider value={{ tokens, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  )
}`}
      />
      <p>
        On app start, we read tokens from SecureStore. If they exist, the
        user is already signed in. No login screen unless they explicitly
        logged out.
      </p>

      <TipBox tone="warning">
        <strong>SecureStore is async.</strong> Don&apos;t render the app until{' '}
        <code>isLoading</code> is false — otherwise the route guard may
        send the user to <code>/login</code> while the token is still being
        loaded. Show a splash screen until ready.
      </TipBox>

      <h2>What NOT to do</h2>
      <CodeBlock
        language="ts"
        code={`// AsyncStorage — plaintext on disk
await AsyncStorage.setItem('token', accessToken)  // bad for JWTs

// localStorage — there isn't one in React Native
// Don't try to use the web localStorage shim`}
      />
      <p>
        AsyncStorage is fine for non-secret preferences (theme, last-used
        filter). For tokens, passwords, API keys — SecureStore only.
      </p>

      <h2>Encryption-at-rest only matters if the device is unlocked-then-compromised</h2>
      <p>
        Once a phone is unlocked, processes that run as you can read your
        Keychain. The defence model is &quot;another app, or someone with file-
        system access (e.g., a backup extraction), can&apos;t read your
        tokens.&quot; SecureStore covers that. It doesn&apos;t protect against
        malware that runs in your app&apos;s sandbox.
      </p>

      <KnowledgeCheck
        question="A teammate uses AsyncStorage to save the JWT. Backups of the phone leak the storage file. What's exposed?"
        choices={[
          {
            label: 'Nothing — AsyncStorage encrypts at rest',
            feedback:
              "Wrong — AsyncStorage is plaintext JSON files in the app's sandbox. Backups (iCloud, Android backups, jailbroken extracts) read them.",
          },
          {
            label: 'The JWT in plaintext — the attacker can impersonate the user',
            correct: true,
            feedback:
              "Right — AsyncStorage stores plaintext. Anyone who reads the file gets the token. SecureStore encrypts via the OS keychain.",
          },
          {
            label: 'A hash of the JWT — they need to crack it',
            feedback:
              "Wrong — AsyncStorage doesn't hash. It writes the raw string you gave it.",
          },
          {
            label: 'AsyncStorage refuses to back up sensitive values',
            feedback:
              "AsyncStorage doesn't know what's sensitive. It backs up everything by default.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Confirm SecureStore is doing its job:</p>
            <ol>
              <li>Sign in on your mobile app</li>
              <li>
                Add a debug button that calls{' '}
                <code>SecureStore.getItemAsync(&apos;grit.refresh_token&apos;)</code>{' '}
                and shows the result in an alert
              </li>
              <li>Confirm the refresh token is there</li>
              <li>
                Force-kill the app, reopen it. You should still be signed in
                (the AuthProvider loaded the tokens on cold start).
              </li>
              <li>Paste a screenshot of the debug alert in <code>notes.md</code>.</li>
            </ol>
          </>
        }
        hint={
          <>
            On iOS, force-kill from the app switcher (swipe up). On Android,
            swipe the app away from recent apps. Both kill the process; cold
            start = a real test.
          </>
        }
        solution={
          <>
            <p>
              You should see a long base64-ish string starting with{' '}
              <code>eyJ</code> (the JWT prefix). After force-kill, reopen
              the app — no login screen, you land on the home screen
              because tokens were restored.
            </p>
            <p>
              That&apos;s the chapter 3 assignment&apos;s middle bit done. Refresh
              is next.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of the chapter — silent <strong>refresh on 401</strong>.
        Access tokens expire; refresh tokens save the user from re-login.
      </p>
    </>
  )
}
