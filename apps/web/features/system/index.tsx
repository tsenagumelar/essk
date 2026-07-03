import { HealthStatus } from "@/features/system/components/health-status";

export function HealthView() {
  return (
    <div>
      <h1 className="text-2xl font-semibold">System Health</h1>
      <div className="mt-4">
        <HealthStatus />
      </div>
    </div>
  );
}
