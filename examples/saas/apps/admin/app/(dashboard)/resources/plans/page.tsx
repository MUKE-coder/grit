"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { planResource } from "@/resources/plans/plans";

export default function PlansPage() {
  return <ResourcePage resource={planResource} />;
}
