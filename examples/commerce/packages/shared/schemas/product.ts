import { z } from "zod";
import { FileRefSchema } from "./file-ref";
import { MoneySchema } from "./money";

export const CreateProductSchema = z.object({
  title: z.string().min(1, "Required"),
  description: z.string(),
  price: MoneySchema.optional(),
  featured_image: FileRefSchema.nullable(),
  images: z.array(FileRefSchema).default([]),
  available: z.boolean().optional(),
  collection_id: z.string().uuid("Invalid ID"),
});

export const UpdateProductSchema = z.object({
  title: z.string().min(1, "Required").optional(),
  description: z.string().optional(),
  price: MoneySchema.optional(),
  featured_image: FileRefSchema.nullable(),
  images: z.array(FileRefSchema).default([]).optional(),
  available: z.boolean().optional(),
  collection_id: z.string().uuid("Invalid ID").optional(),
});

export type CreateProductInput = z.infer<typeof CreateProductSchema>;
export type UpdateProductInput = z.infer<typeof UpdateProductSchema>;
