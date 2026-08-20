import { Button } from "@bookstore/design/ui/button";
import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import heroImg from "@/assets/hero.png";
import { WorkCard } from "./work/WorkCard";
import { WorkCardSkeleton } from "./work/WorkCardSkeleton";
import { workListQueryOptions } from "./work/work-queries";

const TRENDING_SIZE = 4;

export function HeroSection() {
  const { data, isLoading, isError } = useQuery(
    workListQueryOptions({ size: TRENDING_SIZE, sort: "createdAt,desc" }),
  );
  const trending = data?.items ?? [];

  return (
    <div className="space-y-16">
      <section
        className="relative"
        style={{ width: "100vw", marginLeft: "calc(50% - 50vw)" }}
      >
        <div className="absolute inset-0 -z-10 bg-gradient-to-b from-primary/5 via-background to-background" />
        <div className="container mx-auto flex flex-col-reverse items-center gap-8 px-4 py-20 lg:flex-row lg:py-28">
          <div className="flex-1 space-y-6 text-center lg:text-left">
            <h1 className="text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
              Discover Your Next Great Read
            </h1>
            <p className="mx-auto max-w-lg text-lg text-muted-foreground lg:mx-0">
              Browse our curated catalog of books across every genre. From
              timeless classics to the latest releases, find your perfect match.
            </p>
            <div className="flex items-center justify-center gap-3 lg:justify-start">
              <Button size="lg" render={<Link to="/books" />}>
                Browse Books
              </Button>
              <Button
                variant="outline"
                size="lg"
                render={<Link to="/signup" />}
              >
                Create Account
              </Button>
            </div>
          </div>
          <div className="flex-1 flex justify-center">
            <img
              src={heroImg}
              alt="Books illustration"
              width={448}
              height={336}
              className="w-full max-w-md rounded-xl object-cover shadow-lg aspect-[4/3]"
            />
          </div>
        </div>
      </section>

      {!isError && (
        <section className="container mx-auto px-4 pb-16">
          <div className="space-y-6">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-2xl font-bold tracking-tight">
                  Trending Now
                </h2>
                <p className="text-sm text-muted-foreground">
                  Popular picks from our catalog
                </p>
              </div>
              <Button variant="ghost" size="sm" render={<Link to="/books" />}>
                View all
              </Button>
            </div>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
              {isLoading
                ? Array.from({ length: TRENDING_SIZE }).map((_, i) => (
                    // biome-ignore lint/suspicious/noArrayIndexKey: skeleton placeholders
                    <WorkCardSkeleton key={i} />
                  ))
                : trending.map((work) => (
                    <WorkCard key={work.uid} work={work} />
                  ))}
            </div>
          </div>
        </section>
      )}
    </div>
  );
}
