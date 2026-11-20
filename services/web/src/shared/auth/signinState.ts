import { type Locale, toLocale } from "@/shared/i18n/locales";

export interface SigninState {
  locale: Locale;
}

export function signinState(locale: Locale): SigninState {
  return { locale };
}

export function localeFromSigninState(state: unknown): Locale {
  if (state && typeof state === "object" && "locale" in state) {
    return toLocale(state.locale);
  }
  return toLocale(undefined);
}
