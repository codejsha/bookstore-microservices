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
import { useAuth } from "react-oidc-context";
import { cartQueryOptions } from "@/domains/cart";
import { ErrorState } from "@/shared/components/ErrorState";
import {
  HoverStripSidebar,
  StripRow,
} from "@/shared/components/HoverStripSidebar";
import { ThemeToggle } from "@/shared/components/ThemeToggle";
import { ToastProvider } from "@/shared/components/toast";

interface RouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
  errorComponent: ({ reset }) => (
    <RootShell>
      <div className="container mx-auto px-4">
        <ErrorState
          title="Unable to load this page"
          message="An unexpected error occurred. Please try again or head back home."
          onRetry={reset}
        />
      </div>
    </RootShell>
  ),
  notFoundComponent: () => (
    <RootShell>
      <div className="container mx-auto px-4">
        <ErrorState
          title="Page not found"
          message="The page you're looking for doesn't exist or has moved."
        />
      </div>
    </RootShell>
  ),
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
  return (
    <HoverStripSidebar side="left">
      <Link to="/books" className={navLinkClass}>
        <StripRow icon={<LayoutGrid className="size-5" />} label="Books" />
      </Link>
      <Link to="/orders" className={navLinkClass}>
        <StripRow icon={<Package className="size-5" />} label="Orders" />
      </Link>
      <Link to="/wishlist" className={navLinkClass}>
        <StripRow icon={<Heart className="size-5" />} label="Wishlist" />
      </Link>
    </HoverStripSidebar>
  );
}

function RightSidebar() {
  const auth = useAuth();
  const { data: cart } = useQuery({
    ...cartQueryOptions(),
    enabled: auth.isAuthenticated,
    placeholderData: {
      uid: "",
      user_id: "0",
      items: [],
      total_items: 0,
      total_amount: 0,
    },
  });
  const cartCount = cart?.total_items ?? 0;

  return (
    <HoverStripSidebar side="right">
      <Link to="/cart" className={navLinkClass}>
        <StripRow
          icon={<ShoppingCart className="size-5" />}
          label="Cart"
          badge={
            cartCount > 0 ? (
              <span className="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-primary text-[10px] font-bold text-primary-foreground">
                {cartCount}
              </span>
            ) : null
          }
        />
      </Link>
      <Link to="/account" className={navLinkClass}>
        <StripRow icon={<User className="size-5" />} label="Account" />
      </Link>
      <Link to="/login" className={navLinkClass}>
        <StripRow icon={<LogIn className="size-5" />} label="Login" />
      </Link>
      <Link to="/signup" className={navLinkClass}>
        <StripRow icon={<UserPlus className="size-5" />} label="Sign up" />
      </Link>
    </HoverStripSidebar>
  );
}

function Header() {
  return (
    <header className="sticky top-0 z-50 border-b bg-muted">
      <div className="flex h-14 items-center justify-between px-4">
        <Link to="/" className="flex items-center gap-2 text-lg font-bold">
          <BookOpen className="size-5" />
          Bookstore
        </Link>
        <ThemeToggle />
      </div>
    </header>
  );
}

function Footer() {
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
              Discover, browse, and find your next great read.
            </p>
          </div>
          <div className="space-y-2">
            <h4 className="text-sm font-semibold">Browse</h4>
            <nav className="flex flex-col gap-1">
              <Link
                to="/books"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Books
              </Link>
              <Link
                to="/wishlist"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Wishlist
              </Link>
              <Link
                to="/cart"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Cart
              </Link>
            </nav>
          </div>
          <div className="space-y-2">
            <h4 className="text-sm font-semibold">Account</h4>
            <nav className="flex flex-col gap-1">
              <Link
                to="/account"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                My account
              </Link>
              <Link
                to="/orders"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Orders
              </Link>
              <Link
                to="/login"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Login
              </Link>
              <Link
                to="/signup"
                className="text-sm text-muted-foreground hover:text-foreground"
              >
                Sign up
              </Link>
            </nav>
          </div>
        </div>
        <Separator className="my-6" />
        <p className="text-center text-sm text-muted-foreground">
          &copy; {new Date().getFullYear()} Bookstore. All rights reserved.
        </p>
      </div>
    </footer>
  );
}
