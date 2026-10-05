"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { pageResource } from "@/resources/pages/pages";

export default function PagesPage() {
  return <ResourcePage resource={pageResource} />;
}
