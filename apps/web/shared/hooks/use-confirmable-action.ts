"use client";

import { useState } from "react";

export type PendingAction = {
  title: string;
  description: string;
  confirmLabel: string;
  variant?: "primary" | "danger";
  run: () => Promise<unknown>;
};

export function useConfirmableAction() {
  const [pending, setPending] = useState<PendingAction | null>(null);

  async function confirm() {
    if (!pending) {
      return;
    }
    await pending.run();
    setPending(null);
  }

  return {
    pending,
    request: setPending,
    cancel: () => setPending(null),
    confirm,
  };
}
