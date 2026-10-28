import Image from 'next/image'
import Link from 'next/link'
import { ArrowRight } from 'lucide-react'

/**
 * The bento: the batteries, each shown rather than listed.
 *
 * A list of forty feature names is read as marketing even when every line is
 * true. A card that draws the thing (the six-digit field, the permission grid,
 * the request pipeline) is checkable at a glance, and the ones that cannot be
 * drawn honestly are left as a number and a link instead of a fake screenshot.
 *
 * Column spans come straight from the design: 7·5 / 4·4·4 / 6·6 / 5·3·4 / 12.
 * They collapse to halves under 980px and to one column under 640px, which is
 * what `g-bento-cell` carries.
 */

function Eyebrow({ children }: { children: React.ReactNode }) {
  return <p className="g-eyebrow mb-3">{children}</p>
}

function Title({ children }: { children: React.ReactNode }) {
  return (
    <h3 className="m-0 text-[20px] font-semibold leading-[1.25] tracking-[-0.02em] text-fg">
      {children}
    </h3>
  )
}

function Body({ children }: { children: React.ReactNode }) {
  return <p className="mt-3 text-[14px] leading-[1.65] text-fg-muted">{children}</p>
}

/** A row of mono key/value facts, the shape the design uses under a card. */
function Facts({ rows }: { rows: [string, string][] }) {
  return (
    <dl className="m-0 space-y-1.5 font-mono text-[12.5px] leading-[18px]">
      {rows.map(([value, label]) => (
        <div key={label} className="flex gap-2">
          <dt className="text-fg">{value}</dt>
          <dd className="m-0 text-fg-subtle">{label}</dd>
        </div>
      ))}
    </dl>
  )
}

const PERMISSIONS: [string, boolean[]][] = [
  ['admin', [true, true, true, true]],
  ['editor', [true, true, true, false]],
  ['viewer', [true, false, false, false]],
]

function Tick({ on }: { on: boolean }) {
  return (
    <span
      aria-hidden="true"
      className={`inline-grid h-[18px] w-[18px] place-items-center rounded-xs border ${
        on ? 'border-accent bg-accent text-accent-fg' : 'border-border-strong'
      }`}
    >
      {on && (
        <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="2.4" strokeLinecap="round" strokeLinejoin="round" className="h-2.5 w-2.5">
          <path d="M13 4.5 6.5 11 3 7.5" />
        </svg>
      )}
    </span>
  )
}

/** The request pipeline across the bottom of the section. */
const PIPELINE = ['request', 'WAF', 'rate limit', 'CSRF', 'auth', 'RBAC', 'handler', 'GORM', 'audit log']

export function Batteries() {
  return (
    <div className="mt-12 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-12">
      {/* ── Identity ─────────────────────────────────────────────── */}
      <div className="g-bento lg:col-span-7">
        <Image
          src="/images/3d/key.webp"
          alt=""
          width={160}
          height={160}
          className="pointer-events-none absolute -right-2 top-4 h-[150px] w-[150px] select-none opacity-90"
        />
        <div className="relative max-w-[26rem]">
          <Eyebrow>Identity</Eyebrow>
          <Title>Authentication, finished.</Title>
          <Body>
            Access and refresh tokens in HttpOnly cookies, OAuth, TOTP two-factor,
            passkeys and sessions you can revoke.
          </Body>
        </div>
        <div className="relative mt-auto flex flex-wrap items-end gap-6 pt-10">
          <div className="rounded-xl border border-border bg-bg p-4">
            <p className="mb-2.5 text-[13px] text-fg-muted">Enter your 6-digit code</p>
            <div className="flex gap-2" aria-hidden="true">
              {['4', '8', '1', '0', '2', '9'].map((d, i) => (
                <span
                  key={i}
                  className={`grid h-9 w-9 place-items-center rounded-lg border font-mono text-[17px] text-fg ${
                    i === 5 ? 'border-accent' : 'border-border-strong'
                  }`}
                >
                  {d}
                </span>
              ))}
            </div>
          </div>
          <Facts
            rows={[
              ['15 min', 'access token'],
              ['168 h', 'refresh token'],
              ['5 / 15 min', 'login rate limit'],
            ]}
          />
        </div>
      </div>

      {/* ── Roles ────────────────────────────────────────────────── */}
      <div className="g-bento lg:col-span-5">
        <Eyebrow>Roles &amp; permissions</Eyebrow>
        <Title>Roles are rows, not code.</Title>
        <Body>Every generated resource registers its own permissions.</Body>
        <table className="mt-auto w-full table-fixed border-separate border-spacing-0 pt-10">
          <caption className="sr-only">
            Which actions each role has on the Products resource
          </caption>
          <colgroup>
            <col className="w-[34%]" />
          </colgroup>
          <thead>
            <tr className="g-eyebrow">
              <th scope="col" className="pb-2 text-left font-medium">Products</th>
              {['View', 'Create', 'Edit', 'Delete'].map((a) => (
                <th key={a} scope="col" className="pb-2 text-center font-medium">
                  {a}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {PERMISSIONS.map(([role, cells]) => (
              <tr key={role}>
                <th scope="row" className="py-2 text-left text-[13px] font-normal text-fg">
                  {role}
                </th>
                {cells.map((on, i) => (
                  <td key={i} className="py-2 text-center">
                    <Tick on={on} />
                    <span className="sr-only">{on ? 'yes' : 'no'}</span>
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* ── Jobs · Storage · Realtime ────────────────────────────── */}
      <div className="g-bento lg:col-span-4">
        <Eyebrow>Jobs &amp; cron</Eyebrow>
        <Title>Work that happens later.</Title>
        <Body>Background jobs on asynq and Redis, cron schedules, durable workflows.</Body>
        <div className="mt-auto space-y-2 pt-10 font-mono text-[12.5px]">
          <p className="rounded-lg border border-border bg-bg px-3 py-2 text-fg">
            send_welcome_email
          </p>
          <p className="flex items-center justify-between rounded-lg border border-border bg-bg px-3 py-2 text-fg">
            0 3 * * * nightly_backup
            <span className="text-success" aria-hidden="true">✓</span>
          </p>
          <p className="text-fg-subtle">asynq · retries · dead letter</p>
        </div>
      </div>

      <div className="g-bento lg:col-span-4">
        <Eyebrow>File storage</Eyebrow>
        <Title>Uploads that know what they accept.</Title>
        <Body>Presigned URLs, server-enforced accept lists.</Body>
        <div className="mt-auto pt-10">
          <div className="rounded-lg border border-border bg-bg px-3 py-2.5">
            <p className="flex items-center justify-between font-mono text-[12.5px] text-fg">
              contact-photo.webp <span className="text-fg-subtle">2.4 MB</span>
            </p>
            <span aria-hidden="true" className="mt-2 block h-1 rounded-full bg-surface-3">
              <span className="block h-full w-3/4 rounded-full bg-accent" />
            </span>
          </div>
          <p className="mt-2.5 font-mono text-[12.5px] text-fg-subtle">S3 · R2 · MinIO</p>
        </div>
      </div>

      <div className="g-bento lg:col-span-4">
        <Eyebrow>Realtime</Eyebrow>
        <Title>Live, on every client.</Title>
        <Body>A realtime hub that pushes updates to web, admin, mobile and desktop.</Body>
        <div className="mt-auto pt-10">
          <p className="flex items-center gap-2.5 font-mono text-[12.5px] text-fg-muted">
            <span aria-hidden="true" className="flex">
              {['AN', 'MK', 'JO'].map((who, i) => (
                <span
                  key={who}
                  className="-ml-1.5 grid h-7 w-7 place-items-center rounded-full border border-bg bg-surface-3 text-[10px] text-fg first:ml-0"
                  style={{ zIndex: 3 - i }}
                >
                  {who}
                </span>
              ))}
            </span>
            <span className="g-ping h-[7px] w-[7px] rounded-full bg-success" aria-hidden="true" />
            3 online
          </p>
          <p className="mt-2.5 rounded-lg border border-border bg-bg px-3 py-2 font-mono text-[12.5px] leading-[18px] text-fg-muted">
            <span className="text-code-fn">contact.updated</span> · Ada edited
            &ldquo;Kampala office&rdquo;
          </p>
        </div>
      </div>

      {/* ── AI gateway · Pulse ───────────────────────────────────── */}
      <div className="g-bento lg:col-span-6">
        <Eyebrow>AI gateway</Eyebrow>
        <Title>
          One gateway. <span className="g-hl">100+ models.</span>
        </Title>
        <Body>With the same auth, rate limits and logging as the rest of your API.</Body>
        <p className="mt-auto flex flex-wrap items-center gap-3 pt-10 font-mono text-[12.5px]">
          <span className="rounded-lg border border-border bg-bg px-3 py-2 text-fg">
            POST /api/ai/chat
          </span>
          <span aria-hidden="true" className="h-px w-8 bg-border-strong" />
          <span className="rounded-lg border border-accent-line bg-accent-soft px-3 py-2 text-fg">
            Gateway
          </span>
          <span aria-hidden="true" className="h-px w-8 bg-border-strong" />
          <span className="text-fg-subtle">Claude · GPT · Gemini · local</span>
        </p>
      </div>

      <div className="g-bento lg:col-span-6">
        <div className="flex items-start justify-between gap-4">
          <div>
            <Eyebrow>Pulse observability</Eyebrow>
            <Title>Know what every route is doing.</Title>
          </div>
          <span className="g-badge font-mono text-[11px] text-fg-subtle">/pulse/ui</span>
        </div>
        <dl className="mt-6 flex gap-8">
          {[
            ['124.2K', 'requests'],
            ['98.6%', '2xx'],
            ['125 ms', 'fastest'],
          ].map(([value, label]) => (
            <div key={label}>
              <dt className="text-[26px] font-medium leading-none tracking-[-0.02em] text-fg">
                {value}
              </dt>
              <dd className="m-0 mt-1.5 font-mono text-[12px] text-fg-subtle">{label}</dd>
            </div>
          ))}
        </dl>
        <svg
          aria-hidden="true"
          viewBox="0 0 520 110"
          preserveAspectRatio="none"
          className="mt-auto h-[110px] w-full pt-6"
        >
          <defs>
            <linearGradient id="pulse-fill" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stopColor="var(--accent)" stopOpacity="0.22" />
              <stop offset="100%" stopColor="var(--accent)" stopOpacity="0" />
            </linearGradient>
          </defs>
          <path
            d="M0 92 L40 86 L80 94 L120 72 L160 80 L200 58 L240 66 L280 44 L320 52 L360 34 L400 42 L440 24 L480 30 L520 14 L520 110 L0 110 Z"
            fill="url(#pulse-fill)"
          />
          <path
            d="M0 92 L40 86 L80 94 L120 72 L160 80 L200 58 L240 66 L280 44 L320 52 L360 34 L400 42 L440 24 L480 30 L520 14"
            fill="none"
            stroke="var(--accent)"
            strokeWidth="1.5"
          />
        </svg>
      </div>

      {/* ── Admin · Mail · the rest ──────────────────────────────── */}
      <div className="g-bento lg:col-span-5">
        <Eyebrow>Admin panel</Eyebrow>
        <Title>A back office per resource.</Title>
        <Body>Filters, export, grid entry, stats, trash and restore. Nine themes.</Body>
        <table className="mt-auto w-full border-separate border-spacing-0 pt-10 text-left">
          <caption className="sr-only">A generated admin table</caption>
          <thead>
            <tr className="g-eyebrow">
              <th scope="col" className="pb-2 font-medium">Contact</th>
              <th scope="col" className="pb-2 font-medium">Group</th>
              <th scope="col" className="pb-2 font-medium">Status</th>
            </tr>
          </thead>
          <tbody className="text-[13px]">
            {[
              ['Ada Nakato', 'Kampala', 'Active'],
              ['Musa Okello', 'Nairobi', 'Invited'],
              ['Jane Atim', 'Remote', 'Active'],
            ].map(([name, group, status]) => (
              <tr key={name}>
                <td className="border-t border-border py-2 text-fg">{name}</td>
                <td className="border-t border-border py-2 text-fg-muted">{group}</td>
                <td className="border-t border-border py-2">
                  <span
                    className={`rounded-md px-2 py-0.5 text-[11px] ${
                      status === 'Active'
                        ? 'bg-success/15 text-success'
                        : 'bg-warning/15 text-warning'
                    }`}
                  >
                    {status}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="g-bento lg:col-span-3">
        <Eyebrow>Mail &amp; payments</Eyebrow>
        <Title>Email and billing, wired.</Title>
        <ul className="mt-auto list-none space-y-2 p-0 pt-10 text-[13px] text-fg">
          {['Resend', 'Stripe plugin', 'Webhooks'].map((name) => (
            <li key={name} className="rounded-lg border border-border bg-bg px-3 py-2">
              {name}
            </li>
          ))}
        </ul>
      </div>

      <Link href="/docs/batteries" className="g-bento lg:col-span-4">
        <Image
          src="/images/3d/battery.webp"
          alt=""
          width={160}
          height={160}
          className="pointer-events-none mx-auto h-[140px] w-[140px] select-none"
        />
        <p className="mt-auto text-[40px] font-semibold leading-none tracking-[-0.03em] text-fg">
          50+
        </p>
        <p className="mt-3 text-[14px] leading-[1.65] text-fg-muted">
          batteries: feature flags, optimistic locking, i18n, multi-tenant, push,
          offline-first desktop…
        </p>
        <span className="g-card__more mt-4">
          See them all <ArrowRight className="h-3.5 w-3.5" aria-hidden="true" />
        </span>
      </Link>

      {/* ── Security, full width ─────────────────────────────────── */}
      <div className="g-bento sm:col-span-2 lg:col-span-12">
        <div className="flex flex-wrap items-end justify-between gap-6">
          <div className="max-w-[34rem]">
            <Eyebrow>Security &amp; compliance</Eyebrow>
            <Title>Hardened by default.</Title>
            <Body>
              OWASP Top 10:2025, a WAF, a hash-chained audit log, and GDPR export and
              erasure.
            </Body>
          </div>
          <Image
            src="/images/3d/shield.webp"
            alt=""
            width={160}
            height={160}
            className="pointer-events-none h-[120px] w-[120px] select-none"
          />
        </div>

        <p className="g-eyebrow mt-10">Every request, on its way in</p>
        <ol className="mt-4 flex list-none flex-wrap items-center gap-2 p-0">
          {PIPELINE.map((stage, i) => (
            <li key={stage} className="flex items-center gap-2">
              <span
                className={`rounded-lg border px-3 py-1.5 font-mono text-[12.5px] ${
                  i === 0 || i === PIPELINE.length - 1
                    ? 'border-border text-fg-subtle'
                    : 'border-border-strong text-fg'
                }`}
              >
                {stage}
              </span>
              {i < PIPELINE.length - 1 && (
                <span aria-hidden="true" className="g-march h-0.5 w-5" />
              )}
            </li>
          ))}
        </ol>
      </div>
    </div>
  )
}
