"use client";

import { useEffect, useRef } from "react";
import Link from "next/link";
import { Minus, Plus, ShoppingBag, Trash2, X } from "lucide-react";

import { Price } from "./price";
import { useCart } from "./cart-provider";

// The basket, as a drawer.
//
// Hand-written rather than reached for from a primitive library, which is the
// convention in this project, and so the three things a dialog has to do are
// here explicitly:
//
//   * Escape closes it
//   * focus moves into it when it opens and back to the trigger when it closes
//   * Tab cannot leave it while it is open
//
// Those are the parts a hand-rolled dialog usually skips, and skipping them
// fails exactly the people who are not in the room when it is demonstrated.
export function CartDrawer() {
  const { cart, open, setOpen, pending, error, setQuantity, remove } = useCart();
  const panel = useRef<HTMLDivElement>(null);
  const closeButton = useRef<HTMLButtonElement>(null);
  const returnTo = useRef<HTMLElement | null>(null);

  // Remember what had focus, move it in, and put it back on close.
  useEffect(() => {
    if (!open) return;
    returnTo.current = document.activeElement as HTMLElement | null;
    closeButton.current?.focus();
    return () => returnTo.current?.focus();
  }, [open]);

  // Escape to close, and a focus trap while it is open.
  useEffect(() => {
    if (!open) return;

    function onKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        event.preventDefault();
        setOpen(false);
        return;
      }
      if (event.key !== "Tab" || !panel.current) return;

      const focusable = panel.current.querySelectorAll<HTMLElement>(
        'button:not([disabled]), a[href], input:not([disabled]), [tabindex]:not([tabindex="-1"])',
      );
      if (focusable.length === 0) return;
      const first = focusable[0];
      const last = focusable[focusable.length - 1];

      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault();
        last.focus();
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault();
        first.focus();
      }
    }

    document.addEventListener("keydown", onKeyDown);
    // The page behind must not scroll while the drawer is over it.
    const previous = document.body.style.overflow;
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKeyDown);
      document.body.style.overflow = previous;
    };
  }, [open, setOpen]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-[60] flex justify-end">
      {/* The backdrop closes it. aria-hidden and no tab stop: the close
          button and Escape are the accessible ways out, and a focusable
          backdrop is a tab stop with no name. */}
      <button
        type="button"
        aria-hidden="true"
        tabIndex={-1}
        onClick={() => setOpen(false)}
        className="absolute inset-0 cursor-default bg-background/70 backdrop-blur-sm"
      />

      <div
        ref={panel}
        role="dialog"
        aria-modal="true"
        aria-labelledby="cart-heading"
        className="relative flex h-full w-full max-w-md flex-col border-l border-border bg-bg-secondary shadow-2xl"
      >
        <header className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 id="cart-heading" className="flex items-center gap-2 text-base font-semibold">
            <ShoppingBag className="h-4 w-4 text-accent" aria-hidden="true" />
            Your basket
            <span className="text-sm font-normal text-text-secondary">
              ({cart.count} {cart.count === 1 ? "item" : "items"})
            </span>
          </h2>
          <button
            ref={closeButton}
            type="button"
            onClick={() => setOpen(false)}
            aria-label="Close the basket"
            className="flex h-9 w-9 items-center justify-center rounded-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <X className="h-4 w-4" aria-hidden="true" />
          </button>
        </header>

        {/* Failures are announced, because the thing that went wrong happened
            somewhere the shopper may not be looking. */}
        <div role="status" aria-live="polite" className="sr-only">
          {pending ? "Updating your basket" : ""}
        </div>
        {error && (
          <p
            role="alert"
            className="border-b border-danger/30 bg-danger/10 px-5 py-3 text-sm text-danger"
          >
            {error}
          </p>
        )}

        {cart.lines.length === 0 ? (
          <div className="flex flex-1 flex-col items-center justify-center gap-3 px-6 text-center">
            <ShoppingBag className="h-10 w-10 text-text-muted" aria-hidden="true" />
            <p className="text-sm text-text-secondary">Your basket is empty.</p>
            <Link
              href="/search"
              onClick={() => setOpen(false)}
              className="text-sm font-medium text-accent hover:underline"
            >
              Browse everything
            </Link>
          </div>
        ) : (
          <>
            <ul className="flex-1 divide-y divide-border overflow-y-auto">
              {cart.lines.map((line) => (
                <li key={line.id} className="flex gap-4 p-5">
                  <Link
                    href={line.product_path}
                    onClick={() => setOpen(false)}
                    className="shrink-0 overflow-hidden rounded-lg border border-border bg-bg-tertiary"
                  >
                    {line.image_url ? (
                      <img
                        src={line.image_url}
                        alt={line.product_title}
                        width={72}
                        height={88}
                        loading="lazy"
                        className="h-[88px] w-[72px] object-cover"
                      />
                    ) : (
                      <span className="block h-[88px] w-[72px]" />
                    )}
                  </Link>

                  <div className="min-w-0 flex-1">
                    <Link
                      href={line.product_path}
                      onClick={() => setOpen(false)}
                      className="block truncate text-sm font-medium hover:text-accent"
                    >
                      {line.product_title}
                    </Link>
                    {line.variant_label && (
                      <p className="mt-0.5 text-xs text-text-secondary">
                        {line.variant_label}
                      </p>
                    )}
                    <Price
                      value={line.unit_price}
                      className="mt-1 block font-mono text-xs text-text-muted"
                    />

                    <div className="mt-3 flex items-center gap-2">
                      <div className="flex items-center rounded-lg border border-border">
                        <button
                          type="button"
                          disabled={pending}
                          onClick={() => setQuantity(line.id, line.quantity - 1)}
                          aria-label={`Reduce the quantity of ${line.product_title}`}
                          className="flex h-8 w-8 items-center justify-center rounded-l-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                        >
                          <Minus className="h-3.5 w-3.5" aria-hidden="true" />
                        </button>
                        <span className="w-9 text-center font-mono text-sm" aria-live="off">
                          {line.quantity}
                        </span>
                        <button
                          type="button"
                          disabled={pending}
                          onClick={() => setQuantity(line.id, line.quantity + 1)}
                          aria-label={`Increase the quantity of ${line.product_title}`}
                          className="flex h-8 w-8 items-center justify-center rounded-r-lg text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                        >
                          <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                        </button>
                      </div>

                      <button
                        type="button"
                        disabled={pending}
                        onClick={() => remove(line.id)}
                        aria-label={`Remove ${line.product_title} from the basket`}
                        className="flex h-8 w-8 items-center justify-center rounded-lg text-text-muted transition-colors hover:bg-danger/10 hover:text-danger disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                      >
                        <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                      </button>

                      <Price
                        value={line.line_total}
                        className="ml-auto font-mono text-sm font-medium"
                      />
                    </div>
                  </div>
                </li>
              ))}
            </ul>

            <footer className="border-t border-border p-5">
              <div className="flex items-baseline justify-between">
                <span className="text-sm text-text-secondary">Subtotal</span>
                <Price value={cart.subtotal} className="font-mono text-lg font-semibold" />
              </div>
              <p className="mt-1 text-xs text-text-muted">
                Shipping and tax are worked out at checkout.
              </p>
              <button
                type="button"
                disabled={pending}
                className="mt-4 w-full rounded-lg bg-accent px-4 py-3 text-sm font-semibold text-accent-fg transition-colors hover:bg-accent-hover disabled:opacity-60 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
              >
                Checkout
              </button>
              <p className="mt-2 text-center text-xs text-text-muted">
                A demonstration shop: checkout is not wired to a payment provider.
              </p>
            </footer>
          </>
        )}
      </div>
    </div>
  );
}
