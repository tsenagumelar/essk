"use client";

type StatItem = {
  label: string;
  value: React.ReactNode;
  description?: string;
  icon?: React.ReactNode;
};

type StatGridProps = {
  items: StatItem[];
};

export function StatGrid({ items }: StatGridProps) {
  return (
    <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
      {items.map((item) => (
        <div key={item.label} className="rounded-xl border bg-white p-4 shadow-sm">
          <div className="flex items-start justify-between gap-3">
            <div>
              <p className="text-sm text-muted-foreground">{item.label}</p>
              <p className="mt-1 text-2xl font-semibold">{item.value}</p>
            </div>
            {item.icon ? <div className="rounded-lg bg-slate-50 p-2 text-primary">{item.icon}</div> : null}
          </div>
          {item.description ? <p className="mt-3 text-xs text-muted-foreground">{item.description}</p> : null}
        </div>
      ))}
    </div>
  );
}
