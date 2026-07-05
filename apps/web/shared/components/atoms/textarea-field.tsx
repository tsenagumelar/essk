"use client";

import { forwardRef } from "react";
import { cn } from "@/lib/utils";
import { Label } from "@/shared/components/atoms/label";

type TextareaFieldProps = React.TextareaHTMLAttributes<HTMLTextAreaElement> & {
  label?: string;
  error?: string;
};

export const TextareaField = forwardRef<HTMLTextAreaElement, TextareaFieldProps>(
  ({ className, label, error, id, ...props }, ref) => {
    const textarea = (
      <textarea
        ref={ref}
        id={id}
        className={cn(
          "min-h-24 w-full rounded-lg border px-3 py-2 text-sm outline-none ring-primary focus:ring-2 disabled:bg-muted",
          error && "border-destructive",
          className,
        )}
        {...props}
      />
    );

    if (!label) {
      return textarea;
    }

    return (
      <div className="space-y-1">
        <Label htmlFor={id}>{label}</Label>
        {textarea}
        {error ? <p className="text-xs text-destructive">{error}</p> : null}
      </div>
    );
  },
);

TextareaField.displayName = "TextareaField";
