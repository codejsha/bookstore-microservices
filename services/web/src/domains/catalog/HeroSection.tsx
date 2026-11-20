import { Button } from "@bookstore/design/ui/button";
import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import heroImg from "@/assets/hero.png";
import { useLocale } from "@/shared/i18n/useLocale";
import { WorkCard } from "./work/WorkCard";
import { WorkCardSkeleton } from "./work/WorkCardSkeleton";
import { workListQueryOptions } from "./work/work-queries";

const TRENDING_SIZE = 4;

export function HeroSection() {
  const { t } = useTranslation("catalog");
  const locale = useLocale();
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
              {t("hero.title")}
            </h1>
            <p className="mx-auto max-w-lg text-lg text-muted-foreground lg:mx-0">
              {t("hero.subtitle")}
            </p>
            <div className="flex items-center justify-center gap-3 lg:justify-start">
              <Button
                size="lg"
                render={<Link to="/$locale/books" params={{ locale }} />}
              >
                {t("hero.browseBooks")}
              </Button>
              <Button
                variant="outline"
                size="lg"
                render={<Link to="/$locale/signup" params={{ locale }} />}
              >
                {t("hero.createAccount")}
              </Button>
            </div>
          </div>
          <div className="flex-1 flex justify-center">
            <img
              src={heroImg}
              alt={t("hero.imageAlt")}
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
                  {t("trending.title")}
                </h2>
                <p className="text-sm text-muted-foreground">
                  {t("trending.subtitle")}
                </p>
              </div>
              <Button
                variant="ghost"
                size="sm"
                render={<Link to="/$locale/books" params={{ locale }} />}
              >
                {t("trending.viewAll")}
              </Button>
            </div>
            <ul
              aria-label={t("trending.listLabel")}
              className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
            >
              {isLoading
                ? Array.from({ length: TRENDING_SIZE }).map((_, i) => (
                    // biome-ignore lint/suspicious/noArrayIndexKey: skeleton placeholders
                    <li key={i}>
                      <WorkCardSkeleton />
                    </li>
                  ))
                : trending.map((work) => (
                    <li key={work.uid}>
                      <WorkCard work={work} />
                    </li>
                  ))}
            </ul>
          </div>
        </section>
      )}
    </div>
  );
}
