'use client'

import { useState } from 'react'
import type React from 'react'
import { Scale, Building2, ClipboardCheck, FileSearch, Check } from 'lucide-react'
import { CodeBlock } from '@/components/code-block'

/**
 * The compliance surface: GDPR, SSO, access reviews, audit trail.
 *
 * This is the section a buyer's security reviewer reads, so every claim here
 * is one the code can back:
 *
 *  - "tamper-evident" means an actual SHA-256 hash chain with a verify
 *    endpoint, not "we log it and nobody can edit the table".
 *  - Erase is described as what it does — hard-delete children, anonymise the
 *    user row — rather than the softer "removes their data", because a
 *    reviewer will ask which one and the honest answer is both.
 *  - SSO connections are database rows, so "per customer" is literal.
 *
 * If any of that changes, change it here first and let the copy follow.
 */

interface Topic {
  key: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  headline: string
  blurb: string
  points: string[]
  image?: string
  imageLabel?: string
  code?: string
  codeFile?: string
  language?: string
}

const TOPICS: Topic[] = [
  {
    key: 'gdpr',
    label: 'GDPR',
    icon: Scale,
    headline: 'Article 15 and Article 17, as buttons',
    blurb:
      'Right-to-access exports and right-to-erasure, with every erasure written to a hash-chained journal that can prove it has not been edited since.',
    points: [
      'Export returns one JSON file: profile, uploads, sessions, activity, 2FA state',
      'Erase hard-deletes personal records across nine tables and anonymises the account',
      'The user id survives as a tombstone, so the audit trail stays readable',
      'The journal stores no PII — just ids, counts, actor and a SHA-256 chain',
      'A verify pass replays the chain and names the row where it breaks',
    ],
    image: '/images/system/gdpr.png',
    imageLabel: 'localhost:3001/system/gdpr',
  },
  {
    key: 'sso',
    label: 'Enterprise SSO',
    icon: Building2,
    headline: 'One connection per customer, added from the admin',
    blurb:
      'OIDC or SAML 2.0, stored as database rows rather than environment variables — so onboarding an enterprise customer is a form, not a deploy.',
    points: [
      'Routed by email domain: acme.com signs in through Acme’s identity provider',
      'OIDC via discovery, or SAML with a metadata URL or pasted XML',
      'Client secrets and the SAML private key are encrypted at rest (AES-256-GCM)',
      'Just-in-time provisioning, with IdP groups mapped to roles',
      'A copy button hands the ACS and metadata URLs to the customer’s IdP admin',
    ],
    image: '/images/system/sso.png',
    imageLabel: 'localhost:3001/system/sso',
  },
  {
    key: 'reviews',
    label: 'Access reviews',
    icon: ClipboardCheck,
    headline: 'The recertification campaign SOC 2 asks for',
    blurb:
      'Open a review and it snapshots every current role assignment. A reviewer keeps or revokes each one; revocations take effect immediately and land in the audit log.',
    points: [
      'Targets SOC 2 CC6.2 / CC6.3 and ISO 27001 A.9.2.5',
      'Each item copies the user’s email and role at open time, so the record survives deletion',
      'Revoking removes the grant there and then — not a ticket for someone else',
      'Completed reviews cannot be reopened, which is what makes them evidence',
    ],
    code: `// A campaign is a snapshot plus a decision per row.

POST   /api/v1/access-reviews                       // open — snapshots grants
GET    /api/v1/access-reviews/:id                   // the items to decide
POST   /api/v1/access-reviews/:id/items/:itemId/decision
POST   /api/v1/access-reviews/:id/complete          // freeze it

// The decision body is deliberately small:
// { "decision": "keep" }     — no change
// { "decision": "revoke" }   — the grant is removed now`,
    codeFile: 'internal/routes/routes.go',
    language: 'go',
  },
  {
    key: 'audit',
    label: 'Audit trail',
    icon: FileSearch,
    headline: 'Two logs, because they answer different questions',
    blurb:
      'A hash-chained record of every authenticated mutation for evidence, and a readable event timeline for operators. Both export to a SIEM in OCSF.',
    points: [
      'Every authenticated POST/PUT/PATCH/DELETE is chained: prev hash, this hash',
      'Bodies are stored as a SHA-256 digest — evidence of what was sent, without the PII',
      'An integrity endpoint replays the chain and names the first row that fails',
      'Written by a single goroutine off a buffered channel, so the chain cannot race',
      'GET /audit/ocsf streams newline-delimited OCSF 1.3.0 for Splunk, Sentinel or Panther',
    ],
    code: `// Prove the chain has not been edited
GET /api/v1/admin/activity/integrity

{ "valid": true, "total_entries": 12345 }

// Or, when a row was tampered with:
{
  "valid": false,
  "broken_at": 47,
  "broken_at_id": "019fb4d9-…",
  "message": "hash mismatch — row was modified after it was written"
}`,
    codeFile: 'internal/handlers/activity.go',
    language: 'go',
  },
]

export function ComplianceShowcase() {
  const [active, setActive] = useState(TOPICS[0].key)
  const topic = TOPICS.find((t) => t.key === active) ?? TOPICS[0]

  return (
    <div>
      <div role="tablist" aria-label="Compliance topic" className="flex flex-wrap gap-2 mb-8">
        {TOPICS.map((t) => {
          const selected = t.key === active
          return (
            <button
              key={t.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(t.key)}
              className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              <t.icon className="h-4 w-4" />
              {t.label}
            </button>
          )
        })}
      </div>

      <div className="grid lg:grid-cols-[1fr_20rem] gap-8 lg:gap-10 items-start">
        <div>
          {topic.image ? (
            <div className="rounded-xl overflow-hidden border border-border bg-card/40 shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
              <div className="flex items-center gap-2 px-3.5 py-2.5 bg-card/70 border-b border-border/60">
                <div className="flex items-center gap-1.5">
                  <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                  <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                  <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                </div>
                <span className="mx-auto text-[11px] font-mono text-muted-foreground">
                  {topic.imageLabel}
                </span>
              </div>
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                key={topic.key}
                src={topic.image}
                alt={`${topic.label}: ${topic.headline}`}
                className="w-full h-auto block"
                loading="lazy"
              />
            </div>
          ) : (
            <div className="rounded-xl overflow-hidden border border-border bg-card/40">
              <div className="flex items-center gap-2 px-4 py-2.5 bg-card/60 border-b border-border/60">
                <span className="text-[11.5px] font-mono text-muted-foreground">
                  {topic.codeFile}
                </span>
              </div>
              <CodeBlock
                key={topic.key}
                code={topic.code ?? ''}
                language={topic.language ?? 'go'}
                className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
              />
            </div>
          )}
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-3 leading-snug">{topic.headline}</h3>
          <p className="text-sm text-muted-foreground leading-relaxed mb-5">{topic.blurb}</p>
          <ul className="space-y-2.5">
            {topic.points.map((p) => (
              <li key={p} className="flex gap-2 text-[12.5px] text-foreground/80 leading-relaxed">
                <Check className="h-3.5 w-3.5 shrink-0 mt-0.5 text-primary" strokeWidth={2.5} />
                {p}
              </li>
            ))}
          </ul>
        </div>
      </div>

      <p className="text-[11.5px] text-muted-foreground/70 mt-8 pt-5 border-t border-border/40 leading-relaxed max-w-3xl">
        What Grit does not ship: a cookie-consent banner, a signed SBOM, or data-residency
        controls. Field-level encryption, a go-live checklist and dependency scanning in CI are
        included &mdash; but the paperwork side of a certification is still yours.
      </p>
    </div>
  )
}
