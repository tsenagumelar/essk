import { apiPost } from "@/lib/api/client";

export type AuthUser = {
  id: string;
  tenant_id?: string;
  email: string;
  name: string;
  status: string;
};

export type AuthResponse = {
  access_token: string;
  access_token_expires_at: string;
  refresh_token: string;
  refresh_token_expires_at: string;
  user: AuthUser;
};

export type LoginPayload = {
  email: string;
  password: string;
};

export function login(payload: LoginPayload) {
  return apiPost<AuthResponse>("/auth/login", payload);
}
