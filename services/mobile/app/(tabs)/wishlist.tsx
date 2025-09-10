import { Button } from "@bookstore/design/native/button";
import { Card } from "@bookstore/design/native/card";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "expo-router";
import { BookOpen, Heart, Trash2 } from "lucide-react-native";
import { FlatList, Pressable, RefreshControl, Text, View } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { useCurrentUserUid } from "@/domains/auth";
import { bookListQueryOptions } from "@/domains/catalog";
import {
  useRemoveFromWishlist,
  wishlistQueryOptions,
} from "@/domains/wishlist";
import { ErrorState } from "@/shared/components/ErrorState";
import { SkeletonList } from "@/shared/components/Skeleton";
import { useThemeColor } from "@/shared/theme/useThemeColor";
import { useToast } from "@/shared/toast";

export default function WishlistScreen() {
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { color } = useThemeColor();
  const toast = useToast();
  const userUid = useCurrentUserUid();
  const {
    data: wishlist,
    isLoading,
    isError,
    refetch,
    isRefetching,
  } = useQuery(wishlistQueryOptions(userUid));
  const { data: books } = useQuery(bookListQueryOptions({ size: 100 }));
  const remove = useRemoveFromWishlist(userUid);

  if (isError) {
    return (
      <ErrorState
        message="Couldn't load your wishlist."
        onRetry={() => refetch()}
      />
    );
  }

  if (isLoading) {
    return (
      <View className="flex-1 bg-background">
        <SkeletonList count={5} />
      </View>
    );
  }

  const bookUids = wishlist?.book_uids ?? [];
  const wishlistBooks = (books?.items ?? []).filter((b) =>
    bookUids.includes(b.uid),
  );

  const handleRemove = (uid: string) => {
    remove.mutate(
      { book_uids: [uid] },
      {
        onSuccess: () => toast.success("Removed from wishlist"),
        onError: () => toast.error("Couldn't remove item"),
      },
    );
  };

  if (bookUids.length === 0) {
    return (
      <View className="flex-1 items-center justify-center bg-background px-5">
        <Heart color={color("muted-foreground")} size={48} />
        <Text className="mt-3 text-base font-medium text-foreground">
          Your wishlist is empty
        </Text>
        <View className="mt-4">
          <Button onPress={() => router.push("/(tabs)/books")}>
            Browse books
          </Button>
        </View>
      </View>
    );
  }

  return (
    <FlatList
      className="bg-background"
      data={wishlistBooks}
      keyExtractor={(b) => b.uid}
      contentContainerClassName="gap-3 p-5"
      contentContainerStyle={{ paddingBottom: insets.bottom + 20 }}
      refreshControl={
        <RefreshControl refreshing={isRefetching} onRefresh={() => refetch()} />
      }
      renderItem={({ item }) => (
        <Card className="flex-row items-center gap-3">
          <View className="h-20 w-14 items-center justify-center rounded-md bg-muted">
            <BookOpen color={color("muted-foreground")} size={22} />
          </View>
          <Pressable
            className="flex-1"
            onPress={() =>
              router.push({
                pathname: "/(tabs)/books/[uid]",
                params: { uid: item.uid },
              })
            }
          >
            <Text className="font-medium text-foreground" numberOfLines={2}>
              {item.title}
            </Text>
            <Text
              className="mt-1 text-xs text-muted-foreground"
              numberOfLines={1}
            >
              {item.authors.map((a) => a.name).join(", ")}
            </Text>
            <Text className="mt-1 text-sm font-semibold text-foreground">
              ${item.price.toFixed(2)}
            </Text>
          </Pressable>
          <Pressable
            onPress={() => handleRemove(item.uid)}
            className="h-8 w-8 items-center justify-center"
          >
            <Trash2 color={color("destructive")} size={16} />
          </Pressable>
        </Card>
      )}
    />
  );
}
