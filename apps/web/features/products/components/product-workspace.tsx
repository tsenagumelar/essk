"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { FormEvent, useState } from "react";
import {
  createProduct,
  deleteProduct,
  listProducts,
  updateProduct,
  type Product,
} from "@/features/products/api";

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

export function ProductWorkspace() {
  const queryClient = useQueryClient();
  const [editing, setEditing] = useState<Product | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);

  const productsQuery = useQuery({
    queryKey: ["products"],
    queryFn: listProducts,
  });

  const products = productsQuery.data ?? [];
  const activeCount = products.filter((product) => product.is_active).length;

  const createMutation = useMutation({
    mutationFn: createProduct,
    onSuccess: async () => {
      setForm(emptyForm);
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
      setEditing(null);
      setForm(emptyForm);
      await queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteProduct,
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });

  function startEdit(product: Product) {
    setEditing(product);
    setForm({
      code: product.code,
      name: product.name,
      category: product.category,
      price: String(product.price_cents / 100),
      status: product.status,
      isActive: product.is_active,
    });
  }

  function resetForm() {
    setEditing(null);
    setForm(emptyForm);
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (editing) {
      await updateMutation.mutateAsync({
        id: editing.id,
        product: {
          ...editing,
          name: form.name,
          category: form.category,
          price_cents: toPriceCents(form.price),
          status: form.status,
          is_active: form.isActive,
        },
      });
      return;
    }

    await createMutation.mutateAsync({
      code: form.code,
      name: form.name,
      category: form.category,
      price_cents: toPriceCents(form.price),
    });
  }

  const isSaving = createMutation.isPending || updateMutation.isPending;
  const error = productsQuery.error ?? createMutation.error ?? updateMutation.error ?? deleteMutation.error;

  return (
    <div className="grid gap-6 lg:grid-cols-[360px_1fr]">
      <section className="rounded-lg border bg-white p-4">
        <h2 className="text-base font-semibold">{editing ? "Edit Product" : "Create Product"}</h2>
        <form className="mt-4 space-y-4" onSubmit={onSubmit}>
          <div>
            <label className="text-sm font-medium" htmlFor="code">
              Code
            </label>
            <input
              id="code"
              disabled={Boolean(editing)}
              required
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm outline-none ring-primary focus:ring-2 disabled:bg-muted"
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
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm outline-none ring-primary focus:ring-2"
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
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm outline-none ring-primary focus:ring-2"
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
              className="mt-1 w-full rounded-md border px-3 py-2 text-sm outline-none ring-primary focus:ring-2"
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
                  className="mt-1 w-full rounded-md border px-3 py-2 text-sm outline-none ring-primary focus:ring-2"
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
          <div className="flex gap-2">
            <button
              type="submit"
              disabled={isSaving}
              className="rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60"
            >
              {isSaving ? "Saving..." : editing ? "Update" : "Create"}
            </button>
            {editing ? (
              <button type="button" className="rounded-md border px-3 py-2 text-sm font-medium" onClick={resetForm}>
                Cancel
              </button>
            ) : null}
          </div>
        </form>
      </section>

      <section className="rounded-lg border bg-white">
        <div className="flex items-center justify-between border-b px-4 py-3">
          <div>
            <h2 className="text-base font-semibold">Master Products</h2>
            <p className="text-sm text-muted-foreground">
              {products.length} records, {activeCount} active
            </p>
          </div>
        </div>
        {productsQuery.isLoading ? (
          <p className="p-4 text-sm text-muted-foreground">Loading products...</p>
        ) : products.length === 0 ? (
          <p className="p-4 text-sm text-muted-foreground">No products yet. Create the first master data record.</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-sm">
              <thead className="bg-muted text-xs uppercase text-muted-foreground">
                <tr>
                  <th className="px-4 py-3">Code</th>
                  <th className="px-4 py-3">Name</th>
                  <th className="px-4 py-3">Category</th>
                  <th className="px-4 py-3">Price</th>
                  <th className="px-4 py-3">Status</th>
                  <th className="px-4 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                {products.map((product) => (
                  <tr key={product.id} className="border-t">
                    <td className="px-4 py-3 font-medium">{product.code}</td>
                    <td className="px-4 py-3">{product.name}</td>
                    <td className="px-4 py-3">{product.category}</td>
                    <td className="px-4 py-3">{formatCurrency(product.price_cents)}</td>
                    <td className="px-4 py-3">
                      <span className="rounded-full bg-muted px-2 py-1 text-xs">{product.status}</span>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <button className="mr-2 rounded-md border px-2 py-1 text-xs" onClick={() => startEdit(product)}>
                        Edit
                      </button>
                      <button
                        className="rounded-md border px-2 py-1 text-xs text-destructive"
                        disabled={deleteMutation.isPending}
                        onClick={() => deleteMutation.mutate(product.id)}
                      >
                        Delete
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
