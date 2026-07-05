"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Dispatch, FormEvent, SetStateAction, useState } from "react";
import {
  createTenant,
  deleteTenant,
  listTenants,
  updateTenant,
  type Tenant,
} from "@/features/tenants/api";
import { createUser, deleteUser, listUsers, updateUser, type AdminUser } from "@/features/users/api";
import { createRole, deleteRole, listRoles, updateRole, type Role } from "@/features/roles/api";
import { Badge } from "@/shared/components/atoms/badge";
import { Button } from "@/shared/components/atoms/button";
import { CheckboxField } from "@/shared/components/atoms/checkbox-field";
import { SelectField } from "@/shared/components/atoms/select-field";
import { TextInput } from "@/shared/components/atoms/text-input";
import { ConfirmableActionDialog } from "@/shared/components/molecules/confirmable-action-dialog";
import { FilterSelect as SharedFilterSelect } from "@/shared/components/molecules/filter-select";
import { Modal as SharedModal } from "@/shared/components/molecules/modal";
import { MultiSelect } from "@/shared/components/molecules/multi-select";
import { Pagination as SharedPagination } from "@/shared/components/molecules/pagination";
import { RowActions as SharedRowActions } from "@/shared/components/molecules/row-actions";
import { CrudToolbar } from "@/shared/components/organisms/crud-toolbar";
import { DataTable as SharedDataTable } from "@/shared/components/organisms/data-table";
import { PageShell as SharedPageShell } from "@/shared/components/organisms/page-shell";
import { useConfirmableAction } from "@/shared/hooks/use-confirmable-action";
import { exportExcel } from "@/shared/functions/export/export-excel";
import { printPdf } from "@/shared/functions/export/print-pdf";

const pageSizeOptions = [5, 10, 20];

function Modal({
  title,
  subtitle,
  children,
  onClose,
}: Readonly<{ title: string; subtitle: string; children: React.ReactNode; onClose: () => void }>) {
  return (
    <SharedModal title={title} subtitle={subtitle} onClose={onClose}>
      {children}
    </SharedModal>
  );
}

function PageShell({
  title,
  subtitle,
  search,
  onSearch,
  onAdd,
  filters,
  onExportExcel,
  onExportPdf,
  children,
}: Readonly<{
  title: string;
  subtitle: string;
  search: string;
  onSearch: (value: string) => void;
  onAdd: () => void;
  filters?: React.ReactNode;
  onExportExcel: () => void;
  onExportPdf: () => void;
  children: React.ReactNode;
}>) {
  return (
    <SharedPageShell
      title={title}
      subtitle={subtitle}
      toolbar={
        <CrudToolbar
          search={search}
          onSearch={onSearch}
          onAdd={onAdd}
          filters={filters}
          onExportExcel={onExportExcel}
          onExportPdf={onExportPdf}
        />
      }
    >
      {children}
    </SharedPageShell>
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
  const [pageSize, setPageSize] = useState(10);
  const [form, setForm] = useState({ name: "", slug: "", status: "active", is_active: true });
  const tenantsQuery = useQuery({ queryKey: ["tenants"], queryFn: listTenants });
  const tenants = tenantsQuery.data ?? [];
  const filtered = tenants.filter((tenant) => {
    const matchesSearch = `${tenant.name} ${tenant.slug} ${tenant.status}`.toLowerCase().includes(search.toLowerCase());
    const matchesStatus = statusFilter === "all" || tenant.status === statusFilter;
    const matchesActive = activeFilter === "all" || String(tenant.is_active) === activeFilter;
    return matchesSearch && matchesStatus && matchesActive;
  });
  const paginated = paginate(filtered, page, pageSize);
  const activeCount = filtered.filter((tenant) => tenant.is_active).length;
  const firstRecord = filtered.length === 0 ? 0 : (paginated.page - 1) * pageSize + 1;
  const lastRecord = Math.min(paginated.page * pageSize, filtered.length);
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
        subtitle={`${filtered.length} records, ${activeCount} active`}
        search={search}
        onSearch={(value) => {
          setSearch(value);
          setPage(1);
        }}
        onAdd={() => openForm()}
        filters={
          <>
            <FilterSelect withIcon label="Status" value={statusFilter} onChange={(value) => { setStatusFilter(value); setPage(1); }} options={[{ value: "all", label: "All status" }, { value: "active", label: "Active" }, { value: "inactive", label: "Inactive" }]} />
            <FilterSelect label="Active" value={activeFilter} onChange={(value) => { setActiveFilter(value); setPage(1); }} options={[{ value: "all", label: "All active" }, { value: "true", label: "Yes" }, { value: "false", label: "No" }]} />
          </>
        }
        onExportExcel={() => exportExcel("tenants.xls", ["Name", "Slug", "Status", "Active"], filtered.map((tenant) => [tenant.name, tenant.slug, tenant.status, tenant.is_active ? "Yes" : "No"]))}
        onExportPdf={printPdf}
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
        <Pagination
          page={paginated.page}
          pageSize={pageSize}
          totalPages={paginated.totalPages}
          totalItems={filtered.length}
          firstRecord={firstRecord}
          lastRecord={lastRecord}
          onPageChange={setPage}
          onPageSizeChange={(value) => {
            setPageSize(value);
            setPage(1);
          }}
        />
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
      <ConfirmableActionDialog
        pending={confirm.pending}
        isLoading={saveMutation.isPending || deleteMutation.isPending}
        onCancel={confirm.cancel}
        onConfirm={confirm.confirm}
      />
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
  const [pageSize, setPageSize] = useState(10);
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
  const paginated = paginate(filtered, page, pageSize);
  const activeCount = filtered.filter((user) => user.is_active).length;
  const firstRecord = filtered.length === 0 ? 0 : (paginated.page - 1) * pageSize + 1;
  const lastRecord = Math.min(paginated.page * pageSize, filtered.length);
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
        subtitle={`${filtered.length} records, ${activeCount} active`}
        search={search}
        onSearch={(value) => {
          setSearch(value);
          setPage(1);
        }}
        onAdd={() => openForm()}
        filters={
          <>
            <FilterSelect withIcon label="Tenant" value={tenantFilter} onChange={(value) => { setTenantFilter(value); setPage(1); }} options={[{ value: "all", label: "All tenants" }, ...tenants.map((tenant) => ({ value: tenant.id, label: tenant.name }))]} />
            <FilterSelect label="Role" value={roleFilter} onChange={(value) => { setRoleFilter(value); setPage(1); }} options={[{ value: "all", label: "All roles" }, ...roles.map((role) => ({ value: role.id, label: `${role.name} (${role.code})` }))]} />
            <FilterSelect label="Status" value={statusFilter} onChange={(value) => { setStatusFilter(value); setPage(1); }} options={[{ value: "all", label: "All status" }, { value: "active", label: "Active" }, { value: "inactive", label: "Inactive" }, { value: "invited", label: "Invited" }, { value: "suspended", label: "Suspended" }]} />
          </>
        }
        onExportExcel={() =>
          exportExcel(
            "users.xls",
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
        onExportPdf={printPdf}
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
        <Pagination
          page={paginated.page}
          pageSize={pageSize}
          totalPages={paginated.totalPages}
          totalItems={filtered.length}
          firstRecord={firstRecord}
          lastRecord={lastRecord}
          onPageChange={setPage}
          onPageSizeChange={(value) => {
            setPageSize(value);
            setPage(1);
          }}
        />
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
      <ConfirmableActionDialog
        pending={confirm.pending}
        isLoading={saveMutation.isPending || deleteMutation.isPending}
        onCancel={confirm.cancel}
        onConfirm={confirm.confirm}
      />
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
  const [pageSize, setPageSize] = useState(10);
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
  const paginated = paginate(filtered, page, pageSize);
  const activeCount = filtered.filter((role) => role.is_active).length;
  const firstRecord = filtered.length === 0 ? 0 : (paginated.page - 1) * pageSize + 1;
  const lastRecord = Math.min(paginated.page * pageSize, filtered.length);
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
        subtitle={`${filtered.length} records, ${activeCount} active`}
        search={search}
        onSearch={(value) => {
          setSearch(value);
          setPage(1);
        }}
        onAdd={() => openForm()}
        filters={
          <>
            <FilterSelect withIcon label="Tenant" value={tenantFilter} onChange={(value) => { setTenantFilter(value); setPage(1); }} options={[{ value: "all", label: "All tenants" }, { value: "global", label: "Global" }, ...tenants.map((tenant) => ({ value: tenant.id, label: tenant.name }))]} />
            <FilterSelect label="System" value={systemFilter} onChange={(value) => { setSystemFilter(value); setPage(1); }} options={[{ value: "all", label: "All system" }, { value: "true", label: "Yes" }, { value: "false", label: "No" }]} />
            <FilterSelect label="Active" value={activeFilter} onChange={(value) => { setActiveFilter(value); setPage(1); }} options={[{ value: "all", label: "All active" }, { value: "true", label: "Yes" }, { value: "false", label: "No" }]} />
          </>
        }
        onExportExcel={() =>
          exportExcel(
            "roles.xls",
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
        onExportPdf={printPdf}
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
        <Pagination
          page={paginated.page}
          pageSize={pageSize}
          totalPages={paginated.totalPages}
          totalItems={filtered.length}
          firstRecord={firstRecord}
          lastRecord={lastRecord}
          onPageChange={setPage}
          onPageSizeChange={(value) => {
            setPageSize(value);
            setPage(1);
          }}
        />
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
      <ConfirmableActionDialog
        pending={confirm.pending}
        isLoading={saveMutation.isPending || deleteMutation.isPending}
        onCancel={confirm.cancel}
        onConfirm={confirm.confirm}
      />
    </>
  );
}

function DataTable({ headers, rows, loading }: Readonly<{ headers: string[]; rows: React.ReactNode[][]; loading: boolean }>) {
  return <SharedDataTable headers={headers} rows={rows} loading={loading} loadingLabel="Loading..." />;
}

function Pagination({
  page,
  pageSize,
  totalPages,
  totalItems,
  firstRecord,
  lastRecord,
  onPageChange,
  onPageSizeChange,
}: Readonly<{
  page: number;
  pageSize: number;
  totalPages: number;
  totalItems: number;
  firstRecord: number;
  lastRecord: number;
  onPageChange: (page: number) => void;
  onPageSizeChange: (pageSize: number) => void;
}>) {
  return (
    <SharedPagination
      page={page}
      pageSize={pageSize}
      totalPages={totalPages}
      totalItems={totalItems}
      firstRecord={firstRecord}
      lastRecord={lastRecord}
      pageSizeOptions={pageSizeOptions}
      onPageChange={onPageChange}
      onPageSizeChange={onPageSizeChange}
    />
  );
}

function RowActions({ onEdit, onDelete }: Readonly<{ onEdit: () => void; onDelete: () => void }>) {
  return <SharedRowActions onEdit={onEdit} onDelete={onDelete} />;
}

function FilterSelect({
  label,
  value,
  options,
  withIcon,
  onChange,
}: Readonly<{ label: string; value: string; options: { value: string; label: string }[]; withIcon?: boolean; onChange: (value: string) => void }>) {
  return <SharedFilterSelect label={label} value={value} options={options} withIcon={withIcon} onChange={onChange} />;
}

function RoleBadges({ labels }: Readonly<{ labels: string[] }>) {
  if (labels.length === 0) {
    return <span className="text-muted-foreground">-</span>;
  }

  return (
    <div className="flex max-w-80 flex-wrap gap-1">
      {labels.map((label) => (
        <Badge key={label}>{label}</Badge>
      ))}
    </div>
  );
}

function Input({ label, value, disabled, onChange }: Readonly<{ label: string; value: string; disabled?: boolean; onChange: (value: string) => void }>) {
  return <TextInput label={label} value={value} disabled={disabled} required onChange={(event) => onChange(event.target.value)} />;
}

function Select({ label, value, disabled, options, onChange }: Readonly<{ label: string; value: string; disabled?: boolean; options: { value: string; label: string }[]; onChange: (value: string) => void }>) {
  return <SelectField label={label} value={value} disabled={disabled} options={options} onChange={(event) => onChange(event.target.value)} />;
}

function RoleDropdown({ label, values, options, onChange }: Readonly<{ label: string; values: string[]; options: { value: string; label: string }[]; onChange: (values: string[]) => void }>) {
  return <MultiSelect label={label} values={values} options={options} emptyLabel="Select roles" onChange={onChange} />;
}

function ActiveField<T extends { is_active: boolean }>({
  form,
  setForm,
}: {
  form: T;
  setForm: Dispatch<SetStateAction<T>>;
}) {
  return (
    <CheckboxField label="Active" checked={form.is_active} onChange={(event) => setForm((current) => ({ ...current, is_active: event.target.checked }))} />
  );
}

function FormActions({ onCancel, loading }: Readonly<{ onCancel: () => void; loading: boolean }>) {
  return (
    <div className="flex justify-end gap-2 border-t pt-4">
      <Button type="button" variant="outline" onClick={onCancel}>
        Cancel
      </Button>
      <Button type="submit" isLoading={loading} loadingLabel="Saving...">
        Save
      </Button>
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
