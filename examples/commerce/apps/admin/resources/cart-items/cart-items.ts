import { defineResource } from "@/lib/resource";
import custom from "./cart-items.custom";

export const cartItemResource = defineResource({
  name: "CartItem",
  slug: "cart-items",
  endpoint: "/api/cart_items",
  icon: "Database",
  label: { singular: "Cart Item", plural: "Cart Items" },
  table: {
    columns: [
      // grit:cols:auto-start
      { key: "cart.name", label: "Cart" },
      { key: "product.name", label: "Product" },
      { key: "variant_id", label: "Variant ID", sortable: true, searchable: true, onClick: "link" },
      { key: "quantity", label: "Quantity", sortable: true },
      { key: "unit_price", label: "Unit Price", sortable: true, format: "money" },
      { key: "title", label: "Title", sortable: true, searchable: true },
      { key: "variant_label", label: "Variant Label", sortable: true, searchable: true },
      { key: "image_url", label: "Image URL", sortable: true, searchable: true },
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
    { key: "cart_id", label: "Cart", type: "relationship-select", required: true, relatedEndpoint: "/api/carts", displayField: "name" },
    { key: "product_id", label: "Product", type: "relationship-select", required: true, relatedEndpoint: "/api/products", displayField: "name" },
    { key: "variant_id", label: "Variant ID", type: "text", required: true },
    { key: "quantity", label: "Quantity", type: "number", numberKind: "int" },
    { key: "unit_price", label: "Unit Price", type: "money" },
    { key: "title", label: "Title", type: "text", required: true },
    { key: "variant_label", label: "Variant Label", type: "text", required: true },
    { key: "image_url", label: "Image URL", type: "text", required: true },
      // grit:fields:auto-end
    ],
  },
  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total CartItems",
        endpoint: "/api/cart_items",
        icon: "Database",
        color: "accent",
      },
    ],
  },
}, custom);
