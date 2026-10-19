import { createFileRoute } from "@tanstack/react-router";
import { WorkDetail, workDetailQueryOptions } from "@/domains/catalog";
import { NotFoundState } from "@/shared/components/NotFoundState";

export const Route = createFileRoute("/books/$uid")({
  loader: ({ context: { queryClient }, params: { uid } }) =>
    queryClient.ensureQueryData(workDetailQueryOptions(uid)),
  component: WorkDetailPage,
  notFoundComponent: NotFoundState,
});

function WorkDetailPage() {
  const { uid } = Route.useParams();
  return <WorkDetail uid={uid} />;
}
