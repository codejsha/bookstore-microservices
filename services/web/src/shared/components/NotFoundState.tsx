import { ErrorState } from "./ErrorState";

export function NotFoundState() {
  return (
    <div className="container mx-auto px-4">
      <ErrorState
        title="Page not found"
        message="The page you're looking for doesn't exist or has moved."
      />
    </div>
  );
}
