import { z } from "zod";
import { MoneySchema } from "./money";

export const CreateInvoiceSchema = z.object({
  subscription_id: z.string().uuid("Invalid ID"),
  number: z.string().min(1, "Required"),
  amount: MoneySchema.optional(),
  status: z.enum(["draft", "open", "paid", "void"]),
  issued_on: z.string().nullable(),
  paid_on: z.string().nullable(),
  user_id: z.string().uuid("Invalid ID"),
});

export const UpdateInvoiceSchema = z.object({
  subscription_id: z.string().uuid("Invalid ID").optional(),
  number: z.string().min(1, "Required").optional(),
  amount: MoneySchema.optional(),
  status: z.enum(["draft", "open", "paid", "void"]).optional(),
  issued_on: z.string().nullable(),
  paid_on: z.string().nullable(),
  user_id: z.string().uuid("Invalid ID").optional(),
});

export type CreateInvoiceInput = z.infer<typeof CreateInvoiceSchema>;
export type UpdateInvoiceInput = z.infer<typeof UpdateInvoiceSchema>;
