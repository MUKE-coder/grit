import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build a Blog & CMS: Complete Content Management System',
  description: 'Build a complete blog with content management using Grit. Categories, tags, rich text posts, comments with moderation, SEO, RSS feeds, and media library.',
}

export default function BlogCMSCourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Blog & CMS</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build a Blog & CMS: Complete Content Management System
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Build a full-featured blog with categories, tags, rich text posts, comments with moderation,
            SEO-friendly slugs, and a media library. You{"'"}ll scaffold the project, generate 4 resources,
            and wire up a complete content management workflow.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is a CMS? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is a CMS?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            If you{"'"}ve ever used WordPress, you{"'"}ve used a CMS. A CMS lets you create, edit,
            organize, and publish content without writing code every time. The admin panel is where
            you manage everything. The public site is what readers see. We{"'"}re building both.
          </p>

          <Definition term="CMS (Content Management System)">
            A system for creating, managing, and publishing content. It separates content from
            presentation — you write in the admin panel, and the public site renders it beautifully.
            WordPress, Ghost, and Strapi are all CMSes. We{"'"}re building one with Grit.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every CMS needs these core pieces:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Content types</strong> — posts, pages, categories, tags</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Rich text editor</strong> — write with headings, bold, lists, images, code blocks</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Media library</strong> — upload and manage images and files</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">User roles</strong> — admins manage content, readers view it</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">SEO</strong> — slugs, meta descriptions, Open Graph tags</li>
          </ul>

          <Tip>
            Grit{"'"}s code generator handles 80% of the CMS work. You define the data model, generate
            resources, and get full CRUD with admin panel, API, and frontend pages — all in minutes.
          </Tip>

          <Challenge number={1} title="Name 3 Popular CMS Platforms">
            <p>Name 3 popular CMS platforms. For each one, describe what kind of content it manages
            and who its typical users are. Think beyond just WordPress — consider headless CMSes,
            e-commerce platforms, and documentation tools.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Plan the Data Model ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Plan the Data Model</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Before writing any code, plan your data model. A blog CMS needs 4 resources that work
            together: categories organize posts, tags add cross-cutting labels, posts hold the content,
            and comments let readers respond.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the complete data model:
          </p>

          <CodeBlock filename="Resource: Category">
{`grit generate resource Category \\
  name:string:unique \\
  slug:slug:name \\
  description:text:optional`}
          </CodeBlock>

          <CodeBlock filename="Resource: Tag">
{`grit generate resource Tag \\
  name:string:unique \\
  slug:slug:name`}
          </CodeBlock>

          <CodeBlock filename="Resource: Post">
{`grit generate resource Post \\
  title:string \\
  slug:slug:title \\
  content:richtext \\
  excerpt:text:optional \\
  published:bool \\
  featured_image:string:optional \\
  author_id:belongs_to:User \\
  category_id:belongs_to:Category \\
  tag_ids:many_to_many:Tag`}
          </CodeBlock>

          <CodeBlock filename="Resource: Comment">
{`grit generate resource Comment \\
  content:text \\
  author_name:string \\
  author_email:string \\
  post_id:belongs_to:Post \\
  approved:bool`}
          </CodeBlock>

          <Note>
            Order matters. Generate Category and Tag before Post, because Post references them
            with <Code>belongs_to</Code> and <Code>many_to_many</Code>. Generate Comment last
            since it references Post.
          </Note>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Notice the field types:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>slug:slug:name</Code> — auto-generates a URL-friendly slug from the name field</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>content:richtext</Code> — renders a rich text editor in the admin panel</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>belongs_to:User</Code> — creates a foreign key relationship with automatic joins</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>many_to_many:Tag</Code> — creates a join table so posts can have multiple tags</li>
          </ul>

          <Challenge number={2} title="Generate All 4 Resources">
            <p>Run all 4 generate commands in order: Category, Tag, Post, Comment. Check that each
            resource creates the expected files: model, service, handler, Zod schema, TypeScript
            types, and admin page. How many files did Grit generate in total?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Scaffold the Project ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold the Project</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Start by scaffolding a fresh Grit project with the triple architecture — API, web frontend,
            and admin panel. Then generate all resources and seed test data.
          </p>

          <CodeBlock filename="Terminal">
{`# Scaffold the project
grit new myblog --triple --next

# Start infrastructure
cd myblog
docker compose up -d

# Generate resources (in order)
grit generate resource Category name:string:unique slug:slug:name description:text:optional
grit generate resource Tag name:string:unique slug:slug:name
grit generate resource Post title:string slug:slug:title content:richtext excerpt:text:optional published:bool featured_image:string:optional author_id:belongs_to:User category_id:belongs_to:Category tag_ids:many_to_many:Tag
grit generate resource Comment content:text author_name:string author_email:string post_id:belongs_to:Post approved:bool

# Start the API (runs migrations automatically)
cd apps/api && go run cmd/server/main.go`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Once the API is running, open the admin panel and start creating content. Build a realistic
            dataset — real blogs have many categories, tags spread across posts, and a mix of published
            and draft posts.
          </p>

          <Tip>
            Use the API docs at <Code>/docs</Code> for quick data entry. You can POST JSON directly
            without navigating the admin UI. Great for seeding large amounts of test data.
          </Tip>

          <Challenge number={3} title="Create Seed Data">
            <p>Create 5 categories (e.g., Technology, Design, Business, Lifestyle, Tutorials), 10 tags
            (e.g., go, react, api, tutorial, beginner, advanced, tips, tools, review, guide), and 20
            blog posts spread across different categories with various tags. Make 12 published and 8 drafts.</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: The Blog List Page ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Blog List Page</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The public blog page shows published posts with filtering and pagination. The API already
            supports all of this — Grit generates list endpoints with built-in query parameters for
            filtering, sorting, and pagination.
          </p>

          <CodeBlock filename="API Queries">
{`# List published posts (page 1, 10 per page)
GET /api/posts?published=true&page=1&page_size=10

# Filter by category
GET /api/posts?published=true&category_id=3

# Filter by tag (many-to-many)
GET /api/posts?published=true&tag_ids=5

# Sort by newest first
GET /api/posts?published=true&sort=created_at&order=desc

# Search by title
GET /api/posts?published=true&search=getting+started`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The response includes pagination metadata:
          </p>

          <CodeBlock filename="API Response">
{`{
  "data": [
    {
      "id": 1,
      "title": "Getting Started with Go",
      "slug": "getting-started-with-go",
      "excerpt": "Learn the basics of Go programming...",
      "published": true,
      "category": { "id": 1, "name": "Technology" },
      "tags": [
        { "id": 1, "name": "go" },
        { "id": 4, "name": "tutorial" }
      ],
      "created_at": "2026-03-15T10:30:00Z"
    }
  ],
  "meta": {
    "total": 12,
    "page": 1,
    "page_size": 10,
    "pages": 2
  }
}`}
          </CodeBlock>

          <Note>
            The <Code>published=true</Code> filter is critical. Without it, draft posts would appear
            on the public blog. Always filter by published status on public-facing endpoints.
          </Note>

          <Challenge number={4} title="Test the Blog List Endpoint">
            <p>Open the API docs at <Code>/docs</Code> and test these queries: (1) List all published
            posts with <Code>page_size=5</Code>, (2) Filter by a specific category, (3) Sort by
            title ascending. How many total published posts does each query return?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: The Blog Detail Page ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Blog Detail Page</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Each blog post has a detail page showing the full content. The slug field makes URLs
            human-readable: <Code>/blog/getting-started-with-go</Code> instead of <Code>/blog/1</Code>.
          </p>

          <CodeBlock filename="Fetch Post by Slug">
{`# Get a single post by slug
GET /api/posts?slug=getting-started-with-go

# The response includes all related data:
# - Author (User) with name and avatar
# - Category with name and slug
# - Tags array
# - Full rich text content
# - Featured image URL
# - Created/updated timestamps`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The rich text content field stores HTML generated by the editor. It supports headings,
            bold, italic, lists, code blocks, images, and links. The frontend renders this HTML
            directly with proper styling.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A complete blog detail page displays:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Title</strong> — the post headline, large and prominent</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Featured image</strong> — full-width hero image at the top</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Meta info</strong> — author name, publish date, category badge, reading time</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Content</strong> — the full rich text body</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Tags</strong> — clickable tag badges that link to filtered lists</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Comments</strong> — approved comments with a submission form</li>
          </ul>

          <Challenge number={5} title="Create a Rich Content Post">
            <p>Create a blog post with rich text content — include at least 2 headings, bold text,
            a bulleted list, and a code block. View it via the API using the slug. Does the HTML
            content come through correctly?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Comment System ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Comment System</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Comments let readers respond to posts. The key design decision: comments require moderation.
            When someone submits a comment, it{"'"}s created with <Code>approved: false</Code>. An admin
            must approve it before it appears publicly. This prevents spam.
          </p>

          <CodeBlock filename="Submit a Comment">
{`POST /api/comments
Content-Type: application/json

{
  "content": "Great article! Really helped me understand Go basics.",
  "author_name": "Jane Reader",
  "author_email": "jane@example.com",
  "post_id": 1,
  "approved": false
}

# Response: 201 Created
# The comment exists but is not visible publicly until approved`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The public API only returns approved comments:
          </p>

          <CodeBlock filename="Fetch Approved Comments">
{`# Public: only approved comments for a post
GET /api/comments?post_id=1&approved=true

# Admin: all comments (approved + pending)
GET /api/comments?post_id=1

# Admin: approve a comment
PATCH /api/comments/5
{ "approved": true }`}
          </CodeBlock>

          <Tip>
            The admin panel{"'"}s comment list should highlight unapproved comments. Think of it as a
            moderation queue — new comments appear at the top, waiting for approval. One click to
            approve, one click to delete spam.
          </Tip>

          <Challenge number={6} title="Test Comment Moderation">
            <p>Add 5 comments to a post using the API. Approve 3 of them in the admin panel. Then
            query <Code>GET /api/comments?post_id=1&approved=true</Code>. How many comments appear?
            The answer should be 3. If all 5 appear, your filter isn{"'"}t working.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: SEO for Blog Posts ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">SEO for Blog Posts</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Search engine optimization starts with good URLs. The <Code>slug</Code> field type
            automatically generates URL-friendly strings from the title: {'"'}Getting Started with
            Go{'"'} becomes <Code>getting-started-with-go</Code>. No IDs in URLs, no random strings.
          </p>

          <Definition term="Slug">
            A URL-friendly version of a string. Lowercase, hyphens instead of spaces, no special
            characters. {'"'}My First Blog Post!{'"'} becomes <Code>my-first-blog-post</Code>. Slugs
            make URLs readable and help search engines understand your content.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For a complete SEO setup, each blog post page needs:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">URL</strong> — <Code>/blog/getting-started-with-go</Code> (slug-based)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Title tag</strong> — {'"'}Getting Started with Go | MyBlog{'"'}</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Meta description</strong> — uses the excerpt field (keep under 160 characters)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Open Graph tags</strong> — og:title, og:description, og:image (featured image)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Canonical URL</strong> — prevents duplicate content issues</li>
          </ul>

          <CodeBlock filename="Next.js Metadata (Simplified)">
{`// In your blog/[slug]/page.tsx
// Fetch the post data, then export metadata:
//
// title: post.title + " | MyBlog"
// description: post.excerpt
// openGraph.title: post.title
// openGraph.description: post.excerpt
// openGraph.images: [post.featured_image]
// canonical: "https://myblog.com/blog/" + post.slug`}
          </CodeBlock>

          <Challenge number={7} title="Check Your Slugs">
            <p>Look at 5 of your blog posts. Are the slugs URL-friendly? Check for: all lowercase,
            hyphens instead of spaces, no special characters, no trailing hyphens. Try creating a
            post with a title that has special characters (like {'"'}Go & React: A Guide!{'"'}) —
            what does the generated slug look like?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: RSS Feed ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">RSS Feed</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            RSS feeds let readers subscribe to your blog using feed readers like Feedly or Inoreader.
            An RSS feed is an XML document listing your latest posts with title, link, description,
            and publish date.
          </p>

          <Definition term="RSS (Really Simple Syndication)">
            An XML format for distributing content updates. Feed readers poll RSS URLs periodically
            to show new posts. Most blogs, podcasts, and news sites publish RSS feeds. It{"'"}s how
            people follow content without checking websites manually.
          </Definition>

          <CodeBlock filename="RSS Feed Endpoint (Go)">
{`// GET /api/feed/rss
// Returns XML content type
//
// Query published posts, ordered by created_at desc, limit 20
// For each post, include:
//   - title: post.Title
//   - link: "https://myblog.com/blog/" + post.Slug
//   - description: post.Excerpt (or first 200 chars of content)
//   - pubDate: post.CreatedAt in RFC 822 format
//   - category: post.Category.Name
//   - author: post.Author.Name`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The XML structure follows the RSS 2.0 specification:
          </p>

          <CodeBlock filename="RSS XML Structure">
{`<!-- RSS 2.0 Feed -->
<!-- channel: title, link, description, language -->
<!-- For each published post: -->
<!--   item: title, link, description, pubDate, category, author -->`}
          </CodeBlock>

          <Challenge number={8} title="Plan Your RSS Feed">
            <p>Plan an RSS feed for your blog. Which fields from the Post model would you include?
            How many posts should the feed contain? Should draft posts appear in the feed? Write
            out the fields you{"'"}d map: Post.Title maps to RSS title, Post.Slug maps to RSS link, etc.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Admin Content Management ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Admin Content Management</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The admin panel is where content creators spend their time. Grit generates a full admin
            interface for each resource — but a good CMS needs more than basic CRUD. It needs
            workflow tools for managing content at scale.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key admin features for a blog CMS:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Post list with filters</strong> — filter by Published, Draft, or All. Filter by category. Search by title.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Bulk actions</strong> — select multiple posts and publish, unpublish, or delete them at once</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Comment moderation queue</strong> — list of pending comments with approve/reject buttons</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Dashboard stats</strong> — total posts, published vs. draft, comments pending, top categories</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Quick edit</strong> — change title, category, or published status without opening the full form</li>
          </ul>

          <Tip>
            The admin panel{"'"}s DataTable component supports sorting and filtering out of the box.
            Click column headers to sort, use the filter inputs to narrow results. The generated
            admin pages already have this wired up.
          </Tip>

          <Challenge number={9} title="Full Admin Workflow">
            <p>In the admin panel, complete this workflow: (1) Create a new post as a draft,
            (2) Edit it and add rich text content, (3) Assign a category and 3 tags, (4) Publish it,
            (5) Unpublish it, (6) Delete it. Verify each step works correctly.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Media Library ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Media Library</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Blog posts need images. Grit{"'"}s file storage system (backed by S3-compatible storage
            like MinIO, Cloudflare R2, or AWS S3) handles uploads, and the featured image field on
            posts references uploaded files.
          </p>

          <CodeBlock filename="Upload Flow">
{`# 1. Upload an image
POST /api/uploads
Content-Type: multipart/form-data
# file: my-hero-image.jpg

# Response: { "data": { "id": 1, "url": "/uploads/abc123.jpg" } }

# 2. Use the URL as featured_image when creating/updating a post
PATCH /api/posts/5
{ "featured_image": "/uploads/abc123.jpg" }`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The media library in the admin panel shows all uploaded files as a grid of thumbnails.
            Click an image to copy its URL, then paste it into the featured image field or the rich
            text editor.
          </p>

          <Note>
            Docker Compose includes MinIO for local development — an S3-compatible object store
            that runs on your machine. In production, switch to Cloudflare R2 or AWS S3 by
            changing environment variables. No code changes needed.
          </Note>

          <Challenge number={10} title="Upload and Attach Images">
            <p>Upload 5 images using the API or admin panel. Use them as featured images for 5
            different blog posts. Verify the images appear when you fetch the posts via the API.
            Check: does the <Code>featured_image</Code> field contain the correct URL?</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s everything you learned in this course:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'A CMS separates content creation (admin) from content display (public site)',
              'Categories and tags organize content — categories are hierarchical, tags are flat',
              'The slug field type auto-generates URL-friendly strings for SEO',
              'Rich text content stores HTML from the editor, rendered directly on the frontend',
              'Comments with moderation prevent spam — approved:bool controls visibility',
              'List endpoints support filtering, sorting, pagination, and search out of the box',
              'RSS feeds distribute your content to feed readers as XML',
              'File uploads handle featured images via S3-compatible storage',
              'The admin panel provides content workflow: create, edit, publish, unpublish, delete',
              'belongs_to and many_to_many relationships connect posts to categories, tags, and users',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={11} title="Build the Complete Blog (Part 1)">
            <p>Create the full dataset: 5 categories, 10 tags, and 20 blog posts. Make 10 posts
            published and 10 drafts. Assign each post a category and at least 2 tags. Give at
            least 3 posts featured images.</p>
          </Challenge>

          <Challenge number={12} title="Build the Complete Blog (Part 2)">
            <p>Add comments to 5 posts — at least 3 comments per post. Approve some, leave others
            pending. Query the public API to verify only approved comments appear. Check the admin
            panel{"'"}s moderation queue.</p>
          </Challenge>

          <Challenge number={13} title="Build the Complete Blog (Part 3)">
            <p>Test the full public API flow: list published posts with pagination, filter by
            category, filter by tag, fetch a single post by slug, view approved comments. Every
            query should return the expected results.</p>
          </Challenge>

          <Challenge number={14} title="Build the Complete Blog (Part 4)">
            <p>Test the full admin workflow: create a post as draft, write rich text content, upload
            a featured image, assign category and tags, publish, receive comments, approve comments,
            verify on the public API. This is the complete CMS cycle — creation to publication to engagement.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
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
