import { queryOptions } from "@tanstack/react-query";
import { api, retryUnlessNotFound } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import { toQueryString } from "@/shared/api/query-string";
import type { Paged } from "@/shared/api/types";
import type { User, UserListParams } from "./types";

export const userKeys = {
  all: ["user"] as const,
  list: (params: UserListParams) => [...userKeys.all, "list", params] as const,
  detail: (uid: string) => [...userKeys.all, "detail", uid] as const,
};

export const usersQueryOptions = (params: UserListParams) =>
  queryOptions({
    queryKey: userKeys.list(params),
    queryFn: () =>
      api.get<Paged<User>>(
        `${paths.users.list}${toQueryString({ ...params })}`,
      ),
  });

export const userQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: userKeys.detail(uid),
    queryFn: () => api.get<User>(paths.users.user(uid)),
    retry: retryUnlessNotFound,
  });
