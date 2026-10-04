import { defineResource } from "@/lib/resource";
import custom from "./authors.custom";

export const authorResource = defineResource({
  name: "Author",
  slug: "authors",
  endpoint: "/api/authors",
  icon: "Database",
  label: { singular: "Author", plural: "Authors" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "name", label: "Name", sortable: true, searchable: true, onClick: "link" },
      { key: "bio", label: "Bio", searchable: true },
      { key: "website", label: "Website", sortable: true, searchable: true, format: "link" },
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
    { key: "name", label: "Name", type: "text", required: true },
    { key: "bio", label: "Bio", type: "textarea" },
    { key: "website", label: "Website", type: "url" },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Authors",
        endpoint: "/api/authors",
        icon: "Database",
        color: "accent",
      },
    ],
  },
}, custom);
