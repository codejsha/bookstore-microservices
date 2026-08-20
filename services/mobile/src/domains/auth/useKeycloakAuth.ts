import {
  exchangeCodeAsync,
  ResponseType,
  useAuthRequest,
} from "expo-auth-session";
import * as WebBrowser from "expo-web-browser";
import { useEffect } from "react";
import { useAuthStore } from "./auth-store";
import { clientId, discovery, redirectUri } from "./oidc";

WebBrowser.maybeCompleteAuthSession();

interface Options {
  register?: boolean;
}

export function useKeycloakAuth({ register = false }: Options = {}) {
  const setSession = useAuthStore((s) => s.setSession);

  const [request, result, promptAsync] = useAuthRequest(
    {
      clientId,
      scopes: ["openid", "profile", "email"],
      redirectUri,
      responseType: ResponseType.Code,
      usePKCE: true,
      extraParams: register ? { kc_action: "register" } : undefined,
    },
    discovery,
  );

  useEffect(() => {
    if (result?.type !== "success") return;
    const code = result.params.code;
    const codeVerifier = request?.codeVerifier;
    if (!code || !codeVerifier) return;

    exchangeCodeAsync(
      {
        clientId,
        code,
        redirectUri,
        extraParams: { code_verifier: codeVerifier },
      },
      discovery,
    )
      .then((token) => {
        setSession({
          token: token.accessToken,
          refreshToken: token.refreshToken ?? null,
          expiresIn: token.expiresIn ?? null,
        });
      })
      .catch((err) => {
        console.warn("[auth] token exchange failed:", err);
      });
  }, [result, request, setSession]);

  return { request, result, promptAsync };
}
