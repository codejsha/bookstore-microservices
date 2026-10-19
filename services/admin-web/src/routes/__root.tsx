import type { QueryClient } from "@tanstack/react-query";
import {
  createRootRouteWithContext,
  Outlet,
  useRouter,
} from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { Suspense } from "react";
import { ErrorState } from "@/shared/components/ErrorState";

export interface RouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
  errorComponent: RootError,
  notFoundComponent: () => (
    <ErrorState
      title="Page not found"
      message="That page does not exist in the admin console."
    />
  ),
});

function RootError() {
  const router = useRouter();
  return (
    <ErrorState
      title="Unable to load this page"
      message="An unexpected error occurred. Please try again."
      onRetry={() => router.invalidate()}
    />
  );
}

function RootComponent() {
  return (
    <Suspense fallback={null}>
      <Outlet />
      {import.meta.env.DEV ? (
        <TanStackRouterDevtools position="bottom-right" />
      ) : null}
    </Suspense>
  );
}
