import { CheckCircle2, Info, X, XCircle } from "lucide-react";
import { createPortal } from "react-dom";
import type { ToastItem, ToastVariant } from "./toast-context";

const variantConfig: Record<
  ToastVariant,
  { icon: typeof Info; iconClass: string }
> = {
  success: { icon: CheckCircle2, iconClass: "text-green-600" },
  error: { icon: XCircle, iconClass: "text-destructive" },
  info: { icon: Info, iconClass: "text-primary" },
};

export function ToastViewport({
  toasts,
  onDismiss,
}: {
  toasts: ToastItem[];
  onDismiss: (id: string) => void;
}) {
  if (typeof document === "undefined") return null;

  return createPortal(
    <div className="pointer-events-none fixed bottom-4 right-4 z-[100] flex w-full max-w-sm flex-col gap-2">
      {toasts.map((toast) => {
        const { icon: Icon, iconClass } = variantConfig[toast.variant];
        return (
          <div
            key={toast.id}
            role="status"
            aria-live="polite"
            className="pointer-events-auto flex items-start gap-3 rounded-lg border bg-popover p-3 text-popover-foreground shadow-lg animate-in fade-in slide-in-from-bottom-2 duration-200"
          >
            <Icon className={`mt-0.5 size-5 shrink-0 ${iconClass}`} />
            <div className="min-w-0 flex-1">
              <p className="text-sm font-medium">{toast.message}</p>
              {toast.description && (
                <p className="mt-0.5 text-sm text-muted-foreground">
                  {toast.description}
                </p>
              )}
            </div>
            <button
              type="button"
              aria-label="Dismiss"
              onClick={() => onDismiss(toast.id)}
              className="shrink-0 rounded-md p-0.5 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
            >
              <X className="size-4" />
            </button>
          </div>
        );
      })}
    </div>,
    document.body,
  );
}
