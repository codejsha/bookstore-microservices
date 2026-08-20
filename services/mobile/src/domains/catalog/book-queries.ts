import { keepPreviousData, queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Book, BookFindAllParams, BookFindAllResp } from "./types";

const PAGE_SIZE = 12;

function buildBookSearchParams(params?: BookFindAllParams) {
  const search = new URLSearchParams();
  if (params?.title) search.set("title", params.title);
  if (params?.isbn) search.set("isbn", params.isbn);
  if (params?.category_uid) search.set("category_uid", params.category_uid);
  if (params?.publisher_uid) search.set("publisher_uid", params.publisher_uid);
  if (params?.author_uid) search.set("author_uid", params.author_uid);
  if (params?.status) search.set("status", params.status);
  search.set("size", String(params?.size ?? PAGE_SIZE));
  search.set("page", String(params?.page ?? 0));
  if (params?.sort) search.set("sort", params.sort);
  return search.toString();
}

export const bookListQueryOptions = (params?: BookFindAllParams) =>
  queryOptions({
    queryKey: ["books", params],
    queryFn: () => {
      const qs = buildBookSearchParams(params);
      return api.get<BookFindAllResp>(`${paths.books.list}?${qs}`);
    },
    placeholderData: keepPreviousData,
  });

export const bookDetailQueryOptions = (uid: string) =>
  queryOptions({
    queryKey: ["books", uid],
    queryFn: () => api.get<Book>(paths.books.detail(uid)),
  });
