/*
 * The data behind /llms.txt and /llms-full.txt.
 *
 * An agent asked to work in a Grit project has no way in: it can crawl the docs
 * site, which is 171 pages of TSX rendered to HTML, or it can guess. llms.txt is
 * the index it reads instead, and llms-full.txt is the text, both built from the
 * same files the site is built from so neither can drift from it.
 *
 * Nothing here is written twice. The page list is docsMetadata, which the pages
 * already use for their <title> and description. The CLI reference is
 * CLI_COMMANDS, which /docs/cli renders. The working guide is public/skill.md,
 * which is also what `grit init` writes into a project. A change to any of those
 * changes these two files on the next build.
 *
 * Format: llmstxt.org. An H1, a blockquote summary, prose, then H2 sections of
 * links, each "- [Title](url): description".
 */

import { docsMetadata } from '@/config/docs-metadata'
import { CLI_COMMANDS, CLI_CATEGORIES } from '@/config/cli-commands'
import { GRIT_VERSION } from '@/config/site'

const BASE = 'https://gritframework.dev'

/**
 * Section titles by the first path segment under /docs, in reading order.
 *
 * A path whose segment is not here still appears, under its segment title-cased,
 * at the end. That is deliberate: a new docs directory should show up in
 * llms.txt the day it is added, not the day somebody remembers this list.
 */
const SECTIONS: Array<[string, string]> = [
  ['', 'Start here'],
  ['getting-started', 'Getting started'],
  ['prerequisites', 'Prerequisites (Go, Next.js, Docker)'],
  ['concepts', 'Concepts'],
  ['cli', 'The CLI'],
  ['backend', 'Backend (Go, Gin, GORM)'],
  ['frontend', 'Frontend (Next.js, TanStack Router)'],
  ['admin', 'Admin panel'],
  ['batteries', 'Batteries (cache, jobs, mail, storage, AI)'],
  ['infrastructure', 'Infrastructure'],
  ['security', 'Security'],
  ['governance', 'Governance and compliance'],
  ['mobile', 'Mobile (Expo)'],
  ['desktop', 'Desktop (Wails)'],
  ['deployment', 'Deployment'],
  ['testing', 'Testing'],
  ['scaling', 'Scaling'],
  ['stability', 'Stability and upgrades'],
  ['plugins', 'Plugins'],
  ['design', 'Design system'],
  ['tech-kits', 'Tech kits'],
  ['tutorials', 'Tutorials'],
  ['benchmarks', 'Benchmarks'],
  ['learnings', 'Learnings'],
  ['ai-skill', 'Working with AI agents'],
  ['ai-workflows', 'AI agent workflows'],
  ['ai-integration', 'AI integration'],
  ['demo', 'The demo application'],
]

/**
 * Pages the Optional section at the end links by hand, so the grouped list does
 * not also give them a section of one.
 */
const LINKED_ELSEWHERE = new Set(['/docs/changelog'])

function sectionOf(path: string): string {
  const parts = path.replace(/^\/docs\/?/, '').split('/').filter(Boolean)
  return parts.length === 0 ? '' : parts[0]
}

function titleCase(segment: string): string {
  return segment
    .split('-')
    .map((w) => w.charAt(0).toUpperCase() + w.slice(1))
    .join(' ')
}

/** One line of a link section. */
function line(path: string): string {
  const page = docsMetadata[path]
  return `- [${page.title}](${BASE}${path}): ${page.description}`
}

/**
 * The docs pages grouped into sections, in the order SECTIONS gives and then
 * alphabetically for anything it does not name.
 */
function groupedPages(): Array<[string, string[]]> {
  const named = new Set(SECTIONS.map(([key]) => key))
  const buckets = new Map<string, string[]>()
  const counts = new Map<string, number>()
  for (const path of Object.keys(docsMetadata)) {
    if (!path.startsWith('/docs')) continue
    const key = sectionOf(path)
    counts.set(key, (counts.get(key) ?? 0) + 1)
  }
  for (const path of Object.keys(docsMetadata)) {
    if (!path.startsWith('/docs') || LINKED_ELSEWHERE.has(path)) continue
    let key = sectionOf(path)
    // A directory with one page in it and no name of its own is not a section.
    // /docs/start, /docs/stack-selector and /docs/versioning each drew a heading
    // with a single link under it, which is noise where the reader is a model
    // counting tokens.
    if (!named.has(key) && (counts.get(key) ?? 0) <= 1) key = ''
    const bucket = buckets.get(key)
    if (bucket) bucket.push(path)
    else buckets.set(key, [path])
  }

  const out: Array<[string, string[]]> = []
  for (const [key, title] of SECTIONS) {
    const paths = buckets.get(key)
    if (paths && paths.length > 0) out.push([title, paths.sort()])
  }
  for (const key of [...buckets.keys()].sort()) {
    if (named.has(key)) continue
    out.push([titleCase(key), (buckets.get(key) as string[]).sort()])
  }
  return out
}

const SUMMARY =
  'Grit is a full-stack meta-framework: a Go API (Gin + GORM), a Next.js or ' +
  'TanStack Router frontend, and a generated admin panel, in one monorepo with ' +
  'shared Zod schemas and TypeScript types.'

const ORIENTATION = `Grit is a code generator, not a runtime library. \`grit generate resource Post\`
writes a Go model, a service and a handler, a Zod schema, TypeScript types,
React Query hooks and an admin page, and injects the wiring into the router and
the registries. That generated code belongs to the project: read it, edit it,
delete it. Hand-writing those files instead of running the command produces
something that compiles and is missing the injections that make it reachable.

So the loop is: generate, then edit. A scaffolded project already has
authentication, roles and permissions, file storage, email, background jobs,
caching, audit logging, an admin panel and an OpenAPI reference wired together
before a line is written.`

/** /llms.txt: the index. */
export function llmsIndex(): string {
  const parts: string[] = []
  parts.push('# Grit')
  parts.push('')
  parts.push(`> ${SUMMARY}`)
  parts.push('')
  parts.push(ORIENTATION)
  parts.push('')
  parts.push(`Current CLI version: ${GRIT_VERSION}. Source: https://github.com/MUKE-coder/grit`)
  parts.push('')
  parts.push('## Read these first')
  parts.push('')
  parts.push(
    `- [The full text](${BASE}/llms-full.txt): this index, the working guide an agent should follow, and the complete CLI reference, as one file.`,
  )
  parts.push(
    `- [Agent skill](${BASE}/skill.md): the same working guide as a skill file, which \`grit init\` also writes into a project.`,
  )
  parts.push(
    `- [Stack selector](${BASE}/docs/stack-selector): which architecture to scaffold, with the exact command for each.`,
  )
  parts.push('')

  for (const [title, paths] of groupedPages()) {
    parts.push(`## ${title}`)
    parts.push('')
    for (const path of paths) parts.push(line(path))
    parts.push('')
  }

  parts.push('## Optional')
  parts.push('')
  parts.push(`- [Changelog](${BASE}/docs/changelog): every release, newest first.`)
  parts.push(`- [Playground](${BASE}/playground): the scaffolder's output without installing anything.`)
  parts.push(`- [OpenAPI reference](${BASE}/docs/backend/api-docs): a scaffolded API serves its own spec at /docs/openapi.json, which is the machine-readable form of everything it exposes.`)
  parts.push('')
  return parts.join('\n')
}

/** The CLI reference section of /llms-full.txt. */
export function llmsCliReference(): string {
  const parts: string[] = []
  parts.push('# CLI reference')
  parts.push('')
  parts.push(
    'Every command, by category. The file effects were captured from real runs against a\n' +
      'freshly scaffolded project rather than written from memory.',
  )
  parts.push('')

  for (const category of CLI_CATEGORIES) {
    const commands = CLI_COMMANDS.filter((c) => c.category === category)
    if (commands.length === 0) continue
    parts.push(`## ${category}`)
    parts.push('')
    for (const command of commands) {
      parts.push(`### \`${command.name}\``)
      parts.push('')
      parts.push(command.summary)
      parts.push('')
      parts.push('```bash')
      parts.push(command.example)
      parts.push('```')
      parts.push('')
      parts.push(command.purpose)
      if (command.flags && command.flags.length > 0) {
        parts.push('')
        parts.push('Flags:')
        parts.push('')
        for (const f of command.flags) parts.push(`- \`${f.flag}\`: ${f.desc}`)
      }
      if (command.notes && command.notes.length > 0) {
        parts.push('')
        parts.push('Worth knowing:')
        parts.push('')
        for (const n of command.notes) parts.push(`- ${n}`)
      }
      parts.push('')
    }
  }
  return parts.join('\n')
}
