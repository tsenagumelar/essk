import { apiDelete, apiGet, apiPatch, apiPost } from "@/shared/api/client";
import type { Tenant } from "@/features/tenants/types";

export type { Tenant } from "@/features/tenants/types";

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
