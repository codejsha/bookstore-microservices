import type { ErrorComponentProps } from "@tanstack/react-router";
import { ErrorState } from "./ErrorState";

export function RouteError({ reset }: ErrorComponentProps) {
  return (
    <div className="container mx-auto px-4">
      <ErrorState
        title="Unable to load this page"
        message="An unexpected error occurred. Please try again or head back home."
        onRetry={reset}
      />
    </div>
  );
}
