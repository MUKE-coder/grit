import { DocsAutoToc } from '@/components/docs-auto-toc'

// Docs pages still render their own <SiteHeader /> + <DocsSidebar /> + <main>,
// so this layout intentionally just passes children through and augments them
// with the shared, self-scanning "On this page" rail. Mounting it here means
// every /docs page gets a TOC without per-page wiring.
export default function DocsLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      {children}
      <DocsAutoToc />
    </>
  )
}
