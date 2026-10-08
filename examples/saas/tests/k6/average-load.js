// Average-load — how does it behave under normal expected traffic?
// Ramps to 100 VUs, holds 5 minutes. The steady-state baseline you
// compare every future test against.
import { userJourney, defaultThresholds } from './lib/common.js'

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // ramp up
    { duration: '5m', target: 100 },  // hold (steady state)
    { duration: '2m', target: 0 },    // ramp down
  ],
  thresholds: defaultThresholds,
}

export default userJourney
