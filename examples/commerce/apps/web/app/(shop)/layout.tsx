import { getPublicCollections } from "@/lib/collections-public";
import { getPublicPages } from "@/lib/pages-public";
import { getCart } from "@/lib/cart";
import { CartProvider } from "@/components/shop/cart-provider";
import { CartDrawer } from "@/components/shop/cart-drawer";
import { ShopHeader } from "@/components/shop/shop-header";
import { ShopFooter } from "@/components/shop/shop-footer";

// The shop's chrome.
//
// The navigation and the footer links come out of the database on every render,
// so adding a collection in the admin adds it to the header, and adding a page
// adds it to the footer. Neither needs a deploy, which is the difference
// between a storefront and a brochure.
//
// The cart is read here, once, and handed to the provider. Every page under this
// layout then has the basket count without fetching anything, and the drawer
// opens with its contents already there rather than showing a spinner on the
// click that matters most.
export default async function ShopLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [collections, pages, cart] = await Promise.all([
    getPublicCollections({ page_size: 20, sort_by: "title", sort_order: "asc" }),
    getPublicPages({ page_size: 20, sort_by: "title", sort_order: "asc" }),
    getCart(),
  ]);

  const collectionLinks = collections.data.map((c) => ({
    href: `/search/${c.handle}`,
    label: c.title,
  }));
  const pageLinks = pages.data.map((p) => ({
    href: `/${p.handle}`,
    label: p.title,
  }));

  return (
    <CartProvider initialCart={cart}>
      <ShopHeader collections={collectionLinks} />
      <main id="main" className="min-h-screen">
        {children}
      </main>
      <ShopFooter collections={collectionLinks} pages={pageLinks} />
      <CartDrawer />
    </CartProvider>
  );
}
