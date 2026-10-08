import type { Money } from "./money";

export interface Plan {
  id: string;
  name: string;
  slug: string;
  price: Money;
  interval: "monthly" | "yearly";
  seats: number;
  blurb: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}
