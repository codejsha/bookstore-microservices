import { createFileRoute } from "@tanstack/react-router";
import { AccountPage } from "@/domains/account";
import { requireAuth } from "@/shared/auth/requireAuth";

export const Route = createFileRoute("/$locale/account")({
  beforeLoad: ({ params }) => requireAuth(params.locale),
  component: AccountPage,
});
