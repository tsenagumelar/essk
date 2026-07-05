"use client";

type AppNavbarProps = {
  left?: React.ReactNode;
  search?: React.ReactNode;
  right?: React.ReactNode;
};

export function AppNavbar({ left, search, right }: AppNavbarProps) {
  return (
    <header className="sticky top-0 z-20 border-b border-slate-200 bg-white/95 backdrop-blur">
      <div className="flex h-16 items-center gap-4 px-4 lg:px-6">
        {left}
        {search}
        <div className="ml-auto flex items-center gap-3">{right}</div>
      </div>
    </header>
  );
}
