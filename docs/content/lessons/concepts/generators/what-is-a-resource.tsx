import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Before we run the generator, you need a word for what it
        generates. Spoiler: it&apos;s a <strong>resource</strong>. The
        rest of this chapter goes deep on the generator — this 5-minute
        lesson nails down the concept so the next four lessons land
        properly.
      </p>

      <h2>A resource is a complete vertical slice</h2>
      <p>
        For each domain concept in your app (User, Contact, Order,
        Product, Invoice…), Grit produces a slice that spans the whole
        stack. Not just &quot;a model&quot; — every layer that touches
        that concept gets a file:
      </p>

      <CodeBlock
        language="text"
        code={`Model            GORM struct          apps/api/internal/models/contact.go
Service          Business logic       apps/api/internal/services/contact.go
Handler          HTTP layer           apps/api/internal/handlers/contact.go
Routes           Wiring               apps/api/internal/routes/routes.go (injected)
Zod schema       Validation           packages/shared/src/schemas/contact.ts
TypeScript type  Frontend type        packages/shared/src/types/contact.ts
React Query hook Data fetching        apps/web/hooks/use-contacts.ts
Admin resource   Filament-style page  apps/admin/app/resources/contacts/page.tsx`}
      />

      <p>
        That&apos;s <strong>eight artefacts</strong> for one concept —
        too many to write by hand every time you add an entity, and too
        easy to drift out of sync if you do (a field on the Go struct
        that&apos;s missing from the TS type; a route the admin page
        forgot to call). The generator writes all eight in one command
        and keeps them aligned.
      </p>

      <h2>The mental model: CRUD plus a UI</h2>
      <p>A resource exists to do CRUD on a thing:</p>
      <ul>
        <li><strong>Create</strong> — <code>POST /api/contacts</code></li>
        <li><strong>Read one</strong> — <code>GET /api/contacts/:id</code></li>
        <li><strong>Read many</strong> — <code>GET /api/contacts</code> (paginated, searchable)</li>
        <li><strong>Update</strong> — <code>PUT /api/contacts/:id</code></li>
        <li><strong>Delete</strong> — <code>DELETE /api/contacts/:id</code></li>
      </ul>
      <p>
        Five endpoints, one model, one admin page. The handler, service,
        type, schema, hook, and admin page all line up so that flow works
        end-to-end. If a feature in your head expands to all five verbs
        and someone (a user, an admin, support) clicks a button to invoke
        them — it&apos;s probably a resource.
      </p>

      <TipBox tone="info">
        Not every domain concept is a resource. <em>Notification</em> is
        usually a service (you don&apos;t CRUD notifications — you send
        them when something happens). <em>Order</em> is a resource (you
        create, list, update, cancel them). The CRUD-via-UI test is the
        cleanest filter.
      </TipBox>

      <h2>Fields define the shape</h2>
      <p>
        When you generate a resource, you tell the generator what fields
        the model has. One simple example you&apos;ll see again next
        lesson:
      </p>

      <CodeBlock
        terminal
        code={`grit generate resource Contact \\
  --fields "name:string,email:string:unique,phone:string:optional"`}
      />

      <p>
        Each field becomes <em>seven things at once</em>: a GORM column,
        a Go struct field, a Zod validator, a TypeScript property, an
        admin form input, a DataTable column, and a search match (if the
        type is searchable). The generator owns that mapping table — you
        just describe the shape.
      </p>

      <KnowledgeCheck
        question="Which of these is most likely NOT a Grit resource?"
        choices={[
          {
            label: 'BlogPost — readers see them, authors create + edit them, admins delete spam',
            feedback:
              'Resource — all 5 CRUD verbs are needed (Create by authors, Read by anyone, Update by authors, Delete by admins, List for the index).',
          },
          {
            label: 'EmailNotification — fired by the system when an order ships',
            correct: true,
            feedback:
              "Right — you never POST /api/email-notifications from the UI; the system enqueues them. EmailNotification is a service (with a model for the queued state, maybe), not a Grit resource.",
          },
          {
            label: 'Customer — staff manages them, customers update their own profile',
            feedback:
              'Resource — full CRUD with role gates (staff sees all; customers see their own).',
          },
          {
            label: 'Invoice — generated from an order; viewed by customer and staff',
            feedback:
              'Resource — Read (customer + staff), Create (manual + auto from Order), maybe Update for adjustments. Most apps treat invoices as a resource.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In your <code>notes.md</code>, list 5 concepts from a real
              product you know (your job&apos;s app, a side project, an
              app you use daily). For each, write
              &quot;resource&quot; or &quot;not a resource&quot; and one
              sentence why.
            </p>
          </>
        }
        hint={<>The CRUD-via-UI test is the cleanest filter.</>}
        solution={
          <>
            <p>Example for a music-teacher booking app:</p>
            <ul>
              <li>Teacher — resource (admin CRUD, public listing)</li>
              <li>Student — resource (teacher CRUD, student profile edit)</li>
              <li>Lesson — resource (booked, cancelled, completed states)</li>
              <li>Payment — partial resource (created by Stripe webhook; read by staff; rarely manually edited)</li>
              <li>Reminder email — not a resource (sent by cron job; no UI CRUD)</li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You know what a resource is. Next lesson we dissect the command
        itself — the anatomy of{' '}
        <code>grit generate resource Contact --fields &quot;…&quot;</code>{' '}
        — and then we generate Contact end-to-end and watch the eight
        files appear.
      </p>
    </>
  )
}
