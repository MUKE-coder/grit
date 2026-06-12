import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Signup + login are the two most important forms in your app. This
        lesson covers the Grit pattern: server action for submit, Zod for
        validation, secure HTTP-only cookie for the session.
      </p>

      <h2>The signup form</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(auth)/signup/page.tsx"
        code={`'use client'
import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { SignupSchema } from '@workspace/shared/schemas/auth'
import { signupAction } from './actions'

export default function SignupPage() {
  const router = useRouter()
  const [errors, setErrors] = useState<Record<string, string>>({})

  async function onSubmit(form: FormData) {
    const result = SignupSchema.safeParse(Object.fromEntries(form))
    if (!result.success) {
      setErrors(result.error.flatten().fieldErrors as any)
      return
    }
    const res = await signupAction(result.data)
    if (res?.error) {
      setErrors({ form: res.error })
      return
    }
    router.push('/dashboard')
  }

  return (
    <form action={onSubmit} className="space-y-3 w-80">
      <input name="email" placeholder="Email" />
      {errors.email && <p className="text-red-500 text-sm">{errors.email}</p>}
      <input name="password" type="password" placeholder="Password" />
      {errors.password && <p className="text-red-500 text-sm">{errors.password}</p>}
      <input name="name" placeholder="Your name" />
      <button className="rounded-full bg-primary px-4 py-2 text-primary-foreground w-full">
        Create account
      </button>
      {errors.form && <p className="text-red-500 text-sm">{errors.form}</p>}
    </form>
  )
}`}
      />

      <h2>The server action that hits your Go API</h2>
      <CodeBlock
        language="ts"
        filename="apps/web/app/(auth)/signup/actions.ts"
        code={`'use server'
import { cookies } from 'next/headers'
import { z } from 'zod'
import { SignupSchema } from '@workspace/shared/schemas/auth'

export async function signupAction(input: z.infer<typeof SignupSchema>) {
  const res = await fetch(\`\${process.env.API_URL}/api/auth/register\`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  })
  const json = await res.json()
  if (!res.ok) return { error: json.error?.message ?? 'Signup failed' }

  // Set HTTP-only cookies with the tokens
  const cookieStore = await cookies()
  cookieStore.set('customer-access', json.data.access_token, {
    httpOnly: true, secure: true, sameSite: 'lax', path: '/',
    maxAge: 60 * 15,                      // 15 min
  })
  cookieStore.set('customer-refresh', json.data.refresh_token, {
    httpOnly: true, secure: true, sameSite: 'lax', path: '/',
    maxAge: 60 * 60 * 24 * 7,             // 7 days
  })
}`}
      />
      <p>
        Server action runs in the Next.js server process, talks to your Go
        API, then sets HTTP-only cookies. JS in the browser never sees
        the token — XSS can&apos;t steal it.
      </p>

      <TipBox tone="warning">
        <strong>HTTP-only cookies, not localStorage.</strong> JWTs in
        localStorage are readable by any JS on the page. An XSS bug = your
        whole user base&apos;s sessions stolen. HTTP-only cookies are
        unreachable from JS — XSS can&apos;t exfiltrate them.
      </TipBox>

      <h2>Refresh on 401 — server-side</h2>
      <p>
        For server components / route handlers that fetch your Go API, wrap
        the fetch in an interceptor that refreshes when the access cookie
        is expired:
      </p>
      <CodeBlock
        language="ts"
        filename="apps/web/lib/api.ts (server-side)"
        code={`export async function apiFetch(path: string, init?: RequestInit) {
  const access = (await cookies()).get('customer-access')?.value

  let res = await fetch(\`\${process.env.API_URL}\${path}\`, {
    ...init,
    headers: { ...init?.headers, Authorization: \`Bearer \${access}\` },
  })

  if (res.status === 401) {
    const refresh = (await cookies()).get('customer-refresh')?.value
    if (!refresh) redirect('/login')
    const r = await fetch(\`\${process.env.API_URL}/api/auth/refresh\`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refresh }),
    })
    if (!r.ok) redirect('/login')
    const json = await r.json()
    ;(await cookies()).set('customer-access', json.data.access_token, COOKIE_OPTS)
    res = await fetch(\`\${process.env.API_URL}\${path}\`, {
      ...init,
      headers: { ...init?.headers, Authorization: \`Bearer \${json.data.access_token}\` },
    })
  }
  return res
}`}
      />
      <p>
        Server components call <code>apiFetch</code>. Refresh is
        transparent; the user&apos;s session feels infinite.
      </p>

      <h2>The login form is the same, simpler</h2>
      <p>
        Same pattern, calls <code>/api/auth/login</code> instead of{' '}
        <code>/api/auth/register</code>. Same cookie setting. Same redirect
        on success.
      </p>

      <KnowledgeCheck
        question="A teammate suggests storing the JWT in localStorage 'because it's simpler'. What's the strongest argument against?"
        choices={[
          {
            label: 'Cookies are faster',
            feedback:
              "Marginal difference; not the main concern.",
          },
          {
            label: 'localStorage is readable by any JS — an XSS bug exfiltrates every user\'s tokens. HTTP-only cookies block this attack class entirely.',
            correct: true,
            feedback:
              "Right — that's the OWASP advice and Grit's default. The cost of switching to localStorage isn't worth the increase in blast radius.",
          },
          {
            label: 'Cookies are easier to debug',
            feedback:
              "Not really — both are visible in DevTools.",
          },
          {
            label: 'Next.js doesn\'t support localStorage',
            feedback:
              "Client components can. The issue is security, not capability.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Wire signup end-to-end on your scaffolded project:
            </p>
            <ol>
              <li>Add the form + server action.</li>
              <li>Visit <code>localhost:3000/signup</code>.</li>
              <li>
                Submit. Verify (a) the API created a user (check{' '}
                <code>localhost:8080/studio</code>), (b) cookies set in
                DevTools, (c) you land on <code>/dashboard</code>.
              </li>
            </ol>
            <p>Paste the cookies (just names + types, not values) in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            In Chrome DevTools, Application → Cookies → localhost. Click any
            cookie to see HttpOnly + SameSite flags.
          </>
        }
        solution={
          <>
            <p>You should see two cookies:</p>
            <CodeBlock
              language="text"
              code={`customer-access   HttpOnly  SameSite=Lax   Path=/
customer-refresh  HttpOnly  SameSite=Lax   Path=/`}
            />
            <p>
              Both HttpOnly. Neither readable by browser JS. That&apos;s the
              difference from localStorage.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        User is signed up + signed in. Last lesson of this chapter — what
        do they see? <strong>Dashboard widgets</strong>.
      </p>
    </>
  )
}
