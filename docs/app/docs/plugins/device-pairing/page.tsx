import Link from 'next/link'
import { ArrowLeft, ArrowRight, QrCode, ShieldCheck, AlertTriangle, Timer } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'

export const metadata = {
  title: 'Device pairing plugin — Grit',
  description:
    'Sign a browser in by showing a QR code that an already-signed-in device scans and approves. The WhatsApp Web flow.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/device-pairing' },
}

export default function DevicePairingPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <QrCode className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugin
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            Device pairing
          </h1>
          <p className="mb-8 text-lg leading-relaxed text-muted-foreground">
            A browser shows a QR code. A device that is already signed in scans it,
            is shown what it is about to trust, and approves. The browser is signed
            in. This is the WhatsApp Web flow, and how Telegram Web, Discord, Steam
            and most TV apps onboard a second screen.
          </p>

          <CodeBlock language="bash" code={`grit plugin add device-pairing
grit migrate`} />

          {/* The flow */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">How it works</h2>

          <LaneFlow
            id="pairing"
            lanes={['Browser (anonymous)', 'API', 'Signed-in device']}
            nodes={[
              { id: 'start', lane: 0, row: 0, title: 'POST /pair/start', sub: 'gets code + QR', tone: 'blue' },
              { id: 'poll', lane: 0, row: 2, title: 'GET /pair/:code', sub: 'polls every 1s', tone: 'blue' },
              { id: 'mint', lane: 1, row: 0, title: 'Mint code', sub: '32 bytes, 2 min TTL', tone: 'primary' },
              { id: 'claim', lane: 1, row: 2, title: 'Claim once', sub: 'conditional UPDATE', tone: 'primary' },
              { id: 'describe', lane: 2, row: 1, title: 'What am I approving?', sub: 'user agent + IP', tone: 'amber' },
              { id: 'approve', lane: 2, row: 2, title: 'Approve or deny', sub: 'one tap either way', tone: 'green' },
            ]}
            edges={[
              { from: 'start', to: 'mint', tone: 'primary' },
              { from: 'mint', to: 'describe', label: 'scanned', dashed: true, tone: 'amber' },
              { from: 'describe', to: 'approve', tone: 'green' },
              { from: 'approve', to: 'claim', tone: 'primary' },
              { from: 'claim', to: 'poll', label: 'tokens', dashed: true, tone: 'blue' },
            ]}
            legend={[
              { tone: 'amber', label: 'Shown before approving' },
              { tone: 'primary', label: 'Single use' },
            ]}
            caption="The approver sees the browser and address before they approve, and the code can only be spent once"
          />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">What you get</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left">
                  <th className="py-2 pr-4 font-medium text-foreground">Endpoint</th>
                  <th className="py-2 pr-4 font-medium text-foreground">Auth</th>
                  <th className="py-2 font-medium text-foreground">Does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                {[
                  ['POST /pair/start', 'none', 'Mints a code and returns it with a QR PNG data URI'],
                  ['GET /pair/:code', 'none', 'Polled by the browser; returns tokens once approved'],
                  ['GET /pair/:code/request', 'session', 'What the code is asking for: user agent, IP, times'],
                  ['POST /pair/:code/approve', 'session', 'Binds the request to the caller'],
                  ['POST /pair/:code/deny', 'session', 'Refuses it and burns the code'],
                ].map(([ep, auth, what]) => (
                  <tr key={ep} className="border-b border-border/50">
                    <td className="py-2 pr-4 font-mono text-xs text-foreground">{ep}</td>
                    <td className="py-2 pr-4 text-xs">{auth}</td>
                    <td className="py-2">{what}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <p className="mt-4 leading-relaxed text-muted-foreground">
            Plus a <code>/link</code> page in the web app that shows the QR and polls,
            and a <strong>System &rarr; Link a device</strong> screen in the admin for
            approving one.
          </p>

          {/* Security */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            Why it is built this way
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            A pairing code is a bearer credential for a whole account. Most of the
            design here exists because of that, and the parts that look like extra
            work are the parts that make the difference between a QR somebody
            photographed across a room being useless and being a silent takeover.
          </p>

          <div className="space-y-4">
            <div className="rounded-lg border border-amber-500/30 bg-amber-500/5 p-4">
              <div className="mb-2 flex items-center gap-2">
                <AlertTriangle className="h-4 w-4 text-amber-400" />
                <strong className="text-sm text-foreground">
                  You see what you are approving
                </strong>
              </div>
              <p className="text-sm leading-relaxed text-muted-foreground">
                Approval is two steps. The approving device fetches the browser and
                IP behind the code and shows them before it asks. Approving an opaque
                code is not consent: the user has no way to tell their own laptop from
                somebody who read the QR over their shoulder. <strong>Deny</strong> is
                a first-class endpoint next to it, because &quot;that wasn&apos;t
                me&quot; has to be one tap.
              </p>
            </div>

            <div className="rounded-lg border border-primary/20 bg-primary/5 p-4">
              <div className="mb-2 flex items-center gap-2">
                <ShieldCheck className="h-4 w-4 text-primary" />
                <strong className="text-sm text-foreground">Single use, enforced by the database</strong>
              </div>
              <p className="text-sm leading-relaxed text-muted-foreground">
                The claim is a conditional <code>UPDATE</code>, not a read followed by
                a write. Two polls arriving together would otherwise both see
                &quot;approved, unclaimed&quot; and both walk away with a token pair.
                Approval is guarded the same way, so two devices racing on one code
                cannot both bind it. The generated test suite races eight approvals
                and asserts exactly one wins.
              </p>
            </div>

            <div className="rounded-lg border border-border p-4">
              <div className="mb-2 flex items-center gap-2">
                <Timer className="h-4 w-4 text-muted-foreground" />
                <strong className="text-sm text-foreground">Short-lived, and rate limited</strong>
              </div>
              <p className="text-sm leading-relaxed text-muted-foreground">
                Codes are 32 bytes from <code>crypto/rand</code> and expire in two
                minutes. A QR on a screen in a public place is visible to everyone in
                the room, so the window has to be small enough that somebody would
                have to be waiting for it. <code>/pair/start</code> is anonymous, so
                one address may hold five codes in flight; the row is deleted the
                moment it is claimed, denied or expires.
              </p>
            </div>
          </div>

          {/* Sessions */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            A paired browser is an ordinary session
          </h2>
          <p className="leading-relaxed text-muted-foreground">
            There is no <code>Device</code> model. Approval creates a normal session
            row and sets the same HttpOnly cookies a password login does, so a paired
            browser appears under <strong>Account &rarr; Security</strong> alongside
            everything else and is signed out from there. A parallel device table
            would be a second list of the same thing, drifting from the first.
          </p>
          <p className="mt-4 leading-relaxed text-muted-foreground">
            Revoking that session closes its realtime socket immediately, so a device
            you sign out stops receiving pushed events as well as losing its API
            access. See{' '}
            <Link href="/docs/backend/realtime" className="text-primary hover:underline">
              Realtime
            </Link>
            .
          </p>

          {/* Customising */}
          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Customising it</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The plugin writes ordinary files into your repo. Nothing here is
            framework indirection, so change what you need:
          </p>
          <CodeBlock
            language="text"
            code={`apps/api/internal/models/pairing_request.go       the handshake row
apps/api/internal/handlers/device_pairing.go     the five endpoints
apps/api/internal/handlers/device_pairing_test.go the race tests
apps/web/app/link/page.tsx                       the QR screen
apps/web/hooks/use-pairing.ts                    start + poll
apps/admin/app/(dashboard)/system/link-device/   the approval screen`}
          />
          <p className="mt-4 leading-relaxed text-muted-foreground">
            The two constants worth knowing are at the top of the handler:{' '}
            <code>pairingTTL</code> (two minutes) and <code>pairingPerIP</code> (five
            in flight). Lengthen the TTL and you widen the window an onlooker has;
            that is the trade, and it is yours to make.
          </p>

          <div className="mt-4 rounded-lg border border-amber-500/30 bg-amber-500/5 p-4">
            <p className="text-sm leading-relaxed text-muted-foreground">
              <strong className="text-amber-400">If you add a scanner.</strong> The
              admin screen takes a typed or pasted code, and accepts either the bare
              code or the whole URL a QR scanner returns. If you wire a camera to it,
              keep the typed field: a camera is not always available, not always
              permitted, and not always the thing the user reaches for. A pairing
              flow with no fallback strands people.
            </p>
          </div>

          {/* Nav */}
          <div className="mt-16 flex items-center justify-between border-t border-border pt-8">
            <Button variant="ghost" asChild>
              <Link href="/docs/plugins/saved-views">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Saved views
              </Link>
            </Button>
            <Button variant="ghost" asChild>
              <Link href="/docs/plugins/authoring">
                Writing a plugin
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
