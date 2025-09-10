import { Card } from "@bookstore/design/native/card";
import { Input } from "@bookstore/design/native/input";
import { useQuery } from "@tanstack/react-query";
import { useRouter } from "expo-router";
import { BookOpen, Search } from "lucide-react-native";
import { useEffect, useState } from "react";
import { FlatList, Pressable, RefreshControl, Text, View } from "react-native";
import { useSafeAreaInsets } from "react-native-safe-area-context";
import {
  bookListQueryOptions,
  categoryListQueryOptions,
} from "@/domains/catalog";
import { ErrorState } from "@/shared/components/ErrorState";
import { SkeletonList } from "@/shared/components/Skeleton";
import { useThemeColor } from "@/shared/theme/useThemeColor";

const PAGE_SIZE = 20;

export default function BooksScreen() {
  const router = useRouter();
  const insets = useSafeAreaInsets();
  const { color } = useThemeColor();
  const [title, setTitle] = useState("");
  const [debouncedTitle, setDebouncedTitle] = useState("");
  const [categoryUid, setCategoryUid] = useState<string | undefined>();

  useEffect(() => {
    const t = setTimeout(() => setDebouncedTitle(title), 400);
    return () => clearTimeout(t);
  }, [title]);

  const { data: categories } = useQuery(categoryListQueryOptions());
  const { data, isLoading, isError, refetch, isRefetching } = useQuery(
    bookListQueryOptions({
      title: debouncedTitle || undefined,
      category_uid: categoryUid,
      size: PAGE_SIZE,
    }),
  );

  const books = data?.items ?? [];
  const total = data?.total ?? 0;

  return (
    <View className="flex-1 bg-background">
      <View className="border-b border-border bg-background px-5 py-4">
        <View className="relative">
          <View className="absolute left-3 top-3.5 z-10">
            <Search color={color("muted-foreground")} size={18} />
          </View>
          <Input
            value={title}
            onChangeText={setTitle}
            placeholder="Search by title..."
            className="pl-10"
          />
        </View>

        <FlatList
          horizontal
          showsHorizontalScrollIndicator={false}
          className="mt-3"
          data={[{ uid: undefined, name: "All" }, ...(categories?.items ?? [])]}
          keyExtractor={(item) => item.uid ?? "all"}
          contentContainerClassName="gap-2"
          renderItem={({ item }) => {
            const active = categoryUid === item.uid;
            return (
              <Text
                onPress={() => setCategoryUid(item.uid)}
                className={
                  active
                    ? "rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground"
                    : "rounded-md bg-muted px-3 py-1.5 text-sm font-medium text-muted-foreground"
                }
              >
                {item.name}
              </Text>
            );
          }}
        />

        <Text className="mt-2 text-xs text-muted-foreground">
          {total > 0
            ? `${total} book${total !== 1 ? "s" : ""}`
            : isLoading
              ? "Loading..."
              : "No books"}
        </Text>
      </View>

      {isError ? (
        <ErrorState message="Couldn't load books." onRetry={() => refetch()} />
      ) : isLoading ? (
        <SkeletonList count={6} />
      ) : (
        <FlatList
          data={books}
          keyExtractor={(item) => item.uid}
          contentContainerClassName="gap-3 p-5"
          contentContainerStyle={{ paddingBottom: insets.bottom + 20 }}
          refreshControl={
            <RefreshControl
              refreshing={isRefetching}
              onRefresh={() => refetch()}
            />
          }
          ListEmptyComponent={
            <View className="items-center py-20">
              <Text className="text-base font-medium text-foreground">
                No books found
              </Text>
              <Text className="mt-1 text-sm text-muted-foreground">
                Try adjusting your search.
              </Text>
            </View>
          }
          renderItem={({ item }) => (
            <Pressable
              onPress={() =>
                router.push({
                  pathname: "/(tabs)/books/[uid]",
                  params: { uid: item.uid },
                })
              }
            >
              <Card className="flex-row items-center gap-3">
                <View className="h-20 w-14 items-center justify-center rounded-md bg-muted">
                  <BookOpen color={color("muted-foreground")} size={22} />
                </View>
                <View className="flex-1">
                  <Text
                    className="font-medium text-foreground"
                    numberOfLines={2}
                  >
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
                </View>
              </Card>
            </Pressable>
          )}
        />
      )}
    </View>
  );
}
