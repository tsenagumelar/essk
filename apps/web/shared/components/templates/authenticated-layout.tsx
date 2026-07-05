"use client";

import { cn } from "@/lib/utils";

type AuthenticatedLayoutProps = {
  sidebar: React.ReactNode;
  navbar: React.ReactNode;
  children: React.ReactNode;
  sidebarCollapsed?: boolean;
};

export function AuthenticatedLayout({ sidebar, navbar, children, sidebarCollapsed = false }: AuthenticatedLayoutProps) {
  return (
    <div className="min-h-screen bg-slate-100">
      {sidebar}
      <div className={cn("transition-all duration-300", sidebarCollapsed ? "lg:pl-20" : "lg:pl-72")}>
        {navbar}
        <main className="px-4 py-6 lg:px-6">{children}</main>
      </div>
    </div>
  );
}
