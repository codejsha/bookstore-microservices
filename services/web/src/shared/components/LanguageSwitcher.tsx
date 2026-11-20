import { Button } from "@bookstore/design/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@bookstore/design/ui/dropdown-menu";
import { useNavigate } from "@tanstack/react-router";
import { Languages } from "lucide-react";
import { useTranslation } from "react-i18next";
import { isLocale, LOCALES } from "@/shared/i18n/locales";
import { useLocale } from "@/shared/i18n/useLocale";

export function LanguageSwitcher() {
  const { t } = useTranslation();
  const locale = useLocale();
  const navigate = useNavigate();

  const switchTo = (next: unknown) => {
    if (!isLocale(next) || next === locale) return;
    navigate({
      to: ".",
      params: (current) => ({ ...current, locale: next }),
      search: true,
      hash: true,
    });
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger
        render={
          <Button
            variant="ghost"
            size="icon-sm"
            aria-label={t("language.label")}
          />
        }
      >
        <Languages className="size-4" />
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuRadioGroup value={locale} onValueChange={switchTo}>
          {LOCALES.map((entry) => (
            <DropdownMenuRadioItem
              key={entry.code}
              value={entry.code}
              lang={entry.code}
              closeOnClick
            >
              {entry.name}
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
