"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { cartItemResource } from "@/resources/cart-items/cart-items";

export default function CartItemsPage() {
  return <ResourcePage resource={cartItemResource} />;
}
