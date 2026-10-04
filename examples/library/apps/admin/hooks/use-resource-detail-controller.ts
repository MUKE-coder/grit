"use client";

import { useCallback, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import type {
  RelatedResource,
  ResourceDefinition,
  ResourceDetailController,
  TableAction,
} from "@/lib/resource";
import { useResourceItem, useDeleteResource } from "@/hooks/use-resource";
import { apiClient } from "@/lib/api-client";
import { resources } from "@/resources";

/*
 * Everything a resource DETAIL page needs except the markup.
 *
 * The list page got this treatment first, and the reasoning is the same: the
 * data was never the hard part, the rest of the page was. Loading one record
 * is a hook call. Working out which other resources point at this one, pulling
 * the inline line-item fields out of the form definition, fetching a
 * server-rendered PDF through the auth interceptor rather than a bare link,
 * and routing back to the list after a delete are not.
 *
 * const c = useResourceDetailController(resource, id)
 * <MyDetail record={c.record} onEdit={c.edit} sections={c.related} />
 */

/** The first present human-readable field, else the resource's own label. */
function titleOf(resource: ResourceDefinition, record: Record<string, unknown> | undefined): string {
  if (!record) return resource.label?.singular ?? resource.name;
  for (const key of ["number", "title", "name", "label", "reference", "slug", "email"]) {
    const value = record[key];
    if (typeof value === "string" && value) return value;
  }
  return resource.label?.singular ?? resource.name;
}

export function useResourceDetailController<T = Record<string, unknown>>(
  resource: ResourceDefinition,
  id: string,
): ResourceDetailController<T> {
  const router = useRouter();
  const { data, isLoading } = useResourceItem<T>(resource.endpoint, id);
  const record = data?.data;

  const [editing, setEditing] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [isPdfBusy, setPdfBusy] = useState(false);

  // The same actions the list honours. Without this the detail page offered
  // Edit and Delete on a resource that had removed both from its table.
  const actions = resource.table.actions ?? ["create", "view", "edit", "delete"];
  const can = (action: TableAction) => actions.includes(action);

  const singular = resource.label?.singular ?? resource.name;
  const { mutate: deleteItem, isPending: isDeleting } = useDeleteResource(resource.endpoint, singular);

  // Through apiClient, not a bare <a href>: that way the auth cookies, the
  // CSRF header and the 401-refresh interceptor all apply. A link gets none
  // of them and silently downloads a login page.
  const downloadPdf = useCallback(async () => {
    setPdfBusy(true);
    try {
      const res = await apiClient.get(resource.endpoint + "/" + id + "/pdf", {
        responseType: "blob",
      });
      const url = URL.createObjectURL(new Blob([res.data], { type: "application/pdf" }));
      window.open(url, "_blank", "noopener,noreferrer");
      // Revoked late: revoking straight away can race the new tab's load.
      setTimeout(() => URL.revokeObjectURL(url), 60000);
    } finally {
      setPdfBusy(false);
    }
  }, [resource.endpoint, id]);

  const lineItemFields = useMemo(
    () => (resource.form?.fields ?? []).filter((f) => f.type === "line-items"),
    [resource.form?.fields],
  );

  // Other resources in the registry that belong to this one, found by a
  // relationship-select whose relatedEndpoint is this endpoint. Anything
  // already rendered inline as line items is skipped, so it is not shown twice.
  const related = useMemo<RelatedResource[]>(() => {
    const inlineEndpoints = new Set(
      (resource.form?.fields ?? [])
        .filter((f) => f.type === "line-items" && f.itemEndpoint)
        .map((f) => f.itemEndpoint as string),
    );
    const out: RelatedResource[] = [];
    for (const other of resources) {
      if (other.slug === resource.slug || other.hidden || inlineEndpoints.has(other.endpoint)) continue;
      const fkField = (other.form?.fields ?? []).find(
        (f) => f.type === "relationship-select" && f.relatedEndpoint === resource.endpoint,
      );
      if (fkField) out.push({ resource: other, fk: fkField.key });
    }
    return out;
  }, [resource]);

  const columns = useMemo(
    () => resource.table.columns.filter((col) => !col.hidden),
    [resource.table.columns],
  );

  const back = useCallback(() => {
    router.push("/resources/" + resource.slug);
  }, [router, resource.slug]);

  return {
    resource,
    id,

    record,
    isLoading,
    notFound: !isLoading && !record,
    title: titleOf(resource, record as Record<string, unknown> | undefined),

    columns,
    lineItemFields,
    related,

    can,
    edit: () => setEditing(true),
    remove: () => setConfirmOpen(true),
    print: () => window.print(),
    downloadPdf,
    back,
    isPdfBusy,
    isDeleting,

    form: { open: editing, item: (record ?? null) as T | null, close: () => setEditing(false) },
    confirmDelete: {
      open: confirmOpen,
      confirm: () => {
        setConfirmOpen(false);
        // Back to the list: staying on the detail page of a record that no
        // longer exists shows "could not be found", which reads as an error
        // rather than as the delete having worked.
        deleteItem(id, { onSuccess: back });
      },
      cancel: () => setConfirmOpen(false),
    },
  };
}
