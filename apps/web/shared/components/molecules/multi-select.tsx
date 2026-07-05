"use client";

import { Check, ChevronDown } from "lucide-react";
import { useState } from "react";
import type { SelectOption } from "@/shared/components/atoms/select-field";

type MultiSelectProps = {
  label: string;
  values: string[];
  options: SelectOption[];
  emptyLabel?: string;
  onChange: (values: string[]) => void;
};

export function MultiSelect({ label, values, options, emptyLabel = "Select options", onChange }: MultiSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const selectedLabels = options.filter((option) => values.includes(option.value)).map((option) => option.label);

  function toggle(value: string) {
    onChange(values.includes(value) ? values.filter((current) => current !== value) : [...values, value]);
  }

  return (
    <div className="relative text-sm font-medium">
      <p>{label}</p>
      <button
        type="button"
        className="mt-1 flex min-h-10 w-full items-center justify-between gap-3 rounded-lg border bg-white px-3 py-2 text-left text-sm outline-none ring-primary focus:ring-2"
        onClick={() => setIsOpen((current) => !current)}
      >
        <span className={selectedLabels.length === 0 ? "text-muted-foreground" : "text-slate-900"}>
          {selectedLabels.length === 0 ? emptyLabel : selectedLabels.join(", ")}
        </span>
        <ChevronDown className="h-4 w-4 shrink-0 text-muted-foreground" />
      </button>
      {isOpen ? (
        <div className="absolute z-20 mt-2 max-h-56 w-full overflow-auto rounded-lg border bg-white p-2 shadow-lg">
          {options.length === 0 ? (
            <p className="px-2 py-2 text-sm text-muted-foreground">No options available.</p>
          ) : (
            options.map((option) => {
              const selected = values.includes(option.value);
              return (
                <button
                  key={option.value}
                  type="button"
                  className="flex w-full items-center gap-2 rounded-md px-2 py-2 text-left text-sm hover:bg-slate-50"
                  onClick={() => toggle(option.value)}
                >
                  <span className="flex h-4 w-4 items-center justify-center rounded border">
                    {selected ? <Check className="h-3 w-3" /> : null}
                  </span>
                  {option.label}
                </button>
              );
            })
          )}
        </div>
      ) : null}
    </div>
  );
}
