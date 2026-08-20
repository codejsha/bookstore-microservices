import { Button } from "@bookstore/design/ui/button";
import { Card, CardContent } from "@bookstore/design/ui/card";
import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { BookOpen, Heart, Loader2, Trash2 } from "lucide-react";
import { workListQueryOptions } from "@/domains/catalog";
import { useToast } from "@/shared/components/toast";
import { useCurrentUserUid } from "@/shared/lib/currentUser";
import { useRemoveFromWishlist } from "./wishlist-mutations";
import { wishlistQueryOptions } from "./wishlist-queries";

export function WishlistPage() {
  const currentUserUid = useCurrentUserUid();
  const { data: wishlist, isLoading } = useQuery(
    wishlistQueryOptions(currentUserUid),
  );
  const { data: works } = useQuery(workListQueryOptions({ size: 100 }));
  const remove = useRemoveFromWishlist(currentUserUid);
  const toast = useToast();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-20">
        <Loader2 className="size-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const bookUids = wishlist?.book_uids ?? [];
  const wishlistWorks = (works?.items ?? []).filter((w) =>
    bookUids.includes(w.uid),
  );
  const isEmpty = bookUids.length === 0;

  return (
    <div className="container mx-auto max-w-4xl space-y-8 px-4 py-8">
      <div className="flex items-center gap-3">
        <Heart className="size-7" />
        <h1 className="text-3xl font-bold tracking-tight">My Wishlist</h1>
        {!isEmpty && (
          <span className="text-lg text-muted-foreground">
            ({bookUids.length} {bookUids.length === 1 ? "book" : "books"})
          </span>
        )}
      </div>

      {isEmpty ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center space-y-4 py-16">
            <Heart className="size-16 text-muted-foreground" />
            <p className="text-xl font-medium text-muted-foreground">
              Your wishlist is empty
            </p>
            <Button render={<Link to="/books" />}>Browse Books</Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {wishlistWorks.map((work) => (
            <Card key={work.uid}>
              <CardContent className="flex items-center gap-4 py-4">
                <div className="flex h-20 w-14 shrink-0 items-center justify-center rounded-md bg-muted">
                  <BookOpen className="size-6 text-muted-foreground" />
                </div>
                <div className="min-w-0 flex-1">
                  <Link
                    to="/books/$uid"
                    params={{ uid: work.uid }}
                    className="line-clamp-2 font-medium hover:underline"
                  >
                    {work.title}
                  </Link>
                  <p className="mt-1 truncate text-sm text-muted-foreground">
                    {work.authors?.map((a) => a.name).join(", ")}
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="size-8 text-destructive hover:text-destructive"
                  disabled={remove.isPending}
                  onClick={() =>
                    remove.mutate(
                      { book_uids: [work.uid] },
                      {
                        onSuccess: () => toast.success("Removed from wishlist"),
                        onError: () => toast.error("Couldn't update wishlist"),
                      },
                    )
                  }
                >
                  <Trash2 className="size-4" />
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
