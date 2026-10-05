export interface Blog {
  id: string;
  title: string;
  slug: string;
  content: string;
  image: string | null;
  excerpt: string | null;
  published: boolean;
  published_at: string | null;
  created_at: string;
  updated_at: string;
}
