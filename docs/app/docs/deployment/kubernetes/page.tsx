import Link from 'next/link'
import { ArrowLeft, ArrowRight, Boxes } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/deployment/kubernetes')

/**
 * Kubernetes output.
 *
 * Written as the manifests' reasoning rather than as a Kubernetes tutorial:
 * somebody deploying a Grit project to a cluster already knows what a
 * Deployment is, and what they cannot know is why the migration is a separate
 * Job, why the probes are not /api/health, and why there is no worker.
 */
export default function KubernetesPage() {
  return (
    <DocsShell>
      <div className="mb-3 flex items-center gap-2">
        <Boxes className="h-5 w-5 text-accent" />
        <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
          Deployment
        </span>
      </div>

      <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
        Kubernetes
      </h1>
      <p className="mb-6 text-lg text-muted-foreground">
        <code>grit deploy --k8s</code> writes the manifests for your project into{' '}
        <code>deploy/k8s</code>: the API, the frontends it actually has, a migration
        Job, a horizontal autoscaler, a config map, a secret template and a
        kustomization that ties them together. It deploys nothing. The files are
        yours to read, change and commit.
      </p>

      <div className="prose-grit">
        <CodeBlock
          language="bash"
          code={`grit deploy --k8s --k8s-image registry.example.com/contacts-api:v1.4.2`}
        />

        <h2 id="apply">Applying them</h2>
        <p>
          Three steps, and the order matters.
        </p>
        <CodeBlock
          language="bash"
          code={`# 1. The secret, once. It is NOT in kustomization.yaml.
cp deploy/k8s/secret.example.yaml deploy/k8s/secret.yaml
$EDITOR deploy/k8s/secret.yaml
kubectl apply -f deploy/k8s/secret.yaml

# 2. The migration, every release, and wait for it.
kubectl apply -f deploy/k8s/migrate-job.yaml
kubectl wait --for=condition=complete --timeout=5m \\
  job/contacts-migrate-v1-4-2 -n contacts

# 3. Everything else.
kubectl apply -k deploy/k8s`}
        />

        <h2 id="kustomize">Plain manifests, not a chart</h2>
        <p>
          These are plain Kubernetes objects with a kustomization, so{' '}
          <code>kubectl apply -k</code> works with nothing else installed. That is
          partly taste and mostly testability: every file is valid YAML, so Grit can
          check its own output by parsing it. A Helm chart cannot be checked that way,
          because <code>{'{{- if .Values.x }}'}</code> is not YAML, and a template that
          produces broken manifests looks exactly like one that does not until somebody
          runs it against a cluster.
        </p>
        <p>
          If your team standardises on Helm, wrap these or use them as the body of a
          chart. Nothing here depends on how they are applied.
        </p>

        <h2 id="probes">Why the probes are not /api/health</h2>
        <p>
          A Grit API serves three health endpoints, and they answer different
          questions. <code>/api/health</code> probes every dependency it has, which is
          what the System Health page needs. <code>/healthz</code> and{' '}
          <code>/readyz</code> exist for an orchestrator: liveness probes nothing on
          purpose, and readiness gates on the database alone.
        </p>
        <p>
          The manifests use the second pair. Pointing readiness at{' '}
          <code>/api/health</code> means a slow object store fails readiness and pulls
          every pod out of the load balancer over something that does not stop it
          serving requests.
        </p>

        <h2 id="migration">Why the migration is a Job, applied first</h2>
        <p>
          An API pod that starts against an un-migrated database answers{' '}
          <em>relation does not exist</em> to every request, so the migration has to
          finish before the new pods roll.
        </p>
        <p>
          It is not an initContainer on the API Deployment, because that runs once per
          replica, concurrently. It is a Job, applied and waited for.
        </p>
        <p>
          The Job&apos;s name carries the image tag, which is the part people are
          surprised by. A Job&apos;s pod template is immutable: a fixed name means the
          second deploy is rejected by the API server for changing a field that cannot
          change, and the deploy fails on its first step with a message about
          immutability rather than about anything you did.
        </p>

        <h2 id="worker">There is no worker Deployment</h2>
        <p>
          A Grit API process starts its own asynq worker whenever Redis is reachable.
          There is no worker-only binary, so every API replica is also a queue
          consumer. Two things follow:
        </p>
        <ul>
          <li>
            Scaling for HTTP load scales queue consumers with it. Usually what you
            want, occasionally not.
          </li>
          <li>
            A slow or memory-hungry job runs inside a pod that is also serving
            requests.
          </li>
        </ul>
        <p>
          If that is wrong for your workload, say so on the issue tracker: separating
          them needs a worker-only mode in the binary, which does not exist yet.
        </p>

        <h2 id="secret">The secret is a template, and deliberately unreachable</h2>
        <p>
          <code>secret.example.yaml</code> holds placeholders and is not listed in{' '}
          <code>kustomization.yaml</code>, so <code>apply -k</code> cannot install it.
          It is the one file in the set that is dangerous when applied by accident: an
          API that starts with a placeholder signing key mints tokens anybody can
          forge and reports itself perfectly healthy.
        </p>

        <h2 id="dependencies">Postgres, Redis and storage are yours</h2>
        <p>
          These manifests do not install a database, a cache or an object store. On a
          cluster those are nearly always a managed service or your own operator, and a
          generated single-replica Postgres with no backups looks production-ready and
          is not. The secret names yours instead.
        </p>
        <p>
          Redis is not optional. The API starts its worker only when Redis answers, so
          without it no background job ever runs.
        </p>

        <h2 id="what-to-change">What you will want to change</h2>
        <ul>
          <li>
            <code>ingress.yaml</code> assumes an nginx ingress controller and
            cert-manager, and is only written when your <code>.env</code> names a real{' '}
            <code>APP_URL</code>. A development <code>localhost</code> writes no
            Ingress, because an Ingress nobody can reach is worse than none.
          </li>
          <li>
            Resource requests and limits are a starting point, not a measurement. Run{' '}
            <Link href="/docs/batteries/pulse">Pulse</Link> and a load test before you
            trust them.
          </li>
          <li>
            <code>api-hpa.yaml</code> scales on CPU. Delete it to pin replicas.
          </li>
        </ul>

        <h2 id="related">Before you scale out</h2>
        <p>
          <Link href="/docs/deployment/multiple-instances">
            Running more than one instance
          </Link>{' '}
          is the page that says what several copies of a Grit API share and what they
          still count separately. The manifests here assume you have read it: they set
          the same secret on every replica and run migrations once, which are two of
          the things that page is about.
        </p>
      </div>

      <div className="mt-12 flex flex-wrap items-center justify-between gap-4 border-t border-border pt-8">
        <Button asChild variant="ghost" size="sm">
          <Link href="/docs/deployment/multiple-instances">
            <ArrowLeft className="mr-1.5 h-4 w-4" />
            More than one instance
          </Link>
        </Button>
        <Button asChild size="sm">
          <Link href="/docs/deployment/checklist">
            Go-live checklist
            <ArrowRight className="ml-1.5 h-4 w-4" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
