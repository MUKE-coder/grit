import { usersResource } from "./users/users";
import { blogsResource } from "./blogs/blogs";
import { authorResource } from "./authors/authors";
import { bookResource } from "./books/books";
import { loanResource } from "./loans/loans";
// grit:resources

import type { ResourceDefinition } from "@/lib/resource";

export const resources: ResourceDefinition[] = [
  usersResource,
  blogsResource,
  authorResource,
  bookResource,
  loanResource,
  // grit:resource-list
];

export function getResource(slug: string): ResourceDefinition | undefined {
  return resources.find((r) => r.slug === slug);
}

export function getResourceByEndpoint(endpoint: string): ResourceDefinition | undefined {
  return resources.find((r) => r.endpoint === endpoint);
}
