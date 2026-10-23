'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'
import {
  CommandDialog,
  CommandInput,
  CommandList,
  CommandGroup,
  CommandItem,
} from '@/components/ui/command'
import {
  Rocket,
  Box,
  Server,
  Shield,
  Layers,
  Database,
  Settings,
  Palette,
  BookOpen,
  Lightbulb,
  Wand2,
  FileText,
  Monitor,
  Plug,
  GraduationCap,
} from 'lucide-react'
import { SEARCH_INDEX, type SearchItem as GeneratedSearchItem } from '@/config/search-index.generated'
import { COURSES } from '@/config/courses'
import { DEPLOYMENT_PROVIDERS } from '@/config/deployment-providers'

// The shape is declared with the generated index, which is the only thing
// that fills it for doc pages; courses and providers build items too.
type SearchItem = GeneratedSearchItem

// Auto-generated course entries: every course, every chapter, every lesson
// becomes searchable. Pulls from courses.ts so a new lesson just needs to
// land in the data model — no manual search-index updates.
function buildCourseSearchItems(): SearchItem[] {
  const out: SearchItem[] = []
  for (const course of COURSES) {
    const section = `Courses · ${course.shortName}`
    out.push({
      title: `Course: ${course.name}`,
      href: `/courses/${course.slug}`,
      section,
      keywords: `${course.name} ${course.shortName} ${course.tagline} ${course.whatYoullLearn.join(' ')} ${course.emoji} learning path course curriculum`,
    })
    for (const chapter of course.chapters) {
      out.push({
        title: `${course.shortName} · ch.${chapter.number}: ${chapter.title}`,
        href: `/courses/${course.slug}/${chapter.slug}`,
        section,
        keywords: `${chapter.title} ${chapter.tagline} ${chapter.learningGoals.join(' ')} chapter ${chapter.number}`,
      })
      for (const module of chapter.modules) {
        for (const lesson of module.lessons) {
          out.push({
            title: `${course.shortName} · ${lesson.title}`,
            href: `/courses/${course.slug}/${chapter.slug}/${lesson.slug}`,
            section,
            keywords: `${lesson.title} ${lesson.tagline} ${module.title} ${chapter.title} lesson ${lesson.minutes}m ${lesson.difficulty}`,
          })
        }
      }
      if (chapter.assignment) {
        out.push({
          title: `${course.shortName} · ch.${chapter.number} Assignment: ${chapter.assignment.title}`,
          href: `/courses/${course.slug}/${chapter.slug}/assignment`,
          section,
          keywords: `${chapter.assignment.title} ${chapter.assignment.brief} assignment exercise practice ${chapter.title}`,
        })
      }
    }
  }
  return out
}

/*
 * Every provider guide, generated from the same list the pages are generated
 * from. Hand-maintaining these was why none of the eight existed in search:
 * the pages shipped and the index was never touched, so "railway" and "fly"
 * found nothing.
 */
function buildDeploymentSearchItems(): SearchItem[] {
  return DEPLOYMENT_PROVIDERS.map((provider) => ({
    title: provider.name,
    href: `/docs/deployment/${provider.slug}`,
    section: 'Deployment',
    keywords: `deploy deployment host hosting ${provider.name} ${provider.slug} ${provider.tagline} ${provider.kind} ${provider.bestFor} ${provider.costFrom} production ship live`,
  }))
}

/**
 * Everything searchable: the generated page index, plus the two sets built
 * live from the configs they belong to.
 *
 * The deployment providers are in both. The generated entry is there so the
 * sitemap can enumerate pages that have no file of their own, and the one
 * built here also carries what a provider costs, what it runs and who it
 * suits, so that is the one worth keeping. Later wins, hence the map.
 */
const searchIndex: SearchItem[] = Array.from(
  new Map(
    [...SEARCH_INDEX, ...buildDeploymentSearchItems(), ...buildCourseSearchItems()].map(
      (item) => [item.href, item] as const
    )
  ).values()
)

const sectionIcons: Record<string, React.ReactNode> = {
  'Getting Started': <Rocket className="h-3.5 w-3.5" />,
  'Core Concepts': <Box className="h-3.5 w-3.5" />,
  'Backend (Go API)': <Server className="h-3.5 w-3.5" />,
  'Admin Panel': <Shield className="h-3.5 w-3.5" />,
  'Frontend (Next.js)': <Layers className="h-3.5 w-3.5" />,
  'Batteries': <Database className="h-3.5 w-3.5" />,
  'Infrastructure': <Settings className="h-3.5 w-3.5" />,
  'Deployment': <Rocket className="h-3.5 w-3.5" />,
  'Design System': <Palette className="h-3.5 w-3.5" />,
  'Tutorials': <BookOpen className="h-3.5 w-3.5" />,
  'For AI Assistants': <Lightbulb className="h-3.5 w-3.5" />,
  'AI Workflows': <Wand2 className="h-3.5 w-3.5" />,
  'Desktop (Wails)': <Monitor className="h-3.5 w-3.5" />,
  'Plugins': <Plug className="h-3.5 w-3.5" />,
  // Every course's section icon — auto-applied via the COURSES list above.
  ...Object.fromEntries(
    COURSES.map((c) => [`Courses · ${c.shortName}`, <GraduationCap key={c.slug} className="h-3.5 w-3.5" />]),
  ),
}

/*
 * Search scoring.
 *
 * cmdk's default filter scores by subsequence: every character of the query
 * has to appear in the value in order, with nothing in between mattering. Each
 * item's value here was the title plus its whole keyword blob, so "deplo"
 * matched Troubleshooting through "errors fix DEbug issues Problems heLp
 * brOken", and the results looked random.
 *
 * This is a substring match instead, ranked. A term either appears in the
 * title or the keywords or the item is gone, which makes results predictable:
 * what you typed is in what you get back.
 *
 * Multiple words are AND, not OR. "deploy railway" has to hit both, which is
 * what people expect from two words and the opposite of what a subsequence
 * match does with them.
 */
function escapeRegExp(value: string) {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

/** True when the term starts a word, so "form" hits "Form Builder" but not "platform". */
function startsWord(haystack: string, term: string) {
  return new RegExp(`\\b${escapeRegExp(term)}`).test(haystack)
}

/**
 * The singular and plural of a term, so one spelling finds the other.
 *
 * "fields" scored zero against "Field Types", because a substring match is
 * exact and the page is singular. The best the scorer could then do was course
 * lessons with the plural in their titles, which is what people were getting
 * back instead of the reference page about the thing.
 *
 * Deliberately crude: an -s/-es/-ies pair, not a stemmer. A stemmer would also
 * fold "testing" into "test" and "generated" into "generate", which changes
 * what a search means. This only says that one of a thing and several of it
 * are the same word.
 */
function variants(term: string): string[] {
  const forms = [term]
  const add = (value: string) => {
    if (value !== term && value.length > 2 && !forms.includes(value)) forms.push(value)
  }
  if (term.endsWith('ies')) add(term.slice(0, -3) + 'y')
  if (term.endsWith('es')) add(term.slice(0, -2))
  if (term.endsWith('s')) add(term.slice(0, -1))
  if (term.endsWith('y')) add(term.slice(0, -1) + 'ies')
  if (!term.endsWith('s')) {
    add(term + 's')
    add(term + 'es')
  }
  return forms
}

/**
 * How early in the keywords the term appears, as a number from 1 down to 0.9.
 *
 * Keywords are written most important first: the curated words for a page come
 * before the ones derived from its description. So two pages that both mention
 * a word are not equally about it, and without this the tie broke on the
 * alphabet. A search for "upload" put "Deploy to Railway" above "File Storage"
 * because D comes before F.
 *
 * Kept inside a tenth, so it only ever separates pages that already scored the
 * same way. It is a tie-break, not a rank.
 */
function keywordPosition(keywords: string, term: string) {
  const at = keywords.indexOf(term)
  if (at < 0 || keywords.length === 0) return 1
  return 1 - 0.1 * (at / keywords.length)
}

function scoreOneForm(title: string, keywords: string, term: string) {
  if (title === term) return 1
  if (title.startsWith(term)) return 0.9
  if (startsWord(title, term)) return 0.8
  if (title.includes(term)) return 0.6
  const position = keywordPosition(keywords, term)
  if (startsWord(keywords, term)) return 0.4 * position
  if (keywords.includes(term)) return 0.2 * position
  return 0
}

function scoreTerm(title: string, keywords: string, term: string) {
  let best = 0
  for (const [i, form] of variants(term).entries()) {
    // The typed spelling at full value; a singular or plural of it slightly
    // behind, so an exact match still sorts first.
    const score = scoreOneForm(title, keywords, form) * (i === 0 ? 1 : 0.95)
    if (score > best) best = score
  }
  return best
}

/**
 * How much a section's pages are worth, as a tie-breaker.
 *
 * A course lesson and a reference page can both match a word honestly, and
 * when they do the reference page is the one somebody searching a docs site
 * wants. This does not let a weak match beat a strong one: it is a nudge, far
 * smaller than the gap between a title hit and a keyword hit.
 */
function sectionWeight(section: string) {
  if (section === 'Courses') return 0.9
  if (section === 'Learnings' || section === 'Demo') return 0.95
  // Systems (Design) is a design-interview essay per subsystem, and for most
  // words there is a reference page about the same thing. Its titles are the
  // bare product name ("Pulse", "Sentinel"), so an exact-title match beat the
  // reference page that actually tells you how to use the thing: searching
  // "pulse" put the essay above Pulse (Observability). Far enough back that a
  // reference page wins, near enough that the essays are still the first hit
  // for anything only they cover.
  if (section === 'Systems (Design)') return 0.85
  return 1
}

export function scoreItem(title: string, keywords: string, search: string, section = '') {
  const query = search.trim().toLowerCase()
  if (!query) return 1

  const haystackTitle = title.toLowerCase()
  const haystackKeywords = keywords.toLowerCase()
  const terms = query.split(/\s+/).filter(Boolean)

  let total = 0
  for (const term of terms) {
    const score = scoreTerm(haystackTitle, haystackKeywords, term)
    // One unmatched word drops the item. Partial matches are how you end up
    // with a list where only the first word did any work.
    if (score === 0) return 0
    total += score
  }
  return (total / terms.length) * sectionWeight(section)
}

/** Cap on a result list. Past this nobody is reading, they are retyping. */
const MAX_RESULTS = 30

export function SearchDialog() {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const router = useRouter()

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === 'k' && (e.metaKey || e.ctrlKey)) {
        e.preventDefault()
        setOpen((open) => !open)
      }
    }

    document.addEventListener('keydown', down)
    return () => document.removeEventListener('keydown', down)
  }, [])

  const onSelect = useCallback((href: string) => {
    setOpen(false)
    router.push(href)
  }, [router])

  const sections = Array.from(new Set(searchIndex.map((item) => item.section)))

  /*
   * Ranking and filtering happen here rather than in cmdk, because cmdk sorts
   * items inside a group but renders the groups themselves in source order.
   * With the sections kept, a title match on "deploy" sat underneath two
   * keyword-only matches purely because Tech Kits is declared above
   * Infrastructure in the index. Searching wants one ranked list; the sections
   * are for browsing, so they come back when the box is empty.
   */
  const results = useMemo(() => {
    if (!query.trim()) return null
    return searchIndex
      .map((item) => ({ item, score: scoreItem(item.title, item.keywords, query, item.section) }))
      .filter((entry) => entry.score > 0)
      .sort((a, b) => b.score - a.score || a.item.title.localeCompare(b.item.title))
      .slice(0, MAX_RESULTS)
  }, [query])

  /* Clearing on close so reopening starts from browse rather than from
     whatever was typed last time. */
  function onOpenChange(next: boolean) {
    setOpen(next)
    if (!next) setQuery('')
  }

  return (
    <>
      <button
        onClick={() => setOpen(true)}
        className="inline-flex items-center gap-2 rounded-md border border-border/50 bg-accent/30 px-3 py-1.5 text-sm text-muted-foreground hover:bg-accent/50 hover:text-foreground transition-colors"
      >
        <FileText className="h-3.5 w-3.5" />
        <span className="hidden sm:inline">Search docs...</span>
        <kbd className="pointer-events-none hidden h-5 select-none items-center gap-1 rounded border border-border/50 bg-background/80 px-1.5 font-mono text-[10px] font-medium text-muted-foreground/70 sm:inline-flex">
          <span className="text-xs">⌘</span>K
        </kbd>
      </button>
      {/* shouldFilter={false}: the ranking above is the filter. Leaving cmdk's
          own filter on as well would sort the already-sorted list again. */}
      <CommandDialog open={open} onOpenChange={onOpenChange} shouldFilter={false}>
        <CommandInput
          placeholder="Search documentation..."
          value={query}
          onValueChange={setQuery}
        />
        <CommandList>
          {results === null ? (
            sections.map((section) => (
              <CommandGroup key={section} heading={section}>
                {searchIndex
                  .filter((item) => item.section === section)
                  .map((item) => (
                    <CommandItem
                      key={item.href}
                      value={item.href}
                      onSelect={() => onSelect(item.href)}
                      className="cursor-pointer"
                    >
                      <span className="text-muted-foreground/60">
                        {sectionIcons[item.section] || <FileText className="h-3.5 w-3.5" />}
                      </span>
                      <span>{item.title}</span>
                    </CommandItem>
                  ))}
              </CommandGroup>
            ))
          ) : results.length === 0 ? (
            /* Written out rather than using CommandEmpty, which keys off
               cmdk's own filtering and never fires with it turned off. */
            <div className="py-6 text-center text-sm text-muted-foreground">
              No results for &ldquo;{query}&rdquo;.
            </div>
          ) : (
            <CommandGroup heading={`${results.length} result${results.length === 1 ? '' : 's'}`}>
              {results.map(({ item }) => (
                <CommandItem
                  key={item.href}
                  value={item.href}
                  onSelect={() => onSelect(item.href)}
                  className="cursor-pointer"
                >
                  <span className="text-muted-foreground/60">
                    {sectionIcons[item.section] || <FileText className="h-3.5 w-3.5" />}
                  </span>
                  <span>{item.title}</span>
                  {/* The section is the only thing telling two similarly named
                      pages apart once the grouping is gone. */}
                  <span className="ml-auto shrink-0 pl-3 text-xs text-muted-foreground/70">
                    {item.section}
                  </span>
                </CommandItem>
              ))}
            </CommandGroup>
          )}
        </CommandList>
      </CommandDialog>
    </>
  )
}
