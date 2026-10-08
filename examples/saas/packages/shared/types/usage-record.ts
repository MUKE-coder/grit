import type { Subscription } from "./subscription";
import type { User } from "./user";

export interface UsageRecord {
  id: string;
  subscription_id: string;
  subscription?: Subscription;
  metric: string;
  quantity: number;
  recorded_on: string | null;
  user_id: string;
  user?: User;
  created_at: string;
  updated_at: string;
}
