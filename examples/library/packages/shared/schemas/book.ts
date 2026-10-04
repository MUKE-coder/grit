import { z } from "zod";
import { FileRefSchema } from "./file-ref";
import { MoneySchema } from "./money";

export const CreateBookSchema = z.object({
  title: z.string().min(1, "Required"),
  isbn: z.string().min(1, "Required"),
  summary: z.string(),
  published: z.string().nullable(),
  price: MoneySchema.optional(),
  cover: FileRefSchema.nullable(),
  genre: z.enum(["fiction", "history", "science"]),
  author_id: z.string().uuid("Invalid ID"),
});

export const UpdateBookSchema = z.object({
  title: z.string().min(1, "Required").optional(),
  isbn: z.string().min(1, "Required").optional(),
  summary: z.string().optional(),
  published: z.string().nullable(),
  price: MoneySchema.optional(),
  cover: FileRefSchema.nullable(),
  genre: z.enum(["fiction", "history", "science"]).optional(),
  author_id: z.string().uuid("Invalid ID").optional(),
});

export type CreateBookInput = z.infer<typeof CreateBookSchema>;
export type UpdateBookInput = z.infer<typeof UpdateBookSchema>;
