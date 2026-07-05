"use client";

import { Pencil, Trash2 } from "lucide-react";
import { IconButton } from "@/shared/components/atoms/icon-button";

type RowActionsProps = {
  editLabel?: string;
  deleteLabel?: string;
  onEdit: () => void;
  onDelete: () => void;
};

export function RowActions({ editLabel = "Edit record", deleteLabel = "Delete record", onEdit, onDelete }: RowActionsProps) {
  return (
    <div className="flex justify-end gap-2">
      <IconButton size="sm" aria-label={editLabel} onClick={onEdit}>
        <Pencil className="h-4 w-4" />
      </IconButton>
      <IconButton size="sm" variant="danger" aria-label={deleteLabel} onClick={onDelete}>
        <Trash2 className="h-4 w-4" />
      </IconButton>
    </div>
  );
}
