/**
 * Bare layout for iframe previews — nothing that could influence how a block
 * lays out. The root layout still supplies fonts and the Tailwind base.
 */
export default function PreviewLayout({ children }: { children: React.ReactNode }) {
  return <>{children}</>
}
