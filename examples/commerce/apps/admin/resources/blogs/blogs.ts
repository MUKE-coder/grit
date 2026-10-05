import { defineResource } from "@/lib/resource";
import custom from "./blogs.custom";

export const blogsResource = defineResource({
  name: "Blog",
  slug: "blogs",
  endpoint: "/api/admin/blogs",
  icon: "Newspaper",
  label: { singular: "Blog", plural: "Blogs" },
  // A blog is long-form writing, so its form is a page of its own rather than
  // a sheet, and the editor has the width to work in.
  formView: "page",

  table: {
    columns: [
      // grit:cols:auto-start
      // v3.31.5: dropped the raw UUID column. Title + status + author
      // already identify a blog row clearly; the ID lives in the URL when
      // you open the detail view.
      { key: "title", label: "Title", sortable: true, searchable: true },
      { key: "slug", label: "Slug" },
      { key: "image", label: "Image", format: "image" },
      {
        key: "published",
        label: "Status",
        format: "badge",
        badge: {
          true: { color: "success", label: "Published" },
          false: { color: "muted", label: "Draft" },
        },
      },
      { key: "published_at", label: "Published At", format: "relative", sortable: true },
      { key: "created_at", label: "Created", format: "relative", sortable: true },
      // grit:cols:auto-end
    ],
    filters: [
      {
        key: "published",
        label: "Status",
        type: "select",
        options: [
          { label: "Published", value: "true" },
          { label: "Draft", value: "false" },
        ],
      },
    ],
    searchable: true,
    searchPlaceholder: "Search blogs by title...",
    actions: ["create", "view", "edit", "delete"],
    // No "archive": the scaffold's own models predate archived_at. A
    // generated resource gets the column and the full set.
    bulkActions: ["edit", "export", "delete"],
    defaultSort: { key: "created_at", direction: "desc" },
    pageSize: 20,
  },

  form: {
    layout: "single",
    fields: [
      // grit:fields:auto-start
      {
        key: "title",
        label: "Title",
        type: "text",
        required: true,
        placeholder: "Enter blog title",
      },
      {
        key: "excerpt",
        label: "Excerpt",
        type: "textarea",
        placeholder: "Brief summary of the blog post",
      },
      {
        key: "content",
        label: "Content",
        type: "richtext",
      },
      {
        key: "image",
        label: "Cover Image",
        type: "image",
      },
      {
        key: "published",
        label: "Published",
        type: "toggle",
      },
      // grit:fields:auto-end
    ],
  },
}, custom);
