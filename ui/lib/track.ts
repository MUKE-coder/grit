'use client'

import { track as zenithTrack } from 'zenith-analytics/client'

/** What a copy action was for. Kept narrow so the dashboard stays readable. */
export type CopyKind = 'command' | 'code'

/**
 * Records that someone took a block.
 *
 * This is the closest thing to a download count the site can have. A block is
 * never fetched from a server when someone copies it — the source is already in
 * the page — so the only moment we can observe intent is the copy itself.
 *
 * Two kinds are tracked separately on purpose. Copying the CLI command means
 * "install this properly"; copying the source means "paste it into a file". They
 * are different intentions and the ratio is worth knowing.
 *
 * Fires into the existing self-hosted Zenith rather than a database of our own.
 * A public counter endpoint would need a table, a migration, rate limiting and
 * abuse handling to produce a number the analytics install already gives us —
 * and would still need a UI to read it, which Zenith has.
 *
 * Never throws: zenithTrack queues before the snippet loads and swallows its own
 * errors, and this wrapper adds nothing that can fail. A lost metric must never
 * cost someone the copy they came for.
 */
export function trackBlockCopy(input: {
  block: string
  category: string
  subcategory: string
  kind: CopyKind
}) {
  zenithTrack('block_copy', {
    block: input.block,
    category: input.category,
    subcategory: input.subcategory,
    kind: input.kind,
  })
}
