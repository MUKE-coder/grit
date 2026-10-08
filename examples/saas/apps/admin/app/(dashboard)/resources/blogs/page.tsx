"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { blogsResource } from "@/resources/blogs/blogs";

export default function BlogsPage() {
  return <ResourcePage resource={blogsResource} />;
}
