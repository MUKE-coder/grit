import { z } from "zod";

export const CreateLoanSchema = z.object({
  book_id: z.string().uuid("Invalid ID"),
  borrower_id: z.string().uuid("Invalid ID"),
  borrowed_at: z.string().nullable(),
  due_at: z.string().nullable(),
  returned_at: z.string().nullable(),
  status: z.enum(["out", "returned", "overdue"]),
});

export const UpdateLoanSchema = z.object({
  book_id: z.string().uuid("Invalid ID").optional(),
  borrower_id: z.string().uuid("Invalid ID").optional(),
  borrowed_at: z.string().nullable(),
  due_at: z.string().nullable(),
  returned_at: z.string().nullable(),
  status: z.enum(["out", "returned", "overdue"]).optional(),
});

export type CreateLoanInput = z.infer<typeof CreateLoanSchema>;
export type UpdateLoanInput = z.infer<typeof UpdateLoanSchema>;
