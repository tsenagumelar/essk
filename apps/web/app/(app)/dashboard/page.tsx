import Link from "next/link";

export default function DashboardPage() {
  return (
    <div>
      <h1 className="text-2xl font-semibold">Dashboard</h1>
      <p className="mt-2 text-sm text-muted-foreground">Foundation workspace with a working modular master data sample.</p>

      <div className="mt-6 grid gap-4 md:grid-cols-3">
        <div className="rounded-lg border bg-white p-4">
          <p className="text-sm text-muted-foreground">Backend</p>
          <p className="mt-2 text-xl font-semibold">Fiber API</p>
          <p className="mt-1 text-sm text-muted-foreground">CQRS-ready modules, RBAC, audit, and rate limiting.</p>
        </div>
        <div className="rounded-lg border bg-white p-4">
          <p className="text-sm text-muted-foreground">Sample Module</p>
          <p className="mt-2 text-xl font-semibold">Products CRUD</p>
          <p className="mt-1 text-sm text-muted-foreground">Tenant scoped master data with soft delete and audit trail.</p>
        </div>
        <div className="rounded-lg border bg-white p-4">
          <p className="text-sm text-muted-foreground">Frontend</p>
          <p className="mt-2 text-xl font-semibold">Next.js App</p>
          <p className="mt-1 text-sm text-muted-foreground">Feature folders and React Query API integration.</p>
        </div>
      </div>

      <Link
        href="/products"
        className="mt-6 inline-flex rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground"
      >
        Open Products CRUD
      </Link>
    </div>
  );
}
