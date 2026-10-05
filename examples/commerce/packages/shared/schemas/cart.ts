import { z } from "zod";

export const CreateCartSchema = z.object({
  token: z.string().min(1, "Required"),
  currency: z.string().min(1, "Required"),
});

export const UpdateCartSchema = z.object({
  token: z.string().min(1, "Required").optional(),
  currency: z.string().min(1, "Required").optional(),
});

export type CreateCartInput = z.infer<typeof CreateCartSchema>;
export type UpdateCartInput = z.infer<typeof UpdateCartSchema>;
