import type { Metadata } from 'next'
import { ErrorPage, Why } from '@/components/error-page'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Docker is not running',
  description: 'The Docker daemon did not answer. What that means on each operating system, and whether you need Docker at all.',
  alternates: { canonical: 'https://gritframework.dev/e/docker-not-running' },
}

export default function Page() {
  return (
    <ErrorPage
      slug="docker-not-running"
      title="Docker is not running"
      cause="The Docker daemon did not answer. On macOS and Windows that almost always means Docker Desktop is closed."
      fix={
        <>
          <p>
            <strong>macOS and Windows:</strong> open Docker Desktop and wait for its icon to stop
            animating. It is an application, and nothing works until it is running.
          </p>
          <p>
            <strong>Linux:</strong> start the service.
          </p>
          <CodeBlock
            terminal
            code={`sudo systemctl start docker
sudo systemctl enable docker   # and at every boot`}
          />
        </>
      }
      confirm={
        <CodeBlock terminal code={`docker run --rm hello-world`} />
      }
    >
      <Why>
        <p>
          The message rarely says &quot;Docker Desktop is closed&quot;. It says something about a
          socket or a named pipe, because the client found no daemon at the end of one. That is
          why this is the most reported Docker error there is.
        </p>
        <p>
          On Linux there is a second cause with a different message: your user is not in the{' '}
          <code>docker</code> group, so only <code>sudo docker</code> works. Add yourself and start
          a new session.
        </p>
        <CodeBlock
          terminal
          code={`sudo usermod -aG docker $USER
newgrp docker`}
        />
        <p>
          <strong>You may not need Docker at all.</strong> A project on SQLite with local file
          storage runs with nothing but Go and Node. Scaffold one with{' '}
          <code>grit new myapp --db sqlite --storage local</code>, and skip the container checks
          with <code>grit doctor --env --no-docker</code>.
        </p>
      </Why>
    </ErrorPage>
  )
}
