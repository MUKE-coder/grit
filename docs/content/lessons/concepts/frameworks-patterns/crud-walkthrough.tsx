import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { PlaygroundChallenge } from '@/components/course/playground-challenge'

export default function Lesson() {
  return (
    <>
      <p>
        Four operations, four URLs, four handlers, four service
        methods. We&apos;ll do all of them for a <code>Note</code>{' '}
        resource — by hand, no generator. After this lesson, any
        resource you ever build is a variation of this template.
      </p>

      <h2>The model</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/note.go"
        code={`package models

import "time"

type Note struct {
  ID        uint      \`gorm:"primaryKey"\`
  UserID    uint      \`gorm:"index;not null"\`
  Title     string    \`gorm:"not null"\`
  Body      string    \`gorm:"type:text"\`
  CreatedAt time.Time
  UpdatedAt time.Time
}`}
      />
      <p>
        Add it to <code>db.AutoMigrate(&amp;models.Note&#123;&#125;)</code>{' '}
        so the table is created on next startup.
      </p>

      <h2>The four routes</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go (add to the authed group)"
        code={`authed.GET   ("/notes",     s.NoteHandler.List)
authed.POST  ("/notes",     s.NoteHandler.Create)
authed.PATCH ("/notes/:id", s.NoteHandler.Update)
authed.DELETE("/notes/:id", s.NoteHandler.Delete)`}
      />
      <p>
        REST conventions: noun is plural; collection vs. item via{' '}
        <code>/:id</code>; verb encoded as HTTP method. Don&apos;t
        invent <code>/notes/create</code>.
      </p>

      <h2>READ many — GET /api/notes</h2>

      <h3>Handler</h3>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/note_handler.go"
        code={`func (h *NoteHandler) List(c *gin.Context) {
  userID := c.GetUint("user_id")           // set by RequireAuth
  page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
  size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

  notes, total, err := h.notes.List(c.Request.Context(), userID, page, size)
  if err != nil {
    c.JSON(500, gin.H{"error": gin.H{"code": "internal", "message": err.Error()}})
    return
  }
  c.JSON(200, gin.H{
    "data": notes,
    "meta": gin.H{"total": total, "page": page, "page_size": size},
  })
}`}
      />

      <h3>Service</h3>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/note_service.go"
        code={`func (s *NoteService) List(ctx context.Context, userID uint, page, size int) ([]models.Note, int64, error) {
  tx := s.db.WithContext(ctx).Model(&models.Note{}).Where("user_id = ?", userID)

  var total int64
  if err := tx.Count(&total).Error; err != nil { return nil, 0, err }

  var notes []models.Note
  err := tx.Order("created_at DESC").
    Offset((page - 1) * size).
    Limit(size).
    Find(&notes).Error
  return notes, total, err
}`}
      />
      <p>
        Always scope by <code>user_id</code> — a logged-in user reads
        only THEIR notes. Without this filter, your list endpoint
        leaks every user&apos;s notes to every user. This is the IDOR
        vulnerability in the OWASP Top 10.
      </p>

      <h2>CREATE — POST /api/notes</h2>

      <h3>Input struct + validation</h3>
      <CodeBlock
        language="go"
        code={`type CreateNoteInput struct {
  Title string \`json:"title" binding:"required,max=200"\`
  Body  string \`json:"body"  binding:"max=10000"\`
}`}
      />
      <p>
        The <code>binding</code> tags are read by{' '}
        <code>c.ShouldBindJSON</code>. If <code>title</code> is missing
        or longer than 200 chars, Gin returns an error you respond with
        as 400.
      </p>

      <h3>Handler</h3>
      <CodeBlock
        language="go"
        code={`func (h *NoteHandler) Create(c *gin.Context) {
  var in services.CreateNoteInput
  if err := c.ShouldBindJSON(&in); err != nil {
    c.JSON(400, gin.H{"error": gin.H{"code": "invalid_body", "message": err.Error()}})
    return
  }
  userID := c.GetUint("user_id")
  note, err := h.notes.Create(c.Request.Context(), userID, in)
  if err != nil {
    c.JSON(500, gin.H{"error": gin.H{"code": "internal", "message": err.Error()}})
    return
  }
  c.JSON(201, gin.H{"data": note, "message": "Note created"})
}`}
      />

      <h3>Service</h3>
      <CodeBlock
        language="go"
        code={`func (s *NoteService) Create(ctx context.Context, userID uint, in CreateNoteInput) (models.Note, error) {
  note := models.Note{UserID: userID, Title: in.Title, Body: in.Body}
  if err := s.db.WithContext(ctx).Create(&note).Error; err != nil {
    return models.Note{}, fmt.Errorf("create note: %w", err)
  }
  return note, nil
}`}
      />
      <p>
        Status code 201 (Created), not 200, on successful create. Tiny
        detail; pros notice.
      </p>

      <h2>UPDATE — PATCH /api/notes/:id</h2>
      <p>
        Why <code>PATCH</code> and not <code>PUT</code>? PATCH is for
        partial updates — the client sends only fields that change.
        PUT replaces the whole resource. Almost always you want PATCH.
      </p>
      <CodeBlock
        language="go"
        code={`type UpdateNoteInput struct {
  Title *string \`json:"title"\` // pointer = optional
  Body  *string \`json:"body"\`
}

func (h *NoteHandler) Update(c *gin.Context) {
  id, err := strconv.ParseUint(c.Param("id"), 10, 64)
  if err != nil {
    c.JSON(400, gin.H{"error": gin.H{"code": "invalid_id", "message": err.Error()}})
    return
  }
  var in services.UpdateNoteInput
  if err := c.ShouldBindJSON(&in); err != nil {
    c.JSON(400, gin.H{"error": gin.H{"code": "invalid_body", "message": err.Error()}})
    return
  }
  userID := c.GetUint("user_id")
  note, err := h.notes.Update(c.Request.Context(), userID, uint(id), in)
  switch {
  case errors.Is(err, services.ErrNotFound):
    c.JSON(404, gin.H{"error": gin.H{"code": "not_found"}})
  case errors.Is(err, services.ErrForbidden):
    c.JSON(403, gin.H{"error": gin.H{"code": "forbidden"}})
  case err != nil:
    c.JSON(500, gin.H{"error": gin.H{"code": "internal"}})
  default:
    c.JSON(200, gin.H{"data": note, "message": "Note updated"})
  }
}`}
      />
      <CodeBlock
        language="go"
        code={`func (s *NoteService) Update(ctx context.Context, userID, id uint, in UpdateNoteInput) (models.Note, error) {
  var note models.Note
  if err := s.db.WithContext(ctx).First(&note, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) { return note, ErrNotFound }
    return note, err
  }
  if note.UserID != userID {
    return note, ErrForbidden   // IDOR check: this user can't edit someone else's note
  }
  updates := map[string]any{}
  if in.Title != nil { updates["title"] = *in.Title }
  if in.Body  != nil { updates["body"]  = *in.Body  }
  if len(updates) == 0 { return note, nil }
  if err := s.db.WithContext(ctx).Model(&note).Updates(updates).Error; err != nil {
    return note, err
  }
  return note, nil
}`}
      />
      <p>
        Three patterns to absorb:
      </p>
      <ul>
        <li>
          <strong>Pointer fields for optional input.</strong>{' '}
          <code>*string</code> distinguishes &quot;not provided&quot;{' '}
          (nil) from &quot;provided as empty string&quot; (&quot;&quot;).
          Critical for PATCH semantics.
        </li>
        <li>
          <strong>Two-step: load + authorize.</strong> First load by
          ID; then check the owner; only then update. This is the IDOR
          defence.
        </li>
        <li>
          <strong>Map-based Updates.</strong> Pass{' '}
          <code>db.Model(&amp;note).Updates(map)</code> with only the
          set fields — GORM updates exactly those columns.
        </li>
      </ul>

      <h2>DELETE — DELETE /api/notes/:id</h2>
      <CodeBlock
        language="go"
        code={`func (h *NoteHandler) Delete(c *gin.Context) {
  id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
  userID := c.GetUint("user_id")
  err := h.notes.Delete(c.Request.Context(), userID, uint(id))
  switch {
  case errors.Is(err, services.ErrNotFound):
    c.JSON(404, gin.H{"error": gin.H{"code": "not_found"}})
  case errors.Is(err, services.ErrForbidden):
    c.JSON(403, gin.H{"error": gin.H{"code": "forbidden"}})
  case err != nil:
    c.JSON(500, gin.H{"error": gin.H{"code": "internal"}})
  default:
    c.Status(204)   // No Content — successful delete, nothing to return
  }
}

func (s *NoteService) Delete(ctx context.Context, userID, id uint) error {
  var note models.Note
  if err := s.db.WithContext(ctx).First(&note, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) { return ErrNotFound }
    return err
  }
  if note.UserID != userID { return ErrForbidden }
  return s.db.WithContext(ctx).Delete(&note).Error
}`}
      />
      <p>
        <strong>204 No Content</strong> is the right status for a
        delete with nothing to return. The browser/client knows the
        delete worked; no body needed.
      </p>

      <TipBox tone="warning">
        <strong>Always scope by user_id.</strong> Three of the four
        endpoints above include <code>WHERE user_id = ?</code> or check{' '}
        <code>note.UserID != userID</code>. Skip this and your delete
        endpoint becomes &quot;delete ANYONE&apos;s note if you know
        the ID&quot;. This is IDOR — silent, catastrophic, common.
        Always scope. Always check.
      </TipBox>

      <h2>Testing the whole cycle</h2>
      <CodeBlock
        language="bash"
        code={`# Register and capture token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/register \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"a@b.com","password":"secret123","name":"A"}' \\
  | jq -r .data.token)

# Create
curl -X POST http://localhost:8080/api/notes \\
  -H "Authorization: Bearer $TOKEN" \\
  -H 'Content-Type: application/json' \\
  -d '{"title":"First","body":"hello"}'

# List
curl http://localhost:8080/api/notes -H "Authorization: Bearer $TOKEN"

# Update
curl -X PATCH http://localhost:8080/api/notes/1 \\
  -H "Authorization: Bearer $TOKEN" \\
  -H 'Content-Type: application/json' \\
  -d '{"title":"Renamed"}'

# Delete
curl -X DELETE http://localhost:8080/api/notes/1 \\
  -H "Authorization: Bearer $TOKEN"`}
      />

      <PlaygroundChallenge title="Hand-write CRUD for Task">
        <p>
          In the playground, build the four endpoints for a{' '}
          <code>Task</code> resource (id, user_id, title, done).
          Follow the exact pattern from this lesson. Time yourself —
          aim to do it in 20 minutes. The point is the muscle memory
          of handler → service → GORM.
        </p>
      </PlaygroundChallenge>

      <KnowledgeCheck
        question="Why does the Update handler use *string fields in its input struct instead of plain string?"
        choices={[
          {
            label: "Because plain string can't hold UTF-8",
            feedback: 'Plain string holds UTF-8 fine. The reason is something else.',
          },
          {
            label: "Pointers let you tell &quot;not provided&quot; from &quot;provided as empty string&quot; — essential for PATCH semantics where the client only sends fields they want to change",
            correct: true,
            feedback:
              "Right — nil means &quot;don't touch&quot;; empty string means &quot;set to empty&quot;. With plain string you can't tell them apart, so a client trying to update only Body would accidentally blank the Title.",
          },
          {
            label: 'Pointers are faster',
            feedback:
              "Marginal. The reason is semantics, not speed.",
          },
          {
            label: 'GORM requires pointers',
            feedback: "It doesn't. This is a Gin/JSON parsing concern, not a GORM one.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For chapter 5&apos;s assignment, hand-write a Note resource:</p>
            <ol>
              <li>Add the model + auto-migrate.</li>
              <li>Wire all four routes in <code>routes.go</code>.</li>
              <li>
                Write the handler + service for all four operations
                following this lesson&apos;s pattern.
              </li>
              <li>
                Test all four via curl, INCLUDING the IDOR check: log
                in as User A, create a note, then log in as User B and
                try to PATCH it — must return 403.
              </li>
              <li>
                Now run <code>grit generate resource Note user_id:uint title:string body:string</code>{' '}
                in a fresh project. DIFF the generated files against
                yours. Note where the generated version differs (it
                probably does some things better, some things you
                preferred yours).
              </li>
              <li>One paragraph in notes.md on what surprised you.</li>
            </ol>
          </>
        }
        hint={
          <>
            For the IDOR test, the easiest path is two terminals: one
            with User A&apos;s token, one with User B&apos;s. Or use a
            REST client with two environments.
          </>
        }
        solution={
          <>
            <p>
              You now have the pattern that the generator automates and
              an intuition for what it&apos;s automating. From this
              chapter forward, when something looks confusing, you can
              read the generated code and understand every line.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 6 — <strong>The Batteries</strong>. The included
        services that turn this clean CRUD foundation into a
        production-ready API: cache, file storage, email, jobs, AI.
      </p>
    </>
  )
}
