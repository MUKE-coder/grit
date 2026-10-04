"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { authorResource } from "@/resources/authors/authors";

export default function AuthorsPage() {
  return <ResourcePage resource={authorResource} />;
}
