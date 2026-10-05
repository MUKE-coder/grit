"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { productResource } from "@/resources/products/products";

export default function ProductsPage() {
  return <ResourcePage resource={productResource} />;
}
