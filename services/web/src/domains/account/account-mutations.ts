import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { User, UserUpdateRequest } from "./types";

export function useUpdateCustomer(uid: string | undefined) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: UserUpdateRequest) =>
      api.put<User>(paths.users.detail(uid ?? ""), data),
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: ["customer", uid] }),
  });
}
