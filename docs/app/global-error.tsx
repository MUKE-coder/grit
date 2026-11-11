'use client'

/*
 * The boundary of last resort: an error in the root layout itself.
 *
 * This one replaces the layout rather than rendering inside it, so there is no
 * header, no font and, critically, no stylesheet: the layout that imports
 * globals.css is the thing that failed. Every style here is inline for that
 * reason, and the colours are literals rather than CSS variables because the
 * file defining those variables was never loaded.
 *
 * It renders on a server error too, so it must not reach for anything that
 * only exists in a browser at module scope.
 */

// The dark palette, written out. A page that falls back to white here would
// flash white on a site that is dark everywhere else, at the exact moment
// somebody is already unsure whether the site is broken.
const BG = '#0b0d0f'
const FG = '#e8e8f0'
const MUTED = '#9090a8'
const BORDER = '#23272d'
const ACCENT = '#a99bff'

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <html lang="en">
      <body
        style={{
          margin: 0,
          minHeight: '100vh',
          backgroundColor: BG,
          color: FG,
          fontFamily:
            'ui-sans-serif, system-ui, -apple-system, "Segoe UI", Roboto, Helvetica, Arial, sans-serif',
        }}
      >
        <div
          style={{
            display: 'flex',
            minHeight: '100vh',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '1.5rem',
          }}
        >
          <div style={{ maxWidth: '32rem', textAlign: 'center' }}>
            <p
              style={{
                margin: '0 0 1.25rem',
                fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
                fontSize: '0.75rem',
                letterSpacing: '0.08em',
                textTransform: 'uppercase',
                color: ACCENT,
              }}
            >
              Grit
            </p>

            <h1
              style={{
                margin: '0 0 1rem',
                fontSize: '1.875rem',
                lineHeight: 1.2,
                fontWeight: 700,
                letterSpacing: '-0.02em',
              }}
            >
              The site failed to load
            </h1>

            <p style={{ margin: '0 0 1.5rem', fontSize: '1rem', lineHeight: 1.6, color: MUTED }}>
              Something went wrong before the page could be drawn. Reloading
              usually fixes it, because the commonest cause is a tab that was
              open while a new version was deployed.
            </p>

            {error.digest && (
              <p
                style={{
                  margin: '0 0 1.5rem',
                  fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
                  fontSize: '0.75rem',
                  color: MUTED,
                }}
              >
                Reference: {error.digest}
              </p>
            )}

            <div
              style={{
                display: 'flex',
                flexWrap: 'wrap',
                gap: '0.75rem',
                justifyContent: 'center',
              }}
            >
              {/* A full reload, not reset(). If the root layout threw, the tree
                  reset() re-renders is the one that just failed. */}
              <button
                type="button"
                onClick={() => window.location.reload()}
                style={{
                  cursor: 'pointer',
                  border: 'none',
                  borderRadius: '0.5rem',
                  backgroundColor: ACCENT,
                  color: '#14121f',
                  padding: '0.7rem 1.1rem',
                  fontSize: '0.875rem',
                  fontWeight: 600,
                }}
              >
                Reload the page
              </button>

              <button
                type="button"
                onClick={reset}
                style={{
                  cursor: 'pointer',
                  borderRadius: '0.5rem',
                  border: '1px solid ' + BORDER,
                  backgroundColor: 'transparent',
                  color: FG,
                  padding: '0.7rem 1.1rem',
                  fontSize: '0.875rem',
                  fontWeight: 500,
                }}
              >
                Try again
              </button>

              <a
                href="https://github.com/MUKE-coder/grit/issues/new"
                target="_blank"
                rel="noreferrer"
                style={{
                  display: 'inline-flex',
                  alignItems: 'center',
                  borderRadius: '0.5rem',
                  border: '1px solid ' + BORDER,
                  color: FG,
                  padding: '0.7rem 1.1rem',
                  fontSize: '0.875rem',
                  fontWeight: 500,
                  textDecoration: 'none',
                }}
              >
                Report this
              </a>
            </div>
          </div>
        </div>
      </body>
    </html>
  )
}
