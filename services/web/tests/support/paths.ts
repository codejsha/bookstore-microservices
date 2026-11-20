import { DEFAULT_LOCALE, type Locale } from "../../src/shared/i18n/locales";

export function localePath(
  path: string,
  locale: Locale = DEFAULT_LOCALE,
): string {
  return path === "/" ? `/${locale}` : `/${locale}${path}`;
}
