// Soak — does the system leak over hours? A flat latency line that
// slowly tilts upward is the classic memory-leak signature. Run before
// every major release; watch /pulse/ui in another tab.
//
// NOTE: this takes ~4 hours. Schedule it overnight or on a CI runner.
import { userJourney, defaultThresholds } from './lib/common.js'

export const options = {
  stages: [
    { duration: '5m', target: 100 },  // ramp up
    { duration: '4h', target: 100 },  // long hold — the soak
    { duration: '5m', target: 0 },    // ramp down
  ],
  thresholds: defaultThresholds,
}

export default userJourney
