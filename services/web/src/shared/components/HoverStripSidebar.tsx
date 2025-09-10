import { useEffect, useRef, useState } from "react";
import { cn } from "@/shared/lib/utils";

interface HoverStripSidebarProps {
  side: "left" | "right";
  expandDelayMs?: number;
  collapseDelayMs?: number;
  stripWidthClass?: string;
  expandedWidthClass?: string;
  className?: string;
  children: React.ReactNode;
}

export function HoverStripSidebar({
  side,
  expandDelayMs = 350,
  collapseDelayMs = 120,
  stripWidthClass = "w-12",
  expandedWidthClass = "w-56",
  className,
  children,
}: HoverStripSidebarProps) {
  const [expanded, setExpanded] = useState(false);
  const expandTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const collapseTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(
    () => () => {
      if (expandTimer.current) clearTimeout(expandTimer.current);
      if (collapseTimer.current) clearTimeout(collapseTimer.current);
    },
    [],
  );

  const cancelExpand = () => {
    if (expandTimer.current) {
      clearTimeout(expandTimer.current);
      expandTimer.current = null;
    }
  };
  const cancelCollapse = () => {
    if (collapseTimer.current) {
      clearTimeout(collapseTimer.current);
      collapseTimer.current = null;
    }
  };

  const handleEnter = () => {
    cancelCollapse();
    if (expanded || expandTimer.current) return;
    expandTimer.current = setTimeout(() => {
      setExpanded(true);
      expandTimer.current = null;
    }, expandDelayMs);
  };

  const handleLeave = () => {
    cancelExpand();
    if (!expanded || collapseTimer.current) return;
    collapseTimer.current = setTimeout(() => {
      setExpanded(false);
      collapseTimer.current = null;
    }, collapseDelayMs);
  };

  const handleFocus = () => {
    cancelCollapse();
    cancelExpand();
    setExpanded(true);
  };

  const handleBlur: React.FocusEventHandler<HTMLElement> = (event) => {
    if (event.currentTarget.contains(event.relatedTarget as Node | null)) {
      return;
    }
    handleLeave();
  };

  return (
    <aside
      data-expanded={expanded}
      data-side={side}
      onMouseEnter={handleEnter}
      onMouseLeave={handleLeave}
      onFocus={handleFocus}
      onBlur={handleBlur}
      className={cn(
        "group/strip sticky top-14 z-40 hidden h-[calc(100vh-3.5rem)] shrink-0 self-start overflow-hidden bg-background transition-[width] duration-200 ease-out md:block",
        side === "left" ? "border-r" : "border-l",
        expanded ? expandedWidthClass : stripWidthClass,
        className,
      )}
    >
      <div className="flex h-full flex-col gap-1 py-3">{children}</div>
    </aside>
  );
}

interface StripRowProps {
  icon: React.ReactNode;
  label: string;
  badge?: React.ReactNode;
}

export function StripRow({ icon, label, badge }: StripRowProps) {
  return (
    <span className="flex items-center gap-3">
      <span className="relative flex size-8 shrink-0 items-center justify-center">
        {icon}
        {badge}
      </span>
      <span className="whitespace-nowrap text-sm opacity-0 transition-opacity duration-150 group-data-[expanded=true]/strip:opacity-100">
        {label}
      </span>
    </span>
  );
}
