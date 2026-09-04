import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { Suspense } from "react";
import { ErrorState } from "@/shared/components/ErrorState";

export interface RouterContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootComponent,
  errorComponent: ({ reset }) => (
    <ErrorState
      title="Unable to load this page"
      message="An unexpected error occurred. Please try again."
      onRetry={reset}
    />
  ),
  notFoundComponent: () => (
    <ErrorState
      title="Page not found"
      message="That page does not exist in the admin console."
    />
  ),
});

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
