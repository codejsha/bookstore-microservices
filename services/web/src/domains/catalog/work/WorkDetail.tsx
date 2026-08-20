import { Badge } from "@bookstore/design/ui/badge";
import { Button } from "@bookstore/design/ui/button";
import { Card, CardContent } from "@bookstore/design/ui/card";
import { Separator } from "@bookstore/design/ui/separator";
import { useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import {
  ArrowLeft,
  BookOpen,
  Calendar,
  Hash,
  Heart,
  Layers,
  ShoppingCart,
} from "lucide-react";
import { useAddToCart } from "@/domains/cart";
import { BookReviews } from "@/domains/review";
import {
  useAddToWishlist,
  useRemoveFromWishlist,
  wishlistQueryOptions,
} from "@/domains/wishlist";
import { useToast } from "@/shared/components/toast";
import { useCurrentUserUid } from "@/shared/lib/currentUser";
import { numericIdFromUid, syntheticPrice } from "../catalog-utils";
import type { Edition } from "../types";
import {
  workDetailQueryOptions,
  workEditionsQueryOptions,
} from "./work-queries";

interface WorkDetailProps {
  uid: string;
}

export function WorkDetail({ uid }: WorkDetailProps) {
  const { data: work } = useSuspenseQuery(workDetailQueryOptions(uid));
  const { data: editionData } = useQuery(workEditionsQueryOptions(uid));
  const toast = useToast();

  const editions = editionData?.items ?? [];

  const currentUserUid = useCurrentUserUid();

  const { data: wishlist } = useQuery(wishlistQueryOptions(currentUserUid));
  const isWished = wishlist?.book_uids?.includes(uid) ?? false;
  const addToWishlist = useAddToWishlist(currentUserUid);
  const removeFromWishlist = useRemoveFromWishlist(currentUserUid);
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

  const year = work.first_publish_date
    ? work.first_publish_date.slice(0, 4)
    : undefined;

  return (
    <div className="container mx-auto max-w-5xl px-4 py-8 space-y-8">
      <Button variant="ghost" size="sm" render={<Link to="/books" />}>
        <ArrowLeft className="size-4 mr-1" />
        Back to books
      </Button>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
        <div className="lg:col-span-1">
          <Card>
            <CardContent className="flex aspect-[3/4] items-center justify-center bg-muted rounded-lg">
              <BookOpen className="size-16 text-muted-foreground" />
            </CardContent>
          </Card>
        </div>

        <div className="lg:col-span-2 space-y-6">
          <div className="space-y-3">
            <div className="flex flex-wrap items-center gap-2">
              {work.subjects.map((s) => (
                <Badge key={s.uid} variant="secondary">
                  {s.name}
                </Badge>
              ))}
            </div>

            <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">
              {work.title}
            </h1>

            <p className="text-lg text-muted-foreground">
              by {work.authors.map((a) => a.name).join(", ")}
            </p>

            {year && (
              <p className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <Calendar className="size-3.5" />
                First published {year}
              </p>
            )}
          </div>

          <div className="flex items-center gap-4">
            <Button
              variant="outline"
              size="icon"
              aria-label={isWished ? "Remove from wishlist" : "Add to wishlist"}
              disabled={wishlistPending || !currentUserUid}
              onClick={handleToggleWishlist}
            >
              <Heart
                className={
                  isWished ? "size-4 fill-red-500 text-red-500" : "size-4"
                }
              />
            </Button>
          </div>

          {work.description && (
            <div className="space-y-2">
              <h2 className="text-lg font-semibold">Description</h2>
              <p className="text-muted-foreground leading-relaxed whitespace-pre-line">
                {work.description}
              </p>
            </div>
          )}

          {work.authors.some((a) => a.bio) && (
            <div className="space-y-3">
              <h2 className="text-lg font-semibold">
                About the Author{work.authors.length > 1 ? "s" : ""}
              </h2>
              {work.authors.map(
                (author) =>
                  author.bio && (
                    <div key={author.uid} className="space-y-1">
                      <p className="font-medium">{author.name}</p>
                      <p className="text-sm text-muted-foreground">
                        {author.bio}
                      </p>
                    </div>
                  ),
              )}
            </div>
          )}
        </div>
      </div>

      <Separator />
      <div className="space-y-4">
        <div className="flex items-center gap-2">
          <Layers className="size-5" />
          <h2 className="text-xl font-semibold">
            Editions{editions.length > 0 ? ` (${editions.length})` : ""}
          </h2>
        </div>
        {editions.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            No editions available for this work yet.
          </p>
        ) : (
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            {editions.map((edition) => (
              <EditionRow key={edition.uid} edition={edition} />
            ))}
          </div>
        )}
      </div>

      <Separator />
      <BookReviews bookUid={uid} />
    </div>
  );
}

function EditionRow({ edition }: { edition: Edition }) {
  const addToCart = useAddToCart();
  const toast = useToast();
  const price = syntheticPrice(edition.uid);

  const handleAddToCart = () => {
    addToCart.mutate(
      {
        product_id: numericIdFromUid(edition.uid),
        product_name: edition.title,
        quantity: 1,
        currency: "USD",
        price,
      },
      {
        onSuccess: () =>
          toast.success("Added to cart", { description: edition.title }),
        onError: () => toast.error("Couldn't add to cart"),
      },
    );
  };

  return (
    <Card>
      <CardContent className="space-y-3 py-4">
        <div className="space-y-1">
          <p className="font-medium leading-snug">{edition.title}</p>
          <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
            {edition.physical_format && (
              <Badge variant="outline" className="text-xs">
                {edition.physical_format}
              </Badge>
            )}
            {edition.publisher && <span>{edition.publisher.name}</span>}
            {edition.publish_date && <span>{edition.publish_date}</span>}
          </div>
        </div>

        <div className="grid grid-cols-1 gap-1 text-xs text-muted-foreground">
          {edition.isbn13 && (
            <span className="inline-flex items-center gap-1.5">
              <Hash className="size-3" />
              ISBN-13 {edition.isbn13}
            </span>
          )}
          {edition.isbn10 && (
            <span className="inline-flex items-center gap-1.5">
              <Hash className="size-3" />
              ISBN-10 {edition.isbn10}
            </span>
          )}
          {edition.number_of_pages != null && (
            <span>{edition.number_of_pages} pages</span>
          )}
          {edition.languages && edition.languages.length > 0 && (
            <span>{edition.languages.join(", ")}</span>
          )}
        </div>

        <div className="flex items-center justify-between pt-1">
          <span className="text-lg font-bold">${price.toFixed(2)}</span>
          <Button
            size="sm"
            disabled={addToCart.isPending}
            onClick={handleAddToCart}
          >
            <ShoppingCart className="size-4 mr-1.5" />
            Add to Cart
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
