import { ProductWorkspace } from "@/features/products/components/product-workspace";

export default function ProductsPage() {
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-semibold">Products</h1>
        <p className="mt-2 text-sm text-muted-foreground">Sample modular master data CRUD connected to the backend.</p>
      </div>
      <ProductWorkspace />
    </div>
  );
}
