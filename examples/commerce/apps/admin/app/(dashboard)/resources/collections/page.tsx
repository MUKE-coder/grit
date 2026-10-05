"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { collectionResource } from "@/resources/collections/collections";

export default function CollectionsPage() {
  return <ResourcePage resource={collectionResource} />;
}
