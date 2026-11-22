import Link from 'next/link'
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
  title: 'Deploy a Go and React app to a VPS with HTTPS',
  description:
    'What a VPS, systemd and a reverse proxy are, and a first deploy end to end: a server, a domain, a certificate and a service that restarts itself.',
  alternates: { canonical: 'https://gritframework.dev/guides/deploy-go-app-vps-https' },
  openGraph: {
    title: 'Deploy a Go and React app to a VPS with HTTPS',
    description: 'From an empty server to a running app on your own domain, with TLS.',
    url: 'https://gritframework.dev/guides/deploy-go-app-vps-https',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="deploy-go-app-vps-https"
      lead={
        <>
          <p className="mb-4">
            A Go program compiles to one file with no runtime to install, which makes deploying it
            genuinely simple: copy the file to a computer on the internet and run it. Everything
            else, the proxy, the certificate, the service manager, exists to handle what happens
            when that computer reboots or somebody visits over HTTPS.
          </p>
          <p>
            This goes from a fresh server to your app on your own domain with a valid certificate.
            It assumes nothing: if you do not know what systemd is, that is covered.
          </p>
        </>
      }
    >
      <Outcome>
        <p>
          Your API and frontend reachable at <code>https://yourdomain.com</code>, with a
          certificate that renews itself.
        </p>
        <p>A service that starts on boot and restarts if it crashes.</p>
        <p>A server that is not trivially breakable: no password logins, a firewall, no root.</p>
      </Outcome>

      <Section id="vocab" title="The four words you need">
        <Checklist
          items={[
            <>
              A <strong>VPS</strong> is a computer somebody else runs, that you get a root login
              to. Hetzner, DigitalOcean, Vultr and Linode all rent small ones for a few currency
              units a month, which is enough for a real application with real users.
            </>,
            <>
              <strong>systemd</strong> is the thing that starts programs when Linux boots and
              restarts them when they die. Without it, your app stops the moment you close your
              SSH session.
            </>,
            <>
              A <strong>reverse proxy</strong> sits in front of your app on ports 80 and 443,
              terminates HTTPS and forwards to your program on some internal port. Caddy, nginx and
              Traefik all do this; Caddy gets certificates by itself, which is why it is here.
            </>,
            <>
              A <strong>certificate</strong> is what makes the padlock appear. Let&apos;s Encrypt
              issues them free, for ninety days, automatically renewed.
            </>,
          ]}
        />
      </Section>

      <Step n={1} title="Get a server and a domain">
        <p>
          Create the smallest instance your provider offers with Ubuntu LTS. Two vCPU and 2GB is
          comfortable for a Go API, Postgres and a Next.js frontend on one box. One vCPU and 1GB
          works and will make <code>next build</code> painful, which is an argument for building
          elsewhere, later on this page.
        </p>
        <p>Point your domain at it. Two records, at your registrar:</p>
        <CodeBlock
          language="text"
          code={`Type   Name   Value              TTL
A      @      203.0.113.42       300
A      www    203.0.113.42       300`}
        />
        <p>
          <strong>Checkpoint.</strong> Wait a few minutes, then confirm the world agrees:
        </p>
        <CodeBlock
          terminal
          code={`dig +short yourdomain.com
# 203.0.113.42`}
        />
        <p>
          Do not move on until that prints your server&apos;s address. Certificate issuance checks
          DNS, so a deploy attempted before DNS has propagated fails in a way that looks like a
          Caddy problem and is not.
        </p>
      </Step>

      <Step n={2} title="Secure the server before anything is on it">
        <p>
          A new server on a public address starts receiving password-guessing attempts within
          minutes. Ten minutes now, once.
        </p>
        <CodeBlock
          terminal
          code={`# From your own machine: copy your key up, so you never type a password
ssh-copy-id root@203.0.113.42
ssh root@203.0.113.42`}
        />
        <CodeBlock
          terminal
          code={`# On the server: a user that is not root
adduser --disabled-password --gecos "" deploy
usermod -aG sudo deploy
rsync --archive --chown=deploy:deploy ~/.ssh /home/deploy/`}
        />
        <CodeBlock
          language="text"
          filename="/etc/ssh/sshd_config"
          code={`PasswordAuthentication no
PermitRootLogin no
PubkeyAuthentication yes`}
        />
        <CodeBlock
          terminal
          code={`systemctl restart ssh

ufw allow OpenSSH
ufw allow 80
ufw allow 443
ufw --force enable

apt update && apt install -y unattended-upgrades`}
        />
        <Warn title="Open a second terminal before you log out">
          <p>
            Confirm <code>ssh deploy@203.0.113.42</code> works in a new window{' '}
            <strong>while the first one is still connected</strong>. If the key did not copy
            correctly, closing your only session locks you out of your own server, and the recovery
            is your provider&apos;s web console at best.
          </p>
        </Warn>
        <p>
          Note what is not open: nothing but SSH, 80 and 443. Postgres stays unreachable from the
          internet, and you reach it over an SSH tunnel when you need to. A database on a public
          port is found by scanners the same day.
        </p>
      </Step>

      <Step n={3} title="Build the binary, on your machine">
        <p>
          Cross-compiling is one line, and it means the server needs no Go toolchain, no build
          dependencies and no memory to spare.
        </p>
        <CodeBlock
          terminal
          code={`cd apps/api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \\
  go build -ldflags="-s -w" -o api-linux ./cmd/server

# -s -w drop the debug symbols: a smaller file to upload, and you are not
# debugging on the server anyway.
# CGO_ENABLED=0 produces a fully static binary with no libc dependency, so
# it runs on any Linux. Note: it also rules out the cgo SQLite driver, which
# is one more reason production uses Postgres.`}
        />
        <CodeBlock
          terminal
          code={`scp api-linux deploy@203.0.113.42:/home/deploy/
ssh deploy@203.0.113.42 'sudo mv ~/api-linux /usr/local/bin/myapp-api && sudo chmod +x /usr/local/bin/myapp-api'`}
        />
        <p>
          <strong>Checkpoint.</strong> <code>ssh deploy@... /usr/local/bin/myapp-api --version</code>
          , or whatever your binary prints, runs without a &quot;cannot execute&quot; error. If it
          fails with <code>Exec format error</code>, you built for the wrong architecture: an Arm
          server needs <code>GOARCH=arm64</code>.
        </p>
      </Step>

      <Step n={4} title="Make systemd responsible for it">
        <p>
          Running the binary over SSH works until you disconnect. A unit file hands it to the
          operating system.
        </p>
        <CodeBlock
          language="ini"
          filename="/etc/systemd/system/myapp-api.service"
          code={`[Unit]
Description=myapp API
After=network.target postgresql.service

[Service]
Type=simple
User=deploy
WorkingDirectory=/home/deploy/myapp
ExecStart=/usr/local/bin/myapp-api

# Restart on any exit, including a clean one, after five seconds. Without
# this a single panic takes the site down until somebody notices.
Restart=always
RestartSec=5

# Secrets live in a file only this user can read, not in this unit, which
# is world-readable.
EnvironmentFile=/home/deploy/myapp/.env

# Logs go to the journal, so journalctl can read them.
StandardOutput=journal
StandardError=journal

# Cheap hardening: no new privileges, a private /tmp, and a read-only
# system except for the one directory the app writes to.
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=/home/deploy/myapp

[Install]
WantedBy=multi-user.target`}
        />
        <CodeBlock
          terminal
          code={`sudo chmod 600 /home/deploy/myapp/.env     # before it has anything in it

sudo systemctl daemon-reload
sudo systemctl enable --now myapp-api
systemctl status myapp-api`}
        />
        <p>
          <strong>Checkpoint.</strong> <code>active (running)</code>, and{' '}
          <code>curl localhost:8080/health</code> on the server answers. When it does not:
        </p>
        <CodeBlock
          terminal
          code={`journalctl -u myapp-api -n 50 --no-pager   # the last 50 lines
journalctl -u myapp-api -f                 # follow it live`}
        />
        <p>
          Nine times out of ten the journal says a configuration variable is missing or the
          database refused the connection, and it says so in plain words.
        </p>
      </Step>

      <Step n={5} title="Caddy, for the proxy and the certificate">
        <p>
          This is the step that is traditionally painful and is now four lines, because Caddy
          requests and renews certificates on its own.
        </p>
        <CodeBlock
          terminal
          code={`sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' \\
  | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' \\
  | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install -y caddy`}
        />
        <CodeBlock
          language="text"
          filename="/etc/caddy/Caddyfile"
          code={`yourdomain.com {
	# The frontend: everything that is not /api
	reverse_proxy localhost:3000

	# The API
	handle /api/* {
		reverse_proxy localhost:8080
	}

	encode gzip zstd

	header {
		Strict-Transport-Security "max-age=31536000; includeSubDomains"
		X-Content-Type-Options "nosniff"
		X-Frame-Options "DENY"
		Referrer-Policy "strict-origin-when-cross-origin"
	}

	log {
		output file /var/log/caddy/access.log
	}
}

www.yourdomain.com {
	redir https://yourdomain.com{uri} permanent
}`}
        />
        <CodeBlock
          terminal
          code={`sudo systemctl reload caddy
journalctl -u caddy -f      # watch the certificate being issued`}
        />
        <p>
          <strong>Checkpoint.</strong> <code>https://yourdomain.com</code> loads with a padlock, and{' '}
          <code>http://</code> redirects to it. Caddy did that redirect without being asked, and it
          will renew the certificate without being asked either.
        </p>
        <Note title="Why not nginx">
          <p>
            nginx is excellent and is what you will meet in most existing systems. For a first
            deploy Caddy removes certbot, a renewal cron job and about forty lines of
            configuration, and the gap in capability does not matter until you are doing something
            unusual. Learn nginx when something requires it.
          </p>
        </Note>
      </Step>

      <Step n={6} title="The frontend">
        <p>
          A Next.js app with server components needs Node running. Building on a small server is
          the step that runs out of memory, so build where you have memory and ship the output.
        </p>
        <CodeBlock
          terminal
          code={`# locally
cd apps/web && pnpm build

# the standalone output includes only the dependencies actually used
rsync -avz --delete .next/standalone/ deploy@203.0.113.42:/home/deploy/myapp/web/
rsync -avz .next/static/ deploy@203.0.113.42:/home/deploy/myapp/web/.next/static/
rsync -avz public/ deploy@203.0.113.42:/home/deploy/myapp/web/public/`}
        />
        <p>Then a second unit, the same shape as the first:</p>
        <CodeBlock
          language="ini"
          filename="/etc/systemd/system/myapp-web.service"
          code={`[Unit]
Description=myapp web
After=network.target

[Service]
Type=simple
User=deploy
WorkingDirectory=/home/deploy/myapp/web
Environment=PORT=3000
Environment=NODE_ENV=production
ExecStart=/usr/bin/node server.js
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target`}
        />
        <p>
          A purely static frontend is simpler still: no Node, no unit, just{' '}
          <code>root * /home/deploy/myapp/web</code> and <code>file_server</code> in the Caddyfile.
        </p>
      </Step>

      <Step n={7} title="A deploy script, so the next one is one command">
        <p>
          You will deploy again today. Write this now, while every step is still in your head,
          because reconstructing it in three weeks is how a deploy becomes a nervous event.
        </p>
        <CodeBlock
          language="bash"
          filename="deploy.sh"
          code={`#!/usr/bin/env bash
set -euo pipefail      # stop on the first failure rather than carrying on

HOST=deploy@203.0.113.42
APP=/home/deploy/myapp

echo "Building..."
(cd apps/api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \\
  go build -ldflags="-s -w" -o api-linux ./cmd/server)
(cd apps/web && pnpm build)

echo "Uploading..."
scp apps/api/api-linux "$HOST:/tmp/api-new"
rsync -avz --delete apps/web/.next/standalone/ "$HOST:$APP/web/"
rsync -avz apps/web/.next/static/ "$HOST:$APP/web/.next/static/"

echo "Switching over..."
ssh "$HOST" bash -s <<'REMOTE'
set -euo pipefail
sudo mv /tmp/api-new /usr/local/bin/myapp-api
sudo chmod +x /usr/local/bin/myapp-api
sudo systemctl restart myapp-api myapp-web

# Wait for health rather than assuming. A deploy that reports success
# while the service is crash-looping is worse than one that fails.
for i in $(seq 1 30); do
  if curl -fsS localhost:8080/health > /dev/null; then
    echo "healthy"; exit 0
  fi
  sleep 1
done
echo "did not become healthy; check journalctl -u myapp-api" >&2
exit 1
REMOTE`}
        />
        <p>
          There is a gap of a second or two where requests fail. For a first deploy that is
          acceptable. Closing it means running two copies and switching the proxy between them,
          which is a worthwhile second step and not a first one.
        </p>
      </Step>

      <Section id="database" title="The database, and the backup you have not made">
        <p>
          Run Postgres on the same box to begin with. It is one <code>apt install</code>, and the
          network hop disappears. Move it to a managed service when backups and failover matter
          more than the bill.
        </p>
        <CodeBlock
          terminal
          code={`sudo apt install -y postgresql
sudo -u postgres createuser --pwprompt myapp
sudo -u postgres createdb --owner=myapp myapp`}
        />
        <p>
          It listens on localhost only by default. Leave it that way. To connect from your laptop,
          tunnel:
        </p>
        <CodeBlock
          terminal
          code={`ssh -L 5433:localhost:5432 deploy@203.0.113.42
# then connect to localhost:5433 as if it were local`}
        />
        <Warn title="A backup you have not restored is not a backup">
          <p>
            Dump nightly, copy the dump off the server, and once, now, restore it into a scratch
            database and look at the data. People discover their backup was empty, or was the
            wrong database, at the only moment it matters.
          </p>
          <CodeBlock
            language="bash"
            code={`#!/usr/bin/env bash
set -euo pipefail
pg_dump -Fc myapp > /tmp/myapp-$(date +%F).dump
# then copy it OFF this machine: object storage, another provider, anywhere
# that is not the disk you are backing up.`}
          />
        </Warn>
      </Section>

      <Section id="mistakes" title="When it goes wrong">
        <Mistake error="502 Bad Gateway">
          <p>
            Caddy is working and your app is not answering on the port Caddy forwards to. Check the
            service is running, and check what it is actually listening on: an app bound to{' '}
            <code>127.0.0.1:8080</code> is fine here, but one bound to a different port than the
            Caddyfile names is the usual cause.
          </p>
          <CodeBlock
            terminal
            code={`systemctl status myapp-api
sudo ss -tlnp | grep 8080`}
          />
        </Mistake>

        <Mistake error="Job for myapp-api.service failed because the control process exited with error code">
          <p>
            The binary started and immediately stopped. The unit file tells you nothing; the
            journal tells you everything:
          </p>
          <CodeBlock terminal code={`journalctl -u myapp-api -n 50 --no-pager`} />
          <p>
            The frequent causes are a missing environment variable, a database that is not
            running, a port already in use, and the <code>EnvironmentFile</code> path being wrong,
            which fails silently as &quot;variable not set&quot;.
          </p>
        </Mistake>

        <Mistake error="Caddy cannot obtain a certificate">
          <p>
            Three causes, in order of likelihood. DNS does not yet point at this server, so check{' '}
            <code>dig +short yourdomain.com</code> from somewhere other than the server. Port 80 is
            closed, and the challenge needs it even though your site is HTTPS. Or you have been
            retrying and hit Let&apos;s Encrypt&apos;s rate limit, which is five failures per hour
            per domain, so stop, fix DNS properly, and wait.
          </p>
        </Mistake>

        <Mistake error="Permission denied on the binary or the env file">
          <p>
            The unit runs as <code>deploy</code>, so <code>deploy</code> must be able to execute
            the binary and read <code>.env</code>. A file uploaded as root and left at mode 600 is
            unreadable by the service. <code>ls -l</code> on both, and remember that{' '}
            <code>ProtectSystem=strict</code> makes everything except{' '}
            <code>ReadWritePaths</code> read-only to the service, which is deliberate and
            occasionally surprising.
          </p>
        </Mistake>

        <Mistake error="The site works and then stops after a few days">
          <p>
            Usually the disk filled with logs, or the process was killed by the out-of-memory
            killer. Both are visible after the fact:
          </p>
          <CodeBlock
            terminal
            code={`df -h                             # disk
journalctl --disk-usage           # how much of it is logs
sudo journalctl --vacuum-time=7d  # keep a week
dmesg | grep -i "killed process"  # the OOM killer, if it was that`}
          />
        </Mistake>
      </Section>

      <Section id="alternatives" title="If you would rather not do any of that">
        <p>
          A platform does all of the above for you, at a higher price per unit of compute and with
          less to understand.
        </p>
        <Checklist
          items={[
            <>
              <strong>Railway, Render, Fly.</strong> Push to git, get a URL and a database. Pricing
              grows with usage, and at steady load a VPS is several times cheaper. The right choice
              when your time is the scarce resource.
            </>,
            <>
              <strong>Dokploy or Coolify on your own VPS.</strong> A deployment platform you host:
              the git push workflow and a web interface, on hardware you rent. A good middle
              ground, and more moving parts than this page.
            </>,
            <>
              <strong>Kubernetes.</strong> Not for a first deploy, and not for a second. It solves
              problems you do not yet have, at a cost in understanding that you will pay every
              week.
            </>,
          ]}
        />
      </Section>

      <Section id="with-grit" title="One command, for the same thing">
        <p>
          Grit automates exactly the sequence above: cross-compile, upload, write the systemd unit,
          configure Caddy, issue the certificate, run migrations, and health-check before reporting
          success.
        </p>
        <CodeBlock terminal code={`grit deploy --host yourdomain.com`} />
        <p>
          It writes real systemd and Caddy configuration onto the server, which you can read and
          edit afterwards. That matters more than the convenience: the day something is wrong at
          three in the morning, you are reading a unit file like the one above and not a
          platform&apos;s dashboard.{' '}
          <Link href="/docs/deployment/vps" className="text-accent hover:underline">
            The VPS deployment guide
          </Link>{' '}
          has the full flow, including the rollback.
        </p>
      </Section>

      <Faq
        slug="deploy-go-app-vps-https"
        items={[
          {
            q: 'How much server do I need?',
            a: (
              <>
                Less than you think. A Go API serves thousands of requests a second on two cores,
                and the memory is dominated by Postgres and Node rather than by Go. Two vCPU and
                2GB is a comfortable start for a real product. Scale up when a measurement says to,
                not when traffic feels like it is growing.
              </>
            ),
            text:
              'Less than you expect. A Go API handles thousands of requests per second on two cores, and memory is dominated by Postgres and Node rather than Go. Two vCPU and 2GB is a comfortable start for a real product. Scale up on a measurement rather than a feeling.',
          },
          {
            q: 'Should I deploy with Docker instead of a binary?',
            a: (
              <>
                Docker makes the server identical to your development environment and adds a layer
                to debug through. A Go binary already has no runtime dependencies, so containers
                buy less here than in most ecosystems. Use Docker when you are running several
                services together, or when the platform you are deploying to expects an image.
              </>
            ),
            text:
              'Docker makes the server match your development environment at the cost of a layer to debug through. A Go binary already has no runtime dependencies, so containers buy less here than in most ecosystems. Use Docker when running several services together or when the target platform expects an image.',
          },
          {
            q: 'How do I deploy without downtime?',
            a: (
              <>
                Run two copies on different ports and point the proxy at the healthy one, switching
                after the new copy passes its health check. Caddy and nginx both support this. It
                is a worthwhile second step; a one-second gap during a deploy is not what is
                costing you users on day one.
              </>
            ),
            text:
              'Run two copies on different ports and switch the proxy to the new one after it passes a health check. Caddy and nginx both support it. It is a good second step, but a one-second gap during a deploy is not the thing costing you users at the start.',
          },
          {
            q: 'Is a VPS secure enough for real customer data?',
            a: (
              <>
                Yes, with the basics on this page: key-only SSH, a firewall, unattended security
                updates, a non-root service user, and the database not listening on a public
                address. What a managed platform really buys you is that somebody else remembers to
                apply patches. If that somebody is not going to be you, use a platform.
              </>
            ),
            text:
              'Yes, with the basics: key-only SSH, a firewall, unattended security updates, a non-root service user and a database that does not listen publicly. What a managed platform really provides is somebody else remembering to patch. If that will not be you, use a platform.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/docker-compose-postgres-redis',
            label: 'Postgres and Redis with Docker Compose',
            blurb: 'The local stack this is the production counterpart of.',
          },
          {
            href: '/guides/first-go-react-app',
            label: 'Build your first full-stack app',
            blurb: 'The application this deploys.',
          },
          {
            href: '/docs/deployment/vps',
            label: 'VPS deployment',
            blurb: 'The automated flow, including migrations and rollback.',
          },
          {
            href: '/docs/deployment/checklist',
            label: 'Deployment checklist',
            blurb: 'What to confirm before the first real user arrives.',
          },
        ]}
      />
    </GuidePage>
  )
}
