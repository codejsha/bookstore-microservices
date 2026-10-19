import { Button } from "@bookstore/design/ui/button";
import { Input } from "@bookstore/design/ui/input";
import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import {
  ArrowDownAZ,
  ArrowUpAZ,
  Grid3X3,
  List,
  Loader2,
  Search,
} from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";
import { ErrorState } from "@/shared/components/ErrorState";
import { subjectListQueryOptions } from "../subject/subject-queries";
import { WorkCard } from "./WorkCard";
import { WorkCardSkeleton } from "./WorkCardSkeleton";
import { WorkListItem } from "./WorkListItem";
import { workListInfiniteQueryOptions } from "./work-queries";

const SORT_OPTIONS = [
  { value: "title,asc", label: "Title A-Z", icon: ArrowDownAZ },
  { value: "title,desc", label: "Title Z-A", icon: ArrowUpAZ },
  { value: "createdAt,desc", label: "Newest" },
  { value: "firstPublishDate,desc", label: "Recently published" },
] as const;

const PAGE_SIZE = 12;

type ViewMode = "card" | "list";

export function WorkList() {
  const [title, setTitle] = useState("");
  const [debouncedTitle, setDebouncedTitle] = useState("");
  const [subjectUid, setSubjectUid] = useState<string | undefined>();
  const [viewMode, setViewMode] = useState<ViewMode>("card");
  const [sort, setSort] = useState<string>("title,asc");

  const { data: subjectData } = useQuery(subjectListQueryOptions({ size: 12 }));
  const subjects = subjectData?.items ?? [];

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    isFetching,
    isError,
    refetch,
  } = useInfiniteQuery(
    workListInfiniteQueryOptions({
      title: debouncedTitle || undefined,
      subjectUid,
      sort,
      size: PAGE_SIZE,
    }),
  );

  const sentinelRef = useRef<HTMLDivElement>(null);
  const searchTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(
    () => () => {
      if (searchTimer.current) clearTimeout(searchTimer.current);
    },
    [],
  );

  const handleIntersect = useCallback(
    (entries: IntersectionObserverEntry[]) => {
      if (entries[0]?.isIntersecting && hasNextPage && !isFetchingNextPage) {
        fetchNextPage();
      }
    },
    [fetchNextPage, hasNextPage, isFetchingNextPage],
  );

  useEffect(() => {
    const el = sentinelRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(handleIntersect, {
      rootMargin: "200px",
    });
    observer.observe(el);
    return () => observer.disconnect();
  }, [handleIntersect]);

  function handleSearch(value: string) {
    setTitle(value);
    if (searchTimer.current) clearTimeout(searchTimer.current);
    searchTimer.current = setTimeout(() => setDebouncedTitle(value), 400);
  }

  const allWorks = data?.pages.flatMap((p) => p.items) ?? [];
  const total = data?.pages[0]?.total ?? 0;

  return (
    <div className="container mx-auto space-y-6 px-4 py-8">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Books</h1>
        <p className="text-sm text-muted-foreground min-h-5">
          {total > 0 && `${total} book${total !== 1 ? "s" : ""} found`}
        </p>
      </div>

      <div className="space-y-3">
        <div className="flex flex-wrap gap-1.5">
          {[{ uid: undefined, name: "All" }, ...subjects].map((subject) => (
            <button
              key={subject.uid ?? "all"}
              type="button"
              aria-pressed={subjectUid === subject.uid}
              onClick={() => setSubjectUid(subject.uid)}
              className={`inline-flex h-8 items-center rounded-md border px-3 text-sm font-medium transition-colors ${
                subjectUid === subject.uid
                  ? "border-primary bg-primary text-primary-foreground"
                  : "border-transparent bg-muted text-muted-foreground hover:text-foreground"
              }`}
            >
              {subject.name}
            </button>
          ))}
        </div>
        <div className="flex items-center justify-between gap-2">
          <div className="relative w-80">
            <Search className="absolute left-2.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              aria-label="Search by title"
              placeholder="Search by title..."
              value={title}
              onChange={(e) => handleSearch(e.target.value)}
              className="pl-8 h-8"
            />
          </div>
          <div className="flex items-center gap-2">
            <select
              aria-label="Sort by"
              value={sort}
              onChange={(e) => setSort(e.target.value)}
              className="h-8 rounded-md border bg-background px-2 text-sm text-foreground outline-none"
            >
              {SORT_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
            <div className="flex items-center gap-1 rounded-lg border p-0.5">
              <Button
                variant={viewMode === "card" ? "secondary" : "ghost"}
                size="icon-sm"
                aria-label="Grid view"
                aria-pressed={viewMode === "card"}
                onClick={() => setViewMode("card")}
              >
                <Grid3X3 className="size-4" />
              </Button>
              <Button
                variant={viewMode === "list" ? "secondary" : "ghost"}
                size="icon-sm"
                aria-label="List view"
                aria-pressed={viewMode === "list"}
                onClick={() => setViewMode("list")}
              >
                <List className="size-4" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="relative min-h-[300px]">
        <div
          className={`grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 transition-opacity duration-300 ${isLoading ? "opacity-100" : "opacity-0 absolute inset-0 pointer-events-none"}`}
        >
          {Array.from({ length: PAGE_SIZE }).map((_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: skeleton placeholders
            <WorkCardSkeleton key={i} />
          ))}
        </div>

        {!isLoading && isError && (
          <ErrorState
            title="Couldn't load books"
            message="There was a problem loading the catalog. Please try again."
            onRetry={() => refetch()}
            showHome={false}
          />
        )}

        {!isLoading && !isError && (
          <div
            className={`transition-opacity duration-200 ${isFetching && !isFetchingNextPage ? "opacity-40" : "opacity-100"}`}
          >
            {allWorks.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-20 text-center animate-in fade-in duration-300">
                <p className="text-lg font-medium">No books found</p>
                <p className="mt-1 text-sm text-muted-foreground">
                  Try adjusting your search or filters.
                </p>
              </div>
            ) : viewMode === "card" ? (
              <ul
                aria-label="Books"
                className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 animate-in fade-in duration-300"
              >
                {allWorks.map((work) => (
                  <li key={work.uid}>
                    <WorkCard work={work} />
                  </li>
                ))}
              </ul>
            ) : (
              <ul
                aria-label="Books"
                className="space-y-2 animate-in fade-in duration-300"
              >
                {allWorks.map((work) => (
                  <li key={work.uid}>
                    <WorkListItem work={work} />
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}
      </div>

      <div ref={sentinelRef} className="flex justify-center py-4">
        {isFetchingNextPage && (
          <Loader2 className="size-6 animate-spin text-muted-foreground" />
        )}
        {!hasNextPage && allWorks.length > 0 && (
          <p className="text-sm text-muted-foreground">
            All {total} books loaded
          </p>
        )}
      </div>
    </div>
  );
}
