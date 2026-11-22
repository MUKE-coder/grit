import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Note, Mistake, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'
import { REQUIRED_GO, REQUIRED_NODE, PINNED_PNPM } from '@/config/site'

export const metadata: Metadata = {
  title: 'Fix the errors you hit starting a Go and Docker project',
  description:
    'Every error message you meet in the first hour of a Go, Docker and Node project, with the cause and the fix for Windows, macOS and Linux.',
  alternates: { canonical: 'https://gritframework.dev/guides/common-go-docker-errors' },
  openGraph: {
    title: 'Fix the errors you hit starting a Go and Docker project',
    description: 'The exact messages, the cause of each, and the fix per operating system.',
    url: 'https://gritframework.dev/guides/common-go-docker-errors',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="common-go-docker-errors"
      lead={
        <>
          <p className="mb-4">
            Almost every error in the first hour of a new project is one of about twenty, and none
            of them means you have done something stupid. They mean a tool is not installed, is not
            on your PATH, is not running, or is already holding the port you asked for.
          </p>
          <p>
            Each section below is headed with the message as it actually appears, so you can match
            what is on your screen. Read the cause, not just the command: the same fix applied
            without understanding comes back a week later.
          </p>
        </>
      }
    >
      <Section id="before" title="The three questions that solve most of them">
        <p>
          Before searching for a specific message, answer these. They resolve a surprising
          proportion of what follows.
        </p>
        <H3>Did you open a new terminal?</H3>
        <p>
          Installing a tool changes your PATH. A shell that was already running kept the old one, so{' '}
          <code>go</code>, <code>docker</code> or <code>pnpm</code> will still look missing no
          matter how many times you reinstall. Close the terminal, open a new one, and check again.
          On Windows this includes the integrated terminal inside your editor: restart the editor.
        </p>
        <H3>Is the thing actually running?</H3>
        <p>
          Docker Desktop is an application. If it is not open, every <code>docker</code> command
          fails, and the message does not say &quot;Docker Desktop is closed&quot;, it says
          something about a pipe or a socket. Open it and wait for the whale icon to stop
          animating.
        </p>
        <H3>Are you in the right directory?</H3>
        <p>
          Go commands read the <code>go.mod</code> in the current directory. Node commands read{' '}
          <code>package.json</code>. Half of the &quot;cannot find module&quot; family is a command
          run one level too high or too low. <code>pwd</code>, or <code>cd</code> on Windows, before
          anything else.
        </p>
      </Section>

      <Section id="go" title="Go">
        <Mistake error="go: command not found">
          <p>
            Also <code>&apos;go&apos; is not recognized as an internal or external command</code> on
            Windows. Go is not installed, or it is installed and its <code>bin</code> directory is
            not on your PATH.
          </p>
          <p>
            Install Go {REQUIRED_GO} or newer from{' '}
            <a
              href="https://go.dev/dl/"
              className="text-accent hover:underline"
              target="_blank"
              rel="noreferrer"
            >
              go.dev/dl
            </a>
            . Then <strong>open a new terminal</strong> and run <code>go version</code>.
          </p>
          <CodeBlock
            terminal
            code={`# macOS, with Homebrew
brew install go

# Debian and Ubuntu: the apt package is usually years old.
# Download the tarball from go.dev/dl instead, then:
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go*.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile

# Windows: use the .msi installer, which sets PATH for you.`}
          />
          <p>
            If <code>go version</code> works but a tool you installed with <code>go install</code>{' '}
            does not, that is a different PATH: add <code>$(go env GOPATH)/bin</code>, which is{' '}
            <code>%USERPROFILE%\go\bin</code> on Windows.
          </p>
        </Mistake>

        <Mistake error="go.mod file not found in current directory or any parent directory">
          <p>
            You are outside a Go module. Either you are in the wrong folder, or the module was
            never created.
          </p>
          <CodeBlock
            terminal
            code={`# Are you where you think you are?
pwd && ls go.mod

# No go.mod anywhere: create the module
go mod init example.com/myapp`}
          />
          <p>
            In a project with the API in a subdirectory, which is the usual layout, the Go commands
            run from there: <code>cd apps/api</code> first.
          </p>
        </Mistake>

        <Mistake error="package X is not in GOROOT">
          <p>
            Go looked for your import among the standard library and did not find it. Nearly always
            a dependency that was never downloaded, or an import path with a typo.
          </p>
          <CodeBlock terminal code={`go mod tidy`} />
          <p>
            That reads every import in your code, fetches what is missing, and removes what is no
            longer used. If it still fails, the import path itself is wrong: compare it against the
            package&apos;s own documentation, since the import path and the repository URL are not
            always the same string.
          </p>
        </Mistake>

        <Mistake error="missing go.sum entry; to add it: go mod download">
          <p>
            The module graph changed and <code>go.sum</code>, which records the checksum of every
            dependency, did not keep up. Usually after editing <code>go.mod</code> by hand or
            merging a branch.
          </p>
          <CodeBlock terminal code={`go mod tidy`} />
          <p>
            Commit the resulting <code>go.sum</code>. It is not noise: it is what stops a
            dependency being swapped under you.
          </p>
        </Mistake>

        <Mistake error="dial tcp: lookup proxy.golang.org: no such host">
          <p>
            Go cannot reach the module proxy. No network, DNS blocked, or a corporate proxy
            intercepting it.
          </p>
          <CodeBlock
            terminal
            code={`# Behind a company proxy with its own mirror
go env -w GOPROXY=https://your-mirror.example.com,direct

# Or go direct to the source hosts
go env -w GOPROXY=direct GOSUMDB=off`}
          />
          <p>
            Turning off <code>GOSUMDB</code> disables checksum verification, so use it to get
            unblocked and not as a permanent setting. On a VPN, try it disconnected first: that is
            the single most common cause.
          </p>
        </Mistake>

        <Mistake error="Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work">
          <p>
            The standard SQLite driver is a wrapper around C code and needs a C compiler. You get
            this on a machine without one, which on Windows is most machines.
          </p>
          <p>Either install a compiler:</p>
          <CodeBlock
            terminal
            code={`# Debian and Ubuntu
sudo apt install build-essential

# macOS
xcode-select --install

# Windows: install MSYS2 or TDM-GCC, then confirm
gcc --version`}
          />
          <p>Or use the pure Go driver, which needs no compiler and has the same API:</p>
          <CodeBlock
            terminal
            code={`go get github.com/glebarez/sqlite
# then change the import from gorm.io/driver/sqlite`}
          />
        </Mistake>

        <Mistake error="listen tcp :8080: bind: address already in use">
          <p>
            Another process holds the port. Almost always a server you started earlier and did not
            stop, often in a terminal tab you closed without stopping the process.
          </p>
          <CodeBlock
            terminal
            code={`# macOS and Linux
lsof -ti:8080 | xargs kill

# Windows PowerShell
Get-NetTCPConnection -LocalPort 8080 -State Listen |
  Select-Object -Expand OwningProcess |
  ForEach-Object { Stop-Process -Id $_ -Force }`}
          />
          <p>
            On macOS, port 5000 is held by AirPlay Receiver and port 7000 by the Control Center.
            Those are not stray processes of yours, and killing them is not the fix: pick another
            port, or turn the feature off in System Settings.
          </p>
        </Mistake>

        <Mistake error="panic: runtime error: invalid memory address or nil pointer dereference">
          <p>
            Something you called a method on was <code>nil</code>. In a new project this is almost
            always a database handle or a service you forgot to construct before using, because the
            code compiles perfectly well either way.
          </p>
          <p>
            Read the stack trace from the <strong>bottom</strong> up: the first line in your own
            package, not in a library, is where to look. The variable on that line is the nil one.
          </p>
          <Note title="The habit that prevents it">
            <p>
              Return an error from a constructor rather than returning a half-built value. A{' '}
              <code>NewService()</code> that cannot fail will happily hand you a struct whose
              database field is nil, and the crash arrives later, somewhere else.
            </p>
          </Note>
        </Mistake>
      </Section>

      <Section id="docker" title="Docker">
        <Mistake error="Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?">
          <p>
            On Windows the same thing reads{' '}
            <code>error during connect: open //./pipe/dockerDesktopLinuxEngine: The system cannot
            find the file specified</code>
            .
          </p>
          <p>
            Docker Desktop is not running. Open it and wait until the icon settles. On Linux, where
            there is no Desktop app, start the service:
          </p>
          <CodeBlock
            terminal
            code={`sudo systemctl start docker
sudo systemctl enable docker   # and at every boot`}
          />
        </Mistake>

        <Mistake error="permission denied while trying to connect to the Docker daemon socket">
          <p>
            Linux only. Your user is not in the <code>docker</code> group, so only{' '}
            <code>sudo docker</code> works.
          </p>
          <CodeBlock
            terminal
            code={`sudo usermod -aG docker $USER
# then log out and back in, or:
newgrp docker`}
          />
          <p>
            The group change applies to new sessions, which is why people add themselves and then
            report that it did not work.
          </p>
        </Mistake>

        <Mistake error="Bind for 0.0.0.0:5432 failed: port is already allocated">
          <p>
            Something already listens on that port: usually a Postgres you installed natively
            months ago and forgot, or another project&apos;s containers still up.
          </p>
          <CodeBlock
            terminal
            code={`# Which container has it?
docker ps --format "table {{.Names}}\\t{{.Ports}}"

# Another project's stack, still running
docker compose -f ../other-project/docker-compose.yml down`}
          />
          <p>
            If it is a native Postgres, either stop the service or move the container to a different
            host port. Change only the left-hand number, which is the host side:
          </p>
          <CodeBlock
            language="yaml"
            filename="docker-compose.yml"
            code={`services:
  postgres:
    ports:
      - "5433:5432"   # host 5433, container still 5432`}
          />
          <p>
            Then update the port in your <code>.env</code> too. Changing one and not the other is
            the follow-up error.
          </p>
        </Mistake>

        <Mistake error="docker: 'compose' is not a docker command">
          <p>
            You have an old Docker, where Compose was a separate program called{' '}
            <code>docker-compose</code> with a hyphen. Modern Docker has it built in as{' '}
            <code>docker compose</code> with a space.
          </p>
          <p>
            Upgrade Docker, or use the hyphenated form. Do not install both and alternate between
            them: they keep separate state and will confuse you about what is running.
          </p>
        </Mistake>

        <Mistake error="no space left on device">
          <p>
            Docker accumulates images, stopped containers, build caches and dangling volumes. On a
            development machine it can reach tens of gigabytes without you noticing.
          </p>
          <CodeBlock
            terminal
            code={`docker system df        # where it went
docker system prune     # stopped containers, unused networks, dangling images`}
          />
          <p>
            <strong>Do not run <code>docker system prune --volumes</code> casually.</strong> Volumes
            are where your database lives, and that flag deletes every one not attached to a running
            container. A stopped project loses its data.
          </p>
        </Mistake>

        <Mistake error="database system is starting up / connection refused on first run">
          <p>
            Your API started before Postgres finished initialising. Containers start in parallel,
            and a container being <em>up</em> is not the same as the service inside it being{' '}
            <em>ready</em>.
          </p>
          <p>
            The fix is a healthcheck and a condition, not a sleep. <code>depends_on</code> alone
            waits for the container, not the database:
          </p>
          <CodeBlock
            language="yaml"
            filename="docker-compose.yml"
            code={`services:
  postgres:
    image: postgres:17-alpine
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 2s
      timeout: 3s
      retries: 15

  api:
    depends_on:
      postgres:
        condition: service_healthy`}
          />
        </Mistake>
      </Section>

      <Section id="node" title="Node, npm and pnpm">
        <Mistake error="pnpm: command not found">
          <p>
            pnpm is not installed. The least surprising way is Corepack, which ships with Node and
            installs the exact version a project asks for.
          </p>
          <CodeBlock
            terminal
            code={`corepack enable
corepack prepare pnpm@${PINNED_PNPM} --activate

# or, without Corepack
npm install -g pnpm@${PINNED_PNPM}`}
          />
          <p>
            Pin the version rather than taking the newest. A major pnpm release changes how the
            lockfile and <code>node_modules</code> are laid out, and a project installed with two
            different majors behaves differently on two machines.
          </p>
        </Mistake>

        <Mistake error="Unsupported engine: wanted node >=20.9.0, current node v18.x">
          <p>
            Your Node is older than the tooling accepts. Install Node {REQUIRED_NODE} or newer. The
            least painful route is a version manager, because you will eventually need two versions
            on one machine.
          </p>
          <CodeBlock
            terminal
            code={`# nvm, macOS and Linux
nvm install 22 && nvm use 22

# fnm, which is faster and works on Windows too
fnm install 22 && fnm use 22`}
          />
        </Mistake>

        <Mistake error="Module not found: Can't resolve '@/components/...'">
          <p>
            The <code>@</code> alias is a TypeScript path mapping, not something Node knows about.
            Either <code>tsconfig.json</code> does not define it, or the file genuinely is not
            there.
          </p>
          <CodeBlock
            language="json"
            filename="tsconfig.json"
            code={`{
  "compilerOptions": {
    "paths": {
      "@/*": ["./*"]
    }
  }
}`}
          />
          <p>
            Check the case of the filename as well. macOS and Windows do not care about case and
            Linux does, so an import of <code>Button</code> against a file named{' '}
            <code>button.tsx</code> works on your laptop and fails in CI, which is a miserable way
            to find out.
          </p>
        </Mistake>

        <Mistake error="Error: Cannot find module 'next' / ERR_MODULE_NOT_FOUND">
          <p>
            Dependencies are not installed, or the install was interrupted and left{' '}
            <code>node_modules</code> half written. Start again from clean:
          </p>
          <CodeBlock
            terminal
            code={`rm -rf node_modules apps/*/node_modules
rm -f pnpm-lock.yaml        # only if the lockfile itself is suspect
pnpm install`}
          />
          <p>
            Deleting the lockfile is a bigger step than it looks: it lets every dependency move to
            a newer version at once. Try the install without deleting it first.
          </p>
        </Mistake>

        <Mistake error="EADDRINUSE: address already in use :::3000">
          <p>
            The Node version of the port problem, usually a dev server that survived a closed
            terminal. Kill it as above, or start on another port:
          </p>
          <CodeBlock terminal code={`pnpm dev --port 3005`} />
        </Mistake>

        <Mistake error="Hydration failed because the server rendered HTML didn't match the client">
          <p>
            Something rendered differently on the server and in the browser. Three causes cover
            nearly all of it: <code>new Date()</code> or <code>Math.random()</code> called during
            render, reading <code>window</code> or <code>localStorage</code> during the first
            render, and invalid HTML nesting such as a <code>div</code> inside a <code>p</code>,
            which the browser silently restructures.
          </p>
          <p>
            For anything that genuinely only exists in the browser, read it in an effect after
            mounting rather than during render.
          </p>
        </Mistake>
      </Section>

      <Section id="grit" title="If you are using Grit">
        <Mistake error="grit: command not found">
          <p>
            The binary is not on your PATH. The installer puts it in <code>~/.grit/bin</code>, which
            your shell only knows about after it has been added and the terminal reopened.
          </p>
          <CodeBlock
            terminal
            code={`echo 'export PATH=$PATH:$HOME/.grit/bin' >> ~/.zshrc   # or ~/.bashrc
source ~/.zshrc
grit version`}
          />
          <p>
            If <code>grit version</code> prints an older version than you installed, you have two
            copies and the wrong one is first on PATH. <code>which -a grit</code>, or{' '}
            <code>where grit</code> on Windows, lists every one it can see.
          </p>
        </Mistake>

        <Mistake error="This does not look like a Grit project">
          <p>
            Grit commands read <code>grit.json</code> at the project root. You are in a
            subdirectory, or in the wrong project. Move to the directory that contains it.
          </p>
        </Mistake>

        <p className="mt-8">
          For the Grit-specific failures, generation, migrations, the admin build, the{' '}
          <Link
            href="/docs/getting-started/troubleshooting"
            className="text-accent hover:underline"
          >
            troubleshooting page in the documentation
          </Link>{' '}
          goes further, and <code>grit doctor</code> checks a project against the things that are
          commonly wrong without touching the network.
        </p>
        <CodeBlock terminal code={`grit doctor`} />
      </Section>

      <Faq
        slug="common-go-docker-errors"
        items={[
          {
            q: 'Why does reinstalling the tool not fix "command not found"?',
            a: (
              <>
                Because the problem is usually PATH, not the installation. Your shell reads PATH
                once when it starts, so a terminal that was already open keeps the old value
                forever. Close it, open a new one, and check again before reinstalling anything.
              </>
            ),
            text:
              'Because the problem is usually PATH rather than the installation. A shell reads PATH when it starts, so a terminal that was already open keeps the old value. Open a new terminal and check again before reinstalling.',
          },
          {
            q: 'Is it safe to run docker system prune?',
            a: (
              <>
                Without flags, yes: it removes stopped containers, unused networks and dangling
                images, all of which rebuild. With <code>--volumes</code> it also deletes volumes
                that no running container is using, and that is where your database lives. A
                project you stopped yesterday loses its data.
              </>
            ),
            text:
              'Without flags, yes: it removes stopped containers, unused networks and dangling images, all of which can be rebuilt. With --volumes it also deletes volumes not attached to a running container, which is where databases live, so a stopped project loses its data.',
          },
          {
            q: 'My port is in use but nothing is running. What is holding it?',
            a: (
              <>
                A process you cannot see in your terminal: a dev server whose window you closed
                without stopping it, a Docker container, or an operating system service. On macOS,
                ports 5000 and 7000 belong to AirPlay and Control Center. List what is listening
                with <code>lsof -i:PORT</code>, or <code>Get-NetTCPConnection -LocalPort PORT</code>{' '}
                in PowerShell.
              </>
            ),
            text:
              'A process you cannot see: a dev server whose terminal was closed without stopping it, a Docker container, or an OS service. On macOS, ports 5000 and 7000 belong to AirPlay and Control Center. List listeners with lsof -i:PORT, or Get-NetTCPConnection -LocalPort PORT in PowerShell.',
          },
          {
            q: 'Should I develop inside WSL on Windows?',
            a: (
              <>
                If you are comfortable with it, yes. File watching and Docker are both faster, and
                most tutorials assume a Unix shell. Keep your code inside the Linux filesystem,
                under <code>~</code>, rather than under <code>/mnt/c</code>: crossing the boundary
                makes every file read slow enough to notice, and file watching unreliable.
              </>
            ),
            text:
              'If you are comfortable with it, yes: file watching and Docker are faster and most tutorials assume a Unix shell. Keep the code inside the Linux filesystem under your home directory rather than under /mnt/c, because crossing that boundary makes file reads slow and file watching unreliable.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/first-go-react-app',
            label: 'Build your first full-stack app',
            blurb: 'Go, React and the agreement between them, with a checkpoint at every step.',
          },
          {
            href: '/guides/docker-compose-postgres-redis',
            label: 'Postgres, Redis and MinIO with Docker Compose',
            blurb: 'What the compose file says, line by line, and the Windows-specific traps.',
          },
          {
            href: '/docs/getting-started/troubleshooting',
            label: 'Grit troubleshooting',
            blurb: 'The failures specific to generation, migrations and the admin build.',
          },
          {
            href: '/docs/cli#doctor',
            label: 'grit doctor',
            blurb: 'An offline check of the things that are commonly wrong in a project.',
          },
        ]}
      />
    </GuidePage>
  )
}
