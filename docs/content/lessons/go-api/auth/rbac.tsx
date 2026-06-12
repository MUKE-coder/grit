import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Authentication tells you <em>who</em> they are. Authorization tells
        you <em>what they can do</em>. Grit splits this into two primitives:
        roles (system-wide groups) and ownership (per-row access). Both ship
        as middleware.
      </p>

      <h2>Roles — the system-wide group</h2>
      <p>
        Every User has a <code>Role</code> field (default <code>user</code>).
        Grit ships three roles to start: <code>user</code>,{' '}
        <code>staff</code>, <code>admin</code>. Add your own
        (<code>auditor</code>, <code>support</code>, etc.) by editing the
        User model.
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go"
        code={`adminGroup := api.Group("/admin")
adminGroup.Use(middleware.Auth(authService))
adminGroup.Use(middleware.RequireRoles("admin"))
{
    adminGroup.GET("/users", userHandler.AdminList)
    adminGroup.DELETE("/users/:id", userHandler.AdminDelete)
}`}
      />
      <p>
        <code>RequireRoles</code> reads <code>role</code> from the JWT and
        returns 404 (not 403) if the user isn&apos;t in the allow-list. 404 is
        deliberate — it doesn&apos;t leak the route&apos;s existence to
        unauthorized users.
      </p>

      <h2>Ownership — the per-row check (IDOR defence)</h2>
      <p>
        Roles answer &quot;can this user touch admin routes?&quot;. Ownership
        answers &quot;can this user touch <em>this specific row</em>?&quot;.
        Without it, alice can read bob&apos;s order by changing the URL from{' '}
        <code>/api/orders/123</code> to <code>/api/orders/456</code>. That&apos;s
        the OWASP A01:2025 #1 risk class — IDOR.
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/order.go"
        code={`func (h *OrderHandler) GetByID(c *gin.Context) {
    var order models.Order
    if err := authz.MustOwn(c, h.db, &order, c.Param("id")); err != nil {
        return // helper already wrote 404 — never 403
    }
    c.JSON(http.StatusOK, gin.H{"data": order})
}

// Model must implement Ownable
func (o *Order) GetOwnerID() string { return o.UserID }`}
      />
      <p>
        <code>authz.MustOwn</code> loads the row, checks{' '}
        <code>GetOwnerID()</code> matches the authenticated user, returns
        404 on mismatch. The 404 (not 403) closes the existence-leak side
        channel.
      </p>

      <TipBox tone="success">
        <strong>The pattern is so important that the code generator wires
        it for you.</strong> When you run <code>grit generate resource
        Order</code>, the generated handler already calls{' '}
        <code>authz.MustOwn</code>. IDOR is closed by the generator, not by
        the developer remembering.
      </TipBox>

      <h2>Tenant scoping</h2>
      <p>
        Multi-tenant apps go one step further — a row belongs to a tenant,
        and users only see rows for their tenant. Grit&apos;s{' '}
        <code>authz.CheckScope</code> handles this:
      </p>
      <CodeBlock
        language="go"
        code={`func (h *OrderHandler) List(c *gin.Context) {
    tenantID := c.GetString("tenant_id")  // set by middleware
    var orders []models.Order
    h.db.Where("tenant_id = ?", tenantID).Find(&orders)
    c.JSON(200, gin.H{"data": orders})
}`}
      />
      <p>
        Every query is tenant-scoped. There&apos;s no way for a handler to forget
        — code review (or a custom lint rule) catches{' '}
        <code>db.Find()</code> without a tenant filter.
      </p>

      <h2>The invitation flow</h2>
      <p>For team-based products: invite a user to a tenant + role.</p>
      <CodeBlock
        language="text"
        code={`Admin clicks "Invite member" → enters email + role
  → POST /api/invitations  { email, role, tenant_id }
       → Grit creates an invitation row + sends an email with a one-time link
Invitee clicks the link → /api/invitations/:token/accept
  → Grit creates the user (if new), assigns the role + tenant, signs them in`}
      />
      <p>
        Single-use, time-bound (7 days default). Stored in the activity log
        so you can audit who invited whom.
      </p>

      <KnowledgeCheck
        question="Alex (regular user) sends `GET /api/admin/users/123` to inspect another user. What does Grit do?"
        choices={[
          {
            label: 'Returns 403 Forbidden — Alex isn\'t admin',
            feedback:
              "Wrong — 403 leaks that the route exists. Grit returns 404 instead so Alex can't tell whether /api/admin/users/123 is a real route or not.",
          },
          {
            label: 'Returns 404 Not Found — RequireRoles middleware rejects the request',
            correct: true,
            feedback:
              "Right — RequireRoles returns 404 to non-matching users. From Alex's perspective, the admin routes simply don't exist.",
          },
          {
            label: 'Returns 200 OK with the user data',
            feedback:
              "Catastrophically wrong — that's the IDOR vulnerability you're supposed to prevent. Grit's defaults stop it.",
          },
          {
            label: 'Returns 401 Unauthorized',
            feedback:
              "Wrong — Alex IS authenticated (has a valid token). 401 means 'no auth'; this is more nuanced — auth is fine, role is insufficient.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Implement the chapter assignment: add{' '}
              <code>/api/admin/stats</code> that returns user counts. Protect
              it with both Auth and RequireRoles. Verify the right
              responses:
            </p>
            <ol>
              <li>
                Without a token: 401
              </li>
              <li>
                With a regular-user token: 404 (RequireRoles rejected you)
              </li>
              <li>
                With an admin token (manually set <code>role=admin</code> in
                your DB row): 200 + the stats
              </li>
            </ol>
            <p>Paste all three curl responses in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            Promote your test user via SQL:{' '}
            <code>UPDATE users SET role = &apos;admin&apos; WHERE email = &apos;alex@example.com&apos;;</code>{' '}
            Then re-login to get a token with the new role.
          </>
        }
        solution={
          <>
            <p>The three responses you should capture:</p>
            <CodeBlock
              language="text"
              code={`# No token
< HTTP/1.1 401 Unauthorized
{ "error": { "code": "UNAUTHORIZED", "message": "authentication required" } }

# Regular user
< HTTP/1.1 404 Not Found
{ "error": { "code": "NOT_FOUND", "message": "not found" } }

# Admin
< HTTP/1.1 200 OK
{ "data": { "user_count": 5, "admin_count": 1 } }`}
            />
            <p>
              Chapter 3 done. You can now lock down routes by both role and
              row-ownership.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 4 — <strong>Batteries: jobs, mail, storage, AI</strong>. The
        included primitives that ship with every Grit API and save you a
        week of integration work.
      </p>
    </>
  )
}
