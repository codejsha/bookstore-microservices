import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/shared/api/client";
import { paths } from "@/shared/api/paths";
import type {
  Wishlist,
  WishlistAddRequest,
  WishlistRemoveRequest,
} from "./types";

export function useAddToWishlist(uid: string | undefined) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: WishlistAddRequest) =>
      api.post<Wishlist>(paths.customers.wishlistAdd(uid ?? ""), data),
    onSuccess: (data) => {
      queryClient.setQueryData(["wishlist", uid], data);
    },
  });
}

export function useRemoveFromWishlist(uid: string | undefined) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: WishlistRemoveRequest) =>
      api.post<Wishlist>(paths.customers.wishlistRemove(uid ?? ""), data),
    onSuccess: (data) => {
      queryClient.setQueryData(["wishlist", uid], data);
    },
  });
}
