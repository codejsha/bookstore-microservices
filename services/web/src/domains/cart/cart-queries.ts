import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Cart } from "./types";

export const cartQueryOptions = () =>
  queryOptions({
    queryKey: ["cart"],
    queryFn: () => api.get<Cart>(paths.cart.root),
  });
