import { apiDelete, apiGet, apiPatch, apiPost } from "@/shared/api/client";

export type Product = {
  id: string;
  tenant_id: string;
  code: string;
  name: string;
  category: string;
  price_cents: number;
  status: "active" | "inactive" | "draft";
  is_active: boolean;
};

export type CreateProductPayload = {
  code: string;
  name: string;
  category: string;
  price_cents: number;
};

export type UpdateProductPayload = {
  name: string;
  category: string;
  price_cents: number;
  status: Product["status"];
  is_active: boolean;
};

export function listProducts() {
  return apiGet<Product[]>("/products");
}

export function createProduct(payload: CreateProductPayload) {
  return apiPost<Product>("/products", payload);
}

export function updateProduct(id: string, payload: UpdateProductPayload) {
  return apiPatch<Product>(`/products/${id}`, payload);
}

export function deleteProduct(id: string) {
  return apiDelete<{ deleted: boolean }>(`/products/${id}`);
}
