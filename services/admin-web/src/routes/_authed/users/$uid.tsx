import { Badge } from "@bookstore/design/ui/badge";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { userQueryOptions } from "@/domains/user";
import {
  DefinitionList,
  DetailError,
  DetailHeader,
  DetailNotFound,
} from "@/shared/components/detail";
import { formatDateTime, orEmpty } from "@/shared/lib/format";

export const Route = createFileRoute("/_authed/users/$uid")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(userQueryOptions(params.uid)),
  errorComponent: ({ error }) => <DetailError error={error} entity="user" />,
  notFoundComponent: () => <DetailNotFound entity="user" />,
  component: UserDetailPage,
});

function UserDetailPage() {
  const { uid } = Route.useParams();
  const { data: user } = useSuspenseQuery(userQueryOptions(uid));

  return (
    <div className="flex flex-col gap-6">
      <DetailHeader
        back={
          <Link to="/users" search={{ page: 0 }}>
            Back to users
          </Link>
        }
        title={`${user.first_name} ${user.last_name}`.trim() || user.email}
        subtitle={user.uid}
        actions={
          <Link
            to="/orders"
            search={{ page: 0, user_uid: user.uid }}
            className="underline"
          >
            View orders
          </Link>
        }
      />

      <DefinitionList
        items={[
          { label: "Email", value: user.email },
          { label: "First name", value: orEmpty(user.first_name) },
          { label: "Last name", value: orEmpty(user.last_name) },
          { label: "Phone", value: orEmpty(user.phone) },
          {
            label: "Status",
            value: <Badge variant="outline">{user.status}</Badge>,
          },
          { label: "Roles", value: orEmpty(user.roles.join(", ")) },
          { label: "Last login", value: formatDateTime(user.last_login_at) },
          { label: "Created", value: formatDateTime(user.created_at) },
          { label: "Updated", value: formatDateTime(user.updated_at) },
        ]}
      />
    </div>
  );
}
