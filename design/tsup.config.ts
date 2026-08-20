import { defineConfig } from "tsup";

export default defineConfig({
  entry: [
    "components/ui/*.tsx",
    "components/aria/*.tsx",
    "components/lib/*.ts",
    "!components/**/*.stories.tsx",
  ],
  outDir: "dist",
  format: ["esm"],
  target: "es2022",
  dts: true,
  splitting: true,
  treeshake: true,
  sourcemap: true,
  clean: true,
  tsconfig: "tsconfig.build.json",
  external: ["react", "react-dom", "react/jsx-runtime", "tailwindcss"],
});
