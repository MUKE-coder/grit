import { z } from "zod";
import { UrlSchema } from "./field-formats";

export const CreateAuthorSchema = z.object({
  name: z.string().min(1, "Required"),
  bio: z.string().optional(),
  website: UrlSchema.or(z.literal("")).optional(),
});

export const UpdateAuthorSchema = z.object({
  name: z.string().min(1, "Required").optional(),
  bio: z.string().optional(),
  website: UrlSchema.or(z.literal("")).optional(),
});

export type CreateAuthorInput = z.infer<typeof CreateAuthorSchema>;
export type UpdateAuthorInput = z.infer<typeof UpdateAuthorSchema>;
