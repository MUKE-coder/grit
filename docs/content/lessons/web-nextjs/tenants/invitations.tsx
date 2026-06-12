import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Team invitations: admin types an email + role, system emails a
        one-time link, recipient clicks, sets a password, lands in your
        tenant with the right role. Grit ships this end-to-end. This
        lesson covers what&apos;s wired and the extension points.
      </p>

      <h2>The flow</h2>
      <CodeBlock
        language="text"
        code={`Admin clicks "Invite member" → modal with email + role
  → POST /api/invitations { email, role }
       → Grit creates an invitation row (UUID token, 7-day expiry)
       → Job worker sends an email with /accept/<token>

Invitee opens link → /accept/abc123
  → /api/invitations/abc123  returns { email, role, tenant_name }
  → Frontend shows "Join Acme as Staff" + password form
  → POST /api/invitations/abc123/accept { password }
       → Grit creates the user (tenant + role from the invitation)
       → Marks invitation as accepted
       → Returns access + refresh tokens — user is signed in`}
      />

      <h2>The model</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/invitation.go"
        code={`type Invitation struct {
    ID         uuid.UUID  \`gorm:"type:uuid;primaryKey" json:"id"\`
    Token      string     \`gorm:"uniqueIndex;not null" json:"-"\`     // never serialise
    Email      string     \`gorm:"index;not null" json:"email"\`
    Role       string     \`gorm:"not null" json:"role"\`
    TenantID   uuid.UUID  \`gorm:"type:uuid;not null" json:"tenant_id"\`
    InvitedBy  uuid.UUID  \`gorm:"type:uuid;not null" json:"invited_by"\`
    AcceptedAt *time.Time \`json:"accepted_at"\`
    ExpiresAt  time.Time  \`gorm:"index;not null" json:"expires_at"\`
    CreatedAt  time.Time  \`json:"created_at"\`
}`}
      />
      <p>
        Single-use (set <code>AcceptedAt</code> on accept; reject if
        already set), time-bound (default 7 days), bound to email + role +
        tenant.
      </p>

      <h2>The invite UI</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/components/team/invite-modal.tsx"
        code={`'use client'
import { useState } from 'react'

export function InviteMemberModal({ onClose }) {
  const [email, setEmail] = useState('')
  const [role, setRole] = useState<'user' | 'staff' | 'admin'>('staff')

  async function onSubmit() {
    await fetch('/api/invitations', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, role }),
    })
    toast.success('Invitation sent')
    onClose()
  }

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent>
        <DialogTitle>Invite a teammate</DialogTitle>
        <input value={email} onChange={e => setEmail(e.target.value)} placeholder="alex@example.com" />
        <select value={role} onChange={e => setRole(e.target.value as any)}>
          <option value="user">User</option>
          <option value="staff">Staff</option>
          <option value="admin">Admin</option>
        </select>
        <Button onClick={onSubmit}>Send invitation</Button>
      </DialogContent>
    </Dialog>
  )
}`}
      />

      <h2>The accept page</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(auth)/accept/[token]/page.tsx"
        code={`export default async function AcceptPage({ params }) {
  const inv = await fetch(\`\${API}/api/invitations/\${params.token}\`).then(r => r.json())
  if (!inv?.data) return <p>This invitation is invalid or expired.</p>

  return (
    <AcceptForm
      token={params.token}
      email={inv.data.email}
      tenantName={inv.data.tenant_name}
      role={inv.data.role}
    />
  )
}`}
      />
      <p>
        Show context — &quot;Join Acme as Staff&quot; — so the recipient
        knows what they&apos;re accepting. Then a password form. Submit
        creates the user atomically with the right tenant + role.
      </p>

      <TipBox tone="success">
        <strong>Atomic creation:</strong> the API creates the user and
        marks the invitation accepted in a single transaction. If anything
        fails, both roll back. No half-created users with no invitation
        record.
      </TipBox>

      <h2>Re-sending + revoking</h2>
      <p>Sensible affordances:</p>
      <ul>
        <li>
          <strong>Re-send</strong> — generate a new token + email. The old
          one is invalidated.
        </li>
        <li>
          <strong>Revoke</strong> — delete the invitation. Even if the
          recipient still has the email, the link returns 404.
        </li>
        <li>
          <strong>Expired view</strong> — show pending invitations in the
          admin so admins can clean up stale ones.
        </li>
      </ul>

      <h2>Edge cases worth handling</h2>
      <ul>
        <li>
          <strong>User already exists</strong> — if alex@example.com is
          already in the system (different tenant), how do you handle
          joining a second tenant? Usually: create a tenant-membership
          row, let the user switch between tenants.
        </li>
        <li>
          <strong>Email mismatch on accept</strong> — invitation says{' '}
          <code>alex@example.com</code>; recipient is logged in as a
          different account. Force them to log out first.
        </li>
        <li>
          <strong>Email goes to spam</strong> — show admins a &quot;Copy invite
          link&quot; button as a fallback.
        </li>
      </ul>

      <KnowledgeCheck
        question="Bob's invitation expired (7-day window passed). He clicks the link anyway. What should happen?"
        choices={[
          {
            label: 'The link silently regenerates a new token',
            feedback:
              "Wrong — that bypasses the security guarantee. Expired means expired; require an explicit re-invite.",
          },
          {
            label: 'A clear error: \'This invitation has expired. Ask your admin to re-send.\' + nothing else',
            correct: true,
            feedback:
              "Right — clear message, no creep. The admin can revoke + re-invite. The expiry exists exactly to prevent forever-valid links.",
          },
          {
            label: 'Show the password form anyway',
            feedback:
              "Wrong — let the user set a password against an expired invitation and you've broken your own security contract.",
          },
          {
            label: 'Auto-create the account with the original email',
            feedback:
              "Same security failure. The expiry MUST mean expired.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For chapter 5&apos;s final assignment, end-to-end invitation:</p>
            <ol>
              <li>
                As an admin user, POST to{' '}
                <code>/api/invitations</code> with another email +{' '}
                <code>role=staff</code>.
              </li>
              <li>
                Check Mailhog (<code>localhost:8025</code>) — the invitation
                email should be there.
              </li>
              <li>
                Click the accept link. Fill in the password form.
              </li>
              <li>
                Verify you land signed in as the new user, with the staff
                role, in the inviting admin&apos;s tenant.
              </li>
              <li>
                As that staff user, visit a page that gates on{' '}
                <code>admin</code> role — should be hidden (from previous
                lesson).
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If Mailhog doesn&apos;t show the email, the invitation job is
            failing — check the worker logs.
          </>
        }
        solution={
          <>
            <p>
              The captured email looks like:
            </p>
            <CodeBlock
              language="text"
              code={`From: noreply@yourapp.dev
To: bob@example.com
Subject: You're invited to join Acme on Grit Field Service

[Invitation link]
http://localhost:3000/accept/abc123def...`}
            />
            <p>
              Click → password form → signed in as a real user. Multi-
              tenant SaaS shape end-to-end. Chapter 5 done.
            </p>
          </>
        }
      />

      <h2>You finished Building Web with Next.js + Go API 🎉</h2>
      <p>
        Five chapters, 13 lessons. You can now scaffold the triple kit,
        build a marketing landing page with proper SEO, wire signup +
        dashboard widgets, customise an admin panel via{' '}
        <code>defineResource()</code>, and handle multi-tenancy + roles +
        invitations cleanly.
      </p>
      <p>
        From here: try{' '}
        <a href="/courses/multiplatform" className="text-primary hover:underline">
          Building Web + Desktop + Mobile
        </a>{' '}
        to add native surfaces, or go deeper on the API with the{' '}
        <a href="/docs/plugins" className="text-primary hover:underline">
          Grit plugins
        </a>{' '}
        (Stripe, WebSockets, OAuth).
      </p>
    </>
  )
}
