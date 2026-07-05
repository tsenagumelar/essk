"use client";

import { cn } from "@/lib/utils";

type TabItem = {
  value: string;
  label: string;
};

type TabsProps = {
  value: string;
  items: TabItem[];
  onChange: (value: string) => void;
};

export function Tabs({ value, items, onChange }: TabsProps) {
  return (
    <div className="inline-flex rounded-lg border bg-white p-1">
      {items.map((item) => (
        <button
          key={item.value}
          type="button"
          className={cn(
            "rounded-md px-3 py-1.5 text-sm font-medium",
            value === item.value ? "bg-primary text-primary-foreground" : "text-slate-600 hover:bg-slate-50",
          )}
          onClick={() => onChange(item.value)}
        >
          {item.label}
        </button>
      ))}
    </div>
  );
}
