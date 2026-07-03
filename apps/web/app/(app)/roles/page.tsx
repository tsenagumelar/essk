import { RolesWorkspace } from "@/features/admin/components/admin-workspaces";

export default function RolesPage() {
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-semibold">Roles</h1>
        <p className="mt-2 text-sm text-muted-foreground">Manage global and tenant role master data.</p>
      </div>
      <RolesWorkspace />
    </div>
  );
}
