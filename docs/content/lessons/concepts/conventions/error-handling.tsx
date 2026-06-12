import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Errors are normal; how you handle them is what makes the code
        production-ready. Grit has one pattern for Go (explicit + wrapped)
        and one for React (error boundary + toast). Both close the loop so
        no error is silently swallowed.
      </p>

      <h2>Go side — explicit, wrapped, never ignored</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/user.go"
        code={`func (s *UserService) Create(ctx context.Context, in CreateUserInput) (*User, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("hashing password: %w", err)
    }

    user := &User{Email: in.Email, PasswordHash: hashed}
    if err := s.db.Create(user).Error; err != nil {
        return nil, fmt.Errorf("creating user: %w", err)
    }
    return user, nil
}`}
      />
      <p>The three rules:</p>
      <ul>
        <li>
          <strong>Every error is checked.</strong> Never{' '}
          <code>_ := someFunc()</code>.
        </li>
        <li>
          <strong>Errors are wrapped with context.</strong>{' '}
          <code>fmt.Errorf("doing X: %w", err)</code>. The <code>%w</code> makes
          it unwrappable upstream (with <code>errors.Is</code> /{' '}
          <code>errors.As</code>).
        </li>
        <li>
          <strong>Handlers convert errors to HTTP responses.</strong> Services
          return errors; handlers map them to status codes via a tiny helper.
        </li>
      </ul>

      <h2>The handler-side helper</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/user.go"
        code={`func (h *UserHandler) Create(c *gin.Context) {
    var in services.CreateUserInput
    if err := c.ShouldBindJSON(&in); err != nil {
        respond.Error(c, 422, "VALIDATION_ERROR", err)
        return
    }
    user, err := h.svc.Create(c.Request.Context(), in)
    if err != nil {
        respond.Error(c, 500, "INTERNAL_ERROR", err)
        return
    }
    respond.Created(c, user, "User created")
}`}
      />
      <p>
        <code>respond</code> is a tiny package that shapes the envelope you
        learned in the last lesson.
      </p>

      <h2>React side — error boundary + toast</h2>
      <CodeBlock
        language="tsx"
        filename="apps/web/components/use-create-user.ts"
        code={`export function useCreateUser() {
  return useMutation({
    mutationFn: api.users.create,
    onSuccess: () => toast.success("User created"),
    onError: (err: ApiError) => {
      if (err.code === "VALIDATION_ERROR") {
        // form-level handling — see details.email etc.
        return
      }
      toast.error(err.message || "Something went wrong")
    },
  })
}`}
      />
      <p>
        Above the page, an <code>&lt;ErrorBoundary&gt;</code> catches anything
        the mutations / queries don&apos;t — including render-time crashes.
        Users see a friendly fallback, not a white screen.
      </p>

      <TipBox tone="warning">
        <strong>Never log and forget.</strong> If you catch an error and decide
        to continue, that&apos;s a deliberate choice — write a comment explaining
        why. Most errors in services should bubble up to the handler.
      </TipBox>

      <KnowledgeCheck
        question="A teammate writes: `db.Save(&user)` without checking the returned error. What's the right code review feedback?"
        choices={[
          {
            label: 'Fine — Save rarely fails in practice.',
            feedback:
              "Wrong — connection issues, constraint violations, and replication lag all make Save fail. Silent failures here cause 'why didn't my update stick?' bugs.",
          },
          {
            label: 'Wrap with `if err := db.Save(&user).Error; err != nil { return fmt.Errorf("saving user: %w", err) }`',
            correct: true,
            feedback:
              "Right — every error is checked, wrapped with context, returned up. The %w lets callers unwrap with errors.Is.",
          },
          {
            label: 'Log it with log.Printf and continue.',
            feedback:
              "Wrong — if the save failed, the caller's mental model is broken. You can't 'continue' as if it succeeded; you have to return an error.",
          },
          {
            label: 'Use a defer + recover to catch panics.',
            feedback:
              "Wrong — db.Save returns errors, not panics. defer/recover is for genuinely unexpected runtime failures, not normal error paths.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Trigger every error path in a real request. With the API
              running, use curl or the OpenAPI docs to hit:
            </p>
            <ol>
              <li>
                <code>POST /api/auth/login</code> with no body — capture the
                422 response in <code>notes.md</code>.
              </li>
              <li>
                <code>POST /api/auth/login</code> with valid email + wrong
                password — capture the 401.
              </li>
              <li>
                <code>GET /api/users/00000000-0000-0000-0000-000000000000</code>{' '}
                — capture the 404.
              </li>
            </ol>
            <p>You should see Shape 3 in all three cases.</p>
          </>
        }
        hint={
          <>
            For 2 + 3, you need a valid auth token. For 1, you don&apos;t —
            422 fires before auth.
          </>
        }
        solution={
          <>
            <p>Each response follows Shape 3 with a different code:</p>
            <ol>
              <li>
                <code>{`{ "error": { "code": "VALIDATION_ERROR", "message": "...", "details": { "email": "required", "password": "required" } } }`}</code>
              </li>
              <li>
                <code>{`{ "error": { "code": "INVALID_CREDENTIALS", "message": "invalid credentials" } }`}</code>
              </li>
              <li>
                <code>{`{ "error": { "code": "NOT_FOUND", "message": "user not found" } }`}</code>
              </li>
            </ol>
            <p>
              That&apos;s the convention surface end-to-end — same shape for every
              error, code + message + optional details.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 4 — <strong>Code generation</strong>. The most fun command
        in Grit. <code>grit generate resource Product</code> writes 7 files
        and wires them all together. We tour them line by line.
      </p>
    </>
  )
}
