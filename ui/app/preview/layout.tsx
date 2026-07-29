/**
 * Bare layout for iframe previews — no header, no padding, nothing that would
 * change how the component lays out. The root layout still supplies the fonts
 * and the token variables.
 */
export default function PreviewLayout({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen bg-background">{children}</div>
}
