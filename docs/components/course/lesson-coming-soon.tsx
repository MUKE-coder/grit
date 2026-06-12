import Link from 'next/link'
import { Construction, ArrowRight, Bell } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface Props {
  /** Course slug — for the "Back to course" link. */
  courseSlug: string
  /** Course name — for context. */
  courseName: string
  /** Lesson title — what they wanted to read. */
  lessonTitle: string
}

/**
 * LessonComingSoon — placeholder rendered by every stub lesson page in
 * the 6 non-Grit-Concepts courses. Keeps the SEO + sitemap real (every
 * lesson is a route) without forcing the author to write all 80+
 * lessons in one sitting.
 */
export function LessonComingSoon({ courseSlug, courseName, lessonTitle }: Props) {
  return (
    <div className="not-prose my-12 rounded-2xl border border-border/50 bg-card/30 p-8 text-center">
      <div className="inline-flex items-center justify-center h-14 w-14 rounded-full bg-amber-500/10 text-amber-500 mb-5">
        <Construction className="h-7 w-7" />
      </div>
      <p className="text-xs uppercase tracking-wider text-amber-500 font-mono mb-2">
        Lesson in production
      </p>
      <h2 className="text-2xl font-semibold tracking-tight mb-3">
        &quot;{lessonTitle}&quot;
      </h2>
      <p className="text-sm text-muted-foreground max-w-xl mx-auto leading-relaxed mb-6">
        This lesson is part of the <strong>{courseName}</strong> course and
        is being written. The structure (chapters / modules / lesson
        order) is final — the prose is being polished. While you wait, the{' '}
        <Link href="/courses/concepts" className="text-primary hover:underline">
          Grit Concepts course
        </Link>{' '}
        is fully written and covers the foundation every other course
        builds on.
      </p>
      <div className="flex items-center justify-center gap-3 flex-wrap">
        <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground">
          <Link href="/courses/concepts">
            Start with Grit Concepts
            <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
          </Link>
        </Button>
        <Button asChild variant="outline" className="rounded-full border-border/60">
          <Link href={`/courses/${courseSlug}`}>Back to course outline</Link>
        </Button>
      </div>
      <div className="mt-6 inline-flex items-center gap-1.5 text-xs text-muted-foreground">
        <Bell className="h-3 w-3" />
        Bookmark — drop a star on the{' '}
        <Link
          href="https://github.com/MUKE-coder/grit"
          target="_blank"
          rel="noopener noreferrer"
          className="text-primary hover:underline"
        >
          repo
        </Link>{' '}
        to see when this ships.
      </div>
    </div>
  )
}
