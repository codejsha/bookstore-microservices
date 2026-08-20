import { HttpResponse, http } from "msw";
import { paths } from "@/shared/api/paths";
import {
  authors,
  dashboard,
  identity,
  riskEntries,
  subjects,
  works,
} from "./data";

const store = [...works];
const riskStore = [...riskEntries];

export const handlers = [
  http.get(paths.admin.me, () => HttpResponse.json(identity)),
  http.get(paths.admin.dashboard, () => HttpResponse.json(dashboard)),

  http.get(paths.risk.list, () => HttpResponse.json({ entries: riskStore })),
  http.put("/api/v1/admin/risk/:uid", async ({ params, request }) => {
    const body = (await request.json()) as {
      level: string;
      reason: string;
      ttl_seconds?: number;
    };
    const uid = String(params.uid);
    const ttlSeconds = body.ttl_seconds ?? 86_400;
    const entry = {
      user_uid: uid,
      level: body.level,
      reason: body.reason,
      flagged_by: "11111111-1111-1111-1111-111111111111",
      flagged_at: new Date().toISOString(),
      expires_at: new Date(Date.now() + ttlSeconds * 1000).toISOString(),
    };
    const idx = riskStore.findIndex((e) => e.user_uid === uid);
    if (idx >= 0) riskStore[idx] = entry;
    else riskStore.push(entry);
    return HttpResponse.json(entry);
  }),
  http.delete("/api/v1/admin/risk/:uid", ({ params }) => {
    const idx = riskStore.findIndex((e) => e.user_uid === String(params.uid));
    if (idx >= 0) riskStore.splice(idx, 1);
    return new HttpResponse(null, { status: 204 });
  }),

  http.get(paths.catalog.authors, () =>
    HttpResponse.json({ total: authors.length, items: authors }),
  ),
  http.get(paths.catalog.subjects, () =>
    HttpResponse.json({ total: subjects.length, items: subjects }),
  ),

  http.get(paths.catalog.works, ({ request }) => {
    const url = new URL(request.url);
    const title = url.searchParams.get("title")?.toLowerCase();
    const page = Number(url.searchParams.get("page") ?? 0);
    const size = Number(url.searchParams.get("size") ?? 20);
    const sort = url.searchParams.get("sort");

    let rows = title
      ? store.filter((w) => w.title.toLowerCase().includes(title))
      : store;

    if (sort) {
      const [field, direction] = sort.split(",");
      rows = [...rows].sort((a, b) => {
        const left = String(a[field as keyof typeof a] ?? "");
        const right = String(b[field as keyof typeof b] ?? "");
        return direction === "desc"
          ? right.localeCompare(left)
          : left.localeCompare(right);
      });
    }

    return HttpResponse.json({
      total: rows.length,
      items: rows.slice(page * size, page * size + size),
    });
  }),

  http.get("/api/v1/admin/catalog/works/:uid", ({ params }) => {
    const work = store.find((w) => w.uid === params.uid);
    return work
      ? HttpResponse.json(work)
      : HttpResponse.json(
          { code: 404, message: "work not found" },
          { status: 404 },
        );
  }),

  http.post(paths.catalog.works, () => new HttpResponse(null, { status: 201 })),

  http.put("/api/v1/admin/catalog/works/:uid", async ({ params, request }) => {
    const index = store.findIndex((w) => w.uid === params.uid);
    if (index < 0) {
      return HttpResponse.json(
        { code: 404, message: "work not found" },
        { status: 404 },
      );
    }
    const patch = (await request.json()) as Record<string, unknown>;
    const updated = {
      ...store[index],
      ...Object.fromEntries(
        Object.entries(patch).filter(([, v]) => v !== undefined),
      ),
      updated_at: "2026-07-10T00:00:00Z",
    } as (typeof store)[number];
    store[index] = updated;
    return HttpResponse.json(updated);
  }),
];
