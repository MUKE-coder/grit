import Link from "next/link";

// The footer, whose links come out of the database.
//
// Every entry under "Information" is a row in `pages`. Adding a fourth needs no
// deploy and no code: the row appears here and renders at /[handle] through one
// dynamic route. That is the part of a CMS a shop actually uses.
export function ShopFooter({
  collections,
  pages,
}: {
  collections: { href: string; label: string }[];
  pages: { href: string; label: string }[];
}) {
  return (
    <footer className="mt-20 border-t border-border bg-bg-secondary">
      <div className="mx-auto grid max-w-7xl gap-10 px-4 py-14 sm:px-6 md:grid-cols-4">
        <div>
          <Link href="/" className="flex items-center gap-2.5">
            <span className="flex h-8 w-8 items-center justify-center rounded-lg border border-accent/20 bg-accent/15 font-mono text-sm font-bold text-accent">
              A
            </span>
            <span className="text-lg font-bold tracking-tight">Aura</span>
          </Link>
          <p className="mt-3 max-w-xs text-sm leading-relaxed text-text-secondary">
            A demonstration shop, built with Grit. The catalogue, the variants,
            the basket and the admin are all real.
          </p>
        </div>

        <nav aria-labelledby="footer-shop">
          <h2 id="footer-shop" className="mb-3 text-sm font-semibold">
            Shop
          </h2>
          <ul className="space-y-2">
            <li>
              <Link
                href="/search"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Everything
              </Link>
            </li>
            {collections.map((c) => (
              <li key={c.href}>
                <Link
                  href={c.href}
                  className="text-sm text-text-secondary transition-colors hover:text-foreground"
                >
                  {c.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <nav aria-labelledby="footer-info">
          <h2 id="footer-info" className="mb-3 text-sm font-semibold">
            Information
          </h2>
          <ul className="space-y-2">
            {pages.map((p) => (
              <li key={p.href}>
                <Link
                  href={p.href}
                  className="text-sm text-text-secondary transition-colors hover:text-foreground"
                >
                  {p.label}
                </Link>
              </li>
            ))}
            <li>
              <Link
                href="/blog"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Journal
              </Link>
            </li>
          </ul>
        </nav>

        <nav aria-labelledby="footer-built">
          <h2 id="footer-built" className="mb-3 text-sm font-semibold">
            Built with
          </h2>
          <ul className="space-y-2">
            <li>
              <a
                href="https://gritframework.dev"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Grit
              </a>
            </li>
            <li>
              <a
                href="https://gritframework.dev/docs"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Documentation
              </a>
            </li>
            <li>
              <a
                href="https://github.com/MUKE-coder/grit"
                className="text-sm text-text-secondary transition-colors hover:text-foreground"
              >
                Source
              </a>
            </li>
          </ul>
        </nav>
      </div>

      <div className="border-t border-border px-4 py-6 text-center text-xs text-text-muted sm:px-6">
        A demonstration shop. Nothing here can be bought.
      </div>
    </footer>
  );
}
