import type { Plan } from "./plan";
import type { User } from "./user";

export interface Subscription {
  id: string;
  plan_id: string;
  plan?: Plan;
  status: "trialing" | "active" | "past_due" | "canceled";
  seats: number;
  current_period_end: string | null;
  canceled_on: string | null;
  user_id: string;
  user?: User;
  created_at: string;
  updated_at: string;
}
