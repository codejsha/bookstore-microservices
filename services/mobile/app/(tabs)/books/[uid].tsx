import { Button } from "@bookstore/design/native/button";
import { Input } from "@bookstore/design/native/input";
import { useQuery } from "@tanstack/react-query";
import { useLocalSearchParams } from "expo-router";
import { BookOpen, Heart } from "lucide-react-native";
import { useState } from "react";
import { Pressable, ScrollView, Text, View } from "react-native";
import { useCurrentUserUid } from "@/domains/auth";
import { useAddToCart } from "@/domains/cart";
import { bookDetailQueryOptions, numericIdFromUid } from "@/domains/catalog";
import { bookReviewsQueryOptions, useWriteReview } from "@/domains/review";
import {
  useAddToWishlist,
  useRemoveFromWishlist,
  wishlistQueryOptions,
} from "@/domains/wishlist";
import { ErrorState } from "@/shared/components/ErrorState";
import { SkeletonCard } from "@/shared/components/Skeleton";
import { useThemeColor } from "@/shared/theme/useThemeColor";
import { useToast } from "@/shared/toast";

export default function BookDetailScreen() {
  const { uid } = useLocalSearchParams<{ uid: string }>();
  const productId = numericIdFromUid(uid);
  const userUid = useCurrentUserUid();
  const toast = useToast();
  const { color } = useThemeColor();

  const {
    data: book,
    isLoading,
    isError,
    refetch,
  } = useQuery(bookDetailQueryOptions(uid));
  const { data: reviews } = useQuery(bookReviewsQueryOptions(uid));
  const { data: wishlist } = useQuery(wishlistQueryOptions(userUid));

  const addToCart = useAddToCart();
  const addToWishlist = useAddToWishlist(userUid);
  const removeFromWishlist = useRemoveFromWishlist(userUid);

  if (isError) {
    return (
      <ErrorState
        message="Couldn't load this book."
        onRetry={() => refetch()}
      />
    );
  }

  if (isLoading || !book) {
    return (
      <View className="flex-1 bg-background p-5">
        <SkeletonCard />
      </View>
    );
  }

  const isWished = wishlist?.book_uids?.includes(uid) ?? false;
  const wishlistPending =
    addToWishlist.isPending || removeFromWishlist.isPending;

  const handleToggleWishlist = () => {
    if (isWished) {
      removeFromWishlist.mutate(
        { book_uids: [uid] },
        {
          onSuccess: () => toast.success("Removed from wishlist"),
          onError: () => toast.error("Couldn't update wishlist"),
        },
      );
    } else {
      addToWishlist.mutate(
        { book_uids: [uid] },
        {
          onSuccess: () => toast.success("Added to wishlist"),
          onError: () => toast.error("Couldn't update wishlist"),
        },
      );
    }
  };

  const handleAddToCart = () => {
    addToCart.mutate(
      {
        product_id: productId,
        product_name: book.title,
        quantity: 1,
        currency: "USD",
        price: book.price,
      },
      {
        onSuccess: () => toast.success("Added to cart"),
        onError: () => toast.error("Couldn't add to cart"),
      },
    );
  };

  return (
    <ScrollView
      className="flex-1 bg-background"
      contentContainerClassName="pb-12"
    >
      <View className="items-center bg-muted/40 py-10">
        <View className="h-48 w-32 items-center justify-center rounded-lg bg-muted">
          <BookOpen color={color("muted-foreground")} size={48} />
        </View>
      </View>

      <View className="px-5 pt-5">
        <Text className="text-xs uppercase tracking-wider text-muted-foreground">
          {book.category.name}
        </Text>
        <Text className="mt-1 text-2xl font-bold text-foreground">
          {book.title}
        </Text>
        <Text className="mt-1 text-base text-muted-foreground">
          by {book.authors.map((a) => a.name).join(", ")}
        </Text>

        <View className="mt-4 flex-row items-center gap-3">
          <Text className="text-2xl font-bold text-foreground">
            ${book.price.toFixed(2)}
          </Text>
          <View className="flex-1" />
          <Pressable
            onPress={handleToggleWishlist}
            disabled={wishlistPending}
            className="h-12 w-12 items-center justify-center rounded-md border border-border"
          >
            <Heart
              color={
                isWished ? color("destructive") : color("muted-foreground")
              }
              fill={isWished ? color("destructive") : "transparent"}
              size={20}
            />
          </Pressable>
        </View>

        <Button
          className="mt-3"
          disabled={addToCart.isPending}
          onPress={handleAddToCart}
        >
          {addToCart.isPending ? "Adding..." : "Add to cart"}
        </Button>

        {book.description ? (
          <View className="mt-6">
            <Text className="text-base font-semibold text-foreground">
              Description
            </Text>
            <Text className="mt-2 text-sm leading-6 text-muted-foreground">
              {book.description}
            </Text>
          </View>
        ) : null}

        <View className="mt-6">
          <Text className="text-base font-semibold text-foreground">
            Reviews
            {reviews?.total != null ? ` (${reviews.total})` : ""}
          </Text>
          <ReviewComposer bookUid={uid} userUid={userUid} />
          {(reviews?.items?.length ?? 0) === 0 ? (
            <Text className="mt-2 text-sm text-muted-foreground">
              No reviews yet.
            </Text>
          ) : (
            <View className="mt-2 gap-3">
              {reviews?.items?.map((r) => (
                <View
                  key={r.uid}
                  className="rounded-md border border-border bg-card p-3"
                >
                  <Text className="text-sm font-semibold text-foreground">
                    {"★".repeat(r.rating)}
                    {"☆".repeat(Math.max(0, 5 - r.rating))}
                    {r.title ? `  ${r.title}` : ""}
                  </Text>
                  {r.content ? (
                    <Text className="mt-1 text-sm text-muted-foreground">
                      {r.content}
                    </Text>
                  ) : null}
                </View>
              ))}
            </View>
          )}
        </View>

        <View className="mt-6">
          <DetailRow label="Publisher" value={book.publisher.name} />
          {book.isbn13 ? (
            <DetailRow label="ISBN-13" value={book.isbn13} />
          ) : null}
          {book.sku ? <DetailRow label="SKU" value={book.sku} /> : null}
          {book.publishedAt ? (
            <DetailRow
              label="Published"
              value={new Date(book.publishedAt).toLocaleDateString()}
            />
          ) : null}
        </View>
      </View>
    </ScrollView>
  );
}

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <View className="flex-row justify-between border-t border-border py-2">
      <Text className="text-sm text-muted-foreground">{label}</Text>
      <Text className="text-sm font-medium text-foreground">{value}</Text>
    </View>
  );
}

function ReviewComposer({
  bookUid,
  userUid,
}: {
  bookUid: string;
  userUid: string | null;
}) {
  const toast = useToast();
  const writeReview = useWriteReview(userUid);
  const [open, setOpen] = useState(false);
  const [rating, setRating] = useState(0);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  const reset = () => {
    setRating(0);
    setTitle("");
    setContent("");
  };

  const handleSubmit = () => {
    if (rating === 0) return;
    writeReview.mutate(
      {
        book_uid: bookUid,
        rating,
        title: title || undefined,
        content: content || undefined,
      },
      {
        onSuccess: () => {
          toast.success("Review submitted");
          reset();
          setOpen(false);
        },
        onError: () => toast.error("Couldn't submit review"),
      },
    );
  };

  if (!open) {
    return (
      <View className="mt-3">
        <Button variant="outline" size="sm" onPress={() => setOpen(true)}>
          Write a review
        </Button>
      </View>
    );
  }

  return (
    <View className="mt-3 gap-3 rounded-md border border-border bg-card p-3">
      <View className="flex-row gap-1">
        {[1, 2, 3, 4, 5].map((star) => (
          <Pressable key={star} onPress={() => setRating(star)}>
            <Text className="text-2xl text-foreground">
              {star <= rating ? "★" : "☆"}
            </Text>
          </Pressable>
        ))}
      </View>
      <Input
        placeholder="Title (optional)"
        value={title}
        onChangeText={setTitle}
      />
      <Input
        placeholder="Write your review..."
        value={content}
        onChangeText={setContent}
        multiline
        className="h-24 py-2"
        textAlignVertical="top"
      />
      <View className="flex-row justify-end gap-2">
        <Button
          variant="outline"
          size="sm"
          onPress={() => {
            reset();
            setOpen(false);
          }}
        >
          Cancel
        </Button>
        <Button
          size="sm"
          disabled={rating === 0 || writeReview.isPending || !userUid}
          onPress={handleSubmit}
        >
          {writeReview.isPending ? "Submitting..." : "Submit"}
        </Button>
      </View>
    </View>
  );
}
