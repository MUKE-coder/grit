"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { bookResource } from "@/resources/books/books";

export default function BooksPage() {
  return <ResourcePage resource={bookResource} />;
}
