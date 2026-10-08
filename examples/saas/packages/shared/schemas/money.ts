import { z } from "zod";

export const MoneySchema = z.object({
  amount: z.number().int(),
  currency: z.string().length(3),
});

export type Money = z.infer<typeof MoneySchema>;
