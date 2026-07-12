import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/prerequisites/docker')

/* ------------------------------------------------------------------ */
/*  Code strings defined outside JSX to avoid curly-brace parsing     */
/* ------------------------------------------------------------------ */

const dockerRunExample = `# Pull the official PostgreSQL image
$ docker pull postgres:16-alpine

# Run a container from that image
$ docker run -d \\
  --name my-postgres \\
  -e POSTGRES_USER=grit \\
  -e POSTGRES_PASSWORD=grit \\
  -e POSTGRES_DB=myapp \\
  -p 5432:5432 \\
  postgres:16-alpine`;

const dockerfileExample = `# ---------- Build stage ----------
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Download dependencies first (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# ---------- Run stage ----------
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]`;

const simpleComposeExample = `services:
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: grit
      POSTGRES_PASSWORD: grit
      POSTGRES_DB: myapp
    volumes:
      - postgres-data:/var/lib/postgresql/data

volumes:
  postgres-data:`;

const gritComposeExample = `services:
  # --- Primary Database ---
  postgres:
    image: postgres:16-alpine
    container_name: myapp-postgres
    restart: unless-stopped
    ports:
      - "5434:5432"
    environment:
      POSTGRES_USER: grit
      POSTGRES_PASSWORD: grit
      POSTGRES_DB: myapp
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U grit"]
      interval: 5s
      timeout: 5s
      retries: 5

  # --- Cache & Job Queue ---
  redis:
    image: redis:7-alpine
    container_name: myapp-redis
    restart: unless-stopped
    ports:
      - "6380:6379"
    volumes:
      - redis-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  # --- S3-Compatible File Storage ---
  minio:
    image: minio/minio
    container_name: myapp-minio
    restart: unless-stopped
    ports:
      - "9002:9000"
      - "9003:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio-data:/data
    command: server /data --console-address ":9001"

  # --- Email Testing ---
  mailhog:
    image: mailhog/mailhog
    container_name: myapp-mailhog
    restart: unless-stopped
    ports:
      - "1025:1025"
      - "8025:8025"

volumes:
  postgres-data:
  redis-data:
  minio-data:`;

const volumeExample = `volumes:
  postgres-data:    # Database files
  redis-data:       # Cache snapshots
  minio-data:       # Uploaded files`;

const portMappingExample = `ports:
  - "5434:5432"   # HOST:CONTAINER
  #    ^     ^
  #    |     |
  #    |     +-- Port inside the container
  #    +-------- Port on your machine (localhost)

# Your Go API connects to localhost:5434
# which Docker routes to the container's port 5432`;

export default function DockerForGritPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Prerequisites</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Docker for Grit Developers</h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                A practical guide to Docker for developers who want to use Grit. You don&apos;t need to be a Docker
                expert &mdash; just understand enough to start and stop the infrastructure services that power
                your Grit project.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 1. What is Docker? */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>What is Docker?</h2>
              <p>
                Think of Docker like a shipping container for software. Just as shipping containers let you
                move goods between ships, trucks, and warehouses without repacking, Docker containers let you
                run software on any computer without worrying about what&apos;s already installed there.
              </p>
              <p>
                Without Docker, setting up a database like PostgreSQL means downloading an installer, configuring
                paths, managing versions, and troubleshooting OS-specific quirks. With Docker, you run a single
                command and get a fully working database in seconds &mdash; identical on every developer&apos;s machine.
              </p>
              <p>
                This eliminates the &quot;works on my machine&quot; problem entirely. Every developer on your team
                gets the exact same PostgreSQL version, the exact same Redis configuration, and the exact same
                MinIO file storage &mdash; guaranteed.
              </p>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Docker runs your database (PostgreSQL), cache and job queue (Redis), file storage (MinIO),
                and email testing server (Mailhog). Your Go API and Next.js apps run natively on your machine
                for faster development &mdash; only the infrastructure services live in Docker containers.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 2. Installing Docker */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Installing Docker</h2>
              <p>
                Docker Desktop is the easiest way to get started on Windows and macOS. It bundles the Docker engine,
                Docker Compose, and a graphical interface for managing containers. On Linux, you install Docker Engine
                and Docker Compose separately.
              </p>
            </div>

            <div className="space-y-4 mb-6">
              {[
                {
                  platform: 'macOS',
                  instruction: 'Download Docker Desktop from docker.com/products/docker-desktop and run the installer. Or use Homebrew:',
                },
                {
                  platform: 'Windows',
                  instruction: 'Download Docker Desktop from docker.com/products/docker-desktop. Enable WSL 2 when prompted (recommended). Restart your computer after installation.',
                },
                {
                  platform: 'Linux (Ubuntu/Debian)',
                  instruction: 'Install Docker Engine and the Compose plugin from the official repository.',
                },
              ].map((item) => (
                <div key={item.platform} className="p-4 rounded-lg border border-border/30 bg-card/30">
                  <h3 className="text-sm font-semibold mb-1.5">{item.platform}</h3>
                  <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.instruction}</p>
                </div>
              ))}
            </div>

            <div className="mb-4">
              <p className="text-sm font-semibold mb-2 text-foreground/80">
                macOS (Homebrew)
              </p>
              <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden">
                <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                  <span className="text-[11px] font-mono text-muted-foreground/40">terminal</span>
                </div>
                <div className="p-5 font-mono text-sm">
                  <span className="text-primary/50 select-none">$ </span>
                  <span className="text-foreground/80">brew install --cask docker</span>
                </div>
              </div>
            </div>

            <div className="mb-6">
              <p className="text-sm font-semibold mb-2 text-foreground/80">
                Linux (Ubuntu/Debian)
              </p>
              <CodeBlock language="bash" code={`# Install Docker Engine
$ curl -fsSL https://get.docker.com | sh

# Add your user to the docker group (avoids needing sudo)
$ sudo usermod -aG docker $USER

# Log out and back in, then verify
$ docker --version`} />
            </div>

            <div className="mb-6">
              <p className="text-sm font-semibold mb-2 text-foreground/80">
                Verify your installation
              </p>
              <CodeBlock language="bash" code={`$ docker --version
Docker version 27.x.x, build abcdef0

$ docker compose version
Docker Compose version v2.x.x`} />
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Docker must be installed and running before you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code>.
                The scaffolded project includes a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker-compose.yml</code> that
                requires Docker Compose V2 (the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker compose</code> command
                without the hyphen).
              </p>
            </div>

            {/* ============================================================ */}
            {/* 3. Images & Containers */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Images &amp; Containers</h2>
              <p>
                An <strong>image</strong> is a blueprint &mdash; a read-only template that contains everything needed
                to run a piece of software: the operating system, the application, its dependencies, and configuration.
                Think of it like a recipe.
              </p>
              <p>
                A <strong>container</strong> is a running instance of an image &mdash; the actual dish you cook from
                the recipe. You can run multiple containers from the same image, and each one is isolated from the
                others. When you stop a container, it&apos;s like throwing away the dish, but the recipe (image)
                remains for next time.
              </p>
            </div>

            <CodeBlock language="bash" code={dockerRunExample} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit uses official images from Docker Hub: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">postgres:16-alpine</code> for
                PostgreSQL, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">redis:7-alpine</code> for
                Redis, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">minio/minio</code> for file storage,
                and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">mailhog/mailhog</code> for email testing.
                The &quot;alpine&quot; variants are smaller, lightweight images based on Alpine Linux.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 4. Dockerfile Basics */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Dockerfile Basics</h2>
              <p>
                A <strong>Dockerfile</strong> is a text file with instructions for building a custom image. Each line
                is a step: choose a base image, copy files, install dependencies, and specify how to start the
                application. Here are the key instructions:
              </p>
              <ul>
                <li><strong>FROM</strong> &mdash; Sets the base image (e.g., <code>golang:1.21-alpine</code>)</li>
                <li><strong>WORKDIR</strong> &mdash; Sets the working directory inside the container</li>
                <li><strong>COPY</strong> &mdash; Copies files from your machine into the image</li>
                <li><strong>RUN</strong> &mdash; Executes a command during the build (e.g., installing dependencies)</li>
                <li><strong>EXPOSE</strong> &mdash; Documents which port the app listens on</li>
                <li><strong>CMD</strong> &mdash; The command that runs when a container starts</li>
              </ul>
              <p>
                A <strong>multi-stage build</strong> uses multiple <code>FROM</code> statements. The first stage
                compiles your code with all build tools, then the final stage copies only the compiled binary into
                a tiny base image. This keeps production images small and secure.
              </p>
            </div>

            <CodeBlock filename="apps/api/Dockerfile" code={dockerfileExample} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit scaffolds two Dockerfiles: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/api/Dockerfile</code> builds
                the Go binary using a multi-stage build (final image is under 20MB),
                and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/web/Dockerfile</code> builds
                the Next.js app with standalone output. These are used by the production compose file, not during
                everyday development.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 5. Docker Compose */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Docker Compose</h2>
              <p>
                Running individual <code>docker run</code> commands for each service gets tedious fast. <strong>Docker
                Compose</strong> solves this by letting you define all your services in a
                single <code>docker-compose.yml</code> file. One command starts everything.
              </p>
              <p>
                A compose file defines <strong>services</strong> (containers to run), their <strong>images</strong>,
                port mappings, environment variables, volumes for data persistence,
                and <strong>dependencies</strong> between services. It&apos;s like a blueprint for your entire
                infrastructure.
              </p>
            </div>

            <CodeBlock language="yaml" filename="docker-compose.yml (simplified example)" code={simpleComposeExample} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Every Grit project has a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker-compose.yml</code> at
                the project root that defines all four infrastructure services. There&apos;s also
                a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker-compose.prod.yml</code> for
                production that additionally builds and runs the Go API, web app, and admin panel as containers.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 6. Grit's Docker Services */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Grit&apos;s Docker Services</h2>
              <p>
                When you scaffold a new Grit project, the generated <code>docker-compose.yml</code> includes four
                services. Here&apos;s what each one does and why it&apos;s there:
              </p>
            </div>

            <div className="space-y-4 mb-6">
              {[
                {
                  title: 'PostgreSQL (port 5434)',
                  desc: 'Your primary relational database. Stores users, resources, and all application data. Grit uses GORM to auto-migrate your Go models into PostgreSQL tables. Default credentials: grit / grit.',
                },
                {
                  title: 'Redis (port 6380)',
                  desc: 'An in-memory data store used for two things: caching API responses (via the cache middleware) and processing background jobs (via the asynq job queue). No authentication required in development.',
                },
                {
                  title: 'MinIO (ports 9002 / 9003)',
                  desc: 'An S3-compatible object storage server. Handles file uploads like user avatars, documents, and images. Port 9002 is the API endpoint your Go code talks to. Port 9003 is a web console where you can browse files, create buckets, and manage storage. Login: minioadmin / minioadmin.',
                },
                {
                  title: 'Mailhog (ports 1025 / 8025)',
                  desc: 'A fake SMTP server that catches all outgoing emails. Port 1025 receives emails from your Go API. Port 8025 is a web interface where you can view every email your app sent &mdash; perfect for testing password resets, welcome emails, and notifications without sending real mail.',
                },
              ].map((service) => (
                <div key={service.title} className="p-4 rounded-lg border border-border/30 bg-card/30">
                  <h3 className="text-sm font-semibold mb-1.5">{service.title}</h3>
                  <p className="text-xs text-muted-foreground/70 leading-relaxed">{service.desc}</p>
                </div>
              ))}
            </div>

            <CodeBlock language="yaml" filename="docker-compose.yml (generated by grit new)" code={gritComposeExample} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Start all four services with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker compose up -d</code>.
                They&apos;ll run in the background. Your Go API connects to them via <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">localhost</code> and
                the ports listed above. All connection strings are pre-configured in
                the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 7. Essential Commands */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Essential Commands</h2>
              <p>
                You only need a handful of Docker commands for day-to-day Grit development.
                Here&apos;s your cheat sheet:
              </p>
            </div>

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Command</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">What It Does</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  {[
                    { cmd: 'docker compose up -d', desc: 'Start all services in the background (detached mode)' },
                    { cmd: 'docker compose down', desc: 'Stop and remove all containers (data is preserved in volumes)' },
                    { cmd: 'docker compose logs', desc: 'Show logs from all services' },
                    { cmd: 'docker compose logs -f postgres', desc: 'Follow (stream) logs for a specific service' },
                    { cmd: 'docker ps', desc: 'List all running containers' },
                    { cmd: 'docker compose restart redis', desc: 'Restart a specific service' },
                    { cmd: 'docker compose down -v', desc: 'Stop everything AND delete all data volumes (fresh start)' },
                    { cmd: 'docker compose ps', desc: 'Show the status of all compose services' },
                  ].map((item) => (
                    <tr key={item.cmd} className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs text-foreground/90">{item.cmd}</td>
                      <td className="px-4 py-2.5 text-xs">{item.desc}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <CodeBlock language="bash" code={`# Start your infrastructure
$ docker compose up -d

# Check that everything is running
$ docker compose ps

# When you're done for the day
$ docker compose down`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                In practice, you mainly need two commands: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker compose up -d</code> to
                start your infrastructure at the beginning of a coding session,
                and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker compose down</code> when
                you&apos;re done. Everything else is for troubleshooting.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 8. Volumes & Data */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Volumes &amp; Data Persistence</h2>
              <p>
                By default, when you stop and remove a container, all its data disappears. <strong>Volumes</strong> solve
                this by storing data outside the container on your machine&apos;s filesystem. When you restart
                the container, the volume is re-attached and all your data is still there.
              </p>
              <p>
                Docker supports two types of storage: <strong>named volumes</strong> (managed by Docker, recommended)
                and <strong>bind mounts</strong> (map a specific folder on your machine). Grit uses named volumes
                for all services because they&apos;re simpler and work the same on every operating system.
              </p>
            </div>

            <CodeBlock language="yaml" filename="docker-compose.yml (volumes section)" code={volumeExample} />

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Volume</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Stores</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Survives Restart?</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-mono text-xs">postgres-data</td>
                    <td className="px-4 py-2.5">All your database tables, rows, and indexes</td>
                    <td className="px-4 py-2.5">Yes</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-mono text-xs">redis-data</td>
                    <td className="px-4 py-2.5">Cached data and job queue state</td>
                    <td className="px-4 py-2.5">Yes</td>
                  </tr>
                  <tr>
                    <td className="px-4 py-2.5 font-mono text-xs">minio-data</td>
                    <td className="px-4 py-2.5">Uploaded files and bucket metadata</td>
                    <td className="px-4 py-2.5">Yes</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Your PostgreSQL data persists across container restarts thanks to the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">postgres-data</code> volume.
                If you need a completely fresh start (e.g., your schema is out of sync),
                run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker compose down -v</code> to
                delete all volumes, then <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker compose up -d</code> to
                recreate everything from scratch. GORM will auto-migrate your tables on the next API startup.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 9. Port Mapping */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Port Mapping</h2>
              <p>
                Containers run in their own isolated network. To access a service from your machine, you create
                a <strong>port mapping</strong> that connects a port on your computer (the host) to a port inside the
                container. The syntax is <code>HOST:CONTAINER</code>.
              </p>
              <p>
                For example, <code>&quot;5434:5432&quot;</code> means &quot;when something connects to port 5434 on my
                machine, route it to port 5432 inside the container.&quot; Your Go API code
                uses <code>localhost:5434</code> to talk to PostgreSQL, and Docker transparently routes the traffic
                to the container.
              </p>
            </div>

            <CodeBlock filename="port mapping explained" code={portMappingExample} />

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Service</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Mapping</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Access From Your Code</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">PostgreSQL</td>
                    <td className="px-4 py-2.5 font-mono text-xs">5434:5432</td>
                    <td className="px-4 py-2.5 font-mono text-xs">localhost:5434</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">Redis</td>
                    <td className="px-4 py-2.5 font-mono text-xs">6380:6379</td>
                    <td className="px-4 py-2.5 font-mono text-xs">localhost:6380</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">MinIO API</td>
                    <td className="px-4 py-2.5 font-mono text-xs">9002:9000</td>
                    <td className="px-4 py-2.5 font-mono text-xs">localhost:9002</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">MinIO Console</td>
                    <td className="px-4 py-2.5 font-mono text-xs">9003:9001</td>
                    <td className="px-4 py-2.5 font-mono text-xs">localhost:9003 (browser)</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">Mailhog SMTP</td>
                    <td className="px-4 py-2.5 font-mono text-xs">1025:1025</td>
                    <td className="px-4 py-2.5 font-mono text-xs">localhost:1025</td>
                  </tr>
                  <tr>
                    <td className="px-4 py-2.5 font-medium text-foreground/90">Mailhog UI</td>
                    <td className="px-4 py-2.5 font-mono text-xs">8025:8025</td>
                    <td className="px-4 py-2.5 font-mono text-xs">localhost:8025 (browser)</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div className="prose-grit mb-4">
              <p>
                <strong>Port busy?</strong> If a port is already in use, change the host side of the mapping.
                For example, change <code>&quot;5434:5432&quot;</code> to <code>&quot;5435:5432&quot;</code> and
                update the <code>POSTGRES_PORT</code> in your <code>.env</code> file to use port 5435.
              </p>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file
                generated by <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code> contains
                all connection strings pre-configured to point at these default ports. If you change a port
                mapping in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">docker-compose.yml</code>,
                update the corresponding environment variable in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> to match.
              </p>
            </div>

            {/* ============================================================ */}
            {/* 10. Troubleshooting */}
            {/* ============================================================ */}
            <div className="prose-grit mb-10">
              <h2>Troubleshooting</h2>
              <p>
                Docker is generally reliable, but things can go wrong. Here are the most common issues
                and how to fix them:
              </p>
            </div>

            <div className="space-y-4 mb-6">
              {[
                {
                  title: 'Container won\u2019t start',
                  solution: 'Check the logs to see what went wrong. Run docker compose logs <service> (e.g., docker compose logs postgres). Common causes include corrupted data or missing environment variables. If the logs mention data corruption, try docker compose down -v && docker compose up -d for a clean start.',
                },
                {
                  title: 'Port already in use',
                  solution: 'Another process is using the port. On macOS/Linux, find it with lsof -i :5432 and kill it, or change the host port in docker-compose.yml (e.g., "5433:5432") and update your .env file. On Windows, use netstat -ano | findstr :5432 to find the process ID.',
                },
                {
                  title: 'Out of disk space',
                  solution: 'Docker images and volumes accumulate over time. Run docker system prune to remove unused containers, networks, and dangling images. Add the -a flag to also remove unused images. Add --volumes to remove unused volumes (caution: this deletes data).',
                },
                {
                  title: 'WSL issues on Windows',
                  solution: 'Docker Desktop on Windows requires WSL 2. If you see WSL-related errors, open PowerShell as admin and run wsl --update. Make sure WSL 2 is the default: wsl --set-default-version 2. Restart Docker Desktop after any WSL changes.',
                },
                {
                  title: 'Permission errors on Linux',
                  solution: 'If you see "permission denied" when running docker commands, add your user to the docker group: sudo usermod -aG docker $USER. Then log out and back in (or run newgrp docker) for the change to take effect.',
                },
                {
                  title: 'Docker daemon not running',
                  solution: 'On macOS/Windows, make sure Docker Desktop is open and running (check for the whale icon in your system tray). On Linux, start the daemon with sudo systemctl start docker. To make it start automatically: sudo systemctl enable docker.',
                },
              ].map((item) => (
                <div key={item.title} className="p-4 rounded-lg border border-border/30 bg-card/30">
                  <h3 className="text-sm font-semibold mb-1.5">{item.title}</h3>
                  <p className="text-xs text-muted-foreground/70 leading-relaxed">{item.solution}</p>
                </div>
              ))}
            </div>

            <CodeBlock filename="terminal &mdash; common troubleshooting commands" code={`# View logs for a specific service
$ docker compose logs postgres

# Follow logs in real-time
$ docker compose logs -f redis

# Check what's using a port (macOS/Linux)
$ lsof -i :5432

# Clean up unused Docker resources
$ docker system prune

# Nuclear option: remove everything and start fresh
$ docker compose down -v && docker compose up -d`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                For more detailed troubleshooting tips specific to Grit projects, see
                the{" "}
                <Link href="/docs/getting-started/troubleshooting" className="text-primary hover:underline">
                  Troubleshooting
                </Link>{" "}
                page in the Getting Started section. It covers API connection issues, frontend build
                errors, and other common problems beyond Docker.
              </p>
            </div>

            {/* Navigation footer */}
            <div className="mt-16 pt-8 border-t border-border/40 flex items-center justify-between">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/prerequisites/nextjs" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Next.js &amp; React
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/getting-started/quick-start" className="gap-1.5">
                  Quick Start
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
