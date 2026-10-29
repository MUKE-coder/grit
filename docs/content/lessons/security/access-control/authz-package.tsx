import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The last lesson fixed one IDOR by hand. Doing that in every handler
        is how the second one gets forgotten, so Grit puts the answer in one
        package: <code>internal/authz</code>. This lesson is a tour of what
        is in it and, more usefully, which of the four questions each part
        answers.
      </p>

      <h2>Authorization is four different questions</h2>
      <p>
        Most arguments about access control are really people answering
        different questions with the same word. Separate them and the package
        stops looking like a grab bag:
      </p>

      <ul>
        <li>
          <strong>May this user touch widgets at all?</strong> That is a{' '}
          <em>permission</em>, it comes from the route&apos;s middleware, and
          it never looks at a row.
        </li>
        <li>
          <strong>Is this their row?</strong> That is <em>ownership</em>, and
          it is a column.
        </li>
        <li>
          <strong>Is this row in a state where that is allowed?</strong> A
          draft may be edited and a published post may not. That is a{' '}
          <em>policy</em>, and it needs the row and the user together.
        </li>
        <li>
          <strong>Can this admin grant a permission they do not hold?</strong>{' '}
          That is the <em>ceiling</em>, and forgetting it is how a junior
          admin promotes themselves.
        </li>
      </ul>

      <TipBox tone="info">
        <p>
          A rule that cannot be expressed as one of those four is usually two
          rules wearing a coat. Split it before you write it.
        </p>
      </TipBox>

      <h2>Ownership: MustOwn and the 404</h2>
      <p>
        The ownership helpers load the row and check the owner in one call.
        Note the status code:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/note.go"
        code={`func (h *NoteHandler) Update(c *gin.Context) {
    var note models.Note
    // Loads by id AND checks the owner. On any failure it has already
    // written the response, so the handler just returns.
    if err := authz.MustOwn(c, h.DB, &note, c.Param("id")); err != nil {
        return
    }

    // From here the row is theirs. Ordinary handler code.
}`}
      />

      <p>
        <code>MustOwn</code> answers a failure with <strong>404, not 403</strong>.
        That is deliberate, and it is the part people change back without
        thinking. A 403 says &quot;this row exists and is not yours&quot;,
        which is a yes/no oracle: an attacker walks the id space and learns
        exactly which records exist before they ever break in. A 404 says
        nothing. Same code, same user experience for the legitimate caller,
        and no free map of your database.
      </p>

      <p>
        The variants exist because admin screens call the same endpoints:
      </p>

      <ul>
        <li>
          <code>MustOwn</code>, strict ownership.
        </li>
        <li>
          <code>MustOwnUnlessAdmin</code>, the same check with an ADMIN
          bypass. Without it an administrator opening the admin panel sees
          only the rows they personally created, which reads as a bug and
          gets fixed by removing the check entirely.
        </li>
        <li>
          <code>ScopeOwned(ctx, query, column)</code>, for lists. It adds the
          owner filter to a query rather than checking one row.
        </li>
      </ul>

      <TipBox tone="warning">
        <p>
          <code>ScopeOwned</code> returns <code>q.Where(...)</code>, a
          statement clause. That matters: a clause is carried into{' '}
          <code>Find</code>, <code>Scan</code> and <code>Rows</code> alike.
          Scoping implemented as a GORM callback instead is easy to get wrong,
          because <code>Scan</code> does not go through the query callback.
        </p>
      </TipBox>

      <h2>Permissions: the catalogue, not a string</h2>
      <p>
        Permissions are generated into a catalogue rather than typed as
        strings at each call site, because a typo in a permission string
        fails open: <code>Can(db, user, &quot;psots.delete&quot;)</code> is
        simply false for everyone, including the people who should be denied
        loudly.
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/authz"
        code={`authz.Catalog()                      // every module and its permissions
authz.Keys()                         // every permission key, flat
authz.Can(db, userID, "posts.delete")// does this user hold it
authz.Granted(grants, "posts.delete")// against a list you already loaded
authz.Expand(grants)                 // wildcards resolved to concrete keys`}
      />

      <h2>Policies: the rule that needs the row</h2>
      <p>
        Neither a permission nor an owner column can say &quot;a post may be
        edited while it is a draft and not after it is published&quot;. That
        is a policy: a named ability with a rule that receives the actor and
        the record.
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/authz/policies.go"
        code={`authz.Define("post.edit", func(ctx context.Context, record any) authz.Decision {
    post, ok := record.(*models.Post)
    if !ok {
        return authz.Deny("not a post")
    }
    if post.Status == "published" {
        return authz.Deny("a published post is not editable")
    }
    if !authz.Owns(ctx, post) {
        return authz.Deny("not yours")
    }
    return authz.Allow()
})`}
      />

      <p>
        Then <code>authz.Allows(ctx, &quot;post.edit&quot;, &amp;post)</code>{' '}
        at the call site, or <code>Inspect</code> when you want the reason to
        put in the response. The reason is worth having: &quot;forbidden&quot;
        teaches a user nothing, and &quot;a published post is not
        editable&quot; stops a support ticket.
      </p>

      <h2>The ceiling: the one people forget</h2>
      <p>
        An admin screen that lets somebody assign permissions must stop them
        assigning permissions they do not themselves hold. Otherwise the
        lowest-privileged admin account is a route to the highest, in two
        clicks and with a full audit trail showing it was allowed.
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/role.go"
        code={`// Everything in requested that the caller does not hold.
if over := authz.CannotGrant(c, requested); len(over) > 0 {
    respond.Fail(c, respond.CodeForbidden,
        "You cannot grant: "+strings.Join(over, ", "))
    return
}`}
      />

      <TipBox tone="danger">
        <p>
          Privilege escalation almost never looks like an attack in the logs.
          It looks like an administrator using an administration screen
          exactly as designed, which is why the check has to be in the code
          rather than in the reviewer&apos;s head.
        </p>
      </TipBox>

      <h2>Which to reach for</h2>
      <CodeBlock
        language="text"
        code={`Can they use this endpoint at all?   permission, in the route middleware
Is this their row?                   MustOwn / ScopeOwned
Is the row in a state that allows    a policy: authz.Define + Allows
  it, or does the rule involve both
  the user and the record?
Are they handing out power they       authz.CannotGrant
  do not have?`}
      />

      <Exercise
        prompt={
          <>
            <p>
              Take one endpoint in your project that lists rows and one that
              updates a single row.
            </p>
            <p>
              <strong>One.</strong> Put <code>ScopeOwned</code> on the list and{' '}
              <code>MustOwn</code> on the update.
            </p>
            <p>
              <strong>Two.</strong> Sign in as a second user and request the
              first user&apos;s row by id. Write down the status code you get.
            </p>
            <p>
              <strong>Three.</strong> Change <code>MustOwn</code> to return 403
              instead, repeat, and explain in one sentence what an attacker
              now learns that they did not before.
            </p>
          </>
        }
        hint={
          <p>
            For step three, think about what you can discover by trying ids
            you have no access to, without ever succeeding at anything.
          </p>
        }
        solution={
          <>
            <p>
              Step two gives <strong>404</strong>. Step three gives{' '}
              <strong>403</strong>, and the difference is an existence oracle:
              403 means the row is there, 404 means it may or may not be.
              Walking the id space against a 403 produces a list of every
              record in the table, which is useful on its own and is the
              reconnaissance for the next attack.
            </p>
            <p>Put it back to 404.</p>
          </>
        }
      />

      <KnowledgeCheck
        question="An admin panel calls the same endpoints as the app. You add MustOwn to a handler and administrators start seeing only the rows they created themselves. What is the fix?"
        choices={[
          {
            label: 'Remove MustOwn from the handler',
            feedback:
              'This fixes the symptom by removing the security. The endpoint is now open to every signed-in user, which is the bug MustOwn was added for.',
          },
          {
            label: 'Use MustOwnUnlessAdmin, the same check with an ADMIN bypass',
            correct: true,
            feedback:
              'One call, the check intact, and the admin panel works. This variant exists precisely because the panel calls the same endpoints as the app.',
          },
          {
            label: 'Give every administrator ownership of every row',
            feedback:
              'Ownership would stop meaning anything, and the owner column is what your audit trail and your list scoping both read.',
          },
          {
            label: 'Add a second set of endpoints for the admin panel',
            feedback:
              'Twice the surface to secure, and the two copies drift. The second one is where next year\'s missing check lives.',
          },
        ]}
      />
    </>
  )
}
