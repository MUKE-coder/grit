import { usersResource } from "./users/users";
import { blogsResource } from "./blogs/blogs";
import { collectionResource } from "./collections/collections";
import { productResource } from "./products/products";
import { optionResource } from "./options/options";
import { pageResource } from "./pages/pages";
import { cartResource } from "./carts/carts";
import { cartItemResource } from "./cart-items/cart-items";
// grit:resources

import type { ResourceDefinition } from "@/lib/resource";

export const resources: ResourceDefinition[] = [
  usersResource,
  blogsResource,
  collectionResource,
  productResource,
  optionResource,
  pageResource,
  cartResource,
  cartItemResource,
  // grit:resource-list
];

export function getResource(slug: string): ResourceDefinition | undefined {
  return resources.find((r) => r.slug === slug);
}

export function getResourceByEndpoint(endpoint: string): ResourceDefinition | undefined {
  return resources.find((r) => r.endpoint === endpoint);
}
