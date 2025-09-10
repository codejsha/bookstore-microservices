import { HttpResponse, http } from "msw";
import {
  authors,
  cartItems,
  customer,
  editions,
  orders,
  pointHistory,
  points,
  publishers,
  reviews,
  subjects,
  wishlist,
  works,
} from "./data";

function paginate<T>(items: T[], size = 12, page = 0) {
  const start = page * size;
  return {
    total: items.length,
    items: items.slice(start, start + size),
  };
}

export const handlers = [
  // ─── Catalog: Works ──────────────────────────────────────────────────
  http.get("/api/v1/works/search", ({ request }) => {
    const url = new URL(request.url);
    const q = url.searchParams.get("q");
    const authorUid = url.searchParams.get("authorUid");
    const subjectUid = url.searchParams.get("subjectUid");
    const size = Number(url.searchParams.get("size")) || 12;
    const page = Number(url.searchParams.get("page")) || 0;

    let filtered = [...works];
    if (q) {
      const needle = q.toLowerCase();
      filtered = filtered.filter(
        (w) =>
          w.title.toLowerCase().includes(needle) ||
          w.description.toLowerCase().includes(needle),
      );
    }
    if (authorUid)
      filtered = filtered.filter((w) =>
        w.authors.some((a) => a.uid === authorUid),
      );
    if (subjectUid)
      filtered = filtered.filter((w) =>
        w.subjects.some((s) => s.uid === subjectUid),
      );

    const start = page * size;
    const items = filtered.slice(start, start + size).map((w) => ({
      ...w,
      score: 1,
    }));
    return HttpResponse.json({ total: filtered.length, items });
  }),

  http.get("/api/v1/works/:uid/editions", ({ params, request }) => {
    const url = new URL(request.url);
    const size = Number(url.searchParams.get("size")) || 20;
    const page = Number(url.searchParams.get("page")) || 0;
    const matched = editions.filter((e) => e.work.uid === params.uid);
    return HttpResponse.json(paginate(matched, size, page));
  }),

  http.get("/api/v1/works", ({ request }) => {
    const url = new URL(request.url);
    const title = url.searchParams.get("title");
    const authorUid = url.searchParams.get("authorUid");
    const subjectUid = url.searchParams.get("subjectUid");
    const olKey = url.searchParams.get("olKey");
    const size = Number(url.searchParams.get("size")) || 12;
    const page = Number(url.searchParams.get("page")) || 0;

    let filtered = [...works];
    if (title)
      filtered = filtered.filter((w) =>
        w.title.toLowerCase().includes(title.toLowerCase()),
      );
    if (authorUid)
      filtered = filtered.filter((w) =>
        w.authors.some((a) => a.uid === authorUid),
      );
    if (subjectUid)
      filtered = filtered.filter((w) =>
        w.subjects.some((s) => s.uid === subjectUid),
      );
    if (olKey) filtered = filtered.filter((w) => w.ol_key === olKey);

    return HttpResponse.json(paginate(filtered, size, page));
  }),

  http.get("/api/v1/works/:uid", ({ params }) => {
    const work = works.find((w) => w.uid === params.uid);
    if (!work)
      return HttpResponse.json(
        { code: 404, message: "Work not found" },
        { status: 404 },
      );
    return HttpResponse.json(work);
  }),

  // ─── Catalog: Editions ───────────────────────────────────────────────
  http.get("/api/v1/editions", ({ request }) => {
    const url = new URL(request.url);
    const title = url.searchParams.get("title");
    const isbn = url.searchParams.get("isbn");
    const workUid = url.searchParams.get("workUid");
    const publisherUid = url.searchParams.get("publisherUid");
    const size = Number(url.searchParams.get("size")) || 12;
    const page = Number(url.searchParams.get("page")) || 0;

    let filtered = [...editions];
    if (title)
      filtered = filtered.filter((e) =>
        e.title.toLowerCase().includes(title.toLowerCase()),
      );
    if (isbn)
      filtered = filtered.filter((e) => e.isbn13 === isbn || e.isbn10 === isbn);
    if (workUid) filtered = filtered.filter((e) => e.work.uid === workUid);
    if (publisherUid)
      filtered = filtered.filter((e) => e.publisher?.uid === publisherUid);

    return HttpResponse.json(paginate(filtered, size, page));
  }),

  http.get("/api/v1/editions/:uid", ({ params }) => {
    const edition = editions.find((e) => e.uid === params.uid);
    if (!edition)
      return HttpResponse.json(
        { code: 404, message: "Edition not found" },
        { status: 404 },
      );
    return HttpResponse.json(edition);
  }),

  // ─── Catalog: Subjects / Authors / Publishers ────────────────────────
  http.get("/api/v1/subjects", () => {
    return HttpResponse.json({ total: subjects.length, items: subjects });
  }),

  http.get("/api/v1/authors", () => {
    return HttpResponse.json({ total: authors.length, items: authors });
  }),

  http.get("/api/v1/publishers", () => {
    return HttpResponse.json({ total: publishers.length, items: publishers });
  }),

  // ─── Orders ─────────────────────────────────────────────────────────
  http.get("/api/v1/orders", ({ request }) => {
    const url = new URL(request.url);
    const status = url.searchParams.get("status");
    const size = Number(url.searchParams.get("size")) || 10;
    const page = Number(url.searchParams.get("page")) || 0;

    let filtered = [...orders];
    if (status) filtered = filtered.filter((o) => o.status === status);

    return HttpResponse.json(paginate(filtered, size, page));
  }),

  http.get("/api/v1/orders/:uid", ({ params }) => {
    const order = orders.find((o) => o.uid === params.uid);
    if (!order)
      return HttpResponse.json(
        { code: 404, message: "Order not found" },
        { status: 404 },
      );
    return HttpResponse.json(order);
  }),

  http.post("/api/v1/orders/:uid/cancel", ({ params }) => {
    const order = orders.find((o) => o.uid === params.uid);
    if (!order)
      return HttpResponse.json(
        { code: 404, message: "Order not found" },
        { status: 404 },
      );
    if (order.status !== "PENDING") {
      return HttpResponse.json(
        { code: 400, message: "Only PENDING orders can be cancelled" },
        { status: 400 },
      );
    }
    order.status = "CANCELLED";
    return HttpResponse.json(order);
  }),

  // ─── Reviews (keyed by book uid) ─────────────────────────────────────
  http.get("/api/v1/books/:bookUid/reviews", ({ params, request }) => {
    const bookUid = params.bookUid as string;
    const url = new URL(request.url);
    const ratingFilter = url.searchParams.get("rating");

    let filtered = reviews.filter((r) => r.book_uid === bookUid);
    if (ratingFilter)
      filtered = filtered.filter((r) => r.rating === Number(ratingFilter));

    return HttpResponse.json({ total: filtered.length, items: filtered });
  }),

  // ─── Cart ───────────────────────────────────────────────────────────
  http.get("/api/v1/cart", () => {
    const total_items = cartItems.reduce((sum, i) => sum + i.quantity, 0);
    const total_amount = Number(
      cartItems.reduce((sum, i) => sum + i.price * i.quantity, 0).toFixed(2),
    );
    return HttpResponse.json({
      uid: "cart-1",
      user_id: 1,
      items: cartItems.map((i) => ({
        ...i,
        subtotal: Number((i.price * i.quantity).toFixed(2)),
      })),
      total_items,
      total_amount,
    });
  }),

  http.post("/api/v1/cart/items", async ({ request }) => {
    const body = (await request.json()) as {
      product_id: number;
      product_name?: string;
      quantity: number;
      currency: string;
      price: number;
    };
    const existing = cartItems.find((i) => i.product_id === body.product_id);
    if (existing) {
      existing.quantity += body.quantity;
    } else {
      cartItems.push({
        uid: `ci-${Date.now()}`,
        product_id: body.product_id,
        product_name: body.product_name ?? "",
        quantity: body.quantity,
        currency: body.currency,
        price: body.price,
      });
    }
    const total_items = cartItems.reduce((sum, i) => sum + i.quantity, 0);
    const total_amount = Number(
      cartItems.reduce((sum, i) => sum + i.price * i.quantity, 0).toFixed(2),
    );
    return HttpResponse.json({
      uid: "cart-1",
      user_id: 1,
      items: cartItems.map((i) => ({
        ...i,
        subtotal: Number((i.price * i.quantity).toFixed(2)),
      })),
      total_items,
      total_amount,
    });
  }),

  http.put("/api/v1/cart/items/:uid", async ({ params, request }) => {
    const body = (await request.json()) as { quantity: number };
    const idx = cartItems.findIndex((i) => i.uid === params.uid);
    if (idx === -1)
      return HttpResponse.json(
        { code: 404, message: "Item not found" },
        { status: 404 },
      );
    if (body.quantity <= 0) {
      cartItems.splice(idx, 1);
    } else {
      cartItems[idx].quantity = body.quantity;
    }
    const total_items = cartItems.reduce((sum, i) => sum + i.quantity, 0);
    const total_amount = Number(
      cartItems.reduce((sum, i) => sum + i.price * i.quantity, 0).toFixed(2),
    );
    return HttpResponse.json({
      uid: "cart-1",
      user_id: 1,
      items: cartItems.map((i) => ({
        ...i,
        subtotal: Number((i.price * i.quantity).toFixed(2)),
      })),
      total_items,
      total_amount,
    });
  }),

  http.delete("/api/v1/cart/items/:uid", ({ params }) => {
    const idx = cartItems.findIndex((i) => i.uid === params.uid);
    if (idx === -1)
      return HttpResponse.json(
        { code: 404, message: "Item not found" },
        { status: 404 },
      );
    cartItems.splice(idx, 1);
    return HttpResponse.json(null, { status: 204 });
  }),

  http.delete("/api/v1/cart", () => {
    cartItems.length = 0;
    return HttpResponse.json(null, { status: 204 });
  }),

  http.post("/api/v1/cart/checkout", async ({ request }) => {
    const body = (await request.json()) as {
      currency: string;
      idempotency_key: string;
    };
    const orderItems = cartItems.map((i, idx) => ({
      uid: `oi-checkout-${idx}`,
      product_id: i.product_id,
      product_name: i.product_name ?? `Product #${i.product_id}`,
      quantity: i.quantity,
      currency: i.currency,
      price: i.price,
      tax_rate: 0.09,
      subtotal: Number((i.price * i.quantity).toFixed(2)),
      created_at: new Date().toISOString(),
    }));
    const itemsAmount = Number(
      orderItems.reduce((s, i) => s + i.subtotal, 0).toFixed(2),
    );
    const newOrder = {
      uid: `ord-${Date.now()}`,
      user_id: 1,
      order_number: `ORD-${Date.now()}`,
      status: "PENDING",
      currency: body.currency,
      items_amount: itemsAmount,
      discount_amount: 0,
      shipping_amount: 5.99,
      tax_amount: Number((itemsAmount * 0.09).toFixed(2)),
      total_amount: Number(
        (itemsAmount + 5.99 + itemsAmount * 0.09).toFixed(2),
      ),
      idempotency_key: body.idempotency_key,
      items: orderItems,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    orders.push(newOrder);
    cartItems.length = 0;
    return HttpResponse.json(newOrder, { status: 201 });
  }),

  // ─── Wishlist (keyed by book uid) ────────────────────────────────────
  http.get("/api/v1/customers/:uid/wishlist", () => {
    return HttpResponse.json({ book_uids: wishlist.book_uids });
  }),

  http.post("/api/v1/customers/:uid/wishlist/add", async ({ request }) => {
    const body = (await request.json()) as { book_uids: string[] };
    for (const bookUid of body.book_uids) {
      if (!wishlist.book_uids.includes(bookUid))
        wishlist.book_uids.push(bookUid);
    }
    return HttpResponse.json({ book_uids: wishlist.book_uids });
  }),

  http.post("/api/v1/customers/:uid/wishlist/remove", async ({ request }) => {
    const body = (await request.json()) as { book_uids: string[] };
    wishlist.book_uids = wishlist.book_uids.filter(
      (bookUid) => !body.book_uids.includes(bookUid),
    );
    return HttpResponse.json({ book_uids: wishlist.book_uids });
  }),

  // ─── Customer profile ──────────────────────────────────────────────
  http.get("/api/v1/customers/:uid", () => {
    return HttpResponse.json(customer);
  }),

  http.put("/api/v1/customers/:uid", async ({ request }) => {
    const body = (await request.json()) as Partial<typeof customer>;
    Object.assign(customer, body);
    return HttpResponse.json(customer);
  }),

  // ─── Points ─────────────────────────────────────────────────────────
  http.get("/api/v1/customers/:uid/points", () => {
    return HttpResponse.json({
      user_uid: points.user_uid,
      balance: points.balance,
    });
  }),

  http.get("/api/v1/customers/:uid/points/history", ({ request }) => {
    const url = new URL(request.url);
    const size = Number(url.searchParams.get("size")) || 20;
    const page = Number(url.searchParams.get("page")) || 0;
    const start = page * size;
    return HttpResponse.json({
      total: pointHistory.length,
      items: pointHistory.slice(start, start + size),
    });
  }),

  http.post("/api/v1/customers/:uid/points/spend", async ({ request }) => {
    const body = (await request.json()) as { amount: number; reason?: string };
    points.balance = Math.max(0, points.balance - body.amount);
    pointHistory.unshift({
      uid: `ph-${Date.now()}`,
      amount: -body.amount,
      change_type: "SPEND",
      reason: body.reason ?? "Spent points",
      created_at: new Date().toISOString(),
    });
    return HttpResponse.json({
      user_uid: points.user_uid,
      balance: points.balance,
    });
  }),

  http.post("/api/v1/customers/:uid/points/earn", async ({ request }) => {
    const body = (await request.json()) as { amount: number; reason?: string };
    points.balance += body.amount;
    pointHistory.unshift({
      uid: `ph-${Date.now()}`,
      amount: body.amount,
      change_type: "EARN",
      reason: body.reason ?? "Earned points",
      created_at: new Date().toISOString(),
    });
    return HttpResponse.json({
      user_uid: points.user_uid,
      balance: points.balance,
    });
  }),

  http.post("/api/v1/customers/:uid/reviews", async ({ params, request }) => {
    const body = (await request.json()) as {
      book_uid: string;
      rating: number;
      title?: string;
      content?: string;
    };
    const userUid = params.uid as string;
    const newReview = {
      uid: `rev-${Date.now()}`,
      user_uid: userUid,
      book_uid: body.book_uid,
      rating: body.rating,
      title: body.title ?? "",
      content: body.content ?? "",
      created_at: new Date().toISOString(),
    };
    reviews.push(newReview);
    return HttpResponse.json(null, {
      status: 201,
      headers: {
        Location: `/api/v1/customers/${userUid}/reviews/${newReview.uid}`,
      },
    });
  }),
];
