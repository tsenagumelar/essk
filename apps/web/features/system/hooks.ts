import { useQuery } from "@tanstack/react-query";
import { getHealth } from "@/features/system/api";

export function useHealth() {
  return useQuery({ queryKey: ["system", "health"], queryFn: getHealth });
}
