import { defineResource } from "@/lib/resource";
import custom from "./carts.custom";

export const cartResource = defineResource({
  name: "Cart",
  slug: "carts",
  endpoint: "/api/carts",
  icon: "Database",
  label: { singular: "Cart", plural: "Carts" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "token", label: "Token", sortable: true, searchable: true, onClick: "link" },
      { key: "currency", label: "Currency", sortable: true, searchable: true },
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
    { key: "token", label: "Token", type: "text", required: true },
    { key: "currency", label: "Currency", type: "text", required: true },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Carts",
        endpoint: "/api/carts",
        icon: "Database",
        color: "accent",
      },
    ],
  },
}, custom);
