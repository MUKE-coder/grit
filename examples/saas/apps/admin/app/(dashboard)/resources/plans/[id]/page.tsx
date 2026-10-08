"use client";

import { use } from "react";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { planResource } from "@/resources/plans/plans";

export default function PlansDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ResourceDetailPage resource={planResource} id={id} />;
}
