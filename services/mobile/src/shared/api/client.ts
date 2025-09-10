import { SpanKind } from "@opentelemetry/api";
import Constants from "expo-constants";
import { isSessionValid, useAuthStore } from "@/domains/auth";
import { injectTraceHeaders, withSpan } from "@/shared/lib/telemetry";

const API_BASE_URL =
  Constants.expoConfig?.extra?.apiBaseUrl ?? "http://localhost:8080";

export interface ApiErrorBody {
  code: number;
  message: string;
  details?: string[];
}

export class ApiError extends Error {
  status: number;
  body: ApiErrorBody | null;

  constructor(status: number, statusText: string, body: ApiErrorBody | null) {
    super(body?.message ?? `${status} ${statusText}`);
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

function getAuthHeaders(): Record<string, string> {
  const state = useAuthStore.getState();
  return isSessionValid(state)
    ? { Authorization: `Bearer ${state.token}` }
    : {};
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const method = (init?.method ?? "GET").toUpperCase();
  const url = `${API_BASE_URL}${path}`;

  return withSpan(
    `HTTP ${method}`,
    async (span) => {
      const headers = injectTraceHeaders({
        "Content-Type": "application/json",
        ...getAuthHeaders(),
        ...(init?.headers as Record<string, string> | undefined),
      });

      const res = await fetch(url, { ...init, headers });

      span.setAttribute("http.request.method", method);
      span.setAttribute("url.full", url);
      span.setAttribute("http.response.status_code", res.status);

      if (!res.ok) {
        let body: ApiErrorBody | null = null;
        try {
          body = await res.json();
        } catch {}
        throw new ApiError(res.status, res.statusText, body);
      }

      if (res.status === 204) {
        return undefined as T;
      }

      return res.json() as Promise<T>;
    },
    { kind: SpanKind.CLIENT },
  );
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
