"use client";

import { use } from "react";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { subscriptionResource } from "@/resources/subscriptions/subscriptions";

export default function SubscriptionsDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ResourceDetailPage resource={subscriptionResource} id={id} />;
}
