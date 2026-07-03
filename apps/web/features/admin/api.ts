import { apiDelete, apiGet, apiPatch, apiPost } from "@/lib/api/client";

export type Tenant = {
  id: string;
  name: string;
  slug: string;
  status: string;
  is_active: boolean;
};

export type Role = {
  id: string;
  tenant_id?: string;
  name: string;
  code: string;
  description?: string;
  is_system: boolean;
  is_active: boolean;
};

export type Permission = {
  id: string;
  code: string;
  name: string;
  description?: string;
  is_active: boolean;
};

export type AdminUser = {
  id: string;
  tenant_id?: string;
  email: string;
  name: string;
  status: string;
  is_active: boolean;
  role_ids: string[];
};

export function listTenants() {
  return apiGet<Tenant[]>("/tenants");
}

export function createTenant(payload: Pick<Tenant, "name" | "slug">) {
  return apiPost<Tenant>("/tenants", payload);
}

export function updateTenant(id: string, payload: Pick<Tenant, "name" | "status" | "is_active">) {
  return apiPatch<Tenant>(`/tenants/${id}`, payload);
}

export function deleteTenant(id: string) {
  return apiDelete<{ deleted: boolean }>(`/tenants/${id}`);
}

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
