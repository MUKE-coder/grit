import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Before we run the generator, you need to know what it generates.
        Spoiler: it&apos;s a <strong>resource</strong>. This 5-minute lesson
        defines the word.
      </p>

      <h2>A resource is a complete vertical slice</h2>
      <p>
        For each domain concept (User, Order, Invoice, Product), Grit
        produces a slice that spans the whole stack:
      </p>

      <CodeBlock
        language="text"
        code={`Model           GORM struct        apps/api/internal/models/order.go
Service         Business logic     apps/api/internal/services/order.go
Handler         HTTP layer         apps/api/internal/handlers/order.go
Route           Wiring             apps/api/internal/routes/routes.go (injected)
Zod schema      Validation         packages/shared/src/schemas/order.ts
TS type         Frontend type      packages/shared/src/types/order.ts
React Query hook  Data fetching    apps/web/hooks/use-orders.ts
Admin resource  Filament-style page  apps/admin/app/resources/orders/`}
      />

      <p>
        That&apos;s eight artefacts for one concept — too many to write by hand
        every time you add an entity. The generator writes all of them in
        one command.
      </p>

      <h2>The mental model</h2>
      <p>A resource exists to do CRUD on a thing:</p>
      <ul>
        <li><strong>Create</strong> — POST /api/orders</li>
        <li><strong>Read one</strong> — GET /api/orders/:id</li>
        <li><strong>Read many</strong> — GET /api/orders (with pagination)</li>
        <li><strong>Update</strong> — PUT /api/orders/:id</li>
        <li><strong>Delete</strong> — DELETE /api/orders/:id</li>
      </ul>
      <p>
        Five endpoints, one model, one admin page. The handler, service,
        type, schema, hook, and admin page all line up to make that flow
        work end-to-end.
      </p>

      <TipBox tone="info">
        Not every domain concept needs to be a resource. <em>Notification</em>{' '}
        is a service (you don&apos;t CRUD notifications — you send them).{' '}
        <em>Order</em> is a resource (you create, list, update, cancel them).
        If you&apos;d expose all 5 verbs to a UI, it&apos;s a resource.
      </TipBox>

      <h2>Fields define the shape</h2>
      <p>
        When you generate a resource, you tell the generator what fields the
        model has:
      </p>

      <CodeBlock
        terminal
        code={`grit generate resource Order \\
  --field "customerId:uuid:required" \\
  --field "total:decimal:required" \\
  --field "status:string:default=pending"`}
      />
      <p>
        Each field becomes: a GORM column, a Zod validator, a TS type
        property, an admin form input, a DataTable column. The generator
        knows the mapping for every supported type.
      </p>

      <KnowledgeCheck
        question="Which of these is most likely NOT a Grit resource?"
        choices={[
          {
            label: 'BlogPost — readers see them, authors create + edit them, admins delete spam',
            feedback:
              "Resource — all 5 CRUD verbs are needed (Create by authors, Read by anyone, Update by authors, Delete by admins, List for the index).",
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
              "Resource — full CRUD with role gates (staff sees all; customers see their own).",
          },
          {
            label: 'Invoice — generated from an order; viewed by customer and staff',
            feedback:
              "Resource — Read (customer + staff), Create (manual + auto from Order), maybe Update for adjustments. Most apps treat invoices as a resource.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In your <code>notes.md</code>, list 5 concepts from a real
              product you know (your job&apos;s app, a side project, anything).
              For each, write &quot;resource&quot; or &quot;not a resource&quot; and one
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
              <li>Payment — partial resource (created by Stripe webhook; read by staff)</li>
              <li>Reminder email — not a resource (sent by cron job, no UI CRUD)</li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You know what a resource is. Time to make one — next lesson we run{' '}
        <code>grit generate resource Product</code> end-to-end.
      </p>
    </>
  )
}
