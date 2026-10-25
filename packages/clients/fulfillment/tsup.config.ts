import { defineConfig } from "tsup";

export default defineConfig({
  entry: ["src/model/*.ts", "src/constant/*.ts", "src/api/*.ts"],
  outDir: "dist",
  format: ["esm"],
  target: "es2022",
  dts: true,
  splitting: false,
  treeshake: true,
  sourcemap: true,
  clean: true,
  external: ["zod"],
});
