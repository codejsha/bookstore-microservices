import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Customer, PointBalance, PointHistoryResp } from "./types";

export const customerQueryOptions = (uid: string | undefined) =>
  queryOptions({
    queryKey: ["customer", uid],
    queryFn: () => api.get<Customer>(paths.customers.detail(uid ?? "")),
    enabled: !!uid,
  });

export const pointBalanceQueryOptions = (uid: string | undefined) =>
  queryOptions({
    queryKey: ["points", uid, "balance"],
    queryFn: () => api.get<PointBalance>(paths.customers.points(uid ?? "")),
    enabled: !!uid,
  });

export const pointHistoryQueryOptions = (
  uid: string | undefined,
  params?: { size?: number; page?: number },
) =>
  queryOptions({
    queryKey: ["points", uid, "history", params],
    queryFn: () => {
      const search = new URLSearchParams();
      if (params?.size != null) search.set("size", String(params.size));
      if (params?.page != null) search.set("page", String(params.page));
      const qs = search.toString();
      return api.get<PointHistoryResp>(
        `${paths.customers.pointsHistory(uid ?? "")}${qs ? `?${qs}` : ""}`,
      );
    },
    enabled: !!uid,
  });
