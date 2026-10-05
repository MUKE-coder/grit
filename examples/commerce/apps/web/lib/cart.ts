import { cookies } from "next/headers";

import type { Money } from "@repo/shared/types";

// The cart, from the storefront's server.
//
// Every call here runs on the Next.js server, never in the browser. The browser
// talks to a server action; the action keeps the cart token in an HttpOnly
// cookie and talks to the Go API with the publishable key. Three things fall out
// of that, all of them good:
//
//   * no script on the shop's pages can read the token, so an injected script
//     cannot read or empty somebody's cart
//   * the API needs no CORS allowance for the browser and no CSRF token, because
//     no browser ever calls it directly
//   * the cart survives a reload and a new tab without any client state at all
//
// Hand-written, like lib/product-variants-public.ts: Grit's --public generator
// emits read-only clients, and this is the write side.

const API_URL = (
  process.env.API_INTERNAL_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:8080"
).replace(/\/+$/, "");

const API_KEY = process.env.NEXT_PUBLIC_API_KEY ?? "";

/** The cookie the token lives in. */
const COOKIE = "shop_cart";

/** Thirty days, which is longer than anyone's shopping trip and short enough
 *  that an abandoned cart does not haunt the database forever. */
const COOKIE_MAX_AGE = 60 * 60 * 24 * 30;

export interface CartLine {
  id: string;
  product_id: string;
  product_title: string;
  product_path: string;
  variant_id: string;
  variant_label: string;
  image_url: string;
  quantity: number;
  unit_price: Money;
  line_total: Money;
}

export interface Cart {
  token: string;
  currency: string;
  lines: CartLine[];
  /** Items, not lines: three of one thing is three. */
  count: number;
  subtotal: Money;
}

const EMPTY: Cart = {
  token: "",
  currency: "USD",
  lines: [],
  count: 0,
  subtotal: { amount: 0, currency: "USD" },
};

type Method = "GET" | "POST" | "PATCH" | "DELETE";

/**
 * One request to the cart API, carrying the token and storing any new one.
 *
 * `persist` is false for reads. The API answers a read with a fresh token when
 * the request had none, which is correct for it and wrong to write down here: a
 * crawler that loads the homepage would get a Set-Cookie and a row in `carts`.
 * A cart is created when something is added to it.
 *
 * `cache: "no-store"` because the answer is this shopper's own. The Go side
 * says `Cache-Control: private` for the same reason.
 */
async function call(
  method: Method,
  path: string,
  body?: unknown,
  persist = true,
): Promise<Cart> {
  const jar = await cookies();
  const token = jar.get(COOKIE)?.value ?? "";

  const headers: Record<string, string> = {};
  if (API_KEY) headers["X-API-Key"] = API_KEY;
  if (token) headers["X-Cart-Token"] = token;
  if (body !== undefined) headers["Content-Type"] = "application/json";

  let res: Response;
  try {
    res = await fetch(API_URL + "/api/v1/shop" + path, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
      cache: "no-store",
    });
  } catch {
    // The API is unreachable, as it is while a production build runs in CI.
    // An empty cart renders; it is not an error worth failing a page over.
    return EMPTY;
  }

  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(
      "the cart API answered " + res.status + (text ? ": " + text : ""),
    );
  }

  const payload = (await res.json()) as { data: Cart };
  const cart = payload.data ?? EMPTY;

  if (persist && cart.token && cart.token !== token) {
    jar.set(COOKIE, cart.token, {
      httpOnly: true,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      path: "/",
      maxAge: COOKIE_MAX_AGE,
    });
  }

  return cart;
}

/** The current cart, without creating one. */
export async function getCart(): Promise<Cart> {
  const jar = await cookies();
  if (!jar.get(COOKIE)?.value) return EMPTY;
  return call("GET", "/cart", undefined, false);
}

/** Add a quantity of a variant, or of a plain product. */
export async function addToCart(
  productId: string,
  variantId: string,
  quantity = 1,
): Promise<Cart> {
  return call("POST", "/cart/items", {
    product_id: productId,
    variant_id: variantId,
    quantity,
  });
}

/** Set a line's quantity. Zero removes the line. */
export async function setCartQuantity(
  lineId: string,
  quantity: number,
): Promise<Cart> {
  return call("PATCH", "/cart/items/" + encodeURIComponent(lineId), {
    quantity,
  });
}

/** Take a line out of the cart. */
export async function removeFromCart(lineId: string): Promise<Cart> {
  return call("DELETE", "/cart/items/" + encodeURIComponent(lineId));
}
