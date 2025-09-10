import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import { riskKeys } from "./risk-queries";
import type { RiskEntry, RiskFlagRequest } from "./types";

export function useFlagRisk() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ uid, body }: { uid: string; body: RiskFlagRequest }) =>
      api.put<RiskEntry>(paths.risk.entry(uid), body),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: riskKeys.list() }),
  });
}

export function useUnflagRisk() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (uid: string) => api.delete<void>(paths.risk.entry(uid)),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: riskKeys.list() }),
  });
}
