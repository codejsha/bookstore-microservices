export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const NOTIFICATION_UID = __ENV.NOTIFICATION_UID || '';
const TEMPLATE_UID = __ENV.TEMPLATE_UID || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/notifications?size=20',
  '/api/v1/templates?size=20',
  ...(NOTIFICATION_UID ? [`/api/v1/notifications/${NOTIFICATION_UID}`] : []),
  ...(TEMPLATE_UID ? [`/api/v1/templates/${TEMPLATE_UID}`] : []),
];
