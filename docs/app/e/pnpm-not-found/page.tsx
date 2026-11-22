import type { Metadata } from 'next'
import { ErrorPage, Why } from '@/components/error-page'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'pnpm: command not found',
  description: 'Generated projects install with pnpm rather than npm, because the workspace layout depends on it. How to install it, and why the version is pinned.',
  alternates: { canonical: 'https://gritframework.dev/e/pnpm-not-found' },
}

export default function Page() {
  return (
    <ErrorPage
      slug="pnpm-not-found"
      title="pnpm: command not found"
      cause="Generated projects install with pnpm, not npm, because the workspace layout depends on it."
      fix={
        <>
          <p>Corepack ships with Node and installs the version the project asks for.</p>
          <CodeBlock
            terminal
            code={`corepack enable
corepack prepare pnpm@latest --activate

# or, without corepack
npm install -g pnpm`}
          />
        </>
      }
      confirm={
        <CodeBlock terminal code={`pnpm --version`} />
      }
    >
      <Why>
        <p>
          A generated project is a pnpm workspace: the apps and the shared package resolve each
          other through it. Installing with npm or yarn instead produces a tree where{' '}
          <code>packages/shared</code> is not linked, and the first import of a shared type fails
          to resolve.
        </p>
        <p>
          A project&apos;s <code>package.json</code> names a pnpm version in{' '}
          <code>packageManager</code>, and corepack honours it. That pin is deliberate: a major
          pnpm release changes how the lockfile and <code>node_modules</code> are laid out, and a
          project installed with two different majors behaves differently on two machines.
        </p>
      </Why>
    </ErrorPage>
  )
}
