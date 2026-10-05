import { z } from "zod";
import { FileRefSchema } from "./file-ref";

export const CreateCollectionSchema = z.object({
  title: z.string().min(1, "Required"),
  description: z.string().optional(),
  image: FileRefSchema.nullable(),
});

export const UpdateCollectionSchema = z.object({
  title: z.string().min(1, "Required").optional(),
  description: z.string().optional(),
  image: FileRefSchema.nullable(),
});

export type CreateCollectionInput = z.infer<typeof CreateCollectionSchema>;
export type UpdateCollectionInput = z.infer<typeof UpdateCollectionSchema>;
