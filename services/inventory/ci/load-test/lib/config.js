export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const EDITION_UID = __ENV.EDITION_UID || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/stocks?size=20',
  '/api/v1/balance',
  '/api/v1/warehouses?size=20',
  '/api/v1/transfers?size=20',
  '/api/v1/audits?size=20',
  '/api/v1/closings?size=20',
  ...(EDITION_UID
    ? [
        `/api/v1/stocks/${EDITION_UID}`,
        `/api/v1/stocks/${EDITION_UID}/history`,
      ]
    : []),
];
