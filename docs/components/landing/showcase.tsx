import Link from 'next/link'
import { ArrowRight, Heart } from 'lucide-react'
import { Button } from '@/components/grit/button'
import { HomeSponsors } from '@/components/sponsors'

/*
 * Built with Grit, with the sponsors folded in.
 *
 * Who backs it and who uses it are the same question, and they were two
 * sections asking it twice. The projects come first because a reader deciding
 * whether to try the framework wants evidence it has been used, not an
 * invitation to fund it.
 */

interface Project {
  name: string
  desc: string
  url: string
  tag: string
}

const PROJECTS: Project[] = [
  {
    name: 'GritCMS',
    desc: 'Self-hostable creator platform. Website builder, email marketing, courses, community.',
    url: 'https://gritcms.com',
    tag: 'SaaS',
  },
  {
    name: 'Golang Battles',
    desc: 'Competitive Go coding platform with real-time WebSocket battles, ELO ranking, and sandbox execution.',
    url: '#',
    tag: 'Platform',
  },
  {
    name: 'GORM Studio',
    desc: 'Visual database browser for GORM. View tables, run queries, export data. Embedded in every Grit project.',
    url: 'https://github.com/MUKE-coder/gorm-studio',
    tag: 'Tool',
  },
  {
    name: 'Pulse',
    desc: 'Self-hosted observability SDK. Request tracing, DB monitoring, runtime metrics, Prometheus export.',
    url: 'https://github.com/MUKE-coder/pulse',
    tag: 'Library',
  },
  {
    name: 'Sentinel',
    desc: 'WAF, rate limiting and brute-force protection, with an address that keeps attacking blocked automatically.',
    url: 'https://github.com/MUKE-coder/sentinel',
    tag: 'Security',
  },
  {
    name: 'Scout',
    desc: 'Security scanner you mount in your own API. Headers, TLS, DNS, exposed files, secrets, reachable CVEs, and your real route table.',
    url: 'https://github.com/MUKE-coder/scout',
    tag: 'Security',
  },
  {
    name: 'gin-docs',
    desc: 'Zero-annotation API documentation generator for Gin. Auto-generates OpenAPI spec with Scalar UI.',
    url: 'https://github.com/MUKE-coder/gin-docs',
    tag: 'Library',
  },
]

export function Showcase() {
  return (
    <>
      <ul className="mt-12 grid list-none grid-cols-1 gap-4 p-0 md:grid-cols-2 lg:grid-cols-3">
        {PROJECTS.map((project) => {
          const external = project.url.startsWith('http')
          return (
            <li key={project.name}>
              <Link
                href={project.url}
                target={external ? '_blank' : undefined}
                rel={external ? 'noreferrer' : undefined}
                className="g-card g-card--interactive h-full"
              >
                <div className="flex items-start justify-between gap-3">
                  <h3 className="g-card__title">{project.name}</h3>
                  <span className="g-badge font-mono text-[11px] text-fg-subtle">
                    {project.tag}
                  </span>
                </div>
                <p className="g-card__body">{project.desc}</p>
              </Link>
            </li>
          )
        })}

        <li>
          <Link
            href="/showcase"
            className="g-card g-card--interactive h-full justify-center text-center"
          >
            <span className="mx-auto grid h-9 w-9 place-items-center rounded-full border border-border-strong text-fg-subtle">
              +
            </span>
            <h3 className="g-card__title">Submit your project</h3>
            <p className="g-card__body">Built something with Grit? Get it listed.</p>
          </Link>
        </li>
      </ul>

      {/* Sponsors, folded in. */}
      <div className="mt-16 border-t border-border pt-12 text-center">
        <p className="g-eyebrow mb-3">Sponsors</p>
        <p className="mx-auto mb-8 max-w-xl text-[14px] leading-relaxed text-fg-muted">
          Grit is free and MIT licensed. Sponsors fund the features, docs and releases,
          and get their name in front of everyone who builds with it.
        </p>

        <HomeSponsors />

        <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
          <Button href="/sponsor" variant="secondary">
            <Heart className="h-4 w-4" aria-hidden="true" />
            Become a sponsor
          </Button>
          <Button href="/sponsors" variant="ghost">
            See all sponsors
            <ArrowRight className="h-4 w-4" aria-hidden="true" />
          </Button>
        </div>
      </div>
    </>
  )
}
