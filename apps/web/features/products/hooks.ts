import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createProduct, deleteProduct, listProducts, updateProduct } from "@/features/products/api";

export function useProducts() {
  return useQuery({ queryKey: ["products"], queryFn: listProducts });
}

export function useProductMutations() {
  const queryClient = useQueryClient();
  const invalidate = async () => queryClient.invalidateQueries({ queryKey: ["products"] });

  return {
    createProduct: useMutation({ mutationFn: createProduct, onSuccess: invalidate }),
    updateProduct: useMutation({
      mutationFn: ({ id, payload }: { id: string; payload: Parameters<typeof updateProduct>[1] }) =>
        updateProduct(id, payload),
      onSuccess: invalidate,
    }),
    deleteProduct: useMutation({ mutationFn: deleteProduct, onSuccess: invalidate }),
  };
}
