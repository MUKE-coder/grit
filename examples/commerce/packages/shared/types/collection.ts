import type { FileRef } from "../schemas/file-ref";

export interface Collection {
  id: string;
  title: string;
  handle: string;
  description: string;
  image: FileRef | null;
  created_at: string;
  updated_at: string;
}
