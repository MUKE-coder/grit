import { defineResource } from "@/lib/resource";
import custom from "./products.custom";

export const productResource = defineResource({
  name: "Product",
  slug: "products",
  endpoint: "/api/products",
  icon: "Package",
  label: { singular: "Product", plural: "Products" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "title", label: "Title", sortable: true, searchable: true, onClick: "link" },
      { key: "handle", label: "Handle", sortable: true, searchable: true },
      { key: "description", label: "Description", searchable: true, format: "richtext" },
      { key: "price", label: "Price", sortable: true, format: "money" },
      { key: "featured_image", label: "Featured Image", format: "file" },
      { key: "images", label: "Images", format: "files" },
      { key: "available", label: "Available", format: "boolean" },
      { key: "collection.name", label: "Collection" },
      { key: "created_at", label: "Created", sortable: true, format: "relative" },
      // grit:cols:auto-end
    ],
    filters: [
    { key: "available", label: "Available", type: "boolean" },
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
    { key: "description", label: "Description", type: "richtext" },
    { key: "price", label: "Price", type: "money" },
    { key: "featured_image", label: "Featured Image", type: "file", accepts: ["image"], maxSizeMB: 5 },
    { key: "images", label: "Images", type: "files", accepts: ["all"], maxSizeMB: 5, max: 5 },
    { key: "available", label: "Available", type: "toggle" },
    { key: "collection_id", label: "Collection", type: "relationship-select", required: true, relatedEndpoint: "/api/collections", displayField: "name" },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Products",
        endpoint: "/api/products",
        icon: "Package",
        color: "accent",
      },
    ],
  },
}, custom);
