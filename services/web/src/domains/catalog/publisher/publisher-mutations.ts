import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type {
  Publisher,
  PublisherCreateReq,
  PublisherUpdateReq,
} from "../types";

export function useCreatePublisher() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: PublisherCreateReq) =>
      api.post<void>(paths.publishers.list, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["publishers"] });
    },
  });
}

export function useUpdatePublisher() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ uid, data }: { uid: string; data: PublisherUpdateReq }) =>
      api.put<Publisher>(paths.publishers.detail(uid), data),
    onSuccess: (_data, { uid }) => {
      queryClient.invalidateQueries({ queryKey: ["publishers"] });
      queryClient.invalidateQueries({ queryKey: ["publishers", uid] });
    },
  });
}
