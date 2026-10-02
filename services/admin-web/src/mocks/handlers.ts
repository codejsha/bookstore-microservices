import { HttpResponse, http } from "msw";
import { paths } from "@/shared/api/paths";
import {
  authors,
  dashboard,
  identity,
  orders,
  payments,
  refunds,
  riskEntries,
  settlements,
  stocks,
  subjects,
  users,
  warehouses,
  works,
} from "./data";

const store = [...works];
const riskStore = [...riskEntries];

function compareValues(left: unknown, right: unknown): number {
  if (left === right) return 0;
  if (left === undefined || left === null) return 1;
  if (right === undefined || right === null) return -1;
  if (typeof left === "number" && typeof right === "number") {
    return left - right;
  }
  return String(left).localeCompare(String(right));
}

function pageOf<T>(
  rows: T[],
  url: URL,
  sortable: Record<string, keyof T> = {},
) {
  const page = Number(url.searchParams.get("page") ?? 0);
  const size = Number(url.searchParams.get("size") ?? 20);
  const [field, direction] = (url.searchParams.get("sort") ?? "").split(",");
  const key = sortable[field];
  const sorted = key
    ? [...rows].sort(
        (a, b) =>
          compareValues(a[key], b[key]) * (direction === "desc" ? -1 : 1),
      )
    : rows;
  return {
    total: rows.length,
    items: sorted.slice(page * size, page * size + size),
  };
}

function sameNames<T>(...keys: (keyof T & string)[]): Record<string, keyof T> {
  return Object.fromEntries(keys.map((k) => [k, k]));
}

const contains = (value: string | undefined, query: string | null) =>
  !query || (value ?? "").toLowerCase().includes(query.toLowerCase());

const matches = (value: string | undefined, query: string | null) =>
  !query || value === query;

function notFound(entity: string) {
  return HttpResponse.json(
    {
      type: "about:blank",
      title: "Not Found",
      status: 404,
      detail: `${entity} not found`,
    },
    { status: 404 },
  );
}

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

  http.get(paths.users.list, ({ request }) => {
    const url = new URL(request.url);
    const q = url.searchParams;
    const name = q.get("name");
    const rows = users.filter(
      (u) =>
        contains(u.email, q.get("email")) &&
        (contains(u.first_name, name) || contains(u.last_name, name)) &&
        contains(u.phone, q.get("phone")) &&
        matches(u.status, q.get("status")),
    );
    return HttpResponse.json(
      pageOf(rows, url, {
        ...sameNames<(typeof rows)[number]>("email", "phone", "status"),
        firstname: "first_name",
        lastname: "last_name",
        lastloginat: "last_login_at",
        createdat: "created_at",
      }),
    );
  }),
  http.get(paths.users.user(":uid"), ({ params }) => {
    const user = users.find((u) => u.uid === params.uid);
    return user ? HttpResponse.json(user) : notFound("user");
  }),

  http.get(paths.orders.list, ({ request }) => {
    const url = new URL(request.url);
    const q = url.searchParams;
    const rows = orders
      .filter(
        (o) =>
          matches(o.user_uid, q.get("user_uid")) &&
          matches(o.status, q.get("status")),
      )
      .map(({ items: _items, shipping: _shipping, ...summary }) => summary);
    return HttpResponse.json(
      pageOf(
        rows,
        url,
        sameNames<(typeof rows)[number]>(
          "order_number",
          "status",
          "total_amount",
          "created_at",
        ),
      ),
    );
  }),
  http.get(paths.orders.order(":uid"), ({ params }) => {
    const order = orders.find((o) => o.uid === params.uid);
    return order ? HttpResponse.json(order) : notFound("order");
  }),

  http.get(paths.payments.list, ({ request }) => {
    const url = new URL(request.url);
    const q = url.searchParams;
    const rows = payments.filter(
      (p) =>
        matches(p.customer_id, q.get("customer_id")) &&
        matches(p.status, q.get("status")) &&
        matches(p.connector, q.get("connector")),
    );
    return HttpResponse.json(
      pageOf(
        rows,
        url,
        sameNames<(typeof rows)[number]>(
          "payment_id",
          "customer_id",
          "status",
          "amount",
          "payment_method",
          "connector",
          "created_at",
        ),
      ),
    );
  }),
  http.get(paths.payments.payment(":uid"), ({ params }) => {
    const payment = payments.find((p) => p.uid === params.uid);
    return payment ? HttpResponse.json(payment) : notFound("payment");
  }),

  http.get(paths.refunds.list, ({ request }) => {
    const url = new URL(request.url);
    const q = url.searchParams;
    const rows = refunds.filter(
      (r) =>
        matches(r.payment_id, q.get("payment_id")) &&
        matches(r.status, q.get("status")),
    );
    return HttpResponse.json(
      pageOf(
        rows,
        url,
        sameNames<(typeof rows)[number]>(
          "refund_id",
          "payment_id",
          "status",
          "refund_type",
          "amount",
          "created_at",
        ),
      ),
    );
  }),
  http.get(paths.refunds.refund(":uid"), ({ params }) => {
    const refund = refunds.find((r) => r.uid === params.uid);
    return refund ? HttpResponse.json(refund) : notFound("refund");
  }),

  http.get(paths.settlements.list, ({ request }) => {
    const url = new URL(request.url);
    const q = url.searchParams;
    const rows = settlements
      .filter(
        (s) =>
          matches(s.settlement_date, q.get("date")) &&
          matches(s.status, q.get("status")),
      )
      .map(({ details: _details, ...summary }) => summary);
    return HttpResponse.json(
      pageOf(
        rows,
        url,
        sameNames<(typeof rows)[number]>(
          "settlement_date",
          "currency",
          "payment_method",
          "gross_amount",
          "refund_amount",
          "fee_amount",
          "net_amount",
          "payment_count",
          "refund_count",
          "status",
        ),
      ),
    );
  }),
  http.get(paths.settlements.settlement(":uid"), ({ params }) => {
    const settlement = settlements.find((s) => s.uid === params.uid);
    return settlement ? HttpResponse.json(settlement) : notFound("settlement");
  }),

  http.get(paths.inventory.warehouses, ({ request }) => {
    const url = new URL(request.url);
    const rows = warehouses.filter((w) =>
      contains(w.name, url.searchParams.get("name")),
    );
    return HttpResponse.json(
      pageOf(
        rows,
        url,
        sameNames<(typeof rows)[number]>("name", "address", "capacity"),
      ),
    );
  }),
  http.get(paths.inventory.warehouse(":uid"), ({ params }) => {
    const warehouse = warehouses.find((w) => w.uid === params.uid);
    return warehouse ? HttpResponse.json(warehouse) : notFound("warehouse");
  }),

  http.get(paths.inventory.stocks, ({ request }) => {
    const url = new URL(request.url);
    const editionUid = url.searchParams.get("edition_uid");
    const warehouseUid = url.searchParams.get("warehouse_uid");
    const rows = stocks
      .filter((s) => matches(s.edition_uid, editionUid))
      .map((s) => {
        const held = s.warehouses.filter((w) =>
          matches(w.warehouse_uid, warehouseUid),
        );
        return {
          ...s,
          warehouses: held,
          total_quantity: held.reduce((sum, w) => sum + w.quantity, 0),
        };
      })
      .filter((s) => s.warehouses.length > 0);
    return HttpResponse.json(pageOf(rows, url));
  }),
];
