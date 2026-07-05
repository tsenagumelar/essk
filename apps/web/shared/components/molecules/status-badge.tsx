"use client";

import { Badge } from "@/shared/components/atoms/badge";

type StatusBadgeProps = {
  active: boolean;
  activeLabel?: string;
  inactiveLabel?: string;
};

export function StatusBadge({ active, activeLabel = "Active", inactiveLabel = "Inactive" }: StatusBadgeProps) {
  return <Badge variant={active ? "success" : "muted"}>{active ? activeLabel : inactiveLabel}</Badge>;
}
