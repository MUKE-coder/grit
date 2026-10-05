import { defineResource } from "@/lib/resource";
import { OptionLibraryPage } from "@/components/variants/option-library";

/**
 * The shared option library, as a sidebar entry.
 *
 * The table and form below are never rendered — components.Page replaces the
 * whole list view — but they are not decoration either. The registry reads
 * endpoint and label to resolve relationships and breadcrumbs, so they describe
 * the resource honestly rather than being left empty.
 */
export const optionResource = defineResource(
  {
    name: "Option",
    slug: "options",
    endpoint: "/api/options",
    icon: "Layers",
    label: { singular: "Option", plural: "Options" },
    adminOnly: true,
    // The stat cards are counts of a resource the server knows how to count,
    // and /admin/dashboard/resource-stats only knows the resources the
    // generator registered. Options were registered here by hand, so asking is
    // two 400s on every visit. There is nothing to count anyway: the page shows
    // the whole library.
    stats: false,
    table: {
      columns: [
        { key: "name", label: "Name", sortable: true, searchable: true },
        { key: "kind", label: "Kind" },
        { key: "affects_price", label: "Changes price", format: "boolean" },
      ],
      searchable: true,
    },
    form: {
      fields: [
        { key: "name", label: "Name", type: "text", required: true },
        {
          key: "kind",
          label: "How a storefront draws it",
          type: "select",
          options: [
            { value: "select", label: "Dropdown" },
            { value: "swatch", label: "Swatch" },
            { value: "size", label: "Size boxes" },
          ],
        },
        { key: "affects_price", label: "Changes the price", type: "toggle" },
      ],
    },
  },
  { components: { Page: OptionLibraryPage } },
);
