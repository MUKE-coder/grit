import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'

export default function NamingConventionsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Core Concepts</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Naming Conventions
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit enforces consistent naming conventions across the entire stack. When you
                generate a resource, the CLI automatically converts your resource name into the
                correct case for each context -- PascalCase for Go structs, snake_case for database
                columns, kebab-case for TypeScript files, and more.
              </p>
            </div>

            <div className="prose-grit">
              {/* Master Naming Table */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Complete Naming Table
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  This table shows the naming convention used for every context in a Grit project.
                  The code generator handles all conversions automatically.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden overflow-x-auto">
                  <table className="w-full text-sm min-w-[600px]">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Context</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Convention</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Example</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go file name</td>
                        <td className="px-4 py-2.5 font-mono text-xs">snake_case.go</td>
                        <td className="px-4 py-2.5 font-mono text-xs">blog_post.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go struct</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PascalCase</td>
                        <td className="px-4 py-2.5 font-mono text-xs">BlogPost</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go struct fields</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PascalCase</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PublishDate</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go handler struct</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PascalCaseHandler</td>
                        <td className="px-4 py-2.5 font-mono text-xs">BlogPostHandler</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go handler variable</td>
                        <td className="px-4 py-2.5 font-mono text-xs">camelCaseHandler</td>
                        <td className="px-4 py-2.5 font-mono text-xs">blogPostHandler</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go service struct</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PascalCaseService</td>
                        <td className="px-4 py-2.5 font-mono text-xs">BlogPostService</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">JSON field</td>
                        <td className="px-4 py-2.5 font-mono text-xs">snake_case</td>
                        <td className="px-4 py-2.5 font-mono text-xs">publish_date</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Database table</td>
                        <td className="px-4 py-2.5 font-mono text-xs">plural_snake_case</td>
                        <td className="px-4 py-2.5 font-mono text-xs">blog_posts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Database column</td>
                        <td className="px-4 py-2.5 font-mono text-xs">snake_case</td>
                        <td className="px-4 py-2.5 font-mono text-xs">publish_date</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">API route</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/plural_snake</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/blog_posts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">TypeScript file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">kebab-case.ts</td>
                        <td className="px-4 py-2.5 font-mono text-xs">blog-post.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">TypeScript interface</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PascalCase</td>
                        <td className="px-4 py-2.5 font-mono text-xs">BlogPost</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">React Query hook file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">use-plural-kebab.ts</td>
                        <td className="px-4 py-2.5 font-mono text-xs">use-blog-posts.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">React Query hook name</td>
                        <td className="px-4 py-2.5 font-mono text-xs">usePluralPascal</td>
                        <td className="px-4 py-2.5 font-mono text-xs">useBlogPosts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Zod schema</td>
                        <td className="px-4 py-2.5 font-mono text-xs">CreatePascalSchema</td>
                        <td className="px-4 py-2.5 font-mono text-xs">CreateBlogPostSchema</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Resource slug</td>
                        <td className="px-4 py-2.5 font-mono text-xs">plural-kebab-case</td>
                        <td className="px-4 py-2.5 font-mono text-xs">blog-posts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Resource variable</td>
                        <td className="px-4 py-2.5 font-mono text-xs">camelCaseResource</td>
                        <td className="px-4 py-2.5 font-mono text-xs">blogPostResource</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Admin page route</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/resources/plural-kebab</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/resources/blog-posts</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 text-xs">API route constant</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PLURAL_UPPER</td>
                        <td className="px-4 py-2.5 font-mono text-xs">BLOG_POSTS</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* How the Generator Handles Names */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How the Generator Handles Names
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The generator builds a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Names</code> struct
                  with all naming variants from your resource name. You can provide the name in
                  PascalCase, snake_case, or kebab-case -- the generator normalizes it to
                  PascalCase first, then derives all other forms.
                </p>
                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">internal naming struct</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{`type Names struct {
    Pascal       string  // BlogPost
    Camel        string  // blogPost
    Snake        string  // blog_post
    Kebab        string  // blog-post
    Lower        string  // blogpost
    Plural       string  // blog_posts
    PluralPascal string  // BlogPosts
    PluralSnake  string  // blog_posts
    PluralKebab  string  // blog-posts
}`}</pre>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Conversion Flow
                </h3>
                <div className="space-y-3">
                  {[
                    { step: '1', title: 'Input normalization', desc: 'The raw input is converted to PascalCase. Input "blog_post", "blog-post", or "BlogPost" all become "BlogPost".' },
                    { step: '2', title: 'Snake case derivation', desc: 'PascalCase is split at uppercase boundaries. "BlogPost" becomes "blog_post". This is used for file names, JSON fields, and DB columns.' },
                    { step: '3', title: 'Kebab case derivation', desc: 'Snake case with underscores replaced by hyphens. "blog_post" becomes "blog-post". Used for TypeScript files and resource slugs.' },
                    { step: '4', title: 'Camel case derivation', desc: 'PascalCase with the first letter lowered. "BlogPost" becomes "blogPost". Used for handler variables and resource definitions.' },
                    { step: '5', title: 'Pluralization', desc: 'The snake_case form is pluralized using English rules. "blog_post" becomes "blog_posts". This is used for API routes, table names, and hook names.' },
                  ].map((item) => (
                    <div key={item.step} className="flex gap-3 p-3 rounded-lg border border-border/30 bg-card/30">
                      <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-xs font-mono font-semibold text-primary">
                        {item.step}
                      </div>
                      <div>
                        <p className="text-sm font-medium">{item.title}</p>
                        <p className="text-xs text-muted-foreground/70 mt-0.5">{item.desc}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Full Examples */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Full Name Derivation Examples
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Here are complete examples showing how a resource name maps to every file
                  and identifier in the project:
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  Example: &quot;Post&quot;
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden overflow-x-auto">
                  <table className="w-full text-sm min-w-[500px]">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Usage</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Value</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go model file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">models/post.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go struct</td>
                        <td className="px-4 py-2.5 font-mono text-xs">type Post struct</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go handler</td>
                        <td className="px-4 py-2.5 font-mono text-xs">PostHandler / postHandler</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">DB table</td>
                        <td className="px-4 py-2.5 font-mono text-xs">posts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">API routes</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/posts, /api/posts/:id</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">TS type file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">types/post.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Zod schema file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">schemas/post.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">React hooks</td>
                        <td className="px-4 py-2.5 font-mono text-xs">use-posts.ts / usePosts()</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Admin page</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/resources/posts/page.tsx</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 text-xs">Resource definition</td>
                        <td className="px-4 py-2.5 font-mono text-xs">resources/posts.ts / postResource</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Example: &quot;BlogPost&quot;
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden overflow-x-auto">
                  <table className="w-full text-sm min-w-[500px]">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Usage</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Value</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go model file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">models/blog_post.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go struct</td>
                        <td className="px-4 py-2.5 font-mono text-xs">type BlogPost struct</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go handler</td>
                        <td className="px-4 py-2.5 font-mono text-xs">BlogPostHandler / blogPostHandler</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">DB table</td>
                        <td className="px-4 py-2.5 font-mono text-xs">blog_posts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">API routes</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/blog_posts, /api/blog_posts/:id</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">TS type file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">types/blog-post.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Zod schema file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">schemas/blog-post.ts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">React hooks</td>
                        <td className="px-4 py-2.5 font-mono text-xs">use-blog-posts.ts / useBlogPosts()</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Admin page</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/resources/blog-posts/page.tsx</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 text-xs">Resource definition</td>
                        <td className="px-4 py-2.5 font-mono text-xs">resources/blog-posts.ts / blogPostResource</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Example: &quot;Category&quot;
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden overflow-x-auto">
                  <table className="w-full text-sm min-w-[500px]">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Usage</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Value</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go model file</td>
                        <td className="px-4 py-2.5 font-mono text-xs">models/category.go</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Go struct</td>
                        <td className="px-4 py-2.5 font-mono text-xs">type Category struct</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">DB table</td>
                        <td className="px-4 py-2.5 font-mono text-xs">categories</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">API routes</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/api/categories, /api/categories/:id</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">React hooks</td>
                        <td className="px-4 py-2.5 font-mono text-xs">use-categories.ts / useCategories()</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 text-xs">Admin page</td>
                        <td className="px-4 py-2.5 font-mono text-xs">/resources/categories/page.tsx</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  Notice how &quot;Category&quot; correctly pluralizes to &quot;categories&quot; (y
                  preceded by a consonant becomes -ies).
                </p>
              </div>

              {/* Pluralization Rules */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Pluralization Rules
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The Grit CLI includes a built-in English pluralization engine that handles
                  regular and irregular forms. Pluralization is applied to the snake_case form
                  of the resource name.
                </p>

                <h3 className="text-xl font-semibold tracking-tight mt-6 mb-3">
                  Regular Rules
                </h3>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Rule</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Singular</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Plural</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Default: add -s</td>
                        <td className="px-4 py-2.5 font-mono text-xs">post</td>
                        <td className="px-4 py-2.5 font-mono text-xs">posts</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Ends in -y (consonant before): -ies</td>
                        <td className="px-4 py-2.5 font-mono text-xs">category</td>
                        <td className="px-4 py-2.5 font-mono text-xs">categories</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Ends in -y (vowel before): -ys</td>
                        <td className="px-4 py-2.5 font-mono text-xs">key</td>
                        <td className="px-4 py-2.5 font-mono text-xs">keys</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 text-xs">Ends in -s, -ss, -sh, -ch, -x, -z: -es</td>
                        <td className="px-4 py-2.5 font-mono text-xs">address</td>
                        <td className="px-4 py-2.5 font-mono text-xs">addresses</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 text-xs">Ends in -f or -fe: -ves</td>
                        <td className="px-4 py-2.5 font-mono text-xs">leaf</td>
                        <td className="px-4 py-2.5 font-mono text-xs">leaves</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-xl font-semibold tracking-tight mt-8 mb-3">
                  Irregular Plurals
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The generator handles these irregular English plurals:
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Singular</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Plural</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      {[
                        ['person', 'people'],
                        ['child', 'children'],
                        ['mouse', 'mice'],
                        ['goose', 'geese'],
                        ['man', 'men'],
                        ['woman', 'women'],
                        ['tooth', 'teeth'],
                        ['foot', 'feet'],
                        ['ox', 'oxen'],
                        ['datum', 'data'],
                        ['medium', 'media'],
                        ['index', 'indices'],
                        ['matrix', 'matrices'],
                        ['vertex', 'vertices'],
                        ['crisis', 'crises'],
                        ['axis', 'axes'],
                        ['analysis', 'analyses'],
                      ].map(([singular, plural]) => (
                        <tr key={singular} className="border-b border-border/20 last:border-b-0">
                          <td className="px-4 py-2 font-mono text-xs">{singular}</td>
                          <td className="px-4 py-2 font-mono text-xs">{plural}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Icon Guessing */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Automatic Icon Selection
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When generating an admin resource definition, the CLI automatically picks a
                  Lucide icon based on the resource name. It matches the lowercase resource name
                  against a dictionary of common domain terms. If no match is found, it defaults
                  to the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Database</code> icon.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden overflow-x-auto">
                  <table className="w-full text-sm min-w-[500px]">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Resource Name Contains</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Lucide Icon</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      {[
                        ['post, page, document', 'FileText'],
                        ['article, blog', 'Newspaper'],
                        ['comment', 'MessageSquare'],
                        ['category', 'FolderTree'],
                        ['tag', 'Tag'],
                        ['product', 'Package'],
                        ['order', 'ShoppingCart'],
                        ['invoice', 'Receipt'],
                        ['payment, subscription', 'CreditCard'],
                        ['customer', 'UserCircle'],
                        ['user', 'Users'],
                        ['project', 'Briefcase'],
                        ['task', 'CheckSquare'],
                        ['event', 'Calendar'],
                        ['file', 'File'],
                        ['image, media', 'Image'],
                        ['message, email', 'Mail'],
                        ['notification', 'Bell'],
                        ['setting', 'Settings'],
                        ['role', 'Shield'],
                        ['permission', 'Lock'],
                        ['team', 'UsersRound'],
                        ['company, organization', 'Building2'],
                        ['report', 'BarChart3'],
                        ['analytic', 'TrendingUp'],
                        ['log', 'ScrollText'],
                        ['review', 'Star'],
                        ['plan', 'Gem'],
                        ['coupon', 'Ticket'],
                        ['discount', 'Percent'],
                        ['shipping', 'Truck'],
                        ['address, location', 'MapPin'],
                        ['contact', 'Contact'],
                        ['lead', 'Target'],
                        ['deal', 'Handshake'],
                        ['pipeline', 'GitBranch'],
                        ['workflow', 'Workflow'],
                        ['template', 'LayoutTemplate'],
                        ['campaign', 'Megaphone'],
                        ['survey', 'ClipboardList'],
                        ['form', 'FormInput'],
                        ['question', 'HelpCircle'],
                        ['answer', 'MessageCircle'],
                        ['ticket, issue', 'Ticket / AlertCircle'],
                        ['bug', 'Bug'],
                        ['feature', 'Sparkles'],
                        ['release', 'Rocket'],
                        ['version', 'GitCommit'],
                        ['deploy', 'CloudUpload'],
                      ].map(([names, icon]) => (
                        <tr key={names} className="border-b border-border/20 last:border-b-0">
                          <td className="px-4 py-2 text-xs">{names}</td>
                          <td className="px-4 py-2 font-mono text-xs">{icon}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  The matching is substring-based. A resource named &quot;ProductCategory&quot; would
                  match &quot;product&quot; first and get the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Package</code> icon.
                  You can always change the icon in the generated resource definition file.
                </p>
              </div>

              {/* File Naming Conventions */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  File Naming Summary
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Different ecosystems use different naming conventions for files. Grit follows
                  the idiomatic convention for each language:
                </p>
                <div className="space-y-4">
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h3 className="text-sm font-semibold mb-1.5">Go files: snake_case</h3>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      All Go source files use snake_case. This is the idiomatic Go convention.
                      Examples: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">user_handler.go</code>,
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">blog_post.go</code>,
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">auth.go</code>
                    </p>
                  </div>
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h3 className="text-sm font-semibold mb-1.5">TypeScript utility files: kebab-case</h3>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      Non-component TypeScript files use kebab-case. This matches Next.js and modern
                      TypeScript conventions.
                      Examples: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">api-client.ts</code>,
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">use-posts.ts</code>,
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">blog-post.ts</code>
                    </p>
                  </div>
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h3 className="text-sm font-semibold mb-1.5">React components: PascalCase</h3>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      React component files use PascalCase to match the component name inside.
                      Examples: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">DataTable.tsx</code>,
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">StatsCard.tsx</code>,
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">AdminLayout.tsx</code>
                    </p>
                  </div>
                  <div className="p-4 rounded-lg border border-border/30 bg-card/30">
                    <h3 className="text-sm font-semibold mb-1.5">Next.js pages: lowercase</h3>
                    <p className="text-xs text-muted-foreground/70 leading-relaxed">
                      Next.js App Router requires directory names to be lowercase (since they map to URL paths).
                      Examples: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/resources/posts/page.tsx</code>,
                      <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/(auth)/login/page.tsx</code>
                    </p>
                  </div>
                </div>
              </div>

              {/* Column Label Generation */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Column Label Generation
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  When the generator creates admin DataTable columns and form fields, it
                  converts field names to human-readable labels by splitting PascalCase into
                  separate words.
                </p>
                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Field Name</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">PascalCase</th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Generated Label</th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">title</td>
                        <td className="px-4 py-2.5 font-mono text-xs">Title</td>
                        <td className="px-4 py-2.5 text-xs">Title</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">due_date</td>
                        <td className="px-4 py-2.5 font-mono text-xs">DueDate</td>
                        <td className="px-4 py-2.5 text-xs">Due Date</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">is_published</td>
                        <td className="px-4 py-2.5 font-mono text-xs">IsPublished</td>
                        <td className="px-4 py-2.5 text-xs">Is Published</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">author_id</td>
                        <td className="px-4 py-2.5 font-mono text-xs">AuthorId</td>
                        <td className="px-4 py-2.5 text-xs">Author Id</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">total_amount</td>
                        <td className="px-4 py-2.5 font-mono text-xs">TotalAmount</td>
                        <td className="px-4 py-2.5 text-xs">Total Amount</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <p className="text-sm text-muted-foreground/60 mt-3">
                  The label splitting algorithm finds uppercase letter boundaries in PascalCase
                  and inserts spaces. You can always customize labels in the generated resource
                  definition file after generation.
                </p>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/type-system" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Type System
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/concepts/styles" className="gap-1.5">
                  Style Variants
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
