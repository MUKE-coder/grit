import { PAID_BLUEPRINTS } from '@/config/blueprints'

/*
 * The pure half of blueprint access: which repository, and whose account.
 *
 * Separated from the webhook route so it can be read and exercised without a
 * Stripe event, and because a route file may only export the handlers Next
 * knows about.
 */

export interface Target {
  owner: string
  repo: string
}

/** github.com/MUKE-coder/grit-instagram -> { owner, repo } */
export function parseRepo(url: string | undefined): Target | null {
  if (!url) return null
  const m = url.match(/github\.com\/([^/]+)\/([^/#?]+)/)
  if (!m) return null
  return { owner: m[1], repo: m[2].replace(/\.git$/, '') }
}

/** Every repository a membership covers. */
export function membershipTargets(): Target[] {
  return PAID_BLUEPRINTS.map((b) => parseRepo(b.repo)).filter(
    (t): t is Target => t !== null
  )
}

/**
 * A GitHub username, or null.
 *
 * Checkout collects it as a custom field, which means the buyer typed it.
 * People paste a profile URL, an @ handle or a trailing slash, so the obvious
 * shapes are unwrapped rather than rejected: a wrong username here costs a
 * support email, and the point of the webhook is that there is not one.
 *
 * The pattern is GitHub's own rule: 1 to 39 characters, alphanumeric or
 * single hyphens, never starting or ending with a hyphen.
 */
export function cleanUsername(raw: string | null | undefined): string | null {
  if (!raw) return null
  let value = raw.trim()
  value = value.replace(/^https?:\/\/(www\.)?github\.com\//i, '')
  value = value.replace(/^@/, '')
  value = value.split(/[/?#\s]/)[0].trim()
  if (!/^[a-z\d](?:[a-z\d]|-(?=[a-z\d])){0,38}$/i.test(value)) return null
  return value
}

/**
 * True when a failed GitHub call is worth Stripe retrying.
 *
 * A 404 means the username does not exist, and it will not exist any harder
 * on the fourth attempt: that one is logged and acknowledged instead, so the
 * event stops rather than retrying for three days.
 */
export function worthRetrying(status: number): boolean {
  return status === 429 || status >= 500
}
