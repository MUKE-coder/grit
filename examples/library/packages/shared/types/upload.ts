export interface Upload {
  id: string;
  filename: string;
  original_name: string;
  mime_type: string;
  size: number;
  path: string;
  url: string;
  thumbnail_url?: string;
  user_id: string;
  created_at: string;
  updated_at: string;
}
