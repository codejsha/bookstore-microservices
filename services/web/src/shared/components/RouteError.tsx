import type { ErrorComponentProps } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { ErrorState } from "./ErrorState";

export function RouteError({ reset }: ErrorComponentProps) {
  const { t } = useTranslation();

  return (
    <div className="container mx-auto px-4">
      <ErrorState
        title={t("error.routeTitle")}
        message={t("error.routeMessage")}
        onRetry={reset}
      />
    </div>
  );
}
