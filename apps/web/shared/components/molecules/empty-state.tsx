"use client";

import { Button } from "@/shared/components/atoms/button";

type EmptyStateProps = {
  title: string;
  description?: string;
  actionLabel?: string;
  icon?: React.ReactNode;
  onAction?: () => void;
};

export function EmptyState({ title, description, actionLabel, icon, onAction }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center p-10 text-center">
      {icon ? <div className="mb-3 text-muted-foreground">{icon}</div> : null}
      <p className="text-sm font-medium">{title}</p>
      {description ? <p className="mt-1 max-w-md text-sm text-muted-foreground">{description}</p> : null}
      {actionLabel && onAction ? (
        <Button type="button" className="mt-4" onClick={onAction}>
          {actionLabel}
        </Button>
      ) : null}
    </div>
  );
}
