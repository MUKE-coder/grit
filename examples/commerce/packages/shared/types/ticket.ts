export interface Ticket {
  id: string;
  user_id: string;
  subject: string;
  description: string;
  status: string;
  priority: string;
  labels: string;
  assignee_id: string;
  last_reply_at: string | null;
  closed_at: string | null;
  created_at: string;
  updated_at: string;
  user: unknown | null;
  assignee: unknown | null;
  replies: unknown[];
}
