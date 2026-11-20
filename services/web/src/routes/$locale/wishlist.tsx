import { createFileRoute } from "@tanstack/react-router";
import { WishlistPage } from "@/domains/wishlist";
import { requireAuth } from "@/shared/auth/requireAuth";

export const Route = createFileRoute("/$locale/wishlist")({
  beforeLoad: ({ params }) => requireAuth(params.locale),
  component: WishlistPage,
});
