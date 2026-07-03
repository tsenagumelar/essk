import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createUser, deleteUser, listUsers, updateUser } from "@/features/users/api";

export function useUsers() {
  return useQuery({ queryKey: ["users"], queryFn: listUsers });
}

export function useUserMutations() {
  const queryClient = useQueryClient();
  const invalidate = async () => queryClient.invalidateQueries({ queryKey: ["users"] });

  return {
    createUser: useMutation({ mutationFn: createUser, onSuccess: invalidate }),
    updateUser: useMutation({
      mutationFn: ({ id, payload }: { id: string; payload: Parameters<typeof updateUser>[1] }) =>
        updateUser(id, payload),
      onSuccess: invalidate,
    }),
    deleteUser: useMutation({ mutationFn: deleteUser, onSuccess: invalidate }),
  };
}
