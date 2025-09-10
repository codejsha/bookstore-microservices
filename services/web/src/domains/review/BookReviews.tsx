import { Button } from "@bookstore/design/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@bookstore/design/ui/card";
import { Input } from "@bookstore/design/ui/input";
import { Separator } from "@bookstore/design/ui/separator";
import { Textarea } from "@bookstore/design/ui/textarea";
import { useQuery } from "@tanstack/react-query";
import { MessageSquare, Star, User } from "lucide-react";
import { useState } from "react";
import { useCurrentUserUid } from "@/shared/lib/currentUser";
import { useWriteReview } from "./review-mutations";
import { bookReviewsQueryOptions } from "./review-queries";
import type { Review } from "./types";

function StarRating({
  rating,
  onRate,
  interactive = false,
}: {
  rating: number;
  onRate?: (r: number) => void;
  interactive?: boolean;
}) {
  return (
    <div className="flex gap-0.5">
      {[1, 2, 3, 4, 5].map((star) => (
        <button
          key={star}
          type="button"
          disabled={!interactive}
          onClick={() => onRate?.(star)}
          className={interactive ? "cursor-pointer" : "cursor-default"}
        >
          <Star
            className={`size-4 ${
              star <= rating
                ? "fill-yellow-400 text-yellow-400"
                : "text-muted-foreground/30"
            }`}
          />
        </button>
      ))}
    </div>
  );
}

function ReviewCard({ review }: { review: Review }) {
  return (
    <div className="space-y-2 py-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-muted">
            <User className="size-4 text-muted-foreground" />
          </div>
          <div>
            <p className="text-sm font-medium">
              User {review.user_uid.slice(0, 8)}
            </p>
            <p className="text-xs text-muted-foreground">
              {new Date(review.created_at).toLocaleDateString()}
            </p>
          </div>
        </div>
        <StarRating rating={review.rating} />
      </div>
      {review.title && <p className="text-sm font-medium">{review.title}</p>}
      {review.content && (
        <p className="text-sm text-muted-foreground">{review.content}</p>
      )}
    </div>
  );
}

function ReviewForm({
  bookUid,
  onSuccess,
}: {
  bookUid: string;
  onSuccess: () => void;
}) {
  const [rating, setRating] = useState(0);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const uid = useCurrentUserUid();
  const mutation = useWriteReview(uid);

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (rating === 0) return;
    mutation.mutate(
      {
        book_uid: bookUid,
        rating,
        title: title || undefined,
        content: content || undefined,
      },
      {
        onSuccess: () => {
          setRating(0);
          setTitle("");
          setContent("");
          onSuccess();
        },
      },
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="space-y-1">
        <p className="text-sm font-medium">Your rating</p>
        <StarRating rating={rating} onRate={setRating} interactive />
      </div>
      <Input
        placeholder="Review title (optional)"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
      />
      <Textarea
        placeholder="Write your review..."
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={3}
      />
      <Button
        type="submit"
        size="sm"
        disabled={rating === 0 || mutation.isPending || !uid}
      >
        {mutation.isPending ? "Submitting..." : "Submit Review"}
      </Button>
    </form>
  );
}

export function BookReviews({ bookUid }: { bookUid: string }) {
  const { data, refetch } = useQuery(bookReviewsQueryOptions(bookUid));
  const [showForm, setShowForm] = useState(false);

  const reviews = data?.items ?? [];
  const avgRating =
    reviews.length > 0
      ? reviews.reduce((sum, r) => sum + r.rating, 0) / reviews.length
      : 0;

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-base">
            <MessageSquare className="size-4" />
            Reviews ({data?.total ?? 0})
          </CardTitle>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setShowForm(!showForm)}
          >
            {showForm ? "Cancel" : "Write a Review"}
          </Button>
        </div>
        {reviews.length > 0 && (
          <div className="flex items-center gap-2 mt-1">
            <StarRating rating={Math.round(avgRating)} />
            <span className="text-sm text-muted-foreground">
              {avgRating.toFixed(1)} out of 5
            </span>
          </div>
        )}
      </CardHeader>
      <CardContent className="space-y-0">
        {showForm && (
          <>
            <ReviewForm
              bookUid={bookUid}
              onSuccess={() => {
                setShowForm(false);
                refetch();
              }}
            />
            <Separator className="my-4" />
          </>
        )}

        {reviews.length === 0 ? (
          <p className="py-6 text-center text-sm text-muted-foreground">
            No reviews yet. Be the first to review this book!
          </p>
        ) : (
          <div className="divide-y">
            {reviews.map((review) => (
              <ReviewCard key={review.uid} review={review} />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
