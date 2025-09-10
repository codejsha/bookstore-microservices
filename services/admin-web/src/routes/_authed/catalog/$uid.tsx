import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import {
  authorsQueryOptions,
  toUpdateRequest,
  useUpdateWork,
  WorkForm,
  workQueryOptions,
} from "@/domains/catalog";
import { ErrorState } from "@/shared/components/ErrorState";

export const Route = createFileRoute("/_authed/catalog/$uid")({
  loader: ({ context, params }) =>
    Promise.all([
      context.queryClient.ensureQueryData(workQueryOptions(params.uid)),
      context.queryClient.ensureQueryData(authorsQueryOptions()),
    ]),
  component: EditWorkPage,
});

function EditWorkPage() {
  const { uid } = Route.useParams();
  const navigate = useNavigate();
  const { data: work } = useSuspenseQuery(workQueryOptions(uid));
  const { data: authors } = useSuspenseQuery(authorsQueryOptions());
  const updateWork = useUpdateWork(uid);

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="font-semibold text-2xl">{work.title}</h1>
        <p className="text-muted-foreground text-xs">{work.uid}</p>
      </div>

      {updateWork.isError ? (
        <ErrorState
          title="Could not save this work"
          message={updateWork.error.message}
          onRetry={updateWork.reset}
        />
      ) : null}

      <WorkForm
        authors={authors.items}
        work={work}
        submitLabel="Save changes"
        pending={updateWork.isPending}
        onSubmit={(values) => updateWork.mutate(toUpdateRequest(values))}
        onCancel={() => navigate({ to: "/catalog", search: { page: 0 } })}
      />
    </div>
  );
}
