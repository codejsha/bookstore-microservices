import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { ReviewFindAllResp } from "./types";

export const bookReviewsQueryOptions = (bookUid: string) =>
  queryOptions({
    queryKey: ["reviews", "book", bookUid],
    queryFn: () => api.get<ReviewFindAllResp>(paths.reviews.byBook(bookUid)),
  });
