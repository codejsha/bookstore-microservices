import { Button } from "@bookstore/design/ui/button";
import { Moon, Sun } from "lucide-react";
import { useTheme } from "@/shared/hooks/useTheme";

export function ThemeToggle() {
  const { theme, setTheme } = useTheme();

  const resolved =
    theme === "system"
      ? window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light"
      : theme;

  return (
    <Button
      variant="ghost"
      size="icon-sm"
      aria-label={
        resolved === "dark" ? "Switch to light theme" : "Switch to dark theme"
      }
      onClick={() => setTheme(resolved === "dark" ? "light" : "dark")}
    >
      {resolved === "dark" ? (
        <Sun className="size-4" />
      ) : (
        <Moon className="size-4" />
      )}
    </Button>
  );
}
