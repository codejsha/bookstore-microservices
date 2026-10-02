import { useRouter } from "@tanstack/react-router";
import { Fragment, type ReactNode } from "react";
import { ErrorState } from "@/shared/components/ErrorState";
import { EMPTY } from "@/shared/lib/format";

interface DetailHeaderProps {
  title: ReactNode;
  subtitle?: ReactNode;
  back: ReactNode;
  actions?: ReactNode;
}

export function DetailHeader({
  title,
  subtitle,
  back,
  actions,
}: DetailHeaderProps) {
  return (
    <div className="flex flex-col gap-2">
      <p className="text-muted-foreground text-xs">{back}</p>
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="font-semibold text-2xl">{title}</h1>
          {subtitle ? (
            <p className="text-muted-foreground text-xs">{subtitle}</p>
          ) : null}
        </div>
        {actions ? (
          <div className="flex items-center gap-3 text-sm">{actions}</div>
        ) : null}
      </div>
    </div>
  );
}

interface DetailSectionProps {
  title: string;
  children: ReactNode;
}

export function DetailSection({ title, children }: DetailSectionProps) {
  return (
    <section className="flex flex-col gap-3">
      <h2 className="font-semibold text-lg">{title}</h2>
      {children}
    </section>
  );
}

export interface DefinitionItem {
  label: string;
  value: ReactNode;
}

export function DefinitionList({ items }: { items: DefinitionItem[] }) {
  return (
    <dl className="grid gap-x-6 gap-y-2 text-sm sm:grid-cols-[12rem_1fr]">
      {items.map((item) => (
        <Fragment key={item.label}>
          <dt className="text-muted-foreground">{item.label}</dt>
          <dd className="break-all">{item.value ?? EMPTY}</dd>
        </Fragment>
      ))}
    </dl>
  );
}

interface DetailErrorProps {
  error: unknown;
  entity: string;
}

export function DetailNotFound({ entity }: { entity: string }) {
  return (
    <ErrorState
      title={`${entity.charAt(0).toUpperCase()}${entity.slice(1)} not found`}
      message={`No ${entity} exists with this id.`}
    />
  );
}

export function DetailError({ error, entity }: DetailErrorProps) {
  const router = useRouter();
  return (
    <ErrorState
      title={`Unable to load this ${entity}`}
      message={
        error instanceof Error
          ? error.message
          : "An unexpected error occurred. Please try again."
      }
      onRetry={() => router.invalidate()}
    />
  );
}
