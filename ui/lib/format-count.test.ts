/**
 * Run with:  npx tsx lib/format-count.test.ts
 *
 * A plain assertion script rather than a test runner: this is one pure function
 * and adding vitest to the site to check it would be more dependency than the
 * thing being tested.
 */
import assert from 'node:assert/strict'
import { formatCount, installLabel } from './format-count'

const cases: [number, string][] = [
  [0, '0'],
  [1, '1'],
  [912, '912'],
  [999, '999'],
  // Exact thousands drop the pointless ".0".
  [1_000, '1k'],
  [2_000, '2k'],
  [1_234, '1.2k'],
  [1_250, '1.2k'],
  // Rounds DOWN: 1999 must not claim 2k installs that did not happen.
  [1_999, '1.9k'],
  [10_500, '10.5k'],
  [999_999, '999.9k'],
  [1_000_000, '1m'],
  [1_450_000, '1.4m'],
  // Nonsense in, harmless out — never NaN on the page.
  [-5, '0'],
  [Number.NaN, '0'],
  [Number.POSITIVE_INFINITY, '0'],
]

for (const [input, expected] of cases) {
  const got = formatCount(input)
  assert.equal(got, expected, `formatCount(${input}) = ${got}, want ${expected}`)
}

assert.equal(installLabel(1), 'install')
assert.equal(installLabel(0), 'installs')
assert.equal(installLabel(2), 'installs')

console.log(`format-count: ${cases.length + 3} assertions passed`)
