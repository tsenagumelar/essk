"use client";

import {
  Bell,
  Check,
  Home,
  Info,
  LayoutDashboard,
  Package,
  Settings,
  User,
} from "lucide-react";
import { useState } from "react";
import { Avatar } from "@/shared/components/atoms/avatar";
import { Badge } from "@/shared/components/atoms/badge";
import { Button } from "@/shared/components/atoms/button";
import { CheckboxField } from "@/shared/components/atoms/checkbox-field";
import { HelperText } from "@/shared/components/atoms/helper-text";
import { IconButton } from "@/shared/components/atoms/icon-button";
import { Kbd } from "@/shared/components/atoms/kbd";
import { Label } from "@/shared/components/atoms/label";
import { LinkButton } from "@/shared/components/atoms/link-button";
import { RadioGroup } from "@/shared/components/atoms/radio-group";
import { SelectField } from "@/shared/components/atoms/select-field";
import { Separator } from "@/shared/components/atoms/separator";
import { Skeleton } from "@/shared/components/atoms/skeleton";
import { Spinner } from "@/shared/components/atoms/spinner";
import { SwitchField } from "@/shared/components/atoms/switch-field";
import { TextareaField } from "@/shared/components/atoms/textarea-field";
import { TextInput } from "@/shared/components/atoms/text-input";
import { Alert } from "@/shared/components/molecules/alert";
import { Breadcrumbs } from "@/shared/components/molecules/breadcrumbs";
import { ConfirmationDialog } from "@/shared/components/molecules/confirmation-dialog";
import { ConfirmableActionDialog } from "@/shared/components/molecules/confirmable-action-dialog";
import { DescriptionList } from "@/shared/components/molecules/description-list";
import { DropdownMenu, DropdownMenuItem } from "@/shared/components/molecules/dropdown-menu";
import { EmptyState } from "@/shared/components/molecules/empty-state";
import { FilterSelect } from "@/shared/components/molecules/filter-select";
import { LoadingState } from "@/shared/components/molecules/loading-state";
import { Modal } from "@/shared/components/molecules/modal";
import { MultiSelect } from "@/shared/components/molecules/multi-select";
import { Pagination } from "@/shared/components/molecules/pagination";
import { RowActions } from "@/shared/components/molecules/row-actions";
import { SearchBox } from "@/shared/components/molecules/search-box";
import { StatusBadge } from "@/shared/components/molecules/status-badge";
import { Tabs } from "@/shared/components/molecules/tabs";
import { AppNavbar } from "@/shared/components/organisms/app-navbar";
import { AppSidebar } from "@/shared/components/organisms/app-sidebar";
import { CrudToolbar } from "@/shared/components/organisms/crud-toolbar";
import { DataTable } from "@/shared/components/organisms/data-table";
import { DetailPanel } from "@/shared/components/organisms/detail-panel";
import { FormPanel } from "@/shared/components/organisms/form-panel";
import { PageShell } from "@/shared/components/organisms/page-shell";
import { SectionHeader } from "@/shared/components/organisms/section-header";
import { StatGrid } from "@/shared/components/organisms/stat-grid";
import { AuthenticatedLayout } from "@/shared/components/templates/authenticated-layout";
import { CrudPageTemplate } from "@/shared/components/templates/crud-page-template";
import { SplitAuthLayout } from "@/shared/components/templates/split-auth-layout";
import { useConfirmableAction } from "@/shared/hooks/use-confirmable-action";

const options = [
  { value: "active", label: "Active" },
  { value: "inactive", label: "Inactive" },
  { value: "draft", label: "Draft" },
];

const sidebarGroups = [
  {
    label: "Workspace",
    items: [
      { href: "/dashboard", label: "Dashboard", icon: <LayoutDashboard className="h-4 w-4 shrink-0" /> },
      { href: "/components", label: "Components", icon: <Package className="h-4 w-4 shrink-0" /> },
    ],
  },
  {
    label: "Settings",
    items: [{ href: "/health", label: "Health", icon: <Settings className="h-4 w-4 shrink-0" /> }],
  },
];

export function ComponentPreviewView() {
  const [text, setText] = useState("ESSK");
  const [textarea, setTextarea] = useState("Reusable textarea content");
  const [selectValue, setSelectValue] = useState("active");
  const [checked, setChecked] = useState(true);
  const [radio, setRadio] = useState("admin");
  const [switchValue, setSwitchValue] = useState(true);
  const [search, setSearch] = useState("");
  const [multiValues, setMultiValues] = useState(["active", "draft"]);
  const [tab, setTab] = useState("overview");
  const [dropdownOpen, setDropdownOpen] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const confirmable = useConfirmableAction();

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-semibold">Component Preview</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Preview semua base atomic design components yang tersedia di shared components.
        </p>
      </div>

      <PreviewSection title="Atoms" description="Primitive UI components yang menjadi building block utama.">
        <PreviewCard title="Buttons">
          <div className="flex flex-wrap gap-2">
            <Button>Primary</Button>
            <Button variant="secondary">Secondary</Button>
            <Button variant="outline">Outline</Button>
            <Button variant="ghost">Ghost</Button>
            <Button variant="danger">Danger</Button>
            <Button isLoading loadingLabel="Saving...">Loading</Button>
            <IconButton aria-label="Notification">
              <Bell className="h-4 w-4" />
            </IconButton>
            <LinkButton href="/dashboard" variant="outline">Link Button</LinkButton>
          </div>
        </PreviewCard>

        <PreviewCard title="Form Controls">
          <div className="grid gap-4 md:grid-cols-2">
            <TextInput label="Text Input" value={text} onChange={(event) => setText(event.target.value)} />
            <SelectField label="Select Field" value={selectValue} options={options} onChange={(event) => setSelectValue(event.target.value)} />
            <TextareaField label="Textarea Field" value={textarea} onChange={(event) => setTextarea(event.target.value)} />
            <div className="space-y-3">
              <Label>Choice Controls</Label>
              <CheckboxField label="Checkbox field" checked={checked} onChange={(event) => setChecked(event.target.checked)} />
              <SwitchField label="Switch field" checked={switchValue} onCheckedChange={setSwitchValue} />
              <RadioGroup
                name="role-preview"
                value={radio}
                options={[
                  { value: "super_admin", label: "Super Admin" },
                  { value: "admin", label: "Admin" },
                  { value: "user", label: "User" },
                ]}
                onChange={setRadio}
              />
            </div>
          </div>
        </PreviewCard>

        <PreviewCard title="Display Atoms">
          <div className="flex flex-wrap items-center gap-3">
            <Avatar name="ESSK Admin" />
            <Badge>Default</Badge>
            <Badge variant="success">Success</Badge>
            <Badge variant="danger">Danger</Badge>
            <StatusBadge active />
            <Spinner />
            <Kbd>Cmd</Kbd>
            <Kbd>K</Kbd>
            <HelperText>Helper text for field guidance.</HelperText>
          </div>
          <Separator className="my-4" />
          <div className="grid gap-2 md:grid-cols-3">
            <Skeleton className="h-6" />
            <Skeleton className="h-6" />
            <Skeleton className="h-6" />
          </div>
        </PreviewCard>
      </PreviewSection>

      <PreviewSection title="Molecules" description="Gabungan kecil dari atoms untuk pattern UI yang sering dipakai.">
        <PreviewCard title="Search, Filter, Multi Select">
          <div className="flex flex-col gap-3 lg:flex-row lg:items-center">
            <SearchBox value={search} onChange={setSearch} placeholder="Search preview" className="lg:w-72" />
            <FilterSelect label="Status" value={selectValue} options={options} withIcon onChange={setSelectValue} />
            <MultiSelect label="Multi Select" values={multiValues} options={options} onChange={setMultiValues} />
          </div>
        </PreviewCard>

        <PreviewCard title="Feedback">
          <div className="grid gap-3 md:grid-cols-2">
            <Alert title="Info">Informational alert message.</Alert>
            <Alert variant="success" title="Success">Operation completed successfully.</Alert>
            <Alert variant="warning" title="Warning">Review this configuration before production.</Alert>
            <Alert variant="danger" title="Danger">Something requires attention.</Alert>
          </div>
          <LoadingState label="Loading molecule preview..." />
          <EmptyState title="No records" description="Empty state for list and table screens." actionLabel="Create" onAction={() => setModalOpen(true)} />
        </PreviewCard>

        <PreviewCard title="Navigation And Details">
          <div className="space-y-4">
            <Breadcrumbs items={[{ label: "Dashboard", href: "/dashboard" }, { label: "Components" }]} />
            <Tabs
              value={tab}
              items={[
                { value: "overview", label: "Overview" },
                { value: "usage", label: "Usage" },
                { value: "api", label: "API" },
              ]}
              onChange={setTab}
            />
            <DescriptionList
              items={[
                { label: "Component", value: "DescriptionList" },
                { label: "Current Tab", value: tab },
                { label: "Status", value: <StatusBadge active /> },
              ]}
            />
          </div>
        </PreviewCard>

        <PreviewCard title="Actions And Overlays">
          <div className="flex flex-wrap items-center gap-3">
            <RowActions onEdit={() => setModalOpen(true)} onDelete={() => setConfirmOpen(true)} />
            <div className="relative">
              <Button variant="outline" onClick={() => setDropdownOpen((current) => !current)}>Dropdown</Button>
              <DropdownMenu open={dropdownOpen}>
                <DropdownMenuItem icon={<User className="h-4 w-4" />}>Profile</DropdownMenuItem>
                <DropdownMenuItem icon={<Settings className="h-4 w-4" />}>Settings</DropdownMenuItem>
              </DropdownMenu>
            </div>
            <Button onClick={() => setModalOpen(true)}>Open Modal</Button>
            <Button
              variant="danger"
              onClick={() =>
                confirmable.request({
                  title: "Run confirmable action?",
                  description: "This preview uses the shared confirmable action hook and dialog.",
                  confirmLabel: "Run",
                  run: async () => undefined,
                })
              }
            >
              Confirmable Hook
            </Button>
          </div>
        </PreviewCard>
      </PreviewSection>

      <PreviewSection title="Organisms" description="Reusable sections untuk page-level composition.">
        <PreviewCard title="Page Shell And CRUD Toolbar">
          <PageShell
            title="Page Shell"
            subtitle="Header, subtitle, toolbar, and content area."
            toolbar={
              <CrudToolbar
                search={search}
                onSearch={setSearch}
                onAdd={() => setModalOpen(true)}
                filters={<FilterSelect label="Status" value={selectValue} options={options} onChange={setSelectValue} />}
                onExportExcel={() => undefined}
                onExportPdf={() => undefined}
              />
            }
          >
            <DataTable
              headers={["Name", "Status", "Owner", "Actions"]}
              rows={[
                ["Component Preview", <StatusBadge key="status" active />, "ESSK", <RowActions key="actions" onEdit={() => setModalOpen(true)} onDelete={() => setConfirmOpen(true)} />],
                ["Atomic Design", <Badge key="badge">Draft</Badge>, "Platform", <RowActions key="actions-2" onEdit={() => setModalOpen(true)} onDelete={() => setConfirmOpen(true)} />],
              ]}
            />
            <Pagination
              page={page}
              pageSize={pageSize}
              totalPages={3}
              totalItems={25}
              firstRecord={1}
              lastRecord={10}
              onPageChange={setPage}
              onPageSizeChange={setPageSize}
            />
          </PageShell>
        </PreviewCard>

        <PreviewCard title="Panels And Stats">
          <div className="space-y-4">
            <SectionHeader title="Section Header" subtitle="Reusable section title and actions." actions={<Button size="sm">Action</Button>} />
            <StatGrid
              items={[
                { label: "Atoms", value: "17", description: "Primitive components", icon: <Check className="h-4 w-4" /> },
                { label: "Molecules", value: "16", description: "Composed controls", icon: <Info className="h-4 w-4" /> },
                { label: "Organisms", value: "9", description: "Page sections", icon: <Package className="h-4 w-4" /> },
                { label: "Templates", value: "3", description: "Layout patterns", icon: <Home className="h-4 w-4" /> },
              ]}
            />
            <DetailPanel
              title="Detail Panel"
              description="Structured detail view."
              items={[
                { label: "Name", value: "Enterprise SaaS Starter Kit" },
                { label: "Mode", value: <Badge>Preview</Badge> },
              ]}
            />
            <form onSubmit={(event) => event.preventDefault()}>
              <FormPanel title="Form Panel" description="Reusable form shell." onCancel={() => undefined}>
                <TextInput label="Name" defaultValue="Sample Form" />
                <TextareaField label="Description" defaultValue="This form panel can wrap future generated forms." />
              </FormPanel>
            </form>
          </div>
        </PreviewCard>

        <PreviewCard title="App Layout Parts">
          <div className="overflow-hidden rounded-xl border bg-slate-100">
            <AppNavbar
              left={<span className="font-semibold">ESSK</span>}
              search={<SearchBox value={search} onChange={setSearch} placeholder="Navbar search" className="max-w-md flex-1" />}
              right={<IconButton aria-label="Notification"><Bell className="h-4 w-4" /></IconButton>}
            />
            <div className="flex">
              <AppSidebar
                position="static"
                groups={sidebarGroups}
                activePath="/components"
                brand={<div className="border-b px-4 py-4 font-semibold">Sidebar</div>}
                footer={<div className="border-t p-4 text-xs text-muted-foreground">Footer</div>}
              />
              <div className="flex-1 p-4 text-sm text-muted-foreground">Static preview area for AppNavbar and AppSidebar.</div>
            </div>
          </div>
        </PreviewCard>
      </PreviewSection>

      <PreviewSection title="Templates" description="High-level layout composition templates.">
        <PreviewCard title="CRUD Page Template">
          <CrudPageTemplate
            title="CRUD Template"
            subtitle="Template wrapper for generated CRUD modules."
            toolbar={<Button size="sm">Template Action</Button>}
            table={<DataTable headers={["Column", "Value"]} rows={[["Generated", "Ready"]]} />}
            pagination={<Pagination page={1} pageSize={10} totalPages={1} totalItems={1} firstRecord={1} lastRecord={1} onPageChange={() => undefined} onPageSizeChange={() => undefined} />}
          />
        </PreviewCard>

        <PreviewCard title="Split Auth Layout">
          <div className="overflow-hidden rounded-xl border">
            <SplitAuthLayout media={<div className="flex h-full min-h-64 items-center justify-center bg-primary text-primary-foreground">Media Side</div>}>
              <div className="w-full max-w-sm rounded-lg border bg-white p-4">
                <p className="text-sm font-semibold">Auth Form Side</p>
                <TextInput className="mt-3" placeholder="email@example.com" />
              </div>
            </SplitAuthLayout>
          </div>
        </PreviewCard>

        <PreviewCard title="Authenticated Layout">
          <div className="max-h-[520px] overflow-hidden rounded-xl border">
            <AuthenticatedLayout
              sidebarCollapsed
              sidebar={<div className="hidden border-r bg-white p-4 lg:block">Sidebar slot</div>}
              navbar={<div className="border-b bg-white p-4">Navbar slot</div>}
            >
              <div className="rounded-lg bg-white p-4">Authenticated layout content slot.</div>
            </AuthenticatedLayout>
          </div>
        </PreviewCard>
      </PreviewSection>

      {modalOpen ? (
        <Modal title="Preview Modal" subtitle="This modal is rendered from shared components." onClose={() => setModalOpen(false)}>
          <div className="space-y-4 p-5">
            <TextInput label="Modal Field" defaultValue="Editable value" />
            <div className="flex justify-end gap-2 border-t pt-4">
              <Button variant="outline" onClick={() => setModalOpen(false)}>Cancel</Button>
              <Button onClick={() => setModalOpen(false)}>Save</Button>
            </div>
          </div>
        </Modal>
      ) : null}

      <ConfirmationDialog
        open={confirmOpen}
        title="Delete preview item?"
        description="This is the shared confirmation dialog."
        confirmLabel="Delete"
        variant="danger"
        onCancel={() => setConfirmOpen(false)}
        onConfirm={() => setConfirmOpen(false)}
      />
      <ConfirmableActionDialog
        pending={confirmable.pending}
        onCancel={confirmable.cancel}
        onConfirm={confirmable.confirm}
      />
    </div>
  );
}

function PreviewSection({
  title,
  description,
  children,
}: Readonly<{ title: string; description: string; children: React.ReactNode }>) {
  return (
    <section className="space-y-4">
      <div>
        <h2 className="text-xl font-semibold">{title}</h2>
        <p className="mt-1 text-sm text-muted-foreground">{description}</p>
      </div>
      <div className="grid gap-4">{children}</div>
    </section>
  );
}

function PreviewCard({ title, children }: Readonly<{ title: string; children: React.ReactNode }>) {
  return (
    <section className="rounded-xl border bg-white p-4 shadow-sm">
      <h3 className="mb-4 text-sm font-semibold uppercase tracking-wide text-slate-500">{title}</h3>
      {children}
    </section>
  );
}
