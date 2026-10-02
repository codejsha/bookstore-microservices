import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { ColumnDef } from "@tanstack/react-table";
import {
  USER_STATUSES,
  type User,
  type UserStatus,
  usersQueryOptions,
} from "@/domains/user";
import {
  FilterBar,
  FilterInput,
  FilterSelect,
} from "@/shared/components/filters";
import { DataTable } from "@/shared/components/ui/data-table";
import { Pagination } from "@/shared/components/ui/pagination";
import { formatDateTime, orEmpty } from "@/shared/lib/format";
import {
  type ListSearch,
  PAGE_SIZE,
  parseListSearch,
  parseOneOf,
  parseText,
} from "@/shared/lib/list-search";
import { useUrlSorting } from "@/shared/lib/sorting";

interface UsersSearch extends ListSearch {
  email?: string;
  name?: string;
  phone?: string;
  status?: UserStatus;
}

function validateSearch(search: Record<string, unknown>): UsersSearch {
  return {
    ...parseListSearch(search),
    email: parseText(search.email),
    name: parseText(search.name),
    phone: parseText(search.phone),
    status: parseOneOf(search.status, USER_STATUSES),
  };
}

export const Route = createFileRoute("/_authed/users/")({
  validateSearch,
  loaderDeps: ({ search }) => search,
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(
      usersQueryOptions({ ...deps, size: PAGE_SIZE }),
    ),
  component: UsersPage,
});

const columns: ColumnDef<User, unknown>[] = [
  { accessorKey: "email", header: "Email" },
  {
    id: "firstname",
    accessorKey: "first_name",
    header: "First name",
  },
  {
    id: "lastname",
    accessorKey: "last_name",
    header: "Last name",
  },
  {
    accessorKey: "phone",
    header: "Phone",
    cell: ({ row }) => orEmpty(row.original.phone),
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ row }) => <Badge variant="outline">{row.original.status}</Badge>,
  },
  {
    id: "roles",
    header: "Roles",
    enableSorting: false,
    cell: ({ row }) => orEmpty(row.original.roles.join(", ")),
  },
  {
    id: "lastloginat",
    accessorKey: "last_login_at",
    header: "Last login",
    cell: ({ row }) => formatDateTime(row.original.last_login_at),
  },
  {
    id: "createdat",
    accessorKey: "created_at",
    header: "Created",
    cell: ({ row }) => formatDateTime(row.original.created_at),
  },
];

function UsersPage() {
  const search = Route.useSearch();
  const navigate = useNavigate({ from: Route.fullPath });
  const { data } = useSuspenseQuery(
    usersQueryOptions({ ...search, size: PAGE_SIZE }),
  );
  const { sorting, onSortingChange } = useUrlSorting(search.sort, (sort) =>
    navigate({ search: { ...search, sort, page: 0 } }),
  );
  const filter = (patch: Partial<UsersSearch>) =>
    navigate({ search: { ...search, ...patch, page: 0 }, replace: true });

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">Users</h1>

      <FilterBar>
        <FilterInput
          id="users-email"
          label="Email"
          value={search.email}
          onChange={(email) => filter({ email })}
        />
        <FilterInput
          id="users-name"
          label="Name"
          value={search.name}
          onChange={(name) => filter({ name })}
        />
        <FilterInput
          id="users-phone"
          label="Phone"
          value={search.phone}
          onChange={(phone) => filter({ phone })}
        />
        <FilterSelect
          id="users-status"
          label="Status"
          value={search.status}
          options={USER_STATUSES}
          allLabel="All statuses"
          onChange={(status) => filter({ status })}
        />
      </FilterBar>

      <DataTable
        columns={columns}
        rows={data.items}
        sorting={sorting}
        onSortingChange={onSortingChange}
        onRowClick={(user) =>
          navigate({ to: "/users/$uid", params: { uid: user.uid } })
        }
        emptyMessage="No users match these filters."
      />

      <Pagination
        page={search.page}
        size={PAGE_SIZE}
        total={data.total}
        onPageChange={(page) => navigate({ search: { ...search, page } })}
      />
    </div>
  );
}
