import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Go and GORM make SQL injection hard to write by accident. They do not
        make it impossible, and the places it still gets in are specific and
        worth memorising, because they are the places where a parameter is
        not allowed and people reach for string concatenation instead.
      </p>

      <h2>The ordinary path is already safe</h2>
      <p>
        Every value you pass through a placeholder is sent to the database
        separately from the statement. The driver never builds one string out
        of both, so there is nothing for a quote to break out of:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/services/post.go"
        code={`// Safe. The ? is a bind parameter; the value never touches the SQL text.
db.Where("title = ?", userInput).Find(&posts)

// Safe for the same reason.
db.Where(&models.Post{Title: userInput}).Find(&posts)

// Also safe: Raw with placeholders is still parameterised.
db.Raw("SELECT * FROM posts WHERE title = ?", userInput).Scan(&posts)`}
      />

      <p>
        So <code>Raw</code> is not the villain. Concatenation is. The two get
        conflated, and the result is teams that ban <code>Raw</code> and then
        build strings somewhere else.
      </p>

      <CodeBlock
        language="go"
        filename="DELIBERATELY BROKEN"
        code={`// Injectable. The value is part of the statement.
db.Where("title = '" + userInput + "'").Find(&posts)

// Injectable, and harder to spot in review because it looks formatted.
db.Raw(fmt.Sprintf("SELECT * FROM posts WHERE title = '%s'", userInput)).Scan(&posts)`}
      />

      <TipBox tone="danger">
        <p>
          <code>fmt.Sprintf</code> next to SQL is the single highest-signal
          thing to grep for in a Go codebase. It is almost never the right
          way to put a <em>value</em> into a query.
        </p>
      </TipBox>

      <h2>Where it actually gets in: identifiers</h2>
      <p>
        A placeholder works for values. It does not work for column names,
        table names, or the direction of a sort, because those are part of
        the statement&apos;s structure rather than its data. The database
        will not accept <code>ORDER BY ?</code>.
      </p>

      <p>
        Which means any feature that lets the caller choose a sort column or
        a filter field has to build SQL from input. That is the real
        injection surface in a modern Go app, and it is in almost every list
        endpoint ever written:
      </p>

      <CodeBlock
        language="go"
        filename="DELIBERATELY BROKEN"
        code={`sortBy := c.Query("sort")           // "created_at" ... or anything
db.Order(sortBy).Find(&posts)       // and now they choose the SQL`}
      />

      <p>
        The fix is not escaping. It is an allow list, because you know the
        complete set of sortable columns and the caller does not get to
        invent one:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/post.go"
        code={`var sortable = map[string]string{
    "created": "created_at",
    "title":   "title",
    "views":   "view_count",
}

column, ok := sortable[c.Query("sort")]
if !ok {
    column = "created_at"
}

dir := "ASC"
if strings.EqualFold(c.Query("dir"), "desc") {
    dir = "DESC"
}

db.Order(column + " " + dir).Find(&posts)`}
      />

      <p>
        Note what makes this safe: the strings that reach the query are
        <em>yours</em>. The caller&apos;s input is only ever used to pick one.
        An allow list fails closed, and it also gives you a stable public
        name for a column you may want to rename later.
      </p>

      <TipBox tone="warning">
        <p>
          Grit&apos;s generated list endpoints already do this, with the
          sortable set derived from the resource definition. The risk is in
          the endpoint you hand-write next to them, where the habit of the
          generated code does not automatically follow.
        </p>
      </TipBox>

      <h2>The other three</h2>
      <ul>
        <li>
          <strong>LIKE.</strong> <code>Where(&quot;title LIKE ?&quot;,
          &quot;%&quot;+q+&quot;%&quot;)</code> is parameterised and safe
          from injection, but <code>%</code> and <code>_</code> in the input
          are still wildcards. A search for <code>%</code> matches
          everything, which is a performance problem rather than a breach.
        </li>
        <li>
          <strong>IN with a slice.</strong>{' '}
          <code>Where(&quot;id IN ?&quot;, ids)</code> is safe. Building the
          list by joining strings is not.
        </li>
        <li>
          <strong>Migrations and admin tooling.</strong> The code that
          genuinely needs dynamic identifiers. It usually runs with the
          widest database privileges in the system, so it is the worst place
          to take a shortcut.
        </li>
      </ul>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Add a sort parameter to a list endpoint
              the unsafe way, passing <code>c.Query(&quot;sort&quot;)</code>{' '}
              straight to <code>Order</code>.
            </p>
            <p>
              <strong>Two.</strong> Call it with{' '}
              <code>?sort=(CASE WHEN (SELECT COUNT(*) FROM users) &gt; 0 THEN
              title ELSE id END)</code>. It is a legal expression, so it
              sorts. Confirm it does, and note that you just ran a subquery
              the endpoint never offered.
            </p>
            <p>
              <strong>Three.</strong> Replace it with an allow list. Try the
              same input and confirm it falls back to the default.
            </p>
          </>
        }
        hint={
          <p>
            Step two is a demonstration, not an exploit: if an arbitrary
            expression runs at all, the boundary is already gone, and
            whether it can be weaponised depends only on the database and the
            driver.
          </p>
        }
        solution={
          <>
            <p>
              The unsafe version sorts without complaint, which is the whole
              finding: the caller supplied a SQL expression and the database
              evaluated it. From there the usual next step is a boolean
              oracle, flipping the sort on a condition to read data one bit
              at a time.
            </p>
            <p>
              With the allow list the lookup misses, the default applies, and
              there is nothing to probe. Escaping would have been the wrong
              fix: you cannot escape an identifier into safety, and you do
              not need to when you already know the valid set.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="Which of these is injectable?"
        choices={[
          {
            label: 'db.Raw("SELECT * FROM posts WHERE id = ?", id).Scan(&p)',
            feedback:
              'Raw with a placeholder is parameterised like any other query. Raw is not the problem; concatenation is.',
          },
          {
            label: 'db.Order(c.Query("sort")).Find(&posts)',
            correct: true,
            feedback:
              'A placeholder cannot be used for an identifier, so this has to build SQL from input. Use an allow list: the caller picks from a set you define, and the strings in the query are always yours.',
          },
          {
            label: 'db.Where("title LIKE ?", "%"+q+"%").Find(&posts)',
            feedback:
              'Parameterised and safe from injection. The % and _ in q are still wildcards, which is a performance consideration rather than a vulnerability.',
          },
          {
            label: 'db.Where(&models.Post{Title: in.Title}).Find(&posts)',
            feedback:
              'Struct conditions are parameterised by GORM.',
          },
        ]}
      />
    </>
  )
}
