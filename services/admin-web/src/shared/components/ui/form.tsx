import { Label } from "@bookstore/design/ui/label";
import { type ReactNode, useId } from "react";
import { cn } from "@/shared/lib/utils";

interface FormFieldProps {
  label: string;
  errors?: unknown[];
  hint?: string;
  children: (props: {
    id: string;
    describedBy: string | undefined;
  }) => ReactNode;
}

export function FormField({ label, errors, hint, children }: FormFieldProps) {
  const id = useId();
  const messageId = `${id}-message`;
  const messages = (errors ?? []).filter(Boolean).map(String);
  const hasError = messages.length > 0;

  return (
    <div className="flex flex-col gap-1.5">
      <Label htmlFor={id} className={cn(hasError && "text-destructive")}>
        {label}
      </Label>
      {children({ id, describedBy: hasError || hint ? messageId : undefined })}
      {hasError ? (
        <p id={messageId} className="text-destructive text-xs" role="alert">
          {messages.join(", ")}
        </p>
      ) : hint ? (
        <p id={messageId} className="text-muted-foreground text-xs">
          {hint}
        </p>
      ) : null}
    </div>
  );
}

export function FormActions({ children }: { children: ReactNode }) {
  return (
    <div className="flex items-center justify-end gap-2 pt-2">{children}</div>
  );
}
