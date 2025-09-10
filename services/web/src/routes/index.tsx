import { createFileRoute } from "@tanstack/react-router";
import { HeroSection } from "@/domains/catalog";

export const Route = createFileRoute("/")({
  component: HomePage,
});

function HomePage() {
  return <HeroSection />;
}
