package plugin

import "strings"

// tsSource renders an emitted TypeScript template: ~ becomes a backtick, so
// the template can use template literals inside a Go raw string.
func tsSource(src string) string {
	return strings.ReplaceAll(src, "~", "`")
}

// kanbanHook emits hooks/use-boards.ts.
//
// The optimistic update in useMoveCard is why the board feels instant. It
// also has to agree with what the server will do, which is the subtle part:
// applyMove here and RankForCard on the server must put the card in the same
// place, or every drag visibly jumps when the refetch lands.
func kanbanHook() string {
	return tsSource(`"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export interface BoardCard {
  id: string;
  board_id: string;
  column_id: string;
  title: string;
  description: string;
  rank: string;
  assignee_id?: string | null;
  due_date?: string | null;
  created_at: string;
}

export interface BoardColumn {
  id: string;
  board_id: string;
  name: string;
  color: string;
  rank: string;
  wip_limit: number;
  cards?: BoardCard[];
}

export interface Board {
  id: string;
  name: string;
  description: string;
  created_by: string;
  columns?: BoardColumn[];
  created_at: string;
}

export interface MoveCardInput {
  cardId: string;
  columnId: string;
  /** The card the dragged card now sits below. Empty means the top. */
  prevId: string;
  /** The card it now sits above. Empty means the bottom. */
  nextId: string;
}

const boardKey = (id: string) => ["board", id];

export function useBoards() {
  return useQuery<Board[]>({
    queryKey: ["boards"],
    queryFn: async () => (await apiClient.get("/api/boards")).data.data ?? [],
  });
}

export function useBoard(id: string) {
  return useQuery<Board>({
    queryKey: boardKey(id),
    queryFn: async () => (await apiClient.get("/api/boards/" + id)).data.data,
    enabled: Boolean(id),
  });
}

export function useCreateBoard() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: { name: string; description?: string }) =>
      (await apiClient.post("/api/boards", input)).data.data as Board,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["boards"] }),
  });
}

export function useDeleteBoard() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => (await apiClient.delete("/api/boards/" + id)).data,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["boards"] }),
  });
}

export function useCreateCard(boardId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: { column_id: string; title: string; description?: string }) =>
      (await apiClient.post("/api/boards/" + boardId + "/cards", input)).data.data as BoardCard,
    onSuccess: () => qc.invalidateQueries({ queryKey: boardKey(boardId) }),
  });
}

export function useUpdateCard(boardId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, ...input }: { id: string; title?: string; description?: string }) =>
      (await apiClient.patch("/api/board-cards/" + id, input)).data.data as BoardCard,
    onSuccess: () => qc.invalidateQueries({ queryKey: boardKey(boardId) }),
  });
}

export function useDeleteCard(boardId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => (await apiClient.delete("/api/board-cards/" + id)).data,
    onSuccess: () => qc.invalidateQueries({ queryKey: boardKey(boardId) }),
  });
}

export function useCreateColumn(boardId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: { name: string; color?: string; wip_limit?: number }) =>
      (await apiClient.post("/api/boards/" + boardId + "/columns", input)).data.data as BoardColumn,
    onSuccess: () => qc.invalidateQueries({ queryKey: boardKey(boardId) }),
  });
}

/**
 * applyMove is the optimistic half of a drag: it rearranges the cached board
 * exactly the way the server is about to rearrange the real one.
 *
 * It works from the same prev/next neighbours the request carries rather than
 * from an index, for the same reason the request does. If the two disagreed,
 * every drag would land in one place and then visibly jump to another when
 * the refetch arrived.
 */
export function applyMove(board: Board, input: MoveCardInput): Board {
  const columns = (board.columns ?? []).map((c) => ({ ...c, cards: [...(c.cards ?? [])] }));

  let moving: BoardCard | undefined;
  for (const col of columns) {
    const at = col.cards.findIndex((c) => c.id === input.cardId);
    if (at >= 0) {
      moving = col.cards[at];
      col.cards.splice(at, 1);
      break;
    }
  }
  if (!moving) return board;

  const dest = columns.find((c) => c.id === input.columnId);
  if (!dest) return board;

  let index = dest.cards.length;
  if (input.prevId) {
    const at = dest.cards.findIndex((c) => c.id === input.prevId);
    if (at >= 0) index = at + 1;
  } else if (input.nextId) {
    const at = dest.cards.findIndex((c) => c.id === input.nextId);
    if (at >= 0) index = at;
    else index = 0;
  } else {
    index = dest.cards.length;
  }

  dest.cards.splice(index, 0, { ...moving, column_id: dest.id });
  return { ...board, columns };
}

export function useMoveCard(boardId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: MoveCardInput) =>
      (
        await apiClient.post("/api/board-cards/" + input.cardId + "/move", {
          column_id: input.columnId,
          prev_id: input.prevId,
          next_id: input.nextId,
        })
      ).data.data as BoardCard,

    onMutate: async (input) => {
      // An in-flight board refetch would otherwise land on top of the
      // optimistic update and snap the card back for a frame.
      await qc.cancelQueries({ queryKey: boardKey(boardId) });
      const previous = qc.getQueryData<Board>(boardKey(boardId));
      if (previous) {
        qc.setQueryData<Board>(boardKey(boardId), applyMove(previous, input));
      }
      return { previous };
    },

    onError: (_err, _input, context) => {
      // A refused move (a full column, say) has to put the card back where
      // the user took it from, not leave the board showing a move that did
      // not happen.
      if (context?.previous) {
        qc.setQueryData(boardKey(boardId), context.previous);
      }
    },

    onSettled: () => qc.invalidateQueries({ queryKey: boardKey(boardId) }),
  });
}
`)
}

// kanbanBoardView emits components/kanban/board-view.tsx.
func kanbanBoardView() string {
	return tsSource(`"use client";

import { useMemo, useState } from "react";
import {
  DndContext,
  DragOverlay,
  KeyboardSensor,
  PointerSensor,
  closestCorners,
  useDroppable,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import type { Announcements, DragEndEvent, DragStartEvent } from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, Trash2 } from "@/lib/icons";
import { CardComposer, ColumnComposer } from "@/components/kanban/card-composer";
import { useDeleteCard, useMoveCard } from "@/hooks/use-boards";
import type { Board, BoardCard, BoardColumn } from "@/hooks/use-boards";

const COLUMN_PREFIX = "column:";

/** Column accent, by the token name stored on the column. */
const ACCENTS: Record<string, string> = {
  info: "bg-info",
  success: "bg-success",
  warning: "bg-warning",
  danger: "bg-danger",
  accent: "bg-accent",
};

function accentFor(color: string): string {
  return ACCENTS[color] ?? "bg-bg-hover";
}

export function BoardView({ board }: { board: Board }) {
  const columns = useMemo(() => board.columns ?? [], [board.columns]);
  const move = useMoveCard(board.id);
  const [activeId, setActiveId] = useState<string | null>(null);

  const sensors = useSensors(
    // Without a distance threshold the pointer sensor swallows clicks, and a
    // card can never be opened, only dragged.
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } }),
    // Space picks a card up, the arrows move it, Space drops it, Escape
    // cancels. This is the whole reason for taking the dependency.
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  const cardsById = useMemo(() => {
    const out = new Map<string, { card: BoardCard; column: BoardColumn }>();
    for (const column of columns) {
      for (const card of column.cards ?? []) out.set(card.id, { card, column });
    }
    return out;
  }, [columns]);

  const activeCard = activeId ? cardsById.get(activeId)?.card : undefined;

  // Screen readers get the column name, which is the only part of a board
  // that carries meaning. dnd-kit's defaults announce coordinates.
  const announcements: Announcements = {
    onDragStart: ({ active }) => {
      const found = cardsById.get(String(active.id));
      return found ? "Picked up " + found.card.title + " from " + found.column.name : undefined;
    },
    onDragOver: ({ active, over }) => {
      const found = cardsById.get(String(active.id));
      if (!found || !over) return undefined;
      const column = columnFromOver(String(over.id));
      return column ? "Over " + column.name : undefined;
    },
    onDragEnd: ({ active, over }) => {
      const found = cardsById.get(String(active.id));
      if (!found) return undefined;
      if (!over) return "Cancelled. " + found.card.title + " stayed in " + found.column.name;
      const column = columnFromOver(String(over.id));
      return column ? "Dropped " + found.card.title + " in " + column.name : undefined;
    },
    onDragCancel: ({ active }) => {
      const found = cardsById.get(String(active.id));
      return found ? "Cancelled. " + found.card.title + " stayed in " + found.column.name : undefined;
    },
  };

  function columnFromOver(overId: string): BoardColumn | undefined {
    if (overId.startsWith(COLUMN_PREFIX)) {
      const id = overId.slice(COLUMN_PREFIX.length);
      return columns.find((c) => c.id === id);
    }
    return cardsById.get(overId)?.column;
  }

  function onDragStart(event: DragStartEvent) {
    setActiveId(String(event.active.id));
  }

  function onDragEnd(event: DragEndEvent) {
    setActiveId(null);

    const { active, over } = event;
    if (!over) return;

    const cardId = String(active.id);
    const overId = String(over.id);
    const from = cardsById.get(cardId)?.column;
    const dest = columnFromOver(overId);
    if (!from || !dest) return;

    const destCards = dest.cards ?? [];
    let prevId = "";
    let nextId = "";

    if (dest.id === from.id) {
      // Same column: work out the resulting order with dnd-kit's own
      // arrayMove, so the neighbours sent to the server are the ones the user
      // actually sees after the drop.
      const ids = destCards.map((c) => c.id);
      const oldIndex = ids.indexOf(cardId);
      const newIndex = overId.startsWith(COLUMN_PREFIX) ? ids.length - 1 : ids.indexOf(overId);
      if (oldIndex < 0 || newIndex < 0 || oldIndex === newIndex) return;

      const after = arrayMove(ids, oldIndex, newIndex);
      const at = after.indexOf(cardId);
      prevId = after[at - 1] ?? "";
      nextId = after[at + 1] ?? "";
    } else {
      // Another column: the moving card is not in this list yet, so the drop
      // target's own index is where it lands.
      const ids = destCards.map((c) => c.id);
      const at = overId.startsWith(COLUMN_PREFIX) ? ids.length : Math.max(0, ids.indexOf(overId));
      prevId = ids[at - 1] ?? "";
      nextId = ids[at] ?? "";
    }

    move.mutate({ cardId, columnId: dest.id, prevId, nextId });
  }

  return (
    <DndContext
      sensors={sensors}
      collisionDetection={closestCorners}
      accessibility={{ announcements }}
      onDragStart={onDragStart}
      onDragEnd={onDragEnd}
      onDragCancel={() => setActiveId(null)}
    >
      <div className="flex gap-4 overflow-x-auto pb-4">
        {columns.map((column) => (
          <Column key={column.id} boardId={board.id} column={column} />
        ))}
        <ColumnComposer boardId={board.id} />
      </div>

      {/* The card under the cursor, rendered outside the scroll container so
          it is not clipped by a column. */}
      <DragOverlay>
        {activeCard ? <CardFace card={activeCard} dragging /> : null}
      </DragOverlay>
    </DndContext>
  );
}

function Column({ boardId, column }: { boardId: string; column: BoardColumn }) {
  const cards = column.cards ?? [];
  const ids = cards.map((c) => c.id);

  // The whole body is droppable, not just the cards, so a card can be dropped
  // into an empty column and below the last card of a full one.
  const { setNodeRef, isOver } = useDroppable({ id: COLUMN_PREFIX + column.id });

  const atLimit = column.wip_limit > 0 && cards.length >= column.wip_limit;

  return (
    <section
      className="flex w-72 shrink-0 flex-col rounded-xl border border-border bg-bg-secondary"
      aria-label={column.name}
    >
      <header className="flex items-center gap-2 border-b border-border px-3 py-2.5">
        <span className={"h-2 w-2 shrink-0 rounded-full " + accentFor(column.color)} aria-hidden="true" />
        <h2 className="flex-1 truncate text-sm font-semibold text-foreground">{column.name}</h2>
        <span
          className={
            "rounded-full px-2 py-0.5 text-xs tabular-nums " +
            (atLimit ? "bg-warning/15 text-warning" : "bg-bg-tertiary text-muted")
          }
          title={column.wip_limit > 0 ? "Work-in-progress limit: " + column.wip_limit : undefined}
        >
          {column.wip_limit > 0 ? cards.length + " / " + column.wip_limit : cards.length}
        </span>
      </header>

      <div
        ref={setNodeRef}
        className={
          "flex min-h-24 flex-1 flex-col gap-2 p-2 transition-colors " +
          (isOver ? "bg-bg-tertiary/60" : "")
        }
      >
        <SortableContext items={ids} strategy={verticalListSortingStrategy}>
          {cards.map((card) => (
            <SortableCard key={card.id} boardId={boardId} card={card} />
          ))}
        </SortableContext>

        {cards.length === 0 ? (
          <p className="px-1 py-6 text-center text-xs text-muted">Nothing here yet</p>
        ) : null}
      </div>

      <CardComposer boardId={boardId} columnId={column.id} disabled={atLimit} />
    </section>
  );
}

function SortableCard({ boardId, card }: { boardId: string; card: BoardCard }) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: card.id,
  });

  return (
    <div
      ref={setNodeRef}
      style={{ transform: CSS.Transform.toString(transform), transition }}
      className={isDragging ? "opacity-40" : ""}
    >
      <CardFace card={card} boardId={boardId} handleProps={{ ...attributes, ...listeners }} />
    </div>
  );
}

function CardFace({
  card,
  boardId,
  handleProps,
  dragging,
}: {
  card: BoardCard;
  boardId?: string;
  handleProps?: Record<string, unknown>;
  dragging?: boolean;
}) {
  const remove = useDeleteCard(boardId ?? "");

  return (
    <article
      className={
        "group flex items-start gap-2 rounded-lg border border-border bg-bg-elevated p-2.5 " +
        (dragging ? "shadow-lg ring-1 ring-accent" : "")
      }
    >
      <button
        type="button"
        {...handleProps}
        // A real button, so it is in the tab order and the keyboard sensor
        // can receive the Space that picks the card up. A div with a
        // listener here is the version keyboard users cannot reach.
        className="mt-0.5 cursor-grab rounded text-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
        aria-label={"Move " + card.title}
      >
        <GripVertical className="h-4 w-4" aria-hidden="true" />
      </button>

      <div className="min-w-0 flex-1">
        <p className="break-words text-sm text-foreground">{card.title}</p>
        {card.description ? (
          <p className="mt-1 line-clamp-2 break-words text-xs text-muted">{card.description}</p>
        ) : null}
      </div>

      {boardId ? (
        <button
          type="button"
          onClick={() => remove.mutate(card.id)}
          className="rounded p-1 text-muted opacity-0 transition-opacity hover:text-danger focus-visible:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-danger group-hover:opacity-100"
          aria-label={"Delete " + card.title}
        >
          <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
        </button>
      ) : null}
    </article>
  );
}
`)
}

// kanbanCardComposer emits components/kanban/card-composer.tsx.
//
// Adding a card is an inline form at the foot of the column rather than a
// modal. A modal would need a focus trap, and a hand-written focus trap is
// the component that works in a demo and strands a keyboard user in practice.
// The inline form needs none of it, and it is the faster interaction anyway.
func kanbanCardComposer() string {
	return tsSource(`"use client";

import { useState } from "react";
import { Plus, X } from "@/lib/icons";
import { buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";
import { useCreateCard, useCreateColumn } from "@/hooks/use-boards";

export function CardComposer({
  boardId,
  columnId,
  disabled,
}: {
  boardId: string;
  columnId: string;
  disabled?: boolean;
}) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState("");
  const create = useCreateCard(boardId);

  function submit(event: React.FormEvent) {
    event.preventDefault();
    const trimmed = title.trim();
    if (!trimmed) return;
    create.mutate(
      { column_id: columnId, title: trimmed },
      {
        onSuccess: () => {
          // Stay open: adding one card is usually adding three.
          setTitle("");
        },
      },
    );
  }

  if (disabled) {
    return (
      <p className="border-t border-border px-3 py-2 text-xs text-warning">
        At its work-in-progress limit
      </p>
    );
  }

  if (!open) {
    return (
      <div className="border-t border-border p-2">
        <button
          type="button"
          onClick={() => setOpen(true)}
          className={buttonClasses({ variant: "ghost", size: "sm", className: "w-full justify-start" })}
        >
          <Plus className="h-3.5 w-3.5" aria-hidden="true" />
          Add a card
        </button>
      </div>
    );
  }

  return (
    <form onSubmit={submit} className="border-t border-border p-2">
      <label htmlFor={"new-card-" + columnId} className="sr-only">
        Card title
      </label>
      <textarea
        id={"new-card-" + columnId}
        autoFocus
        rows={2}
        value={title}
        onChange={(event) => setTitle(event.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            submit(event);
          }
          if (event.key === "Escape") setOpen(false);
        }}
        placeholder="What needs doing?"
        className={inputClasses({ multiline: true, className: "resize-none" })}
      />
      <div className="mt-2 flex items-center gap-2">
        <button type="submit" disabled={create.isPending} className={buttonClasses({ size: "sm" })}>
          {create.isPending ? "Adding" : "Add card"}
        </button>
        <button
          type="button"
          onClick={() => setOpen(false)}
          className={buttonClasses({ variant: "ghost", size: "sm" })}
          aria-label="Cancel adding a card"
        >
          <X className="h-4 w-4" aria-hidden="true" />
        </button>
      </div>
    </form>
  );
}

export function ColumnComposer({ boardId }: { boardId: string }) {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const create = useCreateColumn(boardId);

  function submit(event: React.FormEvent) {
    event.preventDefault();
    const trimmed = name.trim();
    if (!trimmed) return;
    create.mutate(
      { name: trimmed },
      {
        onSuccess: () => {
          setName("");
          setOpen(false);
        },
      },
    );
  }

  if (!open) {
    return (
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="flex h-fit w-72 shrink-0 items-center gap-2 rounded-xl border border-dashed border-border px-3 py-2.5 text-sm text-muted transition-colors hover:border-accent hover:text-foreground"
      >
        <Plus className="h-4 w-4" aria-hidden="true" />
        Add a column
      </button>
    );
  }

  return (
    <form onSubmit={submit} className="h-fit w-72 shrink-0 rounded-xl border border-border bg-bg-secondary p-2">
      <label htmlFor="new-column" className="sr-only">
        Column name
      </label>
      <input
        id="new-column"
        autoFocus
        value={name}
        onChange={(event) => setName(event.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Escape") setOpen(false);
        }}
        placeholder="Column name"
        className={inputClasses()}
      />
      <div className="mt-2 flex items-center gap-2">
        <button type="submit" disabled={create.isPending} className={buttonClasses({ size: "sm" })}>
          {create.isPending ? "Adding" : "Add column"}
        </button>
        <button
          type="button"
          onClick={() => setOpen(false)}
          className={buttonClasses({ variant: "ghost", size: "sm" })}
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
`)
}

// kanbanBoardsPage emits the boards index page.
func kanbanBoardsPage() string {
	return tsSource(`"use client";

import { useState } from "react";
import Link from "next/link";
import { Plus, Trash2 } from "@/lib/icons";
import { buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";
import { useBoards, useCreateBoard, useDeleteBoard } from "@/hooks/use-boards";

export default function BoardsPage() {
  const { data: boards, isLoading } = useBoards();
  const create = useCreateBoard();
  const remove = useDeleteBoard();
  const [name, setName] = useState("");

  function submit(event: React.FormEvent) {
    event.preventDefault();
    const trimmed = name.trim();
    if (!trimmed) return;
    create.mutate({ name: trimmed }, { onSuccess: () => setName("") });
  }

  return (
    <div className="space-y-6 p-6">
      <header>
        <h1 className="text-2xl font-semibold text-foreground">Boards</h1>
        <p className="mt-1 text-sm text-muted">
          Work in columns. Drag a card with the pointer, a finger, or the keyboard.
        </p>
      </header>

      <form onSubmit={submit} className="flex max-w-xl items-center gap-2">
        <label htmlFor="new-board" className="sr-only">
          Board name
        </label>
        <input
          id="new-board"
          value={name}
          onChange={(event) => setName(event.target.value)}
          placeholder="New board name"
          className={inputClasses()}
        />
        <button type="submit" disabled={create.isPending} className={buttonClasses()}>
          <Plus className="h-4 w-4" aria-hidden="true" />
          {create.isPending ? "Creating" : "Create"}
        </button>
      </form>

      {isLoading ? (
        <p className="text-sm text-muted">Loading boards</p>
      ) : !boards || boards.length === 0 ? (
        <p className="rounded-xl border border-dashed border-border p-8 text-center text-sm text-muted">
          No boards yet. Create one above and it starts with To do, In progress and Done.
        </p>
      ) : (
        <ul className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {boards.map((board) => (
            <li
              key={board.id}
              className="group flex items-center gap-3 rounded-xl border border-border bg-bg-secondary p-4 transition-colors hover:border-accent"
            >
              <Link href={"/boards/" + board.id} className="min-w-0 flex-1">
                <span className="block truncate font-medium text-foreground">{board.name}</span>
                {board.description ? (
                  <span className="mt-0.5 block truncate text-xs text-muted">{board.description}</span>
                ) : null}
              </Link>
              <button
                type="button"
                onClick={() => remove.mutate(board.id)}
                className="rounded p-1.5 text-muted opacity-0 transition-opacity hover:text-danger focus-visible:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-danger group-hover:opacity-100"
                aria-label={"Delete " + board.name}
              >
                <Trash2 className="h-4 w-4" aria-hidden="true" />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
`)
}

// kanbanBoardDetailPage emits the single-board page.
//
// Next reads the id from the route params; TanStack reads it from the router,
// so the two shapes differ by one line and NextToTanStack cannot infer it.
func kanbanBoardDetailPage(ctx Context) string {
	if ctx.Frontend == "tanstack" {
		return tsSource(`import { Link, useParams } from "@tanstack/react-router";
import { ArrowLeft } from "@/lib/icons";
import { BoardView } from "@/components/kanban/board-view";
import { useBoard } from "@/hooks/use-boards";

export default function BoardDetailPage() {
  const { id } = useParams({ strict: false }) as { id: string };
  const { data: board, isLoading, error } = useBoard(id);

  if (isLoading) return <p className="p-6 text-sm text-muted">Loading board</p>;
  if (error || !board) return <p className="p-6 text-sm text-danger">That board could not be loaded.</p>;

  return (
    <div className="space-y-5 p-6">
      <header className="flex items-center gap-3">
        <Link to="/boards" className="rounded p-1.5 text-muted hover:text-foreground" aria-label="Back to boards">
          <ArrowLeft className="h-4 w-4" aria-hidden="true" />
        </Link>
        <div>
          <h1 className="text-2xl font-semibold text-foreground">{board.name}</h1>
          {board.description ? <p className="text-sm text-muted">{board.description}</p> : null}
        </div>
      </header>
      <BoardView board={board} />
    </div>
  );
}
`)
	}

	return tsSource(`"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { ArrowLeft } from "@/lib/icons";
import { BoardView } from "@/components/kanban/board-view";
import { useBoard } from "@/hooks/use-boards";

export default function BoardDetailPage() {
  const params = useParams<{ id: string }>();
  const id = params?.id ?? "";
  const { data: board, isLoading, error } = useBoard(id);

  if (isLoading) return <p className="p-6 text-sm text-muted">Loading board</p>;
  if (error || !board) return <p className="p-6 text-sm text-danger">That board could not be loaded.</p>;

  return (
    <div className="space-y-5 p-6">
      <header className="flex items-center gap-3">
        <Link href="/boards" className="rounded p-1.5 text-muted hover:text-foreground" aria-label="Back to boards">
          <ArrowLeft className="h-4 w-4" aria-hidden="true" />
        </Link>
        <div>
          <h1 className="text-2xl font-semibold text-foreground">{board.name}</h1>
          {board.description ? <p className="text-sm text-muted">{board.description}</p> : null}
        </div>
      </header>
      <BoardView board={board} />
    </div>
  );
}
`)
}
