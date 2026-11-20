import { Button } from "@bookstore/design/ui/button";
import { Link } from "@tanstack/react-router";
import { AlertTriangle, Home, RefreshCw } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useLocale } from "@/shared/i18n/useLocale";

interface ErrorStateProps {
  title?: string;
  message?: string;
  onRetry?: () => void;
  showHome?: boolean;
}

export function ErrorState({
  title,
  message,
  onRetry,
  showHome = true,
}: ErrorStateProps) {
  const { t } = useTranslation();
  const locale = useLocale();

  return (
    <div className="flex flex-col items-center justify-center py-20 text-center animate-in fade-in duration-300">
      <div className="mb-4 flex size-12 items-center justify-center rounded-full bg-destructive/10">
        <AlertTriangle className="size-6 text-destructive" />
      </div>
      <p className="text-lg font-medium">{title ?? t("error.defaultTitle")}</p>
      <p className="mt-1 max-w-md text-sm text-muted-foreground">
        {message ?? t("error.defaultMessage")}
      </p>
      <div className="mt-6 flex items-center gap-2">
        {onRetry && (
          <Button onClick={onRetry}>
            <RefreshCw className="size-4 mr-1.5" />
            {t("error.retry")}
          </Button>
        )}
        {showHome && (
          <Button
            variant="outline"
            render={<Link to="/$locale" params={{ locale }} />}
          >
            <Home className="size-4 mr-1.5" />
            {t("error.goHome")}
          </Button>
        )}
      </div>
    </div>
  );
}
