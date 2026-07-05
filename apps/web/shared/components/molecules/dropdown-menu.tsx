"use client";

import { cn } from "@/lib/utils";

type DropdownMenuProps = {
  open: boolean;
  align?: "left" | "right";
  children: React.ReactNode;
};

export function DropdownMenu({ open, align = "right", children }: DropdownMenuProps) {
  if (!open) {
    return null;
  }

  return (
    <div
      className={cn(
        "absolute z-30 mt-2 min-w-56 rounded-xl border bg-white p-2 shadow-lg",
        align === "right" ? "right-0" : "left-0",
      )}
    >
      {children}
    </div>
  );
}

type DropdownMenuItemProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  icon?: React.ReactNode;
};

export function DropdownMenuItem({ className, icon, children, type = "button", ...props }: DropdownMenuItemProps) {
  return (
    <button
      type={type}
      className={cn("flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm hover:bg-muted", className)}
      {...props}
    >
      {icon}
      {children}
    </button>
  );
}
