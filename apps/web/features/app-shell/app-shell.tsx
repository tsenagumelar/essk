"use client";

import {
  Bell,
  ChevronDown,
  LayoutDashboard,
  LogOut,
  Menu,
  Package,
  Shield,
  Search,
  Settings,
  ShieldCheck,
  Building2,
  Users,
  UserCircle,
} from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { clearSession, getAccessToken, getStoredUser } from "@/features/auth/session";
import { ConfirmationDialog } from "@/features/shared/components/confirmation-dialog";
import { cn } from "@/lib/utils";

const navGroups = [
  {
    label: "Workspace",
    items: [{ href: "/dashboard", label: "Dashboard", icon: LayoutDashboard }],
  },
  {
    label: "Master Data",
    items: [
      { href: "/tenants", label: "Tenants", icon: Building2 },
      { href: "/users", label: "Users", icon: Users },
      { href: "/roles", label: "Roles", icon: Shield },
      { href: "/products", label: "Products", icon: Package },
    ],
  },
  {
    label: "Settings",
    items: [{ href: "/health", label: "Health", icon: ShieldCheck }],
  },
];

export function AppShell({ children }: Readonly<{ children: React.ReactNode }>) {
  const router = useRouter();
  const pathname = usePathname();
  const [isReady, setIsReady] = useState(false);
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [isProfileOpen, setIsProfileOpen] = useState(false);
  const [isLogoutConfirmOpen, setIsLogoutConfirmOpen] = useState(false);
  const [user, setUser] = useState<ReturnType<typeof getStoredUser>>(null);

  useEffect(() => {
    const token = getAccessToken();
    if (!token) {
      router.replace("/login");
      return;
    }
    setUser(getStoredUser());
    setIsReady(true);
  }, [router]);

  const initials = useMemo(() => {
    const source = user?.name || user?.email || "User";
    return source
      .split(" ")
      .map((part) => part[0])
      .join("")
      .slice(0, 2)
      .toUpperCase();
  }, [user]);

  function logout() {
    setIsLogoutConfirmOpen(false);
    clearSession();
    router.replace("/login");
  }

  if (!isReady) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-100 text-sm text-muted-foreground">
        Loading workspace...
      </main>
    );
  }

  return (
    <div className="min-h-screen bg-slate-100">
      <aside
        className={cn(
          "fixed inset-y-0 left-0 z-30 hidden border-r border-slate-200 bg-white shadow-sm transition-all duration-300 lg:flex lg:flex-col",
          isSidebarCollapsed ? "w-20" : "w-72",
        )}
      >
        <div className={cn("flex h-16 items-center border-b px-4", isSidebarCollapsed ? "justify-center" : "gap-3")}>
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-primary to-sky-500 text-sm font-bold text-primary-foreground shadow-sm">
            ES
          </div>
          <div className={cn("min-w-0 transition-opacity", isSidebarCollapsed && "hidden")}>
            <p className="text-sm font-semibold">ESSK</p>
            <p className="text-xs text-muted-foreground">Enterprise SaaS Kit</p>
          </div>
        </div>

        <nav className="flex-1 space-y-5 px-3 py-5">
          {navGroups.map((group) => (
            <div key={group.label}>
              <p
                className={cn(
                  "mb-2 px-3 text-[11px] font-semibold uppercase tracking-wide text-slate-400",
                  isSidebarCollapsed && "sr-only",
                )}
              >
                {group.label}
              </p>
              <div className="space-y-1">
                {group.items.map((item) => {
                  const Icon = item.icon;
                  const isActive = pathname === item.href;
                  return (
                    <Link
                      key={item.href}
                      href={item.href}
                      title={isSidebarCollapsed ? `${group.label}: ${item.label}` : undefined}
                      className={cn(
                        "group relative flex h-11 items-center rounded-xl text-sm font-medium transition-colors",
                        isSidebarCollapsed ? "justify-center px-0" : "gap-3 px-3",
                        isActive
                          ? "bg-primary text-primary-foreground shadow-sm"
                          : "text-slate-600 hover:bg-slate-100 hover:text-slate-950",
                      )}
                    >
                      <Icon className="h-4 w-4 shrink-0" />
                      <span className={cn(isSidebarCollapsed && "hidden")}>{item.label}</span>
                      {isActive ? (
                        <span
                          className={cn(
                            "absolute rounded-full bg-primary-foreground/70",
                            isSidebarCollapsed ? "bottom-1 h-1 w-5" : "right-3 h-1.5 w-1.5",
                          )}
                        />
                      ) : null}
                    </Link>
                  );
                })}
              </div>
            </div>
          ))}
        </nav>

        <div className={cn("border-t px-5 py-4 text-xs text-muted-foreground", isSidebarCollapsed && "px-2 text-center")}>
          {isSidebarCollapsed ? (
            <p>© 2026</p>
          ) : (
            <>
              <p>Copyright © 2026 ESSK.</p>
              <p>All rights reserved.</p>
            </>
          )}
        </div>
      </aside>

      <div className={cn("transition-all duration-300", isSidebarCollapsed ? "lg:pl-20" : "lg:pl-72")}>
        <header className="sticky top-0 z-20 border-b border-slate-200 bg-white/95 backdrop-blur">
          <div className="flex h-16 items-center gap-4 px-4 lg:px-6">
            <Link href="/dashboard" className="flex items-center gap-2 font-semibold lg:hidden">
              <span className="flex h-8 w-8 items-center justify-center rounded-md bg-primary text-xs text-primary-foreground">
                ES
              </span>
              ESSK
            </Link>

            <button
              className="hidden h-10 w-10 items-center justify-center rounded-xl border bg-white text-slate-600 shadow-sm hover:bg-slate-50 hover:text-slate-950 lg:flex"
              aria-label={isSidebarCollapsed ? "Expand sidebar" : "Collapse sidebar"}
              onClick={() => setIsSidebarCollapsed((current) => !current)}
            >
              <Menu className="h-4 w-4" />
            </button>

            <div className="relative max-w-xl flex-1">
              <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type="search"
                placeholder="Search modules, records, or actions"
                className="h-10 w-full rounded-xl border bg-slate-50 pl-9 pr-3 text-sm outline-none ring-primary focus:bg-white focus:ring-2"
              />
            </div>

            <div className="ml-auto flex items-center gap-3">
              <button
                className="relative flex h-10 w-10 items-center justify-center rounded-xl border bg-white shadow-sm hover:bg-slate-50"
                aria-label="Notifications"
              >
                <Bell className="h-4 w-4" />
                <span className="absolute right-2 top-2 h-2 w-2 rounded-full bg-destructive" />
              </button>

              <div className="relative">
                <button
                  className="flex h-10 items-center gap-2 rounded-xl border bg-white px-2 shadow-sm hover:bg-slate-50"
                  onClick={() => setIsProfileOpen((current) => !current)}
                >
                  <span className="flex h-7 w-7 items-center justify-center rounded-full bg-primary text-xs font-semibold text-primary-foreground">
                    {initials}
                  </span>
                  <span className="hidden max-w-32 truncate text-sm font-medium md:block">{user?.name ?? "User"}</span>
                  <ChevronDown className="hidden h-4 w-4 text-muted-foreground md:block" />
                </button>

                {isProfileOpen ? (
                  <div className="absolute right-0 mt-2 w-72 rounded-xl border bg-white p-2 shadow-lg">
                    <div className="border-b px-3 py-3">
                      <p className="font-medium">{user?.name ?? "User"}</p>
                      <p className="truncate text-sm text-muted-foreground">{user?.email ?? "-"}</p>
                      <p className="mt-1 truncate text-xs text-muted-foreground">Tenant: {user?.tenant_id ?? "-"}</p>
                    </div>
                    <Link className="flex items-center gap-2 rounded-md px-3 py-2 text-sm hover:bg-muted" href="/profile">
                      <UserCircle className="h-4 w-4" />
                      Profile detail
                    </Link>
                    <button className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm hover:bg-muted">
                      <Settings className="h-4 w-4" />
                      Settings
                    </button>
                    <button
                      className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm text-destructive hover:bg-muted"
                      onClick={() => {
                        setIsProfileOpen(false);
                        setIsLogoutConfirmOpen(true);
                      }}
                    >
                      <LogOut className="h-4 w-4" />
                      Logout
                    </button>
                  </div>
                ) : null}
              </div>
            </div>
          </div>
        </header>

        <main className="px-4 py-6 lg:px-6">{children}</main>
      </div>

      <ConfirmationDialog
        open={isLogoutConfirmOpen}
        title="Logout from workspace?"
        description="Your local session token will be removed and you will be redirected to the login page."
        confirmLabel="Logout"
        variant="danger"
        onCancel={() => setIsLogoutConfirmOpen(false)}
        onConfirm={logout}
      />
    </div>
  );
}
