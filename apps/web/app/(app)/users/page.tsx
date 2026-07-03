import { UsersWorkspace } from "@/features/admin/components/admin-workspaces";

export default function UsersPage() {
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-semibold">Users</h1>
        <p className="mt-2 text-sm text-muted-foreground">Manage tenant users and assigned roles.</p>
      </div>
      <UsersWorkspace />
    </div>
  );
}
