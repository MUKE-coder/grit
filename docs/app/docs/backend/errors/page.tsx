import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'
import { errorCodeAreas, errorCodeCount } from '@/components/error-codes'

export const metadata = getDocMetadata('/docs/backend/errors')

// The status colours group errors the way a client does: what the caller can fix,
// what needs a sign-in, what is worth retrying.
const statusTone: Record<string, string> = {
  request: 'text-amber-500/90 border-amber-500/20 bg-amber-500/5',
  auth: 'text-sky-500/90 border-sky-500/20 bg-sky-500/5',
  permission: 'text-sky-500/90 border-sky-500/20 bg-sky-500/5',
  notfound: 'text-muted-foreground border-border/40 bg-accent/20',
  conflict: 'text-violet-500/90 border-violet-500/20 bg-violet-500/5',
  state: 'text-violet-500/90 border-violet-500/20 bg-violet-500/5',
  limit: 'text-orange-500/90 border-orange-500/20 bg-orange-500/5',
  server: 'text-red-500/90 border-red-500/20 bg-red-500/5',
  upstream: 'text-red-500/90 border-red-500/20 bg-red-500/5',
  disabled: 'text-muted-foreground border-border/40 bg-accent/20',
}

export default function ErrorCodesPage() {
  return (
    <DocsShell>
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Error codes</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every error this API can return, what it means, and what a client should do about
                it. {errorCodeCount} codes, and each one always arrives with the same status.
              </p>
            </div>

            <div className="prose-grit">
              <p>
                A client needs two things from an error: a code it can branch on, and the
                certainty that the code comes with the same status every time. Grit had the first
                and not the second. <code>VALIDATION_ERROR</code> came back as 422 from
                thirty-eight handlers and 400 from twenty-five, <code>INVALID_TOKEN</code> as both
                401 and 400, and no page listed the codes at all, so the only way to learn one was
                to trigger it.
              </p>
              <p>
                This page is generated from the same catalogue that generates{' '}
                <code>internal/respond/codes.go</code> in your project and{' '}
                <code>packages/shared/types/errors.ts</code> beside it. They cannot disagree, and a
                test in the CLI refuses any code a handler returns that is not listed here, or that
                is returned with a status other than the one shown.
              </p>

              <h2 id="shape">The shape</h2>
              <p>Every error is the same envelope, whatever went wrong:</p>
              <CodeBlock
                language="json"
                filename="422 Unprocessable Entity"
                code={`{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Check the highlighted fields",
    "details": {
      "email": "That address is already in use",
      "price": "must be greater than 0"
    }
  }
}`}
              />
              <p>
                <code>details</code> is present on a validation failure and maps each field to what
                is wrong with it. Branch on <code>code</code> rather than on the status: the status
                groups errors, the code says which one it is.
              </p>

              <h2 id="typescript">In TypeScript</h2>
              <p>
                The union makes the switch exhaustive, so a client that handles five codes and
                forgets the sixth fails to compile instead of falling through to &quot;Something
                went wrong&quot;.
              </p>
              <CodeBlock
                language="typescript"
                code={`import { API_ERRORS, documentedErrorCode, isRetryable } from '@shared/types'

try {
  await api.post('/api/v1/invoices', body)
} catch (err) {
  const code = documentedErrorCode(err)
  if (!code) throw err                        // not one of ours: a proxy, or your own code
  if (isRetryable(code)) return retryLater()  // server, upstream, rate limit, not configured
  toast(API_ERRORS[code].client)              // what the person should do about it
}`}
              />

              <h2 id="go">In Go</h2>
              <p>
                Handlers return a code, and the status comes from the catalogue rather than being
                typed next to it. That is what keeps one code on one status.
              </p>
              <CodeBlock
                language="go"
                code={`// 422, because that is what the catalogue gives VALIDATION_ERROR
respond.Fail(c, respond.CodeValidationError, "Check the highlighted fields",
    map[string]string{"email": "That address is already in use"})

// 403, with a message written for the person reading it
respond.Fail(c, respond.CodeForbidden, "Only an owner can archive an invoice")`}
              />
              <p>
                A rule your own code breaks does not need a new code:{' '}
                <code>respond.Rule(&quot;debits and credits do not balance&quot;)</code> returned
                from a service or a GORM hook becomes a 422 with that sentence. Add a code to the
                catalogue only when a client should branch on it.
              </p>

              <h2 id="catalogue">The catalogue</h2>
              <p>
                Grouped by the part of the API that raises them. The first group is the one every
                endpoint can answer with.
              </p>
            </div>

            {errorCodeAreas.map((area) => (
              <div key={area.area} className="mb-10">
                <h3 id={area.area} className="text-lg font-semibold tracking-tight mb-3 mt-8">
                  {area.label}
                </h3>
                <div className="space-y-3">
                  {area.codes.map((row) => (
                    <div
                      key={row.code}
                      className="rounded-lg border border-border/40 bg-card/30 p-4"
                    >
                      <div className="flex items-center gap-2.5 mb-2 flex-wrap">
                        <code className="text-[13px] font-mono font-semibold text-foreground">
                          {row.code}
                        </code>
                        <span
                          className={`inline-flex items-center rounded-md border px-1.5 py-0.5 text-[11px] font-mono ${
                            statusTone[row.category] ?? statusTone.notfound
                          }`}
                        >
                          {row.status}
                        </span>
                        <span className="text-[11px] font-mono text-muted-foreground/70">
                          {row.category}
                        </span>
                      </div>
                      <p className="text-[13.5px] text-muted-foreground leading-relaxed mb-1.5">
                        {row.meaning}
                      </p>
                      <p className="text-[13.5px] text-foreground/80 leading-relaxed">
                        <span className="text-muted-foreground/70">What to do: </span>
                        {row.client}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            ))}

            <div className="prose-grit">
              <h2 id="your-own">Your own codes</h2>
              <p>
                Grit&apos;s catalogue cannot hold every code a product needs. A plan limit, a
                thesis that has already been submitted, a licence that does not cover this
                region: those are domain facts, and the nearest built-in codes are wrong in
                ways a client can feel. <code>RATE_LIMITED</code> says &ldquo;retry
                shortly&rdquo;, so a client with backoff retries against a wall.{' '}
                <code>FORBIDDEN</code> is indistinguishable from a permission failure, so the UI
                can only tell them apart by matching on the message string.
              </p>
              <p>
                Declare your own in <code>.grit/error-codes.yaml</code>. Not in{' '}
                <code>internal/respond/codes.go</code>: that file is generated and an entry added
                to it is lost on the next upgrade.
              </p>
              <CodeBlock
                language="yaml"
                code={`# .grit/error-codes.yaml
codes:
  - code: PLAN_LIMIT_REACHED          # SHOUTING_SNAKE_CASE
    status: 402                       # a 4xx or a 5xx
    category: limit
    area: billing                     # groups it in the table above
    meaning: The account has used all three of its similarity checks.
    client: Show the upgrade prompt. Retrying will not help.`}
              />
              <CodeBlock language="bash" code={`grit sync`} />
              <p>
                That folds them into the same two generated files the built-in codes come from,
                which is the property worth having: the typed Go constant and the{' '}
                <code>ApiErrorCode</code> union come from one list, so a switch over error codes
                in the frontend stays exhaustive for the errors your product actually raises. A
                Go-only extension would give the API its code and leave the frontend unable to
                type a case for it.
              </p>
              <CodeBlock
                language="go"
                code={`// apps/api/internal/handlers/check.go
if used >= plan.Checks {
    respond.Fail(c, respond.CodePlanLimitReached,
        "You have used all 3 checks on the Student plan",
        map[string]string{"used": "3", "limit": "3", "plan": "student"})
    return
}`}
              />
              <p>
                <code>PLAN_LIMIT_REACHED</code> is in the catalogue already, at 402, because
                nearly every SaaS needs it and 402 is what Stripe and GitHub answer for the same
                situation. A code of your own gets a constant the same way:{' '}
                <code>THESIS_ALREADY_SUBMITTED</code> becomes{' '}
                <code>respond.CodeThesisAlreadySubmitted</code>.
              </p>

              <h3 className="mt-8">A code that is not declared</h3>
              <p>
                <code>respond.Code</code> is a string type, so this compiles:
              </p>
              <CodeBlock
                language="go"
                code={`respond.Fail(c, respond.Code("PLAN_LIMIT_REACHED"), "...")   // a bug`}
              />
              <p>
                It is still a bug. <code>StatusOf</code> does not know the code, so the caller is
                answered <strong>500</strong>: a client error reported as a server fault, which
                is the most expensive kind to find. Two things catch it. The generated{' '}
                <code>internal/respond/codes_project_test.go</code> fails on any{' '}
                <code>respond.Code(&quot;...&quot;)</code> in your project that is not declared,
                so <code>go test ./...</code> finds it before a deploy does. And if one reaches
                production anyway, <code>respond.Fail</code> logs it once by name with what to do
                about it, rather than answering 500 in silence.
              </p>
              <p>
                Two helpers to know when you branch on codes in TypeScript:{' '}
                <code>documentedErrorCode()</code> returns null for a code it does not know,
                while <code>apiErrorCode()</code> from <code>@shared/types</code> returns the raw
                string. After <code>grit sync</code> your codes are documented, so the first one
                works for them too.
              </p>
            </div>

            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/response-format" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  API Response Format
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/request-lifecycle" className="gap-1.5">
                  The Request Lifecycle
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
    </DocsShell>
  )
}
