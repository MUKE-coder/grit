import type { Author } from "./author";
import type { FileRef } from "../schemas/file-ref";
import type { Money } from "./money";

export interface Book {
  id: string;
  title: string;
  isbn: string;
  summary: string;
  published: string | null;
  price: Money;
  cover: FileRef | null;
  genre: "fiction" | "history" | "science";
  author_id: string;
  author?: Author;
  created_at: string;
  updated_at: string;
}
