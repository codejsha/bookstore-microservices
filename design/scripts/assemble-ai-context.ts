import { existsSync, readFileSync, writeFileSync } from "node:fs";
import { basename, dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const rootDir = resolve(__dirname, "..");

interface ComponentMeta {
  name: string;
  description?: string;
  variants: Record<string, { svg: string; tokens: string[] }>;
  sizes?: string[];
  states: string[];
  responsive?: boolean | Record<string, Record<string, string>>;
  codeRef?: string;
  cvaVariants?: boolean;
}

interface AIContext {
  version: string;
  generated: string;
  pipeline: string[];
  tokens: {
    primitive: Record<string, unknown>;
    semantic: {
      light: Record<string, unknown>;
      dark: Record<string, unknown>;
    };
    component: Record<string, unknown>;
  };
  components: Array<{
    name: string;
    description?: string;
    meta: ComponentMeta;
    svgPaths: string[];
    codeRef?: string;
    cvaVariants?: boolean;
  }>;
}

const primitiveDir = resolve(rootDir, "tokens/primitive");
const primitives: Record<string, unknown> = {};
for (const file of [
  "color.json",
  "spacing.json",
  "typography.json",
  "radius.json",
]) {
  const filePath = resolve(primitiveDir, file);
  if (existsSync(filePath)) {
    primitives[basename(file, ".json")] = JSON.parse(
        readFileSync(filePath, "utf-8"),
    );
  }
}

const semanticLight = JSON.parse(
    readFileSync(resolve(rootDir, "tokens/semantic/color.json"), "utf-8"),
);
const semanticDark = JSON.parse(
    readFileSync(resolve(rootDir, "tokens/semantic/color.dark.json"), "utf-8"),
);

const componentDir = resolve(rootDir, "tokens/component");
const componentTokens: Record<string, unknown> = {};
for (const file of ["button.json", "card.json", "input.json", "sidebar.json"]) {
  const filePath = resolve(componentDir, file);
  if (existsSync(filePath)) {
    componentTokens[basename(file, ".json")] = JSON.parse(
        readFileSync(filePath, "utf-8"),
    );
  }
}

const designsDir = resolve(rootDir, "assets/components");
const components: AIContext["components"] = [];

const componentDirs = [
  "button",
  "card",
  "input",
  "label",
  "badge",
  "avatar",
  "dialog",
  "dropdown-menu",
  "separator",
  "select",
  "tabs",
  "textarea",
];

for (const dir of componentDirs) {
  const metaPath = resolve(designsDir, dir, "meta.json");
  if (existsSync(metaPath)) {
    const meta: ComponentMeta = JSON.parse(readFileSync(metaPath, "utf-8"));
    const svgPaths = Object.values(meta.variants).map(
        (v) => `assets/components/${dir}/${v.svg}`,
    );

    components.push({
      name: meta.name,
      description: meta.description,
      meta,
      svgPaths,
      codeRef: meta.codeRef,
      cvaVariants: meta.cvaVariants,
    });
  }
}

const context: AIContext = {
  version: "1.0.0",
  generated: new Date().toISOString(),
  pipeline: [
    "1. Primitive tokens (시스템 이해)",
    "2. Semantic tokens (의미 매핑)",
    "3. Component tokens (해당 컴포넌트 값)",
    "4. SVG (시각적 구조)",
    "5. Meta JSON (상태/변형 명세)",
    "6. 기존 코드 예시 (프로젝트 패턴)",
  ],
  tokens: {
    primitive: primitives,
    semantic: {
      light: semanticLight,
      dark: semanticDark,
    },
    component: componentTokens,
  },
  components,
};

const outputPath = resolve(rootDir, "build/ai-context.json");
writeFileSync(outputPath, JSON.stringify(context, null, 2) + "\n");

console.log(`Context assembled → ${outputPath}`)
console.log(`  ${components.length} components`);
console.log(`  ${Object.keys(primitives).length} primitive categories`);
console.log(`  ${Object.keys(componentTokens).length} component token sets`);
