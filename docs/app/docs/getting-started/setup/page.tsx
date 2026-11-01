import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { PageHelp } from '@/components/page-help'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/setup')

function FlagTable({ rows }: { rows: [string, string][] }) {
  return (
    <div className="my-6 overflow-hidden rounded-xl border border-border/40">
      <table className="w-full text-sm">
        <thead className="bg-card/60">
          <tr>
            <th className="px-4 py-2.5 text-left font-medium">Flag</th>
            <th className="px-4 py-2.5 text-left font-medium">Meaning</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border/20">
          {rows.map(([flag, meaning]) => (
            <tr key={flag}>
              <td className="px-4 py-2 align-top">
                <code className="whitespace-nowrap text-xs font-mono text-accent/80">{flag}</code>
              </td>
              <td className="px-4 py-2 text-[13px] text-muted-foreground">{meaning}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function SetupPage() {
  return (
    <DocsShell
      title="grit setup"
      description="The guided installer: a starter, a shape, the services the project talks to, and the steps that make the result run."
    >
      <div className="prose-grit max-w-none">
        <p className="mb-10 text-lg leading-relaxed text-muted-foreground">
          <code>grit setup</code> creates a new project the way an installer does. It asks for a
          starter and for the services the project will use, then runs the steps that make the
          result runnable: dependencies, local services, migrations, sample data, the first commit.
        </p>

        <Callout type="note" title="It is not grit new with questions">
          <code>grit new</code> is the fast, scriptable scaffold and its behaviour is frozen, so a
          build script that runs it keeps working. <code>grit setup</code> is the guided path, and
          it can ask questions and change what the project starts with. Keeping them apart is why
          a change to one cannot surprise somebody relying on the other. Setup only creates new
          projects; use{' '}
          <Link href="/docs/upgrades" className="text-accent hover:underline">
            grit upgrade
          </Link>{' '}
          on one that exists.
        </Callout>

        <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Quick start</h2>
        <CodeBlock terminal code={`grit setup myapp`} />
        <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">The wizard asks for:</p>
        <ul className="mb-6 space-y-1.5 text-[14px] text-muted-foreground">
          <li>the architecture: single, double, triple, desktop, api, mobile or full</li>
          <li>the frontend: Next.js or TanStack Router, when the shape has one</li>
          <li>the starter: blank or blog</li>
          <li>the database: Postgres, MySQL, SQLite or in-memory</li>
          <li>Redis, storage and mail: local, external or off</li>
          <li>the built-in plugins that fit the chosen shape</li>
          <li>whether to run <code>grit start</code> at the end</li>
        </ul>

        <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Starters</h2>
        <p className="mb-4 leading-relaxed text-muted-foreground">
          A starter is the example content the project begins with. Both are the same scaffold:
          <code>blank</code> is that scaffold with the built-in Blog example removed, and{' '}
          <code>blog</code> keeps it exactly as <code>grit new</code> writes it. There is no
          separate template repository, so a starter cannot fall behind the scaffold.
        </p>
        <p className="mb-4 leading-relaxed text-muted-foreground">
          An <code>api</code> or <code>mobile</code> project always starts blank: it has no web or
          admin pages for the Blog&apos;s screens to live on.
        </p>

        <Callout type="warning" title="A blank project's edited files">
          Removing the Blog means editing a handful of files the scaffold owns: the nav menu, the
          seeder, the landing page. Those files then read as edited, so{' '}
          <code>grit upgrade</code> leaves them alone and reports them. That is deliberate: the
          templates still contain the Blog, so an upgrade that rewrote them would put the Blog
          links back, pointing at pages the starter removed. The trade is a handful of files that
          stop receiving framework fixes, and{' '}
          <code>grit upgrade --diff</code> shows what you are missing in them.
        </Callout>

        <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Non-interactive use</h2>
        <p className="mb-4 leading-relaxed text-muted-foreground">
          Pass <code>--yes</code> with the answers as flags. This is the shape for CI and scripts,
          and a project name is required. Setup refuses to prompt when it has no terminal, rather
          than hanging.
        </p>
        <CodeBlock
          terminal
          code={`grit setup myapp --yes \\
  --arch triple --frontend next --starter blank \\
  --db sqlite --redis off --storage local --mail log \\
  --no-docker --no-git --no-start`}
        />

        <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Services</h2>
        <p className="mb-4 leading-relaxed text-muted-foreground">
          Redis, storage and mail are written into <code>.env</code>. Each one can point at a
          local container, at an external provider, or be off.
        </p>
        <CodeBlock
          terminal
          code={`# Local, started with docker compose
grit setup myapp --redis local --storage minio --mail mailhog

# External
grit setup myapp --redis external --redis-url redis://...
grit setup myapp --storage s3 --storage-endpoint ... --storage-access-key ... --storage-secret-key ...
grit setup myapp --mail resend --mail-api-key ...

# Off
grit setup myapp --redis off --storage local --mail log`}
        />
        <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
          Credentials go into <code>.env</code>, which setup writes with owner-only permissions
          because it now holds real values rather than placeholders. No external secret is printed
          back to the terminal.
        </p>

        <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">What it runs</h2>
        <ol className="mb-6 space-y-1.5 text-[14px] text-muted-foreground">
          <li>Scaffolds the project with the chosen architecture, frontend and theme.</li>
          <li>Applies the starter.</li>
          <li>Configures <code>.env</code> for Redis, storage and mail.</li>
          <li>Installs the built-in plugins that fit the shape.</li>
          <li>Installs dependencies and tidies the Go module.</li>
          <li>Starts local services with Docker, unless <code>--no-docker</code>.</li>
          <li>Migrates the database and seeds sample data.</li>
          <li>Initialises git and makes the first commit.</li>
          <li>Offers to run <code>grit start</code>.</li>
        </ol>
        <p className="mb-4 leading-relaxed text-muted-foreground">
          Every step after the scaffold is a command you could run by hand, and when one fails
          setup prints that command rather than stopping with nothing to do next. A failed{' '}
          <code>pnpm install</code> does not cost you the migration or the first commit: neither
          needs node.
        </p>

        <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Flags</h2>
        <FlagTable
          rows={[
            ['--arch', 'single, double, triple, desktop, api, mobile or full'],
            ['--frontend', 'next or tanstack'],
            ['--starter', 'blank or blog'],
            ['--db', 'postgres, mysql, sqlite or memory'],
            ['--theme', 'atlas, aurora, pulse, coral, amber, sky, mono or emerald'],
            ['--redis', 'local, external or off'],
            ['--redis-url', 'External Redis URL'],
            ['--storage', 'minio, s3, r2, local or off'],
            ['--storage-*', 'S3 or R2 endpoint, keys, bucket, region and public URL'],
            ['--mail', 'mailhog, log, resend, mailgun, postmark, sendgrid, ses or smtp'],
            ['--mail-*', 'From address, API key, domain, endpoint and stream'],
            ['--smtp-*', 'SMTP host, port, username, password and encryption'],
            ['--ses-region, --aws-*', 'Amazon SES region and credentials'],
            ['--plugin', 'A built-in plugin to install; repeat or comma-separate'],
            ['--i18n', 'Add internationalisation'],
            ['--yes', 'Use the flags and defaults, with no prompts'],
            ['--here', 'Scaffold into the current directory'],
            ['--force', 'Allow scaffolding into a directory that is not empty'],
            ['--no-install', 'Skip pnpm install'],
            ['--no-docker', 'Do not start local services'],
            ['--no-migrate', 'Skip migrations'],
            ['--no-seed', 'Skip sample data'],
            ['--no-git', 'Skip git init and the first commit'],
            ['--star', 'Star MUKE-coder/grit without asking'],
            ['--start, --no-start', 'Run, or do not run, grit start at the end'],
          ]}
        />

        <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The logins it prints</h2>
        <p className="mb-4 leading-relaxed text-muted-foreground">
          The closing summary lists what the project exposes, with the login for each:
        </p>
        <CodeBlock
          language="text"
          code={`GORM Studio   http://localhost:8080/studio       | admin / <generated>
Sentinel      http://localhost:8080/sentinel/ui  | admin / <generated>
Admin         http://localhost:3001              | admin@example.com / admin123`}
        />
        <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
          GORM Studio and Sentinel get a generated local password each. The admin password is{' '}
          <code>admin123</code> in development unless <code>SEED_ADMIN_PASSWORD</code> is set, and
          when it is set setup says so rather than echoing it. All of these are development
          credentials. Change them before production.
        </p>
        <p className="mb-4 leading-relaxed text-muted-foreground">
          If the seed step was skipped, no admin login is shown: run <code>grit seed</code> to
          create it.
        </p>

        <Callout type="note" title="Contributed">
          <code>grit setup</code> was contributed by{' '}
          <a
            href="https://github.com/daljo25"
            className="text-accent hover:underline"
            rel="noopener noreferrer"
            target="_blank"
          >
            daljo25
          </a>{' '}
          in{' '}
          <a
            href="https://github.com/MUKE-coder/grit/pull/170"
            className="text-accent hover:underline"
            rel="noopener noreferrer"
            target="_blank"
          >
            #170
          </a>
          .
        </Callout>

        <div className="mt-12 flex flex-wrap gap-3">
          <Button asChild variant="outline">
            <Link href="/docs/getting-started/create-a-project">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Create a project
            </Link>
          </Button>
          <Button asChild>
            <Link href="/docs/getting-started/cli-cheatsheet">
              CLI cheatsheet
              <ArrowRight className="ml-2 h-4 w-4" />
            </Link>
          </Button>
        </div>

        <PageHelp />
      </div>
    </DocsShell>
  )
}
