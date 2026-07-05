"use client";

import { cn } from "@/lib/utils";

type CheckboxFieldProps = Omit<React.InputHTMLAttributes<HTMLInputElement>, "type"> & {
  label: string;
};

export function CheckboxField({ className, label, ...props }: CheckboxFieldProps) {
  return (
    <label className={cn("flex items-center gap-2 text-sm", className)}>
      <input type="checkbox" {...props} />
      {label}
    </label>
  );
}
