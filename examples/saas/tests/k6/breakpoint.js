// Breakpoint — slowly ramp until the system actually breaks.
// The VU count at which error-rate or p95 first breaches the threshold
// is your REAL capacity number for launch planning.
import { userJourney } from './lib/common.js'

export const options = {
  // Slow climb to 5000 VUs over 1h. Adjust the ceiling for your gear.
  stages: [
    { duration: '1h', target: 5000 },
  ],
  thresholds: {
    // Abort the test the moment thresholds breach. The VU count when
    // that happens is the breaking point.
    http_req_duration: [{ threshold: 'p(95)<1000', abortOnFail: true }],
    http_req_failed:   [{ threshold: 'rate<0.05',  abortOnFail: true }],
  },
}

export default userJourney
