import type { Metadata } from 'next'
import { ErrorPage, Why } from '@/components/error-page'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'This is not a Grit project',
  description: 'Grit commands look for grit.json at the project root. What to do when the command cannot find one.',
  alternates: { canonical: 'https://gritframework.dev/e/not-a-grit-project' },
}

export default function Page() {
  return (
    <ErrorPage
      slug="not-a-grit-project"
      title="This is not a Grit project"
      cause="The command looks for grit.json at the project root and did not find it here or in any directory above."
      fix={
        <>
          <p>Move into the directory that grit new created.</p>
          <CodeBlock
            terminal
            code={`cd myapp
ls grit.json      # it should be here`}
          />
          <p>Or make one, if you have not yet:</p>
          <CodeBlock terminal code={`grit new myapp`} />
        </>
      }
      confirm={
        <CodeBlock terminal code={`grit doctor`} />
      }
    >
      <Why>
        <p>
          Almost every Grit command needs to know where the project root is, because it writes
          into several apps at once and the paths differ by architecture. <code>grit.json</code>{' '}
          is how it finds that root, and the search walks upwards, so being in a subdirectory is
          fine.
        </p>
        <p>
          Two cases that look like this and are not. If <code>grit.json</code> exists and the
          command still refuses, it may be unreadable or not valid JSON; open it. And{' '}
          <code>grit doctor --env</code> deliberately does not need a project, because it answers
          whether the machine is ready before there is one.
        </p>
      </Why>
    </ErrorPage>
  )
}
