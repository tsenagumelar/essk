import { apiDelete, apiGet, apiPatch, apiPost } from "@/shared/functions/api/client";
import type { Permission, Role } from "@/features/roles/types";

export type { Permission, Role } from "@/features/roles/types";

export function listRoles(tenantId?: string) {
  return apiGet<Role[]>(tenantId ? `/roles?tenant_id=${tenantId}` : "/roles");
}

export function createRole(payload: {
  tenant_id?: string;
  name: string;
  code: string;
  description?: string;
  is_system: boolean;
}) {
  return apiPost<Role>("/roles", payload);
}

export function updateRole(id: string, payload: Pick<Role, "name" | "description" | "is_active">) {
  return apiPatch<Role>(`/roles/${id}`, payload);
}

export function deleteRole(id: string) {
  return apiDelete<{ deleted: boolean }>(`/roles/${id}`);
}

export function listPermissions() {
  return apiGet<Permission[]>("/permissions");
}

export function assignPermission(roleId: string, permissionId: string) {
  return apiPost<{ assigned: boolean }>(`/roles/${roleId}/permissions`, { permission_id: permissionId });
}

export function removePermission(roleId: string, permissionId: string) {
  return apiDelete<{ removed: boolean }>(`/roles/${roleId}/permissions/${permissionId}`);
}
