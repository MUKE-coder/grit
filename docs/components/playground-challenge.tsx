import Link from 'next/link'
import { Code2 } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface PlaygroundChallengeProps {
  title: string
  description: string
  code: string
}

export function PlaygroundChallenge({ title, description, code }: PlaygroundChallengeProps) {
  const encoded = typeof window !== 'undefined'
    ? btoa(code)
    : Buffer.from(code).toString('base64')

  return (
    <div className="rounded-xl border border-emerald-500/20 bg-emerald-500/5 p-5 mb-8">
      <div className="flex items-start gap-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-emerald-500/10 border border-emerald-500/20">
          <Code2 className="h-4 w-4 text-emerald-400" />
        </div>
        <div className="flex-1 min-w-0">
          <h4 className="text-sm font-semibold text-emerald-400 uppercase tracking-wider mb-1">
            Try This
          </h4>
          <p className="text-[13px] text-muted-foreground/80 leading-relaxed mb-3">
            {description}
          </p>
          <Button size="sm" variant="outline" className="h-7 text-xs border-emerald-500/30 text-emerald-400 hover:bg-emerald-500/10 hover:text-emerald-300" asChild>
            <Link href={`/playground?code=${encodeURIComponent(encoded)}&title=${encodeURIComponent(title)}`}>
              Open in Playground
            </Link>
          </Button>
        </div>
      </div>
    </div>
  )
}
