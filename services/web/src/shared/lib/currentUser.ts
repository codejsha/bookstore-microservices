import { useAuth } from "react-oidc-context";

export function useCurrentUserUid(): string | undefined {
  const auth = useAuth();
  return auth.user?.profile.sub;
}
