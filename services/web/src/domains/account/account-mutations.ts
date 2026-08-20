import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type { Customer, CustomerUpdateRequest } from "./types";

export function useUpdateCustomer(uid: string | undefined) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CustomerUpdateRequest) =>
      api.put<Customer>(paths.customers.detail(uid ?? ""), data),
    onSuccess: (data) => {
      queryClient.setQueryData(["customer", uid], data);
    },
  });
}
