export interface APIKey {
  id: string;
  user_id: string;
  name: string;
  prefix: string;
  kind: string;
  token: string;
  scopes: string[];
  endpoints: string[];
  origins: string[];
  rate_limit: number;
  last_used_at: string | null;
  expires_at: string | null;
  revoked_at: string | null;
  created_at: string;
}
