import { useQuery } from "@tanstack/react-query";
import { getMe } from "@/features/auth/api";
import { getStoredUser } from "@/shared/auth/session";

export function useProfile() {
  const fallbackUser = getStoredUser();
  const meQuery = useQuery({
    queryKey: ["auth", "me"],
    queryFn: getMe,
  });

  return {
    meQuery,
    user: meQuery.data ?? fallbackUser,
  };
}
