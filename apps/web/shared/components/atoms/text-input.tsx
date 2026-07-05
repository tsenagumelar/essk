"use client";

import { forwardRef } from "react";
import { cn } from "@/lib/utils";

type TextInputProps = React.InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
  error?: string;
};

export const TextInput = forwardRef<HTMLInputElement, TextInputProps>(
  ({ className, label, error, id, ...props }, ref) => {
    const input = (
      <input
        ref={ref}
        id={id}
        className={cn(
          "h-10 w-full rounded-lg border px-3 text-sm outline-none ring-primary focus:ring-2 disabled:bg-muted",
          error && "border-destructive",
          className,
        )}
        {...props}
      />
    );

    if (!label) {
      return input;
    }

    return (
      <label className="block text-sm font-medium" htmlFor={id}>
        {label}
        <span className="mt-1 block">{input}</span>
        {error ? <span className="mt-1 block text-xs text-destructive">{error}</span> : null}
      </label>
    );
  },
);

TextInput.displayName = "TextInput";
