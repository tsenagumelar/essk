"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Check, ChevronDown, Download, Pencil, Plus, Search, Trash2, X } from "lucide-react";
import { Dispatch, FormEvent, SetStateAction, useState } from "react";
import {
  createRole,
  createTenant,
  createUser,
  deleteRole,
  deleteTenant,
  deleteUser,
  listRoles,
  listTenants,
  listUsers,
  updateRole,
  updateTenant,
  updateUser,
  type AdminUser,
  type Role,
  type Tenant,
} from "@/features/admin/api";
import { ConfirmationDialog } from "@/features/shared/components/confirmation-dialog";

type PendingAction = {
  title: string;
  description: string;
  confirmLabel: string;
  variant?: "primary" | "danger";
  run: () => Promise<unknown>;
};

function useConfirmableAction() {
  const [pending, setPending] = useState<PendingAction | null>(null);
  return {
    pending,
    request: setPending,
    dialog: (isLoading = false) => (
      <ConfirmationDialog
        open={Boolean(pending)}
        title={pending?.title ?? ""}
        description={pending?.description ?? ""}
        confirmLabel={pending?.confirmLabel ?? "Confirm"}
        variant={pending?.variant ?? "primary"}
        isLoading={isLoading}
        onCancel={() => setPending(null)}
        onConfirm={async () => {
          if (!pending) {
            return;
          }
          await pending.run();
          setPending(null);
        }}
      />
    ),
  };
}

function Modal({
  title,
  subtitle,
  children,
  onClose,
}: Readonly<{ title: string; subtitle: string; children: React.ReactNode; onClose: () => void }>) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
      <section className="w-full max-w-lg rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b px-5 py-4">
          <div>
            <h2 className="text-base font-semibold">{title}</h2>
            <p className="text-sm text-muted-foreground">{subtitle}</p>
          </div>
          <button className="rounded-lg p-2 hover:bg-slate-100" onClick={onClose} aria-label="Close form">
            <X className="h-4 w-4" />
          </button>
        </div>
        {children}
      </section>
    </div>
  );
}

function PageShell({
  title,
  subtitle,
  search,
  onSearch,
  onAdd,
  filters,
  onExport,
  children,
}: Readonly<{
  title: string;
  subtitle: string;
  search: string;
  onSearch: (value: string) => void;
  onAdd: () => void;
  filters?: React.ReactNode;
  onExport: () => void;
  children: React.ReactNode;
}>) {
  return (
    <section className="overflow-hidden rounded-xl border bg-white shadow-sm">
      <div className="flex flex-col gap-3 border-b bg-gradient-to-r from-white to-slate-50 px-4 py-4 lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h2 className="text-base font-semibold">{title}</h2>
          <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
        </div>
        <div className="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:justify-end">
          <div className="relative sm:w-72">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <input
              className="h-10 w-full rounded-lg border bg-white pl-9 pr-3 text-sm outline-none ring-primary focus:ring-2"
              placeholder="Search"
              value={search}
              onChange={(event) => onSearch(event.target.value)}
            />
          </div>
          {filters}
          <button className="inline-flex h-10 items-center justify-center gap-2 rounded-lg border bg-white px-3 text-sm font-medium hover:bg-slate-50" onClick={onExport}>
            <Download className="h-4 w-4" />
            Export
          </button>
          <button className="inline-flex h-10 items-center gap-2 rounded-lg bg-primary px-3 text-sm font-medium text-primary-foreground" onClick={onAdd}>
            <Plus className="h-4 w-4" />
            Add
          </button>
        </div>
      </div>
      {children}
    </section>
  );
}

export function TenantsWorkspace() {
  const queryClient = useQueryClient();
  const confirm = useConfirmableAction();
  const [search, setSearch] = useState("");
  const [editing, setEditing] = useState<Tenant | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [statusFilter, setStatusFilter] = useState("all");
  const [activeFilter, setActiveFilter] = useState("all");
  const [page, setPage] = useState(1);
  const [form, setForm] = useState({ name: "", slug: "", status: "active", is_active: true });
  const tenantsQuery = useQuery({ queryKey: ["tenants"], queryFn: listTenants });
  const tenants = tenantsQuery.data ?? [];
  const filtered = tenants.filter((tenant) => {
    const matchesSearch = `${tenant.name} ${tenant.slug} ${tenant.status}`.toLowerCase().includes(search.toLowerCase());
    const matchesStatus = statusFilter === "all" || tenant.status === statusFilter;
    const matchesActive = activeFilter === "all" || String(tenant.is_active) === activeFilter;
    return matchesSearch && matchesStatus && matchesActive;
  });
  const paginated = paginate(filtered, page);
  const saveMutation = useMutation({
    mutationFn: async () =>
      editing
        ? updateTenant(editing.id, { name: form.name, status: form.status, is_active: form.is_active })
        : createTenant({ name: form.name, slug: form.slug }),
    onSuccess: async () => {
      setEditing(null);
      setIsOpen(false);
      await queryClient.invalidateQueries({ queryKey: ["tenants"] });
    },
  });
  const deleteMutation = useMutation({
    mutationFn: deleteTenant,
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: ["tenants"] }),
  });

  function openForm(tenant?: Tenant) {
    setEditing(tenant ?? null);
    setIsOpen(true);
    setForm(
      tenant
        ? { name: tenant.name, slug: tenant.slug, status: tenant.status, is_active: tenant.is_active }
        : { name: "", slug: "", status: "active", is_active: true },
    );
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    confirm.request({
      title: editing ? "Save tenant changes?" : "Create tenant?",
      description: editing ? `Update ${editing.name}.` : `Create ${form.name}.`,
      confirmLabel: editing ? "Save Edit" : "Save Add",
      run: () => saveMutation.mutateAsync(),
    });
  }

  return (
    <>
      <PageShell
        title="Tenants"
        subtitle="Super admin can manage all tenants."
        search={search}
        onSearch={(value) => {
          setSearch(value);
          setPage(1);
        }}
        onAdd={() => openForm()}
        filters={
          <>
            <FilterSelect label="Status" value={statusFilter} onChange={(value) => { setStatusFilter(value); setPage(1); }} options={[{ value: "all", label: "All status" }, { value: "active", label: "Active" }, { value: "inactive", label: "Inactive" }]} />
            <FilterSelect label="Active" value={activeFilter} onChange={(value) => { setActiveFilter(value); setPage(1); }} options={[{ value: "all", label: "All active" }, { value: "true", label: "Yes" }, { value: "false", label: "No" }]} />
          </>
        }
        onExport={() => exportCsv("tenants.csv", ["Name", "Slug", "Status", "Active"], filtered.map((tenant) => [tenant.name, tenant.slug, tenant.status, tenant.is_active ? "Yes" : "No"]))}
      >
        <DataTable
          headers={["Name", "Slug", "Status", "Active", "Actions"]}
          rows={paginated.items.map((tenant) => [
            tenant.name,
            tenant.slug,
            tenant.status,
            tenant.is_active ? "Yes" : "No",
            <RowActions
              key={tenant.id}
              onEdit={() => openForm(tenant)}
              onDelete={() =>
                confirm.request({
                  title: "Delete tenant?",
                  description: `This will soft delete ${tenant.name}.`,
                  confirmLabel: "Delete",
                  variant: "danger",
                  run: () => deleteMutation.mutateAsync(tenant.id),
                })
              }
            />,
          ])}
          loading={tenantsQuery.isLoading}
        />
        <Pagination page={paginated.page} totalPages={paginated.totalPages} totalItems={filtered.length} onPageChange={setPage} />
      </PageShell>
      {isOpen ? (
        <Modal title={editing ? "Edit Tenant" : "Create Tenant"} subtitle="Manage tenant master data." onClose={() => setIsOpen(false)}>
          <form className="space-y-4 p-5" onSubmit={submit}>
            <Input label="Name" value={form.name} onChange={(value) => setForm((current) => ({ ...current, name: value }))} />
            <Input label="Slug" value={form.slug} disabled={Boolean(editing)} onChange={(value) => setForm((current) => ({ ...current, slug: value }))} />
            {editing ? <ActiveField form={form} setForm={setForm} /> : null}
            <FormActions onCancel={() => setIsOpen(false)} loading={saveMutation.isPending} />
          </form>
        </Modal>
      ) : null}
      {confirm.dialog(saveMutation.isPending || deleteMutation.isPending)}
    </>
  );
}

export function UsersWorkspace() {
  const queryClient = useQueryClient();
  const confirm = useConfirmableAction();
  const [search, setSearch] = useState("");
  const [editing, setEditing] = useState<AdminUser | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [tenantFilter, setTenantFilter] = useState("all");
  const [roleFilter, setRoleFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [page, setPage] = useState(1);
  const [form, setForm] = useState({ tenant_id: "", email: "", name: "", password: "Admin123!", status: "active", is_active: true, role_ids: [] as string[] });
  const tenantsQuery = useQuery({ queryKey: ["tenants"], queryFn: listTenants });
  const rolesQuery = useQuery({ queryKey: ["roles"], queryFn: () => listRoles() });
  const usersQuery = useQuery({ queryKey: ["users"], queryFn: listUsers });
  const users = usersQuery.data ?? [];
  const roles = rolesQuery.data ?? [];
  const tenants = tenantsQuery.data ?? [];
  const roleOptions = roles
    .filter((role) => role.tenant_id === form.tenant_id || form.role_ids.includes(role.id))
    .map((role) => ({ value: role.id, label: `${role.name} (${role.code})` }));
  const roleNamesById = new Map(roles.map((role) => [role.id, `${role.name} (${role.code})`]));
  const filtered = users.filter((user) => {
    const roleLabels = user.role_ids.map((roleID) => roleNamesById.get(roleID) ?? roleID).join(" ");
    const matchesSearch = `${user.name} ${user.email} ${user.status} ${roleLabels}`.toLowerCase().includes(search.toLowerCase());
    const matchesTenant = tenantFilter === "all" || user.tenant_id === tenantFilter;
    const matchesRole = roleFilter === "all" || user.role_ids.includes(roleFilter);
    const matchesStatus = statusFilter === "all" || user.status === statusFilter;
    return matchesSearch && matchesTenant && matchesRole && matchesStatus;
  });
  const paginated = paginate(filtered, page);
  const saveMutation = useMutation({
    mutationFn: async () =>
      editing
        ? updateUser(editing.id, { name: form.name, status: form.status, is_active: form.is_active, role_ids: form.role_ids })
        : createUser({ tenant_id: form.tenant_id, email: form.email, name: form.name, password: form.password, role_ids: form.role_ids }),
    onSuccess: async () => {
      setIsOpen(false);
      setEditing(null);
      await queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
  const deleteMutation = useMutation({
    mutationFn: deleteUser,
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: ["users"] }),
  });

  function openForm(user?: AdminUser) {
    setEditing(user ?? null);
    setIsOpen(true);
    setForm(
      user
        ? { tenant_id: user.tenant_id ?? "", email: user.email, name: user.name, password: "", status: user.status, is_active: user.is_active, role_ids: user.role_ids }
        : { tenant_id: tenants[0]?.id ?? "", email: "", name: "", password: "Admin123!", status: "active", is_active: true, role_ids: [] },
    );
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    confirm.request({
      title: editing ? "Save user changes?" : "Create user?",
      description: editing ? `Update ${editing.email}.` : `Create ${form.email}.`,
      confirmLabel: editing ? "Save Edit" : "Save Add",
      run: () => saveMutation.mutateAsync(),
    });
  }

  return (
    <>
      <PageShell
        title="Users"
        subtitle="Super admin sees all users; tenant admins are scoped to their tenant."
        search={search}
        onSearch={(value) => {
          setSearch(value);
          setPage(1);
        }}
        onAdd={() => openForm()}
        filters={
          <>
            <FilterSelect label="Tenant" value={tenantFilter} onChange={(value) => { setTenantFilter(value); setPage(1); }} options={[{ value: "all", label: "All tenants" }, ...tenants.map((tenant) => ({ value: tenant.id, label: tenant.name }))]} />
            <FilterSelect label="Role" value={roleFilter} onChange={(value) => { setRoleFilter(value); setPage(1); }} options={[{ value: "all", label: "All roles" }, ...roles.map((role) => ({ value: role.id, label: `${role.name} (${role.code})` }))]} />
            <FilterSelect label="Status" value={statusFilter} onChange={(value) => { setStatusFilter(value); setPage(1); }} options={[{ value: "all", label: "All status" }, { value: "active", label: "Active" }, { value: "inactive", label: "Inactive" }, { value: "invited", label: "Invited" }, { value: "suspended", label: "Suspended" }]} />
          </>
        }
        onExport={() =>
          exportCsv(
            "users.csv",
            ["Name", "Email", "Tenant", "Roles", "Status", "Active"],
            filtered.map((user) => [
              user.name,
              user.email,
              tenants.find((tenant) => tenant.id === user.tenant_id)?.name ?? "-",
              user.role_ids.map((roleID) => roleNamesById.get(roleID) ?? roleID).join("; "),
              user.status,
              user.is_active ? "Yes" : "No",
            ]),
          )
        }
      >
        <DataTable
          headers={["Name", "Email", "Tenant", "Roles", "Status", "Active", "Actions"]}
          rows={paginated.items.map((user) => [
            user.name,
            user.email,
            tenants.find((tenant) => tenant.id === user.tenant_id)?.name ?? "-",
            <RoleBadges key={`${user.id}-roles`} labels={user.role_ids.map((roleID) => roleNamesById.get(roleID) ?? roleID)} />,
            user.status,
            user.is_active ? "Yes" : "No",
            <RowActions
              key={user.id}
              onEdit={() => openForm(user)}
              onDelete={() =>
                confirm.request({
                  title: "Delete user?",
                  description: `This will soft delete ${user.email}.`,
                  confirmLabel: "Delete",
                  variant: "danger",
                  run: () => deleteMutation.mutateAsync(user.id),
                })
              }
            />,
          ])}
          loading={usersQuery.isLoading}
        />
        <Pagination page={paginated.page} totalPages={paginated.totalPages} totalItems={filtered.length} onPageChange={setPage} />
      </PageShell>
      {isOpen ? (
        <Modal title={editing ? "Edit User" : "Create User"} subtitle="Manage tenant user account and roles." onClose={() => setIsOpen(false)}>
          <form className="space-y-4 p-5" onSubmit={submit}>
            <Select label="Tenant" value={form.tenant_id} disabled={Boolean(editing)} onChange={(value) => setForm((current) => ({ ...current, tenant_id: value }))} options={tenants.map((tenant) => ({ value: tenant.id, label: tenant.name }))} />
            <Input label="Email" value={form.email} disabled={Boolean(editing)} onChange={(value) => setForm((current) => ({ ...current, email: value }))} />
            <Input label="Name" value={form.name} onChange={(value) => setForm((current) => ({ ...current, name: value }))} />
            {!editing ? <Input label="Password" value={form.password} onChange={(value) => setForm((current) => ({ ...current, password: value }))} /> : null}
            <RoleDropdown label="Roles" values={form.role_ids} options={roleOptions} onChange={(values) => setForm((current) => ({ ...current, role_ids: values }))} />
            {editing ? <ActiveField form={form} setForm={setForm} /> : null}
            <FormActions onCancel={() => setIsOpen(false)} loading={saveMutation.isPending} />
          </form>
        </Modal>
      ) : null}
      {confirm.dialog(saveMutation.isPending || deleteMutation.isPending)}
    </>
  );
}

export function RolesWorkspace() {
  const queryClient = useQueryClient();
  const confirm = useConfirmableAction();
  const [search, setSearch] = useState("");
  const [editing, setEditing] = useState<Role | null>(null);
  const [isOpen, setIsOpen] = useState(false);
  const [tenantFilter, setTenantFilter] = useState("all");
  const [systemFilter, setSystemFilter] = useState("all");
  const [activeFilter, setActiveFilter] = useState("all");
  const [page, setPage] = useState(1);
  const [form, setForm] = useState({ tenant_id: "", name: "", code: "", description: "", is_system: false, is_active: true });
  const tenantsQuery = useQuery({ queryKey: ["tenants"], queryFn: listTenants });
  const rolesQuery = useQuery({ queryKey: ["roles"], queryFn: () => listRoles() });
  const roles = rolesQuery.data ?? [];
  const tenants = tenantsQuery.data ?? [];
  const filtered = roles.filter((role) => {
    const tenantLabel = role.tenant_id ? tenants.find((tenant) => tenant.id === role.tenant_id)?.name ?? role.tenant_id : "Global";
    const matchesSearch = `${role.name} ${role.code} ${tenantLabel}`.toLowerCase().includes(search.toLowerCase());
    const matchesTenant = tenantFilter === "all" || (tenantFilter === "global" ? !role.tenant_id : role.tenant_id === tenantFilter);
    const matchesSystem = systemFilter === "all" || String(role.is_system) === systemFilter;
    const matchesActive = activeFilter === "all" || String(role.is_active) === activeFilter;
    return matchesSearch && matchesTenant && matchesSystem && matchesActive;
  });
  const paginated = paginate(filtered, page);
  const saveMutation = useMutation({
    mutationFn: async () =>
      editing
        ? updateRole(editing.id, { name: form.name, description: form.description, is_active: form.is_active })
        : createRole({ tenant_id: form.tenant_id || undefined, name: form.name, code: form.code, description: form.description, is_system: form.is_system }),
    onSuccess: async () => {
      setIsOpen(false);
      setEditing(null);
      await queryClient.invalidateQueries({ queryKey: ["roles"] });
    },
  });
  const deleteMutation = useMutation({
    mutationFn: deleteRole,
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: ["roles"] }),
  });

  function openForm(role?: Role) {
    setEditing(role ?? null);
    setIsOpen(true);
    setForm(
      role
        ? { tenant_id: role.tenant_id ?? "", name: role.name, code: role.code, description: role.description ?? "", is_system: role.is_system, is_active: role.is_active }
        : { tenant_id: "", name: "", code: "", description: "", is_system: false, is_active: true },
    );
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    confirm.request({
      title: editing ? "Save role changes?" : "Create role?",
      description: editing ? `Update ${editing.code}.` : `Create ${form.code}.`,
      confirmLabel: editing ? "Save Edit" : "Save Add",
      run: () => saveMutation.mutateAsync(),
    });
  }

  return (
    <>
      <PageShell
        title="Roles"
        subtitle="Manage global and tenant scoped roles."
        search={search}
        onSearch={(value) => {
          setSearch(value);
          setPage(1);
        }}
        onAdd={() => openForm()}
        filters={
          <>
            <FilterSelect label="Tenant" value={tenantFilter} onChange={(value) => { setTenantFilter(value); setPage(1); }} options={[{ value: "all", label: "All tenants" }, { value: "global", label: "Global" }, ...tenants.map((tenant) => ({ value: tenant.id, label: tenant.name }))]} />
            <FilterSelect label="System" value={systemFilter} onChange={(value) => { setSystemFilter(value); setPage(1); }} options={[{ value: "all", label: "All system" }, { value: "true", label: "Yes" }, { value: "false", label: "No" }]} />
            <FilterSelect label="Active" value={activeFilter} onChange={(value) => { setActiveFilter(value); setPage(1); }} options={[{ value: "all", label: "All active" }, { value: "true", label: "Yes" }, { value: "false", label: "No" }]} />
          </>
        }
        onExport={() =>
          exportCsv(
            "roles.csv",
            ["Name", "Code", "Tenant", "System", "Active"],
            filtered.map((role) => [
              role.name,
              role.code,
              role.tenant_id ? tenants.find((tenant) => tenant.id === role.tenant_id)?.name ?? role.tenant_id : "Global",
              role.is_system ? "Yes" : "No",
              role.is_active ? "Yes" : "No",
            ]),
          )
        }
      >
        <DataTable
          headers={["Name", "Code", "Tenant", "System", "Active", "Actions"]}
          rows={paginated.items.map((role) => [
            role.name,
            role.code,
            role.tenant_id ? tenants.find((tenant) => tenant.id === role.tenant_id)?.name ?? role.tenant_id : "Global",
            role.is_system ? "Yes" : "No",
            role.is_active ? "Yes" : "No",
            <RowActions
              key={role.id}
              onEdit={() => openForm(role)}
              onDelete={() =>
                confirm.request({
                  title: "Delete role?",
                  description: `This will soft delete ${role.code}.`,
                  confirmLabel: "Delete",
                  variant: "danger",
                  run: () => deleteMutation.mutateAsync(role.id),
                })
              }
            />,
          ])}
          loading={rolesQuery.isLoading}
        />
        <Pagination page={paginated.page} totalPages={paginated.totalPages} totalItems={filtered.length} onPageChange={setPage} />
      </PageShell>
      {isOpen ? (
        <Modal title={editing ? "Edit Role" : "Create Role"} subtitle="Manage RBAC role master data." onClose={() => setIsOpen(false)}>
          <form className="space-y-4 p-5" onSubmit={submit}>
            <Select label="Tenant" value={form.tenant_id} disabled={Boolean(editing)} onChange={(value) => setForm((current) => ({ ...current, tenant_id: value }))} options={[{ value: "", label: "Global" }, ...tenants.map((tenant) => ({ value: tenant.id, label: tenant.name }))]} />
            <Input label="Name" value={form.name} onChange={(value) => setForm((current) => ({ ...current, name: value }))} />
            <Input label="Code" value={form.code} disabled={Boolean(editing)} onChange={(value) => setForm((current) => ({ ...current, code: value }))} />
            <Input label="Description" value={form.description} onChange={(value) => setForm((current) => ({ ...current, description: value }))} />
            {editing ? <ActiveField form={form} setForm={setForm} /> : null}
            <FormActions onCancel={() => setIsOpen(false)} loading={saveMutation.isPending} />
          </form>
        </Modal>
      ) : null}
      {confirm.dialog(saveMutation.isPending || deleteMutation.isPending)}
    </>
  );
}

function DataTable({ headers, rows, loading }: Readonly<{ headers: string[]; rows: React.ReactNode[][]; loading: boolean }>) {
  if (loading) {
    return <p className="p-4 text-sm text-muted-foreground">Loading...</p>;
  }
  if (rows.length === 0) {
    return <p className="p-10 text-center text-sm text-muted-foreground">No records found.</p>;
  }
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-left text-sm">
        <thead className="bg-slate-50 text-xs uppercase text-muted-foreground">
          <tr>{headers.map((header) => <th key={header} className="px-4 py-3">{header}</th>)}</tr>
        </thead>
        <tbody>{rows.map((row, index) => <tr key={index} className="border-t hover:bg-slate-50">{row.map((cell, cellIndex) => <td key={cellIndex} className="px-4 py-3">{cell}</td>)}</tr>)}</tbody>
      </table>
    </div>
  );
}

function Pagination({
  page,
  totalPages,
  totalItems,
  onPageChange,
}: Readonly<{ page: number; totalPages: number; totalItems: number; onPageChange: (page: number) => void }>) {
  if (totalItems === 0) {
    return null;
  }

  return (
    <div className="flex flex-col gap-3 border-t px-4 py-3 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between">
      <p>
        Page {page} of {totalPages} · {totalItems} records
      </p>
      <div className="flex gap-2">
        <button
          type="button"
          className="rounded-lg border bg-white px-3 py-2 font-medium text-slate-700 disabled:cursor-not-allowed disabled:opacity-50"
          disabled={page <= 1}
          onClick={() => onPageChange(page - 1)}
        >
          Previous
        </button>
        <button
          type="button"
          className="rounded-lg border bg-white px-3 py-2 font-medium text-slate-700 disabled:cursor-not-allowed disabled:opacity-50"
          disabled={page >= totalPages}
          onClick={() => onPageChange(page + 1)}
        >
          Next
        </button>
      </div>
    </div>
  );
}

function RowActions({ onEdit, onDelete }: Readonly<{ onEdit: () => void; onDelete: () => void }>) {
  return (
    <div className="flex justify-end gap-2">
      <button className="inline-flex h-8 w-8 items-center justify-center rounded-lg border hover:bg-slate-50" onClick={onEdit}>
        <Pencil className="h-4 w-4" />
      </button>
      <button className="inline-flex h-8 w-8 items-center justify-center rounded-lg border text-destructive hover:bg-slate-50" onClick={onDelete}>
        <Trash2 className="h-4 w-4" />
      </button>
    </div>
  );
}

function FilterSelect({
  label,
  value,
  options,
  onChange,
}: Readonly<{ label: string; value: string; options: { value: string; label: string }[]; onChange: (value: string) => void }>) {
  return (
    <label className="sr-only">
      {label}
      <select
        className="not-sr-only h-10 rounded-lg border bg-white px-3 text-sm outline-none ring-primary focus:ring-2"
        value={value}
        onChange={(event) => onChange(event.target.value)}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </label>
  );
}

function RoleBadges({ labels }: Readonly<{ labels: string[] }>) {
  if (labels.length === 0) {
    return <span className="text-muted-foreground">-</span>;
  }

  return (
    <div className="flex max-w-80 flex-wrap gap-1">
      {labels.map((label) => (
        <span key={label} className="rounded-full bg-slate-100 px-2 py-1 text-xs font-medium text-slate-700">
          {label}
        </span>
      ))}
    </div>
  );
}

function Input({ label, value, disabled, onChange }: Readonly<{ label: string; value: string; disabled?: boolean; onChange: (value: string) => void }>) {
  return (
    <label className="block text-sm font-medium">
      {label}
      <input disabled={disabled} required className="mt-1 h-10 w-full rounded-lg border px-3 text-sm outline-none ring-primary focus:ring-2 disabled:bg-muted" value={value} onChange={(event) => onChange(event.target.value)} />
    </label>
  );
}

function Select({ label, value, disabled, options, onChange }: Readonly<{ label: string; value: string; disabled?: boolean; options: { value: string; label: string }[]; onChange: (value: string) => void }>) {
  return (
    <label className="block text-sm font-medium">
      {label}
      <select disabled={disabled} className="mt-1 h-10 w-full rounded-lg border px-3 text-sm outline-none ring-primary focus:ring-2 disabled:bg-muted" value={value} onChange={(event) => onChange(event.target.value)}>
        {options.map((option) => <option key={option.value} value={option.value}>{option.label}</option>)}
      </select>
    </label>
  );
}

function RoleDropdown({ label, values, options, onChange }: Readonly<{ label: string; values: string[]; options: { value: string; label: string }[]; onChange: (values: string[]) => void }>) {
  const [isOpen, setIsOpen] = useState(false);
  const selectedLabels = options.filter((option) => values.includes(option.value)).map((option) => option.label);

  function toggle(value: string) {
    onChange(values.includes(value) ? values.filter((current) => current !== value) : [...values, value]);
  }

  return (
    <div className="relative text-sm font-medium">
      <p>{label}</p>
      <button
        type="button"
        className="mt-1 flex min-h-10 w-full items-center justify-between gap-3 rounded-lg border bg-white px-3 py-2 text-left text-sm outline-none ring-primary focus:ring-2"
        onClick={() => setIsOpen((current) => !current)}
      >
        <span className={selectedLabels.length === 0 ? "text-muted-foreground" : "text-slate-900"}>
          {selectedLabels.length === 0 ? "Select roles" : selectedLabels.join(", ")}
        </span>
        <ChevronDown className="h-4 w-4 shrink-0 text-muted-foreground" />
      </button>
      {isOpen ? (
        <div className="absolute z-20 mt-2 max-h-56 w-full overflow-auto rounded-lg border bg-white p-2 shadow-lg">
          {options.length === 0 ? (
            <p className="px-2 py-2 text-sm text-muted-foreground">No roles available for this tenant.</p>
          ) : (
            options.map((option) => {
              const selected = values.includes(option.value);
              return (
                <button
                  key={option.value}
                  type="button"
                  className="flex w-full items-center gap-2 rounded-md px-2 py-2 text-left text-sm hover:bg-slate-50"
                  onClick={() => toggle(option.value)}
                >
                  <span className="flex h-4 w-4 items-center justify-center rounded border">
                    {selected ? <Check className="h-3 w-3" /> : null}
                  </span>
                  {option.label}
                </button>
              );
            })
          )}
        </div>
      ) : null}
    </div>
  );
}

function ActiveField<T extends { is_active: boolean }>({
  form,
  setForm,
}: {
  form: T;
  setForm: Dispatch<SetStateAction<T>>;
}) {
  return (
    <label className="flex items-center gap-2 text-sm">
      <input type="checkbox" checked={form.is_active} onChange={(event) => setForm((current) => ({ ...current, is_active: event.target.checked }))} />
      Active
    </label>
  );
}

function FormActions({ onCancel, loading }: Readonly<{ onCancel: () => void; loading: boolean }>) {
  return (
    <div className="flex justify-end gap-2 border-t pt-4">
      <button type="button" className="rounded-lg border px-3 py-2 text-sm font-medium" onClick={onCancel}>
        Cancel
      </button>
      <button type="submit" disabled={loading} className="rounded-lg bg-primary px-3 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60">
        {loading ? "Saving..." : "Save"}
      </button>
    </div>
  );
}

function paginate<T>(items: T[], page: number, pageSize = 10) {
  const totalPages = Math.max(1, Math.ceil(items.length / pageSize));
  const normalizedPage = Math.min(Math.max(page, 1), totalPages);
  const start = (normalizedPage - 1) * pageSize;

  return {
    items: items.slice(start, start + pageSize),
    page: normalizedPage,
    totalPages,
  };
}

function exportCsv(filename: string, headers: string[], rows: Array<Array<string | number | boolean>>) {
  const csv = [headers, ...rows].map((row) => row.map(csvCell).join(",")).join("\n");
  const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}

function csvCell(value: string | number | boolean) {
  const raw = String(value);
  return `"${raw.replaceAll('"', '""')}"`;
}
