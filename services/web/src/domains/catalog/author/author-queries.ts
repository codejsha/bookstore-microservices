import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Author, AuthorFindAllParams, AuthorFindAllResp } from "../types";

export const authorListQueryOptions = (params?: AuthorFindAllParams) =>
  queryOptions({
    queryKey: ["authors", params],
    queryFn: () => {
      const search = new URLSearchParams();
      if (params?.name) search.set("name", params.name);
      if (params?.size != null) search.set("size", String(params.size));
      if (params?.page != null) search.set("page", String(params.page));
      if (params?.sort) search.set("sort", params.sort);
      const qs = search.toString();
      return api.get<AuthorFindAllResp>(
        `${paths.authors.list}${qs ? `?${qs}` : ""}`,
      );
    },
  });

export const authorDetailQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: ["authors", uid],
    queryFn: () => api.get<Author>(paths.authors.detail(uid)),
  });
