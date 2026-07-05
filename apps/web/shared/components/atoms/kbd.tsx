"use client";

import { cn } from "@/lib/utils";

export function Kbd({ children, className }: Readonly<{ children: React.ReactNode; className?: string }>) {
  return (
    <kbd className={cn("rounded border bg-slate-50 px-1.5 py-0.5 font-mono text-[11px] text-slate-600", className)}>
      {children}
    </kbd>
  );
}
