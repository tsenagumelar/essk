"use client";

import { cn } from "@/lib/utils";

export function Skeleton({ className }: Readonly<{ className?: string }>) {
  return <span className={cn("block animate-pulse rounded-md bg-slate-200", className)} />;
}
