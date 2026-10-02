import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import { toQueryString } from "@/shared/api/query-string";
import type { Author, Paged, Subject, Work, WorkListParams } from "./types";

export const WORKS_PAGE_SIZE = 20;

export const catalogKeys = {
  workList: (params?: WorkListParams) =>
    params
      ? (["catalog", "work-list", params] as const)
      : (["catalog", "work-list"] as const),
  work: (uid: string) => ["catalog", "work", uid] as const,
} as const;

export const worksQueryOptions = (params: WorkListParams) =>
  queryOptions({
    queryKey: catalogKeys.workList(params),
    queryFn: () =>
      api.get<Paged<Work>>(
        `${paths.catalog.works}${toQueryString({ ...params })}`,
      ),
  });

export const workQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: catalogKeys.work(uid),
    queryFn: () => api.get<Work>(paths.catalog.work(uid)),
  });

export const authorsQueryOptions = () =>
  queryOptions({
    queryKey: ["catalog", "authors"],
    queryFn: () => api.get<Paged<Author>>(`${paths.catalog.authors}?size=200`),
    staleTime: 1000 * 60 * 5,
  });

export const subjectsQueryOptions = () =>
  queryOptions({
    queryKey: ["catalog", "subjects"],
    queryFn: () =>
      api.get<Paged<Subject>>(`${paths.catalog.subjects}?size=200`),
    staleTime: 1000 * 60 * 5,
  });
