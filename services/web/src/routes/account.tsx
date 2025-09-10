import { createFileRoute } from "@tanstack/react-router";
import { AccountPage } from "@/domains/account";
import { requireAuth } from "@/shared/auth/requireAuth";

export const Route = createFileRoute("/account")({
  beforeLoad: () => requireAuth(),
  component: AccountPage,
});
