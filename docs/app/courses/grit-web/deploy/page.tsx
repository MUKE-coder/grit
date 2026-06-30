import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Deploy to Production — Grit Web Course',
  description: 'Learn how to deploy your Grit app to production using grit deploy, systemd, Caddy reverse proxy with auto-TLS, Docker, and production best practices.',
}

export default function DeployToProductionCourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">Deploy to Production</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 8 of 8</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Deploy to Production
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this final course, you will learn how to take your Grit app from localhost to a live server
            that real users can access. You will learn the <Code>grit deploy</Code> command, how systemd keeps
            your app running, how Caddy provides automatic HTTPS, and how Docker gives you an alternative
            deployment path. By the end, you will have a complete deployment playbook.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is Deployment? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Deployment?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            During development, your app runs on <Code>localhost:3000</Code> or <Code>localhost:8080</Code>.
            Only you can access it from your own machine. Deployment is the process of putting your app on a
            real server with a real domain so that anyone on the internet can use it. It is the final step that
            turns your project into a product.
          </p>

          <Definition term="Deployment">
            The process of transferring your application from your local development machine to a remote server
            where it can be accessed by real users over the internet. It includes building your code, uploading
            it, configuring the server, and starting the application.
          </Definition>

          <Definition term="Production">
            The live environment where real users interact with your application. Unlike development (where errors
            are expected), production must be stable, secure, and performant. You never test in production — you
            test locally and in staging, then deploy to production when everything works.
          </Definition>

          <Definition term="VPS (Virtual Private Server)">
            A remote computer in the cloud that you rent from a provider like DigitalOcean, Hetzner, or Linode.
            You get root access to install software, run your app, and configure your server. It is like having
            your own computer in a data center, but it is a virtual slice of a larger physical machine.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On your local machine, you type <Code>grit dev</Code> and everything just works. But your users need
            to access your app from their browsers — they cannot connect to your laptop. A VPS gives your app a
            permanent home on the internet with a public IP address. You point your domain name to that IP, and
            suddenly <Code>myapp.com</Code> loads your Grit application.
          </p>

          <Challenge number={1} title="Localhost vs Production">
            <p>What{"'"}s the difference between running your app on localhost and deploying it to production?
            Think about who can access it, what URL they use, and what happens when you close your laptop.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: The grit deploy Command ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The grit deploy Command</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit gives you a single command that handles the entire deployment pipeline — building your code,
            uploading it to your server, configuring the process manager, and setting up HTTPS. Here is the
            full command:
          </p>

          <CodeBlock language="bash">{`grit deploy --host deploy@server.com --domain myapp.com`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            That{"'"}s it. One command, and your app is live. Let{"'"}s break down every flag:
          </p>

          <div className="bg-white/[0.03] border border-border/50 rounded-lg p-4 mb-4">
            <div className="space-y-3">
              <div>
                <Code>--host deploy@server.com</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  The SSH connection string. <Code>deploy</Code> is the username on your server, and
                  <Code> server.com</Code> is your server{"'"}s IP address or hostname. Grit connects to this
                  address to upload your app and configure services.
                </p>
              </div>
              <div>
                <Code>--domain myapp.com</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Your domain name. When provided, Grit configures Caddy to serve your app at this domain with
                  automatic HTTPS. If omitted, your app is accessible only via IP address on the app port.
                </p>
              </div>
              <div>
                <Code>--port 22</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  The SSH port on your server. Defaults to 22. Some servers use a custom port like 2222 for
                  extra security.
                </p>
              </div>
              <div>
                <Code>--key ~/.ssh/id_rsa</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Path to your SSH private key file. If your server uses key-based authentication (recommended),
                  point this to your private key.
                </p>
              </div>
              <div>
                <Code>--app-port 8080</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  The port your Go API listens on. Defaults to 8080. Caddy will forward traffic from port 443
                  (HTTPS) to this port internally.
                </p>
              </div>
            </div>
          </div>

          <Definition term="SSH (Secure Shell)">
            A protocol for securely connecting to a remote server. You use it to run commands on your server from
            your local computer. When you type <Code>ssh deploy@server.com</Code>, you get a terminal on the
            remote machine — as if you were sitting in front of it. All traffic is encrypted, so passwords and
            commands cannot be intercepted.
          </Definition>

          <Challenge number={2} title="Custom Deploy Flags">
            <p>Your server is at IP <Code>192.168.1.100</Code>, your domain is <Code>shop.example.com</Code>,
            and you use a custom SSH port <Code>2222</Code>. Write the full <Code>grit deploy</Code> command
            with the correct <Code>--host</Code>, <Code>--domain</Code>, and <Code>--port</Code> flags.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: What Happens During Deploy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What Happens During Deploy</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you run <Code>grit deploy</Code>, a 5-step pipeline executes automatically. Understanding
            each step helps you debug deployment issues and customize the process when needed.
          </p>

          <div className="space-y-4 mb-6">
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/20 text-primary flex items-center justify-center text-sm font-bold">1</div>
              <div>
                <p className="text-foreground font-medium">Cross-compile Go binary for Linux</p>
                <p className="text-muted-foreground text-sm mt-1">
                  Your development machine might run Windows or macOS, but your server runs Linux. Go makes
                  cross-compilation trivial — just set the target OS and architecture.
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/20 text-primary flex items-center justify-center text-sm font-bold">2</div>
              <div>
                <p className="text-foreground font-medium">Build frontend if present</p>
                <p className="text-muted-foreground text-sm mt-1">
                  If your project has a Next.js frontend, Grit runs <Code>pnpm build</Code> to create
                  production-optimized static files.
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/20 text-primary flex items-center justify-center text-sm font-bold">3</div>
              <div>
                <p className="text-foreground font-medium">Upload binary to server via SCP</p>
                <p className="text-muted-foreground text-sm mt-1">
                  The compiled binary is securely copied to <Code>/opt/myapp/</Code> on your server.
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/20 text-primary flex items-center justify-center text-sm font-bold">4</div>
              <div>
                <p className="text-foreground font-medium">Create systemd service with auto-restart</p>
                <p className="text-muted-foreground text-sm mt-1">
                  Grit writes a systemd unit file that keeps your app running 24/7 and automatically restarts
                  it if it crashes.
                </p>
              </div>
            </div>
            <div className="flex items-start gap-3">
              <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/20 text-primary flex items-center justify-center text-sm font-bold">5</div>
              <div>
                <p className="text-foreground font-medium">Configure Caddy reverse proxy with auto-TLS</p>
                <p className="text-muted-foreground text-sm mt-1">
                  If you provided <Code>--domain</Code>, Grit configures Caddy to handle HTTPS with automatic
                  certificate provisioning from Let{"'"}s Encrypt.
                </p>
              </div>
            </div>
          </div>

          <Definition term="Cross-compilation">
            Building a program on one OS (like Windows or Mac) that runs on a different OS (Linux). Go makes
            this easy — just set the <Code>GOOS</Code> and <Code>GOARCH</Code> environment variables. No
            extra tools or virtual machines needed. This is one of Go{"'"}s biggest strengths for deployment.
          </Definition>

          <Definition term="SCP (Secure Copy)">
            A command that copies files between your computer and a remote server over SSH. Like <Code>cp</Code>
            but over the network. When Grit uploads your binary to the server, it uses SCP under the hood to
            securely transfer the file.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is the exact build command Grit runs internally during Step 1:
          </p>

          <CodeBlock language="bash">{`CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/myapp ./cmd/server`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s break this down:
          </p>

          <div className="bg-white/[0.03] border border-border/50 rounded-lg p-4 mb-4">
            <div className="space-y-3">
              <div>
                <Code>CGO_ENABLED=0</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Disables C bindings (cgo). This produces a fully static binary with zero external dependencies.
                  Your binary will run on any Linux machine without needing to install libraries.
                </p>
              </div>
              <div>
                <Code>GOOS=linux</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Target operating system. Even if you are on Windows or macOS, Go will produce a Linux binary.
                </p>
              </div>
              <div>
                <Code>GOARCH=amd64</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Target CPU architecture. Most VPS servers use 64-bit x86 processors (amd64). If your server
                  uses ARM (like some AWS instances), you would use <Code>GOARCH=arm64</Code>.
                </p>
              </div>
              <div>
                <Code>-o bin/myapp</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Output path. The compiled binary will be saved as <Code>bin/myapp</Code> in your project
                  directory.
                </p>
              </div>
            </div>
          </div>

          <Challenge number={3} title="Build Command Breakdown">
            <p>Explain what each part of the build command does: <Code>CGO_ENABLED=0</Code>,
            <Code> GOOS=linux</Code>, <Code>GOARCH=amd64</Code>, <Code>-o bin/myapp</Code>. Why is
            <Code> CGO_ENABLED=0</Code> important for deployment?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: systemd — The Process Manager ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">systemd — The Process Manager</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Once your binary is on the server, something needs to keep it running. You cannot just SSH in and
            type <Code>./myapp</Code> — the process would die the moment you close your terminal. That{"'"}s where
            systemd comes in.
          </p>

          <Definition term="systemd">
            The init system and service manager for most modern Linux distributions (Ubuntu, Debian, CentOS,
            Fedora). It starts and manages long-running processes called {'"'}services.{'"'} When you tell systemd to
            manage your app, it ensures the app starts on boot, restarts on crash, and logs output to the
            system journal.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit automatically generates a systemd service file for your app. Here is what it looks like:
          </p>

          <CodeBlock language="ini">{`[Unit]
Description=myapp
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/myapp
Restart=on-failure
RestartSec=5
EnvironmentFile=/opt/myapp/.env

[Install]
WantedBy=multi-user.target`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s understand each section:
          </p>

          <div className="bg-white/[0.03] border border-border/50 rounded-lg p-4 mb-4">
            <div className="space-y-3">
              <div>
                <p className="text-foreground font-medium">[Unit] — Metadata</p>
                <p className="text-muted-foreground text-sm mt-1">
                  <Code>Description</Code> is a human-readable name. <Code>After=network.target</Code> means
                  {'"'}wait until the network is available before starting.{'"'} Your app needs a network to
                  listen for HTTP requests.
                </p>
              </div>
              <div>
                <p className="text-foreground font-medium">[Service] — How to run</p>
                <p className="text-muted-foreground text-sm mt-1">
                  <Code>Type=simple</Code> means the process runs in the foreground. <Code>User=www-data</Code>
                  runs the app as a non-root user for security. <Code>ExecStart</Code> is the command to launch
                  your binary. <Code>Restart=on-failure</Code> means systemd will restart the app if it exits
                  with an error. <Code>RestartSec=5</Code> waits 5 seconds between restart attempts.
                  <Code> EnvironmentFile</Code> loads your <Code>.env</Code> variables.
                </p>
              </div>
              <div>
                <p className="text-foreground font-medium">[Install] — When to start</p>
                <p className="text-muted-foreground text-sm mt-1">
                  <Code>WantedBy=multi-user.target</Code> means the service starts during normal system boot
                  (when the server has network and multi-user capabilities). This is the standard target for
                  server applications.
                </p>
              </div>
            </div>
          </div>

          <Definition term="Service">
            A long-running program managed by the operating system. Your Go API runs as a systemd service so
            it stays alive 24/7 — even after you disconnect from SSH, even after the server reboots. The
            operating system itself is responsible for keeping your app running.
          </Definition>

          <Tip>
            You can check your service status with <Code>sudo systemctl status myapp</Code>, view logs with
            <Code> sudo journalctl -u myapp -f</Code> (live tail), and manually restart with
            <Code> sudo systemctl restart myapp</Code>.
          </Tip>

          <Challenge number={4} title="systemd Crash Recovery">
            <p>Read the systemd file above. What happens if the app crashes? How does
            <Code> Restart=on-failure</Code> help? How long does systemd wait before restarting
            (<Code>RestartSec=5</Code>)?</p>
          </Challenge>

          <Challenge number={5} title="Security: User Permissions">
            <p>What user does the app run as? Why is <Code>www-data</Code> used instead of <Code>root</Code>?
            What could go wrong if your app ran as root and had a security vulnerability?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Caddy — Reverse Proxy with Auto-TLS ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Caddy — Reverse Proxy with Auto-TLS</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your Go app listens on port 8080, but users expect to visit <Code>https://myapp.com</Code> (port
            443). Something needs to sit between the user{"'"}s browser and your app to handle HTTPS, compress
            responses, and add security headers. That something is Caddy.
          </p>

          <Definition term="Reverse Proxy">
            A server that sits in front of your application and forwards requests to it. The browser talks to
            the proxy (Caddy on port 443), and the proxy talks to your app (port 8080). Users never connect to
            your app directly — they always go through Caddy. This adds a security layer and lets Caddy handle
            HTTPS, compression, and caching.
          </Definition>

          <Definition term="TLS (Transport Layer Security)">
            The encryption protocol that makes HTTPS work. When you see the padlock icon in your browser,
            that{"'"}s TLS in action. It encrypts all data between the browser and the server so that passwords,
            credit cards, and personal data cannot be intercepted by anyone on the network.
          </Definition>

          <Definition term="Let's Encrypt">
            A free, automated certificate authority that provides TLS certificates. Before Let{"'"}s Encrypt, you
            had to buy certificates and manually install them. Caddy integrates with Let{"'"}s Encrypt automatically
            — it obtains, installs, and renews certificates without you doing anything.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is the Caddy configuration that Grit generates:
          </p>

          <CodeBlock language="text">{`myapp.com {
    reverse_proxy localhost:8080
    encode gzip
    header {
        X-Frame-Options "DENY"
        X-Content-Type-Options "nosniff"
        Referrer-Policy "strict-origin-when-cross-origin"
        -Server
    }
    log {
        output file /var/log/caddy/myapp.log {
            roll_size 10mb
            roll_keep 5
        }
    }
}`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Let{"'"}s understand what each directive does:
          </p>

          <div className="bg-white/[0.03] border border-border/50 rounded-lg p-4 mb-4">
            <div className="space-y-3">
              <div>
                <Code>myapp.com</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  The domain name. Caddy automatically obtains a TLS certificate from Let{"'"}s Encrypt for this
                  domain. HTTPS is enabled by default — you do not need to configure it.
                </p>
              </div>
              <div>
                <Code>reverse_proxy localhost:8080</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Forward all incoming requests to your Go app running on port 8080.
                </p>
              </div>
              <div>
                <Code>encode gzip</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Compress responses with gzip. This makes your API responses and pages load faster by reducing
                  the data sent over the network.
                </p>
              </div>
              <div>
                <Code>X-Frame-Options {'"'}DENY{'"'}</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Prevents your site from being embedded in an iframe. Protects against clickjacking attacks.
                </p>
              </div>
              <div>
                <Code>X-Content-Type-Options {'"'}nosniff{'"'}</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Tells browsers to trust the Content-Type header and not try to guess the file type. Prevents
                  MIME-type sniffing attacks.
                </p>
              </div>
              <div>
                <Code>Referrer-Policy {'"'}strict-origin-when-cross-origin{'"'}</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Controls how much referrer information is sent when navigating away from your site. Reduces
                  information leakage.
                </p>
              </div>
              <div>
                <Code>-Server</Code>
                <p className="text-muted-foreground text-sm mt-1">
                  Removes the <Code>Server</Code> response header. By default, Caddy advertises itself in the
                  header. Removing it hides your technology stack from potential attackers — they do not need to
                  know you are using Caddy.
                </p>
              </div>
            </div>
          </div>

          <Note>
            HTTPS is completely automatic with Caddy. You do not need to run <Code>certbot</Code>, buy
            certificates, or configure renewal cron jobs. Caddy handles obtaining, installing, and renewing
            TLS certificates from Let{"'"}s Encrypt. It even redirects HTTP to HTTPS automatically.
          </Note>

          <Challenge number={6} title="Security Headers">
            <p>In the Caddy config, what security headers are set? Why is the <Code>Server</Code> header removed
            with <Code>-Server</Code>? What information would an attacker gain if the Server header was
            present?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Environment Variables for Deploy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Environment Variables for Deploy</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Typing <Code>--host</Code>, <Code>--domain</Code>, and <Code>--key</Code> every time you deploy is
            tedious and error-prone. Instead, you can set environment variables in your <Code>.env</Code> file
            and Grit will read them automatically.
          </p>

          <CodeBlock language="env">{`DEPLOY_HOST=deploy@server.com
DEPLOY_KEY_FILE=~/.ssh/id_rsa
DEPLOY_DOMAIN=myapp.com`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With these variables set, deploying becomes a single word:
          </p>

          <CodeBlock language="bash">{`grit deploy`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            No flags needed. Grit reads <Code>DEPLOY_HOST</Code>, <Code>DEPLOY_KEY_FILE</Code>, and
            <Code> DEPLOY_DOMAIN</Code> from your <Code>.env</Code> file. You can still override any variable
            with a flag — flags take priority over environment variables.
          </p>

          <Tip>
            Keep your <Code>.env</Code> file out of version control (it is in <Code>.gitignore</Code> by
            default). Each developer and server has its own <Code>.env</Code> with different values. Use
            <Code> .env.example</Code> as a template that IS committed to git, showing which variables are
            needed without revealing actual values.
          </Tip>

          <Challenge number={7} title="Configure Deploy Variables">
            <p>Add <Code>DEPLOY_HOST</Code> and <Code>DEPLOY_DOMAIN</Code> to your <Code>.env</Code> file
            for a server at <Code>deploy@198.51.100.42</Code> with domain <Code>mystore.com</Code>. Then run
            <Code> grit deploy</Code> without any flags.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Maintenance Mode During Deploy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Maintenance Mode During Deploy</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Sometimes you need to take your app offline briefly during a deployment — for example, when running
            database migrations that change table structures. Grit provides maintenance mode for this:
          </p>

          <CodeBlock language="bash">{`grit down              # 503 for all requests
grit deploy --host ... # Deploy new version
grit up                # Back online`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When maintenance mode is active, every API request receives a <Code>503 Service Unavailable</Code>
            response with a friendly message. Your frontend can detect this status code and show a
            {'"'}We{"'"}ll be right back{'"'} page instead of confusing error messages.
          </p>

          <Definition term="Maintenance Mode">
            A state where your app returns <Code>503 Service Unavailable</Code> for all requests. Useful during
            deployments so users see a {'"'}we{"'"}ll be right back{'"'} message instead of errors. It is a
            controlled way to take your app offline temporarily, unlike a crash where users see broken pages.
          </Definition>

          <Note>
            For most deployments, you do NOT need maintenance mode. The <Code>grit deploy</Code> command handles
            the binary swap and systemd restart so quickly that there is near-zero downtime. Use maintenance mode
            only when you are running migrations that could break things if the old code and new database schema
            are running simultaneously.
          </Note>

          <Challenge number={8} title="Maintenance Mode Test">
            <p>Run <Code>grit down</Code> in your local project. Try accessing the API — what HTTP status code
            do you get? What message does the response contain? Run <Code>grit up</Code> to bring it back
            online.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Docker Deployment (Alternative) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Docker Deployment (Alternative)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>grit deploy</Code> command is designed for simple VPS deployment — one server, one app. For
            more complex setups where you need to run your API, database, and cache all together in isolated
            containers, Docker with Docker Compose is the way to go.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit scaffolds a <Code>docker-compose.prod.yml</Code> file for production Docker deployment:
          </p>

          <CodeBlock language="yaml">{`services:
  api:
    build: ./apps/api
    ports:
      - "8080:8080"
    env_file: .env
    depends_on:
      - postgres
      - redis
  postgres:
    image: postgres:16
    volumes:
      - pgdata:/var/lib/postgresql/data
  redis:
    image: redis:7-alpine`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This defines your entire production stack: the API service (built from your Go code), PostgreSQL for
            the database, and Redis for caching and job queues. Docker Compose starts all three together and
            handles networking between them automatically.
          </p>

          <Definition term="Docker Compose (production)">
            A configuration file that defines how to run your entire application stack (API, database, cache) in
            Docker containers on a production server. Unlike the development <Code>docker-compose.yml</Code>
            which includes tools like Mailhog and MinIO, the production version only includes what{"'"}s needed to
            serve real users. The <Code>volumes</Code> directive ensures database data persists even if the
            container restarts.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To deploy with Docker, copy your project to the server and run:
          </p>

          <CodeBlock language="bash">{`docker compose -f docker-compose.prod.yml up -d --build`}</CodeBlock>

          <Tip>
            Use <Code>grit deploy</Code> when you have a simple setup — one VPS, one app, and you want the
            fastest path to production. Use Docker when you need reproducible environments, are deploying to
            multiple servers, or want to use container orchestration tools like Kubernetes in the future.
          </Tip>

          <Challenge number={9} title="Dev vs Production Docker">
            <p>Compare <Code>docker-compose.yml</Code> (development) with <Code>docker-compose.prod.yml</Code>
            (production). What services are present in development but missing in production? Why are tools like
            Mailhog and MinIO not needed in production?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Production Checklist ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Production Checklist</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Before going live, walk through this checklist. Missing any of these items could lead to security
            vulnerabilities, data loss, or embarrassing errors in front of real users.
          </p>

          <div className="bg-white/[0.03] border border-border/50 rounded-lg p-4 mb-4">
            <div className="space-y-3">
              <div className="flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-foreground font-medium">Change JWT_SECRET to a strong random string</p>
                  <p className="text-muted-foreground text-sm mt-1">
                    The default secret is for development only. Generate a random 64-character string with
                    <Code> openssl rand -hex 32</Code> and set it in your production <Code>.env</Code>.
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-foreground font-medium">Set APP_ENV=production</p>
                  <p className="text-muted-foreground text-sm mt-1">
                    This disables debug logging, enables stricter security checks, and optimizes performance
                    settings.
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-foreground font-medium">Configure real STORAGE_DRIVER (S3 or R2, not MinIO)</p>
                  <p className="text-muted-foreground text-sm mt-1">
                    MinIO is a development stand-in. In production, use Cloudflare R2, AWS S3, or Backblaze B2
                    for reliable, scalable file storage.
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-foreground font-medium">Set RESEND_API_KEY for real emails</p>
                  <p className="text-muted-foreground text-sm mt-1">
                    Mailhog catches emails in development but does not send them. Set your Resend API key so
                    password reset emails, welcome emails, and notifications actually reach users.
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-foreground font-medium">Remove default passwords</p>
                  <p className="text-muted-foreground text-sm mt-1">
                    Change all default passwords for GORM Studio, Sentinel rate limiter dashboard, and Pulse
                    monitoring. Attackers know the defaults.
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-foreground font-medium">Set up database backups</p>
                  <p className="text-muted-foreground text-sm mt-1">
                    Schedule daily PostgreSQL backups with <Code>pg_dump</Code> and store them off-server (S3
                    bucket, another server). Test restoring from a backup before you need it.
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <CheckCircle2 className="w-5 h-5 text-emerald-400 flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-foreground font-medium">Point your domain DNS to your server{"'"}s IP</p>
                  <p className="text-muted-foreground text-sm mt-1">
                    Create an A record in your DNS provider (Cloudflare, Namecheap, etc.) pointing your domain
                    to your server{"'"}s public IP address. DNS propagation can take up to 48 hours, so do this
                    early.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <Challenge number={10} title="Audit Your Checklist">
            <p>Go through the production checklist above. How many items apply to your project? Which ones would
            you most likely forget without this list? What could go wrong if you deployed with the default
            <Code> JWT_SECRET</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Updating a Deployed App ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Updating a Deployed App</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your app is live and users are using it. Now you need to ship a bug fix or a new feature. How do you
            update? Just run <Code>grit deploy</Code> again:
          </p>

          <CodeBlock language="bash">{`# Make your code changes locally
# Test them thoroughly
grit deploy`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit runs the same 5-step pipeline: rebuild the binary, upload it, and restart the systemd service.
            The restart happens so quickly that users experience near-zero downtime. systemd stops the old
            process and starts the new one in milliseconds.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here is the typical update workflow:
          </p>

          <div className="bg-white/[0.03] border border-border/50 rounded-lg p-4 mb-4">
            <div className="space-y-2">
              <div className="flex items-center gap-2">
                <span className="text-xs font-mono text-primary bg-primary/10 px-1.5 py-0.5 rounded">1</span>
                <p className="text-muted-foreground text-sm">Make code changes and commit to git</p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs font-mono text-primary bg-primary/10 px-1.5 py-0.5 rounded">2</span>
                <p className="text-muted-foreground text-sm">Run tests locally: <Code>go test ./...</Code></p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs font-mono text-primary bg-primary/10 px-1.5 py-0.5 rounded">3</span>
                <p className="text-muted-foreground text-sm">Run <Code>grit deploy</Code></p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs font-mono text-primary bg-primary/10 px-1.5 py-0.5 rounded">4</span>
                <p className="text-muted-foreground text-sm">Verify the update: visit your domain and test the changes</p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs font-mono text-primary bg-primary/10 px-1.5 py-0.5 rounded">5</span>
                <p className="text-muted-foreground text-sm">Check logs: <Code>sudo journalctl -u myapp -f</Code></p>
              </div>
            </div>
          </div>

          <Tip>
            If an update breaks something, you can SSH into your server, replace the binary with the previous
            version (keep backups!), and restart the service with <Code>sudo systemctl restart myapp</Code>.
            This is your manual rollback procedure until you set up CI/CD with automated rollbacks.
          </Tip>

          <Challenge number={11} title="Update Workflow">
            <p>You have fixed a bug in your user registration handler. Describe the exact steps you would take to
            deploy this fix to production, starting from your local machine. Include the commands you would
            run.</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You have completed the entire Grit Web course series. Let{"'"}s recap what you learned in this final
            course:
          </p>

          <div className="bg-white/[0.03] border border-border/50 rounded-lg p-4 mb-6">
            <div className="space-y-2">
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
                <p className="text-muted-foreground text-sm">
                  <strong className="text-foreground">grit deploy</strong> — one command to build, upload, and
                  configure your production server
                </p>
              </div>
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
                <p className="text-muted-foreground text-sm">
                  <strong className="text-foreground">Cross-compilation</strong> — Go builds Linux binaries from
                  any OS with GOOS and GOARCH
                </p>
              </div>
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
                <p className="text-muted-foreground text-sm">
                  <strong className="text-foreground">systemd</strong> — keeps your app running 24/7 with
                  auto-restart on crash and boot
                </p>
              </div>
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
                <p className="text-muted-foreground text-sm">
                  <strong className="text-foreground">Caddy</strong> — reverse proxy with automatic HTTPS from
                  Let{"'"}s Encrypt, security headers, and gzip compression
                </p>
              </div>
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
                <p className="text-muted-foreground text-sm">
                  <strong className="text-foreground">Environment variables</strong> — configure deployment
                  without flags using DEPLOY_HOST, DEPLOY_DOMAIN, DEPLOY_KEY_FILE
                </p>
              </div>
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
                <p className="text-muted-foreground text-sm">
                  <strong className="text-foreground">Maintenance mode</strong> — grit down/up for controlled
                  downtime during major migrations
                </p>
              </div>
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
                <p className="text-muted-foreground text-sm">
                  <strong className="text-foreground">Docker deployment</strong> — an alternative for complex
                  setups with docker-compose.prod.yml
                </p>
              </div>
              <div className="flex items-start gap-2">
                <CheckCircle2 className="w-4 h-4 text-emerald-400 flex-shrink-0 mt-0.5" />
                <p className="text-muted-foreground text-sm">
                  <strong className="text-foreground">Production checklist</strong> — JWT secrets, environment
                  settings, real email/storage drivers, DNS configuration
                </p>
              </div>
            </div>
          </div>

          <Challenge number={12} title="Final Challenge: Complete Deployment Plan">
            <p>You are deploying a bookstore application built with Grit. Write a complete deployment plan that
            covers:</p>
            <ol className="list-decimal list-inside mt-2 space-y-1 text-muted-foreground">
              <li><strong className="text-foreground">VPS provider choice</strong> — which provider would you use
              (DigitalOcean, Hetzner, Linode, etc.) and why?</li>
              <li><strong className="text-foreground">The grit deploy command</strong> — write the exact command
              with all flags for your bookstore domain</li>
              <li><strong className="text-foreground">Production .env changes</strong> — list every variable
              you would change from development defaults</li>
              <li><strong className="text-foreground">Production checklist</strong> — walk through each item
              and explain what you would do</li>
              <li><strong className="text-foreground">First deployment vs updates</strong> — how does the first
              deploy differ from subsequent updates?</li>
            </ol>
          </Challenge>

          <Challenge number={13} title="Bonus: Rollback Strategy">
            <p>Your latest deployment introduced a bug that breaks checkout. Describe your rollback strategy: how
            would you revert to the previous working version? What tools and commands would you use? How would
            you prevent this from happening again?</p>
          </Challenge>

          <Challenge number={14} title="Bonus: Monitoring After Deploy">
            <p>Your app is live. How would you monitor it? Describe what logs you would check
            (<Code>journalctl</Code>), what metrics matter (response times, error rates), and how you would know
            if something breaks at 3 AM (alerts, health checks). Write a monitoring plan for your first week in
            production.</p>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web/ai-features', label: 'AI Features' }}
            next={{ href: '/courses/grit-web', label: 'Back to All Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
