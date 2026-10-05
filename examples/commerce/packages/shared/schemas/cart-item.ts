import { z } from "zod";
import { MoneySchema } from "./money";

export const CreateCartItemSchema = z.object({
  cart_id: z.string().uuid("Invalid ID"),
  product_id: z.string().uuid("Invalid ID"),
  variant_id: z.string().min(1, "Required"),
  quantity: z.number().int().optional(),
  unit_price: MoneySchema.optional(),
  title: z.string().min(1, "Required"),
  variant_label: z.string().min(1, "Required"),
  image_url: z.string().min(1, "Required"),
});

export const UpdateCartItemSchema = z.object({
  cart_id: z.string().uuid("Invalid ID").optional(),
  product_id: z.string().uuid("Invalid ID").optional(),
  variant_id: z.string().min(1, "Required").optional(),
  quantity: z.number().int().optional(),
  unit_price: MoneySchema.optional(),
  title: z.string().min(1, "Required").optional(),
  variant_label: z.string().min(1, "Required").optional(),
  image_url: z.string().min(1, "Required").optional(),
});

export type CreateCartItemInput = z.infer<typeof CreateCartItemSchema>;
export type UpdateCartItemInput = z.infer<typeof UpdateCartItemSchema>;
