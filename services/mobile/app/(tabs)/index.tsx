import { Button } from "@bookstore/design/native/button";
import { Card } from "@bookstore/design/native/card";
import { useQuery } from "@tanstack/react-query";
import { Link, useRouter } from "expo-router";
import { BookOpen, ChevronRight } from "lucide-react-native";
import {
  Pressable,
  RefreshControl,
  ScrollView,
  Text,
  View,
} from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import { bookListQueryOptions } from "@/domains/catalog";
import { ErrorState } from "@/shared/components/ErrorState";
import { SkeletonCard } from "@/shared/components/Skeleton";
import { useThemeColor } from "@/shared/theme/useThemeColor";

export default function HomeScreen() {
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { color } = useThemeColor();
  const { data, isLoading, isError, refetch, isRefetching } = useQuery(
    bookListQueryOptions({ size: 6, sort: "createdAt,desc" }),
  );
  const books = data?.items ?? [];

  return (
    <ScrollView
      className="flex-1 bg-background"
      contentContainerStyle={{ paddingBottom: insets.bottom }}
      refreshControl={
        <RefreshControl refreshing={isRefetching} onRefresh={() => refetch()} />
      }
    >
      <View className="bg-primary/10 px-5 pt-10 pb-12">
        <Text className="text-3xl font-bold text-foreground">Bookstore</Text>
        <Text className="mt-2 text-base text-muted-foreground">
          Browse and discover great reads.
        </Text>
        <View className="mt-6 flex-row gap-2">
          <Button onPress={() => router.push("/(tabs)/books")}>
            Browse books
          </Button>
          <Button
            variant="outline"
            onPress={() => router.push("/(tabs)/wishlist")}
          >
            My wishlist
          </Button>
        </View>
      </View>

      <View className="px-5 py-6">
        <View className="mb-3 flex-row items-center justify-between">
          <Text className="text-lg font-semibold text-foreground">
            New arrivals
          </Text>
          <Link href="/(tabs)/books" className="flex-row items-center">
            <Text className="text-sm text-muted-foreground">See all</Text>
            <ChevronRight color={color("muted-foreground")} size={16} />
          </Link>
        </View>

        {isError ? (
          <ErrorState
            className="items-center py-16"
            message="Couldn't load new arrivals."
            onRetry={() => refetch()}
          />
        ) : isLoading ? (
          <View className="gap-3">
            {Array.from({ length: 4 }).map((_, i) => (
              // biome-ignore lint/suspicious/noArrayIndexKey: static placeholder list
              <SkeletonCard key={i} />
            ))}
          </View>
        ) : (
          <View className="gap-3">
            {books.slice(0, 6).map((book) => (
              <Pressable
                key={book.uid}
                onPress={() =>
                  router.push({
                    pathname: "/(tabs)/books/[uid]",
                    params: { uid: book.uid },
                  })
                }
              >
                <Card className="flex-row items-center gap-3">
                  <View className="h-16 w-12 items-center justify-center rounded-md bg-muted">
                    <BookOpen color={color("muted-foreground")} size={20} />
                  </View>
                  <View className="flex-1">
                    <Text
                      className="font-medium text-foreground"
                      numberOfLines={2}
                    >
                      {book.title}
                    </Text>
                    <Text className="mt-1 text-xs text-muted-foreground">
                      ${book.price.toFixed(2)}
                    </Text>
                  </View>
                </Card>
              </Pressable>
            ))}
          </View>
        )}
      </View>
    </ScrollView>
  );
}
