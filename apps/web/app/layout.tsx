import type { Metadata } from "next";
import "./globals.css";
import { QueryProvider } from "@/lib/query/query-provider";

export const metadata: Metadata = {
  title: "ESSK",
  description: "Enterprise SaaS Starter Kit",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <QueryProvider>{children}</QueryProvider>
      </body>
    </html>
  );
}
