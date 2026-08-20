import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Order, OrderFindAllResp, OrderQueryParams } from "./types";

export const orderListQueryOptions = (params?: OrderQueryParams) =>
  queryOptions({
    queryKey: ["orders", params],
    queryFn: () => {
      const search = new URLSearchParams();
      if (params?.user_id != null)
        search.set("user_id", String(params.user_id));
      if (params?.status) search.set("status", params.status);
      if (params?.size != null) search.set("size", String(params.size));
      if (params?.page != null) search.set("page", String(params.page));
      if (params?.sort) search.set("sort", params.sort);
      const qs = search.toString();
      return api.get<OrderFindAllResp>(
        `${paths.orders.list}${qs ? `?${qs}` : ""}`,
      );
    },
  });

export const orderDetailQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: ["orders", uid],
    queryFn: () => api.get<Order>(paths.orders.detail(uid)),
  });
