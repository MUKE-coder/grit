package scaffold

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// randomHex returns 2*n hex characters of cryptographically random bytes.
// Used to seed JWT, Sentinel and Pulse secrets at scaffold time so a fresh
// `grit new` project boots cleanly in production mode without the user
// having to manually replace placeholder credentials.
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand should never fail on a real OS; if it did the scaffold
		// would be unusable. Fall back to a clearly-fake string so the user
		// at least sees something obvious to change rather than a silent
		// weak default.
		return "REPLACE_ME_crypto_rand_failed"
	}
	return hex.EncodeToString(b)
}

// randomBase64Key returns 32 random bytes in standard base64: the shape
// FIELD_ENCRYPTION_KEY takes, 44 characters.
func randomBase64Key() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// Not base64, so the API refuses to start and names the variable, rather
		// than running with a key nobody generated.
		return "REPLACE_ME_crypto_rand_failed"
	}
	return base64.StdEncoding.EncodeToString(b)
}

// fieldEncryptionKeyComment introduces FIELD_ENCRYPTION_KEY in .env.
const fieldEncryptionKeyComment = `# Encrypts two-factor secrets and every crypto.EncryptedString column
# (AES-256-GCM, 32 bytes in base64). Back it up with the database, but not in the
# same place: losing it locks every two-factor user out and leaves encrypted
# columns unreadable, and changing it does the same. Generate a new one for each
# environment: openssl rand -base64 32
`

// envExampleFile is .env with every generated secret replaced by a placeholder.
//
// .env.example is meant to be committed, and it was a byte-for-byte copy of .env:
// a real JWT secret, database password and dashboard passwords in the repository
// from the first git add, and in production for anybody who deployed with
// cp .env.example .env. The shape stays, so it still documents every variable.
func envExampleFile(opts Options) string {
	out := generatedSecret.ReplaceAllString(envFile(opts), "${1}=CHANGE_ME${2}")
	return encryptionKeyLine.ReplaceAllString(out, "${1}=CHANGE_ME${2}")
}

var (
	// generatedSecret matches a KEY=value line whose value is what randomHex
	// produces: 32 or more lowercase hex characters.
	generatedSecret = regexp.MustCompile(`(?m)^([A-Z][A-Z0-9_]*)=[0-9a-f]{32,}(\r?)$`)
	// encryptionKeyLine covers the field-encryption key, which is base64.
	encryptionKeyLine = regexp.MustCompile(`(?m)^(FIELD_ENCRYPTION_KEY)=\S+(\r?)$`)
)

func writeRootFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, ".env"):                                      envFile(opts),
		filepath.Join(root, ".env.example"):                              envExampleFile(opts),
		filepath.Join(root, ".gitignore"):                                rootGitignore(),
		filepath.Join(root, "README.md"):                                 readmeFile(opts),
		filepath.Join(root, "grit.json"):                                 gritJSON(opts),
		filepath.Join(root, ".claude", "skills", "grit", "SKILL.md"):     gritSkillFile(opts),
		filepath.Join(root, ".claude", "skills", "grit", "reference.md"): gritSkillReference(opts),
	}

	// Biome, the linter and formatter, sits where pnpm runs: the workspace
	// root, or a single app's frontend/.
	if opts.ShouldUseTurborepo() {
		files[filepath.Join(root, biomeConfigFile)] = biomeConfig(false)
	} else if opts.Architecture == ArchSingle {
		files[filepath.Join(root, "frontend", biomeConfigFile)] = biomeConfig(true)
	}

	if opts.ShouldUseTurborepo() {
		files[filepath.Join(root, "pnpm-workspace.yaml")] = pnpmWorkspace(opts.ShouldIncludeDesktop(), opts.ShouldIncludeExpo())
		files[filepath.Join(root, "turbo.json")] = turboJSON()
		files[filepath.Join(root, "package.json")] = rootPackageJSON(opts)
		files[filepath.Join(root, "grit.config.ts")] = gritConfig(opts)
		files[filepath.Join(root, ".npmrc")] = rootNpmrc()
	}

	// Single-binary apps keep their SPA in frontend/ and have no workspace root,
	// so they never received an .npmrc. Without it `pnpm install` exits non-zero
	// with ERR_PNPM_IGNORED_BUILDS on esbuild — a failed install on a freshly
	// generated project.
	if opts.Architecture == ArchSingle {
		files[filepath.Join(root, "frontend", ".npmrc")] = rootNpmrc()
		files[filepath.Join(root, "frontend", "pnpm-workspace.yaml")] = pnpmAllowBuilds()
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

// adminURLEnv is the NEXT_PUBLIC_ADMIN_URL block, which depends on where this
// project's admin panel actually is.
//
// It is a URL only when the panel is its own application. In a double it is a
// route inside the web app, and setting the variable to localhost:3001 pointed
// the navbar at a port nothing listens on; in a single there is no Next.js app
// to read it at all. The web app compiles in the right default, so the variable
// is there to override it, not to restate it.
func adminURLEnv(opts Options) string {
	switch {
	case opts.ShouldIncludeAdmin():
		return `# The admin panel runs as its own app, on its own port. Set this to your
# production admin origin (e.g. https://admin.example.com) before shipping.
NEXT_PUBLIC_ADMIN_URL=http://localhost:3001`
	case opts.ShouldEmbedAdmin():
		return `# The admin panel is a route group inside this web app, so the link to it is a
# path and the app already knows it: /admin/dashboard. Set this only if you move
# the panel to an origin of its own.
# NEXT_PUBLIC_ADMIN_URL=`
	default:
		// No Next.js web app in this shape, so nothing reads it.
		return ""
	}
}

func envFile(opts Options) string {
	// Generated per-scaffold so APP_ENV=production works out of the box.
	// Rotate any of these any time with `openssl rand -hex 32`.
	jwtSecret := randomHex(32)
	sentinelPassword := randomHex(16)
	sentinelSecretKey := randomHex(32)
	sentinelAuditKey := randomHex(32)
	pulsePassword := randomHex(16)
	studioPassword := randomHex(16)
	postgresPassword := randomHex(24) // strong default; new project just works
	provider := opts.DBProvider
	if provider == "" {
		provider = "postgres"
	}

	out := fmt.Sprintf(`# %s: Environment Variables

# App
APP_NAME=%s
# development or production. Unset means production, the strict one: rate limits
# on, the WAF blocking, GORM Studio and /docs off, and a default or short secret
# stops the server from starting.
APP_ENV=development
# The password grit seed gives admin@example.com. Unset, APP_ENV=development
# seeds admin123 and four demo accounts; any other APP_ENV refuses to seed the
# admin without a password of 12 or more characters, and skips the demo accounts.
# SEED_ADMIN_PASSWORD=
# Serve the API reference at /docs in production as well. It maps every route
# and has a console that calls them, so it is off there unless you say so.
API_DOCS_PUBLIC=false
APP_PORT=8080
APP_URL=http://localhost:8080

# ─── Database ───────────────────────────────────────────────────────────
# Which engine this project talks to. One of:
#
#   postgres   the default, and what docker-compose.yml starts for you
#   mysql      MySQL 8 or MariaDB
#   sqlite     one file, pure Go, no CGO and no server to run
#   memory     SQLite in RAM: empty at every boot, for tests and demos
#
# Only the block for the provider you choose is read. The others can stay.
# Picked by 'grit new --db <provider>'; change it here any time.
DB_PROVIDER=%s

# Postgres — read when DB_PROVIDER=postgres
# Single source of truth: docker-compose.yml reads the same values via ${VAR},
# so they cannot drift. POSTGRES_PASSWORD is generated per scaffold, so
# 'grit migrate' works on a fresh project without any editing.
POSTGRES_USER=grit
POSTGRES_PASSWORD=%s
POSTGRES_DB=%s
POSTGRES_HOST=localhost
# 5434 (not the default 5432) avoids host collisions — see docker-compose.yml.
# The container still listens on 5432 inside the Docker network, which is
# why docker-compose.prod.yml overrides POSTGRES_PORT back to 5432 for
# inter-container traffic.
POSTGRES_PORT=5434
# require / verify-full for a managed Postgres that insists on TLS.
POSTGRES_SSLMODE=disable

# MySQL — read when DB_PROVIDER=mysql
# Nothing in docker-compose.yml starts MySQL: point these at your own server,
# or run one with
#   docker run -d --name mysql -p 3306:3306 \
#     -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=%s \
#     -e MYSQL_USER=grit -e MYSQL_PASSWORD=grit mysql:8
MYSQL_USER=grit
MYSQL_PASSWORD=grit
MYSQL_DB=%s
MYSQL_HOST=localhost
MYSQL_PORT=3306

# SQLite — read when DB_PROVIDER=sqlite
# A path, relative to apps/api. Add it to .gitignore.
SQLITE_PATH=./app.db
# ─── Docker host ports ──────────────────────────────────────────────────
#
# What docker-compose binds on your machine. Every Grit project defaults to
# the same numbers, so the second one you start fails with "port is already
# allocated". Change them here: compose binds them, and the API builds its
# Postgres, Redis and MinIO addresses from them unless a URL below is set.
#
# Only the host side moves. Inside the compose network the services keep
# their standard ports, so nothing else needs to know.
REDIS_PORT=6380
# The production Redis (docker-compose.prod.yml) takes this password. The
# development one has none.
REDIS_PASSWORD={{REDIS_PASSWORD}}
MAILHOG_SMTP_PORT=1025
MAILHOG_UI_PORT=8025
MINIO_PORT=9002
MINIO_CONSOLE_PORT=9003
{{MINIO_BIND_ADDRESS}}
# Override the connection string ONLY if you're pointing at an external
# Postgres (Neon, Supabase, RDS) or want to use SQLite. When set, this
# wins over the POSTGRES_* parts above.
#   DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=require
#   DATABASE_URL=sqlite:./app.db           # pure-Go driver, no CGO
#   DATABASE_URL=sqlite::memory:           # gone on restart, great for tests
# DATABASE_URL=

`+envPoolNew+`
# JWT — generated at scaffold time. Rotate with: openssl rand -hex 32
JWT_SECRET=%s
`+fieldEncryptionKeyComment+`FIELD_ENCRYPTION_KEY={{FIELD_ENCRYPTION_KEY}}
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# OAuth2 — Social Login (Google + GitHub)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=
OAUTH_FRONTEND_URL=http://localhost:3001

# API keys are created from the admin at /system/api-keys — there is no env
# var for them. A key is sent as "X-API-Key: grit_..." or
# "Authorization: Bearer grit_...", and acts as the user who created it.

# Account lockout. Sentinel rate-limits by IP; this protects a single account
# from attempts spread across many addresses. Set LOGIN_MAX_ATTEMPTS=0 to
# disable. Only wrong passwords on real accounts count — an unknown email never
# locks anything, or anyone could lock an address they guess.
LOGIN_MAX_ATTEMPTS=10
LOGIN_LOCKOUT_MINUTES=15

# Require a confirmed email address before a password sign-in is accepted.
# Leave false on an existing project: every current user has a NULL
# email_verified_at and would be locked out the moment you switch it on.
REQUIRE_EMAIL_VERIFICATION=false

# Redis
# Unset, the API connects to localhost on REDIS_PORT above, so moving the port
# moves the connection with it. Set REDIS_URL for an external Redis, or to empty
# (REDIS_URL=) to run without Redis at all: cache, background jobs and cron then
# stay off instead of retrying a dial in a loop.
# REDIS_URL=redis://localhost:6380

# Public API URL — baked into Next.js bundles at build time
API_URL=http://localhost:8080

{{ADMIN_URL_ENV}}

`+envStorageHeadNew+`
# Browser-facing storage origin. Uploads are presigned PUTs the browser makes
# straight to object storage, and stored images load from the same host, so
# this origin must be allowed by the frontend Content-Security-Policy. In dev
# the admin/web CSP already defaults to the MinIO endpoint below. In PRODUCTION,
# set NEXT_PUBLIC_STORAGE_URL (Next.js apps) or VITE_STORAGE_URL (Vite apps) in
# the frontend environment to your public storage origin — e.g.
# https://<bucket>.s3.<region>.amazonaws.com or your CDN domain — or presigned
# uploads and image display will be blocked by CSP.

# MinIO (local development — used when STORAGE_DRIVER=minio)
# Unset, built from MINIO_PORT above. Set it only for a MinIO somewhere else.
# MINIO_ENDPOINT=http://localhost:9002
MINIO_ACCESS_KEY={{MINIO_ACCESS_KEY}}
MINIO_SECRET_KEY={{MINIO_SECRET_KEY}}
MINIO_BUCKET=%s-uploads
MINIO_REGION=us-east-1
MINIO_USE_SSL=false

# AWS S3 (used when STORAGE_DRIVER=s3)
# Leave S3_ENDPOINT empty to use the AWS regional default.
# S3_ACCESS_KEY + S3_SECRET_KEY fall back to AWS_ACCESS_KEY_ID +
# AWS_SECRET_ACCESS_KEY (and S3_REGION to AWS_REGION) so an IAM role
# attached to your EC2 / ECS / Lambda Just Works.
S3_ENDPOINT=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=
S3_REGION=us-east-1
# Optional. Set to a CloudFront/CDN domain, or to
# https://<bucket>.s3.<region>.amazonaws.com if the bucket is public.
S3_PUBLIC_URL=

# Cloudflare R2 (used when STORAGE_DRIVER=r2)
#
# R2_PUBLIC_URL is REQUIRED for images to display. The R2 endpoint below is
# the S3 API host and only answers signed requests, so an <img> pointed at it
# gets 401 — uploads succeed and nothing renders. Get a public origin by
# either enabling r2.dev on the bucket (Settings -> Public Development URL)
# or binding a custom domain, then paste it here WITHOUT the bucket name.
#
# R2 also needs a CORS rule for the presigned upload itself:
#   AllowedOrigins: your admin/web origins  AllowedMethods: GET, PUT
#   AllowedHeaders: *                       ExposeHeaders: ETag
R2_ENDPOINT=
R2_ACCESS_KEY=
R2_SECRET_KEY=
R2_BUCKET=
R2_REGION=auto
R2_PUBLIC_URL=

# Browser-facing origin for any driver. Per-driver vars
# (S3_PUBLIC_URL / R2_PUBLIC_URL / B2_PUBLIC_URL / MINIO_PUBLIC_URL) win over
# this one. Leave empty for MinIO in dev — it serves objects from the same
# host it takes API calls on.
STORAGE_PUBLIC_URL=

# Backblaze B2 (used when STORAGE_DRIVER=b2)
B2_ENDPOINT=
B2_ACCESS_KEY=
B2_SECRET_KEY=
B2_BUCKET=
B2_REGION=us-west-004

`+envMailHeadNew+`MAIL_FROM=noreply@%s.dev
`+envMailDriversNew+`
# Support inbox — every ticket opened in /system/support is emailed here
# (when RESEND_API_KEY is set). Leave empty in dev to skip email-out.
SUPPORT_EMAIL=

# CORS
# Browser origins allowed to call the API. The Wails desktop webview does NOT
# need an entry here: its origin varies (http://wails.localhost:34115 in dev,
# http://wails.localhost in a Windows build, wails://wails on macOS/Linux), so
# the CORS middleware matches the wails.localhost host on any port instead.
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# GORM Studio
# Studio browses and edits every table, so its password is generated per
# scaffold like the dashboard ones. Rotate with: openssl rand -hex 16
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin
GORM_STUDIO_PASSWORD=%s
GORM_STUDIO_READ_ONLY=false
GORM_STUDIO_DISABLE_SQL=false
# In production Studio stays off whatever GORM_STUDIO_ENABLED says, because this
# file gets copied to servers whole. This turns it on there, read-only and with
# no SQL editor.
GORM_STUDIO_IN_PRODUCTION=false

# ============================================
# Optional modules
# ============================================
# Grit ships every battery enabled. Switch one off and it mounts no routes,
# starts no workers, creates no tables, and disappears from the admin nav.
# The code stays in the repo — delete it by hand if you want it gone entirely.
# All default to true, so leaving these unset changes nothing.
MODULE_AI=true
MODULE_JOBS=true
MODULE_CRON=true
MODULE_BACKUP=true
MODULE_WEBHOOKS=true
MODULE_REALTIME=true
MODULE_FILES=true
MODULE_MAIL=true
MODULE_AUDIT=true
MODULE_FLAGS=true
MODULE_TWOFACTOR=true

# AI — Vercel AI Gateway (one key, hundreds of models)
#
# The default model needs paid credits. A free-tier key authenticates fine and
# then refuses the request with "Free tier users do not have access to this
# model", so if /api/v1/ai/complete returns AI_FORBIDDEN, this line is why:
# either add credits or set a model your plan covers.
AI_GATEWAY_API_KEY=                           # Get from vercel.com/ai-gateway
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model format
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1

# Two-Factor Authentication (TOTP)
TOTP_ISSUER=%s

# Observability — Pulse performance monitoring dashboard
# Pulse refuses to mount in APP_ENV=production with the literal default
# password "pulse" — the value below is generated at scaffold time so the
# gate is satisfied. Rotate with: openssl rand -hex 16
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=%s

# Security — Sentinel WAF, rate limiting, threat detection
# Same as Pulse: Sentinel refuses to mount in production with default
# credentials. SECRET_KEY needs at least 32 bytes of entropy — both values
# below are generated per scaffold. Rotate with: openssl rand -hex 16 (password)
# / openssl rand -hex 32 (secret key).
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=%s
SENTINEL_SECRET_KEY=%s
# Keys the security audit log's hash chain, so an entry edited by someone with
# database access but not this key fails verification. Keep it: entries
# written under one key verify only with that key.
SENTINEL_AUDIT_KEY=%s
{{TRUSTED_PROXIES}}
# ─── Theme (v3.28+) ────────────────────────────────────────────────────
# Picks the visual identity for auth pages and the dashboard. Options:
#   atlas  — split-screen, team/organisation, Inter (default)
#   aurora — centered Clerk-style, consumer SaaS, Geist
#   pulse  — split-screen with carousel, ecommerce/brand, Onest + DM Serif
# Mirrored into NEXT_PUBLIC_THEME for the web + admin clients in
# next.config — the apps read it at build time so server components render
# the right theme without a flash of unstyled content.
THEME=%s

# Social login buttons (Google + GitHub). The OAuth API routes stay
# registered server-side either way; this only changes what the UI renders.
# NEXT_PUBLIC_SOCIAL_AUTH_ENABLED / EXPO_PUBLIC_SOCIAL_AUTH_ENABLED mirror
# the value to the browser and mobile bundles.
#
# Off by default because no provider is configured yet: GOOGLE_CLIENT_ID and
# GITHUB_CLIENT_ID above are empty, so the buttons would render and then fail
# with "no provider for google exists". Fill in a provider's credentials,
# then flip this to true.
SOCIAL_AUTH_ENABLED=false
`,
		opts.ProjectName, opts.ProjectName, // banner + APP_NAME
		provider,                           // DB_PROVIDER
		postgresPassword, opts.ProjectName, // POSTGRES_PASSWORD + POSTGRES_DB
		opts.ProjectName, opts.ProjectName, // the MySQL database, in the docker run hint and MYSQL_DB
		jwtSecret,
		opts.ProjectName, opts.ProjectName, // MINIO_BUCKET + MAIL_FROM
		studioPassword, opts.ProjectName, // GORM_STUDIO_PASSWORD + TOTP_ISSUER
		pulsePassword, sentinelPassword, sentinelSecretKey, sentinelAuditKey,
		opts.Theme, // THEME — picked by --theme at scaffold time, defaults to atlas
	)

	// Every project encrypts its two-factor secrets and encrypted columns from the
	// start. .env.example gets a placeholder instead: see encryptionKeyLine.
	out = strings.Replace(out, "{{FIELD_ENCRYPTION_KEY}}", randomBase64Key(), 1)

	// MinIO's root credentials, generated like the rest. They were
	// minioadmin/minioadmin, which is the whole bucket to anyone who can reach it.
	out = strings.Replace(out, "{{MINIO_ACCESS_KEY}}", "grit"+randomHex(8), 1)
	out = strings.Replace(out, "{{MINIO_SECRET_KEY}}", randomHex(24), 1)
	out = strings.Replace(out, "{{REDIS_PASSWORD}}", randomHex(24), 1)
	out = strings.Replace(out, "{{MINIO_BIND_ADDRESS}}", minioBindEnv(opts), 1)
	out = strings.Replace(out, "{{TRUSTED_PROXIES}}", envTrustedProxies, 1)

	// The admin panel's URL is not a URL in every shape: see adminURLEnv.
	return strings.Replace(out, "{{ADMIN_URL_ENV}}", adminURLEnv(opts), 1)
}
func envCloudExampleFile(opts Options) string {
	return fmt.Sprintf(`# %s: Cloud Environment Variables
#
# Use this file if you DON'T have Docker and want to use cloud services instead.
# Copy this to .env and fill in your keys:
#
#   cp .env.cloud.example .env
#
# No Docker required — just your API keys.

# ─── App ───────────────────────────────────────────────
APP_NAME=%s
APP_ENV=development
APP_PORT=8080
APP_URL=http://localhost:8080

# ─── Database (Neon — https://neon.tech) ───────────────
# Create a free project at neon.tech, copy the connection string
DATABASE_URL=postgres://user:password@ep-xxx-xxx-123456.us-east-2.aws.neon.tech/neondb?sslmode=require

# ─── JWT ───────────────────────────────────────────────
JWT_SECRET=change-me-to-a-random-string-at-least-32-chars
`+fieldEncryptionKeyComment+`FIELD_ENCRYPTION_KEY=CHANGE_ME
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# ─── Redis (Upstash — https://upstash.com) ─────────────
# Create a free Redis database at upstash.com, copy the Redis URL
REDIS_URL=rediss://default:your-password@your-endpoint.upstash.io:6379

# ─── Storage ──────────────────────────────────────────
# Active driver: r2 or b2 (no minio in cloud mode)
STORAGE_DRIVER=r2

# ─── Cloudflare R2 (https://dash.cloudflare.com) ─────
# Dashboard → R2 → Create Bucket → Manage R2 API Tokens
R2_ENDPOINT=https://your-account-id.r2.cloudflarestorage.com
R2_ACCESS_KEY=your-r2-access-key-id
R2_SECRET_KEY=your-r2-secret-access-key
R2_BUCKET=%s-uploads
R2_REGION=auto

# ─── Backblaze B2 (https://backblaze.com/cloud-storage) ─
# B2 Cloud Storage → Buckets → Create → App Keys → Add Key
B2_ENDPOINT=https://s3.us-west-004.backblazeb2.com
B2_ACCESS_KEY=your-b2-key-id
B2_SECRET_KEY=your-b2-application-key
B2_BUCKET=%s-uploads
B2_REGION=us-west-004

# ─── Email ─────────────────────────────────────────────
# MAIL_MAILER picks how mail is sent: resend, smtp, mailgun, postmark,
# sendgrid, ses or failover. Left empty, a real RESEND_API_KEY means Resend.
# MAIL_MAILER=
RESEND_API_KEY=re_your_api_key_here
MAIL_FROM=noreply@yourdomain.com
MAIL_FROM_NAME=
# MAIL_FAILOVER=resend,smtp

# SMTP (MAIL_MAILER=smtp). SMTP_ENCRYPTION is tls, starttls or none.
# SMTP_HOST=smtp.yourprovider.com
# SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_ENCRYPTION=

# Mailgun (MAIL_MAILER=mailgun). api.eu.mailgun.net for a domain in the EU.
MAILGUN_DOMAIN=
MAILGUN_SECRET=
MAILGUN_ENDPOINT=api.mailgun.net

# Postmark (MAIL_MAILER=postmark)
POSTMARK_TOKEN=
POSTMARK_MESSAGE_STREAM=outbound

# SendGrid (MAIL_MAILER=sendgrid)
SENDGRID_API_KEY=

# Amazon SES API v2 (MAIL_MAILER=ses), with the usual AWS keys.
AWS_SES_REGION=us-east-1

# ─── CORS ──────────────────────────────────────────────
# Browser origins allowed to call the API. The Wails desktop webview does NOT
# need an entry here: its origin varies (http://wails.localhost:34115 in dev,
# http://wails.localhost in a Windows build, wails://wails on macOS/Linux), so
# the CORS middleware matches the wails.localhost host on any port instead.
CORS_ORIGINS=http://localhost:3000,http://localhost:3001

# ─── GORM Studio ──────────────────────────────────────
GORM_STUDIO_ENABLED=true
GORM_STUDIO_USERNAME=admin               # Login username for the Studio UI
GORM_STUDIO_PASSWORD=change-me-in-prod   # Login password — CHANGE THIS in production!
GORM_STUDIO_READ_ONLY=false             # Refuse every write from Studio
GORM_STUDIO_DISABLE_SQL=true            # The SQL editor bypasses every GORM guard

# ─── AI (Vercel AI Gateway) ──────────────────────────
AI_GATEWAY_API_KEY=your-gateway-key  # Get from vercel.com/ai-gateway
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6  # provider/model format
AI_GATEWAY_URL=https://ai-gateway.vercel.sh/v1

# ─── TOTP (Two-Factor Authentication) ────────────────
TOTP_ISSUER=%s

# ─── Observability (Pulse) ────────────────────────────
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=change-me-in-production

# ─── Security (Sentinel) ─────────────────────────────
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=change-me-in-production
SENTINEL_SECRET_KEY=generate-a-random-string-here
SENTINEL_AUDIT_KEY=generate-a-random-string-here
`, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName, opts.ProjectName)
}

func rootGitignore() string {
	return `# Grit local state.
# .grit/manifest.json is deliberately NOT ignored: it records which files Grit
# generated and what it wrote, which is how grit upgrade knows to leave your
# edits alone. Commit it so the whole team gets that. The rest is local.
.grit/backups/
.grit/slots.json

# Dependencies
node_modules/
.pnpm-store/

# Build output
.next/
out/
dist/
build/
tmp/
# TypeScript, Next.js and Expo write these on every build or dev run.
*.tsbuildinfo
next-env.d.ts
.expo/

# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
migrate.exe

# Environment
.env
.env.local
.env.*.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Docker volumes
postgres-data/
redis-data/
minio-data/

# Turborepo
.turbo/

# Local databases. SQLite files hold real data: users with their password
# hashes, sessions and API keys, as well as Sentinel's WAF log.
*.db
*.db-shm
*.db-wal
*.sqlite
*.sqlite3

# Testing
e2e/test-results/
e2e/playwright-report/
coverage/

# Debug
*.log
npm-debug.log*
pnpm-debug.log*
`
}

// prettierConfig and prettierIgnore are what Grit wrote before Biome replaced
// Prettier. New projects do not get them; upgrade compares against them to
// recognise a copy nobody has edited, which is the only kind it removes.
func prettierConfig() string {
	return `{
  "semi": false,
  "singleQuote": true,
  "trailingComma": "es5",
  "tabWidth": 2,
  "printWidth": 100,
  "bracketSpacing": true,
  "arrowParens": "always",
  "endOfLine": "lf",
  "plugins": ["prettier-plugin-tailwindcss"]
}
`
}

func prettierIgnore() string {
	return `node_modules/
.next/
dist/
build/
out/
coverage/
.turbo/
pnpm-lock.yaml
*.min.js
*.min.css
`
}

func pnpmWorkspace(includeDesktop, includeExpo bool) string {
	ws := `packages:
  - "apps/*"
  - "packages/*"
`
	// The Wails desktop client keeps its React app in apps/desktop/frontend,
	// which "apps/*" doesn't match. Add it so a single root ` + "`pnpm install`" + `
	// covers the desktop frontend too (otherwise it only gets installed the
	// first time you run ` + "`wails dev`" + `).
	if includeDesktop {
		ws += `  - "apps/desktop/frontend"
`
	}
	// React/react-dom are pinned to one exact version directly in every
	// frontend package.json (see the *_files.go templates) so the two always
	// match — React 19 hard-errors on a version mismatch. Exact direct pins are
	// more reliable than a pnpm override (pnpm 10 didn't consistently apply the
	// react-dom override, leaving react and react-dom on different 19.x lines).
	ws += "\n" + pnpmAllowBuilds()
	if includeExpo {
		ws += pnpmExpoAudit
	}
	return ws
}

// pnpmExpoAudit keeps the security workflow's pnpm audit passing on a project
// with the Expo app. Every high advisory it reported was in Expo's bundler,
// which is a production dependency of the Expo app, so the audit counts it.
//
// PostCSS is moved to a fixed release: 8.4 to 8.5 is compatible, and the app
// bundles with it. image-size has no fixed 1.x release, only 2.x, which metro
// cannot take, and metro only reads the sizes of the project's own image
// assets at bundle time, never a file a user uploads, so those two are
// accepted here. Remove them when metro moves to image-size 2.
const pnpmExpoAudit = `
# Expo's bundler: see pnpmExpoAudit in Grit's scaffold for the reasoning.
overrides:
  "@expo/metro-config>postcss": "^8.5.28"

auditConfig:
  ignoreGhsas:
    # image-size 1.x in metro: denial of service parsing ICNS, JXL and HEIF.
    # Only the project's own image assets are read, at bundle time.
    - GHSA-w3rx-r6r6-pgpr
    - GHSA-5p2g-fcmc-qvqq
`

// pnpmAllowBuilds allows the dependency install scripts this stack needs.
// pnpm 11 turned an ignored build script into a hard ERR_PNPM_IGNORED_BUILDS
// failure, so esbuild — which fetches its platform binary in postinstall — has
// to be listed or `pnpm install` exits non-zero on a freshly generated project.
// pnpm 11 renamed the setting to allowBuilds and no longer reads the "pnpm"
// field in package.json at all; the pnpm 10 spelling is kept alongside it so
// the project installs on either major.
func pnpmAllowBuilds() string {
	return `allowBuilds:
  esbuild: true

# pnpm 10 spelling of the same setting.
onlyBuiltDependencies:
  - esbuild
`
}

// rootNpmrc pins pnpm to a flat (hoisted) node_modules. Next.js Turbopack
// can't resolve packages that live only in a nested pnpm dependency — most
// visibly @tiptap/starter-kit's transitive extension packages (the rich-text
// editor), which makes `next build` fail with "Can't resolve
// '@tiptap/extension-horizontal-rule'". A hoisted layout puts every dependency
// where Turbopack (and any other tool that assumes an npm-style tree) can find
// it. Vite and the Go tooling are unaffected.
func rootNpmrc() string {
	return `# Flat node_modules so Next.js Turbopack can resolve every dependency,
# including transitive ones (e.g. tiptap starter-kit's extension packages).
# pnpm's default nested layout hides transitive deps from Turbopack's resolver.
node-linker=hoisted
`
}

func turboJSON() string {
	return `{
  "$schema": "https://turbo.build/schema.json",
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": [".next/**", "!.next/cache/**", "dist/**"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    },
    "lint": {},
    "type-check": {
      "dependsOn": ["^build"]
    },
    "test": {
      "dependsOn": ["^build"]
    }
  }
}
`
}

func rootPackageJSON(opts Options) string {
	// dev runs the frontend dev servers via pnpm's own parallel runner rather
	// than `turbo dev`. turbo ships a platform-specific native binary that can
	// fail to load on Windows (exit 0xC0000135 / STATUS_DLL_NOT_FOUND when the
	// VC++ runtime is missing), which would otherwise block `grit start` on a
	// fresh machine. pnpm needs no native binary, so `dev` always works; turbo
	// is kept for build/lint/test where its caching pays off.
	scripts := fmt.Sprintf(`    "dev": "pnpm --parallel --filter \"./apps/*\" --if-present run dev",
    "build": "turbo build",
    "lint": "turbo lint",
    "format": "` + biomeFormatScript + `",
    "type-check": "turbo type-check",
    "dev:api": "cd apps/api && air",`)

	if opts.ShouldIncludeWeb() {
		scripts += `
    "dev:web": "cd apps/web && pnpm dev",`
	}
	if opts.ShouldIncludeAdmin() {
		scripts += `
    "dev:admin": "cd apps/admin && pnpm dev",`
	}
	if opts.ShouldIncludeExpo() {
		scripts += `
    "dev:expo": "cd apps/expo && npx expo start",`
	}
	if opts.ShouldIncludeDesktop() {
		scripts += `
    "dev:desktop": "cd apps/desktop && wails dev",
    "build:desktop": "cd apps/desktop && wails build",`
	}
	if opts.ShouldIncludeDocs() {
		scripts += `
    "dev:docs": "cd apps/docs && pnpm dev",`
	}

	scripts += `
    "docker:up": "docker compose up -d",
    "docker:down": "docker compose down",
    "docker:logs": "docker compose logs -f",
    "test": "turbo test",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui"`

	// react and react-dom MUST resolve to the exact same version — react-dom
	// throws "Incompatible React versions" (React error #527) at mount
	// otherwise, and the app renders a blank screen with no other clue.
	//
	// apps/expo pins react to 19.1.0 (React Native requires an exact match), so
	// React version pinning lives in pnpm-workspace.yaml (see pnpmWorkspace):
	// pnpm 10 ignores the package.json "pnpm.overrides" field, so keeping it
	// here would only print a deprecation warning on every install.
	return fmt.Sprintf(`{
  "name": "%s",
  "private": true,
  "scripts": {
%s
  },
  "devDependencies": {
    `+biomeDevDependency+`,
    "@playwright/test": "^1.48.0",
    "turbo": "^2.0.0"
  },
  `+packageManagerNew+`
}
`, opts.ProjectName, scripts)
}

func gritConfig(opts Options) string {
	style := opts.Style
	if style == "" {
		style = "default"
	}
	config := fmt.Sprintf(`// Grit Framework Configuration
export default {
  name: "%s",
  style: "%s",
  api: {
    port: 8080,
    prefix: "/api",
  },`, opts.ProjectName, style)

	if opts.ShouldIncludeWeb() {
		config += `
  web: {
    port: 3000,
  },`
	}
	if opts.ShouldIncludeAdmin() {
		config += `
  admin: {
    port: 3001,
  },`
	}
	if opts.ShouldIncludeExpo() {
		config += `
  expo: {
    scheme: "` + opts.ProjectName + `",
  },`
	}
	if opts.ShouldIncludeDesktop() {
		config += `
  desktop: {
    framework: "wails",
    port: 5174,
  },`
	}
	if opts.ShouldIncludeDocs() {
		config += `
  docs: {
    port: 3002,
  },`
	}

	config += `
};
`
	return config
}

func readmeFile(opts Options) string {
	return fmt.Sprintf(`# %s

Built with [Grit](https://gritframework.dev) — Go + React. Built with Grit.

## Quick Start

`+"```bash"+`
# 1. Install Air for Go hot reloading
go install github.com/air-verse/air@latest

# 2. Start infrastructure (PostgreSQL, Redis, MinIO, Mailhog)
docker compose up -d

# 3. Install frontend dependencies
pnpm install

# 4. Start all services (API auto-reloads on file changes)
pnpm dev
`+"```"+`

## Cloned this project?

`+"`.env`"+` holds this machine's secrets and is not committed; `+"`.env.example`"+` is, with
`+"`CHANGE_ME`"+` in their place. Create yours before anything else:

`+"```bash"+`
grit env        # .env from .env.example, every secret freshly generated
grit migrate
grit seed
`+"```"+`

To share a database that already holds encrypted data, put the team's
`+"`FIELD_ENCRYPTION_KEY`"+` in `+"`.env`"+` instead of the new one.

## Project Structure

`+"```"+`
%s/
├── apps/
│   ├── api/          # Go backend (Gin + GORM)
│   ├── web/          # Next.js frontend
│   └── admin/        # Next.js admin panel
├── packages/
│   └── shared/       # Shared types, schemas, constants
├── docker-compose.yml
└── turbo.json
`+"```"+`

## Services

| Service       | URL                          |
|---------------|------------------------------|
| API           | http://localhost:8080         |
| GORM Studio   | http://localhost:8080/studio  |
| Web App       | http://localhost:3000         |
| Admin Panel   | http://localhost:3001         |
| PostgreSQL    | localhost:5434               |
| Redis         | localhost:6380               |
| MinIO Console | http://localhost:9003         |
| Mailhog       | http://localhost:8025         |

## Development

`+"```bash"+`
# Run Go API with hot reload
cd apps/api && air

# Run Next.js web app
cd apps/web && pnpm dev

# Run admin panel
cd apps/admin && pnpm dev

# Run all services via Turborepo
pnpm dev
`+"```"+`

## No Docker? No Problem

If you can't run Docker, use cloud services instead:

`+"```bash"+`
cp .env.cloud.example .env
`+"```"+`

Then fill in your keys for:
- **[Neon](https://neon.tech)** — PostgreSQL (free tier)
- **[Upstash](https://upstash.com)** — Redis (free tier)
- **[Cloudflare R2](https://dash.cloudflare.com)** — File storage (free tier)
- **[Resend](https://resend.com)** — Email (free tier)

No Docker needed — just your API keys and `+"``"+`go run`+"``"+`.

## Tech Stack

- **Backend:** Go + Gin + GORM
- **Frontend:** Next.js 14+ (App Router) + React + TypeScript
- **Styling:** Tailwind CSS + shadcn/ui
- **Database:** PostgreSQL
- **Cache:** Redis
- **Monorepo:** Turborepo + pnpm
- **Validation:** Zod (shared schemas)
- **Data Fetching:** React Query (TanStack Query)

---

*Built with Grit v%s*
`, opts.ProjectName, opts.ProjectName, opts.Version)
}

func gritJSON(opts Options) string {
	return fmt.Sprintf(`{
  "architecture": "%s",
  "frontend": "%s",
  "version": "%s",
  "apps": {
    "expo": %t,
    "desktop": %t,
    "docs": %t
  }
}
`, string(opts.Architecture), string(opts.Frontend), opts.Version,
		opts.ShouldIncludeExpo(), opts.ShouldIncludeDesktop(), opts.ShouldIncludeDocs())
}
