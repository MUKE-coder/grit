/*
 * Next's standalone output is a single-threaded Node server. The container has
 * 4 CPUs, so running one process would leave three idle while Go schedules
 * goroutines across all four and PHP runs 32 fpm workers — a strawman rather
 * than a comparison.
 *
 * This forks one standalone server per CPU. They share the listening socket, so
 * it is invisible to the benchmark: same URL, same routes, same responses.
 */

const cluster = require('node:cluster')
const os = require('node:os')

const workers = Number(process.env.NODE_CLUSTER_WORKERS || os.availableParallelism())

if (cluster.isPrimary && workers > 1) {
  console.log(`nextjs-bench: forking ${workers} workers`)
  for (let i = 0; i < workers; i++) cluster.fork()

  cluster.on('exit', (worker, code, signal) => {
    console.error(`worker ${worker.process.pid} exited (${signal || code}) — restarting`)
    cluster.fork()
  })
} else {
  require('./server.js')
}
