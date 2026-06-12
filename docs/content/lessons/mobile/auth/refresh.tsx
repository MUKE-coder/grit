import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Access tokens expire after 15 minutes. Without refresh, the user
        gets logged out every 15 minutes — terrible UX. With refresh,
        they&apos;re logged in for a week. This lesson covers the pattern.
      </p>

      <h2>The flow</h2>
      <CodeBlock
        language="text"
        code={`Mobile sends request with access_token
  → API returns 401 UNAUTHORIZED (token expired)
Mobile catches 401, calls /api/auth/refresh with refresh_token
  → API returns new access_token (+ rotated refresh_token)
Mobile saves new tokens, retries the original request
  → API returns 200 — request succeeds, user sees no interruption`}
      />
      <p>
        The whole dance happens silently. The user doesn&apos;t see a thing.
      </p>

      <h2>The fetch interceptor pattern</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/lib/api.ts (refresh-aware version)"
        code={`async function authedRequest<T>(method: string, path: string, body?: unknown): Promise<T> {
  let { access, refresh } = await loadTokens()

  const doFetch = (token: string) => fetch(\`\${API_URL}\${path}\`, {
    method,
    headers: { 'Content-Type': 'application/json', Authorization: \`Bearer \${token}\` },
    body: body ? JSON.stringify(body) : undefined,
  })

  let res = await doFetch(access!)

  // If access token is dead, try to refresh once
  if (res.status === 401 && refresh) {
    const r = await fetch(\`\${API_URL}/api/auth/refresh\`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refresh }),
    })
    if (r.ok) {
      const json = await r.json()
      await saveTokens(json.data.access_token, json.data.refresh_token)
      res = await doFetch(json.data.access_token)
    }
  }

  const json = await res.json()
  if (!res.ok) throw new ApiError(json.error?.code, json.error?.message, res.status)
  return json as T
}`}
      />
      <p>
        One interceptor, handles every authed call. If refresh succeeds,
        retry; if it fails, the next 401 propagates and the AuthProvider
        kicks the user to <code>/login</code>.
      </p>

      <h2>Single-flight: don&apos;t refresh 10x in parallel</h2>
      <p>
        If 10 React Query queries all fire simultaneously and all return
        401, you don&apos;t want 10 refresh calls — that creates 10 new
        access tokens and 9 are wasted (and refresh-token rotation makes 9
        of them stale). Use a promise lock:
      </p>
      <CodeBlock
        language="ts"
        code={`let refreshing: Promise<{ access: string; refresh: string }> | null = null

async function refreshTokens() {
  if (refreshing) return refreshing
  refreshing = doRefresh().finally(() => { refreshing = null })
  return refreshing
}`}
      />
      <p>
        Every concurrent call hits the same promise — one network refresh,
        all callers get the new token.
      </p>

      <TipBox tone="warning">
        <strong>Refresh tokens rotate.</strong> Every successful refresh
        returns a NEW refresh token; the old one is invalidated server-side.
        Save the new one immediately. If you forget, the next refresh fails
        and the user is logged out.
      </TipBox>

      <h2>What happens when refresh fails</h2>
      <p>
        Two cases:
      </p>
      <ul>
        <li>
          <strong>Network error</strong> — show a toast, let the user retry.
          Don&apos;t log them out.
        </li>
        <li>
          <strong>401 from /refresh</strong> — refresh token is expired or
          revoked. Clear tokens, send to <code>/login</code>.
        </li>
      </ul>

      <h2>React Query plays nicely</h2>
      <p>
        React Query&apos;s built-in retry kicks in if a query throws. With the
        interceptor doing the refresh, you usually want{' '}
        <code>retry: false</code> on auth-required queries — the interceptor
        already retried; a second retry is wasted.
      </p>

      <KnowledgeCheck
        question="Your user has been on the app for 30 minutes. They open a new screen that fires 5 API calls in parallel. All 5 return 401. What's the expected behaviour?"
        choices={[
          {
            label: 'Refresh 5 times in parallel — user gets 5 new access tokens, 4 are stale immediately',
            feedback:
              "Wrong (or at least suboptimal). Refresh token rotation invalidates each previous refresh, so 4 of 5 would be useless. The single-flight pattern handles this.",
          },
          {
            label: 'One refresh; the other 4 queries wait for it, then retry with the new token',
            correct: true,
            feedback:
              "Right — the single-flight promise ensures one network refresh. All 5 queries share the same refresh and retry with the new access token. This is the standard pattern.",
          },
          {
            label: 'All 5 queries fail and the user sees errors',
            feedback:
              "Possible if you didn't implement the interceptor. The 'silent refresh' UX is exactly to avoid this.",
          },
          {
            label: 'The user is force-logged-out',
            feedback:
              "Only if refresh ITSELF fails. A 401 on a normal endpoint should trigger a refresh attempt, not a logout.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Verify your refresh works end-to-end:</p>
            <ol>
              <li>
                Set your API&apos;s <code>JWT_ACCESS_EXPIRY</code> to{' '}
                <code>30s</code> for testing (in <code>.env</code>).
              </li>
              <li>Restart the API.</li>
              <li>Sign in on the mobile app.</li>
              <li>
                Wait 35 seconds. Then tap something that fires a query.
              </li>
              <li>
                The query should silently refresh + retry; you see data, not
                a login screen.
              </li>
              <li>
                Open Metro&apos;s network inspector — you should see one POST
                to <code>/api/auth/refresh</code> sandwiched between the
                two attempts.
              </li>
            </ol>
            <p>
              Paste the network inspector output in <code>notes.md</code>.
              Don&apos;t forget to revert <code>JWT_ACCESS_EXPIRY</code> to{' '}
              <code>15m</code> after.
            </p>
          </>
        }
        hint={
          <>
            If you get logged out instead of silently refreshed, the
            interceptor isn&apos;t wrapping the failing call. Add a console.log
            in the 401 branch to confirm it&apos;s entered.
          </>
        }
        solution={
          <>
            <p>The network log should show:</p>
            <CodeBlock
              language="text"
              code={`POST /api/users     401 Unauthorized
POST /api/auth/refresh  200 OK
POST /api/users     200 OK`}
            />
            <p>
              Three calls, one user-visible result. The user has no idea
              their token expired. That&apos;s chapter 3 done.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 4 — <strong>Push notifications</strong>. Register for push,
        save the token to your Grit API, send a push from a job worker.
      </p>
    </>
  )
}
