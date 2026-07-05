"use client";

type PageShellProps = {
  title: string;
  subtitle?: string;
  toolbar?: React.ReactNode;
  children: React.ReactNode;
};

export function PageShell({ title, subtitle, toolbar, children }: PageShellProps) {
  return (
    <section className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
      <div className="border-b border-slate-200 bg-white px-5 py-4">
        <div className="space-y-4">
          <div>
            <h2 className="text-lg font-semibold tracking-tight text-slate-950">{title}</h2>
            {subtitle ? <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p> : null}
          </div>
          {toolbar}
        </div>
      </div>
      {children}
    </section>
  );
}
