import { defineResource } from "@/lib/resource";
import custom from "./invoices.custom";

export const invoiceResource = defineResource({
  name: "Invoice",
  slug: "invoices",
  endpoint: "/api/invoices",
  icon: "Receipt",
  label: { singular: "Invoice", plural: "Invoices" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "subscription.name", label: "Subscription" },
      { key: "number", label: "Number", sortable: true, searchable: true, onClick: "link" },
      { key: "amount", label: "Amount", sortable: true, format: "money" },
      { key: "status", label: "Status" },
      { key: "issued_on", label: "Issued On", sortable: true, format: "relative" },
      { key: "paid_on", label: "Paid On", sortable: true, format: "relative" },
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
    { key: "number", label: "Number", type: "text", required: true },
    { key: "amount", label: "Amount", type: "money" },
    { key: "status", label: "Status", type: "select", required: true, options: [{ value: "draft", label: "Draft" }, { value: "open", label: "Open" }, { value: "paid", label: "Paid" }, { value: "void", label: "Void" }] },
    { key: "issued_on", label: "Issued On", type: "date" },
    { key: "paid_on", label: "Paid On", type: "date" },
    { key: "user_id", label: "User", type: "relationship-select", required: true, relatedEndpoint: "/api/users", displayField: "first_name" },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Invoices",
        endpoint: "/api/invoices",
        icon: "Receipt",
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
        labels: { "draft": "Draft", "open": "Open", "paid": "Paid", "void": "Void" },
        only: ["draft", "open", "paid", "void"],
      },
    ],
  },
  // The collapsible charts above the table: created over time, and these
  // columns by value. Nothing is requested until somebody opens the panel.
  insights: {
    breakdown: ["status"],
    labels: {
      status: { "draft": "Draft", "open": "Open", "paid": "Paid", "void": "Void" },
    },
  },
}, custom);
