import { CheckIcon, ChevronRightIcon } from "lucide-react";
import type * as React from "react";
import {
  Header as AriaHeader,
  type HeaderProps as AriaHeaderProps,
  Keyboard as AriaKeyboard,
  Menu as AriaMenu,
  MenuItem as AriaMenuItem,
  type MenuItemProps as AriaMenuItemProps,
  type MenuProps as AriaMenuProps,
  MenuSection as AriaMenuSection,
  type MenuSectionProps as AriaMenuSectionProps,
  MenuTrigger as AriaMenuTrigger,
  Popover as AriaPopover,
  type PopoverProps as AriaPopoverProps,
  Separator as AriaSeparator,
  SubmenuTrigger as AriaSubmenuTrigger,
  composeRenderProps,
} from "react-aria-components";

import { cn } from "../lib/utils";

const DropdownMenu = AriaMenuTrigger;
const DropdownMenuSub = AriaSubmenuTrigger;

function DropdownMenuPopover({ className, ...props }: AriaPopoverProps) {
  return (
    <AriaPopover
      data-slot="dropdown-menu-popover"
      className={composeRenderProps(className, (cls) =>
        cn(
          "z-50 max-h-(--trigger-height) min-w-32 origin-(--origin) overflow-x-hidden overflow-y-auto rounded-lg bg-popover text-popover-foreground shadow-md ring-1 ring-foreground/10 outline-none data-[entering]:animate-in data-[entering]:fade-in-0 data-[entering]:zoom-in-95 data-[exiting]:animate-out data-[exiting]:fade-out-0 data-[exiting]:zoom-out-95 data-[placement=bottom]:slide-in-from-top-2 data-[placement=left]:slide-in-from-right-2 data-[placement=right]:slide-in-from-left-2 data-[placement=top]:slide-in-from-bottom-2",
          cls,
        ),
      )}
      {...props}
    />
  );
}

function DropdownMenuContent<T extends object>({
  className,
  ...props
}: AriaMenuProps<T>) {
  return (
    <DropdownMenuPopover>
      <AriaMenu
        data-slot="dropdown-menu-content"
        className={cn("p-1 outline-none", className)}
        {...props}
      />
    </DropdownMenuPopover>
  );
}

function DropdownMenuLabel({
  className,
  inset,
  ...props
}: AriaHeaderProps & { inset?: boolean }) {
  return (
    <AriaHeader
      data-slot="dropdown-menu-label"
      data-inset={inset}
      className={cn(
        "px-1.5 py-1 text-xs font-medium text-muted-foreground data-inset:pl-7",
        className,
      )}
      {...props}
    />
  );
}

function DropdownMenuItem({
  className,
  inset,
  variant = "default",
  ...props
}: AriaMenuItemProps & {
  inset?: boolean;
  variant?: "default" | "destructive";
}) {
  return (
    <AriaMenuItem
      data-slot="dropdown-menu-item"
      data-inset={inset}
      data-variant={variant}
      className={composeRenderProps(className, (cls) =>
        cn(
          "group/dropdown-menu-item relative flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none data-[focused]:bg-accent data-[focused]:text-accent-foreground data-inset:pl-7 data-[variant=destructive]:text-destructive data-[variant=destructive]:data-[focused]:bg-destructive/10 data-[variant=destructive]:data-[focused]:text-destructive data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 data-[variant=destructive]:*:[svg]:text-destructive",
          cls,
        ),
      )}
      {...props}
    />
  );
}

function DropdownMenuCheckboxItem({
  className,
  children,
  ...props
}: AriaMenuItemProps) {
  return (
    <AriaMenuItem
      data-slot="dropdown-menu-checkbox-item"
      className={composeRenderProps(className, (cls) =>
        cn(
          "relative flex cursor-default items-center gap-1.5 rounded-md py-1 pr-8 pl-1.5 text-sm outline-hidden select-none data-[focused]:bg-accent data-[focused]:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
          cls,
        ),
      )}
      {...props}
    >
      {composeRenderProps(children, (resolved, { isSelected }) => (
        <>
          <span className="pointer-events-none absolute right-2 flex items-center justify-center">
            {isSelected ? <CheckIcon /> : null}
          </span>
          {resolved}
        </>
      ))}
    </AriaMenuItem>
  );
}

const DropdownMenuRadioItem = DropdownMenuCheckboxItem;

function DropdownMenuSection<T extends object>({
  className,
  ...props
}: AriaMenuSectionProps<T>) {
  return (
    <AriaMenuSection
      data-slot="dropdown-menu-section"
      className={cn(className)}
      {...props}
    />
  );
}

function DropdownMenuSeparator({
  className,
  ...props
}: React.ComponentProps<typeof AriaSeparator>) {
  return (
    <AriaSeparator
      data-slot="dropdown-menu-separator"
      className={cn("-mx-1 my-1 h-px bg-border", className)}
      {...props}
    />
  );
}

function DropdownMenuShortcut({
  className,
  ...props
}: React.ComponentProps<typeof AriaKeyboard>) {
  return (
    <AriaKeyboard
      data-slot="dropdown-menu-shortcut"
      className={cn(
        "ml-auto text-xs tracking-widest text-muted-foreground group-data-[focused]/dropdown-menu-item:text-accent-foreground",
        className,
      )}
      {...props}
    />
  );
}

function DropdownMenuSubTrigger({
  className,
  children,
  inset,
  ...props
}: AriaMenuItemProps & { inset?: boolean }) {
  return (
    <AriaMenuItem
      data-slot="dropdown-menu-sub-trigger"
      data-inset={inset}
      className={composeRenderProps(className, (cls) =>
        cn(
          "flex cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none data-[focused]:bg-accent data-[focused]:text-accent-foreground data-inset:pl-7 data-[open]:bg-accent data-[open]:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4",
          cls,
        ),
      )}
      {...props}
    >
      {composeRenderProps(children, (resolved) => (
        <>
          {resolved}
          <ChevronRightIcon className="ml-auto" />
        </>
      ))}
    </AriaMenuItem>
  );
}

export {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPopover,
  DropdownMenuRadioItem,
  DropdownMenuSection,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuSub,
  DropdownMenuSubTrigger,
};
