import Link from 'next/link'
import { Code2, Eye } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface PlaygroundChallengeProps {
  title: string
  description: string
  challenge: string
  solution: string
}

function encode(code: string) {
  return typeof window !== 'undefined'
    ? btoa(code)
    : Buffer.from(code).toString('base64')
}

export function PlaygroundChallenge({ title, description, challenge, solution }: PlaygroundChallengeProps) {
  const challengeEncoded = encode(challenge)
  const solutionEncoded = encode(solution)

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
          <div className="flex items-center gap-2">
            <Button size="sm" variant="outline" className="h-7 text-xs border-emerald-500/30 text-emerald-400 hover:bg-emerald-500/10 hover:text-emerald-300" asChild>
              <Link href={`/playground?code=${encodeURIComponent(challengeEncoded)}&title=${encodeURIComponent(title)}`}>
                <Code2 className="h-3 w-3 mr-1" />
                Try It
              </Link>
            </Button>
            <Button size="sm" variant="ghost" className="h-7 text-xs text-muted-foreground/60 hover:text-muted-foreground" asChild>
              <Link href={`/playground?code=${encodeURIComponent(solutionEncoded)}&title=${encodeURIComponent(title + ' (Solution)')}`}>
                <Eye className="h-3 w-3 mr-1" />
                See Solution
              </Link>
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
