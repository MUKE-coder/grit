import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/mobile/first-app')

export default function MobileFirstAppPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Mobile</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Your First Mobile App
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Build a small Notes app end to end: scaffold the project, generate a{' '}
                <code>Note</code> resource, and run the generated screens on your phone. In
                about ten minutes you go from an empty folder to a native CRUD app backed by a
                Go API.
              </p>
            </div>

            <div className="mb-10 rounded-lg border border-primary/20 bg-primary/5 p-4">
              <p className="text-sm text-muted-foreground leading-relaxed">
                <strong className="text-foreground">Prefer a narrative walkthrough?</strong> There&apos;s a
                story-form companion on the blog:{' '}
                <Link href="/blog/build-mobile-app-with-grit" className="text-primary hover:underline">Build a mobile app with Grit &rarr;</Link>
              </p>
            </div>

            {/* Step 1 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  1
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Scaffold the Project</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Create a mobile project. This gives you a Go API in <code>apps/api</code>, an
                  Expo app in <code>apps/expo</code>, and a shared types package — all wired
                  together in a Turborepo monorepo.
                </p>
              </div>
              <CodeBlock terminal code={`grit new notes-app --mobile
cd notes-app`} className="mb-0 glow-purple-sm" />
            </div>

            {/* Step 2 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  2
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Generate the Note Resource</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  One command generates the full stack for a resource — Go model, service,
                  handler, shared Zod schema and TypeScript types, <strong>and</strong> the
                  Expo screens. Give <code>Note</code> a title, a body, and a pinned flag:
                </p>
              </div>
              <CodeBlock terminal code={`grit generate resource Note --fields "title:string,body:text,pinned:bool"`} className="mb-4 glow-purple-sm" />
              <div className="prose-grit mb-4">
                <p>Grit writes and wires these files (nothing to register by hand):</p>
              </div>
              <CodeBlock language="text" filename="generated for Note" code={`apps/api/internal/models/note.go              # GORM model + AutoMigrate
apps/api/internal/services/note.go            # CRUD + pagination
apps/api/internal/handlers/note.go            # Gin handlers (+ note_import.go)
apps/api/internal/routes/...                  # routes registered
packages/shared/schemas/note.ts               # Zod schema
packages/shared/types/note.ts                 # TypeScript type

apps/expo/hooks/use-notes.ts                  # typed React Query hook
apps/expo/app/notes/index.tsx                 # list screen
apps/expo/app/notes/[id].tsx                  # detail screen
apps/expo/app/notes/new.tsx                   # create screen
apps/expo/app/notes/edit/[id].tsx             # edit screen
apps/expo/components/resource-forms/notes-form.tsx   # shared form

# and a "Notes" card is injected into app/(tabs)/explore.tsx (the More tab)`} className="mb-0" />
            </div>

            {/* Step 3 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  3
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Migrate and Seed</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Create the <code>notes</code> table and seed the built-in admin account so you
                  have credentials to log in with.
                </p>
              </div>
              <CodeBlock terminal code={`grit migrate
grit seed`} className="mb-0" />
              <div className="prose-grit mt-4">
                <p>
                  The seeder creates <code>admin@example.com</code> / <code>admin123</code>. IDs
                  are UUID strings and roles are uppercase (<code>ADMIN</code>, <code>EDITOR</code>,{' '}
                  <code>USER</code>) — the admin account is an <code>ADMIN</code>.
                </p>
              </div>
            </div>

            {/* Step 4 */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  4
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Run It</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Start the API in one terminal and Expo in another, then scan the QR code with
                  Expo Go on your phone (phone and computer on the same Wi-Fi).
                </p>
              </div>
              <CodeBlock terminal code={`# terminal 1
grit start server        # Go API on :8080

# terminal 2
grit start expo          # Metro + QR code`} className="mb-4" />
              <div className="prose-grit mb-0">
                <p>
                  Log in with the seeded admin. The app opens on the tab bar — head to the{' '}
                  <strong>More</strong> tab and tap the <strong>Notes</strong> card that the
                  generator added. You are now on the generated list screen.
                </p>
              </div>
            </div>

            {/* Step 5: Tour */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  5
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Tour the Generated Screens</h2>
              </div>
              <div className="prose-grit mb-4">
                <p>
                  Every screen is real and navigable — a complete create / read / update /
                  delete flow the moment you generate the resource. Here is what each one does:
                </p>
              </div>

              <div className="space-y-4 mb-6">
                {[
                  {
                    title: 'List — app/notes/index.tsx',
                    body: 'A scrollable table with a search box, sortable columns, pull-to-refresh, and infinite scroll (it fetches the next page as you reach the end). A + button opens a bottom sheet with the create form. Export and import buttons sit in the header.',
                  },
                  {
                    title: 'Detail — app/notes/[id].tsx',
                    body: 'Tap a row to open the record. Every field is shown as a labelled row, with Edit and Delete actions at the bottom (delete asks for confirmation first).',
                  },
                  {
                    title: 'Create / Edit — new.tsx & edit/[id].tsx',
                    body: 'Both render the shared NotesForm. title becomes a text input, body a multi-line textarea, and pinned a native toggle. Edit pre-fills the form from the record.',
                  },
                ].map((item) => (
                  <div key={item.title} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <p className="text-[15px] font-semibold mb-1 font-mono">{item.title}</p>
                    <p className="text-sm text-muted-foreground leading-relaxed">{item.body}</p>
                  </div>
                ))}
              </div>

              <div className="prose-grit mb-4">
                <p>
                  The list, detail, and forms all talk to the API through one typed hook,{' '}
                  <code>use-notes.ts</code>. It exposes an infinite-scroll list query, a
                  single-record query, and create / update / delete mutations that invalidate
                  the cache on success — so a create instantly shows up in the list:
                </p>
              </div>
              <CodeBlock language="typescript" filename="apps/expo/hooks/use-notes.ts (excerpt)" code={`export function useNotes(search = "", filters = {}, sortBy = "created_at", sortOrder = "desc", pageSize = 20) {
  return useInfiniteQuery({
    queryKey: ["notes", { search, filters, sortBy, sortOrder, pageSize }],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => { /* GET /notes?page=… */ },
    getNextPageParam: (last) =>
      last.meta.page < last.meta.pages ? last.meta.page + 1 : undefined,
  });
}

export function useCreateNote() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input) => api.post("/notes", input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["notes"] }),
  });
}`} className="mb-0" />
            </div>

            {/* Step 6: Make it yours */}
            <div className="mb-10">
              <div className="flex items-center gap-3 mb-4">
                <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary shrink-0">
                  6
                </div>
                <h2 className="text-2xl font-semibold tracking-tight">Make It Yours</h2>
              </div>
              <div className="prose-grit mb-0">
                <p>
                  The generated screens are a starting point, not a cage. Restyle the list rows,
                  swap components, or rework the form in{' '}
                  <code>components/resource-forms/notes-form.tsx</code>. Because the whole stack
                  shares one set of types from <code>packages/shared</code>, a change to the Go
                  model flows to the client — run <code>grit sync</code> to regenerate the
                  TypeScript types after editing a model, and the compiler tells you exactly what
                  to update.
                </p>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/mobile/getting-started" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Getting Started
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/mobile/resource-generation" className="gap-1.5">
                  Resource Generation
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
