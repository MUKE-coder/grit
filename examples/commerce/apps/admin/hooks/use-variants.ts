"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";

/**
 * Data access for options and variants.
 *
 * The endpoints are built from a resource slug rather than hard-coded, because
 * the same matrix editor serves every resource that offers variants. The option
 * endpoints have no slug in them at all: options are shop-wide, and nesting
 * them under a product would suggest each product owns its colours.
 */

export interface AdminOptionValue {
  id: string;
  option_id: string;
  label: string;
  slug: string;
  swatch?: string;
  /** Added to the base price, and only where the OPTION affects price. */
  price_delta: number;
  position: number;
}

export interface AdminOption {
  id: string;
  name: string;
  slug: string;
  /** How a storefront draws it: "swatch", "size" or "select". */
  kind: string;
  affects_price: boolean;
  position: number;
  values?: AdminOptionValue[];
}

export interface AdminVariant {
  id: string;
  sku: string;
  /** null when the price is resolved rather than pinned. */
  price_override: number | null;
  /** Resolved by the server: the override, or the base plus the deltas. */
  price: number;
  stock: number;
  active: boolean;
  position: number;
  /**
   * Photographs of THIS combination, which is why they live on the
   * variant and not on an option value: the picture of the red one is a
   * picture of red-in-XL, while the value's own image is the swatch.
   */
  images?: Array<{ url: string; name?: string }> | null;
  option_values?: AdminOptionValue[];
}

export interface VariantMatrixData {
  options: AdminOption[];
  variants: AdminVariant[];
}

/** What a PATCH to one variant may carry. */
export interface VariantPatch {
  sku?: string;
  stock?: number;
  active?: boolean;
  price_override?: number;
  /**
   * Removes an override. A missing price_override cannot mean "clear it",
   * because that is also what a partial update sends for a field it is not
   * touching, so the server takes a separate flag.
   */
  clear_price?: boolean;
}

export const optionsKey = ["variant-options"];

export function matrixKey(slug: string, id: string) {
  return ["variant-matrix", slug, id];
}

/**
 * The path segment the per-variant update endpoint lives under: "Product"
 * becomes "product-variants".
 *
 * Not nested under the record as /<plural>/:id/variants/:variant, because gin
 * routes a static segment beside a parameter at the same position by panicking
 * at boot, and a bare /variants/:id collides the moment a second resource
 * offers variants. Derived from the resource NAME rather than its slug, because
 * depluralising "products" back to "product" on the client is a guess and
 * kebab-casing the Go name is not.
 */
export function variantPathFor(resourceName: string): string {
  return resourceName.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase() + "-variants";
}

export function useOptions() {
  return useQuery({
    queryKey: optionsKey,
    queryFn: async () => {
      const res = await apiClient.get("/api/options");
      return (res.data.data ?? []) as AdminOption[];
    },
  });
}

export function useVariantMatrix(slug: string, id: string) {
  return useQuery({
    queryKey: matrixKey(slug, id),
    enabled: Boolean(slug) && Boolean(id),
    queryFn: async () => {
      const res = await apiClient.get("/api/" + slug + "/" + id + "/variants");
      const data = res.data.data ?? {};
      return {
        options: data.options ?? [],
        variants: data.variants ?? [],
      } as VariantMatrixData;
    },
  });
}

export function useCreateOption() {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (input: { name: string; kind: string; affects_price: boolean }) => {
      const res = await apiClient.post("/api/options", input);
      return res.data.data as AdminOption;
    },
    successMessage: (option) => option.name + " added",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: optionsKey });
    },
  });
}

/**
 * Deletes an option and its values.
 *
 * The server refuses while anything is built on it and says which way, so there
 * is no pre-flight check here. A client that guessed would be guessing about
 * every other resource in the shop too.
 */
export function useDeleteOption() {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (optionID: string) => {
      await apiClient.delete("/api/options/" + optionID);
    },
    successMessage: "Option deleted",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: optionsKey });
    },
  });
}

export function useCreateOptionValue() {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (input: {
      optionID: string;
      label: string;
      swatch?: string;
      price_delta?: number;
    }) => {
      const { optionID, ...body } = input;
      const res = await apiClient.post("/api/options/" + optionID + "/values", body);
      return res.data.data as AdminOptionValue;
    },
    successMessage: (value) => value.label + " added",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: optionsKey });
    },
  });
}

export function useDeleteOptionValue() {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (valueID: string) => {
      await apiClient.delete("/api/option-values/" + valueID);
    },
    successMessage: "Value deleted",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: optionsKey });
    },
  });
}

/**
 * Sets which options a record offers, in order.
 *
 * Changing the set clears the existing matrix on the server, so the caller is
 * expected to have asked first. The response says how many rows went, which is
 * what the toast reports rather than a flat "saved" that hides the damage.
 */
export function useSetResourceOptions(slug: string, id: string) {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (input: { optionIDs: string[]; valueIDs: string[] }) => {
      const res = await apiClient.put("/api/" + slug + "/" + id + "/options", {
        option_ids: input.optionIDs,
        // Which values of those axes this record offers. An axis with none of
        // its values listed here is offered whole, which is what every record
        // did before it could be narrowed.
        value_ids: input.valueIDs,
      });
      return {
        cleared: (res.data.data?.variants_cleared ?? 0) as number,
        message: (res.data.message ?? "Options set") as string,
      };
    },
    successMessage: (result) =>
      result.cleared > 0
        ? result.message + ", " + result.cleared + " existing combinations cleared"
        : result.message,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: matrixKey(slug, id) });
    },
  });
}

export function useGenerateMatrix(slug: string, id: string) {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (limit?: number) => {
      const res = await apiClient.post(
        "/api/" + slug + "/" + id + "/variants/generate",
        limit ? { limit } : {},
      );
      return (res.data.data?.created ?? 0) as number;
    },
    successMessage: (created) =>
      created === 0 ? "Nothing to add, every combination exists" : created + " combinations added",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: matrixKey(slug, id) });
    },
  });
}

/**
 * Saves every edited row in one go.
 *
 * Sequential rather than parallel, and one toast rather than one per row. A
 * matrix edit is a single act of intent even when it touched nine rows, and
 * nine toasts for it is a notification tray nobody reads.
 */
export function useSaveVariants(resourceName: string, slug: string, id: string) {
  const qc = useQueryClient();
  const path = variantPathFor(resourceName);
  return useToastedMutation({
    mutationFn: async (edits: Array<{ id: string; patch: VariantPatch }>) => {
      for (const edit of edits) {
        await apiClient.patch("/api/" + path + "/" + edit.id, edit.patch);
      }
      return edits.length;
    },
    successMessage: (count) => count + (count === 1 ? " variant saved" : " variants saved"),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: matrixKey(slug, id) });
    },
  });
}
