'use client'

import { useState } from 'react'
import { Check, Copy, Database, HardDrive } from 'lucide-react'

/*
 * A full-stack app, in the commands that build it.
 *
 * Two paths, because the honest answer to "what do I need installed" has two
 * branches and hiding one of them costs somebody an afternoon. With Docker is
 * the default and what the rest of the docs assume. Without it is one extra
 * flag: --db sqlite scaffolds a project that needs no servers at all, and since
 * v3.377.0 it also writes REDIS_URL= so the first boot does not open by failing
 * to reach a Redis nobody started.
 *
 * The app is the same either way: contacts in groups, each with a photo. Small
 * enough to finish inside two minutes, and big enough to involve a
 * relationship, a file upload and two admin screens, which is most of what a
 * CRUD framework is judged on.
 *
 * Every command here is real and in this order. A command that does not work is
 * worse than no command: the first person to paste it finds out, and then
 * doubts the rest of the site.
 */

interface Path {
  key: string
  label: string
  icon: typeof Database
  needs: string
  commands: string
  after: { title: string; body: string }[]
}

const SHARED_RESOURCES = `grit generate resource Group \
  --fields "name:string,description:text"
grit generate resource Contact \
  --fields "name:string,email:email,phone:tel,\
            photo:file:image,group:belongs_to:Group"`

const PATHS: Path[] = [
  {
    key: 'docker',
    label: 'With Docker',
    icon: Database,
    needs: 'Docker Desktop running. Postgres, Redis and MinIO come up with the project.',
    commands: `grit new contacts --triple --next
cd contacts
docker compose up -d
${SHARED_RESOURCES}
grit migrate
grit seed
grit start`,
    after: [
      { title: 'localhost:3000', body: 'The web app' },
      { title: 'localhost:3001', body: 'The admin panel, with a screen for contacts and one for groups' },
      { title: 'localhost:8080/docs', body: 'The API reference, generated from the routes' },
    ],
  },
  {
    key: 'nodocker',
    label: 'Without Docker',
    icon: HardDrive,
    needs: 'Nothing but Go and Node. SQLite is a file, contact photos are written to disk and served back, and caching and background jobs stay off until you want them.',
    commands: `grit new contacts --triple --next --db sqlite
cd contacts
${SHARED_RESOURCES}
grit migrate
grit seed
grit start`,
    after: [
      { title: 'localhost:3001', body: 'The admin panel, with a screen for contacts and one for groups' },
      { title: 'app.db', body: 'The whole database, as one file you can delete and start over' },
      { title: 'storage/app', body: 'Contact photos, on disk: the original, a converted JPEG and a thumbnail' },
    ],
  },
]

export function QuickStart() {
  const [active, setActive] = useState(PATHS[0].key)
  const [copied, setCopied] = useState(false)
  const path = PATHS.find((p) => p.key === active) ?? PATHS[0]

  async function copyAll() {
    try {
      await navigator.clipboard.writeText(path.commands)
      setCopied(true)
      setTimeout(() => setCopied(false), 1800)
    } catch {
      // A browser that refuses the clipboard still shows the commands to read.
    }
  }

  return (
    <section aria-labelledby="quick-start" className="rounded-2xl border border-primary/25 bg-primary/[0.04] p-6">
      <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
        <h2 id="quick-start" className="sr-only">
          Quick start
        </h2>
        <div role="tablist" aria-label="Whether to use Docker" className="inline-flex rounded-lg border border-border bg-background p-1">
          {PATHS.map((p) => {
            const selected = p.key === active
            return (
              <button
                key={p.key}
                type="button"
                role="tab"
                aria-selected={selected}
                onClick={() => setActive(p.key)}
                className={`inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors ${
                  selected ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                <p.icon className="h-3.5 w-3.5" />
                {p.label}
              </button>
            )
          })}
        </div>
        <button
          type="button"
          onClick={() => void copyAll()}
          className="inline-flex items-center gap-1.5 rounded-lg border border-border bg-background px-3 py-1.5 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground"
        >
          {copied ? <Check className="h-3.5 w-3.5 text-emerald-500" /> : <Copy className="h-3.5 w-3.5" />}
          {copied ? 'Copied' : 'Copy all'}
        </button>
      </div>

      <p className="mb-4 text-sm leading-relaxed text-muted-foreground">{path.needs}</p>

      <pre className="overflow-x-auto rounded-xl border border-border bg-[#0d1117] p-4 text-[13px] leading-relaxed">
        <code className="font-mono text-slate-200">{path.commands}</code>
      </pre>

      <div className="mt-5 grid gap-3 sm:grid-cols-3">
        {path.after.map((item) => (
          <div key={item.title} className="rounded-xl border border-border/60 bg-card/40 p-4">
            <p className="font-mono text-[13px] font-semibold text-primary">{item.title}</p>
            <p className="mt-1 text-[12.5px] leading-relaxed text-muted-foreground">{item.body}</p>
          </div>
        ))}
      </div>

      <p className="mt-5 text-sm leading-relaxed text-muted-foreground">
        Sign in to the admin panel with <code className="font-mono text-foreground">admin@example.com</code> and{' '}
        <code className="font-mono text-foreground">admin123</code>, add a group, then add a contact with a photo.
        The photo works on either path. Without Docker there is no MinIO to write to, so the API notices
        and keeps the file on local disk instead: you get the original, a converted JPEG and a thumbnail
        under <code className="font-mono text-foreground">storage/app</code>, and the admin table shows it.
        Nothing above was hand-written: the two generate commands produced the models, migrations, API, Zod
        schemas, TypeScript types, React Query hooks and both admin screens.
      </p>
    </section>
  )
}
