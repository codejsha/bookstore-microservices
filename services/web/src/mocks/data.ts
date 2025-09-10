// ─── Catalog (bibliographic: Work / Edition / Subject / Author / Publisher) ──

const subjects = [
  { uid: "sub-1", name: "Fiction" },
  { uid: "sub-2", name: "Non-Fiction" },
  { uid: "sub-3", name: "Science" },
  { uid: "sub-4", name: "Technology" },
  { uid: "sub-5", name: "History" },
];

const publishers = [
  { uid: "pub-1", name: "O'Reilly Media", address: "Sebastopol, CA" },
  { uid: "pub-2", name: "Penguin Books", address: "London, UK" },
  { uid: "pub-3", name: "Manning Publications", address: "Shelter Island, NY" },
];

const authors = [
  { uid: "auth-1", name: "Martin Fowler", bio: "Software architecture expert" },
  { uid: "auth-2", name: "Robert C. Martin", bio: "Clean Code author" },
  { uid: "auth-3", name: "Eric Evans", bio: "Domain-Driven Design pioneer" },
  { uid: "auth-4", name: "Sam Newman", bio: "Microservices specialist" },
  { uid: "auth-5", name: "Vaughn Vernon", bio: "DDD practitioner" },
];

const WORK_TITLES = [
  "Clean Architecture",
  "Domain-Driven Design",
  "Building Microservices",
  "Designing Data-Intensive Applications",
  "Refactoring",
  "Patterns of Enterprise Application Architecture",
  "The Pragmatic Programmer",
  "Clean Code",
  "Release It!",
  "Implementing Domain-Driven Design",
  "Microservices Patterns",
  "Software Architecture in Practice",
  "Continuous Delivery",
  "Site Reliability Engineering",
  "The Art of Scalability",
  "System Design Interview",
  "Database Internals",
  "Fundamentals of Software Architecture",
  "Head First Design Patterns",
  "Test Driven Development",
  "Working Effectively with Legacy Code",
  "Accelerate",
  "Team Topologies",
  "A Philosophy of Software Design",
];

const works = WORK_TITLES.map((title, i) => ({
  uid: `work-${i + 1}`,
  title,
  ol_key: `/works/OL${1000 + i}W`,
  description: `A comprehensive guide covering essential concepts for modern software development. (${title})`,
  first_publish_date: String(2005 + (i % 18)),
  cover_uids: [],
  subjects: [subjects[i % subjects.length]],
  authors: [
    authors[i % authors.length],
    ...(i % 3 === 0 ? [authors[(i + 1) % authors.length]] : []),
  ],
  created_at: "2024-01-15T10:00:00Z",
  updated_at: "2024-06-01T14:30:00Z",
}));

const FORMATS = ["Hardcover", "Paperback"];

const editions = works.flatMap((work, i) =>
  FORMATS.map((format, j) => {
    const n = i * FORMATS.length + j + 1;
    return {
      uid: `edition-${n}`,
      title: `${work.title} (${format})`,
      isbn13: `978-0-13-${String(235000 + n).padStart(6, "0")}`,
      isbn10: `0-13-${String(235000 + n).padStart(6, "0")}`,
      physical_format: format,
      languages: ["English"],
      number_of_pages: 280 + ((n * 13) % 360),
      ol_key: `/books/OL${20000 + n}M`,
      publish_date: `${work.first_publish_date}`,
      publisher: publishers[n % publishers.length],
      description: work.description,
      work,
      created_at: "2024-01-15T10:00:00Z",
      updated_at: "2024-06-01T14:30:00Z",
    };
  }),
);

const orders = [
  {
    uid: "ord-1",
    user_id: 1,
    order_number: "ORD-2024-0001",
    status: "PENDING",
    currency: "USD",
    items_amount: 45.98,
    discount_amount: 0,
    shipping_amount: 5.99,
    tax_amount: 4.14,
    total_amount: 56.11,
    idempotency_key: "idem-1",
    items: [
      {
        uid: "oi-1",
        product_id: 1,
        product_name: "Clean Architecture",
        quantity: 1,
        currency: "USD",
        price: 25.99,
        tax_rate: 0.09,
        subtotal: 25.99,
        created_at: "2024-06-01T10:00:00Z",
      },
      {
        uid: "oi-2",
        product_id: 2,
        product_name: "Domain-Driven Design",
        quantity: 1,
        currency: "USD",
        price: 19.99,
        tax_rate: 0.09,
        subtotal: 19.99,
        created_at: "2024-06-01T10:00:00Z",
      },
    ],
    shipping: {
      uid: "sh-1",
      recipient_name: "John Doe",
      recipient_phone: "+1-555-0100",
      address_line1: "123 Main St",
      city: "San Francisco",
      state: "CA",
      postal_code: "94102",
      country: "US",
      shipping_method: "Standard",
      created_at: "2024-06-01T10:00:00Z",
    },
    created_at: "2024-06-01T10:00:00Z",
    updated_at: "2024-06-01T10:00:00Z",
  },
  {
    uid: "ord-2",
    user_id: 1,
    order_number: "ORD-2024-0002",
    status: "DELIVERED",
    currency: "USD",
    items_amount: 25.99,
    discount_amount: 5.0,
    shipping_amount: 0,
    tax_amount: 1.89,
    total_amount: 22.88,
    idempotency_key: "idem-2",
    items: [
      {
        uid: "oi-3",
        product_id: 3,
        product_name: "Building Microservices",
        quantity: 1,
        currency: "USD",
        price: 25.99,
        tax_rate: 0.09,
        subtotal: 25.99,
        created_at: "2024-05-15T10:00:00Z",
      },
    ],
    created_at: "2024-05-15T10:00:00Z",
    updated_at: "2024-05-20T14:30:00Z",
  },
  {
    uid: "ord-3",
    user_id: 1,
    order_number: "ORD-2024-0003",
    status: "CANCELLED",
    currency: "USD",
    items_amount: 75.49,
    discount_amount: 0,
    shipping_amount: 5.99,
    tax_amount: 6.79,
    total_amount: 88.27,
    idempotency_key: "idem-3",
    items: [
      {
        uid: "oi-4",
        product_id: 4,
        product_name: "Designing Data-Intensive Applications",
        quantity: 1,
        currency: "USD",
        price: 41.49,
        tax_rate: 0.09,
        subtotal: 41.49,
        created_at: "2024-04-10T10:00:00Z",
      },
      {
        uid: "oi-5",
        product_id: 8,
        product_name: "Clean Code",
        quantity: 1,
        currency: "USD",
        price: 34.0,
        tax_rate: 0.09,
        subtotal: 34.0,
        created_at: "2024-04-10T10:00:00Z",
      },
    ],
    created_at: "2024-04-10T10:00:00Z",
    updated_at: "2024-04-11T09:00:00Z",
  },
];

const reviews = [
  {
    uid: "rev-1",
    user_uid: "user-1",
    book_uid: "work-1",
    rating: 5,
    title: "A must-read for every developer",
    content:
      "This book completely changed how I think about software architecture. The examples are practical and the writing is clear.",
    created_at: "2024-03-15T10:00:00Z",
  },
  {
    uid: "rev-2",
    user_uid: "user-2",
    book_uid: "work-1",
    rating: 4,
    title: "Great concepts, dense read",
    content:
      "Excellent coverage of clean architecture principles. Some chapters are quite dense but well worth the effort.",
    created_at: "2024-04-20T14:30:00Z",
  },
  {
    uid: "rev-3",
    user_uid: "user-3",
    book_uid: "work-1",
    rating: 5,
    title: "Essential reading",
    content:
      "Should be required reading for any software engineer working on large-scale systems.",
    created_at: "2024-05-10T09:00:00Z",
  },
  {
    uid: "rev-4",
    user_uid: "user-1",
    book_uid: "work-2",
    rating: 5,
    title: "The DDD bible",
    content:
      "Eric Evans lays out the foundation for domain-driven design beautifully. A classic that holds up.",
    created_at: "2024-02-20T10:00:00Z",
  },
  {
    uid: "rev-5",
    user_uid: "user-4",
    book_uid: "work-2",
    rating: 3,
    title: "Good but dated",
    content:
      "The core ideas are excellent but some examples feel dated. Still worth reading for the concepts.",
    created_at: "2024-06-01T16:00:00Z",
  },
  {
    uid: "rev-6",
    user_uid: "user-2",
    book_uid: "work-3",
    rating: 4,
    title: "Practical microservices guide",
    content:
      "Sam Newman delivers a comprehensive guide to building microservices. Very practical advice.",
    created_at: "2024-01-10T11:00:00Z",
  },
  {
    uid: "rev-7",
    user_uid: "user-5",
    book_uid: "work-4",
    rating: 5,
    title: "Masterpiece",
    content:
      "Martin Kleppmann's book is the best technical book I've ever read. Covers everything from databases to distributed systems.",
    created_at: "2024-03-25T08:00:00Z",
  },
  {
    uid: "rev-8",
    user_uid: "user-3",
    book_uid: "work-8",
    rating: 4,
    title: "Clean and clear",
    content:
      "Uncle Bob's writing is engaging and the principles are timeless. A few chapters could use updating.",
    created_at: "2024-05-05T13:00:00Z",
  },
];

const cartItems = [
  {
    uid: "ci-1",
    product_id: 1,
    product_name: "Clean Architecture",
    quantity: 1,
    currency: "USD",
    price: 19.99,
  },
  {
    uid: "ci-2",
    product_id: 4,
    product_name: "Designing Data-Intensive Applications",
    quantity: 2,
    currency: "USD",
    price: 41.49,
  },
];

const wishlist = {
  book_uids: ["work-2", "work-5", "work-9", "work-16"] as string[],
};

const customer = {
  uid: "user-1",
  email: "john.doe@example.com",
  first_name: "John",
  last_name: "Doe",
  phone: "+1-555-0100",
  roles: ["PROFILE", "ORDER", "VIEW"],
};

const points = {
  user_uid: "user-1",
  balance: 1240,
};

const pointHistory = [
  {
    uid: "ph-1",
    amount: 100,
    change_type: "EARN",
    reason: "Welcome bonus",
    created_at: "2024-03-01T10:00:00Z",
  },
  {
    uid: "ph-2",
    amount: 56,
    change_type: "EARN",
    reason: "Order ORD-2024-0001 reward",
    created_at: "2024-06-01T10:00:00Z",
  },
  {
    uid: "ph-3",
    amount: -200,
    change_type: "SPEND",
    reason: "Discount applied to order ORD-2024-0002",
    created_at: "2024-05-15T10:00:00Z",
  },
  {
    uid: "ph-4",
    amount: 23,
    change_type: "EARN",
    reason: "Order ORD-2024-0002 reward",
    created_at: "2024-05-15T10:00:00Z",
  },
  {
    uid: "ph-5",
    amount: 88,
    change_type: "EARN",
    reason: "Order ORD-2024-0003 reward",
    created_at: "2024-04-10T10:00:00Z",
  },
  {
    uid: "ph-6",
    amount: 1173,
    change_type: "EARN",
    reason: "Promotional credit",
    created_at: "2024-02-12T10:00:00Z",
  },
];

export {
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
};
