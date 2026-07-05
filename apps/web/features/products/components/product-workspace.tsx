"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { FormEvent, useMemo, useState } from "react";
import {
  createProduct,
  deleteProduct,
  listProducts,
  updateProduct,
  type Product,
} from "@/features/products/api";
import { Badge } from "@/shared/components/atoms/badge";
import { Button } from "@/shared/components/atoms/button";
import { CheckboxField } from "@/shared/components/atoms/checkbox-field";
import { SelectField } from "@/shared/components/atoms/select-field";
import { TextInput } from "@/shared/components/atoms/text-input";
import { ConfirmationDialog } from "@/shared/components/molecules/confirmation-dialog";
import { FilterSelect } from "@/shared/components/molecules/filter-select";
import { Modal } from "@/shared/components/molecules/modal";
import { Pagination } from "@/shared/components/molecules/pagination";
import { RowActions } from "@/shared/components/molecules/row-actions";
import { CrudToolbar } from "@/shared/components/organisms/crud-toolbar";
import { DataTable } from "@/shared/components/organisms/data-table";
import { PageShell } from "@/shared/components/organisms/page-shell";
import { exportExcel as exportExcelFile } from "@/shared/functions/export/export-excel";
import { printPdf } from "@/shared/functions/export/print-pdf";

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
      <PageShell
        title="Master Products"
        subtitle={`${filteredProducts.length} records, ${activeCount} active`}
        toolbar={
          <CrudToolbar
            search={search}
            searchPlaceholder="Search code, name, category"
            onSearch={updateSearch}
            onAdd={openCreateModal}
            addLabel="Add Product"
            filters={
              <>
                <FilterSelect
                  withIcon
                  label="Status"
                  value={statusFilter}
                  options={[
                    { value: "all", label: "All status" },
                    { value: "active", label: "Active" },
                    { value: "inactive", label: "Inactive" },
                    { value: "draft", label: "Draft" },
                  ]}
                  onChange={(value) => updateStatusFilter(value as "all" | Product["status"])}
                />
                <FilterSelect
                  label="Category"
                  value={categoryFilter}
                  options={[{ value: "all", label: "All categories" }, ...categories.map((category) => ({ value: category, label: category }))]}
                  onChange={updateCategoryFilter}
                />
              </>
            }
            onExportExcel={() =>
              exportExcelFile(
                "products.xls",
                ["Code", "Name", "Category", "Price", "Status", "Active"],
                filteredProducts.map((product) => [
                  product.code,
                  product.name,
                  product.category,
                  String(product.price_cents / 100),
                  product.status,
                  product.is_active ? "Active" : "Inactive",
                ]),
              )
            }
            onExportPdf={printPdf}
          />
        }
      >

        {error ? <p className="border-b px-4 py-3 text-sm text-destructive">{error.message}</p> : null}

        <DataTable
          headers={["Code", "Name", "Category", "Price", "Status", "Active", "Actions"]}
          loading={productsQuery.isLoading}
          loadingLabel="Loading products..."
          emptyTitle="No products found"
          emptyDescription="Adjust search/filter or create a new master data record."
          rows={paginatedProducts.map((product) => [
            <span key={`${product.id}-code`} className="font-semibold">{product.code}</span>,
            product.name,
            product.category,
            formatCurrency(product.price_cents),
            <Badge key={`${product.id}-status`} className="capitalize">{product.status}</Badge>,
            <span key={`${product.id}-active`} className={product.is_active ? "text-primary" : "text-muted-foreground"}>
              {product.is_active ? "Yes" : "No"}
            </span>,
            <RowActions
              key={`${product.id}-actions`}
              editLabel={`Edit ${product.name}`}
              deleteLabel={`Delete ${product.name}`}
              onEdit={() => openEditModal(product)}
              onDelete={() => setPendingAction({ type: "delete", product })}
            />,
          ])}
        />

        <Pagination
          page={currentPage}
          pageSize={pageSize}
          totalPages={totalPages}
          totalItems={filteredProducts.length}
          firstRecord={firstRecord}
          lastRecord={lastRecord}
          onPageChange={setPage}
          onPageSizeChange={(value) => {
            setPageSize(value);
            setPage(1);
          }}
        />
      </PageShell>

      {isModalOpen ? (
        <Modal title={editing ? "Edit Product" : "Create Product"} subtitle="Maintain product master data." onClose={closeModal}>
          <form className="space-y-4 p-5" onSubmit={onSubmit}>
            <TextInput
              id="code"
              label="Code"
              disabled={Boolean(editing)}
              required
              value={form.code}
              onChange={(event) => setForm((current) => ({ ...current, code: event.target.value.toUpperCase() }))}
            />
            <TextInput
              id="name"
              label="Name"
              required
              value={form.name}
              onChange={(event) => setForm((current) => ({ ...current, name: event.target.value }))}
            />
            <TextInput
              id="category"
              label="Category"
              required
              value={form.category}
              onChange={(event) => setForm((current) => ({ ...current, category: event.target.value }))}
            />
            <TextInput
              id="price"
              label="Price"
              type="number"
              min="0"
              step="100"
              required
              value={form.price}
              onChange={(event) => setForm((current) => ({ ...current, price: event.target.value }))}
            />
              {editing ? (
                <div className="grid grid-cols-2 gap-3">
                  <SelectField
                    id="status"
                    label="Status"
                    value={form.status}
                    options={[
                      { value: "active", label: "Active" },
                      { value: "inactive", label: "Inactive" },
                      { value: "draft", label: "Draft" },
                    ]}
                    onChange={(event) =>
                      setForm((current) => ({ ...current, status: event.target.value as Product["status"] }))
                    }
                  />
                  <CheckboxField
                    className="items-end pb-2"
                    label="Active"
                    checked={form.isActive}
                    onChange={(event) => setForm((current) => ({ ...current, isActive: event.target.checked }))}
                  />
                </div>
              ) : null}
              {error ? <p className="text-sm text-destructive">{error.message}</p> : null}
              <div className="flex justify-end gap-2 border-t pt-4">
                <Button type="button" variant="outline" onClick={closeModal}>
                  Cancel
                </Button>
                <Button
                  type="submit"
                  isLoading={isSaving}
                  loadingLabel="Saving..."
                >
                  {editing ? "Update" : "Create"}
                </Button>
              </div>
            </form>
        </Modal>
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
