"use client";

import { useQuery } from "@tanstack/react-query";
import { getMe } from "@/features/auth/api";
import { getStoredUser } from "@/features/auth/session";

export default function ProfilePage() {
  const fallbackUser = getStoredUser();
  const meQuery = useQuery({
    queryKey: ["auth", "me"],
    queryFn: getMe,
  });

  const user = meQuery.data ?? fallbackUser;

  return (
    <div>
      <h1 className="text-2xl font-semibold">Profile</h1>
      <p className="mt-2 text-sm text-muted-foreground">Authenticated user context from the backend auth module.</p>

      <section className="mt-6 max-w-xl rounded-lg border bg-white p-4">
        {meQuery.isLoading && !user ? <p className="text-sm text-muted-foreground">Loading profile...</p> : null}
        {meQuery.isError ? <p className="text-sm text-destructive">{meQuery.error.message}</p> : null}
        {user ? (
          <dl className="grid gap-3 text-sm">
            <div>
              <dt className="text-muted-foreground">Name</dt>
              <dd className="font-medium">{user.name}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Email</dt>
              <dd className="font-medium">{user.email}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Tenant ID</dt>
              <dd className="break-all font-medium">{user.tenant_id ?? "-"}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Status</dt>
              <dd className="font-medium">{user.status}</dd>
            </div>
          </dl>
        ) : null}
      </section>
    </div>
  );
}
