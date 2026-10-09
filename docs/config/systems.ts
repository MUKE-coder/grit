import type { SystemDesign } from './systems-types'
import { AUTHENTICATION } from './systems-identity'
import { MFA, PASSKEYS, SESSIONS } from './systems-identity-2'
import { AUTHORIZATION, DATA_ISOLATION } from './systems-access'
import { MULTITENANCY, RATE_LIMITING } from './systems-access-2'
import { CACHING, CONCURRENCY } from './systems-data'
import { ENCRYPTION, MIGRATIONS, PAGINATION } from './systems-data-2'
import { BACKGROUND_JOBS, OUTBOX, SCHEDULING } from './systems-delivery'
import { EMAIL, REALTIME, WEBHOOKS } from './systems-delivery-2'
import { AUDIT, FILE_STORAGE, IMPORT_EXPORT, OBSERVABILITY } from './systems-operations'
import { API_CONTRACT, CODE_GENERATION, FEATURE_FLAGS } from './systems-platform'

/**
 * The registry. Order here is the order the index and the sidebar show.
 *
 * Grouped by what the system is for rather than alphabetically, because
 * somebody reading these in order should move from "who is calling" through
 * "what may they touch" to "how does the work get done", not from Audit to
 * Webhooks.
 */
export const SYSTEMS: SystemDesign[] = [
  AUTHENTICATION,
  SESSIONS,
  MFA,
  PASSKEYS,
  AUTHORIZATION,
  DATA_ISOLATION,
  MULTITENANCY,
  RATE_LIMITING,
  CACHING,
  CONCURRENCY,
  MIGRATIONS,
  PAGINATION,
  ENCRYPTION,
  BACKGROUND_JOBS,
  SCHEDULING,
  OUTBOX,
  REALTIME,
  WEBHOOKS,
  EMAIL,
  OBSERVABILITY,
  AUDIT,
  FILE_STORAGE,
  IMPORT_EXPORT,
  CODE_GENERATION,
  API_CONTRACT,
  FEATURE_FLAGS,
]

export const SYSTEM_GROUPS = [
  'Identity',
  'Access',
  'Data',
  'Delivery',
  'Operations',
  'Platform',
] as const

/** What each group is for, shown above its cards on the index. */
export const GROUP_BLURB: Record<string, string> = {
  Identity: 'Who the caller is, and how they prove it.',
  Access: 'What that caller is allowed to touch, and what stops them touching the rest.',
  Data: 'Where the application state lives, how it changes, and how it stays correct.',
  Delivery: 'Getting work out of the request path and getting results back to the client.',
  Operations: 'Knowing what the system did, and being able to say so afterwards.',
  Platform: 'The parts that generate, migrate and run everything above.',
}

export function getSystem(slug: string): SystemDesign | undefined {
  return SYSTEMS.find((s) => s.slug === slug)
}

export function systemsInGroup(group: string): SystemDesign[] {
  return SYSTEMS.filter((s) => s.group === group)
}
