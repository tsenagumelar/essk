"use client";

type PageShellProps = {
  title: string;
  subtitle?: string;
  toolbar?: React.ReactNode;
  children: React.ReactNode;
};

export function PageShell({ title, subtitle, toolbar, children }: PageShellProps) {
  return (
    <section className="overflow-hidden rounded-xl border bg-white shadow-sm">
      <div className="border-b bg-gradient-to-r from-white to-slate-50 px-4 py-4">
        <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
          <div>
            <h2 className="text-base font-semibold">{title}</h2>
            {subtitle ? <p className="mt-1 text-sm text-muted-foreground">{subtitle}</p> : null}
          </div>
          {toolbar}
        </div>
      </div>
      {children}
    </section>
  );
}
