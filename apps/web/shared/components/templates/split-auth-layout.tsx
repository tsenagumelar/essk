"use client";

type SplitAuthLayoutProps = {
  media: React.ReactNode;
  children: React.ReactNode;
};

export function SplitAuthLayout({ media, children }: SplitAuthLayoutProps) {
  return (
    <main className="grid min-h-screen bg-white lg:grid-cols-[1.1fr_0.9fr]">
      <section className="hidden lg:block">{media}</section>
      <section className="flex items-center justify-center px-6 py-10">{children}</section>
    </main>
  );
}
