import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tutorials/blog')

export default function TutorialBlogPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Tutorial
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Build a Blog
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Build a complete blogging platform with posts, categories, a
                published filter, and a custom frontend &mdash; all in under 30
                minutes. You will use Grit&apos;s code generator, customize Go
                handlers, define admin resources, and wire up a Next.js page to
                display published articles.
              </p>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-semibold tracking-tight mb-4">
                Prerequisites
              </h2>
              <ul className="space-y-2.5">
                {[
                  "Go 1.21+ installed",
                  "Node.js 18+ and pnpm installed",
                  "Docker and Docker Compose installed",
                  "Grit CLI installed globally (go install github.com/MUKE-coder/grit/cmd/grit@latest)",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-1">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            {/* ============================================================ */}
            {/* STEP 1 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                1
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Create the project
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Scaffold a new Grit monorepo called <code>myblog</code>. This
                  creates the Go API, Next.js web app, admin panel, shared
                  package, and Docker configuration in one shot.
                </p>

                <CodeBlock terminal code={`grit new myblog\ncd myblog`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Grit prints an ASCII art logo, creates the folder structure,
                  initializes
                  <code> go.mod</code>, runs <code>pnpm install</code>, and
                  prints the next steps. Your project is ready.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 2 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                2
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Start Docker services
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Spin up PostgreSQL, Redis, MinIO (local S3), and Mailhog.
                  These run in the background and persist data across restarts.
                </p>

                <CodeBlock terminal code="docker compose up -d" className="glow-purple-sm" />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 3 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                3
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Post resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Use the code generator to create a full-stack Post resource.
                  This generates the Go model, handler, service, Zod schema,
                  TypeScript types, React Query hooks, and admin page &mdash;
                  all wired together.
                </p>

                <CodeBlock terminal code={`grit generate resource Post --fields "title:string,slug:string:unique,content:text,excerpt:text,published:bool,views:int"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The generator creates these files:
                </p>

                <CodeBlock language="bash" filename="generated files" code={`apps/api/internal/models/post.go        # GORM model
apps/api/internal/handlers/post.go      # CRUD handler
apps/api/internal/services/post.go      # Business logic
packages/shared/schemas/post.ts         # Zod validation
packages/shared/types/post.ts           # TypeScript types
apps/web/hooks/use-posts.ts             # React Query hooks (web)
apps/admin/hooks/use-posts.ts           # React Query hooks (admin)
apps/admin/app/resources/posts/page.tsx # Admin page
apps/admin/resources/posts.ts           # Resource definition`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Here is the generated Go model:
                </p>

                <CodeBlock filename="apps/api/internal/models/post.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Title     string         \`gorm:"size:255;not null" json:"title" binding:"required"\`
    Slug      string         \`gorm:"size:255;uniqueIndex;not null" json:"slug" binding:"required"\`
    Content   string         \`gorm:"type:text" json:"content"\`
    Excerpt   string         \`gorm:"type:text" json:"excerpt"\`
    Published bool           \`gorm:"default:false" json:"published"\`
    Views     int            \`gorm:"default:0" json:"views"\`
    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 4 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                4
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Generate the Category resource
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Every blog needs categories. Generate a Category resource with
                  a name, slug, and description.
                </p>

                <CodeBlock terminal code={`grit generate resource Category --fields "name:string:unique,slug:string:unique,description:text"`} className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  The generated Category model:
                </p>

                <CodeBlock filename="apps/api/internal/models/category.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Category struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255;uniqueIndex;not null" json:"name" binding:"required"\`
    Slug        string         \`gorm:"size:255;uniqueIndex;not null" json:"slug" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 5 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                5
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Add a relationship &mdash; Post belongs to Category
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  A post should belong to a category. Open the Post model and
                  add a<code> CategoryID</code> foreign key and a{" "}
                  <code>Category</code> relation field. Also add a{" "}
                  <code>Posts</code> slice to the Category model for the inverse
                  relationship.
                </p>

                <CodeBlock filename="apps/api/internal/models/post.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Post struct {
    ID         uint           \`gorm:"primarykey" json:"id"\`
    Title      string         \`gorm:"size:255;not null" json:"title" binding:"required"\`
    Slug       string         \`gorm:"size:255;uniqueIndex;not null" json:"slug" binding:"required"\`
    Content    string         \`gorm:"type:text" json:"content"\`
    Excerpt    string         \`gorm:"type:text" json:"excerpt"\`
    Published  bool           \`gorm:"default:false" json:"published"\`
    Views      int            \`gorm:"default:0" json:"views"\`
    CategoryID uint           \`gorm:"index;not null" json:"category_id" binding:"required"\`
    Category   Category       \`gorm:"foreignKey:CategoryID" json:"category,omitempty"\`
    CreatedAt  time.Time      \`json:"created_at"\`
    UpdatedAt  time.Time      \`json:"updated_at"\`
    DeletedAt  gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />

                <CodeBlock filename="apps/api/internal/models/category.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type Category struct {
    ID          uint           \`gorm:"primarykey" json:"id"\`
    Name        string         \`gorm:"size:255;uniqueIndex;not null" json:"name" binding:"required"\`
    Slug        string         \`gorm:"size:255;uniqueIndex;not null" json:"slug" binding:"required"\`
    Description string         \`gorm:"type:text" json:"description"\`
    Posts       []Post         \`gorm:"foreignKey:CategoryID" json:"posts,omitempty"\`
    CreatedAt   time.Time      \`json:"created_at"\`
    UpdatedAt   time.Time      \`json:"updated_at"\`
    DeletedAt   gorm.DeletedAt \`gorm:"index" json:"deleted_at,omitempty"\`
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Now sync the types to TypeScript so the frontend knows about
                  the relationship:
                </p>

                <CodeBlock terminal code="grit sync" className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  The <code>grit sync</code> command reads every Go model in{" "}
                  <code>internal/models/</code>, parses the structs with Go AST,
                  and regenerates the TypeScript types and Zod schemas in{" "}
                  <code>packages/shared/</code>. The Post type now includes{" "}
                  <code>category_id</code> and an optional <code>category</code>{" "}
                  object.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 6 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                6
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Customize the Post handler to preload Category
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  By default, the generated handler does not preload
                  relationships. Open the Post handler and add a{" "}
                  <code>Preload(&quot;Category&quot;)</code> call so that every
                  post response includes its category data.
                </p>

                <CodeBlock filename="apps/api/internal/services/post.go &mdash; GetAll method" code={`func (s *PostService) GetAll(page, pageSize int, sort, order, search string) ([]models.Post, int64, error) {
    var posts []models.Post
    var total int64

    query := s.db.Model(&models.Post{})

    // Search across title and excerpt
    if search != "" {
        query = query.Where("title ILIKE ? OR excerpt ILIKE ?",
            "%"+search+"%", "%"+search+"%")
    }

    // Count total before pagination
    query.Count(&total)

    // Sorting
    if sort != "" {
        direction := "ASC"
        if order == "desc" {
            direction = "DESC"
        }
        query = query.Order(sort + " " + direction)
    } else {
        query = query.Order("created_at DESC")
    }

    // Preload the Category relationship
    query = query.Preload("Category")

    // Pagination
    offset := (page - 1) * pageSize
    result := query.Offset(offset).Limit(pageSize).Find(&posts)

    return posts, total, result.Error
}

func (s *PostService) GetByID(id uint) (*models.Post, error) {
    var post models.Post
    // Preload Category when fetching a single post
    result := s.db.Preload("Category").First(&post, id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &post, nil
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed">
                  Now every API response for posts includes the nested{" "}
                  <code>category</code> object with its <code>name</code>,{" "}
                  <code>slug</code>, and <code>description</code>.
                </p>
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 7 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                7
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Add a custom endpoint &mdash; published posts only
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The public frontend should only show published posts. Add a
                  new handler method and register it as a public route (no auth
                  required).
                </p>

                <CodeBlock filename="apps/api/internal/handlers/post.go &mdash; add this method" code={`// GetPublished returns only published posts for the public frontend.
func (h *PostHandler) GetPublished(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

    var posts []models.Post
    var total int64

    query := h.db.Model(&models.Post{}).Where("published = ?", true)
    query.Count(&total)

    query.Preload("Category").
        Order("created_at DESC").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&posts)

    pages := int(total) / pageSize
    if int(total)%pageSize != 0 {
        pages++
    }

    c.JSON(http.StatusOK, gin.H{
        "data": posts,
        "meta": gin.H{
            "total":     total,
            "page":      page,
            "page_size": pageSize,
            "pages":     pages,
        },
    })
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Register the new route in <code>routes.go</code>. Place it
                  outside the auth middleware group so anyone can access it:
                </p>

                <CodeBlock filename="apps/api/internal/routes/routes.go &mdash; add this route" code={`// Public routes (no authentication required)
public := router.Group("/api")
{
    // ... existing public routes (auth, health) ...

    // Published blog posts — public access
    public.GET("/posts/published", postHandler.GetPublished)
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 8 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                8
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Customize the admin resource definition
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  The generated admin resource is functional but generic.
                  Let&apos;s make it blog-specific by adding a category filter,
                  a published badge, relative dates, a view count column, and a
                  category selector in the form.
                </p>

                <CodeBlock language="typescript" filename="apps/admin/resources/posts.ts" code={`import { defineResource } from '@grit/admin'

export default defineResource({
  name: 'Post',
  endpoint: '/api/posts',
  icon: 'FileText',

  table: {
    columns: [
      { key: 'id', label: 'ID', sortable: true },
      { key: 'title', label: 'Title', sortable: true, searchable: true },
      { key: 'slug', label: 'Slug', sortable: true },
      { key: 'category.name', label: 'Category', relation: 'category' },
      { key: 'published', label: 'Status', badge: {
        true: { color: 'green', label: 'Published' },
        false: { color: 'yellow', label: 'Draft' },
      }},
      { key: 'views', label: 'Views', sortable: true },
      { key: 'created_at', label: 'Created', format: 'relative' },
    ],
    filters: [
      { key: 'published', type: 'select', options: [
        { label: 'Published', value: 'true' },
        { label: 'Draft', value: 'false' },
      ]},
      { key: 'category_id', type: 'select', resource: 'categories',
        displayKey: 'name', label: 'Category' },
      { key: 'created_at', type: 'date-range' },
    ],
    actions: ['create', 'edit', 'delete', 'export'],
    bulkActions: ['delete', 'export'],
  },

  form: {
    fields: [
      { key: 'title', label: 'Title', type: 'text', required: true },
      { key: 'slug', label: 'Slug', type: 'text', required: true },
      { key: 'category_id', label: 'Category', type: 'relation',
        resource: 'categories', displayKey: 'name', required: true },
      { key: 'excerpt', label: 'Excerpt', type: 'textarea' },
      { key: 'content', label: 'Content', type: 'richtext' },
      { key: 'published', label: 'Published', type: 'toggle', default: false },
    ],
  },
})`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 9 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                9
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Update the web frontend to display posts
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Now build the public-facing blog page. Create a custom React
                  Query hook that calls the <code>/api/posts/published</code>{" "}
                  endpoint, then build a page component that lists the posts
                  with their categories.
                </p>

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  First, add a hook for the published posts endpoint:
                </p>

                <CodeBlock language="typescript" filename="apps/web/hooks/use-published-posts.ts" code={`import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/lib/api-client'
import type { Post } from '@shared/types/post'
import type { PaginatedResponse } from '@shared/types/api'

interface UsePublishedPostsOptions {
  page?: number
  pageSize?: number
}

export function usePublishedPosts({ page = 1, pageSize = 10 }: UsePublishedPostsOptions = {}) {
  return useQuery<PaginatedResponse<Post>>({
    queryKey: ['posts', 'published', page, pageSize],
    queryFn: async () => {
      const { data } = await apiClient.get('/api/posts/published', {
        params: { page, page_size: pageSize },
      })
      return data
    },
  })
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Next, create the blog listing page:
                </p>

                <CodeBlock language="tsx" filename="apps/web/app/(dashboard)/blog/page.tsx" code={`'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePublishedPosts } from '@/hooks/use-published-posts'

export default function BlogPage() {
  const [page, setPage] = useState(1)
  const { data, isLoading } = usePublishedPosts({ page, pageSize: 12 })

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold tracking-tight">Blog</h1>
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="rounded-xl border border-border/40 bg-card/50 p-6 animate-pulse">
              <div className="h-4 w-20 rounded bg-accent/30 mb-3" />
              <div className="h-6 w-3/4 rounded bg-accent/30 mb-2" />
              <div className="h-4 w-full rounded bg-accent/30" />
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold tracking-tight">Blog</h1>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {data?.data.map((post) => (
          <Link
            key={post.id}
            href={\`/blog/\${post.slug}\`}
            className="group rounded-xl border border-border/40 bg-card/50 p-6
                       hover:border-primary/30 hover:bg-card/80 transition-all"
          >
            {post.category && (
              <span className="text-xs font-mono text-primary/70 mb-2 block">
                {post.category.name}
              </span>
            )}
            <h2 className="text-lg font-semibold mb-2 group-hover:text-primary transition-colors">
              {post.title}
            </h2>
            {post.excerpt && (
              <p className="text-sm text-muted-foreground/70 line-clamp-2">
                {post.excerpt}
              </p>
            )}
            <p className="text-xs text-muted-foreground/50 mt-3">
              {new Date(post.created_at).toLocaleDateString('en-US', {
                year: 'numeric', month: 'long', day: 'numeric',
              })}
            </p>
          </Link>
        ))}
      </div>

      {/* Pagination */}
      {data?.meta && data.meta.pages > 1 && (
        <div className="flex items-center justify-center gap-2 pt-4">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-3 py-1.5 text-sm rounded-md border border-border/40
                       hover:bg-accent/50 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <span className="text-sm text-muted-foreground">
            Page {page} of {data.meta.pages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(data.meta.pages, p + 1))}
            disabled={page === data.meta.pages}
            className="px-3 py-1.5 text-sm rounded-md border border-border/40
                       hover:bg-accent/50 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      )}
    </div>
  )
}`} />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-3">
                  Finally, create the single post page to display a full
                  article:
                </p>

                <CodeBlock language="tsx" filename="apps/web/app/(dashboard)/blog/[slug]/page.tsx" code={`'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/lib/api-client'
import type { Post } from '@shared/types/post'
import Link from 'next/link'
import { ArrowLeft } from 'lucide-react'

export default function BlogPostPage() {
  const params = useParams()
  const slug = params.slug as string

  const { data: post, isLoading } = useQuery<Post>({
    queryKey: ['posts', 'slug', slug],
    queryFn: async () => {
      const { data } = await apiClient.get(\`/api/posts/published?slug=\${slug}\`)
      // Find the matching post from the published list
      return data.data[0]
    },
  })

  if (isLoading) {
    return (
      <div className="max-w-2xl mx-auto animate-pulse space-y-4">
        <div className="h-8 w-3/4 rounded bg-accent/30" />
        <div className="h-4 w-1/4 rounded bg-accent/30" />
        <div className="space-y-2 pt-6">
          <div className="h-4 w-full rounded bg-accent/30" />
          <div className="h-4 w-5/6 rounded bg-accent/30" />
          <div className="h-4 w-4/6 rounded bg-accent/30" />
        </div>
      </div>
    )
  }

  if (!post) {
    return (
      <div className="max-w-2xl mx-auto text-center py-20">
        <h1 className="text-2xl font-bold mb-2">Post not found</h1>
        <p className="text-muted-foreground mb-4">This post may have been unpublished or deleted.</p>
        <Link href="/blog" className="text-primary hover:underline">Back to blog</Link>
      </div>
    )
  }

  return (
    <article className="max-w-2xl mx-auto">
      <Link href="/blog" className="flex items-center gap-1.5 text-sm text-muted-foreground/60
                                     hover:text-foreground transition-colors mb-6">
        <ArrowLeft className="h-3.5 w-3.5" />
        Back to blog
      </Link>

      {post.category && (
        <span className="text-xs font-mono text-primary/70 mb-3 block">
          {post.category.name}
        </span>
      )}
      <h1 className="text-3xl font-bold tracking-tight mb-3">{post.title}</h1>
      <p className="text-sm text-muted-foreground/60 mb-8">
        {new Date(post.created_at).toLocaleDateString('en-US', {
          year: 'numeric', month: 'long', day: 'numeric',
        })}
        {' '}&mdash;{' '}{post.views} views
      </p>

      <div className="prose prose-invert max-w-none" dangerouslySetInnerHTML={{ __html: post.content }} />
    </article>
  )
}`} />
              </div>
            </div>

            {/* ============================================================ */}
            {/* STEP 10 */}
            {/* ============================================================ */}
            <div className="relative pl-10 mb-12">
              <div className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                10
              </div>
              <div>
                <h3 className="text-xl font-semibold mb-2">
                  Run and test everything
                </h3>
                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Start the Go API, the web app, and the admin panel in
                  development mode. Grit uses Turborepo to run all services
                  concurrently.
                </p>

                <CodeBlock terminal code="grit dev" className="glow-purple-sm mb-4" />

                <p className="text-[13px] text-muted-foreground/70 leading-relaxed mb-4">
                  Open these URLs in your browser:
                </p>

                <ul className="space-y-2">
                  {[
                    {
                      url: "http://localhost:3000",
                      desc: "Web frontend — register and see the blog",
                    },
                    {
                      url: "http://localhost:3001",
                      desc: "Admin panel — manage posts and categories",
                    },
                    {
                      url: "http://localhost:8080/studio",
                      desc: "GORM Studio — browse your database",
                    },
                    {
                      url: "http://localhost:8080/api/posts/published",
                      desc: "Published posts API",
                    },
                  ].map((item) => (
                    <li
                      key={item.url}
                      className="flex items-start gap-2.5 text-[13px] text-muted-foreground/70"
                    >
                      <code className="text-primary/70 text-xs whitespace-nowrap">
                        {item.url}
                      </code>
                      <span>&mdash; {item.desc}</span>
                    </li>
                  ))}
                </ul>
              </div>
            </div>

            {/* ============================================================ */}
            {/* SUMMARY */}
            {/* ============================================================ */}
            <div className="mb-12 rounded-xl border border-border/40 bg-card/30 p-6">
              <h2 className="text-xl font-semibold tracking-tight mb-4">
                What you&apos;ve built
              </h2>
              <ul className="space-y-2.5">
                {[
                  "A full-stack blog with Go API and Next.js frontend",
                  "Post and Category resources generated with a single CLI command each",
                  "A BelongsTo relationship between Posts and Categories",
                  "A custom public endpoint that returns only published posts",
                  "An admin panel with sortable tables, category filters, and status badges",
                  "A blog listing page with pagination and article detail pages",
                  "Type-safe data fetching with React Query and shared Zod schemas",
                  "Docker-based PostgreSQL, Redis, MinIO, and Mailhog running locally",
                  "GORM Studio for visual database browsing",
                ].map((item) => (
                  <li
                    key={item}
                    className="flex items-start gap-2.5 text-[14px] text-muted-foreground"
                  >
                    <span className="text-primary mt-0.5">&#10003;</span>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            <div className="prose-grit mb-8">
              <h2>Next steps</h2>
              <p>
                Now that you have a working blog, here are some ideas to extend
                it:
              </p>
              <ul>
                <li>
                  <strong>Add tags</strong> &mdash; generate a <code>Tag</code>{" "}
                  resource and create a many-to-many relationship with posts
                  using a join table.
                </li>
                <li>
                  <strong>Markdown rendering</strong> &mdash; store content as
                  Markdown and render it on the frontend with{" "}
                  <code>react-markdown</code>.
                </li>
                <li>
                  <strong>Image uploads</strong> &mdash; use the Grit storage
                  service to upload cover images for each post.
                </li>
                <li>
                  <strong>RSS feed</strong> &mdash; add a custom Go handler that
                  returns an XML RSS feed of published posts.
                </li>
                <li>
                  <strong>SEO metadata</strong> &mdash; use Next.js{" "}
                  <code>generateMetadata</code> to set page titles and
                  descriptions dynamically.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tutorials/learn" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Learn Grit Step by Step
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/tutorials/saas" className="gap-1.5">
                  Build a SaaS
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
