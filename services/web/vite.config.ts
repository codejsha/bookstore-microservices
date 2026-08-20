import path from "node:path";
import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const alloyOtlpUrl = env.VITE_ALLOY_OTLP_URL || "http://localhost:4318";
  const alloyFaroUrl = env.VITE_ALLOY_FARO_URL || "http://localhost:12347";

  return {
    plugins: [tanstackRouter({ quoteStyle: "double" }), react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      proxy: {
        "/otlp": {
          target: alloyOtlpUrl,
          changeOrigin: true,
          rewrite: (p) => p.replace(/^\/otlp/, ""),
        },
        "/faro/collect": {
          target: alloyFaroUrl,
          changeOrigin: true,
          rewrite: (p) => p.replace(/^\/faro\/collect/, "/collect"),
        },
      },
    },
  };
});
