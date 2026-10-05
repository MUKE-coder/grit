import type { Collection } from "./collection";
import type { FileRef } from "../schemas/file-ref";
import type { Money } from "./money";

export interface Product {
  id: string;
  title: string;
  handle: string;
  description: string;
  price: Money;
  featured_image: FileRef | null;
  images: FileRef[];
  available: boolean;
  collection_id: string;
  collection?: Collection;
  created_at: string;
  updated_at: string;
}
