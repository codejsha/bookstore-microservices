import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { warehouseQueryOptions } from "@/domains/inventory";
import {
  DefinitionList,
  DetailError,
  DetailHeader,
  DetailNotFound,
} from "@/shared/components/detail";
import { formatDateTime, orEmpty } from "@/shared/lib/format";

export const Route = createFileRoute("/_authed/inventory/warehouses/$uid")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(warehouseQueryOptions(params.uid)),
  errorComponent: ({ error }) => (
    <DetailError error={error} entity="warehouse" />
  ),
  notFoundComponent: () => <DetailNotFound entity="warehouse" />,
  component: WarehouseDetailPage,
});

function WarehouseDetailPage() {
  const { uid } = Route.useParams();
  const { data: warehouse } = useSuspenseQuery(warehouseQueryOptions(uid));

  return (
    <div className="flex flex-col gap-6">
      <DetailHeader
        back={
          <Link to="/inventory/warehouses" search={{ page: 0 }}>
            Back to warehouses
          </Link>
        }
        title={warehouse.name}
        subtitle={warehouse.uid}
        actions={
          <Link
            to="/inventory/stocks"
            search={{ page: 0, warehouse_uid: warehouse.uid }}
            className="underline"
          >
            View stock
          </Link>
        }
      />

      <DefinitionList
        items={[
          { label: "Name", value: warehouse.name },
          { label: "Address", value: orEmpty(warehouse.address) },
          { label: "Capacity", value: warehouse.capacity.toLocaleString() },
          { label: "Created", value: formatDateTime(warehouse.created_at) },
          { label: "Updated", value: formatDateTime(warehouse.updated_at) },
        ]}
      />
    </div>
  );
}
