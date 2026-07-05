"use client";

import { AlertCircle, CheckCircle2, Info, TriangleAlert } from "lucide-react";
import { cn } from "@/lib/utils";

type AlertVariant = "info" | "success" | "warning" | "danger";

type AlertProps = {
  title?: string;
  children: React.ReactNode;
  variant?: AlertVariant;
  className?: string;
};

const variants = {
  info: "border-blue-100 bg-blue-50 text-blue-900",
  success: "border-emerald-100 bg-emerald-50 text-emerald-900",
  warning: "border-amber-100 bg-amber-50 text-amber-900",
  danger: "border-red-100 bg-red-50 text-red-900",
};

const icons = {
  info: Info,
  success: CheckCircle2,
  warning: TriangleAlert,
  danger: AlertCircle,
};

export function Alert({ title, children, variant = "info", className }: AlertProps) {
  const Icon = icons[variant];
  return (
    <div className={cn("flex gap-3 rounded-lg border p-3 text-sm", variants[variant], className)}>
      <Icon className="mt-0.5 h-4 w-4 shrink-0" />
      <div>
        {title ? <p className="font-semibold">{title}</p> : null}
        <div className={title ? "mt-1" : undefined}>{children}</div>
      </div>
    </div>
  );
}
