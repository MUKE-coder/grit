import { usersResource } from "./users/users";
import { blogsResource } from "./blogs/blogs";
import { planResource } from "./plans/plans";
import { subscriptionResource } from "./subscriptions/subscriptions";
import { invoiceResource } from "./invoices/invoices";
import { usageRecordResource } from "./usage-records/usage-records";
// grit:resources

import type { ResourceDefinition } from "@/lib/resource";

export const resources: ResourceDefinition[] = [
  usersResource,
  blogsResource,
  planResource,
  subscriptionResource,
  invoiceResource,
  usageRecordResource,
  // grit:resource-list
];

export function getResource(slug: string): ResourceDefinition | undefined {
  return resources.find((r) => r.slug === slug);
}

export function getResourceByEndpoint(endpoint: string): ResourceDefinition | undefined {
  return resources.find((r) => r.endpoint === endpoint);
}
