"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  useTransition,
} from "react";

import type { Cart } from "@/lib/cart";
import {
  addItemAction,
  removeItemAction,
  setQuantityAction,
  type CartResult,
} from "@/app/(shop)/cart-actions";

// The cart's client state, which is one thing: the last cart the server sent.
//
// There is no local arithmetic here. Adding something does not increment a
// counter and then reconcile; it calls the server, and the server's answer
// replaces what we had. That costs a round trip per click and buys the thing
// shoppers care about: the number next to the basket is the number of things in
// the basket, always, including after a reload, in a second tab, or when a
// variant sold out between the page loading and the click.

const EMPTY: Cart = {
  token: "",
  currency: "USD",
  lines: [],
  count: 0,
  subtotal: { amount: 0, currency: "USD" },
};

interface CartState {
  cart: Cart;
  /** True while a mutation is in flight, for disabling the controls. */
  pending: boolean;
  /** The last failure, or null. Shown next to whatever caused it. */
  error: string | null;
  open: boolean;
  setOpen: (open: boolean) => void;
  add: (productId: string, variantId: string, quantity?: number) => void;
  setQuantity: (lineId: string, quantity: number) => void;
  remove: (lineId: string) => void;
}

const CartContext = createContext<CartState | null>(null);

export function CartProvider({
  initialCart,
  children,
}: {
  initialCart: Cart;
  children: React.ReactNode;
}) {
  const [cart, setCart] = useState<Cart>(initialCart ?? EMPTY);
  const [error, setError] = useState<string | null>(null);
  const [open, setOpen] = useState(false);
  const [pending, startTransition] = useTransition();

  // One place that applies a result, so every operation handles failure the
  // same way and none of them can forget to clear the previous error.
  const apply = useCallback(
    (run: () => Promise<CartResult>, openOnSuccess: boolean) => {
      setError(null);
      startTransition(async () => {
        const result = await run();
        if (result.ok) {
          setCart(result.cart);
          if (openOnSuccess) setOpen(true);
          return;
        }
        setError(result.message);
      });
    },
    [],
  );

  const value = useMemo<CartState>(
    () => ({
      cart,
      pending,
      error,
      open,
      setOpen,
      // Adding opens the drawer, because a shopper who clicked "add" wants to
      // see that it worked. Changing a quantity does not, because the drawer is
      // already open: that is where the controls are.
      add: (productId, variantId, quantity = 1) =>
        apply(() => addItemAction(productId, variantId, quantity), true),
      setQuantity: (lineId, quantity) =>
        apply(() => setQuantityAction(lineId, quantity), false),
      remove: (lineId) => apply(() => removeItemAction(lineId), false),
    }),
    [cart, pending, error, open, apply],
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
}

export function useCart(): CartState {
  const ctx = useContext(CartContext);
  if (!ctx) {
    throw new Error("useCart must be used inside a CartProvider");
  }
  return ctx;
}
