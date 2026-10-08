import { z } from "zod";

export const CreateSubscriptionSchema = z.object({
  plan_id: z.string().uuid("Invalid ID"),
  status: z.enum(["trialing", "active", "past_due", "canceled"]),
  seats: z.number().int().optional(),
  current_period_end: z.string().nullable(),
  canceled_on: z.string().nullable(),
  user_id: z.string().uuid("Invalid ID"),
});

export const UpdateSubscriptionSchema = z.object({
  plan_id: z.string().uuid("Invalid ID").optional(),
  status: z.enum(["trialing", "active", "past_due", "canceled"]).optional(),
  seats: z.number().int().optional(),
  current_period_end: z.string().nullable(),
  canceled_on: z.string().nullable(),
  user_id: z.string().uuid("Invalid ID").optional(),
});

export type CreateSubscriptionInput = z.infer<typeof CreateSubscriptionSchema>;
export type UpdateSubscriptionInput = z.infer<typeof UpdateSubscriptionSchema>;
