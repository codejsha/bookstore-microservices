import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Work, WorkCreateReq, WorkUpdateReq } from "../types";

export function useCreateWork() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: WorkCreateReq) => api.post<void>(paths.works.list, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["works"] });
    },
  });
}

export function useUpdateWork() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ uid, data }: { uid: string; data: WorkUpdateReq }) =>
      api.put<Work>(paths.works.detail(uid), data),
    onSuccess: (_data, { uid }) => {
      queryClient.invalidateQueries({ queryKey: ["works"] });
      queryClient.invalidateQueries({ queryKey: ["works", uid] });
    },
  });
}
