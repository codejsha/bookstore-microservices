import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { CategoryFindAllResp } from "./types";

export const categoryListQueryOptions = () =>
  queryOptions({
    queryKey: ["categories"],
    queryFn: () => api.get<CategoryFindAllResp>(paths.categories.list),
    staleTime: 1000 * 60 * 5,
  });
