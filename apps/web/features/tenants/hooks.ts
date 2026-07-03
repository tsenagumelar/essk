import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createTenant, deleteTenant, listTenants, updateTenant } from "@/features/tenants/api";

export function useTenants() {
  return useQuery({ queryKey: ["tenants"], queryFn: listTenants });
}

export function useTenantMutations() {
  const queryClient = useQueryClient();
  const invalidate = async () => queryClient.invalidateQueries({ queryKey: ["tenants"] });

  return {
    createTenant: useMutation({ mutationFn: createTenant, onSuccess: invalidate }),
    updateTenant: useMutation({
      mutationFn: ({ id, payload }: { id: string; payload: Parameters<typeof updateTenant>[1] }) =>
        updateTenant(id, payload),
      onSuccess: invalidate,
    }),
    deleteTenant: useMutation({ mutationFn: deleteTenant, onSuccess: invalidate }),
  };
}
