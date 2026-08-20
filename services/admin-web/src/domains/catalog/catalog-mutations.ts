import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import { catalogKeys } from "./catalog-queries";
import type { Work, WorkCreateRequest, WorkUpdateRequest } from "./types";

export function useCreateWork() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: WorkCreateRequest) =>
      api.post<void>(paths.catalog.works, body),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: catalogKeys.workList() }),
  });
}

export function useUpdateWork(uid: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (body: WorkUpdateRequest) =>
      api.put<Work>(paths.catalog.work(uid), body),
    onSuccess: (updated) => {
      queryClient.setQueryData(catalogKeys.work(uid), updated);
      queryClient.invalidateQueries({ queryKey: catalogKeys.workList() });
    },
  });
}
