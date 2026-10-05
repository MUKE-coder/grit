// Spike — can it survive a sudden, massive surge (flash sale, viral
// post)? The KEY question is RECOVERY: does p95 return to baseline
// after the surge passes?
import { userJourney, aggressiveThresholds } from './lib/common.js'

export const options = {
  stages: [
    { duration: '1m',  target: 50 },    // baseline
    { duration: '10s', target: 1000 },  // SLAM
    { duration: '1m',  target: 1000 },  // hold the spike
    { duration: '10s', target: 50 },    // drop back
    { duration: '1m',  target: 50 },    // recovery window — watch p95
    { duration: '10s', target: 0 },
  ],
  thresholds: aggressiveThresholds,
}

export default userJourney
