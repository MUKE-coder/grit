import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  Step,
  H3,
  Outcome,
  Note,
  Warn,
  Mistake,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Build a REST API in Go with Gin',
  description:
    'Routing, JSON binding, validation, pagination and error responses in Gin, written out by hand, with the handler shape that stays readable at fifty endpoints.',
  alternates: { canonical: 'https://gritframework.dev/guides/gin-rest-api-tutorial' },
  openGraph: {
    title: 'Build a REST API in Go with Gin',
    description: 'From one route to a structured API, with the decisions explained.',
    url: 'https://gritframework.dev/guides/gin-rest-api-tutorial',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="gin-rest-api-tutorial"
      lead={
        <>
          <p className="mb-4">
            Gin is a router with conveniences. It matches a method and a path to a function, parses
            JSON into a struct, validates it against tags, and writes a response. That is nearly
            all of it, and you can learn the whole API in an afternoon.
          </p>
          <p>
            What takes longer is the shape: where validation goes, what an error looks like, how
            pagination is expressed, and how a handler stays readable when there are fifty of
            them. That is most of this page.
          </p>
        </>
      }
    >
      <Outcome>
        <p>
          A books API with list, read, create, update and delete, with validation, pagination,
          filtering and an error format you can keep.
        </p>
        <p>A handler shape that does not degrade as the API grows.</p>
        <p>
          Enough of a grip on binding tags to stop writing validation by hand in <code>if</code>{' '}
          statements.
        </p>
      </Outcome>

      <Section id="routing" title="Routing">
        <CodeBlock
          filename="main.go"
          code={`r := gin.Default()   // Logger and Recovery already attached

// Grouping is not cosmetic: middleware attaches to a group, so the group
// is how you say "everything under here needs a signed-in user".
api := r.Group("/api")
{
	books := api.Group("/books")
	{
		books.GET("",        listBooks)
		books.GET("/:id",    getBook)
		books.POST("",       createBook)
		books.PATCH("/:id",  updateBook)
		books.DELETE("/:id", deleteBook)
	}
}

r.Run(":8080")`}
        />
        <p>Two parameter kinds, and they are not interchangeable:</p>
        <CodeBlock
          code={`// Path parameter: part of the address of a thing.
//   /books/42
id := c.Param("id")

// Query parameter: how you want the answer shaped.
//   /books?page=2&search=tolkien
page := c.DefaultQuery("page", "1")
search := c.Query("search")`}
        />
        <Note title="A URL names a thing; a query describes the view of it">
          <p>
            <code>/books/42</code> is a book. <code>/books?author=42</code> is a list of books,
            filtered. Putting a filter in the path, <code>/books/by-author/42</code>, produces a
            new route for every combination, and you will have five of them within a month.
          </p>
        </Note>
      </Section>

      <Section id="binding" title="Binding and validation">
        <p>
          This is the part of Gin worth learning properly, because it replaces a great deal of
          hand-written checking.
        </p>
        <CodeBlock
          code={`type CreateBookRequest struct {
	Title    string \`json:"title" binding:"required,min=1,max=200"\`
	AuthorID uint   \`json:"author_id" binding:"required"\`
	ISBN     string \`json:"isbn" binding:"omitempty,len=13,numeric"\`
	Pages    int    \`json:"pages" binding:"required,gt=0,lte=10000"\`
	Price    int64  \`json:"price" binding:"gte=0"\`
	Genre    string \`json:"genre" binding:"required,oneof=fiction nonfiction poetry"\`
	Email    string \`json:"contact" binding:"omitempty,email"\`
}

func createBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
			"code":    "VALIDATION_ERROR",
			"message": "some fields are invalid",
			"details": fieldErrors(err),
		}})
		return
	}
	// From here, req is known good. No further checking.
}`}
        />
        <Checklist
          items={[
            <>
              <code>required</code> fails on the zero value, so <code>required</code> on a boolean
              rejects <code>false</code>. Use a <code>*bool</code> when false is a real answer.
            </>,
            <>
              <code>omitempty</code> means &quot;if present, validate it&quot;, which is the right
              pairing for optional fields.
            </>,
            <>
              <code>oneof</code> is how you get an enum, and it is worth using rather than checking
              in code: the rule stays next to the field.
            </>,
            <>
              Use <code>ShouldBindJSON</code>, not <code>BindJSON</code>. The second writes a 400
              and sets a header itself, which fights whatever error format you chose.
            </>,
          ]}
        />
        <H3>Turn the error into something a form can use</H3>
        <p>
          The raw error is one string describing the first failure. A frontend wants to put a
          message under each field:
        </p>
        <CodeBlock
          code={`func fieldErrors(err error) map[string]string {
	out := map[string]string{}
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return out   // malformed JSON rather than a field problem
	}
	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		switch fe.Tag() {
		case "required":
			out[field] = "This field is required"
		case "email":
			out[field] = "Enter a valid email address"
		case "min":
			out[field] = "Must be at least " + fe.Param() + " characters"
		case "oneof":
			out[field] = "Must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
		default:
			out[field] = "This value is not valid"
		}
	}
	return out
}`}
        />
        <Warn title="Bind to a request struct, never to the model">
          <p>
            Binding straight into your GORM model lets a caller send{' '}
            <code>{'{"id": 1, "role": "admin", "created_at": "..."}'}</code> and set fields you
            never meant to expose. A separate request struct is an allow-list by construction: a
            field that is not on it cannot be set. This is the mass-assignment vulnerability, and
            the defence is simply not having the field there.
          </p>
        </Warn>
      </Section>

      <Section id="shape" title="The shape of a handler">
        <p>
          A handler has one job: translate HTTP into a function call and the result back into
          HTTP. Business logic belongs below it. This matters less at five endpoints and a great
          deal at fifty.
        </p>
        <CodeBlock
          filename="internal/handlers/book.go"
          code={`// Thin. Reads input, calls the service, writes output. Nothing else.
func (h *BookHandler) Create(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.ValidationError(c, err)
		return
	}

	book, err := h.books.Create(c.Request.Context(), req)
	if err != nil {
		respond.Error(c, err)   // maps a domain error to a status code
		return
	}

	respond.Created(c, book)
}`}
        />
        <CodeBlock
          filename="internal/services/book.go"
          code={`// The rules live here, where a test can reach them without HTTP.
func (s *BookService) Create(ctx context.Context, req CreateBookRequest) (*models.Book, error) {
	if exists, err := s.repo.ISBNExists(ctx, req.ISBN); err != nil {
		return nil, fmt.Errorf("check isbn: %w", err)
	} else if exists {
		return nil, ErrDuplicateISBN
	}

	book := &models.Book{Title: req.Title, AuthorID: req.AuthorID, ISBN: req.ISBN}
	if err := s.repo.Create(ctx, book); err != nil {
		return nil, fmt.Errorf("create book: %w", err)
	}
	return book, nil
}`}
        />
        <p>
          The test for &quot;a duplicate ISBN is rejected&quot; now calls a function with a struct,
          rather than building an HTTP request and parsing a response. That is the practical
          payoff, and it arrives the first time you write the second test.
        </p>
        <p>
          Pass <code>c.Request.Context()</code> down. When the client disconnects, that context is
          cancelled, and a database query holding it stops instead of finishing work nobody is
          waiting for.
        </p>
      </Section>

      <Section id="errors" title="One error shape, decided once">
        <p>
          Pick a response format before the third endpoint, because changing it later touches every
          client you have.
        </p>
        <CodeBlock
          language="json"
          code={`// One thing
{ "data": { "id": 1, "title": "Dune" } }

// A list
{
  "data": [ ... ],
  "meta": { "total": 240, "page": 1, "page_size": 20, "pages": 12 }
}

// A failure
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "some fields are invalid",
    "details": { "isbn": "Must be exactly 13 characters" }
  }
}`}
        />
        <p>
          The <code>code</code> is for your frontend to branch on and the <code>message</code> is
          for a human. Branching on the message means a copy edit breaks a client.
        </p>
        <H3>Map domain errors to status codes in one place</H3>
        <CodeBlock
          filename="internal/respond/respond.go"
          code={`func Error(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		Fail(c, http.StatusNotFound, "NOT_FOUND", "not found")
	case errors.Is(err, services.ErrDuplicateISBN):
		Fail(c, http.StatusConflict, "DUPLICATE_ISBN", "that ISBN already exists")
	case errors.Is(err, services.ErrForbidden):
		Fail(c, http.StatusForbidden, "FORBIDDEN", "you may not do that")
	default:
		// Log the real error with the request id; tell the caller nothing.
		// A driver error names your tables and sometimes your file paths.
		log.Error().Err(err).Str("request_id", RequestID(c)).Msg("unhandled")
		Fail(c, http.StatusInternalServerError, "INTERNAL", "something went wrong")
	}
}`}
        />
        <p>Which status, briefly, because the usual three are not enough:</p>
        <Checklist
          items={[
            <>
              <strong>400</strong> the body is malformed, such as invalid JSON.{' '}
              <strong>422</strong> the body parsed and a field is wrong. Distinguishing them tells
              a client whether to blame the user or itself.
            </>,
            <>
              <strong>401</strong> not signed in. <strong>403</strong> signed in and not allowed.
            </>,
            <>
              <strong>404</strong> no such thing, <em>and</em> the answer for a record outside the
              caller&apos;s scope: a 403 there confirms the record exists.
            </>,
            <>
              <strong>409</strong> a conflict with current state, such as a unique constraint.
            </>,
            <>
              <strong>201</strong> for a create, with the created object in the body.
            </>,
          ]}
        />
      </Section>

      <Section id="pagination" title="Pagination, filtering and sorting">
        <p>
          A list endpoint without pagination is a time bomb with a fuse the length of your table.
          It is fine on day one with forty rows and takes the service down at four hundred
          thousand.
        </p>
        <CodeBlock
          code={`type ListQuery struct {
	Page     int    \`form:"page,default=1" binding:"min=1"\`
	PageSize int    \`form:"page_size,default=20" binding:"min=1,max=100"\`
	Search   string \`form:"search"\`
	Sort     string \`form:"sort,default=created_at"\`
	Order    string \`form:"order,default=desc" binding:"oneof=asc desc"\`
}

func (h *BookHandler) List(c *gin.Context) {
	var q ListQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		respond.ValidationError(c, err)
		return
	}
	// max=100 above is the important tag. Without a ceiling, page_size is
	// a parameter the caller controls that decides how much memory you
	// allocate, which is a denial of service you shipped yourself.
	...
}`}
        />
        <Warn title="Never interpolate a sort column into SQL">
          <p>
            <code>db.Order(q.Sort + &quot; &quot; + q.Order)</code> is SQL injection with extra
            steps: the value came from the query string. Check it against a list of columns you
            allow, and reject anything else.
          </p>
          <CodeBlock
            code={`var sortable = map[string]bool{
	"created_at": true, "title": true, "pages": true, "price": true,
}

if !sortable[q.Sort] {
	q.Sort = "created_at"
}
db = db.Order(clause.OrderByColumn{
	Column: clause.Column{Name: q.Sort},
	Desc:   q.Order == "desc",
})`}
          />
          <p>
            Building the clause rather than concatenating a string also quotes the identifier
            properly, which matters the day a column is called <code>order</code>.
          </p>
        </Warn>
        <p>Then count and fetch:</p>
        <CodeBlock
          code={`var total int64
query := s.db.WithContext(ctx).Model(&models.Book{})
if q.Search != "" {
	query = query.Where("title ILIKE ?", "%"+q.Search+"%")
}
if err := query.Count(&total).Error; err != nil {
	return nil, 0, fmt.Errorf("count books: %w", err)
}

var books []models.Book
err := query.
	Preload("Author").                       // or the list makes one query per row
	Offset((q.Page - 1) * q.PageSize).
	Limit(q.PageSize).
	Find(&books).Error`}
        />
        <Note title="Offset gets slow, eventually">
          <p>
            <code>OFFSET 100000</code> makes the database walk a hundred thousand rows to discard
            them. Fine for an admin table nobody pages that deep into; not fine for an infinite
            scroll. The alternative is a cursor: remember the last row&apos;s sort value and ask
            for rows after it, which stays constant time.
          </p>
        </Note>
      </Section>

      <Section id="testing" title="Testing handlers">
        <p>
          Gin&apos;s router is an <code>http.Handler</code>, so you can exercise it with the
          standard library and no server, no ports and no sleeping.
        </p>
        <CodeBlock
          filename="internal/handlers/book_test.go"
          code={`func TestCreateBookRejectsEmptyTitle(t *testing.T) {
	r := testRouter(t)

	body := \`{"title":"","author_id":1,"pages":100,"genre":"fiction"}\`
	req := httptest.NewRequest("POST", "/api/books", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)

	require.Equal(t, http.StatusUnprocessableEntity, res.Code)
	require.Contains(t, res.Body.String(), "title")
}`}
        />
        <Warn title="In-memory SQLite needs exactly one connection">
          <p>
            Every connection to <code>:memory:</code> opens its own empty database. With a pool,
            your migration runs on one connection and your query on another, which reports{' '}
            <code>no such table</code> and, depending on how the test is written, still passes.
            Pin the pool to one connection in tests.
          </p>
          <CodeBlock
            code={`sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(1)`}
          />
        </Warn>
      </Section>

      <Section id="mistakes" title="When it goes wrong">
        <Mistake error="EOF (from ShouldBindJSON)">
          <p>
            The body was empty. Usually a request sent with no body at all, or a client that set a
            body but no <code>Content-Type: application/json</code>. Also note that{' '}
            <code>ShouldBindJSON</code> consumes the body: calling it twice gives you EOF the
            second time.
          </p>
        </Mistake>

        <Mistake error="json: cannot unmarshal string into Go struct field ... of type int">
          <p>
            The client sent <code>&quot;42&quot;</code> where the struct wants a number. Common
            from HTML forms, where everything is a string. Either fix the client, or accept a
            string in the request struct and convert deliberately.
          </p>
        </Mistake>

        <Mistake error="Key: 'CreateBookRequest.Title' Error:Field validation for 'Title' failed on the 'required' tag">
          <p>
            Validation worked. This is the raw message, and it is fine in a log and wrong in a
            response: it names your Go types. Map it to field errors, as above, before sending it
            anywhere.
          </p>
        </Mistake>

        <Mistake error="A field is always empty after binding, with no error">
          <p>
            The <code>json</code> tag does not match what the client sends, so the field is simply
            never populated and nothing complains. Check the exact spelling, including
            snake_case against camelCase. Add <code>binding:&quot;required&quot;</code> and the
            silence turns into an error, which is how you find it.
          </p>
        </Mistake>

        <Mistake error="404 page not found on a route you registered">
          <p>
            Usually a trailing slash: <code>/api/books/</code> and <code>/api/books</code> are
            different routes, and Gin redirects between them by default but not for every method.
            It is also what you see when a route was registered on a group whose prefix you
            forgot.
          </p>
          <CodeBlock terminal code={`# Print every route Gin knows about at startup
# (gin.Default() already logs them)`} />
        </Mistake>
      </Section>

      <Section id="with-grit" title="The same API, generated">
        <p>
          Everything above is one resource. Grit writes it from a field list, in the same Gin and
          GORM, into your repository:
        </p>
        <CodeBlock
          terminal
          code={`grit generate resource Book \\
  --fields "title:string:required,author:belongs_to:Author,isbn:string:unique,pages:int,price:money,genre:enum(fiction|nonfiction|poetry)"`}
        />
        <p>
          You get the model, the service, the handler with pagination, filtering and sorting
          against an allow-list, the routes, the error mapping, and the matching TypeScript type,
          Zod schema and React Query hook, so the frontend cannot disagree with the backend about
          what a book is.
        </p>
        <p>
          The handler it writes is the thin one from above, because the shape is the point: code
          you would not want to maintain is not worth generating.
        </p>
      </Section>

      <Faq
        slug="gin-rest-api-tutorial"
        items={[
          {
            q: 'Gin, Echo, Fiber or net/http?',
            a: (
              <>
                Gin has the largest ecosystem and the most examples, which is what you want when
                stuck. Echo is very similar, and the choice between them is taste.{' '}
                <code>net/http</code> alone is viable since its router gained method matching and
                path wildcards; you write the binding and validation yourself. Fiber is fast and
                does not implement the standard <code>net/http</code> interfaces, so middleware and
                test helpers from the wider ecosystem do not work with it.
              </>
            ),
            text:
              'Gin has the largest ecosystem and the most examples, which matters when you are stuck. Echo is very similar and the choice is taste. The standard library router is viable now that it matches methods and path wildcards, at the cost of writing binding and validation yourself. Fiber is fast but does not implement the standard net/http interfaces, so wider-ecosystem middleware and test helpers do not work with it.',
          },
          {
            q: 'Do I need a repository layer, or can the service use GORM directly?',
            a: (
              <>
                For most applications the service can use GORM directly, and a repository that only
                forwards calls is indirection with no payoff. Add one when you have a real reason:
                a second data source, a cache in front of the database, or a need to test the
                service with no database at all. Adding it on principle is the most common piece of
                unnecessary structure in Go APIs.
              </>
            ),
            text:
              'For most applications the service can use GORM directly, and a repository that only forwards calls adds indirection with no payoff. Add one for a real reason: a second data source, a cache in front of the database, or testing the service without a database. Adding it on principle is the most common unnecessary structure in Go APIs.',
          },
          {
            q: 'Should PATCH and PUT behave differently?',
            a: (
              <>
                Yes. <code>PUT</code> replaces the whole resource, so an omitted field means
                &quot;set it to empty&quot;. <code>PATCH</code> applies a partial change, so an
                omitted field means &quot;leave it&quot;. For PATCH that distinction needs pointer
                fields in the request struct, because a plain string cannot tell &quot;absent&quot;
                from &quot;empty string&quot;. Most APIs only need PATCH.
              </>
            ),
            text:
              'Yes. PUT replaces the whole resource, so an omitted field means set it to empty. PATCH applies a partial change, so an omitted field means leave it. PATCH needs pointer fields in the request struct, since a plain string cannot distinguish absent from empty. Most APIs only need PATCH.',
          },
          {
            q: 'Where do I put middleware that only some routes need?',
            a: (
              <>
                On the group. <code>r.Group(&quot;/api/admin&quot;, middleware.Auth(),
                middleware.Require(perm))</code> applies to everything in it, and the grouping
                becomes a readable statement of the policy. Note that Gin builds the chain when a
                route is registered, so middleware added after a route does not apply to it, with
                no warning.
              </>
            ),
            text:
              'On the group: r.Group("/api/admin", middleware.Auth(), middleware.Require(perm)) applies to everything inside it and makes the policy readable. Gin builds the handler chain when a route is registered, so middleware added after a route does not apply to it and nothing warns you.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/gorm-tutorial-relationships',
            label: 'GORM relationships',
            blurb: 'The models underneath these handlers, and loading them without N+1.',
          },
          {
            href: '/guides/golang-jwt-authentication',
            label: 'JWT authentication in Go',
            blurb: 'Putting a signed-in user behind the middleware these routes use.',
          },
          {
            href: '/docs/backend/handlers',
            label: 'Handlers reference',
            blurb: 'The generated handler, and what each part of it is doing.',
          },
          {
            href: '/docs/backend/response-format',
            label: 'The response contract',
            blurb: 'The envelope, the status codes and the helpers that write them.',
          },
        ]}
      />
    </GuidePage>
  )
}
