export function numericIdFromUid(uid: string): number {
  const digits = uid.replace(/\D/g, "");
  if (digits) return Number.parseInt(digits, 10);
  let hash = 0;
  for (let i = 0; i < uid.length; i++) {
    hash = (hash * 31 + uid.charCodeAt(i)) | 0;
  }
  return Math.abs(hash);
}

export function syntheticPrice(uid: string): number {
  const n = numericIdFromUid(uid);
  return Number((14.99 + (n % 40) * 1.5).toFixed(2));
}
