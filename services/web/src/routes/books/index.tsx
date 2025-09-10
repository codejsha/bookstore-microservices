import { createFileRoute } from "@tanstack/react-router";
import { WorkList } from "@/domains/catalog";

export const Route = createFileRoute("/books/")({
  component: BooksPage,
});

function BooksPage() {
  return <WorkList />;
}
