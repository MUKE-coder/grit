import { defineResource } from "@/lib/resource";
import custom from "./collections.custom";

export const collectionResource = defineResource({
  name: "Collection",
  slug: "collections",
  endpoint: "/api/collections",
  icon: "Database",
  label: { singular: "Collection", plural: "Collections" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "title", label: "Title", sortable: true, searchable: true, onClick: "link" },
      { key: "handle", label: "Handle", sortable: true, searchable: true },
      { key: "description", label: "Description", searchable: true },
      { key: "image", label: "Image", format: "file" },
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
    { key: "description", label: "Description", type: "textarea" },
    { key: "image", label: "Image", type: "file", accepts: ["image"], maxSizeMB: 5 },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Collections",
        endpoint: "/api/collections",
        icon: "Database",
        color: "accent",
      },
    ],
  },
}, custom);
