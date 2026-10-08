import { z } from "zod";
import { MoneySchema } from "./money";

export const CreatePlanSchema = z.object({
  name: z.string().min(1, "Required"),
  price: MoneySchema.optional(),
  interval: z.enum(["monthly", "yearly"]),
  seats: z.number().int().optional(),
  blurb: z.string().optional(),
  active: z.boolean().optional(),
});

export const UpdatePlanSchema = z.object({
  name: z.string().min(1, "Required").optional(),
  price: MoneySchema.optional(),
  interval: z.enum(["monthly", "yearly"]).optional(),
  seats: z.number().int().optional(),
  blurb: z.string().optional(),
  active: z.boolean().optional(),
});

export type CreatePlanInput = z.infer<typeof CreatePlanSchema>;
export type UpdatePlanInput = z.infer<typeof UpdatePlanSchema>;
