"use client";

import { use } from "react";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { usageRecordResource } from "@/resources/usage-records/usage-records";

export default function UsageRecordsDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ResourceDetailPage resource={usageRecordResource} id={id} />;
}
