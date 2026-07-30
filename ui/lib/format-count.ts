/**
 * Formats an install count for display: 1234 -> "1.2k".
 *
 * Rules chosen so the number never misleads:
 *   - Under 1000 shows exactly, because "0.9k" is worse than "912".
 *   - Thousands keep one decimal, dropping a trailing ".0" so 2000 is "2k"
 *     rather than "2.0k".
 *   - Rounds DOWN, not to nearest. 1999 showing as "2k" claims installs that
 *     did not happen; "1.9k" is honest and nobody minds the difference.
 */
export function formatCount(n: number): string {
  if (!Number.isFinite(n) || n < 0) return '0'

  if (n < 1_000) return String(Math.floor(n))

  if (n < 1_000_000) {
    const k = Math.floor(n / 100) / 10
    return `${trimZero(k)}k`
  }

  const m = Math.floor(n / 100_000) / 10
  return `${trimZero(m)}m`
}

function trimZero(v: number): string {
  return Number.isInteger(v) ? String(v) : v.toFixed(1)
}

/** "1 install" / "2 installs" — a count of one should not read as plural. */
export function installLabel(n: number): string {
  return n === 1 ? 'install' : 'installs'
}
