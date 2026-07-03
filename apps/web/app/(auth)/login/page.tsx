import Image from "next/image";
import { LoginForm } from "@/features/auth/components/login-form";

export default function LoginPage() {
  return (
    <main className="grid min-h-screen bg-white lg:grid-cols-[1.1fr_0.9fr]">
      <section className="relative hidden overflow-hidden lg:block">
        <Image src="/login-hero.png" alt="Enterprise SaaS operations workspace" fill priority className="object-cover" />
        <div className="absolute inset-0 bg-gradient-to-br from-black/60 via-black/20 to-primary/30" />
        <div className="absolute inset-x-0 bottom-0 p-12 text-white">
          <div className="mb-8 flex items-center gap-3">
            <div className="flex h-11 w-11 items-center justify-center rounded-lg bg-white/15 text-sm font-bold backdrop-blur">
              ES
            </div>
            <div>
              <p className="font-semibold">ESSK</p>
              <p className="text-sm text-white/75">Enterprise SaaS Starter Kit</p>
            </div>
          </div>
          <h1 className="max-w-xl text-4xl font-semibold leading-tight">Build enterprise products from a production-ready foundation.</h1>
          <div className="mt-8 grid max-w-2xl grid-cols-3 gap-3 text-sm">
            <div className="rounded-lg border border-white/20 bg-white/10 p-3 backdrop-blur">
              Modular CRUD
            </div>
            <div className="rounded-lg border border-white/20 bg-white/10 p-3 backdrop-blur">
              RBAC & Audit
            </div>
            <div className="rounded-lg border border-white/20 bg-white/10 p-3 backdrop-blur">
              Secure by Default
            </div>
          </div>
        </div>
      </section>

      <section className="relative flex min-h-screen items-center justify-center overflow-hidden px-6 py-10">
        <div className="pointer-events-none absolute -right-20 top-16 select-none text-[180px] font-black leading-none text-primary/[0.04] md:text-[260px]">
          ESSK
        </div>
        <div className="w-full max-w-md">
          <div className="mb-8 lg:hidden">
            <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-sm font-bold text-primary-foreground">
              ES
            </div>
            <p className="text-sm text-muted-foreground">Enterprise SaaS Starter Kit</p>
          </div>
          <div className="rounded-xl border bg-white/95 p-6 shadow-sm backdrop-blur">
            <div className="mb-6">
              <p className="text-sm font-medium text-primary">Welcome back</p>
              <h2 className="mt-2 text-2xl font-semibold">Sign in to ESSK</h2>
              <p className="mt-2 text-sm text-muted-foreground">Use your workspace account to continue.</p>
            </div>
            <LoginForm />
          </div>
          <p className="mt-6 text-center text-xs text-muted-foreground">Protected workspace for enterprise SaaS development.</p>
        </div>
      </section>
    </main>
  );
}
