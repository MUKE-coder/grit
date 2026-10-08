"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { subscriptionResource } from "@/resources/subscriptions/subscriptions";

export default function SubscriptionsPage() {
  return <ResourcePage resource={subscriptionResource} />;
}
