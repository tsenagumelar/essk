import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createRole, deleteRole, listPermissions, listRoles, updateRole } from "@/features/roles/api";

export function useRoles(tenantId?: string) {
  return useQuery({ queryKey: ["roles", tenantId ?? "all"], queryFn: () => listRoles(tenantId) });
}

export function usePermissions() {
  return useQuery({ queryKey: ["permissions"], queryFn: listPermissions });
}

export function useRoleMutations() {
  const queryClient = useQueryClient();
  const invalidate = async () => queryClient.invalidateQueries({ queryKey: ["roles"] });

  return {
    createRole: useMutation({ mutationFn: createRole, onSuccess: invalidate }),
    updateRole: useMutation({
      mutationFn: ({ id, payload }: { id: string; payload: Parameters<typeof updateRole>[1] }) =>
        updateRole(id, payload),
      onSuccess: invalidate,
    }),
    deleteRole: useMutation({ mutationFn: deleteRole, onSuccess: invalidate }),
  };
}
