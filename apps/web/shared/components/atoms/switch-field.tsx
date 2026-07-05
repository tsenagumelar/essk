"use client";

import { cn } from "@/lib/utils";

type SwitchFieldProps = {
  checked: boolean;
  label?: string;
  disabled?: boolean;
  onCheckedChange: (checked: boolean) => void;
};

export function SwitchField({ checked, label, disabled = false, onCheckedChange }: SwitchFieldProps) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      disabled={disabled}
      className="inline-flex items-center gap-2 text-sm disabled:opacity-60"
      onClick={() => onCheckedChange(!checked)}
    >
      <span
        className={cn(
          "relative inline-flex h-6 w-10 rounded-full transition-colors",
          checked ? "bg-primary" : "bg-slate-300",
        )}
      >
        <span
          className={cn(
            "absolute top-1 h-4 w-4 rounded-full bg-white transition-transform",
            checked ? "translate-x-5" : "translate-x-1",
          )}
        />
      </span>
      {label}
    </button>
  );
}
