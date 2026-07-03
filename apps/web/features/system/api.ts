import { apiGet } from "@/shared/api/client";

export type HealthResponse = {
  app: string;
  env: string;
  version: string;
  status: string;
};

export function getHealth() {
  return apiGet<HealthResponse>("/health");
}
