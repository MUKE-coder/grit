"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, Search, ShoppingBag, X } from "lucide-react";

import { useCart } from "./cart-provider";

interface NavLink {
  href: string;
  label: string;
}

export function ShopHeader({ collections }: { collections: NavLink[] }) {
  const [mobileOpen, setMobileOpen] = useState(false);
  const pathname = usePathname();
  const { cart, setOpen } = useCart();

  const links: NavLink[] = [{ href: "/search", label: "All" }, ...collections];

  return (
    <header className="sticky top-0 z-50 border-b border-border/60 bg-background/85 backdrop-blur-lg">
      {/* A skip link, because the collections list sits between the top of the
          page and the products on every single page. */}
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-3 focus:z-10 focus:rounded-lg focus:bg-accent focus:px-4 focus:py-2 focus:text-sm focus:font-medium focus:text-accent-fg"
      >
        Skip to content
      </a>

      <div className="mx-auto flex h-16 max-w-7xl items-center gap-4 px-4 sm:px-6">
        <Link href="/" className="flex shrink-0 items-center gap-2.5">
          <span className="flex h-8 w-8 items-center justify-center rounded-lg border border-accent/20 bg-accent/15 font-mono text-sm font-bold text-accent">
            A
          </span>
          <span className="text-lg font-bold tracking-tight">Aura</span>
        </Link>

        <nav aria-label="Collections" className="hidden md:block">
          <ul className="flex items-center gap-5">
            {links.map((link) => {
              const active = pathname === link.href;
              return (
                <li key={link.href}>
                  <Link
                    href={link.href}
                    aria-current={active ? "page" : undefined}
                    className={
                      active
                        ? "text-sm font-medium text-foreground"
                        : "text-sm text-text-secondary transition-colors hover:text-foreground"
                    }
                  >
                    {link.label}
                  </Link>
                </li>
              );
            })}
          </ul>
        </nav>

        <div className="ml-auto flex items-center gap-2">
          <Link
            href="/search"
            aria-label="Search the shop"
            className="flex h-10 w-10 items-center justify-center rounded-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <Search className="h-[18px] w-[18px]" aria-hidden="true" />
          </Link>

          <button
            type="button"
            onClick={() => setOpen(true)}
            // The count is in the accessible name, so it is announced rather
            // than only drawn. "Basket, 3 items" is the whole point of the
            // badge, and a badge alone says nothing to a screen reader.
            aria-label={`Basket, ${cart.count} ${cart.count === 1 ? "item" : "items"}`}
            className="relative flex h-10 w-10 items-center justify-center rounded-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <ShoppingBag className="h-[18px] w-[18px]" aria-hidden="true" />
            {cart.count > 0 && (
              <span
                aria-hidden="true"
                className="absolute -right-0.5 -top-0.5 flex h-[18px] min-w-[18px] items-center justify-center rounded-full bg-accent px-1 font-mono text-[10px] font-bold text-accent-fg"
              >
                {cart.count}
              </span>
            )}
          </button>

          <button
            type="button"
            onClick={() => setMobileOpen((open) => !open)}
            aria-expanded={mobileOpen}
            aria-controls="shop-mobile-nav"
            aria-label={mobileOpen ? "Close the menu" : "Open the menu"}
            className="flex h-10 w-10 items-center justify-center rounded-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground md:hidden focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            {mobileOpen ? (
              <X className="h-[18px] w-[18px]" aria-hidden="true" />
            ) : (
              <Menu className="h-[18px] w-[18px]" aria-hidden="true" />
            )}
          </button>
        </div>
      </div>

      {mobileOpen && (
        <nav
          id="shop-mobile-nav"
          aria-label="Collections"
          className="border-t border-border bg-bg-secondary md:hidden"
        >
          <ul className="mx-auto max-w-7xl px-4 py-2 sm:px-6">
            {links.map((link) => (
              <li key={link.href}>
                <Link
                  href={link.href}
                  onClick={() => setMobileOpen(false)}
                  aria-current={pathname === link.href ? "page" : undefined}
                  className="block py-2.5 text-sm text-text-secondary transition-colors hover:text-foreground"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>
      )}
    </header>
  );
}
