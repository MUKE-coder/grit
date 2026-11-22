import type { Metadata } from 'next'
import { ErrorPage, Why } from '@/components/error-page'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Port already in use',
  description: 'A port a Grit project needs is taken by something else: how to find what has it, and how to free it or move the project.',
  alternates: { canonical: 'https://gritframework.dev/e/port-in-use' },
}

export default function Page() {
  return (
    <ErrorPage
      slug="port-in-use"
      title="A port the project needs is already taken"
      cause="Something is already listening on it. Usually a server from an earlier run whose terminal was closed without stopping the process."
      fix={
        <>
          <p>Find what has the port and stop it. Replace 8080 with the port in the message.</p>
          <CodeBlock
            terminal
            code={`# macOS and Linux
lsof -ti:8080 | xargs kill

# Windows PowerShell
Get-NetTCPConnection -LocalPort 8080 -State Listen |
  Select-Object -Expand OwningProcess |
  ForEach-Object { Stop-Process -Id $_ -Force }`}
          />
          <p>
            Or move the project instead: change the port in <code>.env</code>. If you change{' '}
            <code>APP_PORT</code>, change <code>APP_URL</code> with it, or links the API
            generates will point at the old one.
          </p>
        </>
      }
      confirm={
        <>
          <CodeBlock terminal code={`grit doctor --env`} />
          <p>Every port it lists should say free.</p>
        </>
      }
    >
      <Why>
        <p>
          A process keeps its ports until it exits, and closing a terminal window does not always
          stop what was running in it. The usual culprits are a previous{' '}
          <code>grit start</code>, a dev server from another project, or containers from another
          project&apos;s compose file.
        </p>
        <p>
          Two ports are special on macOS and are not strays: 5000 belongs to AirPlay Receiver and
          7000 to the Control Center. Killing those is not the fix; pick another port, or turn the
          feature off in System Settings.
        </p>
        <p>
          A generated project deliberately uses 5434 for Postgres and 6380 for Redis rather than
          their defaults, precisely so that a native install of either does not collide. If one of
          those is busy, something else unusual has it.
        </p>
      </Why>
    </ErrorPage>
  )
}
