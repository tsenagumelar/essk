import { TenantsWorkspace } from "@/features/admin/components/admin-workspaces";

export default function TenantsPage() {
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-semibold">Tenants</h1>
        <p className="mt-2 text-sm text-muted-foreground">Manage multi-tenant company workspaces.</p>
      </div>
      <TenantsWorkspace />
    </div>
  );
}
