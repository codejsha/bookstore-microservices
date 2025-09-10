import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Author, Paged, Subject, Work, WorkListParams } from "./types";

export const WORKS_PAGE_SIZE = 20;

function worksSearch({ title, page, size, sort }: WorkListParams): string {
  const params = new URLSearchParams();
  if (title) params.set("title", title);
  params.set("page", String(page));
  params.set("size", String(size));
  if (sort) params.set("sort", sort);
  return `?${params.toString()}`;
}

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
      api.get<Paged<Work>>(`${paths.catalog.works}${worksSearch(params)}`),
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
