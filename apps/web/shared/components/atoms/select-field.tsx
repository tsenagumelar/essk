"use client";

import { forwardRef } from "react";
import { cn } from "@/lib/utils";

export type SelectOption = {
  value: string;
  label: string;
};

type SelectFieldProps = Omit<React.SelectHTMLAttributes<HTMLSelectElement>, "children"> & {
  label?: string;
  options: SelectOption[];
};

export const SelectField = forwardRef<HTMLSelectElement, SelectFieldProps>(
  ({ className, label, id, options, ...props }, ref) => {
    const select = (
      <select
        ref={ref}
        id={id}
        className={cn("h-10 w-full rounded-lg border bg-white px-3 text-sm outline-none ring-primary focus:ring-2 disabled:bg-muted", className)}
        {...props}
      >
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    );

    if (!label) {
      return select;
    }

    return (
      <label className="block text-sm font-medium" htmlFor={id}>
        {label}
        <span className="mt-1 block">{select}</span>
      </label>
    );
  },
);

SelectField.displayName = "SelectField";
