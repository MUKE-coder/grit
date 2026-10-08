import { z } from "zod";

export const CreateUsageRecordSchema = z.object({
  subscription_id: z.string().uuid("Invalid ID"),
  metric: z.string().min(1, "Required"),
  quantity: z.number().int().optional(),
  recorded_on: z.string().nullable(),
  user_id: z.string().uuid("Invalid ID"),
});

export const UpdateUsageRecordSchema = z.object({
  subscription_id: z.string().uuid("Invalid ID").optional(),
  metric: z.string().min(1, "Required").optional(),
  quantity: z.number().int().optional(),
  recorded_on: z.string().nullable(),
  user_id: z.string().uuid("Invalid ID").optional(),
});

export type CreateUsageRecordInput = z.infer<typeof CreateUsageRecordSchema>;
export type UpdateUsageRecordInput = z.infer<typeof UpdateUsageRecordSchema>;
