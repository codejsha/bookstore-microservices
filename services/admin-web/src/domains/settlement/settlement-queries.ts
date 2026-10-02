import { queryOptions } from "@tanstack/react-query";
import { api, retryUnlessNotFound } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import { toQueryString } from "@/shared/api/query-string";
import type { Paged } from "@/shared/api/types";
import type {
  Settlement,
  SettlementListParams,
  SettlementSummary,
} from "./types";

export const settlementKeys = {
  all: ["settlement"] as const,
  list: (params: SettlementListParams) =>
    [...settlementKeys.all, "list", params] as const,
  detail: (uid: string) => [...settlementKeys.all, "detail", uid] as const,
};

export const settlementsQueryOptions = (params: SettlementListParams) =>
  queryOptions({
    queryKey: settlementKeys.list(params),
    queryFn: () =>
      api.get<Paged<SettlementSummary>>(
        `${paths.settlements.list}${toQueryString({ ...params })}`,
      ),
  });

export const settlementQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: settlementKeys.detail(uid),
    queryFn: () => api.get<Settlement>(paths.settlements.settlement(uid)),
    retry: retryUnlessNotFound,
  });
