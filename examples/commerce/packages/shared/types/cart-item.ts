import type { Cart } from "./cart";
import type { Product } from "./product";
import type { Money } from "./money";

export interface CartItem {
  id: string;
  cart_id: string;
  cart?: Cart;
  product_id: string;
  product?: Product;
  variant_id: string;
  quantity: number;
  unit_price: Money;
  title: string;
  variant_label: string;
  image_url: string;
  created_at: string;
  updated_at: string;
}
