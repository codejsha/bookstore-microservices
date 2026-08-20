import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { SubjectFindAllParams, SubjectFindAllResp } from "../types";

export const subjectListQueryOptions = (params?: SubjectFindAllParams) =>
  queryOptions({
    queryKey: ["subjects", params],
    queryFn: () => {
      const search = new URLSearchParams();
      if (params?.name) search.set("name", params.name);
      if (params?.size != null) search.set("size", String(params.size));
      if (params?.page != null) search.set("page", String(params.page));
      if (params?.sort) search.set("sort", params.sort);
      const qs = search.toString();
      return api.get<SubjectFindAllResp>(
        `${paths.subjects.list}${qs ? `?${qs}` : ""}`,
      );
    },
  });
