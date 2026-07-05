"use client";

import { DescriptionList } from "@/shared/components/molecules/description-list";

type DetailPanelProps = {
  title: string;
  description?: string;
  items: { label: string; value: React.ReactNode }[];
  actions?: React.ReactNode;
};

export function DetailPanel({ title, description, items, actions }: DetailPanelProps) {
  return (
    <section className="rounded-xl border bg-white shadow-sm">
      <div className="flex flex-col gap-3 border-b px-5 py-4 md:flex-row md:items-center md:justify-between">
        <div>
          <h2 className="text-base font-semibold">{title}</h2>
          {description ? <p className="mt-1 text-sm text-muted-foreground">{description}</p> : null}
        </div>
        {actions}
      </div>
      <div className="p-5">
        <DescriptionList items={items} />
      </div>
    </section>
  );
}
