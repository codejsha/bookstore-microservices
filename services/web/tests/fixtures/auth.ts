import type { BrowserContext } from "@playwright/test";

const authority =
  process.env.VITE_KEYCLOAK_AUTHORITY ??
  "http://localhost:8080/realms/bookstore";
const clientId = process.env.VITE_KEYCLOAK_CLIENT_ID ?? "bookstore-web";

export const testUser = {
  sub: "user-1",
  email: "john.doe@example.com",
  name: "John Doe",
  preferred_username: "john.doe",
};

interface StoredSession {
  key: string;
  value: string;
}

export function buildStoredSession(now = Date.now()): StoredSession {
  const issuedAt = Math.floor(now / 1000);
  const expiresAt = issuedAt + 60 * 60;
  return {
    key: `oidc.user:${authority}:${clientId}`,
    value: JSON.stringify({
      id_token: "e2e-id-token",
      access_token: "e2e-access-token",
      token_type: "Bearer",
      scope: "openid profile email",
      expires_at: expiresAt,
      profile: {
        ...testUser,
        iss: authority,
        aud: clientId,
        iat: issuedAt,
        exp: expiresAt,
      },
    }),
  };
}

export async function seedStoredSession(
  context: BrowserContext,
): Promise<void> {
  await context.addInitScript((session: StoredSession) => {
    window.sessionStorage.setItem(session.key, session.value);
  }, buildStoredSession());
}
