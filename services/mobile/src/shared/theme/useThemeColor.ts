import { useColorScheme } from "nativewind";
import { type ThemeColorName, themeColors } from "./colors";

export function useThemeColor() {
  const { colorScheme } = useColorScheme();
  const theme = colorScheme === "dark" ? "dark" : "light";
  const colors = themeColors[theme];
  return {
    theme,
    colors,
    color: (name: ThemeColorName) => colors[name],
  };
}
