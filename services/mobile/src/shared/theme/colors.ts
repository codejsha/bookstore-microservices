import { tokens as darkTokens } from "@bookstore/design/tokens/ts/dark";
import { tokens as lightTokens } from "@bookstore/design/tokens/ts/light";

export const themeColors = {
  light: lightTokens,
  dark: darkTokens,
} as const;

export type ThemeName = keyof typeof themeColors;
export type ThemeColorName = keyof typeof lightTokens & keyof typeof darkTokens;
