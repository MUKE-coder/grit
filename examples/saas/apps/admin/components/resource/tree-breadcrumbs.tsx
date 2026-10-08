"use client";

import { useQuery } from "@tanstack/react-query";
import Link from "next/link";
import { ChevronRight } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";

interface Crumb {
  id: string;
  name?: string;
  title?: string;
}

interface TreeBreadcrumbsProps {
  /** The resource endpoint, e.g. "/api/categories". */
  endpoint: string;
  /** The admin route for one record, e.g. "/resources/categories". */
  basePath: string;
  id: string;
}

/**
 * The ancestors of a record, root first, this record last.
 *
 * Renders nothing at all for a root: a single-item breadcrumb is furniture that
 * tells the reader something they can already see in the heading.
 */
export function TreeBreadcrumbs({ endpoint, basePath, id }: TreeBreadcrumbsProps) {
  const { data } = useQuery({
    queryKey: [endpoint, id, "breadcrumbs"],
    queryFn: async () => {
      const res = await apiClient.get<{ data: Crumb[] }>(endpoint + "/" + id + "/breadcrumbs");
      return res.data.data ?? [];
    },
    enabled: Boolean(id),
  });

  const crumbs = data ?? [];
  if (crumbs.length < 2) return null;

  return (
    <nav className="flex items-center gap-1 text-xs text-text-muted" aria-label="Ancestors">
      {crumbs.map((crumb, i) => {
        const label = crumb.name || crumb.title || crumb.id;
        const isLast = i === crumbs.length - 1;
        return (
          <span key={crumb.id} className="flex items-center gap-1">
            {isLast ? (
              <span className="text-foreground">{label}</span>
            ) : (
              <Link href={basePath + "/" + crumb.id} className="hover:text-foreground">
                {label}
              </Link>
            )}
            {!isLast && <ChevronRight className="h-3 w-3" />}
          </span>
        );
      })}
    </nav>
  );
}
