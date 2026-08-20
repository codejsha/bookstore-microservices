import type { StorybookConfig } from "@storybook/react-vite";
import tailwindcss from "@tailwindcss/vite";

const config: StorybookConfig = {
  stories: ["../components/**/*.mdx", "../components/**/*.stories.@(ts|tsx)"],
  addons: ["@storybook/addon-docs", "@storybook/addon-themes"],
  framework: {
    name: "@storybook/react-vite",
    options: {},
  },
  typescript: {
    reactDocgen: "react-docgen-typescript",
  },
  async viteFinal(config) {
    config.plugins ??= [];
    config.plugins.push(tailwindcss());

    config.resolve ??= {};
    const existingAlias = config.resolve.alias;
    const existingArray = Array.isArray(existingAlias)
      ? existingAlias
      : existingAlias
        ? Object.entries(existingAlias).map(([find, replacement]) => ({
            find,
            replacement: replacement as string,
          }))
        : [];
    config.resolve.alias = [
      { find: /^react-native$/, replacement: "react-native-web" },
      {
        find: /^react-native\/Libraries\/Image\/AssetRegistry$/,
        replacement: "react-native-web/dist/modules/AssetRegistry",
      },
      ...existingArray,
    ];
    config.resolve.extensions = [
      ".web.tsx",
      ".web.ts",
      ".web.jsx",
      ".web.js",
      ...(config.resolve.extensions ?? [
        ".mjs",
        ".js",
        ".ts",
        ".tsx",
        ".jsx",
        ".json",
      ]),
    ];

    config.define = {
      ...config.define,
      __DEV__: JSON.stringify(true),
    };

    return config;
  },
};

export default config;
