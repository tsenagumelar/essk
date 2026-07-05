"use client";

import { cn } from "@/lib/utils";

type BadgeVariant = "default" | "success" | "muted" | "danger";

type BadgeProps = {
  children: React.ReactNode;
  variant?: BadgeVariant;
  className?: string;
};

const variants: Record<BadgeVariant, string> = {
  default: "bg-slate-100 text-slate-700",
  success: "bg-emerald-50 text-emerald-700",
  muted: "bg-slate-50 text-muted-foreground",
  danger: "bg-red-50 text-destructive",
};

export function Badge({ children, variant = "default", className }: BadgeProps) {
  return <span className={cn("rounded-full px-2 py-1 text-xs font-medium", variants[variant], className)}>{children}</span>;
}
