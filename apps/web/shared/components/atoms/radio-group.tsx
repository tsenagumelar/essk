"use client";

import { cn } from "@/lib/utils";

type RadioOption = {
  value: string;
  label: string;
};

type RadioGroupProps = {
  name: string;
  value: string;
  options: RadioOption[];
  direction?: "row" | "column";
  onChange: (value: string) => void;
};

export function RadioGroup({ name, value, options, direction = "row", onChange }: RadioGroupProps) {
  return (
    <div className={cn("flex gap-3", direction === "column" && "flex-col")}>
      {options.map((option) => (
        <label key={option.value} className="inline-flex items-center gap-2 text-sm">
          <input
            type="radio"
            name={name}
            value={option.value}
            checked={value === option.value}
            onChange={(event) => onChange(event.target.value)}
          />
          {option.label}
        </label>
      ))}
    </div>
  );
}
