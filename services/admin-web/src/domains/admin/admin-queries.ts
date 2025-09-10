import { queryOptions } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { AdminIdentity, Dashboard } from "./types";
import { ROLE_ADMIN } from "./types";

export const adminMeQueryOptions = () =>
  queryOptions({
    queryKey: ["admin", "me"],
    queryFn: () => api.get<AdminIdentity>(paths.admin.me),
    staleTime: 1000 * 60,
    retry: false,
  });

export const dashboardQueryOptions = () =>
  queryOptions({
    queryKey: ["admin", "dashboard"],
    queryFn: () => api.get<Dashboard>(paths.admin.dashboard),
  });

export function isAdmin(identity: AdminIdentity | undefined): boolean {
  return identity?.roles.includes(ROLE_ADMIN) ?? false;
}
