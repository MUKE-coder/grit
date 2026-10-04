import type { ResourceCustomisation } from "@/lib/resource";
import type { User } from "@repo/shared/types";

/**
 * Customisations for the User resource.
 *
 * This file is yours. The generator creates it once and never writes to it
 * again, so anything you put here survives grit generate and grit sync.
 *
 * Rows are typed as User, so row.id and friends autocomplete and a renamed
 * column fails at compile time rather than at runtime.
 *
 * Three things you can do:
 *
 *   1. Override a single cell, keeping everything else
 *
 *      columns: {
 *        status: { cell: (row) => <StatusPill value={String(row.status)} /> },
 *      }
 *
 *   2. Replace the table, keeping the header, toolbar, filters and pagination.
 *      The props are the same ones DataTable takes, so you can also wrap it:
 *
 *      components: {
 *        Table: (props) => <MyTemplateTable rows={props.data} onSort={props.onSort} />,
 *      }
 *
 *   3. Replace the whole page. Call useResourceController(resource) inside to
 *      get rows, sorting, paging, filters, selection and the CRUD actions:
 *
 *      components: { Page: MyProductsPage }
 *
 * See /docs/admin/custom-pages for the full list of props.
 */
const custom: ResourceCustomisation<User> = {};

export default custom;
