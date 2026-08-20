import { Button } from "@bookstore/design/ui/button";
import { Link } from "@tanstack/react-router";
import { AlertTriangle, Home, RefreshCw } from "lucide-react";

interface ErrorStateProps {
  title?: string;
  message?: string;
  onRetry?: () => void;
  showHome?: boolean;
}

export function ErrorState({
  title = "Something went wrong",
  message = "We couldn't load this content. Please try again.",
  onRetry,
  showHome = true,
}: ErrorStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-20 text-center animate-in fade-in duration-300">
      <div className="mb-4 flex size-12 items-center justify-center rounded-full bg-destructive/10">
        <AlertTriangle className="size-6 text-destructive" />
      </div>
      <p className="text-lg font-medium">{title}</p>
      <p className="mt-1 max-w-md text-sm text-muted-foreground">{message}</p>
      <div className="mt-6 flex items-center gap-2">
        {onRetry && (
          <Button onClick={onRetry}>
            <RefreshCw className="size-4 mr-1.5" />
            Try again
          </Button>
        )}
        {showHome && (
          <Button variant="outline" render={<Link to="/" />}>
            <Home className="size-4 mr-1.5" />
            Go home
          </Button>
        )}
      </div>
    </div>
  );
}
