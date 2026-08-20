import { makeRedirectUri } from "expo-auth-session";
import Constants from "expo-constants";

const extra = Constants.expoConfig?.extra ?? {};

export const authority: string = extra.keycloakAuthority ?? "";
export const clientId: string = extra.keycloakClientId ?? "";

if (!authority || !clientId) {
  console.warn(
    "[auth] keycloakAuthority or keycloakClientId is not set in app.json extra",
  );
}

export const discovery = {
  authorizationEndpoint: `${authority}/protocol/openid-connect/auth`,
  tokenEndpoint: `${authority}/protocol/openid-connect/token`,
  endSessionEndpoint: `${authority}/protocol/openid-connect/logout`,
  revocationEndpoint: `${authority}/protocol/openid-connect/revoke`,
};

export const redirectUri = makeRedirectUri({
  scheme: "bookstore",
  path: "callback",
});
