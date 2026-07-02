import type { AuthResponse } from "@/features/auth/api";

const accessTokenKey = "essk.access_token";
const refreshTokenKey = "essk.refresh_token";
const userKey = "essk.user";

export function storeSession(auth: AuthResponse) {
  window.localStorage.setItem(accessTokenKey, auth.access_token);
  window.localStorage.setItem(refreshTokenKey, auth.refresh_token);
  window.localStorage.setItem(userKey, JSON.stringify(auth.user));
}

export function clearSession() {
  window.localStorage.removeItem(accessTokenKey);
  window.localStorage.removeItem(refreshTokenKey);
  window.localStorage.removeItem(userKey);
}
