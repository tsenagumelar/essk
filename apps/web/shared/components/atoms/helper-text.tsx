"use client";

import { cn } from "@/lib/utils";

type HelperTextProps = {
  children: React.ReactNode;
  variant?: "muted" | "error" | "success";
  className?: string;
};

const variants = {
  muted: "text-muted-foreground",
  error: "text-destructive",
  success: "text-emerald-700",
};

export function HelperText({ children, variant = "muted", className }: HelperTextProps) {
  return <p className={cn("text-xs", variants[variant], className)}>{children}</p>;
}
