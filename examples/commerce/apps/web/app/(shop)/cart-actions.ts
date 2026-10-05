"use server";

import { revalidatePath } from "next/cache";

import {
  addToCart,
  getCart,
  removeFromCart,
  setCartQuantity,
  type Cart,
} from "@/lib/cart";

// The cart's four operations, as server actions.
//
// A server action rather than a route handler, because the cart token lives in
// an HttpOnly cookie that only the server can read, and an action is the one
// way a client component can ask the server to do something without the shop
// having to build and secure its own HTTP endpoint for it.
//
// Each one returns the whole cart, so a click updates the badge, the drawer and
// the subtotal from one round trip. Returning only "ok" would need a second
// request to redraw, on the interaction shoppers repeat most.
//
// They return errors rather than throwing them. An action that throws shows the
// error boundary, which for "that size just sold out" replaces the entire page
// with an apology. The caller gets a message and puts it next to the button.

export type CartResult =
  | { ok: true; cart: Cart }
  | { ok: false; message: string };

/** What a failure is allowed to say.
 *
 * Mapped from the status, not forwarded from the API. A message passed straight
 * through can carry anything the server said, and the server is entitled to say
 * things a shopper should not read. */
function failure(err: unknown): CartResult {
  const text = err instanceof Error ? err.message : "";
  if (text.includes("409")) {
    return { ok: false, message: "That is not available at the moment." };
  }
  if (text.includes("404")) {
    return { ok: false, message: "That line is no longer in your basket." };
  }
  return { ok: false, message: "Something went wrong. Please try again." };
}

export async function addItemAction(
  productId: string,
  variantId: string,
  quantity = 1,
): Promise<CartResult> {
  try {
    const cart = await addToCart(productId, variantId, quantity);
    // The pages that show stock or a price. Not the whole app: the cart itself
    // comes back in this response, so nothing has to re-render to show it.
    revalidatePath("/product/[handle]", "page");
    return { ok: true, cart };
  } catch (err) {
    return failure(err);
  }
}

export async function setQuantityAction(
  lineId: string,
  quantity: number,
): Promise<CartResult> {
  try {
    return { ok: true, cart: await setCartQuantity(lineId, quantity) };
  } catch (err) {
    return failure(err);
  }
}

export async function removeItemAction(lineId: string): Promise<CartResult> {
  try {
    return { ok: true, cart: await removeFromCart(lineId) };
  } catch (err) {
    return failure(err);
  }
}

/** The current cart, for the drawer to load when it first opens. */
export async function getCartAction(): Promise<Cart> {
  return getCart();
}
