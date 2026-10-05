import { z } from "zod";

export const CreatePageSchema = z.object({
  title: z.string().min(1, "Required"),
  body: z.string(),
});

export const UpdatePageSchema = z.object({
  title: z.string().min(1, "Required").optional(),
  body: z.string().optional(),
});

export type CreatePageInput = z.infer<typeof CreatePageSchema>;
export type UpdatePageInput = z.infer<typeof UpdatePageSchema>;
