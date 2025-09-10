import type { Config } from "tailwindcss";

const v = (name: string) => `var(--color-${name})`;

export default {
  darkMode: "class",
  content: [
    "./app/**/*.{ts,tsx}",
    "./src/**/*.{ts,tsx}",
    "./node_modules/@bookstore/design/components/native/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: v("background"),
        foreground: v("foreground"),
        card: { DEFAULT: v("card"), foreground: v("card-foreground") },
        popover: { DEFAULT: v("popover"), foreground: v("popover-foreground") },
        primary: { DEFAULT: v("primary"), foreground: v("primary-foreground") },
        secondary: {
          DEFAULT: v("secondary"),
          foreground: v("secondary-foreground"),
        },
        muted: { DEFAULT: v("muted"), foreground: v("muted-foreground") },
        accent: { DEFAULT: v("accent"), foreground: v("accent-foreground") },
        destructive: {
          DEFAULT: v("destructive"),
          foreground: v("destructive-foreground"),
        },
        border: v("border"),
        input: v("input"),
        ring: v("ring"),
        success: v("success"),
      },
      borderRadius: {
        sm: "calc(0.625rem - 4px)",
        md: "calc(0.625rem - 2px)",
        lg: "0.625rem",
      },
    },
  },
  plugins: [],
} satisfies Config;
