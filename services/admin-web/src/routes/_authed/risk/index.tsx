import { Button } from "@bookstore/design/ui/button";
import { Input } from "@bookstore/design/ui/input";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import type { ColumnDef, SortingState } from "@tanstack/react-table";
import { useMemo, useState } from "react";
import type { RiskEntry, RiskLevel } from "@/domains/risk";
import {
  riskListQueryOptions,
  useFlagRisk,
  useUnflagRisk,
} from "@/domains/risk";
import { DataTable } from "@/shared/components/ui/data-table";

export const Route = createFileRoute("/_authed/risk/")({
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(riskListQueryOptions()),
  component: RiskPage,
});

const DAY_SECONDS = 86_400;

function RiskPage() {
  const { data } = useSuspenseQuery(riskListQueryOptions());
  const flagRisk = useFlagRisk();
  const unflagRisk = useUnflagRisk();

  const [sorting, setSorting] = useState<SortingState>([]);
  const [uid, setUid] = useState("");
  const [level, setLevel] = useState<RiskLevel>("restrict");
  const [reason, setReason] = useState("");
  const [ttlDays, setTtlDays] = useState("1");

  const columns = useMemo<ColumnDef<RiskEntry, unknown>[]>(
    () => [
      { accessorKey: "user_uid", header: "User" },
      {
        accessorKey: "level",
        header: "Level",
        cell: ({ row }) => (
          <span
            className={
              row.original.level === "block"
                ? "font-medium text-destructive"
                : "font-medium text-amber-600 dark:text-amber-400"
            }
          >
            {row.original.level}
          </span>
        ),
      },
      { accessorKey: "reason", header: "Reason" },
      {
        accessorKey: "flagged_by",
        header: "Flagged by",
        cell: ({ row }) => row.original.flagged_by ?? "—",
      },
      {
        accessorKey: "expires_at",
        header: "Expires",
        cell: ({ row }) => new Date(row.original.expires_at).toLocaleString(),
      },
      {
        id: "actions",
        header: "",
        enableSorting: false,
        cell: ({ row }) => (
          <Button
            variant="outline"
            size="sm"
            disabled={unflagRisk.isPending}
            onClick={() => unflagRisk.mutate(row.original.user_uid)}
          >
            Unflag
          </Button>
        ),
      },
    ],
    [unflagRisk],
  );

  const canSubmit =
    uid.trim() !== "" && reason.trim() !== "" && !flagRisk.isPending;

  const submit = () => {
    const days = Number(ttlDays);
    flagRisk.mutate(
      {
        uid: uid.trim(),
        body: {
          level,
          reason: reason.trim(),
          ttl_seconds:
            Number.isFinite(days) && days > 0
              ? Math.round(days * DAY_SECONDS)
              : undefined,
        },
      },
      {
        onSuccess: () => {
          setUid("");
          setReason("");
        },
      },
    );
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="font-semibold text-2xl">Risk</h1>
        <p className="text-muted-foreground text-sm">
          Flagged principals are denied at the mesh authorization layer:
          restrict blocks writes, block denies everything.
        </p>
      </div>

      <div className="flex flex-wrap items-end gap-2 rounded-lg border p-4">
        <div className="w-80 space-y-1">
          <label className="font-medium text-sm" htmlFor="risk-uid">
            User uid
          </label>
          <Input
            id="risk-uid"
            placeholder="00000000-0000-0000-0000-000000000000"
            value={uid}
            onChange={(e) => setUid(e.target.value)}
          />
        </div>
        <div className="space-y-1">
          <label className="font-medium text-sm" htmlFor="risk-level">
            Level
          </label>
          <select
            id="risk-level"
            className="flex h-9 w-28 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-xs"
            value={level}
            onChange={(e) => setLevel(e.target.value as RiskLevel)}
          >
            <option value="restrict">restrict</option>
            <option value="block">block</option>
          </select>
        </div>
        <div className="w-64 space-y-1">
          <label className="font-medium text-sm" htmlFor="risk-reason">
            Reason
          </label>
          <Input
            id="risk-reason"
            placeholder="why this principal is flagged"
            value={reason}
            onChange={(e) => setReason(e.target.value)}
          />
        </div>
        <div className="w-24 space-y-1">
          <label className="font-medium text-sm" htmlFor="risk-ttl">
            TTL (days)
          </label>
          <Input
            id="risk-ttl"
            type="number"
            min="0"
            max="30"
            step="0.5"
            value={ttlDays}
            onChange={(e) => setTtlDays(e.target.value)}
          />
        </div>
        <Button disabled={!canSubmit} onClick={submit}>
          Flag
        </Button>
        {flagRisk.isError ? (
          <p className="text-destructive text-sm">
            {flagRisk.error instanceof Error
              ? flagRisk.error.message
              : "flag failed"}
          </p>
        ) : null}
      </div>

      <DataTable
        columns={columns}
        rows={data.entries}
        sorting={sorting}
        onSortingChange={setSorting}
        emptyMessage="No flagged principals."
      />
    </div>
  );
}
