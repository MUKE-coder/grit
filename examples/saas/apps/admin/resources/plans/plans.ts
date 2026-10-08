import { defineResource } from "@/lib/resource";
import custom from "./plans.custom";

export const planResource = defineResource({
  name: "Plan",
  slug: "plans",
  endpoint: "/api/plans",
  icon: "Gem",
  label: { singular: "Plan", plural: "Plans" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "name", label: "Name", sortable: true, searchable: true, onClick: "link" },
      { key: "slug", label: "Slug", sortable: true, searchable: true },
      { key: "price", label: "Price", sortable: true, format: "money" },
      { key: "interval", label: "Interval" },
      { key: "seats", label: "Seats", sortable: true },
      { key: "blurb", label: "Blurb", searchable: true },
      { key: "active", label: "Active", format: "boolean" },
      { key: "created_at", label: "Created", sortable: true, format: "relative" },
      // grit:cols:auto-end
    ],
    filters: [
    { key: "active", label: "Active", type: "boolean" },
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
    { key: "price", label: "Price", type: "money" },
    { key: "interval", label: "Interval", type: "select", required: true, options: [{ value: "monthly", label: "Monthly" }, { value: "yearly", label: "Yearly" }] },
    { key: "seats", label: "Seats", type: "number", numberKind: "int" },
    { key: "blurb", label: "Blurb", type: "textarea" },
    { key: "active", label: "Active", type: "toggle" },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Plans",
        endpoint: "/api/plans",
        icon: "Gem",
        color: "accent",
      },
    ],
  },
  // One card per value, counted over the rows the table is showing. Generated
  // from the model: interval, active. Delete a line to drop its cards, or set
  // stats: { countBy: [] } to keep only the four defaults.
  stats: {
    countBy: [
      {
        field: "interval",
        label: "Interval",
        labels: { "monthly": "Monthly", "yearly": "Yearly" },
        only: ["monthly", "yearly"],
      },
      {
        field: "active",
        labels: { "true": "Active", "1": "Active", "false": "Inactive", "0": "Inactive" },
        only: ["true", "false"],
        icon: "ToggleLeft",
      },
    ],
  },
  // The collapsible charts above the table: created over time, and these
  // columns by value. Nothing is requested until somebody opens the panel.
  insights: {
    breakdown: ["interval", "active"],
    labels: {
      interval: { "monthly": "Monthly", "yearly": "Yearly" },
      active: { "true": "Active", "1": "Active", "false": "Inactive", "0": "Inactive" },
    },
  },
}, custom);
