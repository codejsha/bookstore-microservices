import * as SecureStore from "expo-secure-store";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

export interface SessionInput {
  token: string;
  refreshToken: string | null;
  expiresIn: number | null;
}

interface PersistedAuthState {
  token: string | null;
  refreshToken: string | null;
  expiresAt: number | null;
}

interface AuthState extends PersistedAuthState {
  setSession: (session: SessionInput) => void;
  setToken: (token: string | null) => void;
  clear: () => void;
}

const secureStorage = createJSONStorage<PersistedAuthState>(() => ({
  getItem: (key: string) => SecureStore.getItemAsync(key),
  setItem: (key: string, value: string) => SecureStore.setItemAsync(key, value),
  removeItem: (key: string) => SecureStore.deleteItemAsync(key),
}));

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      refreshToken: null,
      expiresAt: null,
      setSession: ({ token, refreshToken, expiresIn }) =>
        set({
          token,
          refreshToken,
          expiresAt: expiresIn != null ? Date.now() + expiresIn * 1000 : null,
        }),
      setToken: (token) => set({ token }),
      clear: () => set({ token: null, refreshToken: null, expiresAt: null }),
    }),
    {
      name: "bookstore-auth",
      storage: secureStorage,
      version: 1,
      partialize: (state) => ({
        token: state.token,
        refreshToken: state.refreshToken,
        expiresAt: state.expiresAt,
      }),
      migrate: (persisted, version) => {
        if (version === 0) {
          const prev = persisted as { token?: string | null } | null;
          return {
            token: prev?.token ?? null,
            refreshToken: null,
            expiresAt: null,
          } satisfies PersistedAuthState;
        }
        return persisted as PersistedAuthState;
      },
    },
  ),
);

export function isSessionValid(state: {
  token: string | null;
  expiresAt: number | null;
}): boolean {
  if (!state.token) return false;
  if (state.expiresAt != null && state.expiresAt <= Date.now()) return false;
  return true;
}
