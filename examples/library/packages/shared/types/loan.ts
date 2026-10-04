import type { Book } from "./book";
import type { User } from "./user";

export interface Loan {
  id: string;
  book_id: string;
  book?: Book;
  borrower_id: string;
  borrower?: User;
  borrowed_at: string | null;
  due_at: string | null;
  returned_at: string | null;
  status: "out" | "returned" | "overdue";
  created_at: string;
  updated_at: string;
}
