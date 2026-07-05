"use client";

import { Search } from "lucide-react";
import { cn } from "@/lib/utils";

type SearchBoxProps = {
  value: string;
  placeholder?: string;
  className?: string;
  onChange: (value: string) => void;
};

export function SearchBox({ value, placeholder = "Search records", className, onChange }: SearchBoxProps) {
  return (
    <div className={cn("relative min-w-0", className)}>
      <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
      <input
        type="search"
        className="h-10 w-full rounded-lg border bg-white pl-9 pr-3 text-sm outline-none ring-primary focus:ring-2"
        placeholder={placeholder}
        value={value}
        onChange={(event) => onChange(event.target.value)}
      />
    </div>
  );
}
