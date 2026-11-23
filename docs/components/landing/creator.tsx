import Link from 'next/link'
import { ArrowRight, Gauge, PenTool, ShieldCheck } from 'lucide-react'

/*
 * Word from the Creator.
 *
 * Somebody asked him why he built Grit and he answered at length. This is
 * that answer, condensed, in the first person and attributed, which is the
 * only form in which some of it can honestly be said: "I care a lot what
 * everything looks like" is a fact about a person, and a framework cannot
 * claim it in the third person without sounding like marketing.
 *
 * Placed after the showcase and before the FAQ on purpose. A reader who has
 * come that far has seen what the thing does and is deciding whether to trust
 * the judgement behind it, which is a question about a person.
 *
 * The three cards are the spine of his own answer rather than a list chosen
 * here: performance, interfaces and security are what he said he is obsessed
 * with, in that order, and every technical decision below follows from one of
 * them.
 */

const PORTFOLIO = 'https://jb.desishub.com/'
const AVATAR =
  'https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfAUUBxSZezIN7vwylkF1PXSCqAuseUG0gx8mh'

const OBSESSIONS = [
  {
    icon: Gauge,
    title: 'Performance',
    body: 'I want an app that is fast on my machine and still fast on the server. I tried Node and the MERN stack, Next.js on its own, Hono, Encore.ts, Laravel. None of them gave me what I wanted. Go and Rust are the two that do, and Go is the one a whole team can read.',
  },
  {
    icon: PenTool,
    title: 'Interfaces',
    body: 'I was a graphic designer before I was a programmer, so I care a lot what everything looks like, and clients come back because of it. An admin panel that looks generated is not finished.',
  },
  {
    icon: ShieldCheck,
    title: 'Security',
    body: 'I like working with large companies, and the thing they worry about most is security. They have millions of customers and cannot afford a mistake, so they want control of all of it: the server, the app, the firewall, the logs, the audit trail, roles and permissions.',
  },
]

export function Creator() {
  return (
    <div className="mx-auto max-w-4xl">
      <div className="flex flex-col items-start gap-5 sm:flex-row sm:items-center">
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={AVATAR}
          alt="Muke Johnbaptist, who created Grit"
          className="h-16 w-16 shrink-0 rounded-full object-cover ring-2 ring-accent/20"
          loading="lazy"
        />
        <div>
          <p className="text-[17px] font-semibold text-fg">Muke Johnbaptist</p>
          <p className="text-[14px] leading-relaxed text-fg-muted">
            Creator of Grit, Sentinel, Pulse and GORM Studio. Writes as{' '}
            <a href={PORTFOLIO} target="_blank" rel="noreferrer" className="g-link">
              JB Web Developer
            </a>
            .
          </p>
        </div>
      </div>

      <blockquote className="mt-9 border-l-2 border-accent/40 pl-6">
        <p className="text-[clamp(1.125rem,2vw,1.375rem)] font-medium leading-[1.55] text-fg [text-wrap:balance]">
          I built Grit to fix all of this without abstracting anything away from
          the developer, so you can still overwrite whatever you do not like.
        </p>
      </blockquote>

      <ul className="mt-10 grid list-none grid-cols-1 gap-4 p-0 md:grid-cols-3">
        {OBSESSIONS.map((item) => (
          <li key={item.title} className="g-card h-full">
            <item.icon className="h-5 w-5 text-accent" strokeWidth={1.75} aria-hidden="true" />
            <h3 className="g-card__title mt-3">{item.title}</h3>
            <p className="g-card__body">{item.body}</p>
          </li>
        ))}
      </ul>

      <div className="mt-10 space-y-5 text-[15px] leading-[1.75] text-fg-muted">
        <p>
          I have built more than a hundred full-stack applications, most of them
          Next.js with API routes. I moved to Go because an enterprise project
          needs a backend that carries real load. I kept Next.js because its
          routing, its SEO and the way it lets you write React are still the
          best I have used. I wanted the easiest possible way to build with
          both.
        </p>
        <p>
          Every full-stack app has the same two halves and the same problems: an
          API and a client, types that drift apart, and the question of where a
          mobile or desktop app goes when somebody asks for one. The answer is
          one repository. And for a serious product the web half is usually two
          apps rather than one: the site your customers use, and the admin your
          staff use. That is why{' '}
          <Link href="/docs/concepts/architecture-modes" className="g-link">
            triple is the default
          </Link>
          .
        </p>
        <p>
          Then the boring part. For every table you write a form, a table,
          search, export, pagination, edit, delete and a detail page on the
          client, and a schema, validation, a model, a service and a handler on
          the server. So you copy the last table and edit it. That is what{' '}
          <Link href="/docs/concepts/code-generation" className="g-link">
            grit generate resource
          </Link>{' '}
          is for, and it works because the job is predictable: you already know
          the table is called Product and that it holds a name, a price and an
          image. The generator gives you that starting point, and the business
          logic is yours.
        </p>
        <p>
          Everything else arrived because I needed it twice.{' '}
          <Link href="/docs/security" className="g-link">
            Sentinel
          </Link>{' '}
          because clients ask about security.{' '}
          <Link href="/docs/backend/pulse" className="g-link">
            Pulse
          </Link>{' '}
          because I wanted to see which line was burning the CPU.{' '}
          <Link href="/docs/infrastructure/studio" className="g-link">
            GORM Studio
          </Link>{' '}
          because I came from Prisma Studio and wanted the database in a
          browser. Wails because clients wanted a desktop app and Expo because
          they wanted a phone.{' '}
          <Link href="/plugins" className="g-link">
            Plugins
          </Link>{' '}
          because no framework can ship everything, and you should be able to
          publish your own.
        </p>
        <p>
          I also use AI constantly, so Grit is built to be worked on by it: an
          MCP server, a set of expert skills, and a CLI that does enough in one
          command that an agent spends a fraction of the tokens it would spend
          writing the same files by hand.
        </p>
      </div>

      <div className="mt-9 flex flex-wrap items-center gap-x-6 gap-y-3 text-[14px]">
        <Link href="/blog/why-i-built-grit-the-long-answer" className="g-link inline-flex items-center gap-1.5">
          The whole answer, in his words
          <ArrowRight className="h-3.5 w-3.5" />
        </Link>
        <Link href="/docs/getting-started/philosophy" className="g-link">
          The philosophy, and what each decision costs
        </Link>
        <Link href="/docs/governance" className="g-link">
          Who maintains it
        </Link>
      </div>
    </div>
  )
}
