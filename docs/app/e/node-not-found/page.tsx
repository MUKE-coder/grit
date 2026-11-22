import type { Metadata } from 'next'
import { ErrorPage, Why } from '@/components/error-page'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'node: command not found',
  description: 'Node is not installed or not on PATH. Installing it, why a version manager is worth it, and when you do not need Node at all.',
  alternates: { canonical: 'https://gritframework.dev/e/node-not-found' },
}

export default function Page() {
  return (
    <ErrorPage
      slug="node-not-found"
      title="node: command not found"
      cause="The node command could not be run. The frontends need it; an API-only project does not."
      fix={
        <>
          <p>
            Install Node with a version manager rather than a one-off installer. You will
            eventually need two versions on one machine, and a manager makes that a non-event.
          </p>
          <CodeBlock
            terminal
            code={`# fnm: fast, and works on Windows too
fnm install 22 && fnm use 22

# or nvm, on macOS and Linux
nvm install 22 && nvm use 22`}
          />
          <p>Then open a new terminal.</p>
        </>
      }
      confirm={
        <CodeBlock terminal code={`node --version`} />
      }
    >
      <Why>
        <p>
          Same cause as the Go version of this: installing changes PATH, and a shell that was
          already running keeps the old one.
        </p>
        <p>
          <strong>You may not need Node.</strong> An API-only project has no frontend, so nothing
          needs it. Scaffold one with <code>grit new myapp --api</code>, and skip the frontend
          checks with <code>grit doctor --env --no-node</code>.
        </p>
      </Why>
    </ErrorPage>
  )
}
