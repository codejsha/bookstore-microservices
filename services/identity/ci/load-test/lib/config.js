export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
const TOKEN = __ENV.AUTH_TOKEN || '';
const USER_UID = __ENV.USER_UID || '';
const USER_EMAIL = __ENV.USER_EMAIL || '';

export function headers() {
  const h = { 'Content-Type': 'application/json' };
  if (TOKEN) h.Authorization = `Bearer ${TOKEN}`;
  return h;
}

export const READ_ENDPOINTS = [
  '/api/v1/users?size=20',
  ...(USER_UID ? [`/api/v1/users/${USER_UID}`] : []),
  ...(USER_EMAIL ? [`/api/v1/users/by-email/${encodeURIComponent(USER_EMAIL)}`] : []),
];
