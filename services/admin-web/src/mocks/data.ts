import type { AdminIdentity, Dashboard } from "@/domains/admin";
import type { Author, Subject, Work } from "@/domains/catalog";

export const identity: AdminIdentity = {
  uid: "11111111-1111-1111-1111-111111111111",
  email: "admin@example.com",
  name: "root",
  roles: ["ADMIN", "VIEW"],
};

export const dashboard: Dashboard = {
  works: { count: "1284", available: true },
  orders: { count: "0", available: false },
  warehouses: { count: "3", available: true },
};

export const authors: Author[] = [
  { uid: "33333333-3333-3333-3333-333333333331", name: "Frank Herbert" },
  { uid: "33333333-3333-3333-3333-333333333332", name: "Ursula K. Le Guin" },
  { uid: "33333333-3333-3333-3333-333333333333", name: "Octavia E. Butler" },
];

export const subjects: Subject[] = [
  { uid: "44444444-4444-4444-4444-444444444441", name: "Science Fiction" },
  { uid: "44444444-4444-4444-4444-444444444442", name: "Fantasy" },
];

export const works: Work[] = Array.from({ length: 43 }, (_, i) => ({
  uid: `22222222-2222-2222-2222-${String(i).padStart(12, "0")}`,
  title: `Work ${i + 1}`,
  description: i % 3 === 0 ? `A description for work ${i + 1}.` : undefined,
  first_publish_date: String(1960 + (i % 40)),
  ol_key: i % 2 === 0 ? `OL${i}W` : undefined,
  authors: [authors[i % authors.length]],
  subjects: [subjects[i % subjects.length]],
  created_at: "2026-01-01T00:00:00Z",
  updated_at: undefined,
}));

export const riskEntries = [
  {
    user_uid: "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    level: "restrict",
    reason: "4xx burst on catalog endpoints",
    flagged_by: "11111111-1111-1111-1111-111111111111",
    flagged_at: "2026-08-19T02:10:00Z",
    expires_at: "2026-08-21T02:10:00Z",
  },
  {
    user_uid: "550e8400-e29b-41d4-a716-446655440000",
    level: "block",
    reason: "credential stuffing from shared IP",
    flagged_by: "11111111-1111-1111-1111-111111111111",
    flagged_at: "2026-08-19T09:30:00Z",
    expires_at: "2026-08-26T09:30:00Z",
  },
];
