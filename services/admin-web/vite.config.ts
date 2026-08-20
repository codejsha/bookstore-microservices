import path from "node:path";
import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const adminApiUrl = env.VITE_DEV_API_PROXY || "http://localhost:8080";

  return {
    plugins: [tanstackRouter({ quoteStyle: "double" }), react(), tailwindcss()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    server: {
      proxy: {
        "/api": {
          target: adminApiUrl,
          changeOrigin: true,
        },
      },
    },
  };
});
