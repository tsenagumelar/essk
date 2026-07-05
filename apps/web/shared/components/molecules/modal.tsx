"use client";

import { X } from "lucide-react";
import { IconButton } from "@/shared/components/atoms/icon-button";

type ModalProps = {
  title: string;
  subtitle?: string;
  children: React.ReactNode;
  onClose: () => void;
};

export function Modal({ title, subtitle, children, onClose }: ModalProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4">
      <section className="w-full max-w-lg rounded-xl bg-white shadow-xl">
        <div className="flex items-center justify-between border-b px-5 py-4">
          <div>
            <h2 className="text-base font-semibold">{title}</h2>
            {subtitle ? <p className="text-sm text-muted-foreground">{subtitle}</p> : null}
          </div>
          <IconButton variant="ghost" size="sm" onClick={onClose} aria-label="Close form">
            <X className="h-4 w-4" />
          </IconButton>
        </div>
        {children}
      </section>
    </div>
  );
}
