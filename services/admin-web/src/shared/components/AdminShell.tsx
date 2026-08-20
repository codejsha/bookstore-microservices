import { Separator } from "@bookstore/design/ui/separator";
import { Link } from "@tanstack/react-router";
import { BookOpen, LayoutDashboard, ShieldAlert } from "lucide-react";
import type { ReactNode } from "react";
import type { AdminIdentity } from "@/domains/admin";

interface NavItem {
  to: string;
  label: string;
  icon: typeof LayoutDashboard;
}

const NAV_ITEMS: NavItem[] = [
  { to: "/", label: "Dashboard", icon: LayoutDashboard },
  { to: "/catalog", label: "Catalog", icon: BookOpen },
  { to: "/risk", label: "Risk", icon: ShieldAlert },
];

interface AdminShellProps {
  identity: AdminIdentity;
  children: ReactNode;
}

export function AdminShell({ identity, children }: AdminShellProps) {
  return (
    <div className="flex min-h-screen">
      <nav className="flex w-56 shrink-0 flex-col gap-1 border-border border-r p-4">
        <div className="px-2 pb-2">
          <p className="font-semibold text-sm">Bookstore</p>
          <p className="text-muted-foreground text-xs">Admin console</p>
        </div>
        <Separator />
        <ul className="flex flex-col gap-1 pt-2">
          {NAV_ITEMS.map((item) => (
            <li key={item.to}>
              <Link
                to={item.to}
                activeProps={{ className: "bg-accent text-accent-foreground" }}
                className="flex items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground"
              >
                <item.icon className="size-4" aria-hidden />
                {item.label}
              </Link>
            </li>
          ))}
        </ul>
        <div className="mt-auto px-2 pt-4">
          <p className="truncate text-muted-foreground text-xs">
            {identity.email ?? identity.name ?? identity.uid}
          </p>
          <a
            href="/oauth2/sign_out"
            className="text-muted-foreground text-xs underline hover:text-foreground"
          >
            Sign out
          </a>
        </div>
      </nav>
      <main className="flex-1 p-6">{children}</main>
    </div>
  );
}
