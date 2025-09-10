import { createFileRoute } from "@tanstack/react-router";
import { WorkDetail, workDetailQueryOptions } from "@/domains/catalog";

export const Route = createFileRoute("/books/$uid")({
  loader: ({ context: { queryClient }, params: { uid } }) =>
    queryClient.ensureQueryData(workDetailQueryOptions(uid)),
  component: WorkDetailPage,
});

function WorkDetailPage() {
  const { uid } = Route.useParams();
  return <WorkDetail uid={uid} />;
}
