import { apiGet, apiPost } from "@/shared/api/client";
import type { AuthResponse, AuthUser, LoginPayload } from "@/features/auth/types";

export function login(payload: LoginPayload) {
  return apiPost<AuthResponse>("/auth/login", payload);
}

export function getMe() {
  return apiGet<AuthUser>("/auth/me");
}
