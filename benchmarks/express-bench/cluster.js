/*
 * Node is single-threaded; the container has 4 CPUs. Handing Express one core
 * while Go gets four goroutine-scheduled ones would not be a comparison, it
 * would be a strawman — so this forks one worker per CPU, which is what the
 * Node docs and every production Node deployment guide tell you to do.
 *
 * Workers share the listening socket, so this is transparent to the benchmark:
 * same URL, same routes, same responses.
 */

import cluster from 'node:cluster'
import os from 'node:os'

const workers = Number(process.env.NODE_CLUSTER_WORKERS || os.availableParallelism())

if (cluster.isPrimary && workers > 1) {
  console.log(`express-bench: forking ${workers} workers`)
  for (let i = 0; i < workers; i++) cluster.fork()

  // Restart a worker that dies, so a crash mid-run shows up as a blip rather
  // than as steadily falling throughput nobody notices.
  cluster.on('exit', (worker, code, signal) => {
    console.error(`worker ${worker.process.pid} exited (${signal || code}) — restarting`)
    cluster.fork()
  })
} else {
  await import('./server.js')
}
