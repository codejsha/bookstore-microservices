import { getAccessToken, userManager } from "@/shared/auth/userManager";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "";

export interface ApiErrorBody {
  type?: string;
  title: string;
  status: number;
  detail?: string;
  errors?: string[];
}

export class ApiError extends Error {
  status: number;
  body: ApiErrorBody | null;

  constructor(status: number, statusText: string, body: ApiErrorBody | null) {
    super(body?.detail ?? body?.title ?? `${status} ${statusText}`);
    this.name = "ApiError";
    this.status = status;
    this.body = body;
  }

  get isUnauthorized() {
    return this.status === 401;
  }

  get isForbidden() {
    return this.status === 403;
  }

  get isNotFound() {
    return this.status === 404;
  }
}

async function getAuthHeaders(): Promise<Record<string, string>> {
  const token = await getAccessToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(await getAuthHeaders()),
      ...init?.headers,
    },
  });

  if (!res.ok) {
    let body: ApiErrorBody | null = null;
    try {
      body = await res.json();
    } catch {}

    if (res.status === 401) {
      const user = await userManager.getUser();
      if (user) {
        void userManager.signinRedirect();
      }
    }

    throw new ApiError(res.status, res.statusText, body);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return res.json() as Promise<T>;
}

export const api = {
  get<T>(path: string): Promise<T> {
    return request<T>(path);
  },

  post<T>(path: string, data?: unknown): Promise<T> {
    return request<T>(path, {
      method: "POST",
      body: data !== undefined ? JSON.stringify(data) : undefined,
    });
  },

  put<T>(path: string, data?: unknown): Promise<T> {
    return request<T>(path, {
      method: "PUT",
      body: data !== undefined ? JSON.stringify(data) : undefined,
    });
  },

  delete<T>(path: string): Promise<T> {
    return request<T>(path, { method: "DELETE" });
  },
} as const;
