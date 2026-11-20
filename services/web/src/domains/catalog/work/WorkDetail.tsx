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
import { useTranslation } from "react-i18next";
import { useAddToCart } from "@/domains/cart";
import { BookReviews } from "@/domains/review";
import {
  useAddToWishlist,
  useRemoveFromWishlist,
  wishlistQueryOptions,
} from "@/domains/wishlist";
import { useToast } from "@/shared/components/toast";
import { useFormat } from "@/shared/hooks/useFormat";
import { useLocale } from "@/shared/i18n/useLocale";
import { useCurrentUserUid } from "@/shared/lib/currentUser";
import { numericIdFromUid, syntheticPrice } from "../catalog-utils";
import type { Edition } from "../types";
import {
  workDetailQueryOptions,
  workEditionsQueryOptions,
} from "./work-queries";

const EDITION_CURRENCY = "USD";

interface WorkDetailProps {
  uid: string;
}

export function WorkDetail({ uid }: WorkDetailProps) {
  const { t } = useTranslation(["catalog", "wishlist"]);
  const locale = useLocale();
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
          onSuccess: () => toast.success(t("wishlist:toast.removed")),
          onError: () => toast.error(t("wishlist:toast.updateFailed")),
        },
      );
    } else {
      addToWishlist.mutate(
        { book_uids: [uid] },
        {
          onSuccess: () => toast.success(t("wishlist:toast.added")),
          onError: () => toast.error(t("wishlist:toast.updateFailed")),
        },
      );
    }
  };

  const year = work.first_publish_date
    ? work.first_publish_date.slice(0, 4)
    : undefined;

  return (
    <div className="container mx-auto max-w-5xl px-4 py-8 space-y-8">
      <Button
        variant="ghost"
        size="sm"
        render={<Link to="/$locale/books" params={{ locale }} />}
      >
        <ArrowLeft className="size-4 mr-1" />
        {t("detail.backToBooks")}
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
              {t("detail.byAuthors", {
                authors: work.authors.map((a) => a.name).join(", "),
              })}
            </p>

            {year && (
              <p className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <Calendar className="size-3.5" />
                {t("detail.firstPublished", { year })}
              </p>
            )}
          </div>

          <div className="flex items-center gap-4">
            <Button
              variant="outline"
              size="icon"
              aria-label={
                isWished
                  ? t("wishlist:actions.remove")
                  : t("wishlist:actions.add")
              }
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
              <h2 className="text-lg font-semibold">
                {t("detail.description")}
              </h2>
              <p className="text-muted-foreground leading-relaxed whitespace-pre-line">
                {work.description}
              </p>
            </div>
          )}

          {work.authors.some((a) => a.bio) && (
            <div className="space-y-3">
              <h2 className="text-lg font-semibold">
                {t("detail.aboutAuthor", { count: work.authors.length })}
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
            {editions.length > 0
              ? t("editions.titleWithTotal", { total: editions.length })
              : t("editions.title")}
          </h2>
        </div>
        {editions.length === 0 ? (
          <p className="text-sm text-muted-foreground">{t("editions.empty")}</p>
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
  const { t } = useTranslation("catalog");
  const format = useFormat();
  const addToCart = useAddToCart();
  const toast = useToast();
  const price = syntheticPrice(edition.uid);

  const handleAddToCart = () => {
    addToCart.mutate(
      {
        product_id: String(numericIdFromUid(edition.uid)),
        product_name: edition.title,
        quantity: 1,
        currency: EDITION_CURRENCY,
        price,
      },
      {
        onSuccess: () =>
          toast.success(t("toast.addedToCart"), { description: edition.title }),
        onError: () => toast.error(t("toast.addToCartFailed")),
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
              {t("editions.isbn13", { isbn: edition.isbn13 })}
            </span>
          )}
          {edition.isbn10 && (
            <span className="inline-flex items-center gap-1.5">
              <Hash className="size-3" />
              {t("editions.isbn10", { isbn: edition.isbn10 })}
            </span>
          )}
          {edition.number_of_pages != null && (
            <span>
              {t("editions.pages", { pages: edition.number_of_pages })}
            </span>
          )}
          {edition.languages && edition.languages.length > 0 && (
            <span>{edition.languages.join(", ")}</span>
          )}
        </div>

        <div className="flex items-center justify-between pt-1">
          <span className="text-lg font-bold">
            {format.price(price, EDITION_CURRENCY)}
          </span>
          <Button
            size="sm"
            disabled={addToCart.isPending}
            onClick={handleAddToCart}
          >
            <ShoppingCart className="size-4 mr-1.5" />
            {t("editions.addToCart")}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
