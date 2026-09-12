import type { QueryClient } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { adminMeQueryOptions, isStaff } from "@/domains/admin";
import { ApiError, signInHref } from "@/shared/api/client";
import { AdminShell } from "@/shared/components/AdminShell";

const fetchIdentity = (queryClient: QueryClient) =>
  queryClient.ensureQueryData(adminMeQueryOptions());

export const Route = createFileRoute("/_authed")({
  beforeLoad: async ({ context, location }) => {
    let identity: Awaited<ReturnType<typeof fetchIdentity>>;
    try {
      identity = await fetchIdentity(context.queryClient);
    } catch (error) {
      if (error instanceof ApiError && error.isUnauthorized) {
        throw redirect({
          href: signInHref(location.pathname + location.searchStr),
          reloadDocument: true,
        });
      }
      if (error instanceof ApiError && error.isForbidden) {
        throw redirect({ to: "/unauthorized" });
      }
      throw error;
    }

    if (!isStaff(identity)) {
      throw redirect({ to: "/unauthorized" });
    }
    return { identity };
  },
  component: AuthedLayout,
});

function AuthedLayout() {
  const { identity } = Route.useRouteContext();
  return (
    <AdminShell identity={identity}>
      <Outlet />
    </AdminShell>
  );
}
