import { queryOptions } from "@tanstack/react-query";
import { api, retryUnlessNotFound } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import { toQueryString } from "@/shared/api/query-string";
import type { Paged } from "@/shared/api/types";
import type { Order, OrderListParams, OrderSummary } from "./types";

export const orderKeys = {
  all: ["order"] as const,
  list: (params: OrderListParams) =>
    [...orderKeys.all, "list", params] as const,
  detail: (uid: string) => [...orderKeys.all, "detail", uid] as const,
};

export const ordersQueryOptions = (params: OrderListParams) =>
  queryOptions({
    queryKey: orderKeys.list(params),
    queryFn: () =>
      api.get<Paged<OrderSummary>>(
        `${paths.orders.list}${toQueryString({ ...params })}`,
      ),
  });

export const orderQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: orderKeys.detail(uid),
    queryFn: () => api.get<Order>(paths.orders.order(uid)),
    retry: retryUnlessNotFound,
  });
