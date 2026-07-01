import fs from 'node:fs'
import path from 'node:path'
import matter from 'gray-matter'

// Blog posts live as markdown in docs/content/blog (mirrored from the
// newsletter/ folder at the repo root, minus the internal thumbnail prompts).
// Read + parsed at build time; the /blog routes are fully static.

const BLOG_DIR = path.join(process.cwd(), 'content', 'blog')

export interface BlogPost {
  slug: string
  title: string
  subtitle: string
  date: string // ISO YYYY-MM-DD
  readingTime: string
  author: string
  tags: string[]
  category: string
  accent: string // tailwind gradient classes for the placeholder thumbnail
  content: string
}

const ACCENTS = [
  'from-sky-500 to-indigo-700',
  'from-cyan-500 to-blue-700',
  'from-indigo-500 to-violet-700',
  'from-slate-600 to-slate-900',
  'from-amber-500 to-orange-700',
  'from-emerald-500 to-teal-700',
  'from-sky-500 to-cyan-600',
]

const MONTHS = [
  'January', 'February', 'March', 'April', 'May', 'June',
  'July', 'August', 'September', 'October', 'November', 'December',
]

export function formatDate(iso: string): string {
  const dt = new Date(iso)
  if (Number.isNaN(dt.getTime())) return iso
  return `${MONTHS[dt.getUTCMonth()]} ${dt.getUTCDate()}, ${dt.getUTCFullYear()}`
}

function categoryFor(tags: string[]): string {
  const t = tags.map((x) => x.toLowerCase())
  if (t.includes('comparison')) return 'Comparison'
  if (t.includes('getting-started') || t.includes('tutorial')) return 'Tutorial'
  if (t.includes('auth') || t.includes('security')) return 'Security'
  if (t.includes('deploy') || t.includes('devops')) return 'Deploy'
  if (t.includes('crud') || t.includes('code-generation')) return 'Guide'
  return 'Announcement'
}

function toISO(d: unknown): string {
  if (d instanceof Date) return d.toISOString().slice(0, 10)
  return String(d ?? '')
}

export function getAllPosts(): BlogPost[] {
  if (!fs.existsSync(BLOG_DIR)) return []
  const files = fs.readdirSync(BLOG_DIR).filter((f) => f.endsWith('.md')).sort()

  const posts = files.map((file, i): BlogPost => {
    const raw = fs.readFileSync(path.join(BLOG_DIR, file), 'utf8')
    const { data, content } = matter(raw)
    const tags = (data.tags as string[]) || []
    return {
      slug: file.replace(/^\d{4}-\d{2}-\d{2}-/, '').replace(/\.md$/, ''),
      title: (data.title as string) || file,
      subtitle: (data.subtitle as string) || '',
      date: toISO(data.date),
      readingTime: (data.readingTime as string) || '5 min',
      author: (data.author as string) || 'Muke JohnBaptist',
      tags,
      category: categoryFor(tags),
      accent: ACCENTS[i % ACCENTS.length],
      content,
    }
  })

  // newest first
  return posts.sort((a, b) => (a.date < b.date ? 1 : -1))
}

export function getPost(slug: string): BlogPost | undefined {
  return getAllPosts().find((p) => p.slug === slug)
}

export function getCategories(posts: BlogPost[]): string[] {
  return Array.from(new Set(posts.map((p) => p.category))).sort()
}
