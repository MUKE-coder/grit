"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { usageRecordResource } from "@/resources/usage-records/usage-records";

export default function UsageRecordsPage() {
  return <ResourcePage resource={usageRecordResource} />;
}
