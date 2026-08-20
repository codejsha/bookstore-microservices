export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const TICKET_UID = __ENV.TICKET_UID || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/tickets?size=20',
  '/api/v1/categories?size=20',
  '/api/v1/faqs?size=20',
  ...(TICKET_UID
    ? [
        `/api/v1/tickets/${TICKET_UID}`,
        `/api/v1/tickets/${TICKET_UID}/comments`,
      ]
    : []),
];
