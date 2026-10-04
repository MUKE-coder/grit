"use client";

import { useCallback, useState } from "react";
import type { CustomBulkAction } from "@/lib/resource";

/**
 * Which of a resource list's dialogs are open, and what each is about: the
 * create and edit form, the delete, bulk delete and bulk archive confirms, the
 * bulk editor, a custom action waiting for confirmation, and the importer.
 *
 * State only. What happens on confirm is the controller's business.
 */
export function useResourceDialogs<T>() {
  const [formOpen, setFormOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<T | null>(null);
  // Starting values for the next create. Used by "add a child here" in the tree
  // view, and by anything else that opens a form already scoped to a parent.
  const [formDefaults, setFormDefaults] = useState<Record<string, unknown> | undefined>(undefined);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [bulkDeleteOpen, setBulkDeleteOpen] = useState(false);
  const [bulkArchiveOpen, setBulkArchiveOpen] = useState(false);
  const [bulkEditOpen, setBulkEditOpen] = useState(false);
  const [pendingCustom, setPendingCustom] = useState<CustomBulkAction<T> | null>(null);
  const [importOpen, setImportOpen] = useState(false);

  const openForm = useCallback((item: T | null, defaults?: Record<string, unknown>) => {
    setFormDefaults(defaults);
    setEditingItem(item);
    setFormOpen(true);
  }, []);

  const closeForm = useCallback(() => {
    setFormOpen(false);
    setEditingItem(null);
  }, []);

  return {
    formOpen,
    editingItem,
    formDefaults,
    setFormDefaults,
    openForm,
    closeForm,
    deletingId,
    setDeletingId,
    bulkDeleteOpen,
    setBulkDeleteOpen,
    bulkArchiveOpen,
    setBulkArchiveOpen,
    bulkEditOpen,
    setBulkEditOpen,
    pendingCustom,
    setPendingCustom,
    importOpen,
    setImportOpen,
  };
}
