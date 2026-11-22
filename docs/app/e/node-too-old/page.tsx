import type { Metadata } from 'next'
import { ErrorPage, Why } from '@/components/error-page'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Your Node is older than the frontend needs',
  description: 'Next.js declares a minimum Node version, and an older one fails during the build rather than at install. How to upgrade.',
  alternates: { canonical: 'https://gritframework.dev/e/node-too-old' },
}

export default function Page() {
  return (
    <ErrorPage
      slug="node-too-old"
      title="Your Node is older than the frontend needs"
      cause="Next.js declares a minimum Node version in its engines field, and the one on PATH is below it."
      fix={
        <>
          <CodeBlock
            terminal
            code={`fnm install 22 && fnm use 22
# or
nvm install 22 && nvm use 22`}
          />
          <p>
            Set it as the default too, or a new terminal goes back to the old one:{' '}
            <code>fnm default 22</code>, or <code>nvm alias default 22</code>.
          </p>
        </>
      }
      confirm={
        <CodeBlock terminal code={`node --version
grit doctor --env`} />
      }
    >
      <Why>
        <p>
          The unpleasant part of this one is when it surfaces. An install with an old Node often
          succeeds with a warning, and the failure arrives later during the build, as an error
          that does not mention Node at all.
        </p>
        <p>
          If you have several projects on different versions, put a <code>.nvmrc</code> or{' '}
          <code>.node-version</code> in each and let the manager switch for you on{' '}
          <code>cd</code>.
        </p>
      </Why>
    </ErrorPage>
  )
}
