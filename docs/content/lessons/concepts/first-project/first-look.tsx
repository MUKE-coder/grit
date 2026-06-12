import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Everything is running. This is the &quot;see what you have&quot; lesson — open
        each URL, recognise what it is, log in to the admin panel, and read
        the JSON the API returns.
      </p>

      <h2>The five URLs</h2>
      <CodeBlock
        language="text"
        code={`http://localhost:3000     Web (public site)
http://localhost:3001     Admin panel
http://localhost:8080     API root
http://localhost:8080/api/health   Health check (returns JSON)
http://localhost:8080/docs         OpenAPI docs (auto-generated)
http://localhost:8080/studio       GORM Studio — visual DB browser
http://localhost:8025              Mailhog UI — captures sent emails`}
      />

      <h2>1. The API — the smallest possible win</h2>
      <p>
        Open <code>http://localhost:8080/api/health</code>. You should see:
      </p>
      <CodeBlock
        language="json"
        code={`{ "status": "ok", "version": "0.1.0" }`}
      />
      <p>That&apos;s the API alive and responding. Tiny but important.</p>

      <h2>2. The OpenAPI docs</h2>
      <p>
        Open <code>http://localhost:8080/docs</code>. You see Scalar UI — an
        auto-generated, interactive docs surface for every endpoint. Click
        any endpoint to see the schema, try it out, see the response. No
        OpenAPI YAML to maintain — Grit generates it from the Go code.
      </p>

      <h2>3. The admin panel</h2>
      <p>
        Open <code>http://localhost:3001</code>. Log in with:
      </p>
      <CodeBlock
        language="text"
        code={`Email:    admin@example.com
Password: admin123`}
      />
      <p>
        You land on the admin dashboard. Click <strong>Users</strong> in the
        sidebar — you&apos;ll see a table with sortable columns, a filter bar,
        and pagination. That&apos;s the <code>DataTable</code> primitive doing
        the work. Click a user to edit them — that&apos;s the <code>FormBuilder</code>.
      </p>

      <TipBox tone="info">
        <strong>The admin panel writes itself.</strong> When you generate a
        new resource (chapter 4), an admin page for it appears automatically.
        No HTML, no form code — just <code>defineResource()</code>.
      </TipBox>

      <h2>4. The web app</h2>
      <p>
        Open <code>http://localhost:3000</code>. The public Next.js site —
        scaffolded with a landing page, signup, login, and a basic
        dashboard. We won&apos;t touch this much in the concepts course; the{' '}
        <em>Building Web with Next.js</em> course goes deep here.
      </p>

      <h2>5. GORM Studio — see your data</h2>
      <p>
        <code>http://localhost:8080/studio</code> opens GORM Studio — a
        visual DB browser. Click <strong>users</strong> in the left rail.
        You see the seeded admin row, schema, indexes. Run ad-hoc SQL from
        the query tab. Useful when you&apos;re debugging &quot;is this row really in
        the DB?&quot;.
      </p>

      <KnowledgeCheck
        question="You sign up a new user from the web app. Which URL should you check to see what email the system tried to send them?"
        choices={[
          {
            label: 'http://localhost:8080/docs',
            feedback:
              "Wrong — that's the OpenAPI docs surface. It documents endpoints, doesn't capture emails.",
          },
          {
            label: 'http://localhost:8025',
            correct: true,
            feedback:
              "Right — Mailhog runs at :8025 and captures every email the API tries to send in dev. Click any captured email to see the rendered HTML.",
          },
          {
            label: 'http://localhost:8080/studio',
            feedback:
              "Wrong — that's the DB browser. You'd see the user row, not the email.",
          },
          {
            label: 'Your real inbox',
            feedback:
              "Wrong — Grit's docker-compose dev setup routes mail to Mailhog, never to a real SMTP. Production uses Resend; dev captures everything locally so you can't accidentally email real users.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Take a tour. Open each URL and capture <strong>three
              screenshots</strong> for your <code>notes.md</code>:
            </p>
            <ol>
              <li>The JSON response from <code>/api/health</code></li>
              <li>The admin dashboard after you log in</li>
              <li>The users table in GORM Studio</li>
            </ol>
            <p>
              That&apos;s the chapter 2 assignment done. Open the assignment
              page in the sidebar to tick off the success criteria.
            </p>
          </>
        }
        hint={
          <>
            If the admin login fails with &quot;invalid credentials&quot;, you
            probably skipped <code>grit seed</code>. Run it, then try again.
          </>
        }
        solution={
          <>
            <p>
              No exact solution — the screenshots are the proof. If you can
              see live API JSON + the admin panel logged in + the users
              table in Studio, chapter 2 is done. You now have a working
              Grit project in your browser.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 3 — <strong>The Convention Surface</strong>. Now that you can
        see the project running, we get into the rules every file follows
        and why. Once you internalise these, you can write code that fits
        the codebase before you&apos;ve even read it.
      </p>
    </>
  )
}
