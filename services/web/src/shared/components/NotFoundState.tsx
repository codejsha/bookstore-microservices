import { useTranslation } from "react-i18next";
import { ErrorState } from "./ErrorState";

export function NotFoundState() {
  const { t } = useTranslation();

  return (
    <div className="container mx-auto px-4">
      <ErrorState
        title={t("error.notFoundTitle")}
        message={t("error.notFoundMessage")}
      />
    </div>
  );
}
