"use client";

import Image from "next/image";
import { cn } from "@/lib/utils";

type AvatarProps = {
  name?: string | null;
  src?: string | null;
  size?: "sm" | "md" | "lg";
  className?: string;
};

const sizeClasses = {
  sm: "h-7 w-7 text-xs",
  md: "h-9 w-9 text-sm",
  lg: "h-12 w-12 text-base",
};

function initialsFromName(name?: string | null) {
  const source = name?.trim() || "User";
  return source
    .split(" ")
    .map((part) => part[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();
}

export function Avatar({ name, src, size = "md", className }: AvatarProps) {
  return (
    <span
      className={cn(
        "inline-flex shrink-0 items-center justify-center overflow-hidden rounded-full bg-primary font-semibold text-primary-foreground",
        sizeClasses[size],
        className,
      )}
    >
      {src ? <Image src={src} alt={name ?? "User avatar"} width={48} height={48} className="h-full w-full object-cover" /> : initialsFromName(name)}
    </span>
  );
}
