export interface TicketReply {
  id: string;
  ticket_id: string;
  user_id: string;
  body: string;
  is_admin_reply: boolean;
  created_at: string;
  updated_at: string;
  user: unknown | null;
}
