import Image from 'next/image'
import { CopyButton } from '@/components/grit/terminal'

/*
 * What `grit deploy` actually does, as the five things it does.
 *
 * The old version of this was a bullet list beside a terminal, which made the
 * reader match one to the other. Drawing it as a pipeline means the claim
 * ("one command to a server you own") and its mechanism are the same picture.
 *
 * Every stage here is a real step of the command, named after the tool that
 * performs it, because "no platform lock-in" is only credible if you can see
 * which ordinary Unix pieces are involved.
 */

const COMMAND = 'grit deploy --host deploy@server.com --domain myapp.com'

const STAGES: { tool: string; title: string; body: string }[] = [
  { tool: 'local', title: 'Build', body: 'Linux binary, CGO off, plus the frontend' },
  { tool: 'scp', title: 'Upload', body: 'Over SSH to deploy@server.com' },
  { tool: 'systemd', title: 'Run', body: 'A service with auto-restart' },
  { tool: 'caddy', title: 'Serve', body: "Reverse proxy and Let's Encrypt TLS" },
  { tool: 'live', title: 'myapp.com', body: 'HTTPS, health-checked' },
]

export function DeployPipeline() {
  return (
    <div className="mt-12">
      <div className="flex items-center gap-3 rounded-xl border border-border bg-surface-1 py-3 pl-4 pr-2">
        <span aria-hidden="true" className="font-mono text-[14px] text-fg-subtle">
          $
        </span>
        <code className="min-w-0 flex-1 overflow-x-auto whitespace-nowrap font-mono text-[13.5px] text-code-fg">
          {COMMAND}
        </code>
        <CopyButton text={COMMAND} label="Copy the deploy command" />
      </div>

      <ol className="mt-4 grid list-none grid-cols-1 gap-3 p-0 sm:grid-cols-2 lg:grid-cols-5">
        {STAGES.map((stage, i) => (
          <li
            key={stage.tool}
            className={`rounded-xl border bg-surface-1 p-5 ${
              i === STAGES.length - 1 ? 'border-accent-line' : 'border-border'
            }`}
          >
            <p className="g-eyebrow mb-3">{stage.tool}</p>
            <p className="m-0 text-[16px] font-semibold tracking-[-0.01em] text-fg">
              {stage.title}
            </p>
            <p className="mt-1.5 text-[13px] leading-[1.6] text-fg-muted">{stage.body}</p>
          </li>
        ))}
      </ol>

      <div className="mt-6 flex flex-wrap items-center justify-between gap-6 rounded-xl border border-border px-6 py-5">
        <p className="m-0 max-w-[var(--measure)] text-[14px] leading-[1.65] text-fg-muted">
          <span className="text-fg">It is a Go binary and some containers.</span> Which
          means it runs anywhere: a $5 VPS, a managed platform, or your own Docker host.
          Fly, Railway, Render, Dokploy, Coolify, EC2 and Lightsail each have a guide with
          the actual steps, what it costs, and the thing that catches people out.
        </p>
        <Image
          src="/images/3d/orbit.webp"
          alt=""
          width={160}
          height={160}
          className="pointer-events-none h-[110px] w-[110px] select-none"
        />
      </div>
    </div>
  )
}
