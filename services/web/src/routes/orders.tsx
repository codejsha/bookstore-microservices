import { createFileRoute, Outlet } from "@tanstack/react-router";
import { requireAuth } from "@/shared/auth/requireAuth";

export const Route = createFileRoute("/orders")({
  beforeLoad: () => requireAuth(),
  component: () => <Outlet />,
});
