/*
 * Bun.serve handles requests on one thread per process. The container has 4
 * CPUs, so a single process would leave three of them idle while Go schedules
 * goroutines across all four and PHP runs 32 fpm workers.
 *
 * This forks one server process per CPU. They share the listening socket via
 * reusePort in server.ts, so it is invisible to the benchmark — same URL, same
 * routes, same responses — and it is what Bun's own docs recommend for
 * production.
 */

const workers = Number(process.env.BUN_WORKERS || navigator.hardwareConcurrency || 1)

if (workers > 1 && !process.env.BUN_IS_CHILD) {
  console.log(`bun-bench: forking ${workers} workers`)

  const children = Array.from({ length: workers }, () =>
    Bun.spawn(['bun', 'run', 'server.ts'], {
      env: { ...process.env, BUN_IS_CHILD: '1' },
      stdout: 'inherit',
      stderr: 'inherit',
    }),
  )

  const shutdown = () => {
    for (const c of children) c.kill()
    process.exit(0)
  }
  process.on('SIGTERM', shutdown)
  process.on('SIGINT', shutdown)

  await Promise.all(children.map((c) => c.exited))
} else {
  await import('./server.ts')
}
