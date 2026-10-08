"use client";

import { use } from "react";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { blogsResource } from "@/resources/blogs/blogs";

export default function BlogsDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ResourceDetailPage resource={blogsResource} id={id} />;
}
