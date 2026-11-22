/*
 * The stable channel: the release a team can take without following every day.
 *
 * Grit ships every merged change, several times a day. That is right for the
 * project and wrong for a team with a change-control process, and it has a
 * real cost: a handful of releases have fixed a regression from the days
 * before, so the newest release is occasionally the worst one to be on
 * (grit#245).
 *
 * This file is the answer, and it is written by a machine rather than by
 * hand. scripts/promote-stable.py runs daily after the canary and moves it
 * forward when, and only when, the rule below holds.
 *
 * ── The rule ──────────────────────────────────────────────────────────────
 *
 * stable is the newest release that
 *
 *   1. is at least MIN_AGE_DAYS old,
 *   2. is not a prerelease, and
 *   3. has had every scheduled canary run since it published come back green.
 *
 * ── What that does and does not prove ─────────────────────────────────────
 *
 * Said plainly, because the honest version is less than it sounds. The canary
 * scaffolds from main, not from the candidate tag, so a green week means the
 * scaffolder has been healthy through that week rather than that this exact
 * release was canaried. What it buys is the thing the report asked for: the
 * last week of releases is excluded, and that is where the regressions it
 * named lived, each fixed within days.
 *
 * There are still no backports and no SLA. A fix lands on the edge channel
 * and reaches stable when it has aged. If you need it sooner, take the edge
 * release: that is what it is for.
 */

/** How long a release has to have been out before it can be promoted. */
export const STABLE_MIN_AGE_DAYS = 7

export interface StableChannel {
  /** The version, without a leading v. */
  version: string
  /** The day it was promoted, ISO. */
  promoted: string
  /** The day the release itself published, ISO. */
  released: string
  /** Green scheduled canary runs in its first STABLE_MIN_AGE_DAYS days. */
  canaryRuns: number
}

// Written by scripts/promote-stable.py. Do not edit by hand: the next run
// will overwrite it, and the value is only meaningful because a machine
// checked the rule.
export const STABLE: StableChannel = {
  version: '3.527.0',
  promoted: '2026-11-22',
  released: '2026-11-13',
  canaryRuns: 6,
}
