"use client";

import { cn } from "@/lib/utils";

type SeparatorProps = {
  orientation?: "horizontal" | "vertical";
  className?: string;
};

export function Separator({ orientation = "horizontal", className }: SeparatorProps) {
  return (
    <span
      aria-hidden="true"
      className={cn(
        "block bg-border",
        orientation === "horizontal" ? "h-px w-full" : "h-full w-px",
        className,
      )}
    />
  );
}
