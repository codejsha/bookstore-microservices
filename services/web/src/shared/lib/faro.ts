import {
  ConsoleInstrumentation,
  ErrorsInstrumentation,
  type Faro,
  initializeFaro,
  SessionInstrumentation,
  ViewInstrumentation,
} from "@grafana/faro-web-sdk";

const APP_NAME = "frontend-web";
const APP_VERSION = "0.1.0";
const COLLECT_URL = import.meta.env.VITE_FARO_COLLECT_URL ?? "/faro/collect";

let faroInstance: Faro | undefined;

export function initFaro(): Faro | undefined {
  if (!COLLECT_URL) {
    console.debug("[Faro] No collect URL configured, skipping init");
    return;
  }
  if (faroInstance) return faroInstance;

  faroInstance = initializeFaro({
    url: COLLECT_URL,
    app: {
      name: APP_NAME,
      version: APP_VERSION,
      environment: import.meta.env.MODE,
    },
    instrumentations: [
      new ErrorsInstrumentation(),
      new ConsoleInstrumentation(),
      new ViewInstrumentation(),
      new SessionInstrumentation(),
    ],
  });

  return faroInstance;
}

export function getFaro(): Faro | undefined {
  return faroInstance;
}
