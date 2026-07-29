package plugin

import "strings"

// savedViewModelGo emits internal/models/saved_view.go.
func savedViewModelGo(ctx Context) string {
	src := `package models

import (
	"time"

	"gorm.io/gorm"

	"{{MODULE}}/internal/ids"
)

// SavedView is one user's saved state for one resource table — the URL query
// string that encodes its filters, sort, search and date range. Storing the
// query verbatim means the plugin never has to understand the table's internals:
// applying a view is just navigating to the resource with this query.
type SavedView struct {
	ID       string ~gorm:"primarykey;size:36" json:"id"~
	UserID   string ~gorm:"size:36;index;not null" json:"user_id"~
	Resource string ~gorm:"size:120;index;not null" json:"resource"~
	Name     string ~gorm:"size:160;not null" json:"name" binding:"required"~

	// Query is the raw URL query string (without the leading '?'), e.g.
	// "sort=name&order=asc&status=active".
	Query string ~gorm:"type:text" json:"query"~

	CreatedAt time.Time ~json:"created_at"~
}

func (s *SavedView) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = ids.New()
	}
	return nil
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

// savedViewHandlerGo emits internal/handlers/saved_view.go.
func savedViewHandlerGo(ctx Context) string {
	src := `package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/respond"
)

// SavedViewHandler serves per-user saved table views. Every query is scoped to
// the caller's user_id, so a user only ever sees or deletes their own views.
type SavedViewHandler struct {
	DB *gorm.DB
}

func NewSavedViewHandler(db *gorm.DB) *SavedViewHandler {
	return &SavedViewHandler{DB: db}
}

func currentUserID(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// List returns the caller's saved views for one resource (?resource=<slug>).
func (h *SavedViewHandler) List(c *gin.Context) {
	uid := currentUserID(c)
	resource := c.Query("resource")

	var views []models.SavedView
	q := h.DB.Where("user_id = ?", uid)
	if resource != "" {
		q = q.Where("resource = ?", resource)
	}
	if err := q.Order("name asc").Find(&views).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	respond.OK(c, views)
}

// Create saves the current table state under a name for the caller.
func (h *SavedViewHandler) Create(c *gin.Context) {
	uid := currentUserID(c)

	var in struct {
		Resource string ~json:"resource" binding:"required"~
		Name     string ~json:"name" binding:"required"~
		Query    string ~json:"query"~
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		respond.BadRequest(c, "resource and name are required")
		return
	}

	view := models.SavedView{
		UserID:   uid,
		Resource: in.Resource,
		Name:     in.Name,
		Query:    in.Query,
	}
	if err := h.DB.Create(&view).Error; err != nil {
		respond.Internal(c, err)
		return
	}
	respond.Created(c, view, "View saved")
}

// Delete removes one of the caller's views. Scoping the delete by user_id means
// a crafted id can't delete someone else's view.
func (h *SavedViewHandler) Delete(c *gin.Context) {
	uid := currentUserID(c)

	res := h.DB.Where("id = ? AND user_id = ?", c.Param("id"), uid).Delete(&models.SavedView{})
	if res.Error != nil {
		respond.Internal(c, res.Error)
		return
	}
	if res.RowsAffected == 0 {
		respond.NotFound(c, "View not found")
		return
	}
	respond.OK(c, gin.H{"id": c.Param("id")}, "View deleted")
}
`
	return strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", ctx.Module)
}

// savedViewsHook emits apps/admin/hooks/use-saved-views.ts.
func savedViewsHook() string {
	return `"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export interface SavedView {
  id: string;
  resource: string;
  name: string;
  query: string;
}

export function useSavedViews(resource: string) {
  return useQuery<SavedView[]>({
    queryKey: ["saved-views", resource],
    queryFn: async () =>
      (await apiClient.get("/api/saved-views", { params: { resource } })).data.data ?? [],
  });
}

export function useCreateSavedView(resource: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: { name: string; query: string }) =>
      (await apiClient.post("/api/saved-views", { resource, ...input })).data,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["saved-views", resource] }),
  });
}

export function useDeleteSavedView(resource: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => (await apiClient.delete("/api/saved-views/" + id)).data,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["saved-views", resource] }),
  });
}
`
}

// savedViewsComponent emits apps/admin/components/saved-views.tsx.
func savedViewsComponent() string {
	return `"use client";

import { useState } from "react";
import { useRouter, usePathname, useSearchParams } from "next/navigation";
import { useSavedViews, useCreateSavedView, useDeleteSavedView } from "@/hooks/use-saved-views";

// A row of saved-view chips above a resource table, plus a "Save view" button.
// A saved view is the table's current URL query string — applying one is a
// navigation, so the table restores filters/sort/search with no extra wiring.
export function SavedViews({ resource }: { resource: string }) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const { data: views } = useSavedViews(resource);
  const create = useCreateSavedView(resource);
  const del = useDeleteSavedView(resource);
  const [naming, setNaming] = useState(false);
  const [name, setName] = useState("");

  // The current table state, minus transient params that shouldn't be part of a
  // saved view (an open create/edit drawer).
  function currentQuery(): string {
    const p = new URLSearchParams(searchParams?.toString() ?? "");
    p.delete("action");
    p.delete("edit");
    return p.toString();
  }

  function apply(view: { query: string }) {
    router.push(view.query ? pathname + "?" + view.query : pathname);
  }

  function save() {
    const trimmed = name.trim();
    if (!trimmed) return;
    create.mutate(
      { name: trimmed, query: currentQuery() },
      {
        onSuccess: () => {
          setName("");
          setNaming(false);
        },
      }
    );
  }

  const hasViews = (views?.length ?? 0) > 0;
  if (!hasViews && !naming) {
    return (
      <div className="flex items-center gap-2 border-b border-border px-4 py-2">
        <button
          onClick={() => setNaming(true)}
          className="text-xs font-medium text-accent hover:underline"
        >
          + Save current view
        </button>
      </div>
    );
  }

  return (
    <div className="flex flex-wrap items-center gap-2 border-b border-border px-4 py-2">
      {views?.map((v) => (
        <span
          key={v.id}
          className="group inline-flex items-center gap-1.5 rounded-full border border-border bg-bg-tertiary px-3 py-1 text-xs text-text-secondary"
        >
          <button onClick={() => apply(v)} className="font-medium hover:text-foreground">
            {v.name}
          </button>
          <button
            onClick={() => del.mutate(v.id)}
            aria-label={"Delete view " + v.name}
            className="text-text-muted opacity-0 transition-opacity group-hover:opacity-100 hover:text-danger"
          >
            ×
          </button>
        </span>
      ))}

      {naming ? (
        <span className="inline-flex items-center gap-1.5">
          <input
            autoFocus
            value={name}
            onChange={(e) => setName(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") save();
              if (e.key === "Escape") setNaming(false);
            }}
            placeholder="View name…"
            className="w-32 rounded-lg border border-border bg-bg-tertiary px-2 py-1 text-xs text-foreground outline-none focus:border-accent"
          />
          <button onClick={save} disabled={create.isPending} className="text-xs font-medium text-accent hover:underline disabled:opacity-50">
            Save
          </button>
          <button onClick={() => setNaming(false)} className="text-xs text-text-muted hover:text-foreground">
            Cancel
          </button>
        </span>
      ) : (
        <button onClick={() => setNaming(true)} className="text-xs font-medium text-accent hover:underline">
          + Save current view
        </button>
      )}
    </div>
  );
}
`
}
