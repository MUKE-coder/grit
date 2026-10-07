#!/usr/bin/env python3
"""Build one tutorial page per tier: the same contacts app, shipped six ways.

    python scripts/tier-tutorials.py            # write the pages
    python scripts/tier-tutorials.py --check    # fail if they are out of date

Every tier builds the same two resources with the same two commands, because
that is the point being made: what changes between a web app, a desktop app, a
phone app and an API is what you run and where you put it, not what you write.
The pages are generated from one table so that a correction to a shared step
cannot land on five pages and miss the sixth, which is what happens when six
near-identical pages are maintained by hand.

Each tier's commands were run end to end before this file was written: scaffold,
generate, migrate, seed, build every app it ships, start the API, call the
endpoints the resources were supposed to create, and dry-run the deploy. A
tutorial is a promise that these commands in this order produce the thing
described, and the only way to keep it is to run them.
"""
import io
import os
import sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUT_DIR = os.path.join(ROOT, 'docs', 'app', 'docs', 'tutorials', 'contacts')

INSTALL_NOTE = """
              <p>
                <code>pnpm install</code> is a one-time step and nothing runs it for you: the
                frontends are a pnpm workspace, and <code>grit start</code> has nothing to start
                without it.
              </p>"""

GROUP_CMD = '''grit generate resource Group \\
  --fields "name:string,description:text"'''

CONTACT_CMD = '''grit generate resource Contact \\
  --fields "name:string,email:email,phone:tel,photo:file:image,group:belongs_to:Group"'''


class Tier:
    def __init__(self, key, title, nav, flags, blurb, lead, prereq, structure,
                 run_cmd, run_note, ports, deploy, deploy_note, prev, nxt,
                 vite_flags=None):
        self.key = key
        self.title = title
        self.nav = nav
        self.flags = flags
        self.blurb = blurb
        self.lead = lead
        self.prereq = prereq
        self.structure = structure
        self.run_cmd = run_cmd
        self.run_note = run_note
        self.ports = ports
        self.deploy = deploy
        self.deploy_note = deploy_note
        self.prev = prev
        self.nxt = nxt
        # The same tier built with TanStack Router instead of Next.js. None for
        # the API tier, which has no frontend to choose and no theme to set.
        self.vite_flags = vite_flags


TIERS = [
    Tier(
        key='web',
        vite_flags='--triple --vite',
        title='Contacts: the web app',
        nav='Contacts: web app',
        flags='--triple --next',
        blurb='A Next.js site, an admin panel and a Go API.',
        lead=(
            'Three applications that share one API and one set of types: a public site people '
            'visit, an admin panel your team works in, and the Go service behind both. This is '
            'the shape most Grit projects take.'
        ),
        prereq=[
            ('Go 1.21 or newer', 'The API is a Go module. <code>go version</code> to check.'),
            ('Node 20 or newer, with pnpm', 'The two frontends are a pnpm workspace. <code>corepack enable pnpm</code> if you have none.'),
        ],
        structure=[
            ('apps/api', 'The Go API: models, services, handlers, routes'),
            ('apps/web', 'The public Next.js site'),
            ('apps/admin', 'The admin panel, where the generated screens land'),
            ('packages/shared', 'Zod schemas and TypeScript types, generated from the Go models'),
        ],
        run_cmd='grit start',
        run_note=(
            'One command brings up all three, each line of output prefixed with the app it came '
            'from. Ctrl+C stops them together.'
        ),
        ports=[
            ('http://localhost:8080/api/v1', 'The API'),
            ('http://localhost:8080/docs', 'Its reference, generated from the routes'),
            ('http://localhost:3000', 'The public site'),
            ('http://localhost:3001', 'The admin panel: sign in as admin@example.com / admin123'),
        ],
        deploy='''# A server you have SSH access to: systemd, Caddy, TLS
grit deploy --host user@server.com --domain contacts.example.com

# Or Railway, which takes the variables from .env and gives you a URL
grit deploy --railway''',
        deploy_note=(
            'The first form builds the API, uploads it over SSH, writes a systemd unit and puts '
            'Caddy in front of it with a certificate. The two frontends are Next.js apps and '
            'deploy like any other: <code>pnpm build</code> and your host of choice, with '
            '<code>NEXT_PUBLIC_API_URL</code> pointed at the API you just deployed.'
        ),
        prev=None,
        nxt=('desktop', 'The desktop app'),
    ),
    Tier(
        key='desktop',
        vite_flags='--triple --desktop --vite',
        title='Contacts: the desktop app',
        nav='Contacts: desktop app',
        flags='--triple --desktop',
        blurb='The same app, as a native window, with Wails.',
        lead=(
            'Everything the web tier gives you, and a desktop application beside it that shares '
            'the same API and the same types. One Go binary with a native webview, which '
            'installs the way desktop software installs.'
        ),
        prereq=[
            ('Go 1.21 or newer', 'Both the API and the desktop shell are Go.'),
            ('Node 20 or newer, with pnpm', 'The frontends, including the desktop app&apos;s own.'),
            ('The Wails CLI', 'Install from <a className="underline" href="https://wails.io">wails.io</a>. <code>grit start desktop</code> says so plainly if it is missing, rather than failing obscurely.'),
        ],
        structure=[
            ('apps/api', 'The Go API, shared by every client'),
            ('apps/desktop', 'The Wails application: a Go shell around a React frontend'),
            ('apps/web', 'The public Next.js site'),
            ('apps/admin', 'The admin panel'),
        ],
        run_cmd='grit start desktop',
        run_note=(
            'Wails opens a native window with hot reload. Run <code>grit start</code> in another '
            'terminal for the API it talks to, or <code>grit start</code> on its own brings up '
            'the desktop app too when the toolchain is installed.'
        ),
        ports=[
            ('A native window', 'The desktop app, reloading as you edit'),
            ('http://localhost:8080/api/v1', 'The API it calls'),
            ('http://localhost:3001', 'The admin panel, still a browser app'),
        ],
        deploy='''# The API, as in every other tier
grit deploy --host user@server.com --domain contacts.example.com

# The desktop app is a file people install
cd apps/desktop && wails build''',
        deploy_note=(
            'Deployment splits in two here. The API goes to a server like any other; the desktop '
            'application is an artefact you distribute. <code>wails build</code> produces a '
            'binary for the platform you build on, so a release for three platforms means '
            'building on three, or in CI. Point it at the deployed API by setting '
            '<code>VITE_API_URL</code> before you build, or it will look for localhost on a '
            'stranger&apos;s machine.'
        ),
        prev=('web', 'The web app'),
        nxt=('mobile', 'The mobile app'),
    ),
    Tier(
        key='mobile',
        vite_flags='--triple --expo --vite',
        title='Contacts: the mobile app',
        nav='Contacts: mobile app',
        flags='--triple --expo',
        blurb='The same app on a phone, with Expo.',
        lead=(
            'An Expo client beside the web apps, sharing the Zod schemas and the TypeScript '
            'types the Go models generate. The contact you create on the phone is the contact '
            'the admin panel shows, because there is one API and one definition of what a '
            'contact is.'
        ),
        prereq=[
            ('Go 1.21 or newer', 'For the API.'),
            ('Node 20 or newer, with pnpm', 'The workspace, including apps/expo.'),
            ('Expo Go, or a simulator', 'Expo Go on your phone is the quickest. Xcode or Android Studio if you want a simulator.'),
        ],
        structure=[
            ('apps/api', 'The Go API'),
            ('apps/expo', 'The Expo client: screens, navigation, and the API hooks'),
            ('apps/web', 'The public site'),
            ('apps/admin', 'The admin panel'),
            ('packages/shared', 'The types all three clients import'),
        ],
        run_cmd='grit start expo',
        run_note=(
            'Metro starts and prints a QR code. Scan it with Expo Go, or press <code>i</code> / '
            '<code>a</code> for a simulator. The API needs to be running too, from '
            '<code>grit start server</code> in another terminal.'
        ),
        ports=[
            ('Expo Go, or a simulator', 'The phone app'),
            ('http://localhost:8080/api/v1', 'The API'),
            ('http://localhost:3001', 'The admin panel'),
        ],
        deploy='''# The API
grit deploy --host user@server.com --domain contacts.example.com

# The app, built by Expo's service and submitted to the stores
cd apps/expo
npx eas build --platform all
npx eas submit''',
        deploy_note=(
            'A phone on your network cannot reach <code>localhost</code>, so set '
            '<code>EXPO_PUBLIC_API_URL</code> to the deployed API before building. During '
            'development, point it at your machine&apos;s LAN address rather than 127.0.0.1 for '
            'the same reason.'
        ),
        prev=('desktop', 'The desktop app'),
        nxt=('api', 'The API on its own'),
    ),
    Tier(
        key='api',
        title='Contacts: the API on its own',
        nav='Contacts: API only',
        flags='--api',
        blurb='The Go service, with no frontend at all.',
        lead=(
            'When the clients are somebody else&apos;s: a mobile team, a partner integration, or '
            'a frontend you have already built. You get the Go API with authentication, '
            'migrations, generated reference documentation and nothing else. No Node, no pnpm, '
            'no JavaScript anywhere in the project.'
        ),
        prereq=[
            ('Go 1.21 or newer', 'The only thing this tier needs.'),
        ],
        structure=[
            ('apps/api/internal/models', 'The Go structs the database is built from'),
            ('apps/api/internal/services', 'Where the queries live'),
            ('apps/api/internal/handlers', 'Thin HTTP handlers that call the services'),
            ('apps/api/internal/routes', 'The route table, and the access registry built from it'),
        ],
        run_cmd='grit start',
        run_note=(
            'There are no frontends to start, so this runs the API alone, with hot reload.'
        ),
        ports=[
            ('http://localhost:8080/api/v1', 'The API'),
            ('http://localhost:8080/docs', 'The reference, generated from the routes'),
            ('http://localhost:8080/studio', 'GORM Studio: browse and edit the tables'),
        ],
        deploy='''grit deploy --host user@server.com --domain api.example.com

# Or a container, if that is your shape
docker compose -f docker-compose.prod.yml up --build''',
        deploy_note=(
            'The simplest tier to deploy, because there is one thing to deploy. The Docker form '
            'builds the API image and brings up Postgres and Redis beside it.'
        ),
        prev=('mobile', 'The mobile app'),
        nxt=('single', 'One binary'),
    ),
    Tier(
        key='single',
        vite_flags='--single --next',
        title='Contacts: one binary',
        nav='Contacts: one binary',
        flags='--single',
        blurb='The API and the whole app in one executable.',
        lead=(
            'One folder that looks like the web project it is, with the Go API in '
            '<code>api/</code> beside it and the admin panel at <code>/admin</code>. The build '
            'puts the frontend inside the binary, so what you deploy is a single file with no '
            'runtime to install next to it. Both frontends work this way: TanStack by default, '
            'Next.js with <code>--next</code>, which builds a static export into the same place.'
        ),
        prereq=[
            ('Go 1.21 or newer', 'The API, and the binary that ends up carrying everything.'),
            ('Node 20 or newer, with pnpm', 'To build the frontend. Only at build time: the server needs nothing.'),
        ],
        structure=[
            ('src', 'The Vite app: routes, components, and the admin panel under admin-panel/'),
            ('api', 'The Go module: main.go, internal/, cmd/'),
            ('api/web', 'Where the frontend build lands, so the binary can embed it'),
            ('packages', 'The upload client the frontend shares with the API'),
        ],
        run_cmd='grit start',
        run_note=(
            'The Go API and the frontend dev server, in parallel. In development they are two '
            'processes, because a dev server is what gives you hot reload; in production they '
            'are one file.'
        ),
        ports=[
            ('http://localhost:5173', 'The app, with the admin panel at /admin'),
            ('http://localhost:8080/api/v1', 'The API'),
            ('http://localhost:8080/docs', 'The reference'),
        ],
        deploy='''# Build the frontend into api/web, then the binary around it
pnpm build
cd api && go build -o ../bin/contacts .

# That one file is the deployment
scp bin/contacts user@server.com:/srv/contacts

# Or let grit do the systemd and TLS part
grit deploy --host user@server.com --domain contacts.example.com''',
        deploy_note=(
            'The binary serves the API and the whole frontend, so there is nothing else to put '
            'on the server: no Node, no nginx in front of static files, no second process. Copy '
            'it beside a <code>.env</code> and run it. <code>grit deploy</code> does the same '
            'thing and adds a systemd unit and a certificate. '
            'One difference worth knowing with <code>--next</code>: a static export has to have '
            'a file for every URL it serves, and the id of a record is not known at build time, '
            'so '
            'a detail page is <code>/resources/users/view?id=...</code> rather than '
            '<code>/resources/users/123</code>.'
        ),
        prev=('api', 'The API on its own'),
        nxt=('full', 'Everything at once'),
    ),
    Tier(
        key='full',
        vite_flags='--full --vite',
        title='Contacts: everything at once',
        nav='Contacts: everything',
        flags='--full',
        blurb='Web, admin, desktop, mobile, docs and the API.',
        lead=(
            'Every client Grit can scaffold, against one API and one set of types. Useful when '
            'you know you want more than one client, and useful for looking at what the other '
            'tiers give you before choosing one.'
        ),
        prereq=[
            ('Go 1.21 or newer', 'The API and the desktop shell.'),
            ('Node 20 or newer, with pnpm', 'Four JavaScript applications.'),
            ('The Wails CLI', 'For the desktop app. The rest works without it.'),
            ('Expo Go, or a simulator', 'For the phone app.'),
        ],
        structure=[
            ('apps/api', 'The Go API'),
            ('apps/web', 'The public Next.js site'),
            ('apps/admin', 'The admin panel'),
            ('apps/desktop', 'The Wails desktop app'),
            ('apps/expo', 'The Expo mobile client'),
            ('apps/docs', 'A documentation site for the project'),
        ],
        run_cmd='grit start',
        run_note=(
            'Brings up the API, the web app, the admin panel and the desktop app when Wails is '
            'installed. The Expo client is its own command, <code>grit start expo</code>, '
            'because Metro wants a terminal of its own.'
        ),
        ports=[
            ('http://localhost:8080/api/v1', 'The API'),
            ('http://localhost:3000', 'The public site'),
            ('http://localhost:3001', 'The admin panel'),
            ('A native window', 'The desktop app, when Wails is installed'),
            ('Expo Go', 'The phone app, after grit start expo'),
        ],
        deploy='''# The API, and the web apps behind it
grit deploy --host user@server.com --domain contacts.example.com

# Then each client, on its own schedule
cd apps/desktop && wails build
cd apps/expo && npx eas build --platform all''',
        deploy_note=(
            'The server side deploys once. The desktop and mobile apps are artefacts with their '
            'own release cycles: people update them when they choose to, which is the reason to '
            'keep the API backward compatible rather than assuming every client is current.'
        ),
        prev=('single', 'One binary'),
        nxt=None,
    ),
]


THEMES = [
    ('atlas', 'Split-screen sign-in, Inter. The default: sharp and neutral, for a team tool.'),
    ('aurora', 'Centered sign-in, Geist. Pastel and friendly, for consumer software.'),
    ('pulse', 'Split-screen with a carousel, Onest and DM Serif. Warm and bold, for a brand.'),
    ('coral', 'A sign-in modal over the page. Soft, close up.'),
    ('amber', 'A boxed sign-in card. Warm neutrals.'),
    ('sky', 'A banner above the form. Open and light.'),
    ('mono', 'A showcase panel beside the form. Monochrome and typographic.'),
    ('emerald', 'A quote beside the form. Green and calm: the one the commands above use.'),
]

FRONTEND_TEMPLATE = """              <h3>Next.js or TanStack</h3>
              <p>
                The command above builds the frontend with Next.js. The same tier with TanStack
                Router and Vite instead is one flag:
              </p>
{command}
              <p>
                Everything after this point is identical. The resources, the API, the admin
                screens and the deployment are the same; what changes is the router and the build
                tool, and the admin panel is written in whichever dialect that app speaks.
              </p>

"""

THEME_TEMPLATE = """              <h3>Pick a theme</h3>
              <p>
                <code>--theme</code> sets the sign-in layout, the dashboard tokens, the fonts and
                the brand colours, and it applies to both frontends. Eight ship with Grit:
              </p>

              <table className="w-full text-sm my-6">
                <thead>
                  <tr className="border-b border-border/50 text-left">
                    <th className="px-4 py-2.5 font-medium">Theme</th>
                    <th className="px-4 py-2.5 font-medium">What it looks like</th>
                  </tr>
                </thead>
                <tbody>
{rows}
                </tbody>
              </table>

              <p>
                Nothing is baked in. <code>THEME</code> and <code>VITE_THEME</code> in{{' '}}
                <code>.env</code> both carry the name, so changing it there repaints the app
                without re-scaffolding.
              </p>

"""

THEME_ROW = (
    '                  <tr className="border-b border-border/30">\n'
    '                    <td className="px-4 py-2.5 font-mono text-xs text-primary">{name}</td>\n'
    '                    <td className="px-4 py-2.5 text-muted-foreground">{note}</td>\n'
    '                  </tr>'
)


def frontend_section(tier, code_block):
    """Next.js or TanStack, for a tier that offers both."""
    if not tier.vite_flags:
        return ''
    cmd = code_block('grit new contacts %s --theme emerald --db sqlite' % tier.vite_flags)
    return FRONTEND_TEMPLATE.format(command=cmd)


def theme_section(tier):
    """The themes, for a tier with a frontend to apply one to."""
    if not tier.vite_flags:
        return ''
    rows = '\n'.join(THEME_ROW.format(name=name, note=note) for name, note in THEMES)
    return THEME_TEMPLATE.format(rows=rows)


def esc(text):
    """JSX text: the two characters that end an expression or a tag."""
    return text.replace('{', '&#123;').replace('}', '&#125;')


def code_block(code, terminal=True, language='bash'):
    attrs = ' terminal' if terminal else ''
    return (
        '              <CodeBlock%s language="%s" code={`%s`} />'
        % (attrs, language, code.replace('\\', '\\\\').replace('`', '\\`').replace('${', '\\${'))
    )


def render(tier):
    prereq = '\n'.join(
        '                  <li>\n'
        '                    <strong>%s</strong> &mdash; %s\n'
        '                  </li>' % (name, note)
        for name, note in tier.prereq
    )

    structure = '\n'.join(
        '                  <tr className="border-b border-border/30">\n'
        '                    <td className="px-4 py-2.5 font-mono text-xs text-primary">%s</td>\n'
        '                    <td className="px-4 py-2.5 text-muted-foreground">%s</td>\n'
        '                  </tr>' % (path, note)
        for path, note in tier.structure
    )

    ports = '\n'.join(
        '                  <tr className="border-b border-border/30">\n'
        '                    <td className="px-4 py-2.5 font-mono text-xs text-primary">%s</td>\n'
        '                    <td className="px-4 py-2.5 text-muted-foreground">%s</td>\n'
        '                  </tr>' % (where, what)
        for where, what in tier.ports
    )

    nav_prev = (
        '              <Button variant="ghost" asChild>\n'
        '                <Link href="/docs/tutorials/contacts/%s">\n'
        '                  <ArrowLeft className="mr-2 h-4 w-4" />\n'
        '                  %s\n'
        '                </Link>\n'
        '              </Button>' % tier.prev
        if tier.prev else
        '              <Button variant="ghost" asChild>\n'
        '                <Link href="/docs/tutorials">\n'
        '                  <ArrowLeft className="mr-2 h-4 w-4" />\n'
        '                  All tutorials\n'
        '                </Link>\n'
        '              </Button>'
    )
    nav_next = (
        '              <Button variant="ghost" asChild>\n'
        '                <Link href="/docs/tutorials/contacts/%s">\n'
        '                  %s\n'
        '                  <ArrowRight className="ml-2 h-4 w-4" />\n'
        '                </Link>\n'
        '              </Button>' % tier.nxt
        if tier.nxt else
        '              <Button variant="ghost" asChild>\n'
        '                <Link href="/docs/deployment">\n'
        '                  Deployment in depth\n'
        '                  <ArrowRight className="ml-2 h-4 w-4" />\n'
        '                </Link>\n'
        '              </Button>'
    )

    component = 'Contacts' + tier.key.capitalize() + 'TutorialPage'

    return '''import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

// Generated by scripts/tier-tutorials.py. Edit the table there, not this file.

export const metadata = getDocMetadata('/docs/tutorials/contacts/%(key)s')

export default function %(component)s() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Tutorial</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">%(title)s</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">%(lead)s</p>
            </div>

            <div className="prose-grit">
              <p>
                The app is a contact book: contacts that belong to groups, each with a photo. Two
                resources, two commands, and everything that follows is the same in every tier.
                What differs is what you run and where it ends up, which is what this page is
                about.
              </p>

              <h2>Before you start</h2>
              <ul>
%(prereq)s
              </ul>
              <p>
                And Grit itself: <code>go install github.com/MUKE-coder/grit/v3/cmd/grit@latest</code>.
              </p>

              <h2>1. Create the project</h2>
%(new_cmd)s
              <p>
                Postgres is the default. <code>--db sqlite</code> above means the project runs
                with no database server at all, which is the shortest path to seeing it work; drop
                the flag when you want Postgres, and <code>docker compose up -d</code> brings one
                up along with Redis, MinIO and a mail catcher.
              </p>%(install_note)s

%(frontend_section)s%(theme_section)s
              <table className="w-full text-sm my-6">
                <thead>
                  <tr className="border-b border-border/50 text-left">
                    <th className="px-4 py-2.5 font-medium">Directory</th>
                    <th className="px-4 py-2.5 font-medium">What is in it</th>
                  </tr>
                </thead>
                <tbody>
%(structure)s
                </tbody>
              </table>

              <h2>2. Describe the data</h2>
              <p>
                A group first, because a contact points at one. Each command writes the Go model,
                the migration, the service, the handler, the routes, the Zod schema, the
                TypeScript types and an admin screen, and registers all of it.
              </p>
%(group_cmd)s
              <p>
                Then the contact, with a photo and a group to belong to. The{' '}
                <code>belongs_to:Group</code> field is what makes the admin screen render a picker
                rather than a text box asking for an id.
              </p>
%(contact_cmd)s

              <h2>3. Build the database</h2>
%(migrate_cmd)s
              <p>
                <code>grit migrate</code> creates the tables from the models, and{' '}
                <code>grit seed</code> fills them: an administrator you can sign in as, a few
                users, and an API key. The administrator is{' '}
                <code>admin@example.com</code> with the password <code>admin123</code>. Change it
                before anybody else can reach the machine.
              </p>

              <h2>4. Run it</h2>
%(run_cmd)s
              <p>%(run_note)s</p>

              <table className="w-full text-sm my-6">
                <thead>
                  <tr className="border-b border-border/50 text-left">
                    <th className="px-4 py-2.5 font-medium">Where</th>
                    <th className="px-4 py-2.5 font-medium">What</th>
                  </tr>
                </thead>
                <tbody>
%(ports)s
                </tbody>
              </table>

              <h2>5. What you have</h2>
              <p>
                Two resources, and for each of them: a table with sorting, filtering, search,
                pagination, bulk edit and CSV import and export; a form that validates on both
                sides from one schema; and REST endpoints under <code>/api/v1/groups</code> and{' '}
                <code>/api/v1/contacts</code> with the same rules applied.
              </p>
              <p>
                The photo field is a real upload: the browser asks the API for a presigned URL and
                sends the file straight to storage, so the file never passes through the API
                process. With <code>STORAGE_DRIVER=local</code> that storage is a directory on
                disk, which is the default when there is no MinIO to talk to.
              </p>

              <h2>6. Deploy it</h2>
%(deploy_cmd)s
              <p>%(deploy_note)s</p>
              <p>
                Before any of that, read the{' '}
                <Link href="/docs/deployment/checklist">go-live checklist</Link>: it is the list of
                things that are fine in development and not in production, starting with the
                seeded password above and the secrets in <code>.env</code>.
              </p>
            </div>

            <div className="mt-16 flex flex-col sm:flex-row gap-4 justify-between border-t border-border pt-8">
%(nav_prev)s
%(nav_next)s
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
''' % {
        'key': tier.key,
        'component': component,
        'title': esc(tier.title),
        'lead': tier.lead,
        'prereq': prereq,
        'structure': structure,
        'ports': ports,
        'new_cmd': code_block('grit new contacts %s --theme emerald --db sqlite\ncd contacts%s'
                              % (tier.flags, '' if not tier.vite_flags else '\npnpm install')),
        'install_note': '' if not tier.vite_flags else INSTALL_NOTE,
        'frontend_section': frontend_section(tier, code_block),
        'theme_section': theme_section(tier),
        'group_cmd': code_block(GROUP_CMD),
        'contact_cmd': code_block(CONTACT_CMD),
        'migrate_cmd': code_block('grit migrate\ngrit seed'),
        'run_cmd': code_block(tier.run_cmd),
        'run_note': tier.run_note,
        'deploy_cmd': code_block(tier.deploy),
        'deploy_note': tier.deploy_note,
        'nav_prev': nav_prev,
        'nav_next': nav_next,
    }


def main():
    check = '--check' in sys.argv
    stale = []
    for tier in TIERS:
        path = os.path.join(OUT_DIR, tier.key, 'page.tsx')
        body = render(tier)
        if check:
            try:
                current = io.open(path, encoding='utf-8', newline='').read().replace('\r\n', '\n')
            except IOError:
                stale.append(tier.key + ' (missing)')
                continue
            if current != body:
                stale.append(tier.key)
            continue
        directory = os.path.dirname(path)
        if not os.path.isdir(directory):
            os.makedirs(directory)
        io.open(path, 'w', encoding='utf-8', newline='').write(body)
        print('wrote docs/app/docs/tutorials/contacts/%s/page.tsx' % tier.key)

    if check:
        if stale:
            print('out of date: ' + ', '.join(stale))
            print('run: python scripts/tier-tutorials.py')
            return 1
        print('tier tutorials are up to date (%d pages)' % len(TIERS))
    return 0


if __name__ == '__main__':
    sys.exit(main())
