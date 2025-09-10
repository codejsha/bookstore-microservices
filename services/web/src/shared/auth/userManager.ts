import { UserManager, WebStorageStateStore } from "oidc-client-ts";

const authority = import.meta.env.VITE_KEYCLOAK_AUTHORITY;
const clientId = import.meta.env.VITE_KEYCLOAK_CLIENT_ID;

if (!authority || !clientId) {
  console.warn(
    "[auth] VITE_KEYCLOAK_AUTHORITY or VITE_KEYCLOAK_CLIENT_ID is not set — OIDC will not work",
  );
}

export const userManager = new UserManager({
  authority: authority ?? "",
  client_id: clientId ?? "",
  redirect_uri: `${window.location.origin}/callback`,
  post_logout_redirect_uri: window.location.origin,
  response_type: "code",
  scope: "openid profile email",
  automaticSilentRenew: true,
  loadUserInfo: true,
  userStore: new WebStorageStateStore({ store: window.sessionStorage }),
});

export async function getAccessToken(): Promise<string | null> {
  const user = await userManager.getUser();
  if (!user || user.expired) return null;
  return user.access_token ?? null;
}
