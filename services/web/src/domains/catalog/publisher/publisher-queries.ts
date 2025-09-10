import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type {
  Publisher,
  PublisherFindAllParams,
  PublisherFindAllResp,
} from "../types";

export const publisherListQueryOptions = (params?: PublisherFindAllParams) =>
  queryOptions({
    queryKey: ["publishers", params],
    queryFn: () => {
      const search = new URLSearchParams();
      if (params?.name) search.set("name", params.name);
      if (params?.size != null) search.set("size", String(params.size));
      if (params?.page != null) search.set("page", String(params.page));
      if (params?.sort) search.set("sort", params.sort);
      const qs = search.toString();
      return api.get<PublisherFindAllResp>(
        `${paths.publishers.list}${qs ? `?${qs}` : ""}`,
      );
    },
  });

export const publisherDetailQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: ["publishers", uid],
    queryFn: () => api.get<Publisher>(paths.publishers.detail(uid)),
  });
