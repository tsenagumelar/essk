"use client";

import Link from "next/link";
import { cn } from "@/lib/utils";

type LinkButtonProps = React.ComponentProps<typeof Link> & {
  variant?: "primary" | "outline" | "ghost";
};

const variants = {
  primary: "bg-primary text-primary-foreground shadow-sm hover:opacity-95",
  outline: "border bg-white text-slate-700 hover:bg-slate-50",
  ghost: "text-slate-700 hover:bg-slate-100",
};

export function LinkButton({ className, variant = "primary", ...props }: LinkButtonProps) {
  return (
    <Link
      className={cn("inline-flex h-10 items-center justify-center gap-2 rounded-lg px-3 text-sm font-medium", variants[variant], className)}
      {...props}
    />
  );
}
