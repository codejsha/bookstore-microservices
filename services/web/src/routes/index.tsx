import { createFileRoute, redirect } from "@tanstack/react-router";
import { DEFAULT_LOCALE } from "@/shared/i18n/locales";

export const Route = createFileRoute("/")({
  beforeLoad: ({ location }) => {
    throw redirect({
      to: "/$locale",
      params: { locale: DEFAULT_LOCALE },
      search: location.search,
      hash: location.hash,
      replace: true,
    });
  },
});
