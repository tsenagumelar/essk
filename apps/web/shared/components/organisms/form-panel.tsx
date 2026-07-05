"use client";

import { Button } from "@/shared/components/atoms/button";

type FormPanelProps = {
  title?: string;
  description?: string;
  children: React.ReactNode;
  submitLabel?: string;
  cancelLabel?: string;
  isSubmitting?: boolean;
  onCancel?: () => void;
};

export function FormPanel({
  title,
  description,
  children,
  submitLabel = "Save",
  cancelLabel = "Cancel",
  isSubmitting = false,
  onCancel,
}: FormPanelProps) {
  return (
    <div className="rounded-xl border bg-white shadow-sm">
      {title || description ? (
        <div className="border-b px-5 py-4">
          {title ? <h2 className="text-base font-semibold">{title}</h2> : null}
          {description ? <p className="mt-1 text-sm text-muted-foreground">{description}</p> : null}
        </div>
      ) : null}
      <div className="space-y-4 p-5">{children}</div>
      <div className="flex justify-end gap-2 border-t px-5 py-4">
        {onCancel ? (
          <Button type="button" variant="outline" onClick={onCancel}>
            {cancelLabel}
          </Button>
        ) : null}
        <Button type="submit" isLoading={isSubmitting} loadingLabel="Saving...">
          {submitLabel}
        </Button>
      </div>
    </div>
  );
}
