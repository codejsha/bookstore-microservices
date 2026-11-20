import { Separator } from "@bookstore/design/ui/separator";
import type { QueryClient } from "@tanstack/react-query";
import { useQuery } from "@tanstack/react-query";
import {
  createRootRouteWithContext,
  Link,
  Outlet,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import {
  BookOpen,
  Heart,
  LayoutGrid,
  Loader2,
  LogIn,
  Package,
  ShoppingCart,
  User,
  UserPlus,
} from "lucide-react";
import { Suspense } from "react";
import { useTranslation } from "react-i18next";
import { useAuth } from "react-oidc-context";
import { cartQueryOptions } from "@/domains/cart";
import {
  HoverStripSidebar,
  StripRow,
} from "@/shared/components/HoverStripSidebar";
import { LanguageSwitcher } from "@/shared/components/LanguageSwitcher";
import { NotFoundState } from "@/shared/components/NotFoundState";
import { RouteError } from "@/shared/components/RouteError";
import { ThemeToggle } from "@/shared/components/ThemeToggle";
import { ToastProvider } from "@/shared/components/toast";
import { useLocale } from "@/shared/i18n/useLocale";

interface RouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
  errorComponent: (props) => (
    <RootShell>
      <RouteError {...props} />
    </RootShell>
  ),
  notFoundComponent: NotFoundState,
});

function RootShell({ children }: { children: React.ReactNode }) {
  return (
    <ToastProvider>
      <div className="flex min-h-screen flex-col">
        <Header />
        <div className="flex flex-1">
          <LeftSidebar />
          <main className="min-w-0 flex-1">{children}</main>
          <RightSidebar />
        </div>
        <Footer />
      </div>
    </ToastProvider>
  );
}

function RootComponent() {
  return (
    <RootShell>
      <Suspense
        fallback={
          <div className="flex items-center justify-center py-20">
            <Loader2 className="size-8 animate-spin text-muted-foreground" />
          </div>
        }
      >
        <Outlet />
      </Suspense>
      <TanStackRouterDevtools />
    </RootShell>
  );
}

const navLinkClass =
  "flex items-center rounded-md px-2 mx-2 py-1.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground [&.active]:text-foreground [&.active]:bg-muted";

function LeftSidebar() {
  const { t } = useTranslation();
  const locale = useLocale();

  return (
    <HoverStripSidebar side="left" label={t("nav.browse")}>
      <Link to="/$locale/books" params={{ locale }} className={navLinkClass}>
        <StripRow
          icon={<LayoutGrid className="size-5" />}
          label={t("nav.books")}
        />
      </Link>
      <Link to="/$locale/orders" params={{ locale }} className={navLinkClass}>
        <StripRow
          icon={<Package className="size-5" />}
          label={t("nav.orders")}
        />
      </Link>
      <Link to="/$locale/wishlist" params={{ locale }} className={navLinkClass}>
        <StripRow
          icon={<Heart className="size-5" />}
          label={t("nav.wishlist")}
        />
      </Link>
    </HoverStripSidebar>
  );
}

function RightSidebar() {
  const { t } = useTranslation();
  const locale = useLocale();
  const auth = useAuth();
  const { data: cart } = useQuery({
    ...cartQueryOptions(),
    enabled: auth.isAuthenticated,
    placeholderData: {
      uid: "",
      user_uid: "",
      items: [],
      total_items: 0,
      total_amount: 0,
    },
  });
  const cartCount = cart?.total_items ?? 0;

  return (
    <HoverStripSidebar side="right" label={t("nav.account")}>
      <Link
        to="/$locale/cart"
        params={{ locale }}
        className={navLinkClass}
        aria-label={
          cartCount > 0
            ? t("nav.cartWithCount", { count: cartCount })
            : t("nav.cart")
        }
      >
        <StripRow
          icon={<ShoppingCart className="size-5" />}
          label={t("nav.cart")}
          badge={
            cartCount > 0 ? (
              <span className="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-primary-foreground">
                {cartCount}
              </span>
            ) : null
          }
        />
      </Link>
      <Link to="/$locale/account" params={{ locale }} className={navLinkClass}>
        <StripRow icon={<User className="size-5" />} label={t("nav.account")} />
      </Link>
      <Link to="/$locale/login" params={{ locale }} className={navLinkClass}>
        <StripRow icon={<LogIn className="size-5" />} label={t("nav.login")} />
      </Link>
      <Link to="/$locale/signup" params={{ locale }} className={navLinkClass}>
        <StripRow
          icon={<UserPlus className="size-5" />}
          label={t("nav.signUp")}
        />
      </Link>
    </HoverStripSidebar>
  );
}

function Header() {
  const locale = useLocale();

  return (
    <header className="sticky top-0 z-50 border-b bg-muted">
      <div className="flex h-14 items-center justify-between px-4">
        <Link
          to="/$locale"
          params={{ locale }}
          className="flex items-center gap-2 text-lg font-bold"
        >
          <BookOpen className="size-5" />
          Bookstore
        </Link>
        <div className="flex items-center gap-1">
          <LanguageSwitcher />
          <ThemeToggle />
        </div>
      </div>
    </header>
  );
}

function Footer() {
  const { t } = useTranslation();
  const locale = useLocale();

  return (
    <footer className="border-t">
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 gap-8 sm:grid-cols-3">
          <div className="space-y-2">
            <h3 className="flex items-center gap-2 font-semibold">
              <BookOpen className="size-4" />
              Bookstore
            </h3>
            <p className="text-sm text-muted-foreground">
              {t("footer.tagline")}
            </p>
          </div>
          <div className="space-y-2">
            <h4 className="text-sm font-semibold">{t("nav.browse")}</h4>
            <nav aria-label={t("nav.browse")} className="flex flex-col gap-1">
              <Link
                to="/$locale/books"
                params={{ locale }}
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                {t("nav.books")}
              </Link>
              <Link
                to="/$locale/wishlist"
                params={{ locale }}
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                {t("nav.wishlist")}
              </Link>
              <Link
                to="/$locale/cart"
                params={{ locale }}
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                {t("nav.cart")}
              </Link>
            </nav>
          </div>
          <div className="space-y-2">
            <h4 className="text-sm font-semibold">{t("nav.account")}</h4>
            <nav aria-label={t("nav.account")} className="flex flex-col gap-1">
              <Link
                to="/$locale/account"
                params={{ locale }}
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                {t("nav.myAccount")}
              </Link>
              <Link
                to="/$locale/orders"
                params={{ locale }}
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                {t("nav.orders")}
              </Link>
              <Link
                to="/$locale/login"
                params={{ locale }}
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                {t("nav.login")}
              </Link>
              <Link
                to="/$locale/signup"
                params={{ locale }}
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                {t("nav.signUp")}
              </Link>
            </nav>
          </div>
        </div>
        <Separator className="my-6" />
        <p className="text-center text-sm text-muted-foreground">
          {t("footer.copyright", { year: new Date().getFullYear() })}
        </p>
      </div>
    </footer>
  );
}
