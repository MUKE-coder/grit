import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Role-based UI: hide buttons the user can&apos;t use. Two layers — the
        UI hides actions for clean UX, the API rejects unauthorized
        requests for security. Both matter; neither replaces the other.
      </p>

      <h2>The roles you ship with</h2>
      <CodeBlock
        language="ts"
        filename="packages/shared/src/constants/roles.ts"
        code={`export const ROLES = {
  USER:   'user',     // default — read-only mostly
  STAFF:  'staff',    // can edit
  ADMIN:  'admin',    // full access
  OWNER:  'owner',    // tenant owner — billing too
} as const

export type UserRole = typeof ROLES[keyof typeof ROLES]

export const ROLE_HIERARCHY: Record<UserRole, number> = {
  user: 1, staff: 2, admin: 3, owner: 4,
}

export function hasRole(userRole: UserRole, required: UserRole): boolean {
  return ROLE_HIERARCHY[userRole] >= ROLE_HIERARCHY[required]
}`}
      />
      <p>
        Hierarchical: <code>admin</code> includes <code>staff</code>{' '}
        includes <code>user</code>. Owner gets billing + admin actions.
      </p>

      <h2>The Can component</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/components/can.tsx"
        code={`'use client'
import { useCurrentUser } from '@/hooks/use-current-user'
import { hasRole, type UserRole } from '@workspace/shared/constants/roles'

interface Props {
  role: UserRole
  children: React.ReactNode
}

export function Can({ role, children }: Props) {
  const { data: user } = useCurrentUser()
  if (!user) return null
  return hasRole(user.role, role) ? <>{children}</> : null
}`}
      />
      <p>Use it anywhere:</p>
      <CodeBlock
        language="tsx"
        code={`<Can role="admin">
  <Button onClick={deleteCustomer}>Delete customer</Button>
</Can>

<Can role="staff">
  <Button onClick={refundOrder}>Issue refund</Button>
</Can>`}
      />

      <h2>Gating entire pages</h2>
      <p>
        For pages a role shouldn&apos;t see at all (e.g.,{' '}
        <code>/admin/billing</code> for non-owners), redirect in the
        server component:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(app)/billing/page.tsx"
        code={`import { redirect } from 'next/navigation'
import { getCurrentUser } from '@/lib/auth'

export default async function BillingPage() {
  const user = await getCurrentUser()
  if (!user) redirect('/login')
  if (!hasRole(user.role, 'owner')) redirect('/dashboard')

  return <BillingDashboard />
}`}
      />

      <TipBox tone="warning">
        <strong>UI gates are UX, not security.</strong> Anyone with the
        right URL + a browser can still hit your API. The API must reject
        the unauthorized request independently — that&apos;s the actual
        defence. UI gates just hide what they can&apos;t use anyway.
      </TipBox>

      <h2>useCan hook for conditional logic</h2>
      <CodeBlock
        language="ts"
        code={`export function useCan(role: UserRole) {
  const { data: user } = useCurrentUser()
  if (!user) return false
  return hasRole(user.role, role)
}

// In a component:
const canRefund = useCan('staff')
return (
  <Button disabled={!canRefund || isProcessing}>
    {canRefund ? 'Issue refund' : 'Refunds require staff role'}
  </Button>
)`}
      />
      <p>
        Same outcome as <code>&lt;Can&gt;</code> but lets you use the
        boolean in JS logic (disabled state, conditional styling, etc.).
      </p>

      <h2>Server-side role checks for actions</h2>
      <p>
        For server actions / route handlers, check the role too:
      </p>
      <CodeBlock
        language="ts"
        code={`'use server'
export async function deleteCustomer(id: string) {
  const user = await getCurrentUser()
  if (!user || !hasRole(user.role, 'admin')) {
    throw new Error('Forbidden')
  }
  return apiFetch(\`/api/customers/\${id}\`, { method: 'DELETE' })
}`}
      />
      <p>
        And the Go API has its own <code>RequireRoles</code> middleware
        (from the Go API course ch.3). Defence in depth: web hides, server
        action checks, API rejects.
      </p>

      <KnowledgeCheck
        question={<>The &quot;Delete Customer&quot; button is wrapped in <code>{`<Can role="admin">`}</code>. A staff user opens DevTools and removes the button&apos;s parent element from the DOM. The button now appears. What happens when they click it?</>}
        choices={[
          {
            label: 'Customer is deleted — the JS bypass worked',
            feedback:
              "Only if there's no server-side check. With Grit's RequireRoles middleware on the API, the request returns 404. Defence in depth pays off here.",
          },
          {
            label: 'Server action checks `hasRole(user.role, \"admin\")` AND the API has RequireRoles — both reject the request',
            correct: true,
            feedback:
              "Right — UI is UX, server-side is security. Two independent checks both reject. The DevTools-mucking attacker gets a 403/404.",
          },
          {
            label: 'Next.js prevents the click',
            feedback:
              "Wrong — Next.js can't override a user-modified DOM. The click fires.",
          },
          {
            label: 'The button does nothing because the onClick handler unmounted',
            feedback:
              "Depends on how the user manipulated the DOM. Assuming the button is still wired, the request fires.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Add role gating to your scaffolded project:</p>
            <ol>
              <li>
                In <code>apps/admin/app/resources/customers/page.tsx</code>,
                wrap the Delete action in <code>&lt;Can role=&quot;admin&quot;&gt;</code>.
              </li>
              <li>
                Promote one of your test users to <code>staff</code> (not
                admin). Log in as them. The Delete button should be
                hidden.
              </li>
              <li>
                Open DevTools, manually click the API endpoint that
                deletes (curl or fetch). Confirm the API rejects with{' '}
                404 — that&apos;s Grit&apos;s RequireRoles middleware.
              </li>
            </ol>
            <p>Paste the curl response in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            For the server-side check to work, the API must have{' '}
            <code>RequireRoles(&quot;admin&quot;)</code> middleware on the delete
            route. From the Go API course you know this is one line in{' '}
            <code>routes.go</code>.
          </>
        }
        solution={
          <>
            <p>The expected curl response:</p>
            <CodeBlock
              language="json"
              code={`HTTP/1.1 404 Not Found
{ "error": { "code": "NOT_FOUND", "message": "not found" } }`}
            />
            <p>
              404 not 403 — leaks no info about whether the route exists.
              That&apos;s the convention from the Defender&apos;s Handbook page.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson — the <strong>invitation flow</strong>. Inviting a
        teammate by email, completing signup, role + tenant assigned
        atomically.
      </p>
    </>
  )
}
