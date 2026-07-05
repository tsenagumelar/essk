"use client";

import Link from "next/link";
import { cn } from "@/lib/utils";

type SidebarItem = {
  href: string;
  label: string;
  icon?: React.ReactNode;
};

type SidebarGroup = {
  label: string;
  items: SidebarItem[];
};

type AppSidebarProps = {
  groups: SidebarGroup[];
  activePath: string;
  collapsed?: boolean;
  position?: "fixed" | "static";
  brand?: React.ReactNode;
  footer?: React.ReactNode;
};

export function AppSidebar({ groups, activePath, collapsed = false, position = "fixed", brand, footer }: AppSidebarProps) {
  return (
    <aside
      className={cn(
        "border-r border-slate-200 bg-white shadow-sm transition-all duration-300 lg:flex lg:flex-col",
        position === "fixed" ? "fixed inset-y-0 left-0 z-30 hidden" : "relative flex min-h-96",
        collapsed ? "w-20" : "w-72",
      )}
    >
      {brand}
      <nav className="flex-1 space-y-5 px-3 py-5">
        {groups.map((group) => (
          <div key={group.label}>
            <p className={cn("mb-2 px-3 text-[11px] font-semibold uppercase tracking-wide text-slate-400", collapsed && "sr-only")}>
              {group.label}
            </p>
            <div className="space-y-1">
              {group.items.map((item) => {
                const isActive = activePath === item.href;
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    title={collapsed ? `${group.label}: ${item.label}` : undefined}
                    className={cn(
                      "group relative flex h-11 items-center rounded-xl text-sm font-medium transition-colors",
                      collapsed ? "justify-center px-0" : "gap-3 px-3",
                      isActive ? "bg-primary text-primary-foreground shadow-sm" : "text-slate-600 hover:bg-slate-100 hover:text-slate-950",
                    )}
                  >
                    {item.icon}
                    <span className={cn(collapsed && "hidden")}>{item.label}</span>
                  </Link>
                );
              })}
            </div>
          </div>
        ))}
      </nav>
      {footer}
    </aside>
  );
}
