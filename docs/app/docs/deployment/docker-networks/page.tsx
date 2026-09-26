import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/deployment/docker-networks')

export default function DockerNetworksPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Deployment</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                &quot;All predefined address pools have been fully subnetted&quot;
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every image builds, then the deploy dies at the last step. It is a limit of the
                Docker daemon on the host, not a fault in your project, and there are three ways
                out.
              </p>
            </div>

            <div className="prose-grit">
              <h2 id="symptom">What it looks like</h2>
              <CodeBlock
                language="text"
                code={`Image shop-api  Built
Network shop_default  Creating
Network shop_default  Error  Error response from daemon:
  all predefined address pools have been fully subnetted
failed to create network ...
Error: Docker command failed`}
              />
              <p>
                Builds succeed, nothing starts, and retrying gives the same error. Nothing is said
                about disk, memory or CPU because none of them are the problem, which is why this
                tends to cost an afternoon the first time.
              </p>

              <h2 id="cause">Why it happens</h2>
              <p>
                Docker hands out bridge-network subnets from a fixed pool. By default that is{' '}
                <code>172.17.0.0/12</code> sliced into <code>/16</code> blocks, which is about{' '}
                <strong>sixteen networks for the entire daemon</strong>. Every Compose project takes
                at least one, and a project that declares a network of its own can take two.
              </p>
              <p>
                They are not returned when a deploy fails or a project is deleted. Run half a dozen
                stacks on one host, delete and redeploy a few times, and the pool empties. After
                that no new stack starts, however idle the machine is.
              </p>

              <h2 id="grit">What Grit generates</h2>
              <p>
                Since v3.333.0, <code>docker-compose.prod.yml</code> declares{' '}
                <strong>no network</strong>. Compose&apos;s implicit <code>&lt;project&gt;_default</code>{' '}
                already gives service-name DNS and the same isolation, so a named one took a second
                subnet and bought nothing. It also pins no <code>container_name</code>: names are
                global to the daemon, so a second copy of the stack, staging beside production,
                could not start.
              </p>
              <p>
                <code>grit upgrade</code> applies the same two changes to an existing project, so
                the stacks you have already deployed converge on it rather than keeping the
                arrangement that causes this.
              </p>

              <h2 id="full">If the pool is already empty</h2>
              <p>
                One fewer network per stack does not help a host with none left: even the implicit
                default fails to allocate. Every project ships{' '}
                <code>docker-compose.shared-network.yml</code>, which joins a network that already
                exists and therefore asks for no subnet at all:
              </p>
              <CodeBlock
                terminal
                code={`docker compose -f docker-compose.prod.yml \\
               -f docker-compose.shared-network.yml up -d --build`}
              />
              <p>
                On Dokploy, set the Compose Path to both files in that order. It defaults to{' '}
                <code>dokploy-network</code>; set <code>SHARED_NETWORK</code> if your proxy uses
                another.
              </p>

              <Callout type="warning" title="A shared network is shared">
                Containers on it resolve each other by service name. If another project on the host
                also calls a service <code>postgres</code>, whichever answers first is what you get,
                and a silent connection to somebody else&apos;s database is a worse day than a
                failed deploy. The overlay therefore gives every service a project-prefixed alias
                and points the API at those rather than the bare names. Treat it as a stopgap.
              </Callout>

              <h2 id="host">The fix that belongs on the host</h2>
              <p>First, return what failed deploys left behind. No restart, no downtime:</p>
              <CodeBlock terminal code={`docker network prune -f`} />
              <p>
                That removes only networks with nothing attached. Running stacks are untouched.
                Then widen the pool so it cannot recur, in <code>/etc/docker/daemon.json</code>:
              </p>
              <CodeBlock
                language="json"
                code={`{ "default-address-pools": [ { "base": "172.17.0.0/12", "size": 24 } ] }`}
              />
              <CodeBlock terminal code={`sudo systemctl restart docker`} />
              <p>
                <code>/24</code> blocks give about 4096 networks instead of 16, and a stack needs a
                handful of addresses rather than 65,000. The restart bounces every container on the
                host, so pick a quiet window, and merge the key if the file already exists.
              </p>
              <p>
                Once that is done, go back to <code>docker-compose.prod.yml</code> on its own and
                keep one private network per stack.
              </p>

              <h2 id="check">Before deploying to a busy host</h2>
              <CodeBlock
                terminal
                code={`docker network ls | wc -l          # above ~12 is the smell
docker network ls                  # look for stale <project>_default
docker network prune -f            # after deleting any project`}
              />
            </div>

            <div className="mt-12 flex items-center justify-between border-t border-border pt-6">
              <Button variant="ghost" asChild>
                <Link href="/docs/deployment">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Deployment
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/deployment/checklist">
                  Deployment checklist
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
