"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { cartResource } from "@/resources/carts/carts";

export default function CartsPage() {
  return <ResourcePage resource={cartResource} />;
}
