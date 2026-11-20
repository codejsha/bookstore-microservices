import i18n, { type BackendModule, type ResourceLanguage } from "i18next";
import { initReactI18next } from "react-i18next";
import en from "@/locales/en";
import {
  DEFAULT_LOCALE,
  isLocale,
  LOCALES,
  type Locale,
  toLocale,
} from "./locales";

const DEFAULT_NAMESPACE = "common";

type LazyLocale = Exclude<Locale, typeof DEFAULT_LOCALE>;

const lazyResources: Record<
  LazyLocale,
  () => Promise<{ default: ResourceLanguage }>
> = {
  ko: () => import("@/locales/ko"),
  ja: () => import("@/locales/ja"),
};

const lazyLocales: BackendModule = {
  type: "backend",
  init() {},
  read(language, namespace, callback) {
    if (!isLocale(language) || language === DEFAULT_LOCALE) {
      callback(new Error(`Unsupported language: ${language}`), null);
      return;
    }
    lazyResources[language]().then(
      (resources) => callback(null, resources.default[namespace]),
      (error) => callback(error, null),
    );
  },
};

function syncDocumentLanguage() {
  document.documentElement.lang = currentLocale();
}

export function currentLocale(): Locale {
  return toLocale(i18n.resolvedLanguage);
}

export async function activateLocale(locale: Locale) {
  if (i18n.resolvedLanguage !== locale) {
    await i18n.changeLanguage(locale);
  }
}

export async function initI18n() {
  i18n.on("languageChanged", syncDocumentLanguage);

  await i18n
    .use(lazyLocales)
    .use(initReactI18next)
    .init({
      resources: { [DEFAULT_LOCALE]: en },
      partialBundledLanguages: true,
      lng: DEFAULT_LOCALE,
      supportedLngs: LOCALES.map((locale) => locale.code),
      fallbackLng: DEFAULT_LOCALE,
      load: "languageOnly",
      ns: Object.keys(en),
      defaultNS: DEFAULT_NAMESPACE,
      interpolation: { escapeValue: false },
      react: { useSuspense: false },
    });

  syncDocumentLanguage();
}
