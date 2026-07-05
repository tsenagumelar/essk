"use client";

import { Filter } from "lucide-react";
import { SelectField, type SelectOption } from "@/shared/components/atoms/select-field";

type FilterSelectProps = {
  label: string;
  value: string;
  options: SelectOption[];
  withIcon?: boolean;
  onChange: (value: string) => void;
};

export function FilterSelect({ label, value, options, withIcon = false, onChange }: FilterSelectProps) {
  const select = (
    <SelectField
      aria-label={label}
      className={withIcon ? "h-auto w-full border-0 bg-transparent p-0 ring-0 focus:ring-0" : "h-10 w-44 shrink-0"}
      value={value}
      options={options}
      onChange={(event) => onChange(event.target.value)}
    />
  );

  if (!withIcon) {
    return select;
  }

  return (
    <div className="flex h-10 w-44 shrink-0 items-center gap-2 rounded-lg border bg-white px-3">
      <Filter className="h-4 w-4 text-muted-foreground" />
      {select}
    </div>
  );
}
