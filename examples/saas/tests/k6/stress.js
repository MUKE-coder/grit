// Stress — how does it behave well beyond normal?
// At 4× normal load you're learning the failure mode, not just the
// limit. Loose thresholds so the run completes even when the system
// degrades under load.
import { userJourney, aggressiveThresholds } from './lib/common.js'

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // normal
    { duration: '5m', target: 400 },  // 4× normal — the stress
    { duration: '2m', target: 0 },    // ramp down
  ],
  thresholds: aggressiveThresholds,
}

export default userJourney
