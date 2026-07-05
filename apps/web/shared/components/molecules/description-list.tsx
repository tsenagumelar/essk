"use client";

type DescriptionItem = {
  label: string;
  value: React.ReactNode;
};

type DescriptionListProps = {
  items: DescriptionItem[];
};

export function DescriptionList({ items }: DescriptionListProps) {
  return (
    <dl className="grid gap-3 text-sm sm:grid-cols-2">
      {items.map((item) => (
        <div key={item.label}>
          <dt className="text-xs font-medium uppercase text-muted-foreground">{item.label}</dt>
          <dd className="mt-1 text-slate-900">{item.value}</dd>
        </div>
      ))}
    </dl>
  );
}
