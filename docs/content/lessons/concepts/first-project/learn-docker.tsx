import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Before we boot the dev servers, let&apos;s deal with the elephant in
        the room: <strong>Docker</strong>. Most people get to the{' '}
        <code>docker compose up -d</code> command in the next lesson and stall
        because nobody actually taught them what Docker is. This lesson
        teaches it from first principles. By the end you&apos;ll know what
        Docker is, why it exists, how Grit uses it, how to skip it entirely
        if you want to — and why you should learn it anyway.
      </p>

      <h2>The problem Docker solves</h2>
      <p>
        Software has dependencies. Your app needs Postgres v16, Redis v7, a
        specific version of Node, a specific Go toolchain. On your laptop
        everything works. You hand the project to a teammate — their
        Postgres is v14 and migrations fail. You deploy to a server — its
        Redis is missing a feature. You spend a week chasing
        &quot;works-on-my-machine&quot; bugs.
      </p>
      <p>
        Docker fixes this by bundling each piece of software with its exact
        environment into a <strong>container</strong> — a lightweight,
        self-contained box. Postgres v16 in a container behaves identically
        on your Mac, your teammate&apos;s Windows machine, and the
        production Linux server. The whole &quot;but it worked on my
        machine&quot; era largely ends with this idea.
      </p>

      <h2>The three concepts you need today</h2>
      <Diagram label="Images, containers, volumes" caption="Image = recipe. Container = running instance. Volume = data that survives container restarts.">
{`   ┌─────────────────────────────────────────────────────────┐
   │  IMAGE                                                  │
   │  A read-only template. "postgres:16-alpine" is an image.│
   │  You download it once.                                  │
   └────────────────────┬────────────────────────────────────┘
                        │ docker run / docker compose up
                        ▼
   ┌─────────────────────────────────────────────────────────┐
   │  CONTAINER                                              │
   │  A live, running instance of an image. You can have many│
   │  containers from the same image (3 postgres containers, │
   │  each with their own data).                             │
   └────────────────────┬────────────────────────────────────┘
                        │ writes to
                        ▼
   ┌─────────────────────────────────────────────────────────┐
   │  VOLUME                                                 │
   │  A named, persistent disk attached to a container. When │
   │  you delete the container, the volume's data stays.     │
   │  postgres-data, redis-data are volumes.                 │
   └─────────────────────────────────────────────────────────┘`}
      </Diagram>
      <p>That&apos;s the whole vocabulary for the next 6 months.</p>

      <h2>Installing Docker</h2>
      <ul>
        <li>
          <strong>Mac / Windows:</strong> Install{' '}
          <a href="https://www.docker.com/products/docker-desktop/" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
            Docker Desktop
          </a>
          . It bundles the Docker engine + a GUI + Compose. After install,
          open the Docker Desktop app — it stays running in your menu / system
          tray.
        </li>
        <li>
          <strong>Linux:</strong> Use{' '}
          <a href="https://docs.docker.com/engine/install/" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">
            Docker Engine
          </a>{' '}
          (apt / yum / dnf). No GUI, just the daemon. Add yourself to the{' '}
          <code>docker</code> group so you don&apos;t need <code>sudo</code>{' '}
          for every command:{' '}
          <code>sudo usermod -aG docker $USER</code> then log out and back in.
        </li>
      </ul>
      <p>Verify the install:</p>
      <CodeBlock
        language="bash"
        code={`docker --version
# Docker version 27.x.x

docker compose version
# Docker Compose version v2.x.x

docker run --rm hello-world
# pulls a tiny test image and runs it — confirms the daemon works`}
      />

      <h2>How Grit uses Docker</h2>
      <p>
        Grit&apos;s <code>docker-compose.yml</code> describes 4 services your
        app depends on locally:
      </p>
      <ul>
        <li>
          <strong>Postgres</strong> on port 5432 — the database.
        </li>
        <li>
          <strong>Redis</strong> on port 6379 — cache + job queue.
        </li>
        <li>
          <strong>MinIO</strong> on ports 9000 (S3 API) + 9001 (web console)
          — S3-compatible file storage.
        </li>
        <li>
          <strong>Mailhog</strong> on ports 1025 (SMTP) + 8025 (web UI) — fake
          mail server so you can test emails without sending real ones.
        </li>
      </ul>
      <p>
        <code>docker compose up -d</code> boots all four in the background.
        <code> docker compose down</code> stops them. You never had to install
        Postgres or Redis on your machine — they live in containers and your
        Go API connects to them via <code>localhost:5432</code> +{' '}
        <code>localhost:6379</code>. When you&apos;re done for the day, run{' '}
        <code>down</code> and your laptop is back to normal.
      </p>

      <TipBox tone="info">
        <strong>Grit binds these ports to 127.0.0.1 only.</strong> That means
        only your laptop can connect to them — not anything else on your
        wifi. The default Docker behaviour is to bind to <code>0.0.0.0</code>{' '}
        (every interface). On a coffee-shop wifi that exposes Postgres with{' '}
        <code>grit:grit</code> credentials to the whole room. We pinned it to
        localhost specifically to prevent that.
      </TipBox>

      <h2>Docker commands you&apos;ll actually type</h2>
      <CodeBlock
        language="bash"
        code={`# Start everything in docker-compose.yml in the background
docker compose up -d

# Stop everything (preserves data — volumes stay)
docker compose down

# Stop AND wipe data (use this when you want a fresh DB)
docker compose down -v

# Stop EVERY running container on your machine — even ones from other
# projects you forgot were running. Useful when port 5432 / 6379 is
# already taken and you can't remember by what.
docker stop $(docker ps -q)

# See what's running
docker ps
docker compose ps

# Live logs from one service (Ctrl+C to exit; data keeps flowing)
docker compose logs -f postgres

# Open a shell INSIDE a running container (great for debugging)
docker compose exec postgres bash
docker compose exec postgres psql -U grit -d my-first-grit

# Restart one service after editing its config
docker compose restart redis

# Pull newer image versions (e.g., postgres 16.2 → 16.3)
docker compose pull

# See disk usage — containers + volumes + images
docker system df

# Reclaim disk from stopped containers + dangling images (safe)
docker system prune

# Nuclear: also delete unused volumes (DESTROYS DATA — confirm twice)
docker system prune --volumes`}
      />
      <p>
        These thirteen cover 99% of what you&apos;ll do. Bookmark this lesson —
        when in doubt, scroll up.
      </p>

      <h2>The errors you WILL hit, and what they mean</h2>

      <h3>1. <code>Cannot connect to the Docker daemon</code></h3>
      <p>
        Docker Desktop isn&apos;t running (Mac/Windows) or the daemon
        isn&apos;t started (Linux: <code>sudo systemctl start docker</code>).
        Open Docker Desktop or start the service.
      </p>

      <h3>2. <code>port is already allocated</code></h3>
      <p>
        Something else on your machine already owns port 5432 (or 6379, or
        whatever). Probably a Postgres you installed via Homebrew years ago.
        Either stop it or change the port in <code>docker-compose.yml</code>:
      </p>
      <CodeBlock
        language="yaml"
        code={`# Was:  "127.0.0.1:5432:5432"
# Now:  127.0.0.1:5434:5432
#         host port ─┘   └─ container port (stays 5432)
# Then update DATABASE_URL in .env to use 5434.`}
      />

      <h3>3. <code>no space left on device</code></h3>
      <p>
        Docker disk usage piles up over months. Run{' '}
        <code>docker system df</code> to see usage; <code>docker system prune</code>{' '}
        clears safe stuff (stopped containers, dangling images). For aggressive
        cleanup including unused volumes: <code>docker system prune --volumes</code>
        {' '}— but verify you don&apos;t need that data first.
      </p>

      <h3>4. Container exits immediately (status: Exited (1))</h3>
      <p>
        Check the logs: <code>docker compose logs &lt;service&gt;</code>. The
        last 20 lines usually tell you exactly what went wrong (missing env
        var, password mismatch, port collision).
      </p>

      <h3>5. <code>permission denied</code> on Linux</h3>
      <p>
        Either run with <code>sudo</code> (annoying) or add yourself to the
        docker group: <code>sudo usermod -aG docker $USER</code> then log
        out + back in.
      </p>

      <h3>6. <code>healthcheck failed</code></h3>
      <p>
        Postgres / Redis are still starting up. Wait 5-10 seconds and the
        healthcheck should pass. If it stays red after a minute, check the
        logs — usually an env var typo.
      </p>

      <h3>7. <code>password authentication failed for user &quot;grit&quot; (SQLSTATE 28P01)</code></h3>
      <p>
        Grit 3.26+ projects use a single source of truth for Postgres
        credentials — the <code>POSTGRES_*</code> block in{' '}
        <code>.env</code>. <code>docker-compose.yml</code> reads it via{' '}
        <code>${'{POSTGRES_USER}'}</code> / <code>${'{POSTGRES_PASSWORD}'}</code> /
        <code>${'{POSTGRES_DB}'}</code> substitution; the Go API builds{' '}
        <code>DATABASE_URL</code> from the same parts. So this error
        shouldn&apos;t happen on a fresh scaffold.
      </p>
      <p>You can still hit it if you:</p>
      <ul>
        <li>
          Manually set <code>DATABASE_URL</code> in <code>.env</code> AND
          let it disagree with <code>POSTGRES_*</code>. Solution: set one
          OR the other, not both. <code>DATABASE_URL</code> wins when set —
          it&apos;s the &quot;external Postgres&quot; escape hatch
          (Neon, Supabase, RDS).
        </li>
        <li>
          Change <code>POSTGRES_PASSWORD</code> in <code>.env</code> after
          Postgres has already initialised its data volume.
        </li>
      </ul>
      <TipBox tone="warning">
        <strong>The Postgres volume gotcha:</strong> Postgres reads
        <code> POSTGRES_USER</code> / <code>POSTGRES_PASSWORD</code> /
        <code> POSTGRES_DB</code> from its environment <em>only on
        first-time volume init</em>. Edit <code>.env</code> after the
        container has already run once and the change is silently ignored
        — the password baked into the volume on first boot wins. To pick
        up new credentials you have to wipe the volume:
        <CodeBlock
          language="bash"
          code={`docker compose down -v   # the -v drops the postgres-data volume
docker compose up -d     # Postgres re-initialises with the new password from .env
grit migrate             # recreate the schema
grit seed                # re-seed dev data`}
        />
        Yes, this destroys local DB data. That&apos;s why you commit a
        seeder — five seconds with <code>grit seed</code> and you&apos;re
        back where you were.
      </TipBox>

      <h2>Running Grit WITHOUT Docker</h2>
      <p>
        You don&apos;t actually need Docker to use Grit. You need:
      </p>
      <ul>
        <li>A Postgres database (somewhere — local install, cloud, container)</li>
        <li>A Redis instance (same)</li>
        <li>S3-compatible storage (optional — only if you handle uploads)</li>
        <li>An SMTP server or a transactional email service (only if you send mail)</li>
      </ul>
      <p>
        Grit&apos;s scaffold ships with <code>.env.cloud.example</code> for
        exactly this case. Copy it over <code>.env</code>:
      </p>
      <CodeBlock
        language="bash"
        code={`cp .env.cloud.example .env
# Then fill in the connection strings`}
      />
      <p>
        The cloud picks I recommend for a free, no-Docker setup:
      </p>
      <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden my-5 text-sm">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border/30 bg-accent/20">
              <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Service</th>
              <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Provider</th>
              <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Free tier?</th>
              <th className="text-left px-4 py-2.5 font-medium text-foreground/80">env var</th>
            </tr>
          </thead>
          <tbody className="text-muted-foreground">
            <tr className="border-b border-border/20">
              <td className="px-4 py-2 font-medium text-foreground">Postgres</td>
              <td className="px-4 py-2">
                <a href="https://neon.tech" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Neon</a>
              </td>
              <td className="px-4 py-2">Yes — 0.5 GB</td>
              <td className="px-4 py-2 font-mono text-xs">DATABASE_URL</td>
            </tr>
            <tr className="border-b border-border/20">
              <td className="px-4 py-2 font-medium text-foreground">Redis</td>
              <td className="px-4 py-2">
                <a href="https://upstash.com" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Upstash</a>
              </td>
              <td className="px-4 py-2">Yes — 10k commands/day</td>
              <td className="px-4 py-2 font-mono text-xs">REDIS_URL</td>
            </tr>
            <tr className="border-b border-border/20">
              <td className="px-4 py-2 font-medium text-foreground">Object storage</td>
              <td className="px-4 py-2">
                <a href="https://dash.cloudflare.com" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Cloudflare R2</a>
              </td>
              <td className="px-4 py-2">Yes — 10 GB, no egress fees</td>
              <td className="px-4 py-2 font-mono text-xs">STORAGE_*</td>
            </tr>
            <tr>
              <td className="px-4 py-2 font-medium text-foreground">Email</td>
              <td className="px-4 py-2">
                <a href="https://resend.com" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Resend</a>
              </td>
              <td className="px-4 py-2">Yes — 3k/month, 1 domain</td>
              <td className="px-4 py-2 font-mono text-xs">RESEND_API_KEY</td>
            </tr>
          </tbody>
        </table>
      </div>
      <p>
        Sign up for each (5 minutes), paste the connection string into your
        <code> .env</code>, and your Grit project runs against managed
        services. No Docker, no localhost dependencies. Many production
        deploys actually look like this anyway — your Go binary on a small
        VPS, talking to managed Postgres + Redis. Docker is for{' '}
        <em>local development convenience</em>, not a hard requirement.
      </p>

      <h2>So… should I learn Docker?</h2>
      <p>
        Yes. Even if you start with the cloud-only setup, learn Docker.
        Three reasons:
      </p>
      <ol>
        <li>
          <strong>Onboarding speed.</strong> A new teammate clones the repo,
          runs <code>docker compose up -d</code>, and is productive in 60
          seconds. With managed services they have to sign up for 4 things and
          fill in <code>.env</code> values from a shared password manager. Big
          difference at scale.
        </li>
        <li>
          <strong>Reproducibility.</strong> Same Postgres version, same Redis
          version, same MinIO bucket policy, every dev. Bugs reproduce; fixes
          stay fixed.
        </li>
        <li>
          <strong>Production deployment.</strong> Most modern deploy targets
          (Dokploy, Coolify, Fly, Railway, Render, Kubernetes, ECS) speak
          Docker. The Dockerfile you ship today is your production artifact
          tomorrow. The patterns you learn here scale to every job you&apos;ll
          ever have.
        </li>
      </ol>
      <p>
        The first hour is the hardest. The second hour you&apos;re reading
        logs. After day two you wonder how you ever shipped software without
        it.
      </p>

      <KnowledgeCheck
        question="You ran `docker compose down` last night. This morning you run `docker compose up -d` and the database is empty. What happened?"
        choices={[
          {
            label: 'Docker corrupted the data',
            feedback: "Docker doesn't corrupt data. There's a specific reason this happens — keep going.",
          },
          {
            label: 'You ran `docker compose down -v` instead of plain `down` — the -v flag wipes named volumes (where Postgres stores data)',
            correct: true,
            feedback:
              "Right — `down` stops containers but preserves volumes (so data survives). `down -v` ALSO deletes the volumes. Easy to type by accident. The fix on a fresh `up` is to re-run your migrations + seed.",
          },
          {
            label: 'docker-compose.yml lost the volume definition',
            feedback: "The compose file you're using ships with named volumes. They persist by default.",
          },
          {
            label: 'Containers are stateless so data never persists',
            feedback: "Containers are stateless BUT volumes are persistent. That's exactly why volumes exist.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Get comfortable with the daily Docker loop:</p>
            <ol>
              <li>
                Make sure Docker Desktop (or your Linux daemon) is running.
                Run <code>docker run --rm hello-world</code> — confirm it
                works.
              </li>
              <li>
                In your <code>my-first-grit</code> folder, run{' '}
                <code>docker compose up -d</code>. Then <code>docker compose ps</code>{' '}
                — confirm 4 services are <code>healthy</code> or{' '}
                <code>running</code>.
              </li>
              <li>
                Open Mailhog at <code>http://localhost:8025</code> — confirm
                the inbox loads (empty is fine).
              </li>
              <li>
                Open MinIO at <code>http://localhost:9001</code> — login{' '}
                <code>minioadmin / minioadmin</code>, confirm you see the
                console.
              </li>
              <li>
                Run <code>docker compose logs -f postgres</code>. Wait 5
                seconds, then Ctrl+C. You should see startup logs ending in
                &quot;database system is ready to accept connections&quot;.
              </li>
              <li>
                Run <code>docker compose down</code> (no <code>-v</code>!) and
                confirm the containers stop. Run <code>docker compose ps</code>{' '}
                — list should be empty.
              </li>
            </ol>
            <p>
              Paste a screenshot of step 2&apos;s output into <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            If <code>docker compose up -d</code> fails with a port conflict,
            you have something already running on 5432 / 6379. Either stop
            that service or edit the host port in{' '}
            <code>docker-compose.yml</code> as shown in the &quot;errors&quot;
            section.
          </>
        }
        solution={
          <>
            <p>
              All 4 services should be visible in <code>docker compose ps</code>{' '}
              as <code>healthy</code>. Mailhog + MinIO consoles should load in
              your browser. Logs should be clean. <code>docker compose down</code>{' '}
              should leave you with an empty <code>ps</code>.
            </p>
            <p>
              Congratulations — you can now operate Docker comfortably for the
              entire rest of the course.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Docker boxes are humming. Next lesson — actually start the dev
        servers. <code>docker compose up -d</code> for infra,{' '}
        <code>pnpm dev</code> for the apps, and you&apos;ll have URLs to
        click in 3 minutes.
      </p>
    </>
  )
}
