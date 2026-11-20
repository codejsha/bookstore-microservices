import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { activateLocale } from "@/shared/i18n/i18n";
import { DEFAULT_LOCALE, isLocale } from "@/shared/i18n/locales";

export const Route = createFileRoute("/$locale")({
  beforeLoad: async ({ params, location }) => {
    if (!isLocale(params.locale)) {
      throw redirect({
        to: "/$locale/$",
        params: {
          locale: DEFAULT_LOCALE,
          _splat: location.pathname.replace(/^\/+/, ""),
        },
        search: location.search,
        hash: location.hash,
        replace: true,
      });
    }
    await activateLocale(params.locale);
  },
  component: () => <Outlet />,
});
