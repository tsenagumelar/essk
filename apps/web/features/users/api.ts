import { apiDelete, apiGet, apiPatch, apiPost } from "@/shared/functions/api/client";
import type { AdminUser } from "@/features/users/types";

export type { AdminUser } from "@/features/users/types";

export function listUsers() {
  return apiGet<AdminUser[]>("/users");
}

export function createUser(payload: {
  tenant_id: string;
  email: string;
  name: string;
  password: string;
  role_ids: string[];
}) {
  return apiPost<AdminUser>("/users", payload);
}

export function updateUser(id: string, payload: Pick<AdminUser, "name" | "status" | "is_active" | "role_ids">) {
  return apiPatch<AdminUser>(`/users/${id}`, payload);
}

export function deleteUser(id: string) {
  return apiDelete<{ deleted: boolean }>(`/users/${id}`);
}
