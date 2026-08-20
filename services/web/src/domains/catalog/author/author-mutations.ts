import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Author, AuthorCreateReq, AuthorUpdateReq } from "../types";

export function useCreateAuthor() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: AuthorCreateReq) =>
      api.post<void>(paths.authors.list, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["authors"] });
    },
  });
}

export function useUpdateAuthor() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ uid, data }: { uid: string; data: AuthorUpdateReq }) =>
      api.put<Author>(paths.authors.detail(uid), data),
    onSuccess: (_data, { uid }) => {
      queryClient.invalidateQueries({ queryKey: ["authors"] });
      queryClient.invalidateQueries({ queryKey: ["authors", uid] });
    },
  });
}
