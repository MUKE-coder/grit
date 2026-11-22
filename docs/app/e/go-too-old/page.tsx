import type { Metadata } from 'next'
import { ErrorPage, Why } from '@/components/error-page'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Your Go is older than Grit needs',
  description: 'A generated project pins a Go version for a reason: older patches carry standard-library vulnerabilities the code reaches.',
  alternates: { canonical: 'https://gritframework.dev/e/go-too-old' },
}

export default function Page() {
  return (
    <ErrorPage
      slug="go-too-old"
      title="Your Go is older than Grit needs"
      cause="A generated project\u2019s go.mod asks for a newer Go than the one on PATH."
      fix={
        <>
          <p>Install the current Go from go.dev/dl, then open a new terminal.</p>
          <CodeBlock
            terminal
            code={`# macOS
brew upgrade go

# Anywhere: the tarball or installer from go.dev/dl`}
          />
        </>
      }
      confirm={
        <CodeBlock terminal code={`go version
grit doctor --env`} />
      }
    >
      <Why>
        <p>
          The pin is not about language features. Older patch releases carry standard-library
          vulnerabilities that a generated project reaches before anybody has signed in: several
          in <code>net/http</code>, and one in the XML parsing behind the SAML sign-in path.
        </p>
        <p>
          That is why the <code>go</code> directive in a generated <code>go.mod</code> is as much
          a security decision as any dependency version, and why <code>grit doctor</code> reports
          a toolchain below the floor rather than leaving it to you to notice.
        </p>
        <p>
          A newer Go than the pin is fine. The check is a floor, not an exact match.
        </p>
      </Why>
    </ErrorPage>
  )
}
