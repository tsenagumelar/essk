"use client";

import { Spinner } from "@/shared/components/atoms/spinner";

export function LoadingState({ label = "Loading..." }: Readonly<{ label?: string }>) {
  return (
    <div className="flex items-center gap-2 p-4 text-sm text-muted-foreground">
      <Spinner />
      {label}
    </div>
  );
}
