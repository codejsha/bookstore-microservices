import type { Span } from "@opentelemetry/api";
import { SpanStatusCode } from "@opentelemetry/api";
import { usePathname } from "expo-router";
import { useEffect, useRef } from "react";
import { getTracer } from "@/shared/lib/telemetry";

export function useRouterTelemetry(): void {
  const pathname = usePathname();
  const currentSpan = useRef<Span | null>(null);
  const previous = useRef<string | undefined>(undefined);

  useEffect(() => {
    const tracer = getTracer();
    if (currentSpan.current) {
      currentSpan.current.setStatus({ code: SpanStatusCode.OK });
      currentSpan.current.end();
    }
    currentSpan.current = tracer.startSpan("router.navigate", {
      attributes: {
        "router.from": previous.current,
        "router.to": pathname,
      },
    });
    previous.current = pathname;
  }, [pathname]);

  useEffect(() => {
    return () => {
      currentSpan.current?.end();
      currentSpan.current = null;
    };
  }, []);
}
