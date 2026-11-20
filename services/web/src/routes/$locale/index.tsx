import { createFileRoute } from "@tanstack/react-router";
import { HeroSection } from "@/domains/catalog";

export const Route = createFileRoute("/$locale/")({
  component: HomePage,
});

function HomePage() {
  return <HeroSection />;
}
