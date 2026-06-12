import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        IDOR — Insecure Direct Object Reference — is the most common
        security bug in shipped APIs. It&apos;s also the easiest to fix.
        This lesson: how it happens, how to exploit it on your own
        endpoint, then ship the fix in three lines.
      </p>

      <h2>The attack, in one diagram</h2>
      <Diagram label="IDOR exploitation" caption="The attacker is authenticated. They just guess (or enumerate) another user's ID and access data they shouldn't.">
{`   User A (legit)               API                   DB
        │                        │                     │
        │  PATCH /api/notes/42   │                     │
        ├───────────────────────►│                     │
        │  (token: A)            │  UPDATE notes       │
        │                        │  WHERE id = 42      │
        │                        ├────────────────────►│
        │                        │  ✓                  │
        │◄───────────────────────┤                     │
        │  200 OK                │                     │
                                                       │
   User B (attacker)             │                     │
        │  PATCH /api/notes/42   │                     │
        ├───────────────────────►│                     │
        │  (token: B)            │  UPDATE notes       │
        │                        │  WHERE id = 42      │  ← no ownership check!
        │                        ├────────────────────►│
        │                        │  ✓                  │
        │◄───────────────────────┤                     │
        │  200 OK — note 42 now belongs to B's content`}
      </Diagram>

      <h2>The vulnerable code</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/note_handler.go — DELIBERATELY BROKEN"
        code={`func (h *NoteHandler) Update(c *gin.Context) {
  id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

  var in services.UpdateNoteInput
  c.ShouldBindJSON(&in)

  // 🚨 BUG: no check that the note belongs to the authed user
  var note models.Note
  h.db.First(&note, id)
  h.db.Model(&note).Updates(map[string]any{"body": in.Body})

  c.JSON(200, gin.H{"data": note})
}`}
      />
      <p>
        Three problems in this code:
      </p>
      <ul>
        <li>
          The handler reads <code>c.Param(&quot;id&quot;)</code>{' '}
          without scoping by the authed user.
        </li>
        <li>
          The handler updates the row whose <code>id</code> matches the
          URL, regardless of ownership.
        </li>
        <li>
          The handler doesn&apos;t even check whether the row exists
          before updating.
        </li>
      </ul>

      <h2>Exploit it yourself</h2>
      <CodeBlock
        language="bash"
        code={`# 1. Register User A
curl -s -X POST localhost:8080/api/auth/register \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"a@a.com","password":"secret123","name":"A"}' | jq

# 2. Get A's token
TOKEN_A=$(curl -s -X POST localhost:8080/api/auth/login \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"a@a.com","password":"secret123"}' | jq -r .data.token)

# 3. A creates a private note
curl -X POST localhost:8080/api/notes \\
  -H "Authorization: Bearer $TOKEN_A" \\
  -H 'Content-Type: application/json' \\
  -d '{"title":"Private","body":"My secret"}'
# returns {data: {id: 1, ...}}

# 4. Register User B (the attacker)
curl -s -X POST localhost:8080/api/auth/register \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"b@b.com","password":"secret123","name":"B"}'

TOKEN_B=$(curl -s -X POST localhost:8080/api/auth/login \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"b@b.com","password":"secret123"}' | jq -r .data.token)

# 5. B overwrites A's note
curl -X PATCH localhost:8080/api/notes/1 \\
  -H "Authorization: Bearer $TOKEN_B" \\
  -H 'Content-Type: application/json' \\
  -d '{"body":"hacked"}'
# returns 200 OK — and note 1 now says "hacked"`}
      />
      <p>
        Five steps. Two minutes. B never had access to A&apos;s account
        but rewrote their data. This is IDOR.
      </p>

      <h2>The fix</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/note_service.go (the fix)"
        code={`var ErrNotFound = errors.New("not found")
var ErrForbidden = errors.New("forbidden")

func (s *NoteService) Update(ctx context.Context, userID, id uint, in UpdateNoteInput) (models.Note, error) {
  var note models.Note
  if err := s.db.WithContext(ctx).First(&note, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) { return note, ErrNotFound }
    return note, err
  }

  // 🔒 Ownership check — the fix
  if note.UserID != userID {
    return note, ErrForbidden
  }

  updates := map[string]any{}
  if in.Title != nil { updates["title"] = *in.Title }
  if in.Body  != nil { updates["body"]  = *in.Body  }
  if err := s.db.WithContext(ctx).Model(&note).Updates(updates).Error; err != nil {
    return note, err
  }
  return note, nil
}`}
      />
      <p>
        Three lines of fix:
      </p>
      <ul>
        <li>Load the note first.</li>
        <li>Compare <code>note.UserID</code> to the authed user.</li>
        <li>Return <code>ErrForbidden</code> if they don&apos;t match.</li>
      </ul>

      <h2>The handler maps to 403</h2>
      <CodeBlock
        language="go"
        code={`switch {
case errors.Is(err, services.ErrNotFound):
  c.JSON(404, gin.H{"error": gin.H{"code": "not_found"}})
case errors.Is(err, services.ErrForbidden):
  c.JSON(403, gin.H{"error": gin.H{"code": "forbidden"}})
case err != nil:
  c.JSON(500, gin.H{"error": gin.H{"code": "internal"}})
default:
  c.JSON(200, gin.H{"data": note})
}`}
      />
      <p>
        Re-run the exploit — B gets 403. A still gets 200. Vulnerability
        closed.
      </p>

      <TipBox tone="warning">
        <strong>404 vs 403 — which to return?</strong> Some teams
        return 404 (&quot;not found&quot;) even when the row exists but
        is forbidden — to avoid leaking the existence of the resource.
        For most apps, 403 is fine. For high-stakes apps (medical,
        legal), 404 is paranoid-correct. Pick a policy and apply it
        everywhere.
      </TipBox>

      <h2>The pattern for ALL access control</h2>
      <ol>
        <li><strong>Load the resource by ID.</strong></li>
        <li>
          <strong>Check authorisation.</strong> Compare
          owner/tenant/role to the authed actor.
        </li>
        <li>
          <strong>Act if allowed.</strong> Update / delete / return.
        </li>
      </ol>
      <p>
        Every CRUD endpoint that operates on a single resource follows
        this. Memorise it. <code>grit generate resource</code> emits
        this pattern by default — the bug only appears when humans
        write handlers and forget step 2.
      </p>

      <h2>Defense in depth — the SQL filter</h2>
      <p>
        Belt-and-suspenders: scope the QUERY ITSELF by user_id, not
        just check after loading:
      </p>
      <CodeBlock
        language="go"
        code={`if err := s.db.WithContext(ctx).
  Where("user_id = ?", userID).
  First(&note, id).Error; err != nil { ... }`}
      />
      <p>
        Now even if the post-load check is removed by a future
        refactor, the query can&apos;t return another user&apos;s
        note. Two layers; both should fail to leak data.
      </p>

      <h2>Test that the fix sticks</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/note_handler_test.go"
        code={`func TestUpdateNote_IDORBlocked(t *testing.T) {
  s := testServices(t)
  a := registerTestUser(t, s, "a@a.com")
  b := registerTestUser(t, s, "b@b.com")

  note, _ := s.Notes.Create(ctx, a.ID, services.CreateNoteInput{Title: "x", Body: "y"})

  _, err := s.Notes.Update(ctx, b.ID, note.ID, services.UpdateNoteInput{Body: ptr("hacked")})
  require.ErrorIs(t, err, services.ErrForbidden)
}`}
      />
      <p>
        Now if a future refactor accidentally removes the check, this
        test fails and CI blocks the merge.
      </p>

      <KnowledgeCheck
        question="A reviewer says &quot;just return 404 instead of 403 to be more secure&quot;. What's the actual security benefit?"
        choices={[
          {
            label: "Faster response — saves a few ms",
            feedback: "Same DB hit. No measurable speed difference.",
          },
          {
            label: 'It hides whether the resource exists at all — an attacker can\'t enumerate IDs to map your data layout',
            correct: true,
            feedback:
              "Right — 404 says &quot;nothing here or you can't see it&quot;; 403 confirms the resource exists. For typical apps this matters little; for high-value data (medical records, M&A docs), 404-for-everything-not-yours is the paranoid-correct choice.",
          },
          {
            label: 'HTTP 403 is deprecated',
            feedback: 'Not at all — 403 is a perfectly standard response.',
          },
          {
            label: 'Browsers handle 404 better',
            feedback: 'Browsers handle both identically for API responses.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Find and fix an IDOR in your own code:</p>
            <ol>
              <li>
                Search your handlers for any{' '}
                <code>c.Param(&quot;id&quot;)</code> that&apos;s used in
                a DB query without a corresponding user_id check.
              </li>
              <li>
                Pick one and write the exploit as curl commands (like
                this lesson&apos;s). Confirm you can read or modify
                another user&apos;s data.
              </li>
              <li>
                Move the DB call into the service (if not already
                there). Add the ownership check.
              </li>
              <li>Re-run the exploit. Should now return 403.</li>
              <li>
                Write a unit test that fails if the check is removed.
              </li>
              <li>
                Commit the fix + test in one commit. Reference the
                threat from SECURITY.md.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If <code>grit generate resource</code> emitted your code,
            it&apos;s probably already correct — the bug is most often
            in hand-written handlers. Look for endpoints with custom
            logic (publish, share, transfer ownership).
          </>
        }
        solution={
          <>
            <p>
              You&apos;ve eliminated the most common API vulnerability
              from at least one endpoint. Now do every other endpoint
              with the same pattern. This is what mature teams call
              &quot;authz hardening&quot; and it&apos;s 80% of the
              work.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>The authz package</strong>. Grit ships
        an <code>internal/authz</code> helper that centralises this
        pattern so you write it once, apply it everywhere.
      </p>
    </>
  )
}
