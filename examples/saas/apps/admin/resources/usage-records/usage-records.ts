import { defineResource } from "@/lib/resource";
import custom from "./usage-records.custom";

export const usageRecordResource = defineResource({
  name: "UsageRecord",
  slug: "usage-records",
  endpoint: "/api/usage_records",
  icon: "Database",
  label: { singular: "Usage Record", plural: "Usage Records" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "subscription.name", label: "Subscription" },
      { key: "metric", label: "Metric", sortable: true, searchable: true, onClick: "link" },
      { key: "quantity", label: "Quantity", sortable: true },
      { key: "recorded_on", label: "Recorded On", sortable: true, format: "relative" },
      { key: "user.name", label: "User" },
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
    { key: "subscription_id", label: "Subscription", type: "relationship-select", required: true, relatedEndpoint: "/api/subscriptions", displayField: "plan_id" },
    { key: "metric", label: "Metric", type: "text", required: true },
    { key: "quantity", label: "Quantity", type: "number", numberKind: "int" },
    { key: "recorded_on", label: "Recorded On", type: "date" },
    { key: "user_id", label: "User", type: "relationship-select", required: true, relatedEndpoint: "/api/users", displayField: "first_name" },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total UsageRecords",
        endpoint: "/api/usage_records",
        icon: "Database",
        color: "accent",
      },
    ],
  },
}, custom);
