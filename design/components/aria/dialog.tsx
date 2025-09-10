import { XIcon } from "lucide-react";
import type * as React from "react";
import {
  Dialog as AriaDialog,
  type DialogProps as AriaDialogProps,
  DialogTrigger as AriaDialogTrigger,
  Heading as AriaHeading,
  type HeadingProps as AriaHeadingProps,
  Modal as AriaModal,
  ModalOverlay as AriaModalOverlay,
  type ModalOverlayProps as AriaModalOverlayProps,
  composeRenderProps,
} from "react-aria-components";
import { cn } from "../lib/utils";
import { Button } from "./button";

const Dialog = AriaDialogTrigger;

function DialogOverlay({
  className,
  isDismissable = true,
  ...props
}: AriaModalOverlayProps) {
  return (
    <AriaModalOverlay
      data-slot="dialog-overlay"
      isDismissable={isDismissable}
      className={composeRenderProps(className, (cls) =>
        cn(
          "fixed inset-0 isolate z-50 flex items-center justify-center bg-black/10 duration-100 supports-backdrop-filter:backdrop-blur-xs data-[entering]:animate-in data-[entering]:fade-in-0 data-[exiting]:animate-out data-[exiting]:fade-out-0",
          cls,
        ),
      )}
      {...props}
    />
  );
}

function DialogContent({
  className,
  children,
  showCloseButton = true,
  ...props
}: AriaDialogProps & {
  showCloseButton?: boolean;
}) {
  return (
    <DialogOverlay>
      <AriaModal
        className={cn(
          "z-50 w-full max-w-[calc(100%-2rem)] outline-none duration-100 sm:max-w-sm data-[entering]:animate-in data-[entering]:fade-in-0 data-[entering]:zoom-in-95 data-[exiting]:animate-out data-[exiting]:fade-out-0 data-[exiting]:zoom-out-95",
        )}
      >
        <AriaDialog
          data-slot="dialog-content"
          className={cn(
            "relative grid gap-4 rounded-xl bg-popover p-4 text-sm text-popover-foreground ring-1 ring-foreground/10 outline-none",
            className,
          )}
          {...props}
        >
          {composeRenderProps(children, (resolved) => (
            <>
              {resolved}
              {showCloseButton && (
                <DialogCloseButton aria-label="Close">
                  <XIcon />
                </DialogCloseButton>
              )}
            </>
          ))}
        </AriaDialog>
      </AriaModal>
    </DialogOverlay>
  );
}

function DialogCloseButton({
  className,
  ...props
}: React.ComponentProps<typeof Button>) {
  return (
    <Button
      slot="close"
      variant="ghost"
      size="icon-sm"
      data-slot="dialog-close"
      className={composeRenderProps(className, (cls) =>
        cn("absolute top-2 right-2", cls),
      )}
      {...props}
    />
  );
}

function DialogHeader({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="dialog-header"
      className={cn("flex flex-col gap-2", className)}
      {...props}
    />
  );
}

function DialogFooter({
  className,
  showCloseButton = false,
  children,
  ...props
}: React.ComponentProps<"div"> & { showCloseButton?: boolean }) {
  return (
    <div
      data-slot="dialog-footer"
      className={cn(
        "-mx-4 -mb-4 flex flex-col-reverse gap-2 rounded-b-xl border-t bg-muted/50 p-4 sm:flex-row sm:justify-end",
        className,
      )}
      {...props}
    >
      {children}
      {showCloseButton && (
        <Button slot="close" variant="outline">
          Close
        </Button>
      )}
    </div>
  );
}

function DialogTitle({ className, ...props }: AriaHeadingProps) {
  return (
    <AriaHeading
      slot="title"
      data-slot="dialog-title"
      className={cn(
        "font-heading text-base leading-none font-medium",
        className,
      )}
      {...props}
    />
  );
}

function DialogDescription({ className, ...props }: React.ComponentProps<"p">) {
  return (
    <p
      data-slot="dialog-description"
      className={cn(
        "text-sm text-muted-foreground *:[a]:underline *:[a]:underline-offset-3 *:[a]:hover:text-foreground",
        className,
      )}
      {...props}
    />
  );
}

export {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogOverlay,
  DialogTitle,
};
