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
    <div className={cn("flex h-10 min-w-0 items-center gap-2 rounded-lg border bg-white px-3 ring-primary focus-within:ring-2", className)}>
      <Search className="h-4 w-4 shrink-0 text-muted-foreground" />
      <input
        type="search"
        className="h-full min-w-0 flex-1 border-0 bg-transparent p-0 text-sm leading-none outline-none placeholder:text-muted-foreground"
        placeholder={placeholder}
        value={value}
        onChange={(event) => onChange(event.target.value)}
      />
    </div>
  );
}
