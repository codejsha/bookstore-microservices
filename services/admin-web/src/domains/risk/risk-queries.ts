import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { RiskEntryList } from "./types";

export const riskKeys = {
  all: ["risk"] as const,
  list: () => [...riskKeys.all, "list"] as const,
};

export const riskListQueryOptions = () =>
  queryOptions({
    queryKey: riskKeys.list(),
    queryFn: () => api.get<RiskEntryList>(paths.risk.list),
  });
