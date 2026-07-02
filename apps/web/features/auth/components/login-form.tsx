"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { login } from "@/features/auth/api";
import { loginSchema, type LoginFormValues } from "@/features/auth/schema";
import { storeSession } from "@/features/auth/session";

export function LoginForm() {
  const router = useRouter();
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const mutation = useMutation({
    mutationFn: login,
    onSuccess: (result) => {
      storeSession(result);
      router.push("/dashboard");
    },
  });

  async function onSubmit(values: LoginFormValues) {
    await mutation.mutateAsync(values);
  }

  return (
    <form className="space-y-4" onSubmit={handleSubmit(onSubmit)}>
      <div>
        <label className="text-sm font-medium" htmlFor="email">
          Email
        </label>
        <input
          id="email"
          type="email"
          className="mt-1 w-full rounded-md border px-3 py-2 text-sm outline-none ring-primary focus:ring-2"
          {...register("email")}
        />
        {errors.email ? <p className="mt-1 text-xs text-destructive">{errors.email.message}</p> : null}
      </div>

      <div>
        <label className="text-sm font-medium" htmlFor="password">
          Password
        </label>
        <input
          id="password"
          type="password"
          className="mt-1 w-full rounded-md border px-3 py-2 text-sm outline-none ring-primary focus:ring-2"
          {...register("password")}
        />
        {errors.password ? <p className="mt-1 text-xs text-destructive">{errors.password.message}</p> : null}
      </div>

      {mutation.isError ? <p className="text-sm text-destructive">Invalid email or password.</p> : null}

      <button
        type="submit"
        disabled={isSubmitting || mutation.isPending}
        className="w-full rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground disabled:opacity-60"
      >
        {isSubmitting || mutation.isPending ? "Signing in..." : "Sign in"}
      </button>
    </form>
  );
}
