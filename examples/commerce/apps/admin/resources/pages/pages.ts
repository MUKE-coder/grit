import { defineResource } from "@/lib/resource";
import custom from "./pages.custom";

export const pageResource = defineResource({
  name: "Page",
  slug: "pages",
  endpoint: "/api/pages",
  icon: "FileText",
  label: { singular: "Page", plural: "Pages" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "title", label: "Title", sortable: true, searchable: true, onClick: "link" },
      { key: "handle", label: "Handle", sortable: true, searchable: true },
      { key: "body", label: "Body", searchable: true, format: "richtext" },
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
    { key: "body", label: "Body", type: "richtext" },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Pages",
        endpoint: "/api/pages",
        icon: "FileText",
        color: "accent",
      },
    ],
  },
}, custom);
