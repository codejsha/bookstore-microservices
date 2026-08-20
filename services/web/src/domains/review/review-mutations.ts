import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { ReviewCreateReq } from "./types";

export function useWriteReview(uid: string | undefined) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: ReviewCreateReq) =>
      api.post<void>(paths.customers.reviews(uid ?? ""), data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["reviews", "book", variables.book_uid],
      });
    },
  });
}
