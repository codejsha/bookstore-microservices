import { createFileRoute, Outlet } from "@tanstack/react-router";
import { requireAuth } from "@/shared/auth/requireAuth";

export const Route = createFileRoute("/$locale/orders")({
  beforeLoad: ({ params }) => requireAuth(params.locale),
  component: () => <Outlet />,
});
