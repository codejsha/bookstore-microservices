import {
  infiniteQueryOptions,
  keepPreviousData,
  queryOptions,
} from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type {
  EditionFindAllResp,
  PaginationParams,
  Work,
  WorkFindAllParams,
  WorkFindAllResp,
  WorkSearchParams,
  WorkSearchResp,
} from "../types";

const PAGE_SIZE = 12;

function buildWorkSearchParams(
  params?: Omit<WorkFindAllParams, "page">,
  page = 0,
) {
  const search = new URLSearchParams();
  if (params?.title) search.set("title", params.title);
  if (params?.authorUid) search.set("authorUid", params.authorUid);
  if (params?.subjectUid) search.set("subjectUid", params.subjectUid);
  if (params?.olKey) search.set("olKey", params.olKey);
  search.set("size", String(params?.size ?? PAGE_SIZE));
  search.set("page", String(page));
  if (params?.sort) search.set("sort", params.sort);
  return search.toString();
}

export const workListInfiniteQueryOptions = (
  params?: Omit<WorkFindAllParams, "page">,
) =>
  infiniteQueryOptions({
    queryKey: ["works", "infinite", params],
    queryFn: ({ pageParam = 0 }) => {
      const qs = buildWorkSearchParams(params, pageParam);
      return api.get<WorkFindAllResp>(`${paths.works.list}?${qs}`);
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage, _allPages, lastPageParam) => {
      const size = params?.size ?? PAGE_SIZE;
      const loaded = (lastPageParam + 1) * size;
      return loaded < lastPage.total ? lastPageParam + 1 : undefined;
    },
    placeholderData: keepPreviousData,
  });

export const workListQueryOptions = (params?: WorkFindAllParams) =>
  queryOptions({
    queryKey: ["works", params],
    queryFn: () => {
      const qs = buildWorkSearchParams(params, params?.page);
      return api.get<WorkFindAllResp>(`${paths.works.list}?${qs}`);
    },
  });

export const workFullTextSearchQueryOptions = (params: WorkSearchParams) =>
  queryOptions({
    queryKey: ["works", "search", params],
    queryFn: () => {
      const search = new URLSearchParams();
      if (params.q) search.set("q", params.q);
      if (params.authorUid) search.set("authorUid", params.authorUid);
      if (params.subjectUid) search.set("subjectUid", params.subjectUid);
      if (params.size != null) search.set("size", String(params.size));
      if (params.page != null) search.set("page", String(params.page));
      const qs = search.toString();
      return api.get<WorkSearchResp>(
        `${paths.works.search}${qs ? `?${qs}` : ""}`,
      );
    },
  });

export const workDetailQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: ["works", uid],
    queryFn: () => api.get<Work>(paths.works.detail(uid)),
  });

export const workEditionsQueryOptions = (
  uid: string,
  params?: PaginationParams,
) =>
  queryOptions({
    queryKey: ["works", uid, "editions", params],
    queryFn: () => {
      const search = new URLSearchParams();
      if (params?.size != null) search.set("size", String(params.size));
      if (params?.page != null) search.set("page", String(params.page));
      if (params?.sort) search.set("sort", params.sort);
      const qs = search.toString();
      return api.get<EditionFindAllResp>(
        `${paths.works.editions(uid)}${qs ? `?${qs}` : ""}`,
      );
    },
  });
