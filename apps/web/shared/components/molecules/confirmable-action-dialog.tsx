"use client";

import { ConfirmationDialog } from "@/shared/components/molecules/confirmation-dialog";
import type { PendingAction } from "@/shared/hooks/use-confirmable-action";

type ConfirmableActionDialogProps = {
  pending: PendingAction | null;
  isLoading?: boolean;
  onCancel: () => void;
  onConfirm: () => void | Promise<void>;
};

export function ConfirmableActionDialog({
  pending,
  isLoading = false,
  onCancel,
  onConfirm,
}: ConfirmableActionDialogProps) {
  return (
    <ConfirmationDialog
      open={Boolean(pending)}
      title={pending?.title ?? ""}
      description={pending?.description ?? ""}
      confirmLabel={pending?.confirmLabel ?? "Confirm"}
      variant={pending?.variant ?? "primary"}
      isLoading={isLoading}
      onCancel={onCancel}
      onConfirm={onConfirm}
    />
  );
}
