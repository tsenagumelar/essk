"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Download, FileText, Filter, Pencil, Plus, Search, Trash2, X } from "lucide-react";
import { FormEvent, useMemo, useState } from "react";
import {
  createProduct,
  deleteProduct,
  listProducts,
  updateProduct,
  type Product,
} from "@/features/products/api";
import { ConfirmationDialog } from "@/features/shared/components/confirmation-dialog";

type FormState = {
  code: string;
  name: string;
  category: string;
  price: string;
  status: Product["status"];
  isActive: boolean;
};

const emptyForm: FormState = {
  code: "",
  name: "",
  category: "",
  price: "0",
  status: "active",
  isActive: true,
};

const pageSizeOptions = [5, 10, 20];
const emptyProducts: Product[] = [];

type PendingAction =
  | { type: "create" }
  | { type: "update"; product: Product }
  | { type: "delete"; product: Product };

function formatCurrency(priceCents: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(priceCents / 100);
}

function toPriceCents(value: string) {
  const parsed = Number(value);
  if (Number.isNaN(parsed) || parsed < 0) {
    return 0;
  }
  return Math.round(parsed * 100);
}

function productMatchesSearch(product: Product, search: string) {
  const query = search.trim().toLowerCase();
  if (!query) {
    return true;
  }
  return [product.code, product.name, product.category, product.status].some((value) => value.toLowerCase().includes(query));
}

function exportExcel(products: Product[]) {
  const rows = products.map((product) => `
    <tr>
      <td>${product.code}</td>
      <td>${product.name}</td>
      <td>${product.category}</td>
      <td>${product.price_cents / 100}</td>
      <td>${product.status}</td>
      <td>${product.is_active ? "Active" : "Inactive"}</td>
    </tr>
  `);
  const workbook = `
    <table>
      <thead>
        <tr>
          <th>Code</th>
          <th>Name</th>
          <th>Category</th>
          <th>Price</th>
          <th>Status</th>
          <th>Active</th>
        </tr>
      </thead>
      <tbody>${rows.join("")}</tbody>
    </table>
  `;
  const blob = new Blob([workbook], { type: "application/vnd.ms-excel" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = "products.xls";
  link.click();
  URL.revokeObjectURL(url);
}

function exportPdf() {
  window.print();
}

export function ProductWorkspace() {
  const queryClient = useQueryClient();
  const [editing, setEditing] = useState<Product | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<"all" | Product["status"]>("all");
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [pendingAction, setPendingAction] = useState<PendingAction | null>(null);

  const productsQuery = useQuery({
    queryKey: ["products"],
    queryFn: listProducts,
  });

  const products = productsQuery.data ?? emptyProducts;
  const categories = useMemo(() => Array.from(new Set(products.map((product) => product.category))).sort(), [products]);
  const [categoryFilter, setCategoryFilter] = useState("all");

  const filteredProducts = useMemo(() => {
    return products.filter((product) => {
      const matchesSearch = productMatchesSearch(product, search);
      const matchesStatus = statusFilter === "all" || product.status === statusFilter;
      const matchesCategory = categoryFilter === "all" || product.category === categoryFilter;
      return matchesSearch && matchesStatus && matchesCategory;
    });
  }, [categoryFilter, products, search, statusFilter]);

  const totalPages = Math.max(1, Math.ceil(filteredProducts.length / pageSize));
  const currentPage = Math.min(page, totalPages);
  const paginatedProducts = filteredProducts.slice((currentPage - 1) * pageSize, currentPage * pageSize);
  const activeCount = filteredProducts.filter((product) => product.is_active).length;

  const createMutation = useMutation({
    mutationFn: createProduct,
    onSuccess: async () => {
      closeModal();
      await queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, product }: { id: string; product: Product }) =>
      updateProduct(id, {
        name: product.name,
        category: product.category,
        price_cents: product.price_cents,
        status: product.status,
        is_active: product.is_active,
      }),
    onSuccess: async () => {
      closeModal();
      await queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteProduct,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });

  function openCreateModal() {
    setEditing(null);
    setForm(emptyForm);
    setIsModalOpen(true);
  }

  function openEditModal(product: Product) {
    setEditing(product);
    setForm({
      code: product.code,
      name: product.name,
      category: product.category,
      price: String(product.price_cents / 100),
      status: product.status,
      isActive: product.is_active,
    });
    setIsModalOpen(true);
  }

  function closeModal() {
    setEditing(null);
    setForm(emptyForm);
    setIsModalOpen(false);
  }

  function updateSearch(value: string) {
    setSearch(value);
    setPage(1);
  }

  function updateStatusFilter(value: "all" | Product["status"]) {
    setStatusFilter(value);
    setPage(1);
  }

  function updateCategoryFilter(value: string) {
    setCategoryFilter(value);
    setPage(1);
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (editing) {
      setPendingAction({ type: "update", product: editing });
      return;
    }

    setPendingAction({ type: "create" });
  }

  async function confirmPendingAction() {
    if (!pendingAction) {
      return;
    }

    if (pendingAction.type === "create") {
      await createMutation.mutateAsync({
        code: form.code,
        name: form.name,
        category: form.category,
        price_cents: toPriceCents(form.price),
      });
      setPendingAction(null);
      return;
    }

    if (pendingAction.type === "update") {
      await updateMutation.mutateAsync({
        id: pendingAction.product.id,
        product: {
          ...pendingAction.product,
          name: form.name,
          category: form.category,
          price_cents: toPriceCents(form.price),
          status: form.status,
          is_active: form.isActive,
        },
      });
      setPendingAction(null);
      return;
    }

    await deleteMutation.mutateAsync(pendingAction.product.id);
    setPendingAction(null);
  }

  function cancelPendingAction() {
    setPendingAction(null);
  }

  function getConfirmationContent(action: PendingAction | null) {
    if (!action) {
      return {
        title: "",
        description: "",
        confirmLabel: "Confirm",
        variant: "primary" as const,
      };
    }

    if (action.type === "create") {
      return {
        title: "Save new product?",
        description: `This will add "${form.name || form.code || "new product"}" to master data.`,
        confirmLabel: "Save Add",
        variant: "primary" as const,
      };
    }

    if (action.type === "update") {
      return {
        title: "Save product changes?",
        description: `This will update "${action.product.name}" with the latest form values.`,
        confirmLabel: "Save Edit",
        variant: "primary" as const,
      };
    }

    return {
      title: "Delete product?",
      description: `This will soft delete "${action.product.name}" from master data.`,
      confirmLabel: "Delete",
      variant: "danger" as const,
    };
  }

  const isSaving = createMutation.isPending || updateMutation.isPending;
  const isProcessing = isSaving || deleteMutation.isPending;
  const error = productsQuery.error ?? createMutation.error ?? updateMutation.error ?? deleteMutation.error;
  const firstRecord = filteredProducts.length === 0 ? 0 : (currentPage - 1) * pageSize + 1;
  const lastRecord = Math.min(currentPage * pageSize, filteredProducts.length);
  const confirmation = getConfirmationContent(pendingAction);

  return (
    <>
      <section className="overflow-hidden rounded-xl border bg-white shadow-sm">
        <div className="border-b bg-gradient-to-r from-white to-slate-50 px-4 py-4">
          <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
            <div>
              <h2 className="text-base font-semibold">Master Products</h2>
              <p className="mt-1 text-sm text-muted-foreground">
                {filteredProducts.length} records, {activeCount} active
              </p>
            </div>

            <div className="flex flex-col gap-3 lg:flex-row lg:items-center">
              <div className="relative min-w-0 lg:w-72">
                <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <input
                  type="search"
                  placeholder="Search code, name, category"
                  className="h-10 w-full rounded-lg border bg-white pl-9 pr-3 text-sm outline-none ring-primary focus:ring-2"
                  value={search}
                  onChange={(event) => updateSearch(event.target.value)}
                />
              </div>

              <div className="flex flex-wrap gap-2">
                <div className="flex h-10 items-center gap-2 rounded-lg border bg-white px-3">
                  <Filter className="h-4 w-4 text-muted-foreground" />
                  <select
                    className="bg-transparent text-sm outline-none"
                    value={statusFilter}
                    onChange={(event) => updateStatusFilter(event.target.value as "all" | Product["status"])}
                  >
                    <option value="all">All status</option>
                    <option value="active">Active</option>
                    <option value="inactive">Inactive</option>
                    <option value="draft">Draft</option>
                  </select>
                </div>
                <select
                  className="h-10 rounded-lg border bg-white px-3 text-sm outline-none ring-primary focus:ring-2"
                  value={categoryFilter}
                  onChange={(event) => updateCategoryFilter(event.target.value)}
                >
                  <option value="all">All categories</option>
                  {categories.map((category) => (
                    <option key={category} value={category}>
                      {category}
                    </option>
                  ))}
                </select>
                <button
                  className="inline-flex h-10 items-center gap-2 rounded-lg border bg-white px-3 text-sm font-medium hover:bg-slate-50"
                  onClick={() => exportExcel(filteredProducts)}
                >
                  <Download className="h-4 w-4" />
                  Excel
                </button>
                <button className="inline-flex h-10 items-center gap-2 rounded-lg border bg-white px-3 text-sm font-medium hover:bg-slate-50" onClick={exportPdf}>
                  <FileText className="h-4 w-4" />
                  PDF
                </button>
                <button
                  className="inline-flex h-10 items-center gap-2 rounded-lg bg-primary px-3 text-sm font-medium text-primary-foreground shadow-sm hover:opacity-95"
                  onClick={openCreateModal}
                >
                  <Plus className="h-4 w-4" />
                  Add Product
                </button>
              </div>
            </div>
          </div>
        </div>

        {error ? <p className="border-b px-4 py-3 text-sm text-destructive">{error.message}</p> : null}

        {productsQuery.isLoading ? (
          <p className="p-4 text-sm text-muted-foreground">Loading products...</p>
        ) : filteredProducts.length === 0 ? (
          <div className="p-10 text-center">
            <p className="text-sm font-medium">No products found</p>
            <p className="mt-1 text-sm text-muted-foreground">Adjust search/filter or create a new master data record.</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase text-muted-foreground">
                <tr>
                  <th className="px-4 py-3">Code</th>
                  <th className="px-4 py-3">Name</th>
                  <th className="px-4 py-3">Category</th>
                  <th className="px-4 py-3">Price</th>
                  <th className="px-4 py-3">Status</th>
                  <th className="px-4 py-3">Active</th>
                  <th className="px-4 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {paginatedProducts.map((product) => (
                  <tr key={product.id} className="border-t hover:bg-slate-50/80">
                    <td className="px-4 py-3 font-semibold">{product.code}</td>
                    <td className="px-4 py-3">{product.name}</td>
                    <td className="px-4 py-3">{product.category}</td>
                    <td className="px-4 py-3">{formatCurrency(product.price_cents)}</td>
                    <td className="px-4 py-3">
                      <span className="rounded-full bg-slate-100 px-2 py-1 text-xs capitalize">{product.status}</span>
                    </td>
                    <td className="px-4 py-3">
                      <span className={product.is_active ? "text-primary" : "text-muted-foreground"}>
                        {product.is_active ? "Yes" : "No"}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <button
                        className="mr-2 inline-flex h-8 w-8 items-center justify-center rounded-lg border hover:bg-slate-50"
                        aria-label={`Edit ${product.name}`}
                        onClick={() => openEditModal(product)}
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        className="inline-flex h-8 w-8 items-center justify-center rounded-lg border text-destructive hover:bg-slate-50"
                        aria-label={`Delete ${product.name}`}
                        disabled={deleteMutation.isPending}
                        onClick={() => setPendingAction({ type: "delete", product })}
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <div className="flex flex-col gap-3 border-t px-4 py-3 text-sm text-muted-foreground md:flex-row md:items-center md:justify-between">
          <p>
            Showing {firstRecord}-{lastRecord} of {filteredProducts.length}
          </p>
          <div className="flex items-center gap-2">
            <select
              className="h-9 rounded-lg border bg-white px-2 text-sm outline-none"
              value={pageSize}
              onChange={(event) => {
                setPageSize(Number(event.target.value));
                setPage(1);
              }}
            >
              {pageSizeOptions.map((size) => (
                <option key={size} value={size}>
                  {size} / page
                </option>
              ))}
            </select>
            <button
              className="h-9 rounded-lg border bg-white px-3 disabled:opacity-50"
              disabled={currentPage <= 1}
              onClick={() => setPage((current) => Math.max(1, current - 1))}
            >
              Prev
            </button>
            <span className="px-2">
              Page {currentPage} of {totalPages}
            </span>
            <button
              className="h-9 rounded-lg border bg-white px-3 disabled:opacity-50"
              disabled={currentPage >= totalPages}
              onClick={() => setPage((current) => Math.min(totalPages, current + 1))}
            >
              Next
            </button>
          </div>
        </div>
      </section>

      {isModalOpen ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
          <section className="w-full max-w-lg rounded-xl bg-white shadow-xl">
            <div className="flex items-center justify-between border-b px-5 py-4">
              <div>
                <h2 className="text-base font-semibold">{editing ? "Edit Product" : "Create Product"}</h2>
                <p className="text-sm text-muted-foreground">Maintain product master data.</p>
              </div>
              <button className="rounded-lg p-2 hover:bg-slate-100" onClick={closeModal} aria-label="Close form">
                <X className="h-4 w-4" />
              </button>
            </div>

            <form className="space-y-4 p-5" onSubmit={onSubmit}>
              <div>
                <label className="text-sm font-medium" htmlFor="code">
                  Code
                </label>
                <input
                  id="code"
                  disabled={Boolean(editing)}
                  required
                  className="mt-1 h-10 w-full rounded-lg border px-3 text-sm outline-none ring-primary focus:ring-2 disabled:bg-muted"
                  value={form.code}
                  onChange={(event) => setForm((current) => ({ ...current, code: event.target.value.toUpperCase() }))}
                />
              </div>
              <div>
                <label className="text-sm font-medium" htmlFor="name">
                  Name
                </label>
                <input
                  id="name"
                  required
                  className="mt-1 h-10 w-full rounded-lg border px-3 text-sm outline-none ring-primary focus:ring-2"
                  value={form.name}
                  onChange={(event) => setForm((current) => ({ ...current, name: event.target.value }))}
                />
              </div>
              <div>
                <label className="text-sm font-medium" htmlFor="category">
                  Category
                </label>
                <input
                  id="category"
                  required
                  className="mt-1 h-10 w-full rounded-lg border px-3 text-sm outline-none ring-primary focus:ring-2"
                  value={form.category}
                  onChange={(event) => setForm((current) => ({ ...current, category: event.target.value }))}
                />
              </div>
              <div>
                <label className="text-sm font-medium" htmlFor="price">
                  Price
                </label>
                <input
                  id="price"
                  type="number"
                  min="0"
                  step="100"
                  required
                  className="mt-1 h-10 w-full rounded-lg border px-3 text-sm outline-none ring-primary focus:ring-2"
                  value={form.price}
                  onChange={(event) => setForm((current) => ({ ...current, price: event.target.value }))}
                />
              </div>
              {editing ? (
                <div className="grid grid-cols-2 gap-3">
                  <div>
                    <label className="text-sm font-medium" htmlFor="status">
                      Status
                    </label>
                    <select
                      id="status"
                      className="mt-1 h-10 w-full rounded-lg border px-3 text-sm outline-none ring-primary focus:ring-2"
                      value={form.status}
                      onChange={(event) =>
                        setForm((current) => ({ ...current, status: event.target.value as Product["status"] }))
                      }
                    >
                      <option value="active">Active</option>
                      <option value="inactive">Inactive</option>
                      <option value="draft">Draft</option>
                    </select>
                  </div>
                  <label className="flex items-end gap-2 pb-2 text-sm">
                    <input
                      type="checkbox"
                      checked={form.isActive}
                      onChange={(event) => setForm((current) => ({ ...current, isActive: event.target.checked }))}
                    />
                    Active
                  </label>
                </div>
              ) : null}
              {error ? <p className="text-sm text-destructive">{error.message}</p> : null}
              <div className="flex justify-end gap-2 border-t pt-4">
                <button type="button" className="rounded-lg border px-3 py-2 text-sm font-medium" onClick={closeModal}>
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSaving}
                  className="rounded-lg bg-primary px-3 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60"
                >
                  {isSaving ? "Saving..." : editing ? "Update" : "Create"}
                </button>
              </div>
            </form>
          </section>
        </div>
      ) : null}

      <ConfirmationDialog
        open={Boolean(pendingAction)}
        title={confirmation.title}
        description={confirmation.description}
        confirmLabel={confirmation.confirmLabel}
        variant={confirmation.variant}
        isLoading={isProcessing}
        onCancel={cancelPendingAction}
        onConfirm={confirmPendingAction}
      />
    </>
  );
}
