import type { ConfigContext, ExpoConfig } from "expo/config";

const DEV_API_BASE_URL = "http://localhost:8080";
const DEV_KEYCLOAK_AUTHORITY = "http://localhost:8080/realms/bookstore";
const DEV_KEYCLOAK_CLIENT_ID = "bookstore-mobile";

export default ({ config }: ConfigContext): ExpoConfig => ({
  ...config,
  name: config.name ?? "Bookstore",
  slug: config.slug ?? "bookstore",
  extra: {
    ...config.extra,
    apiBaseUrl: process.env.EXPO_PUBLIC_API_BASE_URL ?? DEV_API_BASE_URL,
    otelCollectorUrl: process.env.EXPO_PUBLIC_OTEL_COLLECTOR_URL ?? "",
    keycloakAuthority:
      process.env.EXPO_PUBLIC_KEYCLOAK_AUTHORITY ?? DEV_KEYCLOAK_AUTHORITY,
    keycloakClientId:
      process.env.EXPO_PUBLIC_KEYCLOAK_CLIENT_ID ?? DEV_KEYCLOAK_CLIENT_ID,
  },
});
