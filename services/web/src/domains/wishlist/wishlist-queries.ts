import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Wishlist } from "./types";

export const wishlistQueryOptions = (uid: string | undefined) =>
  queryOptions({
    queryKey: ["wishlist", uid],
    queryFn: () => api.get<Wishlist>(paths.customers.wishlist(uid ?? "")),
    enabled: !!uid,
  });
