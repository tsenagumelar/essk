import { env } from "@/lib/env";

export type ApiEnvelope<T> = {
  success: boolean;
  message: string;
  data?: T;
  meta?: unknown;
  errors?: unknown;
};

export class ApiError extends Error {
  status: number;
  errors: unknown;

  constructor(message: string, status: number, errors: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.errors = errors;
  }
}

export async function apiGet<T>(path: string): Promise<T> {
  const response = await fetch(`${env.apiBaseUrl}${path}`, {
    headers: {
      Accept: "application/json",
    },
    cache: "no-store",
  });

  const envelope = (await response.json()) as ApiEnvelope<T>;

  if (!response.ok || !envelope.success) {
    throw new ApiError(envelope.message || "API request failed", response.status, envelope.errors);
  }

  return envelope.data as T;
}
