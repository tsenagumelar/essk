"use client";

import { Bell, LayoutDashboard, LogOut, Package, Search, Settings, ShieldCheck, UserCircle } from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";
import { clearSession, getAccessToken, getStoredUser } from "@/features/auth/session";
import { cn } from "@/lib/utils";

const navItems = [
  { href: "/dashboard", label: "Dashboard", icon: LayoutDashboard },
  { href: "/products", label: "Products", icon: Package },
  { href: "/health", label: "Health", icon: ShieldCheck },
  { href: "/profile", label: "Profile", icon: UserCircle },
];

export function AppShell({ children }: Readonly<{ children: React.ReactNode }>) {
  const router = useRouter();
  const pathname = usePathname();
  const [isReady, setIsReady] = useState(false);
  const [isProfileOpen, setIsProfileOpen] = useState(false);
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
    clearSession();
    router.replace("/login");
  }

  if (!isReady) {
    return <main className="flex min-h-screen items-center justify-center bg-muted text-sm text-muted-foreground">Loading workspace...</main>;
  }

  return (
    <div className="min-h-screen bg-muted">
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-72 border-r bg-white lg:flex lg:flex-col">
        <div className="flex h-16 items-center gap-3 border-b px-5">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary text-sm font-bold text-primary-foreground">
            ES
          </div>
          <div>
            <p className="text-sm font-semibold">ESSK</p>
            <p className="text-xs text-muted-foreground">Enterprise SaaS Kit</p>
          </div>
        </div>

        <nav className="flex-1 space-y-1 px-3 py-4">
          {navItems.map((item) => {
            const Icon = item.icon;
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-muted-foreground hover:bg-muted hover:text-foreground",
                  isActive && "bg-muted text-foreground",
                )}
              >
                <Icon className="h-4 w-4" />
                {item.label}
              </Link>
            );
          })}
        </nav>

        <div className="border-t px-5 py-4 text-xs text-muted-foreground">
          <p>Copyright © 2026 ESSK.</p>
          <p>All rights reserved.</p>
        </div>
      </aside>

      <div className="lg:pl-72">
        <header className="sticky top-0 z-20 border-b bg-white">
          <div className="flex h-16 items-center gap-4 px-4 lg:px-6">
            <Link href="/dashboard" className="flex items-center gap-2 font-semibold lg:hidden">
              <span className="flex h-8 w-8 items-center justify-center rounded-md bg-primary text-xs text-primary-foreground">ES</span>
              ESSK
            </Link>

            <div className="relative max-w-xl flex-1">
              <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input
                type="search"
                placeholder="Search modules, records, or actions"
                className="h-10 w-full rounded-md border bg-muted/50 pl-9 pr-3 text-sm outline-none ring-primary focus:bg-white focus:ring-2"
              />
            </div>

            <button className="relative flex h-10 w-10 items-center justify-center rounded-md border bg-white hover:bg-muted" aria-label="Notifications">
              <Bell className="h-4 w-4" />
              <span className="absolute right-2 top-2 h-2 w-2 rounded-full bg-destructive" />
            </button>

            <div className="relative">
              <button
                className="flex h-10 items-center gap-2 rounded-md border bg-white px-2 hover:bg-muted"
                onClick={() => setIsProfileOpen((current) => !current)}
              >
                <span className="flex h-7 w-7 items-center justify-center rounded-full bg-primary text-xs font-semibold text-primary-foreground">
                  {initials}
                </span>
                <span className="hidden max-w-32 truncate text-sm font-medium md:block">{user?.name ?? "User"}</span>
              </button>

              {isProfileOpen ? (
                <div className="absolute right-0 mt-2 w-72 rounded-lg border bg-white p-2 shadow-lg">
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
                  <button className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm text-destructive hover:bg-muted" onClick={logout}>
                    <LogOut className="h-4 w-4" />
                    Logout
                  </button>
                </div>
              ) : null}
            </div>
          </div>
        </header>

        <main className="px-4 py-6 lg:px-6">{children}</main>
      </div>
    </div>
  );
}
