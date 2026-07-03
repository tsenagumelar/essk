"use client";

import { useState } from "react";
import { ConfirmationDialog } from "@/shared/molecules/confirmation-dialog";

export type PendingAction = {
  title: string;
  description: string;
  confirmLabel: string;
  variant?: "primary" | "danger";
  run: () => Promise<unknown>;
};

export function useConfirmableAction() {
  const [pending, setPending] = useState<PendingAction | null>(null);
  return {
    pending,
    request: setPending,
    dialog: (isLoading = false) => (
      <ConfirmationDialog
        open={Boolean(pending)}
        title={pending?.title ?? ""}
        description={pending?.description ?? ""}
        confirmLabel={pending?.confirmLabel ?? "Confirm"}
        variant={pending?.variant ?? "primary"}
        isLoading={isLoading}
        onCancel={() => setPending(null)}
        onConfirm={async () => {
          if (!pending) {
            return;
          }
          await pending.run();
          setPending(null);
        }}
      />
    ),
  };
}
