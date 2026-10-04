import { defineResource } from "@/lib/resource";
import custom from "./loans.custom";

export const loanResource = defineResource({
  name: "Loan",
  slug: "loans",
  endpoint: "/api/loans",
  icon: "Database",
  label: { singular: "Loan", plural: "Loans" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "book.name", label: "Book" },
      { key: "borrower.name", label: "Borrower" },
      { key: "borrowed_at", label: "Borrowed At", sortable: true, format: "relative", onClick: "link" },
      { key: "due_at", label: "Due At", sortable: true, format: "relative" },
      { key: "returned_at", label: "Returned At", sortable: true, format: "relative" },
      { key: "status", label: "Status" },
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
    { key: "book_id", label: "Book", type: "relationship-select", required: true, relatedEndpoint: "/api/books", displayField: "name" },
    { key: "borrower_id", label: "Borrower", type: "relationship-select", required: true, relatedEndpoint: "/api/users", displayField: "name" },
    { key: "borrowed_at", label: "Borrowed At", type: "datetime" },
    { key: "due_at", label: "Due At", type: "datetime" },
    { key: "returned_at", label: "Returned At", type: "datetime" },
    { key: "status", label: "Status", type: "select", required: true, options: [{ value: "out", label: "On loan" }, { value: "returned", label: "Returned" }, { value: "overdue", label: "Overdue" }] },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Loans",
        endpoint: "/api/loans",
        icon: "Database",
        color: "accent",
      },
    ],
  },
}, custom);
