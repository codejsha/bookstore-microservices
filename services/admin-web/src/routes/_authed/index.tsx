import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { dashboardQueryOptions } from "@/domains/admin";
import { MetricCard } from "@/domains/admin/MetricCard";

export const Route = createFileRoute("/_authed/")({
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(dashboardQueryOptions()),
  component: DashboardPage,
});

function DashboardPage() {
  const { data: dashboard } = useSuspenseQuery(dashboardQueryOptions());

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">Dashboard</h1>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <MetricCard label="Works" metric={dashboard.works} />
        <MetricCard label="Orders" metric={dashboard.orders} />
        <MetricCard label="Warehouses" metric={dashboard.warehouses} />
      </div>
    </div>
  );
}
