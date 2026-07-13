import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Auto-Generated API Docs: Scalar & Swagger in Grit',
  description: 'Learn how Grit auto-generates interactive API documentation using gin-docs. Covers OpenAPI 3.1, Scalar UI, request testing, authentication, response formats, and spec export.',
}

export default function ApiDocsScalarCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Auto-Generated API Docs</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Auto-Generated API Docs: Scalar & Swagger in Grit
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Every API needs documentation. Other developers (and your future self) need to know what
            endpoints exist, what parameters they accept, and what they return. In this course, you will
            learn how Grit auto-generates interactive API documentation using gin-docs — no annotations
            required. You{"'"}ll explore the Scalar UI, test endpoints live, and export the OpenAPI spec.
          </p>
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/backend/api-docs" className="text-primary hover:underline">API Documentation →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is API Documentation? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is API Documentation?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Imagine you inherit a project with 50 API endpoints. No documentation. No comments.
            You have to read every handler, every route, every model to figure out what the API does.
            That{"'"}s a nightmare — and it{"'"}s surprisingly common. API documentation solves this by
            providing a clear, browsable reference of every endpoint your API offers.
          </p>

          <Definition term="API Documentation">
            A structured reference that describes every endpoint in your API: the URL, HTTP method,
            request parameters, request body schema, response format, authentication requirements,
            and example responses. Good documentation lets a developer use your API without reading
            a single line of source code.
          </Definition>

          <Definition term="OpenAPI Specification">
            A standard format (formerly called Swagger Specification) for describing REST APIs.
            It{"'"}s a JSON or YAML file that machines can read and tools can use to generate
            documentation UIs, client SDKs, and test suites. The current version is OpenAPI 3.1.
          </Definition>

          <Definition term="Swagger">
            The original name for the OpenAPI Specification, and also a set of tools built around it.
            When people say {'"'}Swagger,{'"'} they usually mean the interactive documentation UI that
            lets you browse endpoints and send test requests. Swagger was donated to the OpenAPI
            Initiative in 2015, but the name stuck.
          </Definition>

          <Definition term="Scalar">
            A modern, beautiful alternative to Swagger UI for displaying OpenAPI documentation.
            It offers a cleaner design, better search, syntax-highlighted request/response examples,
            and a built-in request builder. Grit uses Scalar as its documentation frontend.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Why does documentation matter? Three reasons:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">For other developers:</strong> Frontend developers, mobile developers, and third-party integrators need to know your API{"'"}s contract without reading Go code</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">For your future self:</strong> You will forget what endpoints exist, what parameters are required, and what the response looks like. Documentation is your memory</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">For testing:</strong> Interactive docs let you test endpoints directly in the browser — no Postman, no curl, no frontend needed</li>
          </ul>

          <Challenge number={1} title="Reflect on Undocumented APIs">
            <p>Have you ever used an API without documentation? How was the experience? Think about
            how you figured out what endpoints existed, what parameters to send, and what the response
            looked like. Write down 3 problems you encountered (or would encounter) when using an
            undocumented API.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: How gin-docs Works ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">How gin-docs Works</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Most API documentation tools require you to write annotations — special comments above
            every handler function describing the endpoint. That{"'"}s tedious and error-prone because
            the annotations can drift out of sync with the actual code. Grit takes a different approach.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit uses <a href="https://github.com/MUKE-coder/gin-docs" target="_blank" rel="noreferrer" className="text-primary hover:underline">gin-docs</a> to
            auto-generate OpenAPI 3.1 specifications. No annotations needed. It works by <strong className="text-foreground">introspecting</strong> your
            Gin routes and GORM models at startup.
          </p>

          <Definition term="Introspection">
            The ability of a program to examine its own structure at runtime. gin-docs reads your
            route registrations (which URL maps to which handler) and your model definitions (which
            fields exist, what types they are) to build the API specification automatically. When you
            add a new route or model field, the documentation updates itself.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s what gin-docs does when your API starts:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span><strong className="text-foreground">Scans your Gin router.</strong> It reads every registered route — the HTTP method, the URL path, the middleware chain, and the handler function.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span><strong className="text-foreground">Reads your GORM models.</strong> It inspects the struct fields, their types, JSON tags, and GORM tags to generate request/response schemas.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">3.</span>
              <span><strong className="text-foreground">Generates the OpenAPI 3.1 spec.</strong> A complete JSON specification describing every endpoint, parameter, request body, and response schema.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">4.</span>
              <span><strong className="text-foreground">Serves the Scalar UI.</strong> The interactive documentation is available at <Code>/docs</Code> in your browser.</span>
            </li>
          </ol>

          <CodeBlock filename="What You Get — Zero Configuration">
{`# Start your Grit API
cd apps/api && go run cmd/server/main.go

# Open your browser
http://localhost:8080/docs

# That's it. No annotations. No swagger comments.
# No yaml files to maintain. Just your routes and models.`}
          </CodeBlock>

          <Tip>
            Every time you generate a new resource with <Code>grit generate resource</Code>, the
            documentation updates automatically on the next server restart. The new endpoints,
            request schemas, and response schemas all appear in <Code>/docs</Code> without any
            additional work.
          </Tip>

          <Challenge number={2} title="Explore Your API Docs">
            <p>Open <Code>http://localhost:8080/docs</Code> in your browser. How many endpoint groups
            do you see? List the group names (e.g., Auth, Users, etc.). How many total endpoints
            are documented?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: The Scalar UI ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Scalar UI</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you open <Code>/docs</Code>, you see the Scalar UI — a modern, interactive documentation
            interface. Let{"'"}s walk through its key features.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The left sidebar shows all your <strong className="text-foreground">endpoint groups</strong>,
            organized by route prefix. Click a group to expand it and see individual endpoints.
            Each endpoint shows:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">HTTP method badge</strong> — color-coded: green for GET, blue for POST, orange for PUT, red for DELETE</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">URL path</strong> — the full endpoint path including path parameters like <Code>/api/users/:id</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Description</strong> — auto-generated from the handler name and route</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Request builder</strong> — fill in parameters and body, then click Send to test the endpoint live</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Response preview</strong> — shows the actual response with syntax highlighting</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Authentication indicator</strong> — shows whether the endpoint requires a JWT token</li>
          </ul>

          <Note>
            The Scalar UI is not just documentation — it{"'"}s a fully functional API client. You can
            send real requests to your running API, see real responses, and debug issues without
            leaving your browser. Think of it as Postman built into your docs.
          </Note>

          <Challenge number={3} title="Test an Endpoint Live">
            <p>Find the <Code>POST /api/auth/login</Code> endpoint in the Scalar UI. Fill in the
            email and password fields in the request builder. Click {'"'}Send.{'"'} Does it return
            a token? What other fields are in the response?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Request & Response Schemas ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Request & Response Schemas</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            One of gin-docs{"'"} most powerful features is automatic schema generation. It reads your
            GORM models and generates JSON schemas that describe the shape of request bodies and
            response payloads.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For example, consider a simple User model:
          </p>

          <CodeBlock filename="apps/api/internal/models/user.go">
{`type User struct {
    ID        uint           \`gorm:"primaryKey" json:"id"\`
    Name      string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Email     string         \`gorm:"uniqueIndex;not null" json:"email" binding:"required,email"\`
    Password  string         \`gorm:"not null" json:"-"\`
    Role      string         \`gorm:"default:USER" json:"role"\`
    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            From this model, gin-docs automatically generates:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> A <strong className="text-foreground">request schema</strong> for POST/PUT — includes name, email, role (excludes id, password, timestamps)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> A <strong className="text-foreground">response schema</strong> for GET — includes id, name, email, role, timestamps (excludes password)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Required field markers</strong> — from the <Code>binding:{'"'}required{'"'}</Code> tag</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Type information</strong> — string, integer, boolean, datetime for each field</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Fields with <Code>json:{'"'}-{'"'}</Code> (like Password and DeletedAt) are excluded from
            the documentation entirely. This is important — your docs never expose sensitive fields.
          </p>

          <Tip>
            The <Code>binding</Code> tag does double duty. It validates incoming requests AND tells
            gin-docs which fields are required. <Code>binding:{'"'}required{'"'}</Code> means the field
            shows as required in the docs, and the API returns a 422 error if it{"'"}s missing.
          </Tip>

          <Challenge number={4} title="Watch Schemas Appear">
            <p>Generate a new resource with several fields:</p>
            <CodeBlock filename="Terminal">
{`grit generate resource Product --fields "name:string,price:float,description:text:optional,sku:string,in_stock:bool"`}
            </CodeBlock>
            <p className="mt-3">Restart your API and refresh <Code>/docs</Code>. Does the Product
            schema appear? Can you see which fields are marked as required vs optional?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Authentication in Docs ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Authentication in Docs</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Many endpoints require authentication — they expect a JWT token in the request headers.
            When you try to call a protected endpoint without a token, you get a 401 Unauthorized
            error. The Scalar UI lets you set your authentication token so you can test protected
            endpoints directly.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the workflow:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span><strong className="text-foreground">Login via the docs.</strong> Find <Code>POST /api/auth/login</Code>, enter your email and password, click Send. Copy the <Code>access_token</Code> from the response.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span><strong className="text-foreground">Set the token.</strong> Look for the {'"'}Authentication{'"'} or {'"'}Auth{'"'} section in Scalar. Select {'"'}Bearer Token{'"'} and paste your access token.</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">3.</span>
              <span><strong className="text-foreground">Test protected endpoints.</strong> Now every request you send from the docs will include the <Code>Authorization: Bearer &lt;token&gt;</Code> header automatically.</span>
            </li>
          </ol>

          <Note>
            Access tokens expire after 15 minutes. If you start getting 401 errors, you need to
            log in again and update the token. In a real app, the frontend handles this automatically
            using refresh tokens. In the docs, you do it manually.
          </Note>

          <Challenge number={5} title="Authenticate and Test">
            <p>Login via <Code>/docs</Code> using <Code>POST /api/auth/login</Code>. Copy the
            access token from the response. Set it in the Scalar authentication section. Then call
            <Code> GET /api/users</Code>. Do you get a list of users? What happens if you remove
            the token and try again?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Testing Endpoints ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Testing Endpoints</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Scalar UI is a complete API testing tool. Let{"'"}s walk through testing all five
            CRUD operations for a resource.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">POST — Create a Resource</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Find the POST endpoint for your resource (e.g., <Code>POST /api/products</Code>).
            The request builder shows the JSON body schema with all required fields. Fill in the
            values and click Send. A successful creation returns a <Code>201 Created</Code> status
            with the new resource in the response body.
          </p>

          <CodeBlock filename="Example: Create a Product">
{`POST /api/products
Content-Type: application/json
Authorization: Bearer <your-token>

{
  "name": "Wireless Keyboard",
  "price": 79.99,
  "description": "Bluetooth mechanical keyboard",
  "sku": "KB-001",
  "in_stock": true
}

// Response: 201 Created
{
  "data": {
    "id": 1,
    "name": "Wireless Keyboard",
    "price": 79.99,
    "description": "Bluetooth mechanical keyboard",
    "sku": "KB-001",
    "in_stock": true,
    "created_at": "2026-03-27T10:00:00Z",
    "updated_at": "2026-03-27T10:00:00Z"
  },
  "message": "Product created successfully"
}`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">GET — List Resources</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            The list endpoint supports pagination. Use query parameters <Code>?page=1&page_size=10</Code> to
            control which page of results you see. The response includes a <Code>meta</Code> object
            with total count, current page, and total pages.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">GET — Single Resource</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Use the ID from a previous creation. For example, <Code>GET /api/products/1</Code> returns
            just that one product. If the ID does not exist, you get a <Code>404 Not Found</Code> error.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">PUT — Update a Resource</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Send only the fields you want to change. For example, to update just the price,
            send <Code>{'{"price": 69.99}'}</Code> to <Code>PUT /api/products/1</Code>. Other
            fields remain unchanged.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">DELETE — Remove a Resource</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Send <Code>DELETE /api/products/1</Code>. The resource is soft-deleted (the DeletedAt
            timestamp is set). It disappears from list queries but remains in the database for
            recovery if needed.
          </p>

          <Challenge number={6} title="Full CRUD Through Docs">
            <p>Using only the Scalar UI at <Code>/docs</Code> (no frontend, no curl, no Postman),
            perform the complete CRUD cycle:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Create a new resource</li>
              <li>List all resources to verify it appears</li>
              <li>Get the single resource by ID</li>
              <li>Update one field</li>
              <li>Get it again to verify the update</li>
              <li>Delete it</li>
              <li>Try to get it again — what status code do you receive?</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Section 7: Response Format ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Response Format</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit enforces a consistent JSON response format across all endpoints. This predictability
            makes frontend development much easier — you always know the shape of the response.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Success — Single Item</h3>
          <CodeBlock filename="Response Format">
{`{
  "data": {
    "id": 1,
    "name": "Wireless Keyboard",
    "price": 79.99
  },
  "message": "Product created successfully"
}`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Success — List with Pagination</h3>
          <CodeBlock filename="Response Format">
{`{
  "data": [
    { "id": 1, "name": "Keyboard" },
    { "id": 2, "name": "Mouse" }
  ],
  "meta": {
    "total": 47,
    "page": 1,
    "page_size": 20,
    "pages": 3
  }
}`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Error Response</h3>
          <CodeBlock filename="Response Format">
{`{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Name is required",
    "details": {
      "name": "This field is required",
      "price": "Must be greater than 0"
    }
  }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The three shapes are always the same:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Single item:</strong> <Code>data</Code> (object) + <Code>message</Code> (string)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">List:</strong> <Code>data</Code> (array) + <Code>meta</Code> (pagination object)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Error:</strong> <Code>error.code</Code> (string) + <Code>error.message</Code> (string) + optional <Code>error.details</Code> (object)</li>
          </ul>

          <Challenge number={7} title="Trigger Each Response Format">
            <p>Trigger all three response formats and observe the structure:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Call <Code>POST /api/auth/register</Code> with valid data — observe the {'"'}data + message{'"'} format</li>
              <li>Call <Code>GET /api/users?page=1&page_size=5</Code> — observe the {'"'}data + meta{'"'} format</li>
              <li>Call <Code>POST /api/auth/register</Code> with an empty body — observe the {'"'}error{'"'} format. What error code and message do you see?</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Section 8: Exporting the Spec ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Exporting the Spec</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The OpenAPI specification is not just for the Scalar UI — it{"'"}s a portable, machine-readable
            file that many tools can consume. You can export it and import it into:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Postman</strong> — import as a collection with all endpoints pre-configured</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Insomnia</strong> — another popular API client that reads OpenAPI specs</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">API client generators</strong> — tools like openapi-generator can create TypeScript, Python, or Java clients from your spec</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">API testing tools</strong> — tools like Dredd or Schemathesis can validate your API against the spec</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The raw OpenAPI JSON is available at:
          </p>

          <CodeBlock filename="OpenAPI Spec URL">
{`# Raw JSON specification
http://localhost:8080/docs/openapi.json

# Copy this URL or download the file
# Import into Postman: File → Import → URL → paste the URL`}
          </CodeBlock>

          <Tip>
            When sharing your API with other teams, you can send them just the OpenAPI spec file.
            They can import it into their tool of choice and have a complete, interactive reference
            of your API — even if they do not have access to your running server.
          </Tip>

          <Challenge number={8} title="Export and Import">
            <p>Open <Code>http://localhost:8080/docs/openapi.json</Code> in your browser. You should
            see the raw JSON specification. Try importing it into Postman: File → Import → paste
            the URL. Does Postman create a collection with all your endpoints?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Customization ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Customization</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            gin-docs provides several configuration options to customize the generated documentation.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3">Authentication Configuration</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            By default, gin-docs detects your JWT middleware and adds Bearer Token authentication
            to the spec. You can customize the auth configuration:
          </p>

          <CodeBlock filename="gin-docs Auth Config">
{`// gin-docs reads your middleware to determine which
// endpoints require authentication.
//
// Protected routes (behind AuthMiddleware) are marked
// with a lock icon in the Scalar UI.
//
// Public routes (like /api/auth/login) show no lock.`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Route Grouping</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Endpoints are grouped by their route prefix. For example, all <Code>/api/users/*</Code> endpoints
            appear under the {'"'}Users{'"'} group. All <Code>/api/auth/*</Code> endpoints appear
            under {'"'}Auth.{'"'} This grouping happens automatically based on your Gin route groups.
          </p>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Excluding Endpoints</h3>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Some internal endpoints (like health checks or metrics) should not appear in public
            documentation. gin-docs excludes internal routes that are not part of your API groups.
            The <Code>/docs</Code> endpoint itself is also excluded — documentation does not document itself.
          </p>

          <Challenge number={9} title="Explore the Configuration">
            <p>Look at how gin-docs is configured in your project. Find the setup code in your
            routes file or main.go. What options are being passed? Can you identify where route
            groups are defined and how they map to the documentation sections?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s review everything you{"'"}ve learned:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Concept</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Key Point</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">gin-docs</td>
                  <td className="px-4 py-3">Auto-generates OpenAPI 3.1 specs by introspecting routes and models</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Scalar UI</td>
                  <td className="px-4 py-3">Interactive docs + API client at /docs</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Schemas</td>
                  <td className="px-4 py-3">Generated from GORM models, respects json:{'"'}-{'"'} and binding tags</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Authentication</td>
                  <td className="px-4 py-3">Set Bearer token in Scalar to test protected endpoints</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Response Format</td>
                  <td className="px-4 py-3">Consistent: data+message, data+meta, or error object</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Export</td>
                  <td className="px-4 py-3">Download OpenAPI JSON, import into Postman/Insomnia</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={10} title="Three Resources, Docs Only">
            <p>Generate three related resources:</p>
            <CodeBlock filename="Terminal">
{`grit generate resource Category --fields "name:string,description:text:optional"
grit generate resource Product --fields "name:string,price:float,sku:string,category_id:belongs_to"
grit generate resource Order --fields "product_id:belongs_to,quantity:int,total:float,status:string"`}
            </CodeBlock>
            <p className="mt-3">Restart your API and open <Code>/docs</Code>. Using ONLY the
            Scalar UI (no frontend, no curl, no code), complete these tasks:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Create 5 categories</li>
              <li>Create 10 products, each linked to a category</li>
              <li>Create 3 orders for different products</li>
              <li>List all products filtered by category</li>
              <li>Update an order status</li>
            </ol>
          </Challenge>

          <Challenge number={11} title="Export Your Spec">
            <p>After creating all three resources, download the OpenAPI spec from
            <Code> /docs/openapi.json</Code>. How many endpoints does it contain? How many schemas?
            Try importing it into Postman or Insomnia.</p>
          </Challenge>

          <Challenge number={12} title="Teach Someone Else">
            <p>Explain to a teammate (or write down in your own words) how Grit{"'"}s API documentation
            works. Cover: (1) how gin-docs generates the spec, (2) how to test endpoints in Scalar,
            (3) how to authenticate in the docs, and (4) how to export the spec. If you can explain
            it clearly, you understand it.</p>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
