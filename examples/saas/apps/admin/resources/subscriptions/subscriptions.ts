import { defineResource } from "@/lib/resource";
import custom from "./subscriptions.custom";

export const subscriptionResource = defineResource({
  name: "Subscription",
  slug: "subscriptions",
  endpoint: "/api/subscriptions",
  icon: "CreditCard",
  label: { singular: "Subscription", plural: "Subscriptions" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "plan.name", label: "Plan" },
      { key: "status", label: "Status", onClick: "link" },
      { key: "seats", label: "Seats", sortable: true },
      { key: "current_period_end", label: "Current Period End", sortable: true, format: "relative" },
      { key: "canceled_on", label: "Canceled On", sortable: true, format: "relative" },
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
    { key: "plan_id", label: "Plan", type: "relationship-select", required: true, relatedEndpoint: "/api/plans", displayField: "name" },
    { key: "status", label: "Status", type: "select", required: true, options: [{ value: "trialing", label: "Trialing" }, { value: "active", label: "Active" }, { value: "past_due", label: "Past Due" }, { value: "canceled", label: "Canceled" }] },
    { key: "seats", label: "Seats", type: "number", numberKind: "int" },
    { key: "current_period_end", label: "Current Period End", type: "date" },
    { key: "canceled_on", label: "Canceled On", type: "date" },
    { key: "user_id", label: "User", type: "relationship-select", required: true, relatedEndpoint: "/api/users", displayField: "first_name" },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Subscriptions",
        endpoint: "/api/subscriptions",
        icon: "CreditCard",
        color: "accent",
      },
    ],
  },
  // One card per value, counted over the rows the table is showing. Generated
  // from the model: status. Delete a line to drop its cards, or set
  // stats: { countBy: [] } to keep only the four defaults.
  stats: {
    countBy: [
      {
        field: "status",
        label: "Status",
        labels: { "trialing": "Trialing", "active": "Active", "past_due": "Past Due", "canceled": "Canceled" },
        only: ["trialing", "active", "past_due", "canceled"],
      },
    ],
  },
  // The collapsible charts above the table: created over time, and these
  // columns by value. Nothing is requested until somebody opens the panel.
  insights: {
    breakdown: ["status"],
    labels: {
      status: { "trialing": "Trialing", "active": "Active", "past_due": "Past Due", "canceled": "Canceled" },
    },
  },
}, custom);
