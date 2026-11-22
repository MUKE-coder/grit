import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  Step,
  H3,
  Outcome,
  Note,
  Warn,
  Mistake,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Run Postgres, Redis and MinIO locally with Docker Compose',
  description:
    'What a container actually is, a compose file explained line by line, named volumes, healthchecks, and the port collisions that catch everyone on Windows.',
  alternates: { canonical: 'https://gritframework.dev/guides/docker-compose-postgres-redis' },
  openGraph: {
    title: 'Run Postgres, Redis and MinIO locally with Docker Compose',
    description: 'A compose file for local development, explained line by line.',
    url: 'https://gritframework.dev/guides/docker-compose-postgres-redis',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="docker-compose-postgres-redis"
      lead={
        <>
          <p className="mb-4">
            Installing Postgres on your laptop works until the second project, which wants a
            different version, or until you wipe the machine and cannot remember what you
            configured. Docker Compose replaces that with one file you can read, delete and
            recreate, and which every person on the project runs identically.
          </p>
          <p>
            This explains a working compose file line by line: Postgres for data, Redis for cache
            and queues, MinIO as a local stand-in for S3, and Mailhog so you can see your
            application&apos;s email without sending any.
          </p>
        </>
      }
    >
      <Outcome>
        <p>
          Four services running locally with one command, and gone again with another, leaving
          nothing installed on your machine.
        </p>
        <p>
          Data that survives restarts, because the volumes are named rather than anonymous, and an
          understanding of why that distinction is the one that loses people&apos;s databases.
        </p>
        <p>Ports chosen so they do not fight whatever you already have installed.</p>
      </Outcome>

      <Section id="concept" title="What a container actually is">
        <p>
          Not a virtual machine. A container is a normal process on your kernel, with its own view
          of the filesystem, the network and the process list. That is why it starts in a second
          rather than a minute, and why a Postgres container uses roughly the memory Postgres uses.
        </p>
        <p>Three words you need, and they are the ones people conflate.</p>
        <Checklist
          items={[
            <>
              An <strong>image</strong> is a read-only template: Postgres 16, built and published by
              somebody. You never change an image.
            </>,
            <>
              A <strong>container</strong> is a running instance of an image. Delete it and you lose
              everything it wrote, unless that something was written to a volume.
            </>,
            <>
              A <strong>volume</strong> is storage that outlives the container. This is where your
              database actually lives.
            </>,
          ]}
        />
        <p>
          Compose is a file that describes several containers and how they are wired together, so
          that <code>docker compose up</code> starts the lot in the right order.
        </p>
      </Section>

      <Step n={1} title="Install Docker">
        <p>
          <strong>macOS and Windows:</strong> Docker Desktop from{' '}
          <a
            href="https://docs.docker.com/get-started/get-docker/"
            className="text-accent hover:underline"
            target="_blank"
            rel="noreferrer"
          >
            docker.com
          </a>
          . It is an application, and it must be open for any <code>docker</code> command to work.
        </p>
        <p>
          <strong>Windows specifically:</strong> enable WSL2 when the installer offers it. The
          alternative, Hyper-V, is slower and has more edge cases. If the installer complains that
          virtualisation is disabled, it is a BIOS setting rather than anything in Windows.
        </p>
        <p>
          <strong>Linux:</strong> install Docker Engine from your distribution&apos;s instructions,
          not the <code>docker.io</code> package, which is usually behind. Then add yourself to the
          docker group so you are not typing <code>sudo</code> all day:
        </p>
        <CodeBlock
          terminal
          code={`sudo usermod -aG docker $USER
newgrp docker`}
        />
        <p>
          <strong>Checkpoint.</strong> <code>docker run --rm hello-world</code> prints a greeting
          and exits. If it does, everything below will work.
        </p>
      </Step>

      <Step n={2} title="Write the compose file">
        <p>
          One file at the root of your project. Read the comments: every one of them is a decision
          somebody made after being bitten.
        </p>
        <CodeBlock
          language="yaml"
          filename="docker-compose.yml"
          code={`services:
  postgres:
    image: postgres:16-alpine
    restart: unless-stopped
    ports:
      # Host 5434, not the Postgres default 5432. Two collisions this avoids:
      # a Postgres installed on the host, and Windows reserving 5432 inside
      # its dynamic Hyper-V port range. The container still listens on 5432
      # inside the Docker network, which is what other containers use.
      - "127.0.0.1:\${POSTGRES_PORT:-5434}:5432"
    environment:
      # :- supplies a fallback, so the file still works with no .env
      POSTGRES_USER: \${POSTGRES_USER:-app}
      POSTGRES_PASSWORD: \${POSTGRES_PASSWORD:-app}
      POSTGRES_DB: \${POSTGRES_DB:-app}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U \${POSTGRES_USER:-app}"]
      interval: 5s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    restart: unless-stopped
    ports:
      # Host 6380, because a native Redis or Memurai install takes 6379.
      - "127.0.0.1:\${REDIS_PORT:-6380}:6379"
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  minio:
    image: minio/minio
    restart: unless-stopped
    ports:
      # Host 9002 and 9003. Portainer and several monitoring stacks take
      # 9000 and 9001.
      - "127.0.0.1:9002:9000"
      - "127.0.0.1:9003:9001"
    environment:
      # :? refuses to start with a clear message rather than booting with a
      # default password, which is the right behaviour for a credential.
      MINIO_ROOT_USER: \${MINIO_ACCESS_KEY:?set MINIO_ACCESS_KEY in .env}
      MINIO_ROOT_PASSWORD: \${MINIO_SECRET_KEY:?set MINIO_SECRET_KEY in .env}
    volumes:
      - minio-data:/data
    command: server /data --console-address ":9001"

  mailhog:
    image: mailhog/mailhog
    restart: unless-stopped
    ports:
      - "127.0.0.1:1025:1025"   # SMTP, what your app connects to
      - "127.0.0.1:8025:8025"   # the web inbox you read

volumes:
  postgres-data:
  redis-data:
  minio-data:`}
        />
        <p>And the credentials beside it, which is the file you never commit:</p>
        <CodeBlock
          language="bash"
          filename=".env"
          code={`POSTGRES_USER=app
POSTGRES_PASSWORD=choose-something
POSTGRES_DB=app
POSTGRES_PORT=5434

REDIS_PORT=6380

MINIO_ACCESS_KEY=local-dev-key
MINIO_SECRET_KEY=local-dev-secret`}
        />
        <Warn title="Add .env to .gitignore before the first commit">
          <p>
            Once a credential is in git history, removing it from the current file does nothing:
            the old commit still has it, and on a public repository it is scraped within minutes.
            Getting it out afterwards means rewriting history for everybody.
          </p>
        </Warn>
      </Step>

      <Step n={3} title="Start it, and check it">
        <CodeBlock
          terminal
          code={`docker compose up -d        # -d detaches, so you get your terminal back
docker compose ps           # what is running, and whether it is healthy`}
        />
        <p>
          <strong>Checkpoint.</strong> <code>docker compose ps</code> lists four services, and
          Postgres and Redis say <code>healthy</code> rather than just <code>running</code>. If
          Postgres says <code>starting</code> for more than ten seconds, read its logs:
        </p>
        <CodeBlock
          terminal
          code={`docker compose logs -f postgres    # -f follows; Ctrl-C stops following,
                                   # not the container`}
        />
        <p>Then prove you can actually talk to each one:</p>
        <CodeBlock
          terminal
          code={`# Postgres: a shell inside the container
docker compose exec postgres psql -U app -d app -c '\\l'

# Redis
docker compose exec redis redis-cli ping     # PONG

# MinIO console, in a browser
open http://localhost:9003                   # xdg-open on Linux, start on Windows

# Mailhog inbox
open http://localhost:8025`}
        />
      </Step>

      <Section id="lines" title="The four lines that matter most">
        <H3>
          <code>127.0.0.1:5434:5432</code>
        </H3>
        <p>
          Three parts: bind address, host port, container port. Writing <code>5434:5432</code>{' '}
          without the address binds to <strong>every</strong> interface, which means your database
          is reachable by anyone on the coffee shop wifi, with the password you chose while in a
          hurry. The <code>127.0.0.1</code> prefix costs nothing and closes that.
        </p>
        <p>
          Containers do not need this at all to reach each other. Inside the Docker network they
          use the service name and the real port: <code>postgres:5432</code>. The ports line exists
          only so that programs on your machine, your API running outside Docker, a database GUI,
          can get in.
        </p>
        <H3>
          <code>postgres-data:/var/lib/postgresql/data</code>
        </H3>
        <p>
          A <strong>named</strong> volume. The name on the left is declared in the{' '}
          <code>volumes:</code> block at the bottom, and it is what makes the data survive{' '}
          <code>docker compose down</code> and come back on the next <code>up</code>.
        </p>
        <p>
          Omit the name and you get an anonymous volume, which survives too, but under a
          random identifier that nothing reconnects to. That is how people conclude Docker loses
          their data: it did not, they just cannot find it.
        </p>
        <H3>
          <code>healthcheck</code>
        </H3>
        <p>
          A container being up does not mean the service inside it is ready. Postgres spends a few
          seconds initialising on first run, and an API that connects during that window fails with{' '}
          <code>connection refused</code> and dies. The healthcheck is what lets other services
          wait for <em>ready</em> rather than <em>started</em>:
        </p>
        <CodeBlock
          language="yaml"
          code={`  api:
    depends_on:
      postgres:
        condition: service_healthy`}
        />
        <p>
          <code>depends_on</code> on its own, without the condition, waits only for the container
          to exist, which is almost never what you wanted.
        </p>
        <H3>
          <code>{'${POSTGRES_PASSWORD:-app}'}</code> against{' '}
          <code>{'${MINIO_SECRET_KEY:?set it}'}</code>
        </H3>
        <p>
          The first supplies a default. The second refuses to start, with your message, when the
          variable is missing. Use the second for anything that is a credential in production: a
          default password is a password somebody will still be using in a year.
        </p>
      </Section>

      <Section id="daily" title="The six commands you will actually use">
        <CodeBlock
          terminal
          code={`docker compose up -d           # start everything
docker compose ps              # what is running and how healthy
docker compose logs -f api     # follow one service's output
docker compose restart redis   # restart one service
docker compose down            # stop and remove containers, KEEP volumes
docker compose down -v         # ... and delete the volumes too`}
        />
        <Note title="down against down -v">
          <p>
            <code>down</code> is safe and routine: it stops everything and leaves your data.{' '}
            <code>down -v</code> deletes the volumes, which means the database. It is the right
            command when you want a genuinely fresh start, and the wrong one by accident. There is
            no confirmation prompt.
          </p>
        </Note>
      </Section>

      <Section id="mistakes" title="When it goes wrong">
        <Mistake error="Bind for 0.0.0.0:5432 failed: port is already allocated">
          <p>
            Something else holds the port. Usually a Postgres installed natively, or another
            project&apos;s compose stack still up.
          </p>
          <CodeBlock
            terminal
            code={`docker ps --format "table {{.Names}}\\t{{.Ports}}"`}
          />
          <p>
            Change only the host side of the mapping, the left-hand number, and update the matching
            variable in <code>.env</code>. Changing one without the other is the next error you
            will see.
          </p>
        </Mistake>

        <Mistake error="Cannot connect to the Docker daemon. Is the docker daemon running?">
          <p>
            Docker Desktop is not open, or on Linux the service is not started. The message never
            says so plainly, which is why it is the most reported Docker error there is.
          </p>
        </Mistake>

        <Mistake error="dial tcp 127.0.0.1:5432: connect: connection refused (from inside a container)">
          <p>
            Inside a container, <code>127.0.0.1</code> means <em>that container</em>, not your
            machine. A service looking for the database at localhost is looking at itself.
          </p>
          <p>
            Use the service name: the connection string becomes{' '}
            <code>postgres://app:app@postgres:5432/app</code>. And note the port is the container
            port 5432, not the host port 5434, because this connection never leaves the Docker
            network.
          </p>
        </Mistake>

        <Mistake error="role &quot;app&quot; does not exist / database &quot;app&quot; does not exist">
          <p>
            The Postgres image creates the user and database from the environment variables{' '}
            <strong>only on the very first start</strong>, when the data directory is empty. If you
            started the container once with different credentials, the volume already holds a
            database that does not know about your new ones, and changing the compose file does
            nothing.
          </p>
          <CodeBlock
            terminal
            code={`docker compose down -v        # deletes the volume, so the next start initialises
docker compose up -d`}
          />
          <p>That deletes the data in it. On a local development database that is usually fine.</p>
        </Mistake>

        <Mistake error="no space left on device">
          <p>
            Images and build caches accumulate. <code>docker system df</code> shows where it went
            and <code>docker system prune</code> reclaims the safe parts. Leave{' '}
            <code>--volumes</code> off unless you mean to delete databases.
          </p>
        </Mistake>

        <Mistake error="Everything is extremely slow on Windows or macOS">
          <p>
            Almost always a bind mount crossing the boundary between the host filesystem and the
            container. On Windows, keep the project inside the WSL2 filesystem, under{' '}
            <code>~</code>, rather than under <code>/mnt/c</code>. On macOS, enable VirtioFS in
            Docker Desktop&apos;s settings if it is not already on.
          </p>
          <p>
            This does not apply to the file above, which uses named volumes rather than bind
            mounts. It bites when you start mounting your source code into a container.
          </p>
        </Mistake>
      </Section>

      <Section id="with-grit" title="The same file, written for you">
        <p>
          A Grit project ships this compose file already, with the ports chosen as above, the
          healthchecks wired, and the application&apos;s environment variables pointing at the
          right service names.
        </p>
        <CodeBlock
          terminal
          code={`grit new myapp
cd myapp
docker compose up -d
grit migrate
grit start`}
        />
        <p>
          The generated file has one difference worth knowing about: nothing in the production
          compose file publishes a port at all. Databases are reachable only from inside the
          Docker network, and you reach them over an SSH tunnel. Development is the only place a
          database port is on your machine, which is the right way round.
        </p>
      </Section>

      <Faq
        slug="docker-compose-postgres-redis"
        items={[
          {
            q: 'Do I need Docker to develop locally?',
            a: (
              <>
                No. You can install Postgres and Redis natively, and for a single project that is
                perfectly workable. Docker earns its place at the second project, at the first
                version conflict, and on the first day somebody new joins and has to reproduce your
                setup from a README.
              </>
            ),
            text:
              'No. Installing Postgres and Redis natively works fine for a single project. Docker earns its place at the second project, at the first version conflict, and the first time somebody else has to reproduce your setup.',
          },
          {
            q: 'Will docker compose down delete my database?',
            a: (
              <>
                Not with named volumes, which is what the file above uses: <code>down</code> removes
                the containers and leaves the data, and the next <code>up</code> finds it again.{' '}
                <code>down -v</code> does delete it, deliberately and without asking.
              </>
            ),
            text:
              'Not with named volumes, which the file above uses. docker compose down removes containers and leaves the data, and the next up finds it. docker compose down -v deletes the volumes, deliberately and with no confirmation prompt.',
          },
          {
            q: 'Why does my API say connection refused when the database is clearly running?',
            a: (
              <>
                Two likely causes. If the API runs inside Docker, it must use the service name and
                container port (<code>postgres:5432</code>), because <code>127.0.0.1</code> inside a
                container means that container. If the API runs on your machine, it uses the host
                port (<code>localhost:5434</code>). Mixing the two up is the single most common
                Compose mistake.
              </>
            ),
            text:
              'Two likely causes. An API running inside Docker must use the service name and container port, postgres:5432, because 127.0.0.1 inside a container refers to that container. An API running on the host uses the published host port, localhost:5434. Mixing the two is the most common Compose mistake.',
          },
          {
            q: 'Is MinIO good enough to stand in for S3?',
            a: (
              <>
                For development, yes: it speaks the S3 API, so the same SDK and the same code work
                against both, and you change an endpoint and some credentials to move to the real
                thing. The differences that matter appear at scale and around IAM policies, neither
                of which you are exercising locally.
              </>
            ),
            text:
              'For development, yes. MinIO speaks the S3 API, so the same SDK and application code work against both, and moving to real S3, R2 or B2 is an endpoint and a credential change. The differences appear at scale and around IAM policies, which local development does not exercise.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/common-go-docker-errors',
            label: 'Fix the errors you hit starting out',
            blurb: 'The exact messages, with the cause and the fix per operating system.',
          },
          {
            href: '/guides/first-go-react-app',
            label: 'Build your first full-stack app',
            blurb: 'A Go API and a React frontend, with a checkpoint at every step.',
          },
          {
            href: '/docs/infrastructure/docker',
            label: 'Docker in a Grit project',
            blurb: 'The generated development and production stacks, and what differs between them.',
          },
          {
            href: '/docs/deployment',
            label: 'Deploying',
            blurb: 'What changes between the compose file you run locally and the one on a server.',
          },
        ]}
      />
    </GuidePage>
  )
}
