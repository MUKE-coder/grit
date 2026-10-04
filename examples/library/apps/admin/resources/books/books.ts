import { defineResource } from "@/lib/resource";
import custom from "./books.custom";

export const bookResource = defineResource({
  name: "Book",
  slug: "books",
  endpoint: "/api/books",
  icon: "Database",
  label: { singular: "Book", plural: "Books" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "title", label: "Title", sortable: true, searchable: true, onClick: "link" },
      { key: "isbn", label: "Isbn", sortable: true, searchable: true },
      { key: "summary", label: "Summary", searchable: true, format: "richtext" },
      { key: "published", label: "Published", sortable: true, format: "relative" },
      { key: "price", label: "Price", sortable: true, format: "money" },
      { key: "cover", label: "Cover", format: "file" },
      { key: "genre", label: "Genre" },
      { key: "author.name", label: "Author" },
      { key: "created_at", label: "Created", sortable: true, format: "relative" },
      // grit:cols:auto-end
    ],
    filters: [
    ],
    defaultSort: { key: "created_at", direction: "desc" },
    searchable: true,
    pageSize: 20,
    // Shown once rows are ticked. Drop "archive" here and the Archived tab
    // goes with it; the model keeps its archived_at either way.
    bulkActions: ["edit", "archive", "restore", "export", "delete"],
  },
  form: {
    fields: [
      // grit:fields:auto-start
    { key: "title", label: "Title", type: "text", required: true },
    { key: "isbn", label: "Isbn", type: "text", required: true },
    { key: "summary", label: "Summary", type: "richtext" },
    { key: "published", label: "Published", type: "date" },
    { key: "price", label: "Price", type: "money" },
    { key: "cover", label: "Cover", type: "file", accepts: ["image"], maxSizeMB: 5 },
    { key: "genre", label: "Genre", type: "select", required: true, options: [{ value: "fiction", label: "Fiction" }, { value: "history", label: "History" }, { value: "science", label: "Science" }] },
    { key: "author_id", label: "Author", type: "relationship-select", required: true, relatedEndpoint: "/api/authors", displayField: "name" },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Books",
        endpoint: "/api/books",
        icon: "Database",
        color: "accent",
      },
    ],
  },
}, custom);
