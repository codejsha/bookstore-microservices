import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import {
  authorsQueryOptions,
  toCreateRequest,
  useCreateWork,
  WorkForm,
} from "@/domains/catalog";
import { ErrorState } from "@/shared/components/ErrorState";

export const Route = createFileRoute("/_authed/catalog/new")({
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(authorsQueryOptions()),
  component: NewWorkPage,
});

function NewWorkPage() {
  const navigate = useNavigate();
  const { data: authors } = useSuspenseQuery(authorsQueryOptions());
  const createWork = useCreateWork();

  return (
    <div className="flex flex-col gap-6">
      <h1 className="font-semibold text-2xl">New work</h1>

      {createWork.isError ? (
        <ErrorState
          title="Could not create this work"
          message={createWork.error.message}
          onRetry={createWork.reset}
        />
      ) : null}

      <WorkForm
        authors={authors.items}
        submitLabel="Create work"
        pending={createWork.isPending}
        onSubmit={(values) =>
          createWork.mutate(toCreateRequest(values), {
            onSuccess: () => navigate({ to: "/catalog", search: { page: 0 } }),
          })
        }
        onCancel={() => navigate({ to: "/catalog", search: { page: 0 } })}
      />
    </div>
  );
}
