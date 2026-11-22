import type { Metadata } from 'next'
import { ErrorPage, Why } from '@/components/error-page'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'go: command not found',
  description: 'Go is not installed, or is installed and not on PATH. The install for each operating system, and the reason reinstalling often does not help.',
  alternates: { canonical: 'https://gritframework.dev/e/go-not-found' },
}

export default function Page() {
  return (
    <ErrorPage
      slug="go-not-found"
      title="go: command not found"
      cause="The go command could not be run. Either Go is not installed, or it is installed and its bin directory is not on PATH."
      fix={
        <>
          <p>
            Install Go, then <strong>open a new terminal</strong>. A shell that was already
            running kept the old PATH and will not see it.
          </p>
          <CodeBlock
            terminal
            code={`# macOS
brew install go

# Windows: use the .msi from go.dev/dl, which sets PATH for you.

# Debian and Ubuntu: the apt package is usually years behind, so take
# the tarball from go.dev/dl instead.
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go*.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile`}
          />
        </>
      }
      confirm={
        <CodeBlock terminal code={`go version`} />
      }
    >
      <Why>
        <p>
          Reinstalling usually does not help, because the problem is normally PATH rather than the
          installation. A shell reads PATH once, when it starts, so a terminal that was already
          open keeps the old value no matter how many times you install.
        </p>
        <p>
          On Windows this includes the terminal inside your editor: restart the editor, not just
          the panel.
        </p>
        <p>
          A related case: <code>go version</code> works but a tool you installed with{' '}
          <code>go install</code> does not. That is a different directory. Add{' '}
          <code>$(go env GOPATH)/bin</code>, which is <code>%USERPROFILE%\go\bin</code> on
          Windows.
        </p>
      </Why>
    </ErrorPage>
  )
}
