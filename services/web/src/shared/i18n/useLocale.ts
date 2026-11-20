import { useTranslation } from "react-i18next";
import { type Locale, toLocale } from "./locales";

export function useLocale(): Locale {
  const { i18n } = useTranslation();

  return toLocale(i18n.resolvedLanguage);
}
