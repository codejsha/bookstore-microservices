import { useMemo } from "react";
import { createFormatters } from "@/shared/i18n/format";
import { useLocale } from "@/shared/i18n/useLocale";

export function useFormat() {
  const locale = useLocale();

  return useMemo(() => createFormatters(locale), [locale]);
}
