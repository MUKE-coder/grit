#import "../template.typ": *

= Project: a publication, from draft to public page

Every project in Part V so far has been behind a login. This one is not, and
that changes almost everything: the audience is the internet, the content is
written by people rather than generated, and the pages have to be fast for
somebody who has never signed in and never will.

We are building a publication: posts with a rich text body, images, a draft and
published workflow, a public API and a web front end that renders it.

== The model you already have

A new project ships with a `Blog` model, and it is worth reading before
replacing it, because it is a worked example of several chapters at once:

```go
type Blog struct {
    ID          string     `gorm:"primarykey;size:36"`
    Title       string     `gorm:"size:255;not null"
                            binding:"required"`
    Slug        string     `gorm:"size:255;uniqueIndex"`
    Content     string     `gorm:"type:text" sanitize:"html"`
    Image       string     `gorm:"size:500"`
    Excerpt     string     `gorm:"size:500"`
    Published   bool       `gorm:"default:false"`
    PublishedAt *time.Time
    Version     int        `gorm:"not null;default:1"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

Five decisions in one struct:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`uniqueIndex` on slug], [The URL is a key, and the database enforces it],
  [`sanitize:"html"`], [The body is cleaned on every write. See below],
  [`Published` plus `PublishedAt`], [State and the moment it changed],
  [`Version`], [Optimistic locking, from chapter thirty three],
  [`DeletedAt`], [Soft delete, so an unpublish is never a loss],
)

The slug is generated in a `BeforeCreate` hook, which is the right place for it:
a title is typed once and the slug must exist whichever surface created the
row, including a seeder and an import.

== The one that is not optional: sanitising HTML

A rich text editor produces HTML, and HTML from a user is an attack surface
with a text box attached. The tag says what to do about it:

```go
Content string `gorm:"type:text" sanitize:"html"`
```

Every write that has a model passes through the sanitiser, so a `<script>` in a
draft does not survive being saved. Not on render, on *write*, which matters
because the content will be read by more things than your web page: an RSS
feed, a mobile app, a newsletter.

#gotcha[
  Sanitising on render rather than on write is the mistake that keeps coming
  back, because it works perfectly until the second reader exists. The database
  holds the attack, and every new consumer has to remember to defuse it.

  Sanitise once, at the boundary where the data arrives. Then the thing in the
  database is safe, and the number of places that have to be careful is one
  rather than one per feature.
]

== Draft and published, as a service rule

Publishing is not a column edit. It is a transition with consequences:

```go
func (s *BlogService) Publish(ctx context.Context,
    id string) (*models.Blog, error) {

    post, err := s.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if strings.TrimSpace(post.Content) == "" {
        return nil, ErrEmptyPost
    }
    if post.Published {
        return post, nil  // idempotent on purpose
    }

    now := time.Now().UTC()
    post.Published = true
    post.PublishedAt = &now

    if err := s.db(ctx).Save(post).Error; err != nil {
        return nil, fmt.Errorf("publishing %s: %w", id, err)
    }
    events.Emit(ctx, events.BlogPublished{ID: post.ID})
    return post, nil
}
```

Three things in that function are worth copying into whatever you build next.

*It refuses an empty post.* The rule lives here, so the admin, the API and a
scheduled publish all obey it.

*It is idempotent.* Publishing something already published returns it rather
than erroring. Retries happen, double clicks happen, and a second publish
should not move `PublishedAt`.

*It emits an event* rather than doing the downstream work inline. Chapter
forty four is about the bus; what matters here is that "tell the newsletter" is
not the publishing function's business.

== The public surface

The generated endpoints are behind authentication, which is right for an admin
and useless for readers. A public-facing resource asks for a second, narrower
surface:

```bash
grit generate resource Article \
  --fields "title:string, slug:slug, body:richtext, \
cover:file, published:bool" --public
```

That adds a list and a get-by-slug under `/api/v1/public/`, mounted outside the
auth middleware. The protected routes are untouched, because the admin still
calls those.

Two properties of that surface are the whole design.

*The response is an allowlist, not the model.* A public response is its own
struct, built from the fields that are publishable by default: names, slugs,
descriptions, images. Anything whose name contains `cost`, `margin`, `profit`,
`internal`, `note`, `secret`, `token`, `supplier` or `wholesale` is excluded
whatever its type, and foreign keys and stock counts are left out too.

*A column added next month is private until somebody says otherwise.* That is
the opposite default to the admin surface, and it is the right one when the
audience is the internet: the alternative, publishing everything and trusting
the developer to remove what should not be there, puts a `cost_price` column on
the web the first time somebody adds one in a hurry.

To publish a field that was excluded, you edit one struct, deliberately, in a
diff somebody can review.

#gotcha[
  The public group is guarded by an API key, which is not secrecy: a
  publishable key ships inside your app where anybody can read it.

  What it buys is identification, a rate limit bucket per key, per-endpoint and
  per-origin narrowing, and the ability to turn one client off without a
  deploy. Treat it as a name badge rather than a password.
]

== The cache, and the one rule about it

Responses on the public group are cached, and nowhere else. The cache key is
the URL and nothing else:

```go
publicAPI.Use(middleware.CacheResponse(svc.Cache, ttl))
```

On a public endpoint that is exactly right. Every caller gets the same answer,
so one cached copy serves all of them and an article page stops hitting the
database on every visit.

#gotcha[
  *Never mount a writing endpoint on the public group, and never move this
  cache onto an authenticated one.*

  The second is the dangerous one: with the URL as the key, `GET /api/v1/
  invoices` cached for one user is served to the next user who asks. One line
  of middleware in the wrong group is a data breach that passes every test,
  because tests usually run one user at a time.

  `--public` is read-only by design. A public *write*, a comment form, a
  newsletter signup, is a hand-written endpoint with its own rate limit and its
  own validation, and it does not belong on a cached group.
]

== Rendering it in the web app

The public list, in a server component:

```tsx
async function getArticles() {
  const res = await fetch(`${API}/api/v1/public/articles`, {
    headers: { "X-API-Key": process.env.PUBLIC_API_KEY! },
    next: { revalidate: 60 },
  })
  if (!res.ok) throw new Error("articles unavailable")
  return (await res.json()) as { data: Article[]; meta: Meta }
}
```

And the article page, by slug rather than ID, which is why the public endpoint
looks up by slug:

```tsx
export default async function ArticlePage(
  { params }: { params: { slug: string } },
) {
  const article = await getArticle(params.slug)
  return (
    <article className="prose prose-invert mx-auto">
      <h1>{article.title}</h1>
      <div dangerouslySetInnerHTML={{ __html: article.body }} />
    </article>
  )
}
```

#bridge[
  `dangerouslySetInnerHTML` is doing what its name says, and it is acceptable
  here for one reason only: the HTML was sanitised on write.

  If you are ever tempted to render untrusted HTML that has not been through a
  sanitiser, the answer is not a careful regular expression. It is a library,
  on the way in.
]

Metadata comes from the same fetch, so the social preview and the page agree:

```tsx
export async function generateMetadata({ params }) {
  const a = await getArticle(params.slug)
  return {
    title: a.title,
    description: a.excerpt,
    openGraph: { images: [a.cover] },
  }
}
```

== The things a publication needs that nothing generates

Four small pieces, each about twenty lines, none of them framework specific:

*A sitemap.* `app/sitemap.ts` fetching the public list and returning the URLs
with their `published_at` as `lastModified`.

*An RSS feed.* A route handler returning XML. Readers still use them, and so do
aggregators that will send you traffic.

*Related posts.* The public surface can expose a related endpoint; otherwise it
is a query by tag with the current article excluded.

*A reading time.* Word count divided by two hundred, computed on write and
stored, because computing it on every render is work repeated for no reason.

#project(title: "Project 36: ship the publication")[
  Build it properly, including the parts that are only visible when something
  goes wrong.

  ```bash
  grit new press && cd press
  grit generate resource Article \
    --fields "title:string, slug:slug, body:richtext, \
cover:file, excerpt:text, published:bool, \
author:belongs_to:User, tags:many_to_many:Tag" --public
  grit generate resource Tag --fields "name:string, slug:slug" --public
  grit migrate && grit seed Article --count 60
  ```

  *One: the workflow.* Add `Publish` and `Unpublish` to the service with the
  rules above, plus endpoints. Publishing twice must not move `published_at`.

  *Two: the public surface, inspected.* Fetch
  `/api/v1/public/articles` and compare the JSON with the model. List in
  `NOTES.md` every field that did not come out, and say for each whether you
  agree with the default.

  *Three: drafts must not leak.* Create a draft, note its slug, and request it
  from the public endpoint. You should get a 404. Then check the list's
  `meta.total` and confirm the draft is not counted. A total that leaks the
  number of unpublished drafts is a smaller version of the same bug.

  *Four: the cache, observed.* Request the list twice and compare response
  times. Then publish a new article and watch how long it takes to appear. Note
  the TTL and decide whether sixty seconds is right for your publication.

  *Five: the attack.* Paste this into a draft body through the admin editor,
  save, and then read the stored row in `/studio`:

  ```html
  <img src=x onerror="alert(document.cookie)">
  ```

  Record exactly what the database holds. That is the sanitiser, and seeing it
  work is more convincing than reading that it exists.

  *Six: the front end.* A list page, an article page by slug, metadata, a
  sitemap and an RSS feed. Run Lighthouse on the article page and keep the
  score.
]

#tryit[
  Turn the API off and reload your article page.

  Whatever happens next is what a reader sees on your worst day. A blank page
  with a stack trace, a spinner forever, or a cached page with a quiet notice
  are three different products, and only one of them is a decision rather than
  an accident.
]

#assignment(name: "Publication - Assignment One")[
  In `part5-building/39-publication`, build the publication and prove the
  boundary between private and public.

  *One.* The `Article` resource with `--public`, the publish and unpublish
  rules in the service, and tests for both, including the idempotent publish.

  *Two.* `PUBLIC.md` containing: the full model field list, the fields the
  public response actually returns, and one line each on why the difference is
  correct or what you changed.

  *Three.* Evidence that a draft is invisible publicly: the 404, and the
  `meta.total` with and without drafts.

  *Four.* The XSS test from Project 36, with the stored row pasted in. If your
  sanitiser let something through, that is the finding, and it is worth more
  than a passing test.

  *Five.* A front end with a list, an article page, metadata, a sitemap and an
  RSS feed, plus one paragraph on what you would cache and for how long if this
  publication had a hundred thousand readers a day.
]
