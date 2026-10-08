import type { Subscription } from "./subscription";
import type { User } from "./user";
import type { Money } from "./money";

export interface Invoice {
  id: string;
  subscription_id: string;
  subscription?: Subscription;
  number: string;
  amount: Money;
  status: "draft" | "open" | "paid" | "void";
  issued_on: string | null;
  paid_on: string | null;
  user_id: string;
  user?: User;
  created_at: string;
  updated_at: string;
}
