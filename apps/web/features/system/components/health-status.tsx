"use client";

import { useQuery } from "@tanstack/react-query";
import { getHealth } from "@/features/system/api";

export function HealthStatus() {
  const query = useQuery({
    queryKey: ["system", "health"],
    queryFn: getHealth,
  });

  if (query.isLoading) {
    return <div className="rounded-lg border p-4 text-sm text-muted-foreground">Loading health status...</div>;
  }

  if (query.isError) {
    return <div className="rounded-lg border border-destructive p-4 text-sm text-destructive">Backend health check failed.</div>;
  }

  if (!query.data) {
    return <div className="rounded-lg border p-4 text-sm text-muted-foreground">No health data available.</div>;
  }

  return (
    <div className="rounded-lg border p-4">
      <div className="text-sm font-medium">Backend is {query.data.status}</div>
      <dl className="mt-3 grid gap-2 text-sm text-muted-foreground sm:grid-cols-3">
        <div>
          <dt className="font-medium text-foreground">App</dt>
          <dd>{query.data.app}</dd>
        </div>
        <div>
          <dt className="font-medium text-foreground">Environment</dt>
          <dd>{query.data.env}</dd>
        </div>
        <div>
          <dt className="font-medium text-foreground">Version</dt>
          <dd>{query.data.version}</dd>
        </div>
      </dl>
    </div>
  );
}
