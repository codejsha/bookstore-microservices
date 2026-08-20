import { QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { AuthProvider } from "react-oidc-context";
import "./index.css";
import { queryClient } from "@/shared/api/queryClient";
import { userManager } from "@/shared/auth/userManager";
import { routeTree } from "./routeTree.gen";

const router = createRouter({
  routeTree,
  context: { queryClient },
  defaultPreloadStaleTime: 0,
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

async function bootstrap() {
  if (import.meta.env.DEV) {
    const { worker } = await import("./mocks/browser");
    await worker.start({ onUnhandledRequest: "bypass" });
  }

  try {
    const { initTelemetry, instrumentRouter, initWebVitals } = await import(
      "@/shared/lib/telemetry"
    );
    initTelemetry();
    instrumentRouter(router);
    initWebVitals();
  } catch (e) {
    console.warn("[OTel] Failed to initialize telemetry:", e);
  }

  try {
    const { initFaro } = await import("@/shared/lib/faro");
    initFaro();
  } catch (e) {
    console.warn("[Faro] Failed to initialize Faro:", e);
  }

  const rootElement = document.getElementById("root");
  if (rootElement && !rootElement.innerHTML) {
    createRoot(rootElement).render(
      <StrictMode>
        <AuthProvider
          userManager={userManager}
          onSigninCallback={() => {
            window.history.replaceState(
              {},
              document.title,
              window.location.pathname,
            );
          }}
        >
          <QueryClientProvider client={queryClient}>
            <RouterProvider router={router} />
            <ReactQueryDevtools initialIsOpen={false} />
          </QueryClientProvider>
        </AuthProvider>
      </StrictMode>,
    );
  }
}

bootstrap();
