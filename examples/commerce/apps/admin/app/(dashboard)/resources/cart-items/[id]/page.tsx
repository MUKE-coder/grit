"use client";

import { use } from "react";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { cartItemResource } from "@/resources/cart-items/cart-items";

export default function CartItemsDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ResourceDetailPage resource={cartItemResource} id={id} />;
}
