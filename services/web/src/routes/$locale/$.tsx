import { createFileRoute } from "@tanstack/react-router";
import { NotFoundState } from "@/shared/components/NotFoundState";

export const Route = createFileRoute("/$locale/$")({
  component: NotFoundState,
});
