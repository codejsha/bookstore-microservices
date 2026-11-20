import { createFileRoute } from "@tanstack/react-router";
import { CartPage } from "@/domains/cart";
import { requireAuth } from "@/shared/auth/requireAuth";

export const Route = createFileRoute("/$locale/cart")({
  beforeLoad: ({ params }) => requireAuth(params.locale),
  component: CartPage,
});
