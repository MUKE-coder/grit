package scaffold

// The admin's tree view.
//
// Generated once per project (not per resource) and driven by
// ResourceDefinition.tree, the same way DataTable is driven by table.columns.
// A resource generated with --tree declares tree: true and gets this for free.
//
// Native HTML5 drag and drop, deliberately. dnd-kit or react-dnd would be
// nicer to write against and would put a dependency into every scaffolded admin
// for one screen. The browser has had draggable since IE5, and a tree is the
// one case where the native API is enough: single-item drag, no sorting
// animation, no multi-select.
//
// Three drop targets per row, and this is the part worth understanding, because
// a tree with only one is unusable:
//
//   - the row itself      -> become a child of this node
//   - a thin bar above it -> become its previous sibling
//   - a thin bar below    -> become its next sibling (last row only)
//
// Without the sibling bars there is no way to reorder within a parent, and no
// way to promote a node back to the root once it has a parent. Those are the
// two things somebody actually does with a category tree.

func adminResourceTree() string {
	return `"use client";

import { useCallback, useMemo, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  ChevronDown,
  ChevronRight,
  GripVertical,
  Pencil,
  Plus,
  RefreshCw,
} from "@/lib/icons";
import { apiClient } from "@/lib/api-client";
import { Button } from "@/components/ui/button";
import type { ResourceDefinition } from "@/lib/resource";

/** A node as the tree endpoint returns it: the row, plus its children. */
interface TreeNode {
  id: string;
  name?: string;
  title?: string;
  label?: string;
  parent_id?: string;
  path?: string;
  depth?: number;
  position?: number;
  children?: TreeNode[] | null;
  [key: string]: unknown;
}

/** Where a dragged node is about to land. */
type DropWhere = "inside" | "before" | "after";

interface DropTarget {
  id: string;
  where: DropWhere;
}

interface ResourceTreeProps {
  resource: ResourceDefinition;
  /**
   * Handed the whole node rather than its id, because the page's edit form is
   * populated from the row it is given and a tree node already is that row.
   */
  onEdit?: (node: TreeNode) => void;
  /**
   * Open a create form with the parent already chosen.
   *
   * This is why the controller has createWith(): calling plain create() here
   * would open a form with no parent set, and the new row would be born at the
   * root, which is a button that appears to do one thing and does another.
   */
  onAddChild?: (parentID: string) => void;
}

/** The first of these a node has is what we render as its label. */
function labelOf(node: TreeNode): string {
  return (
    (typeof node.name === "string" && node.name) ||
    (typeof node.title === "string" && node.title) ||
    (typeof node.label === "string" && node.label) ||
    node.id
  );
}

/**
 * A second line under the label, when the row has one worth showing.
 *
 * Slug first because it is short, stable and the thing you actually look for
 * when checking a URL. A description is trimmed hard: two lines of prose in a
 * tree row destroys the scannability the tree exists for.
 */
function subtitleOf(node: TreeNode): string {
  const slug = typeof node.slug === "string" ? node.slug : "";
  if (slug) return slug;
  const description = typeof node.description === "string" ? node.description : "";
  return description.length > 60 ? description.slice(0, 60) + "..." : description;
}

/**
 * The row's own image, if the resource has one.
 *
 * A file field arrives as an object with a url, and a files field as an array
 * of them. Both are worth showing: a category tree with pictures is far faster
 * to scan than one without.
 */
function imageOf(node: TreeNode): string {
  for (const key of ["image", "images", "photo", "thumbnail", "logo", "avatar", "icon"]) {
    const value = node[key];
    if (!value) continue;
    if (typeof value === "string" && value.startsWith("http")) return value;
    const first = Array.isArray(value) ? value[0] : value;
    if (first && typeof first === "object" && typeof (first as { url?: string }).url === "string") {
      return (first as { url: string }).url;
    }
  }
  return "";
}

/** Two letters for the fallback tile, so a row without a picture still reads. */
function initialsOf(node: TreeNode): string {
  return labelOf(node)
    .split(/\s+/)
    .slice(0, 2)
    .map((word) => word.charAt(0).toUpperCase())
    .join("");
}

/** Counts a node and everything under it, for the "N inside" badge. */
function descendantCount(node: TreeNode): number {
  return (node.children ?? []).reduce((total, child) => total + 1 + descendantCount(child), 0);
}

/** Flattens a tree into ids, used to stop a node being dropped inside itself. */
function subtreeIDs(node: TreeNode): string[] {
  const out = [node.id];
  for (const child of node.children ?? []) {
    out.push(...subtreeIDs(child));
  }
  return out;
}

/** Finds a node anywhere in the forest. */
function findNode(nodes: TreeNode[], id: string): TreeNode | null {
  for (const node of nodes) {
    if (node.id === id) return node;
    const found = findNode(node.children ?? [], id);
    if (found) return found;
  }
  return null;
}

/** The siblings of a node, in their rendered order. */
function siblingsOf(nodes: TreeNode[], parentID: string): TreeNode[] {
  if (!parentID) return nodes;
  const parent = findNode(nodes, parentID);
  return parent?.children ?? [];
}

export function ResourceTree({ resource, onEdit, onAddChild }: ResourceTreeProps) {
  const queryClient = useQueryClient();
  const [expanded, setExpanded] = useState<Record<string, boolean>>({});
  const [dragging, setDragging] = useState<string | null>(null);
  const [target, setTarget] = useState<DropTarget | null>(null);

  const treeKey = useMemo(() => [resource.slug, "tree"], [resource.slug]);

  const { data, isLoading, error } = useQuery({
    queryKey: treeKey,
    queryFn: async () => {
      const res = await apiClient.get<{ data: TreeNode[] }>(resource.endpoint + "/tree");
      return res.data.data ?? [];
    },
  });

  const nodes = data ?? [];

  /** Every id in the forest, for expand all. */
  const allIDs = useMemo(() => {
    const out: string[] = [];
    const walk = (list: TreeNode[]) => {
      for (const node of list) {
        out.push(node.id);
        walk(node.children ?? []);
      }
    };
    walk(nodes);
    return out;
  }, [nodes]);

  /** How deep the tree goes, shown in the header so the shape is legible. */
  const maxDepth = useMemo(() => {
    let deepest = 0;
    const walk = (list: TreeNode[], depth: number) => {
      for (const node of list) {
        deepest = Math.max(deepest, depth);
        walk(node.children ?? [], depth + 1);
      }
    };
    walk(nodes, 1);
    return deepest;
  }, [nodes]);

  // Collapsed by default would hide the structure the page exists to show, and
  // expanded by default makes a deep tree unreadable. Two levels is the
  // compromise: you see the shape without scrolling past it.
  const isExpanded = useCallback(
    (node: TreeNode) => expanded[node.id] ?? (node.depth ?? 0) < 1,
    [expanded]
  );

  const invalidate = () => {
    queryClient.invalidateQueries({ queryKey: treeKey });
    // The table view of the same resource is now stale too: a move changed
    // parent_id, and that is a column in it.
    queryClient.invalidateQueries({ queryKey: [resource.slug] });
  };

  const move = useMutation({
    mutationFn: async (vars: { id: string; parentID: string; position: number }) => {
      await apiClient.patch(resource.endpoint + "/" + vars.id + "/move", {
        parent_id: vars.parentID,
        position: vars.position,
      });
    },
    onError: (err: unknown) => {
      // The server refuses a move into a node's own subtree. The UI blocks that
      // too, so reaching here means something raced, or a rule the client does
      // not know about. Either way the message is worth showing verbatim.
      const message =
        (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error
          ?.message ?? "That move was refused.";
      toast.error(message);
      invalidate();
    },
  });

  const reorder = useMutation({
    mutationFn: async (vars: { parentID: string; ids: string[] }) => {
      await apiClient.post(resource.endpoint + "/reorder", {
        parent_id: vars.parentID,
        ids: vars.ids,
      });
    },
    onError: () => toast.error("Could not save the new order."),
  });

  const rebuild = useMutation({
    mutationFn: async () => {
      await apiClient.post(resource.endpoint + "/rebuild-tree", {});
    },
    onSuccess: () => {
      toast.success("Paths rebuilt.");
      invalidate();
    },
    onError: () => toast.error("Could not rebuild the paths."),
  });

  /** True when dropping onto this target would put a node inside itself. */
  const wouldCycle = useCallback(
    (draggedID: string, targetID: string) => {
      const dragged = findNode(nodes, draggedID);
      if (!dragged) return false;
      return subtreeIDs(dragged).includes(targetID);
    },
    [nodes]
  );

  const handleDrop = useCallback(
    async (draggedID: string, drop: DropTarget) => {
      setTarget(null);
      setDragging(null);
      if (draggedID === drop.id) return;

      const dragged = findNode(nodes, draggedID);
      const targetNode = findNode(nodes, drop.id);
      if (!dragged || !targetNode) return;

      // Refused here as well as on the server, because a toast after a failed
      // request is a worse answer than a cursor that says no.
      if (subtreeIDs(dragged).includes(drop.id)) {
        toast.error("A " + resource.name.toLowerCase() + " cannot go inside itself.");
        return;
      }

      if (drop.where === "inside") {
        const children = targetNode.children ?? [];
        await move.mutateAsync({
          id: draggedID,
          parentID: targetNode.id,
          position: children.length,
        });
        // Opened, or the node appears to vanish into a collapsed parent.
        setExpanded((prev) => ({ ...prev, [targetNode.id]: true }));
        invalidate();
        return;
      }

      // Sibling drop. Move first when the parent changes, then write the whole
      // sibling order: Move sets one node's position and deliberately does not
      // shift the others, so the insert is Reorder's job.
      const newParent = targetNode.parent_id ?? "";
      const siblings = siblingsOf(nodes, newParent).filter((s) => s.id !== draggedID);
      const at = siblings.findIndex((s) => s.id === targetNode.id);
      const index = drop.where === "before" ? at : at + 1;
      const ordered = [
        ...siblings.slice(0, index).map((s) => s.id),
        draggedID,
        ...siblings.slice(index).map((s) => s.id),
      ];

      if ((dragged.parent_id ?? "") !== newParent) {
        await move.mutateAsync({ id: draggedID, parentID: newParent, position: index });
      }
      await reorder.mutateAsync({ parentID: newParent, ids: ordered });
      invalidate();
    },
    [nodes, move, reorder, resource.name] // eslint-disable-line react-hooks/exhaustive-deps
  );

  if (isLoading) {
    return (
      <div className="rounded-xl border border-border bg-background-secondary p-6">
        <div className="space-y-2">
          {[0, 1, 2, 3, 4].map((i) => (
            <div
              key={i}
              className="h-[52px] animate-pulse rounded-xl border border-border bg-background-tertiary"
              style={{ marginLeft: (i % 3) * 24 }}
            />
          ))}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-xl border border-border bg-background-secondary p-6 text-sm text-danger">
        Could not load the tree.
      </div>
    );
  }

  return (
    <div className="rounded-xl border border-border bg-background-secondary">
      <div className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-4 py-3">
        <div className="min-w-0">
          <p className="text-sm font-medium text-foreground">
            {nodes.length} top level
            {maxDepth > 1 && (
              <span className="text-text-muted">
                {" "}
                &middot; {maxDepth} levels deep
              </span>
            )}
          </p>
          <p className="text-xs text-text-muted">
            Drag a row onto another to nest it, or between rows to reorder. Drop it at
            the very top to make it a root.
          </p>
        </div>

        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="sm"
            onClick={() =>
              setExpanded(Object.fromEntries(allIDs.map((id) => [id, true])))
            }
          >
            Expand all
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() =>
              setExpanded(Object.fromEntries(allIDs.map((id) => [id, false])))
            }
          >
            Collapse all
          </Button>
        <Button
          variant="ghost"
          size="sm"
          onClick={() => rebuild.mutate()}
          loading={rebuild.isPending}
          // The case for this button: --tree was added to a table that already
          // had rows, so every one of them has an empty path and the tree
          // renders flat. One click fixes it.
          title="Recompute every path and depth from parent_id. Safe to run any time."
        >
          {!rebuild.isPending && <RefreshCw className="h-3.5 w-3.5" />}
          Rebuild paths
        </Button>
        </div>
      </div>

      {/* A drop bar above the first root, which is the only way to promote a
          nested node back to the top level. */}
      <RootDropBar
        active={target?.where === "before" && target.id === (nodes[0]?.id ?? "")}
        onEnter={() => nodes[0] && setTarget({ id: nodes[0].id, where: "before" })}
        onLeave={() => setTarget(null)}
        onDrop={() => nodes[0] && dragging && handleDrop(dragging, { id: nodes[0].id, where: "before" })}
        enabled={Boolean(dragging) && nodes.length > 0}
      />

      <ul className="space-y-1 p-3">
        {nodes.map((node, i) => (
          <TreeRow
            key={node.id}
            node={node}
            depth={0}
            last={i === nodes.length - 1}
            dragging={dragging}
            target={target}
            expanded={isExpanded(node)}
            onToggle={(id) =>
              setExpanded((prev) => ({ ...prev, [id]: !(prev[id] ?? (node.depth ?? 0) < 1) }))
            }
            isExpanded={isExpanded}
            setDragging={setDragging}
            setTarget={setTarget}
            onDrop={handleDrop}
            wouldCycle={wouldCycle}
            onEdit={onEdit}
            onAddChild={onAddChild}
          />
        ))}
      </ul>

      {nodes.length === 0 && (
        <div className="px-4 py-12 text-center">
          <p className="text-sm font-medium text-foreground">Nothing here yet</p>
          <p className="mt-1 text-xs text-text-muted">
            Create one with the New button above, then drag rows onto each other to
            build the hierarchy.
          </p>
        </div>
      )}
    </div>
  );
}

function RootDropBar({
  active,
  enabled,
  onEnter,
  onLeave,
  onDrop,
}: {
  active: boolean;
  enabled: boolean;
  onEnter: () => void;
  onLeave: () => void;
  onDrop: () => void;
}) {
  if (!enabled) return null;
  return (
    <div
      onDragOver={(e) => {
        e.preventDefault();
        onEnter();
      }}
      onDragLeave={onLeave}
      onDrop={(e) => {
        e.preventDefault();
        onDrop();
      }}
      className="relative mx-3 mt-3 h-2"
      aria-hidden
    >
      <span
        className={
          "absolute inset-x-0 top-1/2 h-0.5 -translate-y-1/2 rounded-full transition-all " +
          (active ? "bg-accent" : "bg-border")
        }
      />
      <span
        className={
          "absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 rounded-full px-2 text-[10px] font-medium transition-colors " +
          (active ? "bg-accent text-white" : "bg-background-secondary text-text-muted")
        }
      >
        drop here to make it a root
      </span>
    </div>
  );
}

interface TreeRowProps {
  node: TreeNode;
  /** Rendered depth, 0 for a root. Read from the tree shape rather than from
      the row's own depth column, so a stale column cannot misdraw the view. */
  depth: number;
  last: boolean;
  dragging: string | null;
  target: DropTarget | null;
  expanded: boolean;
  onToggle: (id: string) => void;
  isExpanded: (node: TreeNode) => boolean;
  setDragging: (id: string | null) => void;
  setTarget: (t: DropTarget | null) => void;
  onDrop: (draggedID: string, drop: DropTarget) => void;
  wouldCycle: (draggedID: string, targetID: string) => boolean;
  onEdit?: (node: TreeNode) => void;
  onAddChild?: (parentID: string) => void;
}

function TreeRow({
  node,
  depth,
  last,
  dragging,
  target,
  expanded,
  onToggle,
  isExpanded,
  setDragging,
  setTarget,
  onDrop,
  wouldCycle,
  onEdit,
  onAddChild,
}: TreeRowProps) {
  const children = node.children ?? [];
  const hasChildren = children.length > 0;
  const total = descendantCount(node);
  const subtitle = subtitleOf(node);
  const image = imageOf(node);
  const isDragging = dragging === node.id;
  const forbidden = Boolean(dragging) && dragging !== node.id && wouldCycle(dragging as string, node.id);
  const isInsideTarget = target?.id === node.id && target.where === "inside";
  const isBeforeTarget = target?.id === node.id && target.where === "before";
  const isAfterTarget = target?.id === node.id && target.where === "after";

  return (
    <li className="relative mt-1 first:mt-0">
      <DropBar
        active={isBeforeTarget}
        enabled={Boolean(dragging) && !isDragging}
        onEnter={() => setTarget({ id: node.id, where: "before" })}
        onLeave={() => setTarget(null)}
        onDrop={() => dragging && onDrop(dragging, { id: node.id, where: "before" })}
      />

      {/* The elbow that joins this card to its parent's guide line. Indentation
          alone reads as "further right"; a connector reads as "belongs to". */}
      {depth > 0 && (
        <span
          aria-hidden
          className="pointer-events-none absolute -left-3 top-[26px] h-px w-3 bg-border"
        />
      )}

      <div
        draggable
        onDragStart={(e) => {
          e.stopPropagation();
          setDragging(node.id);
          e.dataTransfer.effectAllowed = "move";
          // Firefox refuses to start a drag without data on the transfer.
          e.dataTransfer.setData("text/plain", node.id);
        }}
        onDragEnd={() => {
          setDragging(null);
          setTarget(null);
        }}
        onDragOver={(e) => {
          if (!dragging || isDragging || forbidden) return;
          e.preventDefault();
          e.stopPropagation();
          setTarget({ id: node.id, where: "inside" });
        }}
        onDragLeave={(e) => {
          e.stopPropagation();
          if (isInsideTarget) setTarget(null);
        }}
        onDrop={(e) => {
          if (!dragging || forbidden) return;
          e.preventDefault();
          e.stopPropagation();
          onDrop(dragging, { id: node.id, where: "inside" });
        }}
        className={
          "group flex items-center gap-2.5 rounded-xl border px-2.5 py-2 transition-all " +
          (isDragging ? "opacity-40 " : "") +
          (forbidden ? "cursor-no-drop opacity-40 " : "") +
          (isInsideTarget
            ? "border-accent bg-accent/10 shadow-sm ring-1 ring-accent "
            : "border-border bg-background hover:border-accent/50 hover:shadow-sm ")
        }
      >
        {/* Visible at rest, not on hover. A drag handle nobody can see is a
            feature nobody finds: the whole panel read as an unstyled list
            because every affordance was hidden until the pointer arrived. */}
        <GripVertical className="h-4 w-4 shrink-0 cursor-grab text-text-muted/50 transition-colors group-hover:text-text-muted" />

        <button
          type="button"
          onClick={() => onToggle(node.id)}
          className={
            "flex h-6 w-6 shrink-0 items-center justify-center rounded-md text-text-muted transition-colors hover:bg-background-tertiary hover:text-foreground " +
            (hasChildren ? "" : "invisible")
          }
          aria-label={expanded ? "Collapse" : "Expand"}
          aria-expanded={hasChildren ? expanded : undefined}
        >
          {expanded ? (
            <ChevronDown className="h-4 w-4" />
          ) : (
            <ChevronRight className="h-4 w-4" />
          )}
        </button>

        {/* The picture if there is one, initials if not. Either way the row has
            a fixed leading block, so labels line up down the whole tree. */}
        {image ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={image}
            alt=""
            className="h-9 w-9 shrink-0 rounded-lg border border-border object-cover"
          />
        ) : (
          <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-border bg-background-tertiary text-[11px] font-semibold text-text-muted">
            {initialsOf(node)}
          </span>
        )}

        <span className="min-w-0 flex-1">
          <span className="flex items-center gap-2">
            <span className="truncate text-sm font-medium text-foreground">
              {labelOf(node)}
            </span>
            {/* Stated, not merely implied by indentation. On a wide screen a
                third-level row sits a long way from its parent, and counting
                pixels is not reading. */}
            <span className="shrink-0 rounded bg-background-tertiary px-1.5 py-0.5 font-mono text-[10px] leading-none text-text-muted">
              L{depth + 1}
            </span>
          </span>
          {subtitle && (
            <span className="block truncate text-xs text-text-muted">{subtitle}</span>
          )}
        </span>

        {hasChildren && (
          <span
            className="shrink-0 rounded-full border border-border px-2 py-0.5 text-[11px] font-medium text-text-muted"
            title={total + " inside, " + children.length + " directly under this one"}
          >
            {total}
          </span>
        )}

        {/* Actions fade in, which is fine: they are secondary, and the row
            already looks like something before the pointer arrives. */}
        <span className="flex shrink-0 items-center gap-1 opacity-0 transition-opacity focus-within:opacity-100 group-hover:opacity-100">
          {onAddChild && (
            <button
              type="button"
              onClick={() => onAddChild(node.id)}
              title={"Add a child under " + labelOf(node)}
              className="rounded-md p-1.5 text-text-muted transition-colors hover:bg-accent/15 hover:text-accent"
            >
              <Plus className="h-4 w-4" />
            </button>
          )}
          {onEdit && (
            <button
              type="button"
              onClick={() => onEdit(node)}
              title={"Edit " + labelOf(node)}
              className="rounded-md p-1.5 text-text-muted transition-colors hover:bg-background-tertiary hover:text-foreground"
            >
              <Pencil className="h-4 w-4" />
            </button>
          )}
        </span>
      </div>

      {expanded && hasChildren && (
        <ul className="ml-6 space-y-1 border-l-2 border-border/70 pl-3">
          {children.map((child, i) => (
            <TreeRow
              key={child.id}
              node={child}
              depth={depth + 1}
              last={i === children.length - 1}
              dragging={dragging}
              target={target}
              expanded={isExpanded(child)}
              onToggle={onToggle}
              isExpanded={isExpanded}
              setDragging={setDragging}
              setTarget={setTarget}
              onDrop={onDrop}
              wouldCycle={wouldCycle}
              onEdit={onEdit}
              onAddChild={onAddChild}
            />
          ))}
        </ul>
      )}

      {/* Only the last row in a list gets an "after" bar. Every other gap is
          already covered by the next row's "before" bar, and two bars in one
          gap means the highlight flickers between them as the pointer moves. */}
      {last && (
        <DropBar
          active={isAfterTarget}
          enabled={Boolean(dragging) && !isDragging}
          onEnter={() => setTarget({ id: node.id, where: "after" })}
          onLeave={() => setTarget(null)}
          onDrop={() => dragging && onDrop(dragging, { id: node.id, where: "after" })}
        />
      )}
    </li>
  );
}

function DropBar({
  active,
  enabled,
  onEnter,
  onLeave,
  onDrop,
}: {
  active: boolean;
  enabled: boolean;
  onEnter: () => void;
  onLeave: () => void;
  onDrop: () => void;
}) {
  if (!enabled) return null;
  return (
    <div
      onDragOver={(e) => {
        e.preventDefault();
        e.stopPropagation();
        onEnter();
      }}
      onDragLeave={(e) => {
        e.stopPropagation();
        onLeave();
      }}
      onDrop={(e) => {
        e.preventDefault();
        e.stopPropagation();
        onDrop();
      }}
      className="relative -my-1 h-2"
      aria-hidden
    >
      <span
        className={
          "absolute inset-x-0 top-1/2 h-0.5 -translate-y-1/2 rounded-full transition-all " +
          (active ? "bg-accent" : "bg-transparent")
        }
      />
      {active && (
        <span className="absolute left-0 top-1/2 h-2 w-2 -translate-y-1/2 rounded-full bg-accent" />
      )}
    </div>
  );
}
`
}

// adminTreeBreadcrumbs renders the ancestors of one record, for a detail page.
//
// Reads the path from the breadcrumbs endpoint rather than walking parents in
// the client: the ids are already in the stored path, so the server answers in
// one query however deep the tree goes.
func adminTreeBreadcrumbs() string {
	return `"use client";

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
`
}

// Exported so the release-testing helper can write these into an existing
// project without a full re-scaffold. The upgrade path uses the same values.
func AdminResourceTree() string    { return adminResourceTree() }
func AdminTreeBreadcrumbs() string { return adminTreeBreadcrumbs() }
func AdminResourcePage() string    { return adminResourcePage() }
func AdminIconMap() string         { return adminIconMap() }
func AdminResourceTypes() string   { return adminResourceTypes() }
