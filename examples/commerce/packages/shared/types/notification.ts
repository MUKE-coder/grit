export interface Notification {
  id: string;
  user_id: string;
  source: string;
  severity: string;
  title: string;
  body: string;
  link: string;
  count: number;
  read_at: string | null;
  created_at: string;
  updated_at: string;
}
