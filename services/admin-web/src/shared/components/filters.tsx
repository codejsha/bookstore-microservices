import { Input } from "@bookstore/design/ui/input";
import { Label } from "@bookstore/design/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@bookstore/design/ui/select";
import { type ReactNode, useState } from "react";
import { isUuid } from "@/shared/lib/list-search";

export function FilterBar({ children }: { children: ReactNode }) {
  return <div className="flex flex-wrap items-end gap-3">{children}</div>;
}

interface FilterFieldProps {
  id: string;
  label: string;
  children: ReactNode;
}

function FilterField({ id, label, children }: FilterFieldProps) {
  return (
    <div className="flex flex-col gap-1">
      <Label htmlFor={id} className="text-muted-foreground text-xs">
        {label}
      </Label>
      {children}
    </div>
  );
}

const parseTrimmed = (raw: string) => raw.trim() || undefined;

export const parseUuidInput = (raw: string) => {
  const value = raw.trim();
  return isUuid(value) ? value : undefined;
};

interface FilterInputProps {
  id: string;
  label: string;
  value: string | undefined;
  onChange: (value: string | undefined) => void;
  parse?: (raw: string) => string | undefined;
  placeholder?: string;
  type?: "text" | "date";
  className?: string;
}

export function FilterInput({
  id,
  label,
  value,
  onChange,
  parse = parseTrimmed,
  placeholder,
  type = "text",
  className = "w-56",
}: FilterInputProps) {
  const [draft, setDraft] = useState(value ?? "");
  const [synced, setSynced] = useState(value);
  if (value !== synced) {
    setSynced(value);
    if (value !== parse(draft)) setDraft(value ?? "");
  }
  const invalid = draft.trim() !== "" && parse(draft) === undefined;

  return (
    <FilterField id={id} label={label}>
      <Input
        id={id}
        type={type}
        className={className}
        placeholder={placeholder}
        value={draft}
        aria-invalid={invalid || undefined}
        onChange={(e) => {
          setDraft(e.target.value);
          const next = parse(e.target.value);
          if (next !== value) onChange(next);
        }}
      />
    </FilterField>
  );
}

interface FilterSelectProps<T extends string> {
  id: string;
  label: string;
  value: T | undefined;
  options: readonly T[];
  allLabel: string;
  onChange: (value: T | undefined) => void;
}

export function FilterSelect<T extends string>({
  id,
  label,
  value,
  options,
  allLabel,
  onChange,
}: FilterSelectProps<T>) {
  const items = [
    { value: null, label: allLabel },
    ...options.map((option) => ({ value: option, label: option })),
  ];

  return (
    <FilterField id={id} label={label}>
      <Select
        items={items}
        value={value ?? null}
        onValueChange={(next) => onChange((next as T | null) ?? undefined)}
      >
        <SelectTrigger id={id} className="w-56">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {items.map((item) => (
            <SelectItem key={item.label} value={item.value}>
              {item.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </FilterField>
  );
}
