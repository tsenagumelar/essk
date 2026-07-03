import { AppShell } from "@/features/app-shell/app-shell";

export default function AppLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <AppShell>{children}</AppShell>;
}
