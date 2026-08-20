export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/works?size=20',
  '/api/v1/works/search?q=ring',
  '/api/v1/editions?size=20',
  '/api/v1/authors?size=20',
  '/api/v1/publishers?size=20',
  '/api/v1/subjects?size=20',
];
