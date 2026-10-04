export interface User {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  role: "ADMIN" | "EDITOR" | "USER"; // grit:role-union
  avatar: string;
  job_title: string;
  bio: string;
  active: boolean;
  provider: string;
  email_verified_at: string | null;
  ip_address: string;
  mac_address: string;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  mac_address?: string; // optional — provided by client if obtainable
}

export interface AuthResponse {
  user: User;
  tokens: {
    access_token: string;
    refresh_token: string;
    expires_at: number;
  };
}
